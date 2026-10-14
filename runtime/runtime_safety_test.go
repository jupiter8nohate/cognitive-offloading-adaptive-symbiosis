package runtime

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type brokenResultMemory struct {
	err error
}

func (m brokenResultMemory) AppendEvidence(context.Context, Evidence) error {
	return nil
}

func (m brokenResultMemory) AppendResult(context.Context, Result) error {
	return m.err
}

func (m brokenResultMemory) EvidenceForTask(context.Context, string) ([]Evidence, error) {
	return nil, nil
}

func TestRegisterRejectsInvalidCapabilityAtomically(t *testing.T) {
	engine, err := NewEngine(BoundedGovernor{}, SequencePlanner{}, NewInMemoryMemory(), StrictVerifier{})
	if err != nil {
		t.Fatal(err)
	}
	worker := FuncWorker{
		WorkerName:         "PARTIAL_WORKER",
		WorkerCapabilities: []string{"verify", " "},
		Run: func(context.Context, Task, State) (StepResult, []Evidence, error) {
			return StepResult{Succeeded: true}, nil, nil
		},
	}
	if err := engine.Register(worker); err == nil {
		t.Fatal("expected invalid capability to reject registration")
	}
	if _, err := engine.pickWorker("verify"); err == nil {
		t.Fatal("invalid registration partially modified worker registry")
	}
}

func TestRegisterRejectsDuplicateCapability(t *testing.T) {
	engine, err := NewEngine(BoundedGovernor{}, SequencePlanner{}, NewInMemoryMemory(), StrictVerifier{})
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Register(FuncWorker{
		WorkerName:         "DUPLICATE_WORKER",
		WorkerCapabilities: []string{"verify", "verify"},
	}); err == nil {
		t.Fatal("expected duplicate capability registration to fail")
	}
}

func TestEngineRejectsCanceledContext(t *testing.T) {
	engine, err := NewEngine(BoundedGovernor{}, SequencePlanner{}, NewInMemoryMemory(), StrictVerifier{})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = engine.Run(ctx, Task{ID: "cancelled", Goal: "do not execute", MaxSteps: 1})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

func TestEngineRejectsInvalidTaskBeforeExecution(t *testing.T) {
	engine, err := NewEngine(BoundedGovernor{}, SequencePlanner{}, NewInMemoryMemory(), StrictVerifier{})
	if err != nil {
		t.Fatal(err)
	}
	called := false
	err = engine.Register(FuncWorker{
		WorkerName:         "W",
		WorkerCapabilities: []string{"verify"},
		Run: func(context.Context, Task, State) (StepResult, []Evidence, error) {
			called = true
			return StepResult{Succeeded: true}, nil, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = engine.Run(context.Background(), Task{
		ID: "bad-caps", Goal: "reject blank capability",
		CognitiveLoad: 90, Stakes: 10, Reversibility: 90,
		Capabilities: []string{"verify", " "}, MaxSteps: 2,
		HumanConfirmed: true,
	})
	if err == nil || called {
		t.Fatalf("invalid task was allowed: err=%v called=%v", err, called)
	}
}

func TestVerifierRequiresAllRequestedCapabilitiesAttempted(t *testing.T) {
	task := Task{ID: "two-steps", Capabilities: []string{"first", "second"}}
	result := Result{TaskID: task.ID, Steps: []StepResult{{Succeeded: true}}}
	err := (StrictVerifier{}).Verify(context.Background(), task, result)
	if err == nil || !strings.Contains(err.Error(), "required capabilities") {
		t.Fatalf("verifier accepted incomplete task: %v", err)
	}
}

func TestEnginePropagatesManualResultStorageError(t *testing.T) {
	sentinel := errors.New("storage unavailable")
	engine, err := NewEngine(BoundedGovernor{}, SequencePlanner{}, brokenResultMemory{err: sentinel}, StrictVerifier{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = engine.Run(context.Background(), Task{
		ID: "manual", Goal: "a low-load task", CognitiveLoad: 5,
		Stakes: 5, Reversibility: 100, MaxSteps: 1,
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("error = %v, want storage failure", err)
	}
}

func TestEnginePropagatesConfirmationResultStorageError(t *testing.T) {
	sentinel := errors.New("storage unavailable")
	engine, err := NewEngine(BoundedGovernor{}, SequencePlanner{}, brokenResultMemory{err: sentinel}, StrictVerifier{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = engine.Run(context.Background(), Task{
		ID: "confirm", Goal: "a reversible task", CognitiveLoad: 90,
		Stakes: 10, Reversibility: 90, MaxSteps: 1,
		Capabilities: []string{"verify"},
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("error = %v, want storage failure", err)
	}
}
