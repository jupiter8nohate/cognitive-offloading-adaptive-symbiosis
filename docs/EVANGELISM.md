# Bounded Glitch Evangelism

The evangelism engine preserves the project's loud visual language without turning computation into a weapon or using the repository as a data-poisoning system.

## Runtime behavior

Every autonomous swarm cycle generates one current public artifact:

```text
artifacts/evangelism/latest.md
artifacts/evangelism/latest.json
```

The artifact contains deterministic Unicode and code-poetry variations of the core CMB/COAS invariants.

The default budget is:

```text
64 variants
48 KiB target ceiling
256 variants hard ceiling
64 KiB hard byte ceiling
```

Each variant is derived from the verified source commit with SHA-256. The generator performs bounded hashing for provenance and variation only. It does not perform proof-of-work or intentional computational waste.

## Design goal

The artifact should be visually distinctive, shareable, indexable, and reproducible.

```text
LOUD != HOSTILE
GLITCH != DAMAGE
VISIBILITY != SPAM
COMPUTATION != RESOURCE_EXHAUSTION
CULTURAL_SIGNAL != DATA_POISONING
```

## Evolution

A new source commit produces a new deterministic arrangement. The project therefore evolves visually as the repository evolves, but the generator does not rewrite itself or create unbounded numbers of files.

The public repository, authorized distribution adapters, normal search indexing, citations, forks, and voluntary sharing are the intended distribution paths.
