---
id: A-002-independent-checkpoints-b-c-d
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-003-r2-codec-and-descriptor-m1-m2 · 检查点 B/C/D + D-021 F-I-005 residual 复审触发 · commit c69ee93d
verdict: conditional
open_required: 1
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-003-r2-codec-and-descriptor-m1-m2
version: 0.1.0
---

# A-002 · independent · GOAL-003 检查点 B/C/D 与 F-I-005 residual 复审

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：execution-facts + finding-closure（F-I-005 residual 复审触发）· 检查点 B（15 个 SQLite `Apply` + 目录断言）/ C（PG 显式 DDL）/ D（真实 `MigrationChecksum` 与 `D-021` 复审）
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`id` 匹配；`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-003-r2-codec-and-descriptor-m1-m2`。同批 commit `c69ee93d` 的仓储改造见 sibling `GOAL-004` A-002；本条不改 GOAL-004 状态。
- **未读其他工作区作为审计上下文。** 未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree / 任何 Go 代码。
- 对照基线：本目标 A-001（self · `conditional`）；权威合同 = Root `D-004`–`D-017`、GOAL-002 `D-017`–`D-021` 与 C2 冻结附件；实现 = `c69ee93d`。
- 本轮可执行证据：`apps/api` 下 `go test ./internal/store -run TestCompiledMigrationCatalogOwnership|TestCompleteFingerprintTracksCatalogHead -count=1` → **ok**；`go test ./internal/temporal -count=1` → **ok**。未在本会话重跑全仓 `go test ./...` 或真实 PG 15.4（self 已报绿；PG 执行标 **unverified here**）。

## 成果（有证据）

1. **Checksum 算法与 D-017 选项 A 同文。** `kernel.MigrationChecksum`（`apps/api/kernel/persistence.go:14-28`）= `sha256(normalizeSQL(join(stmts,"\n"))+"\n"+transformID)`；`normalizeSQL` 逐行 TrimSpace、丢空行。v73–v87 的 `Checksum` 均取 `MigrationChecksum(vp040V*Statements(), transformID)`，`vp040V*Statements()` = `temporalmigrate.Ordered(preflight, rebuild, verify)`（m0 SELECT → m1–m3 字面 rebuild → m4 校验 SQL）。PG `ApplyPostgres` **不进哈希**；15 个 `transform_id` 均无 `:sqlite`/`:pg` 后缀，且与台账 §1 逐字相同（例 `0073:vp040-temporal-core-persistence:v1`）。
2. **m0/m4 进入 checksum 输入符合台账 §2 + D-017 括号定义，不是静默新约定。** D-017 L26「SQLite canonical DDL 切片」的括号把切片定义为 **Apply 语句序列、按 m0→m5**；台账 §2 把 m0 写成 preflight SELECT、m4 写成校验语句。把 m0/m4 排除出哈希会修订 D-017，从而触发 `D-021` 失效条件——本审**不**作该修订。m5（ledger INSERT）正确地未进切片（循环哈希）。
3. **15 个真实哈希已记录，且与编译期 descriptor 一致。** 冻结表 `migrate_test.go:768-782` 与 runtime `kernel.MigrationChecksum(...)` 由 `TestCompiledMigrationCatalogOwnership` 锁住；本轮该测试 **ok**。清单：

