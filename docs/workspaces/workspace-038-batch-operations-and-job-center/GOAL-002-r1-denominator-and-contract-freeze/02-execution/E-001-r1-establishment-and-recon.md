---
id: E-001-r1-establishment-and-recon
doc: execution-entry
parent: GOAL-002-r1-denominator-and-contract-freeze
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-001 · R1 立项与只读侦察

## 事实（2026-09-19）

### 立项

按 P-001，Root `GOAL-001` 的纲领路线图（R1→R5）已就位，故在 R1 阶段创建子目标 `GOAL-002-r1-denominator-and-contract-freeze`（五件套 + 三个 ledger 目录齐全），承接 Root 的 R1 检查点。

`GOAL-002` 的 4 个显式检查点：

| 检查点 | 内容 | 承接信息项 |
|--------|------|-----------|
| C1 | 分母与作用域矩阵 | `I-038-001` |
| C2 | 契约形态冻结 | `I-038-002` |
| C3 | 首波分母冻结 | `I-038-003` |
| C4 | R1 审计与投影 | — |

审计模式在实施前按 P-002 判定为 **`cross`**（协议/跨边界高影响门禁）。

### 只读侦察（未改动 `apps/**`）

以三个并行只读侦察任务收集 `I-038-001`～`003` 的证据，产出三份 `status: draft` 报告：

| 报告 | 承接 | 行数 | 路径 |
|------|------|------|------|
| Job 种类与作用域 | `I-038-001` | 452 | `attachments/R1-recon-I-038-001-job-kinds-and-scopes.md` |
| 批量异步契约 | `I-038-002` | 180 | `attachments/R1-recon-I-038-002-batch-async-contract.md` |
| 批量操作与长操作清单 | `I-038-003` | 470 | `attachments/R1-recon-I-038-003-batch-operation-inventory.md` |

编排器对关键主张做了独立复核（见 Root `../GOAL-001-batch-operations-and-job-center/02-execution/E-002-r1-recon.md` F-1～F-16），本条目补充本轮**新核验**的三项：

| # | 事实 | 证据 |
|---|------|------|
| G-1 | Job 运行时在组合根**无条件构造**，但 `enabled` 仅在 `plan.HasModule("admin.wallet")` 时置 true；`Start()`/`Stop()` 在 `enabled=false` 时为 no-op ⇒ 不含 `admin.wallet` 的 Profile 下 Job runner **不会启动** | `internal/composition/composition.go:170-197`（`jobRuntime` + `Start`/`Stop` 守卫）、`:582-588`（`jobRuntime.enabled.Store(true)`）、`:1101`（`jobs.Start()`） |
| G-2 | `core.jobs` 迁移的 DDL/索引变更会打破**两处冻结断言**：迁移 catalog checksum（`55e1d3f8…`）与描述符 `Version == 42` / `Name == "async_jobs"` | `internal/store/migrate_test.go:692-693`；`modules/jobs/migration/migration_test.go:15` |
| G-3 | `admin.jobs` 模块尚不存在；现有 `modules/jobs/` 只有 `migration/`（`ModuleID = "core.jobs"`，`Register` 空实现） | `modules/jobs/migration/provider.go:12-20`；`modules/compiled/persistence.go:14,44` |

### 侦察发现的三个影响 R1 冻结的新事实

1. **批量 UI 分母为 0**：生产页面 schema **无一**声明 `batchMapping` / `requiresSelection`；唯一的 schema 级批量声明是 dev 范例 `dev/examples/schema/admin-list-batch.json`（27 个 `schema/*.json` 全量扫描，仅此一处命中）。Root 与 VP-038 文本中「批量操作已在 users/roles/data-dictionary/scheduled-tasks 落地」指的是**后端路由 + 协议能力**，不是**页面 UI**——两者需在 R1 分母中区分。
2. **Job 读面无任何跨 actor 路径，且索引不支持管理列表形状**：`GetForActor` 是 `id + kind + actor_id` 三重限定；`Repository` 无任何通用列表方法（仅 worker 用的 `ListRunnable`）；三个既有索引均为运行期状态机索引，`ORDER BY updated_at DESC` 跨 actor **不被索引覆盖**。⇒ R2 需要**新增 repository 查询方法 + 新增迁移索引**，后者触发 G-2 的两处冻结断言更新。
3. **通用资源工厂给每个非只读资源都挂了 `batch-delete`**，但只有 `users` / `roles` 实现原子 `DeleteBatch`；`dict-types` / `dict-entries` / `scheduled-tasks` 走**顺序回退**（非原子、首败即停、已删不回滚）。

### 本轮**未**做的事（边界）

- **未**改动 `apps/**`（侦察为纯只读；`admin.jobs` 模块建立属 R2/R3）。
- **未**关闭 `I-038-001`～`003`：证据已收集，但冻结口径需 C1～C3 决策落盘；C2/C3 涉及方案选型，须经用户 P-004 裁决。
- **未**执行审计（C4 待 C1～C3 完成后进行）。
- **未**创建 R2/R3 子目标（R1 冻结后再按 P-001 立项）。
- **未**做 Git checkpoint：本轮写入均为新建治理文档（`docs/workspaces/workspace-038-…/GOAL-002-…/` + 三份侦察报告），`apps/**` 无变更；checkpoint 将在 R1 冻结决策落盘后与 C1～C3 交付物一并提交。
