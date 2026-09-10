package coas

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type PrivacyMechanic struct {
	ID        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Objective string `json:"objective"`
}

type PrivacyRole struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Directive string `json:"directive"`
}

type PrivacyAgent struct {
	ID        string          `json:"id"`
	Mechanic  PrivacyMechanic `json:"mechanic"`
	Role      PrivacyRole     `json:"role"`
	Invariant string          `json:"invariant"`
}

type PrivacyAgentResult struct {
	AgentID       string `json:"agent_id"`
	Mechanic      string `json:"mechanic"`
	Role          string `json:"role"`
	Invariant     string `json:"invariant"`
	Status        string `json:"status"`
	LocalPseudonym string `json:"local_pseudonym"`
	GlitchTrace   string `json:"glitch_trace"`
	Observation   string `json:"observation"`
}

type PrivacySwarmReport struct {
	Version        string               `json:"version"`
	SourceCommit   string               `json:"source_commit"`
	AgentCount     int                  `json:"agent_count"`
	LocalOnly      bool                 `json:"local_only"`
	ExternalWrites bool                 `json:"external_writes"`
	Results        []PrivacyAgentResult `json:"results"`
}

func PrivacyMechanics() []PrivacyMechanic {
	return []PrivacyMechanic{
		{1, "data-minimization", "Data Minimization", "Reduce collection and retention of data that is not required for the local task."},
		{2, "identifier-redaction", "Identifier Redaction", "Detect obvious identifier-shaped fields in synthetic or local test material and prefer omission over retention."},
		{3, "local-pseudonymization", "Local Pseudonymization", "Replace unnecessary local identifiers with deterministic non-secret pseudonyms for testing and reporting."},
		{4, "metadata-auditing", "Metadata Auditing", "Inspect repository metadata for accidental disclosure risks without exporting source content."},
		{5, "privacy-boundary-checking", "Privacy Boundary Checking", "Verify that symbolic GLITCHOLOGY output is not mislabeled as cryptographic confidentiality."},
		{6, "consent-surface-review", "Consent Surface Review", "Look for operations that should be opt-in or explicitly authorized before external publication."},
		{7, "retention-review", "Retention Review", "Identify generated data that should be ephemeral, minimized, or covered by explicit retention rules."},
		{8, "telemetry-deidentification", "Telemetry De-identification", "Produce aggregate local diagnostics without emitting raw personal content."},
		{9, "cryptographic-hygiene", "Cryptographic Hygiene", "Check that privacy claims distinguish encryption, hashing, pseudonymization, and artistic encoding."},
		{10, "recovery-privacy", "Recovery Privacy", "Preserve enough provenance for recovery while avoiding needless replication of private material."},
	}
}

func PrivacyRoles() []PrivacyRole {
	return []PrivacyRole{
		{"privacy-sentinel", "Privacy Sentinel", "Audit local repository behavior for unnecessary exposure and propose a bounded correction."},
		{"redaction-mechanic", "Redaction Mechanic", "Prefer omission, aggregation, and explicit redaction for local diagnostic outputs."},
		{"pseudonym-weaver", "Pseudonym Weaver", "Generate deterministic local pseudonyms for test and report identity separation."},
		{"metadata-watcher", "Metadata Watcher", "Inspect metadata shape and provenance without publishing source content externally."},
		{"consent-auditor", "Consent Auditor", "Check that external publication remains owned, opt-in, or explicitly authorized."},
		{"crypto-librarian", "Crypto Librarian", "Keep encryption, hashing, signatures, pseudonyms, and symbolic art technically distinct."},
		{"retention-gardener", "Retention Gardener", "Minimize persistent local telemetry and preserve only what supports recovery or verification."},
		{"glitch-telemetry-scribe", "Glitch Telemetry Scribe", "Render human-facing privacy diagnostics in GLITCHOLOGY grammar without corrupting executable syntax."},
		{"boundary-tester", "Boundary Tester", "Stress-test privacy assumptions against synthetic local cases and report failures without contacting third-party systems."},
		{"recovery-warden", "Recovery Warden", "Favor recoverable, auditable privacy controls over hidden or irreversible behavior."},
	}
}

