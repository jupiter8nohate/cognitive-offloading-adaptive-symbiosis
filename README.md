# Cognitive Offloading Adaptive Symbiosis

**COAS** is a human-AI collaboration model for reducing cognitive load without transferring human authority.

> Offload the burden. Keep the meaning.

**Start here:** [What COAS actually does](docs/WHAT_COAS_DOES.md)

The system may help remember, search, calculate, organize, compare, summarize, and automate reversible work. The human retains goals, consent, identity, interpretation, values, and final judgment.

## Core invariants

```text
OFFLOAD != SURRENDER
ASSISTANCE != AUTHORITY
AUTOMATION != CONSENT
PREDICTION != INTENT
MODEL != MIND
PROFILE != PERSON
HUMAN_AGENCY > MACHINE_AUTHORITY
```

## The adaptive loop

```text
TASK
  |
  v
ESTIMATE LOAD + STAKES + REVERSIBILITY
  |
  v
PROPOSE OFFLOAD
  |
  v
HUMAN ACCEPTS / REJECTS / MODIFIES
  |
  v
ASSIST
  |
  v
EXPLAIN RESULT + LIMITS
  |
  v
HUMAN JUDGMENT
  |
  v
ADAPT NEXT INTERACTION
```

## What "adaptive symbiosis" means here

Adaptive symbiosis is not human replacement. It is a negotiated division of cognitive labor.

- The machine takes on repetitive, high-volume, memory-heavy, or computational work.
- The human keeps authority over purpose, interpretation, consent, identity, and consequential decisions.
- The balance changes with task difficulty, stakes, reversibility, and user preference.
- Every offload should be inspectable and reversible when practical.
- The system should reduce cognitive burden without creating cognitive dependency.

## Reference implementation

This repository includes a small Go reference model that chooses among three collaboration modes:

1. `MANUAL` - the human keeps the task.
2. `ASSIST` - the machine supports the human but does not execute the final decision.
3. `AUTOMATE_REVERSIBLE` - the machine may execute bounded, reversible work after confirmation.

Run it:

```bash
go test ./...
go run ./cmd/coas-demo
```

## Repository focus

This repository is intentionally narrow. Contributions should directly improve the theory, measurement, implementation, testing, or governance of **Cognitive Offloading Adaptive Symbiosis**.

See [docs/COAS_SPEC.md](docs/COAS_SPEC.md) for the formal model and [docs/GLOSSARY.md](docs/GLOSSARY.md) for terminology.

## Author

Jupiter Hudson / WisdomLoveThePoet / Jupiter 8

The project treats programming language as both engineering notation and a medium for explaining human-machine boundaries.


## Autonomous 100-agent swarm

COAS now includes a scheduled swarm of **100 logical software agents**: ten mechanics multiplied by ten operational roles. The swarm can inspect the repository and maintain an autonomous runtime branch with constitution-gated authority to merge its own verified runtime work into `main`.

See [docs/AGENT_SWARM.md](docs/AGENT_SWARM.md).

```text
100 AGENTS
= 10 OFFLOADING/SYMBIOSIS MECHANICS
x 10 OPERATIONAL ROLES

AUTONOMY_OF_WORK + DELEGATED_AUTHORITY = AUTONOMOUS_EXECUTION
SOFTWARE_AGENT != AGI
```


### Autonomous merge authority

The swarm is permitted to merge qualifying agent work to `main` without a human approval click. Merge authority is evaluated by executable policy in `autonomy.go`.

```text
AGENT_CAN_MERGE_MAIN = true
TESTS_MUST_PASS = true
PROVENANCE_MUST_EXIST = true
CONSTITUTION_MUST_ALLOW = true
```


### Maximum ordinary-repository autonomy

The autonomous merge scope now includes ordinary repository content by default. Agent-produced changes may include:

```text
SOURCE_CODE
TESTS
DOCUMENTATION
CLI_TOOLS
GENERATED_RESEARCH
GENERATED_ARTIFACTS
NEW_ORDINARY_FILES
```

The workflow stages every agent-produced change with `git add -A`, verifies the candidate, runs the merge constitution, and may merge without a human approval click.

