---
id: A-034-r1-independent-e030-per-table-rebuild-ddl
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · ad-hoc design-plan + finding-closure · E-030 / commit 726f62c1 逐表 exact rebuild DDL 正文对照 A-032 F-I-002.1 / F-I-021 / F-I-022 · 不是实施审计 · freeze-candidate ≠ 已实施
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-034 · R1 independent · E-030 逐表 exact rebuild DDL

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：ad-hoc + finding-closure（用户指定核对本轮 `726f62c1` / E-030 新落盘的逐表 exact rebuild DDL 是否构成对 A-032 三条 required 的合法闭合、收窄或不足：F-I-002.1 / F-I-021 / F-I-022。对照 A-032，当时 open required = 6。用户已 P-004 裁决 F-5 + 两次重建，落盘 child `D-019`。附件为设计材料，不是实施证据）。
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree / `apps/`。
- 实测：本机 `sqlite3` 3.51.2，全部在 `:memory:`（脚本暂存 `%TEMP%\a034-fk-audit\`），未触碰仓库数据。
- 本条**不**复审 A-031 对 F-I-006 谓词方向的修正；F-I-006 维持 open。F-I-004 / F-I-005 本轮未触及。

## 核对方法

1. 通读 E-030、`r1-c2-per-table-rebuild-ddl-v1.0-fc.md`、`D-019`、`D-014` / descriptor ledger、mechanism §2、A-032/A-033。
2. 对位 `apps/` 现行源码：authsession / account / notifications / corepersistence / jobs / datadictionary / datapermission / settings / telegram / wallet / operationlog / mfa / identity.go / migrate.go / postgres.go / `kernel.MigrationChecksum`。
3. 独立用 sqlite 3.51.2 复现：live 列序（users 17 列、mail_outbox、telegram_config、dict_entries、site_settings 版本序 vs 引用序）、完整子表集合、§1.2 v74 F-5 整链、v81 二次重建、D0 / voucher 包裹。
4. `git show --stat 726f62c1`：3 个文件均在 workspace-040；`apps/` 无变更。

## 成果（有证据）

1. **本轮未实施 DDL/codec。** commit `726f62c1` `--stat` 为 `02-execution.md` + `E-030` + 本附件；`apps/` 无变更。附件 `status: freeze-candidate`。E-030 自承未闭合任何 required。本审同意该实施边界。
2. **F-I-021 关闭要求 1–5 均已满足**（见 C）。12 张子表是 `users`/`roles`/`permissions`/`menu_items` 的完整集合，无第 13 张。§1.2 批量子女先行是 `D-019` §1 在四父表交织下的正确展开。本机 F-5 复现与 E-030 自测一致。
3. **F-I-022 关闭要求已满足**（见 D）。5 处生产写入路径完整准确；无第六处生产写入 `schema_migrations.applied_at`。v1 `schemaMigrationsDDL` 禁止改、restore 字面不进 `0001:r2-baseline` 哈希输入，两条边界成立。
4. **F-I-002.1 仍不能闭合**（见 A）。v73–v87 的 **v78 整 descriptor 缺席**；`site_settings` 15 列 new CREATE 未写且 ALTER 引用序 ≠ live 列序。表达式层与 F-5 编排层已收窄。
5. **§7 append-only 论证成立**（见 E）。`MigrationChecksum` 哈希的是规范化 `stmts` 切片 + `transform_id`，不是 Go `Apply` 函数体。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C2 物理合同 / 逐表 SQLite rebuild | **仍不可冻结** | F-I-002.1 仍缺 v78 正文 + `site_settings` live CREATE；PG 显式 DDL 仍缺 |
| F-I-021 FK 处置 | **本审接受 closed** | `D-019` P-004 + 完整 12 子表 + 本机 F-5 复现 |
| F-I-022 v73 写入路径 | **本审接受 closed** | 附件 §6 五处与源码对位 |
| C3 / C4 / R2 放行 | **未满足** | F-I-004/005/006 未触及；本轮 freeze-candidate ≠ 实施 |

## 对用户 A–F 的直接判定

### A. 主交付物 `r1-c2-per-table-rebuild-ddl-v1.0-fc.md`

#### A1. 是否覆盖 `D-014` 的 15 个 descriptor / 时间列表

**未覆盖完整。** Descriptor ledger（`r1-c2-descriptor-ledger-v1.0-fc.md` §1）与 `D-014` 接受的 allocation 是 **15 个 descriptor、44 张时间列表**（3+12+2+1+2+2+2+2+1+1+2+1+6+4+3）。附件自称「15 个 descriptor / 20 张时间列表」——**「20 张」计数不成立**，且 **v78 整段缺席**。

| v | 表 | 本附件 |
|--:|----|--------|
| 73 | `schema_migrations`, `mail_outbox`, `mail_config` | 有（§2.1–2.3） |
| 74 | 12 张转换表 + 3 张 FK-preserve + 3 张跨 descriptor 子表 | 有（§1 / §3） |
| 75 | `operation_log`, `operation_log_archive` | 有（§4；event 枚举走源码行号） |
| 76 | `jobs` | 有（§2.4；CREATE 省略、INSERT/INDEX 写出） |
| 77 | `dict_types`, `dict_entries` | 有（§2.5；new CREATE 省略） |
| **78** | **`data_scope_policies`, `user_data_scopes`（#50/#51）** | **无。缺失，不是多余** |
| 79–87 | captcha / mfa / notifications / recycle / scheduled-tasks / settings / wallet / telegram / digital-offer | 有（§2.6–2.11 / §5） |

多余：无。FK-preserve 三张联接表与 C 组三张跨 descriptor 子表是 `D-019` 要求的附加，不是分母扩大。

`user_data_scopes.user_id` **没有** `REFERENCES users`（`datapermission/migration/migration.go:27-33`），故 v78 不是第 13 张 FK 子表；缺的是 **时间列转换 DDL**，不是 F-5 清单。

#### A2. §2/§3 DDL 与 `apps/` 抽核

本机按全局版本序 apply ALTER 后 `PRAGMA table_info` / `sqlite_master.sql`（3.51.2 **会把 ALTER 列折进 `sql`**）：

| 抽核项 | 判定 |
|--------|------|
| `users` 17 列列序 | **一致**。cid 0–16 = `id, username, name, roles, password_hash, created_at, updated_at, token_version, failed_login_count, locked_until, enabled, notifications_enabled, avatar_url, must_change_password, email, email_status, last_login_failure_at`。与 §3.2 逐字同序。权威路径 = v1 CREATE → v11 → v12（2 列）→ v13 `enabled` → v17 `notifications_enabled` → v35 `avatar_url` → v38 → v54（2 列）→ v61 |
| `jobs` 六态 CHECK | **源码一致**（`jobs/migration/migration.go:18,36-43`：`queued/running/succeeded/failed/cancelled/expired`）。附件未重抄 CHECK 正文，只写「逐字保留」+ INSERT 五列包裹。INSERT 可空性与 conversion contract #41–#45 一致（`lease_expires_at`/`finished_at`/`expires_at` NULL 包裹；`created_at`/`updated_at` NN）。四索引含 v72 `idx_jobs_created_at`（`:127`，`IF NOT EXISTS`） |
| `dict_entries.badge_style` | **列序对**：live cid 9 = `badge_style`（v39 ALTER 在末尾）；`sqlite_master.sql` 已折入。INSERT 把 `badge_style` 放最后。**new CREATE 仍是省略号** |
| `telegram_config` 两个 v0067 列 | **一致**。live = `id, bot_token_enc, webhook_secret_enc, updated_at, mode, webhook_public_base_url`（`:29-30`，checksum `0067:telegram-config-connection:v1`）。§5.5 legacy CREATE 列序正确 |
| `mail_outbox` 两个 v0060 列 | **一致**。live = `id, to_addr, subject, body, created_at, channel, delivery_status`（`:82-84`，`0060:mail-outbox-channels:v1`）。§2.2 new CREATE 写出 |
| `site_settings` 15 列 | **new CREATE 未写**。按版本序 live cid = `id, site_title, logo_url, updated_at, logo_url_light, logo_url_dark, favicon_url, default_locale, site_timezone, default_theme, copyright_text, icp_number, operation_log_retention_days, operation_log_expiration_action, default_currency`。§2.11 引用 `:86-91 / :121-122 / :207 / :222-223` 把 **v62 `default_currency`（:207）写在 v46 retention（:222）之前**。本机按该引用序 apply 后 `default_currency` 落在 cid 12、retention 在 13–14，**与 live 不一致**。这不是「漏列」，是「按引用序拼会拼错列序」 |

#### A3. 表达式与三种包裹

**与 mechanism §2 及 A-032 复证一致。** §0 秒族 / 毫秒族字面与 mechanism §2.2/§2.3 同一；`%f` / `/1000.0` 禁用维持；负值 `< 0` 不进表达式。

| 包裹 | 本附件使用 | 判定 |
|------|------------|------|
| NN 直出 | `schema_migrations.applied_at`、`mail_outbox.created_at`、`users.created_at/updated_at` 等 | 对 |
| NULL | `jobs.lease_expires_at/finished_at/expires_at`、`notifications.read_at`、`task_runs.finished_at`、`refresh_tokens.revoked_at` | 对 |
| D0 | `mail_config.updated_at`（毫秒）、`users.locked_until` / `last_login_failure_at`、`login_failures.locked_until`、`telegram_config.updated_at`（秒）；新 DDL 去 `NOT NULL`/`DEFAULT 0` | 对。本机：`locked_until=0` → NULL；`1758320002` → `2025-09-19T22:13:22.000000Z` |
| voucher | §5.4 写明 `IS NULL OR = 0`；本机三行探针 NULL/0 → NULL、正秒 → 27 字符 TEXT | 对（正文未展开 INSERT，见 A5） |

`#61 task_runs.finished_at` 是 **N + 运行时去掉 `COALESCE(…,0)`**，不是 D0 列；附件用 NULL 包裹，正确。

