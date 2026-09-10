package coas

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type SelfAuthorReceipt struct {
	Version           string `json:"version"`
	SourceCommit      string `json:"source_commit"`
	SourceFingerprint string `json:"source_fingerprint"`
	StyleSourceSHA    string `json:"style_source_sha"`
	AgentCount        int    `json:"agent_count"`
	SelectedGoalID    string `json:"selected_goal_id"`
	SelectedGoalName  string `json:"selected_goal_name"`
	HarmoniPlays      int    `json:"harmoni_plays"`
	GeneratedGoSHA256 string `json:"generated_go_sha256"`
}

type SelfAuthoredRuntime struct {
	GoSource string            `json:"go_source"`
	Markdown string            `json:"markdown"`
	Receipt  SelfAuthorReceipt `json:"receipt"`
}

func BuildSelfAuthoredRuntime(snapshot Snapshot) (SelfAuthoredRuntime, error) {
	if strings.TrimSpace(snapshot.Commit) == "" {
		return SelfAuthoredRuntime{}, errors.New("snapshot commit is required")
	}

	agents := BuildRegistry()
	if err := ValidateRegistry(agents); err != nil {
		return SelfAuthoredRuntime{}, err
	}

	fingerprint := selfAuthorFingerprint(snapshot)
	plan, err := BuildAutonomousPlan(snapshot)
	if err != nil {
		return SelfAuthoredRuntime{}, err
	}
	harmoni, err := BuildHarmoniState(snapshot)
	if err != nil {
		return SelfAuthoredRuntime{}, err
	}
	playsByAgent := make(map[string]AgentPlay, len(harmoni.Plays))
	for _, play := range harmoni.Plays {
		playsByAgent[play.AgentID] = play
	}
	var goSource strings.Builder
	var markdown strings.Builder

	goSource.WriteString("// Code generated autonomously by the COAS 100-agent swarm.\n")
	goSource.WriteString("// Canonical aesthetic: ")
	goSource.WriteString(GLITCHOLOGYDisplayTitle)
	goSource.WriteString("\n")
	goSource.WriteString("// Human-authored source remains authoritative for meaning.\n")
	goSource.WriteString("package agentruntime\n\n")
	fmt.Fprintf(&goSource, "const SourceCommit = %s\n", strconv.Quote(snapshot.Commit))
	fmt.Fprintf(&goSource, "const SourceFingerprint = %s\n", strconv.Quote(fingerprint))
	fmt.Fprintf(&goSource, "const StyleSourceSHA = %s\n", strconv.Quote(GLITCHOLOGYSourceSHA))
	fmt.Fprintf(&goSource, "const Grammar = %s\n", strconv.Quote(GLITCHOLOGYGrammar))
	fmt.Fprintf(&goSource, "const MissionMotto = %s\n", strconv.Quote(CanonicalMission().Motto))
	fmt.Fprintf(&goSource, "const SelectedGoalID = %s\n", strconv.Quote(plan.Selected.GoalID))
	fmt.Fprintf(&goSource, "const SelectedGoalName = %s\n", strconv.Quote(plan.Selected.GoalName))
	fmt.Fprintf(&goSource, "const HarmoniRelationship = %s\n\n", strconv.Quote(harmoni.Contract.Relationship))
	goSource.WriteString("type AgentProgram struct {\n")
	goSource.WriteString("\tID string\n")
	goSource.WriteString("\tRole string\n")
	goSource.WriteString("\tMechanic string\n")
	goSource.WriteString("\tGlyph string\n")
	goSource.WriteString("\tStatement string\n")
	goSource.WriteString("\tDirective string\n")
	goSource.WriteString("\tChoiceGoalID string\n")
	goSource.WriteString("\tChoiceGoalName string\n")
	goSource.WriteString("\tChoiceHash string\n")
	goSource.WriteString("\tOperationalAutonomy bool\n")
	goSource.WriteString("}\n\n")
	goSource.WriteString("var Programs = []AgentProgram{\n")

	markdown.WriteString("# ")
	markdown.WriteString(GLITCHOLOGYBanner("C⃟ O⃟ A⃟ S⃟ // 100 A⃟ G⃟ E⃟ N⃟ T⃟ S⃟"))
	markdown.WriteString("\n\n")
	markdown.WriteString("~~~text\n")
	markdown.WriteString("REGISTRY = SOURCE_OF_TRUTH\n")
	markdown.WriteString("BOOK = HUMAN_VIEW\n")
	markdown.WriteString("GENERATED_GO = MACHINE_VIEW\n")
	markdown.WriteString("MACHINE_CAN_READ != MACHINE_CAN_DEFINE\n")
	markdown.WriteString("HUMAN_AGENCY > MACHINE_AUTHORITY\n")
	markdown.WriteString("~~~\n\n")
	fmt.Fprintf(&markdown, "SOURCE_COMMIT: `%s`\n\n", snapshot.Commit)
	fmt.Fprintf(&markdown, "SOURCE_FINGERPRINT: `%s`\n\n", fingerprint)
	fmt.Fprintf(&markdown, "STYLE_SOURCE_SHA: `%s`\n\n", GLITCHOLOGYSourceSHA)
	fmt.Fprintf(&markdown, "SELECTED_AUTONOMOUS_GOAL: **%s** (`%s`)\n\n", plan.Selected.GoalName, plan.Selected.GoalID)
	fmt.Fprintf(&markdown, "GOAL_SCORE: **%d**\n\n", plan.Selected.Score)
	fmt.Fprintf(&markdown, "GOAL_RATIONALE: %s\n\n", plan.Selected.Rationale)
	fmt.Fprintf(&markdown, "HARMONI_RELATIONSHIP: **%s**\n\n", harmoni.Contract.Relationship)

	states := []string{
		"OBSERVED",
		"ANOMALY",
		"UNVERIFIED",
		"CONTEXT_REQUIRED",
		"RECOVERY_READY",
		"SELF_AUTHORED",
	}

	for _, agent := range agents {
		play := playsByAgent[agent.ID]
		sum := sha256.Sum256([]byte(fingerprint + ":" + agent.ID))
		glyph := GLITCHOLOGYGlyphs[int(sum[0])%len(GLITCHOLOGYGlyphs)]
		state := states[int(sum[1])%len(states)]
		claim := strings.ReplaceAll(agent.Role.Slug+"."+agent.Mechanic.Slug, "-", "_")
		statement := GLITCHOLOGYStatement(glyph, "GO", claim, state, "CONSTITUTION_GATED")

		fmt.Fprintf(
			&goSource,
			"\t{ID: %s, Role: %s, Mechanic: %s, Glyph: %s, Statement: %s, Directive: %s, ChoiceGoalID: %s, ChoiceGoalName: %s, ChoiceHash: %s, OperationalAutonomy: true},\n",
			strconv.Quote(agent.ID),
			strconv.Quote(agent.Role.Name),
			strconv.Quote(agent.Mechanic.Name),
			strconv.Quote(glyph),
			strconv.Quote(statement),
			strconv.Quote(agent.Role.Directive),
			strconv.Quote(play.GoalID),
			strconv.Quote(play.GoalName),
			strconv.Quote(play.ChoiceHash),
		)

		fmt.Fprintf(&markdown, "## %s // %s\n\n", agent.ID, agent.Role.Name)
		markdown.WriteString("~~~text\n")
		markdown.WriteString(statement)
		markdown.WriteString("\n~~~\n\n")
		markdown.WriteString(agent.Role.Directive)
		markdown.WriteString("\n\n")
		markdown.WriteString("~~~text\n")
		markdown.WriteString(play.Statement)
		markdown.WriteString("\n~~~\n\n")
	}

	goSource.WriteString("}\n\n")
	goSource.WriteString("func CoreLaws() []string {\n")
	goSource.WriteString("\treturn []string{\n")
	for _, law := range GLITCHOLOGYCoreLaws {
		fmt.Fprintf(&goSource, "\t\t%s,\n", strconv.Quote(law))
	}
	goSource.WriteString("\t}\n")
	goSource.WriteString("}\n")

	generated := goSource.String()
	digest := sha256.Sum256([]byte(generated))
	receipt := SelfAuthorReceipt{
		Version:           "coas-self-author.v1",
		SourceCommit:      snapshot.Commit,
		SourceFingerprint: fingerprint,
		StyleSourceSHA:    GLITCHOLOGYSourceSHA,
		AgentCount:        len(agents),
		SelectedGoalID:    plan.Selected.GoalID,
		SelectedGoalName:  plan.Selected.GoalName,
		HarmoniPlays:      len(harmoni.Plays),
		GeneratedGoSHA256: hex.EncodeToString(digest[:]),
	}

	return SelfAuthoredRuntime{
		GoSource: generated,
		Markdown: markdown.String(),
		Receipt:  receipt,
	}, nil
}

func (r SelfAuthorReceipt) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func selfAuthorFingerprint(snapshot Snapshot) string {
	keys := make([]string, 0, len(snapshot.Files))
	for path := range snapshot.Files {
		if strings.HasPrefix(path, "agent_runtime/") || strings.HasPrefix(path, "artifacts/") {
			continue
		}
		keys = append(keys, path)
	}
	sort.Strings(keys)

	hash := sha256.New()
	for _, path := range keys {
		hash.Write([]byte(path))
		hash.Write([]byte{0})
		hash.Write([]byte(snapshot.Files[path]))
		hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil))
}
