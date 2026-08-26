---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-014-w13-account-lockout-redesign
version: 0.1.0
---

# A-003 · 编排器对 A-002 的响应（S6 闭合记录）

- **source**: self（编排器响应，非独立意见；独立原文见 [A-002](A-002-independent.md)）
- **日期**: 2026-08-26
- **响应范围**: A-002 verdict pass；recommended R-F001 / R-F002 / R-F003 全部响应

## 逐项响应

### R-F001（台账未跟上 HEAD 并发兜底）→ **fixed（部分已预先兑现）**

- E-002 已于提交 `3391f0fb` 补记并发首插兜底与真实 Postgres 方言复核（审计读取快照早于该笔，本条确认台账现态已含）。
- 并发首插回归测试：**经裁不补**——该竞态在 sqlite/PG 上无法确定性地强制触发（需精确交错 UPDATE-miss 与 INSERT），钱包 `TestGetOrCreateUserAccountConcurrent` 先例覆盖的是同形 get-or-create 而非此三行哨兵回退路径；逻辑本身为 sentinel 分支 + 新事务重跑同函数，风险与复杂度不成比例。残余风险登记于本条留痕。

### R-F002（全局锁窗内 Refresh 先轮换后校验）→ **accepted-residual**

接受为 D-002 范围内既有契约残余：攻击成本已 ×20（100 次失败 + 触发全局管理员通知），且仅"锁窗内主动出示的那一张令牌"受影响；密码修改/token_version/管理员禁用等真实入侵响应路径未动。**复审触发**：若产品提出"全局锁窗内出示也不丢会话"需求，另开决策改 Refresh 校验顺序（超出本次冻结范围）。范围：仅全局锁窗 ≤15min；来源锁不受影响。

### R-F003（注释语义 + 负向断言缺口）→ **fixed**

- `TestAccountSourceLockKeepsSessions` 注释已改为准确语义（来源锁不写 users.locked_until、Refresh 不受影响）——`12b5a7e7`。
- 补负向/正向断言：`TestLoginSourceScopedLockout` 增加"OnLockOpened 对来源锁恒 0"与"来源锁下既存刷新令牌仍可轮换（会话不波及）"。全部通过。

### 附带

- 审计指出"E-002/A-001 写 `26655b55` 未跟上 HEAD"：E-002 已由 `3391f0fb` 补记两笔后续提交（PG 复核 + `cf5675f1` 兜底）；A-001 按台账追加原则不改写，以本条为准。

## 台账状态（本条目后）

| 来源 | 条目 | verdict | 开放 required |
|------|------|---------|----------------|
| self | A-001 | pass | 0 |
| independent | A-002 | pass | 0 |
| self | A-003（本条） | —（响应记录） | 0 |

**开放 required findings = 0。** S5 完成；S6 剩余动作为用户书面关门确认（按 GOAL-013 D-003 两目标一并处理）。
