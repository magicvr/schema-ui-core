---
id: E-002-r4-feedback-surfaces
doc: execution-entry
status: recorded
goal_id: GOAL-005-r4-unified-feedback-recovery
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-005-r4-unified-feedback-recovery
version: 1.0.0
---

# E-002 · R4 分类与共享反馈表面

## 已发生事实

- 新增 `apps/web/src/renderer/feedback-policy.ts`：把 HTTP status、既有 error code、`AbortError` 和浏览器网络 `TypeError` 映射为 validation/auth/forbidden/not-found/conflict/rate-limited/maintenance/unavailable/offline/timeout/unknown；只有调用方显式传入读取 retry 时才附加 retry。
- `ActionResult` 及 request/batch/custom action 的非 2xx 结果保留 status、messageKey/params 与 correlationId；catalog 文案优先，`safeParams` 只允许 string/number，原始请求体不进入反馈表面。
- 新增 `apps/web/src/components/ui/feedback.tsx` 并从 `components/ui/index.ts` 导出：成功使用短暂 `role=status`，错误使用持久 `role=alert`；toast 可 dismiss，inline 默认保持；retry 是键盘可达按钮并由 ref 防止 in-flight 双击。
- SchemaTable、statCard、chart、recordSource prefill 和 FormInner 已接入共享表面；列表/展示/recordSource 只对幂等读取提供显式 retry，表单写失败不附加 retry，表单诊断 code 保留为非主文案信息。
- `en-US` 与 `zh-CN` 均补齐 R4 分类文案、retrying、dismiss 文案；`HostFailureScreen` 未被普通 resource feedback 替代，继续保留 Host 终态恢复边界。

## 证据

- `src/renderer/feedback-policy.test.ts` 覆盖 status/code/transport 分类、catalog key、correlation 与读取/写入 retry 边界。
- `src/components/ui/feedback.test.tsx` 覆盖 status/alert、success auto-dismiss、error 持久、dismiss、键盘按钮与双击防重试。
