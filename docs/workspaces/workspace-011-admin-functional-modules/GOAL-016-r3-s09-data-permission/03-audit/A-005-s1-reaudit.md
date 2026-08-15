---
id: A-005
goal: GOAL-016-r3-s09-data-permission
title: S1 方案 A-004 required 闭合复审
date: 2026-08-15
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-004 F-001 required 闭合验证（D-002 §2/§3/§6 + D-003；行访问全覆盖 / ScopeAware / 导出面 / default_scope / 组合根）
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-016-r3-s09-data-permission
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# A-005 · 独立复审意见（S1 · A-004 required 闭合）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：finding-closure · A-004 F-001（required med）闭合验证；对照 D-002 §2/§3/§6、D-003 与现网 handler
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`root_goal: GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`plan_refs`/`primary_plan` = `VP-011-admin-functional-modules`；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001、D-002、D-003、E-002、`03-audit.md`、A-001～A-004。
- **代码核对**：`handler/resources.go`（list L310–376 / create L501–546 / detail L549–559 / update L563–590 / delete L594–624 / batchDelete L636–）、`handler/export.go`（L1–5、L126–128、L147–160）、`composition/composition_test.go` L465（admin 24/12）、`kernel/provider.go` L403–416（导航 12 项）、`operationlog/migration/migration.go` L225–226（max Version 26）。
- **covered**：A-004 F-001 关闭证据是否真实、充分、可实施；ScopeAware 是否消除「零代码接入」矛盾；导出面必办；default_scope 必填；组合根 24→26；方案自洽（迁移号 / 未选方案 / S2 清单）。
- **excluded**：S2 实现、S3～S5；不改 `status` / `progress` / goal-tree / D-002 正文 / `00-meta`。
- **保证等级**：L0。不得解读为第三方鉴证。

## 成果（有证据）

| A-004 主张 / 本复审核对项 | 闭合证据 |
|---------------------------|----------|
| Get/Update/Delete 按 owner；self 不属本人 → 404（不泄露存在性） | D-002 §2 L29；现网 `detail` L549–559 仅 `Entity.Get(id)`、`update` L582、`delete` L610 无 owner 检查——属 S2 待实施，方案已写死语义 |
| BatchDelete 仅删本人行（跳过） | D-002 §2 L29 |
| Create 时 self 资源 `owner_column` 写入当前 actor | D-002 §2 L29；现网 `create` L527 已传 `user`，可实施 |
| 导出面：本波无登记资源无暴露面；登记时须评估 data-transfer 并施加同约束 | D-002 §2 L29；`export.go` L126–128 仅 `users`/`roles` 硬编码，工厂登记新资源不会自动进导出 |
| ScopeAware：已登记实体必须消费 `filter.Scope`；工厂登记未实现 → 拒绝 | D-002 §2 L30；「零代码接入」原文已删除 |
| `default_scope` 必填；PATCH 省略 → 400 `INVALID_PATCH_FIELD` | D-002 §3 L38；§4 L48 `default_scope TEXT NOT NULL CHECK in (all,self)` |
| 组合根 admin 权限 24→26、导航 12→13 | D-002 §6 L62；`composition_test.go` L465 `wantPermissions: 24, wantNavigation: 12`；`provider.go` L403–416 现长 12 |
| 迁移 0027/0028 相对 max 26 | D-002 §4；`operationlog` `Version: 26`；与 017 的 0029/0030 互斥 |
| D-003 声明路径 | D-003 L23–29：五项修正 + 闭合路径 **fixed** |

## 对照 A-004 required 闭合标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 行访问全覆盖写入方案（Get/Update/Delete 404、BatchDelete 过滤、Create owner） | **闭合** | D-002 §2 L29 |
| 导出面是否本波已裁定 | **闭合** | 本波不接线；登记规范 + 后续目标必办（D-002 §2 L29） |
| ScopeAware 消除「登记即过滤 / 零代码」矛盾 | **主体闭合** | 主张已改为两项动作；强制点仍有实施歧义（本意见 F-001 recommended） |
| `default_scope` 必填、无隐式列默认 | **闭合** | D-002 §3 L38、§4 L48 |
| 组合根 24→26 与 live snapshot 一致 | **闭合** | L465 = 24/12 |
| I-001/I-002 required 信息项未重新打开 | **满足** | 设计层仍成立；无到期未证 required 信息项 |

## Findings

### F-001 · ScopeAware 强制点写在「工厂登记」，与运行时 PATCH 登记面未对齐

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | D-002 §2 L30「工厂登记时校验（未实现 → 拒绝登记）」。现网 `registerResource`（`resources.go` L195–199）挂载 users/roles 等全部通用资源，实体今日均无 ScopeAware。字面实施会拒绝既有资源挂载；若只校验新测试资源，则 `PATCH /api/data-permission/policies/{resource}`（§3 L38）仍可给未消费 `filter.Scope` 的实体写入策略行，list（L367 `Entity.List(filter)`）不会过滤——「登记行 = 过滤」在管理面上仍不 fail-closed。 |
| closure | — |
| 影响门禁 | 不阻断 S2；S2 实施时写死强制点 |

建议 S2：对 **已写入 `data_scope_policies` 的资源**，附加约束前断言实体实现 ScopeAware，否则不附加且 PATCH 返回 400；**不要**把 ScopeAware 作为所有 `registerResource` 的前置。v1 不登记生产资源，故当前无生产泄露。

### F-002 · Create「默认写入」owner 未禁止客户端伪造

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | D-002 §2 L29「Create 时 self 作用域资源 owner_column 默认写入当前 actor」。`decodeResourceCreate`（L394–438）按 `CreateFields` 收录客户端字段；若 owner 列在白名单内且仅「缺省补 actor」，调用方可写他人 id。 |
| closure | — |

建议：self 资源 Create **强制覆盖**为 actor，忽略/拒绝客户端 owner。

### F-003 · S2 清单未回写 F-001 全覆盖动作

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | D-002 §8 第 3 条仅「RowScopeProvider 接口 + resourceFilter.Scope 扩展 + 端点」。Get/Update/Delete 404、BatchDelete 过滤、Create owner、ScopeAware 校验、`list()` 须保留 actor（现网 L312 丢弃 user）未列入清单。§2 正文已写清，清单过窄可能回退 list-only。 |
| closure | — |

另：BatchDelete「跳过非本人」与现网整批原子（L627–635 / L690+ `BatchDeleter`）需在 S2 写清（先按 owner 过滤再提交剩余 id）。`owner_column`（DB 列）与 JSON map 键的对应未写，工厂 Get 后核对时需约定。

## 必改项汇总

无 required。A-004 F-001 已合法闭合（fixed）。

## 与既有意见的异同

- A-004 independent conditional：开放 F-001 required（list-only IDOR + 零代码矛盾 + 组合根 22→24 过时）。本意见核对 D-002 修正后：**F-001 required 闭合**；组合根已改为 24→26，与 L465 一致。
- A-004 F-002～F-004（recommended）未再升级：主张已改；本意见 F-001/F-003 是其实施精度残留，不恢复 required。
- A-003 self pass 与本复审不冲突。
- D-003 声称全 fixed：就 A-004 **required** 范围成立；recommended 精度项不否定闭合。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass**。A-004 F-001（required med）关闭证据充分、可重复核对：行访问全覆盖、导出面本波裁定、`default_scope` 必填、组合根 24→26 均已写入 D-002。无 high/med required；无到期 required 信息项。

**可放行 S2 实施。** 本意见 F-001～F-003 为 recommended，随 S2 一并处理，不单独阻断。

建议 `/govern`：响应本意见（记录 A-004 F-001 → fixed，本 A-005 recommended 带入 S2 清单）后开 S2。勿用 `progress: 1/5` 作为放行依据。

## 声明

本意见不修改 `status` / `progress` / goal-tree / 方案正文。响应由 `/govern` 处理。保证等级 L0。
