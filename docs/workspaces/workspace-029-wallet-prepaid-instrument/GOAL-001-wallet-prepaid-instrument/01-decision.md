---
id: GOAL-001-wallet-prepaid-instrument
doc: decision
status: done
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.3.0
---

# 决策记录 · GOAL-001-wallet-prepaid-instrument

## 信息需求与阶段门禁

> 权威表在 `00-meta.md`。全部 9 项 required/non-blocking 信息项均已 closed。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-029-001 | required | 主体落点 + `owner_type` + `OwnerExists` 迁出 `UserByID` | 方案冻结 + 判据 1 | R1 | 用户裁决 | closed | — | D-002: 独立主体表 subjects，owner_type 扩充 subject，OwnerExists 校验主体 |
| I-029-002 | required | 核销 `entry_type` | 判据 3/4 | R1 | 用户裁决 | closed | — | D-002: 复用 adjust + ref_type='voucher' |
| I-029-003 | required | 生成权限键 | 判据 5 | R1 | 用户裁决 | closed | — | D-002: 新增细粒度 wallet.voucher.issue 权限键 |
| I-029-004 | non-blocking | 导出格式 | 判据 5 | R3 | 用户确认 | closed | — | D-001: API 一次性返回明文数组作为导出数据源（明文不入库） |
| I-029-005 | non-blocking | HTTP 自助核销是否**首波**；若是则本 VP 做限流评估 | 当时判据面 | R1 | 用户裁决 | closed | — | D-002: 首波仅模块 API；R5 见 I-029-007/008 |
| I-029-006 | required | 哈希算法 / 码熵 / 常时比较 / UNIQUE+同事务双花 | 判据 2/3 | R1 | 用户裁决 | closed | — | D-002: 高熵码 + SHA-256 + 单事务 CAS 原子核销入金 |
| I-029-007 | required | R5 HTTP 路径 / 函数形状 | 判据 8 · GOAL-005 S1 | R5 S1 | GOAL-005 D-002 | closed | — | POST /api/wallet/me/redeem + RedeemForUser |
| I-029-008 | required | R5 已登录核销限流评估 | 判据 10 · GOAL-005 S1 | R5 S1 | GOAL-005 D-002 | closed | — | 内存专用桶 15min/10/user id |
| I-029-009 | required | R5 权限模型 | 判据 8 | R5 重开 | 用户确认 A | closed | — | D-003: identity-only |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-02 | 工作区 / Root 建立与开区决策 | accepted | `01-decision/D-001-workspace-root-establishment.md` |
| D-002 | 2026-09-02 | R1 核心合同冻结与技术选型裁决 | accepted | `01-decision/D-002-r1-contract-freeze.md` |
| D-003 | 2026-09-02 | R5 重开：我的钱包预付凭证自助核销 | accepted | `01-decision/D-003-r5-reopen-my-wallet-self-redeem.md` |
| D-004 | 2026-09-02 | 根目标 GOAL-001 与工作区 029 结项关门裁决 | accepted | `01-decision/D-004-root-r5-closeout.md` |
