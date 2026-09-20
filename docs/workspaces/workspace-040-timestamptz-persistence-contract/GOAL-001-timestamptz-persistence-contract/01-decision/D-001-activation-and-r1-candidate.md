---
id: D-001-activation-and-r1-candidate
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-001 · 激活边界与 R1 默认候选登记

## 用户指令与范围

用户 2026-09-20 指令为「/vision 走流程激活 vp-040，没有问题的话，交 /govern 开设工作区」。本条记录该指令在实现层的交接边界：VP-040 激活，随后 scaffold delivery 工作区与 Root；不把激活误写成实现或关门。

## 本轮决定

1. VP-040 `planned → active` v0.2.0；绑定 `workspace-040-timestamptz-persistence-contract`，Root 为 `GOAL-001-timestamptz-persistence-contract`，`vision_role: delivery`。
2. 为满足激活门禁，登记一个 **R1 默认候选**：SQLite 使用 `INTEGER`、PostgreSQL 使用 `BIGINT`，两者承载 UTC Unix seconds 的同一绝对时刻语义。
3. 该候选不是 R1 最终冻结。现有代码存在秒与毫秒并存的事实；所有毫秒字段、精度、NULL/零值、编解码、备份与迁移分母必须在 R1 再核对并由目标台账冻结。

## 取舍与未选方案

- `TEXT`/RFC3339 可作为 R1 的对照方案，但本轮不把它写成已选合同；其索引、排序、精度与存量转换影响待 R1 评估。
- 直接把 PostgreSQL 改成 `timestamptz`、SQLite 继续原样 `INTEGER` 不可接受，因为会破坏 VP-013 的合同平等要求。
- 不把所有 `INTEGER` 归为时间列；金额、flag、version、bot_id 等必须显式排除。

## 证据与下一门禁

- 激活资格、VP-039 前置、架构 freshness 与 slug： [VRev-104](../../../../vision/reviews/VRev-104-vp040-timestamptz-persistence-contract-activation.md)。
- 组合/方向权威： [VP-040](../../../../vision/plans/VP-040-timestamptz-persistence-contract.md)。
- 下一步唯一主动作：先在 R1 冻结 `I-040-001`～`I-040-003` 与列分母，再进入任何不可逆 schema/迁移实施。

本决定不接受 residual、不关闭 required 信息项，也不替代 Goal `03-audit`。
