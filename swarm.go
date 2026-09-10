package coas

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Mechanic struct {
	ID        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Objective string `json:"objective"`
	Safeguard string `json:"safeguard"`
}

type Role struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Directive string `json:"directive"`
	ReadOnly  bool   `json:"read_only"`
}

type Agent struct {
	ID          string   `json:"id"`
	Mechanic    Mechanic `json:"mechanic"`
	Role        Role     `json:"role"`
	CanMerge    bool     `json:"can_merge"`
	RiskCeiling string   `json:"risk_ceiling"`
}

type Snapshot struct {
	Commit string            `json:"commit"`
	Files  map[string]string `json:"files"`
}

type AgentResult struct {
	AgentID       string `json:"agent_id"`
	Mechanic      string `json:"mechanic"`
	Role          string `json:"role"`
	Status        string `json:"status"`
	Observation   string `json:"observation"`
	ProposedWork  string `json:"proposed_work"`
	AuthorityNote string `json:"authority_note"`
}

type SwarmReport struct {
	Version      string        `json:"version"`
	SourceCommit string        `json:"source_commit"`
	AgentCount   int           `json:"agent_count"`
	Results      []AgentResult `json:"results"`
}

func Mechanics() []Mechanic {
	return []Mechanic{
		{1, "syntax-valid-mind-dumping", "Syntax-Valid Mind Dumping", "Transform raw ideas into inspectable structure without requiring polished prose.", "Structure assists expression; it must not redefine the author's meaning."},
		{2, "multi-channel-code-switching", "Multi-Channel Code-Switching", "Support movement between natural language, code, and symbolic notation.", "Treat channel changes as communication choices, not diagnostic signals."},
		{3, "algorithmic-filtering", "Algorithmic Filtering", "Externalize sorting, indexing, clustering, and retrieval burden.", "Filtering must preserve provenance and allow the human to inspect what was omitted."},
		{4, "visual-cipher-anchoring", "Visual Cipher Anchoring", "Bind user-defined symbols to explicit user-defined meanings.", "Never infer emotional or clinical meaning from a symbol without explicit mapping."},
		{5, "semantic-deescalation", "Semantic De-escalation", "Translate dense communication into concise factual structure.", "Simplification must preserve uncertainty, context, and the option to inspect the original."},
		{6, "metacognitive-error-tracking", "Metacognitive Error-Tracking", "Log stalls, revisions, and corrections as process signals.", "A pause or typo is not proof of a mental state, diagnosis, or hidden intent."},
		{7, "predictive-interface-recalibration", "Predictive Interface Recalibration", "Adapt interfaces to explicit preferences and local friction signals.", "Do not label friction as stress or infer private mental state from interaction patterns."},
		{8, "co-compiled-evidence-checks", "Co-Compiled Evidence Checks", "Test claims against explicit evidence and logical constraints.", "The system offers evidence calibration, not authority over reality or personal perception."},
		{9, "asynchronous-processing-buffers", "Asynchronous Processing Buffers", "Hold raw drafts until the human chooses to publish or discard them.", "Buffers default to private, reversible, and user-controlled retention."},
		{10, "autonomous-sovereignty-shields", "Autonomous Sovereignty Shields", "Reduce unnecessary tracking through local processing, minimization, and explicit consent.", "Privacy protection must not rely on deception, harmful interference, or false claims of invisibility."},
	}
}

func Roles() []Role {
	return []Role{
		{"observer", "Observer", "Inventory relevant repository signals and detect missing structure.", true},
		{"mapper", "Mapper", "Map concepts, dependencies, inputs, outputs, and boundaries.", true},
		{"normalizer", "Normalizer", "Convert inconsistent structure into stable machine-readable form.", false},
		{"critic", "Critic", "Search for category errors, overclaims, unsafe inference, and ambiguity.", true},
		{"verifier", "Verifier", "Check claims and implementation against explicit invariants.", true},
		{"planner", "Planner", "Produce bounded next actions with acceptance criteria.", true},
		{"tester", "Tester", "Design or execute verification for the mechanic's behavior.", false},
		{"archivist", "Archivist", "Preserve provenance, definitions, and decision history.", false},
		{"provenance", "Provenance Keeper", "Track source, authorship, transformations, and uncertainty.", true},
		{"steward", "Steward", "Enforce authority boundaries and reject unsafe autonomous escalation.", true},
	}
}

