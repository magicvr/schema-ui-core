---
id: A-002-r1-r2-independent
doc: audit-entry
parent_goal: GOAL-001-nav-group-collapsible
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
type: design-plan
scope: R1 分组基线与 R2 契约优先方案
date: 2026-09-07
verdict: conditional
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# A-002 · R1 分组基线与 R2 契约优先方案独立审计

## 范围与区间

本意见由项目指定的本地 grok build（grok-4.6 · reasoning high）只读复核以下范围：

- `workspace-034-nav-group-collapsible` 的 workspace、goal-tree、Root 五件套与 D-002/A-001；
- `apps/api/kernel` 导航贡献/校验/Profile、system-data；
- `apps/api/internal/manifest`、composition 传参、模块 provider/Manifest fragment；
- `apps/web` app-manifest、navigation projection、App Shell 与相关测试；
- I-034-001/I-034-002 的阶段门禁与 D-002 用户确认事实。

本意见未修改任何文件、status、progress、goal-tree 或方案正文；未把未实现内容写成事实。

## Verdict

**`conditional`**；开放 required = **4**（F-001～F-004）。

方向正确：用户确认的五组与 Dashboard/top-user 边界成立；契约优先优于 Manifest-only 或 composition 中央映射；现码尚未实现分组，A-001 未把未实现写成事实。不能无条件冻结 R2 实施契约。

## 成果（有证据）

1. D-002 确实记录了用户确认的五组顺序 `identity-access → content-data → operations → communications → commerce`、Dashboard 顶层单例、Examples 保留及 top/user slot 边界。
2. D-002 选择结构化 `NavigationContribution` 作为分组语义权威，由 Manifest 聚合跨模块归一化；现码仍为平铺贡献、fragment 槽位作者与 `SortNavigation(NodeID)`。
3. Web 协议已经接受标准 `NavGroup`（label/labelKey/icon/items/visibleWhen/permissions），但 App Shell 当前只静态渲染组；折叠、键盘交互、自动展开尚未实施。
4. 现有 top/user 边界可核对：Dashboard 在 sidebar；Settings/Account/My wallet 在 user；通知为铃面且无 user-slot 链接；Examples 组和 Overview top 仅由 demo fragment 提供，无 NavigationContribution。
5. Vision required = 0；Goal 历史仅 A-001 self `pass`；R2/R3/R4 尚未完成。

## 对照成功标准

| 项 | 结论 |
|---|---|
| R1 五组顺序与 Dashboard/top-user 边界 | 部分：产品选择有 D-002；R1 分母矩阵尚未由运行时/代码证据冻结 |
| 模块可声明可选 group | 方向成立，类型/校验/聚合未实施 |
| 跨模块共组与冲突 fail-closed | 方向成立，比较字段/失败点/测试预言不足 |
| 未声明 group 向后兼容 | 方向成立，混排规则不足 |
| top/user 不误改 | 原则明确，但贡献类型不含 slot，聚合匹配方法未冻结 |
| Manifest 兼容 NavGroup | 可行；禁止向现行协议额外输出 `key`/`id` |
| R3 直接 URL 自动展开 | 已注册 sidebar 自身 route 可验证；内页/动态路径尚未纳入方案矩阵 |

## Findings

### F-001 · I-034-001 在 R1 仍未关闭，分母矩阵未冻结

- 级别：`required`，严重度 med；状态：open；关联 I-034-001。
- `I-034-001` 最晚阶段为 R1，影响 R1 分母/R4 回归；`00-meta` 仍为 `collecting`。当前只有可由 `profile.go`、provider、fragment 重建的盘点，没有经核对的 NodeID × PageID × slot × Profile 冻结证据，运行时矩阵仍待核。
- 在本 finding 合法闭合前，不得将 R1 检查点写为完成，不得宣称 R2 实施分母已冻结。
- 证据：`00-meta.md` I-034-001、VP-034 §信息需求、Profile/provider/fragment 代码。

### F-002 · 组序合成算法未定义，直接包组会得到错误顺序

- 级别：`required`，严重度 high；状态：open。
- D-002 要求冻结五组顺序，但未定义 `Group.Order` 是否权威、聚合发生在 Sort 前后、未分组项如何与组混排。
- 若先按现行 `DefaultNavigationOrder` 排序再按成员包组，首成员位次会产生 `identity-access → commerce → operations → communications → content-data`，违反用户确认顺序。R2 必须明确使用独立显式 GroupOrder，不得把 NodeID 首成员全局位次当作组序。
- 证据：D-002、`apps/api/kernel/provider.go` 的 `DefaultNavigationOrder`、`manifest.go` 的 `SortNavigation`、`navigation_order_test.go`。

### F-003 · 同 key 元数据一致性与 fail-closed 位置未定义

- 级别：`required`，严重度 med；状态：open。
- D-002 要求同一 key 的跨模块声明使用一致元数据，但未冻结比较字段（Key/Order/Label/LabelKey/Icon）、literal/key 等价规则、空 icon 规则、first-writer 与全等规则、失败发生层级和错误码。
- 没有这些定义，R2 无法写出可重复的冲突测试，也无法证明 fail-closed。
- 证据：D-002 结构化贡献方案；现码无 group 冲突路径。

### F-004 · 未分组项、Dashboard、已有协议组的混排规则未定义