#### A4. `<秒表达式(col)>` / `<毫秒表达式(col)>` 占位

**可接受为 exact DDL 的占位写法，不必逐处展开。** 展开式在同文件 §0 唯一冻结、与 mechanism §2 同一；使用处标明了族（秒/毫秒）与包裹。逐处展开只会增加转录漂移。R2 落码必须机械替换，不得改写。

#### A5. 源码行号 + 差异说明是否足以闭合 F-I-002.1

**不足以闭合。** 允许的例外与仍须重抄的分界：

**允许源码行号 + 差异（live 形状 = 单一现行 contiguous CREATE，差异仅为时间列 INTEGER→TEXT 及 D0 空性）：**

- `operation_log.event` 超长枚举：钉死 `operationLogDigitalOfferDDL[0]`（`operationlog/migration/migration.go:513-521`），R2 **逐字节**复制 CHECK；重抄 60+ 事件名风险更高。
- 已被 v33/v64 重建的 wallet 表：`wallet_ledger_entries` `:145-162`、`wallet_accounts` `:304-318`（及 `subjects` `:279-285`、`vouchers` `:287-301`、`voucher_batches` `:391-395`、`wallet_reconciliation_runs` `:56-64`）。
- `jobs` 六态 CHECK：钉死 `:15-44`。

