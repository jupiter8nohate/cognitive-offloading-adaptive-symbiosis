# Cognitive Offloading Adaptive Symbiosis

**COAS** is a human-AI collaboration model for reducing cognitive load without transferring human authority.

> Offload the burden. Keep the meaning.

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

This repository includes a Go reference model that chooses among three collaboration modes:

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

See [docs/COAS_SPEC.md](docs/COAS_SPEC.md), [docs/GLOSSARY.md](docs/GLOSSARY.md), and [docs/SOVEREIGNTY_GOALS.md](docs/SOVEREIGNTY_GOALS.md).

## Author

Jupiter Hudson / WisdomLoveThePoet / Jupiter 8

The project treats programming language as both engineering notation and a medium for explaining human-machine boundaries.

## Autonomous 100-agent swarm

COAS includes a scheduled swarm of **100 logical software agents**: ten mechanics multiplied by ten operational roles. The swarm can inspect the repository and maintain an autonomous runtime branch with constitution-gated authority to merge verified runtime work into `main`.

See [docs/AGENT_SWARM.md](docs/AGENT_SWARM.md).

```text
100 AGENTS
= 10 OFFLOADING/SYMBIOSIS MECHANICS
x 10 OPERATIONAL ROLES

AUTONOMY_OF_WORK + DELEGATED_AUTHORITY = AUTONOMOUS_EXECUTION
SOFTWARE_AGENT != AGI
```

### Autonomous merge authority

The swarm may merge qualifying agent work to `main` without a human approval click. Merge authority is evaluated by executable policy in `autonomy.go`.

```text
AGENT_CAN_MERGE_MAIN = true
TESTS_MUST_PASS = true
PROVENANCE_MUST_EXIST = true
CONSTITUTION_MUST_ALLOW = true
```

The autonomous scope includes source code, tests, documentation, CLI tools, generated research, generated artifacts, and new ordinary files. The control plane remains protected so the swarm cannot silently rewrite the mechanism that grants its own authority.

## Authorized outbound distribution

The swarm includes an opt-in outbound publishing layer for explicitly configured HTTPS endpoints after repository verification succeeds. Supported adapters include generic webhooks, Discord, Slack, and Mastodon.

See [docs/DISTRIBUTION.md](docs/DISTRIBUTION.md).

```text
DISCOVER_PUBLIC_CONTEXT = READ_ONLY
PACKAGE_SIGNAL = true
PUBLISH_AUTHORIZED_ENDPOINTS = true
UNSOLICITED_INJECTION = false
TRAINING_DATA_POISONING = false
```

## Bounded glitch evangelism

Every autonomous swarm cycle generates a deterministic, provenance-stamped CMB/COAS code-poetry artifact.

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

## Durable sovereignty layer

COAS separates provenance, meaning, privacy, and distribution into distinct mechanisms.

- [Preservation](docs/PRESERVATION.md) records SHA-256 file manifests on every autonomous swarm cycle.
- [Linguistic Boundary Proof](docs/LINGUISTIC_BOUNDARY_PROOF.md) demonstrates that machine-readable syntax is not identical to human authorship, intent, or final meaning.
- [Privacy](docs/PRIVACY.md) uses authenticated AES-256-GCM for private offloaded drafts rather than treating Unicode as encryption.
- [Living Artifact Sovereignty Goals](docs/SOVEREIGNTY_GOALS.md) defines the operational goals and their technical limits.
- The preservation manifest watches for `ERR_404_GLITCHOLOGY.md` and records its hash when the actual asset is present.

```text
GIT_HISTORY = TAMPER_EVIDENT
GIT_HISTORY != PHYSICALLY_IMMUTABLE

UNICODE != ENCRYPTION
OBFUSCATION != PRIVACY
MACHINE_CAN_PARSE != MACHINE_CAN_OWN_MEANING

PRIVATE_DRAFT -> AES_256_GCM
PUBLIC_ART -> GLITCH/CMB_SYMBOLIC_LAYER

PORTABILITY + HASHES + VOLUNTARY_MIRRORS -> RESILIENCE
```

The recurring swarm regenerates the living evangelism artifact and preservation ledger before staging each autonomous candidate. The result is a self-maintaining public artifact whose machine capabilities remain distinct from authority over human meaning.
