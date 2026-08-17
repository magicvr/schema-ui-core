---
id: GOAL-015-w14-user-perspective-review
doc: audit-entry
record_id: A-008
source: independent
scope: S5 关门工作结果（execution-facts + close-out）· F-01～F-14 as-built 对照 D-003 / 子目标冻结方案
verdict: conditional
status: recorded
auditor: grok-build (grok-4.6 · reasoning high)
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# A-008 · W14 工作结果独立交叉审计（2026-08-17）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：execution-facts + close-out · [workspace-010 / GOAL-015] S1～S5 与整改子目标 R1～R4 声称的 F-01～F-14 工作结果
- **verdict**：**conditional**
- **工作区**：`workspace-010-design-implementation-conformance`（`root_goal` = `GOAL-001-design-implementation-conformance`；`canonical_scope` 已核对；`shared_materials_catalog: none`；`primary_plan` = `VP-010-design-implementation-conformance`）

## 范围与区间

- **covered**：GOAL-015 五件套与 ledger（D-001～D-003、E-001～E-010、A-001～A-007）；同区子目标 GOAL-016～GOAL-019 的冻结方案 / 实施记录 / 关门审计（作为 GOAL-015 R1～R4 证据，不改子目标状态）；F-01～F-14 对应 as-built 抽验；本轮独立复跑 `apps/api` `go test ./internal/handler/ ./internal/errorcatalog/`（handler 包 **ok**）与 `apps/web` 针对性 vitest **23/23 pass**。
- **不 covered**：完整 Web 1041 与 `go test ./...` 全量（E-010 / A-007 主张存在，本轮未整包复跑）；浏览器活栈点验；其他工作区上下文。
- **本意见不修改** `00-meta` / `status` / `progress` / 方案正文 / `goal-tree.md` / `workspace.md`。

## 成果（有证据）

S1～S4 审视/台账/I-001 裁决路径可追溯（E-001、D-001、A-002、D-003、E-003/E-005）。整改按 D-003 分批 A→C→D→B 挂在 GOAL-015 下，四子目标均标 done。对照冻结方案，**大部分 F 项 as-built 成立**：

| 台账项 | 独立核对 |
|--------|----------|
| F-01 | `GET /api/scheduled-tasks/handlers`（`tasks.read`）+ create/edit `handler` select `optionsSource` 成立。`HandlerKeys()` 仍仅 `system.noop`——与 GOAL-016 D-001「v1 只有 system.noop」一致，属冻结范围内残余（见 recommended F-004）。 |
| F-02 | `data-permission-scopes` 自定义组件：选用户 → GET scopes → PATCH `{userId,scopes}`；schema body 已挂 custom 节点；i18n 键齐。 |
| F-03 | ExtraQuery `event/actorName/from/to` + `GET /api/operations/export` + activity 搜索字段 + `activity-export` 组件成立。非法日期走 DomainError `INVALID_DATE_FILTER`（测试 `operations_test.go` 覆盖 400）。**目录缺口见 F-002**。 |
| F-04 | 迁移 Version 37 增 `title_key`/`body_key`；`NotifyAccountEvent` 写 key、title/body 空串；center/bell 优先 `t(titleKey)`；en/zh `notification.account.*` 齐。 |
| F-05 | recycle-bin / wallet accounts·entries·reconcile-runs / per-task runs 非法 `page`/`pageSize` 返回 400；policies 内存切片分页并回传真实 `page`/`pageSize`。 |
| F-06 | `OPERATION_NOT_FOUND` 入目录与契约冻结集；`INVALID_SCOPE_BODY` / `INVALID_WALLET_OWNER` / `INVALID_WALLET_ACCOUNT` / `INVALID_WALLET_STATUS` 已接线。GOAL-019 A-001 required 在当前代码上可复核为已修。 |
| F-07（部分） | 通知 `q` 已 `ToLower`；wallet 账户搜索扩到 owner_id/owner_type/currency；ledger `q` + `entryType` 贯通 handler→repository，schema 有筛选。recycle-bin **API** 接受 `sort`/`order`（deletedAt/resource/actorName）。**UI 未暴露见 F-001**。 |
| F-08 | `App.tsx` 标题区仅面包屑 + 标题，无 pageId/route 技术框。 |
| F-09 | toast 可见文本不再前缀错误码；`data-feedback-code`/`title` 保留码；renderer 英文反馈带 `messageKey`。 |
| F-10 | `PageSchemaErrorSurface` 友好标题 + 重新加载 + `<details>` 技术详情；`shell.pageSchemaError.*` 双目录存在。 |
| F-11 | 必填 label 追加 `*`，主要控件 `aria-required`。 |
| F-12 | `confirm.tsx`：初始焦点 Cancel、ESC、Tab 圈。 |
| F-13 | 桌面 `<tr>` `tabIndex=0` + Enter/Space；交互单元格不触发行选中。冻结未要求 `role`（桌面用原生 row）。 |
| F-14 | 移动卡片 `fields.slice(1)` 全列；select 空值占位；第三次点击清排序；Tabs 方向键/`aria-controls`；禁用按钮 `title`；`aria-describedby`；`<sm` LocaleSwitcher 不再 hidden；通知空态 `schema.notifications.empty`。 |

