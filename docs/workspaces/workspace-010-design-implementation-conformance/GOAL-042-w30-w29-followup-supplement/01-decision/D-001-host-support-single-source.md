---
id: GOAL-042-w30-w29-followup-supplement
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# D-001 · claim↔host-support 单源真值与一致性契约

## 决定

1. **单源真值**：新建 `apps/web/src/host/host-support.json`，包含 `supportedPageVersions` 与 `supportedCapabilities`（19 能力全集）。这是 host 支持集的**唯一权威**：
   - `host-support.ts` 从该 JSON 导入并导出 `HOST_SUPPORTED_PAGE_VERSIONS` / `HOST_SUPPORTED_CAPABILITIES`（load-page 门禁、boot.ts HOST_SUPPORT 消费）。
   - `generate-claim.mjs` 构建期读取同一 JSON 生成 claim 的 `support.pageVersions` 与 `support.capabilities`（manifestVersions 保持既有常量）。
2. **一致性测试**（`host-support-consistency.test.ts`）：
   - claim 工件 `support.capabilities` / `support.pageVersions` 与 host-support.json 集合相等；
   - 每个已声明能力的 `mandatorySuites`（capability-registry.json）⊆ claim `conformance.suites`，且 suites 均 `result: pass`；
   - host-support.json 能力全集 ⊆ 上游 capability-registry 且 ID 合法（防拼写漂移）。
3. **F-001 遗留的同步注释**：`generate-claim.mjs` 与 `host-support.ts` 的「keep in sync」注释改为指向 JSON 单源。

## 未选方案

- 仅靠注释约定同步（S4 现状）：已由 W29 A-005 记录为 non-blocking 缺口，本决定闭环。
- 解析源文件文本比对：脆弱，JSON 单源更可维护。
