---
id: GOAL-002-r1-contract-and-denominator-freeze
doc: decision
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.1
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 状态 | 证据 / 决策 |
|----|------|-----------------|----------|------|-------------|
| I-040-001 | required | PG `timestamptz(6)` / SQLite fixed-6 RFC3339 TEXT 的精度、编解码、排序与 NULL 规则 | C2/R2 | **collecting（显著收窄）** | 90 列逐列合同、谓词 exact SQL、逐表 SQLite/PG DDL 均已落盘；**负值政策按列分档、PG 毫秒式经 15/16/17 实测更正**；剩余见 `03-audit.md` A-044（F-I-002 未闭） |
| I-040-002 | required | 全部绝对时刻列与排除列分母 | C1/C2/R2 | **collecting（分母已冻结口径）** | inventory v0.3：**90 列 / 44 张表**（「倒数第二段 = 表名」机械去重，A-044 独立复算一致）；排除列清单已在 inventory 列明 |
| I-040-003 | required | SQLite/PG 原地转换、失败恢复、备份依赖 | C3/R2/R3 | **collecting（C3 边界已落盘）** | `r1-c3-backup-recovery-boundary-v1.0-fc.md`（三类产物区分、双 token、调用点、harness、错误分类）经 A-042 判 F-I-004 closed；PG 跨版本兼容矩阵仍待 R3 |
| I-040-004 | required | VP-020 展示/输入回归矩阵 | R3 | open | Root I-040-004；R1 仅登记接口，R3 执行 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-20 | R1 合同冻结（承接 Root 用户裁决） | accepted | `01-decision/D-001-r1-contract-freeze.md` |
| D-002 | 2026-09-20 | R2 migration 归属承接 | accepted | `01-decision/D-002-r2-migration-ownership.md` |
| D-003 | 2026-09-20 | C2 wire 输入兼容承接 | accepted | `01-decision/D-003-wire-input-compat.md` |
| D-004 | 2026-09-20 | C3 Backup SPI/Service 合同承接 | accepted | `01-decision/D-004-backup-spi-contract.md` |
| D-005 | 2026-09-20 | C2/C3 guardrails 草案 | proposed | `01-decision/D-005-c2-c3-guardrails-proposed.md` |
| D-006 | 2026-09-20 | 最小 Backup/RecoveryPoint Port 承接 | accepted | `01-decision/D-006-backup-port-surface.md` |
| D-007 | 2026-09-20 | C2 精度、D0 与 PG backup provider 承接 | accepted | `01-decision/D-007-c2-precision-zero-backup.md` |
| D-008 | 2026-09-20 | C2 wire 非 DB 例外承接 | accepted | `01-decision/D-008-wire-nondb-exceptions.md` |
| D-009 | 2026-09-20 | Backup Port 最小方法承接 | accepted | `01-decision/D-009-backup-port-methods.md` |
| D-010 | 2026-09-20 | schema_migrations owner 承接 | accepted | `01-decision/D-010-schema-ledger-owner.md` |
| D-011 | 2026-09-20 | voucher / monotonic time policies 承接 | accepted | `01-decision/D-011-voucher-monotonic-policies.md` |
| D-012 | 2026-09-20 | v73 allocation / negative truncation 承接 | accepted（**毫秒式已于 2026-09-20 更正**） | `01-decision/D-012-v73-allocation-negative-truncation.md` |
| D-017 | 2026-09-20 | v73+ checksum 约定（用户 P-004 选 A：单 checksum / SQLite DDL 切片） | accepted | `01-decision/D-017-v73-checksum-convention.md` |
| D-018 | 2026-09-20 | R2 行拷贝机制与扫描器兼容窗口（选项 C + 无过渡期） | accepted | `01-decision/D-018-r2-row-copy-and-compat-window.md` |
| D-019 | 2026-09-20 | FK 父表重建模式与子表处置（F-5 子女先行 + 两次重建） | accepted | `01-decision/D-019-fk-parent-rebuild-mode.md` |
| D-020 | 2026-09-20 | F-I-002 可执行测试的载体与重定向（一次性验证库） | accepted | `01-decision/D-020-fi002-executable-test-scope.md` |
| D-021 | 2026-09-20 | F-I-005 关门口径（拆分）与 PG 临时容器验证授权 | accepted（**F-I-005 闭合待用户 residual 裁决，见 A-044**） | `01-decision/D-021-fi005-gate-and-pg-verification.md` |

> **编号说明**：`D-013`～`D-016` 在本 child 目标内**不存在**（历史未使用）；child 编号的连续性由 A-028/A-029 的 E-ID/D-ID 卫生要求维护为「已用编号单调不复用、允许空洞」。引用 Root 决策时一律写 **Root D-0NN**（见各文件的编号限定块）。
>
> 用户裁决原文与范围记录在 Root `D-002-r1-contract-freeze-user-decisions.md`；本子目标承接并将其转为可验证 C1～C4 交付物。
