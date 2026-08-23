---
id: GOAL-001-outbound-mail
doc: execution-entry
record_id: E-004
goal: GOAL-001-outbound-mail
status: recorded
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-004 · R3 执行：默认 sink 落地 + composition 接线 + sweep（2026-08-22）

## 已发生事实

- 开设子目标 `GOAL-004-default-sink-surface-sweep`（五件套齐全，parent = 本 Root），承接 R3 治理上下文。
- 落地代码：`internal/mail/capture.go`（capture sink，R1 冻结形态）+ composition `newMailSender` 接线进 fx 容器（缺省 capture / 显式 SMTP 分叉）+ `config.MailSMTPConfigured()`。
- 公共面 sweep 完成：`net/smtp`/`tls.Dialer` 仅存在于 `internal/mail/smtp.go`；handler/modules 对邮件类型零引用（rg 证据留痕子目标 E-001）。
- 实跑证据：mail/config/kernel/handler/composition 五包测试全绿；vet 干净；未配置 SMTP 的双 profile lifecycle 启动/停止测试继续通过（启动不变量成立）；fx 容器解析断言 + Send→Last 回读测试通过。
- Root 路线图 R3 → 已完成；progress 3/4。

## 证据

| 主张 | 路径 |
|------|------|
| 子目标决策正文 | `../GOAL-004-default-sink-surface-sweep/01-decision/D-001-default-sink-wiring.md` |
| 本 Root 决策 | [D-004-r3-default-sink-wiring.md](../01-decision/D-004-r3-default-sink-wiring.md) |
| 子目标执行条目（含 sweep 证据） | `../GOAL-004-default-sink-surface-sweep/02-execution/E-001-capture-sink-and-wiring.md` |

## 未做

- 未动 `readyz`（R4）；未做显式投递 live harness（R4）；I-005（HTML/MIME 关门叙事）留待 R4 复核。
