---
id: GOAL-006-dual-path-evidence
doc: execution-entry
record_id: E-001
status: recorded
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## E-001 · R5 立项与证据方案冻结（checkpoint）

### 事实

- 立项 `GOAL-006-dual-path-evidence`（五件套 + ledger 目录 + attachments），`parent: GOAL-001-observability`，承载 Root 纲领阶段 R5。
- D-001 落盘：判据映射（VP 退出 3/4 ↔ 证据形态）、`otlp-sink` 工具决策、判定标准、关门顺序（GOAL-006 self → Root self → grok independent → 合并响应 → Root done）。
- 新增工具 `apps/api/cmd/otlp-sink`（极简 OTLP/HTTP 接收缘，`OTLP_SINK_ADDR` 可配，仅证据/排障用途）。
- 同步更新 goal-tree（树 + 表）。

### Git checkpoint

| hash | scope |
|------|-------|
| `8ddbb60` | workspace-015 GOAL-006 五件套（5 文件）+ goal-tree |
| `cf9df6c` | `apps/api/cmd/otlp-sink/main.go` |

### 备注

- 审计模式判定：R5 证据面向 VP「生产向验收」判据 → 关门审计走项目级独立审计路径（grok build /audit），self 先行。