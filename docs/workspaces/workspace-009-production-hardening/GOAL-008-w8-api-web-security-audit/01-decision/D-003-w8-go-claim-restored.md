---
id: D-003-w8-go-claim-restored
doc: decision-entry
goal: GOAL-008-w8-api-web-security-audit
date: 2026-08-20
status: accepted
parent: GOAL-001-production-hardening
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# D-003 · VP-008 go 宣称恢复（W8 F-001/F-002 闭合复核通过）

### 触发

D-002 规定：在 F-001/F-002 两条 required 闭合前，不对外宣称 VP-008 `go` 消费有效性；闭合后恢复宣称前应复核。

A-002 self pass + A-003 independent pass（grok build · grok-4.6 · reasoning high）确认 A-001 F-001/F-002 在现行代码中 genuine fixed；API `go test ./...`、Web `npm test`（1072）、`npm run build` 全绿。

### 决定

1. **VP-008 `go` 消费有效性从暂挂恢复为有效。**
2. I-002 状态更新为「已恢复」；证据 = A-002 + A-003 + D-002（原暂挂裁决）。
3. GOAL-008 成功标准 S4 标记完成；本目标按用户目标轮次指令关闭（status=done）。
4. 后续业务 VP 激活前仍需完成 VP-008 既有的消费前 freshness review；本恢复不改变 VP-008 `closed` 状态。

### 为什么

- A-003 为项目级独立 provider（grok build）出具的 close-out pass，并独立重跑全量回归；A-002 自审同向 pass。required 按 P-003 三路径以 `fixed` 合法闭合，开放 required = 0。
- 目标轮次指令要求推进至闭门，且 D-002 已裁决闭合后恢复前应复核；复核算术完成，故恢复 go 宣称。
- F-003/F-004 维持 A-001 原非阻断处置；不因 recommended 延迟恢复。

### 未选方案

- 维持暂挂直至 recommended 全清：用户未选；recommended 不阻断 required 闭合门禁。
- 仅依赖 A-002 self 恢复：A-003 independent 已存在，且为项目级既定 provider，无需降级。