func BuildRegistry() []Agent {
	mechanics := Mechanics()
	roles := Roles()
	agents := make([]Agent, 0, len(mechanics)*len(roles))
	for _, mechanic := range mechanics {
		for i, role := range roles {
			agents = append(agents, Agent{
				ID:          fmt.Sprintf("COAS-%02d-%02d", mechanic.ID, i+1),
				Mechanic:    mechanic,
				Role:        role,
				CanMerge:    true,
				RiskCeiling: "constitution-gated",
			})
		}
	}
	return agents
}

func ValidateRegistry(agents []Agent) error {
	if len(agents) != 100 {
		return fmt.Errorf("expected 100 agents, got %d", len(agents))
	}
	seen := make(map[string]struct{}, len(agents))
	for _, agent := range agents {
		if agent.ID == "" {
			return fmt.Errorf("agent id is required")
		}
		if !agent.CanMerge {
			return fmt.Errorf("%s is missing autonomous merge capability", agent.ID)
		}
		if _, ok := seen[agent.ID]; ok {
			return fmt.Errorf("duplicate agent id: %s", agent.ID)
		}
		seen[agent.ID] = struct{}{}
	}
	return nil
}

func LoadSnapshot(root, commit string) (Snapshot, error) {
	files := map[string]string{}
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
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".go", ".md", ".json", ".yml", ".yaml", ".txt":
		default:
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Size() > 1<<20 {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Commit: commit, Files: files}, nil
}

func RunSwarm(ctx context.Context, snapshot Snapshot) (SwarmReport, error) {
	agents := BuildRegistry()
	if err := ValidateRegistry(agents); err != nil {
		return SwarmReport{}, err
	}

	results := make(chan AgentResult, len(agents))
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
				results <- agent.Evaluate(snapshot)
			}
		}()
	}
	wg.Wait()
	close(results)

	report := SwarmReport{
		Version:      "coas-swarm.v1",
		SourceCommit: snapshot.Commit,
		AgentCount:   len(agents),
		Results:      make([]AgentResult, 0, len(agents)),
	}
	for result := range results {
		report.Results = append(report.Results, result)
	}
	sort.Slice(report.Results, func(i, j int) bool {
		return report.Results[i].AgentID < report.Results[j].AgentID
	})
	if err := ctx.Err(); err != nil {
		return SwarmReport{}, err
	}
	return report, nil
}

func (a Agent) Evaluate(snapshot Snapshot) AgentResult {
	corpus := strings.ToLower(joinSnapshot(snapshot))
	needle := strings.ToLower(a.Mechanic.Name)
	status := "attention"
	observation := fmt.Sprintf("%d tracked repository files inspected; mechanic needs stronger direct coverage", len(snapshot.Files))
	if strings.Contains(corpus, needle) {
		status = "covered"
		observation = fmt.Sprintf("%d tracked repository files inspected; mechanic is represented in repository text", len(snapshot.Files))
	}

	work := fmt.Sprintf("%s For %s: %s", a.Role.Directive, a.Mechanic.Name, a.Mechanic.Objective)
	return AgentResult{
		AgentID:       a.ID,
		Mechanic:      a.Mechanic.Name,
		Role:          a.Role.Name,
		Status:        status,
		Observation:   observation,
		ProposedWork:  work,
		AuthorityNote: "Agent may analyze, draft, test, update the runtime branch, and merge qualifying work under the MergeConstitution. Human-authored meaning remains attributable to the human source.",
	}
}

func joinSnapshot(snapshot Snapshot) string {
	keys := make([]string, 0, len(snapshot.Files))
	for key := range snapshot.Files {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, key := range keys {
		b.WriteString(key)
		b.WriteByte('\n')
		b.WriteString(snapshot.Files[key])
		b.WriteByte('\n')
	}
	return b.String()
}

func (r SwarmReport) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r SwarmReport) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# COAS Autonomous Swarm Report\n\n")
	fmt.Fprintf(&b, "- Version: %s\n", r.Version)
	fmt.Fprintf(&b, "- Source commit: %s\n", r.SourceCommit)
	fmt.Fprintf(&b, "- Logical agents: %d\n", r.AgentCount)
	fmt.Fprintf(&b, "- Merge authority: constitution-gated autonomous main merges are enabled\n\n")
	fmt.Fprintf(&b, "| Agent | Mechanic | Role | Status | Proposed work |\n")
	fmt.Fprintf(&b, "|---|---|---|---|---|\n")
	for _, result := range r.Results {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n",
			escapeTable(result.AgentID),
			escapeTable(result.Mechanic),
			escapeTable(result.Role),
			escapeTable(result.Status),
			escapeTable(result.ProposedWork),
		)
	}
	return b.String()
}

func escapeTable(value string) string {
	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\n", " ")
	return value
}
