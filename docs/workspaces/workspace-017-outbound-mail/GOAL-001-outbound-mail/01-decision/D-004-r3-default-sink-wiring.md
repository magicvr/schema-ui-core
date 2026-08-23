---
id: GOAL-001-outbound-mail
doc: decision-entry
record_id: D-004
status: accepted
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-004 · R3 默认 sink 接线与公共面 sweep 规则

### 触发

Root 路线图 R3 门禁：默认 sink 落地（未配置 SMTP 仍能启动；测试可取出最后一封）+ 公共面去 SMTP 客户端类型。

### 决定

采纳子目标 `GOAL-004-default-sink-surface-sweep` D-001：

1. `CaptureSink` 按 R1 冻结形态落地（容量 1 环 / Last / Reset / slog 结构化日志），取报文 API 不进 kernel 端口。
2. composition 经 `newMailSender` 向 fx 容器提供唯一 `kernel.MailSender`：`MailSMTPConfigured()` 分叉 capture/SMTP，构造器防御性二次校验。
3. 本波端口零消费者（handler/模块不注册邮件面）；公共面 sweep 判据 = SMTP 客户端类型仅限 `internal/mail`，grep 证据留痕子目标 E-001。

### 为什么

单例 + 每连接独立拨号天然并发安全；双保险校验防手拼 Config 绕过 ValidateProd。

### 未选方案

见 `GOAL-004` D-001「未选方案」节（Last 进端口 / 全局单例注入）。
