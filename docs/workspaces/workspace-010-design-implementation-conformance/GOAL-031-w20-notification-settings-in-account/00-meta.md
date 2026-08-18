---
id: GOAL-031-w20-notification-settings-in-account
title: W20 · 通知设置迁入个人中心
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-18
updated: 2026-08-18
version: 0.2.0
progress: 4/4
---

# GOAL-031 · W20 · 通知设置迁入个人中心

VP-010 / workspace-010 的**第二十波**：把「启用站内通知」从通知列表页挪到个人中心。列表页只做收件箱。API 路径仍是 `/api/notifications/settings`。

## 当前边界

- **范围**：`notifications.json` 去掉设置表单；`account.json` 增加「通知」Tab（switch + 回读/保存）；列表页能力声明去掉仅设置表单需要的 `form.record.load`。
- **非范围**：改开关语义、通道/邮件推送、公告；不改 Profile / 模块矩阵 / Manifest 装配。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：D-001。
- [x] **S2 · 实施**：schema + i18n + 契约测试（E-002）。
- [x] **S3 · 定向验证**：dval/schema-keys + Go 通知套件 + `tsc`（E-003）。
- [x] **S4 · 自审与关门**：A-001 self pass；goal-tree / workspace 同步（E-004）。

progress: 四个等权检查点；当前 **4/4**。

## 审计策略

S4 关门 `self`（可逆 IA；无安全门禁语义变化）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required | 开关放个人中心哪一档 | S1 | S1 | 对照 account 既有 Tabs | **verified** | D-001：新 Tab「通知」，在资料之后、安全之前 |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 溯源

- 用户 2026-08-18：列表页不应承载通知设置；应在个人中心
- [workspace-011] GOAL-006 D-002 §5：列表页冻结为收件箱，未规定设置表单上列表页
