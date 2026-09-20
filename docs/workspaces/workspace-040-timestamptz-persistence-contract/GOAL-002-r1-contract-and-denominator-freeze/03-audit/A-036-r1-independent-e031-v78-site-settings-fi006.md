---
id: A-036-r1-independent-e031-v78-site-settings-fi006
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · ad-hoc design-plan + finding-closure · 两块：(1) E-031 / commit b8d157a0 对照 A-034 开放 required=6，核 F-I-023/F-I-024 与 F-I-002 收窄；(2) A-031 / commit 92bf74ef 对照 A-030/A-032 的 F-I-006 关闭要求（A-032/A-034 曾要求另安排、一直未做）· 不是实施审计 · freeze-candidate ≠ 已实施
verdict: conditional
open_required: 4
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-036 · R1 independent · E-031 v78/site_settings/PG DDL + F-I-006 复审

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：ad-hoc + finding-closure（用户指定两块都做。第一块：commit `b8d157a0` / E-031 是否补齐 A-034 点名的 v78 与 `site_settings` 缺口，F-I-002 / F-I-023 / F-I-024 能否闭合或收窄。第二块：A-031 对 F-I-006 的修正对照 A-030 关闭要求。附件为设计材料，不是实施证据）。
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree / `apps/`。
- 实测：`OpenSeeded`（`modernc.org/sqlite` 报 `sqlite_version()=3.53.3`）对临时库跑完整 catalog 后 dump `PRAGMA table_info` / `sqlite_master`；观测测试已删除，`apps/` 无残留。
- `git show --stat b8d157a0`：7 个文件均在 workspace-040；`apps/` 无变更。`git show --stat 92bf74ef`：9 个文件均在 workspace-040；`apps/` 无变更。

## 核对方法

1. 通读 A-034/A-035、E-031、`D-019` §6、rebuild DDL §0/§2.6/§2.12 与目录、PG DDL 全文、A-030 F-I-006 关闭要求、A-031、exact SQL、readwrite spec。
2. 对位源码：`datapermission/migration/migration.go:20-33`、`settings/migration/migration.go`（v7/v10/v40/v46/v62 行号与 Version）、`users_repository.go:496-503`、`accounts.go:194-201,235`、`account_operations.go:181`、`jobs/migration/migration.go:36-47,127`、`password_policy.go:148`、`scheduledtasks/store/repository.go:289-290`。
3. 独立用 `OpenSeeded` 复现 v78 两表 live CREATE、REFERENCES 计数、显式索引、`site_settings` cid 0–14 与 `sqlite_master.sql` 折入文本。
4. 脚本核 §5 六类列号并集：90 且无重复、无 1–90 空洞。

## 成果（有证据）

1. **本轮未实施 DDL/codec。** 两 commit `--stat` 均无 `apps/**`。两份 DDL 附件 `status: freeze-candidate`。E-031 / A-035 自承未闭合任何 required。本审同意该实施边界。
2. **F-I-023 关闭要求已满足**（见 A）。v78 两表 legacy/new CREATE + 裸四步已写；本机 live `sqlite_master` 与附件逐字一致；REFERENCES count=0；无显式索引。
3. **F-I-024 关闭要求已满足**（见 B）。15 列 new CREATE / INSERT / `idx_site_settings_updated_at` 按 live cid 写出。A-035「sqlite_master 文本序 ≠ cid 序」**不成立**（二者都把 `default_currency` 放在末列）；真正会误导的是 **Go 源文件行号序**（v62 `:207` 在 v46 `:222` 之前）。附件实际按 cid 写，关闭要求仍满足。
4. **§2 子节 2.1–2.12 连续；v73–v87 15 个 descriptor 在 SQLite 附件内全部有载体**（见 C）。
5. **PG 显式 DDL 形式可接受为 R1 冻结交付，完备性不足**（见 D）：骨架写法合法，但 **v78 整 descriptor 缺席**（新 **F-I-026**）。「PG 侧不需要 F-5」论证成立。
6. **A-034 接受的占位/行号例外本轮被保持**（见 E）。
7. **F-I-006 关闭要求 1–5 均已满足**（见 F）。`#34` E3 已与毫秒族 `date_trunc` 同一；exact SQL 的 `#72/#73` 负值单路径成立。转换合同表细胞仍写「`= 0` 与 `< 0` 双分支」，归 F-I-002 残留，**不阻断 F-I-006**。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C2 物理合同 / 逐表 SQLite rebuild | **仍不可冻结** | F-I-023/024 closed，但 F-I-002 仍缺 `dict_entries` 可粘贴 new CREATE、PG v78、可执行测试 |
| F-I-006 谓词 exact SQL | **本审接受 closed** | A-030 五条关闭要求均对位 |
| C3 / C4 / R2 放行 | **未满足** | F-I-004/005 未触及；freeze-candidate ≠ 实施 |

