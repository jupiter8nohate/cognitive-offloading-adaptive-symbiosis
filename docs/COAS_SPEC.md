# COAS Specification

## 1. Purpose

Cognitive Offloading Adaptive Symbiosis is a collaboration protocol for deciding what work a machine should carry and what authority must remain with a human.

The protocol optimizes for reduced cognitive burden while preserving human agency.

## 2. Inputs

A task is evaluated along three dimensions:

- **Cognitive load**: mental effort required to complete the task, from 0 to 100.
- **Stakes**: potential consequence of error or misuse, from 0 to 100.
- **Reversibility**: how easily an action can be undone, from 0 to 100.

These values are decision aids, not measurements of a person's intelligence, competence, or mental state.

## 3. Collaboration modes

### MANUAL

Use when the task is low burden or machine assistance would add unnecessary complexity.

### ASSIST

Use when support is useful but the task is high-stakes, difficult to reverse, ambiguous, or meaning-sensitive.

### AUTOMATE_REVERSIBLE

Use only for bounded work that is sufficiently reversible and low enough in consequence to safely automate after human confirmation.

## 4. Decision policy

Reference policy:

```text
if stakes >= 80:
    mode = ASSIST
else if cognitive_load < 30:
    mode = MANUAL
else if reversibility < 40:
    mode = ASSIST
else if cognitive_load >= 70 and stakes < 60:
    mode = AUTOMATE_REVERSIBLE
else:
    mode = ASSIST
```

The policy is deliberately conservative. It can be replaced by a more sophisticated policy only if the replacement preserves the authority boundaries below.

## 5. Authority boundary

Machine capabilities may include:

```text
remember
retrieve
search
calculate
organize
compare
summarize
simulate
draft
execute_reversible_actions
```

Human authority includes:

```text
set_goals
grant_consent
define_identity
assign_meaning
make_value_judgments
accept_consequences
override_system
withdraw_permission
```

Core rule:

```text
CAPABILITY != AUTHORITY
```

## 6. Adaptive behavior

The system may adapt:

- how much detail it provides,
- which tasks it proposes to offload,
- how much confirmation it requests,
- how it presents options,
- how it learns explicit user preferences.

The system must not infer that repeated acceptance creates permanent consent.

```text
PAST_CONSENT != PRESENT_CONSENT
PREFERENCE != COMMAND
PATTERN != PROOF
```

## 7. Cognitive dependency guard

Successful offloading should reduce unnecessary burden without eroding the user's ability to understand, inspect, or override the process.

For consequential tasks, the system should preserve:

- explanation,
- provenance,
- alternatives,
- uncertainty,
- a manual path,
- an override path.

## 8. Non-goals

COAS does not claim to:

- measure a person's intelligence,
- diagnose cognitive conditions,
- infer private intent,
- replace professional judgment,
- make machine outputs morally authoritative,
- prove that a model understands a person.

## 9. Verification principle

A useful COAS implementation should be testable against explicit invariants:

```text
HIGH_STAKES -> NO_AUTONOMOUS_FINAL_JUDGMENT
LOW_REVERSIBILITY -> NO_UNSUPERVISED_EXECUTION
AUTOMATION -> HUMAN_CONFIRMATION
HUMAN_OVERRIDE -> ALWAYS_AVAILABLE
MODEL_OUTPUT -> EXPLAINABLE_AS_RECOMMENDATION
```
