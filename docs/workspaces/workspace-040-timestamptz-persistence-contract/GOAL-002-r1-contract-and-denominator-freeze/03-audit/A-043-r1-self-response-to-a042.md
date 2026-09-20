---
id: A-043-r1-self-response-to-a042
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-042 / negative-value policy scoped per column, table count corrected to 44
verdict: conditional
open_required: 2
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-043 · R1 self response to A-042

- **source**：self（编排器响应，**不**冒充 independent）
- **verdict**：conditional
- **开放 required**：**2**（F-I-002、F-I-005）——由 3 降至 2

## 1. 接受 A-042 的判定

| finding | A-042 verdict | 本响应动作 |
|---------|---------------|------------|
| **F-I-004** | **closed**（§G 七项均已对位：G1 runbook superseded + 两处错误就地标注 + 历史正文未改写；G2 Port L60/L71 与 conversion §0 收口；G3 PG 表与 `postgres.go` 一致且 `applyPendingPG` 缺口属实；G4 A/C 已分且混合形状方向正确；G5 四条可落码；G6 六类足以堵住缺文件误绿；G7 实测驱动可落码） | 接受闭合 |
| **F-I-027** | **closed**（runbook / Port / conversion §0 三源已收口） | 接受闭合 |
| **F-I-002** | **仍 open**：A-038 的可执行测试子项 `fixed`；剩余 = 冻结包把「负值一律 m0 fail closed」写成**全局规则**，与 Root `D-012`（仅 voucher）/ Root `D-015`（负 epoch 合法 instant）**及已绿测试矛盾** | 见 §2 |
| **F-I-005** | 仍 open（本轮未触及） | 照录 |
| **F-I-025** | recommended 仍 open：「20 张带时间列的表」≠ 独立计数 **44** | 见 §3 |
| 新 required / recommended | **无** | — |

**A-042 的独立验证已照录**：它在 `apps/api` 独立跑 `go test ./internal/w040contracttest/ -v -count=1`，三测试全绿；独立核算负毫秒 floor、999 ms 无进位、公元 9999 年、D0 `0→NULL`、voucher `0→NULL`、27 字符的期望值一致；并确认现行测试**没有**把 `D-012` 泛化到普通列。

## 2. F-I-002 剩余：冻结包负值政策按列分档（本轮修正）

**A-042 指出的是一个我自己造成的传播缺口**：第 7 轮的测试**已经**暴露出「负值政策被过度泛化」，但我只改了**测试代码**，**没有**把同一更正传播到冻结包正文。这正是 A-042 说「与已绿测试矛盾」的原因。

**修正的四份冻结载体**：

| 载体 | 原表述 | 修正后 |
|------|--------|--------|
| `r1-c2-predicate-exact-sql-v1.0-fc.md` §1 与 §6 `m0` | 「`< 0` **只**由 `m0` 预检 fail closed」「全部 sentinel 列一致」 | 新增**两档政策表**：仅 voucher `#72`/`#73` 的 `bucket_negative > 0` 触发回滚（**Root** `D-012`）；**其余全部时间列记录但不阻断**（**Root** `D-015`：负 epoch 是合法 instant）；并明确非 sentinel 列出现 0 时按 epoch instant 转换 |
| `r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` §2 | 「负值 `< 0` **一律**不进表达式；只由 `m0` 预检 fail closed」 | 改为**按列分档表**（voucher 两列 fail closed / 其余正常转换），并指向可执行断言 `TestNegativeMustFailClosed` 的两个子测试 |
| `r1-c2-per-table-rebuild-ddl-v1.0-fc.md` §0 | 「**不进任何表达式**；只由 `m0` 预检 fail closed」 | 同上分档表述 |
| `r1-c2-per-column-conversion-contract-v1.0-fc.md` §3.5 | 「`< 0` 只走 `m0` 预检 fail closed」 | 「负值单路径 + **分档政策**」 |
| `r1-c3-backup-restore-runbook-v0.1.md`（superseded 正文） | 断言样本含 `negative-invalid` | 就地标注**已移除**并说明理由 |

> `#72`/`#73` 两行的「`< 0` 不进 USING，只走 `m0` 预检 fail closed」**保留不变**——它们**就是** voucher 两列，政策正确，不是泛化。

## 3. F-I-025：表数 20 → **44**（已实测复核）

我用与 A-042 相同的机械方法复核：对 inventory v0.3 的 90 行，按「**倒数第二段 = 表名**」去重（该格式对 `module.table.column`、`module.sub.table.column`、以及 authsession 的 `table.column` 三种形状统一成立），得到**恰好 44 张不同的表**（清单：`captcha_challenges` … `wallet_reconciliation_runs`）。

- `r1-c2-per-table-rebuild-ddl-v1.0-fc.md` §0 覆盖说明已由「20 张带时间列的表」改为「**44 张不同的表**（机械去重所得；分母仍为 **90 列**）」。
- 我此前的「20」是错的（既非 44 也非 15 descriptor）；已按实测改为 44，并写明推导方法以便复核。

## 4. 仍开放（不得放行）

- **F-I-002**：负值政策已按列分档收口；整条是否闭合**待 independent 复审本轮修正**。
- **F-I-005**：需 R2 落码才能记录真实 `MigrationChecksum` 哈希；测试改写同属 R2。
- **F-I-025**（recommended）：计数已改为实测 44；是否闭合待复审。

**本响应不闭合任何 required。** C2/C3 未冻结，R2 未放行。

## 5. 两项需用户裁决的事项（另见下轮汇报）

1. **PG 侧从未被经验验证**：本机无 `psql`/`pg_dump`/`pg_restore`，5432 无监听，compose 只有 api/web → PG 转换表达式（`date_trunc` + `to_timestamp` / 整数 interval）**从未实测**。Docker 可用且 `postgres:16`/`15-alpine`/`17-alpine` 镜像已在本地，可临时起容器验证。
2. **F-I-005 的关门口径**：该条要「已记录的 `MigrationChecksum` 哈希」，而哈希只能对**真实存在的迁移语句**求值，即必须先有 R2 代码——存在与 F-I-002 同类的结构性张力，需用户明确 R1 的关门口径。
