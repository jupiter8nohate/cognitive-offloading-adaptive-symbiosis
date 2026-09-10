package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	root := flag.String("root", ".", "repository root")
	commit := flag.String("commit", "unknown", "source commit used for provenance")
	cycle := flag.Uint64("cycle", 0, "experiment cycle")
	outDir := flag.String("out", "artifacts/robot-bible", "output directory")
	flag.Parse()

	snapshot, err := coas.LoadSnapshot(*root, *commit)
	if err != nil {
		fatal(err)
	}

	report, err := coas.RunRobotBibleExperiment(snapshot, *cycle)
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
	if err := os.WriteFile(filepath.Join(*outDir, "latest.json"), append(jsonData, '\n'), 0o644); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(filepath.Join(*outDir, "latest.md"), []byte(report.Markdown()), 0o644); err != nil {
		fatal(err)
	}

	fmt.Printf("robot bible experiment complete: %d agents, cycle %d\n", report.AgentCount, report.Cycle)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
