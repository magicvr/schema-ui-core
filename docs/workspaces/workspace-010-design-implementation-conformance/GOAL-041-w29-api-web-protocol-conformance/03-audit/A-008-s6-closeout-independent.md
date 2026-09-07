---
id: GOAL-041-w29-api-web-protocol-conformance
doc: audit-entry
record_id: A-008
source: independent
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# A-008 · S6 关门 close-out · independent 审计

## A-008 · S6 关门复审（2026-09-06）
- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`）
- **类型 / scope**：close-out（GOAL-041 全目标 S1～S6）；S2 分类与 D-002 冻结如实性；S3 C-009 custom 边界（D-001 §4）；S4 实现整改（页面级门禁 / 19 能力 / digitaloffer / dogfood / C-010）；S5 运行时验证（35/35、5 组合快照、失败路径、go 影响）；A-001～A-007 required 闭合；I-001～I-009
- **verdict**：conditional

### 范围与区间

工作区：`workspace-010-design-implementation-conformance`（`root_goal` = `GOAL-001-design-implementation-conformance`；canonical = `docs/workspaces/workspace-010-design-implementation-conformance/`；`shared_materials_catalog: none`；`primary_plan` = `VP-010-design-implementation-conformance`）。GOAL-041 `parent` = `GOAL-001-design-implementation-conformance`，在 canonical 范围内。无共享资料引用被当成事实。本条只审本工作区本目标；不改 `00-meta` status/progress、goal-tree 或方案正文。

本轮独立核验：对照代码、附件、快照与定向复跑，不把 self A-007 的绿灯当作完成证明。定向复跑（本机 win32，`apps/web`）：`claim-artifact` / `load-page` / `stage3-fixtures` / `capability-declaration.guard` / `dogfood-manifest.guard` = **5 files / 336 tests PASS**。未复跑全量 vitest 1307、`denominator-render.test.tsx` 36 tests、Go `./...`；S5 分母与快照以源码 + 落盘工件 + walker 实测计数核对。

### 成果（有证据）

1. **工作区绑定合格**：workspace id / Root / canonical 与 GOAL-041 `parent` 一致；共享资料目录 `none`。
2. **S2 分类与 D-002 冻结如实（主类）**：C-001～C-014 均有唯一处置类别；统计 implementation-gap ×6 / custom ×1 / explicitly-out ×1 / excluded ×2 / no-gap ×3 / **upstream-protocol-gap = 0** 与 D-002、I-004「不适用 / 不建空报告」一致。C-004 汇总行已改标 implementation-gap（A-002 F-004 响应）；详细节标题仍写 no-gap，见本条 F-002。
3. **S3 C-009 边界符合 D-001 §4**：用户 P-004 书面裁决在 D-003（路径 = 本仓合法 custom；namespace = 保留 15 键 + 登记 + 新键规范；C-005 子项 = 删除）。`attachments/custom-extension-boundary.md` 覆盖 namespace（§3.2）、capability（§3.3 不新增）、schema/validator（§3.4）、failure（§3.5）、compatibility（§3.6）、fixtures（§3.7）、退出/迁移（§3.8）。15 键均有 `registerCustomComponent`（含 `telegram-admin-tab.tsx:1089`）。
4. **S4 实现整改真实（产品路径）**：
   - **F-001 页面级门禁**：`load-page.ts` 在 D-VAL 之后 fail-closed `UNSUPPORTED_PROTOCOL_VERSION` / `MISSING_REQUIRED_CAPABILITY`；`App.tsx:560` 生产路径调用 `loadPageDocument`。`load-page.test.ts` 负例（2.10 版本 / `host.virtual-reality`）本轮复跑通过。
   - **19 能力全集**：`host-support.ts` `HOST_SUPPORTED_CAPABILITIES` 19 项 = 上游 `docs/schemas/capability-registry.json` 全集；`boot.ts:71` 复用该常量；claim `support.capabilities` 同序 19 项、`conformance.suites` 12 个、`protocolArtifact.artifactVersion` = `2.9.0`；`generate-claim.mjs` `ARTIFACT_VERSION = "2.9.0"`；`conformance-local-report.json` `pinnedUpstream.artifactVersion` = `2.9.0`。`claim-artifact.test.ts` C0/C1 本轮 5/5 PASS。
   - **C-001 provenance**：`provenance-v2.9.json` fixture-suite 路径 = `conformance/schemas/fixture-suite.schema.json`；stage3 守卫断言该路径；本轮 stage3 278 tests PASS。
   - **C-003**：`provenance.json` note 标注 R3 v2.7.0 兼容基线 / 现行权威 = v2.9；`provenance-v2.8.json` 保留为历史；README 版本协商节已改。
   - **C-005 digitaloffer**：`digitaloffer-entitlements.json` / `digitaloffer-offers.json` 已无 `data.route-binding` / `form.controls.readonly`。capability-declaration 守卫本轮 **35 tests PASS**。
   - **C-006 dogfood**：守卫规范化路径分隔符（跨平台）；本轮 3/3 PASS。
   - **C-010**：`render.tsx` 未知 custom → `role="alert"` 红框占位 + `console.error`（含 component 键与 node id）；`render.test.tsx` 有对应用例。符合用户「明显占位 + console.error」裁决。
5. **S5 分母与快照（Windows 计数成立；Linux CI 不成立，见 F-001）**：本机 walker 用 `/\\schema\\/` 对原生 win32 路径收集 **35/35**，pageId 与 `S5-coverage-matrix.md` 全分母表一致。落盘快照页数：mvp 6 / admin 22 / demo 14 / admin+digitaloffer 25 / admin+telegram 24；mvp 快照 `protocolVersion` = `2.7`、`requiredCapabilities` 含 `app.manifest` + `app.navigation`。`s5_manifest_snapshot_test.go` 冻结断言与上述页数一致（Go 测试本身不依赖 Windows 路径正则）。
6. **失败路径有测试**：load-page 负例、C-010 占位、`CUSTOM_HANDLER_NOT_FOUND`（`render.tsx:353` + `download-behavior.test.tsx`）、D-VAL `PAGE_SCHEMA_INVALID`、manifest 错误码类型在 `app-manifest.ts`。覆盖矩阵失败路径表与代码对得上。
7. **I-007 go 无影响**：W29 产品提交（`9d00d16a` S4 / `cbfb8508` S5）未改 `apps/api/kernel/profile.go`（`profileDefaults` / `BuiltinModules`）；变更面为 claim/provenance、host-support/load-page/renderer、digitaloffer 元数据声明、守卫与快照测试。判定「无影响不暂挂」成立。
8. **历史 required 闭合链可核对**：A-001/A-002 F-001 先 accepted-residual（A-003 用户书面 + S4 复审触发）后以 S4 完成证据 **fixed**（接线 + 19 能力）；A-002 F-002～F-005 在 Windows 证据下为 fixed。本条不推翻该链的产品事实；F-002 的递归 walker 在 Linux CI 上不可复跑，作为**新** required 提出（见下），不是宣布原闭合声明作伪。

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 分母与候选目录 | 达成 | E-002；`attachments/S1-*`；本轮 35 个 schema pageId 与目录一致 |
| S2 分类唯一、upstream gap = 0 | 达成（台账卫生见 F-002） | D-002；`S2-candidate-classification.md` |
| S3 custom 边界 + 用户书面裁决 | 达成 | D-003；`custom-extension-boundary.md` §3.2–§3.8 |
| S4 按已冻结契约整改 + 防复发 | 产品路径达成 | load-page / host-support / claim / provenance / C-010 / digitaloffer / dogfood |
| S5 35/35 + 5 组合快照 + 失败路径 + 可复跑 | **部分** | Windows 35/35 与快照成立；项目 CI `ubuntu-latest` 上 schema walker 收集 0 页（F-001） |
| S6 required findings / 信息项合法闭合 | **未满足** | 本条 F-001 required 开放；I-006「可复跑」在 CI 上不成立 |
| I-007 go 影响 | 达成 | 无 Profile/模块矩阵/Manifest 装配改动 |
| 历史 A-001～A-007 required | 产品闭合成立 | F-001 门禁已接线；F-002～F-005 台账/walker 在 win32 可核对 |

### Findings

#### F-001 · S5/S4 schema walker 使用 Windows 反斜杠正则，Linux CI 无法收集 35 页分母（high · required）
- **严重度**：high
- **建议**：required（闭合前不得将 I-006「可复跑」与 S6 关门视为满足）
- **关联**：I-006；C-007 / D-002 §3；A-002 F-002 的递归 walker 修复
- **描述**：四个分母/守卫测试用 `/\\schema\\/` 匹配绝对路径中的 `\schema\`：
  - `apps/web/src/protocol/all-module-schemas-dval.test.ts:41`
  - `apps/web/src/renderer/denominator-render.test.tsx:53`
  - `apps/web/src/protocol/capability-declaration.guard.test.ts:36`
  - `apps/web/src/renderer/custom-components.schema.test.ts:54`
  本机 `process.platform = win32`：该正则收集 **35** 份 schema（pageId 全集与覆盖矩阵一致）；把同一路径规范化为 POSIX `/` 后再套该正则收集 **0**。项目 CI `.github/workflows/r6-basic-matrix.yml` web job 为 `runs-on: ubuntu-latest`（Node 路径为 `/schema/`）。在 Linux 上：`denominator-render` 的 `expect(REFS.length).toBe(35)` 必红；`capability-declaration` 的 `>= 30` 必红；D-VAL `docs.length > 10` 必红；custom 守卫 `refs.length > 0` 必红。C-006 dogfood 守卫已做 `replace(/\\/g, "/")` 跨平台处理，对照说明作者知道分隔符问题，但分母 walker 未采用。
  产品门禁与 35 页文档本身在 Windows 上可核对，**不是**「页面不存在」；不成立的是 I-006 / S5 / S6 的「回归可复跑」在项目 CI 上的主张。A-002 F-002 把一层 walker 改成递归时引入该正则：旧一层 `join(module, "schema")` 在 Linux 可工作，新递归正则在 Linux 收集 0——F-002 的「35/35 绿」仅在 win32 上可重复。
- **证据**：上述四处 file:line；本轮 node 实测 nativeWinCount=35 / posixNormalizedWinRegexCount=0 / posixSlashCount=35；`.github/workflows/r6-basic-matrix.yml:11-24`；`dogfood-manifest.guard.test.ts:59-62`（正确规范化对照）；`denominator-render.test.tsx:127-128`。
- **状态**：open。
- **建议修复**：walker 先 `abs.replace(/\\/g, "/")` 再匹配 `/\/schema\//`（或按 `path.sep` / 目录名 `schema` 判断）；补一条「规范化后仍为 35」的断言，使 ubuntu CI 与 win32 同分母。修复后须在 Linux 语义下复跑 D-VAL / denominator-render / capability-declaration / custom-components 四处。

