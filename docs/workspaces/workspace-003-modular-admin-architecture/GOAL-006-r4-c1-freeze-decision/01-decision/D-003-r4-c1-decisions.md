---
id: D-003-r4-c1-decisions
doc: decision-entry
goal: GOAL-006-r4-c1-freeze-decision
source: user
date: 2026-08-05
status: accepted
---

# D-003 · R4-C1 三项 P-004 裁决

## 用户裁决

- **Provider**：接受 framework-agnostic `Provider` + `Registrar` surface；Persistence
  由 compiled-global Provider catalog 收集，不由运行时启用 Profile 决定迁移是否
  编译或执行。
- **Records**：历史演示的无语义实体，`historical-only`；删除其当前产品实现，不
  恢复 CRUD/API/seed/RBAC/menu/manifest/专属前端面。
- **operationlog**：接受 Option A。业务写入成功后 best-effort append；append 失败
  记录服务日志但不翻转业务成功；R4 不自动 purge/archive/delete，Activity UI 开关
  不改变 writer。

## Provider 精确契约（整包接受）

用户整包接受冻结包作为本 D-003 的 Provider 精确契约正文：
`GOAL-005/attachments/r4-c1-freeze-package-draft.md`（`status: accepted`）。该包
含 Contribution 最小字段与 `ContributionIdentity.Key` 校验、Plan 元数据与运行时
双检、注册/发布/生命周期顺序、`CompiledPersistence()` 与 compiled-global
Persistence 规则、Authorization/seed/security owner matrix、operationlog Option
A 冻结边界与兼容性切换顺序。C2 实施不得在未记录的情况下改变身份、冲突键、安全
语义或顺序；`ConfigNamespaces` 不在 R4 新增独立 Registrar 方法。

- Provider surface：framework-agnostic `Provider` + Plan-owned `Registrar`；
  Persistence 由 compiled-global Provider catalog 收集，不由运行时启用 Profile
  决定迁移是否编译或执行。
- R4-I002 / C1-I001 以该冻结包为 `verified` 证据；A-004 `F-IND-006-FR-001`
  以此闭合为 `fixed`。

## Accepted residual

| 字段 | 接受内容 |
|------|----------|
| residual | operationlog append 失败可能产生审计缺口；长期 duration/archive 尚未定义 |
| scope | R4 Users/Roles/Auth/Settings 写入和既有历史 events |
| owner | `magicvr` |
| review trigger | 合规/运营 retention 要求、日志规模阈值、恢复演练发现缺口，或进入 R5 数据生命周期决策 |
| review date | `2026-08-05 08:32:22 +08:00` |
| closure route | `accepted-residual`，不把接受解释为 retention 已永久定义 |

## 约束

`0003`/`0006` 迁移 ledger、历史 `records.*` operation-log 合法值、通用
`record_id`/`recordView`/`recordSource` 能力继续保留。GOAL-007 承接 Records 运行面
核验。C1 仍需最终 self + Grok independent review，D-003 不直接放行 C2。
