---
id: D-004-w9-go-restore
doc: decision-entry
goal: GOAL-009-w9-api-web-security-audit
record_id: D-004
status: accepted
parent: GOAL-001-production-hardening
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# D-004 · 恢复 VP-008 go 消费有效性宣称（2026-08-21）

### 触发

用户 2026-08-21 /govern 书面指令："先修正 A-005 的三条 recommended，然后写 D-004 恢复 VP-008 go 宣称，并关闭 GOAL-009"。

### 决定

1. **恢复 VP-008 go 消费有效性宣称**（解除 [D-003](D-003-w9-scope-and-go-hold.md) §4 的暂挂）。依据链：
   - S3：[E-004](../02-execution/E-004-w9-s3-implementation.md) —— D-002 消费 12 条 required 全部实施，API/Web 回归全绿；
   - S4 self：[A-004](../03-audit/A-004-w9-self.md) pass；
   - S4 independent：[A-005](../03-audit/A-005-w9-s4-independent.md)（grok-build · grok-4.6 · reasoning high · /audit）**pass**，12/12 genuine fixed、无新引入阻断缺陷、回归本会话复跑一致；
   - 闭合记录：[A-006](../03-audit/A-006-w9-a005-response.md) —— fixed ×12，开放 required = 0；
   - 加固收尾：[E-005](../02-execution/E-005-w9-recommended-hardening.md) —— A-005 三条 recommended 全部实施并锁定（L2 校验器接线生产路径、恢复码 CAS 换值令牌、6 组原缺陷形状回归锁），回归再次全绿。
2. A-005 三条 recommended 已在本波内实施完毕（E-005），不再留开放 recommended 项；后续如出现新缺陷按常规波次处理。
3. 本决策不改变 Charter primary workspace、不重开 Root、不触碰 VP-009 状态。

### 未选方案

- **维持暂挂至下一波**：用户未选；独立复核与加固收尾均已闭环，继续暂挂缺乏事实依据。
- **不写决策、口头恢复**：违反 W7/W8 以来的书面宣称纪律。

### 影响

- VP-008 go 消费有效性宣称恢复有效（自本决策落盘时点起）。
- GOAL-009 具备关门条件：S1–S4 全勾、开放 required = 0、无到期 required 信息项。

### 后续

GOAL-009 关门（status=done, 4/4）；Root 保持 active。