| v | transform_id | frozen checksum |
|--:|--------------|-----------------|
| 73 | `0073:vp040-temporal-core-persistence:v1` | `c2d2218e327b6ae859b897ed55ed01793822a6a6fdf52dc9deeb96c565920683` |
| 74 | `0074:vp040-temporal-authsession:v1` | `3ce1174a07300011f182baf93005f8fd2ce4b584f4b02d3f6cbcbc3e0106905c` |
| 75 | `0075:vp040-temporal-operationlog:v1` | `075d9f69a48c718872f1a3ed148b8b9c6e0da7e2d7fbf46edd9f9d026d4b1581` |
| 76 | `0076:vp040-temporal-jobs:v1` | `6b3649579cc6aedc713fab8c7f9f6dafbec548317f7395082d9e5ddb730aca56` |
| 77 | `0077:vp040-temporal-dictionary:v1` | `ae55a63a4366667f05410a76d15308504cf805540163e85e139b5737c44f3e73` |
| 78 | `0078:vp040-temporal-data-permission:v1` | `db2565e2da3ea3c5fe879bed5a9914a0cff3c0d0d3e59dafd9b68ba5b37055aa` |
| 79 | `0079:vp040-temporal-captcha:v1` | `493d66f70e904249d36bc593d070d0ebef0394a7b855d45a824e60a797599449` |
| 80 | `0080:vp040-temporal-mfa:v1` | `6537215f79af2ed02e4bbe2218733f8134f7af477b84fb50cc5b931c583762f8` |
| 81 | `0081:vp040-temporal-notifications:v1` | `a7565dc641f3c3291ff25cbefa52199ca06a8f94f7a978efac5b907615442b37` |
| 82 | `0082:vp040-temporal-recycle:v1` | `0132f6a873dd427b42c3a668bc88badc9b50a6c6729601e7cd1c5d5e4a1568fe` |
| 83 | `0083:vp040-temporal-scheduled-tasks:v1` | `e5df9a9e46bb6d8c259d8e134cb95bd1b7d043f9b108387800b23ee5c0ad6487` |
| 84 | `0084:vp040-temporal-settings:v1` | `bb3041a3d3fbeb5b3d706209f53cc578dc0e5d15016502919aac040b6bec2112` |
| 85 | `0085:vp040-temporal-wallet:v1` | `e1b5140669cfe3a7360787578a50d5a978c63f0b3b6ca0e1a1338ba5f33a8b60` |
| 86 | `0086:vp040-temporal-telegram:v1` | `80d5ad96381abd85cccf5022530e277b62c3333eb39aad5b5999754904a80973` |
| 87 | `0087:vp040-temporal-digital-offer:v1` | `31cef809f68578758bb2bf6544a158eaafaea191335b3a9e334933f4e4f56cc2` |

   附件 `r2-v73-v87-generated-statements-v0.1.md` 的「canonical」围栏是 **m1–m3**；m0/m4 以 Go struct 转储。第三方**不能**只对 markdown SQL 围栏做 sha256 得到上表。复算命令：`go test ./internal/store -run TestCompiledMigrationCatalogOwnership -count=1`。
