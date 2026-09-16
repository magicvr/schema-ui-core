---
id: A-003-r4-f001-recheck
doc: audit-opinion
status: recorded
source: independent
verdict: pass
scope: A-002 required F-001 finding-closure recheck (transport timeout/offline classification, write vs recordSource retry, E-004 tests/tsc)
audit_type: finding-closure
goal_id: GOAL-005-r4-unified-feedback-recovery
auditor: grok-build (grok-4.6 · reasoning high)
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-005-r4-unified-feedback-recovery
version: 1.0.0
---

# A-003 · A-002 F-001 finding-closure independent recheck（2026-09-17）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure；仅复核 A-002 required F-001 是否已按 `fixed` 合法闭合（`render.tsx` transport catch、写入无 retry、recordSource 显式 retry、E-004 测试/tsc）。不含 C4 关门、Git checkpoint、Root R4 投影、R5。
- **verdict**：pass

## 范围与区间

工作区：`workspace-037-admin-workflow-continuity`（`root_goal: GOAL-001-admin-workflow-continuity`，`canonical_scope` 匹配，`shared_materials_catalog: none`，`primary_plan: VP-037-admin-workflow-continuity`）。未读取其他工作区。无固定共享资料引用。

对照 A-002 F-001 闭合要求：transport catch 保留 `AbortError` → `REQUEST_TIMEOUT` / `feedback.timeout`，仅将网络 `TypeError` / `Failed to fetch` 标为 `REQUEST_FAILED` / `feedback.offline`；写入无 retry；recordSource 读取显式 retry 一次。本意见不修改 `status` / `progress` / 方案正文 / goal-tree。

## 成果（有证据）

### 生产 catch 不再扁平化为 `REQUEST_FAILED`

`transportFailureResult`（`apps/web/src/renderer/render.tsx` 661–671 行）把原始 error 交给 `feedbackFromError(error)`，**不**传 retry，再写出 `code` / `message` / `messageKey`。A-002 点名的五处 catch 现为：

| 调用方 | 现行处理 | 行号 |
|--------|----------|------|
| custom `library.preview` fetch catch | `transportFailureResult(error)` | 413–415 |
| custom download fetch catch | `transportFailureResult(error)` | 441–444 |
| `runRequest` fetch catch | `transportFailureResult(error)` | 606–609 |
| `runBatchRequest` fetch catch | `transportFailureResult(error)` | 754–755 |
| `useRecordSourcePrefill` fetch catch | `feedbackFromError(error, { retry })` | 1718–1725 |
| FormInner 防御性 submit catch | `feedbackFromError(error)`（无 retry） | 2000–2006 |

`failureLikeOf` 仍先识别 `AbortError` → `REQUEST_TIMEOUT`，再识别 `TypeError` / `Failed to fetch` → `REQUEST_FAILED`（`feedback-policy.ts` 45–61、85–89 行）。写入结果经 `errorFeedback(result)` / `feedbackFromError(result)` 二次分类时保留 `messageKey`，不附加 retry。

### 写入无 retry、recordSource 读取显式 retry

- 表单提交：`runRequest` → `transportFailureResult` → FormInner `feedbackFromError(result)`，不传 `retry`。回归断言 AbortError 文案含 `The request timed out`、字段值保留、无 Retry 按钮、fetcher 仍为 1 次（`render.test.tsx` 1213–1230 行）。
- 离线写入：同一路径对 `TypeError("Failed to fetch")` 显示 `network is unavailable` 且按钮重新可用（1232–1254 行）；因未传 retry，不会出现写重试。
- recordSource：`feedbackFromError(error, { retry })`；AbortError 与 Failed to fetch 均显示对应 catalog 文案，Retry 点击后 fetcher 从 1 次变为 2 次并恢复字段（1016–1031 行）。
- 行/页 action：`runRowAction` 失败走 `errorFeedback(result)`，同样不传 retry。外层 `.catch` 的 `ROW_ACTION_FAILED` / `BATCH_FAILED` 仅覆盖 `runRequest`/`runBatchRequest` 已 resolve 之后的未预见拒绝；现行 fetch catch 不再抛出 AbortError。

### E-004 测试 / tsc 本轮独立复跑

在 `apps/web`、vitest 3.2.7：

| 集合 | 本轮结果 |
|------|----------|
| `render.test.tsx` / `feedback-policy.test.ts` / `feedback.test.tsx` / `schema-table.test.tsx` | **4 files / 92 passed** |
| `npx tsc -p tsconfig.app.json --noEmit` | exit 0 |