## 第一块：A-034 缺口的修正

### A. §2.6 v78（F-I-023）

本机 `OpenSeeded` live：

```text
CREATE TABLE data_scope_policies (
  resource      TEXT PRIMARY KEY,
  owner_column  TEXT NOT NULL,
  default_scope TEXT NOT NULL CHECK (default_scope IN ('all','self')),
  enabled       INTEGER NOT NULL DEFAULT 1,
  updated_at    INTEGER NOT NULL
)
CREATE TABLE user_data_scopes (
  user_id    TEXT NOT NULL,
  resource   TEXT NOT NULL,
  scope_type TEXT NOT NULL CHECK (scope_type IN ('all','self')),
  updated_at INTEGER NOT NULL,
  PRIMARY KEY (user_id, resource)
)
```

与附件 §2.6 legacy **逐字一致**（含空白）。new 仅 `updated_at INTEGER NOT NULL` → `TEXT NOT NULL`。裸四步 + `<秒表达式(updated_at)>` 直出（NN，非 D0）正确。源码 `datapermission/migration/migration.go:20-33` 两表均无 `REFERENCES`。

| 检查项 | 本审结果 |
|--------|----------|
| `sqlite_master.sql LIKE '%REFERENCES data_scope_policies%'` | **count=0** |
| `… REFERENCES user_data_scopes%` | **count=0** |
| 其它表 SQL 提及这两名 | **0** |
| 显式 `CREATE INDEX` | **无**；仅 `sqlite_autoindex_*`（PK 自动） |

裸四步安全；不必进 v74 F-5。**F-I-023：本审接受 closed。**

### B. §2.12 `site_settings`（F-I-024）与「文本序 ≠ cid 序」

本机 live `PRAGMA table_info` cid 0–14：

```text
0 id, 1 site_title, 2 logo_url, 3 updated_at,
4 logo_url_light, 5 logo_url_dark, 6 favicon_url, 7 default_locale, 8 site_timezone,
9 default_theme, 10 copyright_text, 11 icp_number,
12 operation_log_retention_days, 13 operation_log_expiration_action, 14 default_currency
```

**与 A-035 §2.2 / 附件 §2.12 列序一致。** new CREATE 15 列、仅 `updated_at TEXT NOT NULL`、其余 DEFAULT/CHECK 与 live 一致；INSERT 列清单同序；`CREATE INDEX IF NOT EXISTS idx_site_settings_updated_at ON site_settings (updated_at)` 与 v63 源码同文（live `sqlite_master` 折成无 `IF NOT EXISTS`，重建语句保留源码 `IF NOT EXISTS` 正确）。

**独立复现「文本顺序 ≠ cid 顺序」：**

| 权威 | `default_currency` 位置 | retention/expiration 相对位置 |
|------|-------------------------|--------------------------------|
| live cid | cid **14（末列）** | cid 12/13 **在其前** |
| live `sqlite_master.sql` 折入文本 | ALTER 追加段**末尾** | 紧挨其前 |
| Go 源文件行号 | v62 `:207` | v46 `:222`（**文件中更靠后**） |

