---
id: GOAL-001-account-email-identity
doc: decision
status: active
parent: null
created: 2026-08-24
updated: 2026-08-24
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

> 状态以 `00-meta.md` 信息表为准（本表为镜像，须保持同号同状态）。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 非空邮箱唯一性细则 | R1 方案 | R1 冻结 | R1 决策 | collecting | — | I-018-001 |
| I-002 | required | 校验投递形态：验证码 vs 链接 | R1 方案 | R1 冻结 | R1 决策 | collecting | — | I-018-002 |
| I-003 | required | 账号可否长期无邮箱 | R1 方案 | R1 | 投影 VP | **registered**（D-001） | — | VP 冻结可空 |
| I-004 | required | 换绑是否进本波 | R1 方案 | R1 | 投影 VP | **registered**（D-001） | — | VP 冻结换绑进分母 |
| I-005 | required | 过期与重发冷却 | R3 方案 | R3 接入前 | R3 决策 | collecting | — | I-018-005 |
| I-006 | non-blocking | 管理员代填待校验 | R3 方案 | R3 | R3 或不进分母 | collecting | — | I-018-006 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-24 | 开区 scaffold 与 Admin 邮箱身份纲领路线图 | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
