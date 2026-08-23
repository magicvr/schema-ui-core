---
id: GOAL-001-outbound-mail
doc: decision-entry
record_id: D-002
status: accepted
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-002 · R1 发送合同冻结（关闭 I-001 / I-002）

### 触发

Root 路线图 R1 门禁：I-001（默认 sink 形态与测试取报文）、I-002（单次 `Send` 的 To 基数）须在 R1 方案冻结时由决策关闭。VP-017（`active`，VRev-037/038 均 pass）对两项均给出方向级建议；本决策将其落成可实施合同。

### 决定

采纳子目标 `GOAL-002-port-contract-freeze` D-001 的全部冻结项（该文件为合同正文）：

1. **端口**：`kernel.MailSender.Send(ctx, kernel.MailMessage)` 同步发送；`MailMessage = {To 单收件人, Subject, TextBody}`；From 不进消息体，由适配器按配置投递时补齐；端口级校验仅 To（非空、bare RFC 5322 addr-spec）。端口代码已落地：`apps/api/internal/kernel/mail.go`。
2. **I-001 → verified**：默认 sink = 进程内 capture sink（容量 1 环）+ 结构化日志双写；测试经 `internal/mail.CaptureSink.Last()` 取最后一封。落地归 R3。
3. **I-002 → verified**：单收件人（`To string`）；多收件人走将来加法演进。
4. **公共面规则**：公共契约仅可引用 `kernel.MailSender` / `kernel.MailMessage`；SMTP 客户端类型封在 `internal/mail`。核对归 R3 sweep。

### 为什么

- 与 VP-014 objectstore 先例同构（kernel 端口 + internal 适配器），薄内核零新依赖。
- From 由配置注入从端口层消除发件人伪造面（VP-009 关注点前移）。
- 两项均沿 VP 已公示建议方向冻结，无 P-004 冲突/residual 情形。

### 未选方案

见 `GOAL-002` D-001「未选方案」节（只写日志 sink / Mailhog 常驻 / 集合 To / 第三方邮件库）。
