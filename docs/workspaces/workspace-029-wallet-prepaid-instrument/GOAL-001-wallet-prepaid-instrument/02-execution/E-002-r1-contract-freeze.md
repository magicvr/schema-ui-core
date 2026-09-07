---
id: GOAL-001-wallet-prepaid-instrument
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-wallet-prepaid-instrument
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-002 · R1 核心合同冻结与选型用户裁决留痕

## 2026-09-02 · R1 合同冻结

### 已发生事实

1. 依据 AGENTS P-004 及用户指令中“方案选型等关键决策请询问用户，禁止擅自静默代替用户决策”的硬约束，针对 I-029-001、I-029-002、I-029-003、I-029-005、I-029-006 共 5 项方案选型完成调研并向用户发起决策提问。
2. 用户已逐项书面确认裁决结果：
   - I-029-001：主体落点为独立通道无关 `subjects` 表，`owner_type` 扩充 `subject`，`OwnerExistsFunc` 校验主体。
   - I-029-002：核销入金复用 `adjust` + `ref_type='voucher'`。
   - I-029-003：新增细粒度权限键 `wallet.voucher.issue`。
   - I-029-005：本波仅交付 Go 模块内部 API `Redeem(ctx, subjectID, code)`。
   - I-029-006：高熵码 + SHA-256 + 单事务 CAS 原子核销入金，并发失败 fail-closed。
3. 决策落盘至 `01-decision/D-002-r1-contract-freeze.md`，更新 Root `00-meta.md` 中的信息需求表格，关闭相关 required 门禁。
4. Root progress 更新为 1/4（R1 合同冻结完成）。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 决策落盘 | `docs/workspaces/workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument/01-decision/D-002-r1-contract-freeze.md` |
| 信息需求闭合 | `docs/workspaces/workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument/00-meta.md` |
| 用户裁决记录 | 会话问答交互（`ask_user_question` 确认） |
