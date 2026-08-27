---
id: GOAL-008-mail-admin-surface
doc: decision
status: active
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 0.1.0
---

# 决策记录 · GOAL-008

## 信息需求与阶段门禁

> 状态以 Root `00-meta.md` 为准。本目标消费 I-009（D-007）与 I-012（D-006），无开放 required 信息项。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-009 | required | 热切换密钥/失败语义/单进程 | 本目标 C1/C2 | 已关闭 | — | **verified** | — | Root D-007（用户裁决） |
| I-012 | required | 设置「邮件」tab 形状 | 本目标 C2 | 已关闭 | — | **verified** | — | Root D-006；独立 API + mock 记录表 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-24 | 开设 R7 子目标（消费 I-009/I-012；实现细节授权在目标内决策） | accepted | [D-001-r7-goal-establishment.md](01-decision/D-001-r7-goal-establishment.md) |
