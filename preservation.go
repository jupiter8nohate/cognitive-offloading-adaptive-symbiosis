package coas

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const MaxPreservedFileBytes int64 = 5 << 20

type PreservationEntry struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Bytes  int64  `json:"bytes"`
}

type CanonicalAssetStatus struct {
	Name    string `json:"name"`
	Present bool   `json:"present"`
	Path    string `json:"path,omitempty"`
	SHA256  string `json:"sha256,omitempty"`
}

type PreservationManifest struct {
	Version         string                 `json:"version"`
	SourceCommit    string                 `json:"source_commit"`
	Entries         []PreservationEntry    `json:"entries"`
	CanonicalAssets []CanonicalAssetStatus `json:"canonical_assets"`
	ManifestHash    string                 `json:"manifest_hash"`
}

func BuildPreservationManifest(root, commit string, canonicalNames []string) (PreservationManifest, error) {
	commit = strings.TrimSpace(commit)
	if commit == "" {
		return PreservationManifest{}, errors.New("commit is required")
	}

	entries := make([]PreservationEntry, 0)
	assets := make(map[string]CanonicalAssetStatus, len(canonicalNames))
	for _, name := range canonicalNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		assets[name] = CanonicalAssetStatus{Name: name}
	}

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == "artifacts" {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Size() > MaxPreservedFileBytes {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		hash := sha256.New()
		_, copyErr := io.Copy(hash, file)
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		digest := hex.EncodeToString(hash.Sum(nil))
		entries = append(entries, PreservationEntry{Path: rel, SHA256: digest, Bytes: info.Size()})

		base := filepath.Base(rel)
		if asset, ok := assets[base]; ok {
			asset.Present = true
			asset.Path = rel
			asset.SHA256 = digest
			assets[base] = asset
		}
		return nil
	})
	if err != nil {
		return PreservationManifest{}, err
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	canonical := make([]CanonicalAssetStatus, 0, len(assets))
	for _, asset := range assets {
		canonical = append(canonical, asset)
	}
	sort.Slice(canonical, func(i, j int) bool { return canonical[i].Name < canonical[j].Name })

	manifest := PreservationManifest{
		Version:         "coas-preservation.v1",
		SourceCommit:    commit,
		Entries:         entries,
		CanonicalAssets: canonical,
	}
	data, err := json.Marshal(manifest)
	if err != nil {
		return PreservationManifest{}, err
	}
	sum := sha256.Sum256(data)
	manifest.ManifestHash = hex.EncodeToString(sum[:])
	return manifest, nil
}

func (m PreservationManifest) JSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

func (m PreservationManifest) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# COAS Preservation Manifest\n\n")
	fmt.Fprintf(&b, "- Version: %s\n", m.Version)
	fmt.Fprintf(&b, "- Source commit: `%s`\n", m.SourceCommit)
	fmt.Fprintf(&b, "- Manifest SHA-256: `%s`\n", m.ManifestHash)
	fmt.Fprintf(&b, "- Tracked files: %d\n\n", len(m.Entries))
	b.WriteString("## Canonical assets\n\n")
	for _, asset := range m.CanonicalAssets {
		if asset.Present {
			fmt.Fprintf(&b, "- `%s`: present at `%s`, SHA-256 `%s`\n", asset.Name, asset.Path, asset.SHA256)
		} else {
			fmt.Fprintf(&b, "- `%s`: not present in this repository state\n", asset.Name)
		}
	}
	b.WriteString("\n## File ledger\n\n")
	b.WriteString("| Path | Bytes | SHA-256 |\n|---|---:|---|\n")
	for _, entry := range m.Entries {
		fmt.Fprintf(&b, "| `%s` | %d | `%s` |\n", entry.Path, entry.Bytes, entry.SHA256)
	}
	return b.String()
}
