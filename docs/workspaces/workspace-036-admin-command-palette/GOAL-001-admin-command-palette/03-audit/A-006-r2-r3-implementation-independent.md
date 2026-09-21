---
doc_type: goal-audit
id: A-006-r2-r3-implementation-independent
status: recorded
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-14
scope: R2/R3 实现交叉审计 / SearchableProvider v1、Manifest 聚合、Command Palette、App handoff、programmatic gate、missing-context 表达式、i18n/theme、测试与 mvp/admin SQLite smoke、R1 16/18 分母、fail-closed、无实体搜索/协议加宽、ARIA/focus/shortcut、可否进入 R4
verdict: pass
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# A-006 · R2/R3 实现独立交叉审计（2026-09-14）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型**：execution-facts
- **scope**：已提交的 R2 SearchableProvider v1 与 R3 Command Palette / 权限动作 handoff / ARIA·focus·shortcut / i18n·theme / 测试与 mvp/admin SQLite browser smoke；核对修正后的 R1 矩阵（admin 16 / custom 18）、Profile×permission×route fail-closed、无实体搜索与 pinned AppManifest 加宽，以及 R3 是否可进入 R4
- **verdict**：pass

## 范围与区间

工作区：`[workspace-036-admin-command-palette]`，canonical `docs/workspaces/workspace-036-admin-command-palette/`，Root `GOAL-001-admin-command-palette`。`shared_materials_catalog: none`，本轮无共享资料引用。未读取其他工作区上下文。

本意见只审已提交实现（`6ccec237` SearchableProvider / Palette 接线；`c5308b0b` 交互切片、表达式 fail-closed 与 e2e spec）是否满足 R3 进入 R4 的实施门禁。不把 self A-005 的 Playwright 结果当作本会话复跑证据；不勾选 R3 检查点、不改 status/progress。R4 全 Profile×权限×路由证据矩阵与关门不在本条放行。

## 成果（有证据）

- 工作区绑定成立：`workspace.md` 的 `id` / `root_goal` / `canonical_scope` / `plan_refs` / `primary_plan` 与 Root、VP-036 `active` v0.2.0、Charter `schema-ui-core-admin-foundation@0.4.0` 一致；共享资料目录 `none`。
- **SearchableProvider v1**（`apps/web/src/app/searchable.ts`）：稳定 `id`/`kind`/`label`/`keywords`/`group`/`href`/action context；provider 按 id 排序聚合；重复 id 双方排除并报非敏感 `DUPLICATE_ITEM_ID`；大小写/重音不敏感；精确 → 前缀 → 包含；上限 12。本会话 targeted Vitest `searchable.test.ts` 5/5 通过。
- **Manifest provider**：目标来自 `projectNavigation` 的可见 top/sidebar/user 叶子；认证上下文补 `notifications`；`page.route` 含 `{`、未挂导航、未解析 href、行级 `table.props.actions`、`requiresSelection` / `batchMapping`、`type: upload`、未知 custom handler、无 inline content 的 modal、含 `{` 或无法 `matchRoute` 的 navigate 均不进全局命令。Schema 读取走 host 注入的 `loadPageDocument`（D-VAL + 认证 transport），失败只隐藏该页动作并给出非敏感 `PAGE_SCHEMA_UNAVAILABLE`。
- **R1 修正分母独立复核**：按矩阵 §1.3 与当前模块 Schema/fragment 重数可见页直接 toolbar/actionButton。admin 仍为 16：users×4（create/invites/export/import）、roles×2、settings `reset`、file-library `openUpload`（**modal**，不是顶层 `upload`）、data-dictionary `openCreate`、scheduled-tasks `openCreate`、recycle-bin `purgeAll`、data-permission `openRegister`、my-wallet `openRedeem`、wallet-vouchers `openGenerate`、wallet `openCreate`/`runReconcile`。custom 另加 telegram-settings `telegram-operator-entry-button`（navigate `/telegram-settings/operator`，与已注册页路由一致）与 digitaloffer-offers `openCreate`，共 18。`dictionary-entries` 无公开导航叶且路由含 `{dictKey}`，其 toolbar 不进分母。activity `openDetail`、account `revokeSession` 为行级 `actions[]`，未收录。
- **Command Palette UI**（`CommandPalette.tsx`）：topbar 全断点入口；`Ctrl+K`/`Meta+K` 忽略 input/textarea/select/contentEditable 与 composition；dialog + combobox/listbox；Arrow/Home/End/Enter/Escape；backdrop 外点击关闭；Tab trap；关闭恢复焦点、选择后不把焦点强行拉回触发器；loading / partial-error / empty；结果 cap 12。
- **App handoff**（`App.tsx`）：页面项走既有 `onNavigate` History API；动作项写入 `pendingPaletteAction`，等 owner page `document.meta.pageId` 匹配后再 `invokeAction(trigger, null)`；用户中途离开 owner 页丢弃 pending；同页选择聚焦 `#page-title`。导航分组 `CollapsibleNavigationGroup` 在 `item.active` 时自动展开（VP-034 既有语义），Palette 未另建第二套路由。
- **Programmatic gate**（`render.tsx`）：modal/navigate/custom/request 统一复核 declared permission、L2 结构、rowless `tableActionGate`；`runRequest` 的 custom 分支在 gate 之后；actionButton 回退 `node.id` 为 `key`；无 mapping 的 `{…}` navigate 拒绝；row navigate 构造异常转反馈。`programmatic-action-gate.test.tsx` 3/3：denied modal 不打开、允许路径打开 modal、denied custom 不发请求。
- **Missing-context / malformed literal**：`evaluateExpression` 对缺失 path 一律 `false`（不再让 `undefined != x` 为真）；`isValidExpression` 拒绝 JSON 无法解析的引号字面量。Go `account.Evaluate` 对空/null path 同样 fail-closed。AppManifest 根字段仍仅为 `protocolVersion` / `requiredCapabilities` / `app` / `pages` / `navigation`，未知字段 `UNKNOWN_MANIFEST_FIELD`。本会话 `app-manifest.test.ts` 15/15、`go test ./internal/account` ok。
- **无实体搜索、无协议加宽**：provider 不读 list/DB 行；查询只在已聚合候选项上客户端过滤，关闭时清空，无 localStorage。`searchableProviders` 仅 App 可选注入 seam，生产入口未接线额外 provider。`APP_MANIFEST_SOURCE` pin 未改。
- **i18n/theme**：`en-US`/`zh-CN` 均有完整 `commandPalette.*`；Palette 使用 `bg-card` / `text-foreground` / `bg-overlay` 等语义 token。`command-palette.test.tsx` 覆盖 zh-CN chrome、Escape 焦点恢复、Tab trap、Arrow/Enter。
- **本会话复跑**：Web targeted Vitest **5 files / 30 tests passed**；Go `internal/account` ok。未复跑 Playwright（见下）。

