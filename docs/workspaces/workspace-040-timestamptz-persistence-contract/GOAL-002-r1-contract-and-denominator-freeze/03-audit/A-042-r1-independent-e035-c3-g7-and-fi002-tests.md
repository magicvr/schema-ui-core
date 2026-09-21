---
id: A-042-r1-independent-e035-c3-g7-and-fi002-tests
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · ad-hoc design-plan + finding-closure + execution-facts · E-035 / commit 4462e73d A-040 §G 七项收口对照 F-I-004 · commit a2db84ae / D-020 / apps/api/internal/w040contracttest 对照 F-I-002 · F-I-025/F-I-027 recommended 复核 · 不是实施审计 · freeze-candidate ≠ 已实施
verdict: conditional
open_required: 2
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-042 · R1 independent · E-035 §G 七项 + F-I-002 可执行测试 + F-I-025/027

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：ad-hoc + finding-closure + execution-facts（用户指定三块：① commit `4462e73d` / E-035 对照 A-040 §G 七项能否收口 F-I-004；② commit `a2db84ae` 的 `w040contracttest` 对照 F-I-002 剩余「非法/越界可执行测试」，并独立跑测试；③ F-I-025 / F-I-027 recommended 复核。对照基线 A-040：开放 required = 3，F-I-002 / F-I-004 / F-I-005）。
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区作为审计上下文**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree / `apps/`。
- `git show --stat 4462e73d`：9 个文件均在 workspace-040（execution / A-041 / C3 与 C2 附件）；**`apps/` 无变更**。
- `git show --stat a2db84ae`：新增 `apps/api/internal/w040contracttest/contract_boundaries_test.go`（仅测试包）+ `D-020` + E-034 + A-040。
- **证据窗口**：以上述两 commit 已提交材料 + 现行 `postgres.go` / `migrate.go` + 本机 `go test ./internal/w040contracttest/ -v -count=1` 为准。不把 freeze-candidate 当实施证据。

## 核对方法

1. 通读 A-040 §G 原文、A-041 / E-035、C3 边界 v1.0-fc、旧 runbook、Port draft、conversion contract §0、rebuild DDL §0。
2. 对位 `internal/store/migrate.go`（`actionNoop` `:52-53`、`applyPending` `:81-103`、`snapshotBeforePending` `:279-300`、`verifyIntegrity` `:345-363`）与 `internal/store/postgres.go`（`migrate` `:72-118`、`actionNoop` `:94-95`、`applyPendingPG` `:120-127`、`applyMigrationPG` `:154-173`）。
3. 独立跑 `go test ./internal/w040contracttest/ -v -count=1`（`apps/api`）；用 Python `datetime` 独立计算负毫秒 floor / 999 ms / 公元 9999 年等期望值；对照 Root `D-012` / Root `D-015` / child `D-020`。
4. 从 conversion contract 90 行独立抽出唯一表名；对照 descriptor ledger 表范围与 rebuild §1.1 B/C 组。全扫 attachments 的 `D-0NN` 与「尚未落盘」。

## 成果（有证据）

1. **A-040 §G 七项在 R1 设计层均已对位收口**（见第一块）。F-I-004 **本审接受 closed**（设计层；freeze-candidate ≠ 实施，与 F-I-003 / F-I-006 同一尺子）。
2. **F-I-002 的 A-038 可执行测试剩余已落地且本机全绿**（见第二块）。期望值与 Root D-012 / D-015 政策正确。**整条仍不能闭合**：冻结包仍把「负值一律 m0 fail closed、不进表达式」写成全局规则，与测试及 Root 决策矛盾。
3. **F-I-027 三源已收口，本审接受 closed**。**F-I-025 不能闭合**：rebuild §0 改成「20 张带时间列的表」，独立计数仍是 **44** 张。
4. **无新 required。无新 recommended 编号。** 开放 required **3 → 2**（F-I-002 / F-I-005）。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C3 原地转换 / 备份回滚冻结 | **设计层可冻；实施未发生** | F-I-004 closed；`apps/` 无 Backup 实现 |
| C2 逐列 USING / rebuild / codec 冻结 | **仍不可冻结** | F-I-002 负值政策冻结包不唯一 |
| F-I-002 关闭要求（A-038 原文：可执行测试） | **测试项 `fixed`；整条仍 open** | 见第二块 |
| F-I-004 关闭要求（A-040 §G 七项） | **满足（设计层）** | 见第一块 |
| C2 / C4 / R2 放行 | **未满足** | F-I-002 / F-I-005 仍开 |

