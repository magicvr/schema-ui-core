---
id: A-001-self-r3d-exit-matrix
doc: audit-entry
status: active
parent: GOAL-008-r3-exit-matrix-and-root-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# A-001 · self 审计：R3-D 退出判据矩阵与 Root 关门就绪度

- **source**: `self`
- **日期**: 2026-09-21
- **scope**: `GOAL-008` 检查点 A（退出判据证据矩阵、判据 5 反向核验、残留清账）与 **Root 关门就绪度**（判据 1–6 的证据充分性）；不含对已关门目标技术结论的重审。
- **verdict**: `conditional`（**开放 required = 0**；3 条 recommended 提交 independent 复核；2 条限定已写入矩阵）

## 成果（逐项可核对）

| # | 判据 | 结论 | 证据 |
|--:|------|------|------|
| 1 | 判据 1 证据充分 | **成立** | Root `D-002`（`GOAL-001/01-decision/D-002-r1-contract-freeze-user-decisions.md`，R1 承接 `GOAL-002/01-decision/D-001-r1-contract-freeze.md`）与 Root `D-003`/`D-005`/`D-009`/`D-015`；`r1-time-column-inventory-v0.3.md`；`internal/temporalcontract`（`Count = 90`）；`internal/temporal` codec |
| 2 | 判据 2 证据充分 | **成立** | 15 个 v73–v87 descriptor；本轮实测 **v73–v87 的 15/15 checksum 冻结在 `internal/store/migrate_test.go`**；真实迁移边界矩阵（`internal/w040contracttest/migration_boundaries_test.go`） |
| 3 | 判据 3 证据充分（含真实 PG 路径） | **成立（本环境）** | 本轮 `TestPGRestoreToNewDB` 13.51s / `TestSQLiteRestoreToNewDB` 0.74s / `TestC3RecoveryAnchorsOnPostgresUpgrade` 19.59s / `TestCompositionPostgresStartup` 14.13s 全 PASS、**无一 skip** |
| 4 | 判据 4 证据充分（升级后恢复） | **成立** | C3 harness + 升级路径测试 + R3-C 跨版本矩阵（18 个 supported 格含 15→16/17，形状校验全通过） |
| 5 | 判据 5 反向核验自身可核对 | **成立** | 矩阵 §5 逐项列出命令与结果：依赖零变更、无 ORM/Redis/MQ/第三库、驱动类型未进公共契约、kernel 仅 `backup.go` 被触碰、108 个生产文件范围 |
| 6 | 判据 6 未被越读 | **成立** | 矩阵把判据 6 记为**部分满足**并显式写「用户确认仍待 `I-041-010`」；无「以矩阵代替用户确认」的表述 |
| 7 | 跨目标开放 required 汇总 | **成立** | `GOAL-002` A-046/A-047、`GOAL-003` A-003、`GOAL-004` A-003、`GOAL-005` A-004/A-005、`GOAL-006` A-004/A-005、`GOAL-007` A-002/A-003 —— 六个已关门目标收口意见均为 **开放 required = 0** |
| 8 | residual 未重记 | **成立** | R1 `D-021` 的 `F-I-005` 在矩阵中记为 `fixed`（`GOAL-002/A-048`），未记为 open |
| 9 | 台账同步 | **成立** | `GOAL-008` `00-meta`/`02-execution`/`E-002`/附件 + `goal-tree` 树与表与叙述 + Root `00-meta` R3 行均已对齐 `1/3` |

## 限定（已写入矩阵，供审计复核）

| ID | 限定 | 位置 |
|----|------|------|
| `L-1` | 判据 3/4 的「真实路径」结论**绑定本环境**存在常驻 PG 15.4 且 `PG_TEST_*` 已配置；缺该配置时相关测试会 skip 并记录原因 | 矩阵 §3「诚实边界」 |
| `L-2` | 判据 5 的「多实例」判断基于**本工作区变更范围**（108 个生产文件均落在 temporal/store/migration/wire/backup/jobs），不是全仓架构复审 | 矩阵 §5 最后一行 |

## 提交 independent 复核的事项

| ID | 级别 | 事项 | 复核点 |
|----|------|------|--------|
| `F-S-001` | recommended | 矩阵判据 5 中「`sql.NullTime` 属 `database/sql` 而非驱动类型」的判定是否站得住（`authsession.User.LockedUntil` 被 handler 直接读取） | 请对照 Root 红线原文（禁止 `pgtype`/驱动时间类型/`*sql.Tx` 进入公共契约）与 `GOAL-004` 前置 `D-001` §2 #5，确认不构成残留泄漏或应升为 finding |
| `F-S-002` | recommended | `internal/testsupport/store.go` 在内部辅助函数签名中使用 `*sql.Tx` 是否越界 | 该文件是测试支撑包（非 handler/模块公共契约）；请判定是否需要收窄 |
| `F-S-003` | recommended | 判据 2 只说 v73–v87 的 15 个 **SQLite** checksum 被冻结；PG 变体按 `D-017` 不进哈希 | 请确认这与已冻结决策一致，且矩阵未把 PG 变体误述为「已哈希」 |
| `F-S-004` | note | 判据 6 的「用户确认」是 Root 关门的硬门禁；本 self 审计**不**代替它，也不预判用户会接受 | 请确认矩阵/本审均未越权放行 |
| `F-S-005` | note | `docs/vision/roadmap.md` 的 VP-040 行为立项时投影，未随 R3 进展更新 | 属 `/vision` 层动作，计划在检查点 C（用户确认后）同步；请确认这属已知差异而非矛盾 |

## 自审可复跑证据

- 判据 3/4：`go test -count=1 -v -run "TestPGRestoreToNewDB|TestSQLiteRestoreToNewDB|TestC3RecoveryAnchorsOnPostgresUpgrade|TestCompositionPostgresStartup" ./internal/backup/ ./internal/store/ ./internal/composition/` → 全 PASS、无 SKIP。
- 判据 2：`internal/store/migrate_test.go` 含 v73–v87 的 **15/15** checksum 断言（本轮以脚本核对哈希文本）。
- 判据 5：依赖 diff / 依赖字符串 / 驱动类型 / kernel 触碰范围四条扫描（矩阵 §5 记录命令与结果）。
- 全仓基线：`go build ./...` exit 0；`go test -count=1 ./...` 见 `02-execution.md`。

## 自审结论

退出矩阵的六条结论均有可核对产物，未发现开放 required；两条限定已显式写入矩阵而非隐去；`F-S-001`～`F-S-005` 提交 grok independent 复核。**Root 关门仍取决于用户确认（判据 6）**，本审不放行。
