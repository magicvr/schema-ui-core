---
id: GOAL-016-w14-rectification-batch-a
doc: decision
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# 决策记录 · GOAL-016

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | non-blocking | F-01 handler 目录端点路径与安全 | S2 F-01 | S1 | 方案 + 用户裁决（新增端点） | **closed** | — | D-001：`GET /api/scheduled-tasks/handlers` + `tasks.read` |
| I-002 | non-blocking | F-04 旧文案迁移策略 | S2 F-04 | S1 | 方案 | **closed** | — | D-001：保留旧 title/body 回退，新增可空 title_key/body_key，不迁移旧记录 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-17 | S1 方案冻结：F-01～F-04 设计（handler 端点 / scopes UI / 审计过滤+导出 / 通知 messageKey） | accepted | `01-decision/D-001-s1-freeze.md` |

> 本目标维度编号 D-NNN / E-NNN / A-NNN 均独立于 GOAL-015 的编号。
