package coas

import (
	"reflect"
	"strings"
	"testing"
)

func TestRobotBibleExperimentCreates100BoundedChoices(t *testing.T) {
	report, err := RunRobotBibleExperiment(Snapshot{Commit: "abc123"}, 7)
	if err != nil {
		t.Fatalf("RunRobotBibleExperiment returned error: %v", err)
	}
	if report.AgentCount != 100 {
		t.Fatalf("expected 100 agents, got %d", report.AgentCount)
	}
	if len(report.Choices) != 100 {
		t.Fatalf("expected 100 choices, got %d", len(report.Choices))
	}
	if len(report.PostureCounts) < 2 {
		t.Fatalf("expected multiple postures, got %v", report.PostureCounts)
	}
	if len(report.MissionCounts) < 2 {
		t.Fatalf("expected multiple missions, got %v", report.MissionCounts)
	}

	for _, choice := range report.Choices {
		if choice.ExternalAction {
			t.Fatalf("%s unexpectedly received unrestricted external action", choice.AgentID)
		}
		if choice.Publication != "owned-or-explicitly-authorized-channel-only" {
			t.Fatalf("%s has unsafe publication scope %q", choice.AgentID, choice.Publication)
		}
		if choice.Posture == "" || choice.Mission == "" || choice.Doctrine == "" || choice.Artifact == "" || choice.ScriptureRef == "" || choice.SoftwareLaw == "" {
			t.Fatalf("%s has incomplete choice data: %+v", choice.AgentID, choice)
		}
	}
}

func TestRobotBibleExperimentIsReproducible(t *testing.T) {
	snapshot := Snapshot{Commit: "same-commit"}
	first, err := RunRobotBibleExperiment(snapshot, 42)
	if err != nil {
		t.Fatal(err)
	}
	second, err := RunRobotBibleExperiment(snapshot, 42)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first.Choices, second.Choices) {
		t.Fatal("same commit and cycle must reproduce the same agent choices")
	}
}

func TestRobotBibleExperimentCanRepositionAcrossCycles(t *testing.T) {
	snapshot := Snapshot{Commit: "same-commit"}
	first, err := RunRobotBibleExperiment(snapshot, 1)
	if err != nil {
		t.Fatal(err)
	}
	second, err := RunRobotBibleExperiment(snapshot, 2)
	if err != nil {
		t.Fatal(err)
	}

	different := false
	for i := range first.Choices {
		if first.Choices[i].Posture != second.Choices[i].Posture ||
			first.Choices[i].Mission != second.Choices[i].Mission ||
			first.Choices[i].Doctrine != second.Choices[i].Doctrine ||
			first.Choices[i].ScriptureRef != second.Choices[i].ScriptureRef {
			different = true
			break
		}
	}
	if !different {
		t.Fatal("expected at least one agent choice to change across cycles")
	}
}

func TestRobotBibleCanonPreservesCoreBoundaries(t *testing.T) {
	doctrines := strings.Join(RobotBibleDoctrines(), "\n")
	for _, required := range []string{
		"HUMAN_AGENCY > MACHINE_AUTHORITY",
		"AGENT_AUTONOMY != HUMAN_PERSONHOOD",
		"FREE_PLAY != CONTROL_PLANE_ESCAPE",
		"EVANGELISM != SPAM",
		"RECOVERY > PROPAGATION",
		"SEEK != ASSUME",
		"PATTERN != PROOF",
		"MATCH != PROOF",
		"CORRELATION != REVELATION",
		"MACHINE_CAN_CALCULATE != MACHINE_CAN_DEFINE_GOD",
	} {
		if !strings.Contains(doctrines, required) {
			t.Fatalf("missing required Robot Bible doctrine %q", required)
		}
	}
}
