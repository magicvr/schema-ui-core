---
id: A-038-r1-independent-e032-a036-response-fi026
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · ad-hoc design-plan + finding-closure · E-032 / commit d1fdb4cc 对照 A-036 开放 required=4，核 F-I-026 与 F-I-002 剩余三项（dict_entries 可粘贴 CREATE、#72/#73 单路径、sqlite_master 更正）· 不是实施审计 · freeze-candidate ≠ 已实施
verdict: conditional
open_required: 3
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-038 · R1 independent · E-032 / A-036 响应复审（F-I-026 + F-I-002 收窄）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：ad-hoc + finding-closure（用户指定核对 commit `d1fdb4cc` / E-032 / A-037 是否补齐 A-036 点名的三项 F-I-002 剩余与新增 F-I-026，并给出 F-I-002 / F-I-026 的闭合或收窄判定。对照基线 A-036：开放 required = 4，F-I-002 / F-I-004 / F-I-005 / F-I-026。附件为设计材料，不是实施证据）。
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区作为审计上下文**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree / `apps/`。
- 实测：`OpenSeeded`（`modernc.org/sqlite` 报 `sqlite_version()=3.53.3`）对临时库跑完整 catalog 后 dump `PRAGMA table_info` / `sqlite_master`；观测测试已删除，`apps/` 无残留。
- `git show --stat d1fdb4cc`：10 个文件均在 workspace-040；`apps/` 无变更。
- **证据窗口**：以 `d1fdb4cc` 已提交附件 + 本机 live dump 为准。工作区另有未跟踪文件 `attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md`，**不在本 commit、不在本 scope**；本条不把它当 F-I-004 证据，也不当 C3 已落盘。

## 核对方法

1. 通读 A-036/A-037、E-032、`D-019` §1/§3/§6、PG DDL 全文、rebuild DDL §2.5/§3.11、conversion contract §1/§2 `#72/#73`/§3.5、exact SQL §1/§2。
2. 对位源码：`datadictionary/migration/migration.go:20-41,76`、`settings/migration/migration.go`（v7/v10/v40/v46/v62 行号与 `Descriptors()` Version）。
3. 独立用 `OpenSeeded` 复现 `dict_types` / `dict_entries` / `site_settings` 的 live `sqlite_master.sql` 与 `PRAGMA table_info` cid。
4. 脚本核 conversion contract 90 行 `table.column` 是否均在 PG 附件中被点名。

## 成果（有证据）

1. **本轮未实施 DDL/codec。** `d1fdb4cc` `--stat` 无 `apps/**`。三份附件仍 `status: freeze-candidate`。E-032 / A-037 自承未闭合任何 required。本审同意该实施边界。
2. **F-I-026 关闭要求已满足**（见 A）。PG 附件 §4.1 写出 v78 两列 NN 秒族骨架；§4 标题含 v78；未误入 F-5。v73–v87 15 个 descriptor 的 90 列全部被点名。
3. **F-I-002.1 第 1 项已满足**（见 B）。`dict_entries` new CREATE 为完整可粘贴 10 列；与 live cid 列/类型/约束/顺序一致（仅时间列 `INTEGER`→`TEXT`）。`dict_types` new CREATE 同样完整。
4. **F-I-002.1 第 3 项已满足**（见 C）。conversion contract §2 `#72/#73` 已改为与 exact SQL §1 及同文件 §3.5 同一的单路径；该文件内无其它「双分支」USING 残留。
5. **A-037 §2.4 / E-031 更正块关于三序关系的主张正确**（见 D）；补一句完整性：`Descriptors()` 的 Version 字段序与 cid 同序。
6. **A-036 接受的形式本轮被保持**（见 E）。
7. **F-I-002 整条仍不能闭合**（见 F）：A-036 点名的三项设计剩余已修；仍余「用例仍是 ID；非法/越界可执行测试未发生」。
8. **本轮未引入新的 required finding**（见 G）。`dict_entries` 新 CREATE 与 §3.11 C 组 / `D-019` 两次重建不冲突。F-I-025 recommended 收窄、仍 open。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C2 物理合同 / 逐表 SQLite rebuild / PG 显式 DDL | **仍不可冻结** | F-I-026 closed；F-I-002 设计三项 `fixed`，整条仍 open（可执行测试） |
| F-I-026 PG v78 | **本审接受 closed** | §4.1 NN 秒族骨架 + 15 descriptor 90 列均点名 |
| C3 / C4 / R2 放行 | **未满足** | F-I-004/005 未触及；freeze-candidate ≠ 实施 |

