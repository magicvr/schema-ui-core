---
doc_type: goal-audit
id: A-009-r4-closeout-independent
status: recorded
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-14
scope: R4 final close-out / VP-036 全部方向级退出判据、修正后 R1 §2.1 分母、跨模块聚合、权限·Profile·路由 fail-closed、页面级动作 handoff/确认、Palette UX·ARIA·focus·shortcut、导航分组联动、i18n/theme、浏览器与自动化证据、全量测试构建、非目标/红线
verdict: pass
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# A-009 · R4 关门独立交叉审计（2026-09-14）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型**：close-out
- **scope**：`[workspace-036-admin-command-palette]` Root `GOAL-001-admin-command-palette` 的 R4 最终关门交叉审计；核对 VP-036 全部方向级退出判据与当前已提交实现/证据。不把 `progress: 3/4`、A-008 self 或未复跑的 Postgres 记录当作本会话完成证明。
- **verdict**：pass

## 范围与区间

工作区：`[workspace-036-admin-command-palette]`，canonical `docs/workspaces/workspace-036-admin-command-palette/`，Root `GOAL-001-admin-command-palette`。`workspace.md` 的 `id` / `root_goal` / `canonical_scope` / `plan_refs` / `primary_plan` 与 VP-036 `active` v0.2.0、Charter `schema-ui-core-admin-foundation@0.4.0` 一致。`shared_materials_catalog: none`，本轮无共享资料引用。未读取其他工作区上下文。

审的是当前已提交树 `e2eb838d`（工作区干净）。本会话独立复跑：Web Vitest **104 files / 1366 tests passed**；`tsc -b` exit 0；Vite production build exit 0（仅既有 chunk size warning）；`go test ./...` exit 0；mvp/SQLite 与 admin/SQLite `e2e/command-palette.spec.ts` 各 **2 passed**。Postgres 双方言组合本会话未复跑，只核对其 spec 与 E-005 记录，不升为 required。

本意见不修改 `status` / 检查点 / `progress` / goal-tree / 决策 / 代码。用户确认关门仍是 VP-036 退出判据 7 的剩余编排步骤。

## 成果（有证据）

- **工作区绑定**：delivery workspace 唯一绑定 VP-036；Root `parent: null`；共享资料目录 `none`。
- **分母与契约**：独立用当前模块 Schema corpus + 与 R1 §2.1 对齐的四 Profile Manifest 重跑 provider，item id **精确等于**修正后的 §2.1 清单，无 missing/extra：mvp 5/6、admin 18/16、demo 13/6、custom 22/18。inner/参数页（`dictionary-entries`、`wallet-entries`、`task-runs`、`users-invites`、`telegram-operator`）不进页面分母；`admin-list-batch` 的 `requiresSelection` 不进动作分母。`SearchableItem`/provider v1：稳定 id、冲突双方排除 `DUPLICATE_ITEM_ID`、精确→前缀→包含、上限 12。
- **跨模块聚合**：内置 `manifest` provider 只读 `projectNavigation` 投影 + 认证 Schema；App 以 `[manifestSearchProvider, ...(searchableProviders ?? [])]` 注入 seam，生产入口未接线额外 provider。不增加 Shell 中央模块注册分支。pinned AppManifest envelope 仍仅为 `protocolVersion` / `requiredCapabilities` / `app` / `pages` / `navigation`；`APP_MANIFEST_SOURCE` 仍 pin `81aa1d8`。
- **权限 / Profile / 路由 fail-closed**：可见页来自 `projectNavigation`（`navigation.test.ts` 对 `menu_*` false/非布尔 fail-closed）；动作收录复核 `visibleWhen` / `disabledWhen` / `permissionIntent` / cascade；声明权限无法解析则拒绝；Schema 失败只隐藏该页动作并报非敏感 `PAGE_SCHEMA_UNAVAILABLE`。`invokeAction` 对 modal/navigate/custom/request 统一前置 gate；custom 白名单在 gate 之后；未绑定 `{…}` navigate 与畸形 row mapping 拒绝且不调用 host。缺失 context path 的 `!=` 在 Web/Go 均为 false。后端 401/403 仍是最终边界。
- **页面级动作 handoff / 确认**：Palette 选择动作写入 `pendingPaletteAction`，等 owner `document.meta.pageId` 匹配后 `invokeAction(trigger, null)`；中途离开 owner 页丢弃 pending。本会话 mvp/admin SQLite e2e 均打开 Users 检索并经同一 executor 弹出 New user modal。分母内唯一带 `confirm`/`confirmKey` 的页面级触发器是 settings `reset`；`invokeAction` 会走既有 `setPendingConfirm`，Palette 未另建确认通道。
- **Command Palette UX / ARIA / shortcut**：topbar 全断点入口；`Ctrl+K`/`Meta+K` 忽略 editable/composition；dialog + combobox/listbox；Arrow/Home/End（输入框不抢 caret）/Enter/Escape；backdrop 外点击；Tab trap；关闭恢复焦点、选择后不强制拉回触发器；loading 时 listbox 仍在并以 `aria-busy` 标记；option 为 `role=option`；`en-US`/`zh-CN` 均有完整 `commandPalette.*`；语义 token（`bg-card` / `text-foreground` / `bg-overlay`）。
- **导航分组联动**：页面选择走 History API `onNavigate`；`CollapsibleNavigationGroup` 在 `item.active` 时自动展开。`command-palette-app.test.tsx` 断言 Palette 打开 Users 后 Admin group `aria-expanded=true`。
- **红线**：`searchable.ts` / `CommandPalette.tsx` 不读实体行、不引入 RT-X01/RT-X02/Redis/MQ/跨进程索引；查询关闭即清空，无 Palette 持久化；无 Saved Views / recent / pinned / 批量结果中心 / 未保存保护 / Toast 全局重做 / 第二业务域。

