---
title: S2 · C-001～C-014 证据分类矩阵（唯一处置类别）
status: active
created: 2026-09-06
updated: 2026-09-06
parent: GOAL-041-w29-api-web-protocol-conformance
version: 0.1.0
---

# S2 · C-001～C-014 证据分类矩阵

> 本附件把 S1 登记的 14 个候选逐项补齐「上游协议 + 本仓代码/测试 + 运行时」证据并给出**唯一处置类别**。分类口径沿用 D-001 §2 五类（implementation-gap / upstream-protocol-gap / custom-extension-candidate / explicitly-out / excluded），并对「证据证实符合、无偏差」的候选给出第 6 类结论 **no-gap（已核验符合，无处置类别适用）**——它不是放行信号，而是带证据的分类结论。处置决策见 [D-002-s2-classification-and-scheme-freeze](../01-decision/D-002-s2-classification-and-scheme-freeze.md)。

被审身份：上游 `schema-ui-docs@v2.9.0`（tag `463d563a…` → commit `81aa1d8954717f4ebdcc695eed6fafaeafcebe8d`，已复核一致）；本仓 HEAD `1427f9b7`（产品代码快照 = S1 的 `ce66abab`）。

---

## 分类汇总

| ID | 范围 | 处置类别 | 证据要点 | 后续 |
|----|------|----------|----------|------|
| C-001 | v2.9 provenance 路径与数量口径 | **implementation-gap** | 上游 tag 中 fixture-suite schema 真实路径 `conformance/schemas/fixture-suite.schema.json`；本仓 provenance-v2.9.json 登记为 `docs/schemas/fixture-suite.schema.json`；LF digest 一致（`a93f500b…`，本次重算双路径均匹配）；note 声称「20 个 fixture 套件」与实际 vendor 19 不符 | S4 修正 provenance 路径/note + stage3 守卫 |
| C-002 | claim 证据报告 artifactVersion | **implementation-gap** | `generate-claim.mjs:99` `report.pinnedUpstream.artifactVersion = "2.8.0"`，claim `protocolArtifact.artifactVersion = "2.9.0"`；已生成 `public/protocol/conformance-local-report.json` 现为 2.8.0 | S4 改 2.9.0 + 重生成 claim 三件套 |
| C-003 | legacy provenance / 测试双权威 | **implementation-gap** | `upstream-fixtures.test.ts` 固定 v2.7 provenance.json（含 2.9 线 digest 混标）；`provenance-v2.8.json` 无代码消费者；stage3 固定 v2.9；README 称「production host 仅接受 2.7」过时 | S4 台账卫生：标注兼容基线/清理死数据/澄清 current pin 权威 |
| C-004 | Manifest 协议版本协商 | **implementation-gap（主类；模型层 no-gap 为上下文）** | 全部 17 fragment 发布 2.7 envelope；Web 适配器支持 2.7/2.8/2.9；claim manifestVersions [2.7,2.8,2.9]；上游 version-negotiation 含 decoupledVersions 语义；8 页声明 2.9 且渲染层按页做版本+能力门禁。**S2 新发现（A-001/A-002 F-001）**：页面级能力协商（MISSING_REQUIRED_CAPABILITY）仅存在于 fixture 适配器 `version-negotiate.ts`，未接线生产 `RenderPage`；claim/HOST_SUPPORT 只声明 7 个能力，而 35 页声明的页面能力并集相对 HOST_SUPPORT 缺 **11 项**（actions.batch.request / actions.page.trigger / actions.row.navigate / actions.row.request / actions.upload / form.controls.advanced / form.controls.extended / form.record.load / permissions.inheritance / table.selection / table.sort） | 版本协商模型合法（上下文）；S4 必须「接线页面级版本+能力门禁」与「扩展 claim/HOST_SUPPORT 至能力全集」同步做（否则现网页面会 fail-closed） |
| C-005 | 能力声明与实际使用 | **no-gap** | wallet-entries / dictionary-entries / wallet 均声明且实际使用（`$context.route.params.*`、`readOnly:true`）；digitaloffer 两页声明未使用 = 上游允许的保守声明（单向约束：用时须声明） | digitaloffer 两页未使用声明 → recommended 清理（随 C-009/S3 custom 边界一并处置） |
| C-006 | 静态 dogfood Manifest | **implementation-gap** | mvp-dogfood 13 页 / admin-dogfood 26 页，均 2.7；与真实 profile 投影（MVP 6 / Admin 22 / Demo 14）不同；手维护、sha 固定、非运行时权威 | S4 加 dogfood ⊆ 模块联合守卫（route/schemaUrl 一致），或改生成 |
| C-007 | 页面链路测试覆盖 | **excluded（→S5 验收分母；D-VAL 守卫已修复）** | D-VAL 结构验证 **35/35 绿**（`all-module-schemas-dval.test.ts` 递归 walker 修复后实测 38 tests，A-002 F-002 纠正并修复）；渲染级 representative 10/35；v2.9 页有定向测试（dictionary-entries / wallet-entries / inner nav） | S2 只定义 S5 分母契约（D-002 §3），S5 按 I-006 执行 |
| C-008 | profile / optional 运行时矩阵 | **excluded（→S5 验收矩阵）** | 真实 `ForModulesWithFragments` 投影 MVP 6 / Admin 22 / Demo 14；Telegram/Digital Offer 仅显式 custom 启用 | S2 定义验收矩阵（D-002 §3），S5 执行并留 HTTP 快照 |
| C-009 | custom component 扩展面 | **custom-extension-candidate** | 15 键（12 body + 1 action content + 2 afterComponent）均有注册；`runtime-schema-validate.ts` 以本地扩展 `component` 字段保持上游 schema 字节不变；key 无统一 namespace；上游 renderer-spec §1.1 Host Extension 模型。**W25 守卫已递归化并补 telegram-admin-tab import（A-002 F-003 修复，实测绿）** | S3 逐键固定 namespace/capability/schema/validator/failure/fixtures + 用户 P-004 |
| C-010 | 未注册 custom 失败语义 | **implementation-gap** | 未知标准 node → `RENDER_UNKNOWN_NODE_TYPE` fail closed；未知 custom → 内联文字 fallback（无 console.error、不明显）；上游 §1.1/§1.3 要求 UNKNOWN_COMPONENT_TYPE 拒绝或明显占位 + console.error；conformance 适配器已实现 UNKNOWN_COMPONENT_TYPE | S4 产品渲染器对齐（明显占位 + console.error，或 fail-closed） |
| C-011 | custom action allowlist | **no-gap** | 5 handler 全在 `CUSTOM_HANDLER_URLS` allowlist；非白名单 → `CUSTOM_HANDLER_NOT_FOUND` fail closed；上游 07-actions-contract §6 明确要求白名单+拒绝 | 无需改动；namespace 归 C-009 记账 |
| C-012 | 内页与导航边界 | **no-gap** | 内页全部登记在 manifest 且经 navigate action / Host 铃铛可达；orphan/重复/缺失 = 0（E-002）；breadcrumb 父页映射就位 | 无需改动；notifications 为 Host 铃铛入口合法 |
| C-013 | API-only / Host-only 表面 | **explicitly-out** | admin.data-transfer / admin.mfa / admin.login-captcha 均为无页面模块（profile.go 描述符无 Pages）；login/session/branding/shell/failure 为 Host 层（GOAL-004 已处置） | UI 经 custom（C-009）/ action（C-011）面记账；不伪装成页面协议能力 |
| C-014 | 历史 GOAL-004 证据边界 | **excluded（方法论规则）** | GOAL-004 D-002 曾对 95/95 Host/App 候选按 ADR-0034 D10/D6 处置；当前 35 页/15 custom/v2.9 增量已扩大 | 本波分类矩阵为现行权威；GOAL-004 仅作历史边界语义 |

