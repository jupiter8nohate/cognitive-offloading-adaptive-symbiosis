package coas

import "testing"

func validMergeReport() SwarmReport {
	results := make([]AgentResult, 0, 100)
	for _, agent := range BuildRegistry() {
		results = append(results, AgentResult{
			AgentID: agent.ID,
			Status:  "covered",
		})
	}
	return SwarmReport{
		Version:      "coas-swarm.v1",
		SourceCommit: "abc123",
		AgentCount:   100,
		Results:      results,
	}
}

func TestDecideMergeAllowsVerifiedSwarmArtifacts(t *testing.T) {
	decision := DecideMerge(
		validMergeReport(),
		[]string{"artifacts/swarm/latest.json", "artifacts/swarm/latest.md"},
		true,
		DefaultMergeConstitution(),
	)
	if !decision.Allow {
		t.Fatalf("DecideMerge() denied valid candidate: %v", decision.Reasons)
	}
}

func TestDecideMergeRejectsSourceCodeChanges(t *testing.T) {
	decision := DecideMerge(
		validMergeReport(),
		[]string{"swarm.go"},
		true,
		DefaultMergeConstitution(),
	)
	if decision.Allow {
		t.Fatal("DecideMerge() allowed source code mutation outside autonomous scope")
	}
}

func TestDecideMergeRejectsFailedChecks(t *testing.T) {
	decision := DecideMerge(
		validMergeReport(),
		[]string{"artifacts/swarm/latest.json"},
		false,
		DefaultMergeConstitution(),
	)
	if decision.Allow {
		t.Fatal("DecideMerge() allowed failed verification checks")
	}
}

func TestDecideMergeRejectsIncompleteSwarm(t *testing.T) {
	report := validMergeReport()
	report.AgentCount = 99
	report.Results = report.Results[:99]

	decision := DecideMerge(
		report,
		[]string{"artifacts/swarm/latest.json"},
		true,
		DefaultMergeConstitution(),
	)
	if decision.Allow {
		t.Fatal("DecideMerge() allowed an incomplete swarm")
	}
}