## 对照成功标准

| VP-036 退出判据 / 信息项 | 状态 | 证据 |
|---|---|---|
| 1. 分母与契约 | 达成 | R1 matrix §1/§2.1；本会话独立 ID dump 与 16/18 完全一致；`searchable.test.ts` |
| 2. 跨模块聚合 | 达成 | `searchable.ts` aggregate + Manifest provider；注入 seam；未改 pinned envelope |
| 3. 权限与 Profile 安全 | 达成 | 四 Profile fixture 分母；`programmatic-action-gate.test.tsx`；`searchable.test.ts` denied action；Web/Go missing-context / `\q` |
| 4. Command Palette 体验 | 达成 | `CommandPalette.tsx`；`command-palette.test.tsx`；i18n keys；本会话 mvp/admin SQLite e2e |
| 5. 导航联动与回归 | 达成 | `command-palette-app.test.tsx` group expansion；e2e Users History API |
| 6. 基础设施与范围保持 | 达成 | 代码无实体搜索/协议加宽；D-002/E-005 红线与本会话检索一致 |
| 7. 证据与审计 | 达成（用户确认仍待编排） | A-001～A-008 台账；本条 independent；开放 required = 0。`progress: 3/4` 不是本条证据 |
| I-036-001 | verified | §2.1 ID 清单 + 本会话 dump |
| I-036-002 | verified | gate + handoff + 测试/e2e |
| I-036-003 | verified | D-002 口径已实现；双语/ARIA/shortcut 有代码与测试 |
| I-036-004 / I-036-006 | verified | 无实体搜索；无 freshness 回退 |
| I-036-005 | deferred · non-blocking | 未到期；不进首波 |

## Findings

### F-001 · 四 Profile 矩阵测试仍只断言计数，未钉死 §2.1 ID 清单（recommended · low · open）

- **描述**：`searchable-profile-matrix.test.ts` 对 mvp/admin/demo/custom 只断言 page/action **长度**（5/6、18/16、13/6、22/18）以及 inner page / `requiresSelection` 排除，**不断言** R1 §2.1 的稳定 ID 数组。A-007 将 A-006 F-001 记为已用「R1 修正 ID」处理，过述了测试钉死程度。本会话独立 dump 证实当前实现 ID 与 §2.1 **完全一致**，因此这不是分母错误，也不阻断关门；若未来 Schema 在保持条数不变时替换触发器，计数测试不会失败。
- **证据**：`apps/web/src/app/searchable-profile-matrix.test.ts`；`attachments/r1-searchable-item-matrix.md` §2.1；本会话 dump（admin 16 / custom 18 ID 全匹配）。
- **关联**：I-036-001；A-006 F-001。
- **建议闭合**：在矩阵测试中 `toEqual` 四 Profile 的 page/action ID 列表；或书面接受计数 oracle 为残余。
- **状态**：open

