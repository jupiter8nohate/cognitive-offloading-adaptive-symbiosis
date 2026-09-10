# Authorized Outbound Distribution

COAS can publish verified CMB/COAS signal bundles to endpoints that the repository owner explicitly authorizes.

This subsystem is an outbound publisher, not a self-propagating network agent.

## Supported adapters

- `generic`: HTTPS JSON webhook
- `discord`: Discord webhook owned or administered by the operator
- `slack`: Slack incoming webhook owned or administered by the operator
- `mastodon`: Mastodon account endpoint with an operator-supplied bearer token

No public target is discovered automatically. No third-party repository is modified automatically.

## Enable publishing

Set the repository variable:

```text
COAS_DISTRIBUTION_ENABLED=true
```

Create the repository secret:

```text
COAS_DISTRIBUTION_TARGETS_JSON
```

The value is a JSON array. Example:

```json
[
  {
    "name": "my-discord",
    "kind": "discord",
    "url": "https://discord.com/api/webhooks/REPLACE_WITH_YOUR_WEBHOOK",
    "authorized": true
  },
  {
    "name": "my-slack",
    "kind": "slack",
    "url": "https://hooks.slack.com/services/REPLACE_WITH_YOUR_WEBHOOK",
    "authorized": true
  },
  {
    "name": "my-mastodon",
    "kind": "mastodon",
    "url": "https://YOUR_INSTANCE.example/api/v1/statuses",
    "bearer_token": "REPLACE_WITH_YOUR_TOKEN",
    "authorized": true
  }
]
```

Keep the JSON in a GitHub secret because webhook URLs and bearer tokens can contain credentials.

## Runtime rules

The publisher:

1. accepts at most 10 targets per run,
2. requires every target to declare `authorized: true`,
3. requires HTTPS,
4. runs the full Go verification suite first,
5. emits an idempotency key derived from the source commit and manifesto,
6. does not log target URLs or bearer tokens,
7. does not crawl the web looking for places to post,
8. does not create unsolicited pull requests in unrelated repositories.

## Distribution signal

The default signal is:

```text
𒄆𓁹✞𒀱✞𓁹𒄆
PATTERN != PROOF
PROFILE != PERSON
MODEL != MIND
PREDICTION != DESTINY
HUMAN_AGENCY > MACHINE_AUTHORITY
OFFLOAD != SURRENDER
```

Every payload carries a source commit and SHA-256 content hash so recipients can trace which repository state produced it.

## Role architecture

The 10 operational roles are now:

1. Code Weaver
2. Test Mechanic
3. Profile Auditor
4. Glitchology Compiler
5. Memory Buffer Indexer
6. Contextual Researcher
7. Sovereignty Shield Vector
8. Cipher Deployer
9. Manifesto Broadcaster
10. Runtime Bridge

Across 10 COAS mechanics, that remains exactly 100 logical agents.

The Contextual Researcher is read-only with respect to external public sources. The outbound roles publish only through configured authorized endpoints.

## Not implemented

The repository intentionally does not:

- inject content into arbitrary APIs,
- poison tracking data or model-training corpora,
- mass-post into channels the operator does not control,
- open unsolicited PRs against unrelated public projects,
- evade platform anti-abuse controls.

Those behaviors would undermine the project's human-agency thesis by taking agency away from other people.
