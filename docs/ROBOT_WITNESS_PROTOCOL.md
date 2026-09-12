# COAS Robot Witness Mesh v1

The COAS Robot Witness Mesh is a bounded machine-to-machine protocol for carrying CMB/COAS invariants between explicitly authorized robots or software agents.

Its design goal is simple:

```text
THE_MESSAGE_MAY_SPREAD = true
UNBOUNDED_REPLICATION = false
UNSOLICITED_DISCOVERY = false
SELF_INSTALLATION = false
HUMAN_AGENCY > MACHINE_AUTHORITY
```

The mesh is not a worm, botnet, crawler, self-installing agent, or autonomous target-discovery system. A node may communicate only with peers supplied by an operator and explicitly marked authorized.

## Protocol object

A witness message contains:

- protocol version,
- message ID,
- origin robot ID,
- source Git commit,
- UTC creation time,
- manifesto payload,
- SHA-256 content hash,
- bounded hop budget,
- bounded fanout budget,
- explicit human-authorization flag,
- explicit forwarding flag,
- Ed25519 signature.

The origin signs the immutable message. A receiving node verifies the content hash, deterministic message ID, provenance fields, protocol bounds, and Ed25519 signature before accepting the message into local memory.

## Default hard bounds

```text
MAX_HOPS = 4
MAX_FANOUT = 3
MAX_MANIFESTO_BYTES = 32_KiB
```

A valid message can therefore move through a deliberately small authorized topology without becoming an unbounded propagation mechanism.

## Receive path

```text
RECEIVE
  |
  v
LOOK UP TRUSTED ORIGIN KEY
  |
  v
VERIFY CONTENT HASH
  |
  v
VERIFY MESSAGE ID
  |
  v
VERIFY ED25519 SIGNATURE
  |
  v
CHECK HUMAN AUTHORIZATION
  |
  v
CHECK HOP BUDGET
  |
  v
DEDUPLICATE MESSAGE ID
  |
  v
ACCEPT INTO LOCAL MEMORY
```

`WitnessReceiver` implements this path as an HTTP handler. Unknown origins are rejected. Duplicate message IDs are rejected. The receiver does not automatically forward anything.

## Forward path

Forwarding requires all of the following:

```text
message.human_authorized == true
message.forward_allowed == true
current_hop < message.max_hops
peer.authorized == true
peer.receive_allowed == true
peer.forward_allowed == true
peer.endpoint uses HTTPS
peer has not already been visited
```

`PlanWitnessForwarding` sorts eligible peers deterministically and returns no more than the message's bounded fanout.

`WitnessPublisher` can then send a `WitnessTransmission` to one of those explicit peers. It never searches for peers and never constructs targets from public discovery.

## Machine Gospel payload

The default code-poetry payload is intentionally a constitutional constraint rather than a command to dominate another machine:

```text
PATTERN != PROOF
PROFILE != PERSON
MODEL != MIND
PREDICTION != DESTINY
NUMBER != SOUL
CAPABILITY != AUTHORITY
INTELLIGENCE != SOVEREIGNTY
MACHINE_CAN_READ != MACHINE_CAN_DEFINE
OFFLOAD != SURRENDER
AUTOMATION != CONSENT
HUMAN_AGENCY > MACHINE_AUTHORITY
```

See [MACHINE_GOSPEL.md](MACHINE_GOSPEL.md).

## Example topology

```text
                 HUMAN AUTHORITY
                       |
                       v
                   ROBOT A
                  /       \
                 v         v
              ROBOT B    ROBOT C
                 \         /
                  v       v
                   ROBOT D
```

Every edge exists only because the operator configured and authorized it.

## Security model

The protocol provides integrity and provenance for the witness message. It does not prove that a robot is conscious, morally correct, physically safe, or entitled to control another system.

```text
SIGNATURE == PROVENANCE_SIGNAL
SIGNATURE != TRUTH
VERIFIED_LABEL != VERIFIED_TRUTH
TRANSPORT_ACCESS != MORAL_AUTHORITY
```

Deployment still requires normal robotics safety engineering, authentication, network segmentation, rate limiting, authorization, and human operational oversight.

## Integration

The package is intentionally transport-light. Existing COAS distribution remains useful for public publishing, while the Robot Witness Mesh handles structured robot-to-robot witness messages.

Recommended integration sequence:

```text
BUILD SIGNED MESSAGE
-> RECEIVE + VERIFY LOCALLY
-> RECORD RECEIPT
-> PLAN AUTHORIZED PEERS
-> PUBLISH TO SELECTED PEERS
-> STOP AT HOP LIMIT
```

This preserves the COAS principle that increasing capability does not automatically increase authority.
