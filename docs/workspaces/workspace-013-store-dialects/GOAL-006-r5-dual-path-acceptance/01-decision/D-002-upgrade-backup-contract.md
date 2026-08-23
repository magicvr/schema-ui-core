---
id: D-002
doc: decision-entry
goal: GOAL-006-r5-dual-path-acceptance
status: accepted
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# D-002 · 升级策略（Root I-001）与 PG 备份/恢复合同（Root I-004）

## 触发

- Root I-001 / I-004 到达 R5 最晚需要阶段。R1–R4 已交付双方言 catalog、postgres 完整启动；本拍给升级与备份以书面结论 + 可执行证据。

## 决定

### U1 · SQLite→PostgreSQL 升级策略（I-001 → verified）

1. **in-place（原地改库文件）不可行**：SQLite 与 PostgreSQL 是不同引擎、无共享物理格式；且 R1 v1.4 明确要求 PG 时间列 `BIGINT`、`COLLATE NOCASE→CITEXT` 等 DDL 差异，不存在「同一文件直接升级」路径。
2. **支持的升级路径（已实证）**：
   - **fresh bootstrap（verified）**：在 PG 建空库 → 双方言 compiled catalog（48 迁移）apply（live：`TestFullCatalogPostgresBootstrapIntegration` / `TestCompositionPostgresStartup`）→ 逻辑迁移业务数据。
   - **逻辑数据迁移最小原型（verified）**：`TestPostgresDataMigrationPrototype`——从 seeded SQLite store 读用户，经同一领域仓库写回 fresh PG，round-trip 一致（`go test -run TestPostgresDataMigrationPrototype ./internal/store/` PASS，live PG）。
3. **有界残余（VP-013 退出判据 2 形式，书面记录）**：本 VP **不提供**自动化 SQLite→PG 搬运器（批量 dump/import 工具由 fork / 运维选型；既有存量 = fresh bootstrap + 运维自备搬运）。此残余随本决策留痕，后续若出现 in-place 或内置搬运需求另立目标。
4. **不写死唯一迁移器**：本 VP 是「存储方言」架构；内核契约只保证「双方言 schema + 同仓库 API 可读写」，不承诺无人值守自动 in-place。

### U2 · PostgreSQL 备份/恢复合同（I-004 → verified）

1. **替代 SQLite `VACUUM INTO` 的 PG 生产路径**：
   - **逻辑备份**：`pg_dump -F c <db> > backup.dump`（custom 格式，可 `pg_restore -d <new> backup.dump`）——适合常规/跨版本/选择性恢复，**本次已可执行验证**（round-9：r5u2 建表+2 行 → pg_dump → pg_restore → count=2；2026-08-21 独立审在 fresh-bootstrap 库上 catalog 级 round-trip：48 迁移/35 表 checksum 一致）。
   - **客户端主版本（A-001 F-004）**：`pg_dump`/`pg_restore` 客户端主版本应 **≤ 服务器主版本**（同机或容器内同版工具）；否则恢复可能遇到未知 GUC 的 SET 告警（如 `transaction_timeout`，`errors ignored: 1`，功能与 checksum 仍成立）。剧本应钉工具镜像版本或显式允许该类告警。
   - **物理/PITR**（高可用规模）：`pg_basebackup` + WAL 归档；运维按部署规模选型，最小合同要求「至少逻辑备份可恢复」。
2. **恢复验证要求**：任何备份合同上线前必须跑一次「dump → restore 到新库 → 数据/台账可核对」（本 VP 以账台账计数 + 代表性行计数为核对面）。
3. **非目标**：不在本 VP 实现自动备份调度/KMS/TLS；属 VP-009/运维。

## 为什么

- fresh bootstrap + 回放是双方言架构的自然路径（catalog 已双写、PG boot 已证），比发明 in-place 迁移器更稳。
- pg_dump/restore 是 PG 官方逻辑备份，合同最小且可验证；`VACUUM INTO` 无 PG 等价，故显式换为 dump/restore。
- 备份/恢复合同在 R5 是「关闭要求」（退出判据 4），给出可执行验证即可闭门。

## 未选方案

- 自研 SQLite→PG 迁移适配器：scope 过大、非内核；不选。
- 只写文案不做可执行验证：premium 风险；本拍已做 dump/restore round-trip。
- 把 PG 备份合同做成常驻调度：非本 VP；不选。

## 影响范围

- 文档决策 + 运维剧本；无运行时代码变更（R1–R4 契约冻结）。

## 用户书面追认（2026-08-23）

用户经 ask_user_question 书面追认 §U1.3 的有界残余条款（「本 VP 不提供自动化 SQLite→PG 搬运器；存量路径 = fresh bootstrap + 运维自备导出/回放」）。自此本条构成 P-003 `accepted-residual`：范围 = VP-013 退出判据 2；复审触发 = 出现 in-place 或内置搬运需求时另立目标。裁决留痕：[workspace-010 GOAL-033 D-002](../../../../../workspace-010-design-implementation-conformance/GOAL-033-w22-residual-closeout/01-decision/D-002-w22-p004-adjudications.md)。
