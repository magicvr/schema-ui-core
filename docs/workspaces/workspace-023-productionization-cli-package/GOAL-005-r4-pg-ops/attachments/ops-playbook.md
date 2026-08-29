# Ops Playbook · 包形态下游运维路径（VP-023 R4 · 2026-08-29）

> 对照主仓契约（module-architecture / VP-013/015/016/021 已交付面）。**包形态与 fork 形态的运维差异 = 升级路径（bump 而非 merge）+ 进程/配置由下游组合根承担**。

## 1. 启动

- SQLite（内嵌默认）：`go run ./cmd/server -dialect sqlite -dsn ./data.db`
- PostgreSQL（生产权威方言）：`go run ./cmd/server -dialect postgres -dsn "postgres://user:pass@host:5432/db?sslmode=disable"`
- 迁移随 `assembly.OpenStore` 自动 apply（checksum 校验 fail-closed；未知/冲突迁移拒绝启动）。
- ⚠️ HTTP server 壳与 config 装载 = **go 后正式化**（assembly 扩展服务器面）；本波运维面 = 组合根/迁移/方言/备份。

## 2. 升级（核心差异）

| 形态 | 动作 | 冲突面 |
|------|------|--------|
| fork | pull upstream + merge/rebase | 文件级冲突、基线漂移 |
| **包** | `schema-ui upgrade`（go get @latest + pnpm add @latest + 探针回归） | **零冲突**（R2 实证：v0.1.0→v0.2.0；R3：renderer 0.1.0→0.2.0） |

- 迁移说明 = 上游 changelog（breaking 迁移节必读）；新迁移随启动自动 apply。
- breaking 大版本：先读 `semver-breaking-policy` 迁移说明，再 bump。

## 3. 迁移

- 全自动（Open 时）；台账 = 全局不可变 checksum（与主仓同）。
- 人工核对：`schema-ui` 无专门命令——库内 `schema_migrations` 视图（store 契约）核对版本/校验和。

## 4. 备份

| 方言 | 命令（主仓契约） | 本波核对 |
|------|------------------|----------|
| SQLite | 升级前快照 `VACUUM INTO '<path>.pre-vX.sqlite'` | 主仓先例（升级前快照链） |
| PG | `pg_dump -Fc` / `pg_restore`（VP-013 A1 契约） | 命令即库 |
- 密钥轮换后的恢复语义：VP-016（JWT current+previous）——包形态同契约。

## 5. 优雅停机 / 连接排空

- 主仓 RT-D02 合同（停机顺序 / HTTP drain / 运行中 Job 语义 / 双方言 Store 排空）——**下游 server 壳落地时按合同接线**（go 后）；本波注明引用。

## 6. 观察信号

- `fresh=true` = 空库首跑；`fresh=false` = 已初始化（幂等重入实证）。
- 升级后探针（probe.mjs/probe-render/probe-six + 组合根冒烟）为消费回归基线。

## 7. 残余（go 后）

HTTP 壳 + config 装载 + compose 运行时（样例已给）的正式化；`schema-ui` 缺 `serve` 命令（R5 建议项）。