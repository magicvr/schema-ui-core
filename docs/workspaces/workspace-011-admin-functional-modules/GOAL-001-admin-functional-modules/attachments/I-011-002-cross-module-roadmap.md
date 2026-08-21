---
id: I-011-002
goal: GOAL-001-admin-functional-modules
title: 跨模块能力地图与路线图登记（四档分层 · 非立即实施）
date: 2026-08-18
status: registered
parent: GOAL-001-admin-functional-modules
created: 2026-08-18
updated: 2026-08-18
version: 1.1.0
---

# I-011-002 · 跨模块能力地图与路线图登记

> 本附件登记 2026-08-18 审视意见形成的**四档能力地图与推进顺序**。性质为**路线图/backlog 登记**，不是立即实施清单；具体立项须按触发条件、P-001 与既有波次纪律单独创建子目标。
> 与 `I-011-001-tiered-inventory.md` 的关系：`I-011-001` 是 R1 调研产出的“可见功能页”三档分档；本附件补充 R1 清单遗漏的三类跨模块能力，避免后续重复分析。
> **2026-08-18 更新**：其中“横切契约”部分已由新 VP-012 / workspace-012 承接；本附件继续作为四档能力地图与 Tier D 业务域的登记源。

## 1. 为什么需要四档而不是继续扩 S/B 编号

- 当前 `I-011-001` 已覆盖大量“可见功能页”：用户/角色/RBAC、个人中心、MFA、操作日志、设置、Dashboard、导入导出、文件库、数据字典、系统监控、定时任务、验证码、回收站、数据权限、钱包等。
- 真正遗漏的是三类**跨模块能力**：
  1. 横切基架能力：身份、安全、审计、并发、运维、配置迁移；
  2. 未来扩展必须提前稳定的接缝：事件、通知通道、异步任务、认证提供商、能力/租户上下文；
  3. 业务领域中尚未列出的实体和流程：支付、退款、价格、计费、发票、退货、评论、内容发布等。
- 若不分开，容易把通用 Admin 做成电商/交易单体，或过早建设过度抽象的框架。

## 2. 四档能力地图

### Tier A：Admin 基架规划（通用能力，按触发/门禁实现）

1. 安全审计增强：失败认证、结果、IP/UA、correlation、结构化 diff；
2. 管理员全局会话治理；
3. 密码策略、邀请和账号恢复状态机；
4. API Token / Service Credential；
5. 组织/部门/岗位，以及数据权限 `org` 扩展；
6. 乐观并发、冲突和幂等契约；
7. 异步 Job / 长操作契约；
8. 配置包导出、diff、dry-run、导入；
9. correlation/metrics/诊断；
10. maintenance/degraded/read-only 模式；
11. 文件扫描/隔离接缝；
12. 时间、时区、数字、货币格式语义。

### Tier B：扩展接缝（预留接口，出现触发条件再实现）

1. typed domain event / outbox；
2. Notification Transport Provider；
3. OIDC/SSO/SCIM provider 接缝；
4. Approval Gate；
5. Entitlement/feature availability；
6. 多组织/workspace context；
7. 实时 SSE/WebSocket 和缓存失效；
8. 外部连接器/Secret Provider；
9. 自定义 metadata/tags；
10. 文件预览、分片上传和恶意内容扫描。

### Tier C：Admin 后续体验增强

1. 全局搜索/Command Palette；
2. Saved Views、筛选器书签、列偏好；
3. 批量操作的异步进度与结果中心；
4. 未保存变更保护；
5. 通知/Toast/错误恢复的统一体验；
6. 版本更新、维护提示和诊断报告。

### Tier D：真实业务成立后再立项

1. Catalog、商品、SKU、价格、税；
2. 库存、仓库、预留、调拨；
3. 订单、支付、退款、退货；
4. 物流、包裹、履约；
5. 营销、优惠券、促销规则；
6. 订阅、计费、发票、用量；
7. 工单、客服、CRM；
8. CMS、内容发布、知识库；
9. 支付网关、ERP、物流、CRM 连接器。

## 3. 最值得马上补入路线图的八项

1. 管理员全局会话治理；
2. 密码策略/邀请/账号恢复；
3. 安全审计增强与结构化变更 Diff；
4. correlation、metrics、maintenance/degraded 模式；
5. 乐观并发、幂等和异步 Job 契约；
6. 配置包导出、环境 Diff 和 dry-run；
7. 领域事件 / 通知 Transport / SSO 的扩展接缝；
8. 业务领域中缺失的价格、税、支付、退款、退货、发票和结算。

> 其中前五项明显属于通用基架；后三项分别是平台接缝、运维工具和业务域能力，不能用统一“Admin 功能模块”标签处理。

## 4. 推进顺序（登记用，非立即执行）

```text
0. Profile 边界收敛：通用 Admin 与 Commerce / SaaS / Support 领域包分开

1. 可靠性与安全横切契约
   correlation / 错误恢复 / 乐观并发 / 幂等 / 审计事件模型

2. IAM 运维增强
   安全审计增强 / 全局会话治理 / 密码策略·邀请·恢复 / API Token

3. 企业组织与平台管理
   组织·部门·岗位 / 数据权限 org 扩展 / 公告 / 消息模板·通知通道 / Webhook

4. 运维与多环境
   maintenance·degraded·read-only / metrics·diagnostics / 配置包迁移 / 审计 retention·归档 / 异步 Job

5. Admin 生产力体验
   全局搜索 / Saved Views / 批量任务结果 / 时区·数字·货币 / 一致错误恢复

6. 真实业务领域
   Commerce / SaaS / Support / CMS 等按触发条件独立立项
```

## 5. 处置规则

- **不**继续把 Tier B/C/D 全部堆进 workspace-011 的 S/B 模块编号。
- **不**为上述计划批量创建子目标；仅在触发条件成立、P-001 路线图就位后按阶段立项。
- 共享基架/横切问题默认分流到 **VP-009（生产加固）/ VP-010（设计—实现符合性）** 或未来平台 VP，不默认塞进 VP-011。
- 具体业务领域（Tier D）保持领域问题留领域台账；Charter 非目标（多租户、白标等）不变。
