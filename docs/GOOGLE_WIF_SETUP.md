# Google Workload Identity Federation for COAS Search Auditing

This document configures the COAS autonomous swarm to authenticate from GitHub Actions to Google Cloud without storing a long-lived Google service account key in GitHub.

The resulting path is:

```text
GITHUB ACTIONS
      |
      v
GITHUB OIDC TOKEN
      |
      v
GOOGLE WORKLOAD IDENTITY FEDERATION
      |
      v
RESTRICTED COAS SERVICE ACCOUNT
      |
      v
GOOGLE SECRET MANAGER
      |
      v
WEB SEARCH SERVICE CREDENTIALS
      |
      v
100-AGENT INDEX AUDIT
```

Security invariants:

```text
HUMAN_AGENCY > MACHINE_AUTHORITY
SHORT_LIVED_IDENTITY > STATIC_SERVICE_ACCOUNT_KEY
SEARCH != SPAM
OBSERVATION != INDEX_CONTROL
SEARCH_RESULT != GUARANTEED_INDEX_STATE
```

## Important Google requirement

Workload Identity Federation secures the machine identity. It does not grant access to Google Web Search Service by itself.

Google Web Search Service requires:

- an active Google Cloud project;
- a valid Web Search Service API key;
- a designated partner client ID;
- the required end-user IP context for each request.

The API key and partner client ID come from Google's Web Search Service program. Connecting a normal personal Google account does not create those values automatically.

## Repository variables

The workflow reads these non-secret GitHub repository variables:

```text
GOOGLE_CLOUD_PROJECT
GOOGLE_WIF_PROVIDER
GOOGLE_WIF_SERVICE_ACCOUNT
GOOGLE_WSS_API_KEY_SECRET
GOOGLE_WSS_CLIENT_ID_SECRET
GOOGLE_WSS_USER_IP_SECRET
```

Recommended Secret Manager secret IDs:

```text
coas-wss-api-key
coas-wss-client-id
coas-wss-user-ip
```

The current workflow also retains the older GitHub Secrets path as a fallback. The preferred production path is Workload Identity Federation plus Google Secret Manager.

## One-time Google Cloud setup

Set local shell values for the Google Cloud project and this repository:

```bash
export PROJECT_ID="YOUR_GOOGLE_CLOUD_PROJECT_ID"
export PROJECT_NUMBER="$(gcloud projects describe "$PROJECT_ID" --format='value(projectNumber)')"
export POOL_ID="github-coas"
export PROVIDER_ID="coas-main"
export SERVICE_ACCOUNT_ID="coas-search-audit"
export SERVICE_ACCOUNT_EMAIL="${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com"
export REPO="jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
```

Enable the Google Cloud services used by the identity and secret path:

```bash
gcloud services enable \
  iamcredentials.googleapis.com \
  sts.googleapis.com \
  secretmanager.googleapis.com \
  --project="$PROJECT_ID"
```

Create the restricted service account:

```bash
gcloud iam service-accounts create "$SERVICE_ACCOUNT_ID" \
  --project="$PROJECT_ID" \
  --display-name="COAS autonomous search audit"
```

Create the Workload Identity Pool:

```bash
gcloud iam workload-identity-pools create "$POOL_ID" \
  --project="$PROJECT_ID" \
  --location="global" \
  --display-name="COAS GitHub Actions"
```

Create a GitHub OIDC provider restricted to this repository and the main branch:

```bash
gcloud iam workload-identity-pools providers create-oidc "$PROVIDER_ID" \
  --project="$PROJECT_ID" \
  --location="global" \
  --workload-identity-pool="$POOL_ID" \
  --issuer-uri="https://token.actions.githubusercontent.com/" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository,attribute.repository_owner=assertion.repository_owner,attribute.ref=assertion.ref" \
  --attribute-condition="assertion.repository=='${REPO}' && assertion.ref=='refs/heads/main'"
```

Resolve the pool and provider resource names:

