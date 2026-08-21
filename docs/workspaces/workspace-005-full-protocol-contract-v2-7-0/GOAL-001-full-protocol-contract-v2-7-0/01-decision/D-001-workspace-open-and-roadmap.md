---
id: GOAL-001-full-protocol-contract-v2-7-0
doc: decision-entry
record_id: D-001
status: accepted
parent: null
created: 2026-08-08
updated: 2026-08-08
version: 0.1.0
---

## D-001 · 开区与纲领路线图采纳

### 触发

- 用户 2026-08-08 经 `/vision` 确认：响应 VRev-013、激活 VP-006、执行 `/govern` 开启新工作区。
- VRev-013 **pass**；本 VP scope 无 open required；仓库级 `F-V018` 仅阻断 VP-005。
- VP-006 禁止在 closed VP-003/004 工作区吸收本意图。

### 决定

1. 新建 delivery 工作区 **`workspace-005-full-protocol-contract-v2-7-0`**（编号 = 现有最大工作区 + 1；slug 对齐 VP-006 id 惯例）。
2. Root **`GOAL-001-full-protocol-contract-v2-7-0`**，`parent: null`，`primary_plan` / `plan_refs` = `VP-006-full-protocol-contract-v2-7-0`，`vision_role: delivery`。
3. 采纳 VP-006 建议阶段 **S0–S5** 为本 Root 等权纲领检查点（P-001）；建区当日均未完成。
4. 覆盖表权威落点：S1 在本 Root `attachments/` 落盘 **`I-PROTO-FULL-001`**（新文件）；**禁止**就地改写 workspace-001 的 `I-PROTO-001 v0.1.3`。
5. disposition 纪律（继承 VP-006 exit 1）：默认 `include`；`include-partial` 仅保真/边角；范围收缩 → exclude + 用户 residual 或用户书面范围收缩。
6. **不**将激活/建区读成全量兼容已交付；**不**启动 VP-005。

### 为什么

- 结构选型：新纲领波次 → 新 VP + 独立 delivery 工作区（P-006）。
- 组合门闩：协议优先于视觉；MVP 子集不是终态成功条件。
- VRev-013 已确认方向已稳，可激活开区。

### 未选方案

| 方案 | 未选原因 |
|------|----------|
| 在 workspace-003/004 继续吸收协议扩张 | 违反 VP-006 硬约束；closed 区默认不接新区 |
| 先做 VP-005 视觉再补协议 | 用户 2026-08-08 书面否决；VP-005 实施冻结 |
| 就地升版改写 I-PROTO-001 v0.1.3 | 破坏历史 MVP 冻结证据；VRev-012 F-V022 已要求新文件 |
