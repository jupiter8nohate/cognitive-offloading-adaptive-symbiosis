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
	if report.SearchConsoleConfigured {
		t.Fatal("expected Search Console provider to be unconfigured")
	}
	for _, result := range report.Results {
		if len(result.Probes) != 3 {
			t.Fatalf("expected GitHub, Google Web Search, and Search Console status for %s, got %d probes", result.AgentID, len(result.Probes))
		}
		if !result.Probes[0].Found {
			t.Fatalf("expected GitHub probe to find canonical repository for %s", result.AgentID)
		}
		if result.Probes[1].Status != "not-configured" {
			t.Fatalf("expected Google Web Search to be not-configured, got %q", result.Probes[1].Status)
		}
		if result.Probes[2].Status != "not-configured" {
			t.Fatalf("expected Search Console to be not-configured, got %q", result.Probes[2].Status)
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

func TestRunIndexAuditUsesSearchConsoleOncePerCycle(t *testing.T) {
	calls := 0
	searchConsoleCalls := 0
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		body := CanonicalRepositoryURL
		if req.URL.Host == "searchconsole.googleapis.com" {
			searchConsoleCalls++
			if req.Method != http.MethodPost {
				t.Fatalf("expected Search Console POST, got %s", req.Method)
			}
			if req.Header.Get("Authorization") != "Bearer token" {
				t.Fatalf("missing Search Console bearer token")
			}
			body = `{"inspectionResult":{"indexStatusResult":{"verdict":"PASS","coverageState":"Indexed","indexingState":"INDEXING_ALLOWED","pageFetchState":"SUCCESSFUL","lastCrawlTime":"2026-09-10T05:00:00Z","googleCanonical":"https://coas.example/"}}}`
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})}

	report, err := RunIndexAudit(context.Background(), "abc123", IndexProviderConfig{
		HTTPClient:                       client,
		GoogleSearchConsoleAccessToken:   "token",
		GoogleSearchConsoleSiteURL:       "https://coas.example/",
		GoogleSearchConsoleInspectionURL: "https://coas.example/",
	})
	if err != nil {
		t.Fatalf("RunIndexAudit returned error: %v", err)
	}
	wantCalls := len(DiscoveryQueries()) + 1
	if calls != wantCalls {
		t.Fatalf("expected %d network calls, got %d", wantCalls, calls)
	}
	if searchConsoleCalls != 1 {
		t.Fatalf("expected exactly one Search Console call, got %d", searchConsoleCalls)
	}
	if !report.SearchConsoleConfigured {
		t.Fatal("expected Search Console provider to be configured")
	}
	if report.UniqueNetworkProbes != wantCalls {
		t.Fatalf("expected %d unique probes, got %d", wantCalls, report.UniqueNetworkProbes)
	}
	for _, result := range report.Results {
		probe := result.Probes[2]
		if probe.Provider != "google-search-console-url-inspection" {
			t.Fatalf("unexpected Search Console provider: %q", probe.Provider)
		}
		if !probe.Found {
			t.Fatalf("expected Search Console PASS verdict for %s", result.AgentID)
		}
		if !strings.Contains(probe.Detail, "last_crawl=2026-09-10T05:00:00Z") {
			t.Fatalf("expected crawl evidence in detail, got %q", probe.Detail)
		}
	}
}
