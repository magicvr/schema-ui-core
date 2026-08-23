---
id: GOAL-004-default-sink-surface-sweep
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-001 · R3 默认 sink 接线与公共面 sweep 规则

### 触发

Root R3 门禁：默认 sink 落地（未配置 SMTP 仍能启动；测试可取出最后一封）+ 公共面去 SMTP 客户端类型。R1/R2 已冻结端口与适配器合同。

### 决定

1. **CaptureSink 形态**（按 R1 冻结实现）：容量 1 内存环 + `Last() (kernel.MailMessage, bool)` + `Reset()`；每封接受的消息打一行 slog 结构化日志；`Send` 先按端口合同校验再入环。取报文 API 留在具体适配器上，不进 kernel 端口。
2. **composition 选择规则**：`newMailSender(cfg, logger)` 单例提供进 fx 容器——`cfg.MailSMTPConfigured()` 为真 → `*mail.SMTP`；否则 → `*mail.CaptureSink`。完整性已由 `config.ValidateProd` fail-closed 把关，构造器再做一次防御性校验（双保险，测试留证）。
3. **本波无消费者**：handler 与模块 Provider 不注册邮件路由/能力；端口消费归后续 Admin VP。容器解析即接线证据。
4. **公共面 sweep 判据**：`net/smtp`、`tls.Dialer`、`smtp.PlainAuth` 仅允许出现在 `internal/mail/*`；handler / modules 包对 mail 类型零引用。以 grep 全仓证据留痕于 E-001。

### 为什么

- 单例 + 每连接独立拨号：SMTP 适配器天然并发安全，capture 有锁；无需连接池（同步发送、无队列）。
- 防御性二次校验：防手拼 Config 绕过 ValidateProd 的路径静默落到半配置端点。

### 未选方案

- **把 Last()/Reset() 提升进 kernel.MailSender**：污染生产端口，测试关注点泄漏进合同。未选。
- **fx 以外全局单例 / init 注入**：破坏既有显式装配风格。未选。