**必须写出完整 new `CREATE TABLE`（live 形状由 CREATE + 后续 ALTER 拼出，或本附件尚未给出任何 CREATE）：**

1. **`data_scope_policies` / `user_data_scopes`（v78）** — 目前完全没有（F-I-023）。
2. **`site_settings` 15 列 new CREATE** — 必须按上表 live cid 顺序写全；禁止按 §2.11 文件行号序拼 ALTER（F-I-024）。
3. **`dict_entries` new CREATE** — 须含末列 `badge_style TEXT NOT NULL DEFAULT 'default'`（与 INSERT 已写的列清单对齐）。legacy+ALTER 已在 §2.5 给出，缺的是可粘贴的 new 正文。

`roles`/`permissions`/`menu_items`/`system_data_reconcile`/A 组挑战表等单一源码 CREATE，维持「行号 + 只改时间列」可接受。

### B. §1.2 v74 F-5 编排 vs `D-019`；独立复现

**与 `D-019` §1/§2/§3 一致。** 四父表子表交织，必须先摘全部子表、再重建全部父表、最后统一建回——这是对逐父表十步的正确批量化，不是偏离。C 组时间列保持 INTEGER、v80/v81 再转，与用户裁决「两次重建」同文。

本机 sqlite 3.51.2、`PRAGMA foreign_keys=ON`、单事务 F-5 后 v81 裸四步：

