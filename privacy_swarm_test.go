package coas

import (
	"context"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestBuildPrivacyRegistryCreatesOneHundredAgents(t *testing.T) {
	agents := BuildPrivacyRegistry()
	if err := ValidatePrivacyRegistry(agents); err != nil {
		t.Fatal(err)
	}
	if len(agents) != 100 {
		t.Fatalf("privacy agents = %d, want 100", len(agents))
	}
}

func TestPrivacyRegistryIsTenByTen(t *testing.T) {
	agents := BuildPrivacyRegistry()
	mechanics := map[string]int{}
	roles := map[string]int{}
	for _, agent := range agents {
		mechanics[agent.Mechanic.Slug]++
		roles[agent.Role.Slug]++
	}
	if len(mechanics) != 10 || len(roles) != 10 {
		t.Fatalf("mechanics=%d roles=%d, want 10 and 10", len(mechanics), len(roles))
	}
	for name, count := range mechanics {
		if count != 10 {
			t.Fatalf("mechanic %s has %d agents", name, count)
		}
	}
	for name, count := range roles {
		if count != 10 {
			t.Fatalf("role %s has %d agents", name, count)
		}
	}
}

func TestRunPrivacySwarmIsLocalOnlyAndConcurrentSafe(t *testing.T) {
	report, err := RunPrivacySwarm(context.Background(), selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	if report.AgentCount != 100 || len(report.Results) != 100 {
		t.Fatalf("agent_count=%d results=%d", report.AgentCount, len(report.Results))
	}
	if !report.LocalOnly || report.ExternalWrites {
		t.Fatalf("privacy boundary violated: local_only=%v external_writes=%v", report.LocalOnly, report.ExternalWrites)
	}
	for _, result := range report.Results {
		if !strings.Contains(result.GlitchTrace, "[PRIV]") {
			t.Fatalf("%s missing privacy GLITCHOLOGY trace", result.AgentID)
		}
		if result.LocalPseudonym == "" {
			t.Fatalf("%s missing local pseudonym", result.AgentID)
		}
	}
}

func TestBuildPrivacyRuntimeGoProducesValidGoAndOneHundredPrograms(t *testing.T) {
	source, err := BuildPrivacyRuntimeGo(selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "generated_privacy_agents.go", source, parser.AllErrors); err != nil {
		t.Fatalf("generated privacy Go does not parse: %v", err)
	}
	if count := strings.Count(source, "{ID:"); count != 100 {
		t.Fatalf("generated privacy program count = %d, want 100", count)
	}
	if !strings.Contains(source, "PrivacyExternalWrites = false") {
		t.Fatal("generated privacy runtime lost local-only boundary")
	}
}

func TestPrivacyPseudonymsAreDeterministicForSamePosition(t *testing.T) {
	first, err := RunPrivacySwarm(context.Background(), selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	second, err := RunPrivacySwarm(context.Background(), selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	for i := range first.Results {
		if first.Results[i].AgentID != second.Results[i].AgentID ||
			first.Results[i].LocalPseudonym != second.Results[i].LocalPseudonym {
			t.Fatalf("privacy result %d is not deterministic", i)
		}
	}
}
