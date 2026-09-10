package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	commit := flag.String("commit", "unknown", "source commit recorded in the report")
	cycleRaw := flag.String("cycle", "1", "experiment cycle number")
	outDir := flag.String("out", "artifacts/god-search", "output directory")
	flag.Parse()

	cycle, err := strconv.ParseUint(*cycleRaw, 10, 64)
	if err != nil {
		fatal(fmt.Errorf("parse cycle: %w", err))
	}

	report, err := coas.RunGodSearchExperiment(coas.Snapshot{Commit: *commit}, cycle)
	if err != nil {
		fatal(err)
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fatal(err)
	}
	jsonData, err := report.JSON()
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile(filepath.Join(*outDir, "latest.json"), jsonData, 0o644); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(filepath.Join(*outDir, "latest.md"), []byte(report.Markdown()), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("god search experiment complete: %d agent views, cycle %d\n", len(report.AgentViews), report.Cycle)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
