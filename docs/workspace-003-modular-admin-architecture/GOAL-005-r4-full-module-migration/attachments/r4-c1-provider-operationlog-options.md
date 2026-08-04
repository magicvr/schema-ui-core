---
id: r4-c1-provider-operationlog-options
doc: proposal
goal: GOAL-005-r4-full-module-migration
date: 2026-08-05
status: draft
decision_state: pending_user
---

# R4 C1 Provider 与 operationlog 方案选项

本附件是 C1 的待裁决材料，不是 accepted decision。它承接当前已核实的
`Module`/`Plan`、组合根、operationlog 和迁移事实；用户确认后才可转为 D-003
并关闭相应信息门禁。

## 一、Provider contract 候选

### 推荐候选：框架无关 Provider + Plan-owned Registrar

建议在现有 `kernel.Module` 元数据之上增加结构化 provider 注册入口，但不把
Fx 类型带入 kernel 或模块公共 API。概念形状如下，名称和精确类型仍待 C2
方案冻结：

```go
type Provider interface {
    Descriptor() Module
    Register(context.Context, Registrar) error
}

type Registrar interface {
    HTTP(RouteContribution) error
    Schema(PageContribution) error
    Authorization(PermissionContribution) error
    Navigation(NavigationContribution) error
    Manifest(FragmentContribution) error
    Persistence(MigrationContribution) error
}
```

每条 contribution 必须带稳定 owner/module ID 和 key；Registrar 先写入临时
`ContributionSet`，完成全部注册后再统一校验和 finalize。推荐消费顺序：

1. 解析 Profile/`modules.enabled`，得到唯一 immutable `kernel.Plan`。
2. 校验启用模块与 provider descriptor 一一对应、版本/API range/依赖一致。
3. 按 Plan 拓扑顺序调用 provider 的 `Register`，只写临时集合。
4. 对 route/page/schema/permission/navigation/config、Manifest fragment、
   migration version/name/checksum 执行全局冲突和完整性校验。
5. 校验失败时丢弃临时集合，不创建或发布部分 HTTP/Schema/Manifest surface。
6. composition 将已 finalize 的结果适配为 `http.ServeMux`、store 和 Fx graph；
   `fx.In`、`fx.Out`、`fx.Option` 只出现在 composition。
7. 继续沿用现有 Plan 顺序启动、Ready 和反向 Stop；provider 注册失败或生命周期
   失败必须 fail closed 并清理已启动资源。

### 六类贡献的当前映射

| 贡献 | 当前代码缺口 | C2 最小验证 |
|---|---|---|
| HTTP | `handler.Register`/settings/activity wrapper 仍由中心控制 | 重复 method+pattern、禁用模块路由不存在、注册失败无部分路由 |
| Schema | `schema.go` 有中心 embed 和 owner map | 重复 page ID、禁用模块 Schema 404、owner 必须来自 provider |
| Authorization | permission metadata 与 `seedRBAC` 分离 | 重复 permission、seed/manifest/handler key 一致、禁用权限 fail closed |
| Navigation | profile metadata + central seed/menu + static manifest | 重复 navigation ID、权限投影与 Plan 一致 |
| Manifest | `manifest.ForModules` 有硬编码 `adminModules` | fragment/page/navigation/protocol/app identity 冲突和确定性排序 |
| Persistence | `compiledMigrations` 仍集中在 store | 全局 version/name/checksum、tombstone、reconcile 和快照/恢复约束 |

### 必须保留的边界

- `core.server-registration`、`core.auth-session`、`core.operationlog` 可以继续是
  core cross-cutting capabilities，但其边界必须在 provider matrix 中显式登记。
- Persistence provider 不能自行重排全局 migration ledger；模块贡献必须汇入
  全局严格递增、不可变 checksum 台账。
- `admin.users`、`admin.roles` 的领域 Store/Resource 适配必须通过 provider
  暴露，不能继续由 composition 按 module ID 写分散 `if` 分支作为最终架构。
- Settings/Activity 当前 module wrapper 不能被误称为完整 provider；其中心委托
  只有在统一注册结果覆盖后才可移除。

### 未选方案与理由

- 模块直接返回 `fx.Option`：实现短，但违反 Fx 仅限组合根的架构边界，公共契约
  被 Fx 绑定。
- 继续按 `plan.HasModule` 增加中心分支：改动最小，但保留 central registration、
  Schema owner 和 Manifest ownership 特例，不能证明模块真正拥有六类贡献。

## 二、operationlog 边界选项

当前事实：业务写入成功后调用 append-only `RecordOperation`；失败只记录服务日志，
不回滚业务写入；没有 retention duration、archive、purge、restore contract 或
失败注入测试。

### 选项 A：保留 best-effort，R4 不引入自动 retention

- 兼容性最高，保持既有 API 响应和历史 operationlog schema。
- R4 明确承诺：不新增自动 purge/archive，不删除现有 operation_log 行；迁移、
  重启和 Activity UI 开关都不改变日志保留。长期 retention duration 作为独立
  运营/合规决策触发项，不得在 R4 中默认为永久保留。
- 需要的证据：Users/Roles/Auth/Settings 日志失败注入后业务响应仍成功；失败被
  记录；append/read/sort、迁移保留、重启保留继续通过。
- 残余风险：日志失败可能造成审计缺口；必须在 R4 close-out 明确该边界和复审
  触发条件。

### 选项 B：业务写入与 operationlog 原子提交

- 需要把业务事务和 `RecordOperation` 改为同一 `*sql.Tx`，不能只修改单条 INSERT。
- 失败会改变现有 API/业务语义，并扩大锁、重试、错误映射和并发影响。
- 需要成功共存、日志故障共同回滚、重复 ID、迁移/重启以及所有 writer 一致性
  测试；Auth/Settings 也不能保留另一种语义。

### 选项 C：引入归档/保留系统

- 需要定义 retention duration、触发机制、存储格式/位置、查询边界、purge、
  restore、失败重试、幂等和审计完整性。
- 需要新增 persistence migration、archive/restore round-trip、时间边界、失败
  重试、Activity 查询一致性和恢复测试。
- 这是显著扩大 R4 范围的方案，不能以“先做一个表”替代完整数据生命周期契约。

## 三、建议与待裁决项

工程建议是：Records 先保持 historical-only 退役；operationlog 选择 A，以兼容性
为优先并明确 R4 不新增自动清理/归档。两项建议都不是当前决策。用户需要书面
确认或改选后，编排器才会写 D-003、关闭 R4-I003/R4-I004，并冻结 C2 provider
方案。

证据：`apps/api/internal/kernel/module.go:13-56,232-390`、
`apps/api/internal/composition/composition.go:25-111,163-195`、
`apps/api/internal/store/operations.go:44-119`、
`apps/api/internal/store/migrate.go:60-90,249-288,348-384`、
`docs/architecture/module-architecture.md:27-35,41-54,77-92,103-114`。
