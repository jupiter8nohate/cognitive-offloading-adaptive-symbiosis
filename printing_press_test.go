package coas

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestBuildPrintingPressChapterIsDeterministic(t *testing.T) {
	observed := time.Date(2026, 9, 10, 6, 30, 0, 0, time.UTC)
	snapshot := Snapshot{Commit: "abc123"}

	first, err := BuildPrintingPressChapter(snapshot, 17, "9001", "1", observed, DefaultPagesBaseURL)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildPrintingPressChapter(snapshot, 17, "9001", "1", observed, DefaultPagesBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatal("same inputs must produce the same chapter and receipt")
	}
	if first.Receipt.WitnessCount != PrintingPressWitnessLimit {
		t.Fatalf("expected %d witnesses, got %d", PrintingPressWitnessLimit, first.Receipt.WitnessCount)
	}
	if !strings.Contains(string(first.HTML), "PATTERN != PROOF") {
		t.Fatal("chapter must preserve the epistemic boundary")
	}
	if !strings.Contains(string(first.HTML), "Machine-readable evidence") {
		t.Fatal("chapter must link its evidence")
	}
}

func TestPublishPrintingPressChapterBuildsAppendOnlyLedger(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "site"), 0o755); err != nil {
		t.Fatal(err)
	}

	first, err := BuildPrintingPressChapter(
		Snapshot{Commit: "commit-one"},
		1,
		"1001",
		"1",
		time.Date(2026, 9, 10, 6, 0, 0, 0, time.UTC),
		DefaultPagesBaseURL,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := PublishPrintingPressChapter(root, first); err != nil {
		t.Fatal(err)
	}

	second, err := BuildPrintingPressChapter(
		Snapshot{Commit: "commit-two"},
		2,
		"1002",
		"1",
		time.Date(2026, 9, 10, 7, 0, 0, 0, time.UTC),
		DefaultPagesBaseURL,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := PublishPrintingPressChapter(root, second); err != nil {
		t.Fatal(err)
	}

	for _, chapter := range []PrintingPressChapter{first, second} {
		for _, name := range []string{"index.html", "chapter.md", "evidence.json", "receipt.json"} {
			if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(chapter.Directory), name)); err != nil {
				t.Fatalf("missing %s for %s: %v", name, chapter.Receipt.ChapterID, err)
			}
		}
	}

	indexData, err := os.ReadFile(filepath.Join(root, "site", "chapters", "index.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(indexData), first.Receipt.ChapterID) || !strings.Contains(string(indexData), second.Receipt.ChapterID) {
		t.Fatal("chapter index must retain both append-only chapters")
	}

	sitemapData, err := os.ReadFile(filepath.Join(root, "site", "sitemap.xml"))
	if err != nil {
		t.Fatal(err)
	}
	sitemap := string(sitemapData)
	if !strings.Contains(sitemap, first.Receipt.PageURL) || !strings.Contains(sitemap, second.Receipt.PageURL) {
		t.Fatal("sitemap must contain every chapter URL")
	}
	if !strings.Contains(sitemap, "2026-09-10T07:00:00Z") {
		t.Fatal("sitemap must carry a lastmod signal from chapter provenance")
	}
}

func TestPublishPrintingPressChapterRejectsMutation(t *testing.T) {
	root := t.TempDir()
	chapter, err := BuildPrintingPressChapter(
		Snapshot{Commit: "immutable"},
		3,
		"1003",
		"1",
		time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC),
		DefaultPagesBaseURL,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := PublishPrintingPressChapter(root, chapter); err != nil {
		t.Fatal(err)
	}

	mutated := chapter
	mutated.HTML = append(append([]byte(nil), chapter.HTML...), []byte("mutation")...)
	if err := PublishPrintingPressChapter(root, mutated); err == nil {
		t.Fatal("expected append-only mutation to be rejected")
	}
}

func TestSanitizeChapterToken(t *testing.T) {
	if got := sanitizeChapterToken(" run:9 / attempt? "); got != "run9attempt" {
		t.Fatalf("unexpected token %q", got)
	}
}