#### F-002 · 关门台账若干字段仍停在 S2/S3 中途表述（low · recommended）
- **严重度**：low
- **建议**：recommended（不阻断 F-001 实质；编排器响应时一并改）
- **描述**：(a) `01-decision.md` I-008 仍写「open（S2 腿：A-001 self conditional；grok build independent 待合并）」——S2 腿实际已由 A-001/A-002/A-003 完成。(b) `attachments/S2-candidate-classification.md` 结论「未完成」仍把 F-001 列为 S4 待做，且 C-004 详细节标题仍为「no-gap」，与汇总行 / D-002 主类 implementation-gap 不一致。(c) 本文件更新前 `03-audit.md` 信息就绪仍写「custom 用户裁决 \| 未到期」。这些不否定已完成的产品工作，但会让 S6 读者误判门禁仍停在 S2。
- **证据**：`01-decision.md` I-008 行；`S2-candidate-classification.md` C-004 标题与文末「未完成」；本审计开始时的 `03-audit.md` 信息就绪表。
- **状态**：open。

### 必改项汇总

- **required**：F-001（四处 schema walker 改为跨平台路径匹配，并在 Linux 语义下复跑 35/35 分母守卫；闭合前不得 `status: done`）。
- **recommended**：F-002（I-008 / S2 附件 / C-004 标题与信息就绪表与事实对齐）。

