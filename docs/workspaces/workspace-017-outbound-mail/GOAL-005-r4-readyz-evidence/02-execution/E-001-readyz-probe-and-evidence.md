---
id: GOAL-005-r4-readyz-evidence
doc: execution-entry
record_id: E-001
goal: GOAL-005-r4-readyz-evidence
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-001 · R4 落地：Ping probe + readyz 扩依赖 + 显式路径证据（2026-08-22）

## 已发生事实

- C2 代码：
  - `internal/mail/smtp.go`：新增 `Ping(ctx)`（复用冻结拨号形态的最小 ESMTP 往返；`tlsConfig()` 抽取为 Send/Ping 共用）；`smtp_test.go` 新增 Ping 对 TLS fake 全绿 / 明文端点必败两测。
  - `internal/mail/smtp_live_test.go`：env-gated live 投递测试（镜像 `s3_live_test.go` 先例，六个 `MAIL_SMTP_TEST_*` 未设即 skip）。
  - `internal/composition/composition.go`：`newMailSender` 改三元返回 `(sender, probe, error)` 并在 `newMuxWithExtraProviders` 内直接构造（fx.Provide 移除）；probe 经 `RegisterWithMFAProbes` 变参传入；`newMux`/`newMuxWithExtraProviders` 增补 `logger *slog.Logger` 参数（三个既有测试调用点同步更新）。
  - `internal/composition/composition_mail_test.go`：改为直接构造断言——缺省 capture 且 probe=nil、显式 SMTP probe 非 nil、部分块防御性拒收。
  - `README.md`：readyz 探测面说明 + 「出站邮件」配置节。
- 测试实跑（go1.26.0 windows/amd64）：`go test ./internal/mail/ ./internal/composition/` 全绿（composition 20.0s 含双 profile lifecycle 启动不变量）；vet 干净。

## 证据

| 主张 | 路径 |
|------|------|
| R4 决策 | 本目标 `01-decision/D-001-readyz-probe-and-evidence.md` |
| 关门叙事 | 本目标 `01-decision/D-002-closeout-narrative-i005-i006.md` |
| Ping 实现 | `apps/api/internal/mail/smtp.go` |
| live harness | `apps/api/internal/mail/smtp_live_test.go` |

## 未做

- 未实际跑 live env-gated 测试（无真实端点凭据；离线证据已覆盖协议形态，live 测试留给 operator 按需触发）。
