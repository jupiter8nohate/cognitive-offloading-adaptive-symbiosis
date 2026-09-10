package main

import (
	"fmt"
	"log"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	task := coas.Task{
		Name:          "organize a large research archive",
		CognitiveLoad: 85,
		Stakes:        25,
		Reversibility: 95,
	}

	rec, err := coas.Recommend(task)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TASK: %s\n", task.Name)
	fmt.Printf("MODE: %s\n", rec.Mode)
	fmt.Printf("WHY: %s\n", rec.Reason)
	fmt.Printf("CONFIRMATION_REQUIRED: %t\n", rec.HumanConfirmationRequired)
	fmt.Printf("HUMAN_OVERRIDE: %t\n", rec.HumanOverrideAvailable)
}
