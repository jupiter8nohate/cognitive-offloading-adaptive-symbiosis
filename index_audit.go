package coas

import (
	"bytes"
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
	GitHubToken                      string
	GoogleAPIKey                     string
	GoogleClientID                   string
	GoogleUserIP                     string
	GoogleSearchConsoleAccessToken   string
	GoogleSearchConsoleSiteURL       string
	GoogleSearchConsoleInspectionURL string
	HTTPClient                       *http.Client
	RequestTimeout                   time.Duration
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
	Version                 string                 `json:"version"`
	SourceCommit            string                 `json:"source_commit"`
	RepositoryURL           string                 `json:"repository_url"`
	LogicalAgentCount       int                    `json:"logical_agent_count"`
	UniqueQueryCount        int                    `json:"unique_query_count"`
	UniqueNetworkProbes     int                    `json:"unique_network_probes"`
	GoogleConfigured        bool                   `json:"google_configured"`
	SearchConsoleConfigured bool                   `json:"search_console_configured"`
	Results                 []DiscoveryAgentResult `json:"results"`
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
	searchConsoleConfigured := cfg.GoogleSearchConsoleAccessToken != "" &&
		cfg.GoogleSearchConsoleSiteURL != "" &&
		cfg.GoogleSearchConsoleInspectionURL != ""

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

	searchConsoleProbe := IndexProbe{
		Provider: "google-search-console-url-inspection",
		Query:    cfg.GoogleSearchConsoleInspectionURL,
		Status:   "not-configured",
		Found:    false,
		Detail:   "Configure Google Workload Identity Federation and a verified Search Console property to obtain verified Google index telemetry.",
	}
	if searchConsoleConfigured {
		var err error
		searchConsoleProbe, err = probeGoogleSearchConsole(ctx, client, cfg)
		uniqueNetworkProbes++
		if err != nil {
			searchConsoleProbe.Status = "error"
			searchConsoleProbe.Detail = err.Error()
		}
	}

	results := make([]DiscoveryAgentResult, 0, len(agents))
	for i, agent := range agents {
		query := queries[i%len(queries)]
		probes := append([]IndexProbe(nil), cache[query]...)
		probes = append(probes, searchConsoleProbe)
		results = append(results, DiscoveryAgentResult{
			AgentID:  agent.ID,
			Mechanic: agent.Mechanic.Name,
			Role:     agent.Role.Name,
			Query:    query,
			Probes:   probes,
		})
	}

	return IndexAuditReport{
		Version:                 "coas-index-audit.v2",
		SourceCommit:            sourceCommit,
		RepositoryURL:           CanonicalRepositoryURL,
		LogicalAgentCount:       len(agents),
		UniqueQueryCount:        len(queries),
		UniqueNetworkProbes:     uniqueNetworkProbes,
		GoogleConfigured:        googleConfigured,
		SearchConsoleConfigured: searchConsoleConfigured,
		Results:                 results,
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
	req.Header.Set("User-Agent", "coas-index-audit/2.0")
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
	req.Header.Set("User-Agent", "coas-index-audit/2.0")

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

func probeGoogleSearchConsole(ctx context.Context, client *http.Client, cfg IndexProviderConfig) (IndexProbe, error) {
	probe := IndexProbe{
		Provider: "google-search-console-url-inspection",
		Query:    cfg.GoogleSearchConsoleInspectionURL,
		Status:   "checked",
	}
	payload, err := json.Marshal(map[string]string{
		"inspectionUrl": cfg.GoogleSearchConsoleInspectionURL,
		"siteUrl":       cfg.GoogleSearchConsoleSiteURL,
		"languageCode":  "en-US",
	})
	if err != nil {
		return probe, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://searchconsole.googleapis.com/v1/urlInspection/index:inspect",
		bytes.NewReader(payload),
	)
	if err != nil {
		return probe, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.GoogleSearchConsoleAccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "coas-index-audit/2.0")

	body, status, err := executeProbe(client, req)
	if err != nil {
		return probe, err
	}
	if status < 200 || status >= 300 {
		return probe, fmt.Errorf("Google Search Console URL Inspection returned HTTP %d", status)
	}

	var response struct {
		InspectionResult struct {
			IndexStatusResult struct {
				Verdict         string `json:"verdict"`
				CoverageState   string `json:"coverageState"`
				IndexingState   string `json:"indexingState"`
				LastCrawlTime   string `json:"lastCrawlTime"`
				PageFetchState  string `json:"pageFetchState"`
				GoogleCanonical string `json:"googleCanonical"`
			} `json:"indexStatusResult"`
		} `json:"inspectionResult"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return probe, fmt.Errorf("decode Google Search Console response: %w", err)
	}

	statusResult := response.InspectionResult.IndexStatusResult
	probe.Found = statusResult.Verdict == "PASS"
	probe.Detail = searchConsoleDetail(
		statusResult.Verdict,
		statusResult.CoverageState,
		statusResult.IndexingState,
		statusResult.PageFetchState,
		statusResult.LastCrawlTime,
		statusResult.GoogleCanonical,
	)
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

func searchConsoleDetail(verdict, coverage, indexing, fetch, lastCrawl, canonical string) string {
	parts := []string{
		"verdict=" + emptyAsUnknown(verdict),
		"coverage=" + emptyAsUnknown(coverage),
		"indexing=" + emptyAsUnknown(indexing),
		"fetch=" + emptyAsUnknown(fetch),
	}
	if lastCrawl != "" {
		parts = append(parts, "last_crawl="+lastCrawl)
	}
	if canonical != "" {
		parts = append(parts, "google_canonical="+canonical)
	}
	return strings.Join(parts, "; ")
}

func emptyAsUnknown(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
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
	fmt.Fprintf(&b, "Google Web Search Service configured: %t\n\n", r.GoogleConfigured)
	fmt.Fprintf(&b, "Google Search Console configured: %t\n\n", r.SearchConsoleConfigured)
	b.WriteString("The 100 logical agents share deduplicated provider responses. This preserves autonomous analysis while preventing duplicate search traffic. Search Console URL Inspection is also deduplicated to one verified-property request per audit cycle.\n\n")
	b.WriteString("| Provider | Agent observations | Surfaced or indexed |\n")
	b.WriteString("|---|---:|---:|\n")
	for _, provider := range providers {
		fmt.Fprintf(&b, "| %s | %d | %d |\n", provider, providerCounts[provider], providerFound[provider])
	}
	b.WriteString("\n## Query families\n\n")
	for _, query := range DiscoveryQueries() {
		fmt.Fprintf(&b, "- %s\n", query)
	}
	b.WriteString("\nPATTERN != PROOF\n\nPUBLIC_SEARCH_RESULT != VERIFIED_INDEX_STATUS\n\nSEARCH_CONSOLE_VERDICT != PERMANENT_INDEXING\n\nHUMAN_AGENCY > MACHINE_AUTHORITY\n")
	return b.String()
}
