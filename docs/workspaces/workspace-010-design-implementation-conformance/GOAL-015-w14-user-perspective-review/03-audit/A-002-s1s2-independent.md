---
id: GOAL-015-w14-user-perspective-review
doc: audit-entry
record_id: A-002
source: independent
scope: S1 审视证据 + S2 台账完整性 + 关门边界（D-002 / I-001 deferred）
verdict: pass
status: recorded
auditor: grok-build (grok-4.6 · reasoning high)
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# A-002 · S1/S2 与关门边界独立交叉审计（2026-08-17）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：ad-hoc / close-out 前置 · S1 审视证据真实性 + S2 台账完整性 + D-002 关门边界诚实性（含 I-001 P-005 deferred）
- **verdict**：**pass**
- **工作区**：`workspace-010-design-implementation-conformance`（`root_goal` = `GOAL-001-design-implementation-conformance`；`canonical_scope` 已核对；`shared_materials_catalog: none`；`primary_plan` = `VP-010-design-implementation-conformance`）

## 范围与区间

- **covered**：GOAL-015 五件套（`00-meta` / `01-decision` + D-001/D-002 / `02-execution` + E-001 / `03-audit` + A-001）；`goal-tree.md` GOAL-015 行；`workspace.md` W14 行；D-001 F-01～F-14 指定抽验路径 + E-001 两条误判修正复核；git 工作区是否越界改业务代码 / Profile 默认集 / 模块矩阵 / Manifest。
- **不 covered**：F-01～F-14 的整改实施（本波明确非范围）；浏览器活栈点验；其他工作区上下文。
- **本意见不修改** `00-meta` / `01-decision` / D-* / `02-execution` / E-* / `status` / `progress` / `goal-tree.md` / `workspace.md`。

## 成果（有证据）

### S1 审视证据（抽验成立）

| 台账项 | 独立核对 |
|--------|----------|
| F-01 | `scheduled-tasks.json`：`createTask`/`updateTask` `bodyMapping` 含 `handler`；`openCreate`/`openEdit` `fields` 无 handler。`scheduler.go:45-49` 仅注册 `system.noop`。`scheduledtasks.go:90-92` 缺省 handler 回落 `system.noop`。UI 建任务会空转。成立。 |
| F-02 | `data-permission.json:32-43` 声明 `updateScopes`；`body` 仅 `openRegister` 工具栏，`policies` 表无 `actions`。成立。 |
| F-03 | `activity.json` `ops-search` 仅 `q`；`operations.go:13-25` Resource 仅 `QSearch`；`operationlog/repository.go` `ListOperationsFiltered` 只按 `operationsWhere(filter.Q)` 过滤。活动模块无 export。成立。 |
| F-04 | `notifications.go:273-285` 英文硬编码 title/body；文件头写明本地化模板为 R3/B-09 deferred。成立。 |
| F-05 | `recyclebin.go:53-54`、`wallet.go:57-58,276-277,303-304` 丢弃 `intParam` 失败；wallet 仓库 `pageSize < 1` 才回落、无上限；`datapermission.go:75-79` 用 `len(items)` 伪造 `pageSize`；`scheduledtasks.go:276-285` 硬编码 Page 1 / PageSize 50。成立。 |
| F-06 | `datapermission.go:98-110,132-135` 多种校验复用 `INVALID_SCOPE`；`wallet.go:80-82` 等复用 `INVALID_WALLET_BODY`；`errorcatalog.go` 各只一条泛化文案；`OPERATION_NOT_FOUND`（`operations.go:24` / `systemmonitoring.go:95`）不在错误目录。成立。 |
| F-07 | 通知：`notifications.go:118` 未 `ToLower(q)`，仓库 `notifications_repository.go:136-138` 只 `lower(title/body)`、查询串原样 → 大写查询 miss（台账「大小写敏感」概括可接受）。钱包搜索仅 `owner_id LIKE`（`wallet/store/repository.go:148-151`）。回收站 list 不读 sort/order。`wallet-entries.json` 有 `entryType` 列、无筛选字段。成立。 |
| F-08 | `App.tsx:669-672` 无条件渲染 `pageId` + `route` 技术信息框。成立。 |
| F-09 | `render.tsx:1169` 前缀 `feedback.code`；316/391/401/426/555/748/855/897 等为英文硬编码；`schema-table.tsx:548-568` 英文 schema 错误。成立。 |
| F-10 | `App.tsx:451-475` 以 `error.code` 作主标题，并展示 `error.url` / `issue.path`。成立。 |
| F-11 | `form-controls.tsx` 各 `Label` 只渲染 `{label}`，无 `*` / `required` / `aria-required`；必填只在提交时由 `form-controls.ts:485-501` + `render.tsx:1516-1528` 拦截。成立。 |
| F-12 | `confirm.tsx:24-49` 有 `role="dialog"`，无焦点圈 / ESC / 初始焦点。成立。 |
| F-13 | `data-table.tsx:329-356` 桌面 `<tr>` 仅 `onClick`，无 `tabIndex`/`role`/`onKeyDown`；移动卡片 `137-151` 已有。成立。 |
| F-14（主体） | 移动卡片 `fields.slice(1, 3)`（121-124）；select 空值走第一项（292-310）；排序只在 asc/desc 间切换（221-228）；Tabs 无方向键/`aria-controls`（1882-1910）；`<sm` `LocaleSwitcher` `hidden`（`App.tsx:902`）；禁用按钮无原因（`schema-table.tsx:714-731`）；`form-controls.tsx` 无 `aria-describedby`。成立。通知空态子项见 F-004。 |

