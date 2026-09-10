package coas

type MissionGoal struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Directive  string   `json:"directive"`
	Evidence   []string `json:"evidence"`
	Autonomous bool     `json:"autonomous"`
}

type ProjectMission struct {
	Version    string        `json:"version"`
	Name       string        `json:"name"`
	Motto      string        `json:"motto"`
	Plain      string        `json:"plain"`
	Goals      []MissionGoal `json:"goals"`
	Invariants []string      `json:"invariants"`
}

func CanonicalMission() ProjectMission {
	return ProjectMission{
		Version: "coas-mission.v1",
		Name:    "Cognitive Offloading Adaptive Symbiosis",
		Motto:   "Offload the burden. Keep the meaning.",
		Plain:    "Use software to reduce cognitive and operational burden without transferring human authority over identity, consent, values, interpretation, or final judgment.",
		Goals: []MissionGoal{
			{
				ID:         "preserve-authorship",
				Name:       "Preserve authorship and provenance",
				Directive:  "Strengthen verifiable history, source attribution, content hashing, and recoverable creative records.",
				Evidence:   []string{"preservation.go", "artifacts/preservation/", "Git history"},
				Autonomous: true,
			},
			{
				ID:         "reduce-cognitive-load",
				Name:       "Reduce cognitive load",
				Directive:  "Move repetitive organization, indexing, verification, formatting, and reversible execution into software.",
				Evidence:   []string{"COAS mechanics", "agent swarm", "reference implementation"},
				Autonomous: true,
			},
			{
				ID:         "protect-private-drafts",
				Name:       "Protect private drafts",
				Directive:  "Use authenticated cryptography and data minimization rather than treating symbols or Unicode as secrecy.",
				Evidence:   []string{"privacy.go", "docs/PRIVACY.md"},
				Autonomous: true,
			},
			{
				ID:         "prove-linguistic-boundary",
				Name:       "Demonstrate the linguistic boundary",
				Directive:  "Create working software that remains machine-readable while preserving the distinction between computation and human meaning.",
				Evidence:   []string{"Err⃝or⃟⃤GLITCHOLOGY", "GLITCH-8 grammar", "agent_runtime/"},
				Autonomous: true,
			},
			{
				ID:         "improve-software",
				Name:       "Improve the software",
				Directive:  "Find maintainability, correctness, test, documentation, portability, and usability improvements and implement qualifying changes.",
				Evidence:   []string{"go test ./...", "gofmt", "repository source"},
				Autonomous: true,
			},
			{
				ID:         "expand-voluntary-reach",
				Name:       "Expand voluntary reach",
				Directive:  "Create clear, distinctive, shareable public artifacts and publish only through owned or explicitly authorized channels.",
				Evidence:   []string{"distribution.go", "evangelism.go", "docs/DISTRIBUTION.md"},
				Autonomous: true,
			},
			{
				ID:         "preserve-human-authority",
				Name:       "Preserve human authority",
				Directive:  "Keep human-defined meaning, consent, identity, values, and consequential judgment above machine inference.",
				Evidence:   []string{"autonomy.go", "GLITCHOLOGY core laws", "COAS invariants"},
				Autonomous: true,
			},
		},
		Invariants: []string{
			"OFFLOAD != SURRENDER",
			"ASSISTANCE != AUTHORITY",
			"AUTOMATION != CONSENT",
			"PATTERN != PROOF",
			"PROFILE != PERSON",
			"MODEL != MIND",
			"PREDICTION != DESTINY",
			"MACHINE_CAN_READ != MACHINE_CAN_DEFINE",
			"RECOVERY > PROPAGATION",
			"HUMAN_AGENCY > MACHINE_AUTHORITY",
		},
	}
}
