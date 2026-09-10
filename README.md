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

COAS now includes a scheduled swarm of **100 logical software agents**: ten mechanics multiplied by ten operational roles. The swarm can inspect the repository and maintain an autonomous runtime branch without granting itself authority to merge into `main`.

See [docs/AGENT_SWARM.md](docs/AGENT_SWARM.md).

```text
100 AGENTS
= 10 OFFLOADING/SYMBIOSIS MECHANICS
x 10 OPERATIONAL ROLES

AUTONOMY_OF_WORK != AUTONOMY_OF_AUTHORITY
SOFTWARE_AGENT != AGI
```