结论：

1. cid 序「retention/expiration 在 `default_currency` 之前」**成立**。
2. `sqlite_master` 存储文本里 `default_currency` 在末尾 **成立**，且与 cid **同一顺序**（折入顺序 = 版本 apply 序 7→10→40→46→62）。
3. A-035「sqlite_master 文本序 ≠ cid 序」**不成立**。真正 ≠ cid 的是 **模块文件行号序**（A-034 按 `:207` 再 `:222` 拼会把 `default_currency` 放到 cid 12）。附件已改以 cid 为唯一权威，关闭要求满足。无需再做 P-004。

**F-I-024：本审接受 closed。**

### C. §2 编号与 15 个 descriptor 载体

§2 子节现为 **2.1–2.12 连续**，无重复编号。逐一点名：

| v | 表 | 载体 |
|--:|----|------|
| 73 | `schema_migrations`, `mail_outbox`, `mail_config` | §2.1 / §2.2 / §2.3 |
| 74 | 12 张转换表 + FK-preserve + 跨 descriptor 子表 | §1 + §3 |
| 75 | `operation_log`, `operation_log_archive` | §4 |
| 76 | `jobs` | §2.4 |
| 77 | `dict_types`, `dict_entries` | §2.5 |
| 78 | `data_scope_policies`, `user_data_scopes` | **§2.6（本轮补）** |
| 79 | `captcha_challenges`, `captcha_config` | §2.7 |
| 80 | `user_mfa`, `mfa_proofs` | §2.8 |
| 81 | `notifications` | §2.9 |
| 82 | `recycle_items` | §2.10 |
| 83 | `scheduled_tasks`, `task_runs` | §2.11 |
| 84 | `site_settings` | **§2.12（本轮改写）** |
| 85 | wallet 六表 | §5.1–5.4 |
| 86 | telegram 四表 | §5.5–5.6 |
| 87 | digital-offer 三表 | §5.7 |

SQLite 侧 15/15 有载体。`dict_entries` **new CREATE 仍省略**（§2.5 仍是 `CREATE TABLE dict_entries ( …new… badge_style 保留… )`），A-034 A5.3 / F-I-025.1 未修，**阻断 F-I-002 整条闭合**，不重开 F-I-023/024。

### D. PG 显式 DDL（`D-019` §6）

**形式足以作为 R1 冻结交付，不必逐列逐字展开到每张表。** 理由与 A-034 接受 `<秒/毫秒表达式>` 占位相同：§0.1 两条表达式 + §0.2 四种骨架（NN / D0 / voucher / 可空）是唯一形态；逐 descriptor 点名目标列与包裹类别后，R2 可机械实例化。90 段重复 ALTER 只会增加转录漂移。附件 §8 已声明若本审要求逐字展开则下一轮补——**本审不要求**。

v74 31 列均被点名（展开的 D0/NN 样例 + 「同形」清单 + 可空清单），与 conversion contract `#2–#32` 对得上。

**完备性不足：全文零处出现 `data_scope_policies` / `user_data_scopes`。** 标题写「v73–v87 每个 descriptor」，§4/§5 跳过 v78（`#50/#51`，秒族 NN）。这与 A-034 时 SQLite 侧 v78 整段缺席同类。升级为 **F-I-026**。

**「PG 侧不需要 F-5」成立。** `ALTER COLUMN TYPE` 不 `RENAME` 表；`pg_constraint.confrelid` 仍指向同一 OID。本分母时间列不是 FK 目标（子表引用的是 `users.id` 等）。SQLite F-5 的根因是父表 rename 永久改写子表 `REFERENCES` 文本，PG 无此机制。禁止 `pgRebuild` / `pgTimeDDL` 派生的边界正确。若将来对某表改走 rewrite，须另审。

### E. A-034 接受的形式本轮是否保持

