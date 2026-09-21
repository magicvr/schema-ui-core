---
id: E-029-fk-rebuild-mode-decided
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-029 · FK 重建模式定案（A-032 响应与 P-004 裁决）

## 事实

1. **独立审计 A-032 完成并落盘**（`source: independent`，`auditor: grok-build (grok-4.6 · reasoning high)`）：verdict `conditional`，**开放 required 由 4 升至 6**（新增 **F-I-021** FK 父表重建处置与子表清单未冻结、**F-I-022** v73 ledger 写入路径不只两个字面）。
2. **A-032 独立复现并确认** E-028 的 F-1～F-4（parent rename 永久改写子表 `REFERENCES`；`legacy_alter_table` 无效；事务内 `foreign_keys=OFF` 为 no-op；朴素重建 CASCADE 删子行或 FK fail），并确认 F-5 在单事务内可行。
3. **A-032 推翻本编排器的两处主张**，均已更正：
   - mechanism §1 的「`<t>_old` 被 DROP 后引用重新指向同名新表」——**为假**，该句已删除，§1 整节改写为「裸四步仅对无 FK 子表的表安全 + 先决条件」；
   - 引用 `SetMaxOpenConns(1)` **不精确**——它只用于 in-memory（`store.go:107`），文件库默认 pool=4（`:29,:109-113`）；已更正并补注「官方 12 步在引擎层可行，不可用的是其前置条件」。
   - `D-018` 中以「既有 rename 重建含被 FK 引用的表也能工作」为依据的理由句**已作废**；**选项 C 本身不变**。
   - finding 的「authsession 6 个 ALTER 列」笔误更正为 **7 列**（合计追加 10 列，live 共 17 列）。
4. **A-032 补充的事实已并入**：F-8 的 `applied_at` 写入路径共**四处**（两个 CREATE 字面 + `stampCatalog` / `applyMigration` / `applyMigrationPG` 三处 INSERT）；v1 checksum 边界；live `PRAGMA table_info` 为列序主权威。
5. **用户 2026-09-20 P-004 裁决两项并落盘 child `01-decision/D-019-fk-parent-rebuild-mode.md`**：
   - **FK 重建模式 = F-5 子女先行 + TEMP 快照**（含十步序列、裸四步的适用前提与先决条件）；
   - **跨 descriptor 子表 = 两次重建**：v74 先以 live DDL（含 `notifications.title_key`/`body_key`）修复 FK、时间列仍 INTEGER，再由 v80/v81 转换类型；**不**把 v80/v81 的时间列转换并进 v74。
   - `D-019` 另冻结：每张父表的**完整子表清单**（10 张引用 `users` 的表逐张列出，含三张无时间列联接表）、v73 的 DDL 与四处写入路径目标形态、`rebuildOperationLog` 的最小安全改动与 append-only 边界论证、PG 显式 DDL 属 R1 冻结交付。
6. **本轮已完成的 A-032 点名更正**：删除 mechanism §1 错误句与 `D-018` 过时理由句（A-033 §2.1）；`A-033` 响应落盘并更新 `03-audit.md` 索引。

## 证据

- A-032 全文：`03-audit/A-032-r1-independent-e028-fk-parent-rebuild-blocker.md`（313 行，含 A–H 逐项实测判定）。
- A-033 全文：`03-audit/A-033-r1-self-response-to-a032.md`。
- `D-019`：`01-decision/D-019-fk-parent-rebuild-mode.md`。
- 代码对位：`internal/store/store.go:29,57,107,109-113`；`internal/store/migrate.go:32,108-132,122-123,253-265,349,367`；`internal/store/identity.go:58,65,317,365-371`；`internal/store/postgres.go:166-167`；`modules/operationlog/migration/migration.go:250,296-301,544-551,709-754,756-776`；`modules/authsession/migration/migration.go:18-23,55-88`；`modules/notifications/migration/migration.go:46`；`modules/account/migration/migration.go:21,28`。
- 本轮 `apps/**` **未修改**（仅文档）。

## 状态评估

- **开放 required = 6**（F-I-002、F-I-004、F-I-005、F-I-006、F-I-021、F-I-022）；C2/C3 未冻结，R2 未放行。
- **未闭合任何 required**：`D-019` 使 F-I-021 具备可复审载体，但接受与否待 independent 复审。
- **下一步**（按 A-032 建议顺序）：写逐表 exact `CREATE` / `INSERT SELECT` / 索引正文（按 `D-019` 的 F-5 与子表清单 + live DDL 权威 + PG 显式 DDL），连同 v73 的 DDL 与四处写入路径、`rebuildOperationLog` 断言进冻结包；然后跑 `/audit` 复审 F-I-002 / F-I-021 / F-I-022。
