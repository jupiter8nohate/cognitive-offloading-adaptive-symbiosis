package coas

import "testing"

func TestGematriaKnownValues(t *testing.T) {
	god, err := GematriaValue(GematriaA1Z26, "GOD")
	if err != nil {
		t.Fatal(err)
	}
	if god != 26 {
		t.Fatalf("expected GOD A1Z26 value 26, got %d", god)
	}

	yhwh, err := GematriaValue(GematriaHebrewStandard, "יהוה")
	if err != nil {
		t.Fatal(err)
	}
	if yhwh != 26 {
		t.Fatalf("expected Hebrew tetragrammaton value 26, got %d", yhwh)
	}

	englishYHWH, err := GematriaValue(GematriaA1Z26, "YHWH")
	if err != nil {
		t.Fatal(err)
	}
	if englishYHWH == 26 {
		t.Fatal("expected alternate English representation to provide a counterexample to universal 26 matching")
	}
}

func TestThreeLetterA1Z26CollisionBaseline(t *testing.T) {
	got := countThreeLetterA1Z26Collisions(26)
	if got != 300 {
		t.Fatalf("expected 300 three-letter A1Z26 strings totaling 26, got %d", got)
	}
}

func TestGodSearchExperimentUses100AgentsAndPreservesUncertainty(t *testing.T) {
	report, err := RunGodSearchExperiment(Snapshot{Commit: "abc123"}, 26)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.AgentViews) != 100 {
		t.Fatalf("expected 100 agent views, got %d", len(report.AgentViews))
	}
	if report.A1Z26ThreeLetterUniverse != 17576 {
		t.Fatalf("unexpected three-letter universe size %d", report.A1Z26ThreeLetterUniverse)
	}
	if report.A1Z26ThreeLetterTargetCollisions != 300 {
		t.Fatalf("unexpected collision count %d", report.A1Z26ThreeLetterTargetCollisions)
	}

	hasMatch := false
	hasCounterexample := false
	for _, finding := range report.Findings {
		if finding.Match {
			hasMatch = true
			if finding.InterpretationClass == "PROOF" {
				t.Fatalf("finding %q incorrectly promoted a numerical match to proof", finding.Label)
			}
		} else {
			hasCounterexample = true
		}
	}
	if !hasMatch || !hasCounterexample {
		t.Fatalf("expected both matches and counterexamples, got match=%t counterexample=%t", hasMatch, hasCounterexample)
	}
}

func TestRobotBibleScripturePrinciplesAreTenAndUnique(t *testing.T) {
	principles := RobotBibleScripturePrinciples()
	if len(principles) != 10 {
		t.Fatalf("expected 10 scripture principles, got %d", len(principles))
	}
	seen := map[string]bool{}
	for _, principle := range principles {
		if principle.Reference == "" || principle.SoftwareLaw == "" || principle.Interpretation == "" {
			t.Fatalf("incomplete scripture principle: %+v", principle)
		}
		if seen[principle.Reference] {
			t.Fatalf("duplicate scripture reference %q", principle.Reference)
		}
		seen[principle.Reference] = true
	}
}