本轮复跑：`go test ./internal/handler/ ./internal/errorcatalog/` exit 0（handler ~190s）；vitest 5 文件 23 tests pass（data-permission-scopes / activity-export / notification-center / error-localization / data-table）。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1～S4 审视 + 台账 + I-001 裁决 + 结构修正可追溯 | 达成 | E-001～E-005；D-003；A-002/A-006 |
| R1～R4 子目标按 D-003 批次完成 F-01～F-14 | **部分** | 多数 as-built 成立；F-07 回收站排序 UI 未接线（F-001）；F-03 新码未入目录（F-002） |
| S5 关门主张与证据、信息项一致 | **部分** | 子目标 done 与代码大体一致；A-007「信息门禁均 closed」与 I-002 `collecting`、00-meta 残留「active · 4/8」矛盾（F-003） |
| 无到期 required 信息项阻断本 scope | 达成 | I-001 closed；I-002 为 non-blocking（但应变 closed，见 F-003） |
| 工作区绑定 / 无共享资料冒充证据 | 达成 | workspace.md；`shared_materials_catalog: none` |

## Findings

### F-001 · F-07 回收站排序仅 API 接线，页面列未 `sortable`（用户仍不能排序）

- 严重度：med
- 建议：**required**
- 状态：open
- 描述：原台账 F-07 与 GOAL-019 D-001 要求 recycle-bin **暴露** `sort`/`order`（deletedAt / resource / actorName）。handler 已校验并下传仓库（`recyclebin.go:63-79`，`repository.go:98-104`）。但 `recycle-bin.json` 各列**没有** `sortable: true`；渲染器只在 `column.sortable === true` 时打开排序（`schema-table.tsx:695`）。activity 页同类列已标 `sortable`，回收站没有。用户在回收站页仍看不到排序控件，「排序缺口」未从产品面闭合。GOAL-019 A-001/A-002 与 GOAL-015 A-007 把 F-07 写成已完成，过宽。
- 证据：`apps/api/internal/modules/recyclebin/schema/recycle-bin.json:88-113`；`apps/web/src/renderer/schema-table.tsx:695`；对照 `activity.json:142-154`；GOAL-019 `01-decision/D-001-s1-freeze.md` F-07 条。

### F-002 · F-03 新增 `INVALID_DATE_FILTER` 未入错误目录，且契约扫描扫不到 DomainError

- 严重度：med
- 建议：**required**
- 状态：open
- 描述：非法 `from`/`to` 经 `DomainError` → `writeLocalizedError`。未入 `errorcatalog.Catalog` 时信封只有英文 `message`、无 `messageKey`（`localize.go:16-21`）。`error_contract_test.go` 的字面量扫描匹配 `writeLocalizedError(..., "CODE")`，**扫不到** `domainErr.Code`；`INVALID_DATE_FILTER` 也不在 frozen 集。结果：活动日志日期过滤失败在 zh-CN 下仍是英文硬编码——与 F-06「缺目录则无法本地化」同类，且是本波 F-03 **新引入**的码。
- 证据：`apps/api/internal/handler/operations.go:112`；`apps/api/internal/errorcatalog/errorcatalog.go`（无该码）；`apps/api/internal/handler/error_contract_test.go:86-87,69-84`；`apps/api/internal/handler/localize.go:16-21`；`apps/api/internal/handler/resources.go:425-428`。

### F-003 · S5 关门台账与权威状态不一致；I-002 仍 collecting；A-007 过述「信息门禁均 closed」

