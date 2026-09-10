package coas

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepositoryContainsNoEmDash(t *testing.T) {
	allowed := map[string]bool{
		".go": true, ".md": true, ".yml": true, ".yaml": true,
		".json": true, ".txt": true, ".cff": true, ".html": true,
		".css": true, ".js": true, ".ts": true,
	}

	err := filepath.WalkDir(".", func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !allowed[strings.ToLower(filepath.Ext(path))] {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), "\u2014") {
			t.Errorf("%s contains forbidden U+2014 punctuation", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("repository punctuation scan failed: %v", err)
	}
}