4. **v1–v72 冻结哈希在 `c69ee93d` 未被改写。** `git diff c69ee93d^ c69ee93d -- apps/api/internal/store/migrate_test.go`：删除 4 行（`len==72` 尾断言）、新增 23 行（`len==87` 尾断言 + 15 条 v73–v87 行）。前 72 条 checksum 元组**不在 diff 中**。`sideTablesAbsent` 加在 `rebuildOperationLog` 的 Go 控制流上（`operationlog/migration/migration.go:785-788`），按 `D-019` §5 不进 0004–0071 的 `stmts` 切片。
5. **SQLite F-5 子女清单对 D-019 §2 完整。** v74 对 `users` 的 10 个子表（含跨 descriptor 的 `notifications`/`user_mfa`/`mfa_proofs` 与 FK-preserve 联接表）均 TEMP 备份 → DROP → 父表重建后再建回；跨 descriptor 三表在 v74 以 **INTEGER** 原样回填，v80/v81 再转 TEXT。`roles`/`permissions`/`menu_items`/`operation_log`/`dict_types`/`scheduled_tasks` 的子表均在对应 descriptor 内。v85 wallet 六表无 `REFERENCES` 子表，走裸四步。全部 `CREATE INDEX` 位于该 descriptor 最后一次 `DROP *_old`/`*_bak` 之后。
6. **CREATE TABLE 正文是 live `sqlite_master` + 仅逐列 INTEGER→TEXT（D0/voucher 去 NOT NULL/DEFAULT 0）。** 证据：`users` CREATE 保留 ALTER 追加的逗号粘贴列（`authsession/migration/vp040_temporal.go:77-85`）；`notifications` 保留 v37 的 `title_key`/`body_key`。未发现静默改列序/发明约束。
7. **PG 侧 v73+ 为显式 DDL，毫秒族为整数拆分式，无 `pgTimeColRe` 派生。** `pgTimeColRe`/`pgRebuild` 仅残留在 v1–v72 operationlog 路径（`D-019` 禁止改历史 DDL，允许）。秒族保留 E1 `to_timestamp(double)+date_trunc`（冻结裁决 B）。D0/voucher 列先 `DROP DEFAULT`/`DROP NOT NULL`。
8. **D0 / voucher / 负瞬间分层正确。** 仅 `vouchers.expires_at`/`redeemed_at` 设 `Voucher: true`（`wallet/migration/vp040_temporal.go:28-29`）；`RunPreflight` 只在 `Voucher && Negative>0` fail closed（`temporalmigrate.go:128-132`）。其余列负 epoch 走转换表达式（Root `D-015`）。sentinel `= 0 → NULL` 用于 #5/#6/#20/#34/#78；voucher 另含 `IS NULL OR = 0`。#34 `mail_config.updated_at` 的 ELSE 用毫秒整数拆分式，不是秒族。
9. **`D-019` §5 改动 3：v75 走通用字面执行器是合同允许的等价路径。** 合同原文是「必须调用 `rebuildOperationLogWithSessions`（**或与 F-5 同一通用 helper**）；禁止 `pgRebuild`」。v75 使用 `temporalmigrate.Exec` + 14 步 F-5 字面（先 DROP 两张子表再 rename），**未**调用 `pgRebuild`。现有 `rebuildOperationLogWithSessions` 的 copy 是 identity `SELECT created_at`（`migration.go:795-799`），接上会**不转换**毫秒列——冻结尾注已写「不可复用」。`sideTablesAbsent` 因此不在 v75 路径上：结构上子表在 rename 前已 DROP，与改动 1 要防的 CASCADE 隐患同构。本审判为**可接受等价**，不升 required。
10. **Codec 对拍基线存在。** `FromUnix`/`FromUnixMilli` 严格等于 `time.Unix`/`time.UnixMilli` UTC；`TestFromUnixMilliFloorSemantics` 钉死 `-1/-999/-1000/-1001/999 ms 无进位/公元 9999`（`internal/temporal/temporal_test.go:28-54`）。真实迁移边界矩阵用**同一组期望串**（`w040contracttest/migration_boundaries_test.go:147-153` vs codec 测试）。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| B：15 个 SQLite `Apply` + v1–v72 不变 | **部分（差 records 断言）** | 15 文件落码；冻结表前 72 行未改；缺台账点名的 `records` 不存在断言（F-I-001） |
| C：PG 显式 DDL、毫秒整数拆分、类型/精度断言 | **源码满足；PG 15.4 本轮未复跑** | 无 `pgTimeColRe`；E2 整数拆分；金额列拆分在 GOAL-004/`postgres_test.go` |
| D：真实 checksum 记录 + D-021 复审 | **哈希已记录且与 D-017 一致；residual 三项源码满足；本条即复审** | 见成果 1–3；F-I-005 ①②③ 可按 `fixed` 交 `/govern` 闭合 |
| `rebuildOperationLog` 改动 1–3 | **改动 1/2 落实；改动 3 等价** | `sideTablesAbsent`；v75 通用 F-5 执行器 |
| I-041-001 | **实现+D-001 已定稿；00-meta 仍写 collecting** | F-I-006 |

## F-I-005 residual 是否可闭合（问题 1，本条核心）