### 与既有意见的异同（self A-007）

| 点 | self A-007 | independent A-008 |
|---|---|---|
| verdict | conditional→pass（待本条 + 用户确认） | **conditional**（不能升 pass） |
| S2 分类 / D-002 | 达成 | **同意** |
| S3 custom 边界 / D-001 §4 | 达成 | **同意**（要素与用户书面裁决齐备） |
| S4 F-001 接线 + 19 能力 + digitaloffer + dogfood + C-010 | 达成 | **同意**（代码与定向复跑 336 tests 核对） |
| S5 35/35 + 5 快照 | 达成；回归可复跑 | **快照与 win32 35 页同意**；**不同意「可复跑」无条件成立**——CI Linux 分母 walker 为 0 |
| 历史 F-001～F-005 闭合 | 全闭合 | 产品闭合 **同意**；F-002 walker 修复在 CI 上不成立 → **新 required F-001** |
| I-007 go | 无影响不暂挂 | **同意** |
| I-001～I-005 / I-009 | verified / 不适用 / open non-blocking | **同意**（I-009 保持 open 合法） |
| I-006 / I-008 | verified / S6 腿待本条 | I-006 因 F-001 **不能维持无条件 verified**；I-008 S6 腿 = A-007 + 本条，required 未闭合 |
| 新 findings | 0 | required F-001 + recommended F-002 |

