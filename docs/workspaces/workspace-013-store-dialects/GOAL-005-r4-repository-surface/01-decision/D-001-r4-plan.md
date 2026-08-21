---
id: D-001
doc: decision-entry
goal: GOAL-005-r4-repository-surface
status: accepted
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# D-001 · R4 方案：收口形状、运行时 SQL 债规则、postgres 启动

## 触发

- Root R3（GOAL-004）done：48 迁移双方言对写、postgres store 级 apply 可用。R4 的阻碍只剩下模块/Handler/jobs 公共签名仍是 `*sql.Tx` 与运行时方言 SQL。
- R1 v1.4 §3/§6 已把运行时债点名：`INSERT OR IGNORE`（operationlog retention）、`LIKE`（wallet/recyclebin）、`ORDER BY … COLLATE NOCASE`（users/roles）、布尔与 `LastInsertId`、`jobs.CommitFunc`。

## 决定

1. **接缝**：各模块 `TxRunner`/`WithTx(ctx, func(*sql.Tx))` 接口统一迁移到 `kernel.Store`（或窄 `kernel.TxRunner`，方法 = `Run(ctx, func(kernel.Tx) error) error`）。仓库方法保持业务语义，仅换端口类型；`store.Store.WithTx` 在全部调用方迁移后降为测试 helper 或删除（R5 前）。
2. **公共签名**：handler/jobs 中出现的 `*sql.Tx`/`*sql.DB` 参数一律改为 `kernel.Tx`；`jobs.CommitFunc` → `func(kernel.Tx) error`；公共面禁止 import 具体驱动。
3. **运行时 SQL 债改写**（每处独立核对 + 测试，遵守 R1 v1.4 §3）：
   - `INSERT OR IGNORE` → `INSERT … ON CONFLICT (...) DO NOTHING`（两方言等价）；
   - `LIKE` → 显式选择并落盘：`ILIKE`（大小写不敏感）或写入前规范化；
   - `ORDER BY … COLLATE NOCASE` → 依赖 PG 的 `CITEXT` 列（v44 已建）或 `LOWER(col)`；
   - 插入取 id → `RETURNING`（`QueryRow`）；
   - 布尔列保持 `INTEGER` 0/1（Go 写 0/1 整数；R3 已按此落盘）——不引入 PG `BOOLEAN` 迁移。
4. **sentinel**：模块用 `errors.Is(err, kernel.ErrNoRows)`；公共面切断 `database/sql` 后 `kernel.ErrNoRows` 不再依赖 `sql.ErrNoRows` 别名（R5 前改为独立 sentinel + 文档）。
5. **composition postgres 启动（S4）**：`openStore` 在 `DBDialect=postgres` 时返回 postgres store（apply 后），仓库已讲 kernel.Tx → 全栈可用；`readyz` 模块门禁（SystemDataReady/Reconcile）在 apply+reconcile 后全绿。补一份 postgres-DSN 启动运行证据。
6. **审计模式**：`compatibility`/`production` 高影响门禁 → S5 关门走 **self + independent**（grok build）。

## 为什么

- 内核端口（R1/R3）已就绪；R4 把它贯彻到调用面，才能让 postgres 真正成为「生产 fork 推荐与验收权威」。
- 收口是一次性全仓改造；用窄 `TxRunner.Run` 保持各仓库只依赖内核端口、不碰方言。
- `CITEXT` 基建（v44）让 `COLLATE NOCASE` 查询侧在 PG 上等价；`LIKE`/`INSERT OR IGNORE` 逐处改写避免行为漂移。

## 未选方案

- 保留 `*sql.Tx` 只换实现：postgres 无法 rebind/无法进公共面；不选。
- 用 `sqlx`/ORM 统一：违反 RT-P03/VP-013；不选。
- R4 内顺带做升级策略/备份：R5 范围；不选。
- 把布尔改成 PG `BOOLEAN`：引入双形态迁移与 Go 类型变化，违反「整数布尔」惯例；不选。

## 影响范围

- `apps/api`：`internal/kernel`（`kernel.ErrNoRows` 独立 sentinel 计划）、`internal/store`（`WithTx` 收口）、`internal/jobs`、`internal/handler`、`internal/composition`、全部模块 `**/store` 仓库与 `systemdata`。迁移包签名已在 R3 完成，不动。

## 后续

- S0 泄漏面扫描 → S1/S2 逐模块迁移（每步 sqlite 回归）→ S3 jobs/handler → S4 postgres 启动证据 → S5 self+independent 关门。