## 对用户 A–G 的直接判定

### A. F-I-026（PG 附件缺 v78）

`r1-c2-per-table-pg-ddl-v1.0-fc.md` §4 现拆为 §4.1 / §4.2。§4.1 正文：

```sql
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" DROP NOT NULL;
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" TYPE timestamptz(6)
  USING date_trunc('microseconds', to_timestamp("updated_at"::double precision));
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" SET NOT NULL;
-- user_data_scopes.updated_at 同形（完整三行已写出）
```

| 检查项 | 本审结果 |
|--------|----------|
| 与 §0.2 NN 骨架一致 | **是。** DROP NOT NULL → TYPE timestamptz(6) USING … → SET NOT NULL |
| 秒族 `to_timestamp` + `date_trunc` | **是。** `date_trunc('microseconds', to_timestamp("updated_at"::double precision))`，与 §0.1 秒族及 conversion E1 同一家族 |
| 误当 F-5 | **否。** §4.1 明示不属于 F-5；PG 侧无 rename；SQLite 侧本机 `REFERENCES data_scope_policies` / `user_data_scopes` count=0 |
| ledger `#50`/`#51` | **对位。** 文内写明 `#50` = `data_scope_policies.updated_at`、`#51` = `user_data_scopes.updated_at`，均 `S-NN` |
| §4 标题含 v78 | **是。** `v76 / v77 / **v78 admin.data-permission** / v79 / v82 / v84` |

**15 个 descriptor 的 PG 目标列是否全部点名**（对照 conversion contract 90 行，脚本核 `table`+`column` 均出现）：**90/90，无缺失。** 逐 descriptor：

| v | 列 | PG 附件位置 |
|--:|----|-------------|
| 73 | `#1/#33/#34` | §1 展开 |
| 74 | `#2–#32`（31 列） | §2 展开样例 + 「同形」清单 + 可空清单 |
| 75 | `#35–#37` | §3 |
| 76 | `#41–#45` | §4.2 `jobs` 清单 |
| 77 | `#46–#49` | §4.2 `dict_types` / `dict_entries` |
| **78** | **`#50/#51`** | **§4.1（本轮补）** |
| 79 | `#52–#55` | §4.2 captcha |
| 80 | `#63–#66` | §5 mfa |
| 81 | `#38–#39` | §5 notifications |
| 82 | `#56–#57` | §4.2 recycle |
| 83 | `#58–#62` | §5 scheduled-tasks |
| 84 | `#40` | §4.2 `site_settings.updated_at` |
| 85 | `#67–#77` | §5 wallet + voucher 样例 |
| 86 | `#78–#84` | §5 telegram + D0 样例 |
| 87 | `#85–#90` | §5 digital-offer |

无第 16 个 descriptor；无漏列。A-036 接受的「骨架 + 列清单」形式维持；v78 按关闭要求写出完整 NN 三行，不是过度展开全 90 段 ALTER。

**F-I-026：本审接受 closed。**

### B. F-I-002.1 第 1 项（`dict_entries` 可粘贴 new CREATE）

本机 `OpenSeeded` live `sqlite_master`（3.53.3；ALTER 已折入）：

```text
CREATE TABLE dict_entries (
  id         TEXT PRIMARY KEY,
  dict_key   TEXT NOT NULL REFERENCES dict_types(key) ON DELETE CASCADE,
  entry_key  TEXT NOT NULL,
  label      TEXT NOT NULL,
  enabled    INTEGER NOT NULL DEFAULT 1,
  sort       INTEGER NOT NULL DEFAULT 0,
  remark     TEXT,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL, badge_style TEXT NOT NULL DEFAULT 'default',
  UNIQUE (dict_key, entry_key)
)
```

cid 0–9：`id, dict_key, entry_key, label, enabled, sort, remark, created_at, updated_at, badge_style`。

附件 §2.5 new CREATE：10 列、`badge_style` 末列、`UNIQUE (dict_key, entry_key)`、FK `REFERENCES dict_types(key) ON DELETE CASCADE` 逐字保留；仅 `created_at`/`updated_at` → `TEXT NOT NULL`。空白/折行（live 把 `badge_style` 折在 `updated_at` 同一行）不计入差异。索引 `idx_dict_entries_dict_key ON dict_entries(dict_key, sort)` 与 live 一致。