| 形式 | 本轮 |
|------|------|
| `<秒表达式>` / `<毫秒表达式>` 占位 | **保持**。§2.6/§2.12 INSERT 仍用占位，未过度展开 |
| `operation_log.event` 超长枚举走源码行号 | **保持**。§4.1 仍钉 `operationLogDigitalOfferDDL[0]`（`:513-521`），不重抄 |
| wallet 既有重建表走「源码行号 + 差异」 | **保持**。§5.1–5.4 仍用 `:304-318` / `:145-162` 等 |
| 不得用行号代替 v78 / `site_settings` 完整 CREATE | **遵守**。这两处已写出完整 CREATE |

未把允许例外扩到 ALTER 拼出的 `dict_entries`（仍缺可粘贴 new CREATE）。

## 第二块：F-I-006 复审（A-031 / `92bf74ef`）

对照 A-030 关闭要求（A-032/A-034 明确未复审）：

| # | 要求 | 判定 |
|--:|------|------|
| 1 | readwrite spec 标 superseded 或删 old/new 表 | **满足。** `r1-c2-readwrite-predicate-spec-v0.1.md` `status: superseded` v0.1.1；声明唯一权威为 exact SQL；点名 L26 锁谓词方向错误 |
| 2 | 修正 `#5` 锁谓词方向 | **满足。** exact §1 `#5`：`locked==true` → `IS NOT NULL AND > ?`；`false` → `IS NULL OR <= ?`。对位 `users_repository.go:496-503` 现行 `> ?` / `<= ?` |
| 3 | ORDER BY vs none 与散文一致，并集覆盖 90 列 | **满足。** `#22/#31/#33/#60/#67/#74/#79` 在 ORDER BY（21）；`#62` 在 none（44）。抽核：`password_policy.go:148` `ORDER BY created_at DESC`；`scheduledtasks/store/repository.go:289-290` 排序键是 `started_at` 不是 `created_at`。计数 `5+3+2+15+21+44=90`；脚本并集 1–90 无重复无空洞 |
| 4 | jobs 四索引 exact old/new `CREATE INDEX` | **满足。** 四行与 `jobs/migration/migration.go:45-47,127` 逐字同（含 `DESC`、`IF NOT EXISTS`）；六态 CHECK `:36-43` 声明逐字保留 |
| 5 | `#6` callsite = `accounts.go:194-201` 并列入写 0 | **满足。** 衰减比较在 `:194-201`；写 0：`accounts.go:235`、`account_operations.go:181`；目标「去除写 0，改写 NULL」 |

附加（用户点名，原属 F-I-002.2，本轮一并核）：

- **`#34` E3 与毫秒族同一**：**是。** conversion contract §1 L45 + 行 `#34` 写明 ELSE 主体 = E2（`date_trunc('microseconds', TIMESTAMPTZ 'epoch' + col * INTERVAL '1 millisecond')`）。A-030 所记「无 `date_trunc`/`to_timestamp`」已改。
- **`#72/#73` 负值单路径**：**exact SQL 内成立。** §1 说明块 + §2 USING 只写 `= 0 → NULL` 与正值；`< 0` 只走 `m0` 预检。**转换合同表 `#72/#73` 细胞仍写「E3（`= 0` 与 `< 0` 双分支）」**，与同文件 §3.5「只走 m0」及 exact SQL 不一致。这是 F-I-002 残留，不是 F-I-006 关闭要求未满足。

**F-I-006：本审接受 closed**（设计层；`apps/` 仍为整数谓词，实施在 R2）。

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-002 · 逐列 USING/rebuild 仍不足以为 C2 冻结

- **严重度**：high · **建议**：required
- **状态**：open（维持；**本轮显著收窄，仍不关闭**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修 / 收窄**：
  1. SQLite 15/15 descriptor 有载体（v78、`site_settings` 已补）。
  2. F-I-023 / F-I-024 closed。
  3. `#34` E3 毫秒 ELSE = E2/`date_trunc`（F-I-002.2 该子项 `fixed`）。
  4. exact SQL `#72/#73` 负值单路径已唯一（m0）。
  5. PG 显式 DDL **形式**接受（骨架 + 列清单）；「PG 无 F-5」成立。
