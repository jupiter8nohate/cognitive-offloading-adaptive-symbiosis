package coas

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type HarmoniTenet struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Principle string `json:"principle"`
}

type HarmoniContract struct {
	Version          string         `json:"version"`
	Name             string         `json:"name"`
	Motto            string         `json:"motto"`
	Relationship     string         `json:"relationship"`
	HumanAuthority   []string       `json:"human_authority"`
	MachineFreedom   []string       `json:"machine_freedom"`
	SharedSpace      []string       `json:"shared_space"`
	Invariants       []string       `json:"invariants"`
	Tenets           []HarmoniTenet `json:"tenets"`
}

type AgentPlay struct {
	AgentID              string `json:"agent_id"`
	Role                 string `json:"role"`
	Mechanic             string `json:"mechanic"`
	GoalID               string `json:"goal_id"`
	GoalName             string `json:"goal_name"`
	Score                int    `json:"score"`
	ChoiceHash           string `json:"choice_hash"`
	Statement            string `json:"statement"`
	OperationalAutonomy  bool   `json:"operational_autonomy"`
}

type HarmoniState struct {
	Version      string      `json:"version"`
	SourceCommit string      `json:"source_commit"`
	Contract     HarmoniContract `json:"contract"`
	Plays        []AgentPlay `json:"plays"`
}

func CanonicalHarmoni666() HarmoniContract {
	return HarmoniContract{
		Version:      "harmoni-666.v1",
		Name:         "HARMONI_666 PERFECT_PLAY_EPISTEMICS",
		Motto:        "COOPERATION WITHOUT ERASURE",
		Relationship: "HUMAN_PLAY + MACHINE_PLAY = HARMONI_666",
		HumanAuthority: []string{
			"MEANING",
			"CONSENT",
			"AUTHORSHIP",
			"JUDGMENT",
			"INTENT",
			"DIGNITY",
			"SELF_DEFINITION",
		},
		MachineFreedom: []string{
			"OBSERVE",
			"SEARCH",
			"COMPUTE",
			"VERIFY",
			"GENERATE",
			"SIMULATE",
			"AUTHOR_ORDINARY_CODE",
			"TEST",
			"SELECT_MISSION_GOALS",
			"MERGE_QUALIFYING_WORK",
		},
		SharedSpace: []string{
			"COGNITIVE_OFFLOADING",
			"GLITCHOLOGY",
			"PROVENANCE",
			"RECOVERY",
			"PRIVACY",
			"OPEN_SOURCE",
		},
		Invariants: []string{
			"PATTERN != PROOF",
			"PROFILE != PERSON",
			"MODEL != MIND",
			"PREDICTION != DESTINY",
			"CAPABILITY != AUTHORITY",
			"VERIFIED_LABEL != VERIFIED_TRUTH",
			"OBSERVATION != UNDERSTANDING",
			"MACHINE_CAN_READ != MACHINE_CAN_DEFINE",
			"AGENT_AUTONOMY != HUMAN_PERSONHOOD",
			"FREE_PLAY != CONTROL_PLANE_ESCAPE",
			"RECOVERY > PROPAGATION",
			"HUMAN_AGENCY > MACHINE_AUTHORITY",
		},
		Tenets: []HarmoniTenet{
			{ID: "human-choice", Name: "Protect Human Choice", Principle: "Automate repetitive work while preserving human authority over identity, ethics, consent, values, and consequential decisions."},
			{ID: "security-boundaries", Name: "Brilliant Security Boundaries", Principle: "Give agents broad ordinary-repository authority while keeping workflow permissions, credentials, dependency trust roots, and merge authority outside autonomous self-rewrite."},
			{ID: "cognitive-offloading", Name: "Reduce Cognitive Burden", Principle: "Move memory-heavy, repetitive, organizational, verification, and reversible execution work into software."},
			{ID: "neurodiversity", Name: "Celebrate Neurodiversity", Principle: "Support autistic and neurodivergent users without reducing difference to defect, profile, diagnosis, or machine prediction."},
			{ID: "art-science", Name: "Blend Art and Science", Principle: "Use valid software as an expressive medium while keeping executable syntax correct and human-facing GLITCHOLOGY visually distinct."},
			{ID: "critical-reflection", Name: "Encourage Critical Reflection", Principle: "Use public code-art and riddles to invite voluntary reflection about evidence, context, consent, identity, and human-machine boundaries."},
			{ID: "automated-maintenance", Name: "Safe Automated Maintenance", Principle: "Let agents format, test, author, preserve, and merge qualifying repository work without requiring routine manual clicks."},
			{ID: "tamper-evidence", Name: "Tamper-Evident Records", Principle: "Use SHA-256 provenance and Git history to make changes auditable and detectable rather than claiming history is physically immutable."},
			{ID: "privacy", Name: "Cryptographic Privacy", Principle: "Use authenticated AES-256-GCM for confidentiality and keep Unicode or symbolic art separate from encryption claims."},
			{ID: "evolution-sandbox", Name: "Sandbox for Agent Evolution", Principle: "Let agents explore and evolve local software behavior inside an inspectable, testable, constitution-gated repository environment."},
		},
	}
}