`dict_types` new CREATE：8 列，与 live `sqlite_master` 逐列一致，仅两时间列 `INTEGER`→`TEXT NOT NULL`。完整可粘贴。

**F-I-002.1 第 1 项：本审接受 `fixed`。**

### C. F-I-002.1 第 3 项（`#72/#73` 细胞）

`d1fdb4cc` 仅改 conversion contract 两行：由「E3（`= 0` 与 `< 0` 双分支）」改为「E3（`= 0` 单分支；**`< 0` 不进 USING**，只走 `m0` 预检 fail closed）」。

与 exact SQL §1 说明块（`< 0` 只走 `m0`；USING 只处理 `= 0 → NULL` 与正值）及同文件 §3.5 **同一**。

同文件其它残留：

| 表述 | 判定 |
|------|------|
| §2 `#72/#73` 的 `old` 列「legacy ≤0 视为缺失」 | **不是 USING 双分支**。这是 inventory 对现行运行时的事实（v0.3 L97–98）；USING 细胞已单路径 |
| §1 E3 模板仍写 `= 0 THEN NULL`（无 `< 0`） | 与单路径一致 |
| 全文已无「双分支」字样 | 已核 |

PG 附件 voucher 样例仍用 `IS NULL OR = 0`（§0.2 voucher 骨架），**不含 `< 0`**，与 rebuild 附件 voucher 包裹同一；`NULL` 走 `= 0` 的 ELSE 时 `to_timestamp(NULL)` 仍为 NULL，与 exact SQL 功能等价。不重开本项。

`#72/#73` 仍标 E3，而 §1 写 E3「D0 列专用」——这是 A-030 已记录的特化，本轮未恶化；USING 语义已与 exact SQL 对齐。**不升为新 finding。**

**F-I-002.1 第 3 项：本审接受 `fixed`。**

### D. 更正是否准确（cid / sqlite_master 文本序 / Go 行号序）

本机独立复现 `site_settings`：

| 权威 | `default_currency` 位置 | retention / expiration |
|------|-------------------------|------------------------|
| live cid | cid **14（末列）** | cid 12 / 13 **紧挨其前** |
| live `sqlite_master.sql` 折入文本 | ALTER 追加段**末尾**（`…expiration_action …, default_currency …`） | 紧挨其前 |
| Go 文件物理行号（ALTER 字面） | v62 `siteCurrencyDDL` **`:206-207`** | v46 `siteOperationLogRetentionDDL` **`:221-223`（文件中更靠后）** |
| Go `Descriptors()` Version 字段 | Version **62** at `:169` | Version **46** at `:162`（**先于 62**） |

结论：

1. **cid 序 = sqlite_master 折入文本序**：成立。二者都把 `default_currency` 放末列；折入顺序 = apply 序 7→10→40→46→62。A-035「文本序 ≠ cid 序」不成立。A-037 / E-031 更正块这一句 **正确**。
2. **真正会误导的是按文件物理行号拼接 ALTER 字面**：v62 `:207` 在 v46 `:222` 之前，按行号拼会把 `default_currency` 排到 cid 12。这一句 **正确**。
3. **完整性补一句**：误导面**仅限**「按文件行号拼 ALTER 变量字面」。`Descriptors()` 的 Version 字段序（7, 10, 40, **46**, **62**, 63）**与 cid/apply 序一致**。不是「Go 源码整体序 ≠ cid」。更正**正确、基本完整**；缺的是这句限定，不构成 finding，无需 P-004。F-I-024 closed 维持。

### E. A-036 接受的形式是否保持

| 形式 | 本轮 |
|------|------|
| `<秒表达式>` / `<毫秒表达式>` 占位 | **保持。** §2.5 INSERT 仍用 `<秒表达式(created_at)>`，未展开 codec |
| PG 骨架 + 列清单 | **保持。** 仅 v78 两列按关闭要求写出 NN 三行；其余同形列仍清单 |
| `operation_log.event` 超长枚举走源码行号 | **保持。** §4.1 仍钉 `operationLogDigitalOfferDDL[0]`（`:513-521`） |
| wallet 既有重建表走「源码行号 + 差异」 | **保持。** §5.1–5.4 未改 |

未把允许例外扩回 `dict_entries`（该表已写出完整 CREATE）。无过度展开，也无把 v78 再缩回行号。

### F. F-I-002 整条

A-036「仍不闭合」四项的本审判定：

