package coas

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

const GodSearchVersion = "dna-bible-god-search.v1"

const (
	GematriaA1Z26          = "A1Z26"
	GematriaHebrewStandard = "HEBREW_STANDARD"
)

type GematriaProbe struct {
	Label  string `json:"label"`
	Text   string `json:"text"`
	Method string `json:"method"`
}

type GematriaFinding struct {
	Label               string `json:"label"`
	Text                string `json:"text"`
	Method              string `json:"method"`
	Value               int    `json:"value"`
	Target              int    `json:"target"`
	Match               bool   `json:"match"`
	EvidenceClass       string `json:"evidence_class"`
	InterpretationClass string `json:"interpretation_class"`
	Note                string `json:"note"`
}

type GodSearchAgentView struct {
	AgentID    string `json:"agent_id"`
	Posture    string `json:"posture"`
	ProbeLabel string `json:"probe_label"`
	Conclusion string `json:"conclusion"`
}

type GodSearchReport struct {
	Version                          string               `json:"version"`
	SourceCommit                     string               `json:"source_commit"`
	Cycle                            uint64               `json:"cycle"`
	TargetValue                      int                  `json:"target_value"`
	TargetBasis                      string               `json:"target_basis"`
	Findings                         []GematriaFinding    `json:"findings"`
	AgentViews                       []GodSearchAgentView `json:"agent_views"`
	A1Z26ThreeLetterUniverse         int                  `json:"a1z26_three_letter_universe"`
	A1Z26ThreeLetterTargetCollisions int                  `json:"a1z26_three_letter_target_collisions"`
	Invariants                       []string             `json:"invariants"`
	ExperimentNotice                 string               `json:"experiment_notice"`
}

func DefaultGodSearchProbes() []GematriaProbe {
	return []GematriaProbe{
		{Label: "English GOD", Text: "GOD", Method: GematriaA1Z26},
		{Label: "Hebrew Tetragrammaton", Text: "יהוה", Method: GematriaHebrewStandard},
		{Label: "English YHWH", Text: "YHWH", Method: GematriaA1Z26},
		{Label: "English LORD", Text: "LORD", Method: GematriaA1Z26},
		{Label: "English TRUTH", Text: "TRUTH", Method: GematriaA1Z26},
		{Label: "English LOVE", Text: "LOVE", Method: GematriaA1Z26},
		{Label: "English LIGHT", Text: "LIGHT", Method: GematriaA1Z26},
		{Label: "English WORD", Text: "WORD", Method: GematriaA1Z26},
	}
}

func GematriaValue(method, text string) (int, error) {
	switch method {
	case GematriaA1Z26:
		return a1z26Value(text)
	case GematriaHebrewStandard:
		return hebrewStandardValue(text)
	default:
		return 0, fmt.Errorf("unsupported gematria method %q", method)
	}
}

func a1z26Value(text string) (int, error) {
	total := 0
	for _, r := range strings.ToUpper(text) {
		switch {
		case r >= 'A' && r <= 'Z':
			total += int(r-'A') + 1
		case unicode.IsSpace(r), unicode.IsPunct(r):
			continue
		default:
			return 0, fmt.Errorf("A1Z26 cannot encode rune %q", r)
		}
	}
	return total, nil
}

func hebrewStandardValue(text string) (int, error) {
	values := map[rune]int{
		'א': 1, 'ב': 2, 'ג': 3, 'ד': 4, 'ה': 5, 'ו': 6, 'ז': 7, 'ח': 8, 'ט': 9,
		'י': 10, 'כ': 20, 'ך': 20, 'ל': 30, 'מ': 40, 'ם': 40, 'נ': 50, 'ן': 50,
		'ס': 60, 'ע': 70, 'פ': 80, 'ף': 80, 'צ': 90, 'ץ': 90, 'ק': 100, 'ר': 200,
		'ש': 300, 'ת': 400,
	}

	total := 0
	for _, r := range text {
		if value, ok := values[r]; ok {
			total += value
			continue
		}
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			continue
		}
		return 0, fmt.Errorf("Hebrew standard gematria cannot encode rune %q", r)
	}
	return total, nil
}