## 对照成功标准 / 信息门禁

| 标准或信息项 | 状态 | 证据 |
|---|---|---|
| 工作区绑定与 VP-036 delivery | 达成 | `workspace.md`；VP-036 `lead_workspace` |
| I-036-001 精确分母 admin 16 / custom 18 | 达成（规则与 Schema 一致） | 修正矩阵 §2/§2.1；本条独立重数；`searchable.ts` 收录规则 |
| I-036-002 权限/URL/动作 fail-closed | 达成（实现已落地；00-meta 证据句仍旧） | `render.tsx` invokeAction；`programmatic-action-gate.test.tsx`；`command-palette-app.test.tsx` denied action |
| I-036-003 字段/排序/快捷键/ARIA | 大部分达成；加载期 ARIA 与若干交互无测试 | D-002；`CommandPalette.tsx`；`command-palette.test.tsx`；F-002 |
| I-036-004 / I-036-006 | verified，无实体搜索、无 RT-X01/X02 | `searchable.ts` 无 list fetch；协议 envelope 未加字段 |
| I-036-005 | deferred non-blocking，不进首波 | 无 recent/pinned/Saved Views 实现 |
| Profile×permission×route 浏览器全矩阵 | 未开始（R4） | E-004 / A-005 仅 mvp/admin SQLite Users+New user；本会话未复跑 e2e |
| 进入 R4 | **可以** | R3 实施与安全 handoff 有代码+单元证据；开放 required = 0 |

## Findings

### F-001 · R4 仍须用真实 Manifest 钉死 16/18 与分组联动（recommended · med · open）

- **描述**：R2/R3 测试用合成 Manifest，不断言 admin/custom 的矩阵 §2.1 ID 清单，也不断言 Palette 跳转后侧栏分组 `aria-expanded`。独立按 Schema/fragment 重数后，实现规则与 16/18 一致，因此**不阻断进入 R4**。但 R4 不得把 fixture 30 测当成 Profile×权限×路由回归完成。
- **证据**：`apps/web/src/app/searchable.test.ts`；`command-palette-app.test.tsx`（只断言 `/users` 与 `#page-title`）；`attachments/r1-searchable-item-matrix.md` §2.1；`App.tsx` `CollapsibleNavigationGroup` 的 `item.active → setOpen(true)`。
- **关联**：I-036-001（影响 R4 回归）；VP-036 退出判据 3 与 5。
- **建议闭合**：R4 用当前运行时 Manifest + 认证 Schema 生成 item id 并对照 §2.1；至少一条 Palette → 分组页断言分组展开。
- **状态**：open

