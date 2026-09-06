---
id: GOAL-040-w28-admin-passwd-convention
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# 决策记录 · GOAL-040

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账见 `00-meta.md`；长决策与独立决策记录在 `01-decision/D-NNN-<slug>.md`，每条记录可独立阅读。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | `ADMIN_PASSWD` 权威 `.env` 位置 | 方案冻结 | S1 | 用户裁决 | verified | — | `apps/api/configs/.env`（用户 2026-09-06 确认） |
| I-002 | required | 是否本轮移除 TEST_ADMIN | 方案冻结 | S1 | 用户裁决 | verified | — | 包含（用户 2026-09-06 确认） |
| I-003 | required | `ADMIN_PASSWD` 是否被 API 读取/重置 | 方案冻结/实施 | S1 | 设计决策 | verified | — | 不读取、不重置（D-001） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-06 | ADMIN_PASSWD 声明约定 + TEST_ADMIN 机制退役 | accepted | `01-decision/D-001-admin-passwd-convention.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。
