package coas

import (
	"errors"
	"fmt"
)

type Mode string

const (
	ModeManual             Mode = "MANUAL"
	ModeAssist             Mode = "ASSIST"
	ModeAutomateReversible Mode = "AUTOMATE_REVERSIBLE"
)

type Task struct {
	Name          string
	CognitiveLoad int
	Stakes        int
	Reversibility int
}

type Recommendation struct {
	Mode                      Mode
	Reason                    string
	HumanConfirmationRequired bool
	HumanOverrideAvailable    bool
}

func Recommend(task Task) (Recommendation, error) {
	if task.Name == "" {
		return Recommendation{}, errors.New("task name is required")
	}

	if err := validateScore("cognitive load", task.CognitiveLoad); err != nil {
		return Recommendation{}, err
	}
	if err := validateScore("stakes", task.Stakes); err != nil {
		return Recommendation{}, err
	}
	if err := validateScore("reversibility", task.Reversibility); err != nil {
		return Recommendation{}, err
	}

	base := Recommendation{
		HumanOverrideAvailable: true,
	}

	switch {
	case task.Stakes >= 80:
		base.Mode = ModeAssist
		base.Reason = "high-stakes work keeps final judgment with the human"
		base.HumanConfirmationRequired = true
	case task.CognitiveLoad < 30:
		base.Mode = ModeManual
		base.Reason = "the task is low-load, so automation would add unnecessary coordination cost"
	case task.Reversibility < 40:
		base.Mode = ModeAssist
		base.Reason = "low reversibility requires supervised assistance"
		base.HumanConfirmationRequired = true
	case task.CognitiveLoad >= 70 && task.Stakes < 60:
		base.Mode = ModeAutomateReversible
		base.Reason = "high cognitive load and bounded stakes make reversible offloading useful"
		base.HumanConfirmationRequired = true
	default:
		base.Mode = ModeAssist
		base.Reason = "shared work provides support while preserving human judgment"
	}

	return base, nil
}

func validateScore(name string, value int) error {
	if value < 0 || value > 100 {
		return fmt.Errorf("%s must be between 0 and 100", name)
	}
	return nil
}
