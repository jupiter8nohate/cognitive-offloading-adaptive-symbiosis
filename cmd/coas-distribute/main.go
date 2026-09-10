package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	commit := flag.String("commit", "", "source commit identifier")
	dryRun := flag.Bool("dry-run", false, "validate and render without network publication")
	flag.Parse()

	payload, err := coas.BuildDistributionPayload(*commit)
	if err != nil {
		log.Fatal(err)
	}

	targets, err := coas.ParseDistributionTargets(os.Getenv("COAS_DISTRIBUTION_TARGETS_JSON"))
	if err != nil {
		log.Fatal(err)
	}
	if len(targets) == 0 {
		fmt.Println("COAS distribution: no authorized targets configured")
		return
	}

	if *dryRun {
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("COAS distribution dry run: %d authorized targets\n%s\n", len(targets), data)
		return
	}

	if !strings.EqualFold(os.Getenv("COAS_DISTRIBUTION_ENABLED"), "true") {
		log.Fatal("COAS_DISTRIBUTION_ENABLED must be true for network publication")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	publisher := coas.Publisher{}
	for _, target := range targets {
		result, err := publisher.Publish(ctx, target, payload)
		if err != nil {
			log.Printf("distribution failed target=%q kind=%q: %v", target.Name, target.Kind, err)
			continue
		}
		fmt.Printf("distribution success target=%q kind=%q status=%d\n", result.TargetName, result.Kind, result.StatusCode)
	}
}
