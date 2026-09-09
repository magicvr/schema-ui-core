---
doc_type: vision-plan
id: VP-034-nav-group-collapsible
title: Admin 导航分组折叠体验
status: closed
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-034-nav-group-collapsible
created: 2026-09-07
updated: 2026-09-09
version: 0.4.0
parent: null
---

# VP-034 · Admin 导航分组折叠体验

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`closed`**（2026-09-09 · v0.4.0 · 用户书面确认关门） |
| lead_workspace | `workspace-034-nav-group-collapsible`（唯一 lead delivery） |
| Vision required | 计划阶段 [VRev-082](../reviews/VRev-082-vp034-nav-group-collapsible-planned.md) self `pass`；激活就绪 [VRev-083](../reviews/VRev-083-vp034-nav-group-collapsible-activation.md) self `pass`；scope 修正 [VRev-084](../reviews/VRev-084-vp034-existing-navigation-scope-correction.md) self `pass`；关门就绪 = [VRev-085](../reviews/VRev-085-vp034-nav-group-collapsible-close-out.md) self `pass`（七条判据 verified；open required = 0） |
| 组合位置 | **Admin 功能分支** · 导航分组折叠体验。不作为 VP-010 子目标；不改 Charter / VP-008 `go` |

## 意图

随着 Admin 左侧导航菜单项目增多，纯平铺列表已难以操作。本意图为 Admin Shell 左侧导航引入**可折叠/展开的分组（Group）**能力，并保持**分组与模块解耦**：每个模块自行声明自己的导航条目属于哪个具名分组，而不是按“一个模块一个组”硬编码。

**本次修正（用户书面确认，2026-09-07）**：本 VP 不仅验证新增 `group` 注册能力，还必须把当前已经注册的左侧导航纳入迁移范围，按产品语义合理分组；不能以“已有模块不迁移”作为退出条件或默认排除项。现有 top/user slot 导航不属于左侧菜单分组面，但仍须保留其原 slot 语义并纳入回归，不能静默丢失。

关键设计约束：

1. **分组与模块解耦**：分组不是按模块边界强制划分；模块可以把自己的不同导航条目注册到不同分组，也可以对明确的单例/非侧栏条目保持不分组。
2. **已有导航全量纳入**：当前已注册的 sidebar 导航节点必须完成合理分组迁移，或以有明确 UX 理由的顶层单例例外留痕；不能因其属于既有模块而排除在本 VP 外。
3. **向后兼容**：未声明 `group` 的未来模块导航条目继续以当前方式平铺展示；既有模块的权限、路由、slot 和 NodeID 不因分组迁移改变。
4. **激活态感知**：用户直接通过 URL 进入某个导航项时，其所属分组自动处于展开状态（不要求用户手动展开）。
5. **折叠/展开 UI**：分组标题可点击折叠或展开子项；折叠状态至少在本次会话内保持（具体存储方案由执行阶段冻结）。

**不进本 VP**：强制所有未来模块必须使用分组；top/user slot 导航全部搬到左侧；分组的服务端持久化；分组权限单独管控；多级嵌套分组；拖拽排序。

## 现有注册导航的初始分组基线

以下是基于当前代码主线已注册导航的**实现阶段初始分组方案**。它是本 VP 的覆盖基线，不是按模块切组；执行阶段可以基于真实 UI 复核调整组名或组内顺序，但不得无审计地减少覆盖范围。

| group key | 分组标题（初始） | 纳入的现有 sidebar 节点 | 所属模块 |
|-----------|------------------|--------------------------|----------|
| `identity-access` | Identity & access | `menu_users`、`menu_roles`、`menu_data_permission` | `admin.users`、`admin.roles`、`admin.data-permission` |
| `content-data` | Content & data | `menu_files`、`menu_dictionary` | `admin.file-library`、`admin.data-dictionary` |
| `operations` | Operations | `menu_activity`、`menu_monitoring`、`menu_scheduled_tasks`、`menu_recycle_bin` | `admin.activity`、`admin.system-monitoring`、`admin.scheduled-tasks`、`admin.recycle-bin` |
| `communications` | Communications | `menu_mail`、`menu_mail_outbox`、`menu_telegram` | `admin.settings`、`channel.telegram` |
| `commerce` | Commerce | `menu_wallet`、`menu_wallet_vouchers`、`menu_digitaloffer_offers`、`menu_digitaloffer_entitlements`、`menu_digitaloffer_purchases` | `admin.wallet`、`biz.digital-offer` |