### F-002 · ARIA/快捷键语义有可核对残余（recommended · med · open）

- **描述**：冻结口径的主路径已实现，但有三处实现/证据缺口：① loading 时 combobox 的 `aria-controls` 指向尚未挂载的 listbox id；② dialog 级 Home/End 始终 `preventDefault`，查询框内光标移动被列表选择抢走；③ listbox option 是 native `button`。另外：外部点击、Home/End、`aria-activedescendant`、`Meta+K`、Windows 下仍显示 `⌘K` 均无测试。这些不是权限绕过，但 R3「键盘/ARIA」证据不完整，应在 R4 浏览器矩阵关闭或书面 residual。
- **证据**：`CommandPalette.tsx` loading 分支 vs `PALETTE_RESULTS_ID`；`handleKeyDown` Home/End；option `button`；`command-palette.test.tsx` 未覆盖上述项；`App.tsx` 触发器 `<kbd>⌘K</kbd>`。
- **关联**：I-036-003（最晚 R1 口径，R3/R4 收集可用性证据）。
- **状态**：open

### F-003 · 信息项/矩阵仍写「gate 未实现」（recommended · low · open）

- **描述**：代码已统一 programmatic gate，但 `00-meta.md` I-036-002 证据列仍写「程序化 gate 修正仍待实现」，I-036-003 仍写「双语/主题与浏览器行为仍待验证」，矩阵 §4 仍把 invokeAction 缺口写成进入 R2/R3 的待办。这是投影漂移，不是运行时 fail-open。
- **证据**：`00-meta.md` 信息表；`attachments/r1-searchable-item-matrix.md` §4；对比 `render.tsx` 与 A-004/A-005 的 fixed 声明。
- **建议闭合**：`/govern` 响应时只更新证据句（本独立意见不改 meta）。
- **状态**：open

### F-004 · D-003 的 navigate 拒绝路径缺少专用测试（recommended · low · open）

- **描述**：未绑定 `{…}` navigate 拒绝、row navigate 构造异常转反馈、Go malformed quoted literal 的 fail-closed 在代码中存在；`programmatic-action-gate.test.tsx` 只覆盖 modal/custom；`permission_test.go` 的 invalid 表不含 `"\q"`。R4 应补 oracle，避免回归时只靠阅读实现。
- **证据**：`render.tsx` `INVALID_NAVIGATE_URL` / `ROW_NAVIGATION_FAILED`；`apps/api/internal/account/permission.go` `parseLiteral`；`permission_test.go` `TestEvaluateInvalid`。
- **状态**：open

## 必改项汇总

无 required finding。本 scope 开放 required = 0。

## 与既有意见的异同

- A-005 self `pass`：独立同意 R3 实施已达到「可进入 R4」的门槛；不把 self 的 Playwright 2/2 当作本会话复跑。
- A-004 F-001（R3 UI 未完成）：独立视为 **fixed**（有 CommandPalette、App 接线、targeted 测试）。
- A-001 F-001 / A-002 F-003（programmatic gate）：独立复核代码与 3 条 gate 测试后视为 **fixed**；执行路径不再绕过 modal/custom 前置权限。
- A-002 F-001（17/19 vs 16/18）：不重开。独立按 §1.3 重数仍为 admin 16 / custom 18；`file-library` 的全局项是 modal `openUpload`，不会被 `type: upload` 排除。
- 本条新增均为 **recommended**，指向 R4 证据与 a11y 残余，不与 self 的 R3 放行结论冲突。无 P-004 冲突需用户在进入 R4 前裁决。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass。** R2 provider/聚合与 R3 Palette/handoff/gate 的已提交实现可核对；修正后的 R1 16/18 分母与收录规则一致；权限/路由 fail-closed 成立；未见实体搜索或 pinned AppManifest 加宽。R3 实现**可以进入 R4**。

本会话未复跑 `apps/web/e2e/command-palette.spec.ts`（mvp/admin SQLite）。该 spec 只覆盖登录、Users 检索、New user modal、移动入口与 Escape，本身也不是 R4 全矩阵。

建议 `/govern`：

1. 记录对本条的响应（recommended 可进 R4 清单，不必阻塞勾选 R3）。
2. 勾选 Root R3 检查点并启动 R4：真实 Manifest 的 16/18 oracle、Palette→分组展开、mvp/admin 双方言回归，以及 F-002/F-004 的 a11y 与 navigate 拒绝路径。
3. 顺手刷新 I-036-002/003 与矩阵 §4 的过时证据句（F-003）。

## 声明

本意见 `source: independent`，不修改 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree / 代码。响应、finding 闭合与是否勾选 R3 由 `/govern` 处理。
