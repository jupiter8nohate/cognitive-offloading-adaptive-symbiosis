package coas

import "testing"

func TestCanonicalMissionPreservesCoreAuthority(t *testing.T) {
	mission := CanonicalMission()
	if mission.Motto != "Offload the burden. Keep the meaning." {
		t.Fatalf("unexpected motto: %q", mission.Motto)
	}
	found := false
	for _, invariant := range mission.Invariants {
		if invariant == "HUMAN_AGENCY > MACHINE_AUTHORITY" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("canonical mission is missing human authority invariant")
	}
}

func TestBuildAutonomousPlanSelectsMissionGoal(t *testing.T) {
	snapshot := Snapshot{
		Commit: "abc123",
		Files: map[string]string{
			"README.md": "# COAS",
		},
	}
	plan, err := BuildAutonomousPlan(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Selected.GoalID == "" || plan.Selected.Score < 1 {
		t.Fatalf("invalid selected move: %#v", plan.Selected)
	}
	if len(plan.Candidates) != len(CanonicalMission().Goals) {
		t.Fatalf("candidate count = %d, want %d", len(plan.Candidates), len(CanonicalMission().Goals))
	}
}

func TestBuildAutonomousPlanIsDeterministic(t *testing.T) {
	snapshot := Snapshot{
		Commit: "abc123",
		Files: map[string]string{
			"README.md": "# COAS",
			"swarm.go":  "package coas",
		},
	}
	first, err := BuildAutonomousPlan(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildAutonomousPlan(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if first.Selected.GoalID != second.Selected.GoalID ||
		first.Selected.Score != second.Selected.Score ||
		first.Selected.SourceHash != second.Selected.SourceHash {
		t.Fatalf("selected moves differ: %#v != %#v", first.Selected, second.Selected)
	}
}