**统计**：implementation-gap ×6（C-001/002/003/004/006/010）· no-gap ×3（C-005/011/012，C-004 模型层 no-gap 为上下文）· custom-extension-candidate ×1（C-009）· explicitly-out ×1（C-013）· excluded ×2（C-007/008 与 C-014，C-014 为方法论规则）· **upstream-protocol-gap ×0**。

> no-gap 计数口径：C-004、C-005、C-011、C-012 为「已核验符合」4 项；C-007/C-008 的「S5 执行契约」与 C-014 的方法论规则在 D-002 §3/§5 固化，不另计为 no-gap。

---

## 逐项证据

### C-001 · v2.9 provenance 路径与数量口径 — implementation-gap

**上游事实**：上游 tag `v2.9.0` 中 fixture-suite schema 的真实路径是 `conformance/schemas/fixture-suite.schema.json`（本机 clone 实查）；`docs/schemas/` 下共 10 个 schema 工件（S1 分母表逐行列 10 件）。

**本仓事实**：`apps/web/src/protocol/upstream/provenance-v2.9.json` 将 `docs/schemas/fixture-suite.schema.json` 登记为 `docs/schemas/…`（line 24），与上游路径不一致；LF 归一化 digest 本次重算 = `a93f500bcb94f93e7a289af700227a407688f8205218fe5681269956b91fd1a3`，上游与本仓副本字节一致（仅路径别名不同）；`provenance-v2.9.json` 的 note 声称「docs/schemas 全量 11 件 + 20 个 conformance fixture 套件（scenarios 未 vendor）」，而 artifacts 实际登记 11 schema + 19 fixture cases = 30 项，与「20 suites」自相矛盾（scenarios 13 cases 未 vendor，vendored 为 19）。

