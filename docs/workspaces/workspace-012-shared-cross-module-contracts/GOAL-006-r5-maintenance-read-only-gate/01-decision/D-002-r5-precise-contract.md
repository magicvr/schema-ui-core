---
id: D-002-r5-precise-contract
goal: GOAL-006-r5-maintenance-read-only-gate
status: proposed
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# D-002 · R5 S0 精确模式、写门禁与 Host 投影

## 1. 运行态来源

`runtime.mode` 是唯一运行态配置，支持 `normal`、`maintenance`、`degraded`、`read-only`。YAML 由 `RUNTIME_MODE` 环境变量覆盖；首版仅在启动时读取，不提供热切换或管理写 API。缺省为 `normal`；未知值、空的显式值或解析错误 fail closed。

## 2. 写门禁矩阵

| mode | 已注册业务 `POST/PUT/PATCH/DELETE` | auth/session recovery allowlist | GET/HEAD/probe | 拒绝 |
|------|-----------------------------------|-------------------------------|----------------|------|
| normal | 放行到既有认证/权限 | 放行 | 放行 | — |
| maintenance | 拒绝 | 放行 `/api/auth/login`、`/api/auth/refresh`、`/api/auth/logout`、`/api/auth/mfa/verify`、`/api/account/password` | 放行 | `503 SERVICE_MAINTENANCE` |
| degraded | 拒绝 | 同 maintenance | 放行 | `503 SERVICE_DEGRADED` |
| read-only | 拒绝 | 同 maintenance | 放行 | `423 SERVICE_READ_ONLY` |

门禁位于 request-id/CORS 之后、既有 mux handler 之前。只对 mux 已注册的当前方法判定；未知路径与方法不匹配继续由 `WithJSONRouteErrors` 产生 `404/405`。因此业务写请求在受控态不会泄露未认证/无权限差异，allowlist 内仍保留原有认证语义。OPTIONS 由 CORS 层先处理。

## 3. 状态投影

- `/api/system-monitoring/status` 在现有单行 list envelope 中追加 `availabilityMode`，值为后端原始 mode；原有 `status`/`ready` 字段仍只表达存储与模块 readiness。
- bootstrap `normal` → `normal`，`maintenance` → `maintenance`，`degraded` → `degraded` + `disabledCapabilities:["form.controls.readonly"]`，`read-only` → `degraded` + 同一 capability + `messageKey:"host.readOnly"`。不扩展上游 bootstrap mode，不产生未知 capability。
- `/healthz` 仍是不访问数据库的存活探针；`/readyz` 仍只表达数据库与模块图 readiness，不随运行态改写。

## 4. 错误与兼容

新增 catalog codes：`SERVICE_MAINTENANCE`、`SERVICE_DEGRADED`、`SERVICE_READ_ONLY`。响应沿用 R1 localization/error envelope 与 request-id；不写入 operation log。maintenance/degraded 的 503 不承诺 Retry-After，客户端使用 bootstrap 的 maintenance/manual 或 degraded recovery 语义。

## 5. 不变式与验证

- 不新增 Profile/module ID，不修改 profile defaults、module dependency closure、Manifest 聚合算法、Manifest bytes、protocol pin 或 readiness semantics。
- 必须用 server/handler 黑盒测试覆盖：三种受控态的 core/provider 写、allowlist、GET/HEAD、404/405、认证优先级、错误包络与 request-id；bootstrap/status 正常态兼容和 read-only/degraded 映射；配置 YAML/env precedence 与非法 mode fail closed。

## 6. 门禁

本决定在 cross 模式下先经 self + independent 设计审计；A-002 pass 且 I-002～I-004 关闭后，才允许进入 S1 实施。
