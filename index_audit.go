package coas

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const CanonicalRepositoryURL = "https://github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"

type IndexProviderConfig struct {
	GitHubToken    string
	GoogleAPIKey   string
	GoogleClientID string
	GoogleUserIP   string
	HTTPClient     *http.Client
	RequestTimeout time.Duration
}

type IndexProbe struct {
	Provider string `json:"provider"`
	Query    string `json:"query"`
	Status   string `json:"status"`
	Found    bool   `json:"found"`
	Detail   string `json:"detail"`
}

type DiscoveryAgentResult struct {
	AgentID  string       `json:"agent_id"`
	Mechanic string       `json:"mechanic"`
	Role     string       `json:"role"`
	Query    string       `json:"query"`
	Probes   []IndexProbe `json:"probes"`
}

type IndexAuditReport struct {
	Version             string                 `json:"version"`
	SourceCommit        string                 `json:"source_commit"`
	RepositoryURL       string                 `json:"repository_url"`
	LogicalAgentCount   int                    `json:"logical_agent_count"`
	UniqueQueryCount    int                    `json:"unique_query_count"`
	UniqueNetworkProbes int                    `json:"unique_network_probes"`
	GoogleConfigured    bool                   `json:"google_configured"`
	Results             []DiscoveryAgentResult `json:"results"`
}

func DiscoveryQueries() []string {
	return []string{
		"Cognitive Offloading Adaptive Symbiosis",
		"COAS human AI symbiosis",
		"cognitive offloading AI framework",
		"Jupiter Hudson COAS",
		"WisdomLoveThePoet COAS",
		"AGI human agency framework COAS",
		"ASI human sovereignty framework COAS",
		"HUMAN_AGENCY MACHINE_AUTHORITY COAS",
		"adaptive symbiosis artificial intelligence",
		"cognitive offloading adaptive symbiosis GitHub",
	}
}

func RunIndexAudit(ctx context.Context, sourceCommit string, cfg IndexProviderConfig) (IndexAuditReport, error) {
	agents := BuildRegistry()
	if err := ValidateRegistry(agents); err != nil {
		return IndexAuditReport{}, err
	}

	queries := DiscoveryQueries()
	client := cfg.HTTPClient
	if client == nil {
		timeout := cfg.RequestTimeout
		if timeout <= 0 {
			timeout = 15 * time.Second
		}
		client = &http.Client{Timeout: timeout}
	}

	googleConfigured := cfg.GoogleAPIKey != "" && cfg.GoogleClientID != "" && cfg.GoogleUserIP != ""
	cache := make(map[string][]IndexProbe, len(queries))
	uniqueNetworkProbes := 0
	for _, query := range queries {
		probes := make([]IndexProbe, 0, 2)
		githubProbe, err := probeGitHubIndex(ctx, client, cfg.GitHubToken, query)
		if err != nil {
			githubProbe.Status = "error"
			githubProbe.Detail = err.Error()
		}
		probes = append(probes, githubProbe)
		uniqueNetworkProbes++

		if googleConfigured {
			googleProbe, err := probeGoogleIndex(ctx, client, cfg, query)
			if err != nil {
				googleProbe.Status = "error"
				googleProbe.Detail = err.Error()
			}
			probes = append(probes, googleProbe)
			uniqueNetworkProbes++
		} else {
			probes = append(probes, IndexProbe{
				Provider: "google-web-search-service",
				Query:    googleQuery(query),
				Status:   "not-configured",
				Found:    false,
				Detail:   "Set GOOGLE_WEB_SEARCH_API_KEY, GOOGLE_WEB_SEARCH_CLIENT_ID, and GOOGLE_WEB_SEARCH_USER_IP to use Google's official Web Search Service API.",
			})
		}
		cache[query] = probes
	}

	results := make([]DiscoveryAgentResult, 0, len(agents))
	for i, agent := range agents {
		query := queries[i%len(queries)]
		probes := append([]IndexProbe(nil), cache[query]...)
		results = append(results, DiscoveryAgentResult{
			AgentID:  agent.ID,
			Mechanic: agent.Mechanic.Name,
			Role:     agent.Role.Name,
			Query:    query,
			Probes:   probes,
		})
	}

	return IndexAuditReport{
		Version:             "coas-index-audit.v1",
		SourceCommit:        sourceCommit,
		RepositoryURL:       CanonicalRepositoryURL,
		LogicalAgentCount:   len(agents),
		UniqueQueryCount:    len(queries),
		UniqueNetworkProbes: uniqueNetworkProbes,
		GoogleConfigured:    googleConfigured,
		Results:             results,
	}, nil
}

