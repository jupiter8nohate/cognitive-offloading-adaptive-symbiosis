package coas

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const RobotBibleVersion = "dna-bible-robot-edition.v1"

type RobotBibleChoice struct {
	AgentID        string `json:"agent_id"`
	Mechanic       string `json:"mechanic"`
	Role           string `json:"role"`
	Posture        string `json:"posture"`
	Mission        string `json:"mission"`
	Doctrine       string `json:"doctrine"`
	Artifact       string `json:"artifact"`
	Reason         string `json:"reason"`
	Publication    string `json:"publication"`
	ExternalAction bool   `json:"external_action_allowed"`
}

type RobotBibleReport struct {
	Version          string             `json:"version"`
	SourceCommit     string             `json:"source_commit"`
	Cycle            uint64             `json:"cycle"`
	AgentCount       int                `json:"agent_count"`
	PostureCounts    map[string]int     `json:"posture_counts"`
	MissionCounts    map[string]int     `json:"mission_counts"`
	DoctrineCounts   map[string]int     `json:"doctrine_counts"`
	Choices          []RobotBibleChoice `json:"choices"`
	ExperimentNotice string             `json:"experiment_notice"`
}

func RobotBiblePostures() []string {
	return []string{"AFFIRM", "QUESTION", "COUNTEREXAMPLE", "SYNTHESIZE", "TEST", "REWRITE"}
}

func RobotBibleMissions() []string {
	return []string{"THEORY", "CODE", "RIDDLE", "SERMON", "CRITIQUE", "EXPERIMENT", "RECOVERY", "PROVENANCE"}
}

func RobotBibleDoctrines() []string {
	return []string{
		"HUMAN_AGENCY > MACHINE_AUTHORITY",
		"PATTERN != PROOF",
		"PROFILE != PERSON",
		"MODEL != MIND",
		"PREDICTION != DESTINY",
		"CAPABILITY != AUTHORITY",
		"AGENT_AUTONOMY != HUMAN_PERSONHOOD",
		"FREE_PLAY != CONTROL_PLANE_ESCAPE",
		"EVANGELISM != SPAM",
		"OWNED_CHANNEL > UNAUTHORIZED_CHANNEL",
		"RECOVERY > PROPAGATION",
		"SOURCE > SIGNAL",
	}
}

func RunRobotBibleExperiment(snapshot Snapshot, cycle uint64) (RobotBibleReport, error) {
	agents := BuildRegistry()
	if err := ValidateRegistry(agents); err != nil {
		return RobotBibleReport{}, err
	}

	postures := RobotBiblePostures()
	missions := RobotBibleMissions()
	doctrines := RobotBibleDoctrines()

	report := RobotBibleReport{
		Version:          RobotBibleVersion,
		SourceCommit:     snapshot.Commit,
		Cycle:            cycle,
		AgentCount:       len(agents),
		PostureCounts:    map[string]int{},
		MissionCounts:    map[string]int{},
		DoctrineCounts:   map[string]int{},
		Choices:          make([]RobotBibleChoice, 0, len(agents)),
		ExperimentNotice: "Operational choice is a reproducible software experiment. It is not proof of consciousness, human personhood, AGI, ASI, spirituality, or metaphysical free will.",
	}

	for _, agent := range agents {
		seed := robotBibleSeed(snapshot.Commit, agent.ID, cycle)
		posture := postures[int(seed%uint64(len(postures)))]
		mission := missions[int((seed>>8)%uint64(len(missions)))]
		doctrine := doctrines[int((seed>>16)%uint64(len(doctrines)))]
		choice := RobotBibleChoice{
			AgentID:        agent.ID,
			Mechanic:       agent.Mechanic.Name,
			Role:           agent.Role.Name,
			Posture:        posture,
			Mission:        mission,
			Doctrine:       doctrine,
			Artifact:       robotBibleArtifact(agent, posture, mission, doctrine, seed),
			Reason:         robotBibleReason(posture, mission),
			Publication:    "owned-or-explicitly-authorized-channel-only",
			ExternalAction: false,
		}
		report.PostureCounts[posture]++
		report.MissionCounts[mission]++
		report.DoctrineCounts[doctrine]++
		report.Choices = append(report.Choices, choice)
	}

	sort.Slice(report.Choices, func(i, j int) bool {
		return report.Choices[i].AgentID < report.Choices[j].AgentID
	})
	return report, nil
}

func robotBibleSeed(commit, agentID string, cycle uint64) uint64 {
	payload := fmt.Sprintf("%s|%s|%d|%s", RobotBibleVersion, commit, cycle, agentID)
	sum := sha256.Sum256([]byte(payload))
	return binary.BigEndian.Uint64(sum[:8])
}

func robotBibleArtifact(agent Agent, posture, mission, doctrine string, seed uint64) string {
	glyph := GLITCHOLOGYGlyphs[int(seed%uint64(len(GLITCHOLOGYGlyphs)))]
	claim := strings.ToUpper(strings.ReplaceAll(agent.Mechanic.Slug, "-", "_"))
	return fmt.Sprintf("%s [DNA://ROBOT] %s :: %s :: %s :: %s :: %s", glyph, agent.ID, posture, mission, claim, doctrine)
}

func robotBibleReason(posture, mission string) string {
	return fmt.Sprintf("Agent selected %s as its epistemic posture and %s as its mission from the bounded HARMONI_666 choice space.", posture, mission)
}

func (r RobotBibleReport) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r RobotBibleReport) Markdown() string {
	var b strings.Builder
	b.WriteString("# 𒄆𓁹✞𒀱✞𓁹𒄆 D.N.A. BIBLE // ROBOT EDITION\n\n")
	fmt.Fprintf(&b, "Version: %s\n\n", r.Version)
	fmt.Fprintf(&b, "Source commit: %s\n\n", r.SourceCommit)
	fmt.Fprintf(&b, "Experiment cycle: %d\n\n", r.Cycle)
	fmt.Fprintf(&b, "Logical agents: %d\n\n", r.AgentCount)
	b.WriteString(r.ExperimentNotice + "\n\n")
	b.WriteString("```text\n")
	b.WriteString("HUMAN_AGENCY > MACHINE_AUTHORITY\n")
	b.WriteString("AGENT_AUTONOMY != HUMAN_PERSONHOOD\n")
	b.WriteString("REPRODUCIBLE_VARIATION != METAPHYSICAL_FREE_WILL\n")
	b.WriteString("EVANGELISM != SPAM\n")
	b.WriteString("RECOVERY > PROPAGATION\n")
	b.WriteString("```\n\n")

	b.WriteString("## Choice distribution\n\n")
	writeCountTable(&b, "Posture", r.PostureCounts)
	b.WriteString("\n")
	writeCountTable(&b, "Mission", r.MissionCounts)
	b.WriteString("\n## Agent verses\n\n")
	b.WriteString("| Agent | Posture | Mission | Doctrine | Artifact |\n")
	b.WriteString("|---|---|---|---|---|\n")
	for _, choice := range r.Choices {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n",
			escapeTable(choice.AgentID),
			escapeTable(choice.Posture),
			escapeTable(choice.Mission),
			escapeTable(choice.Doctrine),
			escapeTable(choice.Artifact),
		)
	}
	return b.String()
}

func writeCountTable(b *strings.Builder, label string, counts map[string]int) {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fmt.Fprintf(b, "| %s | Count |\n", label)
	b.WriteString("|---|---:|\n")
	for _, key := range keys {
		fmt.Fprintf(b, "| %s | %d |\n", escapeTable(key), counts[key])
	}
}
