---
id: A-002-r2-saved-view-independent
doc: audit-opinion
status: recorded
source: independent
verdict: pass
scope: R2 Saved Views C1-C3 implementation, fail-closed contract, and regression evidence
audit_type: execution-facts
goal_id: GOAL-003-r2-saved-views
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-003-r2-saved-views
version: 0.1.0
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-002 · R2 Saved Views 独立审计

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：execution-facts；GOAL-003 R2 C1～C3 实现与回归（不含关门、不含 R3/R4/R5）
- **verdict**：pass

## 范围与区间

工作区：`workspace-037-admin-workflow-continuity`（`root_goal: GOAL-001-admin-workflow-continuity`，`canonical_scope` 匹配，`shared_materials_catalog: none`，`primary_plan: VP-037-admin-workflow-continuity`）。未读取其他工作区。

本意见只核对 GOAL-003 的 C1～C3、R2-I-001～003、R1 F-002～F-004 入口响应、D-003 存储合同，以及当前实现/测试是否把状态越过 allowlist。不关闭 finding，不改 status/progress，不把测试结果当作唯一证据。C4（本独立意见的编排响应 + Git checkpoint）、R3 dirty-state、R4 跨页面统一反馈、R5 组合验收均尚未完成。

## 成果（有证据）

### C1 · 24 个 `type: table` 分母与 R1 F-002～F-004

独立对照 `apps/api/modules/*/schema/*.json`（35 个 Schema 文件）与 `GOAL-002.../attachments/r1-denominator-matrix.json` v0.2.0：

- 活 Schema `type: table` 节点 = 24；矩阵 `listSurfaces` = 24；`derivedCounts.listSurfaces` = 24。集合一致，无多无少。
- 24 张表的 columns / sortableFields / tableFilterFields / rowKey 与活 Schema 一致。
- 4 张表的矩阵 `dataSource` 字符串对应活 Schema 的 `data.url`（`dictionary-entries`、`my-wallet`、`wallet-entries`、`wallet-vouchers`），不是额外 table 表面，也不进入 Saved View 序列化。
- F-002：custom `hiddenRoutePageIds` 含 `telegram-operator`；`registered \ discoverable` = 5 且与 hidden 集合相等；`derivedCounts.profiles.custom.hiddenRoutePages` = 5；admin hidden 仍为 4。
- F-003：`data-permission/policies` 的 `tableFilterFields` 为 `[]`；活 Schema 无 `props.filters`，`resource` 是 `rowKey`。
- F-004：`notifications`（search form + `notification-center`）、`mail`（`mail-admin-tab`）、`telegram-operator`（`telegram-admin-tab`）均无 `type: table`。Saved Views UI 只挂在 `SchemaTable`（`type: table` → App `tableRenderer`），自定义列表不会进入首波。

D-001 与 R2 验收矩阵把分母冻结为这 24 张表。运行时 allowlist 从活 columns / sortable / table filters / search-form `targetTable` 字段派生，不把历史矩阵当 runtime 权威。

### D-003 存储合同

`saved-views.ts` / `schema-table.tsx` / `render.tsx`：

- 键：`schema-ui.saved-views.v1` + `encodeURIComponent(userId)` + `pageId` + `tableId`。`user.id` 为空则 `userId === null`，Saved Views 禁用，不写入共享空键。
- allowlist：只持久化 `q` / `filters` / `sort` / `order` / `pageSize` / `visibleColumns`；query、record、document 遇未知字段 fail closed。`page`、selection、record view、modal draft 不入库。恢复 `setQuery({ ...saved, page: 1 })`；选择视图时 `clearSelection`。
- `activeViewId` 只作 UI 指针：不进入 `ResourceQuery` / `buildResourceQuery`。缺失或不在剩余 views 中则丢弃指针。
- 同浏览器同设备：只经 `window.localStorage`；无 API、无库表、无迁移。跨设备/协作仍属 R2-I-004 / I-037-005 deferred。

当前用户 id 形态为 `user-` + hex 或 `user-admin` 一类，不含 `.`；pageId/tableId 为连字符 slug。现有字母表下键隔离成立。

### Fail closed（不污染当前 query）

