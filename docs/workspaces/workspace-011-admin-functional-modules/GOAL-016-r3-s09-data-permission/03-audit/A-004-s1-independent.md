---
id: A-004
goal: GOAL-016-r3-s09-data-permission
title: S1 方案冻结独立审计（数据权限 · 行级/数据范围）
date: 2026-08-15
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S1 方案冻结（D-002 + I-001~I-004 闭合 + 过滤下推/协议/迁移数字）
audit_type: design-plan
verdict: conditional
status: recorded
parent: GOAL-016-r3-s09-data-permission
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# A-004 · 独立审计意见（S1 方案冻结 · S-09 数据权限）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：design-plan · S1 方案冻结（D-002 全文、I-001~I-004 证据、过滤下推点、fail-open 取舍、协议对照、迁移/组合根、无越界）
- **verdict**：**conditional**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`root_goal: GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`plan_refs`/`primary_plan` = `VP-011-admin-functional-modules`；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001、D-002、E-002、`03-audit.md`、A-001～A-003；I-011-001 §4 S-09、§7；protocol-inventory D-PERM / ADR-0004；capability-registry（无 data-scope 键）。
- **代码核对**：`handler/resources.go`（resourceFilter / ExtraQuery / list / detail / update / delete / batchDelete / requirePermission）、`handler/users.go`（usersEntity.List）、`authsession/users_repository.go`（usersWhere）、`kernel/profile.go`（ProfileAdmin）、`kernel/provider.go`（DefaultNavigationOrder）、`composition/composition_test.go`（admin 24/12）、`modules/operationlog/migration/migration.go`（max Version 26）。
- **covered**：方案可实施性、data 门禁（fail-open / 行级完整面）、信息项闭合、协议对照、迁移与组合根、无越界、未选方案与 S2 清单。
- **excluded**：S2 实现、S3～S5、不改 status/progress/goal-tree/方案正文。
- **保证等级**：L0（入口分离）。不得解读为第三方鉴证。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 过滤下推主路径与代码一致：工厂 `list()` 构造 `resourceFilter` 后走 `ResourceEntity.List` | D-002 §2 引 `resources.go:363-370`；实际 ExtraQuery 环为 L359–366、`Entity.List` 为 L367；filter 构造 L355。路径正确，行号略偏（见 F-004） |
| ExtraQuery 是客户端查询白名单先例（dictKey），不是范围注入面 | `resources.go` L49–53、L147–151；D-002 拟新增服务端 `Scope` 字段，未把 scope 放进 ExtraQuery（正确：客户端可省略 Extra） |
| `usersEntity.List` 是 ResourceEntity 边界；`usersWhere` 是 repository where 组装先例 | `users.go` L108–111；`users_repository.go` L18、L365–370 |
| `RowScopeProvider` nil = 未启用，与 `CaptchaVerifier` 可选注入同构 | D-002 §2；`handler/auth.go` L16–23、L35、L93（nil 跳过） |
| I-001/I-004 为设计裁定（all/self；org 本波不纳入），未把假设写成已验证实现 | D-002 §1、§7；`00-meta` I-001/I-004 → D-002 |
| I-002 闭合指向服务端查询面，优于信息项原文「renderer/schema-render」 | D-002 §2；行级过滤必须在 handler/store，不能在 renderer |
| I-003 Profile 为 admin 默认集内容扩展，先例可核对 | D-002 §3；`profile.go` L47–82（file-library / data-dictionary 等注释段）；mvp/demo 不含 |
| 协议对照未声称协议覆盖：D-PERM = UI 显隐/鉴权意图；ADR-0004 = 表格/行操作 UI 概念 | protocol-inventory L113（ADR-0004「表格/行操作」）、L171（D-PERM）；capability-registry 无 data-scope 键；D-002 §5 处置 = 本地鉴权扩展、不改 pin v2.8.0 |
| 迁移预留 0027/0028 相对当前 max 26 自洽 | `operationlog/migration/migration.go` L225–226 `Version: 26`；无 27+ |
| 导航 12→13 与 `DefaultNavigationOrder` 现长 12 一致 | `provider.go` L403–416（12 项） |
| 无越界声明完整：不改 Profile 默认集语义 / 模块矩阵 / Manifest 装配 / 协议 pin；go 内容扩展不触发失效 | D-001 L20；D-002 §5 |
| 未选方案与 S2 清单覆盖 middleware / 角色级继承 / org / fail-closed / 生产资源登记 | D-002 §7–§8 |
| v1 不登记生产资源已文档化 | D-002 §2、§7；能力面交付 |

## 对照成功标准（S1）

| 标准 | 状态 | 证据 |
|------|------|------|
| 数据范围模型（作用域 / 继承 / 与 RBAC 合成）可实施 | 部分 | all/self + 用户级赋值 + 未赋值走 default_scope 已写清；org 按 I-004 排除。Get/Update/Delete 未纳入范围强制（F-001） |
| 过滤下推路径有代码锚点 | 部分 | list 路径可核对；「零代码接入」与实体 List 不消费 Scope 矛盾（F-002） |
| 权限键 / 端点 / Profile 归属 | 满足 | data-permission.read/write；admin 默认集；mvp/demo 不含 |
| 协议独立对照（§7 口径）且不声称协议覆盖 | 满足 | D-002 §5 |
| I-001/I-002 required 有决策/代码证据，未伪装实现已落地 | 满足（设计层） | 闭合对象是方案，不是运行过滤 |
| 默认 all + 显式登记（fail-open）取舍已留痕 | 满足（可接受，见 F-003） | D-002 §1、§7；v1 无生产登记 |

## Findings

### F-001 · 范围强制只覆盖 list，Get/Update/Delete/batchDelete 未设计（登记后 IDOR）

| 字段 | 值 |
|------|-----|
| level | **required**（med） |
| status | open |
| evidence | D-002 §2 仅注入 `resourceHandler.list()`。代码：`list` L310–376 走 filter；`detail` L548–559 仅 `Entity.Get(id)`；`update` L562–590、`delete` L593–624、`batchDelete` L636– 均按 id 操作、无 owner/scope 检查。`list()` 且丢弃 actor（L312 `if _, ok := requirePermission`）。导出面 `GET /api/export/{resource}`（`profile.go` L162）不经该 list。 |
| closure | — |
| 影响门禁 | S1 方案冻结 / S2 实施（data） |

v1 不登记生产资源，故**当前**无生产 IDOR。但本波交付的是「行级/数据范围」框架：测试资源或后续 S-13/S-14 一旦登记，list 被过滤、按 id 读写仍放行，是经典 IDOR。fail-open 默认（§7）已文档化；**list-only 未写入未选方案或残余**。

S1 必须二选一并留痕后再放行 S2：

1. 把 Scope 扩到 Get/Update/Delete/batchDelete（工厂在 Get 后核对 owner，或 store WHERE）；或
2. 书面声明「v1 仅 list 过滤；按 id 读写与导出不在本波」，范围/复审触发走用户 `accepted-residual`（P-004.3）。

### F-002 · 「零代码接入」与「SQL 在各 repository」矛盾

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | D-002 §2「零代码接入」。`usersEntity.List`（`users.go` L108–111）只转发 `Q/Sort/Order/Page/PageSize`，忽略任何新 `Scope` 字段；`usersWhere`（L365）只接受 search `query`。file-library `List`（`filelibrary.go` L46–76）同样不读 Extra/Scope。工厂无法改写各库 SQL。 |
| closure | — |

登记策略行**不会**自动过滤。S2 须把主张改为：登记 + 该资源 `List`/where（或可选接口）必须消费 `filter.Scope`；测试资源示范接线。不阻断在修正主张后开工。

### F-003 · 默认 all（fail-open）在 v1 可接受；PATCH 须强制显式 default_scope

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | D-002 §1「未赋值 → 资源登记的 default_scope（登记时显式声明，默认 all）」；§7 否决 fail-closed（全量迁移风险）。v1 不登记生产资源 → 与今日行为一致。 |
| closure | — |

作为**已文档化的设计取舍**（非未声明残余），独立审接受 v1 fail-open。建议：PATCH 省略 `default_scope` 时 400，禁止列默认值静默 all；登记视为安全敏感（已有 `data-permission.policy-update`）。后续首次登记生产资源时复审 fail-closed。

### F-004 · 组合根权限基数过时；list 行号与 actor 丢弃

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | D-002 §6「admin 权限 22→24」。实测 `composition_test.go` L465：admin **24** 权限 / **12** 导航。016 增量应为 **24→26**（+data-permission.read/write）、导航 12→13。S 系列先例即 S2 按 live snapshot 改断言（GOAL-012 D-002 也曾写 22→24，落地后 captcha D-003 维持 24/12）。`resources.go` ExtraQuery 为 L359–366 非 L363–370；L312 丢弃 user，self 约束必须保留 actor。 |
| closure | — |

迁移 0027/0028 vs max 26 **自洽**。权限数字不阻断 S2，实施时按 live 值改测试。

## 必改项汇总

1. **F-001（required · med）**：补 Get/Update/Delete（及导出是否本波）的范围强制，或用户书面 residual。未闭合前**不可无条件放行 S2**。

无 high required。I-001/I-002 设计层证据充分，未把未实现过滤写成已验证运行事实。

## 与既有意见的异同

- A-003（self · pass）认为方案可进独立审与 S2，无 findings。本意见同意：注入点、协议本地扩展、I-004、迁移预留、nil 门、无越界成立。
- 不同意「可无条件进 S2」：self 未覆盖 list-only IDOR（F-001）。A-002 立项 F-003（org/B-10）已由 I-004 + D-002 §1 闭合，本 scope 不再开放。
- 不与 A-001/A-002 verdict 冲突（立项 vs 方案是不同 scope）。

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional**。S1 方案主体可实施，协议对照诚实，信息项设计层闭合成立；**开放 F-001 required**，不可无条件放行 S2。

建议 `/govern`：

1. 展示 F-001；建议采纳路径 1（工厂对 Get/Update/Delete 复用同一 ScopeConstraint），或请用户书面 residual（路径 2）。
2. F-002～F-004 随 D-002 勘误或 S2 清单补记，不单独阻断。
3. 响应本意见后，再开 S2。勿把 `progress: 1/5` 当作放行依据。

## 声明

本意见不修改 `status` / `progress` / goal-tree / 方案正文。响应由 `/govern` 处理。保证等级 L0。
