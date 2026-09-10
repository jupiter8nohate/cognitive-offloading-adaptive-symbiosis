package coas

import (
	"context"
	"testing"
)

func TestBuildRegistryCreatesExactlyOneHundredBoundedAgents(t *testing.T) {
	agents := BuildRegistry()
	if err := ValidateRegistry(agents); err != nil {
		t.Fatalf("ValidateRegistry() error = %v", err)
	}
	if got := len(agents); got != 100 {
		t.Fatalf("len(BuildRegistry()) = %d, want 100", got)
	}
	for _, agent := range agents {
		if agent.CanMerge {
			t.Fatalf("%s unexpectedly has merge authority", agent.ID)
		}
		if agent.RiskCeiling != "reversible-only" {
			t.Fatalf("%s risk ceiling = %q", agent.ID, agent.RiskCeiling)
		}
	}
}

func TestRegistryIsTenMechanicsByTenRoles(t *testing.T) {
	agents := BuildRegistry()
	mechanics := map[string]int{}
	roles := map[string]int{}
	for _, agent := range agents {
		mechanics[agent.Mechanic.Slug]++
		roles[agent.Role.Slug]++
	}
	if len(mechanics) != 10 {
		t.Fatalf("mechanic count = %d, want 10", len(mechanics))
	}
	if len(roles) != 10 {
		t.Fatalf("role count = %d, want 10", len(roles))
	}
	for mechanic, count := range mechanics {
		if count != 10 {
			t.Fatalf("mechanic %s has %d agents, want 10", mechanic, count)
		}
	}
	for role, count := range roles {
		if count != 10 {
			t.Fatalf("role %s has %d agents, want 10", role, count)
		}
	}
}

func TestRunSwarmReturnsHundredResults(t *testing.T) {
	snapshot := Snapshot{
		Commit: "test",
		Files: map[string]string{
			"README.md": "# Cognitive Offloading Adaptive Symbiosis",
		},
	}
	report, err := RunSwarm(context.Background(), snapshot)
	if err != nil {
		t.Fatalf("RunSwarm() error = %v", err)
	}
	if report.AgentCount != 100 || len(report.Results) != 100 {
		t.Fatalf("report has agent_count=%d results=%d", report.AgentCount, len(report.Results))
	}
	for i := 1; i < len(report.Results); i++ {
		if report.Results[i-1].AgentID >= report.Results[i].AgentID {
			t.Fatal("results are not deterministically sorted")
		}
	}
}
