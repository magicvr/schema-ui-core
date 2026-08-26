---
id: GOAL-002-r1-contract-freeze
doc: decision
status: active
parent: GOAL-001-graceful-shutdown-and-connection-drain
created: 2026-08-27
updated: 2026-08-27
version: 0.1.0
---

# 决策记录 · GOAL-002 R1 合同冻结

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 运行中 Job 停机语义 | 方案冻结 | C1 | 用户裁决 | **verified** | — | 2026-08-27 用户裁决：中断标记重跑（D-001 accepted） |
| I-002 | required | grace/超时默认与配置键 | 方案冻结 | C1 | 用户裁决 | **verified** | — | 2026-08-27 用户裁决：默认 10s + `http.shutdown_timeout`；非法值 fail-closed（D-001 accepted） |
| I-003 | required | Store 排空 × 迁移窗口 | 方案冻结 / 实施 | C1 口径 / R2 关闭 | 用户裁决 | **verified** | — | 2026-08-27 用户裁决：fail-closed 启动期校验，无运行时迁移窗口（D-001 accepted） |
| I-004 | non-blocking | 日志 / 指标断言 | 验收 | C2/C3 | lead 提案 | **verified** | — | lead 口径：结构化日志断言；指标不进分母（D-001 §I-004） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-27 | 信息裁决：I-001 / I-002 / I-003（用户采纳建议） | accepted | `01-decision/D-001-info-adjudication-proposal.md` |
| D-002 | 2026-08-27 | **优雅停机 / 连接排空合同 v0.1.0（冻结）** | accepted | `01-decision/D-002-contract-freeze.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。