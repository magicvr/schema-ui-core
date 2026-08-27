---
id: E-002
doc: execution-entry
goal: GOAL-006-channel-provider-contract
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-002 · R5 渠道合同冻结执行（2026-08-24）

## 已发生事实

1. 四个合同决策点经用户书面裁决（会话问答留痕，2026-08-24）：mock 持久化 = **DB 表 + 迁移**；解析规则 = **显式键 `mail.channel`**；保留策略 = **有界默认 500 条、管理员可经 mock 渠道配置调整**；**I-010 一并预冻**。
2. 已落盘 D-002（可实施条款 §1～§4），关闭 I-011、预冻 I-010。
3. 对照代码现状核对了条款落点：`internal/composition/composition.go` `newMailSender`（解析唯一实现点）、`internal/config` `mail.smtp.*` / `validateMail` fail-closed 先例、`kernel.Store` 双方言 + 编译迁移目录（现目录头 0050）、模块路由注册先例。条款与现有机制同构，无结构性冲突。
4. 未修改应用代码；未创建迁移；未实施 Resend/mock 产品面（R6 边界不变）。

## 证据

| 主张 | 路径 |
|------|------|
| 合同冻结节 | [D-002-r5-channel-contract-freeze.md](../01-decision/D-002-r5-channel-contract-freeze.md) |
| 解析现行实现 | `apps/api/internal/composition/composition.go`（`newMailSender`） |
| 配置 fail-closed 先例 | `apps/api/internal/config/config.go`（`MailSMTPConfigured` / `validateMail`） |
| 用户裁决留痕 | 本条「已发生事实」第 1 点（会话问答记录） |

## 未做

- 未改 `apps/api` / `apps/web`；未动 VP-018 冻结；未创建 R6 子目标（P-001 按阶段，R5 关门后另开）。
