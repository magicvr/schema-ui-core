---
id: D-002-w10-scope-and-go-hold
goal: GOAL-010-w10-api-web-security-audit
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-010-w10-api-web-security-audit
version: 0.1.0
---

# D-002 · W10 修复范围裁决与 go 宣称暂挂（2026-08-21）

## 决定（用户书面指令）

用户指令原文："整单采纳 F-001 + 六条 recommended，暂挂 go 宣称，开始修复"。

1. **required/recommended 修复范围 = 7 条**：F-001（HIGH · env.example 硬编码凭据）+ F-002～F-007（recommended）。informational F-008～F-012 不在本波实施范围（F-008/F-009 为已接受设计取舍；F-010～F-012 属环境/构建卫生，另行处理）。
2. **VP-008 go 消费有效性宣称：暂挂**——自本决定起至本波 required+recommended 全部合法闭合并复核通过为止，不宣称 VP-008 go 消费有效性（对齐 W7/W8/W9 D-00x 先例）。
3. **实施顺序**：S3 按上述 7 条实施 + API/Web 回归全绿 → A-002 self 审计 → S4 复核关门。
4. I-002 状态：open → **verified**（证据 = 本文件）。

## 未选方案

- **仅修 F-001 最小范围**：用户明确选择整单采纳，不采用最小范围。
- **informational 一并实施**：F-008/F-009 已有文档化取舍记录，重开需新决策；F-010～F-012 与安全门禁无关，避免波次范围蔓延。

## 对 go 宣称的影响

- 暂挂起点：2026-08-21（本决定落盘时）。
- 恢复条件：7 条 finding 全部 fixed 且 self/independent 复核确认无新引入缺陷后，由用户书面恢复（对齐 W9 D-004 模式）。