**判定**：上游路径明确，本仓 provenance 登记路径与 note 口径错误，属于本仓证据台账卫生问题。无需上游增补（上游无缺口）。修正 = 将 provenance 中该条目路径改为 `conformance/schemas/fixture-suite.schema.json`（保留「本仓 vendor 于 docs/schemas」的映射说明），note 改为「docs/schemas 10 + conformance/schemas 1 机器工件 + 19 vendored suites（scenarios 13 cases 未 vendor）」，并在 `stage3-fixtures.test.ts` 增加「provenance 路径 = 上游 tag 规范路径」的守卫断言。属 S4 证据卫生整改。

### C-002 · claim 证据报告 artifactVersion — implementation-gap

**上游事实**：`host-conformance-claim.schema.json` 约束 `protocolArtifact.artifactVersion`（三部分版本）与 `evidence[].sha256`，不约束 `local-report` 内部的 `pinnedUpstream`（报告为 claim 的 sha256-pinned 证据工件）。

**本仓事实**：`apps/web/scripts/generate-claim.mjs`：`report.pinnedUpstream.artifactVersion = "2.8.0"`（line 99），而 `claim.protocolArtifact.artifactVersion = "2.9.0"`（line 115）、`pinnedUpstream.sourceCommit = "81aa1d8"`、`fixtureSha256 = 89baddbc…`、`protocolContentSha256 = c87c22ad…` 均为 v2.9.0 正式发布绑定（注释 line 31-37 也写明「Formal 2.9.0 release bindings」）。已生成的 `apps/web/public/protocol/conformance-local-report.json` 当前即为 2.8.0。无任何测试断言 `pinnedUpstream.artifactVersion === "2.8.0"`（`claim-artifact.test.ts` 只读 sourceCommit/fixtureSha256/protocolContentSha256）。

**判定**：同一证据报告对同一个 pinned upstream 给出 2.8.0 身份而 claim 为 2.9.0，是事实不一致（历史残留），且与 provenance-v2.9.json、claim、`APP_MANIFEST_SOURCE`（`tree/81aa1d8`）形成第 4 处身份源。属本仓 evidence 工件错误。修正 = `artifactVersion: "2.9.0"` + 重生成 claim 三件套（含 buildId 刷新），加断言防回归。S4 整改。