| 检查项 | 本审结果 |
|--------|----------|
| v74 后 `foreign_key_check` | **0 行** |
| 行数 users/roles/user_roles/notifications/user_mfa | **2 / 1 / 1 / 1 / 1**（另：refresh_tokens/mfa_proofs/user_invites 各 1，无丢行） |
| `users.locked_until` D0 | `0` → NULL；`1758320002` → `2025-09-19T22:13:22.000000Z` |
| C 组 `notifications.created_at` | 仍为 **integer** `1758320000` |
| 子表 `REFERENCES` | 指向 `users`/`roles`/`permissions`/`menu_items`，`instr(sql,'_old')=0` |
| v81 二次重建 | `fk_check=0`；`created_at` 27 字符 TEXT；FK 仍指向 `users(id)`；`integrity_check=ok` |

**12 张子表是完整集合，无第 13 张。** 本机对 `sqlite_master.sql LIKE '%REFERENCES <parent>%'`：

| 父表 | 子表 |
|------|------|
| `users`（10） | `refresh_tokens`, `user_roles`, `email_verification_challenges`, `password_recovery_challenges`, `login_failures`, `user_password_history`, `user_invites`, `notifications`, `user_mfa`, `mfa_proofs` |
| `roles`（3） | `user_roles`, `role_permissions`, `role_menu_items` |
| `permissions`（1） | `role_permissions` |
| `menu_items`（1） | `role_menu_items` |

去重 = **12**。显式排除：`service_credentials.created_by` 无 FK；`user_data_scopes.user_id` 无 FK；`password_policy` / `system_data_grants` 无 FK。authsession `:355` 等是 **PG 镜像**，不是第 13 张 SQLite 表。

轻微文档不一致（不阻断 F-I-021）：§1.2「阶段 6 全部索引最后建」vs §3.2 在 `DROP users_old` 后立刻建 `idx_users_email_lower`。只要遵守「索引在 DROP `_old` 之后、且只建一次」，引擎层安全。

### C. F-I-021 关闭要求 1–5

| # | 要求 | 判定 |
|--:|------|------|
| 1 | P-004 选定模式 | **满足。** `D-019` `status: accepted`：F-5 子女先行 + TEMP；不采用 runner 12 步 / 合并 descriptor |
| 2 | 完整子表清单 + live DDL 权威 | **满足。** `D-019` §2 + 附件 §1.1/§3.10/§3.11。12 张无漏。B/C 组 CREATE 与本机 folded `sqlite_master` 一致（含 `notifications.title_key/body_key`）。权威声明为 v72 `sqlite_master`；本审按源码 CREATE+ALTER 复现，未跑 `OpenWithCatalog`，重建结果与声明一致 |
| 3 | v74 可 F-5 非转换表且不抢 v80/v81 | **满足。** §1.3 / `D-019` §3；C 组回填 INTEGER；本机证实 |
| 4 | 删除机制附件错误句 | **满足。** mechanism §1 已改写（错误句删除，裸四步仅无 FK 子表）。`D-018` **理由句**已作废。残留：`D-018`「决定的来源」第 1 点仍写「含被 FK 引用的表」（F-I-025 recommended） |
| 5 | `rebuildOperationLog` fail-closed 断言进冻结包 | **满足。** 附件 §7 改动 1–3；实施在 R2 |

**F-I-021：本审接受 closed。**

### D. F-I-022 五处写入路径

**完整且准确。无第六处生产写入。**

| # | 位置 | 源码核对 |
|--:|------|----------|
| 1 | `identity.go:58` `sqliteLedgerDDL` | `applied_at INTEGER NOT NULL` |
| 2 | `identity.go:65` `postgresLedgerDDL` | `applied_at BIGINT NOT NULL` |
| 3 | `identity.go:317,319-321` `stampCatalog` | `now := time.Now().UTC().Unix()`；sqlite **与** PG `restoreLedger`（`:365-371` / `:390-395`）共用这一处 INSERT |
| 4 | `migrate.go:121-123` `applyMigration` | `time.Now().UTC().Unix()` |
| 5 | `postgres.go:165-167` `applyMigrationPG` | 同上 |

