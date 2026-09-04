---
id: GOAL-003-r2-connection-settings
doc: decision
status: active
parent: GOAL-001-telegram-operator-console
created: 2026-09-04
updated: 2026-09-04
version: 0.3.1
---

# GOAL-003 · R2 决策索引

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-033-014 | required | mode/URL 的 seed、DB authoritative 与 Admin PATCH 优先级 | 方案 / C2 | C1 | 用户裁决并写 D-001 | **verified** | 未延期 | D-001 |
| I-033-015 | required | heartbeat 引用计数/单 lease 与 TTL | 方案 / C3/C4 | C1 | 用户裁决并写 D-001 | **verified** | 未延期 | D-001 |
| I-033-016 | required | getUpdates 长轮询 timeout 与 client 余量 | 方案 / C3 | C1 | 用户裁决并写 D-001 | **verified** | 未延期 | D-001 |
| I-033-017 | non-blocking | disabled profile 的 route/module surface 语义 | 实施 / C4 | C3 | provider/composition disabled profile 404 测试 | **verified** | 未延期 | A-014；Telegram settings、lease、webhook、schema 均未注册 |
| I-033-018 | non-blocking | HasBusinessHandlers 的 adapter 放置 | 实施 / C3 | C3 | 实现决策与测试 | **verified** | 已由 C3 记录 | A-009/A-012 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| [D-001-r2-implementation-parameters](01-decision/D-001-r2-implementation-parameters.md) | 2026-09-04 | R2 实施参数与接缝裁决 | accepted | `01-decision/D-001-r2-implementation-parameters.md` |

未决方案不写成 accepted；后续决策从 `D-002-*.md` 起递增落盘。
