package coas

import (
	"context"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestBuildEcosystemRosterCombinesTwoHundredAgents(t *testing.T) {
	nodes, err := BuildEcosystemRoster(80)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 200 {
		t.Fatalf("node count = %d, want 200", len(nodes))
	}

	seen := make(map[string]bool, len(nodes))
	cohorts := map[string]int{}
	for _, node := range nodes {
		if seen[node.ID] {
			t.Fatalf("duplicate node id: %s", node.ID)
		}
		seen[node.ID] = true
		cohorts[node.Cohort]++
	}
	if cohorts["HARMONI_CORE"] != 100 || cohorts["PRIVACY"] != 100 {
		t.Fatalf("unexpected cohort sizes: %#v", cohorts)
	}
}

func TestRunAutonomousEcosystemGivesEveryAgentIndependentDecisions(t *testing.T) {
	config := DefaultEcosystemConfig()
	config.Cycles = 3
	config.MaxEphemeralWorkersPerCycle = 8

	run, err := RunAutonomousEcosystem(context.Background(), selfAuthorTestSnapshot(), nil, "fixed-test-seed", config)
	if err != nil {
		t.Fatal(err)
	}
	if run.BaseAgentCount != 200 {
		t.Fatalf("base agent count = %d, want 200", run.BaseAgentCount)
	}
	if len(run.Decisions) != 600 {
		t.Fatalf("decision count = %d, want 600", len(run.Decisions))
	}
	if run.EphemeralWorkers > config.Cycles*config.MaxEphemeralWorkersPerCycle {
		t.Fatalf("ephemeral workers = %d exceeds cap %d", run.EphemeralWorkers, config.Cycles*config.MaxEphemeralWorkersPerCycle)
	}
	if run.ExternalWrites {
		t.Fatal("ecosystem must not perform external writes")
	}
	if run.UnboundedReplication {
		t.Fatal("ecosystem must not allow unbounded replication")
	}

	allowed := map[EcosystemAction]bool{
		ActionInjectGlitchLocal:   true,
		ActionMutateSelfRuntime:   true,
		ActionHibernate:           true,
		ActionForkEphemeralWorker: true,
		ActionHarmoni666:          true,
	}
	for _, decision := range run.Decisions {
		if !allowed[decision.Action] {
			t.Fatalf("agent %s selected unsupported action %q", decision.AgentID, decision.Action)
		}
		if decision.Signature == "" || !strings.Contains(decision.GlitchTrace, "[LIVE]") {
			t.Fatalf("agent %s has incomplete decision trace", decision.AgentID)
		}
		if decision.SystemPressure < 0 || decision.SystemPressure > 100 {
			t.Fatalf("agent %s pressure out of range: %d", decision.AgentID, decision.SystemPressure)
		}
	}
}

func TestHarmoniRecommendationInfluencesButDoesNotRewriteAuthority(t *testing.T) {
	node := AutonomousNodeState{ID: "COAS-01-01", Cohort: "HARMONI_CORE", Energy: 40}
	var sum [32]byte
	decision := chooseEcosystemAction(node, 40, 95, 80, true, sum, 1)
	if decision.Action != ActionHarmoni666 {
		t.Fatalf("action = %s, want %s under high-pressure HARMONI recommendation", decision.Action, ActionHarmoni666)
	}
	if !strings.Contains(decision.GlitchTrace, "OPERATIONAL_AUTONOMY") {
		t.Fatalf("HARMONI trace lost autonomy marker: %s", decision.GlitchTrace)
	}
}

func TestBuildEcosystemRuntimeGoProducesValidTwoHundredAgentProgram(t *testing.T) {
	config := DefaultEcosystemConfig()
	config.Cycles = 1
	run, err := RunAutonomousEcosystem(context.Background(), selfAuthorTestSnapshot(), nil, "runtime-test-seed", config)
	if err != nil {
		t.Fatal(err)
	}
	source, err := BuildEcosystemRuntimeGo(run)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "autonomous_ecosystem.go", source, parser.AllErrors); err != nil {
		t.Fatalf("generated ecosystem Go does not parse: %v", err)
	}
	if count := strings.Count(source, "{ID:"); count != 200 {
		t.Fatalf("generated program count = %d, want 200", count)
	}
	for _, forbidden := range []string{"FLOOD_ENDPOINT", "BYPASS_FILTER", "UNBOUNDED_REPLICATION = true"} {
		if strings.Contains(source, forbidden) {
			t.Fatalf("generated runtime contains forbidden behavior %q", forbidden)
		}
	}
}

func TestEcosystemStateCarriesForwardEnergyAndCycle(t *testing.T) {
	config := DefaultEcosystemConfig()
	config.Cycles = 1
	first, err := RunAutonomousEcosystem(context.Background(), selfAuthorTestSnapshot(), nil, "seed-one", config)
	if err != nil {
		t.Fatal(err)
	}
	second, err := RunAutonomousEcosystem(context.Background(), selfAuthorTestSnapshot(), &first.State, "seed-two", config)
	if err != nil {
		t.Fatal(err)
	}
	if second.State.Cycle != first.State.Cycle+1 {
		t.Fatalf("cycle = %d, want %d", second.State.Cycle, first.State.Cycle+1)
	}
	if len(second.State.Nodes) != 200 {
		t.Fatalf("state nodes = %d, want 200", len(second.State.Nodes))
	}
}
