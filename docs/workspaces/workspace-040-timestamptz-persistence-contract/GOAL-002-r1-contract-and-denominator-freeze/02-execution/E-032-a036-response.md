---
id: E-032-a036-response
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-032 · A-036 响应：四项闭合、PG v78 与 dict_entries 补齐

## 事实

1. **独立审计 A-036 完成并落盘**（`source: independent`，`auditor: grok-build (grok-4.6 · reasoning high)`）：verdict `conditional`，**开放 required 由 6 降至 4**——这是本工作区自 A-027 以来首次净下降。
   - **closed**：**F-I-023**（v78 两表 DDL）、**F-I-024**（`site_settings` 15 列 CREATE）、**F-I-006**（谓词 exact SQL，A-030 五条关闭要求逐条对位）、**F-I-020**（freeze-candidate 卫生）。
   - **新增 required**：**F-I-026**（PG 显式 DDL 缺席 v78 两列）。
   - 维持 open：F-I-002（显著收窄）、F-I-004、F-I-005；recommended open：F-I-008、F-I-009、F-I-025。
   - 明确接受：PG DDL 的**骨架 + 列清单形式**可作 R1 冻结交付，不必逐列展开；「PG 侧不需要 F-5」论证成立（`ALTER COLUMN TYPE` 不改 OID，`pg_constraint.confrelid` 不变）。
2. **A-036 独立实测环境**：`OpenSeeded`（`modernc.org/sqlite` 报 `sqlite_version()=3.53.3`）对临时库跑完整 catalog 后 dump `PRAGMA table_info` / `sqlite_master`；观测后删除，`apps/` 无残留。两 commit（`b8d157a0` / `92bf74ef`）的 `--stat` 均无 `apps/**`。
3. **本轮修正（A-037 响应）**：
   - **F-I-026**：PG 附件 §4 拆为 §4.1（**v78 `admin.data-permission`**）与 §4.2，补 `data_scope_policies.updated_at` / `user_data_scopes.updated_at` 两列 NN 秒族骨架；§4 标题 descriptor 枚举加入 v78；明确不属 F-5。
   - **F-I-002.1 第 1 项**：`dict_entries` 由省略式改为**完整可粘贴** `CREATE TABLE`（10 列，`badge_style` 末列，`UNIQUE` 与 FK 逐字保留）；同步写出 `dict_types` 完整 new CREATE 与 copy。
   - **F-I-002.1 第 3 项**：转换合同 §2 的 `#72`/`#73` 两行由「`= 0` 与 `< 0` 双分支」改为「`= 0` **单分支**；`< 0` 不进 USING，只走 `m0` 预检」，与 exact SQL §1 及同文件 §3.5 同一。
4. **更正编排器自身的一处不准确表述**（A-036 §B 指出）：A-035 §2.2 与 E-031 §3 曾写「`sqlite_master` 存储文本序 ≠ live cid 序」。A-036 复现证明**二者同序**（均 `default_currency` 末列）；真正与 cid 不同的是 **Go 源文件行号序**（`settings/migration.go` v62 `:207` 在 v46 `:222` 之前）。原文**保留不改写**，两处各追加更正块，并记入 A-037 §2.4。**F-I-024 的关闭要求不受影响**（附件实际按 cid 写出）。

## 证据

- A-036 全文：`03-audit/A-036-r1-independent-e031-v78-site-settings-fi006.md`（321 行，含两块独立复现与逐条对位表）。
- A-037 响应：`03-audit/A-037-r1-self-response-to-a036.md`。
- 附件变更：`attachments/r1-c2-per-table-pg-ddl-v1.0-fc.md` §4；`attachments/r1-c2-per-table-rebuild-ddl-v1.0-fc.md` §2.5；`attachments/r1-c2-per-column-conversion-contract-v1.0-fc.md` §2 `#72/#73`。
- 本轮 `apps/**` **未修改**。

## 状态评估

- **开放 required = 4**（F-I-002、F-I-004、F-I-005、F-I-026）；**本条未闭合任何 required**（四项 closed 由 A-036 判定）。
- 四项 closed 的判定方均为 independent；编排器只接受并留痕，未自证。
- **编排器错误统计**：本工作区已有 3 次由 independent 纠正编排器的表述/设计错误——A-030（`#5` 锁谓词方向写反）、A-032（「DROP 后引用回同名表」为假 + `SetMaxOpenConns(1)` 引用不精确）、A-036（`sqlite_master` 文本序表述不准确）。三次均已主动更正并留痕。
- 下一步：`/audit` 复审本轮修正（F-I-002 剩余三项 + F-I-026）；随后转 **F-I-004（C3 备份/回滚边界）**，该条自 A-027 以来无实质收窄，且是 R2 的更硬前置。