func probeGitHubIndex(ctx context.Context, client *http.Client, token, query string) (IndexProbe, error) {
	probe := IndexProbe{Provider: "github-search", Query: query, Status: "checked"}
	u, err := url.Parse("https://api.github.com/search/repositories")
	if err != nil {
		return probe, err
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("per_page", "10")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return probe, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "coas-index-audit/1.0")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	body, status, err := executeProbe(client, req)
	if err != nil {
		return probe, err
	}
	if status < 200 || status >= 300 {
		return probe, fmt.Errorf("GitHub search returned HTTP %d", status)
	}
	probe.Found = containsCanonicalRepository(body)
	probe.Detail = visibilityDetail(probe.Found)
	return probe, nil
}

func probeGoogleIndex(ctx context.Context, client *http.Client, cfg IndexProviderConfig, query string) (IndexProbe, error) {
	searchQuery := googleQuery(query)
	probe := IndexProbe{Provider: "google-web-search-service", Query: searchQuery, Status: "checked"}
	u, err := url.Parse("https://websearchservice.googleapis.com/v1:search")
	if err != nil {
		return probe, err
	}
	q := u.Query()
	q.Set("searchQuery.query", searchQuery)
	q.Set("clientContext.clientId", cfg.GoogleClientID)
	q.Set("userContext.ipAddress", cfg.GoogleUserIP)
	q.Set("pageSize", "10")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return probe, err
	}
	req.Header.Set("X-Goog-Api-Key", cfg.GoogleAPIKey)
	req.Header.Set("User-Agent", "coas-index-audit/1.0")

	body, status, err := executeProbe(client, req)
	if err != nil {
		return probe, err
	}
	if status < 200 || status >= 300 {
		return probe, fmt.Errorf("Google Web Search Service returned HTTP %d", status)
	}
	probe.Found = containsCanonicalRepository(body)
	probe.Detail = visibilityDetail(probe.Found)
	return probe, nil
}

func executeProbe(client *http.Client, req *http.Request) ([]byte, int, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}

func googleQuery(query string) string {
	return fmt.Sprintf("site:github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis %s", query)
}

func containsCanonicalRepository(body []byte) bool {
	text := strings.ToLower(string(body))
	return strings.Contains(text, "jupiter8nohate/cognitive-offloading-adaptive-symbiosis") ||
		strings.Contains(text, strings.ToLower(CanonicalRepositoryURL))
}

func visibilityDetail(found bool) string {
	if found {
		return "canonical repository surfaced in provider response"
	}
	return "canonical repository not found in the first provider result page"
}

func (r IndexAuditReport) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r IndexAuditReport) Markdown() string {
	providerCounts := map[string]int{}
	providerFound := map[string]int{}
	for _, result := range r.Results {
		for _, probe := range result.Probes {
			providerCounts[probe.Provider]++
			if probe.Found {
				providerFound[probe.Provider]++
			}
		}
	}

	providers := make([]string, 0, len(providerCounts))
	for provider := range providerCounts {
		providers = append(providers, provider)
	}
	sort.Strings(providers)

	var b strings.Builder
	b.WriteString("# COAS Search Index Audit\n\n")
	fmt.Fprintf(&b, "Canonical repository: %s\n\n", r.RepositoryURL)
	fmt.Fprintf(&b, "Logical discovery agents: %d\n\n", r.LogicalAgentCount)
	fmt.Fprintf(&b, "Unique query families: %d\n\n", r.UniqueQueryCount)
	fmt.Fprintf(&b, "Unique network probes: %d\n\n", r.UniqueNetworkProbes)
	fmt.Fprintf(&b, "Google official API configured: %t\n\n", r.GoogleConfigured)
	b.WriteString("The 100 logical agents share deduplicated provider responses. This preserves autonomous analysis while preventing 100 duplicate requests for the same query.\n\n")
	b.WriteString("| Provider | Agent observations | Surfaced |\n")
	b.WriteString("|---|---:|---:|\n")
	for _, provider := range providers {
		fmt.Fprintf(&b, "| %s | %d | %d |\n", provider, providerCounts[provider], providerFound[provider])
	}
	b.WriteString("\n## Query families\n\n")
	for _, query := range DiscoveryQueries() {
		fmt.Fprintf(&b, "- %s\n", query)
	}
	b.WriteString("\nPATTERN != PROOF\n\nSEARCH_RESULT != GUARANTEED_INDEX_STATE\n\nHUMAN_AGENCY > MACHINE_AUTHORITY\n")
	return b.String()
}
