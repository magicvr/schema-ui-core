---
id: GOAL-003-r2-memory-provider
doc: decision
status: active
parent: GOAL-001-cache-port
created: 2026-09-01
updated: 2026-09-01
version: 0.1.0
---

# 决策记录 · GOAL-003 R2 内存供应商

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-026-004 | non-blocking | mail `cachedAdapter` 迁移评估（版本戳 vs 通用 TTL 语义） | 判据 #2（评估面） | R3 | lead 评估 + 用户确认 | 待确认 | R3 目标承载（GOAL-004） | — |

（R2 无新 required 信息项；I-026-001/002/003 已 verified。）

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-01 | R2 方案冻结 v0.1.1（FIFO 驱逐 · maxEntries 进程总预算（A-002 F-001 用户裁决）· Typed 形态 · 配置键 · 审计模式；勘误 §1～§3） | accepted | `01-decision/D-001-r2-plan-freeze.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。