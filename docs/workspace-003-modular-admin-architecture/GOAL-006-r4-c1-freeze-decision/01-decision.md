---
id: GOAL-006-r4-c1-freeze-decision
doc: decision
status: active
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# 决策记录 · GOAL-006

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| C1-I001 | required | Provider/Registrar 精确契约、Contribution 类型和 compiled-global Persistence 是否接受 | C1 close / C2 entry | C1.2 | 用户书面裁决、D-003、最终审计 | collecting | 无延期；未裁决不得进入 C2 | 父目标 freeze package、A-005 |
| C1-I002 | required | Records historical-only 或恢复产品 CRUD | C1 close / C4 scope | C1.2 | 用户书面选择并更新 canonical 范围 | collecting | 信息冲突不得静默推断 | 父目标 R4-I003、`0006 records_retire` |
| C1-I003 | required | operationlog Option A/B/C 及 A 的 residual 接受条件 | C1 close / C3/C5 data gate | C1.2 | 用户书面选择、failure/retention evidence | collecting | owner/review trigger/date 未确认 | 父目标 R4-I004、A-005 FP-003 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-05 | 建立 R4-C1 冻结裁决子目标与路线图 | accepted | [01-decision/D-001-r4-c1-stage-scope.md](01-decision/D-001-r4-c1-stage-scope.md) |
| D-002 | 2026-08-05 | R4-C1 候选冻结包继承与待裁决轴 | proposed | [01-decision/D-002-r4-c1-freeze-candidate.md](01-decision/D-002-r4-c1-freeze-candidate.md) |

## 当前约束

- 父目标 A-005 是独立候选复审证据，不是本目标的用户裁决或 D-003。
- Provider 的工程建议、Records historical-only 建议和 operationlog Option A 建议均
  保持 `pending_user`，不得以 continuation 或默认偏好代替书面裁决。
- 任何 required finding 只能以 `fixed`、`accepted-residual` 或 `user-overruled`
合法关闭；本目标当前没有已接受 residual。
