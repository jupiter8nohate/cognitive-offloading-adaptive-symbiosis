# Privacy and Algorithmic Shielding

Unicode art is not encryption.

```text
GLYPH_COMPLEXITY != CONFIDENTIALITY
UNICODE != CRYPTOGRAPHY
```

Public glyphs may change how a message is perceived, but a crawler can still store and analyze them.

COAS therefore separates two layers:

## Public symbolic layer

Use Unicode, code-poetry, and CMB invariants for human-readable expression, identity, and cultural signaling.

## Private cryptographic layer

Use `SealDraft` for private offloaded drafts. It uses AES-256-GCM authenticated encryption with a 32-byte key supplied by the human operator.

Example:

```bash
export COAS_PRIVATE_DRAFT_KEY_HEX='<64 hex characters>'
go run ./cmd/coas-seal -in private-note.txt -out private-note.sealed.json
```

Do not commit the key. A sealed draft is private only while the key remains private and the plaintext is not separately published.

Encryption can protect stored or transmitted private content. It cannot stop profiling of material that the user intentionally publishes in plaintext.

Practical shielding therefore also relies on data minimization, local processing, reduced telemetry, explicit consent, retention limits, and careful choice of what is made public.