| 输入 | 行为 | 证据 |
|------|------|------|
| malformed JSON / 非支持文档形状 | 整文档 `ok: false`，不 apply | `saved-views.test.ts`；`readSavedViews` |
| 未知 document 字段 | 整文档失败 | 同上 |
| 未知 query 字段（含 `page`） | `SAVED_VIEW_STATE_INVALID` | `normalizeQuery` |
| duplicate / Schema 外列或 filter | 丢弃该条，`droppedCount++`，不 apply | storage test；UI 未知 filter 测试 |
| active 指向已丢弃/缺失记录 | 指针清空 | storage test |
| storage 读抛错 / 不可用 | `SAVED_VIEW_STORAGE_UNAVAILABLE` | storage test；UI `status: error` + toast |
| storage 写抛错 / 超容量 | `WRITE_FAILED` 或 `LIMIT_REACHED`，不伪装成功 | `writeSavedViews` + storage test |
| 权限/Schema 边界 | 无独立 saved-view permission target；越权等价于活 Schema allowlist 失效并丢弃 | 矩阵 `permissionBaseline`；未知 filter UI 测试 |

选择非法视图时先校验，失败则 toast 并从不合法集合中剔除，不把未知 filter 写入 fetch URL。Toast 走 `messageKey: feedback.savedViewError`，不回显原始 payload。

### UI（C3）

`schema-table.tsx` 有保存、名称确认、选择、更新、删除确认、列可见性 checkbox、空列表（仅「当前筛选」+ 保存）、错误条与 toast。`saved-views.ui.test.tsx` 覆盖：保存 allowlist 状态、重挂载恢复第一页筛选、选择恢复、更新保留 id、删除清空 views/activeViewId、无效记录 toast 且 query 不含未知字段。i18n 中英文键齐全。列可见性有实现（`displayColumns` 过滤 + 随 view 保存/恢复），但 UI 测试未切换 checkbox。

### 测试与 TypeScript（辅助，非唯一证据）

在 `apps/web` 使用项目 vitest 3.2.7 / 本地 `tsc`（不是仓库根 npx 的 vitest 5）：

- `npm test -- src/renderer`：28 files，**374 passed**。其中 Saved View 定向 7+3 通过。
- 与表格/查询/Saved View 相关的源码 `it()` 子集合计 **116**：`schema-table.test.tsx` 35、`saved-views.test.ts` 7、`saved-views.ui.test.tsx` 3、`schema-crud.test.tsx` 22、`resource.test.ts` 33、`search-form-filters.test.tsx` 1、`schema-dictionary-entries.test.tsx` 5、`denominator-render.test.tsx` 2、`representative-pages.test.tsx` 8。该子集包含在上述 374 中。
- `tsc -p tsconfig.app.json --noEmit`：exit 0。

denominator-render 使用 `context={{}}`，无 `user.id`，Saved Views 在该烟测中禁用，只证明 35 页仍能渲染，不能单独证明 Saved View 行为。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1：24 个 `type: table`；F-002～F-004 入口修订 | 达成 | 活 Schema 对照 + 矩阵 v0.2.0 + D-001 + E-002 |
| C2：编码键、allowlist、失效丢弃、读写失败单元证据 | 达成 | `saved-views.ts` + 7 项 storage tests |
| C3：保存/选择/恢复/更新/删除/列可见性/空态/错误反馈 | 达成（列可见性以代码为主，见 F-001） | `schema-table.tsx` + 3 项 UI tests + i18n |
| C4：self + independent + finding 响应 + Git checkpoint | 未开始/未完成 | 本意见刚落盘；尚未 /govern 响应或 checkpoint |
| R2-I-001 | 信息已可核对；`00-meta` 登记仍 `collecting`（F-002） | 矩阵 + 活 Schema |
| R2-I-002 | 信息已可核对 | 活 Schema allowlist + 失效测试 |
| R2-I-003 | 信息已可核对 | storage 读写失败 + UI 无效记录反馈 |
| R2-I-004 | deferred non-blocking | D-003 / I-037-005 |
| R3 / R4 / R5 | 未完成 | workspace 纲领；本目标不覆盖 |

## Findings

### F-001 · UI 回归未覆盖列可见性切换、写失败 UI 与跨用户键隔离

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：C3 要求列可见性、空态与错误反馈。列 checkbox、空下拉和读失败 toast 在代码中存在；3 项 UI 测试覆盖保存/选择/重挂载/更新/删除/无效记录，但未切换列可见性、未断言写失败/容量不足 UI、未跑两个 `user.id` 的 localStorage 隔离。空态体现为仅「当前筛选」选项，无额外空文案。
- 证据：`saved-views.ui.test.tsx`；`schema-table.tsx` 列 checkbox / `persistSavedViews`；`saved-views.test.ts` 无双用户用例。
- 影响：不推翻已实现的 fail-closed 合同；C3 的列可见性主要靠代码而非 UI 测试。
- 建议：补列切换保存/恢复、写失败 toast、以及 `user-a` vs `user-b` 键互不读取的测试。不要把 374/116 通过写成这些路径已测。

