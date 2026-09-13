package runtime

import (
	"context"
	"errors"
	"testing"
)

func TestBoundedGovernor(t *testing.T) {
	governor := BoundedGovernor{}

	tests := []struct {
		name string
		task Task
		want Mode
	}{
		{
			name: "low load stays manual",
			task: Task{ID: "t1", Goal: "small task", CognitiveLoad: 10, Stakes: 10, Reversibility: 100, MaxSteps: 1},
			want: ModeManual,
		},
		{
			name: "high stakes stays assisted",
			task: Task{ID: "t2", Goal: "important task", CognitiveLoad: 95, Stakes: 90, Reversibility: 90, MaxSteps: 1},
			want: ModeAssist,
		},
		{
			name: "reversible bounded work can automate",
			task: Task{ID: "t3", Goal: "reversible task", CognitiveLoad: 90, Stakes: 20, Reversibility: 95, MaxSteps: 1},
			want: ModeAutomateReversible,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision, err := governor.Evaluate(tt.task)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}
			if decision.Mode != tt.want {
				t.Fatalf("mode = %s, want %s", decision.Mode, tt.want)
			}
			if !decision.HumanOverrideAvailable {
				t.Fatal("human override must always remain available")
			}
		})
	}
}

func TestEngineRequiresConfirmationBeforeExecution(t *testing.T) {
	memory := NewInMemoryMemory()
	engine, err := NewEngine(BoundedGovernor{}, SequencePlanner{}, memory, StrictVerifier{})
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	called := false
	if err := engine.Register(FuncWorker{
		WorkerName:         "GUARDED_WORKER",
		WorkerCapabilities: []string{"verify"},
		Run: func(context.Context, Task, State) (StepResult, []Evidence, error) {
			called = true
			return StepResult{Succeeded: true}, nil, nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	result, err := engine.Run(context.Background(), Task{
		ID:            "confirm-1",
		Goal:          "require human confirmation",
		CognitiveLoad: 80,
		Stakes:        20,
		Reversibility: 100,
		Capabilities:  []string{"verify"},
		MaxSteps:      1,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.AwaitingConfirmation {
		t.Fatal("result should be awaiting confirmation")
	}
	if called {
		t.Fatal("worker executed before confirmation")
	}
}

func TestEngineRunsRegisteredCapabilities(t *testing.T) {
	memory := NewInMemoryMemory()
	engine, err := NewEngine(BoundedGovernor{}, SequencePlanner{}, memory, StrictVerifier{})
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	worker := FuncWorker{
		WorkerName:         "VERIFY_WORKER",
		WorkerCapabilities: []string{"verify"},
		Run: func(_ context.Context, task Task, _ State) (StepResult, []Evidence, error) {
			return StepResult{
				Input:     task.Goal,
				Output:    "verified",
				Succeeded: true,
			}, []Evidence{{
				Classification: ClassFact,
				Claim:          "test worker completed",
				Source:         "runtime unit test",
			}}, nil
		},
	}
	if err := engine.Register(worker); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	result, err := engine.Run(context.Background(), Task{
		ID:             "verify-1",
		Goal:           "verify clean-room runtime",
		CognitiveLoad:  80,
		Stakes:         20,
		Reversibility:  100,
		Capabilities:   []string{"verify"},
		MaxSteps:       3,
		HumanConfirmed: true,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Completed {
		t.Fatal("result should be completed")
	}
	if len(result.Steps) != 1 {
		t.Fatalf("steps = %d, want 1", len(result.Steps))
	}
	if len(result.Evidence) != 1 {
		t.Fatalf("evidence = %d, want 1", len(result.Evidence))
	}

	stored, err := memory.EvidenceForTask(context.Background(), "verify-1")
	if err != nil {
		t.Fatalf("EvidenceForTask() error = %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("stored evidence = %d, want 1", len(stored))
	}
}

func TestEngineRecoversByContinuingAfterWorkerFailure(t *testing.T) {
	memory := NewInMemoryMemory()
	engine, err := NewEngine(BoundedGovernor{}, SequencePlanner{}, memory, StrictVerifier{})
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	failing := FuncWorker{
		WorkerName:         "A_FAILING_WORKER",
		WorkerCapabilities: []string{"first"},
		Run: func(context.Context, Task, State) (StepResult, []Evidence, error) {
			return StepResult{Input: "first attempt"}, nil, errors.New("synthetic failure")
		},
	}
	passing := FuncWorker{
		WorkerName:         "B_RECOVERY_WORKER",
		WorkerCapabilities: []string{"second"},
		Run: func(context.Context, Task, State) (StepResult, []Evidence, error) {
			return StepResult{Output: "recovered", Succeeded: true}, nil, nil
		},
	}

	if err := engine.Register(failing); err != nil {
		t.Fatal(err)
	}
	if err := engine.Register(passing); err != nil {
		t.Fatal(err)
	}

	result, err := engine.Run(context.Background(), Task{
		ID:             "recover-1",
		Goal:           "demonstrate bounded recovery",
		CognitiveLoad:  75,
		Stakes:         10,
		Reversibility:  100,
		Capabilities:   []string{"first", "second"},
		MaxSteps:       4,
		HumanConfirmed: true,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Completed {
		t.Fatal("expected recovery result to verify")
	}
	if len(result.Steps) != 2 || result.Steps[0].Succeeded || !result.Steps[1].Succeeded {
		t.Fatalf("unexpected recovery step sequence: %+v", result.Steps)
	}
}

func TestStrictVerifierRejectsAllFailedSteps(t *testing.T) {
	err := (StrictVerifier{}).Verify(context.Background(), Task{ID: "x"}, Result{
		TaskID: "x",
		Steps:  []StepResult{{Succeeded: false}},
	})
	if err == nil {
		t.Fatal("expected verifier to reject all-failed execution")
	}
}