- **仍不闭合**：
  1. **`dict_entries` new CREATE 仍省略**（A-034 A5.3；ALTER 拼出的 live 形状必须可粘贴）。
  2. **PG 附件缺 v78**（升级 F-I-026）。
  3. 转换合同 `#72/#73` 细胞仍写双分支，与 exact SQL / 同文件 §3.5 不一致。
  4. 用例仍是 ID；非法/越界可执行测试未发生。
- **关闭要求**：补 `dict_entries` 可粘贴 new CREATE（含末列 `badge_style`）；补 PG v78 两列（F-I-026）；改写 conversion contract `#72/#73` 细胞与 exact SQL 同一；freeze-candidate ≠ 实施。

### F-I-003 · 90 列 mapping

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-004 · Backup Port ≠ 可执行备份/回滚

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮未触及 C3）

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮无新哈希/测试。`D-017` 单 checksum 约定维持）

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med · **建议**：required
- **状态**：**closed**（本审接受 `fixed`，设计层）
- **关闭证据**：exact SQL `#5` 方向对位 `users_repository.go:496-503`；readwrite spec `superseded`；§5 六类 5+3+2+15+21+44=90、并集 1–90 无空洞，ORDER BY/none 与散文及抽核 SQL 一致；jobs 四索引 exact `CREATE INDEX` 对位 `:45-47,:127`；`#6` callsite `accounts.go:194-201` + 写 0 `:235` / `account_operations.go:181`。代码改动在 R2。

### F-I-007 … F-I-019

- 维持既有 closed/recommended。F-I-018 / F-I-019 **closed** 维持。

### F-I-020 · freeze-candidate 卫生：悬空 C3 引用、readwrite 未 superseded、无限定 D-011

- **严重度**：low · **建议**：recommended
- **状态**：**closed**（本审接受 `fixed`）
- **关闭证据**：转换合同已声明 C3 附件尚未落盘、权威回退到 runbook+Port draft；readwrite spec `superseded`；column-contract-draft L73 改为 Root D-011；「D-015 待落盘」句已删。

### F-I-021 / F-I-022

- **状态**：**closed**（维持 A-034）

### F-I-023 · v78 两表无逐表 rebuild DDL

- **严重度**：high · **建议**：required
- **状态**：**closed**（本审接受 `fixed`）
- **关闭证据**：附件 §2.6 完整 legacy/new CREATE + 裸四步；本机 `OpenSeeded` 3.53.3 live `sqlite_master` 逐字一致；REFERENCES count=0；无显式索引。实施仍在 R2。

### F-I-024 · `site_settings` 15 列 new CREATE 未冻结

- **严重度**：high · **建议**：required
- **状态**：**closed**（本审接受 `fixed`）
- **关闭证据**：§2.12 15 列 new CREATE + INSERT + 索引，列序 = 本机 live cid 0–14。权威是 cid 不是文件行号。A-035 对 sqlite_master≠cid 的表述不准确，但不影响关闭。

### F-I-025 · 冻结包卫生（dict_entries new CREATE / D-018 来源句 / 计数）

- **严重度**：low · **建议**：recommended · **状态**：open（维持）
- **本轮**：`dict_entries` new CREATE 仍省略；rebuild 附件自称「20 张时间列表」未改。新增：conversion contract `#72/#73` 细胞与 exact SQL 双分支残留。

### F-I-026 · PG 显式 DDL 缺席 v78 两表

