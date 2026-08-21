---
id: A-005
goal: GOAL-019-r3-s14-wallet-ledger
title: A-004 闭合复审 · admin.wallet
date: 2026-08-16
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-004 F-001~F-006 闭合核验
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-019-r3-s14-wallet-ledger
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# A-005 · 独立审计意见（A-004 闭合复审 · S-14 钱包/账务）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：finding-closure · A-004 F-001/F-002 required 闭合核验 + F-003~F-006 勘误核验
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`root_goal: GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`plan_refs`/`primary_plan` = `VP-011-admin-functional-modules`；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。
- **对照原文**：`03-audit/A-004-s1-independent.md`（verdict conditional；F-001/F-002 required · med；F-003~F-006 recommended）。
- **闭合主张**：`01-decision/D-003-a004-response.md`（路径 = **fixed**，拒 residual）；`01-decision/D-002-s1-plan-freeze.md` **v1.1.0**；`02-execution/E-003-a004-response.md`；`03-audit.md` 响应记录。
- **已通读**：本目标 `00-meta`（progress 0/5、S1 检查点、I-001~I-004）、`01-decision.md`、D-002 v1.1.0、D-003、`02-execution.md`、E-001/E-002/E-003、`03-audit.md`、A-001～A-004；`goal-tree.md` GOAL-019 行（树 L39 / 表 L68，均为 0/5）。
- **代码核对**（按 D-002 引用验 F-005 基线，非本复审实施证据）：`composition/composition_test.go` L466–471（admin **27** 权限 / **13** 导航）；`operationlog/repository.go` L14–59（EventXxx `{domain}.{action}` kebab）；`store/store.go` L67–81（`WithTx`）；`datapermission/migration/migration.go` L41（Version **27**）；`operationlog/migration/migration.go` L262–266（max Version **30**）；`kernel/provider.go` L403–418（DefaultNavigationOrder **13** 项）。
- **P-005**：I-001~I-003 required、最晚阶段 S1、现标 `verified`；本复审只核验 A-004 指出的 I-001 缺口（金额原语 + 幂等范围）是否已由 D-002 v1.1.0 补满。无共享资料引用。
- **covered**：A-004 六条 finding 的闭合证据是否可重复核对；语义是否唯一可解释；S2 是否可按表实现与测试。
- **excluded**：全目标重审、S2 实现、S3～S5、不改 status/progress/goal-tree/方案正文。本条 **不是** S1 方案重审。
- **保证等级**：L0（入口分离）。不得解读为第三方鉴证。

## F-001 / F-002 闭合核验

### F-001 · amount_delta 语义互否 + 快照链可执行规则 → **closed（fixed）**

| 字段 | 值 |
|------|-----|
| A-004 level | required（med） |
| 声称路径 | D-003 L17：`fixed`（拒 residual） |
| 本复审 | **闭合成立** |

| A-004 要求 | 落点 | 判定 |
|------------|------|------|
| 按 `entry_type` 的语义/符号/作用列/拒绝条件 | D-002 §1 L24–30 apply 表：adjust = 总额变动、`d ≠ 0`、`total += d` / `available += d` / frozen 不变，拒任一余额 &lt; 0、恒等式破、disabled；freeze = available→frozen、`d > 0`、total 不变、`available -= \|d\|` / `frozen += \|d\|`，拒 available 不足、disabled；unfreeze = frozen→available、`d > 0`、total 不变、反向转移，拒 frozen 不足、disabled | 语义唯一。原「总额变动 + freeze 总额不变 + CHECK ≠ 0」互否已解除：freeze/unfreeze 的 `amount_delta` 是**转移额度**（必正），不是总额变动 |
| `apply(prev, entry) == after_*`；首笔 prev=(0,0,0)；末笔 == 账户三余额；每笔恒等式 | D-002 L32：`apply(prev_after, entry) == entry.balance_after_*`；首笔 (0,0,0)；末笔 == 当前三余额；每笔 `total = available + frozen` | 可执行。S2 测试可按表重放 |
| 链序 `(created_at ASC, id ASC)` | D-002 L32；§6 L122 对账测试复述同一序 | 秒级并列键已写死 |
| §4 快照恒等式 CHECK 与 `CHECK(amount_delta != 0)` 自洽 | D-002 L84 `amount_delta != 0`；L88 `balance_after_total = balance_after_available + balance_after_frozen` | **自洽**：adjust/freeze/unfreeze 均禁止 0；freeze/unfreeze 应用层更严（`> 0`），DDL 是子集约束，不再互否。恒等式在三行 apply 下代数成立（freeze/unfreeze 的 ±\|d\| 对冲） |

S2 可按表实现与测试。空账本（无末笔）与账户 (0,0,0) 的衔接属实现边界，表已给出首笔 prev，不构成语义互否。

### F-002 · 幂等键跨账户泄露 → **closed（fixed）**

| 字段 | 值 |
|------|-----|
| A-004 level | required（med） |
| 声称路径 | D-003 L18：`fixed` |
| 本复审 | **闭合成立**；跨账户泄露路径已排除 |

| A-004 要求 | 落点 | 判定 |
|------------|------|------|
| UNIQUE = `(account_id, idempotency_key)` 或显式全局并禁跨账户返回 | D-002 L33 裁定复合；§4 L92 `UNIQUE (account_id, idempotency_key)` | 账户 B 复用账户 A 的 key **不再**命中全局 UNIQUE，A-004 所述「命中后返回他户流水」路径不成立 |
| 同账户 + 同 key + 同载荷 → 返回既有；同 key + 异载荷 → `LEDGER_IDEMPOTENCY_CONFLICT` | D-002 L33：同载荷 = `entry_type` / `amountDelta` / `memo` 一致；异载荷 → `LEDGER_IDEMPOTENCY_CONFLICT`；§6 L123 将该码列入候选 | 分流已裁定，可测 |
| 查找必须带 `account_id`，禁止裸 key 取他户 | D-002 L33 明示 | 跨账户按裸 key 读取被禁止 |

可选键 + SQLite 复合 UNIQUE 对多行 `NULL` key 的允许多重（`NULL ≠ NULL`）与「可选」一致，不是缺口。

**残留文案（不重开 required）**：D-002 §1 L21 仍写「可选 `idempotency_key`（**UNIQUE**；**重复提交返回既有流水**）」，未带账户范围与异载荷分流。权威裁定在同节 L33 + DDL L92；S2 须按 L33/L92 实现，不得按 L21 做全局 UNIQUE 或一律返回。建议 `/govern` 随手删 L21 旧句，**不**阻断 S2。

## F-003~F-006 勘误核验

### F-003 · 快照 CHECK + 链序 → **closed（fixed）**

- 链序：D-002 L32 / L122 = `(created_at ASC, id ASC)`。
- 快照恒等式：D-002 L88 `CHECK (balance_after_total = balance_after_available + balance_after_frozen)`。
- 列表索引仍为 `(account_id, created_at DESC)`（L98），服务读序，不替代对账链序。不重开。

### F-004 · disabled 拒解冻 → **closed（fixed）**

- D-002 L35：disabled 拒绝调账/冻结/**解冻**，「冻结资金随停用锁定，流水只读」。
- apply 表 unfreeze 拒绝条件含「账户 disabled」（L30）。
- §6 L121：disabled 账户拒绝写。裁定唯一。

### F-005 · 组合根 26→29 → 27→30 → **部分闭合**（recommended 残留）

- **§6 已勘误**：D-002 L124「admin 权限键 **27→30**（+wallet.read/write/adjust；实测基线 27）」「导航 **13→14**」；迁移 **30→32**。
- **代码基线与勘误一致**（非实施完成证据）：`composition_test.go` L471 `wantPermissions: 27, wantNavigation: 13`；`kernel/provider.go` L403–418 共 13 项（末项 `menu_data_permission`）；operationlog 迁移 max Version **30**（L263）。27+3=30、13+1=14、30+2=32 自洽。
- **未完全擦净**：D-002 §8 L143 S2 清单仍写「组合根快照（权限 **26→29**、导航 13→14）」，与 §6 / 实测互否。S2 若照抄 §8 会写错断言。
- 级别保持 recommended。A-004 已声明本条不单独阻断 S2；实施以 **live snapshot + §6（27→30）** 为准，禁止采用 §8 的 26→29。

### F-006 · 台账投影 → **主体闭合**（recommended 轻残留）

| A-004 指摘 | 现态 | 判定 |
|------------|------|------|
| `03-audit.md` 信息表写 I-00N **open** vs meta **verified** | `03-audit.md` L17：I-001~I-004 均 **verified**（D-002 v1.1.0 + D-003）；与 `00-meta` L42–45 一致 | 已勘误 |
| S1 勾选含未执行 independent；`progress: 1/5` | `00-meta` L8 / L36：`progress: 0/5`；L30 S1 文案「冻结稿 + self 完成；**independent 闭合中**」「独立审 pass 落盘前 S1 检查点不计数」 | 符合 A-004「闭合前不投 1/5」建议。检查点未拆成两个 `[ ]`，仍用一条 `[x]` +「不计数」——投影略含混，但派生公式未把 S1 计入 0/5 |
| goal-tree 与 meta 不一致 | `goal-tree.md` L39 树、L68 表均为 **0/5**，与 `00-meta` L8 一致 | 一致 |

**轻残留（不阻断）**：`03-audit.md` L34「结论状态」仍写「S1 grok build independent **待执行**」，与同文件 L36–45（A-004 已落盘、待本复审）不同步。历史结论段未刷新，响应记录已覆盖。本审计员不代改。

progress 0/5 **不得**证明本复审通过；本条只确认投影与 A-004 建议方向一致。

## 必改项汇总

**无。** scope 内无未关闭 required。F-001 / F-002 按 P-003 **fixed** 路径合法闭合（决策 D-003 + 产物 D-002 v1.1.0 + 执行 E-003 + 本复审可重复核对）。

开放 recommended（不阻断 S2）：

1. **F-005 残留**：D-002 §8 L143 仍写权限 26→29。
2. **F-002 文案残留**：D-002 §1 L21 旧「UNIQUE / 返回既有」句。
3. **F-006 轻残留**：`03-audit.md` 结论状态「independent 待执行」过时。

## 与既有意见

- **A-004**（independent · conditional）：本条只核验其六条 finding，不重开方案全量。required 两条闭合；recommended 三条闭合、一条部分、台账主体闭合。
- **A-003**（self · pass）与 A-004 对 S2 放行的门禁互否（P-004.2）因 F-001/F-002 已 fixed **解除**。本复审不代替 `/govern` 改 status。
- **A-001 / A-002** 立项 scope，与本条无冲突。A-002 的 019-F-001/F-002（关联单据 / 权限拆分）仍由 D-002 §1/§3 覆盖，本复审不重开。

## 结论 + S2 放行条件

**verdict: pass**。A-004 阻断 S2 的两条 required 已有可重复核对的勘误；金额原语唯一可解释，幂等复合范围排除跨账户泄露。

本意见视角的 **S2 放行条件**（编排器执行，本条不改状态）：

1. F-001 + F-002 视为已合法闭合；**不再**以 A-004 required 阻断 S2。
2. I-001 设计层 `verified` 现与 apply 表 / 幂等裁定一致，可作为 S1 信息门禁关闭证据（仍不是实现已落地）。
3. 实施按 D-002 **§1 apply 表 + L33 幂等 + §6 27→30 / 13→14 / 30→32**；忽略 §8 L143 的 26→29。
4. 迁移仍从 31/32 起；实施前复核 max version（本复审核实当前 max = 30）。
5. `/govern` 可将 S1 检查点计入、progress 0/5→1/5 并同步 goal-tree；顺手清 L21 / §8 / 结论状态残留更干净，非门禁。

## 声明

本意见不修改 `status` / `progress` / goal-tree / 方案正文。响应由 `/govern` 处理。保证等级 L0。