### C-003 · legacy provenance / 测试双权威 — implementation-gap

**上游事实**：无——这是本仓测试台账布局问题。

**本仓事实**：
- `upstream-fixtures.test.ts` 固定 `upstream/provenance.json`（`artifactVersion "2.7.0"`、sourceCommit `ca9e5fe…`，R3 基线）作为**兼容回归**，但同一文件已被多次重 pin：note 记载 app-manifest/app-navigation cases 算法线升 2.9，`byPath` 断言用的 `APP_MANIFEST_SCHEMA_SHA256`/`APP_MANIFEST_FIXTURE_SHA256`/`APP_NAVIGATION_FIXTURE_SHA256` 均为 v2.9 线 digest（`34a3354e…`/`5f14de6…`/`d56d933…`），即一个标为「2.7.0」的文件装着 2.9 线条目——标签与内容不符。
- `upstream/provenance-v2.8.json` 存在但**无任何代码/测试消费者**（仅 `host/bootstrap.ts` 头注释提及），是死数据。
- `stage3-fixtures.test.ts` 固定 `provenance-v2.9.json` 为 current pin（`artifactVersion "2.9.0"`、sourceCommit `81aa1d8…`，断言逐项 digest）。
- `apps/web/src/protocol/README.md` 仍写「negotiation cases 是 fixture-only；production host 仍只接受 2.7」（R3 时代表述），与 `app-manifest.ts`（支持 2.7/2.8/2.9）和 claim 相矛盾。

**判定**：v2.7 兼容回归与 v2.9 current pin 并存本身合法（兼容基线与现行证据是不同概念），但命名/标注/文档未把边界讲清：`provenance.json` 标签 2.7.0 却含 2.9 线条目、`provenance-v2.8.json` 死数据、README 过时。属台账卫生 implementation-gap。修正 = `provenance.json` 顶部标注「R3 v2.7.0 兼容基线（含后续重 pin 条目，仅回归用；现行 pin 见 provenance-v2.9.json）」；删除或归档 `provenance-v2.8.json`；更新 README 描述 production 版本协商（2.7/2.8/2.9）。S4 整改。

### C-004 · Manifest 协议版本（2.7 envelope + 2.9 页面）— no-gap

**上游事实**：`08-renderer-spec.md`/version-negotiation fixtures 提供**解耦协商**语义（manifest 版本与页面版本分开判定；fixture `decoupledVersions`）；`page.schema.json` 的 `meta.protocolVersion` 是页面文档自身的版本锚点；`app-manifest.schema.json` 不要求页面与 envelope 同版本。

**本仓事实**：17 个 API fragment/manifest provider 全部 `ProtocolVersion: "2.7"`（profile.go 描述符 + provider.go 逐模块核实），`manifest.go` 聚合时要求各 fragment 版本一致（冲突即报错）→ served Manifest envelope 恒为 2.7。Web 适配器 `app-manifest.ts` 支持 `["2.7","2.8","2.9"]`；claim `support.manifestVersions: [2.7,2.8,2.9]`、`pageVersions: [2.7,2.9]`；35 页中 8 页声明 `meta.protocolVersion: 2.9`（digitaloffer ×3、wallet ×4、dictionary-entries）。渲染层按页执行版本+能力门禁：`form-controls.types.ts`（extended/advanced/readonly 字段的版本下限+能力）、`permissions.ts`（PROTOCOL_VERSION_TOO_LOW / CAPABILITY_REQUIRED）、`render.tsx:2830-2860 gateDataRouteBinding`（`$context.route.*` 需 2.9 + `data.route-binding`）。

**判定**：2.7 envelope + 2.9 页面是上游允许的解耦协商形态；本仓 host 是支持 2.9 页面与两个 v2.9 能力的合法 host，运行时门禁与 claim 一致。**版本协商模型无协议缺口**。