- **严重度**：high · **建议**：required · **状态**：open（本轮新增）
- **影响门禁**：C2 冻结、F-I-002 PG 侧；关联 `I-040-001`
- **描述**：`D-019` §6 要求 v73–v87 PG DDL 显式书写并作为 R1 冻结交付。新附件 `r1-c2-per-table-pg-ddl-v1.0-fc.md` 自称覆盖每个 descriptor，但 §1–§5 **零处**列出 `data_scope_policies.updated_at` / `user_data_scopes.updated_at`（ledger v78、`#50/#51`，秒族 NN）。骨架形式本身可接受，缺的是这一 descriptor 的列清单。
- **关闭要求**：按 §0.2 NN 骨架写出这两列（秒族 `to_timestamp`+`date_trunc`），列入对应 descriptor 节。不要把它们误加进 F-5。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮 |
|----|------|------------|------|
| F-I-002 | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec | **收窄**：SQLite 15/15 载体；`#34` E3 子项 `fixed`；PG 形式接受。**仍缺** `dict_entries` new CREATE、PG v78、conversion `#72/#73` 细胞、可执行测试 |
| F-I-004 | C3、R2/R3 | 不得把 Port/runbook 当 C3 冻结 | **无新收窄** |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL | **无新收窄** |
| F-I-006 | C2/C3、R2 | — | **closed** |
| F-I-023 | C2、F-I-002.1 | — | **closed** |
| F-I-024 | C2、F-I-002.1 | — | **closed** |
| **F-I-026** | C2、F-I-002 PG | 不得声称 15 descriptor 的 PG DDL 已写完 | **新增**；v78 两列缺席 |

F-I-001、F-I-003、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015、F-I-016、F-I-017、F-I-018、F-I-019、**F-I-006**、**F-I-020**、F-I-021、F-I-022、**F-I-023**、**F-I-024** 为 closed。F-I-008、F-I-009、**F-I-025** 为 recommended open。

**开放 required = 4**（F-I-002、F-I-004、F-I-005、**F-I-026**）。在这些合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-034 independent | E-031 / A-035 自称 | A-036 independent（本条） |
|----|-------------------|--------------------|---------------------------|
| verdict | conditional；open required=6 | 不自证闭合；待本审 | **conditional**；open required=**4** |
| F-I-023 | open（v78 缺席） | 已补 §2.6 | **closed** |
| F-I-024 | open（15 列 CREATE 未写） | 已按 cid 写 | **closed**；sqlite_master≠cid 表述不采纳 |
| F-I-002.1 | 14/15 载体 | 15/15 SQLite | SQLite 15/15 有载体；整条仍 open（dict_entries + PG v78） |
| PG DDL | 未落盘 | 骨架+列清单候选 | **形式接受**；v78 缺席 → F-I-026 |
| PG 无 F-5 | — | 声称成立 | **接受** |
| F-I-006 | 未复审 | （A-031 自称已修、未闭合） | **closed** |
| 占位/行号例外 | 接受 | 保持 | **保持** |
| R2 | 禁止 | 禁止 | **禁止** |

无「一要一否」需用户在 finding 之间裁。F-I-023/024/006 由本审接受 closed，响应仍走 `/govern` 留痕。无需本轮 P-004。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | SQLite 逐表 DDL 因 dict_entries 省略与 PG v78 仍开 |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列分母不因 v78 无 FK 而扩大 |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放 |
| I-040-004 | required | R3 | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001 仍开放，阻断 C2/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** E-031 使 F-I-023 / F-I-024 可闭合，并使 F-I-002.1 从「14/15 + 列序陷阱」收窄为「SQLite 15/15 有载体、缺 dict_entries 可粘贴 CREATE」。A-031 使 F-I-006 可闭合。PG 骨架形式可冻结，但漏 v78（F-I-026）。不足以冻结 C2/C3 或放行 R2。

建议 `/govern`：

1. 响应本 A-036；接受 F-I-023 / F-I-024 / F-I-006 / F-I-020 `fixed`；接受 F-I-026 为新 required；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. 补 PG 附件 v78 两列（NN 秒族骨架）；补 `dict_entries` 可粘贴 new CREATE；改写 conversion contract `#72/#73` 细胞，去掉「`< 0` 双分支」。
3. 之后才谈 F-I-004（C3）与 F-I-005（记录 canonical SQL + checksum + 测试改写）。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree / `apps/`。响应、finding 闭合与是否推进由 `/govern` 处理。