### F-002 · 分母内唯一页面级 confirm 动作缺少 Palette/programmatic 专用测试（recommended · low · open）

- **描述**：admin/custom 分母含 `action:settings:reset`，其 actionButton 带 `confirm`/`confirmKey`。`invokeAction` 会 `setPendingConfirm`，与页面按钮同一通道，Palette 未绕过确认。但 `programmatic-action-gate.test.tsx`、`command-palette-app.test.tsx` 与 e2e 都只覆盖 modal/navigate/custom 拒绝，没有 Palette → settings reset → ConfirmDialog 的 oracle。e2e 覆盖的是 Users「New user」modal。这不是权限绕过，只是确认分支的回归钉不足。
- **证据**：`apps/api/modules/settings/schema/settings.json` actionButton `reset`；`render.tsx` `invokeAction` confirm 分支；测试文件无 `confirm` 断言。
- **关联**：I-036-002；VP-036 退出判据 3。
- **建议闭合**：补一条 programmatic/Palette 测试，断言 reset 弹出既有 ConfirmDialog 且确认前不发请求；或用户接受残余。
- **状态**：open

## 必改项汇总

无 required finding。本 scope 开放 required = 0。无到期且影响关门的 required 信息项。I-036-005 仍为 deferred non-blocking。

## 与既有意见的异同

- A-008 self `pass`：独立同意 R4 实现与证据足以支持关门准备；不把 self 的 Postgres 四组合当作本会话复跑。本条独立复跑了全量 Web/Go/tsc/Vite 与 mvp/admin SQLite Playwright。
- A-002 F-001（17/19 vs 16/18）：独立 ID dump 仍为 admin 16 / custom 18，保持 **fixed**，不重开。
- A-001 F-001 / A-002 F-003 / A-006 对 programmatic gate 的复核：代码与 5 条 gate 测试仍成立，保持 **fixed**。
- A-006 F-001：产品规则与 16/18 ID 一致；测试仍为计数 oracle，本条降为 recommended F-001，不升 required。
- A-006 F-002：loading listbox / `role=option` / Home-End 不抢 caret / outside-click / Meta+K / 双语已在代码与测试中；剩余 Home/End/`aria-activedescendant` 覆盖不足，不单独新开 finding。
- A-006 F-003：`00-meta.md` 与矩阵 §4 证据句已刷新，视为 **fixed**。
- A-006 F-004：unbound navigate、malformed row mapping、Go `"\q"` 已有专用测试，视为 **fixed**。
- 与 A-008 无结论冲突，无 P-004 需用户在关门前裁决的 required 冲突。两条 recommended 不阻断用户确认关门。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass。** 当前已提交实现满足 VP-036 首波方向级退出判据：修正后的 R1 §2.1 分母可独立核对，跨模块聚合、fail-closed、handoff、Palette/ARIA、分组联动、红线与全量测试/构建均有仓库事实。开放 required = 0。`progress: 3/4` 与 Root `active` 不得被本意见改写或当作已关门。

建议 `/govern`：

1. 记录对本条的响应。两条 recommended 可在关门后补测，或由用户书面 `accepted-residual`；不因它们阻断用户确认。
2. 向用户展示 R4 证据与本独立 `pass`，确认后将 Root 标为 `done` 并同步 goal-tree；VP-036 关门记录仍走 `/vision`，本条不改 VP status。
3. 本会话未复跑 mvp/admin Postgres Playwright；若用户要求双方言再钉一次，可在响应前补跑，不是 required 门禁。

## 声明

本意见 `source: independent`，不修改 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree / 代码。响应、finding 闭合与是否将 Root 标为 `done` 由 `/govern` 与用户确认处理。
