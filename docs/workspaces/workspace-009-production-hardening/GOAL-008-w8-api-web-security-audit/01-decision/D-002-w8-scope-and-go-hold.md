---
id: D-002-w8-scope-and-go-hold
doc: decision-entry
goal: GOAL-008-w8-api-web-security-audit
date: 2026-08-20
status: accepted
parent: GOAL-001-production-hardening
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# D-002 · W8 required 修复范围与 I-002 go 暂挂裁决

### 触发

2026-08-20 目标轮次指令要求「推进工作区9目标8，直到顺利闭门」。A-001 已落盘并判定 `fail`：F-001（分页整数溢出/切片 panic DoS，High）与 F-002（生产 CSP 阻止 inline 首屏主题脚本，Low）为 2 条开放 required；I-002 要求裁决 required findings 是否影响 VP-008 `go` 消费有效性及修复优先级。

### 用户书面裁决（目标轮次指令）

1. **S2 修复范围**：整单采纳 A-001 F-001、F-002 为本波 required 修复范围；F-003/F-004 按 A-001 原判定处理（非 required，不作为关闭本目标的必改项）。
2. **I-002**：在 F-001/F-002 闭合前，**不对外宣称 VP-008 `go` 消费有效性**；闭合后恢复宣称前应复核（沿用 W7 D-002/D-003 的 go 暂挂-恢复纪律）。

### 决定

1. 本波实施范围 = A-001 两条 required（F-001、F-002）。不逐条做 residual/overruled；修复为可核对代码改动。
2. `00-meta.md` 成功标准 S2 标记为完成；S3 开始实施；S4 待 required=0 + self/independent 复核后关门。
3. I-002 在 `01-decision.md` 信息表与 meta 中状态改为 `verified/closed`（go 暂挂），证据 = 本 D-002 + F-001/F-002 闭合后的独立复核。
4. 审计模式按 workspace.md：security 高影响默认 `cross`：实施后 self；关门前本地 grok build（grok-4.6 · high）independent。

### 为什么

- 目标轮次指令显式要求推进并闭门，且 W7 先例为整单采纳 required + required 闭合前 go 暂挂；本波延续同一纪律，保证共享基架信任声明不被未闭合安全缺陷背书。
- F-001 是真实可触发的 DoS/panic 路径，必须修复；F-002 是生产 CSP 与功能不一致，必须修复。
- F-003/F-004 已由 A-001 明确为非阻断建议/条件风险，不纳入 required 关门范围。

### 未选方案

- 只修 F-001 不修 F-002：用户未选；F-002 同为 required，范围不成立。
- 不暂挂 go 宣称：用户未选；fail-closed 语义要求 required 闭合后再恢复宣称。
- 将 F-003/F-004 升为本波 required：A-001 未列 required，且 I-003 已登记不阻断本波。