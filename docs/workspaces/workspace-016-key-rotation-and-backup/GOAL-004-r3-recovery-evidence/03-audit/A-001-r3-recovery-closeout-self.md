---
id: A-001
doc: audit-entry
goal: GOAL-004-r3-recovery-evidence
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# A-001 · R3 恢复证据关门自审（2026-08-22）

- **source**：self
- **auditor**：编排器（/govern）
- **类型 / scope**：close-out · GOAL-004 全部（R3 轮换后恢复证据，双方言）
- **verdict**：pass

## 成果（有证据）

| 成果 | 证据 |
|------|------|
| I-004 verified（最小剧本：T0–T5 时序、两方言命令、A1/A2/A3 断言） | GOAL-004 D-001；Root I-004 行 verified |
| 可重复自动化（composition 级真实启动门禁） | `post_rotation_recovery_test.go`：SQLite 循环无条件必跑；PG 循环按 pgtest+工具双门控，不伪造跳过 |
| SQLite 全循环 PASS | `TestSQLitePostRotationRecovery -count=1 -v` PASS（VACUUM INTO → 备份即库 → 轮换启动 → A1/A2/A3） |
| PG 全循环 PASS | `TestPostgresPostRotationRecovery -count=1` PASS（pg_dump 17→15.4 记录在案组合 + ledger 指纹核对 + 恢复库完整启动 + A1/A2/A3） |
| 无越界 | 零产品代码新增（纯测试文件）；无 Backup API、无第二套 dump、无 PITR |

## 对照成功标准（Root 方向级 3）

| 标准 | 状态 | 证据 |
|------|------|------|
| SQLite `VACUUM INTO` 路径：轮换后从备份启动且鉴权可核对 | 达成 | SQLite 循环 T4/T5 |
| PG `pg_dump`/`pg_restore` 路径：同上 | 达成 | PG 循环 + `assertRestoredLedger`（count+checksum 指纹一致） |
| 不重做 dump | 达成 | 仅消费既有合同与官方客户端 |

## Findings

无 required。备注两条（非 finding）：

1. PG restore 的 17 客户端 → 15.4 服务端组合按 workspace-013 GOAL-006 D-002 允许未知 GUC 告警；测试以 ledger 指纹一致性补强，不盲信 exit code。
2. `R16_PG_DUMP_CONTAINER` 为新引入的可选测试环境变量（指向带 pg_dump/pg_restore 的容器）；未设置且本机无客户端时 PG 循环显式 skip——CI 若需强制执行可在 PG job 中注入该变量。

## 必改项汇总（required 列表）

空。

## 结论 + 建议下一步

GOAL-004 达成关门条件：检查点 4/4、0 required finding、I-004 门禁已闭。建议：GOAL-004 done；Root R3 完成、progress 3/5；git checkpoint。下一阶段 R4（默认单密钥仍可用）：以证据整合为主，开 GOAL-005 承接。