---

## 第一块：A-040 §G 七项（E-035 / `4462e73d`）

### G1 · 旧 runbook — **收口**

`attachments/r1-c3-backup-restore-runbook-v0.1.md`：

| 要求 | 本审 |
|------|------|
| `status: superseded` | **是。** frontmatter `status: superseded`，`version: 0.2.0` |
| 就地标注两处错误点 | **是。** 错误点 1：第 1 步 per-migration snapshot 是批次边界产物非常量旧合同；错误点 2：步骤 1 `<artifact>` 被步骤 5 restore 后再由步骤 6 断言 `timestamptz(6)` 必失败 |
| 正文历史保留未被改写 | **是。** 原 SQLite/PG 步骤与 `<artifact>` 命令仍在「历史正文」下；新增的是顶部移交声明与两处标注，不是改写步骤本身 |
| 唯一权威移交边界文件 | **是。** |

满足 A-040「标 superseded 并声明本边界文件为唯一 C3 权威」的或路径。

### G2 / F-I-027 三源 — **收口**

| 载体 | A-040 要求 | 本审 |
|------|------------|------|
| Port draft L60 | 删「`snapshotBeforePending` or」并明令其不是 RecoveryPoint 源 | **是。** 现为「C3-specific native snapshot **after** conversion committed」；`snapshotBeforePending` **explicitly NOT** a RecoveryPoint source |
| Port draft L71 | 标已解决 | **是。** 划掉并指向边界文件 §2/§3 |
| Port 权威范围声明 | 新增 | **是。** 尾注：本文件只辖 kernel Port 类型表面与后置条件 |
| conversion §0 | 收回「尚未落盘」 | **是。** 改为已落盘且为 C3 **唯一权威**；runbook 已 `superseded` |

### G3 · PG 对称调用点 — **收口**（现状缺口属实）

独立核对 `postgres.go`：

| 文件主张 | 实测 | 本审 |
|----------|------|------|
| `postgres.migrate` `:72-118` | **是** | 行号正确 |
| `applyPendingPG` `:120-127` | 循环只调 `applyMigrationPG`，然后 `return nil` | **无 snapshot、无 `verifyIntegrity`。缺口属实** |
| `applyMigrationPG` `:154-173` | `p.Run` 单事务 | 与「事务内禁止调用」一致 |
| `actionNoop` `:94-95` | `return nil` | 连 `verifyIntegrity` 都没有。属实 |

§4.2 PG 表把 A′（循环前一次）/ C′（循环内）/ B′（成功之后）/ 事务内禁止写成**新增设计**，并显式标注上述缺口。作为 R1 设计锚点足够；不是已实施。

### G4 · A 与 C 分离 — **收口**（混合形状方向正确，有一处偏宽）

§2 现为：A = 批次**循环外一次**、旧合同；C = **循环内 per-migration**、**混合形状**；并写「不可合并为批次前一次」。

`snapshotBeforePending` 调用语义：`migrate.go:92-97` 在 `pendingMigrations` 循环内、`version >= 2` 的**每条迁移之前**调用。批次推进后，中途 C 含**部分**新合同列。**混合形状成立。**

偏宽：§2 写「批次末为全部新合同」。C 取在迁移**之前**，v73–v87 批次的最后一次 C 在 v87 **之前**，v87 列仍旧。全部新合同只出现在（i）转换后再有 pending 时的下一次 C，或（ii）B 所对应的转换后库。不阻断 G4，见下方残留。

### G5 · 重试四条 — **收口**（可执行，与现状一致）

现状：SQLite `actionNoop` `:52-53` = `return s.verifyIntegrity()`，不调 `CreateRecoveryPoint`；PG `:94-95` = `return nil`。与 §4.3 问题陈述一致。

四条要求可落码：①「已到目标版本但无 B」必须是显式动作而非 `actionNoop`；② 形状已新 → 补创 B，仍旧 → 走 A/C 回滚；③ 每次启动最多一次补创；④ 禁止把「无 B」当门禁已满足。

用词残留：「不因补创失败而拒绝启动」被写成「沿用现行 fail-closed」。这是**继续服务**（对 B 门禁偏 open），不是 fail-closed。行为本身写清楚了，不阻断 G5。

