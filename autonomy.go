package coas

import (
	"fmt"
	"sort"
	"strings"
)

type MergeConstitution struct {
	RequiredAgentCount int      `json:"required_agent_count"`
	ProtectedPrefixes  []string `json:"protected_prefixes"`
	ProtectedExact     []string `json:"protected_exact"`
	RequireChecks      bool     `json:"require_checks"`
	RequireProvenance  bool     `json:"require_provenance"`
}

type MergeDecision struct {
	Allow   bool     `json:"allow"`
	Reasons []string `json:"reasons"`
}

func DefaultMergeConstitution() MergeConstitution {
	return MergeConstitution{
		RequiredAgentCount: 100,
		ProtectedPrefixes: []string{
			".github/",
			".git/",
			"cmd/coas-merge-policy/",
		},
		ProtectedExact: []string{
			"autonomy.go",
			"autonomy_test.go",
			"go.mod",
			"go.sum",
		},
		RequireChecks:     true,
		RequireProvenance: true,
	}
}

func DecideMerge(report SwarmReport, changedPaths []string, checksPassed bool, constitution MergeConstitution) MergeDecision {
	reasons := make([]string, 0)

	if report.AgentCount != constitution.RequiredAgentCount {
		reasons = append(reasons, fmt.Sprintf("agent count %d does not equal required count %d", report.AgentCount, constitution.RequiredAgentCount))
	}
	if len(report.Results) != constitution.RequiredAgentCount {
		reasons = append(reasons, fmt.Sprintf("result count %d does not equal required count %d", len(report.Results), constitution.RequiredAgentCount))
	}
	if constitution.RequireChecks && !checksPassed {
		reasons = append(reasons, "verification checks did not pass")
	}
	if constitution.RequireProvenance && strings.TrimSpace(report.SourceCommit) == "" {
		reasons = append(reasons, "source commit provenance is missing")
	}
	if len(changedPaths) == 0 {
		reasons = append(reasons, "candidate contains no changed paths")
	}

	seenAgents := make(map[string]struct{}, len(report.Results))
	for _, result := range report.Results {
		if strings.TrimSpace(result.AgentID) == "" {
			reasons = append(reasons, "result contains an empty agent id")
			continue
		}
		if _, exists := seenAgents[result.AgentID]; exists {
			reasons = append(reasons, "duplicate agent result: "+result.AgentID)
		}
		seenAgents[result.AgentID] = struct{}{}
		if strings.TrimSpace(result.Status) == "" {
			reasons = append(reasons, "agent result has empty status: "+result.AgentID)
		}
	}

	for _, path := range changedPaths {
		clean := strings.TrimSpace(path)
		if clean == "" {
			continue
		}
		if isProtectedPath(clean, constitution) {
			reasons = append(reasons, "path is in immutable control plane: "+clean)
		}
	}

	sort.Strings(reasons)
	return MergeDecision{
		Allow:   len(reasons) == 0,
		Reasons: reasons,
	}
}

func isProtectedPath(path string, constitution MergeConstitution) bool {
	for _, exact := range constitution.ProtectedExact {
		if path == exact {
			return true
		}
	}
	for _, prefix := range constitution.ProtectedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
