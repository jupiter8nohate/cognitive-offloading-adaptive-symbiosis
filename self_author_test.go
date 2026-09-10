package coas

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func selfAuthorTestSnapshot() Snapshot {
	return Snapshot{
		Commit: "abc123",
		Files: map[string]string{
			"README.md":  "# COAS",
			"swarm.go":   "package coas",
			"policy.txt": "PATTERN != PROOF",
		},
	}
}

func TestBuildSelfAuthoredRuntimeProducesValidGo(t *testing.T) {
	runtime, err := BuildSelfAuthoredRuntime(selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "generated_agents.go", runtime.GoSource, parser.AllErrors); err != nil {
		t.Fatalf("generated Go does not parse: %v", err)
	}
	if count := strings.Count(runtime.GoSource, "{ID:"); count != 100 {
		t.Fatalf("generated program count = %d, want 100", count)
	}
}

func TestBuildSelfAuthoredRuntimeUsesCanonicalGLITCHOLOGY(t *testing.T) {
	runtime, err := BuildSelfAuthoredRuntime(selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		GLITCHOLOGYDisplayTitle,
		GLITCHOLOGYGrammar,
		"PATTERN != PROOF",
		"MODEL != MIND",
		"HUMAN_AGENCY > MACHINE_AUTHORITY",
		"MACHINE_CAN_READ != MACHINE_CAN_DEFINE",
	} {
		if !strings.Contains(runtime.GoSource+runtime.Markdown, required) {
			t.Fatalf("self-authored runtime missing %q", required)
		}
	}
	if runtime.Receipt.StyleSourceSHA != GLITCHOLOGYSourceSHA {
		t.Fatalf("style source sha = %q, want %q", runtime.Receipt.StyleSourceSHA, GLITCHOLOGYSourceSHA)
	}
}

func TestSelfAuthorFingerprintIgnoresItsOwnGeneratedNamespace(t *testing.T) {
	base := selfAuthorTestSnapshot()
	first := selfAuthorFingerprint(base)

	base.Files["agent_runtime/generated_agents.go"] = "package agentruntime\nconst X = 1\n"
	base.Files["artifacts/self-author/latest.json"] = "{}"
	second := selfAuthorFingerprint(base)

	if first != second {
		t.Fatal("self-authored outputs changed their own source fingerprint")
	}
}

func TestBuildSelfAuthoredRuntimeIsDeterministic(t *testing.T) {
	first, err := BuildSelfAuthoredRuntime(selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildSelfAuthoredRuntime(selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	if first.GoSource != second.GoSource || first.Markdown != second.Markdown {
		t.Fatal("self-authored runtime is not deterministic for the same repository position")
	}
	if first.Receipt.GeneratedGoSHA256 != second.Receipt.GeneratedGoSHA256 {
		t.Fatal("generated source receipt is not deterministic")
	}
}