- `system_data_reconcile.applied_at` 是 **另一张表**（v74 `#2`），不是 ledger。
- 测试夹具（`migrate_test.go:348` 等）写入整数 `1`，属 R2 测试改写，**不是**第六处生产路径。
- **v1 `schemaMigrationsDDL`（`authsession/migration/migration.go:18-23`）禁止改**：属 v1 Apply 正文。即使改了也不会进入 v1 checksum（见 E），所以禁改是对的。
- **restore 字面不进 `0001:r2-baseline` 哈希**：`Descriptors()` 用 `kernel.MigrationChecksum(r2BaselineDDL, "0001:r2-baseline:v1")`；`r2BaselineDDL` 只有 `users` + `refresh_tokens` + 一个 INDEX（`:25-44`），**不含** `schemaMigrationsDDL`，也不含 `identity.go` 两个 restore 字面。边界成立。

**F-I-022：本审接受 closed**（设计层；代码仍在 R2）。

### E. §7 append-only 边界

**成立。** `apps/api/kernel/persistence.go:14-17`：

```text
input = normalizeSQL(strings.Join(stmts, "\n")) + "\n" + transformID
```

`normalizeSQL` 只 TrimSpace / 丢空行。哈希输入 = **DDL 切片字面 + transform_id**。改 `rebuildOperationLog` 的 Go 断言 / 控制流 **不进** checksum。禁止改 `0004–0071` 切片字面与 `transform_id` → 0001–0072 canonical SQL/checksum 不变。与 `D-017` 单 checksum（PG 不进哈希）相容。F-I-005 的「已记录哈希值 / 可执行测试改写」仍开放，本轮无新哈希。

### F. 本轮新 finding；ALTER 列序；wallet 是否需 F-5

- `mail_outbox` / `telegram_config` 追加列列序 **以 live 为准写对**。
- `dict_entries.badge_style` INSERT 列序对；new CREATE 未写（见 A5.3）。
- `site_settings` **未以 live 为准写出 15 列 CREATE**，且引用序会导错列序（F-I-024）。
- **wallet 当前无需 F-5。** 全仓 `modules/wallet/**/*.go` 无任何 `REFERENCES wallet_*`。`account_id` / `batch_id` 都是无 FK 的 TEXT。§5.8「落码前再盘点」作为 R2 前置足够；本审已做源码负向盘点。建议把「现行 v72 源码无 `REFERENCES wallet_*`」写成冻结事实，而不仅是过程句。

新 finding：F-I-023 required（v78 缺席）、F-I-024 required（`site_settings` live CREATE / 引用序陷阱）、F-I-025 recommended（卫生）。

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-002 · 逐列 USING/rebuild 仍不足以为 C2 冻结

- **严重度**：high · **建议**：required
- **状态**：open（维持；**本轮收窄 F-5 编排 + 多数表 INSERT，不关闭 F-I-002.1**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修 / 收窄**：
  1. 逐表 SQLite rebuild 正文已落盘为 freeze-candidate（15 descriptor 中 14 个有载体）。
  2. 表达式占位合法；D0 / NULL / voucher 在已写出的 INSERT 上用对。
  3. 父表走已裁决的 F-5（F-I-021 closed）。
- **仍不闭合**：
  1. **v78 两表无任何 CREATE/INSERT**（升级为 F-I-023）。
  2. **`site_settings` 15 列 new CREATE 未写**，引用序 ≠ live（升级为 F-I-024）。
  3. **`dict_entries` new CREATE 仍省略**（ALTER 拼出的 live 形状必须可粘贴）。
  4. **PG 侧显式 DDL 仍未落盘**（`D-019` §6 / A-032 E；本附件只说「由配套附件承载」，仓库无该附件）。
  5. 用例仍是 ID；非法/越界可执行测试未发生。
- **关闭要求**：补 A5 必抄三项 + PG 显式 DDL（禁止 `pgTimeColRe`）。源码行号例外维持 A5 允许清单。freeze-candidate ≠ 实施。

