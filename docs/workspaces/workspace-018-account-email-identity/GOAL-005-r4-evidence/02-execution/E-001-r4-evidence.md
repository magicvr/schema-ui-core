---
id: E-001
doc: execution-entry
goal: GOAL-005-r4-evidence
status: recorded
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-001 · R4 端到端证据落地（2026-08-24）

## 已发生事实

- 新增 `email_identity_e2e_test.go::TestR4EndToEndBindVerifyThroughMockChannel`：
  绑定 → 校验信经**真实** `mail.OutboxSink`（017 默认渠道适配器）写入出站记录 →
  从记录取码 → VerifyEmail → verified；唯一性 fail-closed 与无邮箱账号不受影响同链断言。
- **测试驱动出一处实现缺陷并修正**：渠道适配器自持存储事务，原「事务内发信」触发
  嵌套 Run 禁止（R1 v1.4 §2）。Bind/Resend 重构为两阶段派发 + 失败补偿
  （恢复先前身份/挑战快照），既有服务测试矩阵全数保持绿。
- 验证：authsession 包全量 ok（含新 e2e 用例）。

## 证据

| 主张 | 路径 |
|------|------|
| e2e 测试 | `apps/api/internal/modules/authsession/email_identity_e2e_test.go` |
| 证据包 | 本目标 `attachments/r4-evidence.md` |