`D-021` 残余范围穷举三项。复审触发 =「R2 首次记录任一 v73+ 哈希」——**已触发**（本批）。D-017 约定**未被修订**。

| # | 残余项 | 本审判定 |
|--:|--------|----------|
| ① | 15 个真实 `MigrationChecksum` | **满足。** 算法 = D-017；输入 = Ordered(m0→m4)；PG 不进哈希；冻结表与编译期 descriptor 由 `TestCompiledMigrationCatalogOwnership` 锁住（本轮 ok） |
| ② | `migrate_test.go`/`postgres_test.go` v73+ 追加 + 金额列拆分 | **源码满足。** 冻结表 15 行；`postgres_test.go:307-314` 断言 `wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 保持 `bigint` |
| ③ | leftover 21 名补入 PG 断言 | **源码满足。** `postgres_test.go:319-325` 的 21 名 = allocation L38 同一集合（顺序不同）；另有 `datetime_precision=6` 断言 |

**结论：F-I-005 的 `accepted-residual` 范围内三项现已有可核对实现证据，本审判定可走 `fixed` 闭合。** 本意见**不**改 GOAL-002 台账；闭合由 `/govern` 落盘。不得把本条读成「R2 已放行」或「M4/Backup Port 已完成」。

若编排器/用户坚持「DDL 切片」= 排除 SELECT/PRAGMA，则 ① 失败、15 个哈希须重算、且构成对 D-017 的修订（`D-021` 失效条件）。那是 **P-004**，本审不代选。

## Findings

### F-I-001 · v73 未实现台账点名的「retired `records` 不存在」断言

- **严重度**：med
- **建议**：required
- **状态**：open
- **影响门禁**：检查点 B/D；台账 §1 v73 表范围；allocation「assert retired records absent」
- **描述**：`r1-c2-descriptor-ledger-v1.0-fc.md` L22 把 v73 表范围写成「`schema_migrations`、`mail_outbox`、`mail_config`；**断言 retired `records` 不存在**」。`corepersistence/migration/vp040_temporal.go` 的 preflight/verify/rebuild **零处**出现 `records`。正常 catalog 路径上 v6 `DROP TABLE IF EXISTS records` 已执行，故新鲜库不会留下该表；缺的是冻结尾注要求的 fail-closed 守卫——若异常路径把 `records` 带回（跳过 v6、错误 restore），v73 不会拒绝，`records.updated_at` 会保持 INTEGER。
- **证据**：`vp040_temporal.go` 全文无 `records`；v6 退休 DDL 在 `corepersistence/migration/migration.go:40-41`。
- **反例 / 关闭要求**：在 v73 SQLite m0 或 m4 增加 `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='records'` 且要求 0（PG 对 `information_schema.tables` 同类断言）；重新计算 v73 checksum（哈希输入将变）。关闭前不得声称台账 §1 表范围已逐字落实。

### F-I-002 · 哈希中的 m4 PRAGMA 基数与 `RunVerify` 执行不一致

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`VerifyStatements`（`temporalmigrate.go:87-108`）在**每张**被转换表循环内追加 `PRAGMA foreign_key_check` 与 `PRAGMA integrity_check`。`RunVerify`（L194-197）在循环**之后**各执行一次。checksum 覆盖的 SQL 文本 ≠ 实际执行次数。不改变转换语义；削弱「hashed SQL == executed SQL」口号。
- **关闭要求**：让 `VerifyStatements` 与 `RunVerify` 同构（每表一次或全局一次），并重算受影响 checksum；或书面接受「校验 PRAGMA 去重执行」。

### F-I-003 · D-018 的 T-#-RT 族是抽样硬编码，不是同进程 codec↔SQL 对拍

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：D-018 要求把「Go codec `FromUnix`/`FromUnixMilli` 与 SQL 结果逐行等价」写成 R2 强制验收（`T-<#>-RT`）。逐列合同给了 T-1-RT…T-90-RT。落地是：(a) codec 单测钉死期望串；(b) `migration_boundaries_test.go` 对**抽样列**用**同一组字面期望串**断言真实迁移。两文件没有互相调用。D-020 §2 把 R2 重定向范围收成边界矩阵（负值/999 ms/历元 0/9999/NULL/sentinel/voucher/词法序）——该矩阵**已**接到真实 descriptor。90 列全覆盖与「同一次调用里 codec()==SQL」仍缺。
- **证据**：`w040contracttest` 无 `temporal.` 导入；jobs 期望串与 `temporal_test.go:33-42` 相同。
- **关闭要求**：至少对秒族与毫秒族各取若干列，在同一测试里 `MustFormat(FromUnix/FromUnixMilli(x))` vs 真实迁移扫描结果；或用户书面接受「共享夹具 + D-020 边界矩阵」为 D-018 的有界实现。