### F-I-003 · 90 列 mapping

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-004 · Backup Port ≠ 可执行备份/回滚

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮未触及 C3）

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮无新哈希/测试）
- **本轮相关**：§7 边界论证成立，**收窄「改 Go 会不会破坏 0001–0072 checksum」**；仍缺已记录哈希值与可执行测试改写。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med · **建议**：required · **状态**：open（维持）
- **本轮**：不复审 A-031 对 `#5` 方向 / readwrite superseded / ORDER BY 归属的修正。

### F-I-007 … F-I-020

- 维持 A-032/A-030 状态。F-I-018 / F-I-019 **closed** 维持。F-I-008 / F-I-009 / F-I-020 recommended 仍 open。

### F-I-021 · FK 父表重建阻塞：处置与子表清单未冻结

- **严重度**：high · **建议**：required
- **状态**：**closed**（本审接受 `fixed`）
- **关闭证据**：`D-019`（P-004 F-5 + 两次重建）+ 附件 §1.2/§3.10/§3.11 + 本机 3.51.2 复现（fk_check=0、无丢行、D0、C 组 INTEGER、无 `_old`、v81 `integrity_check=ok`）+ 12 子表完整（无第 13 张）+ mechanism §1 错误句已删 + §7 断言设计。实施仍在 R2。

### F-I-022 · v73 ledger 写入路径不只 identity.go 两个字面

- **严重度**：med · **建议**：required
- **状态**：**closed**（本审接受 `fixed`，设计层）
- **关闭证据**：附件 §6 五处与现行源码对位；无第六处生产写入；v1 `schemaMigrationsDDL` 禁改；restore 字面不在 `r2BaselineDDL` 哈希输入内。代码改动在 R2。

### F-I-023 · v78 `data_scope_policies` / `user_data_scopes` 无逐表 rebuild DDL

- **严重度**：high · **建议**：required · **状态**：open（本轮新增）
- **影响门禁**：C2 冻结、F-I-002.1；关联 `I-040-001`
- **描述**：`D-014` / descriptor ledger 明确 v78 = `admin.data-permission`、表 `data_scope_policies`/`user_data_scopes`、列 #50/#51 `S-NN`。E-030 声称「v73–v87 全部 15 个 descriptor」。附件无 §、无 CREATE、无 INSERT。源码 `datapermission/migration/migration.go:20-33` 两表均无 FK，应走裸四步，`updated_at INTEGER NOT NULL` → `TEXT NOT NULL` + 秒族直出。
- **关闭要求**：补这两表的 exact `CREATE` / `INSERT … SELECT`（无索引）。不要把它们误加入 v74 F-5 子表清单。

### F-I-024 · `site_settings` 15 列 new CREATE 未冻结，ALTER 引用序 ≠ live 列序

- **严重度**：high · **建议**：required · **状态**：open（本轮新增）
- **影响门禁**：C2 冻结、F-I-002.1；关联 `I-040-001`
- **描述**：live 15 列由 v7 CREATE + v10(6) + v40(2) + v46(2) + v62(1) 拼出，`default_currency` 在最后。§2.11 只给源码行号，且把 `:207`（v62）写在 `:222-223`（v46）之前。本机按引用序 apply 得到错误 cid。R2 若按该序拼 CREATE/`INSERT SELECT` 会错列。Go seeder（`:69-72` / `:52-55`）须与 v84 同批改格式的要求本身正确。
- **关闭要求**：写出 15 列 new `CREATE TABLE site_settings`（仅 `updated_at TEXT NOT NULL`，其余列与约束按 live cid）+ 对应 `INSERT … SELECT` + `CREATE INDEX IF NOT EXISTS idx_site_settings_updated_at`。列序以 v72 `PRAGMA table_info` 为准，不以文件行号为准。

### F-I-025 · 冻结包卫生（dict_entries new CREATE / D-018 来源句 / 计数）

- **严重度**：low · **建议**：recommended · **状态**：open（本轮新增）
- **描述**：
  1. `dict_entries` new CREATE 仍省略；INSERT 已含 `badge_style`，漏 CREATE 会在 R2 直接 fail，但仍应可粘贴。
  2. `D-018`「决定的来源」第 1 点仍写既有 rebuild「含被 FK 引用的表」；理由节已作废，来源节未划掉。
  3. 附件自称「20 张时间列表」与 ledger 的 44 张不符。
  4. §1.2 与 §3.2 的索引时机表述不统一。
