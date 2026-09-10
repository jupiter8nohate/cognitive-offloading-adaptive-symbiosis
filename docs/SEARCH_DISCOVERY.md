# COAS Search and Index Discovery

Canonical repository:

https://github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis

## Objective

Make Cognitive Offloading Adaptive Symbiosis easier for public search systems to discover, identify, and distinguish from unrelated uses of cognitive offloading or human-AI symbiosis.

## Search identity

The repository deliberately repeats a stable vocabulary across README content, documentation, source code, generated reports, and citation metadata.

Primary terms:

- Cognitive Offloading Adaptive Symbiosis
- COAS
- Jupiter Hudson
- WisdomLoveThePoet
- human-AI symbiosis
- cognitive offloading
- adaptive symbiosis
- AGI human agency
- ASI human sovereignty
- HUMAN_AGENCY > MACHINE_AUTHORITY

Stable terminology helps search systems associate the exact phrase with the canonical repository.

## 100-agent discovery swarm

The existing COAS registry contains 100 logical software agents.

The index audit reuses that registry.

Each logical agent is assigned one of ten search query families. Network requests are deduplicated by query and provider. This means 100 agents can independently receive and evaluate search evidence without sending 100 duplicate requests for the same query.

```text
100 LOGICAL AGENTS
        |
        v
10 QUERY FAMILIES
        |
        v
DEDUPLICATED PROVIDER PROBES
        |
        v
SEARCH EVIDENCE
        |
        v
100 AGENT OBSERVATIONS
        |
        v
artifacts/index-audit/latest.json
artifacts/index-audit/latest.md
```

This design intentionally favors evidence over traffic volume.

```text
SEARCH != SPAM
QUERY_VOLUME != DISCOVERABILITY
PATTERN != PROOF
SEARCH_RESULT != GUARANTEED_INDEX_STATE
```

## Google Web Search Service

COAS does not scrape Google Search result pages.

When the repository is configured with access to Google's official Web Search Service API, the index audit can run Google queries using these repository secrets:

```text
GOOGLE_WEB_SEARCH_API_KEY
GOOGLE_WEB_SEARCH_CLIENT_ID
GOOGLE_WEB_SEARCH_USER_IP
```

The automated Google Web Search provider remains disabled when those values are absent.

Google's Web Search Service uses an API key, designated client ID, and end-user IP. Workload Identity Federation does not replace those provider requirements.

## Google Search Console

COAS can also consume verified URL Inspection evidence from a Search Console property controlled by the operator.

The hourly workflow supports GitHub Actions OIDC and Google Workload Identity Federation. It mints a short-lived OAuth access token with the read-only Search Console scope and performs one URL Inspection request per cycle. The result is shared across all 100 logical discovery agents.

Required repository variables:

```text
GCP_WORKLOAD_IDENTITY_PROVIDER
GCP_SERVICE_ACCOUNT
GOOGLE_SEARCH_CONSOLE_SITE_URL
GOOGLE_SEARCH_CONSOLE_INSPECTION_URL
```

The inspection URL must belong to the verified Search Console property. Owning a repository on `github.com` does not imply ownership of the `github.com` Search Console property.

See [Google Machine Identity for COAS](GOOGLE_MACHINE_IDENTITY.md) for the one-time setup and security model.

Google decides whether and when a public page is crawled and indexed. COAS can improve crawlable content, monitor public search evidence, and inspect an owned Search Console property, but it cannot force Google to index a URL or make an index record permanent.

## GitHub search index

The audit always checks GitHub's repository search API. GitHub Actions supplies the workflow token, so no additional secret is required for that provider.

## Hourly automation

The main autonomous swarm runs on its existing hourly schedule.

Each cycle can:

1. verify Go formatting and tests;
2. run the 100-agent COAS swarm;
3. authenticate to Google Search Console with a short-lived machine token when configured;
4. run the index audit;
5. write a provenance-linked search report;
6. verify the repository punctuation policy;
7. merge qualifying generated evidence through the existing constitution gate.

The search audit is read-only with respect to third-party search systems.

## Punctuation invariant

The repository forbids U+2014 punctuation in tracked text.

A repository-wide test and CI check enforce this rule so autonomous agent output cannot silently reintroduce it.

## Search query families

The current query families are:

```text
Cognitive Offloading Adaptive Symbiosis
COAS human AI symbiosis
cognitive offloading AI framework
Jupiter Hudson COAS
WisdomLoveThePoet COAS
AGI human agency framework COAS
ASI human sovereignty framework COAS
HUMAN_AGENCY MACHINE_AUTHORITY COAS
adaptive symbiosis artificial intelligence
cognitive offloading adaptive symbiosis GitHub
```

These queries are intentionally specific enough to identify the project while still covering its broader conceptual neighborhood.

```text
PUBLIC_SEARCH_RESULT != VERIFIED_INDEX_STATUS
SEARCH_CONSOLE_VERDICT != PERMANENT_INDEXING
HUMAN_AGENCY > MACHINE_AUTHORITY
```