无与 self 相反的产品事实结论（门禁未接线、19 能力不足、快照页数造假等均未发现）。分歧限于 **S5/S6「可复跑」主张的环境范围**。无 P-004 冲突项。

### 信息项核对（P-005）

| ID | 台账状态 | 本条判定 |
|----|----------|----------|
| I-001 / I-002 | verified | 维持 |
| I-003 | verified | 维持（分类如实；C-004 标题卫生见 F-002） |
| I-004 | 不适用 | 维持（upstream gap = 0） |
| I-005 | verified | 维持 |
| I-006 | verified | **降为有条件**：win32 证据成立；项目 CI 可复跑不成立（F-001）。编排器应在 F-001 fixed 前不要把 I-006 当关门放行依据 |
| I-007 | verified（无影响不暂挂） | 维持 |
| I-008 | open（S6 腿） | S6 腿 = A-007 + 本条已落盘；因 F-001 required 开放，I-008 **尚未**可标 verified |
| I-009 | open non-blocking | 维持；S6 已复核，不阻断（责任人/触发仍有效） |

### 结论 + 建议给编排器/用户的下一步

S1～S4 产品整改与 S3 custom 裁决可独立核对；S5 快照矩阵与 win32 35 页分母成立；历史 A-001～A-007 required 的产品闭合链成立。S6 **不能无条件关门**：I-006「可复跑」在项目 Linux CI 上被四处 Windows 路径正则证伪。

建议 `/govern`：

1. 响应本条 F-001：跨平台修复 walker → 复跑四处分母守卫（最好在 Linux 或用 POSIX 规范化断言）→ 按 `fixed` 闭合。
2. 顺手做 F-002 台账卫生（I-008 表述、S2 附件「未完成」、C-004 详细标题）。
3. F-001 合法闭合前 **不得** 将 GOAL-041 标 `done`。闭合后需用户书面确认（meta S6 要求）再改 status/progress/goal-tree。
4. I-007 保持不暂挂。

### 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