**S2 新发现（回流，P-005）**：① **生产页面级能力协商缺位**——上游 `08-renderer-spec.md` 定义页面级 `MISSING_REQUIRED_CAPABILITY`（页面 requiredCapabilities ⊆ host 支持集），`version-negotiation.cases.json` 亦含页面级协商用例；本仓 `version-negotiate.ts` 实现了该逻辑且 fixtures 全绿，但 **`RenderPage` 生产路径未调用**（`render.tsx:3026` 起按节点走特性级门禁：form-controls.types.ts / permissions.ts / `gateDataRouteBinding`，无页面级能力集校验）。② **claim/HOST_SUPPORT 能力覆盖不足**——claim `support.capabilities` 与 `boot.ts HOST_SUPPORT` 仅 7 项（app.manifest / app.navigation / host.bootstrap / host.failure-recovery / host.conformance-claim / data.route-binding / form.controls.readonly），而服务页面 `meta.requiredCapabilities` 还要求 permissions.inheritance / actions.row.request / actions.page.trigger / actions.row.navigate / table.sort / form.controls.extended / form.controls.advanced / form.record.load / actions.upload 等；这些能力本仓**均已实现**（排序/权限/扩展控件/recordView/上传等），vendored 对应 suite 亦全绿，但 claim 未声明 → 符合性证据与服务内容不一致。

**处置**：C-004 = no-gap（版本协商模型）+ **implementation-gap 残余项**：①生产页面加载/渲染路径接线页面级版本+能力协商（fail-closed `UNSUPPORTED_PROTOCOL_VERSION` / `MISSING_REQUIRED_CAPABILITY`）；②claim `support.capabilities` 与 `HOST_SUPPORT` 扩展至实际实现能力全集（逐能力 mandatorySuites 已全绿），并加一致性守卫。README 过时表述并入 C-003。S4 整改。

### C-005 · capability 声明与实际使用 — no-gap

**上游事实**：`page.schema.json` `meta.requiredCapabilities` 是「页面依赖的 Renderer 执行能力列表」（可选字段）；协议约束是**单向**的——使用某能力时必须声明（L2 拒绝「用时未声明」），未规定「声明必须使用」。node.schema.json `params` 描述：「使用 `$context.route.query.*/params.*` 时声明 `data.route-binding` 且 protocolVersion >= 2.9」。

**本仓事实（按页核验）**：
- `wallet-entries.json`：声明 `data.route-binding`，**使用** `params: { accountId: "$context.route.params.id" }`（line 115-116）。
- `dictionary-entries.json`：声明 `data.route-binding` + `form.controls.readonly`，**使用** `readOnly: true` ×4（lines 80/86/164/170）+ `params: { dictKey: "$context.route.params.dictKey" }`（line 330-331）。
- `wallet.json`：声明 `form.controls.readonly`，**使用** `readOnly: true`（line 366）。
- `digitaloffer-entitlements.json`：声明 `data.route-binding`（line 12），无任何 `$context.route.*` 绑定。
- `digitaloffer-offers.json`：声明 `form.controls.readonly`（line 15），无任何 `readOnly` 字段。

**判定**：wallet/dictionary 页「声明+使用」合规；digitaloffer 两页「声明未使用」是上游允许的保守声明（不构成违规，也不产生功能偏差——host 支持这两个能力）。**无协议缺口**。留 recommended 卫生项：digitaloffer 两页的未使用能力声明应删除（最诚实）或补上预期用途；因 digitaloffer 是 custom-only 模块（非默认 profile），并入 C-009/S3 custom 边界时一并处置。

### C-006 · 静态 dogfood Manifest — implementation-gap

**上游事实**：无——这是本仓测试夹具管理问题。

