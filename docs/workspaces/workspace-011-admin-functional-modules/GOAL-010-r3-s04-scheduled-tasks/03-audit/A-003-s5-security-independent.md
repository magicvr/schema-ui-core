---
id: A-003
goal: GOAL-010-r3-s04-scheduled-tasks
source: independent
date: 2026-08-14
scope: S5 关门 · 安全/数据门禁（admin.scheduled-tasks vs D-002 冻结方案）
verdict: conditional
auditor: grok-build
audit_type: close-out
status: recorded
parent: GOAL-010-r3-s04-scheduled-tasks
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-003 · independent 安全/数据审计（S-04 实现）

## 范围与区间

- **auditor**：grok-build（independent cross-audit · grok-4.6）
- **type**：close-out / security-data gate
- **workspace**：`workspace-011-admin-functional-modules`（`root_goal` = `GOAL-001-admin-functional-modules`；`canonical_scope` 已核对；`shared_materials_catalog: none`）
- **covered**：
  - `apps/api/internal/modules/scheduledtasks/store/cron.go`（5 字段解析 + Next）
  - `apps/api/internal/modules/scheduledtasks/store/repository.go`
  - `apps/api/internal/modules/scheduledtasks/scheduler.go`
  - `apps/api/internal/handler/scheduledtasks.go`（及工厂 `resources.go` 门禁 / PATCH 解码 / 排序白名单）
  - `apps/api/internal/modules/scheduledtasks/migration/migration.go`（0021）
  - `apps/api/internal/modules/operationlog/migration/migration.go`（0022 CHECK 扩展）
  - `apps/api/internal/modules/scheduledtasks/provider.go` + `schema/*.json` + `manifest/fragment.json`
  - `apps/api/internal/kernel/profile.go`、`apps/api/internal/composition/composition.go`
  - 计划契约 `01-decision/D-002-s1-plan-freeze.md`（及 D-001 / D-003）
- **excluded**：端到端浏览器手测、生产多实例部署实测、其它工作区上下文、实现代码改动
- **信息项**：I-001 / I-002 / I-003 均已 closed（S1）；本 scope 无到期未关闭 required 信息项；无共享资料引用

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 任务 CRUD / 手动触发 / 运行历史均 fail-closed | 工厂 `PermissionRead/Write = tasks.read/tasks.write`（scheduledtasks.go:186–187、201–202）；`list/detail` 走 `requirePermission(read)`（resources.go:286–288、514–518）；`create/update/delete/batch-delete` 走 `requirePermission(write)`（resources.go:469–470、531–532、562–563、589–590）。`POST /{id}/run` 显式 `tasks.write`（scheduledtasks.go:211–214）；`GET /{id}/runs` 显式 `tasks.read`（234–237）。匿名无 Bearer → Middleware 401（auth.go:386–395）；已认证无 key → 403（resources.go:231–241）。测试：`TestScheduledTasksPermissionGates` 覆盖 list/create 403 与 list 401 |
| 手动触发需要 tasks.write | scheduledtasks.go:212–214；与 D-002 §4 一致 |
| 非法/越界 cron 在写入端拒绝 | `ParseCron` 要求恰好 5 字段（cron.go:26–30）；`*/0`、空列表项、越界数字/范围均 error（42–88）。Create/Update 映射 400 `INVALID_CRON`（scheduledtasks.go:67–69、99–101）。单测拒绝 `60 * * * *`、`* 24 * * *`、`* * * * 7`、`*/0` 等（cron_test.go:17–24） |
| Next 含当前槽 + 分钟去重防双跑 | `Next` 从 `from.Truncate(minute)` 起搜（cron.go:108–116）；tick 在 `next.After(now)` 时跳过，`lastRun[task.ID]==slot` 时跳过，先记 slot 再 `Execute`（scheduler.go:84–98）。`TestSchedulerExecutesDueTasks` 同槽二次 tick 仍 1 行 |
| 5 年搜索有界 | `limit := start.AddDate(5, 0, 0)` + `t.Before(limit)`（cron.go:110–111）；无无界循环 |
| 运行行记录 status/detail；handler 失败写 failed | `Execute` 在 handler error 时 `detail=err.Error()` 且 `status=failed`，否则 `ran`（scheduler.go:112–126）；`RecordRun` 失败 slog 并返回（127–130）。noop 成功路径有测试 |
| 无每任务 goroutine；Start 幂等 | `Start` 用 `sync.Once` 启一个 loop（scheduler.go:52–57）；`Execute` 同步调用；tick 串行扫任务 |
| PATCH 缺省保持原值 | Update 对 cron/name 空串回落 existing；description/handler 按 key 是否在 body；enabled 用 `boolField(..., existing.Enabled)`（scheduledtasks.go:90–114）。工厂 PATCH 缺 key 不写入 body（resources.go:425–429） |
| 删除级联 run 行 | `task_runs.task_id REFERENCES scheduled_tasks(id) ON DELETE CASCADE`（migration.go:31）；`DeleteTask` 只删父行（repository.go:181–192）；连接 `PRAGMA foreign_keys=ON` 并断言（store/migrate.go:213–224） |
| UNIQUE key | DDL `key TEXT NOT NULL UNIQUE`（migration.go:20）；`CreateTask` 将 unique 错误映射 `ErrKeyTaken` → 409 `TASK_KEY_TAKEN`（repository.go:153–155；scheduledtasks.go:279–280） |
| 审计事件与 0022 CHECK 一致 | 常量 `scheduled-tasks.create/update/delete`（operationlog/repository.go:43–45）；handler 只写这三事件（scheduledtasks.go:82、118、131）；0022 CHECK 精确追加这三项（operationlog/migration/migration.go:101） |
| SQL 参数化 + sort 白名单 | 仓库 `?` 绑定；`ListTasks` sort 仅 `key/name/updatedAt`，order 仅 desc 否则 asc（repository.go:87–94）。工厂先校验 `SortFields`（resources.go:291–300；scheduledtasks.go:180） |
| 0021/0022 与 1..22 连续 | 0021 建表（scheduledtasks/migration/migration.go:17–51）；0022 `rebuildOperationLog` rename→create→INSERT SELECT→drop→index（operationlog/migration/migration.go:226–246）。`TestCompiledMigrationCatalogOwnership` 断言 catalog[i].Version==i+1 且 22 条以 scheduled_tasks / operation_log_tasks 结尾 |
| Profile 仅为 admin 内容扩展 | `profileDefaults[ProfileAdmin]` 追加 `"admin.scheduled-tasks"`（profile.go:74–76）；mvp/demo 未加入；`ResolveProfile` 逻辑未改；`plan.HasModule` 才挂 provider（composition.go:238–240）。admin 权限 18→20、导航 10→11，mvp 8/5 不变（composition_test.go:449–459） |

