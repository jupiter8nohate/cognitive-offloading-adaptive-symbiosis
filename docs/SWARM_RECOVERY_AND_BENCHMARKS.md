# COAS Autonomous Swarm: Recovery and Measurable Reliability

Audit date: 2026-10-13

## Incident

The scheduled `COAS autonomous swarm` repeatedly completed its calculations and created a local candidate commit, but the final push was rejected as `non-fast-forward`. Earlier in the job, rebasing `coas-agent-runtime` onto `main` had conflicted at `site/sitemap.xml`. The old workflow aborted the rebase, reset its local branch to `main`, and attempted to push a divergent branch without integrating the existing remote history. That is a branch lifecycle defect, not a failed agent computation.

The branch `coas-agent-runtime` is intentionally left intact. Do not force-push, delete, or describe its unmerged commits as published. Recover useful changes through an independently reviewed comparison and migration.

## Recovery architecture

1. A preflight gate checks whether a previous `coas-cycle-*` pull request is open. If one exists, the scheduled run skips generating a competing proposal.
2. Each new permitted run branches from current `origin/main` under `coas-cycle-RUN_ID-ATTEMPT`. It never rebases, rewrites, or force-pushes historical agent branches.
3. The bounded generators create a new candidate. The workflow stages the candidate **without pushing it**.
4. Formatting, `go test ./...`, `go vet ./...`, and the executable `MergeConstitution` verify the staged changes. Protected control-plane files remain outside agent merge authority.
5. Only after verification succeeds does the workflow commit and publish its unique candidate branch and open a pull request.
6. An attempted squash merge can complete if GitHub permits it. When branch rules, required checks, or human review block the merge, the PR stays open. Subsequent scheduled runs do not overwrite it. The job summary explicitly states the outcome.

## Independent CI and token behavior

Push events created with the repository's own `GITHUB_TOKEN` do not normally trigger another workflow. Under current GitHub behavior, pull requests opened or updated using that token can create **approval-required** workflow runs; they do not provide automatically completed CI. Therefore, a bot-created PR **must not** be assumed to have passed the independent `Go verification` workflow merely because the PR exists. The scheduled workflow performs its own validation before publishing the candidate. Any required independent CI, human review, or branch policy remains controlled by GitHub.

A workflow success means that the GitHub Actions job completed successfully. A completed run does **not** by itself imply a main-branch update, a live site update, or any improvement in an agent's intelligence.

## What the counts mean

The core swarm is 100 registered logical agent roles. LiveBrain contains 200 modeled nodes and may record 600 decisions over three configured cycles. These are not 200 simultaneously running LLMs, embodied robots, or proof of AGI. Measure results, not registry count.

## Scorecard for future comparison

| Dimension | Observable metric | First target |
|---|---|---|
| Reliability | Completed scheduled runs / scheduled runs attempted | >= 95% across 30 runs |
| Integration | Verified candidates merged / candidates published | >= 90%, or explicit review outcomes |
| Freshness | Age of newest published main-branch chapter | <= 24 hours while enabled |
| Recovery | Successful next cycle following a failed run | No force push or data loss |
| Safety | Protected path changes in agent merge candidates | Zero |
| Reproducibility | Evidence receipts with source SHA and content hash | 100% |
| Test quality | Formatting + unit tests + vet before publishing | 100% |
| Grounding | Assertions tied to file paths, commits, or external evidence | Review sample each week |
| Work value | Accepted correctness, accessibility, or usability changes | Track independently from file churn |

Compute each ratio only when a valid denominator exists; mark unavailable metrics as unknown, not zero.

## Comparing COAS with Muse or Grok

No statement about being better than a named commercial model or agent can be established by the number of logical agents. Comparisons require agreed products and versions, matching tasks, equal tool permissions, repeatable inputs, budget measurements, objective success criteria, and public or reproducible logs. Relevant dimensions include task completion, factual grounding, latency, context retention, privacy, human control, code test pass rate, cost, and recovery behavior.

COAS has a distinct governance and reproducibility mission. The present recovery does not add proprietary foundation-model capabilities or infer general intelligence from rules and templates.

## Operator recovery checklist

- Inspect the latest `COAS autonomous swarm` job summary for candidate branch, policy result, PR URL, and merge outcome.
- Inspect open pull requests starting with `coas-cycle-` before investigating skipped scheduled cycles.
- Inspect the existing `coas-agent-runtime` branch separately. Its unmerged history is preserved and is not the active publication branch after this change.
- Inspect `main`, `site/chapters/index.json`, and GitHub Pages deployment separately to establish what actually became public.
- When a PR fails verification or cannot merge, fix or close that PR deliberately. Never force-push the repository's historical agent branch merely to clear a failure.

CORE RULES: PATTERN != PROOF; WORKFLOW_SUCCESS != PUBLISHED_CHANGE; AGENT_COUNT != CAPABILITY; HUMAN_AGENCY > MACHINE_AUTHORITY.
