---
id: GOAL-006-channel-provider-contract
doc: decision
status: active
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 0.1.0
---

# 决策记录 · GOAL-006

## 信息需求与阶段门禁

> 状态以 `00-meta.md` 为准。本目标关闭 Root I-011；I-007/I-008/I-012 已由 Root D-006 verified，本目标只消费。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-011 | required | mock 持久化：Store vs 扩容进程内；管理端列表/详情 | R5 方案冻结 | R5 | 本目标决策 | **verified**（D-002 §3，2026-08-24 用户裁决 DB 表 + 迁移） | — | Root I-011 / I-017-011 |
| I-010（预冻） | required | Resend 配置键与 fail-closed | R6 方案/实施 | R6 接入前 | 提前于本目标冻结 | **verified**（D-002 §4） | — | Root I-010 / I-017-010 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-24 | 开设 R5 子目标（合同尚未冻结） | accepted | [D-001-r5-goal-establishment.md](01-decision/D-001-r5-goal-establishment.md) |
| D-002 | 2026-08-24 | R5 渠道供应商合同冻结（关闭 I-011 / 预冻 I-010） | accepted | [D-002-r5-channel-contract-freeze.md](01-decision/D-002-r5-channel-contract-freeze.md) |
