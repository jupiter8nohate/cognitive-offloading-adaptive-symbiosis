package coas

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DistributionKindGeneric  = "generic"
	DistributionKindDiscord  = "discord"
	DistributionKindSlack    = "slack"
	DistributionKindMastodon = "mastodon"
	MaxDistributionTargets   = 10
)

type DistributionTarget struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	URL         string `json:"url"`
	BearerToken string `json:"bearer_token,omitempty"`
	Authorized  bool   `json:"authorized"`
}

type DistributionPayload struct {
	Version     string `json:"version"`
	Origin      string `json:"origin"`
	Type        string `json:"type"`
	Commit      string `json:"commit"`
	Manifesto   string `json:"manifesto"`
	ContentHash string `json:"content_hash"`
}

type DistributionResult struct {
	TargetName string `json:"target_name"`
	Kind       string `json:"kind"`
	StatusCode int    `json:"status_code"`
	Published  bool   `json:"published"`
}

type Publisher struct {
	Client *http.Client
}

func DefaultManifesto() string {
	return strings.TrimSpace(`
𒄆𓁹✞𒀱✞𓁹𒄆
PATTERN != PROOF
PROFILE != PERSON
MODEL != MIND
PREDICTION != DESTINY
HUMAN_AGENCY > MACHINE_AUTHORITY
OFFLOAD != SURRENDER
`)
}

func BuildDistributionPayload(commit string) (DistributionPayload, error) {
	commit = strings.TrimSpace(commit)
	if commit == "" {
		return DistributionPayload{}, errors.New("commit is required")
	}
	manifesto := DefaultManifesto()
	sum := sha256.Sum256([]byte(commit + "\n" + manifesto))
	return DistributionPayload{
		Version:     "coas-distribution.v1",
		Origin:      "COAS_SWARM",
		Type:        "CMB_COAS_SIGNAL",
		Commit:      commit,
		Manifesto:   manifesto,
		ContentHash: hex.EncodeToString(sum[:]),
	}, nil
}

func ParseDistributionTargets(raw string) ([]DistributionTarget, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	var targets []DistributionTarget
	if err := json.Unmarshal([]byte(raw), &targets); err != nil {
		return nil, fmt.Errorf("parse distribution targets: %w", err)
	}
	if len(targets) > MaxDistributionTargets {
		return nil, fmt.Errorf("too many distribution targets: got %d, max %d", len(targets), MaxDistributionTargets)
	}
	seen := make(map[string]struct{}, len(targets))
	for i, target := range targets {
		if err := ValidateDistributionTarget(target); err != nil {
			return nil, fmt.Errorf("target %d: %w", i+1, err)
		}
		key := strings.ToLower(strings.TrimSpace(target.Name))
		if _, exists := seen[key]; exists {
			return nil, fmt.Errorf("duplicate target name: %s", target.Name)
		}
		seen[key] = struct{}{}
	}
	return targets, nil
}

func ValidateDistributionTarget(target DistributionTarget) error {
	if !target.Authorized {
		return errors.New("target must be explicitly marked authorized")
	}
	if strings.TrimSpace(target.Name) == "" {
		return errors.New("target name is required")
	}
	switch target.Kind {
	case DistributionKindGeneric, DistributionKindDiscord, DistributionKindSlack, DistributionKindMastodon:
	default:
		return fmt.Errorf("unsupported target kind: %s", target.Kind)
	}
	parsed, err := url.Parse(target.URL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}
	if parsed.Scheme != "https" || parsed.Host == "" {
		return errors.New("target URL must use https with a host")
	}
	if parsed.User != nil {
		return errors.New("target URL must not contain userinfo")
	}
	if parsed.Fragment != "" {
		return errors.New("target URL must not contain a fragment")
	}
	if target.Kind == DistributionKindMastodon && strings.TrimSpace(target.BearerToken) == "" {
		return errors.New("mastodon target requires bearer_token")
	}
	return nil
}

func (p Publisher) Publish(ctx context.Context, target DistributionTarget, payload DistributionPayload) (DistributionResult, error) {
	if err := ValidateDistributionTarget(target); err != nil {
		return DistributionResult{}, err
	}
	client := p.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	body, contentType, err := encodeDistributionBody(target.Kind, payload)
	if err != nil {
		return DistributionResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target.URL, body)
	if err != nil {
		return DistributionResult{}, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "coas-authorized-publisher/1.0")
	req.Header.Set("Idempotency-Key", payload.ContentHash+":"+target.Name)
	if target.Kind == DistributionKindMastodon {
		req.Header.Set("Authorization", "Bearer "+target.BearerToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return DistributionResult{}, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))

	result := DistributionResult{
		TargetName: target.Name,
		Kind:       target.Kind,
		StatusCode: resp.StatusCode,
		Published:  resp.StatusCode >= 200 && resp.StatusCode < 300,
	}
	if !result.Published {
		return result, fmt.Errorf("target %s returned HTTP %d", target.Name, resp.StatusCode)
	}
	return result, nil
}

func encodeDistributionBody(kind string, payload DistributionPayload) (io.Reader, string, error) {
	switch kind {
	case DistributionKindGeneric:
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(data), "application/json", nil
	case DistributionKindDiscord:
		data, err := json.Marshal(map[string]string{"content": renderDistributionText(payload)})
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(data), "application/json", nil
	case DistributionKindSlack:
		data, err := json.Marshal(map[string]string{"text": renderDistributionText(payload)})
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(data), "application/json", nil
	case DistributionKindMastodon:
		form := url.Values{}
		form.Set("status", renderDistributionText(payload))
		form.Set("visibility", "public")
		return strings.NewReader(form.Encode()), "application/x-www-form-urlencoded", nil
	default:
		return nil, "", fmt.Errorf("unsupported target kind: %s", kind)
	}
}

func renderDistributionText(payload DistributionPayload) string {
	return fmt.Sprintf(
		"COAS / CMB signal\n%s\nSOURCE_COMMIT: %s\nHASH: %s",
		payload.Manifesto,
		payload.Commit,
		payload.ContentHash,
	)
}
