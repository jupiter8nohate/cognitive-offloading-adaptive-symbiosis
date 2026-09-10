package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	reportPath := flag.String("report", "artifacts/swarm/latest.json", "swarm report JSON path")
	changedPath := flag.String("changed", "", "file containing changed paths, one per line")
	checksPassed := flag.Bool("checks-passed", false, "whether verification checks passed")
	flag.Parse()

	if *changedPath == "" {
		log.Fatal("-changed is required")
	}

	reportData, err := os.ReadFile(*reportPath)
	if err != nil {
		log.Fatal(err)
	}
	var report coas.SwarmReport
	if err := json.Unmarshal(reportData, &report); err != nil {
		log.Fatal(err)
	}

	changedData, err := os.ReadFile(*changedPath)
	if err != nil {
		log.Fatal(err)
	}
	changed := make([]string, 0)
	for _, line := range strings.Split(string(changedData), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			changed = append(changed, line)
		}
	}

	decision := coas.DecideMerge(report, changed, *checksPassed, coas.DefaultMergeConstitution())
	if !decision.Allow {
		for _, reason := range decision.Reasons {
			fmt.Println("DENY:", reason)
		}
		os.Exit(2)
	}

	fmt.Println("ALLOW: autonomous merge constitution satisfied")
}