### F-I-004 · 生成器留在生产树（`VP040_GENERATE=1` 才跑）

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`internal/store/vp040_generate_test.go` 会写 `modules/*/migration/vp040_temporal.go`。默认 skip，可复现性是优点。风险是「测试写源码」工具误开环境变量会改 checksum 输入。非门禁。

### F-I-005 · v86 `ModuleID` 台账写 `admin.channel.telegram`，代码用 `channel.telegram`

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：台账 §1 的 ModuleID 列**不在** §4「R2 唯一允许输入」（name / transform_id / 表范围 / 列分配）里。运行时 `ModuleID = "channel.telegram"`（`modules/channel/telegram/migration/provider.go:10`）与冻结 catalog 行 `migrate_test.go:781` 一致。checksum 不含 ModuleID。属展示层更正，应回写台账以免对照漂移。

### F-I-006 · `00-meta` 仍把 I-041-001 标为 collecting

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`01-decision.md` 与 `D-001` 已把 I-041-001 标 verified；`00-meta.md` L67 与 `03-audit.md` 信息就绪表仍写 collecting。过程不一致，不否定 codec 已落码。由编排器同步元数据（本审不改 status/progress）。

## 必改项汇总

1. **F-I-001（required）**：v73 补 `records` 不存在断言，并重算 v73 checksum。

## 与 A-001（self）的异同

| 项 | self A-001 | 本条 |
|----|------------|------|
| 改动 3 / v75 通用执行器 | 待判 | **可接受等价**（合同有「或同一通用 helper」；现有 WithSessions copy 不能转换） |
| m0/m4 进哈希 | 待判 | **符合**台账 §2 / D-017 括号定义 |
| 生成器留树 | 待判 | recommended F-I-004 |
| F-I-005 residual | 不自证闭合 | **①②③ 可 `fixed`**，交 `/govern` |
| `records` 断言 | 未提 | **新 required F-I-001** |
| 逐行 15 文件 | 自承未逐行 | 本条核对了 F-5 子女、索引位置、PG 显式 DDL、D0/voucher 分层；未对 15 个 CREATE 与 live sqlite_master 做逐字节差（生成器自检存在，本轮未重跑 `VP040_GENERATE=1`） |

## 结论 + 建议给编排器/用户的下一步

检查点 B/C 的主体实现与 checksum 约定**可核对**；D-021 residual 的复审触发已满足，三项残余**可以**按 `fixed` 闭合。开放 required = **F-I-001**（v73 `records` 断言）。在该条闭合或用户书面 `accepted-residual`/`user-overruled` 之前，不得把 GOAL-003 标 `done`，也不得把「台账 §1 已逐字落地」写成事实。

建议 `/govern` 下一句：响应 GOAL-003 A-002（闭合 F-I-001 或书面 residual）并同步把 F-I-005 三项从 residual 改为 `fixed`；GOAL-004 的 PG 截断 required 见该目标 A-002。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
