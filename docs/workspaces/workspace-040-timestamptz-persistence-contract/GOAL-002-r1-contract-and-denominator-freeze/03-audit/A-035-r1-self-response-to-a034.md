---
id: A-035-r1-self-response-to-a034
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-034 / F-I-021 F-I-022 closed, F-I-023 F-I-024 intake and fixes
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-035 · R1 self response to A-034

- **source**：self（编排器响应，**不**冒充 independent）
- **verdict**：conditional
- **开放 required**：6（F-I-002、F-I-004、F-I-005、F-I-006、**F-I-023**、**F-I-024**）

## 1. 接受 A-034 的判定

| finding | A-034 verdict | 本响应动作 |
|---------|---------------|------------|
| **F-I-021**（FK 父表重建处置与完整子表清单） | **closed**（P-004 五要求均满足；12 张子表完整、无第 13 张；独立复现与 E-030 自测一致） | 接受闭合 |
| **F-I-022**（v73 ledger 写入路径） | **closed**（设计层；5 处生产写入完整准确、无第六处；两条边界成立） | 接受闭合 |
| **F-I-002.1** | **仍 open**：14/15 descriptor 有载体，**v78 整段缺席**，`site_settings` 15 列 new CREATE 未写 | 见 §2 |
| **F-I-023**（v78 两表无任何 rebuild DDL） | 新增 required | 见 §2.1 |
| **F-I-024**（`site_settings` 按文件行号拼 ALTER 会错序） | 新增 required | 见 §2.2 |
| F-I-004 / F-I-005 / F-I-006 | 维持 open（本轮未触及 / 未复审 A-031） | 照录 |

**A-034 明确接受的形式**（本响应照录，不再改动）：`<秒表达式>` / `<毫秒表达式>` 占位写法**可接受**，不必逐处展开；`operation_log.event` 超长枚举与 wallet 既有重建表允许「源码行号 + 差异」，但**不能**代替 v78 与 `site_settings` 的完整 CREATE。

## 2. 本轮修正（响应 F-I-023 / F-I-024）

### 2.1 F-I-023 · 补 v78 两表的完整 rebuild DDL

新增 `r1-c2-per-table-rebuild-ddl-v1.0-fc.md` **§2.6**，含 `data_scope_policies` 与 `user_data_scopes` 的 legacy / new 完整 `CREATE TABLE` 与**裸四步**正文。

- 两表 DDL 均**逐字取自 live `sqlite_master` 实测**（非按模块文件推断）。
- `data_scope_policies`：5 列（`resource` PK / `owner_column` / `default_scope` CHECK / `enabled` DEFAULT 1 / `updated_at`）；仅 `updated_at` → `TEXT NOT NULL`。
- `user_data_scopes`：4 列（`user_id` + `resource` 复合 PK / `scope_type` CHECK / `updated_at`）；仅 `updated_at` → `TEXT NOT NULL`。
- **子表盘点（实测）**：**无任何表**的存留 DDL 含 `REFERENCES data_scope_policies` 或 `REFERENCES user_data_scopes`（count = 0）→ 裸四步安全。
- 两表**无显式 `CREATE INDEX`**；PK 自动索引由 `PRIMARY KEY` 子句自动重建。

### 2.2 F-I-024 · `site_settings` 改为 live 实测列序 + 补完整 15 列 CREATE

改写为 §2.12，给出 **cid 0–14 实测列序**与完整 15 列 legacy / new `CREATE TABLE` 与 copy 正文。

**实测结论（须如实记录）**：live `PRAGMA table_info` 的 cid 序为

```text
0 id, 1 site_title, 2 logo_url, 3 updated_at,
4 logo_url_light, 5 logo_url_dark, 6 favicon_url, 7 default_locale, 8 site_timezone,
9 default_theme, 10 copyright_text, 11 icp_number,
12 operation_log_retention_days, 13 operation_log_expiration_action, 14 default_currency
```

即 **`operation_log_retention_days`(12) / `operation_log_expiration_action`(13) 确实在 `default_currency`(14) 之前**；同时 `sqlite_master` 的存储文本显示 ALTER 追加段的**文本顺序**里 `default_currency` 位于末尾——**文本顺序 ≠ cid 顺序**，二者不一致。

- 因此 A-034 所担心的错序（把 `default_currency` 放到 retention 之前）在本实现中**并未发生**；但该风险是**真实存在**的——`sqlite_master` 文本顺序确实会误导，故本节改为以 **cid 实测序**为唯一权威并写明理由。
- 按 `D-019` 的列序权威规则（live `PRAGMA table_info`）执行；本次实测已满足该规则，无需再做 P-004。

> **2026-09-20 更正（A-036 §B 指出，本响应照录不改写）**：上文「`sqlite_master` 的存储文本**顺序**里 `default_currency` 位于末尾——文本顺序 ≠ cid 顺序」这一表述**不准确**。A-036 独立复现结论为：live cid 序与 `sqlite_master.sql` 折入文本序**是同一顺序**（都把 `default_currency` 放末列），二者**不矛盾**；真正与 cid 序不同、会误导的是 **Go 模块源文件的行号序**（`settings/migration.go` 的 v62 `:207` 出现在 v46 `:222` **之前**，按行号拼会把 `default_currency` 排到 cid 12）。附件实际按 cid 写出，故 **F-I-024 的关闭要求仍满足**；本条更正只涉及本响应的论证表述，不改变结论。

### 2.3 附带补强

- 新增 `r1-c2-per-table-pg-ddl-v1.0-fc.md`（`D-019` §6 要求的 **PG 显式 DDL**）：给出 §0 的两条 PG 表达式与四种语句骨架（NN / D0 / voucher / 可空），逐 descriptor 列出目标列与语句形态，并声明 **PG 侧不需要 F-5**（`ALTER COLUMN TYPE` 不重命名表，子表 FK 不受影响）。
- §2 子节编号已重排为 `2.1`–`2.12` 连续（修正本轮自身引入的编号重复）；v73–v87 **15 个 descriptor 全部有载体**（v74/v75 在 §1/§3/§4，v86/v87 在 §5）。

## 3. 仍开放（不得放行）

- **F-I-002**：v78 与 `site_settings` 已补；整条是否闭合待 independent 复审本轮修正。**注意 F-I-002 还含 PG 侧**（新附件为候选，未复审）。
- **F-I-004 / F-I-005 / F-I-006**：本轮未触及；A-031 对 F-I-006 的修正仍**未被复审**（A-032 曾要求另安排）。
- **F-I-023 / F-I-024**：修正已落盘，**接受与否待 independent 复审**。

**本响应不闭合任何 required。** C2/C3 未冻结，R2 未放行。

## 4. 下一步

1. 跑 `/audit` 复审本轮修正（重点 F-I-002 / F-I-023 / F-I-024，并附 A-031 的 F-I-006 修正一并复审）。
2. 之后才谈 F-I-004（C3 备份/回滚边界）与 F-I-005（记录 canonical SQL + checksum + 测试改写）。
