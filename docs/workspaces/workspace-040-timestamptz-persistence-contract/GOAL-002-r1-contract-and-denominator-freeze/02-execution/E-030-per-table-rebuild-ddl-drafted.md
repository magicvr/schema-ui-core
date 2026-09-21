---
id: E-030-per-table-rebuild-ddl-drafted
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-030 · 逐表 exact rebuild DDL 正文落盘（F-I-002.1 / F-I-021 / F-I-022 交付物）

## 事实

1. **新增附件** `attachments/r1-c2-per-table-rebuild-ddl-v1.0-fc.md`（`status: freeze-candidate`，565 行，10 个 `##` 节 / 35 个子节 / 23 个 SQL 块）：按用户 `D-019` 的 **F-5 子女先行**模式与 `D-018` 选项 C 的行拷贝机制，写出 v73–v87 全部 15 个 descriptor 的逐表 exact `CREATE TABLE` / `INSERT … SELECT` / `CREATE INDEX` 正文，并覆盖 `D-019` §2/§3 的 FK-preserve 子表。
2. **v74 的 F-5 编排已整链实测通过**（本机 SQLite 3.51.2，`PRAGMA foreign_keys=ON`），探针覆盖 4 张父表 + A/B/C 三组各 1 张子表：

| 检查项 | 结果 |
|--------|------|
| v74 后 `PRAGMA foreign_key_check` | **0 行** |
| v74 后行数（users/roles/user_roles/notifications/user_mfa） | **2 / 1 / 1 / 1 / 1**（无丢行） |
| v74 后 `users.locked_until`（D0） | `0` → **NULL**；非 0 → `2025-09-19T22:13:22.000000Z` |
| v74 后 C 组 `notifications.created_at` 列类型 | 仍为 **INTEGER**（按 `D-019` §3 设计） |
| v74 后 B/C 组子表 `REFERENCES` 文本 | 指向 `users(id)`，不含 `_old` → **OK** |
| v81 后 `notifications` 二次重建（裸四步） | `fk_check = 0`；`created_at` 27 字符 TEXT；FK 文本仍 **OK**；`integrity_check = ok` |

3. **F-I-022 的 v73 写入路径已列全 5 处**（2 处 DDL 字面 + 3 处 INSERT），并标明 v1 `schemaMigrationsDDL` **禁止改**、restore 字面**不进** `0001:r2-baseline` 哈希输入。
4. **F-I-021 关闭要求第 5 项已落盘**：`rebuildOperationLog` 的 fail-closed 断言设计（改动 1–3）与 append-only 边界论证（§7）。
5. **文档过程中发生一次自伤并已修复**：用 PowerShell 拼接/写回该附件时，`Set-Content` 以 ANSI 码页重编码 UTF-8 内容，导致全文中文出现乱码，且行号切片错位产生重复节。已**删除该文件后用 write 工具一次性重写**，并核验：无 U+FFFD、无 GBK 伪影、节号 0–9 唯一、23 个 SQL 块完整。**教训**：本仓库文档含中文，禁止用 `Set-Content`/`Add-Content` 写回正文。

## 证据

- 附件：`attachments/r1-c2-per-table-rebuild-ddl-v1.0-fc.md`（含 §8 整链实测记录表）。
- 模式来源：`D-019-fk-parent-rebuild-mode.md` §1/§2/§3；表达式来源：`r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` §2（A-032 复证）。
- 代码对位：`internal/store/identity.go:58,65,317,319-321`；`internal/store/migrate.go:121-123`；`internal/store/postgres.go:165-167`；`modules/operationlog/migration/migration.go:513-562,703-776,782-786`；`modules/authsession/migration/migration.go:18-23,35-42,47-54,55-88,61-67,74-82,92-100,171-177,200-207,231-237,268-273,294-306,462-476`；`modules/corepersistence/migration/migration.go:53-106`；`modules/jobs/migration/migration.go:15-44`；`modules/wallet/migration/migration.go:145-163,279-303,304-318,391-395`；`modules/datadictionary/migration/migration.go:20-42,76`；`modules/logincaptcha/migration/migration.go:19-30`；`modules/mfa/migration/migration.go:21-36`；`modules/notifications/migration/migration.go:16-25,53-54`；`modules/recyclebin/migration/migration.go:20-31`；`modules/scheduledtasks/migration/migration.go:18-38`；`modules/settings/migration/migration.go:26-31,52-55,69-72,86-91,121-122,189,207,222-223`；`modules/channel/telegram/migration/migration.go:11-16,29-30,34-46,47-63,100-116`；`modules/digitaloffer/migration/migration.go:21-75`。
- 本轮 `apps/**` **未修改**（仅文档）。

## 状态评估

- **未闭合任何 required**；开放 required 仍为 6（F-I-002、F-I-004、F-I-005、F-I-006、F-I-021、F-I-022）。
- 本条目使 **F-I-002.1 / F-I-021 / F-I-022** 具备可复审载体；**接受与否待 independent 复审**，编排器不自证。
- §4/§5 对**超长既有 DDL**（`operation_log.event` 枚举、已被 v33/v64 重建过的 wallet 表）以「源码行号 + 差异说明」代替重抄；若复审要求逐字重抄，下一轮补齐。
- 下一步：跑 `/audit` 复审 F-I-002 / F-I-021 / F-I-022（以及 A-031 的 F-I-006 修正，A-032 曾要求另安排）。
