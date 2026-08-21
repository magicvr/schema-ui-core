---
id: D-002
goal: GOAL-016-r3-s09-data-permission
title: 方案冻结：数据权限（行级/数据范围）设计（S1）
date: 2026-08-15
status: accepted
parent: GOAL-016-r3-s09-data-permission
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# D-002 · 方案冻结（S-09 数据权限 · 行级/数据范围）

> 依据：I-011-001 §4 S-09、§7 协议对照口径；GOAL-016 00-meta 边界与 I-001~I-004；A-002 016-F-003（B-10 依赖裁定）。
> 证据：handler/resources.go（资源工厂与 requirePermission）、handler/users.go（ResourceEntity.List 边界）、authsession/users_repository.go（SQL where 组装）、kernel/profile.go（BuiltinModules/profileDefaults）、kernel/provider.go（DefaultNavigationOrder）、protocol-inventory-v2.7.0.md（D-PERM / ADR-0004）、docs/schemas/capability-registry.json。

## 1. 数据范围模型（I-001 闭合 · I-004 裁定）

- **v1 作用域类型**：all（全部数据）/ self（本人数据）。**组织范围本波不纳入**（I-004 裁定：B-10 组织/部门/岗位未立项，依赖不成立；作用域类型枚举留扩展位，后续目标按 B-10 立项再扩 org）。
- **语义**：范围按**用户 × 资源**赋值；合成规则 = 用户显式赋值优先，未赋值 → 资源登记的 default_scope（登记时显式声明，默认 all——与现状一致零回归，文档化残余：默认 fail-closed 需全量范围迁移，风险高，不采纳）。
- **继承**：沿用 RBAC 继承面（PolicyID→rolesForPolicy 已把管理键带给角色集）；数据范围本身为**用户级赋值**，角色级范围继承 v1 不实现（文档化残余，后续扩展）。
- **owner 列**：登记资源时声明 owner_column（白名单校验，防注入）；self → 查询附加 WHERE owner_column = :actorID。

## 2. 过滤下推（I-002 闭合）

- **注入点**：resourceHandler.list() 构造 resourceFilter 时（resources.go:363-370，ExtraQuery 白名单先例 dictKey）附加 Scope 约束 → 经 ResourceEntity.List(filter) 单一 store 边界下传 → 各 repository where/args 组装（users_repository.go:365 usersWhere 同款）。**不落 middleware**（权限/查询面证据：检查层在 handler，SQL 在各 repository）。
- **机制**：新可选接口 RowScopeProvider { ScopeFor(userID, resource string) (ScopeConstraint, error) } 注入资源工厂（nil = 未启用，行为逐字节不变——captcha 模式先例）；工厂仅对**已登记资源**附加约束。
- **覆盖全部行访问路径（A-004 F-001 响应）**：list 之外，Get（按 id 单行）/ Update / Delete / BatchDelete 同样按 owner 约束校验——self 作用域下行不属本人 → 404（Get/Update/Delete，不泄露存在性）/ 跳过（BatchDelete 仅删本人行）；Create 时 self 作用域资源 owner_column 默认写入当前 actor。导出面（data-transfer）不经 resourceHandler.list：**登记规范要求**——登记资源时须同步评估导出路径并施加同约束（v1 无登记资源，无暴露面；导出面接入点登记为后续目标必办）。
- **ScopeAware 契约（A-004 响应）**：已登记资源的实体必须实现可选 ScopeAware 接口（消费 filter.Scope 的 where 段），工厂登记时校验（未实现 → 拒绝登记，配置面 fail-closed）。因此「登记新资源」= 实体实现 ScopeAware + 登记行两项动作，**非零代码**（修正原表述）。
- **v1 落地目标**：交付框架 + 管理面 + 工厂接线 + 测试资源验证过滤；**不登记任何生产资源**（领域资源 S-13/S-14 未立项，无真实 owner 语义；登记入口留给后续目标）。文档化：本波无生产数据行被过滤（能力面交付），审计意见把关。

## 3. 端点与权限

| 端点 | 门禁 | 说明 |
|------|------|------|
| GET /api/data-permission/policies | data-permission.read | 资源登记列表（resource / owner_column / default_scope / enabled） |
| PATCH /api/data-permission/policies/{resource} | data-permission.write | 登记/更新资源策略（owner_column + **default_scope 必填**，省略 → 400 INVALID_PATCH_FIELD）；审计 data-permission.policy-update |
| GET /api/data-permission/scopes?userId= | data-permission.read | 用户范围赋值列表 |
| PATCH /api/data-permission/scopes | data-permission.write | 批量赋值（userId × resource × scopeType all/self）；审计 data-permission.scope-update |

