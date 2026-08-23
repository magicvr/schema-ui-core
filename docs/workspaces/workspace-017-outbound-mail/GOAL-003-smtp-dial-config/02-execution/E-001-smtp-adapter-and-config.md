---
id: GOAL-003-smtp-dial-config
doc: execution-entry
record_id: E-001
goal: GOAL-003-smtp-dial-config
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-001 · R2 落地：SMTP 适配器 + 配置面（2026-08-22）

## 已发生事实

- C1：D-001 冻结拨号路径与配置键；Root D-003 同步关闭 I-003/I-004（→ verified）。
- C2 代码：
  - `apps/api/internal/mail/smtp.go`：`NewSMTP`（fail-closed 校验 host/username/password/from、port 0→465、from bare addr-spec）+ `Send`（隐式 TLS、MinVersion TLS1.2、证书校验恒开、AUTH PLAIN over TLS、Subject 控制字符拒发守卫、正文 CRLF 规范化）。
  - `internal/mail/smtp_test.go`：本地 throwaway CA + 叶证书的 TLS loopback harness，实跑 EHLO/AUTH PLAIN/MAIL FROM/RCPT TO/DATA 全会话断言；明文端点必须失败（无 STARTTLS/明文回退）；header 注入拒发；取消 context 拒发。
  - `internal/config/config.go`：新增 `MailSMTP{Host,Port,Username,Password,From}` 字段 + `yaml mail.smtp.*` 解析 + `MAIL_SMTP_*` env 覆盖（非法显式端口在 Load 即拒）+ `validateMail()` 挂入 `ValidateProd` 链（每个环境含 development）。
  - `internal/config/config_mail_test.go`：未配置放行 / YAML 解析 / env 覆盖 / 部分块 fail-closed 报键名不回显值 / 端口范围 / display-name From 拒绝。
  - 运维文档面：`config.default.yaml`（embed）、`configs/config.yaml`、`configs/env.example` 增加 mail.smtp 块与 MAIL_SMTP_* 说明。
- 测试实跑（go1.26.0 windows/amd64）：`go test ./internal/mail/ -v` 5 测试全 PASS；`go test ./internal/config/ ./internal/kernel/` 全绿；`go build ./...` OK；vet 干净；`go test ./internal/composition/` 全绿（20.3s，确认 config 消费方无回归）。

## 证据

| 主张 | 路径 |
|------|------|
| 决策正文 | 本目标 `01-decision/D-001-smtp-dial-and-config-keys.md` |
| Root 关闭 I-003/I-004 | `../../GOAL-001-outbound-mail/01-decision/D-003-r2-dial-and-config-freeze.md` |
| SMTP 适配器 | `apps/api/internal/mail/smtp.go` |
| 配置面 | `apps/api/internal/config/config.go`（validateMail）|

## 未做

- 未接 composition（适配器选择归 R3）；未动 `readyz`（R4）；未实现 capture sink（R3）；未加第二拨号路径。