```bash
export POOL_RESOURCE="$(gcloud iam workload-identity-pools describe "$POOL_ID" \
  --project="$PROJECT_ID" \
  --location="global" \
  --format='value(name)')"

export PROVIDER_RESOURCE="$(gcloud iam workload-identity-pools providers describe "$PROVIDER_ID" \
  --project="$PROJECT_ID" \
  --location="global" \
  --workload-identity-pool="$POOL_ID" \
  --format='value(name)')"
```

Allow only identities from this repository mapping to impersonate the COAS service account:

```bash
gcloud iam service-accounts add-iam-policy-binding "$SERVICE_ACCOUNT_EMAIL" \
  --project="$PROJECT_ID" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${POOL_RESOURCE}/attribute.repository/${REPO}"
```

## Secret Manager

Create three Secret Manager secrets. Do not commit their values to this repository and do not put a Google password, 2FA code, recovery code, or service account JSON key in the repository.

```bash
gcloud secrets create coas-wss-api-key --project="$PROJECT_ID" --replication-policy="automatic"
gcloud secrets create coas-wss-client-id --project="$PROJECT_ID" --replication-policy="automatic"
gcloud secrets create coas-wss-user-ip --project="$PROJECT_ID" --replication-policy="automatic"
```

Add the actual Web Search Service values as secret versions. The values must come from your authorized Google Web Search Service configuration.

```bash
printf '%s' "$GOOGLE_WEB_SEARCH_API_KEY" | gcloud secrets versions add coas-wss-api-key --project="$PROJECT_ID" --data-file=-
printf '%s' "$GOOGLE_WEB_SEARCH_CLIENT_ID" | gcloud secrets versions add coas-wss-client-id --project="$PROJECT_ID" --data-file=-
printf '%s' "$GOOGLE_WEB_SEARCH_USER_IP" | gcloud secrets versions add coas-wss-user-ip --project="$PROJECT_ID" --data-file=-
```

Grant the COAS service account permission to read only these three secrets:

```bash
for SECRET_ID in coas-wss-api-key coas-wss-client-id coas-wss-user-ip; do
  gcloud secrets add-iam-policy-binding "$SECRET_ID" \
    --project="$PROJECT_ID" \
    --role="roles/secretmanager.secretAccessor" \
    --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}"
done
```

Do not grant broad project Editor or Owner roles to the service account.

## GitHub repository variables

Set the six non-secret variables in GitHub repository settings, or use GitHub CLI after authenticating locally:

```bash
gh variable set GOOGLE_CLOUD_PROJECT --repo "$REPO" --body "$PROJECT_ID"
gh variable set GOOGLE_WIF_PROVIDER --repo "$REPO" --body "$PROVIDER_RESOURCE"
gh variable set GOOGLE_WIF_SERVICE_ACCOUNT --repo "$REPO" --body "$SERVICE_ACCOUNT_EMAIL"
gh variable set GOOGLE_WSS_API_KEY_SECRET --repo "$REPO" --body "coas-wss-api-key"
gh variable set GOOGLE_WSS_CLIENT_ID_SECRET --repo "$REPO" --body "coas-wss-client-id"
gh variable set GOOGLE_WSS_USER_IP_SECRET --repo "$REPO" --body "coas-wss-user-ip"
```

These values identify resources. The API key and partner values remain in Google Secret Manager.

## Autonomous behavior after setup

The hourly `COAS autonomous swarm` workflow receives a GitHub OIDC token. Google validates that the token belongs to this repository and main branch. The workflow then obtains a short-lived Google identity, reads the three authorized Secret Manager values, and executes the 100-agent search index audit.

If WIF is not configured or the Google credential path is unavailable, the Google authentication steps do not prevent the core swarm and GitHub search audit from running. Google search remains unavailable until valid Web Search Service credentials are present.

No workflow receives your personal Google password.

## What this does not do

This configuration does not:

- force Google to crawl or index a page;
- access Google's private indexing logs;
- guarantee search ranking;
- create a permanent or unerasable Google record;
- authorize posting into arbitrary third-party services;
- make the present COAS software agents AGI or ASI.

It provides keyless machine authentication to Google Cloud resources and controlled retrieval of the credentials required for an authorized public search audit.