**本仓事实**：`apps/web/src/test-fixtures/app-manifest.mvp-dogfood.json`（13 页，协议 2.7，含 demo 页）与 `app-manifest.admin-dogfood.json`（26 页，协议 2.7）为手维护静态夹具，sha 固定于 `upstream-fixtures.test.ts`（`STATIC_MANIFEST_SHA256`），用于 app-manifest/app-navigation 回归。真实 profile 投影为 MVP 6 / Admin 22 / Demo 14（E-002，真实 `ForModulesWithFragments`），与两个 dogfood 均不同；dogfood 同时保存 pageId/route/schemaUrl 副本，存在与模块 manifest 漂移风险。目前无「dogfood ⊆ 模块联合」守卫（E-002 已确认无 HTTP 快照保留）。

**判定**：夹具合法（回归价值），但命名暗示「运行时权威」与实际角色不符，且无防漂移守卫。属测试证据卫生 implementation-gap。修正 = 新增守卫测试：从 `apps/api/modules/*/manifest/fragment.json` 推导页面联合，断言 dogfood 的每个 pageId/route/schemaUrl 与该联合一致（或标注为「静态回归夹具，非运行时权威」并加校验）。S4 整改。

### C-007 · 页面链路测试覆盖 — excluded（→S5 验收分母）

**上游事实**：无对应上游缺口。

**本仓事实**：D-VAL 结构验证经 **A-002 F-002 纠正并修复**：`all-module-schemas-dval.test.ts` 原 walker 只扫 `modules/<top>/schema` 一层（漏 `channel/telegram/schema` 2 页 + `dev/examples/schema` 8 页 = 实扫 25/35）；本轮已改为递归遍历全部 `**/schema/*.json` 并实测 **35/35 绿（38 tests）**。渲染级 `representative-pages.test.tsx` 覆盖 10/35（data-table、search-form-table、form-controls、form-with-reactions、users、users-invites、roles、settings、activity、data-permission），直接读取真实 API schema 文档 + 注入 fixture fetcher（非真实 HTTP/handler）；v2.9 页另有定向测试（`schema-dictionary-entries.test.tsx` 驱动真实 dictionary-entries 文档含 route-binding、`wallet-navigate.test.tsx`、`App.integration.test.tsx` 内页深链、`row-action-bindings.test.ts`）。

**判定**：非产品缺口，而是 S5 验收分母契约定义项。S2 在 D-002 §3 固化验收矩阵（每页 validator + `loadPageDocument → RenderPage` + 代表性 handler 触达 + 真实 HTTP Manifest 快照），S5 按 I-006 执行。**排除（本波不补产品覆盖；责任=GOAL-041 S5；触发=进入 S5）**。

### C-008 · profile / optional 运行时矩阵 — excluded（→S5 验收矩阵）

**上游事实**：无对应上游缺口。

**本仓事实**：E-002 已用真实 `ForModulesWithFragments` 按 profile 过滤 fragment 得到 MVP 6 / Admin 22 / Demo 14；Telegram / Digital Offer 不在默认 profile，仅显式 custom 模块配置进入（profile.go：channel.telegram / biz.digital-offer 均注明「not enabled in defaults」）；未保留实际 HTTP Manifest 快照。

**判定**：非产品缺口，是验收矩阵设计项。S2 在 D-002 §3 固化矩阵（默认三 profile × 显式 custom 组合；每个组合留 HTTP Manifest 快照），S5 执行。**排除（本波不实施；责任=GOAL-041 S5；触发=进入 S5）**。

### C-009 · custom component 扩展面 — custom-extension-candidate

**上游事实**：`08-renderer-spec.md §1.1` 明确定义 **Host Extension 模型**：核心注册表只含 component-registry 声明的类型；宿主可用私有 API 安装扩展；扩展不属于页面协议核心身份、不得跨 Renderer 以核心身份互操作；未安装扩展在标准入口以 `UNKNOWN_COMPONENT_TYPE` 拒绝。上游未规定扩展 key 的 namespace 格式（无强制）。