### E-001 误判修正（诚实）

| 修正 | 独立核对 |
|------|----------|
| markAllRead 非死代码 | `apps/web/src/components/notification-center.tsx:120-156`：组件内按钮直接 `POST /api/notifications/read-all`。schema 未引用 `markAllRead` 仅为卫生问题。修正成立。 |
| 钱包 adjust/freeze/unfreeze/deductFrozen 非「无确认」 | `wallet.json`：上述动作为 `type: modal` 表单；工具栏 `Reconcile` → `runReconcile` 为直接 request。修正成立。 |

### 台账完整性

F-01～F-14 每条均有文件:行证据、用户影响与优先级分组（A/B/C/D）。条目语气是「改进候选」，不是「已完成整改」。E-001 / D-002 / `00-meta` 均写明无代码改动。git 工作区仅 `GOAL-015/` 新目录 + `goal-tree.md` / `workspace.md` 文档行，无 `apps/api` / `apps/web` 改动。

### 关门边界（D-002 / I-001）

- **用户原始请求**：`00-meta` 与 D-002 一致引述为「审视 api/web 已实现功能 → 找出并非很小的改进点 → 在工作区 10 新建子项目进行落盘」。本审计任务书亦书面确认「本波 = 审视 + 台账落盘 + 审计会签 + 台账同步；整改实施 deferred」。未见「本波必须实施整改」的书面指示。
- **I-001 P-005 deferred 字段**：延期理由（原始请求仅要求落盘）、责任人（用户）、下一复核触发（用户指示开始整改）均已落盘；原级别保留为 `required（对整改波次）`；影响门禁与最晚需要阶段均指向**未来整改波次**，不指向本波 S4。
- **不是借重组绕过门禁**：D-002 把本波交付收窄到与原始请求一致，F-01～F-14 完整留作 backlog，未静默实施、未静默丢弃、未立刻拆不可执行子目标（符合 P-001 / P-004 / P-005）。
- **本波闭门合理**：S1/S2 已有可核对产物；对本波关门无到期 required 信息项。S4 仍须响应本意见并做台账同步 / 关门自审（A-003），但**无 required finding 阻断放行**。

### 六处一致性（status / progress / deferred）

| 位置 | status | progress | 整改 |
|------|--------|----------|------|
| `00-meta` | active | 2/4 | deferred |
| `01-decision` | active | — | I-001 deferred；S3 推进中 |
| `02-execution` / E-001 | active | — | 只记 S1 审视事实 |
| `03-audit`（本意见写入前） | active | 结论曾写 1/4（过时，本索引已改） | A-001 写 I-001 归 S2 裁决 |
| `goal-tree.md` | active | 2/4 | deferred |
| `workspace.md` W14 | active | 2/4 | deferred；S3 推进中 |

权威状态列一致为 `active` / `2/4` / 整改 deferred。`00-meta` 检查点预勾与 D-001 §3 阶段号撞名见 recommended findings。

### 无越界

- 工作区绑定：`workspace_id` / `root_goal` / `canonical_scope` / `plan_refs`+`primary_plan` 合格；`shared_materials_catalog: none`，无资料引用被当作事实。
- 无跨区 `GOAL-*` 裸引用；无共享资料行。
- 未见 Profile 默认集 / 模块矩阵 / Manifest 装配语义改动（本波无业务代码）。go 判定「无影响、不暂挂」与 as-built 一致。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 审视有证据、非编造 | 达成 | E-001 + 上表抽验 |
| S2 台账 F-01～F-14 完整（证据/影响/优先级）且为候选非已改 | 达成 | D-001 §2；F-14 一子项过述见 F-004 |
| 本波不实施整改成文诚实、可追溯 | 达成 | D-002 + I-001 deferred 三字段 |
| I-001 非借重组绕过 / 非静默丢弃 | 达成 | 台账仍在；门禁移至未来整改波次并保留 required 级别 |
| 六处 status/progress/deferred 一致 | 达成（有 hygiene 缺口） | 权威列一致；检查点预勾见 F-001 |
| 无越界改代码 / 矩阵 / 跨区 | 达成 | git 短状态仅本区文档 |

