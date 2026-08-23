---
id: GOAL-001-outbound-mail
doc: execution-entry
record_id: E-003
goal: GOAL-001-outbound-mail
status: recorded
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-003 · R2 执行：SMTP 适配器与配置面落地（2026-08-22）

## 已发生事实

- 开设子目标 `GOAL-003-smtp-dial-config`（五件套齐全，parent = 本 Root），承接 R2 治理上下文。
- 子目标 D-001 冻结拨号路径（隐式 TLS 465 唯一路径）与 `mail.smtp.*` 配置键；Root D-003 关闭 I-003/I-004（均 → verified）。
- 落地代码：`apps/api/internal/mail/smtp.go`（隐式 TLS 适配器，证书校验恒开，AUTH PLAIN over TLS，Subject 注入守卫）；`internal/config` 邮件配置面（YAML + env 三层、`validateMail()` 挂 ValidateProd、部分块 fail-closed 报键名不回显值）。
- 运维文档：`config.default.yaml` / `configs/config.yaml` / `configs/env.example` 增加 mail.smtp 与 MAIL_SMTP_* 说明。
- 实跑证据：mail/config/kernel 包测试全 PASS；`go build ./...` OK；vet 干净；composition 全绿（20.3s 无回归）。
- Root 路线图 R2 → 已完成；progress 2/4。composition / handler / readyz 行为面零改动。

## 证据

| 主张 | 路径 |
|------|------|
| 子目标决策正文 | `../GOAL-003-smtp-dial-config/01-decision/D-001-smtp-dial-and-config-keys.md` |
| 本 Root 决策 | [D-003-r2-dial-and-config-freeze.md](../01-decision/D-003-r2-dial-and-config-freeze.md) |
| 子目标执行条目 | `../GOAL-003-smtp-dial-config/02-execution/E-001-smtp-adapter-and-config.md` |

## 未做

- 未接 composition（默认 sink 接线归 R3）；未动 readyz；未实现 capture sink；未做 live 投递 harness（R4）。
