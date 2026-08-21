---
id: D-002-r6-precise-contract
goal: GOAL-007-r6-api-token-service-credential
status: superseded
created: 2026-08-19
updated: 2026-08-19
parent: GOAL-007-r6-api-token-service-credential
version: 0.1.0
---

# D-002 · R6 精确机器凭据契约

## 1. Secret 与认证分流

1. raw secret 格式固定为 `sui_sc_<random>`；`random` 为 `crypto/rand` 生成的 32 bytes（256-bit）并用 unpadded Base64URL 编码。
2. 创建响应只返回一次完整 raw secret。持久层只保存完整 raw 的 SHA-256 hex hash，以及安全展示用 `tokenPrefix`（`sui_sc_` + 前 8 个随机字符）；不得保存或再次返回 raw。
3. `Authorization: Bearer` 值以大小写敏感的 `sui_sc_` 起始时进入 service credential 分支；其余值保持现有 JWT 解析。任一分支失败都返回相同 `401 UNAUTHENTICATED`，不暴露不存在/过期/吊销差异。
4. service credential 不签 JWT、不生成 refresh token、不读取用户 `token_version`，也不触发浏览器 refresh 行为。

## 2. 数据模型与生命周期

新增全局 migration `0044 service_credentials`，归 `core.auth-session`：

| 字段 | 约束 |
|------|------|
| `id` | 随机稳定主键 |
| `name` | 非空、大小写不敏感唯一 |
| `token_prefix` | 非空安全展示值 |
| `token_hash` | 64-char SHA-256 hex、unique |
| `scopes` | JSON string array；创建时去重排序 |
| `expires_at` | Unix seconds；必填 |
| `revoked_at` | nullable Unix seconds |
| `last_used_at` | nullable Unix seconds |
| `created_by` | human user id，FK users，`ON DELETE RESTRICT` |
| `created_at` / `updated_at` | Unix seconds |

- 创建 `expiresAt` 必须是 RFC3339、晚于当前时刻且不超过 365 天；这是首版有界行业实践，避免无期限 bearer secret。
- 首版不提供 PATCH/extend。轮换 = 创建新 credential 后吊销旧 credential。
- 吊销使用 `UPDATE ... WHERE revoked_at IS NULL`；已吊销重复请求返回相同 `204`，不存在返回 `404 NOT_FOUND`。无可变 PATCH，故不引入 ETag/expectedVersion；并发安全由 unique constraints 与 guarded revoke 提供。

## 3. Principal 与用户自助隔离

1. `account.User` 增加不序列化的 principal kind/credential id；现有空 kind 兼容解释为 user，service credential 明确标为 `service-credential`。
2. 有效 credential 投影：`ID=service-credential:<credential-id>`、`Name=<credential name>`、`Roles=[]`、`Permissions=<frozen scopes>`，并写入现有 request identity context，使统一 `requirePermission` 可复用。
3. 新增 `UserIdentityFrom`：只接受 user principal。`/api/accounts/me`、account profile/password/session/avatar、MFA self-service、notifications 与 wallet self-service 改用该 helper，机器凭据进入这些路径统一 `401 UNAUTHENTICATED`。
4. 资源、导入导出、上传等 permission-gated API 继续使用通用 identity；service credential 仅能通过其 scope 中的既有 permission key。

## 4. 管理 API 与授权上限

| method | path | 权限 | 结果 |
|--------|------|------|------|
| GET | `/api/service-credentials` | `service-credentials.read` | metadata list；不含 secret/hash |
| GET | `/api/service-credentials/{id}` | `service-credentials.read` | metadata detail；不存在 404 |
| POST | `/api/service-credentials` | `service-credentials.write` | 201；metadata + 一次性 `secret` |
| POST | `/api/service-credentials/{id}/revoke` | `service-credentials.write` | 204；幂等吊销 |

- 四条管理路由除 permission 外必须是 human user principal；service credential 不可管理 credential。
- 两个管理 permission 由 composition 以 `core.auth-session` system-data contribution 注入，默认仅 `system.admin`；不新增 module、page、navigation 或 Manifest fragment。
- `name` trim 后 1～100 chars；`scopes` 为 1～64 个非空唯一 permission key。
- 每个 scope 必须存在于 reconciled `permissions` catalog，且必须是创建者当前 effective permissions 的子集；`service-credentials.read/write` 为保留管理权限，不得写入 credential scopes。未知 key → `400 INVALID_PERMISSION_REF`，超出创建者权限或包含保留 key → `403 FORBIDDEN`。
- 管理 mutation 由 R5 final operational gate 自动覆盖；GET metadata 在受控态仍保持读取语义。

## 5. 错误、审计与敏感数据

- 复用冻结错误码：`UNAUTHENTICATED`、`FORBIDDEN`、`INVALID_CREATE_BODY`、`INVALID_CREATE_FIELD`、`INVALID_PERMISSION_REF`、`NOT_FOUND`、`INTERNAL` 及 R5 `SERVICE_*`；本目标不扩展 error catalog。
- operation log 新增事件 `service-credentials.create/use/revoke`。
- create/revoke 记录 human actor、credential record id、name/scopes/tokenPrefix/expiry/revokedAt 与 correlation id；use 记录 service actor、credential id、scope count、请求 method/path 与 correlation id。
- `secret`、`token_hash`、Authorization header 永不进入 response metadata、operation detail、service log 或 public Manifest。operation log 失败与 `last_used_at` 更新失败均为 best-effort，不改变已经通过的认证结果，但写 service log。

## 6. 组合与不变式

- repository 与 migration 留在 `core.auth-session`；管理路由由 composition 中央挂载，因为该能力随 `core.auth-session` 存在而非新 Profile module。
- permission contribution 只增加 reconciled authorization rows；不改 `kernel.BuiltinModules()`、Profile IDs/default set、Provider/module matrix、Manifest fragments/bytes、Host bootstrap、protocol version/pin 或 readiness 定义。
- 认证 middleware 保持 request-id/error envelope；service credential use audit 通过 composition 注入的 best-effort recorder，避免 auth core 直接依赖 operationlog 实现包。

## 7. 验证矩阵

1. migration/repository：hash-only、唯一名/hash、list/detail 不泄密、过期、原子幂等吊销、last-used。
2. auth：JWT 非回归；credential success、unknown、malformed、expired、revoked；scope identity；user-only helper 拒绝 service。
3. handler：human/permission gates、scope catalog/subset/reserved、一次性 secret、metadata、404/idempotent revoke、审计/correlation。
4. composition：真实管理路由 + service credential 调用至少一个 core/Provider permission surface；R5 controlled-mode mutation；Profile/Manifest bytes 与协议 pin 不变。
5. 全量：`go test -timeout 15m ./...`、Web targeted/build（若生成 claim 变化则核查内容）、docscheck、whitespace。

## 未选方案

- 独立 Profile/module + 管理页面：改变装配矩阵，超出 R6 最小管理面。
- service credential 继承创建者实时角色：会使 scope 随用户会话漂移，破坏冻结与独立生命周期。
- 永不过期或服务端可再次 reveal：扩大 bearer secret 风险。
- service credential 自举创建/吊销：首版增加持久化越权链，故 human-only fail closed。

## 门禁

D-002 经 A-002 independent 复核后由 D-003 修订并 supersede；I-002～I-004 保持 `collecting`，S1/S2 不放行，直至 finding-closure independent 通过。
