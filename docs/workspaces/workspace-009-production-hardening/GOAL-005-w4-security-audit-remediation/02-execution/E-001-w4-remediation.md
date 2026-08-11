---
id: E-001
goal: GOAL-005-w4-security-audit-remediation
title: W4 八项安全审计发现修复与回归完成
date: 2026-08-11
status: recorded
---

# E-001 · W4 八项安全审计发现修复

承接 D-001 冻结范围（P0×3 / P1×2 / P2×3 + 回归），全部实施并回归。审计基线：2026-08-11 四路并行 agent + 主路径第一手核对。

## 审计输入（四路 agent 报告 + 第一手验证）

| 面 | agent | 报告要点 |
|----|-------|----------|
| api 认证/授权/上传/输入 | aaaae409469763014 | F-1 限流驱逐死代码；F-3 上传无权限门；F-4 改密不吊销 access token 等 |
| api 内核/数据层 | a3ea8120e97a56d13 | store 单连接饿死 readyz；dev 密钥；TOCTOU 等 |
| web 协议/会话/上传 | a772356f1577bcdf7 | pageTrigger URL 弱校验；invocationId 未接通；__proto__ key 等 |
| web 组件/渲染/i18n | a03e40b237699430a | recordSource prefill 无 try/catch 白屏；login {status} 字面占位符；initTheme 白屏等 |

主路径对关键 finding 逐一人工复核（rate_limit 驱逐、upload 授权、recordSource prefill、login 文案、pageTrigger URL、initTheme 调用点），全部属实。

## P0-1 · 限流器容量驱逐修复

- **根因**：`allow()` 无条件 `l.attempts[key] = kept` 预建 key；`record()` 仅 key 不存在时驱逐 → `exists` 恒真 → 驱逐分支死代码 → 喷洒用户名 map 无界增长 → OOM。
- **修复**（`rate_limit.go`）：`allow()` 只读已存在条目，key 不存在直接放行且**不创建**；驱逐统一由 `record()` 登记时触发（含 `order` 空切片防御）。map 有界。
- **回归**：`TestLoginRateLimiterAllowDoesNotRegisterKey`（allow 不建条目 + 真实登录路径 allow→record 驱逐后 map 保持容量 2、最旧 key 可再试、最新 key 仍持失败）。

## P0-2 · 上传端点权限门

- **根因**：`POST /api/upload` 仅要求已登录，不检查任何 permission key；低权限账户可无限上传 8MiB 文件灌满磁盘。
- **修复**：
  - `composition.go` 中央注入 `files.write` PermissionContribution（ModuleID `core.server-registration`，PolicyID `PolicyAdmin`，默认仅 admin 持有）→ Reconcile 前追加。
  - `upload.go` 新增 `uploadPermissionGate`（`requirePermission` fail-closed 403）包裹 `POST /api/upload`；下载保持 owner-only 不变。
  - `testsupport/store.go`、`account/session.go`（dev session）同步 `files.write`。
  - **配额层**（D-001 范围内）：`upload.go` 新增每用户配额——`UPLOAD_MAX_FILES_PER_USER`（默认 1000 个文件）与 `UPLOAD_MAX_BYTES_PER_USER`（默认 256 MiB）环境变量可调；`quotaReached` 扫描上传目录 owner meta（含实际文件 Stat 字节），超限返回 `UPLOAD_QUOTA_EXCEEDED`（429）。新错误码入 error_contract frozen 集 + errorcatalog 双语 + web en/zh 目录。
- **回归**：`TestUploadRequiresFilesWritePermission`（viewer 403 且 0 文件落盘；admin 200）+ `TestUploadPerUserQuota`（第 3 个文件 429）。

## P0-3 · 改密即吊销 access token

- **根因**：改密只吊销 refresh token，已签发 access token（15min TTL）仍有效 → 被窃 token 在改密后仍可调用。
- **修复**（token_version 方案）：
  - 迁移 v11 `access_token_revocation`：`ALTER TABLE users ADD COLUMN token_version INTEGER NOT NULL DEFAULT 0`（additive，legacy fingerprint 路径兼容）。
  - `authsession.User` + 全部 users 表 SELECT/scan 加 `TokenVersion`。
  - `UpdateUser` 改密时同事务 `token_version = current + 1`（与 refresh 吊销同一事务原子）。
  - `auth.go`：`accessClaims{TokenVersion, RegisteredClaims}`；`SignAccessToken` 带版本；`ParseAccessToken` 返回 `ParsedAccessToken{UserID, TokenVersion}`；`Middleware` 比对版本不一致 → 401 UNAUTHENTICATED（superseded）。
  - 同步 `auth_test.go` 签名（SignAccessToken 增参、sub→sub.UserID）。
  - 迁移计数断言更新：store 三处 10→11（含 `access_token_revocation` checksum 冻结）。