- 级别：`required`，严重度 med；状态：open。
- 未声明 group 的平铺项仍需兼容，但尚未定义 Dashboard 是否固定最先、未来未分组项插入组间还是组后、Examples 作者组如何保持、user/top 误标 Group 是拒绝还是忽略、settings/wallet 模块的 per-node 语义如何约束。
- 没有混排规则，R2 无法证明未声明 group 向后兼容和 top/user 不受影响。
- 证据：D-002、`dev/examples/provider.go`、settings/wallet/account/notifications fragments、composition 当前只传 fragments。

### F-005 · R3 直接 URL 自动展开对内页/动态路径尚未定义

- 级别：`recommended`，严重度 med；状态：open；关联 I-034-004，尚未到期。
- `navigation.ts` 的 `linkActive` 只匹配 NavLink 自身 page/url；`BREADCRUMB_PAGE_PARENTS` 未参与 nav active。需要在 R3 矩阵中明确是否要求 users-invites、wallet-entries、dictionary-entries、task-runs、telegram-operator 等内页点亮父组。

### F-006 · NavGroup 无稳定 id；不得用额外 key/id 破坏当前协议

- 级别：`recommended`，严重度 med；状态：open。
- 现行 schema 与 `parseItem/ensureKeys` 不允许 NavGroup 的额外 `key`/`id`。若为了折叠持久化加入字段，必须协议 bump；建议 R2/R3 先不改协议字段，使用 `labelKey` 或 child active 作为 UI 稳定标识。

### F-007 · I-034-002 台账状态不一致

- 级别：`recommended`，严重度 low；状态：open；关联 I-034-002。
- D-002/`01-decision` 已记录用户决策，但 `00-meta` 与 VP-034 的 collecting 表述没有统一。建议明确为“产品决策已验证；Manifest/Shell 行为仍待后续证据”，避免读取冲突。

### F-008 · Parent、Group 与 system-data checksum 的边界未落盘

- 级别：`recommended`，严重度 low；状态：open。
- `Parent` 是现有 NodeID 层级/排序字段，不应被复用为产品 Group；需写明 Group 不改变 `menu_items` 身份，并明确 group 元数据是否进入 navigation checksum/SystemDataVersion。分组若只影响 Manifest，可明确与 system-data 解耦。

### F-009 · 残组、空组、fragment 作者组与结构化组的共存规则未写

- 级别：`recommended`，严重度 low；状态：open。
- mvp/admin 下会出现只有部分成员的残组；未启用的 optional modules 不应造成空组；零成员组不能输出（协议拒绝空 items）；Examples 组应保留且不与 Admin 五组合并。

### F-010 · 组内顺序权威源与契约测试清单未写

- 级别：`recommended`，严重度 low；状态：open。
- D-002 同时写了“VP 表格顺序”和“现行 NodeID 排序”，需选唯一权威；R2 必须加入一致/冲突、跨模块共组、未声明平铺、top/user 不变、Examples 保留、Dashboard 单例、optional modules 聚合、组序不跟随 NodeID 首成员的测试预言。

## 必改项汇总

1. F-001：冻结并落盘 I-034-001 的 NodeID/PageID/slot/Profile 矩阵；运行时核验可作为 R4 回归，但 R1 分母本身必须形成可追踪状态。
2. F-002：写死显式 GroupOrder 的组序算法，确保输出是用户确认的五组顺序。
3. F-003：冻结同 key 元数据比较字段、等价规则、失败层与错误码。
4. F-004：冻结 Dashboard/未分组平铺/结构化组/Examples 的混排规则，以及只归一化 sidebar 普通链接的匹配和误标 Group 行为。

## 信息门禁

| ID | 级别 | 最晚阶段 | 本轮判定 |
|----|---|---|---|
| I-034-001 | required | R1 | 未关闭（F-001） |
| I-034-002 | required | R1 | 产品决策有 D-002 证据；台账不一致，行为证据留待后阶段（F-007） |
| I-034-003 | non-blocking | R2 | 不阻断 R1/R2，R3 前再决 |
| I-034-004 | required | R3 | 未到期；可验证性缺口见 F-005 |
| I-034-005 | required | R4 | 未到期 |

共享资料：`none`，未作为事实或关闭证据。

## 与 A-001 self 的异同

- 同意：用户确认的产品基线、契约优先方向、无中央业务映射；未把 R2/R3 写成已完成。
- 补充：现行 `SortNavigation` 首成员位次会导致错误组序；Examples 无结构化 NavigationContribution；R2 还缺组冲突、混排、slot 匹配的操作定义。
- 分歧：A-001 self 为 pass 且 0 finding；本 independent 为 conditional，4 required + 6 recommended。两者不是产品选择冲突，而是“是否已具备无条件实施冻结”的门禁判断差异；按 P-003 应响应 required，不得静默取 pass。

## 结论与建议

建议 `/govern`：

1. 响应 A-002；F-001～F-004 走 `fixed`（补 D-002/矩阵）前不要开 R2 编码。
2. 用当前代码盘点形成 I-034-001 初稿并补运行时核对策略。
3. 统一 I-034-002 台账为“决策已验证；Manifest/Shell 行为待后续证据”。
4. 将 F-005/F-006 作为 R3/R2 约束附言，不在本轮私自 bump 协议。
5. required 闭合后可再次 `/audit` 复审本 scope，再进入 R2。

## 声明

本意见为 `source: independent`，不修改目标 status/progress；响应由 `/govern` 处理。
