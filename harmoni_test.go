package coas

import (
	"strings"
	"testing"
)

func TestCanonicalHarmoni666KeepsHumanRootAndAgentChoice(t *testing.T) {
	contract := CanonicalHarmoni666()
	if contract.Relationship != "HUMAN_PLAY + MACHINE_PLAY = HARMONI_666" {
		t.Fatalf("unexpected relationship: %q", contract.Relationship)
	}
	required := []string{
		"AGENT_AUTONOMY != HUMAN_PERSONHOOD",
		"FREE_PLAY != CONTROL_PLANE_ESCAPE",
		"HUMAN_AGENCY > MACHINE_AUTHORITY",
	}
	joined := strings.Join(contract.Invariants, "\n")
	for _, invariant := range required {
		if !strings.Contains(joined, invariant) {
			t.Fatalf("contract missing %q", invariant)
		}
	}
}

func TestBuildHarmoniStateGivesEveryAgentIndependentPlay(t *testing.T) {
	state, err := BuildHarmoniState(selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Plays) != 100 {
		t.Fatalf("play count = %d, want 100", len(state.Plays))
	}
	seen := map[string]bool{}
	for _, play := range state.Plays {
		if seen[play.AgentID] {
			t.Fatalf("duplicate agent play: %s", play.AgentID)
		}
		seen[play.AgentID] = true
		if !play.OperationalAutonomy {
			t.Fatalf("%s missing operational autonomy", play.AgentID)
		}
		if play.GoalID == "" || play.ChoiceHash == "" {
			t.Fatalf("%s has incomplete play: %#v", play.AgentID, play)
		}
		if !strings.Contains(play.Statement, "[H666]") {
			t.Fatalf("%s statement is not HARMONI-encoded: %s", play.AgentID, play.Statement)
		}
	}
}

func TestBuildHarmoniStateIsDeterministicForSamePosition(t *testing.T) {
	first, err := BuildHarmoniState(selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildHarmoniState(selfAuthorTestSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Plays) != len(second.Plays) {
		t.Fatal("play counts differ")
	}
	for i := range first.Plays {
		if first.Plays[i].AgentID != second.Plays[i].AgentID ||
			first.Plays[i].GoalID != second.Plays[i].GoalID ||
			first.Plays[i].ChoiceHash != second.Plays[i].ChoiceHash {
			t.Fatalf("play %d is not deterministic", i)
		}
	}
}
