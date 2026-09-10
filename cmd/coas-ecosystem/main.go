package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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
	commit := flag.String("commit", "local", "source commit identifier")
	statePath := flag.String("state", "artifacts/ecosystem/state.json", "persistent ecosystem state path")
	reportJSONPath := flag.String("report-json", "artifacts/ecosystem/latest.json", "ecosystem JSON report path")
	reportMarkdownPath := flag.String("report-md", "artifacts/ecosystem/latest.md", "ecosystem Markdown report path")
	runtimePath := flag.String("runtime", "agent_runtime/autonomous_ecosystem.go", "generated Go runtime path")
	seed := flag.String("seed", "", "optional reproducible decision seed; cryptographic entropy is used when empty")
	cycles := flag.Int("cycles", 3, "bounded decision cycles per run")
	maxWorkers := flag.Int("max-ephemeral-workers", 32, "maximum ephemeral child workers per cycle")
	timeout := flag.Duration("timeout", 45*time.Second, "maximum ecosystem runtime")
	flag.Parse()

	decisionSeed := *seed
	if decisionSeed == "" {
		var entropy [32]byte
		if _, err := rand.Read(entropy[:]); err != nil {
			log.Fatalf("decision entropy: %v", err)
		}
		decisionSeed = hex.EncodeToString(entropy[:])
	}

	snapshot, err := coas.LoadSnapshot(*root, *commit)
	if err != nil {
		log.Fatal(err)
	}

	var previous *coas.EcosystemState
	fullStatePath := filepath.Join(*root, filepath.FromSlash(*statePath))
	if data, err := os.ReadFile(fullStatePath); err == nil {
		var state coas.EcosystemState
		if err := json.Unmarshal(data, &state); err != nil {
			log.Fatalf("read previous ecosystem state: %v", err)
		}
		previous = &state
	} else if !os.IsNotExist(err) {
		log.Fatal(err)
	}

	config := coas.DefaultEcosystemConfig()
	config.Cycles = *cycles
	config.MaxEphemeralWorkersPerCycle = *maxWorkers

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	run, err := coas.RunAutonomousEcosystem(ctx, snapshot, previous, decisionSeed, config)
	if err != nil {
		log.Fatal(err)
	}
	runtimeGo, err := coas.BuildEcosystemRuntimeGo(run)
	if err != nil {
		log.Fatal(err)
	}
	reportJSON, err := run.JSON()
	if err != nil {
		log.Fatal(err)
	}
	stateJSON, err := run.State.JSON()
	if err != nil {
		log.Fatal(err)
	}

	outputs := map[string][]byte{
		filepath.Join(*root, filepath.FromSlash(*reportJSONPath)):     append(reportJSON, '\n'),
		filepath.Join(*root, filepath.FromSlash(*reportMarkdownPath)): []byte(run.Markdown()),
		fullStatePath:                                                append(stateJSON, '\n'),
		filepath.Join(*root, filepath.FromSlash(*runtimePath)):        []byte(runtimeGo),
	}
	for path, data := range outputs {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf(
		"%s\nLIVEBRAIN://base_agents=%d cycles=%d decisions=%d ephemeral_workers=%d harmoni_cycles=%d external_writes=%v\n",
		coas.GLITCHOLOGYBanner("L⃟ I⃟ V⃟ E⃟_B⃟ R⃟ A⃟ I⃟ N⃟_2⃟0⃟0⃟"),
		run.BaseAgentCount,
		run.Cycles,
		len(run.Decisions),
		run.EphemeralWorkers,
		run.HarmoniCycles,
		run.ExternalWrites,
	)
}
