---
id: A-033-r1-self-response-to-a032
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-032 / FK rebuild mode pending P-004, mechanism corrections, F-I-021 / F-I-022 intake
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-033 · R1 self response to A-032

- **source**：self（编排器响应，**不**冒充 independent）
- **verdict**：conditional
- **开放 required**：6（F-I-002、F-I-004、F-I-005、F-I-006、**F-I-021**、**F-I-022**）

## 1. 接受 A-032 的核心判定

| 项 | A-032 判定 | 本响应动作 |
|----|-----------|------------|
| E-028 FK 父表重建阻塞（F-1～F-4） | **成立**，独立复现 | 接受；不作为夸大处理 |
| F-5 TEMP 子女先行 | 单事务内**可行** | 接受为**候选**；**定案待用户 P-004**（§3） |
| F-6 跨模块列 | **属实** | 接受；已更正「6 个 ALTER 列」的笔误为 **7 个** |
| F-7 PG 正则派生 | **属实**，显式 PG DDL 属 **R1 冻结交付** | 接受；写入冻结包 |
| F-8 无中途读/无自插入 | **属实**，但写入路径**不完整** | 接受，升级为 **F-I-022**；已更正 finding |
| G 既有 `rebuildOperationLog` | **函数层隐患成立**；≠ 当前库已损坏；最小改动 = fail-closed 断言；**不改 0001–0072 stmts** | 接受；设计进 GOAL-002 冻结包、代码进 R2 |
| H 转换表达式 | 抽测范围内正确，未见反例 | 接受 |

## 2. 本轮已执行的更正（响应 A-032 点名项）

### 2.1 删除我方的两处**错误主张**（A-032 第 43/65/83 行点名）

1. **`r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` §1**：删除「`<t>_old` 被 DROP 后引用重新指向同名新表」——**该句为假**（A-032 实测：DROP 后引用**不会**回到同名新表）。§1 已整节改写为：
   - 裸四步**仅对无 FK 子表的表**安全；加**先决条件**（无存留 DDL 含 `REFERENCES <t>` 的子表）；
   - 有 FK 子表的父表**必须**走 F-5；
   - 删除「`rebuildOperationLog` 可作为被引用表也能重建的先例」的说法；指明真正安全的先例是 `rebuildOperationLogWithSessions`。
2. **`SetMaxOpenConns(1)` 引用不精确**（finding F-3 与 mechanism §1）：已改为「`SetMaxOpenConns(1)` **只用于 in-memory**（`store.go:107`）；文件库默认 `sqlitePoolDefault = 4`（`:29,:109-113`），pragma 是 per-connection，pool=4 上无法保证 `Exec(PRAGMA …)` 与 `Begin()` 同连接」。并补注：**官方 12 步在引擎层可行**，不可用的是其**前置条件**在现行 runner 里——不得读成「SQLite 12 步无效」。
3. **`D-018` 的过时理由句**：其中以「既有 rename 重建含被 FK 引用的表也能工作」为可行性依据的句子**已划掉/作废**。**选项 C 本身不变**（纯整数 SQL + 强制 Go 对拍），只是不再依赖那句错误论证。
4. **finding 的「6 个 ALTER 列」笔误**：更正为 authsession 侧 **7 列**；加 3 个跨模块列 = 基线 7 列后共追加 **10 列**，合计 **17 列**（与实测一致）。

### 2.2 已按 A-032 补充的事实

- **F-8 更正**：`applied_at` 的写入路径除 `identity.go:58/:65` 两个 CREATE 字面外，还有**三处 INSERT** 写 Unix 秒——`stampCatalog`（`identity.go:317`）、`applyMigration`（`migrate.go:122-123`）、PG `applyMigrationPG`（`postgres.go:166-167`）。若只改 CREATE 字面，restore 库与 migrate 库形状分叉，且 v73 行本身会把整数写进 TEXT 列（**静默破坏合同**）。
- **边界**：与 v73 同批改「DDL 字面 + 三处 INSERT 格式」；**不**改 v1 checksum（`identity.go` 的 restore 字面不在 `0001:r2-baseline` 的 `r2BaselineDDL` 哈希输入内）；authsession `schemaMigrationsDDL`（`migration.go:18-23`）属 **v1，禁止改**。
- **column 序权威**（A-032 建议已采纳）：主权威 = v72 已 apply 库的 `PRAGMA table_info(users)`（测试夹具 / `OpenWithCatalog`，**不用生产库**）；交叉核对 = 全局版本序；**禁止**按模块文件分组拼列序。

## 3. 已由用户 P-004 裁决（不再是开放项）

**F-I-021 的 FK 重建模式**（A-032 明确「须用户书面处置，本审不代选，仅给建议」）：

- **用户 2026-09-20 裁决 = 选项 A（F-5 子女先行 + TEMP 快照）**；并同时裁决**跨 descriptor 子表 = 选项 A（两次重建：v74 先修 FK，v80/v81 再转类型）**。
- **落盘**：child `01-decision/D-019-fk-parent-rebuild-mode.md`（`status: accepted`）。其中含：
  - §1 F-5 十步序列（单事务）与裸四步的适用前提；
  - §2 **每张父表的完整子表清单**（10 张引用 `users` 的表逐张列出，含三张无时间列联接表）与 live DDL 权威；
  - §3 跨 descriptor 子表两次重建的切法（明确**不**把 v80/v81 的时间列转换并进 v74）；
  - §4 v73 的 DDL 字面 + 四处 `applied_at` 写入路径目标形态（F-I-022 关闭要求）；
  - §5 `rebuildOperationLog` 最小安全改动（rename 前 fail-closed 断言）+ append-only 边界论证；
  - §6 PG 显式 DDL 属 R1 冻结交付。
- **未选**：runner 级官方 12 步 / 并进同一 descriptor / 拆独立子目标（理由均已记入 `D-019`「未选方案」）。
- **边界**：该裁决满足 **F-I-021 关闭要求第 1/3/4 项**；第 2 项（完整子表清单 + live DDL 权威）以 `D-019` §2/§3 为载体，**接受与否由 independent 复审判定**；第 5 项见 §5。**F-I-021 本条不自行闭合。**

## 4. 仍开放（不得放行）

- **F-I-002**：逐表 exact DDL 正文（已抽取、受 F-I-021 阻塞未落盘）；PG 侧须显式 DDL（不得走 `pgTimeColRe`）。
- **F-I-004 / F-I-005 / F-I-006**：本轮未触及 / 无新收窄；**A-031 对 F-I-006 的修正本轮未被复审**，不得据此闭合。
- **F-I-021**：待 P-004 定案 + 每张父表完整子表清单 + live DDL 权威 + 删除机制附件错误句（本轮已完成最后一项）。
- **F-I-022**：v73 设计须列出 CREATE 字面**与**三处 INSERT 的目标类型/格式。

**本响应不闭合任何 required。** C2/C3 未冻结，R2 未放行。

## 5. 下一步（A-032 建议顺序）

1. 用户 P-004 选定 FK 重建模式（§3）。
2. 定案后写逐表 exact DDL 正文（含每张父表的完整子表清单 + live DDL 权威 + F-5 步骤 + PG 显式 DDL）。
3. 把 `rebuildOperationLog` 的 fail-closed 断言与 v73 三处 INSERT 写入进冻结包。
4. 再跑 `/audit` 看 F-I-002/F-I-021/F-I-022。
5. A-031 的 F-I-006 修正**另安排复审**，不搭本条顺手闭合。
