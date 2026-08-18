---
id: D-003-r5-a002-corrected-contract
goal: GOAL-006-r5-maintenance-read-only-gate
status: proposed
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# D-003 · R5 A-002 finding 修订契约

## 修订依据

本稿响应 A-002 F-001/F-002；D-002 原文保留为历史提案，不静默改写。A-002 指出的 F-003～F-006 为 recommended，转为 S1/S2 实施验证门。

## 运行态与写门禁

运行态仍为启动配置 `runtime.mode`：`normal`、`maintenance`、`degraded`、`read-only`；`RUNTIME_MODE` 覆盖 YAML。已注册的业务 `POST/PUT/PATCH/DELETE` 在后三种 mode 一律拒绝，auth/session recovery 与强制改密 allowlist 保持放行；GET/HEAD、探针、未知路径和方法不匹配继续走既有 handler/error envelope。门禁位于 request-id/CORS 之后、mux handler 之前，并只按精确注册方法判定。

拒绝响应统一为 HTTP `503`，分别使用 `SERVICE_MAINTENANCE`、`SERVICE_DEGRADED`、`SERVICE_READ_ONLY`；客户端必须以 error code 而不是 HTTP status 做运行态分流。这样避免 `423` 与现有 `ACCOUNT_LOCKED` 客户端语义冲突。受控态写拒绝是应用内 API 包络错误，不触发 Host failure recovery、账号锁定或 `Retry-After`；只有 bootstrap `maintenance` 仍可产生 Host `MAINTENANCE` 与既有 manual recovery。

## Host 与 status 投影

- system-monitoring status 在现有单行 envelope 中追加 `availabilityMode`，表达后端原始 mode；`status`/`ready` 仍只表达存储与模块图 readiness。
- bootstrap `normal` → `normal`，`maintenance` → `maintenance`，`degraded` → `degraded`，`read-only` → `degraded`。后三者均**省略** `disabledCapabilities`，不把任何协议 capability 当作运行态开关；read-only 与 degraded 的权威区别只在 status 的 `availabilityMode`。
- 不新增 `host.readOnly` 文案依赖；若后续 UI 需要横幅，读取 status 的原始 mode。`READY_DEGRADED` 继续进入应用，不应被当作 Host 终态。

## 配置 fail-closed 与验证门

- `RUNTIME_MODE` 的显式空字符串必须 fail closed，不能复用把空值视为未设置的 `envOr`；YAML 缺省才使用 `normal`。
- S1 必须覆盖 YAML/env precedence、空值、未知值与解析错误；S2 必须覆盖核心/Provider 真实写路由、精确 allowlist、未知路径 404/405、request-id、本地化 envelope 与客户端按 code 分流。
- 不改变 Profile 默认集、模块依赖闭包、Manifest bytes/聚合算法、协议 pin 或 health/readiness 语义。

## 门禁

本稿须经 A-004 independent closure 后将 D-003 标为 `accepted`，I-002～I-004 才可从 collecting 改为 verified 并放行 S1。