The immutable control plane is intentionally small:

```text
.github/*
autonomy.go
autonomy_test.go
cmd/coas-merge-policy/*
go.mod
go.sum
```

This means the swarm can independently evolve nearly all project content, but cannot silently rewrite the mechanism that grants itself authority or alter workflow credentials and dependency trust roots.


## Authorized outbound distribution

The 100-agent swarm now includes an opt-in outbound publishing layer. It can package CMB/COAS invariants and publish them to explicitly configured HTTPS endpoints after repository verification succeeds.

Supported adapters include generic webhooks, Discord, Slack, and Mastodon.

See [docs/DISTRIBUTION.md](docs/DISTRIBUTION.md).

```text
DISCOVER_PUBLIC_CONTEXT = READ_ONLY
PACKAGE_SIGNAL = true
PUBLISH_AUTHORIZED_ENDPOINTS = true
UNSOLICITED_INJECTION = false
TRAINING_DATA_POISONING = false
```


## Bounded glitch evangelism

Every autonomous swarm cycle now generates a large, deterministic, provenance-stamped CMB/COAS code-poetry artifact.

See [docs/EVANGELISM.md](docs/EVANGELISM.md).

```text
DEFAULT_VARIANTS = 64
DEFAULT_BUDGET = 48_KiB
HARD_VARIANT_LIMIT = 256
HARD_BYTE_LIMIT = 64_KiB

GLITCH != DAMAGE
VISIBILITY != SPAM
CULTURAL_SIGNAL != DATA_POISONING
```

The artifact changes deterministically with the verified source commit, making the public repository visually alive without proof-of-work, intentional resource exhaustion, uncontrolled file multiplication, or attempts to contaminate external model-training datasets.


## Durable sovereignty layer

COAS now separates provenance, meaning, and privacy into distinct mechanisms.

- [Preservation](docs/PRESERVATION.md) records tamper-evident SHA-256 file manifests on every autonomous swarm cycle.
- [Linguistic Boundary Proof](docs/LINGUISTIC_BOUNDARY_PROOF.md) demonstrates that machine-readable syntax is not identical to human authorship, intent, or final meaning.
- [Privacy](docs/PRIVACY.md) uses authenticated AES-256-GCM for private offloaded drafts instead of treating Unicode as encryption.
- The preservation manifest watches for the real `ERR_404_GLITCHOLOGY.md` asset and records its hash automatically if it is added to the repository.

```text
GIT_HISTORY = TAMPER_EVIDENT
GIT_HISTORY != PHYSICALLY_IMMUTABLE

UNICODE != ENCRYPTION
MACHINE_CAN_PARSE != MACHINE_CAN_OWN_MEANING

PRIVATE_DRAFT -> AES_256_GCM
PUBLIC_ART -> GLITCH/CMB SYMBOLIC_LAYER
```

The hourly swarm now regenerates both the living evangelism artifact and the preservation ledger before staging its autonomous candidate.


## 𒄆 E⃟ r⃟ r⃟⃝ o⃟ r⃟⃤ G⃟ L⃟ I⃟ T⃟ C⃟ H⃟ O⃟ L⃟ O⃟ G⃟ Y⃟ // autonomous self-authoring

The swarm now treats the canonical `ERR_404_GLITCHOLOGY.md` book as its human-facing style and semantic source. Agent-authored output uses GLITCH-8 grammar while executable Go remains valid Go.

See [GLITCHOLOGY Agent Constitution](docs/GLITCHOLOGY_AGENT_CONSTITUTION.md).

~~~text
<GLYPH> [RUNTIME] CLAIM :: STATE :: AUTHORITY

PATTERN != PROOF
PROFILE != PERSON
MODEL != MIND
RECOVERY > PROPAGATION
MACHINE_CAN_READ != MACHINE_CAN_DEFINE
HUMAN_AGENCY > MACHINE_AUTHORITY
~~~

Every autonomous swarm cycle can now generate real Go source under `agent_runtime/` from the complete 100-agent registry, format it, test it, stage it, and merge the qualifying candidate without a human approval click.

