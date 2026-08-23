---
id: D-003-w9-scope-and-go-hold
doc: decision-entry
goal: GOAL-009-w9-api-web-security-audit
record_id: D-003
status: accepted
parent: GOAL-001-production-hardening
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# D-003 · W9 required 修复范围与 I-002 go 暂挂

### 触发

用户在 D-002 清单调和后，于 2026-08-21 `/govern` 书面选择：**整单 12 条 required + 闭合前暂挂 VP-008 go 消费有效性宣称**（编排器推荐项，对齐 W7/W8）。

消费清单权威为 [D-002](D-002-w9-finding-inventory.md)，不是 A-001 原文「F-001～F-012」。

### 决定

1. **本波 required 实施范围（12）**：F-001、F-002（high）；F-004、F-005、F-006、F-007、F-008、F-009、F-010、F-011、F-012、F-025（med）。全部按 A-002 已确认的代码事实修复为可核对改动。不逐条 residual / overruled。
2. **F-003** 保持作废，不实施、不复用。
3. **F-013～F-023** 保持 recommended；**F-024** info。不阻断本波关门。本波可不修；若顺手修须在 execution 单列，不得冒充 required 闭合证据。
4. **I-002 / VP-008 go**：在上述 12 条 required 按 P-003 合法闭合（S4 self + independent）之前，**不对外宣称 VP-008 `go` 消费有效性**。恢复宣称另写 D-00N（沿用 W7/W8 D-003 纪律）。
5. **A-002 F-003/F-004 recommended**（把若干条降级、弱化 F-004 影响）**不采纳为本波范围裁剪**。分级说明可留在修复注释里，但不把 F-004/F-007/F-008/F-009/F-010 移出 required。
6. **审计模式**：security 高影响默认 `cross`。S3 后 self；S4 关门前 grok build（grok-4.6 · high · `/audit`）。I-003：A-002 已是 grok 对 A-001 的交叉复审，**不**替代 S4 修复后独立复核。
7. S2 勾选完成；S3 为下一步实施。

### 为什么

- 与 W7/W8 同一纪律：独立审 fail 且含 high 时，整单采纳 required，闭合前不让共享基架 `go` 背书未修缺陷。
- A-002 确认两条 high 与 MFA/方言/nginx 等主干成立；推荐降级未获用户采纳。
- F-025 已进入调和清单，整单即包含 cron DOM/DOW。

### 未选方案

- **A-002 最少集**（只修 F-001/F-002/F-005/F-006/F-011/F-012/F-025）：用户未选。
- **仅两条 high**：用户未选。
- **整单但不暂挂 go**：用户未选；两条 high 仍开放时对外宣称 `go` 有效，与 fail-closed 消费纪律冲突。

### 影响

- I-002 → `verified`（范围 = 12 条 + go 暂挂）。
- 成功标准 S2 勾选；progress 2/4。
- 不实施代码（S3 另步）。不恢复 go。

### 后续

S3：按本范围实施 + API/Web 回归。S4：self + grok independent；required=0 后再议 go 恢复。
