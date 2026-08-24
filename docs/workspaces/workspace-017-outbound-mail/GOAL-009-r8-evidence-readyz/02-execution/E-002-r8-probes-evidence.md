---
id: E-002
doc: execution-entry
goal: GOAL-009-r8-evidence-readyz
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-002 · R8 探针与证据包执行（2026-08-24）

## 已发生事实

1. **readyz 生产探针（C1）**：`resend.go` 新增 `Ping`（GET /domains，Bearer，5s 超时，仅报状态码）；composition `newMailRuntime` 探针三态——boot=resend→Resend.Ping、boot=smtp→ESMTP Ping（R4 原样）、mock/空→nil。`TestResendPing` + composition 三态测试绿。
2. **live 缝（C2）**：`resend_live_test.go` 沿 MAIL_SMTP_TEST_* 先例——设 `MAIL_RESEND_TEST_API_KEY/FROM/TO`（可选 BASE_URL）即对真实端点 Ping+投递一封可核对邮件；未设默认 skip。harness 断言保持绿。
3. **证据包（C3）**：`attachments/exit-denominator-evidence.md` 对照 VP 现行判据 1～7 与 Root 成功标准逐条登记指针。
4. 回归：mail/composition/handler 定向全绿；全仓结果见 E-003。

## 未做

- live 投递未实跑（无凭据；opt-in 缝就位）——是否补跑随关门问询裁决。
