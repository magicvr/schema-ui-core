---
id: E-011
doc: execution-entry
goal: GOAL-001-account-email-identity
status: recorded
parent: null
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# E-011 · A-003 响应：recommended F-001 顺手修正（2026-08-25 · 关门后维护）

## 已发生事实

- A-003（independent · 2026-08-25 · pass · 0 required）新发现 **F-001 recommended/low**：未绑定账号（email IS NULL）调用 verify 时，users 无行错误被当作硬存储错误上抛，HTTP 面返回 500 INTERNAL 而非受控 `EMAIL_NOT_PENDING` 409。
- **fixed**：`evaluateVerification` 将 `kernel.ErrNoRows` 判为受控 outcome（真实存储错误仍上抛）；新增回归用例 `TestVerifyUnboundAccountIsControlledNotPending`。
- 复跑：authsession 全包 ok。Root 保持 `done`（关门状态不重开；本条为关门后低危契约修缮）。
- N-1 / N-2 维持 known-boundary 台账（A-003 建议路径）。

## 证据

| 主张 | 路径 |
|------|------|
| 意见原文 | 本目标 `03-audit/A-003-independent-workspace-code-closeout.md` |
| 修复与用例 | `apps/api/internal/modules/authsession/email_identity.go` + `email_identity_test.go` |
