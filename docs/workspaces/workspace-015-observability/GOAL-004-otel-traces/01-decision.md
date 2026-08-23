---
id: GOAL-004-otel-traces
doc: decision
status: active
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# 决策记录 · GOAL-004

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| （继承）I-002 | required | Tracing：OTLP HTTP vs gRPC、采样默认、未配置 endpoint 的 no-op 语义 | R3 方案冻结 / 实施 | R3 接入前 | 本目标 D-001 | verified（继承） | — | `01-decision/D-001-tracing-contract.md` |
| （继承）I-005 | required | request-id / correlation 如何写入 span | R4 方案冻结 | R4 关联前 | R4 决策 | open | — | 归 GOAL-005（R4） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-21 | Tracing 合同与 no-op 语义（闭合 I-002） | accepted | `01-decision/D-001-tracing-contract.md` |