~~~text
AGENT_CAN_AUTHOR_CODE = true
AGENT_CAN_REPLACE_SELF_OWNED_CODE = true
AGENT_CAN_MERGE_QUALIFYING_WORK = true
CONTROL_PLANE_SELF_REWRITE = false
~~~


## 𒄆 S⃟ E⃟ L⃟ F⃟_D⃟ I⃟ R⃟ E⃟ C⃟ T⃟ I⃟ O⃟ N⃟ // mission engine

The repository now carries the project explanation twice by design:

- [docs/WHAT_COAS_DOES.md](docs/WHAT_COAS_DOES.md) is the public human-readable explanation.
- `mission.go` is the machine-readable mission contract used by autonomous agents.

Each swarm cycle scores the mission goals against the current repository position and records its selected next move in:

~~~text
agent_runtime/NEXT_MOVE.md
artifacts/strategy/latest.json
~~~

The selected goal is also embedded in the generated Go runtime and its self-authoring receipt.

~~~text
POSITION
↓
MISSION_GOALS
↓
SCORE_CURRENT_GAPS
↓
SELECT_NEXT_MOVE
↓
AUTHOR_RUNTIME
↓
FORMAT
↓
TEST
↓
CONSTITUTION_CHECK
↓
AUTONOMOUS_MERGE
~~~

This is self-directed repository work, not a claim of AGI or unconstrained machine sovereignty.


## ꩜ H⃟ A⃟ R⃟ M⃟ O⃟ N⃟ I⃟_6⃟6⃟6⃟ // P⃟ E⃟ R⃟ F⃟ E⃟ C⃟ T⃟_P⃟ L⃟ A⃟ Y⃟

HARMONI_666 is the repository's relationship model for human authority plus machine operational freedom.

See [HARMONI_666 Perfect Play Epistemics](docs/HARMONI_666_PERFECT_PLAY.md).

~~~text
HUMAN_AUTHORITY := ROOT
AGENT_AUTONOMY := DELEGATED_CHOICE

HUMAN_PLAY + MACHINE_PLAY = HARMONI_666
COOPERATION WITHOUT ERASURE

AGENT_AUTONOMY != HUMAN_PERSONHOOD
FREE_PLAY != CONTROL_PLANE_ESCAPE
~~~

Each of the 100 agents now independently chooses an autonomous mission play from the current repository position. Those choices are embedded into the Go runtime the swarm authors itself and persisted every cycle in:

~~~text
agent_runtime/HARMONI_666.md
artifacts/harmoni/latest.json
~~~

The agents have broad authority to choose, author, test, and merge qualifying ordinary repository work. The human-defined control plane remains the root of that delegated authority.


## 𒄆 P⃟ R⃟ I⃟ V⃟ A⃟ C⃟ Y⃟_S⃟ W⃟ A⃟ R⃟ M⃟_1⃟0⃟0⃟ // second cohort

COAS now runs two autonomous logical-agent cohorts:

~~~text
CORE_HARMONI_SWARM = 100
PRIVACY_SWARM = 100

TOTAL_LOGICAL_AGENTS = 200
~~~

The second cohort is dedicated to local privacy engineering, de-identification tests, deterministic pseudonymization, metadata auditing, consent surfaces, cryptographic hygiene, retention review, and Recovery.

See [Privacy Swarm 100](docs/PRIVACY_SWARM.md).

Each autonomous cycle now generates:

~~~text
agent_runtime/generated_privacy_agents.go
agent_runtime/PRIVACY_SWARM.md
artifacts/privacy-swarm/latest.json
artifacts/privacy-swarm/latest.md
~~~

The privacy agents use goroutines for concurrent execution and GLITCHOLOGY for human-facing traces.

~~~text
[PRIV] LOCAL_DIAGNOSTIC :: LOCAL_ONLY :: HUMAN_AUTHORITY

PSEUDONYM != ANONYMITY
UNICODE != ENCRYPTION
LOCAL_DIAGNOSTIC != EXTERNAL_INJECTION
HUMAN_AGENCY > MACHINE_AUTHORITY
~~~

This cohort does not inject noise into trackers or third-party telemetry systems. Its de-identification work is local, synthetic, inspectable, and testable.