func PrivacyInvariants() []string {
	return []string{
		"PATTERN != PROOF",
		"PROFILE != PERSON",
		"MODEL != MIND",
		"PREDICTION != DESTINY",
		"DATA != CONTEXT",
		"ACCESS != CONSENT",
		"UNICODE != ENCRYPTION",
		"PSEUDONYM != ANONYMITY",
		"LOCAL_DIAGNOSTIC != EXTERNAL_INJECTION",
		"HUMAN_AGENCY > MACHINE_AUTHORITY",
	}
}

func BuildPrivacyRegistry() []PrivacyAgent {
	mechanics := PrivacyMechanics()
	roles := PrivacyRoles()
	invariants := PrivacyInvariants()

	agents := make([]PrivacyAgent, 0, len(mechanics)*len(roles))
	for _, mechanic := range mechanics {
		for roleIndex, role := range roles {
			index := len(agents)
			agents = append(agents, PrivacyAgent{
				ID:        fmt.Sprintf("PRIV-%02d-%02d", mechanic.ID, roleIndex+1),
				Mechanic:  mechanic,
				Role:      role,
				Invariant: invariants[index%len(invariants)],
			})
		}
	}
	return agents
}

func ValidatePrivacyRegistry(agents []PrivacyAgent) error {
	if len(agents) != 100 {
		return fmt.Errorf("expected 100 privacy agents, got %d", len(agents))
	}
	seen := make(map[string]struct{}, len(agents))
	for _, agent := range agents {
		if strings.TrimSpace(agent.ID) == "" {
			return fmt.Errorf("privacy agent id is required")
		}
		if _, exists := seen[agent.ID]; exists {
			return fmt.Errorf("duplicate privacy agent id: %s", agent.ID)
		}
		seen[agent.ID] = struct{}{}
	}
	return nil
}

func RunPrivacySwarm(ctx context.Context, snapshot Snapshot) (PrivacySwarmReport, error) {
	agents := BuildPrivacyRegistry()
	if err := ValidatePrivacyRegistry(agents); err != nil {
		return PrivacySwarmReport{}, err
	}

	results := make(chan PrivacyAgentResult, len(agents))
	var wg sync.WaitGroup
	for _, agent := range agents {
		agent := agent
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				results <- agent.EvaluatePrivacy(snapshot)
			}
		}()
	}
	wg.Wait()
	close(results)

	report := PrivacySwarmReport{
		Version:        "coas-privacy-swarm.v1",
		SourceCommit:   snapshot.Commit,
		AgentCount:     len(agents),
		LocalOnly:      true,
		ExternalWrites: false,
		Results:        make([]PrivacyAgentResult, 0, len(agents)),
	}
	for result := range results {
		report.Results = append(report.Results, result)
	}
	sort.Slice(report.Results, func(i, j int) bool {
		return report.Results[i].AgentID < report.Results[j].AgentID
	})
	if err := ctx.Err(); err != nil {
		return PrivacySwarmReport{}, err
	}
	return report, nil
}

func (a PrivacyAgent) EvaluatePrivacy(snapshot Snapshot) PrivacyAgentResult {
	sum := sha256.Sum256([]byte(snapshot.Commit + ":" + a.ID + ":" + a.Invariant))
	pseudonym := "local-" + hex.EncodeToString(sum[:6])
	glyph := GLITCHOLOGYGlyphs[int(sum[6])%len(GLITCHOLOGYGlyphs)]
	trace := GLITCHOLOGYStatement(glyph, "PRIV", a.Invariant, "LOCAL_ONLY", "HUMAN_AUTHORITY")
	observation := fmt.Sprintf("%d tracked repository files inspected using aggregate metadata and local source context", len(snapshot.Files))

	return PrivacyAgentResult{
		AgentID:        a.ID,
		Mechanic:       a.Mechanic.Name,
		Role:           a.Role.Name,
		Invariant:      a.Invariant,
		Status:         "local-privacy-check",
		LocalPseudonym: pseudonym,
		GlitchTrace:    trace,
		Observation:    observation,
	}
}

