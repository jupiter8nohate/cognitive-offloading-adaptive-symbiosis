package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	root := flag.String("root", ".", "repository root")
	commit := flag.String("commit", "", "source commit identifier")
	flag.Parse()

	snapshot, err := coas.LoadSnapshot(*root, *commit)
	if err != nil {
		log.Fatal(err)
	}
	runtime, err := coas.BuildSelfAuthoredRuntime(snapshot)
	if err != nil {
		log.Fatal(err)
	}
	plan, err := coas.BuildAutonomousPlan(snapshot)
	if err != nil {
		log.Fatal(err)
	}
	planJSON, err := plan.JSON()
	if err != nil {
		log.Fatal(err)
	}
	harmoni, err := coas.BuildHarmoniState(snapshot)
	if err != nil {
		log.Fatal(err)
	}
	harmoniJSON, err := harmoni.JSON()
	if err != nil {
		log.Fatal(err)
	}

	outputs := map[string][]byte{
		filepath.Join(*root, "agent_runtime", "generated_agents.go"):      []byte(runtime.GoSource),
		filepath.Join(*root, "agent_runtime", "GENERATED_GLITCHOLOGY.md"): []byte(runtime.Markdown),
		filepath.Join(*root, "artifacts", "self-author", "latest.json"):   nil,
		filepath.Join(*root, "artifacts", "strategy", "latest.json"):      append(planJSON, '\n'),
		filepath.Join(*root, "agent_runtime", "NEXT_MOVE.md"):             []byte(plan.Markdown()),
		filepath.Join(*root, "agent_runtime", "HARMONI_666.md"):            []byte(harmoni.Markdown()),
		filepath.Join(*root, "artifacts", "harmoni", "latest.json"):        append(harmoniJSON, '\n'),
	}
	receipt, err := runtime.Receipt.JSON()
	if err != nil {
		log.Fatal(err)
	}
	outputs[filepath.Join(*root, "artifacts", "self-author", "latest.json")] = append(receipt, '\n')

	for path, data := range outputs {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf(
		"%s\nSELF_AUTHOR://agents=%d harmoni_plays=%d source=%s style=%s generated=%s\n",
		coas.GLITCHOLOGYBanner("S⃟ E⃟ L⃟ F⃟_A⃟ U⃟ T⃟ H⃟ O⃟ R⃟"),
		runtime.Receipt.AgentCount,
		runtime.Receipt.HarmoniPlays,
		runtime.Receipt.SourceFingerprint,
		runtime.Receipt.StyleSourceSHA,
		runtime.Receipt.GeneratedGoSHA256,
	)
}
