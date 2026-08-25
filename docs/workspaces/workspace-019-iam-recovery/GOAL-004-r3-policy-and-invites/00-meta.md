---
id: GOAL-004-r3-policy-and-invites
title: R3 密码策略 + 邀请入职
status: active
parent: GOAL-001-iam-recovery
created: 2026-08-25
updated: 2026-08-25
version: 0.2.0
progress: 2/5
plan_refs:
  - VP-019-iam-recovery
primary_plan: VP-019-iam-recovery
serves_summary: 承接 Root R3：按 R1 合同（GOAL-002 D-001 §2/§3）实施密码策略（迁移单例策略行 + 历史 + 四口强制 + admin.settings 配置面）与邀请入职（迁移邀请表 + 管理面 + 公开激活 + 邮件/双形态）。不做 SMS、模板中心、存量扫描。
---

# GOAL-004 · R3 密码策略 + 邀请入职

## 概述

本目标承接 Root `GOAL-001-iam-recovery` 纲领阶段 **R3**。合同输入全部冻结（GOAL-002 D-001 §2/§3 + I-003/I-004/I-005/I-007 verified）；本轮新增一项用户裁决：**受邀账号初始角色 = 管理员发布邀请时指定**；同时用户指出 Web 新建用户表单应支持直接选角色（后端已支持，表单缺字段，列入本目标 Web 检查点）。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | R3 方案冻结决策落盘（本目标 D-001，含角色裁决） | **完成**：D-001 §1～§5（f65f9b01） |
| C2 | 迁移 0057/0058 落地（策略行/历史/邀请表；checksum 台账 + 黄金断言同步） | **完成**：v57/58/59 checksum 入台账（bfc7c4f2…/04f77fdc…/a35bbb21…）；identity head 56→59 + 三表清单；四处黄金断言同步；store 全绿 |
| C3 | 后端实施：策略域 ValidateNewPassword + 四口强制 + 配置 API；邀请域 + 管理 API + 公开激活 + 测试绿 | **进行中**：策略域 + 四口接线完成（E-002）；配置 API 与邀请域待做 |
| C4 | Web 面：设置页策略 tab、用户页邀请管理、公开激活页、新建用户表单补角色选择 + i18n | 待完成 |
| C5 | independent 审计开放 required = 0 + self 关门审 | 待完成 |

`progress` = 已完成检查点 / 5。当前 **0/5**。

## 边界

- 渐进生效：策略仅在四个设密时刻强制；不扫描存量、不强登出（I-007 verified）。
- 不做 SMS / 模板中心 / 多邮箱 / 组织权限；不改 Profile 默认集。
- 审计模式：实施切片 independent（grok build · grok-4.6 · high）+ 关门 self。

## 成功标准

1. 策略可配置且四口统一强制：minLength/categories/history 生效可核对；默认值保持现行行为。
2. 邀请全链可核对：指定角色建邀（邮件或链接）→ 激活建号带角色 → 一次性/撤销/过期/重发冷却语义正确。
3. independent 开放 required = 0；Root progress 推进至 3/4。
