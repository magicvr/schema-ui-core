---
id: GOAL-001-iam-recovery
doc: decision
status: active
parent: null
created: 2026-08-25
updated: 2026-08-25
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

> 状态以 `00-meta.md` 信息表为准（本表为镜像，须保持同号同状态）。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 自助恢复证明形态（默认候选 = VP-018 6 位码） | 方案冻结 | R1 | 用户裁决 | collecting | — | 待确认 |
| I-002 | required | 恢复令牌/验证码 TTL 与重发冷却 | 方案冻结 | R2 方案冻结前 | 用户裁决 | collecting | — | 待确认 |
| I-003 | required | 密码策略默认参数与配置边界 | R2 方案冻结 | R2 | 用户裁决 | collecting | — | 待确认 |
| I-004 | required | 邀请形态 / 预置账号 | R3 方案冻结 | R3 | 用户裁决 | collecting | — | 待确认 |
| I-005 | required | 邀请有效期 / 撤销 / 一次性 | R3 方案 / 实施 | R3 接入前 | 用户裁决 | collecting | — | 待确认 |
| I-006 | required | 无邮箱账号边界：仅管理员重置 | 方案冻结 | R1 | 产品事实投影 | **registered**（2026-08-22） | — | 无邮箱不自助 |
| I-007 | required | 策略对既有账号生效边界 | R2 方案冻结 | R2 | 用户裁决 | collecting | — | 待确认 |
| I-008 | non-blocking | 改密后会话语义 | 方案冻结 | R4 | 用户裁决 | collecting | — | 待确认 |
| I-009 | required | MFA 账号如何走自助恢复（防旁路） | 方案冻结 | R1 合同冻结 | 用户裁决 | collecting | — | 待确认 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-25 | 开区 scaffold 与 IAM 纲领路线图（含 Admin 类 freshness） | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |