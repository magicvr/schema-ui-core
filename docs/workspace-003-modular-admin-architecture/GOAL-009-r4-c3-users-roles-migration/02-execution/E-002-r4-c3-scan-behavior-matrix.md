---
id: E-002-r4-c3-scan-behavior-matrix
doc: execution-entry
goal: GOAL-009-r4-c3-users-roles-migration
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-002 · R4-C3 迁移扫描与行为矩阵（C3.1）

## 中心注册 / 特例状态（C3-I001 verified）

| 特例 | 位置 | C3 动作 |
|------|------|---------|
| 中心业务 Register | `handler/health.go` `Register`：`admin.users`→`registerResource(usersResource)`、`admin.roles`→`registerResource(rolesResource)` | 改为 provider 自注册（冻结 §7 步骤 2/3） |
| Schema owner map | `handler/schema.go` `schemaDocumentsForPlan` 硬编码 `owners{users,roles,settings,activity}` + `fixtures/schema/*.json` | schema 文档归 module-owned provider；移除硬编码 owner map |
| Manifest adminModules | `manifest/manifest.go` `ForModules` 硬编码 `adminModules{users,roles,settings,activity}` | fragment 由 provider 贡献；移除硬编码 map |
| 通用资源工厂 | `handler/resources.go` `registerResource` + `Resource`/`ResourceEntity` | users/roles provider 复用工厂产出 HTTP contribution |
| 权限/seed | `store/seed.go` seedRBAC + handler `requirePermission` | 权限 key 由 provider Authorization contribution 声明，seed/reconcile 一致性（冻结 §5） |

## 保留行为矩阵（C3-I002 verified，枚举自冻结 §7 兼容清单 + 现有测试）

| 维度 | 保留内容 |
|------|----------|
| HTTP | users/roles 五路由 list/create/detail/update/delete 的成功 status、错误分类、`{error,message}` 信封、body 上限 4KB、pageSize≤100、q/sort/order 校验 |
| 授权 | `users.read`/`users.write`、`roles.read`/`roles.write`/`roles.assign` 权限键门禁；无/无效/过期 access → 401；缺权限 → 403 |
| 语义 | users 敏感字段隔离（`password_hash` 永不出响应）；角色分配双写（legacy JSON + `user_roles`，不隐式建角色）；self/last-admin 保护；roles system 角色不可改删、in-use 保护 |
| 持久化 | `0001`-`0008` 迁移账本、users/roles 表、user_roles 关联、checksum 不变；module-owned 迁移只追加全局版本 |
| operationlog | 业务成功后 best-effort append；失败记服务日志不翻转业务；`users.*`/`roles.*`/`auth.*` 事件保留；Activity disabled 仍写入 |
| Schema | `users`、`roles` page/resource/action ID 不变；Renderer 通用渲染 |
| Manifest | users/roles 菜单、页面投影行为不变；登录前无 secret |

## 迁移路径（冻结 §7 顺序）

1. 建立 `admin.users`/`admin.roles` provider 元数据 + 无发布 contract tests。
2. typed provider 生成 surface，在测试中与现有中心输出兼容比较（不永久双注册）。
3. 切换 composition 消费 `RegisterContributions`/`CollectPersistence`；移除中心
   users/roles 分支、Schema owner map、Manifest `adminModules`。
4. 静态扫描 + 运行测试证明中心业务 Register、全局 Schema fixture 占用与永久双路径
   已删除；tombstone 只保留数据兼容所需事实。