## 对照成功标准（本 scope）

| 标准 | 结论 |
|------|------|
| 全部端点 tasks.read/write 401/403 fail-closed；手动触发 tasks.write | **满足**（代码路径完整；测试覆盖 list/create，自定义 /run 与 /runs 靠代码审查） |
| cron 拒绝非法/越界；Next 含当前槽不双跑；5 年有界 | **部分满足** — 文档化的 `*` / 数字 / `*/n` / 列表 / 范围正确；标量 `n/step` 被接受但步进被丢弃（F-001） |
| 分钟槽去重；run 行 status/detail；失败 detail；无无界 goroutine；Stop 安全 | **部分满足** — 去重/记录/单 loop 满足；Stop 非幂等（F-004）；不可调度不落 detail（F-002） |
| PATCH merge；删除级联；UNIQUE key | **满足** |
| 审计事件与 0022 CHECK 精确匹配 | **满足** |
| 参数化查询；sort 白名单 | **满足** |
| 0021 表 + 0022 重建保行；版本 1..22 | **满足** |
| Profile 仅追加 admin 默认集 | **满足**（装配语义未改） |

## Findings

### F-001 · 标量 `n/step` 被接受但步进被丢弃，导致少触发

| 字段 | 值 |
|------|-----|
| level | required |
| status | open |
| evidence | `parseCronField` 在看到 `/` 后解析 `step`（cron.go:50–59）。`base=="*"` 与 `a-b` 分支按 `v += step` 填集合（60–78）。**裸数字分支只写 `values[n]=true`，不再使用 `step`**（80–87）。因此 `0/5 * * * *` 与 `5 * * * *` 等价（仅第 0 分钟），而 `*/5` 与 `0-59/5` 才是每 5 分钟。D-002 §2 支持 `*/n` 步进；`n/step` 不在清单内，却通过 Create/Update 的 `ParseCron` 校验（scheduledtasks.go:67–69）以 201/200 入库。单测未覆盖 `0/5`（cron_test.go:8–24） |
| severity | med（非授权绕过。管理员按常见 cron 惯用写法 `0/5` 保存成功后，调度会按小时而非每 5 分钟触发，属于冻结计划「不能 skip」语义下的少触发） |

**建议修复**（二选一，须 fail-closed，不可再静默忽略 step）：

