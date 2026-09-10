package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	root := flag.String("root", ".", "repository root")
	commit := flag.String("commit", "", "source commit identifier")
	out := flag.String("out", "artifacts/preservation/latest.json", "manifest JSON output path")
	markdown := flag.String("markdown", "artifacts/preservation/latest.md", "manifest Markdown output path")
	flag.Parse()

	manifest, err := coas.BuildPreservationManifest(*root, *commit, []string{"ERR_404_GLITCHOLOGY.md"})
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range []string{*out, *markdown} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			log.Fatal(err)
		}
	}
	data, err := manifest.JSON()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, append(data, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*markdown, []byte(manifest.Markdown()), 0o644); err != nil {
		log.Fatal(err)
	}
}
