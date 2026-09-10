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

func TestDecideMergeAllowsOrdinaryRepositoryChanges(t *testing.T) {
	decision := DecideMerge(
		validMergeReport(),
		[]string{
			"swarm.go",
			"swarm_test.go",
			"README.md",
			"docs/generated/research.md",
			"cmd/new-agent-tool/main.go",
			"artifacts/swarm/latest.json",
		},
		true,
		DefaultMergeConstitution(),
	)
	if !decision.Allow {
		t.Fatalf("DecideMerge() denied ordinary repository changes: %v", decision.Reasons)
	}
}

func TestDecideMergeRejectsControlPlaneChanges(t *testing.T) {
	protected := []string{
		".github/workflows/swarm.yml",
		"autonomy.go",
		"autonomy_test.go",
		"cmd/coas-merge-policy/main.go",
		"go.mod",
	}
	for _, path := range protected {
		t.Run(path, func(t *testing.T) {
			decision := DecideMerge(
				validMergeReport(),
				[]string{path},
				true,
				DefaultMergeConstitution(),
			)
			if decision.Allow {
				t.Fatalf("DecideMerge() allowed immutable control-plane change: %s", path)
			}
		})
	}
}

func TestDecideMergeRejectsFailedChecks(t *testing.T) {
	decision := DecideMerge(
		validMergeReport(),
		[]string{"swarm.go"},
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
		[]string{"README.md"},
		true,
		DefaultMergeConstitution(),
	)
	if decision.Allow {
		t.Fatal("DecideMerge() allowed an incomplete swarm")
	}
}
