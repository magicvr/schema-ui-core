---
id: GOAL-001-timestamptz-persistence-contract
doc: decision
status: active
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-040-001 | required | SQLite 与 PG 的合同平等物理类型、精度、UTC 表示、NULL/零值、编解码 | R1/R2 | R1 | 对照 VP-013 与现有时间列；形成 R1 冻结决策 | collecting | R1 冻结前复核 | 用户 D-002：PG `timestamptz(6)`；SQLite fixed-6 UTC RFC3339 TEXT；sentinel 0 → NULL；待逐列证据 |
| I-040-002 | required | 首波时间列分母与非时间 INTEGER 排除清单 | R1/R2 | R1 | 全仓扫描并冻结列清单 | collecting | R1 冻结前复核 | 用户 D-002：全部绝对时刻列纳入，ID/duration/step/version/计数/金额/flag 排除 |
| I-040-003 | required | 存量升级与备份 residual | R1/R3 | R1 | 对照 dump/restore 路径并取得必要书面裁决 | collecting | R1 冻结前复核 | 用户 D-002：SQLite/PG 各自原地转换；不提供 SQLite→PG 产品搬运器 |
| I-040-004 | required | VP-020 展示/输入与 UTC 存储回归矩阵 | R3 | R1 | 复用 VP-020 验收用例，补存储形状对照 | open | — | 待 R3 |
| I-040-005 | required | 激活前置与工作区绑定 | 激活 | 激活前 | `/vision` self Review + `/govern` scaffold | verified | — | VRev-104；workspace/Root 已建立 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-20 | 激活边界与 R1 默认候选登记 | accepted | `01-decision/D-001-activation-and-r1-candidate.md` |
| D-002 | 2026-09-20 | R1 合同选型与用户裁决 | accepted | `01-decision/D-002-r1-contract-freeze-user-decisions.md` |
| D-003 | 2026-09-20 | R1 公共时间输出合同裁决 | accepted | `01-decision/D-003-r1-public-wire-contract-user-decision.md` |
| D-004 | 2026-09-20 | R2 migration 归属与 catalog 形态裁决 | accepted | `01-decision/D-004-r2-migration-ownership-user-decision.md` |
| D-005 | 2026-09-20 | 公共 wire 输入兼容裁决 | accepted | `01-decision/D-005-r1-wire-input-compat.md` |
| D-006 | 2026-09-20 | C3 统一 Backup SPI/Service 裁决 | accepted | `01-decision/D-006-r1-backup-spi-user-decision.md` |
| D-007 | 2026-09-20 | 最小 Backup/RecoveryPoint Port 公共边界裁决 | accepted | `01-decision/D-007-backup-port-surface-user-decision.md` |
| D-008 | 2026-09-20 | C2 精度、D0 与 PG backup provider 裁决 | accepted | `01-decision/D-008-c2-precision-zero-backup-decisions.md` |
| D-009 | 2026-09-20 | 公共 wire 非 DB 例外范围裁决 | accepted | `01-decision/D-009-wire-nondb-exceptions.md` |
| D-010 | 2026-09-20 | Backup Port 最小方法裁决 | accepted | `01-decision/D-010-backup-port-methods-user-decision.md` |
| D-011 | 2026-09-20 | schema_migrations conversion owner 裁决 | accepted | `01-decision/D-011-schema-ledger-owner.md` |
| D-012 | 2026-09-20 | Voucher 时间异常值策略裁决 | accepted | `01-decision/D-012-voucher-invalid-value-policy.md` |
| D-013 | 2026-09-20 | 微秒单调 updated_at 裁决 | accepted | `01-decision/D-013-monotonic-updated-at-policy.md` |
| D-014 | 2026-09-20 | v73–v87 allocation baseline 裁决 | accepted | `01-decision/D-014-v73-allocation-baseline.md` |
| D-015 | 2026-09-20 | 负时间与微秒截断解释 | accepted | `01-decision/D-015-negative-instant-truncation.md` |

> 新决策从 `01-decision/D-NNN-<slug>.md` 写入；编号在本目标内单调不复用。
