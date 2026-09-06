---
id: GOAL-041-w29-api-web-protocol-conformance
doc: execution-entry
record_id: E-005
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# E-005 · S4 执行 — API/Web 实现整改（C-001/002/003/004/006/010 + C-005 子项）

## 已执行事实（2026-09-06）

1. **C-001 provenance 路径与数量口径**：`apps/web/src/protocol/upstream/provenance-v2.9.json` 将 fixture-suite schema 条目改为规范上游路径 `conformance/schemas/fixture-suite.schema.json`（本仓仍 vendor 于 `docs/schemas/`），note 计数修正为「docs/schemas 10 + conformance/schemas 1 + 19 vendored suites」；`stage3-fixtures.test.ts` 增加「provenance 条目路径 = 上游规范路径」守卫 + 支持 `conformance/schemas/` 前缀的本地映射。
2. **C-002 claim 报告 artifactVersion**：`generate-claim.mjs` `pinnedUpstream.artifactVersion` `2.8.0` → `2.9.0`；重生成 `public/protocol/` 三件套（buildId `git:13670154…`）；`claim-artifact.test.ts` 现对 19 能力/12 suites 的 claim 通过 C0/C1 校验。
3. **C-003 台账卫生**：`provenance.json` note 前缀标注「R3 v2.7.0 兼容基线，仅回归；现行 pin 权威 = provenance-v2.9.json」；`provenance-v2.8.json` 因 vision/Charter 历史记录引用而**保留为历史 pin 记录**（README 注明非现行权威，未删除——A-002 的「死数据」仅指无代码消费者，删除会破坏历史引用）；`protocol/README.md` 增加「版本协商」节并删除「production host 仅接受 2.7」过时表述；`host/bootstrap.ts` 头注释指向 v2.9 pin。
4. **C-004 / F-001 生产页面级版本+能力门禁 + 支持集扩展**（用户裁决：严格 fail-closed + 同步扩支持集）：
   - 新增 `apps/web/src/host/host-support.ts`（单源真值：`HOST_SUPPORTED_PAGE_VERSIONS = [2.7,2.8,2.9]`、`HOST_SUPPORTED_CAPABILITIES` = 19 能力全集）；`boot.ts` HOST_SUPPORT 复用之。
   - `load-page.ts` `loadPageDocument` 在 D-VAL 后新增页面级协商：`UNSUPPORTED_PROTOCOL_VERSION` / `MISSING_REQUIRED_CAPABILITY`（含缺失能力 issues），fail-closed；`load-page.test.ts` 增 2 条负例。
   - `generate-claim.mjs` claim `support.capabilities` 扩至 19 能力、`conformance.suites` 扩至 12 个 mandatory suites 并集（app-manifest/app-navigation/host-bootstrap/host-failure/host-conformance-claim/request-construction/component-format/response-mapping/search-table/permissions-inheritance/table-sort/uploads），与 host-support.ts 同步；重生成 claim。
   - 验证：35 页全部通过页面级门禁（D-VAL 35/35 + representative/集成测试绿），claim C1 校验通过。
5. **C-005 子项（用户裁决：删除未使用能力声明）**：`digitaloffer-entitlements.json` 移除 `data.route-binding`、`digitaloffer-offers.json` 移除 `form.controls.readonly`；新增 `src/protocol/capability-declaration.guard.test.ts`——仅对 v2.9 能力对（data.route-binding / form.controls.readonly）强制「声明 ⊆ 使用」（防止 C-005 类漂移复发）。**范围说明**：legacy 能力（permissions.inheritance / actions.row.request / form.controls.extended 等）在本仓广泛保守声明（上游单向约束合法），全量声明-使用审计为 **non-blocking 后续项**（见 E-005 §附注），不属 C-005 决策范围。
6. **C-006 dogfood 防漂移守卫**：新增 `src/protocol/dogfood-manifest.guard.test.ts`——从全部模块 `manifest/fragment.json`（递归）推导页面联合，断言 mvp/admin-dogfood 每个 pageId/route/schemaUrl 与联合一致。**守卫当场捕获一处真实漂移**：admin-dogfood 的 `dictionary-entries` route 为旧值 `/dictionary-entries`，模块 fragment 为参数化 `/dictionary-entries/{dictKey}`；已修正 fixture 并重算 `STATIC_MANIFEST_SHA256`（`0efb4054…` → `ec00f6f1…`）。
7. **C-010 未知 custom 呈现对齐**（用户裁决：明显占位 + console.error）：`render.tsx` custom 分发未命中时渲染 `role="alert"` + 红色边框占位（含 component 键）+ `console.error`（含 node id）；`render.test.tsx` 增 C-010 用例；`representative-pages.test.tsx` 补 6 个 custom 组件 side-effect import（activity-export/data-permission-scopes/import-template-download/invite-issue-card/invite-resend-dialog/password-policy-tab）使代表页真实渲染而非回退占位。

## 验证

- 定向批次：13 files / **604 tests PASS**（stage3 / load-page / render / dval 35/35 / custom guard / capability-declaration / dogfood / claim-artifact / app-manifest / upstream-fixtures / upstream-host-fixtures / representative / boot·bootstrap）。
- 全量回归：`npm test`（vitest 全量）、`npm run build`（tsc + vite build，prebuild 重生成 claim）、`go test ./...`（api 全量）——结果见下方补录。
- Go 定向：`go test ./internal/manifest ./internal/composition ./modules/digitaloffer/...` PASS（digitaloffer schema JSON 改动不影响 Go 行为）。

## 产物

- 产品/测试改动：`host/host-support.ts`（新）、`host/boot.ts`、`protocol/load-page.ts`、`protocol/load-page.test.ts`、`renderer/render.tsx`、`renderer/render.test.tsx`、`renderer/representative-pages.test.tsx`、`protocol/capability-declaration.guard.test.ts`（新）、`protocol/dogfood-manifest.guard.test.ts`（新）、`protocol/conformance/stage3-fixtures.test.ts`、`protocol/upstream-fixtures.test.ts`、`test-fixtures/app-manifest.admin-dogfood.json`、`protocol/upstream/provenance-v2.9.json`、`protocol/upstream/provenance.json`、`protocol/README.md`、`host/bootstrap.ts`、`scripts/generate-claim.mjs`、`public/protocol/*`（重生成）、`api/modules/digitaloffer/schema/*.json` ×2。
- 治理记录：`02-execution/E-005-…`（本条）、`03-audit/A-005-s4-implementation-self.md`（自审）、D-002 §4b F-001 闭合、progress 4/6。

## I-007 说明（go 影响，预判，全量回归后定稿）

本波改动面：provenance/claim 证据工件、claim 生成脚本、README、host-support/load-page/renderer 运行时门禁与呈现、schema 能力声明清理（digitaloffer，仅元数据）、守卫测试。**未改 Profile 默认集、模块矩阵、Manifest 装配语义或共同门禁解释** → VP-008 `go` 消费有效性判定为**无影响不暂挂**（以 S5 记录定稿）。

## 附注（non-blocking 后续项）

- **legacy 能力保守声明全量审计**：account/activity/notifications 等多页保守声明 permissions.inheritance / actions.row.request / form.controls.extended 等且无对应使用；上游单向约束下合法，未在本波清理。后续波次可评估「删除未使用声明」或补文档化意图；复核触发 = 本波 S6 或新增页面/控件时。
