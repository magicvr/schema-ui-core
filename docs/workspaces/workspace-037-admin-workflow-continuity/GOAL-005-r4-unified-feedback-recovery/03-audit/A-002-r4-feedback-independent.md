---
id: A-002-r4-feedback-independent
doc: audit-opinion
status: recorded
source: independent
verdict: conditional
scope: R4 C1-C3 implementation, regression evidence, information gates, and C4 close-out readiness
audit_type: execution-facts
goal_id: GOAL-005-r4-unified-feedback-recovery
auditor: grok-build (grok-4.6 · reasoning high)
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-005-r4-unified-feedback-recovery
version: 1.0.0
---

# A-002 · R4 统一反馈与恢复独立审计

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：execution-facts；GOAL-005 R4 C1～C3 实现、回归证据、R4-I-001～005 信息门禁与 C4 关门准备（不含 Root/VP 关门，不含 R5）
- **verdict**：conditional

## 范围与区间

工作区：`workspace-037-admin-workflow-continuity`（`root_goal: GOAL-001-admin-workflow-continuity`，`canonical_scope` 匹配，`shared_materials_catalog: none`，`primary_plan: VP-037-admin-workflow-continuity`）。未读取其他工作区。无固定共享资料引用。

本意见只核对 GOAL-005 的 C1～C3 实现、D-001 / R1 D-005 合同、R4-I-001～005、A-001 self 主张，以及 C4 是否具备关门条件。不修改 `status` / `progress` / 方案正文 / goal-tree。C4（本意见的编排响应 + Git checkpoint）、Root R4 投影与 R5 均尚未完成。

## 成果（有证据）

### C1 · 分类合同

`apps/web/src/renderer/feedback-policy.ts` 把 status/code/AbortError/`Failed to fetch` 映射到 validation / authentication / forbidden / not-found / conflict / rate-limited / maintenance / unavailable / offline / timeout / unknown；catalog key 优先；`retry` 仅在调用方显式传入 **且** `retryableRead === true` 时附加。`en-US` / `zh-CN` 均有对应 `feedback.*` 文案。Host 终态仍由 `HostFailureScreen` + `host/failure.ts` 处理，普通 resource 表面未接管 `main.tsx` 的 Host 恢复动作。

### C2 · 共享反馈表面

`FeedbackNoticeView`：成功 `role=status` 且 4s 自动消失；错误 `role=alert` 持久；toast 可 dismiss；retry/dismiss 为 `type=button`；ref 防止 in-flight 双击。SchemaTable / statCard / chart / recordSource 走共享表面；FormInner 写失败不传 retry，并 `showDiagnosticCode`。`safeParams` 只保留 string/number。

### C3 · 回归（本独立审计复跑）

在 `apps/web` 使用项目 vitest 3.2.7 与本地 `tsc`：

| 集合 | 结果 |
|------|------|
| 定向 8 文件（policy / FeedbackNoticeView / DataTable / SchemaTable / render / representative-pages / r3-dirty / error-localization） | **8 files / 126 passed** |
| 前端全量 | **110 files / 1401 passed** |
| `.\node_modules\.bin\tsc -p tsconfig.app.json --noEmit` | exit 0 |

已核对：列表 500 显示 unavailable catalog 且 `[data-table-retry]` 再拉一次成功；statCard 503 点一次 Retry 恢复值、`calls === 2`；representative `/data-table` 500 使用 R4 catalog 文案；FormInner 400/transport throw 保留 dirty 且无写 retry；Host fixtures（`upstream-host-fixtures.test.ts` 99）与 page schema/manifest 既有恢复测试仍在全量中通过。

E-003 所记「定向 8 文件 / 136 测试」与本次同一 8 文件集合的 **126** 不一致；全量 110/1401 可复核。该数字差不推翻已通过的实现路径，但不得把 136 当作本次独立复跑结果。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1：错误分类、catalog、诊断 code/correlation、恢复策略可核对 | 部分 | policy 纯函数 + 单测成立；**写入/action/recordSource 的 transport throw 未走该映射**（F-001） |
| C2：共享表面 status/alert、dismiss/retry、键盘路径、去重 | 达成 | `feedback.tsx` + `feedback.test.tsx`；写路径不附加 retry |
| C3：读重试、写失败保留、maintenance/unavailable/offline/timeout 与 Host 边界的跨页面自动化 | 部分 | 读重试与写失败保留有页面/组件证据；unavailable 有跨页面证据；maintenance 无自动化；offline/timeout 仅 policy 单测，且生产写入路径会丢失 timeout（F-001/F-002） |
| C4：self + independent + required 响应 + Git checkpoint | 未完成 | A-001 已落盘；本意见待 `/govern` 响应；尚无 R4 checkpoint |
| R4-I-001 required · 最晚 C1 | 证据不完整 | 映射函数存在，但 C1 主张的 timeout/offline 分裂在写入路径未落地（F-001） |
| R4-I-002 required · 最晚 C2 | 可核对 | 仅读取调用方传 retry；FormInner/action `errorFeedback()` 不传 retry |
| R4-I-003 required · 最晚 C2 | 可核对 | role/dismiss/retry/双击 guard 有组件测试 |
| R4-I-004 required · 最晚 C3 | 证据不完整 | 边界在组件分层上成立；跨页面自动化未覆盖 maintenance，offline/timeout 被 F-001 削弱 |
| R4-I-005 non-blocking | deferred | 仅安全 code/correlation；不阻断 C4 |