补充边界：

- `menu_dashboard` 保持左侧**顶层主入口**，这是单例首页的有意例外，不是排除在 VP 外；必须验证其位置、激活态和直接 URL 行为。
- `dev.examples` 已有 `Examples` 组，保留为 demo profile 的现有分组，并纳入组折叠与回归验证。
- `menu_settings`、`menu_account`、`menu_wallet_self` 属于现有 `user` slot；`menu_notifications` 由通知铃/用户面消费。它们不搬到 sidebar，但必须验证 slot 不丢失、不被 sidebar 分组逻辑误处理。
- `admin.data-transfer`、`admin.login-captcha`、`admin.mfa` 当前没有独立导航贡献，不构成“已注册 sidebar 节点”遗漏。
- `channel.telegram` 与 `biz.digital-offer` 当前是已编译候选、非 `mvp`/`admin` 默认启用；它们的 sidebar 分组声明须在对应 custom/profile 启用时验证。

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在工作区目标内）：

1. **Shell 分组渲染**：Admin Shell 左侧导航支持 `group` 节点渲染，包含折叠/展开交互，键盘可访问（至少支持回车/空格切换）。
2. **模块注册协议兼容**：模块导航注册 API 新增可选 `group` 字段（字符串 key，并能表达组内顺序/组定义所需的最小元数据）；既有模块的权限、路由、slot、NodeID 不变。
3. **跨模块聚合**：来自不同模块的导航条目可以注册到同一个 group，聚合结果稳定、无重复/丢失；不得要求 Shell 为每个模块维护中央业务注册分支。
4. **现有 sidebar 全量分组迁移**：按“现有注册导航的初始分组基线”完成当前已注册 sidebar 节点的合理分组：`identity-access`、`content-data`、`operations`、`communications`、`commerce` 五组全部有证据；`menu_dashboard` 顶层单例例外与 `dev.examples` 既有组均有明确验证。不得以“已有模块”作为排除理由。
5. **激活态感知**：用户直接通过 URL 导航到某个已分组模块页面时，其所属分组自动展开（若当前处于折叠态）；顶层单例和 user slot 项目不受影响。
6. **回归与向后兼容**：未声明 `group` 的未来/测试模块导航条目继续以当前方式平铺展示；既有 top/user slot 导航继续可用；既有 `mvp`、`admin`、`demo` 及 custom/profile 组合的导航覆盖、权限和直接 URL 路由验证通过。
7. **playbook 更新**：`module-contribution-playbook.md` 更新导航注册规范，包含 `group` 字段、跨模块共组、key 命名空间、组内排序、未分组例外与迁移建议。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 |
|----|--------------|------|----------|----------|------------------|------|
| I-034-001 | 当前已注册导航的完整清单、所属 slot、默认 Profile / optional Profile 覆盖是什么？ | required | R1 分母冻结 / R4 回归 | R1 | 对照 `apps/api/kernel/profile.go`、各模块 `provider.go` 与 `manifest/fragment.json`；冻结 NodeID、PageID、slot 与 profile 矩阵 | **verified**（R1 静态分母 + R4 runtime matrix） |
| I-034-002 | 五个初始分组的标题、组内顺序，以及 `menu_dashboard` 顶层单例是否符合实际使用语义？ | required | R1 方案冻结 | R1 | 依据本 VP 初始基线进行 Shell 实现前 UI 复核并留存用户/审计决策 | **verified**（D-002；R1 顶层单例为历史基线；现行展示见 GOAL-003 `workspace` 组） |
| I-034-003 | 折叠状态采用会话内状态还是浏览器持久化？ | non-blocking | R2/R3 交互实现 | R2 | 实现阶段选择并测试；不引入服务端存储 | **verified**（D-005 sessionStorage + R3 测试） |
| I-034-004 | 每个已分组 NodeID 的直接 URL、动态路径和激活态展开行为是否完整覆盖？ | required | 判据 5 / R4 验收 | R3 | 建立 sidebar NodeID → route 的 e2e 矩阵，覆盖直接 URL 与组折叠状态 | **verified**（R3/R4 route matrix） |
| I-034-005 | optional compiled modules（`channel.telegram`、`biz.digital-offer`、`dev.examples`）在 custom/demo profile 中的跨模块分组聚合是否成立？ | required | 判据 3/4/6 / R4 回归 | R4 | 使用 custom/demo profile manifest harness 验证可选模块启用、共组、权限和无丢失 | **verified**（R4 custom/demo runtime matrix） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-034-nav-group-collapsible | GOAL-001-nav-group-collapsible | lead delivery | 2026-09-07 | `/govern` scaffold；Root `done` 5/5；结项后 GOAL-002～005 均 `done`；工作区 `done` |

