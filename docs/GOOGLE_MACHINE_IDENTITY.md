# Google Machine Identity for COAS

COAS can operate against Google services without storing a personal Google password or a long-lived service account key.

The preferred architecture is:

```text
GITHUB ACTIONS OIDC
        |
        v
GOOGLE WORKLOAD IDENTITY FEDERATION
        |
        v
SHORT-LIVED SERVICE ACCOUNT ACCESS TOKEN
        |
        v
SEARCH CONSOLE URL INSPECTION
        |
        v
ONE VERIFIED INDEX SIGNAL PER CYCLE
        |
        v
100 COAS LOGICAL AGENTS SHARE THE EVIDENCE
```

Core boundary:

```text
MACHINE_IDENTITY != PERSONAL_ACCOUNT_PASSWORD
PUBLIC_SEARCH_RESULT != VERIFIED_INDEX_STATUS
SEARCH_CONSOLE_VERDICT != PERMANENT_INDEXING
HUMAN_AGENCY > MACHINE_AUTHORITY
```

## Two Google channels

COAS intentionally keeps two Google integrations separate.

### 1. Google Web Search Service

Purpose: query public Google Search results for COAS discovery evidence.

Required repository secrets:

```text
GOOGLE_WEB_SEARCH_API_KEY
GOOGLE_WEB_SEARCH_CLIENT_ID
GOOGLE_WEB_SEARCH_USER_IP
```

Google's Web Search Service requires an API key, designated client ID, and end-user IP. Workload Identity Federation does not replace those requirements.

### 2. Google Search Console URL Inspection

Purpose: obtain index-status telemetry for a URL in a Search Console property that the operator controls.

Authentication: GitHub Actions OIDC to Google Workload Identity Federation to a short-lived OAuth access token.

Required repository variables:

```text
GCP_WORKLOAD_IDENTITY_PROVIDER
GCP_SERVICE_ACCOUNT
GOOGLE_SEARCH_CONSOLE_SITE_URL
GOOGLE_SEARCH_CONSOLE_INSPECTION_URL
```

The workflow requests only this OAuth scope:

```text
https://www.googleapis.com/auth/webmasters.readonly
```

No service account JSON key is stored in the repository.

## Ownership requirement

`GOOGLE_SEARCH_CONSOLE_INSPECTION_URL` must be under the Search Console property identified by `GOOGLE_SEARCH_CONSOLE_SITE_URL`.

A normal GitHub repository URL is hosted under `github.com`. Do not assume that owning a repository gives ownership of the `github.com` Search Console property.

For full Search Console telemetry, use a COAS website or other URL-prefix/domain property that you can verify in Search Console. A future COAS GitHub Pages or custom-domain site can serve this purpose once it is enabled and verified.

## One-time Google Cloud setup

Use a Google Cloud project you control. Enable the Search Console API, create a service account for read-only telemetry, and configure a Workload Identity Pool and GitHub OIDC provider.

Recommended trust restriction:

```text
repository == jupiter8nohate/cognitive-offloading-adaptive-symbiosis
```

Google's GitHub provider issuer is:

```text
https://token.actions.githubusercontent.com/
```

Grant the GitHub repository identity permission to impersonate only the dedicated telemetry service account. Do not grant broad project roles that are unrelated to Search Console telemetry.

Then add that service account as an authorized user for the Search Console property that COAS will inspect.

## GitHub variables

After Google Cloud and Search Console are ready, configure these non-secret repository variables:

```text
GCP_WORKLOAD_IDENTITY_PROVIDER=projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL/providers/PROVIDER
GCP_SERVICE_ACCOUNT=coas-search-telemetry@PROJECT_ID.iam.gserviceaccount.com
GOOGLE_SEARCH_CONSOLE_SITE_URL=https://YOUR-VERIFIED-SITE.example/
GOOGLE_SEARCH_CONSOLE_INSPECTION_URL=https://YOUR-VERIFIED-SITE.example/
```

The hourly swarm will then mint a short-lived token automatically and perform one URL Inspection request per cycle. All 100 logical discovery agents share that response.

## What the agents can learn

When Search Console is configured, COAS records evidence such as:

- high-level index verdict
- coverage state
- indexing state
- page fetch state
- last Google crawl time when available
- Google-selected canonical URL when available

A `PASS` verdict is recorded as indexed evidence for that cycle. It is not treated as a permanent guarantee.

## What the agents cannot do

This integration does not:

- expose a personal Google password to GitHub
- grant access to Gmail, Drive, Photos, or unrelated Google data
- force Google to crawl or index a page
- make an index entry permanent
- use the general Indexing API for ordinary COAS pages
- scrape Google result HTML
- multiply one Search Console check into 100 duplicate requests

Google's general Indexing API is intentionally not used because Google restricts it to eligible JobPosting pages and livestream BroadcastEvent pages embedded in VideoObject markup.

## Recovery behavior

If Workload Identity Federation or Search Console variables are missing, the Google Search Console provider reports `not-configured` and the remaining search audit continues.

If Google returns an API error, the report records the error rather than inventing an index state.

```text
AUTH_FAILURE != INDEX_FAILURE
NO_RESULT != NOT_INDEXED
SEARCH_RESULT != GUARANTEED_INDEX_STATE
```