## Findings

### F-001 · 写入/action/recordSource 将 transport throw 压成 `REQUEST_FAILED`，C1 的 timeout/offline 分裂未落地

- 严重度：med
- 建议：required
- 状态：open
- 影响门禁：C1 分类合同、R4-I-001、R4-I-004、C3 验收、C4 关门
- 描述：`classifyFeedbackFailure` / `failureLikeOf` 把 `AbortError` 映射为 `REQUEST_TIMEOUT` → `feedback.timeout`，把浏览器 `TypeError` / `Failed to fetch` 映射为 `REQUEST_FAILED` → `feedback.offline`。policy 单测覆盖了把 **原始 error** 交给 `feedbackFromError` 的路径。生产里列表/statCard/chart 的 `catch` 确实传入原始 error。但以下调用方先把任意 throw 写成 `{ code: "REQUEST_FAILED" }`，再进入 `feedbackFromError`，因此 **客户端超时（`withTimeout` 的 AbortError，这是现行默认超时路径）会被显示成离线文案**：
  - `runRequest` fetch catch（表单/行/页 request action）
  - `runBatchRequest` fetch catch
  - custom download/preview fetch catch
  - `useRecordSourcePrefill` catch（幂等读取，仍会带 retry，但文案变成 offline）
  - FormInner 的防御性 catch（同一扁平化）
- HTTP 504 仍会按 status 分成 timeout，因为那是 `response.ok === false` 而不是 throw。现行 `authFetch`/`withTimeout` 的 30s 中止是 throw，不会变成 504。
- 写失败仍然不自动 retry、dirty 仍保留，所以这不是 retry-loop 漏洞；它使 C1 已冻结的 timeout vs offline 用户文案在默认超时路径上名不副实，也使 R4-I-001/004 不能按「已 verified」无条件关闭。
- 证据：`apps/web/src/renderer/feedback-policy.ts` 第 45–51、85–89 行；`feedback-policy.test.ts` AbortError 用例；`apps/web/src/renderer/render.tsx` 约 413–415、441–444、606–609、749–750、1717–1722、1998–2006 行；`apps/web/src/lib/fetch-timeout.ts`；`apps/web/src/account/auth-client.ts` 使用 `timeoutFetch`。
- 闭合要求（`fixed`）：transport catch 保留 `AbortError` → `REQUEST_TIMEOUT`（或把原始 error 交给 `feedbackFromError`），仅将网络 `TypeError`/`Failed to fetch` 标为 `REQUEST_FAILED`。至少补：① 表单/action AbortError 显示 `feedback.timeout` 且 **无** retry；② recordSource AbortError 显示 timeout **且** 显式 retry 一次。在 `fixed` / `accepted-residual` / `user-overruled` 之一留痕前，不得关闭 C4、不得将 GOAL-005 标为 `done`、不得投影 Root R4 完成。

### F-002 · C3/R4-I-004 的跨页面自动化未覆盖 maintenance，offline/timeout 也不是页面级证据

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：C3 检查点与 A-001 将 maintenance/unavailable/offline/timeout 与 Host 边界写成已有跨页面自动化。独立对照：unavailable 有 SchemaTable 500、statCard 503、representative list 500；`SERVICE_MAINTENANCE` / `feedback.maintenance` **无单测也无页面测**；offline/timeout 只有 policy 单测，且被 F-001 从写入/recordSource 生产路径上拆掉。Host 边界是组件分层事实（`HostFailureScreen` 未被 resource toast 替代，fixtures 仍绿），但没有「resource 503 不挂 Host 终态 / Host maintenance 不走 FeedbackNoticeView」的对照断言。
- 证据：`feedback-policy.test.ts`（无 maintenance 用例）；`schema-table.test.tsx` / `render.test.tsx` / `representative-pages.integration.test.tsx`；`host/failure.ts` 的 bootstrap `MAINTENANCE` vs resource `SERVICE_MAINTENANCE` 是不同触发源。
- 建议：补 maintenance 分类 + 至少一条资源读取 maintenance 文案/retry；若 F-001 修复，再补页面或组件级 timeout/offline。不要把全量 1401 绿扩写成这四类都已跨页面覆盖。不阻断在 F-001 闭合后的 C4，除非 `/govern` 把本条升级为 required。