## Findings

### F-001 · `00-meta` 路线图 S3/S4 预勾，与派生 progress 2/4 不一致

- 严重度：low
- 建议：**non-blocking**
- 状态：open
- 描述：P-001 要求 `progress` 由检查点确定性派生。`00-meta.md` 路线图 S1～S4 四个检查点均标 `[x]`，同时正文写「**2/4**（S1、S2 完成落盘），S3/S4 推进中」。`当前边界` 把 S3/S4 写在「已完成的事实性工作」标题下。权威 `status`/`progress`/`goal-tree`/`workspace` 仍为 active · 2/4，**不是**预写 `done` 或预登记本意见 pass（与 workspace-010 GOAL-013 A-002 F-001 不同量级）。S4 响应须把未完成检查点改回 `[ ]`，并把「已完成的事实性工作」与未完成阶段分开写。
- 证据：`00-meta.md` 路线图 S3/S4 行；同文件 `progress: 2/4` 段。

### F-002 · D-001 §3 用本波检查点号描述「未来整改波次」阶段

- 严重度：low
- 建议：**non-blocking**
- 状态：open
- 描述：D-002 已把本波 S2 定义为台账落盘、整改 deferred。D-001 §3 仍写「**S2 冻结**：用户对 F-01～F-14 逐项 in-scope / defer 裁决」「**S3 分批**实施」「**S4** 回归」。读者会把这组 S2/S3/S4 当成 GOAL-015 检查点，与 `00-meta` 的 S2=落盘 / S3=independent / S4=关门冲突。建议在 D-001 §3 或决策索引注明：该节阶段号属于**未来整改波次**，不是本波检查点。
- 证据：`01-decision/D-001-w14-scope-and-findings.md` §3；对照 `01-decision/D-002-w14-delivery-boundary.md`、`00-meta.md` 路线图。

### F-003 · F-14 将通知空态写成硬编码英文，证据不成立

- 严重度：low
- 建议：**non-blocking**
- 状态：open
- 描述：D-001 F-14 写「通知空态文案『No items match.』」，引用 `notification-bell.tsx:179`。实际文件为 `apps/web/src/app/notification-bell.tsx:179`，文案走 `t("feedback.noItemsMatch")`；`notification-center.tsx:164` 同样。`zh-CN.json` 已是「没有匹配的条目。」。该子项应降级或删除，不宜作为「英文硬编码」证据。F-14 其余子项抽验成立。同类误判 E-001 已纠正 2 条，本条漏网。
- 证据：`apps/web/src/app/notification-bell.tsx:179`；`apps/web/src/i18n/messages/zh-CN.json`（`feedback.noItemsMatch`）；`01-decision/D-001-w14-scope-and-findings.md` F-14 行。

## 必改项汇总

无 required / 必改 finding。

## 与既有意见的异同

- **A-001 self · S1 · pass**：同意 S1 只读、台账为改进候选、无本阶段 required。本意见把范围扩到 S2 台账 + D-002 关门边界，并补 3 条 non-blocking hygiene / 精度项。不与 A-001 冲突。

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** S1 抽验与 E-001 两条修正属实；F-01～F-14 是带证据的改进候选而非已整改；D-002 把整改 deferred 与用户「审视+落盘」请求一致，I-001 的 P-005 deferred 三字段齐备，不是借重组绕过 required 门禁，也不是静默丢弃。

**可放行 W14 关门（S4）**：无未闭合 required finding，无对本波关门到期的 required 信息项。`/govern` 响应本意见后即可做 S4（goal-tree / workspace 已基本同步；git 提交；关门自审 A-003）。响应时建议顺手处理 F-001～F-003（改检查点措辞、给 D-001 §3 加「未来整改波次」标注、从 F-14 去掉通知空态子项）。

未来整改波次启动时，I-001 恢复为该波次 S2 的开放 required，须用户书面裁决 in-scope / 优先级（含 F-01 handler 目录、F-04 本地化方案、F-08 调试框去留）。

## 声明

本意见不修改 status / progress / 方案正文 / goal-tree / workspace。响应由 `/govern` 处理。
