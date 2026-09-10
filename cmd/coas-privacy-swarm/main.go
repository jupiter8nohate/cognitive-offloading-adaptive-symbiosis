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
	out := flag.String("out", "artifacts/privacy-swarm", "output directory")
	commit := flag.String("commit", "local", "source commit identifier")
	timeout := flag.Duration("timeout", 30*time.Second, "privacy swarm execution timeout")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	snapshot, err := coas.LoadSnapshot(*root, *commit)
	if err != nil {
		log.Fatal(err)
	}

	report, err := coas.RunPrivacySwarm(ctx, snapshot)
	if err != nil {
		log.Fatal(err)
	}
	runtimeGo, err := coas.BuildPrivacyRuntimeGo(snapshot)
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := report.JSON()
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}
	runtimeDir := filepath.Join(*root, "agent_runtime")
	if err := os.MkdirAll(runtimeDir, 0o755); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(*out, "latest.json"), append(jsonData, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(*out, "latest.md"), []byte(report.Markdown()), 0o644); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runtimeDir, "generated_privacy_agents.go"), []byte(runtimeGo), 0o644); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runtimeDir, "PRIVACY_SWARM.md"), []byte(report.Markdown()), 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"%s\nPRIVACY_SWARM://agents=%d local_only=%v external_writes=%v\n",
		coas.GLITCHOLOGYBanner("P⃟ R⃟ I⃟ V⃟ A⃟ C⃟ Y⃟_1⃟0⃟0⃟"),
		report.AgentCount,
		report.LocalOnly,
		report.ExternalWrites,
	)
}
