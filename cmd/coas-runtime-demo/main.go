package main

import (
	"context"
	"fmt"
	"log"

	coasruntime "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis/runtime"
)

func main() {
	memory := coasruntime.NewInMemoryMemory()
	engine, err := coasruntime.NewEngine(
		coasruntime.BoundedGovernor{},
		coasruntime.SequencePlanner{},
		memory,
		coasruntime.StrictVerifier{},
	)
	if err != nil {
		log.Fatal(err)
	}

	observe := coasruntime.FuncWorker{
		WorkerName:         "COAS_OBSERVER",
		WorkerCapabilities: []string{"observe"},
		Run: func(_ context.Context, task coasruntime.Task, _ coasruntime.State) (coasruntime.StepResult, []coasruntime.Evidence, error) {
			return coasruntime.StepResult{
				Input:     task.Goal,
				Output:    "observation recorded",
				Succeeded: true,
			}, []coasruntime.Evidence{{
				Classification: coasruntime.ClassFact,
				Claim:          "the clean-room observer executed successfully",
				Source:         "coas-runtime-demo",
			}}, nil
		},
	}

	review := coasruntime.FuncWorker{
		WorkerName:         "COAS_REVIEWER",
		WorkerCapabilities: []string{"review"},
		Run: func(_ context.Context, _ coasruntime.Task, state coasruntime.State) (coasruntime.StepResult, []coasruntime.Evidence, error) {
			return coasruntime.StepResult{
				Input:     fmt.Sprintf("prior evidence=%d", len(state.PriorEvidence)),
				Output:    "bounded review completed",
				Succeeded: true,
			}, nil, nil
		},
	}

	if err := engine.Register(observe); err != nil {
		log.Fatal(err)
	}
	if err := engine.Register(review); err != nil {
		log.Fatal(err)
	}

	task := coasruntime.Task{
		ID:             "demo-001",
		Goal:           "demonstrate original COAS orchestration",
		CognitiveLoad:  85,
		Stakes:         20,
		Reversibility:  100,
		Capabilities:   []string{"observe", "review"},
		MaxSteps:       4,
		HumanConfirmed: true,
	}

	result, err := engine.Run(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}

	receipt, err := coasruntime.SealResult(result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("mode=%s completed=%t steps=%d evidence=%d\n", result.Mode, result.Completed, len(result.Steps), len(result.Evidence))
	fmt.Printf("receipt=%s\n", receipt.Digest)
	fmt.Println("HUMAN_AGENCY > MACHINE_AUTHORITY")
}
