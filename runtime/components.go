package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// BoundedGovernor implements COAS authority rules without external framework code.
type BoundedGovernor struct{}

func (BoundedGovernor) Evaluate(task Task) (AuthorityDecision, error) {
	if err := validateTask(task); err != nil {
		return AuthorityDecision{}, err
	}

	decision := AuthorityDecision{HumanOverrideAvailable: true}
	switch {
	case task.Stakes >= 80:
		decision.Mode = ModeAssist
		decision.HumanConfirmationRequired = true
		decision.Reason = "high-stakes work keeps final judgment with the human"
	case task.CognitiveLoad < 30:
		decision.Mode = ModeManual
		decision.Reason = "low-load task does not justify automation overhead"
	case task.Reversibility < 40:
		decision.Mode = ModeAssist
		decision.HumanConfirmationRequired = true
		decision.Reason = "low-reversibility work requires supervised assistance"
	case task.CognitiveLoad >= 70 && task.Stakes < 60 && task.Reversibility >= 70:
		decision.Mode = ModeAutomateReversible
		decision.HumanConfirmationRequired = true
		decision.Reason = "high-load, bounded-stakes, reversible work is eligible for automation"
	default:
		decision.Mode = ModeAssist
		decision.Reason = "shared execution preserves human judgment"
	}
	return decision, nil
}

// SequencePlanner is deterministic by design. It walks the task's requested
// capabilities in order and stops when each requested capability has run once.
type SequencePlanner struct{}

func (SequencePlanner) Next(_ context.Context, task Task, state State) (string, bool, error) {
	index := len(state.PriorSteps)
	if index >= len(task.Capabilities) {
		return "", true, nil
	}
	capability := task.Capabilities[index]
	if capability == "" {
		return "", false, errors.New("task contains an empty capability")
	}
	return capability, false, nil
}

// InMemoryMemory is a concurrency-safe reference memory implementation.
type InMemoryMemory struct {
	mu       sync.RWMutex
	evidence map[string][]Evidence
	results  map[string][]Result
}

func NewInMemoryMemory() *InMemoryMemory {
	return &InMemoryMemory{
		evidence: make(map[string][]Evidence),
		results:  make(map[string][]Result),
	}
}

func (m *InMemoryMemory) AppendEvidence(_ context.Context, evidence Evidence) error {
	if evidence.TaskID == "" {
		return errors.New("evidence task id is required")
	}
	if evidence.Claim == "" {
		return errors.New("evidence claim is required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.evidence[evidence.TaskID] = append(m.evidence[evidence.TaskID], evidence)
	return nil
}

func (m *InMemoryMemory) AppendResult(_ context.Context, result Result) error {
	if result.TaskID == "" {
		return errors.New("result task id is required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.results[result.TaskID] = append(m.results[result.TaskID], result)
	return nil
}

func (m *InMemoryMemory) EvidenceForTask(_ context.Context, taskID string) ([]Evidence, error) {
	if taskID == "" {
		return nil, errors.New("task id is required")
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := m.evidence[taskID]
	copyOfItems := make([]Evidence, len(items))
	copy(copyOfItems, items)
	return copyOfItems, nil
}

// StrictVerifier requires at least one successful step and rejects unsupported
// evidence classifications. It intentionally does not infer truth from memory.
type StrictVerifier struct{}

func (StrictVerifier) Verify(_ context.Context, task Task, result Result) error {
	if result.TaskID != task.ID {
		return errors.New("result task id does not match task")
	}
	if len(result.Steps) == 0 {
		return errors.New("no execution steps were produced")
	}

	successes := 0
	for _, step := range result.Steps {
		if step.Succeeded {
			successes++
		}
	}
	if successes == 0 {
		return errors.New("all execution steps failed")
	}

	for _, item := range result.Evidence {
		switch item.Classification {
		case ClassFact, ClassInference, ClassHypothesis, ClassUnknown:
		default:
			return fmt.Errorf("unsupported evidence classification %q", item.Classification)
		}
	}
	return nil
}

// FuncWorker is a small adapter for registering project-specific functions as workers.
type FuncWorker struct {
	WorkerName         string
	WorkerCapabilities []string
	Run                func(context.Context, Task, State) (StepResult, []Evidence, error)
}

func (w FuncWorker) Name() string { return w.WorkerName }

func (w FuncWorker) Capabilities() []string {
	out := make([]string, len(w.WorkerCapabilities))
	copy(out, w.WorkerCapabilities)
	return out
}

func (w FuncWorker) Execute(ctx context.Context, task Task, state State) (StepResult, []Evidence, error) {
	if w.Run == nil {
		return StepResult{}, nil, errors.New("worker run function is required")
	}
	return w.Run(ctx, task, state)
}