- **关闭要求**：补 `dict_entries` new CREATE；划掉 `D-018` 来源中的过时句；改计数；统一「索引在 DROP `_old` 之后且只建一次」。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮 |
|----|------|------------|------|
| F-I-002 | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec | **收窄**：F-5 编排 + 14/15 descriptor 载体 + 表达式占位合法。**仍缺** v78、site_settings live CREATE、dict_entries new CREATE、PG 显式 DDL |
| F-I-004 | C3、R2/R3 | 不得把 Port/runbook 当 C3 冻结 | **无新收窄** |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL | **收窄**：§7 哈希输入核对成立；仍无已记录哈希 |
| F-I-006 | C2/C3、R2 | 不得在 exact old/new 未冻时 table-rebuild | **本轮未复审 A-031** |
| F-I-021 | C2、R2 | — | **closed** |
| F-I-022 | C2 v73、R2 | — | **closed**（设计层） |
| **F-I-023** | C2、F-I-002.1 | 不得声称 15 descriptor 已写完 | **新增**；v78 两表缺席 |
| **F-I-024** | C2、F-I-002.1 | 不得按 §2.11 行号序拼 `site_settings` | **新增**；15 列 live CREATE 未写 |

F-I-001、F-I-003、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015、F-I-016、F-I-017、F-I-018、F-I-019、**F-I-021**、**F-I-022** 为 closed。F-I-008、F-I-009、F-I-020、**F-I-025** 为 recommended open。

**开放 required = 6**（F-I-002、F-I-004、F-I-005、F-I-006、**F-I-023**、**F-I-024**）。在这些合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-032 independent | E-030 / 附件自称 | A-034 independent（本条） |
|----|-------------------|------------------|---------------------------|
| verdict | conditional | 未自证闭合；待本审 | **conditional** |
| F-I-002.1 | 未落盘 | 已落盘；§4/§5 用行号代替重抄 | **仍 open**；14/15 有载体；v78 与 site_settings 阻断 |
| F-I-021 | open；待 P-004 | D-019 已裁决；待本审 | **closed** |
| F-I-022 | open | §6 列出 5 处；待本审 | **closed**（设计层） |
| F-5 复现 | 单父+单子可行 | 四父+A/B/C 探针通过 | **独立复现通过**；12 子表完整 |
| 新 required | F-I-021/022 | 无 | **F-I-023、F-I-024** |
| R2 | 禁止 | 禁止 | **禁止** |

无「一要一否」需用户在 finding 之间裁。F-I-021/022 由本审接受 closed，响应仍走 `/govern` 留痕。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 表达式/F-5/v73 路径可核；逐表 DDL 因 v78 与 site_settings 仍开 |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列分母不因联接表或 v78 无 FK 而扩大 |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放 |
| I-040-004 | required | R3 | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001 仍开放，阻断 C2/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** E-030 使 F-I-021 / F-I-022 可闭合，并使 F-I-002.1 从「未落盘」收窄为「14/15 descriptor 有载体」。不足以冻结 C2/C3 或放行 R2。v78 缺席与 `site_settings` 列序陷阱是本轮新的 required。

建议 `/govern`：

1. 响应本 A-034；接受 F-I-021 / F-I-022 `fixed`；接受 F-I-023 / F-I-024 为新 required；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. 补 v78 两表裸四步正文；写出 `site_settings` 15 列 live CREATE（按 `PRAGMA table_info` cid，不按文件行号）；补 `dict_entries` new CREATE。
3. 落下 PG 显式 DDL 配套附件（`D-019` §6），再 `/audit` 看 F-I-002。
4. 可选：划掉 `D-018` 来源过时句；把 wallet「现行无 `REFERENCES wallet_*`」写成冻结负向盘点。
5. A-031 的 F-I-006 修正**另安排复审**，不搭本条顺手闭合。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree / `apps/`。响应、finding 闭合与是否推进由 `/govern` 处理。
