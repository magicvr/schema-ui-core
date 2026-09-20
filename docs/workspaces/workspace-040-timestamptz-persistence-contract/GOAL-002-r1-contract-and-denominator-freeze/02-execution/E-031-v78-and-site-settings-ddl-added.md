---
id: E-031-v78-and-site-settings-ddl-added
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-031 · A-034 响应：补 v78 与 site_settings DDL，起草 PG 显式 DDL

## 事实

1. **独立审计 A-034 完成并落盘**（`source: independent`，`auditor: grok-build (grok-4.6 · reasoning high)`）：verdict `conditional`，并作出两项闭合与两项新增：
   - **F-I-021 closed**（P-004 五要求均满足；12 张子表完整、无第 13 张；本机 sqlite 3.51.2 独立复现与 E-030 自测一致）；
   - **F-I-022 closed**（设计层；5 处生产写入完整准确、无第六处；v1 禁改边界与 `0001:r2-baseline` 哈希边界成立）；
   - **F-I-002.1 仍 open**：14/15 descriptor 有载体，**v78 整段缺席**，`site_settings` 15 列 new CREATE 未写；
   - **新增 F-I-023**（v78 两表无任何 rebuild DDL）、**新增 F-I-024**（`site_settings` 按文件行号拼 ALTER 会错序）。
   - 明确接受：`<秒/毫秒表达式>` 占位写法可接受；`operation_log.event` 超长枚举与 wallet 既有重建表允许「源码行号 + 差异」——但**不能**代替 v78 与 `site_settings` 的完整 CREATE。
2. **本轮修正（A-035 响应）**：
   - **§2.6 新增 v78 两表**（`data_scope_policies`、`user_data_scopes`）的 legacy / new 完整 `CREATE TABLE` 与裸四步正文，DDL **逐字取自 live `sqlite_master` 实测**；子表盘点实测确认**无任何表**引用这两张表（count = 0）→ 裸四步安全；两表无显式索引。
   - **§2.12 改写 `site_settings`**：改为 **live `PRAGMA table_info` cid 0–14 实测列序**，补完整 15 列 legacy / new `CREATE TABLE` 与 copy 正文。
   - **§2 子节编号重排**为 `2.1`–`2.12` 连续（修正本轮自身引入的编号重复）；v73–v87 **15 个 descriptor 全部有载体**。
   - **新增 `r1-c2-per-table-pg-ddl-v1.0-fc.md`**：`D-019` §6 要求的 PG 显式 DDL（两条 PG 表达式 + 四种语句骨架 + 逐 descriptor 目标列 + 「PG 侧不需要 F-5」的论证 + 非 DDL 联动项清单）。
3. **F-I-024 的实测结论（须如实记录）**：live cid 序为 `…, 12 operation_log_retention_days, 13 operation_log_expiration_action, 14 default_currency`，即 **retention/expiration 在 default_currency 之前**，与「按文件行号推断」的结论**一致**；但 `sqlite_master` 的**存储文本**顺序里 `default_currency` 位于末尾——**文本顺序 ≠ cid 顺序**。故 A-034 担心的错序在本实现中**未发生**，该风险本身真实存在，已改为以 cid 实测序为唯一权威。
4. **观测手段**：临时在 `apps/api/internal/store/` 放置观测测试，用 `OpenSeeded` 跑完整迁移后 dump `PRAGMA table_info` 与 `sqlite_master`；**观测后已删除**，`apps/` 无残留变更。

## 证据

- A-034 全文：`03-audit/A-034-r1-independent-e030-per-table-rebuild-ddl.md`。
- A-035 响应：`03-audit/A-035-r1-self-response-to-a034.md`。
- 附件变更：`attachments/r1-c2-per-table-rebuild-ddl-v1.0-fc.md` §2.6 / §2.12；新增 `attachments/r1-c2-per-table-pg-ddl-v1.0-fc.md`。
- 实测数据（v78 两表 + site_settings 的 cid 与存储 DDL）见 `E-031` 上文 §3 与附件 §2.6/§2.12。
- 本轮 `apps/**` **未修改**（观测测试已删除，`git status -- apps/` 为空）。

## 状态评估

- **开放 required = 6**（F-I-002、F-I-004、F-I-005、F-I-006、F-I-023、F-I-024）；**未闭合任何 required**（F-I-021/F-I-022 的 closed 由 A-034 判定，非本条自证）。
- 下一步：跑 `/audit` 复审本轮修正（重点 F-I-002 / F-I-023 / F-I-024），并**一并复审 A-031 对 F-I-006 的修正**（A-032 曾要求另安排）。
