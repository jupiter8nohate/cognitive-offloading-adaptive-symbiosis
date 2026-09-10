# Preservation and Digital Footprints

COAS treats repository history as tamper-evident provenance, not as an unalterable ledger.

Git commits form a hash-linked history. A repository administrator can rewrite history, so the strongest claim is:

```text
COMMIT_HISTORY = TAMPER_EVIDENT
COMMIT_HISTORY != PHYSICALLY_IMMUTABLE
```

Every autonomous swarm cycle generates a preservation manifest at:

```text
artifacts/preservation/latest.json
artifacts/preservation/latest.md
```

The manifest records the verified source commit, file paths, byte sizes, SHA-256 hashes, and a SHA-256 hash of the manifest itself.

## Canonical asset tracking

The preservation engine currently watches for:

```text
ERR_404_GLITCHOLOGY.md
```

If the real file exists anywhere in the repository tree, its path and SHA-256 hash are recorded automatically. If it is absent, the manifest records that absence instead of inventing content.

## Decentralized survival

Durability improves when independent people or systems clone, fork, archive, or mirror the repository. A single GitHub repository is not a guarantee of permanent survival.

Recommended preservation pattern:

```text
PRIMARY_REPO
  |
  +-> NORMAL_GIT_CLONES
  +-> VOLUNTARY_FORKS
  +-> AUTHORIZED_MIRRORS
  +-> HASHED_PRESERVATION_MANIFESTS
  +-> SIGNED/VERIFIABLE RELEASE COPIES WHEN AVAILABLE
```

The repository's authorized distribution layer can broadcast source commit and content hashes to endpoints the operator controls, giving external observers additional reference points for provenance.
