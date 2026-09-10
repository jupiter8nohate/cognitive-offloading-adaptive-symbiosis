package coas

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPreservationManifestIsDeterministic(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	first, err := BuildPreservationManifest(root, "abc123", []string{"ERR_404_GLITCHOLOGY.md"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildPreservationManifest(root, "abc123", []string{"ERR_404_GLITCHOLOGY.md"})
	if err != nil {
		t.Fatal(err)
	}
	if first.ManifestHash != second.ManifestHash {
		t.Fatal("manifest hash is not deterministic")
	}
}

func TestPreservationManifestDetectsCanonicalAsset(t *testing.T) {
	root := t.TempDir()
	book := filepath.Join(root, "books", "ERR_404_GLITCHOLOGY.md")
	if err := os.MkdirAll(filepath.Dir(book), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(book, []byte("canonical text"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest, err := BuildPreservationManifest(root, "abc123", []string{"ERR_404_GLITCHOLOGY.md"})
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.CanonicalAssets) != 1 || !manifest.CanonicalAssets[0].Present {
		t.Fatalf("canonical asset status = %#v", manifest.CanonicalAssets)
	}
}

func TestPreservationManifestExcludesRuntimeArtifacts(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "artifacts", "swarm"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "artifacts", "swarm", "latest.json"), []byte("runtime"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "coas.go"), []byte("package coas"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest, err := BuildPreservationManifest(root, "abc123", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Entries) != 1 || manifest.Entries[0].Path != "coas.go" {
		t.Fatalf("unexpected entries: %#v", manifest.Entries)
	}
}