**本仓事实**：15 个注册键（12 页面 body custom、1 action content custom、2 afterComponent），当前 schema 使用键全部有注册（`custom-components.schema.test.ts` W25 防复发守卫）；`runtime-schema-validate.ts` 以本地扩展 `component` 属性叠加在 vendored node schema 之上（上游 schema 字节不变，注释明确「NOT an upstream protocol change」）；`docs/architecture/module-contribution-playbook.md §6.2/6.3` 文档化注册要求与校验守卫，但**无统一 namespace 约定**（键为 `wallet-ensure`、`mail-admin-tab` 等扁平名）；各扩展的 capability/schema/validator/failure/fixtures 边界未逐键成文。

**判定**：机制形态符合上游 Host Extension 模型（显式注册、校验时扩展、防注册缺失守卫），但按 D-001 §4 custom 门禁，必须逐键固定 namespace、能力声明、schema/validator、失败语义、兼容策略、fixtures 与退出/迁移触发，并**经用户 P-004 书面裁决**（S3）。**custom-extension-candidate**。

### C-010 · 未注册 custom 失败语义 — implementation-gap

**上游事实**：`01-node-protocol.md`「未知 type 时应渲染明显『未识别组件』占位，而不是静默失败」；`08-renderer-spec.md §1.1`「未安装扩展 type 在标准入口以 `UNKNOWN_COMPONENT_TYPE` 拒绝，不得静默降级」、§1.3「未知 type：渲染明显占位 + `console.error`（含 type 与 node id）」；`conformance/fixtures/runtime-defaults` 的 `uninstalled-extension-is-unknown-component` 用例期望 `{ok:false, code:"UNKNOWN_COMPONENT_TYPE"}`（本仓 vendor 同 digest）。

**本仓事实**：未知标准 node → `RENDER_UNKNOWN_NODE_TYPE` fail closed（`render.types.ts:397-407 parseRenderNode`，`render.test.ts:82` 断言）；未知 **custom component key** → `render.tsx:2915-2927` 渲染 `<p>unknown custom component: {node.component}</p>` 内联文字，无 `console.error`、无明显占位样式、不阻断页面其余渲染。conformance 适配器 `runtime-defaults.ts` 已按上游实现 `UNKNOWN_COMPONENT_TYPE`（fixture 全绿）。

**判定**：产品渲染器对未知 custom 的呈现弱于上游契约（§1.1/§1.3 要求明显占位 + console.error，或标准入口 fail-closed `UNKNOWN_COMPONENT_TYPE`）；当前「有文字但不清除」属合法降级与契约要求之间的偏差。**implementation-gap（低严重度，recommended 级）**。修正 = 渲染器对未知 custom 输出明显占位（含 type/component 与 node id）+ `console.error`，或与 parseRenderNode 一致 fail-closed；同时保留现有「有注册校验测试」防线。S4 整改。

### C-011 · custom action allowlist — no-gap

**上游事实**：`07-actions-contract.md §6`：`type: custom + handler`，handler 仅允许前端白名单预注册函数名（不接受任意代码/表达式）；「前端必须维护白名单，handler 不在白名单中时 Renderer 应拒绝执行并报错」。

**本仓事实**：5 个 schema custom handler（`export.users` / `export.roles` / `library.download` / `library.preview` / `library.copyLink`）全部在 `render.tsx:332-343 CUSTOM_HANDLER_URLS` 白名单；非白名单 → `CUSTOM_HANDLER_NOT_FOUND` fail closed（line 353）；`{id}` 槽位缺 row id → `CUSTOM_HANDLER_MISSING_ROW_ID` fail closed（line 359-362）；下载/预览/复制有安全边界（blob + sandboxed iframe + opener=null + 文件名 scrub + origin-absolute copy）；`download-behavior.test.tsx:168` 断言未知 handler fail closed。

**判定**：与上游 custom action 契约一致，白名单 + fail-closed 拒绝已实现。**无协议缺口**。namespace/边界细节并入 C-009 custom 记账。

### C-012 · 内页与导航边界 — no-gap

**上游事实**：`action.schema.json` NavigateAction `url` 为应用路由根下相对路径；manifest 页面 route 模板允许参数化；`09-app-manifest.md` 页面可达性由 manifest pages + navigation 决定（Host-owned 入口合法）。

