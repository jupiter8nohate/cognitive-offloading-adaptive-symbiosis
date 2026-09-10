package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	root := flag.String("root", ".", "repository root")
	out := flag.String("out", "artifacts/swarm", "output directory")
	commit := flag.String("commit", "local", "source commit identifier")
	timeout := flag.Duration("timeout", 30*time.Second, "swarm execution timeout")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	snapshot, err := coas.LoadSnapshot(*root, *commit)
	if err != nil {
		log.Fatal(err)
	}

	report, err := coas.RunSwarm(ctx, snapshot)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}

	jsonData, err := report.JSON()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(*out, "latest.json"), append(jsonData, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(*out, "latest.md"), []byte(report.Markdown()), 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("COAS swarm completed: %d agents\n", report.AgentCount)
}
