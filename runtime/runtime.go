package runtime

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Mode defines the maximum authority granted to a runtime task.
type Mode string

const (
	ModeManual             Mode = "MANUAL"
	ModeAssist             Mode = "ASSIST"
	ModeAutomateReversible Mode = "AUTOMATE_REVERSIBLE"
)

// Classification separates observation from interpretation.
type Classification string

const (
	ClassFact       Classification = "FACT"
	ClassInference  Classification = "INFERENCE"
	ClassHypothesis Classification = "HYPOTHESIS"
	ClassUnknown    Classification = "UNKNOWN"
)

// Task is the unit of delegated cognitive work.
type Task struct {
	ID            string
	Goal          string
	CognitiveLoad int
	Stakes        int
	Reversibility int
	Capabilities  []string
	MaxSteps      int
}

// AuthorityDecision is an inspectable runtime permission decision.
type AuthorityDecision struct {
	Mode                      Mode
	HumanConfirmationRequired bool
	HumanOverrideAvailable    bool
	Reason                    string
}

// Evidence is a provenance-carrying observation produced by a worker.
type Evidence struct {
	TaskID         string
	Worker         string
	Classification Classification
	Claim          string
	Source         string
	RecordedAt     time.Time
}

// Result is the terminal state of a task execution.
type Result struct {
	TaskID     string
	Mode       Mode
	Completed  bool
	Summary    string
	Steps      []StepResult
	Evidence   []Evidence
	StartedAt  time.Time
	FinishedAt time.Time
}

// StepResult records one bounded worker action.
type StepResult struct {
	Worker    string
	Input     string
	Output    string
	Succeeded bool
	Err       string
}

// Worker performs one bounded, reversible action.
type Worker interface {
	Name() string
	Capabilities() []string
	Execute(ctx context.Context, task Task, state State) (StepResult, []Evidence, error)
}

// State is the read-only view given to workers.
type State struct {
	PriorSteps    []StepResult
	PriorEvidence []Evidence
}

// Memory stores runtime evidence and results. Implementations may be in-memory,
// file-backed, database-backed, or remote, but memory never becomes truth by itself.
type Memory interface {
	AppendEvidence(context.Context, Evidence) error
	AppendResult(context.Context, Result) error
	EvidenceForTask(context.Context, string) ([]Evidence, error)
}

// Verifier evaluates candidate results without granting itself execution authority.
type Verifier interface {
	Verify(context.Context, Task, Result) error
}

// Planner chooses the next worker capability based on the task and accumulated state.
type Planner interface {
	Next(context.Context, Task, State) (capability string, done bool, err error)
}

// Governor decides how much authority the machine receives.
type Governor interface {
	Evaluate(Task) (AuthorityDecision, error)
}

// Engine coordinates policy, planning, workers, memory, verification, and recovery.
type Engine struct {
	Governor Governor
	Planner  Planner
	Memory   Memory
	Verifier Verifier

	mu      sync.RWMutex
	workers map[string][]Worker
}

func NewEngine(governor Governor, planner Planner, memory Memory, verifier Verifier) (*Engine, error) {
	if governor == nil || planner == nil || memory == nil || verifier == nil {
		return nil, errors.New("governor, planner, memory, and verifier are required")
	}
	return &Engine{
		Governor: governor,
		Planner:  planner,
		Memory:   memory,
		Verifier: verifier,
		workers:  make(map[string][]Worker),
	}, nil
}

// Register adds a worker under each declared capability.
func (e *Engine) Register(worker Worker) error {
	if worker == nil || worker.Name() == "" {
		return errors.New("worker with a non-empty name is required")
	}
	caps := worker.Capabilities()
	if len(caps) == 0 {
		return fmt.Errorf("worker %q declares no capabilities", worker.Name())
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	for _, cap := range caps {
		if cap == "" {
			return fmt.Errorf("worker %q declares an empty capability", worker.Name())
		}
		e.workers[cap] = append(e.workers[cap], worker)
		sort.SliceStable(e.workers[cap], func(i, j int) bool {
			return e.workers[cap][i].Name() < e.workers[cap][j].Name()
		})
	}
	return nil
}

// Run executes a task only within the authority granted by the Governor.
func (e *Engine) Run(ctx context.Context, task Task) (Result, error) {
	if err := validateTask(task); err != nil {
		return Result{}, err
	}

	decision, err := e.Governor.Evaluate(task)
	if err != nil {
		return Result{}, fmt.Errorf("authority evaluation failed: %w", err)
	}

	result := Result{
		TaskID:    task.ID,
		Mode:      decision.Mode,
		StartedAt: time.Now().UTC(),
	}

	if decision.Mode == ModeManual {
		result.Summary = "manual mode selected; no machine execution performed"
		result.FinishedAt = time.Now().UTC()
		_ = e.Memory.AppendResult(ctx, result)
		return result, nil
	}

	state := State{}
	for step := 0; step < task.MaxSteps; step++ {
		capability, done, planErr := e.Planner.Next(ctx, task, state)
		if planErr != nil {
			return result, fmt.Errorf("planning failed at step %d: %w", step+1, planErr)
		}
		if done {
			break
		}

		worker, pickErr := e.pickWorker(capability)
		if pickErr != nil {
			return result, pickErr
		}

		stepResult, evidence, execErr := worker.Execute(ctx, task, state)
		stepResult.Worker = worker.Name()
		if execErr != nil {
			stepResult.Succeeded = false
			stepResult.Err = execErr.Error()
		}
		result.Steps = append(result.Steps, stepResult)
		state.PriorSteps = append(state.PriorSteps, stepResult)

		for _, item := range evidence {
			item.TaskID = task.ID
			item.Worker = worker.Name()
			if item.RecordedAt.IsZero() {
				item.RecordedAt = time.Now().UTC()
			}
			if err := e.Memory.AppendEvidence(ctx, item); err != nil {
				return result, fmt.Errorf("append evidence: %w", err)
			}
			result.Evidence = append(result.Evidence, item)
			state.PriorEvidence = append(state.PriorEvidence, item)
		}

		if execErr != nil {
			continue
		}
	}

	result.FinishedAt = time.Now().UTC()
	if err := e.Verifier.Verify(ctx, task, result); err != nil {
		result.Summary = "verification failed"
		_ = e.Memory.AppendResult(ctx, result)
		return result, fmt.Errorf("verification failed: %w", err)
	}

	result.Completed = true
	result.Summary = "verified within delegated authority"
	if err := e.Memory.AppendResult(ctx, result); err != nil {
		return result, fmt.Errorf("append result: %w", err)
	}
	return result, nil
}

func (e *Engine) pickWorker(capability string) (Worker, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	workers := e.workers[capability]
	if len(workers) == 0 {
		return nil, fmt.Errorf("no worker registered for capability %q", capability)
	}
	return workers[0], nil
}

func validateTask(task Task) error {
	if task.ID == "" {
		return errors.New("task id is required")
	}
	if task.Goal == "" {
		return errors.New("task goal is required")
	}
	for name, value := range map[string]int{
		"cognitive load": task.CognitiveLoad,
		"stakes":         task.Stakes,
		"reversibility":  task.Reversibility,
	} {
		if value < 0 || value > 100 {
			return fmt.Errorf("%s must be between 0 and 100", name)
		}
	}
	if task.MaxSteps <= 0 || task.MaxSteps > 100 {
		return errors.New("max steps must be between 1 and 100")
	}
	return nil
}