- **回归**：`TestUsersPasswordChangeRevokesAccessToken`（改密前 token 200 → 改密后 401；新密码可登录、旧密码 401）。既有 `TestUsersPasswordPolicyPreservesBytesAndRevokesRefresh` 断言更新为期望 access token 同步吊销（原断言即本波修复目标行为）。

## P1-1 · 前端异常统一捕获

- **根因**：`constructRequest` 内 `serializeQueryValue` 对非标量 throw；recordSource prefill 的 useEffect、`runRequest`/`runBatchRequest` 均无 try/catch → effect 内 throw 整树卸载白屏（无 ErrorBoundary）；动作路径 unhandled rejection 静默无反馈。
- **修复**（`render.tsx`）：
  - recordSource prefill：`constructRequest` 包 try/catch → error state（不卸载）。
  - `runRequest`/`runBatchRequest`：`constructRequest` 包 try/catch → 返回 `{ok:false, code:REQUEST_CONSTRUCTION_FAILED}` 走既有反馈通道。
  - `invokeBatchAction`/`invokeAction` 的 `.then()`/`void` 调用补 `.catch` 兜底（BATCH_FAILED / ROW_ACTION_FAILED 反馈）。
  - 导入 `type RequestConstructionResult`。
- **回归**：web tsc 通过；renderer 14 文件 205 测试全绿。

## P1-2 · pageTrigger/outcomeNavigate URL 校验统一

- **根因**：`buildPageTriggerRequest`/`buildPageTriggerNavigate` 只查 `startsWith("/") && !startsWith("//")`（不拒反斜杠/空白）；`buildOutcomeNavigate` 完全不校验（绝对 URL 透传）。`/\evil.com` 可被 WHATWG 解析为外站 → 若经 authFetch 带 Bearer 外泄（当前被 schema 结构性校验屏蔽，属防御纵深缺口）。
- **修复**（`request-construction.ts`）：三处统一用 `isRelativeProtocolUrl`（单斜杠、拒 `//`、拒反斜杠/空白，可带 query）；`buildOutcomeNavigate` 拒绝绝对 URL（开放重定向面）。
- **回归**：protocol 5 文件 343 测试全绿（含 vendored `request-construction.cases.json` 现有相对 URL 用例不变）。

## P2-1 · web 启动路径加固

- **根因**：`main.tsx:72` 顶层调用 `applyStoredTheme()→initTheme()`，后者 `localStorage.getItem` 无 try/catch；站点存储禁用时模块求值失败 → React 不挂载 → 整页白屏（连 ManifestFailure 都不渲染）。index.html 内联脚本有 try/catch 而 main.tsx 无。
- **修复**（`theme.ts`）：`initTheme`/`applySystemDefaultTheme`/`setTheme` 的 localStorage 读写与 matchMedia 均包 try/catch（与 tokens.ts 一致）；存储禁用降级为「本次页面 best-effort 主题」。

## P2-2 · 登录文案 + 密码 autoComplete

- **根因**：`loginErrorKey("LOGIN_FAILED")` → `login.error.failed` = `"login failed: HTTP {status}"`，`t()` 未传 params → 渲染字面 `{status}`（`interpolate` 对缺失占位符原样保留）；`AuthError` 其实已带真实 HTTP status 但 LoginPage 忽略 `err.message`。
- **修复**：`AuthError` 增加可选 `status`；LOGIN_FAILED 构造时带 `response.status`；LoginPage `handleSubmit` 传 `{status}` 给 `t()`（`MessageParams` 结构兼容，无需显式 import 类型）。
- **autoComplete**：`form-controls.tsx` BaseInput 对 `type="password"` 加 `autoComplete="new-password"`（改密/新密码表单防浏览器误填已保存登录密码；登录页自身 `current-password` 不变）。
- **回归**：`LoginPage.test.tsx` 新增「LOGIN_FAILED 插值 503、无字面 {status}」用例。

## 回归证据

- `go test ./...`（apps/api）：23 包全 ok（含新增 P0-1/P0-2（权限门+配额）/P0-3 回归；迁移计数 10→11 同步）。
- `go vet ./...`：clean。
- `npx tsc --noEmit`（apps/web）：clean。
- `npx vitest run`（apps/web）：44 文件 / 747+ 测试全过（新增 LoginPage 1 + error key 目录；renderer/protocol/i18n 全绿）。
- e2e：未覆盖上传操作；demo profile 的 form-with-upload 由 admin（持有 files.write）访问，无行为变化。

## 范围外 / residual（见 D-001）

- store 单连接 + `context.Background()` 事务饿死 `/readyz`、dev 公开常量密钥、`page` OFFSET 上限、`WithTx` defer 回滚、health 探针指纹、上传 413 状态码、非 admin 改 admin name、operation_log/refresh token 清理、错误目录化丢弃 message —— 标记 recommended / 观察项，下波或按需处理。
