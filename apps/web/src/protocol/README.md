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
- Pinned artifact record: `upstream/provenance.json` (**R3 v2.7.0 兼容基线**；现行
  pin 权威见 `upstream/provenance-v2.9.json`；`provenance-v2.8.json` 为 v2.8 历史
  pin 记录，非现行权威)

## 版本协商（C-004 / F-001，GOAL-041 S2）

- 生产 host 接受 Manifest `2.7` / `2.8` / `2.9`（`APP_MANIFEST_SUPPORTED_PROTOCOL_VERSIONS`），
  页面文档接受 `2.7` / `2.8` / `2.9`（`HOST_SUPPORTED_PAGE_VERSIONS`）。
- 页面级版本 + 能力协商在 `load-page.ts` `loadPageDocument` 中 fail-closed
  （`UNSUPPORTED_PROTOCOL_VERSION` / `MISSING_REQUIRED_CAPABILITY`）；支持集见
  `apps/web/src/host/host-support.ts`（与 claim `support.capabilities` 一致）。
- Manifest envelope 与页面版本解耦（上游 decoupledVersions）：API 服务 2.7 envelope、
  35 页中 8 页声明 2.9，属合法协商形态。

The fixture test executes all 35 app-manifest cases that can be mapped to the
R3 host subset and all 16 app-navigation cases. Two upstream M1 validation
cases remain explicitly excluded because the upstream aggregate error envelope
uses `CAPABILITY_REQUIRED`, while this fail-fast host validator exposes
`MISSING_REQUIRED_CAPABILITY`.

R4 permissions and R5 page rendering live outside the original R3 goal.

R5 stage 3 (I-PROTO-004 = vendor) extends this package:

- Additional schemas in `docs/schemas/` (`node`, `page`, `action`, `reaction`,
  `component-registry`) plus the existing app-manifest schema
- Vendored fixture suites under `upstream/*.cases.json` with SHA pins in
  `upstream/provenance.json`
- Conformance adapters and tests under `conformance/` (Ajv structural checks +
  behavior suite runners). MVP D-EXPR reactions remain `$context` only; the
  upstream multi-round `$deps` reactions suite is accounted as excluded.
