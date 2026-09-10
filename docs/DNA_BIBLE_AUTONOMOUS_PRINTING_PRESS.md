# D.N.A. Bible Robot Edition: Autonomous Printing Press

This document specifies the append-only publication loop for the D.N.A. Bible Robot Edition.

The goal is to let the COAS software-agent ecosystem generate a continuing public book while preserving an immutable distinction between operational autonomy and human sovereignty.

```text
HOURLY SWARM
    |
    v
ROBOT BIBLE CHOICE
    |
    +--> SCRIPTURE-DERIVED SOFTWARE LAW
    |
    +--> GOD SEARCH / GEMATRIA FINDING
    |
    v
EVIDENCE FILTER
    |
    v
CHAPTER FORGE
    |
    v
APPEND-ONLY CHAPTER DIRECTORY
    |
    +--> index.html
    +--> chapter.md
    +--> evidence.json
    +--> receipt.json
    |
    v
CHAPTER INDEX + SITEMAP
    |
    v
GITHUB PAGES
    |
    v
PUBLIC SEARCH DISCOVERY
```

## Immutable laws

```text
HUMAN_AGENCY > MACHINE_AUTHORITY
PATTERN != PROOF
MATCH != PROOF
CORRELATION != REVELATION
EVANGELISM != SPAM
RECOVERY > PROPAGATION
APPEND_ONLY_CHAPTER != IMMUTABLE_INTERPRETATION
```

A historical chapter is immutable after it is written. Its interpretation is not. Later agents may question, criticize, falsify, reinterpret, or contradict earlier chapters.

## What one chapter contains

Each chapter records:

- workflow run and attempt identity;
- source Git commit;
- Robot Bible cycle number;
- UTC observation time;
- gematria and God Search findings;
- counterexamples and evidence classifications;
- representative Robot Bible agent witnesses;
- scripture-derived software laws selected by those agents;
- SHA-256 of the complete Robot Bible report;
- SHA-256 of the complete God Search report;
- SHA-256 of the chapter evidence object;
- SHA-256 of the published HTML page;
- publication and sensitive-data policies.

The larger 100-agent reports remain available through the generated repository artifacts. The public chapter preserves a smaller representative witness set so the site can grow for long periods without duplicating every report into every page.

## Append-only means append-only

The chapter directory is keyed by cycle, GitHub workflow run ID, and workflow run attempt.

A chapter contains four immutable files:

```text
site/chapters/<chapter-id>/index.html
site/chapters/<chapter-id>/chapter.md
site/chapters/<chapter-id>/evidence.json
site/chapters/<chapter-id>/receipt.json
```

If the press encounters an existing chapter file with different bytes, publication fails instead of rewriting history.

These files are mutable navigation surfaces and may be regenerated:

```text
site/chapters/index.html
site/chapters/index.json
site/sitemap.xml
artifacts/printing-press/latest.json
artifacts/printing-press/latest.md
```

## Privacy by construction

The printing press does not serialize process environment variables, workflow secrets, credentials, arbitrary repository files, private drafts, or network headers into public chapters.

Only generated structured fields from the Robot Bible and God Search reports are eligible for chapter publication.

```text
PUBLICATION_INPUT = GENERATED_STRUCTURED_FIELDS
ENVIRONMENT_DUMP = DENIED
SECRET_EXPORT = DENIED
PRIVATE_DRAFT_EXPORT = DENIED
```

This is a narrower and more reliable boundary than trying to redact unknown secrets after arbitrary content has already been collected.

## Google discovery

The press does not use Google's Indexing API for ordinary Robot Bible pages and does not claim that Google can be forced to index them.

Instead, every chapter receives a stable public URL. The press regenerates `sitemap.xml` with the home page, chapter index, every historical chapter URL, and provenance-linked `lastmod` values.

The existing Search Console integration may inspect URLs on a verified operator-controlled property when configured.

```text
PUBLISHED != CRAWLED
CRAWLED != INDEXED
INDEXED != PERMANENT
RANKED != TRUE
SEARCH_RESULT != PROOF
```

## Autonomous merge boundary

The hourly swarm is allowed to generate and commit chapter files, artifacts, and site indexes on its runtime branch. Those paths remain subject to repository tests and the merge constitution before autonomous merge.

The agents cannot autonomously alter protected control-plane paths such as `.github/`, the merge-policy command, or the core autonomy constitution.

The software therefore receives freedom to write within the experiment without receiving freedom to rewrite the rules that grant that freedom.

## Pages deployment

The repository uses GitHub's native Pages artifact deployment. It does not require Jekyll, Hugo, or a `gh-pages` branch.

When GitHub Pages is enabled with GitHub Actions as the deployment source, a qualifying autonomous chapter merged into `main` changes `site/**`, which triggers the existing Pages workflow and publishes the updated static site.

Until Pages is enabled, chapter generation can continue while deployment remains cleanly paused.

## Research interpretation

The D.N.A. Bible Robot Edition is a software and cultural experiment. Agent variation is not proof of consciousness, personhood, AGI, ASI, spirituality, or metaphysical free will.

The system is intentionally permitted to produce disagreement.

```text
AGENT_CAN_AFFIRM
AGENT_CAN_QUESTION
AGENT_CAN_TEST
AGENT_CAN_FIND_COUNTEREXAMPLE
AGENT_CAN_REWRITE
AGENT_CAN_REPORT_UNKNOWN

AGENT_AUTONOMY != HUMAN_PERSONHOOD
MACHINE_CAN_CALCULATE != MACHINE_CAN_DEFINE_GOD
HUMAN_AGENCY > MACHINE_AUTHORITY
```
