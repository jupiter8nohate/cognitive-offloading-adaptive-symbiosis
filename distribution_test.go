package coas

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseDistributionTargetsRequiresAuthorization(t *testing.T) {
	raw := `[{"name":"owned","kind":"generic","url":"https://example.com/hook","authorized":false}]`
	if _, err := ParseDistributionTargets(raw); err == nil {
		t.Fatal("expected unauthorized target to be rejected")
	}
}

func TestParseDistributionTargetsLimitsFanout(t *testing.T) {
	targets := make([]DistributionTarget, 0, MaxDistributionTargets+1)
	for i := 0; i < MaxDistributionTargets+1; i++ {
		targets = append(targets, DistributionTarget{
			Name:       "target-" + string(rune('a'+i)),
			Kind:       DistributionKindGeneric,
			URL:        "https://example.com/hook",
			Authorized: true,
		})
	}
	data, err := json.Marshal(targets)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseDistributionTargets(string(data)); err == nil {
		t.Fatal("expected fanout limit error")
	}
}

func TestPublisherSendsOnlyExplicitTarget(t *testing.T) {
	var gotBody string
	var gotKey string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		gotBody = string(data)
		gotKey = r.Header.Get("Idempotency-Key")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	payload, err := BuildDistributionPayload("abc123")
	if err != nil {
		t.Fatal(err)
	}
	target := DistributionTarget{
		Name:       "owned-test-endpoint",
		Kind:       DistributionKindGeneric,
		URL:        server.URL,
		Authorized: true,
	}
	result, err := (Publisher{Client: server.Client()}).Publish(context.Background(), target, payload)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Published {
		t.Fatal("expected publication success")
	}
	if !strings.Contains(gotBody, "HUMAN_AGENCY") {
		t.Fatalf("payload missing manifesto: %s", gotBody)
	}
	if gotKey == "" {
		t.Fatal("expected idempotency key")
	}
}

func TestMastodonRequiresToken(t *testing.T) {
	err := ValidateDistributionTarget(DistributionTarget{
		Name:       "mastodon",
		Kind:       DistributionKindMastodon,
		URL:        "https://social.example/api/v1/statuses",
		Authorized: true,
	})
	if err == nil {
		t.Fatal("expected missing token error")
	}
}
