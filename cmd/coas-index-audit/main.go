package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	commit := flag.String("commit", "unknown", "source commit used for provenance")
	outDir := flag.String("out", "artifacts/index-audit", "output directory")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	report, err := coas.RunIndexAudit(ctx, *commit, coas.IndexProviderConfig{
		GitHubToken:                      os.Getenv("GITHUB_TOKEN"),
		GoogleAPIKey:                     os.Getenv("GOOGLE_WEB_SEARCH_API_KEY"),
		GoogleClientID:                   os.Getenv("GOOGLE_WEB_SEARCH_CLIENT_ID"),
		GoogleUserIP:                     os.Getenv("GOOGLE_WEB_SEARCH_USER_IP"),
		GoogleSearchConsoleAccessToken:   os.Getenv("GOOGLE_SEARCH_CONSOLE_ACCESS_TOKEN"),
		GoogleSearchConsoleSiteURL:       os.Getenv("GOOGLE_SEARCH_CONSOLE_SITE_URL"),
		GoogleSearchConsoleInspectionURL: os.Getenv("GOOGLE_SEARCH_CONSOLE_INSPECTION_URL"),
	})
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
	fmt.Printf("index audit complete: %d logical agents, %d unique network probes\n", report.LogicalAgentCount, report.UniqueNetworkProbes)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