1. **拒收**（更贴 D-002）：裸数字后带 `/step` 返回 error → 400 `INVALID_CRON`（仅允许 `*/n` 与 `a-b/n`）。
2. **按 start/step 实现**：`for v := n; v <= max; v += step { values[v] = true }`（`0/5` ≡ `*/5`）。

补测：`0/5 * * * *`、`5/2 * * * *` 必须拒绝或展开为完整集合；现有 `*/5` / `0-59/5` 回归不得坏。

### F-002 · 5 年窗口内无匹配时不记录「不可调度」detail

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | D-002 §2：「搜索失败视为任务不可调度（记录 detail）」。`tick` 在 `!ok \|\| next.After(now)` 时直接 `continue`（scheduler.go:90–93），不写 `task_runs`。永不可能匹配的表达式（如 `0 0 31 2 *`，2 月 31 日）会每 30s 空转且无 run 行。闰年路径有测试（cron_test.go:69–77），无「窗口内无匹配 → 落盘」断言 |
| severity | low（不造成双跑或越权；可观测性/计划偏差。每 tick 对无匹配表达式做最多约 2.6e6 次分钟扫描，有界） |

**建议修复**：对 `!ok`（窗口内无匹配）写一条 `status=failed`、`detail` 标明 unschedulable / 5-year bound，并按 task+reason 去重以免每 30s 刷行。`next.After(now)`（尚未到期）保持跳过。

### F-003 · 未注册 handler 静默回落 system.noop 并记 `ran`

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `Execute`：`handler, ok := s.handlers[task.Handler]; if !ok { handler = s.handlers["system.noop"] }`（scheduler.go:105–108）。Create 接受任意非空 `handler` 字符串，缺省才是 `system.noop`（scheduledtasks.go:70–73），写入端不对照 `HandlerKeys()`。拼写错误的 handler 会以成功 noop 落 `ran`，无 detail |
| severity | low（v1 仅 noop，无 RCE/注入面。运行历史 status 不能区分「真正 noop」与「未知处理器」） |

**建议修复**：未知 handler 不要回落成功路径；记 `status=failed`、`detail=unknown handler`（或 Create/Update 直接 400）。后续接入业务处理器时同样 fail-closed。

### F-004 · `Stop` 非幂等；分钟去重仅进程内内存

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `Stop` 直接 `close(s.stop)`（scheduler.go:73–75），无 `sync.Once`；二次 Stop panic。生产 `New` 只 `Start`、不 Stop（provider.go:33–36），测试也未调 Stop，故当前路径未触发。`lastRun` 为进程内 `map`（scheduler.go:34、48、94–97）：进程重启后同一分钟槽可再执行一次；多实例会各自执行。与 A-002 self 建议及 D-002 §3 单实例 best-effort 文档一致 |
| severity | low（Stop 为测试/关停 API；双跑窗口为已文档化的单实例残余） |

**建议修复**：`Stop` 用 `sync.Once`（或等价）只 close 一次。多实例/重启双跑保持 D-002 文档化残余即可；若以后要硬去重，需持久化「task_id + minute slot」唯一约束。

## 必改项汇总

- **required / 必改**：F-001（标量 `n/step` 语义错误 / 未拒收）
- **recommended**：F-002、F-003、F-004（不阻断授权、SQL、迁移、审计 CHECK、Profile 项）

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-002 self（S2–S4，verdict pass） | **同意**授权 fail-closed、0021/0022、PATCH merge、级联删除、UNIQUE、审计 CHECK、Profile 内容扩展、Next 含当前槽 + 内存去重。**不同意** cron「数字/步进」已完全正确——self 未审 `n/step` 与 `*/n` 的分叉。F-004 与 A-002 已写的多实例/重启建议同向，本意见保持 recommended |
| A-001 self（S1） | 自研 5 字段、单实例 best-effort、0022 CHECK、admin 默认集与实现一致 |

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional** — 安全核心（授权 fail-closed、手动触发 write、SQL 参数化与 sort 白名单、UNIQUE/级联、审计事件与 0022 CHECK 对齐、迁移 1..22 连续、Profile 仅 admin 内容扩展）有代码与测试证据。不可无条件放行本门禁：F-001 使部分会被 API 接受的 cron 按错误集合调度（少触发）。

建议 `/govern`：响应本意见；先修 F-001（拒收或实现 start/step，并补单测）；F-002/F-003/F-004 可同波或作维护。闭合 required 前不要把本目标标为 `done`。

### 声明

本意见 `source: independent`，**不修改**目标 `status` / `progress` / goal-tree / 方案正文或实现代码；响应与状态变更由 `/govern` 与用户裁决处理。
