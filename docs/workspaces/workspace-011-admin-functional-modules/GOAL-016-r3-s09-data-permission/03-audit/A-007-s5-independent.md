---
id: A-007
goal: GOAL-016-r3-s09-data-permission
title: S5 关门独立审计（数据权限 · 行级/数据范围）
date: 2026-08-15
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S5 关门（成功标准 S1~S5 + D-002 §8 S2 清单 + A-004 F-001 闭合核验 + 实现/验证/安全/协议）
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-016-r3-s09-data-permission
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# A-007 · 独立关门审计意见（S5 · S-09 数据权限）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · S5 关门（00-meta S1~S5、D-002 §8、A-001~A-006 链、实现与验证证据、安全路径、协议/go、I-001~I-004）
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`root_goal: GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`plan_refs`/`primary_plan` = `VP-011-admin-functional-modules`；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001~D-004、`02-execution.md`、E-001~E-004、`03-audit.md`、A-001~A-006。
- **代码核对**：`handler/resources.go`（`ScopeConstraint` / `RowScopeProvider` / `resolveScope` / `scopeOwned` / list·create·detail·update·delete·batchDelete）、`handler/datapermission.go`、`handler/export.go`、`modules/datapermission/`（provider / store / migration / schema / manifest）、`kernel/profile.go`、`kernel/provider.go`、`composition/composition.go` + `composition_test.go`、`testsupport/store.go`、`store/migrate_test.go`、`errorcatalog.go`。
- **独立复跑（2026-08-15）**：`apps/api` `go test -p 1 -count=1` 覆盖 `handler` / `datapermission` / `composition` / `store` / `kernel` **全绿**；web 定向 `upstream-fixtures` + `app-manifest` + `schema-keys` + `s5-denominator` **82/82** 全绿。未复跑全量 `./...` 与全量 vitest 969（采信 E-004 记录 + 本轮定向复跑）。
- **covered**：成功标准对照、审计链、实现与方案一致性、验证充分性、安全路径、协议/go、信息项、波次级 e2e/冒烟是否可接受。
- **excluded**：不改 `status` / `progress` / goal-tree / `00-meta` / D-002 正文；不跑 e2e 双 profile 与容器冒烟（见 §8）。
- **保证等级**：L0。不得解读为第三方鉴证。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| A-004 F-001（list-only IDOR）已落地为行访问全覆盖 | `resources.go` list L415–422 下传 `filter.Scope`；detail L625–634 self 非本人 → 404；update L660–671 先 Get 再 404；delete L697–712 先 Get 再 404；batchDelete L796–814 仅保留本人 id；create L576–585 **强制覆盖** `body[OwnerColumn]=ActorID` |
| 404 无存在性 oracle | `resources.go` L632–634 / L669–671 / L710–712 与 `notFoundCode` 同码；`resources_test.go` L386–396 `ORDERS_NOT_FOUND` |
| batch-delete 仅删本人行 | `resources.go` L803–814；`resources_test.go` L434–445 `deleted:1` 且 `o-4` 仍在 |
| Create 强制 owner（A-005 F-002） | `resources.go` L583–585 覆盖而非缺省补；`resources_test.go` L422–432 客户端 `owner:user-b` → 入库 `user-a` |
| nil scoper 字节不变 | `resolveScope` L302–304；`TestResourceFactoryUnscopedNilScoper` L449–469 列出他人行且 detail 200 |
| ScopeFor 合成：enforceable ∩ policy ∩ assignment | `provider.go` L47–72：未接线 / 无策略 / disabled / 有效 `all` → nil；仅 `self` 出约束；赋值覆盖 default |
| 强制点在 PATCH（A-005 F-001 recommended） | `Service.UpsertPolicy` L83–86 / `UpsertAssignments` L96–99；composition L325–327 `NewService(..., nil)`；生产 PATCH 一律 `SCOPE_NOT_ENFORCEABLE` |
| default_scope 必填、无隐式 all | `datapermission.go` L87–90 空 → 400；store `CHECK (default_scope IN ('all','self'))`；测试 L144–148 |
| v1 无生产资源登记 | composition `enforceable=nil`；全库 `Scoper:` 仅测试 `resources_test.go`；export.go L127 仍硬编码 `users`/`roles` |
| 迁移 0027/0028 | `datapermission/migration/migration.go` Version 27 checksum `f3ce4c71…`；operationlog Version 28 checksum `a18c42e7…`；`compiled/persistence.go` 注册；`migrate_test.go` L592–593 |
| S-09 组合增量 +2 权限 / +1 导航 | `profile.go` L83–85 / L164–165；`provider.go` L417 `menu_data_permission`；`composition_test.go` L466–468 注释；`testsupport/store.go` L61–63 / L106 |
| 协议未扩 capability、pin 未动 | `capability-registry.json` 无 data-scope 键；D-002 §5 / D-004 本地鉴权扩展、v2.8.0 |
| I-001~I-004 均 verified，无到期 required | `00-meta` 信息表；最晚阶段均为 S1；本 scope 无 `deferred` |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结（D-002 + A-003/A-004/A-005） | 满足 | D-002 accepted；A-004 F-001 经 D-003 **fixed** + A-005 pass |
| S2 实现（D-002 §8 1–7） | 满足（1 项 recommended 残留） | E-003；§8.1–.5/.7 可核对；§6 白名单测试未落地（本意见 F-001） |
| S3 验证（go + web；e2e 波次级） | 满足 | E-004；本轮定向 go/web 复跑全绿；e2e/冒烟见 §8 |
| S4 go 判定 + 自审 | 满足 | D-004 内容扩展不失效；A-006 self pass |
| S5 本独立关门审 | 本条 | 无 required；无到期 required 信息项 |
| A-004 F-001 行访问全覆盖可核对 | 闭合 | 工厂六条路径 + `TestResourceFactoryEnforcesSelfScope` |
| A-005 recommended 随实施 | 主体闭合 | PATCH 强制点、Create 覆盖、全路径测试均已落地 |

## Findings

### F-001 · `owner_column` 白名单未实现（登记门仍关，无现网注入面）

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | D-002 §1 L23、§6 L60 要求「owner_column 白名单校验，防注入」。`store/repository.go` `UpsertPolicy` L117–141 将 `ownerColumn` 原样绑定写入，无标识符白名单；`datapermission.go` L97 透传 `body.OwnerColumn`。现网工厂用 `stringField(row, OwnerColumn)` 作 **JSON map 键**（`resources.go` L314），不拼 SQL；composition `enforceable=nil` 使生产 PATCH 无法落策略行。 |
| closure | — |
| 影响门禁 | 不阻断本目标关门；**首次生产资源登记前**必办（与导出面必办同级） |

v1 无 ScopeAware 生产实体、无列名插值，故无现网注入。后续目标把资源加入 enforceable 且实体把 `OwnerColumn` 拼进 `WHERE` 时，缺白名单即注入。建议：`UpsertPolicy` 只接受 `[A-Za-z_][A-Za-z0-9_]*`（或资源声明的列名集合），非法 → 400。

### F-002 · 省略 `default_scope` 的错误码是 `INVALID_SCOPE`，不是方案写的 `INVALID_PATCH_FIELD`

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | D-002 §3 L38 / D-003 L26：省略 → 400 `INVALID_PATCH_FIELD`。实现 `datapermission.go` L87–90 与 `datapermission_test.go` L144–148 为 `INVALID_SCOPE`。语义（必填、400、无隐式 all）已满足。 |
| closure | — |

契约码名漂移。不构成安全缺口；若管理 UI 未依赖该码，可在后续勘误 D-002 或改码二选一。

### F-003 · 台账「组合根 26/13」相对现网 27/13 过时（兄弟目标 S-10 增量）

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | D-002 §6 / E-003 / A-006 写 admin **26** 权限 / **13** 导航。现网 `composition_test.go` L471 `wantPermissions: 27, wantNavigation: 13`；L469–470 将 +1 权限记在 S-10 `users.mfa-reset`。S-09 自身增量仍是 +2 权限 / +1 导航。迁移快照现为 30（0029/0030 = MFA），0027/0028 行仍在。 |
| closure | — |

不是 S-09 功能缺陷。关门后若回写快照，按 live 27/13（+ 注明 S-10 +1）即可。

## 必改项汇总

无 required。无到期且影响本 scope 的 required 信息项。

## 波次级事项（可接受）

| 项 | 本目标 | 先例与说明 |
|----|--------|------------|
| 双 profile e2e | 未跑 | 00-meta S3 已写「归 S5 波次」。E-004 引 GOAL-009 E-003 略不准：该条 **跑了** e2e 8/8，**推迟的是冒烟**。更贴切的推迟先例是 GOAL-012 E-004（V-007/V-008 批末统一）。v1 无登记资源、无 e2e 新断言（`apps/web/e2e` 无 `data-permission`），零影响假设成立。 |
| 容器冒烟 V-007/V-008 | 未跑 | 与 R3 第二批相同：批末统一。`scripts/smoke.sh` admin `required_pages` 已含 `data-permission`。 |

**本目标关门可接受**这些项留到 R3 第三批收尾统一验证；不构成 required。批末必须补跑，失败则回流，不得用本 pass 代替波次证据。

## 与既有意见的异同

- A-006（self · pass）：同意 S2–S4 主体落地、A-004 F-001 已实施、可进 S5。本意见用代码行号与定向复跑独立核对后维持 pass。
- 不同意「§8 清单 1–7 全部落地」的字面：`owner_column` 白名单（§1/§6）未实现（F-001 recommended）。不升级为 required：强制点 `enforceable=nil` 已关掉登记面。
- A-004 F-001 required：本轮按实现复核为 **fixed**（与 A-005 / D-003 一致），不重开。
- A-005 F-001~F-003 recommended：PATCH 强制点、Create 覆盖、全路径测试均已可核对。
- 不与 A-001~A-003 / D-004 冲突。无意见冲突需 P-004。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass**。S-09 作为「框架 + 管理面、v1 不登记生产资源」的 data 门禁目标，成功标准、审计链、行级安全路径与验证证据充分；A-004 required 闭合可重复核对。无 high/med required；I-001~I-004 均 verified。

**可关门**（`status: done` 与 progress 5/5 由 `/govern` 执行；本意见不改状态）。

建议 `/govern`：

1. 响应本意见：记录 0 required；F-001~F-003 recommended 可带入后续「首次生产资源登记」目标或批末勘误，不阻断关门。
2. 将 S5 检查点勾选、progress 重算为 5/5、goal-tree 同步。
3. R3 第三批收尾统一跑 e2e 双 profile + V-007/V-008；失败回流，不把本 pass 当波次证据。

勿用 `progress: 3/5` 作为放行或拒绝依据。

## 声明

本意见不修改 `status` / `progress` / goal-tree / `00-meta` / D-002 正文。响应由 `/govern` 处理。保证等级 L0。
