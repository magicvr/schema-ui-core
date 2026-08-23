---
id: GOAL-004-default-sink-surface-sweep
doc: audit-entry
record_id: A-001
source: self
goal: GOAL-004-default-sink-surface-sweep
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## A-001 · R3 默认 sink 落地与公共面 sweep 自审

- **source**: `self`
- **日期**: 2026-08-22
- **scope**: capture sink 合同一致性、composition 接线、公共面 sweep、边界遵守
- **verdict**: **pass**（无开放 required）

### 核对结果

| # | 核对面 | 结论 | 证据 |
|---|--------|------|------|
| 1 | R1 冻结兑现 | 容量 1 环 + `Last()` 取最后一封 + slog 双写，取报文 API 不进 kernel 端口 | `capture.go` / `capture_test.go` |
| 2 | 启动不变量 | 未配置 SMTP 时 fx 全图照常 Start/Stop（双 profile lifecycle 测试继续全绿）；容器解析出 capture sink 并可经端口 Send→Last 回读 | composition_mail_test + lifecycle 测试 |
| 3 | 选择规则 | `MailSMTPConfigured()` 分叉；ValidateProd 把关 + 构造器防御性校验双保险（部分块拒收测试留证） | composition_mail_test.go |
| 4 | 公共面 sweep | SMTP 客户端类型零泄漏：`net/smtp`/`tls.Dialer` 仅 internal/mail；handler/modules 对 mail 类型零引用（rg 留痕 E-001） | grep 输出 |
| 5 | 行为面回归 | mail/config/kernel/handler/composition 五包测试全绿；vet 干净 | 实跑记录 |

### Findings

无。实现与 D-001 决策逐条对应；R1/R2 冻结件未被改动。

### 结论

R3 门禁达成且证据可回指；无开放 required finding。C1/C2/C3 完成，本目标可关门（done 3/3）。R4 将做显式路径证据 + `readyz` 扩依赖，并处理 I-005（HTML/MIME 关门叙事）复核。