与 E-004 所记「4 个测试文件、92 项与 tsc 通过」一致，可重复核对。

## 对照成功标准

| 标准（F-001 闭合要求） | 状态 | 证据 |
|------------------------|------|------|
| AbortError → `feedback.timeout`（或原始 error 交给 `feedbackFromError`） | 达成 | `transportFailureResult` + `failureLikeOf`；表单回归显示 timeout catalog |
| 仅网络 TypeError / Failed to fetch → `REQUEST_FAILED` / `feedback.offline` | 达成 | policy 单测 + FormInner C5 + recordSource `it.each` |
| ① 表单/action AbortError 显示 timeout **且无** retry | 达成 | `render.test.tsx` 1213–1230；action 共用 `transportFailureResult` / `errorFeedback` |
| ② recordSource AbortError 显示 timeout **且** 显式 retry 一次 | 达成 | `render.test.tsx` 1016–1031 |
| 写入不自动 retry | 达成 | FormInner / `errorFeedback` 不传 retry；AbortError 用例 fetcher 仍为 1 |
| R4-I-001（本 finding 缺口） | 可维持 verified | timeout/offline 分裂已回到默认超时写入路径 |
| R4-I-004（本 finding 削弱部分） | 本缺口已补 | transport catch 已保留分类；Host 分层不在本条范围 |
| C4 / GOAL-005 `done` / Root R4 | 仍未完成 | 本意见只闭合 F-001；checkpoint 与投影仍属 `/govern` |

## Findings

无新的 required finding。无新的 recommended finding。

A-002 既有条目在本 scope 内的独立判断：

| 既有 finding | 本轮判断 | 说明 |
|--------------|----------|------|
| A-002 F-001 required | **fixed** | 五处点名 catch + FormInner 防御 catch 已保留 timeout/offline；最低回归两条均存在且本轮复跑通过 |
| A-002 F-002 recommended | 仍 open（不阻断本 recheck） | 表单/recordSource 现已有组件级 timeout/offline；maintenance 分类与 Host 对照断言仍缺，不升格 required |
| A-002 F-003 recommended | 已由决策索引对齐（不阻断） | `01-decision.md` 信息表现为 verified，与 `00-meta` / 本索引不再冲突；本意见不改决策正文 |
| A-002 F-004 recommended | 仍 open（不阻断） | 403/401 列表无 retry、chart 一次 retry 仍无专项断言 |

未将 batch/custom preview-download 缺少独立 AbortError 用例升格为新 finding：二者已走同一 `transportFailureResult`，F-001 最低测试门槛已满足。

## 必改项汇总

无。本 scope 开放 required = 0。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 独立判断 |
|----|------|----------|----------|
| R4-I-001 | required | C1 | F-001 缺口已补；可维持 verified |
| R4-I-002 | required | C2 | 写入仍无 retry；recordSource 仍仅显式 retry；可维持 verified |
| R4-I-003 | required | C2 | 不在本 finding 范围；未发现回退 |
| R4-I-004 | required | C3 | F-001 对 timeout/offline 生产路径的削弱已解除；maintenance 自动化仍属 F-002 recommended |
| R4-I-005 | deferred non-blocking | R5/支持需求 | 仍成立；owner=`/vision` |

无到期未接受残余的 required 信息项。`shared_materials_catalog: none`，无资料引用被当作关闭证据。

## 与既有意见的异同

同意 A-002 对 F-001 的原判定（当时 `REQUEST_FAILED` 扁平化使默认 30s `withTimeout` AbortError 显示成离线文案）。不同意继续把该条视为开放 required：E-004 代码与本轮复跑已满足其书面闭合要求。

与 A-001 self `pass` 不再在 timeout/offline 生产路径上冲突。F-002～F-004 保持 recommended，默认不阻断 C4。无 residual / overrule 请求，无 P-004 冲突。

## 结论 + 建议给编排器/用户的下一步

独立审计 **pass**。A-002 required F-001 已按 `fixed` 合法闭合；本 scope 无新的 required finding。

建议 `/govern`：登记本意见对 F-001 的闭合确认，建立 R4 Git checkpoint，完成 C4 并投影 Root R4。不要用 progress `3/4` 代替 checkpoint。F-002～F-004 可在 C4 响应中决定是否顺手补测，不阻断关门。

## 声明

本意见不修改 status/progress/检查点/方案正文/goal-tree。响应、C4 关门与 Root 投影由 `/govern` 处理。
