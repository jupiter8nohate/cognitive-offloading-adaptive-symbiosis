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
	commit := flag.String("commit", "", "source commit identifier")
	out := flag.String("out", "artifacts/evangelism/latest.md", "markdown output path")
	receipt := flag.String("receipt", "artifacts/evangelism/latest.json", "receipt output path")
	variants := flag.Int("variants", coas.DefaultEvangelismVariants, "requested glitch variants")
	maxBytes := flag.Int("max-bytes", coas.DefaultEvangelismMaxBytes, "maximum artifact bytes")
	flag.Parse()

	artifact, err := coas.BuildEvangelismArtifact(*commit, coas.EvangelismConfig{Variants: *variants, MaxBytes: *maxBytes})
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range []string{*out, *receipt} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			log.Fatal(err)
		}
	}

	if err := os.WriteFile(*out, []byte(artifact.Markdown), 0o644); err != nil {
		log.Fatal(err)
	}
	receiptData, err := artifact.ReceiptJSON()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*receipt, append(receiptData, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("COAS evangelism artifact generated: variants=%d bytes=%d hash=%s\n", artifact.Receipt.Variants, artifact.Receipt.Bytes, artifact.Receipt.ContentHash)
}