### F-003 · `01-decision.md` 信息表仍为 collecting，与 `00-meta` / `03-audit` 的 verified 冲突

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`00-meta.md` 与 `03-audit.md` 将 R4-I-001～004 标为 verified；决策索引 `01-decision.md` 仍全部 `collecting`，当前投影仍写「C1～C3 门禁关闭前不把跨页面一致性写成已完成事实」。这是台账漂移，不是「信息尚未发现」。F-001 另说明 I-001/I-004 的 verified 本身还不能无条件成立。
- 证据：`GOAL-005.../01-decision.md` 信息表 vs `00-meta.md` 信息表。
- 建议：`/govern` 响应时把决策索引改到与当时真实门禁一致（F-001 闭合前不要把 I-001/I-004 写成无保留 verified）。本意见不改决策正文。

### F-004 · 403/401 列表不出现 retry、chart retry 仅有结构证据

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：policy 对 401/403/409/422 的 `retryableRead: false` 有单测；SchemaTable 把 retry 传入 `feedbackFromError`，非 retryable 时 notice 不带 retry，DataTable 使用该 notice，故 403 列表按代码不会画 Retry。缺少「403/401 列表无 `[data-table-retry]`」断言。chart 与 statCard 共用 `useDisplayData`，仅 statCard 有一次 Retry 回归。
- 证据：`schema-table.tsx` `errorFeedback` + `onRetry`；`data-table.tsx` 优先 `errorFeedback`；`render.test.tsx` 仅 statCard retry。
- 建议：补 403 列表无 retry、可选 chart 一次 retry。不阻断 C4。

## 必改项汇总

| ID | 严重度 | 闭合前阻断 |
|----|--------|------------|
| **F-001** | med · required | C4 关门、GOAL-005 `done`、Root R4 投影 |

F-002～F-004 为 recommended，默认不单独阻断 C4。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 独立判断 |
|----|------|----------|----------|
| R4-I-001 | required | C1 | 不能维持无保留 verified：timeout/offline 映射在默认超时写入路径未应用（F-001） |
| R4-I-002 | required | C2 | 可维持 verified（读才带 retry，写不带） |
| R4-I-003 | required | C2 | 可维持 verified |
| R4-I-004 | required | C3 | 受 F-001 与 F-002 限制；Host 分层未发现被 Toast 替代 |
| R4-I-005 | deferred non-blocking | R5/支持需求 | 仍成立；owner=`/vision` |
| I-037-004（Root/VP） | required · R1 已冻结 | R4 实施/验收 | R1 D-005 冻结仍有效；R4 实现证据需在 F-001 闭合后才能作为 Root 投影依据。Root `00-meta` 仍写「R4 另留跨页面回归证据」，在 C4 前属预期，不在本目标改 |

无到期未接受残余的 required 信息项被本意见改写成「已验证成功事实」。`shared_materials_catalog: none`，无资料引用被当作关闭证据。

## 与既有意见的异同

A-001 self `pass`、无 required。本意见同意：共享表面、读 retry、写不自动 retry、dirty 保留、catalog-first、Host 组件未被替换、全量 110/1401 与 `tsc` 可复跑。不同意把 R4-I-001/004 与 C3 检查点写成已无缺口完成，因为默认超时路径的分类被 `REQUEST_FAILED` 扁平化。

这不是同范围 pass/fail 冲突：A-001 未声称「AbortError 在 runRequest catch 中仍保持 timeout」。无 residual/overrule 请求。若要以不修 F-001 而关门，才需要用户按 P-004 选择 `accepted-residual` 或 `user-overruled`。

## 结论 + 建议给编排器/用户的下一步

独立审计 **conditional**。C2 与读路径 C3 主路径成立；C4 **可以进入响应流程**，但 **F-001 合法闭合前不得** 将 C4 / GOAL-005 / Root R4 标为完成，也不得把 Git checkpoint 当作放行证据。

建议 `/govern`：按 `fixed` 修正 transport catch 的 timeout/offline 代码并补回归，再写响应节闭合 F-001；顺手对齐 `01-decision.md` 信息表（F-003）并考虑 F-002/F-004 测试。不要用 progress `3/4` 证明 C3 已满足检查点原文。

## 声明

本意见不修改 status/progress/检查点/方案正文/goal-tree。响应、finding 闭合与是否关门由 `/govern` 处理。