### G6 · 错误分类 — **收口**

§5.1 六类 + 硬要求：`TestLegacyArtifactMustFail` 必须属于 `{TimeContractMismatch, TemporalColumnSetIncomplete}` **且不属于** `{ArtifactNotFound, ArtifactUnreadable, ToolFailure}`。缺文件会落到 `ArtifactNotFound`，测试必须红。足以堵住 A-040「缺文件也会绿」。

§5 断言 4：**删除 B 正向 round-trip 上的「负值非法」是正确的。** 理由句「转换 fail-closed 后 B 内不应再有负值」**不正确**——普通列负 epoch 转换后以 pre-1970 TEXT 出现在 B 中（本轮测试已证明）。负值非法只属于 voucher 的 `m0`，不是 A/C 形状测试的主断言。不把 G6 整项打回；该理由句并入 F-I-002 剩余（见第二块）。

### G7 · C→B 机械身份 — **收口**（可落码）

§3.1：A/C **不得**生成 `RecoveryPoint`；B 必须含 `TimeContract` + `CatalogVersion` + `ChecksumSet` + `Verification` + 批次末 version；`contract_shape` 由对 restore 目标的**实测**得出；不符 → `TimeContractMismatch`。

- **实测**挡住 A 与早期 C（INTEGER/BIGINT ≠ 目标形状）。
- **转换后的 C** 形状可以过校验（A-040 原文）；挡住它的是「A/C 不得生成 RecoveryPoint / CreateRecoveryPoint 自己产 dump、不吞 C 文件」——这是记录类型/API 身份，比文件名强，可落码。

残留：把 A/C 一律标 `contract_shape = legacy` 与 §2 混合形状矛盾。不阻断 G7。

### F-I-004 判定

**本审接受 F-I-004 closed（设计层）。** A-040 §G 七项全部对位；包路径候选、双 token、harness 规格、禁止事务内调用维持 A-040 已接受。freeze-candidate ≠ 实施，与 F-I-003 / F-I-006 同一尺子，不单独阻断本条。实施与 harness 落地仍在 R2，不把本条重新打开。

残留（不新开号、不挡本条闭合）：§2「批次末全部新合同」偏宽；§3.1 给 C 打 `legacy` 不准；§4.3「fail-closed」用词；§5 断言 4 理由句（并入 F-I-002）。

---

## 第二块：F-I-002 可执行测试（`a2db84ae` / D-020）

### 独立跑测

在 `apps/api`：

```text
go test ./internal/w040contracttest/ -v -count=1
```

| 测试 | 结果 |
|------|------|
| `TestRebuildConversionBoundaries` | PASS |
| `TestNegativeMustFailClosed/voucher_negative_fails_closed` | PASS |
| `TestNegativeMustFailClosed/ordinary_negative_is_valid_instant` | PASS |
| `TestFixedSixLexicalOrder` | PASS |
| 包 | `ok` · ~0.73s |

### 期望值独立核对（Python `datetime`，UTC）

全部与测试 `want*` **一致**：

| 样本 | 输入 | 独立计算 | 测试期望 |
|------|------|----------|----------|
| epoch 0 | sec 0 / ms 0 | `1970-01-01T00:00:00.000000Z` | 同 |
| 正值 | 1758320000 / 1758320000123 | `2025-09-19T22:13:20.000000Z` / `.123000Z` | 同 |
| 999 ms 无进位 | 1758320000999 | `2025-09-19T22:13:20.999000Z` | 同（秒未进位） |
| −1 ms | −1 | `1969-12-31T23:59:59.999000Z` | 同（floor，非向零） |
| −999 ms | −999 | `1969-12-31T23:59:59.001000Z` | 同 |
| −1000 ms | −1000 | `1969-12-31T23:59:59.000000Z` | 同 |
| −1001 ms | −1001 | `1969-12-31T23:59:58.999000Z` | 同 |
| 负一天 | −86400 / −86400000 | `1969-12-31T00:00:00.000000Z` | 同 |
| 1914 | −1758320000123 | `1914-04-14T01:46:39.877000Z` | 同 |
| 公元 9999 | 253402300799 / 253402300799999 | `9999-12-31T23:59:59.000000Z` / `.999000Z` | 同 |
| 27 字符 | 含 pre-1970 与 9999 | 长度均为 27 | 同 |
| D0 0→NULL | `sec_d0`/`ms_d0` = 0 | 断言 `NULL` | 同 |
| voucher 0→NULL | `voucher_sec` = 0 | 断言 `NULL` | 同 |

