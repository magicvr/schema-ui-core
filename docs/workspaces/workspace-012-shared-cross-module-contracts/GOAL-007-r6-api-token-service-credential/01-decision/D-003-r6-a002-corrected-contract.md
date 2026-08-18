---
id: D-003-r6-a002-corrected-contract
goal: GOAL-007-r6-api-token-service-credential
status: proposed
created: 2026-08-19
updated: 2026-08-19
parent: GOAL-007-r6-api-token-service-credential
version: 0.1.0
supersedes: D-002-r6-precise-contract
responds_to: A-002
---

# D-003 · R6 A-002 finding 修订契约

本决策继承 D-002 未被下列条款覆盖的全部约束，并以本条作为实现时的最终解释。A-002 F-001～F-003 required 与 F-004～F-007 recommended 均在此响应。

## 1. Migration 与 operation log（F-001）

1. `0044 service_credentials` 归 `core.auth-session`，只建立 credential 数据表及索引。
2. `0045 operation_log_service_credentials` 归 `core.operationlog`，在 0043 事件 CHECK 全集上追加：
   - `service-credentials.create`
   - `service-credentials.use`
   - `service-credentials.revoke`
3. 0045 必须沿用 0043 的 correlation-safe rebuild：
   - 建立临时 `operation_log_correlation_backup`；
   - 删除引用旧 `operation_log` 的 correlation table；
   - rebuild `operation_log` 并复制全部历史行；
   - 重建 correlation table/index、恢复全部映射并删除 backup。
4. migration tests 固定 0044/0045 顺序、历史事件/相关性保留、新三事件可写、未知事件仍被 CHECK 拒绝。

## 2. `created_by` 与用户删除（F-002）

- `created_by` 为非空 human actor id，但**不设 users 外键**。它与 operation log 的 `actor_id` 一样是持久历史标识；创建者用户之后被删除不影响机器凭据生命周期或审计可读性。
- `DeleteUser` / batch delete 保持现有行为，不级联、不阻断、不吊销 service credential。凭据只能由 human 管理 API 显式吊销。
- 管理 metadata 仅返回 `createdBy` id，不在读取时强制 join 用户；避免历史主体删除后 list/detail 失败。

## 3. 名称唯一与稳定错误（F-003）

- DDL：`name TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) BETWEEN 1 AND 100)`；handler 存储 trim 后名称。
- repository 将 `service_credentials.name` 唯一约束映射为 `ErrCredentialNameTaken`。
- handler 返回既有 `400 INVALID_CREATE_FIELD`，并附 `fieldErrors=[{field:"name", reason:"name already exists"}]`；不扩 error catalog。
- repository 并发测试必须证明 `Build` 与 `build` 只有一个创建成功，失败方不是 `INTERNAL`。

## 4. Prefix 与 `devSession`（F-004）

- `sui_sc_` 判断必须在 `ParseAccessToken` 与任何 `devSession` fallback **之前**。
- 匹配 service prefix 后，unknown/malformed/expired/revoked 一律直接 `401 UNAUTHENTICATED`；禁止回落 JWT 或 `injectDevSession`。
- 非 service prefix 保持现有 JWT + 显式 dev-session 行为。

## 5. 管理审计事务与 use best-effort（F-005）

- `service-credentials.create/revoke` 是管理 mutation，其 credential row 变更和 operation log 写入必须在**同一个 Store transaction**内完成；审计失败回滚 mutation 并返回 `500 INTERNAL`。
- `authsession.Repository` 的 create/revoke 方法接收不依赖 operationlog 包的 transaction callback；`operationlog.Repository` 提供使用既有 transaction 的 recorder。这样保持包边界，同时保证原子性。
- 重复 revoke 若 row 已有 `revoked_at`，直接 `204` 且不追加第二条 revoke audit；首次状态转换与唯一 audit 同事务提交。
- `service-credentials.use` 与 `last_used_at` 是请求热路径，仍为 best-effort；失败写 service log，但不撤销已经通过的认证。

## 6. 用户自助面与 data-permission `self`（F-006）

S2 必须把下列现有 `IdentityFrom` 用户自助消费者改为 `UserIdentityFrom` 并逐项测试机器凭据返回 401：

- `handler/account.go`（`/api/accounts/me`）
- `handler/account_self.go`
- `handler/account_avatar.go` 的 authenticated mutation
- `handler/mfa.go` self-service
- `handler/notifications.go`
- `handler/wallet_self.go`

permission-gated resource/export/import/upload 等路径保留通用 identity。data-permission 的 `self` 以 synthetic machine principal id `service-credential:<id>` 为作用对象，**不**模拟或继承 `created_by` 用户的 row scope。

## 7. Metadata/list 与 ID（F-007）

- credential `id` 固定为 16 random bytes 的 32-char lowercase hex；路径按该形状校验，非法/不存在均 `404 NOT_FOUND`。
- list 默认包含 active、expired、revoked 全部行，并返回 `revokedAt` / `lastUsedAt`；不返回 secret/hash。
- list 使用 `page` / `pageSize`（默认 1/20，最大 100）、`createdAt DESC, id DESC`，响应为 `{items,total,page,pageSize}`；detail 与 list item 使用同一 metadata shape。

## 8. 修订后的 0044 DDL 约束

`service_credentials` 字段：

- `id TEXT PRIMARY KEY CHECK (length(id)=32)`
- `name TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) BETWEEN 1 AND 100)`
- `token_prefix TEXT NOT NULL CHECK (length(token_prefix)=15)`
- `token_hash TEXT NOT NULL UNIQUE CHECK (length(token_hash)=64)`
- `scopes TEXT NOT NULL`
- `expires_at INTEGER NOT NULL`
- `revoked_at INTEGER`、`last_used_at INTEGER`
- `created_by TEXT NOT NULL`（无 FK）
- `created_at INTEGER NOT NULL`、`updated_at INTEGER NOT NULL`
- index：`created_at DESC, id DESC` 与 `expires_at`

其余 secret、scope ceiling、human-only 管理、错误复用、Profile/Manifest/protocol 不变式和测试矩阵继续以 D-002 为准。

## 门禁

D-003 在 A-002 finding-closure independent 通过前保持 `proposed`；I-002～I-004 继续 `collecting`，S1/S2 不放行。
