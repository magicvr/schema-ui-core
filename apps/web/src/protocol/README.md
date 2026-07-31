# Application manifest protocol

R3 implements the app-manifest and navigation subset pinned to
`schema-ui-docs@2.7.0` at commit
`ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`.

- Default fetch endpoint: `/.well-known/schema-ui/app-manifest.json`
- Local validation entry: `validateAppManifest()` in `app-manifest.ts`
- Route semantics: D4a literal count, template length, then declaration order
- Navigation projection: `top`, `sidebar`, and `user` in `../app/navigation.ts`
- Failure behavior: invalid or unavailable manifests stop page rendering; unknown
  routes render the shell fallback without guessing a page
- Pinned artifact record: `upstream/provenance.json`

The fixture test executes all 35 app-manifest cases that can be mapped to the
R3 host subset and all 16 app-navigation cases. Two upstream M1 validation
cases remain explicitly excluded because the upstream aggregate error envelope
uses `CAPABILITY_REQUIRED`, while this fail-fast host validator exposes
`MISSING_REQUIRED_CAPABILITY`. Negotiation and decoupled-version cases are
fixture-only adapters; the production host still accepts exactly protocol
`2.7`.

R4 permissions, R5 page rendering, and full protocol conformance remain outside
this goal.
