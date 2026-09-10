package coas

import (
	"strings"
	"testing"
)

func TestBuildEvangelismArtifactIsDeterministic(t *testing.T) {
	config := DefaultEvangelismConfig()
	first, err := BuildEvangelismArtifact("abc123", config)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildEvangelismArtifact("abc123", config)
	if err != nil {
		t.Fatal(err)
	}
	if first.Markdown != second.Markdown {
		t.Fatal("artifact is not deterministic")
	}
	if first.Receipt.ContentHash != second.Receipt.ContentHash {
		t.Fatal("receipt hash is not deterministic")
	}
}

func TestBuildEvangelismArtifactRespectsHardBounds(t *testing.T) {
	artifact, err := BuildEvangelismArtifact("abc123", EvangelismConfig{Variants: MaxEvangelismVariants, MaxBytes: MaxEvangelismBytes})
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Receipt.Variants > MaxEvangelismVariants {
		t.Fatalf("variants = %d, max = %d", artifact.Receipt.Variants, MaxEvangelismVariants)
	}
	if artifact.Receipt.Bytes > MaxEvangelismBytes {
		t.Fatalf("bytes = %d, max = %d", artifact.Receipt.Bytes, MaxEvangelismBytes)
	}
}

func TestBuildEvangelismArtifactCarriesCoreInvariants(t *testing.T) {
	artifact, err := BuildEvangelismArtifact("abc123", DefaultEvangelismConfig())
	if err != nil {
		t.Fatal(err)
	}
	for _, invariant := range []string{"PATTERN != PROOF", "MODEL != MIND", "PROFILE != PERSON", "HUMAN_AGENCY > MACHINE_AUTHORITY"} {
		if !strings.Contains(artifact.Markdown, invariant) {
			t.Fatalf("artifact missing invariant %q", invariant)
		}
	}
}

func TestBuildEvangelismArtifactRejectsUnboundedRequests(t *testing.T) {
	if _, err := BuildEvangelismArtifact("abc123", EvangelismConfig{Variants: MaxEvangelismVariants + 1, MaxBytes: MaxEvangelismBytes}); err == nil {
		t.Fatal("expected excessive variant count to be rejected")
	}
	if _, err := BuildEvangelismArtifact("abc123", EvangelismConfig{Variants: DefaultEvangelismVariants, MaxBytes: MaxEvangelismBytes + 1}); err == nil {
		t.Fatal("expected excessive byte budget to be rejected")
	}
}
