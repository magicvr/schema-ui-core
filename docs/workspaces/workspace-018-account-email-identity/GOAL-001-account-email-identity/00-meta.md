---
id: GOAL-001-account-email-identity
title: 账号邮箱身份（绑定与校验）
status: active
parent: null
created: 2026-08-24
updated: 2026-08-24
version: 0.1.0
progress: 0/4
plan_refs:
  - VP-018-account-email-identity
primary_plan: VP-018-account-email-identity
serves_summary: 交付 Admin 功能账号邮箱身份面：users email 可空、绑定/校验状态机、换绑；校验信消费已 closed 的 VP-017 MailSender。不承载自助恢复、邀请、密码策略、SMS 或业务域。
---

# GOAL-001 · 账号邮箱身份（绑定与校验）

## 概述

本 Root 承载 [VP-018-account-email-identity](../../../vision/plans/VP-018-account-email-identity.md)（**`active`**）的实现：在已有 `core.auth-session` 与内核 `MailSender` 之上，补齐账号邮箱字段、绑定与校验状态机、唯一性与换绑。无邮箱账号必须继续能登录；无 SMTP 时用 capture sink 测通校验信。

**边界**：不承接登录页「忘记密码」状态机、邀请入职、密码策略产品化、管理员重置冒充自助恢复、SMS、消息模板或业务域。安全 finding → VP-009；符合性 gap → VP-010。

## 纲领路线图（P-001）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **身份合同冻结**：非空邮箱唯一性细则（I-001 / I-018-001）；校验投递形态验证码 vs 链接（I-002 / I-018-002）；可空与换绑已由 VP 冻结（I-003 / I-004）。 | 起点 | 未开始 |
| R2 | **双方言 schema + 唯一性**：`users` 加列与约束；SQLite 与 PostgreSQL 同一逻辑 schema、成对物理 SQL、checksum 台账。 | 依赖 R1 | 未开始 |
| R3 | **绑定/校验消费 `MailSender`**：发起绑定、投递、完成校验、过期/重发（I-005）；缺省 capture；显式 SMTP 可实投。管理员代填（I-006）若纳入则在本阶段。 | 依赖 R2 | 未开始 |
| R4 | **证据**：capture `Last()` 可取出校验信；唯一性 fail-closed 可核对；无 IAM 恢复 / 邀请 / 密码策略产品进入本波。 | 依赖 R3 | 未开始 |

`progress` = 已完成阶段数 / 4。当前 **0/4**。

## 成功标准（方向级）

1. `users` 可持有可空邮箱与可核对的校验状态；无邮箱账号仍能登录。
2. 绑定与校验流落地；校验信经 `kernel.MailSender` 发送；未配置 SMTP 时 capture sink 可取出最后一封。
3. 非空邮箱唯一性 fail-closed 可核对；换绑走同一校验合同。
4. 未引入忘记密码状态机、邀请、密码策略产品、SMS、第二运输、模板中心；未改 Charter；未改 Profile 默认集作为本波成功条件。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 非空邮箱唯一性细则：未校验地址是否占用唯一槽；大小写/规范化 | R1 方案冻结 | R1 合同冻结 | R1 决策 | collecting | — | 对应 I-018-001；未验证不得改 DDL |
| I-002 | required | 校验投递形态：验证码 vs 魔术链接 | R1 方案冻结 | R1 合同冻结 | R1 决策 | collecting | — | 对应 I-018-002；用户可在 R1 前书面选定 |
| I-003 | required | 账号可否长期无邮箱 | R1 方案冻结 | R1 | 投影 VP 冻结 | **registered**（D-001） | — | VP 已冻结可空；兼容现有 `users` |
| I-004 | required | 换绑是否进本波 | R1 方案冻结 | R1 | 投影 VP 冻结 | **registered**（D-001） | — | VP 已冻结换绑进分母 |
| I-005 | required | 校验令牌/验证码过期与重发冷却 | R3 方案 / 实施 | R3 接入前 | R3 决策 | collecting | — | 对应 I-018-005 |
| I-006 | non-blocking | 管理员可否代填邮箱并保持待校验 | R3 方案 | R3 | R3 或不进分母 | collecting | — | 对应 I-018-006 |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-018）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。