`millisExpr` 的 `(x-999)/1000` floor 与 `(%1000+1000)%1000` 余数归一化与 mechanism §2.3 及本机计算一致。公元 9999 的毫秒值是 `…799999` 不是秒值（D-020 记载的单位错配已被测试拦住）。

### 负值政策 vs Root D-012 / D-015

| 权威 | 原文 | 测试是否同一 |
|------|------|--------------|
| Root `D-012-voucher-invalid-value-policy.md` | **仅** `vouchers.expires_at` / `redeemed_at`：0→NULL；**负值 fail closed**；正值转换 | `TestNegativeMustFailClosed/voucher_negative_fails_closed`：`negativeIsError=true` 必须报错。**同一** |
| Root `D-015-negative-instant-truncation.md` | **负 epoch 不是 sentinel**；合法 instant，正常转换 | `ordinary_negative_is_valid_instant` + 边界矩阵中的 −1/−999/−1000/−1001/负一天/1914。**同一** |

第 7 轮把 fail-closed 泛化到所有时间列、被测试拦住——**现行测试版本正确**，没有把 D-012 扩到普通列。

`D-020` 选项 A（一次性验证库、不碰生产 schema、R2 重定向到真实 conversion）是用户书面 P-004，足以解除 A-038「必须绑真实迁移则与 C2 门禁死锁」的结构性问题。本审接受该载体。

### F-I-002 能否闭合？

**不能。** A-038 点名的「用例仍是 ID；非法/越界可执行测试未发生」本审接受 **`fixed`**。整条仍 open。**全部剩余项**：

1. **冻结包负值政策不唯一（本轮新点名，required）**。测试与 Root D-012/D-015 已分清 voucher vs 普通列，但冻结包仍写全局「负值一律不进表达式 / 只由 m0 fail closed」：
   - `r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` L100：「负值 `< 0` **一律**不进表达式」——与同文件 §2.3「毫秒族负值安全」、`(col-999)/1000` 死代码化互相否定；
   - `r1-c2-per-table-rebuild-ddl-v1.0-fc.md` §0 同行；
   - `r1-c2-predicate-exact-sql-v1.0-fc.md` §1：「USING/rebuild 只承担 `=0→NULL` 与正值」写成全局（A-031 只把单路径落到 `#72/#73`）；
   - conversion contract §3.5「非法值单路径」同样全局化；
   - C3 边界 §5 断言 4 理由句「B 内不应再有负值」。
   R2 若按「一律 fail closed」落码会违反 Root D-015 与已绿测试。C2 冻结前必须把上述句子收成：**仅 voucher（`#72/#73`）负值 m0 fail closed；普通列负 epoch 走表达式（毫秒须 floor）。**
2. **freeze-candidate ≠ 实施**（边界，非新缺口；与 F-I-003/F-I-006 同一尺子，**单独不阻断**本条设计层闭合）。

不列为剩余：一次性库而非 90 列生产表（D-020 明文允许）；未跑 PG USING（D-020 载体是 SQLite 冻结表达式；PG 整数 interval 负值不需要 SQLite 那套 floor 修正）；`d0_positive` 未断言 `ms_d0`、`nullable_present` 未断言转换后字符串（卫生，不挡）。

---

## 第三块：recommended 复核

### F-I-025 · 计数 — **不能闭合**

rebuild 附件覆盖说明现为「15 个 descriptor，含 **20 张带时间列的表**（分母 90 列），加 3 张无时间列联接表 + 3 张跨 descriptor 子表」。

独立计数：

| 口径 | 事实 | 附件 |
|------|------|------|
| 90 列分母 | conversion contract `#1`–`#90` 连续；ledger 列分配 3+31+3+5+4+2+4+4+2+2+5+1+11+7+6=90 | **对** |
| 带时间列的表 | 从 90 行抽出 **44** 张唯一表（与 descriptor ledger 表范围 3+12+2+1+2+2+2+2+1+1+2+1+6+4+3=**44** 一致；A-034 原文即此数） | **仍写 20。不对** |
| 3 张无时间列联接表 | `user_roles` / `role_permissions` / `role_menu_items`（§1.1 B 组） | **对** |
| 3 张跨 descriptor 子表 | `notifications` / `user_mfa` / `mfa_proofs`（§1.1 C 组；时间列由 v80/v81 转） | **对** |