| # | A-036 剩余 | 本审 |
|--:|------------|------|
| 1 | `dict_entries` new CREATE 省略 | **`fixed`**（见 B） |
| 2 | PG 附件缺 v78 | **`fixed`**（升级为 F-I-026；本审接受 F-I-026 closed） |
| 3 | conversion `#72/#73` 双分支 | **`fixed`**（见 C） |
| 4 | 用例仍是 ID；非法/越界可执行测试未发生 | **仍开放** |

**整条不能闭合。** 全部剩余项（不要只列一项）：

1. **用例仍是 ID；非法/越界可执行测试未发生。** `T-72-NEG` / `P-72-NEG` 等仍是标识符；`apps/` 未改，没有可重复执行的负值/越界迁移测试。这是 A-036 点名且本轮未触及的剩余。
2. **freeze-candidate ≠ 实施**（边界，非新缺口）：本轮无 DDL/codec 落地。F-I-006 曾在设计层闭合，故本项单独不阻断把设计子项标 `fixed`；它阻断把 F-I-002 当成 C2 已实施或放行 R2。

本审**不**把 `jobs` / `mail_config` 的省略式 new CREATE 升为 F-I-002 剩余项：A-036 关闭要求只点名 ALTER 拼出的 `dict_entries`；A-034 已接受源码行号例外。`jobs` 仍 `CREATE TABLE jobs ( …new… )`，属既有卫生，归 F-I-025 口径而非本条新缺口。

### G. 本轮是否引入新 finding

**无新 required。无新 recommended。**

1. **PG §4/§5 编号**：原 §4 正文原样成为 §4.2；§5 标题与正文未改。§4 标题加入 v78，并去掉已不准确的全节「（秒族）」标签（§4.2 内 jobs 本就是毫秒族，属既有混排）。§4.1 插入不破坏后续编号。注记「v77/v78/v82/… 部分唯一索引」把 v78 写进去略宽（v78 仅 PK），**不升 finding**。
2. **`dict_entries` vs §3.11 C 组 vs `D-019` 两次重建**：**无冲突。** `D-019` §3 两次重建只覆盖 `notifications` / `user_mfa` / `mfa_proofs`（C 组；v74 先修 FK、时间列仍 INTEGER，v80/v81 再转类型）。`dict_entries` 是 v77 **同 descriptor** 子表（`dict_types` 父、两边都有时间列），走 `D-019` §1 一次 F-5：TEMP 快照 → DROP 子表（不 rename 成 `_old`）→ 重建父表 TEXT → 重建子表 TEXT+FK。§2.5 正是该切法。§3.11 仍是 C 组 INTEGER 建回，未被本轮改写。

F-I-025 recommended 收窄：`dict_entries` new CREATE 与 `#72/#73` 双分支两项已修；**仍 open** 的是 rebuild 附件自称「20 张时间列表」（ledger 为 44 张时间列表）。conversion contract §3.4 仍写「逐表 exact SQLite rebuild DDL 仍未写出」——陈旧句，**本轮未引入**，不新开号。

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-002 · 逐列 USING/rebuild 仍不足以为 C2 冻结

- **严重度**：high · **建议**：required
- **状态**：open（维持；**本轮再收窄，仍不关闭**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修 / 收窄**：
  1. `dict_entries` 可粘贴 10 列 new CREATE（badge_style 末列；与 live cid 一致）；`dict_types` 同步完整。F-I-002.1 第 1 项 `fixed`。
  2. PG v78 两列 NN 秒族骨架写出（F-I-026 closed）。F-I-002.1 第 2 项 `fixed`。
  3. conversion `#72/#73` 细胞与 exact SQL / §3.5 单路径同一。F-I-002.1 第 3 项 `fixed`。
- **仍不闭合（全部剩余项）**：
  1. 用例仍是 ID；非法/越界可执行测试未发生。
  2. freeze-candidate ≠ 实施（边界；`apps/` 未改）。
- **关闭要求**：可执行负值/越界/sentinel 测试（或用户书面 residual）；freeze-candidate 不得被当成已实施。设计三项不再阻断本条闭合，但测试项仍阻断。

### F-I-003 · 90 列 mapping

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-004 · Backup Port ≠ 可执行备份/回滚

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮未触及 C3。工作区未跟踪的 C3 文件名不在 `d1fdb4cc` 证据窗）

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮无新哈希/测试。`D-017` 单 checksum 约定维持）

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med · **建议**：required · **状态**：**closed**（维持 A-036）