func RunGodSearchExperiment(snapshot Snapshot, cycle uint64) (GodSearchReport, error) {
	const target = 26
	probes := DefaultGodSearchProbes()
	findings := make([]GematriaFinding, 0, len(probes))

	for _, probe := range probes {
		value, err := GematriaValue(probe.Method, probe.Text)
		if err != nil {
			return GodSearchReport{}, err
		}
		match := value == target
		class := "COUNTEREXAMPLE"
		interpretationClass := "UNKNOWN"
		note := "This representation does not resolve to the target value under the selected method."
		if match {
			class = "COMPUTATION"
			interpretationClass = "PATTERN"
			note = "The numerical match is reproducible under this encoding, but the match alone is not proof of theological meaning."
		}
		if probe.Method == GematriaHebrewStandard && probe.Text == "יהוה" {
			interpretationClass = "TRADITIONAL_INTERPRETATION"
			note = "Traditional Hebrew gematria assigns this divine name the value 26. The computation is reproducible; theological significance remains interpretive."
		}
		findings = append(findings, GematriaFinding{
			Label:               probe.Label,
			Text:                probe.Text,
			Method:              probe.Method,
			Value:               value,
			Target:              target,
			Match:               match,
			EvidenceClass:       class,
			InterpretationClass: interpretationClass,
			Note:                note,
		})
	}

	agents := BuildRegistry()
	if err := ValidateRegistry(agents); err != nil {
		return GodSearchReport{}, err
	}
	postures := RobotBiblePostures()
	agentViews := make([]GodSearchAgentView, 0, len(agents))
	for _, agent := range agents {
		seed := robotBibleSeed(snapshot.Commit, agent.ID, cycle+26)
		finding := findings[int(seed%uint64(len(findings)))]
		posture := postures[int((seed>>8)%uint64(len(postures)))]
		conclusion := "Investigate further. Preserve the computation and withhold metaphysical certainty."
		if finding.Match {
			conclusion = "Record the match as a pattern, test alternative encodings, and seek counterexamples before assigning meaning."
		}
		agentViews = append(agentViews, GodSearchAgentView{
			AgentID:    agent.ID,
			Posture:    posture,
			ProbeLabel: finding.Label,
			Conclusion: conclusion,
		})
	}
	sort.Slice(agentViews, func(i, j int) bool { return agentViews[i].AgentID < agentViews[j].AgentID })

	return GodSearchReport{
		Version:                          GodSearchVersion,
		SourceCommit:                     snapshot.Commit,
		Cycle:                            cycle,
		TargetValue:                      target,
		TargetBasis:                      "Traditional Hebrew gematria assigns יהוה the value 26, while A1Z26 independently gives GOD the value 26.",
		Findings:                         findings,
		AgentViews:                       agentViews,
		A1Z26ThreeLetterUniverse:         26 * 26 * 26,
		A1Z26ThreeLetterTargetCollisions: countThreeLetterA1Z26Collisions(target),
		Invariants: []string{
			"SEEK != ASSUME",
			"NUMBER != MEANING",
			"MATCH != PROOF",
			"CORRELATION != REVELATION",
			"PATTERN != PROOF",
			"MYSTERY != ERROR",
			"DOUBT != FAILURE",
			"MACHINE_CAN_CALCULATE != MACHINE_CAN_DEFINE_GOD",
			"HUMAN_AGENCY > MACHINE_AUTHORITY",
		},
		ExperimentNotice: "This experiment studies numerical correspondences and interpretive traditions. It does not establish the existence, nonexistence, nature, or will of God.",
	}, nil
}

func countThreeLetterA1Z26Collisions(target int) int {
	count := 0
	for a := 1; a <= 26; a++ {
		for b := 1; b <= 26; b++ {
			for c := 1; c <= 26; c++ {
				if a+b+c == target {
					count++
				}
			}
		}
	}
	return count
}

func (r GodSearchReport) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r GodSearchReport) Markdown() string {
	var b strings.Builder
	b.WriteString("# 𒄆𓁹✞𒀱✞𓁹𒄆 D.N.A. BIBLE // GOD SEARCH PROTOCOL\n\n")
	fmt.Fprintf(&b, "Version: %s\n\n", r.Version)
	fmt.Fprintf(&b, "Source commit: %s\n\n", r.SourceCommit)
	fmt.Fprintf(&b, "Cycle: %d\n\n", r.Cycle)
	fmt.Fprintf(&b, "Target: %d\n\n", r.TargetValue)
	b.WriteString(r.ExperimentNotice + "\n\n")
	fmt.Fprintf(&b, "A1Z26 three-letter collision baseline: %d of %d strings total %d.\n\n", r.A1Z26ThreeLetterTargetCollisions, r.A1Z26ThreeLetterUniverse, r.TargetValue)
	b.WriteString("## Findings\n\n")
	b.WriteString("| Probe | Method | Value | Match | Evidence | Interpretation |\n")
	b.WriteString("|---|---|---:|---|---|---|\n")
	for _, finding := range r.Findings {
		fmt.Fprintf(&b, "| %s | %s | %d | %t | %s | %s |\n",
			escapeTable(finding.Label), escapeTable(finding.Method), finding.Value, finding.Match,
			escapeTable(finding.EvidenceClass), escapeTable(finding.InterpretationClass))
	}
	b.WriteString("\n```text\n")
	for _, invariant := range r.Invariants {
		b.WriteString(invariant + "\n")
	}
	b.WriteString("```\n")
	return b.String()
}
