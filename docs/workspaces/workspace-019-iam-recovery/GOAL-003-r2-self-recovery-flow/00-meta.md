---
id: GOAL-003-r2-self-recovery-flow
title: R2 自助恢复全链（后端 + Web）
status: active
parent: GOAL-001-iam-recovery
created: 2026-08-25
updated: 2026-08-25
version: 0.1.0
progress: 0/5
plan_refs:
  - VP-019-iam-recovery
primary_plan: VP-019-iam-recovery
serves_summary: 承接 Root R2：按 R1 冻结合同（Root D-002 + GOAL-002 D-001 §1/§4/§5）实施自助恢复全链——迁移 0056 挑战表（双方言）、core.auth-session 恢复域逻辑、公开 API（start/complete）、MFA 第二因子门、Web 登录页恢复流、经 mock 渠道取码的端到端证据。不做密码策略产品面（R3）与邀请（R3）。
---

# GOAL-003 · R2 自助恢复全链（后端 + Web）

## 概述

本目标承接 Root `GOAL-001-iam-recovery` 纲领阶段 **R2**：把「忘记密码」从空转收成可核对的闭环。合同输入全部冻结：

- **证明形态**：6 位数字验证码（VP-018 同构；sha256 落库 + 恒时比较）；投递只经 `kernel.MailSender` 现行渠道。
- **时效**：TTL 10 分钟 / 重发冷却 60 秒 / 连续错 5 次作废挑战。
- **MFA 门**：已登记 active TOTP 的账号在邮箱码通过后、设新密码前须再过一次第二因子（TOTP 或恢复码）；无法提供者走管理员重置。
- **无邮箱边界**：未绑定/未校验邮箱的账号受控负响应、不发信，走管理员重置。
- **会话**：完成设密走现行 UpdateUser 语义（token_version+1 → 其余会话失效；refresh 全撤销；must_change_password 清除）。

对齐递归：GOAL-003 → Root GOAL-001（R2）→ VP-019 → Charter @0.2.0。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | R2 方案冻结决策落盘（本目标 D-001）：API 形状、错误码、MFA 门接口、会话语义 | 待完成 |
| C2 | 迁移 0056 双方言落地（checksum 台账 + 黄金断言同步 + 升级/全新库测试） | 待完成 |
| C3 | 域逻辑 + HTTP 面实施（start/complete、MFA 门、限流、审计/通知）+ 包测试绿 | 待完成 |
| C4 | 端到端证据（经 mock 渠道取码全链测试）+ Web 登录页恢复流 + i18n | 待完成 |
| C5 | independent 审计（grok build · grok-4.6 · high）意见落盘且开放 required = 0；self 关门审 | 待完成 |

`progress` = 已完成检查点 / 5。当前 **0/5**。

## 边界

- 密码策略产品面（配置 tab / 策略引擎参数化）归 R3；本阶段设密点沿用现行基线校验（8–72 字节非空白）。
- 邀请入职归 R3；证据关门归 R4。
- 不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；不动既有迁移 checksum。
- 安全滥用面（枚举/重放/open-relay/邀请滥用）归 VP-009/VP-010；本面只做 fail-closed 受控负响应。

## 成功标准

1. 全新库与升级路径干净应用 0056；`go build ./...` 与相关包测试绿。
2. 全链可核对：start 经 mock 渠道出站取码 → complete 设新密码 → 旧会话失效 → 新密码可登录；MFA 账号被第二因子门拦住直至提供有效因子。
3. independent 审计开放 required = 0；Root progress 推进至 2/4。
