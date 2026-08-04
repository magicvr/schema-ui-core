---
id: r4-c1-capability-inventory
doc: evidence
goal: GOAL-005-r4-full-module-migration
date: 2026-08-05
status: recorded
---

# R4 C1 能力与边界事实盘点

## 盘点结论

本附件记录 C1 阶段截至 2026-08-05 的代码事实和迁移边界。它是
`R4-I001` 的核验材料，不把尚未形成决策的范围写成已冻结方案；R4-I002、
R4-I003、R4-I004 仍需决策或契约证据后才能关闭。

## 能力映射

| 能力/资料面 | 当前所有权或注册事实 | C1 处理结论 | 关键证据 |
|---|---|---|---|
| HTTP server、health、readyz、auth/session、accounts | Composition Root 调用中心 `handler.Register`; core 路由在中心 handler | 保留为 core composition boundary，需在 C2 明确其与模块 provider 的边界 | `apps/api/internal/composition/composition.go:82-111`; `apps/api/internal/handler/health.go:20-35` |
| Manifest route、navigation、page projection | Manifest 聚合器从中心 `adminModules` 和静态 manifest 投影；Web 从 manifest 投影导航 | 作为统一 provider 聚合的现有基线，C2 需移除一方模块的硬编码 owner 映射 | `apps/api/internal/manifest/manifest.go:38-63,67-159`; `apps/web/src/app/navigation.ts:137-152` |
| Schema render | 中心 handler embed schema，并用 `schemaDocumentsForPlan` 的 owner map 过滤 | C2/C3 迁移到模块拥有的 Schema provider；当前 owner map 是待删除中心特例 | `apps/api/internal/handler/schema.go:13-18,36-58` |
| Users | `admin.users` 声明 page/navigation/permission；HTTP Resource 和 Store 仍由中心 handler/store 提供 | C3 迁移，行为与协议保持兼容 | `apps/api/internal/kernel/profile.go:99-100`; `apps/api/internal/handler/users.go:34-51`; `apps/api/internal/store/users.go:19-52` |
| Roles | `admin.roles` 声明 page/navigation/permission；HTTP Resource 和 Store 仍由中心 handler/store 提供 | C3 迁移，行为与协议保持兼容 | `apps/api/internal/kernel/profile.go:99-100`; `apps/api/internal/handler/roles.go:21-39`; `apps/api/internal/store/roles.go:16-60` |
| Settings | 已有 `settings` module 包入口，但入口仍委托中心 `handler.RegisterSettings` | C2/C4 迁移为真实 module-owned provider | `apps/api/internal/modules/settings/module.go:9-14`; `apps/api/internal/handler/settings.go` |
| Activity/operation query | 已有 `activity` module 包入口，但入口仍委托中心 `handler.RegisterActivity`; operation query 只读 | C2/C4 迁移为真实 module-owned provider；operationlog 写入保持横切基础设施 | `apps/api/internal/modules/activity/module.go:9-14`; `apps/api/internal/handler/operations.go:14-26,54-85`; `docs/architecture/module-architecture.md:105` |
| Operationlog persistence | `RecordOperation` append-only；业务写入后 best-effort 记录，失败不回滚业务写入；未发现 retention/archival contract | R4-I004 必须明确保留当前语义及 retention 边界，或形成变更决策和新测试 | `apps/api/internal/store/operations.go:44-62,75-119`; `apps/api/internal/handler/users.go:273-303`; `apps/api/internal/handler/roles.go:228-256` |
| Records product CRUD | migration `0006 records_retire` 删除表、权限和菜单；当前无 `admin.records`、`/api/records`、handler 或 fixture；历史 operation-log 事件保留 | R4-I003 信息冲突仍开放；不可把退役事实解释成 VP-003 的最终范围决定 | `apps/api/internal/store/migrate.go:291-310`; `apps/api/README.md:109-110`; `attachments/r4-initial-boundary-scan.md:37-47` |

## Provider 缺口

`kernel.Module` 当前提供 capability 依赖、贡献 key 和生命周期 hooks，但没有
结构化 HTTP、Schema、Authorization、Navigation、Manifest、Persistence provider
字段。现有冲突校验只覆盖字符串贡献 key。C2 进入条件是冻结框架无关 provider
形状、Plan 消费顺序、冲突规则、依赖注入和生命周期失败清理测试。

证据：`apps/api/internal/kernel/module.go:28-56,322-386`、
`apps/api/internal/kernel/lifecycle.go:17-65`。

## Operationlog 事实

当前业务写入和 operationlog INSERT 不在同一可核对事务契约中；日志写入失败只
记录服务日志。已有 append、排序、事件 CHECK、迁移保留和 reopen 测试，但未发现
日志失败注入、原子回滚、retention duration、归档表或归档恢复测试。该事实不
替代 R4-I004 的保留/变更决策。

证据：`apps/api/internal/store/operations_test.go:11-246`、
`apps/api/internal/store/migrate.go:222-238,249-286,348-382`。

## 未关闭门禁

- `R4-I001`：能力映射已形成，但“全部现有 Admin 能力”的 C1 核验仍需审计确认。
- `R4-I002`：provider gap 已核实，contract 尚未冻结。
- `R4-I003`：VP-003 与 `records_retire` 的范围冲突必须由用户或 canonical 决策解决。
- `R4-I004`：operationlog 失败语义和 retention 边界必须形成明确决策与证据。