### F-I-007 … F-I-024

- 维持既有 closed/recommended。F-I-018 / F-I-019 / F-I-020 / F-I-021 / F-I-022 / F-I-023 / F-I-024 **closed** 维持。

### F-I-025 · 冻结包卫生（dict_entries new CREATE / D-018 来源句 / 计数）

- **严重度**：low · **建议**：recommended · **状态**：open（**收窄**）
- **本轮**：`dict_entries` new CREATE 已补；`#72/#73` 双分支已改。**仍 open**：rebuild 附件 L16「20 张时间列表」计数不成立（ledger 44 张）。

### F-I-026 · PG 显式 DDL 缺席 v78 两表

- **严重度**：high · **建议**：required
- **状态**：**closed**（本审接受 `fixed`）
- **关闭证据**：PG 附件 §4.1 按 §0.2 NN 骨架写出 `data_scope_policies.updated_at` / `user_data_scopes.updated_at`；秒族 `to_timestamp`+`date_trunc`；ledger `#50/#51`；明示非 F-5；§4 标题含 v78；conversion 90 列在 PG 附件全部点名。实施仍在 R2。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮 |
|----|------|------------|------|
| F-I-002 | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec | **再收窄**：A-036 三项设计剩余 `fixed`。**仍缺**可执行非法/越界测试 |
| F-I-004 | C3、R2/R3 | 不得把 Port/runbook 当 C3 冻结 | **无新收窄** |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL | **无新收窄** |
| **F-I-026** | C2、F-I-002 PG | — | **closed** |

F-I-001、F-I-003、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015、F-I-016、F-I-017、F-I-018、F-I-019、F-I-006、F-I-020、F-I-021、F-I-022、F-I-023、F-I-024、**F-I-026** 为 closed。F-I-008、F-I-009、**F-I-025** 为 recommended open。

**开放 required = 3**（F-I-002、F-I-004、F-I-005）。在这些合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-036 independent | E-032 / A-037 自称 | A-038 independent（本条） |
|----|-------------------|--------------------|---------------------------|
| verdict | conditional；open required=4 | 不自证闭合；待本审 | **conditional**；open required=**3** |
| F-I-026 | open（PG 缺 v78） | 已补 §4.1 | **closed** |
| F-I-002.1 第 1 项 | dict_entries 省略 | 已写 10 列 CREATE | **`fixed`**（live cid 对位） |
| F-I-002.1 第 3 项 | `#72/#73` 双分支 | 已改单路径 | **`fixed`** |
| F-I-002 整条 | open | 三项已修，整条待审 | **仍 open**（测试项） |
| sqlite_master≠cid 更正 | 指出不成立 | 更正块：cid=文本序，≠的是 Go 行号 | **接受**；补 Version 字段序与 cid 同序 |
| 形式 | 占位/行号/骨架 | 保持 | **保持** |
| 新 finding | F-I-026 | — | **无** |
| R2 | 禁止 | 禁止 | **禁止** |

无「一要一否」需用户在 finding 之间裁。F-I-026 由本审接受 closed，响应仍走 `/govern` 留痕。无需本轮 P-004。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | SQLite/PG 设计三项已齐；可执行测试仍开，阻断 C2 冻结 |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列分母不因本轮补丁而扩大 |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放 |
| I-040-004 | required | R3 | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001 仍开放，阻断 C2/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** `d1fdb4cc` 使 **F-I-026 closed**，并使 F-I-002 从「缺 dict_entries CREATE + PG v78 + 双分支细胞 + 测试」收窄为「只余可执行测试」。不足以冻结 C2/C3 或放行 R2。A-037 对三序关系的更正正确。

建议 `/govern`：

1. 响应本 A-038；接受 F-I-026 `fixed`；接受 F-I-002 三项设计剩余 `fixed`、**整条仍 open**；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. F-I-002 剩余只剩可执行测试——这与 F-I-005（canonical SQL / checksum / 测试改写）一样结构性依赖 R2 落码。R1 内可把 F-I-002 标为「设计子项已齐、实施/测试仍开」，但**不得**因此闭合整条或冻结 C2。
3. 下一硬缺口仍是 **F-I-004（C3）**。未跟踪的 C3 文件名须先经 `/govern` 纳入证据窗再另安排 `/audit`。
4. 顺手改 rebuild 附件「20 张时间列表」计数（F-I-025）。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree / `apps/`。响应、finding 闭合与是否推进由 `/govern` 处理。
