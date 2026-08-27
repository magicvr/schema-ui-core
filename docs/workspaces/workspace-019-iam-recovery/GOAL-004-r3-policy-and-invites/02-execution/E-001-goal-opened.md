---
id: E-001
doc: execution-entry
goal: GOAL-004-r3-policy-and-invites
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# E-001 · R3 开题 + D-001 方案冻结（C1 满）

2026-08-25 完成：

- 用户裁决：受邀账号初始角色 = 管理员发布邀请时指定；另指出 Web 新建用户表单应直接支持角色选择（后端已支持、表单缺字段 → 列入 C4）。
- D-001 落盘：0057 策略行+历史 / 0058 邀请表；策略域四口强制 + admin.settings 配置 API；邀请管理/激活全链（角色随邀请、INVITE_ROLE_GONE fail-closed、不签发会话）。
- 未改动任何产品代码。

后续：C2 迁移 → C3 后端 → C4 Web → C5 审计关门。
