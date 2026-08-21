---
id: GOAL-004-otel-traces
doc: execution-entry
record_id: E-001
status: recorded
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
---

## E-001 · R3 立项与 I-002 闭合（checkpoint）

### 事实

- 立项 `GOAL-004-otel-traces`（五件套 + ledger 目录 + attachments），`parent: GOAL-001-observability`，承载 Root 纲领阶段 R3。
- D-001 落盘（闭合 Root I-002 / VP I-015-002）：OTLP/HTTP protobuf（不做 gRPC）、`observability.traces.{enabled,endpoint,sample_ratio}` 三键 + fail-closed 规则、no-op 语义（缺省无 provider/无 span）、`ParentBased(TraceIDRatioBased)` 采样、span 面复用 R2 Wrap 拦截点（attrs method/route/url.path/status、≥500 Error）、W3C 传播、Resource 与生命周期、隐含备选（gRPC/全局 provider/同步导出）拒绝理由。
- 同步更新 goal-tree 与 workspace.md（R3 待立项 → 进行中先不表；scaffold 只登记在 tree）。

### Git checkpoint

| hash | scope |
|------|-------|
| `0470307` | `docs/workspaces/workspace-015-observability/GOAL-004-otel-traces/`（5 文件新建） |

### 备注

- 审计模式判定（P-004 §3.1）：常规、边界清楚、可逆的非平凡实施 → **`self`**；无冲突意见、无 residual 裁决需求。