### F-002 · `00-meta` 信息登记与检查点不一致

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`00-meta.md` 将 C1～C3 标为完成、`progress: 3/4`，但 R2-I-001～003 仍为 `collecting`、证据栏仍写「待 D-001/测试」。`01-decision.md` 与本文件的信息就绪表已写 `verified`。独立对照认为这三项信息已有证据，登记字段过期，不是信息缺失。
- 证据：`GOAL-003.../00-meta.md` 信息表 vs 检查点；`01-decision.md`；`03-audit.md` 信息就绪表。
- 影响：读者若只看 `00-meta` 会以为 C1～C3 门禁信息未关。不构成实现合同失败。
- 建议：`/govern` 响应时把 `00-meta` 的 I-001～003 与证据栏改到与事实一致；本意见不改 meta。

### F-003 · goal-tree / workspace / VP 投影仍为 `0/4`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：GOAL-003 `00-meta` 为 `active · 3/4`，`goal-tree.md`、`workspace.md`、VP-037 R2 行仍写 `0/4` 且「C1～C4 待完成」。C4 确实未完成，因此不是宣称 R2 已关门；但是 AGENTS 要求改 progress 时同步 goal-tree。
- 证据：`goal-tree.md` 树与表；`workspace.md` 纲领 R2 行；`docs/vision/plans/VP-037-admin-workflow-continuity.md` R2 行。
- 影响：投影滞后。C4 关闭前不得把 Root R2 标完成。
- 建议：等本意见经 `/govern` 响应且 checkpoint 后，再同步树/区/VP；本意见不改这些文件。

### F-004 · 存储键用 `.` 拼接且未编码段内 `.`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`savedViewStorageKey` 对段做 `encodeURIComponent` 后以 `.` 连接。`encodeURIComponent` 不编码 `.`。若未来 `user.id` 或 pageId 含 `.`，可能与另一组分段碰撞。当前 id（`user-*` hex / `user-admin`）与 page/table slug 不含 `.`，未见跨用户泄露。
- 证据：`saved-views.ts` `savedViewStorageKey`；`saved-views.test.ts` 仅覆盖 `/` 与空格；`randomHexID("user")`。
- 影响：现网字母表下不是泄露。字母表变化前应换不可逆碰撞的分隔或显式编码 `.`。
- 建议：保持观察；若用户 id 规则变化，先改键编码再扩数据。不要当作已发生的跨用户泄露。

无其他 finding。未把 F-001～F-004 接受为 residual，也未 overrule。

## 必改项汇总

无 required / 必改 finding。

## 与既有意见的异同

| 项 | A-001 self | 本条 independent |
|----|------------|------------------|
| verdict | pass | pass |
| C1 24 表 + F-002～F-004 | 达成 | 独立复算活 Schema 后同意 |
| D-003 合同 | 达成 | 同意；补充键分隔 hardening（F-004） |
| fail closed | 达成 | 同意；测试非唯一证据 |
| C3 列可见性 | 宣称 UI 已接入 | 代码成立，UI 测试未切列（F-001） |
| 文档一致性 | 未列 | F-002、F-003 recommended |
| C4 / R3 / R4 / R5 | 未完成 | 同意，明确未完成 |

无冲突。self 无开放 required；本条也无。编排器无需 P-004 冲突裁决。

## 结论 + 建议给编排器/用户的下一步

R2 C1～C3 的实现与回归经独立核对成立：首波严格 24 个 Schema table；R1 F-002～F-004 已在矩阵/方案入口处理；localStorage 用户/页面/表格边界、allowlist、`activeViewId`、同设备语义与 fail-closed 与 D-003 一致；未见跨用户泄露或未知字段进入当前 query。本意见 **pass**，无必改项。

未完成且不得在本意见中宣称完成：

- C4：`/govern` 响应本条 + A-001、Git checkpoint、再考虑投影 Root R2
- R3 未保存变更保护
- R4 跨页面统一 Toast/恢复
- R5 组合验收与 VP 关门

建议下一步：用 `/govern` 响应 A-002（及 A-001），处理 F-001～F-004 recommended（可修测试/同步 `00-meta` 与 goal-tree，或书面留下不修理由），然后做 checkpoint。不要把 recommended 静默当成 residual/overrule。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