- 严重度：med
- 建议：**required**
- 状态：open
- 描述：`00-meta` frontmatter 为 `done` / `8/8`，但「信息就绪」段仍写「GOAL-015 **active · 4/8**」；`01-decision.md` 待决节仍写「active · 4/8，整改承接中」。I-002（handler 目录暴露方式）在 D-003 已裁决「新增端点」且 GOAL-016 已实现后仍标 **collecting**。A-007 写「信息门禁均 closed」与 `03-audit.md` 索引表 I-002=collecting 直接矛盾。四整改子目标 `00-meta` 的 I-001/I-002 仍为 `open` /「待 S1」，而其 D-001 已 closed——子目标关门包装同样未刷权威表。I-002 为 non-blocking，不单独构成 P-005 到期 required 阻断，但 **S5 终审声明不实 / 权威段过期** 使关门包不可无条件采信。
- 证据：GOAL-015 `00-meta.md` L18 与 L61；`01-decision.md` L30；`03-audit.md` L20；`03-audit/A-007-s5-closeout-self.md` L23；GOAL-016～019 `00-meta.md` 信息表 vs 各自 `01-decision/D-001-s1-freeze.md`。

### F-004 · F-01 产品残余：目录仍只有 `system.noop`，建任务仍空转

- 严重度：low
- 建议：**recommended**
- 状态：open
- 描述：D-003 / GOAL-016 D-001 明确 v1 只列已注册 handler，现仅 `system.noop`。F-01「可指定 handler」按冻结范围已完成；原用户影响「建了不执行」作为产品能力仍在。不是本波违约，但不应写成「定时任务页核心功能已可用」。
- 证据：`apps/api/internal/modules/scheduledtasks/scheduler.go:45-49`；GOAL-016 D-001 F-01 条。

### F-005 · `error.emptySelection` 在 en/zh 目录各出现两次

- 严重度：low
- 建议：**recommended**
- 状态：open
- 描述：`en-US.json` / `zh-CN.json` 在约 L450（API 批量 ids 语义）与 L700（表格「请先选行」）重复同一 key。`JSON.parse` 后键胜出，F-09 toast 碰巧用到后写文案；目录不自洽，后续改一处易漏。
- 证据：`apps/web/src/i18n/messages/zh-CN.json:450,700`；`en-US.json:450,700`。

### F-006 · S5 回归证据仅为断言；本轮只复跑子集

- 严重度：low
- 建议：**recommended**
- 状态：open
- 描述：E-010 / A-007 写「Go 全量、Web 1041/1041、tsc、build 通过」，无日志附件。本轮独立复跑 handler 包与 23 个相关 vitest 均绿，**不否定**回归主张，但不能把 1041/1041 整包当作本意见已复核事实。
- 证据：`02-execution/E-010-s5-closeout.md`；本轮命令记录见「成果」。

## 必改项汇总

1. **F-001**：给回收站表 `resource` / `actorName` / `deletedAt`（及冻结声明的字段）加 `sortable: true`，或书面收窄 F-07「暴露」为仅 API 并改冻结/用户影响表述。
2. **F-002**：将 `INVALID_DATE_FILTER` 写入 errorcatalog + 契约冻结集（建议扩扫描以覆盖 DomainError 码）+ en/zh `messageKey`。
3. **F-003**：刷新 GOAL-015（及子目标）权威段：I-002 按 D-003/实现标 closed；删/改「active · 4/8」残留；A-007 响应里纠正「信息门禁均 closed」。

## 与既有意见的异同

- **A-002 independent · S1/S2 · pass**：同意当时审视台账属实、整改尚未实施。本意见审的是 **整改后工作结果**，不回改 A-002。
- **A-007 self · S5 · pass**：同意四子目标已立项关门、多数代码与冻结一致。**不同意**「F-01～F-14 已全部从用户面闭合」和「信息门禁均 closed」。
- **GOAL-016 A-001 / GOAL-019 A-001**：其子项 required（迁移号、错误码细分、ledger `q`、台账回填）在当前树可复核为已修；本意见不重开那些编号。F-001/F-002 是关门后仍在的新缺口（回收站 UI 排序、新错误码目录）。
- 与 A-007 **无冲突裁决需求**（A-007 为 self pass，本意见为 independent conditional）：编排器须响应本意见 required，不得以 A-007 覆盖本条。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** W14 审视与分批整改主体工作真实存在，F-01～F-06 / F-08～F-14 抽验大体与冻结方案一致；独立复跑的 handler 包与相关前端测试通过。不能无条件维持「工作结果已全部闭合」：回收站排序对用户仍不可用、本波新错误码未本地化、关门权威段与 I-002 状态过期。

建议 `/govern`：先修 F-001/F-002 并刷 F-003 台账，再复审本意见；F-004～F-006 可一并 hygiene 或用户书面 residual。在 required 合法闭合前，**不要**把 A-007 pass 当作 W14 工作结果的最终独立会签。

## 声明

本意见不修改 status / progress / 方案正文 / goal-tree / workspace。响应由 `/govern` 处理。