只把「20 张时间列表」改成「20 张**带时间列**的表」没有修掉 A-034/A-038 点名的 20 vs 44。C 组三张**已经在 44 张之内**，不能再加到 20 上凑数。

另：conversion contract §3.4 仍写「逐表 exact SQLite rebuild DDL **仍未写出**」（A-038 已记，陈旧）；§3.6 仍写 checksum 二选一「未选定」（与 child `D-017` / ledger §2 已裁决选项 A 矛盾）。不新开号，并入本条剩余。

### F-I-027 · 三源 — **闭合**

A-040 四项关闭要求：conversion §0、runbook superseded、Port L60/L71、边界 §4.1 Root/child 限定——均已兑现。三源不再互相否定。

全扫残留（不挡本条闭合、不新开号）：

- Port draft L14 仍写无限定 `User D-010`（Root D-010 = kernel 只暴露 `CreateRecoveryPoint`；child D-010 = `schema_migrations` owner）。内容同向，编号会误导；**不再**授权把 `snapshotBeforePending` 当 RecoveryPoint 源。
- 旧草案（guardrails / public-wire / allocation draft）仍有无限定 `D-0NN`；A-030 已接受点名五份冻结载体收口后旧草案不单独挡。
- C3 §7 / rebuild 的 `D-019`/`D-018` 无 Root 对号（仅 child），误读风险低。

---

## Findings

### F-I-002 · 逐列 USING / rebuild / codec

- **严重度**：high · **建议**：required
- **状态**：open（**再收窄，仍不关闭**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修**：A-038「用例仍是 ID；非法/越界可执行测试未发生」→ **`fixed`**（`D-020` + `w040contracttest` 三测试本机全绿；期望值独立核对通过；负值政策与 Root D-012 / D-015 同一）。
- **仍不闭合（全部剩余项）**：
  1. 冻结包仍全局化「负值一律 m0 fail closed、不进表达式」，与 Root D-012（voucher 专属）、Root D-015（负 epoch 合法 instant）及已绿测试矛盾。点名：mechanism L100、rebuild §0 负值行、exact SQL §1 全局句、conversion §3.5、C3 §5 断言 4 理由句。
  2. freeze-candidate ≠ 实施（边界；单独不阻断）。
- **关闭要求**：把上述全局句收成 voucher 专属 vs 普通列走表达式（毫秒 floor）；与测试及 Root D-012/D-015 同一。不要把测试当 C2 已实施。

### F-I-003 · 90 列 mapping

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-004 · Backup Port ≠ 可执行备份/回滚方案

- **严重度**：high · **建议**：required
- **状态**：**closed**（本审接受 `fixed`，设计层）
- **关闭证据**：A-040 §G 1–7 均已对位（G1 superseded+两处标注+正文保留；G2 三源收口；G3 PG 表与 `postgres.go` 一致且缺口属实；G4 A/C 分开、混合形状方向对；G5 四条可落码且与 `actionNoop` 现状一致；G6 六类+必须属于 X 且不属于 Y；G7 实测+记录类型身份可落码）。包路径候选维持 A-040 接受。实施在 R2。
- **残留（不重新打开）**：§2 批次末措辞偏宽；§3.1 C=`legacy` 不准；§4.3 fail-closed 用词；§5 理由句归 F-I-002。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮无新哈希/测试改写）

### F-I-006 … F-I-024 / F-I-026

- 维持既有 closed。

### F-I-025 · 冻结包卫生（计数）

- **严重度**：low · **建议**：recommended
- **状态**：open（**收窄，仍不关闭**）
- **本轮**：措辞加上「带时间列 / 90 列 / 3+3」。**仍 open**：20 ≠ 独立计数 **44** 张带时间列表（ledger 表范围与 conversion 90 行唯一表名一致）。conversion §3.4 / §3.6 陈旧句仍在。
- **关闭要求**：把「20 张」改成 **44 张带时间列的表**（或给出与 44 对得上的另一口径并证明）；划掉 §3.4「DDL 仍未写出」与 §3.6「checksum 二选一未选定」。

