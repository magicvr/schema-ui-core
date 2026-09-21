---
id: D-001-workspace-root-establishment
doc: decision-entry
parent: GOAL-001-batch-operations-and-job-center
status: accepted
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# D-001 · 工作区与 Root 建立：VP-038 激活落盘、`I-038-004` 裁决与 freshness 记录

## 背景

用户 2026-09-19 指令「激活，然后开设工作区」，承接同日已确认的组合层下一拍 [VP-038-batch-operations-and-job-center](../../../../vision/plans/VP-038-batch-operations-and-job-center.md)（`planned` v0.1.0，计划 self = `VRev-098` `pass`）。

`/vision` 已完成激活事务：VP-038 `planned → active` v0.2.0，激活 self = `VRev-099` `pass`（0 required），并在 `roadmap.md` / `reviews.md` / `revisions.md`（`VR-085`）/ `workspaces.md` 同步组合投影。本决策记录 `/govern` 侧的建区落地与该激活事务中两项 P-004/P-005 门禁的最终结论。

## 决策

### 1 · 用户 P-004 裁决：`I-038-004` Profile / 模块矩阵边界

**裁决 = 方案 A**：新建 **`admin.jobs`** 模块，进入 **admin 默认集**；`mvp` / `demo` 不启用。

**定性 = Profile 内容扩展，不改装配语义**，依据：

| 项 | 事实 | 锚点 |
|----|------|------|
| Job 运行时已是 Profile 无关设施 | 组合根**无条件**构造 `jobs.NewRepository(st)` + `jobs.NewRunner(...)`，不依赖 Profile 判定 | `apps/api/internal/composition/composition.go:171`–`178`；启动/停止 `:1101`、`:1123`、`:1160` |
| Job schema 由编译目录全局拥有 | `core.jobs` 为 migration-only 模块，`Register` 空实现；迁移随 `compiled.PersistenceProviders` 全局装配 | `apps/api/modules/jobs/migration/provider.go:19`；`apps/api/modules/compiled/persistence.go:14` |
| 既有「内容扩展」裁定先例 | 往 `profileDefaults[ProfileAdmin]` 追加模块 ID 被反复裁定为内容扩展（非装配语义变更）：F-01 dashboard / F-02 data-transfer / F-03 account / F-04 notifications / S-01 / S-02 / S-03 / S-04 / S-09 / S-10 / S-11 / S-12 / S-14 | `apps/api/kernel/profile.go:46`–`93`；playbook M4 |
| `go` 失效触发指向装配语义 | `ResolveProfile` / `ParseModuleList` 逻辑、`BuiltinModules` 机制、Manifest 聚合规则、协议 pin、共同门禁语义**零改动** | VP-008 §`go` 消费有效性失效触发表；`kernel/profile.go:132`–`156` |

**结论**：**不暂挂 VP-008 `go`**。

**约束**：若实施期从「追加模块 ID」改为改动装配语义或默认集**结构**，须按 freshness / `go` 规则暂停与复核。

### 2 · Admin 类 freshness（`I-038-005`）

区间 = `0c29c08`（VP-037 激活基线 · `VRev-095` PASS）→ HEAD `7e5ce891`。

| 域 | 结果 |
|----|------|
| 协议 pin（`v2.9.0` / `81aa1d8`） | **PASS**（零变更） |
| 依赖锁（`go.mod` / `go.sum` / `package.json` / pnpm 锁） | **PASS**（`go.mod`/`go.sum` 零变更；`package.json` 仅 +1 行 `typecheck` 脚本，无依赖增删） |
| 迁移台账（`internal/store` + `modules/**`） | **PASS**（零变更） |
| Profile 默认集与装配（`kernel/**`） | **PASS**（零变更） |
| provenance（`protocol/upstream/`） | **PASS**（零变更） |
| 工作树 `apps/**` | **PASS**（无 staged/unstaged 变更） |

区间内 `apps/**` 变更（39 文件，+5288/−479）全部为 VP-037 已审结目（R2/R3/R4/R6 与整改子目标 GOAL-008/009/010/011）。`go.mod` 无 redis / kafka / amqp / rabbit / gorm / ent 命中。

**结果**：**freshness PASS，不暂挂 `go`**。

### 3 · 工作区与 Root 建立

| 项 | 值 |
|----|----|
| workspace_id | `workspace-038-batch-operations-and-job-center` |
| canonical 范围 | `docs/workspaces/workspace-038-batch-operations-and-job-center/` |
| vision_role | `delivery` |
| `plan_refs` / `primary_plan` | `VP-038-batch-operations-and-job-center` |
| Root | `GOAL-001-batch-operations-and-job-center`（`parent: null`） |
| Root 初始状态 | `active · 0/5`（纲领 R1→R5） |

slug 由用户 2026-09-19 确认，非静默占位。

## 未选方案

| 方案 | 为何未选 |
|------|----------|
| B：新建 `admin.jobs` 但不进任何默认集 | 沿用 `channel.telegram` / `biz.digital-offer` 先例，最保守；代价是默认 admin Profile 看不到作业中心，需 `modules.list`/preset 显式启用，与「Admin 体验增强」的产品意图不符 |
| C：不新建模块，挂在既有模块（如 `admin.system-monitoring`） | Job 运行时是 Profile 无关的跨模块设施，挂在可观测面会把「异步作业管理」与「系统监控」耦合；且违背「作业中心是独立产品面」的意图 |

## 影响

- VP-038 R1 前仍须关闭 `I-038-001`～`003`（required）；`I-038-003` 同时承接 `V-F126`（首波批量操作分母）。
- 本轮**未**改动 `apps/**`；`admin.jobs` 模块的建立属 R2/R3 实现范围。
- 不改变 Charter `primary_workspace`；不改 Charter 目的/边界/非目标或 `vision_id@version`（仍 `@0.4.0`）。
