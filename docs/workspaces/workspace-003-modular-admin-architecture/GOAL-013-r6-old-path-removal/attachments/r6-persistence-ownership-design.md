# R6 Persistence 所有权设计（R6-I002 可实施基线）

> 有界设计：冻结平台 vs 模块 ownership、0001-0008 descriptor 归属、CollectPersistence
> 生产接线顺序、fresh seed vs contribution-driven reconcile。**不**据此宣称 VP 退出
> #2/#3/#5 已取证；取证须等实现 + 审计。

## 1. Ownership 分层

| 层 | 职责 | 拥有 | 不拥有 |
|----|------|------|--------|
| **平台（platform runner）** | DB 连接生命周期（Open/Close/Ping）、事务 runner、迁移 **执行器**（按 catalog 事务应用 + ledger version/name/checksum 校验 + tombstone + 缺口/漂移/未知 fail-closed）、ledger 读写 | 迁移执行引擎、schema_migrations 账本 | 任何领域 SQL、seed、operation_log 写逻辑、业务仓储 |
| **core.auth-session** | 账号/认证/RBAC 域 | 0001（users、refresh_tokens）、0002（RBAC：roles、user_roles、role_menu_items、permissions/menus）descriptor + 对应仓储（用户 CRUD/密码、角色、RBAC seed 最小集） | 业务资源 surface（HTTP/Schema 归 provider） |
| **core.operationlog** | 操作日志追加 | 0004/0005/0008（operation_log）descriptor + `RecordOperation` writer | Activity 查询/UI（归 admin.activity） |
| **admin.settings** | 站点设置域 | 0007（site_settings）descriptor + 仓储 | 横切 writer |
| **core.persistence（历史）** | 仅承载历史 Records 迁移 | 0003（records_persist）、0006（records_retire）descriptor（immutable 历史；records 已退场，不恢复） | 当前产品 CRUD |
| **admin.users / admin.roles** | 资源 surface（Provider） | HTTP/Schema/Auth/Nav/Manifest 贡献；领域逻辑委托 auth/RBAC 仓储 | 直接持表（除非后续独立仓储） |

## 2. 0001-0008 descriptor 归属（CollectPersistence 输入）

| 迁移 | 名称 | 表/动作 | moduleID | 说明 |
|------|------|---------|----------|------|
| 0001 | r2_baseline | users、refresh_tokens | `core.auth-session` | 账号域 |
| 0002 | rbac_expand | roles、role_menu_items、user_roles、permissions | `core.auth-session` | RBAC 域 |
| 0003 | records_persist | records | `core.persistence`（历史） | immutable legacy |
| 0004 | operation_log | operation_log | `core.operationlog` | |
| 0005 | operation_log_expand | operation_log 重建 | `core.operationlog` | |
| 0006 | records_retire | records DROP | `core.persistence`（历史） | 退场 |
| 0007 | site_settings | site_settings | `admin.settings` | |
| 0008 | operation_log_settings | operation_log 重建（settings.update） | `core.operationlog` | |

约束：0001-0008 **不重编号、不改名、不改 checksum、不改 Apply 语义**（冻结包 §4.1）；
仅迁移 **归属登记**（从 store 硬编码 `core.persistence` 改为各模块 descriptor）。
后续新迁移只允许 Provider `CompiledPersistence()` 追加全局版本。

## 3. CollectPersistence 生产接线顺序

1. **收集**：composition 对全部 compiled provider 调 `CompiledPersistence()` →
   `kernel.CollectPersistence` → 校验唯一/连续/checksum/tombstone → compiled-global
   catalog（版本排序）。
2. **应用**：`store.Open`（平台）改收 catalog 而非硬编码 `compiledMigrations`；
   逐条事务应用未应用迁移，ledger 校验（未知已应用、缺口、checksum drift、partial
   baseline fail-closed）。`ModuleID: "core.persistence"` 硬编码从 store 移除。
3. **追加**：新迁移仅经 Provider `CompiledPersistence()`；禁用模块迁移仍执行（数据保留）。
4. **验证**：既有 `migrate_test.go`（unknown/gap/drift/partial/recovery）改为以
   catalog 输入驱动；新增「composition catalog = store 应用 catalog」一致性测试。

## 4. fresh seed vs contribution-driven reconcile

| 路径 | 触发 | 内容 | 所有权 |
|------|------|------|--------|
| **fresh bootstrap** | 空库首次 Open | 最小 admin 用户 + system roles（admin/editor/viewer） | `core.auth-session`（不可覆盖用户字段） |
| **versioned reconcile** | 每次 Open（幂等） | 消费 finalize 后 Authorization/Navigation 贡献 → 补齐系统权限/菜单；只填明确归属 system keys，不覆盖用户字段、不删除用户数据 | RBAC reconcile（core.auth-session / 贡献驱动） |

- 现 `store.seedRBAC` 的权限/菜单硬编码列表改为读 Provider `PermissionContribution`/
  `NavigationContribution`（finalize 后）。
- fresh 与 reconcile **分离**：bootstrap 只在空库；reconcile 版本化、幂等、可重放。

## 5. 实施顺序（设计冻结后，切片推进）

1. 内核已有 `MigrationContribution`；把 0001-0008 Apply+DDL 迁入各模块
   `migration/` 包，产出 descriptor。
2. `store/migrate.go` 收窄为平台 runner（收 catalog）。
3. composition 收集 + 接线 `store.Open(catalog)`。
4. 各模块 `CompiledPersistence()` 返回对应 descriptor（core 模块也 provider 化
   Persistence）。
5. seed/reconcile 改造为贡献驱动。
6. 全量回归 + Grok 复审（C6.2 关门）。

**边界**：本设计不新建独立 `internal/platform` 目录（避免过度工程）；平台 runner 收敛
在 `store` 包但仅剩通用执行引擎；领域 SQL 迁出到模块仓储。若实现中发现 runner 需独立，
再评估是否提为 `internal/platform`。