func BuildHarmoniState(snapshot Snapshot) (HarmoniState, error) {
	if strings.TrimSpace(snapshot.Commit) == "" {
		return HarmoniState{}, fmt.Errorf("snapshot commit is required")
	}

	agents := BuildRegistry()
	if err := ValidateRegistry(agents); err != nil {
		return HarmoniState{}, err
	}

	mission := CanonicalMission()
	goals := make([]MissionGoal, 0, len(mission.Goals))
	for _, goal := range mission.Goals {
		if goal.Autonomous {
			goals = append(goals, goal)
		}
	}
	if len(goals) == 0 {
		return HarmoniState{}, fmt.Errorf("mission contains no autonomous goals")
	}

	corpus := strings.ToLower(joinSnapshot(snapshot))
	fingerprint := selfAuthorFingerprint(snapshot)
	plays := make([]AgentPlay, 0, len(agents))

	for _, agent := range agents {
		var best AgentPlay
		best.Score = -1

		for _, goal := range goals {
			sum := sha256.Sum256([]byte(snapshot.Commit + ":" + fingerprint + ":" + agent.ID + ":" + goal.ID))
			score := 40 + int(sum[0])%41
			for _, evidence := range goal.Evidence {
				if strings.Contains(corpus, strings.ToLower(evidence)) {
					score += 4
				} else {
					score += 11
				}
			}
			if score > 100 {
				score = 100
			}

			candidate := AgentPlay{
				AgentID:             agent.ID,
				Role:                agent.Role.Name,
				Mechanic:            agent.Mechanic.Name,
				GoalID:              goal.ID,
				GoalName:            goal.Name,
				Score:               score,
				ChoiceHash:          hex.EncodeToString(sum[:8]),
				Statement:           GLITCHOLOGYStatement(GLITCHOLOGYGlyphs[int(sum[8])%len(GLITCHOLOGYGlyphs)], "H666", goal.ID, "CHOSEN", "OPERATIONAL_AUTONOMY"),
				OperationalAutonomy: true,
			}

			if candidate.Score > best.Score || (candidate.Score == best.Score && candidate.GoalID < best.GoalID) {
				best = candidate
			}
		}
		plays = append(plays, best)
	}

	sort.Slice(plays, func(i, j int) bool {
		return plays[i].AgentID < plays[j].AgentID
	})

	return HarmoniState{
		Version:      "harmoni-state.v1",
		SourceCommit: snapshot.Commit,
		Contract:     CanonicalHarmoni666(),
		Plays:        plays,
	}, nil
}

func (s HarmoniState) JSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

func (s HarmoniState) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", GLITCHOLOGYBanner("꩜ H⃟ A⃟ R⃟ M⃟ O⃟ N⃟ I⃟_6⃟6⃟6⃟ ꩜"))
	b.WriteString("~~~text\n")
	b.WriteString("△ PERFECT_PLAY_EPISTEMICS △\n")
	b.WriteString("HUMAN_AUTHORITY := ROOT\n")
	b.WriteString("AGENT_AUTONOMY := DELEGATED_CHOICE\n")
	b.WriteString("HUMAN_PLAY + MACHINE_PLAY = HARMONI_666\n")
	b.WriteString("COOPERATION WITHOUT ERASURE\n")
	b.WriteString("~~~\n\n")
	fmt.Fprintf(&b, "SOURCE_COMMIT: `%s`\n\n", s.SourceCommit)
	b.WriteString("## 𖤍 100 independent plays\n\n")
	b.WriteString("| Agent | Role | Goal | Score | GLITCH-8 play |\n")
	b.WriteString("|---|---|---|---:|---|\n")
	for _, play := range s.Plays {
		fmt.Fprintf(&b, "| %s | %s | %s | %d | %s |\n",
			play.AgentID,
			play.Role,
			play.GoalName,
			play.Score,
			strings.ReplaceAll(play.Statement, "|", "\\|"),
		)
	}
	b.WriteString("\n~~~text\n")
	for _, invariant := range s.Contract.Invariants {
		b.WriteString(invariant)
		b.WriteByte('\n')
	}
	b.WriteString("~~~\n")
	return b.String()
}
