---
id: GOAL-001-outbound-mail
doc: decision-entry
record_id: D-003
status: accepted
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-003 · R2 拨号路径与配置键冻结（关闭 I-003 / I-004）

### 触发

Root 路线图 R2 门禁：I-003（拨号路径）、I-004（配置键与凭证注入）须在 SMTP 接入实施前关闭。VP-017 将路径选择授权给 lead Root 方案冻结。

### 决定

采纳子目标 `GOAL-003-smtp-dial-config` D-001 的全部冻结项（该文件为正文）：

1. **I-003 → verified**：唯一拨号路径 = **隐式 TLS 465**（`tls.Dial` → `smtp.NewClient`；ServerName=host、MinVersion TLS1.2、证书校验强制、仅 PlainAuth over TLS）。
2. **I-004 → verified**：`mail.smtp.{host,port,username,password,from}`（env `MAIL_SMTP_*`）；全空合法（默认 sink），任一非空则 host/username/password/from 全必填，port 缺省 465 且须 1–65535，from 须 bare addr-spec；校验挂 `ValidateProd` 链；password 为 secret，只走 env / `configs/.env` 插值，不回显值。

### 为什么

- 隐式 TLS 自首字节起加密，无 STARTTLS 降级歧义，"可核对路径"最强形态（RFC 8314 方向）；stdlib 原生支持零新依赖。
- 配置面沿用 RT-K01 YAML+env fail-closed 三层机制，与 db/objects 键规则同构。

### 未选方案

见 `GOAL-003` D-001「未选方案」节（STARTTLS 587 / 双端口 / 无认证放行）。
