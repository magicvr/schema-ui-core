---
id: E-001
goal: GOAL-004-w3-security-audit-remediation
title: W3 八项安全审计发现修复与回归完成
date: 2026-08-11
status: recorded
---

# E-001 · W3 八项安全审计发现修复

承接 D-001 冻结范围（P0×2 / P1×3 / P2×3 + 落盘），全部实施并回归。

## P0-1 · batch-delete 整批原子

- `authsession` 仓库新增单事务 `DeleteUsersBatch(ids, actorID)` / `DeleteRolesBatch(ids)`：先对全批执行存在性/self/last-admin（users）与存在性/system/in-use（roles）守卫，再统一删除；任一失败整批回滚。
- 通用资源工厂新增 `BatchDeleter` 可选接口（`DeleteBatch(ids, user) (int, error)`）；`POST {path}/batch-delete` 优先走接口，实体不实现时回退顺序删除。
- users/roles 实体实现 `DeleteBatch` 委托仓库；成功后在 handler 统一写 operation-log 删除事件。
- 回归：仓库层 `TestDeleteUsersBatchAtomicRollback` / `TestRolesRepositoryBatchDeleteAtomicRollback`（自删、last-admin、not-found、in-use、system、dedupe、整批回滚）+ HTTP 层 `TestUsersBatchDeleteAtomicRollbackHTTP`（批含 last admin → 409 LAST_ADMIN，先序键全部保留）。

## P0-2 · recordSource.url 拒协议相对 URL

- `apps/web/src/protocol/conformance/request-construction.ts` `buildRecordSource` 由 `(!isProtocolRelativeUrl(url) && !url.startsWith("/"))` 收紧为仅 `isRelativeProtocolUrl`（单斜杠、拒 `//`、拒反斜杠、可带 query），与 rowAction/upload 一致；`//host` 不再可能经 `authFetch` 携带 Bearer 跨源。
- 回归：新建本地 `request-construction.test.ts`（vendored `request-construction.cases.json` 由 sha256 钉死，不改）覆盖合法单斜杠/query、`//`、`/\evil`、绝对 https 四类。

## P1-1 · nginx 对齐与安全头

- `apps/web/nginx.conf`：`client_max_body_size 8m`（对齐 API `maxUploadBytes = 8<<20`）；`server_tokens off`；全站 `X-Content-Type-Options: nosniff`、`X-Frame-Options: DENY`、`Referrer-Policy: strict-origin-when-cross-origin`、保守 CSP（`default-src 'self'`、`img-src 'self' data: https:`、`style-src 'self' 'unsafe-inline'`、`script-src 'self'`、`connect-src 'self'`、`frame-ancestors 'none'`）；`location /api/` 前缀精确匹配（原 `location /api` 会匹配 `/apix`）。

## P1-2 · 登录限流真实客户端 + 成功清桶

- `rate_limit.go` 重构：`loginClientIP(r)` 仅当直接 peer 为 loopback/private（可信反代）时信任 `X-Real-IP`，否则用 peer 地址（不可伪造）；失败键 = `IP|username`（同 IP 喷多用户名不互锁、不锁无辜用户）；成功登录 `clear(key)` 清桶；map 容量上限（默认 2^16）驱逐最旧键防内存耗尽。
- `auth.go` login 接新 API；`TestLoginRateLimiterUnit`（窗口/隔离/过期/清桶/驱逐）+ `TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer` + 既有 `TestLoginRateLimit` 保持。

## P1-3 · 非 admin 不得重置 admin 密码 / demote admin

- `usersEntity.Update` 前置 `authorizeAdminTargetBoundary`：非 admin 对 admin 目标改密 → `ADMIN_ACCOUNT_FORBIDDEN`（403）；非 admin 从 admin 角色集移除 admin → 同码拒绝；自身操作仍由仓库 `SELF_OPERATION` 管辖；非 admin 账户的改密/角色维护不受影响（委派继续可用）。
- 新错误码 `ADMIN_ACCOUNT_FORBIDDEN` 加入 `error_contract_test.go` frozen 集与 `errorcatalog`（中英双语）。
- 回归：`TestUsersAdminTargetBoundary`（改密 403、无 roles.assign demote 403、有 roles.assign demote 403、非 admin 目标正常管理、admin 密码/角色不受损）。

## P2-1 · logo 路径拒反斜杠 + logout/refresh 竞态

- API：`normalizeLogoURL`（settings repository）与 `validateLogoURL`（configuration）同源分支拒绝 `\`（浏览器将 `/\evil.com` 解析为外站 host）；测试补充 `/\evil.com/logo.svg`、`/assets\evil.png`、`\evil.com` 用例。
- Web：`auth-client.ts` 引入单调 session `refreshGeneration`；`logout()` 先 bump 再清 token；`doRefresh` 在写 token 前校验 generation，in-flight refresh 在 logout 后返回也不写回；`TestAuthClient` 新增「logout 后 401 不再打 refresh」「in-flight refresh 跨 logout 不写回」两用例。

## P2-2 · HTTP Serve 失败 fail-closed

- `composition.go` `registerLifecycle` Serve goroutine：非 `ErrServerClosed` 错误 → `logger.Error` + `os.Exit(1)`（配合 compose restart 恢复），不再静默吞错留半死实例。

## 附带修复

- `docscheck` 工作区迁移遗留坏链（`docs/workspace-004-…` → `docs/workspaces/workspace-004-…`），`TestWorkspace004RootCloseoutArtifacts` 恢复绿。

## 回归证据

- `go test ./...`（apps/api）：22 包全 ok。
- `vitest run`（apps/web）：44 文件 / 746 测试全过（含新增 5 + 2）。
- Playwright e2e（APP_PROFILE=admin）：3 passed / 1 skipped（mvp-only），覆盖 shell 登录链路与 users/roles 授权管理。
