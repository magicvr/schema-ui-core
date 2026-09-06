---
id: GOAL-041-w29-api-web-protocol-conformance
doc: execution-entry
record_id: E-006
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# E-006 · S5 执行 — 运行时符合性验证

## 已执行事实（2026-09-06）

1. **全分母渲染链验证（35/35）**：新增 `apps/web/src/renderer/denominator-render.test.tsx`——递归收集全部 35 份模块 schema，逐页执行 `validatePageDocument`（D-VAL）→ `loadPageDocument`（含 F-001 页面级版本+能力协商）→ `RenderPage`（AuthProvider + 15 个 custom 组件 side-effect import + fixture fetcher + SchemaTable），断言非空渲染且无 schema-error 表面/未知 custom 占位。**36 tests PASS（含 my-wallet 需 AuthProvider 的修正）**。至此 D-VAL + Load + Render 覆盖 35/35。
2. **HTTP Manifest 快照 × profile 矩阵**：新增 `apps/api/internal/composition/s5_manifest_snapshot_test.go`——真实组合根（ResolvePlan → store → newMux，telegram 组合装配真实 TelegramRuntime）经 httptest GET `/.well-known/schema-ui/app-manifest.json`，断言页面数冻结投影（mvp 6 / admin 22 / demo 14 / admin+biz.digital-offer 25 / admin+channel.telegram 24）与 envelope（protocolVersion 2.7 + app.manifest）；`S5_SNAPSHOT_DIR` 时写证据快照到 `attachments/S5-manifest-snapshots/<组合>/app-manifest.json`（5 组合已生成）。**5/5 PASS**。期间修复测试自身 append 别名 bug（共享 backing array 导致 digitaloffer 组合被 telegram 覆盖）。
3. **覆盖矩阵落盘**：`attachments/S5-coverage-matrix.md`——35 页 ×（D-VAL/Load+Negotiate/Render/行为触达）+ 失败路径证据 + 快照表。
4. **失败路径**：S4 已补页面级协商负例（UNSUPPORTED_PROTOCOL_VERSION / MISSING_REQUIRED_CAPABILITY）、C-010 占位、未知 handler fail-closed；S5 全量回归复跑确认。
5. **全量回归（复跑）**：Web vitest 全量（含新增 denominator-render 36 tests）**PASS**；`tsc -b` + `vite build` 0；Go `go test ./...` 全量 0 FAIL（含新增 s5_manifest_snapshot_test.go）。

## I-007 · VP-008 `go` 消费影响判定（定稿）

本波（S2～S5）变更面：治理文档、provenance/claim 证据工件、claim 生成脚本、README、host-support/load-page/renderer（运行时门禁与失败呈现，页面行为对已支持能力无变化）、digitaloffer schema 能力声明清理（元数据）、守卫测试与分母测试、manifest 快照测试。**未改 Profile 默认集（profileDefaults）、模块矩阵（BuiltinModules）、Manifest 装配语义（Aggregate/StampHomePageRef/导航序）或共同门禁解释**；页面级门禁接入后 35 页全过（零现网页面 fail-closed）。→ **VP-008 `go` 消费有效性：无影响，不暂挂**（对照 workspace-009/010 惯例，记录于本 E 条目）。

## 产物

- `apps/web/src/renderer/denominator-render.test.tsx`（新）
- `apps/api/internal/composition/s5_manifest_snapshot_test.go`（新）
- `attachments/S5-coverage-matrix.md`（新）
- `attachments/S5-manifest-snapshots/<mvp|admin|demo|admin+digitaloffer|admin+telegram>/app-manifest.json`（新，5 份）

## 阶段结论

S5 完成：35/35 分母（D-VAL + Load+Negotiate + Render）、5 组合 HTTP Manifest 快照、失败路径、全量回归均绿；I-006 → verified；I-007 → 无影响不暂挂（证据如上）。S6 关门审计（self + grok build independent + 用户确认）待执行。
