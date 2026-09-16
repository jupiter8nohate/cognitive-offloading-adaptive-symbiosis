# COAS Clean-Room Runtime

This runtime is an original COAS implementation written from scratch around general software architecture concepts: policy gating, task planning, capability routing, bounded workers, memory, verification, provenance, and recovery.

It does not copy source code from OpenHands, PentAGI, Drosophila_brain_model, or other external agent frameworks.

## Design boundary

```text
HUMAN GOAL
    |
    v
COAS GOVERNOR
    |
    +-- MANUAL
    +-- ASSIST
    +-- AUTOMATE_REVERSIBLE
    |
    v
PLANNER
    |
    v
CAPABILITY ROUTER
    |
    v
BOUNDED WORKER
    |
    +-- evidence -> MEMORY
    |
    v
VERIFIER
    |
    v
RESULT
    |
    v
HUMAN JUDGMENT
```

## Core invariants

```text
OFFLOAD != SURRENDER
CAPABILITY != AUTHORITY
MEMORY != TRUTH
MODEL != MIND
PATTERN != PROOF
AUTOMATION != CONSENT
HUMAN_AGENCY > MACHINE_AUTHORITY
```

## Why this is independent

The runtime uses ordinary, non-exclusive software concepts rather than external implementation code. COAS defines its own types, authority decisions, execution loop, worker registry, memory contract, verification contract, and evidence classifications.

The implementation intentionally avoids external agent-framework dependencies. The initial reference runtime uses only the Go standard library.

## Runtime contracts

### Governor

Evaluates cognitive load, stakes, and reversibility. It determines the maximum delegated authority for the task.

### Planner

Selects the next capability required to advance the task. The initial planner is deterministic so behavior can be tested before model-backed planning is added.

### Worker

Performs one bounded action. Workers do not receive authority to expand their own permissions.

### Memory

Stores evidence and results. Retrieved memory is context, not truth, and must remain subject to verification.

### Verifier

Evaluates the candidate result after execution. Verification is separate from execution so a worker cannot silently certify itself.

## Extension path

Future adapters can add:

- model-backed planning
- file and repository workers
- sandboxed command execution
- structured persistent storage
- semantic retrieval
- browser/research workers
- webhook and scheduled task dispatch
- human approval workflows
- cryptographic provenance records

Those adapters should implement COAS interfaces rather than becoming the authority layer themselves.

## Dependency rule

```text
COAS_GOVERNANCE > EXECUTION_BACKEND
```

External models, tools, databases, or sandboxes may increase capability. They do not receive authority to redefine the COAS constitution, human meaning, consent, identity, or final judgment.
