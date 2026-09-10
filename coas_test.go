package coas

import "testing"

func TestRecommend(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want Mode
	}{
		{
			name: "high stakes stays assist",
			task: Task{Name: "medical decision", CognitiveLoad: 90, Stakes: 95, Reversibility: 20},
			want: ModeAssist,
		},
		{
			name: "low load stays manual",
			task: Task{Name: "simple note", CognitiveLoad: 10, Stakes: 5, Reversibility: 100},
			want: ModeManual,
		},
		{
			name: "low reversibility stays assist",
			task: Task{Name: "irreversible change", CognitiveLoad: 70, Stakes: 40, Reversibility: 10},
			want: ModeAssist,
		},
		{
			name: "high load bounded stakes can automate",
			task: Task{Name: "organize files", CognitiveLoad: 85, Stakes: 20, Reversibility: 90},
			want: ModeAutomateReversible,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Recommend(tt.task)
			if err != nil {
				t.Fatalf("Recommend() error = %v", err)
			}
			if got.Mode != tt.want {
				t.Fatalf("Recommend() mode = %s, want %s", got.Mode, tt.want)
			}
			if !got.HumanOverrideAvailable {
				t.Fatal("human override must always be available")
			}
		})
	}
}

func TestRecommendRejectsInvalidScores(t *testing.T) {
	_, err := Recommend(Task{
		Name:          "bad input",
		CognitiveLoad: 101,
		Stakes:        10,
		Reversibility: 10,
	})
	if err == nil {
		t.Fatal("expected invalid score error")
	}
}
