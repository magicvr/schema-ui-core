---
id: GOAL-004-r3-binding-flow
title: R3 绑定/校验消费 MailSender
status: done
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
progress: 4/4
plan_refs:
  - VP-018-account-email-identity
primary_plan: VP-018-account-email-identity
serves_summary: 承接 Root R3：迁移 0055 验证挑战表；绑定（占槽）/校验/过期重发流消费 kernel.MailSender（I-005 数值已冻结：TTL 10 分钟、冷却 60 秒）；I-006 允许管理员代填待校验；交付身份面 API 与账号页最小绑定 UI；independent 审计后关门。
---

# GOAL-004 · R3 绑定/校验消费 MailSender

## 概述

本目标承接 Root 纲领阶段 **R3**：按 R1 合同（GOAL-002 D-001 §2/§4/§5/§6）实现绑定与校验状态机。用户裁决已关闭 I-005 / I-006（2026-08-24 四项选定），分母确认**含最小 Admin 页面**。

对齐递归：GOAL-004 → Root GOAL-001（R3）→ VP-018 → Charter @0.2.0。不进入忘记密码状态机、邀请、密码策略、SMS、模板中心或业务域。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | 迁移 0055 挑战表落地（双方言成对 DDL）+ 黄金断言全套同步 | **完成**：v55 checksum `1556bda2…` 独立复算一致（A-001 核对）；head 55 全链 |
| C2 | 绑定/校验/重发流实现且合同映射可核对（占槽冲突 fail-closed、覆写换绑、三态迁移、常量时间比较、尝试上限） | **完成**：服务流 9 测试；F-002 同址重绑冷却补齐 |
| C3 | 身份面 API + 账号页最小绑定 UI 可用；测试绿（sqlite 全量 + PG 集成） | **完成**：三端点 + email-identity 卡片；store/authsession/handler/composition 全绿 + PG 集成 + web vitest/build；I-006 HTTP 链路修复后全通（bd1cdff9） |
| C4 | independent 审计（grok build · grok-4.6 · high）落盘且开放 required = 0 | **完成**：A-001 conditional → F-001 fixed（bd1cdff9），其余响应 E-003，开放 required 归零 |

`progress` = 已完成检查点 / 4。当前 **4/4**（已关门）。

## 边界

- 只消费 `kernel.MailSender` 端口；不改渠道模型、不锁历史 CaptureSink 叙事。
- 不做忘记密码/恢复语义；验证码不得当恢复证明。
- 承接清单兑现：F-001 PG 语义 harness（可选，若做则并入本目标证据）、F-002 配对不变量进仓储层、N-1 归一补偿在写入/比较路径落实。
- 审计模式：实施 independent（Root D-001 既定）。

## 成功标准

1. 绑定即占槽可核对：他号绑定同址（大小写折叠）被拒 `email_already_in_use`；多 NULL 共存不受影响。
2. 校验流可核对：正确码 → verified；错码计数、超限作废；过期/缺失拒绝；冷却期内重发拒绝。
3. 管理员代填 → pending；无任何路径绕过投递直达 verified。
4. 最小页面可用（绑定入口 + 码输入 + 状态徽标）；未引入模板中心/设置页邮件 tab 复制品。
