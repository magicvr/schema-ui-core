---
id: A-037-r1-self-response-to-a036
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-036 / F-I-023 F-I-024 F-I-006 F-I-020 closed, F-I-026 intake and fix, dict_entries and #72/#73 corrections
verdict: conditional
open_required: 4
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-037 · R1 self response to A-036

- **source**：self（编排器响应，**不**冒充 independent）
- **verdict**：conditional
- **开放 required**：**4**（F-I-002、F-I-004、F-I-005、**F-I-026**）——由 6 降至 4

## 1. 接受 A-036 的判定（四项闭合）

| finding | A-036 verdict | 本响应动作 |
|---------|---------------|------------|
| **F-I-023**（v78 两表无 rebuild DDL） | **closed**（live `sqlite_master` 逐字一致；REFERENCES count=0；无显式索引） | 接受闭合 |
| **F-I-024**（`site_settings` 15 列 new CREATE） | **closed**（列序 = live cid 0–14；权威是 cid 不是文件行号） | 接受闭合 |
| **F-I-006**（CHECK / 部分索引 / 谓词列表） | **closed**（A-030 五条关闭要求逐条对位；`#5` 方向、superseded、§5 计数与散文、jobs 四索引、`#6` callsite 均满足） | 接受闭合 |
| **F-I-020**（freeze-candidate 卫生） | **closed**（C3 悬空引用、readwrite superseded、Root D-011、「待落盘」句四项均已修） | 接受闭合 |
| **F-I-026**（PG 显式 DDL 缺席 v78） | **新增 required** | 见 §2.1 |
| F-I-002 | 维持 open（显著收窄） | 见 §2.2/§2.3 |
| F-I-004 / F-I-005 | 维持 open | 照录 |
| F-I-025 | recommended open | 见 §2.2 |

## 2. 本轮修正

### 2.1 F-I-026 · PG 附件补 v78 两列

`r1-c2-per-table-pg-ddl-v1.0-fc.md` §4 拆为 §4.1（**v78 `admin.data-permission`**）与 §4.2，新增该 descriptor 的两列秒族 NN 骨架正文：

```sql
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" DROP NOT NULL;
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" TYPE timestamptz(6)
  USING date_trunc('microseconds', to_timestamp("updated_at"::double precision));
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" SET NOT NULL;
-- user_data_scopes.updated_at 同形
```

- 对位 ledger `#50` / `#51`（`S-NN`）。
- 明确**不**属于 F-5：PG 侧本无 F-5 需求；两表在 SQLite 侧亦无 FK 子表（§2.6 实测 count=0）。
- 同时把 §4 标题段的 descriptor 枚举改为含 v78，并说明 v78 在 PG 侧无显式索引、无需重建。

### 2.2 F-I-002.1 第 1 项 · `dict_entries` 补**可粘贴** new CREATE

`r1-c2-per-table-rebuild-ddl-v1.0-fc.md` §2.5 的 `dict_entries` 由「省略式」改为**完整可粘贴** `CREATE TABLE`（10 列，`badge_style` 末列，`UNIQUE (dict_key, entry_key)` 与 FK 逐字保留，仅 `created_at`/`updated_at` → `TEXT NOT NULL`），并同步写出 `dict_types` 的完整 new CREATE 与 copy 正文。

### 2.3 F-I-002.1 第 3 项 · 转换合同 `#72/#73` 细胞改单分支

`r1-c2-per-column-conversion-contract-v1.0-fc.md` §2 的 `#72`/`#73` 两行由「E3（`= 0` 与 `< 0` 双分支）」改为「E3（`= 0` **单分支**；**`< 0` 不进 USING**，只走 `m0` 预检 fail closed）」，与 exact SQL §1 及同文件 §3.5 同一。

### 2.4 更正本编排器自己的一处**不准确表述**（A-036 §B 指出）

A-035 §2.2 与 E-031 §3 曾写「`sqlite_master` 的存储文本顺序 ≠ live cid 顺序」。**A-036 独立复现证明该表述不准确**：live cid 序与 `sqlite_master.sql` 折入文本序**是同一顺序**（均 `default_currency` 末列），二者不矛盾；真正与 cid 不同的是 **Go 源文件行号序**（`settings/migration.go` 的 v62 `:207` 在 v46 `:222` **之前**，按行号拼会把 `default_currency` 排到 cid 12）。

- 处理：**保留原文不改写**（自审条目按落盘事实保留），在两处各追加**更正块**，并记入本响应。**F-I-024 的关闭要求不受影响**（附件实际按 cid 写出）。
- 这是本工作区第 3 次由 independent 纠正编排器的表述/设计错误（前两次：A-030 的 `#5` 锁谓词方向、A-032 的「DROP 后引用回同名表」与 `SetMaxOpenConns(1)`）。三次均由编排器主动更正并留痕。

## 3. 仍开放（不得放行）

- **F-I-002**：A-036 点名的三项剩余（`dict_entries` 可粘贴 CREATE、PG v78、`#72/#73` 细胞）**本轮已全部修正**；整条是否闭合**待 independent 复审本轮修正**。仍含「用例仍是 ID、非法/越界可执行测试未发生」。
- **F-I-004**：C3 备份/回滚边界，本轮未触及（PG 转换后 artifact 独立 token、`CreateRecoveryPoint` 包路径与调用点、restore-to-new-db harness）。
- **F-I-005**：无 v73+ canonical SQL / `MigrationChecksum` 哈希落盘；测试未改写。`D-017` 单 checksum 约定维持。
- **F-I-026**：修正已落盘，**接受与否待 independent 复审**。
- **F-I-025**（recommended）：`dict_entries` 本轮已补；附件自称「20 张时间列表」的计数表述尚未改；`#72/#73` 双分支残留本轮已改。

**本响应不闭合任何 required。** C2/C3 未冻结，R2 未放行。

## 4. 下一步

1. 跑 `/audit` 复审本轮修正（重点 F-I-002 剩余三项与 F-I-026）。
2. 之后转 **F-I-004（C3 备份/回滚边界）**——这是 R2 的更硬前置，且自 A-027 以来一直无实质收窄。
3. F-I-005 需在 R2 落码时才能记录真实哈希，属结构性依赖。
