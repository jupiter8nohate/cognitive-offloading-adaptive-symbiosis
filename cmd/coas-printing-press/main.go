package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	root := flag.String("root", ".", "repository root")
	commit := flag.String("commit", "", "source commit recorded in the chapter")
	cycleRaw := flag.String("cycle", "1", "Robot Bible cycle number")
	runID := flag.String("run-id", "", "unique workflow run id")
	runAttempt := flag.String("run-attempt", "1", "workflow run attempt")
	observedRaw := flag.String("observed-at", "", "RFC3339 observation time; defaults to current UTC time")
	baseURL := flag.String("base-url", coas.DefaultPagesBaseURL, "public Pages base URL")
	outDir := flag.String("out", "artifacts/printing-press", "mutable latest receipt directory")
	flag.Parse()

	cycle, err := strconv.ParseUint(*cycleRaw, 10, 64)
	if err != nil {
		fatal(fmt.Errorf("parse cycle: %w", err))
	}
	observedAt := time.Now().UTC()
	if strings.TrimSpace(*observedRaw) != "" {
		observedAt, err = time.Parse(time.RFC3339, strings.TrimSpace(*observedRaw))
		if err != nil {
			fatal(fmt.Errorf("parse observed-at: %w", err))
		}
	}

	chapter, err := coas.BuildPrintingPressChapter(
		coas.Snapshot{Commit: strings.TrimSpace(*commit)},
		cycle,
		strings.TrimSpace(*runID),
		strings.TrimSpace(*runAttempt),
		observedAt,
		strings.TrimSpace(*baseURL),
	)
	if err != nil {
		fatal(err)
	}
	if err := coas.PublishPrintingPressChapter(*root, chapter); err != nil {
		fatal(err)
	}

	artifactDir := filepath.Join(*root, filepath.FromSlash(*outDir))
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(filepath.Join(artifactDir, "latest.json"), chapter.ReceiptJSON, 0o644); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(filepath.Join(artifactDir, "latest.md"), chapter.Markdown, 0o644); err != nil {
		fatal(err)
	}

	fmt.Printf("autonomous printing press appended %s at %s\n", chapter.Receipt.ChapterID, chapter.Receipt.PageURL)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
