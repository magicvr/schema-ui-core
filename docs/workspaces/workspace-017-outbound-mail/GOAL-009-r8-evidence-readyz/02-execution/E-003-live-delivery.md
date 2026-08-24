---
id: E-003
doc: execution-entry
goal: GOAL-009-r8-evidence-readyz
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-003 · live 投递实跑记录与全量回归（2026-08-24）

## 已发生事实

1. **live 投递实跑成功**（用户裁决补跑并提供本地凭据 `apps/api/configs/.env`）：
   - 第一次尝试：`from=no-reply@eshowy.top` → Ping 通过（密钥有效、账户端点可达），Send 被 Resend 拒绝：**403 域名未验证**（"The eshowy.top domain is not verified"）。这是账户运营状态，非代码缺陷。
   - 第二次尝试：`from=onboarding@resend.dev`（Resend 官方沙箱发件人，仅可投递到账户本人邮箱）→ **PASS**：Ping 2xx + POST /emails 2xx，一封真实投递至 `magicvr@hotmail.com` 可核对。
   - 运营备注：生产启用 `eshowy.top` 发件前须在 resend.com/domains 完成域名验证。
2. 全量回归：`go test ./...` 全绿（含本条 live 测试在无 env 时默认 skip 的 CI 路径）。
3. 判据 3 以「live」面满足（高于等价 harness 下限）；证据包 attachments 已同步更新。

## 证据

| 主张 | 路径 |
|------|------|
| live 测试缝 | `apps/api/internal/mail/resend_live_test.go` |
| 本次运行 | 上方运行日志（PASS；第一次尝试 403 报文原文见审计附件引用） |
