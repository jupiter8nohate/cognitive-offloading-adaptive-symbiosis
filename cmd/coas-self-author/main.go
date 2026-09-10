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

	outputs := map[string][]byte{
		filepath.Join(*root, "agent_runtime", "generated_agents.go"):           []byte(runtime.GoSource),
		filepath.Join(*root, "agent_runtime", "GENERATED_GLITCHOLOGY.md"):     []byte(runtime.Markdown),
		filepath.Join(*root, "artifacts", "self-author", "latest.json"):       nil,
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
		"%s\nSELF_AUTHOR://agents=%d source=%s style=%s generated=%s\n",
		coas.GLITCHOLOGYBanner("S⃟ E⃟ L⃟ F⃟_A⃟ U⃟ T⃟ H⃟ O⃟ R⃟"),
		runtime.Receipt.AgentCount,
		runtime.Receipt.SourceFingerprint,
		runtime.Receipt.StyleSourceSHA,
		runtime.Receipt.GeneratedGoSHA256,
	)
}
