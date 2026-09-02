---
id: GOAL-001-wallet-prepaid-instrument
doc: decision
status: active
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# 决策记录 · GOAL-001-wallet-prepaid-instrument

## 信息需求与阶段门禁

> 权威表在 `00-meta.md`。R1 合同冻结前，I-029-001 / I-029-002 / I-029-003 / I-029-006 为开放 required，阻断方案冻结与受影响实施。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-029-001 | required | 主体落点 + `owner_type` + `OwnerExists` 迁出 `UserByID` | 方案冻结 + 判据 1 | R1 | 用户裁决 | closed | — | D-002: 独立主体表 subjects，owner_type 扩充 subject，OwnerExists 校验主体 |
| I-029-002 | required | 核销 `entry_type` | 判据 3/4 | R1 | 用户裁决 | closed | — | D-002: 复用 adjust + ref_type='voucher' |
| I-029-003 | required | 生成权限键 | 判据 5 | R1 | 用户裁决 | closed | — | D-002: 新增细粒度 wallet.voucher.issue 权限键 |
| I-029-004 | non-blocking | 导出格式 | 判据 5 | R3 | 用户确认 | open | — | 待确认（R3 前置） |
| I-029-005 | non-blocking | HTTP 自助核销是否本波；若是则本 VP 做限流评估 | 判据面 | R1 | 用户裁决 | closed | — | D-002: 本波仅交付 Go 模块内部 API Redeem |
| I-029-006 | required | 哈希算法 / 码熵 / 常时比较 / UNIQUE+同事务双花 | 判据 2/3 | R1 | 用户裁决 | closed | — | D-002: 高熵码 + SHA-256 + 单事务 CAS 原子核销入金 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-02 | 工作区 / Root 建立与开区决策 | accepted | `01-decision/D-001-workspace-root-establishment.md` |
| D-002 | 2026-09-02 | R1 核心合同冻结与技术选型裁决 | accepted | `01-decision/D-002-r1-contract-freeze.md` |
