package coas

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type AutonomousMove struct {
	GoalID     string   `json:"goal_id"`
	GoalName   string   `json:"goal_name"`
	Score      int      `json:"score"`
	Rationale  string   `json:"rationale"`
	Evidence   []string `json:"evidence"`
	SourceHash string   `json:"source_hash"`
}

type AutonomousPlan struct {
	Version      string           `json:"version"`
	SourceCommit string           `json:"source_commit"`
	Selected     AutonomousMove   `json:"selected"`
	Candidates   []AutonomousMove `json:"candidates"`
}

func BuildAutonomousPlan(snapshot Snapshot) (AutonomousPlan, error) {
	if strings.TrimSpace(snapshot.Commit) == "" {
		return AutonomousPlan{}, fmt.Errorf("snapshot commit is required")
	}

	mission := CanonicalMission()
	candidates := make([]AutonomousMove, 0, len(mission.Goals))
	corpus := strings.ToLower(joinSnapshot(snapshot))

	for _, goal := range mission.Goals {
		if !goal.Autonomous {
			continue
		}
		h := sha256.Sum256([]byte(snapshot.Commit + ":" + goal.ID + ":" + selfAuthorFingerprint(snapshot)))
		score := 50 + int(h[0])%31
		for _, evidence := range goal.Evidence {
			if strings.Contains(corpus, strings.ToLower(evidence)) {
				score += 5
			} else {
				score += 12
			}
		}
		if score > 100 {
			score = 100
		}
		candidates = append(candidates, AutonomousMove{
			GoalID:     goal.ID,
			GoalName:   goal.Name,
			Score:      score,
			Rationale:  goal.Directive,
			Evidence:   append([]string(nil), goal.Evidence...),
			SourceHash: hex.EncodeToString(h[:8]),
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			return candidates[i].GoalID < candidates[j].GoalID
		}
		return candidates[i].Score > candidates[j].Score
	})
	if len(candidates) == 0 {
		return AutonomousPlan{}, fmt.Errorf("mission contains no autonomous goals")
	}

	return AutonomousPlan{
		Version:      "coas-autonomous-plan.v1",
		SourceCommit: snapshot.Commit,
		Selected:     candidates[0],
		Candidates:   candidates,
	}, nil
}

func (p AutonomousPlan) JSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

func (p AutonomousPlan) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", GLITCHOLOGYBanner("N⃟ E⃟ X⃟ T⃟_M⃟ O⃟ V⃟ E⃟"))
	fmt.Fprintf(&b, "SOURCE_COMMIT: `%s`\n\n", p.SourceCommit)
	fmt.Fprintf(&b, "SELECTED_GOAL: **%s** (`%s`)\n\n", p.Selected.GoalName, p.Selected.GoalID)
	fmt.Fprintf(&b, "SCORE: **%d**\n\n", p.Selected.Score)
	fmt.Fprintf(&b, "~~~text\n%s\n~~~\n\n", GLITCHOLOGYStatement("⁇", "G8", p.Selected.GoalID, "SELECTED", "AUTONOMOUS_EXECUTION"))
	fmt.Fprintf(&b, "%s\n\n", p.Selected.Rationale)
	b.WriteString("## Candidate board\n\n")
	b.WriteString("| Goal | Score | Source hash |\n|---|---:|---|\n")
	for _, move := range p.Candidates {
		fmt.Fprintf(&b, "| %s | %d | `%s` |\n", move.GoalName, move.Score, move.SourceHash)
	}
	return b.String()
}
