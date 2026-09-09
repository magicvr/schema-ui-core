---
doc_type: goal-decision
id: D-001-r1-denominator-freeze
parent: GOAL-002-r1-denominator-freeze
date: 2026-09-09
status: accepted
version: 0.1.0
---

# D-001 · R1 对照分母冻结

## 上下文

用户指令「先做 R1 分母冻结」。口径继承：

- VP-035 首波冻结表（内核/组合根/已交付端口/Profile/architecture 文档；不进 Admin 体验增强与 gated 实现）
- Root D-001：「现在修」**默认另立**
- I-035-002 四类业界参照集已 verified，本决策不重开

未选：把本评估做成 VP-010 波次；在 R1 开始改代码；把 Redis/MQ「没实现」预写成缺陷。

## 冻结

正文见 [r1-denominator-freeze.md](../attachments/r1-denominator-freeze.md)。摘要：

1. **I-035-001**：R2 必须覆盖模块契约、Profile、组合根、Store/Object/Cache/RateLimiter/EventBus/Mail/Obs/Shutdown/Job/Telegram 端口与内存（或已交付）供应商、Manifest 聚合、architecture 权威文档、handler/模块 import 泄漏抽检。排除 Admin 体验增强、业务域产品规则、009/010 波次史、Web 产品页（除 NavGroup 消费）、未实现 gated 供应商、VP-024 分发残余、VP-034 Dashboard 展示残余。
2. **I-035-004**：搬运器、otlp-sink、指标分母、JWT 立即失效、MFA wrapping、停机 harness、Redis/outbox gated、Telegram keyfile、RT-T03、RT-P04 **入册**供 R2/R3 分类。VP-024 分发四项、VP-034 Dashboard、VP-017 历史 SMTP **不入册**。
3. **I-035-005**：默认另立。本 VP 内只做 R4 文档卫生与只读对照断言。禁止在本 VP 改公开端口或实现 gated 供应商。

I-035-003 保持 collecting，最晚 R3。

## 对后续阶段的门禁

- 未按 §1 包含行取证，不得宣称 R2 矩阵完成。
- 把 §1.2 排除项写成「架构缺口」= 越界。
- 任何端口实现变更须新 VP 或 009/010 波次，并再走 `/vision`。
