---
id: GOAL-041-w29-api-web-protocol-conformance
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# E-002 · S1 协议分母、页面控件目录与候选矩阵

## 已执行事实

1. 固定上游正式身份：本机只读 `schema-ui-docs` clone 的 annotated tag `v2.9.0` 解引用为 `81aa1d8954717f4ebdcc695eed6fafaeafcebe8d`；本仓被审快照为 `ce66abab0897b644c695b8365f8324598bfe2bc8`。
2. 建立机器协议分母：11 个 schema/registry 工件、24 个 component、19 个 capability、20 个上游 fixture suites / 450 cases。
3. 核验本仓 `provenance-v2.9.json`：30 个登记 artifact；本地 pinned 字节均与对应上游 tag 工件 digest 一致。发现 `fixture-suite.schema.json` 的 provenance 来源路径与 tag 真实路径不同，且 note 的 fixture 数量口径与 19 个实际 vendor suites 不一致，已登记 C-001。
4. 从 17 个 API Manifest fragments、35 个 schema 文档建立 35 页完整静态目录：重复 pageId 0、缺失 schema 0、orphan schema 0。
5. 用临时 Go test 调用真实 `apps/api/internal/manifest.ForModulesWithFragments` 并按 profile 过滤 fragments，得到 MVP 6 页、Admin 22 页、Demo 14 页、全 fragment 宇宙 35 页，Manifest 聚合版本均为 2.7。临时测试文件在取得输出后删除，未留下产品代码变更。
6. 建立 Web 分母：11 个 Renderer node types、14 个 form control types、15 个 custom component 注册键（12 个页面 `body` custom、1 个 action `content` custom、2 个 `afterComponent`）、5 个 custom action handlers；当前 schema 使用键均有注册/allowlist。
7. 记录页面触达现状：representative 测试覆盖 10/35 页；两份 dogfood Manifest 为 13/26 页，与当前真实 profile 投影不同。
8. 汇总 14 个 `collecting / pending-S2` 候选；未把任何候选提前写成 implementation gap、upstream protocol gap、custom 或 explicitly-out。

## 产物

- `attachments/S1-protocol-denominator-v2.9.md`
- `attachments/S1-protocol-denominator-v2.9.json`
- `attachments/S1-api-web-page-control-catalog.md`
- `attachments/S1-api-web-page-control-inventory.json`
- `attachments/S1-candidate-matrix.md`

## 定向验证

- `go test -count=1 ./internal/manifest ./internal/composition`（从 `apps/api` 执行）：PASS。
- `npm test -- src/protocol/conformance/stage3-fixtures.test.ts src/protocol/upstream-fixtures.test.ts src/protocol/upstream-host-fixtures.test.ts src/protocol/app-manifest.test.ts src/renderer/representative-pages.test.tsx`（从 `apps/web` 执行）：5 files / 458 tests PASS。
- v2.9 provenance/digest 只读重算：本仓 30 个登记 artifact 与对应上游 tag 字节一致；上游 20 suites / 450 cases，本仓 vendor 19 suites / 437 cases。

## 阶段结论

S1 的目标是冻结分母并登记未知，而不是给出协议符合性 verdict。I-001 与 I-002 已由证据验证（`verified`）；I-003 进入 `collecting`，C-001～C-014 留给 S2 逐项分类与 cross 方案审视。未修改 `apps/api` 或 `apps/web` 产品实现，未创建上游协议增补报告。
