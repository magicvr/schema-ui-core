---
id: GOAL-004-default-sink-surface-sweep
doc: execution-entry
record_id: E-001
goal: GOAL-004-default-sink-surface-sweep
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-001 · R3 落地：capture sink + composition 接线 + sweep（2026-08-22）

## 已发生事实

- C2 代码：
  - `apps/api/internal/mail/capture.go`：`CaptureSink`（容量 1 环 / `Last()` / `Reset()` / slog 结构化日志；Send 先按端口合同校验）。
  - `internal/mail/capture_test.go`：只留最后一封 / 校验拒收不入环 / Reset / 日志行内容断言 / nil logger 兜底。
  - `internal/config/config.go`：新增 `MailSMTPConfigured()`（validateMail 的 touched 判定重构为复用该方法，行为不变）。
  - `internal/composition/composition.go`：新增 `newMailSender(cfg, logger)` 并注册进 `fx.Provide`——未配置 → `*mail.CaptureSink`，显式配置 → `*mail.SMTP`。
  - `internal/composition/composition_mail_test.go`：fx 容器解析断言（缺省 capture、显式 SMTP）、端口 Send 后 `Last()` 可取回、构造器对部分块防御性 fail-closed。
- 公共面 sweep 证据（rg 全仓 apps/api）：
  - `net/smtp` / `tls.Dialer` / `smtp.PlainAuth` **仅**出现在 `internal/mail/smtp.go`；
  - `net/mail` 出现在 `internal/mail/smtp.go`、`internal/kernel/mail.go`（合同校验）、`internal/config/config.go`（from 校验）——均为地址解析器而非 SMTP 客户端类型；
  - `MailSender` / `CaptureSink` 引用仅存在于 kernel（端口）/ internal/mail（适配器）/ composition（接线）；handler 与 modules 包零引用。
- 测试实跑（go1.26.0 windows/amd64）：`go test ./internal/mail/ ./internal/config/ ./internal/kernel/ ./internal/handler/` 全绿（handler 131.9s）；`go test ./internal/composition/` 全绿（20.4s）；vet 干净。

## 证据

| 主张 | 路径 |
|------|------|
| 接线决策 | 本目标 `01-decision/D-001-default-sink-wiring.md` |
| capture sink | `apps/api/internal/mail/capture.go` |
| 容器接线 | `apps/api/internal/composition/composition.go`（newMailSender + fx.Provide） |

## 未做

- 未动 `readyz`（R4）；未做显式投递 live harness（R4）；未加 handler/模块消费方（本波无）。