- 权限键：data-permission.read / data-permission.write（PolicyAdmin）；被作用域资源的既有权限键不变（过滤是查询面增强，不改变授权入口语义）。
- 页面 data-permission：策略登记表 + 用户赋值表（schema 驱动，search-table + 表单）；导航 menu_data_permission（visibleWhen features.menu_data_permission，Permission data-permission.read）。
- Profile：admin.data-permission 进入 **admin 默认集**（S 系列内容扩展先例 file-library/data-dictionary）；mvp/demo 不含（profile.go ProfileAdmin 段插入 + GOAL-016 注释）；DefaultNavigationOrder 尾部追加。

## 4. 迁移

- **0027**（admin.data-permission）：data_scope_policies（resource TEXT PK ｜ owner_column TEXT NOT NULL ｜ default_scope TEXT NOT NULL CHECK in (all,self) ｜ enabled INTEGER NOT NULL DEFAULT 1 ｜ updated_at）+ user_data_scopes（user_id TEXT NOT NULL ｜ resource TEXT NOT NULL ｜ scope_type TEXT NOT NULL CHECK in (all,self) ｜ updated_at ｜ PK(user_id, resource)）。
- **0028**（core.operationlog）：CHECK + data-permission.policy-update / data-permission.scope-update。
- compiled/persistence.go 注册 datapermission migration provider。

## 5. 协议对照（本地鉴权扩展，不声明协议覆盖）

- D-PERM（权限继承/intent）primary 能力已覆盖 UI 显隐与鉴权意图（permissions-inheritance fixtures）；ADR-0004 row-level-scope 为 UI 表格/行操作概念，**非后端数据行鉴权**；capability-registry 无 data-scope 语义。
- 处置：S-09 = **本地鉴权扩展**（管理面 + 查询面增强），不新增协议 capability、不改协议 pin（v2.8.0）、不改 Manifest 装配语义；呈现自由不适用（鉴权语义非呈现）；留痕于本决策与 03-audit。

## 6. 测试与验证

- 工厂：无 RowScopeProvider 时逐字节不变；已登记资源 self → 附加 owner 约束；default_scope 生效；未登记资源不受影响。
- scope 服务：赋值 CRUD、合成规则、owner_column 白名单（防注入）、非法 scope_type 拒绝。
- 端点：门禁 401/403、审计事件、分页。
- 组合根：admin 权限 **24→26**、导航 **12→13**（composition_test.go L465 当前 24/12）；迁移 26→28。
- web：fixture/schema-keys/s5-denominator/e2e 回归（默认无登记资源零影响）+ smoke.sh 页面集 + data-permission 页。

## 7. 未选方案（留痕）

- middleware 层过滤：权限/查询证据显示检查层在 handler、SQL 在各 repository，middleware 无法访问资源级过滤语义。
- 角色级范围继承：v1 用户级赋值（RBAC 继承面已覆盖管理键），角色级 scope 后续扩展。
- 组织（org）范围：B-10 依赖未立项（I-004 裁定，本波不纳入）。
- 默认 fail-closed：需全量既有数据范围迁移，回归风险高；显式登记 + default_scope 声明 + 审计把关替代。
- 生产资源本波登记：领域资源未立项，无真实 owner 语义（能力面交付，登记留后续目标）。

## 8. S2 实现清单

1. migration 0027/0028 + compiled/persistence.go 注册。
2. modules/dataprovider（或 admin.data-permission）：provider（Descriptor + 路由 + 权限 + fragment）+ store（策略/赋值 CRUD）+ schema（data-permission 页）。
3. handler：RowScopeProvider 接口 + resourceFilter.Scope 扩展 + /api/data-permission/* 端点。
4. 工厂接线：composition 注入 RowScopeProvider（nil = 未启用）。
5. web：data-permission 页 + i18n + fixture。
6. 测试：工厂/scope/端点/组合根/web fixture。
7. A-005 recommended：ScopeAware 强制点落在策略 PATCH（非 registerResource）；Create 强制 owner 覆盖（self 资源）；测试覆盖全部行访问路径（Get/Update/Delete/BatchDelete + 导出面评估）。
