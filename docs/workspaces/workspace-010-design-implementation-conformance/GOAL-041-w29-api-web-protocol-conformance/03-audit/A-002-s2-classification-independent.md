---
id: GOAL-041-w29-api-web-protocol-conformance
doc: audit-entry
record_id: A-002
source: independent
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# A-002 · S2 分类与方案冻结 · independent 审计

## A-002 · S2 证据分类与方案冻结（2026-09-06）
- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`）
- **类型 / scope**：design-plan（S2 方案冻结）；C-001～C-014 证据分类 + D-002 冻结范围/S4·S5 契约 + I-003 / I-004 / I-008
- **verdict**：conditional

### 范围与区间

工作区：`workspace-010-design-implementation-conformance`（`root_goal` = `GOAL-001-design-implementation-conformance`；canonical = `docs/workspaces/workspace-010-design-implementation-conformance/`；`shared_materials_catalog: none`）。本条只审 GOAL-041 S2 方案冻结门禁，不改 `00-meta` status/progress、goal-tree 或方案正文。

本意见独立核验：14 个候选的处置类别是否证据充分、唯一且如实；D-002 冻结范围与 S4/S5 契约是否成立；I-003（分类 verified）/ I-004（upstream gap 不适用）/ I-008（cross 审计）是否一致；是否存在漏判候选或类别。对照上游 `schema-ui-docs@v2.9.0`（本机 clone tag 解引用 `81aa1d8954717f4ebdcc695eed6fafaeafcebe8d`）与本仓代码/测试。

### 成果（有证据）

1. **工作区绑定合格**：workspace id / Root / canonical 与 GOAL-041 `parent` 一致；无共享资料引用被当成事实。
2. **分类完整性**：C-001～C-014 均给出处置类别与 file:line 证据；**upstream-protocol-gap = 0** 与 I-004「不适用 / 不建空报告」一致；custom 候选仅 C-009，I-005 仍属 S3 未到期。
3. **逐项实查通过（除 C-007 定量主张，见 F-002）**：
   - C-001：本仓 `docs/schemas/fixture-suite.schema.json` LF digest = `a93f500bcb94f93e7a289af700227a407688f8205218fe5681269956b91fd1a3`（本次重算）；`provenance-v2.9.json` 将该路径登记为 `docs/schemas/…`（上游 tag 真实路径为 `conformance/schemas/fixture-suite.schema.json`）；note 写「11 件 + 20 suites」，artifacts 实为 11 schema + 19 fixture = 30。**implementation-gap** 成立。
   - C-002：`generate-claim.mjs:99` `pinnedUpstream.artifactVersion = "2.8.0"` vs `claim.protocolArtifact.artifactVersion = "2.9.0"`（`ARTIFACT_VERSION`）；`public/protocol/conformance-local-report.json:55` 现为 `2.8.0`；`claim-artifact.test.ts` 不断言该字段。**implementation-gap** 成立。
   - C-003：`upstream-fixtures.test.ts` 固定 `provenance.json`（`artifactVersion "2.7.0"` / `ca9e5fe…`）同时硬编码 2.9 线 digest（`34a3354e…` / `5f14de6…` / `d56d933…`，与 `provenance-v2.9.json` 同值）；`provenance-v2.8.json` 无测试/代码消费者（仅 `boot.ts`/`bootstrap.ts` 头注释与历史文档）；`apps/web/src/protocol/README.md:20-21` 仍写 production host 只接受 2.7。**implementation-gap** 成立。
   - C-004 模型层：17 个 `fragment.json` 全部 `"protocolVersion": "2.7"`；Web `APP_MANIFEST_SUPPORTED_PROTOCOL_VERSIONS = ["2.7","2.8","2.9"]`；35 页中 8 页 `meta.protocolVersion: "2.9"`（wallet×4、digitaloffer×3、dictionary-entries）；上游 `08-renderer-spec.md` §3 与 version-negotiation fixtures 允许 envelope/页面解耦。**模型 no-gap 成立**。残余实施项见 F-001。
   - C-005：`wallet-entries.json:116` 使用 `$context.route.params.id`；`dictionary-entries.json` 使用 `readOnly: true` ×4 + `$context.route.params.dictKey`；`wallet.json:366` 使用 `readOnly: true`；digitaloffer 两页声明 `data.route-binding` / `form.controls.readonly` 且全文无 `$context.route.*` / `readOnly`。上游 `page.schema.json` 为单向约束。**no-gap** 成立。
   - C-006：`app-manifest.mvp-dogfood.json` 13 个 `pageId`；`app-manifest.admin-dogfood.json` 26 个 `pageId`；与 E-002 真实投影 MVP 6 / Admin 22 / Demo 14 不同。**implementation-gap** 成立。
   - C-008：E-002 真实 `ForModulesWithFragments` 投影与「无 HTTP 快照」作为 S5 矩阵定义项；**excluded→S5** 类别成立（本波未核跑 Go 投影，沿用 E-002）。
   - C-009：15 个 `registerCustomComponent(...)` 均存在；`runtime-schema-validate.ts:35-51` 本地扩展 `component` 且注释写明不改上游 schema 字节；key 无 namespace；上游 `08-renderer-spec.md` §1.1 Host Extension 模型允许。**custom-extension-candidate** 成立（S3 P-004）。守卫覆盖缺口见 F-003。
   - C-010：未知标准 node → `parseRenderNode` `RENDER_UNKNOWN_NODE_TYPE` fail-closed（`render.types.ts:397-407`）；未知 custom → `render.tsx:2915-2927` 内联 `<p>unknown custom component: …</p>`，无 `console.error`。上游 §1.1 要求未安装扩展以 `UNKNOWN_COMPONENT_TYPE` 拒绝；§1.3 要求明显占位 + `console.error`；vendor fixture `uninstalled-extension-is-unknown-component` 期望 `{ok:false, code:"UNKNOWN_COMPONENT_TYPE"}`。**implementation-gap** 成立。
   - C-011：5 handler 均在 `CUSTOM_HANDLER_URLS`（`render.tsx:332-343`）；非白名单 `CUSTOM_HANDLER_NOT_FOUND`（:353）。上游 `07-actions-contract.md` §6 要求白名单+拒绝。**no-gap** 成立。
   - C-012：内页 navigate action 存在（`openEntries` / `openRuns` / `openTelegramOperator` / `openInvites`）；`App.tsx:203-212` breadcrumb 父页映射含 dictionary-entries / task-runs / wallet-entries / users-invites / telegram-operator；Host 铃铛 `App.tsx:982-985` → `/notifications` 与 `?open=`；notifications fragment `navigation.user: []`。**no-gap** 成立。
   - C-013：`profile.go` 中 `admin.data-transfer` / `admin.mfa` / `admin.login-captcha` 的 `ContributionKeys` 均无 `Pages`。**explicitly-out** 成立。
   - C-014：GOAL-004 95/95 仅历史边界；D-002 §5 方法论规则。**excluded** 成立。
4. **self A-001 F-001 判定成立**（独立加强证据，见本条 F-001）：生产路径 `loadPageDocument`（`load-page.ts:75-155`，仅 D-VAL + pageId）→ `App.tsx:598 RenderPage` 未调用 `negotiateVersion`；`negotiateVersion` 仅 `stage3-fixtures.test.ts` 消费。上游 `08-renderer-spec.md` §3.4：页面 `requiredCapabilities` 任一缺失则拒绝渲染。

### 对照成功标准

| 标准（S2） | 状态 | 证据 |
|---|---|---|
| 每项候选有协议/代码/测试证据 | 部分 | 13/14 定量主张可重复核对；C-007「D-VAL 35/35」不实（F-002） |
| 归入唯一处置类别 | 部分 | 13 项唯一；C-004 为 no-gap + implementation-gap 双标（F-004，不否定残余项） |
| 完成 self + independent 方案级 cross 审视 | 进行中 | A-001 self + 本条 independent；合并响应前 I-008 仍 open |
| I-003 关闭（分类证据） | 有条件 | 类别齐全；C-007 证据行须纠正后才算「如实」 |
| I-004 不适用（无 upstream gap） | 达成 | 独立同意 upstream-protocol-gap = 0 |
| D-002 S4/S5 契约成立 | 部分 | S4 清单覆盖 C-001/002/003/004-残余/006/010 成立；S5「D-VAL 已全量 35/35」不成立（F-002） |

### Findings

#### F-001 · 生产页面级能力协商缺位 + claim/HOST_SUPPORT 覆盖不足（med · required）
- **严重度**：med
- **建议**：required（并入 S4 工作清单执行闭合；不得把「已纳入计划」当完成）
- **描述**：与 self A-001 F-001 **同意**。上游 `08-renderer-spec.md` §3.2–§3.4 要求 Renderer 按页做精确版本匹配，再校验 `requiredCapabilities ⊆ supportedCapabilities`，缺失则拒绝渲染（`MISSING_REQUIRED_CAPABILITY`）。本仓 `version-negotiate.ts:154-176` 实现该逻辑且 version-negotiation fixtures 全绿，但生产加载/渲染链未接线：`load-page.ts` 无协商；`RenderPage`（`render.tsx:3026`）只包 `SchemaCrudProvider` 后走节点级门禁（`form-controls.types.ts` / `permissions.ts` / `gateDataRouteBinding`）。claim `support.capabilities` 与 `boot.ts:66-78 HOST_SUPPORT` 仅 7 项。本次对 35 页 `meta.requiredCapabilities` 求并集，**不在 HOST_SUPPORT 内的页面能力共 11 项**：`actions.batch.request`、`actions.page.trigger`、`actions.row.navigate`、`actions.row.request`、`actions.upload`、`form.controls.advanced`、`form.controls.extended`、`form.record.load`、`permissions.inheritance`、`table.selection`、`table.sort`（self 列举 9 项并写「等」；独立补上 `table.selection` 与 `actions.batch.request`）。这些能力本仓均有实现痕迹，但若只接线页面级门禁、不同步扩展 HOST_SUPPORT/claim，现网多数页面会立即 fail-closed——S4 必须两项一起做。
- **证据**：`docs/08-renderer-spec.md` §3.4（上游 clone）；`apps/web/src/protocol/conformance/version-negotiate.ts:154-176`；`apps/web/src/protocol/load-page.ts:75-155`；`apps/web/src/app/App.tsx:560-609`；`apps/web/src/renderer/render.tsx:3026-3050`；`apps/web/scripts/generate-claim.mjs:118-131`；`apps/web/src/host/boot.ts:66-78`；`apps/web/public/protocol/conformance-claim.json:22-30`；35 页 schema `meta.requiredCapabilities` 并集（本次枚举）。
- **状态**：open。已由 D-002 列入 S4，但本审计不把「已纳入 S4 计划」当闭合。
- **关联**：C-004 残余项；I-003 / I-008；self A-001 F-001。

#### F-002 · C-007 / D-002 §3「D-VAL 全量 35/35」名不副实（med · required）
- **严重度**：med
- **建议**：required（纠正 S2 证据行与 S5 验收契约基线后再接受冻结；walker 递归化属 S5 分母工作，不是已完成事实）
- **描述**：S2 与 D-002 §3 写「D-VAL 已全量 35/35 绿（`all-module-schemas-dval.test.ts`）」。该测试 `collectSchemaFiles()` 只读 `apps/api/modules/<top>/schema/*.json`，**不递归**嵌套模块。本次实查：一层 walker = **25** 份；递归全部 schema json = **35** 份。漏扫：`channel/telegram/schema/` 2 页（`telegram-operator`、`telegram-settings`）+ `dev/examples/schema/` 8 页。其中 4 个 example 页另由 `representative-pages.test.tsx` 的 `MIGRATED_PAGE_IDS` 做 D-VAL（data-table / search-form-table / form-controls / form-with-reactions），故**完全未进入 D-VAL 的生产 schema 至少 6 页**（telegram×2 + `admin-list-batch` / `data-display` / `form-with-upload` / `overview`）。C-007 归 **excluded→S5** 的类别仍合理（验收分母，非产品实现缺口），但「已全量 35/35」不能作为 S5 已完成基线写入冻结契约。未发现因此需要新增第 15 个候选：纠正 C-007 证据与 D-002 §3 即可。
- **证据**：`apps/web/src/protocol/all-module-schemas-dval.test.ts:19-34,102-112`；本次对 `apps/api/modules/**/schema/*.json` 的一层 vs 递归计数（25 vs 35）；`representative-pages.test.tsx:56-67`（10 个 pageId）；D-002 §3 第 1 款；`attachments/S2-candidate-classification.md` C-007。
- **状态**：open。
- **关联**：C-007；I-003（证据如实）；I-006 / D-002 §3 S5 契约。

#### F-003 · C-009 W25 守卫未覆盖嵌套模块（low · recommended）
- **严重度**：low
- **建议**：recommended（S3 固定 custom 边界或 S5 分母时一并修守卫；不改变 C-009 类别）
- **描述**：`custom-components.schema.test.ts` 使用与 D-VAL 相同的一层 `modules/<top>/schema` 遍历，且测试文件未 side-effect import `telegram-admin-tab`。15 个 `registerCustomComponent` 在组件模块中确实存在（含 `telegram-admin-tab.tsx:1089`，`main.tsx` 有生产 import），故「当前键已注册」的事实成立，**custom-extension-candidate** 不改。但「W25 守卫覆盖 15 键」过宽：telegram 两页的 custom 引用不在该测试扫描范围内，守卫无法防止其注册被删后仍绿。
- **证据**：`apps/web/src/renderer/custom-components.schema.test.ts:18-33,63-84`；`apps/web/src/components/telegram-admin-tab.tsx:1089`；`apps/web/src/main.tsx:29`。
- **状态**：open。
- **关联**：C-009；F-002 同源 walker。

#### F-004 · C-004 双类别与「唯一处置」口径不一致（low · recommended）
- **严重度**：low
- **建议**：recommended（编排器响应时二选一澄清，不阻塞 F-001 实质）
- **描述**：D-001 §2 / D-002 §1 要求每个候选落入**唯一**类别；C-004 写成「no-gap（模型）+ implementation-gap（残余）」。残余项真实（F-001），双标是诚实的，但形式上不唯一。建议在响应中显式拆成「C-004 模型 no-gap」与「C-004-R / 并入 S4 的 implementation-gap」，或将 C-004 整体改标 implementation-gap 并保留「模型合法」说明。不要求新建 C-015，除非编排器认为拆分更清晰。
- **证据**：D-002 §1「必须且只能落入六类之一」；D-002 §2 C-004 行；`S2-candidate-classification.md` 汇总表 C-004。
- **状态**：open。
- **关联**：C-004；F-001。

#### F-005 · 台账内部计数/去向不一致（low · recommended）
- **严重度**：low
- **建议**：recommended
- **描述**：(a) E-003 阶段结论写 implementation-gap **×5**（C-001/002/003/006/010），漏 C-004 残余；S2/D-002/A-001 为 **×6**。若有人按 E-003 列 S4 清单会漏 F-001。(b) C-005 recommended 卫生项：汇总表写「S4」，正文与 D-002 写「并入 C-009/S3」。以 D-002 为准即可，但应改汇总表避免执行分叉。
- **证据**：`02-execution/E-003-s2-classification-and-cross-review.md`「阶段结论」；`S2-candidate-classification.md` 汇总 C-005 行 vs 逐项 C-005；D-002 §2 C-005。
- **状态**：open。

### 必改项汇总

- **required**：F-001（S4：生产页面级版本+能力门禁 + claim/HOST_SUPPORT 扩至页面实际声明且已实现的能力全集，含 `table.selection` / `actions.batch.request`；一致性守卫）。
- **required**：F-002（纠正 C-007 / D-002 §3：D-VAL 基线改为 walker 实扫 25/35，并写明 6 页尚未进入任何 D-VAL；S5 契约改为递归遍历全部 `**/schema/*.json` 后再称全量）。
- **recommended**：F-003（C-009 守卫递归 + telegram-admin-tab import）、F-004（C-004 唯一类别澄清）、F-005（E-003 计数与 C-005 去向）。

### 与既有意见的异同（self A-001）

| 点 | self A-001 | independent A-002 |
|---|---|---|
| verdict | conditional | **conditional**（同意尺度） |
| F-001 生产页面级能力协商 / claim 覆盖 | required · med；并入 C-004 残余与 S4 | **成立并加强**：生产链 `load-page → RenderPage` 无 `negotiateVersion`；上游 §3.4 明文拒绝；页面能力并集相对 HOST_SUPPORT 缺 **11** 项（self 列 9 +「等」） |
| F-002 台账卫生（claim buildId / README） | recommended；并入 C-002/C-003 | **同意并入 C-002/C-003**，本条不重复编号 |
| C-007 35/35 | 未质疑 | **新增 required F-002**：一层 walker 25/35，主张不实 |
| C-009 15 键注册 | 接受 W25 守卫 | 类别同意；**新增 recommended F-003** 守卫未扫 telegram/examples |
| C-004 唯一性 | 接受双标 | **recommended F-004** |
| upstream-protocol-gap = 0 | 同意 | **同意** |
| 漏判候选 | 无 | 无新候选号；C-007 证据纠正即可，不必 C-015 |

无与 self 相反的 required 结论，故无 P-004 冲突项。F-001 不是否决 self，而是独立确认。

### 结论 + 建议给编排器/用户的下一步

S2 分类总体可核对：implementation-gap（C-001/002/003/004-残余/006/010）、custom-extension-candidate（C-009）、explicitly-out（C-013）、excluded（C-007/008→S5，C-014 方法论）、no-gap（C-005/011/012 + C-004 模型）与 **upstream-protocol-gap = 0** 成立。I-004「不适用」成立。I-003「14 项均有类别」成立，但 C-007 定量证据必须按 F-002 纠正后才算如实。I-008 S2 腿现为 A-001 + A-002，**required findings 未闭合前不得把 S2 标完成，也不得进入 S3/S4 实施**（D-002 §4）。

建议 `/govern`：

1. 合并响应 A-001 + A-002。
2. 纠正 D-002 §3 / S2 C-007 的 D-VAL 基线（F-002）；澄清 C-004 唯一口径（F-004）与 E-003 计数（F-005）。
3. F-001 保持 S4 必做；F-002 为方案契约修正（可在进入 S3 前用决策补丁完成，无需等产品代码）。
4. C-009 仍走 S3 P-004，不把 custom 当默认逃生口。

### 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
