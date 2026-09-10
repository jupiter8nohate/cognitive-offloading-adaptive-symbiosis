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

## Google

COAS does not scrape Google Search result pages.

When the repository is configured with access to Google's official Web Search Service API, the index audit can run Google queries using these runtime values:

```text
GOOGLE_WEB_SEARCH_API_KEY
GOOGLE_WEB_SEARCH_CLIENT_ID
GOOGLE_WEB_SEARCH_USER_IP
```

Google Web Search Service requires a valid API key and designated partner client ID. Workload Identity Federation does not replace those requirements.

### Preferred keyless credential path

The preferred production architecture uses GitHub OIDC, Google Workload Identity Federation, a restricted Google service account, and Google Secret Manager.

```text
GITHUB ACTIONS
      |
      v
GITHUB OIDC
      |
      v
GOOGLE WORKLOAD IDENTITY FEDERATION
      |
      v
RESTRICTED COAS SERVICE ACCOUNT
      |
      v
SECRET MANAGER
      |
      v
WEB SEARCH SERVICE VALUES
      |
      v
100-AGENT INDEX AUDIT
```

The hourly workflow requests a short-lived Google identity only when the required repository variables are configured. It then reads the Web Search Service values from Secret Manager. The older GitHub Secrets path remains available as a fallback.

See [Google Workload Identity Federation for COAS Search Auditing](GOOGLE_WIF_SETUP.md) for the one-time setup.

The automated Google provider remains disabled when the required Web Search Service values are absent.

This is intentional. The project does not pretend that a missing credential means the page is absent from Google, and it does not bypass provider controls with HTML scraping.

Google decides whether and when a public GitHub page is crawled and indexed. COAS can improve crawlable content and monitor public search evidence, but it cannot force Google to index a URL.

### Recovery behavior

Google authentication failure is isolated from the core swarm. WIF and Secret Manager steps are allowed to fail closed for the Google provider while the GitHub search audit and the rest of the bounded COAS workflow continue. A report only marks Google as configured when all three required Web Search Service runtime values are present.

## GitHub search index

The audit always checks GitHub's repository search API. GitHub Actions supplies the workflow token, so no additional secret is required for that provider.

## Hourly automation

The main autonomous swarm runs on its existing hourly schedule.

Each cycle can:

1. verify Go formatting and tests;
2. obtain a short-lived Google identity when WIF is configured;
3. retrieve authorized Web Search Service values from Google Secret Manager;
4. run the 100-agent COAS swarm;
5. run the index audit;
6. write a provenance-linked search report;
7. verify the repository punctuation policy;
8. merge qualifying generated evidence through the existing constitution gate.

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
