package coas

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRunIndexAuditUses100LogicalAgentsAndDeduplicatesQueries(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(CanonicalRepositoryURL)),
			Header:     make(http.Header),
		}, nil
	})}

	report, err := RunIndexAudit(context.Background(), "abc123", IndexProviderConfig{HTTPClient: client})
	if err != nil {
		t.Fatalf("RunIndexAudit returned error: %v", err)
	}
	if report.LogicalAgentCount != 100 {
		t.Fatalf("expected 100 logical agents, got %d", report.LogicalAgentCount)
	}
	if report.UniqueQueryCount != len(DiscoveryQueries()) {
		t.Fatalf("expected %d unique queries, got %d", len(DiscoveryQueries()), report.UniqueQueryCount)
	}
	if calls != len(DiscoveryQueries()) {
		t.Fatalf("expected %d network calls, got %d", len(DiscoveryQueries()), calls)
	}
	if report.UniqueNetworkProbes != len(DiscoveryQueries()) {
		t.Fatalf("expected %d unique probes, got %d", len(DiscoveryQueries()), report.UniqueNetworkProbes)
	}
	for _, result := range report.Results {
		if len(result.Probes) != 2 {
			t.Fatalf("expected GitHub plus Google status for %s, got %d probes", result.AgentID, len(result.Probes))
		}
		if !result.Probes[0].Found {
			t.Fatalf("expected GitHub probe to find canonical repository for %s", result.AgentID)
		}
		if result.Probes[1].Status != "not-configured" {
			t.Fatalf("expected Google to be not-configured, got %q", result.Probes[1].Status)
		}
	}
}

func TestRunIndexAuditUsesOfficialGoogleProviderWhenConfigured(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		if req.URL.Host == "websearchservice.googleapis.com" && req.Header.Get("X-Goog-Api-Key") != "key" {
			t.Fatalf("missing Google API key header")
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(CanonicalRepositoryURL)),
			Header:     make(http.Header),
		}, nil
	})}

	report, err := RunIndexAudit(context.Background(), "abc123", IndexProviderConfig{
		HTTPClient:     client,
		GoogleAPIKey:   "key",
		GoogleClientID: "client",
		GoogleUserIP:   "203.0.113.10",
	})
	if err != nil {
		t.Fatalf("RunIndexAudit returned error: %v", err)
	}
	wantCalls := len(DiscoveryQueries()) * 2
	if calls != wantCalls {
		t.Fatalf("expected %d network calls, got %d", wantCalls, calls)
	}
	if !report.GoogleConfigured {
		t.Fatal("expected Google provider to be configured")
	}
	if report.UniqueNetworkProbes != wantCalls {
		t.Fatalf("expected %d unique probes, got %d", wantCalls, report.UniqueNetworkProbes)
	}
}