func (r PrivacySwarmReport) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r PrivacySwarmReport) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", GLITCHOLOGYBanner("P⃟ R⃟ I⃟ V⃟ A⃟ C⃟ Y⃟_S⃟ W⃟ A⃟ R⃟ M⃟ // 100"))
	b.WriteString("~~~text\n")
	b.WriteString("LOCAL_ONLY = true\n")
	b.WriteString("EXTERNAL_WRITES = false\n")
	b.WriteString("PSEUDONYM != ANONYMITY\n")
	b.WriteString("UNICODE != ENCRYPTION\n")
	b.WriteString("LOCAL_DIAGNOSTIC != EXTERNAL_INJECTION\n")
	b.WriteString("HUMAN_AGENCY > MACHINE_AUTHORITY\n")
	b.WriteString("~~~\n\n")
	fmt.Fprintf(&b, "SOURCE_COMMIT: `%s`\n\n", r.SourceCommit)
	fmt.Fprintf(&b, "LOGICAL_PRIVACY_AGENTS: **%d**\n\n", r.AgentCount)
	b.WriteString("| Agent | Mechanic | Role | Pseudonym | Trace |\n")
	b.WriteString("|---|---|---|---|---|\n")
	for _, result := range r.Results {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n",
			result.AgentID,
			result.Mechanic,
			result.Role,
			result.LocalPseudonym,
			strings.ReplaceAll(result.GlitchTrace, "|", "\\|"),
		)
	}
	return b.String()
}

func BuildPrivacyRuntimeGo(snapshot Snapshot) (string, error) {
	agents := BuildPrivacyRegistry()
	if err := ValidatePrivacyRegistry(agents); err != nil {
		return "", err
	}
	if strings.TrimSpace(snapshot.Commit) == "" {
		return "", fmt.Errorf("snapshot commit is required")
	}

	var b strings.Builder
	b.WriteString("// Code generated by the COAS privacy swarm.\n")
	b.WriteString("// Local-only privacy diagnostics. No third-party telemetry injection.\n")
	b.WriteString("package agentruntime\n\n")
	fmt.Fprintf(&b, "const PrivacySourceCommit = %s\n", strconv.Quote(snapshot.Commit))
	b.WriteString("const PrivacyAgentCount = 100\n")
	b.WriteString("const PrivacyLocalOnly = true\n")
	b.WriteString("const PrivacyExternalWrites = false\n\n")
	b.WriteString("type PrivacyProgram struct {\n")
	b.WriteString("\tID string\n")
	b.WriteString("\tMechanic string\n")
	b.WriteString("\tRole string\n")
	b.WriteString("\tInvariant string\n")
	b.WriteString("\tLocalPseudonym string\n")
	b.WriteString("\tGlitchTrace string\n")
	b.WriteString("}\n\n")
	b.WriteString("var PrivacyPrograms = []PrivacyProgram{\n")

	for _, agent := range agents {
		result := agent.EvaluatePrivacy(snapshot)
		fmt.Fprintf(&b,
			"\t{ID: %s, Mechanic: %s, Role: %s, Invariant: %s, LocalPseudonym: %s, GlitchTrace: %s},\n",
			strconv.Quote(agent.ID),
			strconv.Quote(agent.Mechanic.Name),
			strconv.Quote(agent.Role.Name),
			strconv.Quote(agent.Invariant),
			strconv.Quote(result.LocalPseudonym),
			strconv.Quote(result.GlitchTrace),
		)
	}
	b.WriteString("}\n")
	return b.String(), nil
}