## 关门记录

- 2026-09-09 · **`active → closed` v0.4.0**（用户书面确认；[VRev-085](../reviews/VRev-085-vp034-nav-group-collapsible-close-out.md) self `pass`；open required = 0）。
- 七条方向级退出判据 #1～#7 全部 verified；lead [workspace-034](../../workspaces/workspace-034-nav-group-collapsible/workspace.md) `done`；Root [GOAL-001](../../workspaces/workspace-034-nav-group-collapsible/GOAL-001-nav-group-collapsible/00-meta.md) `done` 5/5。
- Root [A-010](../../workspaces/workspace-034-nav-group-collapsible/GOAL-001-nav-group-collapsible/03-audit/A-010-r5-closeout-self.md) self `pass` + [A-011](../../workspaces/workspace-034-nav-group-collapsible/GOAL-001-nav-group-collapsible/03-audit/A-011-r5-closeout-independent.md) grok-build independent `pass`；[A-012](../../workspaces/workspace-034-nav-group-collapsible/GOAL-001-nav-group-collapsible/03-audit/A-012-a011-recommended-response.md) 将 A-011 F-001～F-003 recommended 全部 `fixed`；Goal open required/recommended = 0。
- I-034-001～005 全部 verified（A-002 F-007 计划层 collecting 投影随本轮同步闭合）。
- **有界 residual**：`menu_dashboard` 的 R1–R5 证据是顶层单例；现行展示由 GOAL-003 按用户指令注册到 `workspace` 组。不重开 R1 分母，不把该增量写成判据 4 未交付。

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-07 | 初创（v0.1.0）；用户确认结构选型（新 VP + 新工作区，不作为 VP-010 子目标）；退出判据 5 条（含激活态感知）；自 VP-033 编号递增 |
| 2026-09-07 | v0.2.0 · 激活（VRev-083 self `pass`；Admin freshness PASS `dd1edade`→`f2044cf3`；五域零变更；VP-008 `go` 不暂挂）；lead `workspace-034-nav-group-collapsible`；交 `/govern` 开区 |
| 2026-09-07 | v0.3.0 · 用户修正 scope：当前已注册 sidebar 导航必须纳入迁移与验证；新增 `identity-access`、`content-data`、`operations`、`communications`、`commerce` 初始分组基线；`menu_dashboard` 保留为有意的顶层单例；top/user slot 保留原语义，不再以“已有模块不迁移”作为非目标 |
| 2026-09-07 | `/govern` 已创建 `workspace-034-nav-group-collapsible` 与 Root `GOAL-001-nav-group-collapsible`；R1-R5 路线已就位，Root active 0/5；代码实现待从 R1 开始 |
| 2026-09-09 | 用户书面确认关门：VRev-085 self `pass`（七条判据 verified、open required = 0）；workspace-034 结项（Root done 5/5 + GOAL-002～005 done）；I-034-001～005 计划层同步为 verified；`active → closed` v0.4.0。residual = Dashboard 现行 `workspace` 组（GOAL-003），R1 顶层单例为历史基线 |