**本仓事实**：内页全部登记在 manifest pages：`dictionary-entries`（/dictionary-entries/{dictKey}）、`task-runs`、`wallet-entries`（/wallet-entries/{id}）、`telegram-operator`、`users-invites`，分别经 `openEntries`/`openRuns`/`openEntries`/`openTelegramOperator`/`openInvites` navigate action 到达（data-dictionary/scheduled-tasks/wallet/telegram-settings/users 页面）；`notifications` 页经 Host 铃铛（`notification-bell.tsx` `onViewAll → /notifications`、`/notifications?open=<id>`）到达，manifest `navigation.user` 为空的空槽声明（协议允许）；`App.tsx:204-207` breadcrumb 父页映射就位；E-002 核实重复 pageId 0 / 缺失 schema 0 / orphan schema 0。`App.integration.test.tsx` 覆盖内页深链可达性。

**判定**：全部为合法内页 route / Host-owned 入口，均登记可达，无 Manifest/导航断链。**无协议缺口**。

### C-013 · API-only / Host-only 表面 — explicitly-out

**上游事实**：无——这是归属判定（本仓层责任）。

**本仓事实**：`profile.go` 模块描述符确认 `admin.data-transfer`（路由 GET /api/export/{resource} 等 + 权限 data.export/import，无 Pages）、`admin.mfa`（路由 + 权限 users.mfa-reset，无 Pages/Navigation/Fragments；provider.go 注释「No page/navigation/fragment」）、`admin.login-captcha`（路由 + 权限，无 Pages）——均为**无页面模块**；其 UI 承载为 Web custom 组件（mfa-manager / password-policy-tab / captcha 开关等）与 action（exportUsers 等），已分别计入 C-009/C-011。login/session/branding/shell/failure 属 Host 层（App.tsx / AuthGate / boot.ts / failure.ts / HostFailureScreen），GOAL-004 已按 Host/App 候选处置。

**判定**：**explicitly-out**——API-only 模块与 Host 层责任面，记录归属与边界，不伪装成页面协议能力；其 Web 呈现通过 custom（C-009）/ action（C-011）面完整记账。

### C-014 · 历史 GOAL-004 证据边界 — excluded（方法论规则）

**上游事实**：GOAL-004 曾按上游 ADR-0034 D10/D6 对 95/95 Host/App 候选做 adopt-now / reserve-extension / explicitly-out 处置（GOAL-004 D-002）。

**本仓事实**：当前模块集、35 页、15 custom、v2.9 增量（data.route-binding / form.controls.readonly）已远超 v2.8 时的 95 候选面；GOAL-004 的处置语义（adopt-now=按上游修本仓等）可作为边界规则复用，但其结论不能作为 W29 现行符合性证明。

**判定**：**excluded（方法论规则落盘，见 D-002 §5）**——本波分类矩阵为现行权威；GOAL-004 仅提供历史处置语义；责任=本波 S5/S6 证据与后续波次；复核触发=进入 S5 或生产 Manifest 新增页面/控件时。

---

## S2 结论边界

- 已完成：14 候选全部给出带证据的唯一处置类别；upstream-protocol-gap = 0（I-004 以「不适用」证据收口，不建空报告）；custom 候选 = C-009（S3 触发 P-004）；S5 验收契约（C-007/C-008）与历史边界规则（C-014）在 D-002 固化；**A-002 independent 意见已合并**（F-001 保持 S4 必做，F-002/F-003 walker 守卫已修复并实测绿，F-004/F-005 台账已纠正，见 A-003 响应）。
- 未完成：F-001 生产页面级能力门禁 + claim/HOST_SUPPORT 覆盖（S4 必做，需 P-004 裁决 S2 关门口径）；S4 其余整改（C-001/002/003/006/010）；S3 custom 裁决（C-009）；S5 运行时验证（I-006）与 go 影响判定（I-007）；S6 cross 关门。