### F-I-027 · C3 落地后冻结包未收口

- **严重度**：med · **建议**：recommended
- **状态**：**closed**（本审接受 `fixed`）
- **关闭证据**：runbook `superseded`；Port L60/L71 与硬规则 1 同一；conversion §0 收回「尚未落盘」；边界 §4.1 已限定 Root D-007 + Root D-010。Port L14 无限定 `D-010` 为残留卫生，不再恢复三源矛盾。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮 |
|----|------|------------|------|
| F-I-002 | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec | **再收窄**：可执行测试 `fixed`。**仍缺**冻结包负值政策与 Root D-012/D-015/测试同一 |
| F-I-004 | C3、R2/R3 | — | **closed**（设计层） |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL | **未触及** |
| F-I-025 | 卫生 | 不单独挡冻结 | **仍 open**（20≠44） |
| F-I-027 | C3 卫生 | — | **closed** |

F-I-001、F-I-003、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015、F-I-016、F-I-017、F-I-018、F-I-019、F-I-006、F-I-020、F-I-021、F-I-022、F-I-023、F-I-024、F-I-026、**F-I-004**、**F-I-027** 为 closed。F-I-008、F-I-009、**F-I-025** 为 recommended open。

**开放 required = 2**（F-I-002、F-I-005）。在这些合法闭合前：不得冻结 C2、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。C3 **设计层**已可冻（F-I-004 closed）；Backup/RecoveryPoint **实施**仍在 R2，本条不授权提前改 `apps/` 生产代码。

## 与既有意见的异同

| 项 | A-040 independent | A-041 / E-035 自称 | A-042 independent（本条） |
|----|-------------------|--------------------|---------------------------|
| verdict | conditional；open required=3 | 不自证闭合；待本审 | **conditional**；open required=**2** |
| F-I-004 | 收窄、仍 open；§G 七项 | 七项已收口，待复审 | **closed**（设计层） |
| F-I-002 | 未触及 | 测试已绿，待复审 | 测试项 `fixed`；**整条仍 open**（负值政策不唯一） |
| F-I-005 | open | 结构性依赖 R2 | **仍 open** |
| F-I-025 | open（「20 张时间列表」） | 已改措辞 | **仍 open**（20≠44） |
| F-I-027 | 新增 recommended | 三源已修 | **closed** |
| 新 finding | F-I-027 | 无 | **无新 required / 无新 recommended 编号** |
| R2 | 禁止 | 禁止 | **禁止**（F-I-002/005） |

无「一要一否」需用户在 finding 之间裁。无需本轮新的 P-004（D-020 已落盘；负值政策以已有 Root D-012/D-015 为准，不是新裁决）。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | F-I-002 仍开（负值政策不唯一），阻断 C2 |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列分母不因本轮扩大；44 张表口径须写进 rebuild 覆盖说明 |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 设计层 closed；实施/harness 仍在 R2，信息项未关 |
| I-040-004 | required | R3 | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001 仍开放，阻断 C2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** `4462e73d` 把 A-040 §G 七项收成可核对的 C3 设计冻结，**F-I-004 closed**。`a2db84ae` 把 A-038 的可执行测试剩余做成可复跑的绿测，且负值政策在**测试里**是对的；**F-I-002 整条仍 open**，因为冻结包还在全局化 D-012。F-I-027 closed；F-I-025 因 20≠44 仍 open。开放 required **2**。

建议 `/govern`：

1. 响应本 A-042；**接受 F-I-004 / F-I-027 closed**；**不要**把 F-I-002 或 F-I-005 标 closed。
2. 最小文档收口（仍在 R1 设计面）：把 mechanism L100 / rebuild §0 / exact SQL §1 / conversion §3.5 / C3 §5 理由句改成 voucher 专属 vs 普通列走表达式；rebuild 覆盖说明 20→**44**；顺手划掉 conversion §3.4/§3.6 陈旧句与 Port L14 无限定 `D-010`。
3. **不要**冻结 C2，**不要**启动 R2 生产 schema/codec，**不要**改 formatter/历史 DDL。C3 设计层已可冻，实施仍等 R2。
4. 测试包按 D-020：冻结 DDL/表达式再变时必须同步更新；R2 把同一批用例重定向到真实 conversion。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree / `apps/`。响应、finding 闭合与是否推进由 `/govern` 处理。
