---
id: D-001-r1-correlation-contract
doc: decision-entry
status: accepted
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-002-r1-correlation-error-contract
version: 0.1.0
---

# D-001 · R1 请求链路与 correlation 契约方案冻结

- **日期**：2026-08-18
- **状态**：accepted
- **工作区**：`workspace-012-shared-cross-module-contracts`
- **决定**：在 API 入口由独立 `internal/requestid` 中间件校验/生成 `X-Request-ID`，把值写入请求 context 和响应头；所有 API 错误包络追加 `correlation_id`；Web `ResourceApiError` 同时读取错误体与响应头中的标识并把它带入用户可见错误消息；operationlog 以独立关联表保存请求与操作的 correlation id，先接入 auth 与 settings 两条真实写路径。
- **理由**：请求入口位于 `apps/api/internal/server.New`，统一包络写入点是 handler/auth 两侧的 `writeLocalizedError`，前端统一错误解析点是 `apps/web/src/renderer/resource.ts`，operationlog 的稳定写入边界是 `operationlog.Repository.RecordOperation`。该切片能覆盖 R1 四项成功标准，同时不改变 Profile、Manifest 或业务领域边界。
- **未选方案**：
  - 将 correlation id 拼入既有 `detail` JSON：会破坏现有 auth 事件的 `{username}` 稳定形状和敏感字段审计约束。
  - 在每个业务 handler 内自行生成/透传：容易遗漏未注册路由与 auth middleware 错误，无法保证“每个 API 请求”。
  - 引入 OpenTelemetry/分布式 tracing：超出 R1 非目标，留给后续可观测性波。
- **影响范围**：`apps/api/internal/requestid`、`apps/api/internal/server`、handler/auth 错误包络、operationlog 关联持久化、Web `renderer/resource.ts` 与定向测试；不修改跨 VP 的模块矩阵或协议 pin。
- **关联信息项/门禁**：`I-001`（required）已由本决策和 `E-001` 扫描证据验证，解除 S1 方案冻结门禁。
- **后续动作**：实现代码与定向测试；通过后记录 Git checkpoint，再做 R1 阶段自审与关门审计。
