---
id: A-002
goal: GOAL-005-w4-security-audit-remediation
title: W4 八项修复独立交叉审计
source: independent
date: 2026-08-11
verdict: pass
status: recorded
auditor: grok-build · grok-4.5 · high · audit skill
parent: GOAL-001-production-hardening
created: 2026-08-11
updated: 2026-08-11
version: 0.1.0
---

# A-002 · W4 八项修复独立交叉审计

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build · grok-4.5 · high · `audit` skill |
| **类型** | execution-facts / close-out 前复核 |
| **scope** | GOAL-005 W4 八项修复：实现正确性 + 回归风险 + 是否有遗漏的绕过面（限流驱逐、上传权限门+配额、改密吊销 access token、前端异常捕获、URL 校验统一、web 启动加固、登录文案、autoComplete） |
| **verdict** | **pass**（开放 required = 0；N-001～N-005 均为 recommended） |

## 范围与区间

- **工作区**：`workspace-009-production-hardening` · Root `GOAL-001-production-hardening` · canonical 仅本区 · `shared_materials_catalog: none` · `primary_plan: VP-009-production-hardening`
- **代码**：`apps/api`（`handler/rate_limit.go`、`upload.go`、`auth/auth.go`、`modules/authsession/*`、`composition.go`、`errorcatalog`、相关 `*_test.go`）+ `apps/web`（`renderer/render.tsx`、`protocol/conformance/request-construction.ts`、`theme/theme.ts`、`LoginPage.tsx`、`auth-client.ts`、`form-controls.tsx`、i18n）
- **过程**：`00-meta.md`、`01-decision/D-001-w4-scope.md`、`02-execution/E-001-w4-remediation.md`、`03-audit/A-001-w4-self.md`
- **信息项**：I-001 / I-002 均为 `verified`；本 scope 无到期 open required 信息门禁
- **不覆盖**：其他工作区；D-001「明确不做」项（readyz 饿死、dev 密钥、Compose TLS 等）的实现验收；e2e 全量

## 成果（有证据 · 独立核对）

| # | 项 | independent | 证据摘要 |
|---|----|-------------|----------|
| 1 | P0-1 限流容量驱逐 | **pass** | `allow()` 仅读已有 key、不预建（`rate_limit.go:40-65`）；`record()` 建 key + 容量驱逐 + `order` 空切片守卫（:70-84）；`TestLoginRateLimiterAllowDoesNotRegisterKey` 验证 spray 后 map=2 |
| 2 | P0-2 上传权限门 | **pass** | `composition.go` 中央注入 `files.write` + `PolicyAdmin`；`uploadPermissionGate` → `requirePermission(..., "files.write")` fail-closed；GET 下载仍 owner-only；`TestUploadRequiresFilesWritePermission` viewer 403 / 0 落盘 |
| 3 | P0-2 上传配额 | **pass** | `quotaReached` 按 owner meta + Stat 计文件数/字节；超限 `UPLOAD_QUOTA_EXCEEDED` 429；errorcatalog + frozen 集 + en/zh i18n；`TestUploadPerUserQuota` 第 3 文件 429 |
| 4 | P0-3 改密吊销 access | **pass** | 迁移 v11 `token_version`；`UpdateUser` 改密同事务 +1 并吊销 refresh；`SignAccessToken`/`Middleware` claims `tv` 比对；Refresh 路径 `UserByID` + `issue` 带当前版本；旧 refresh 已吊销无法换新 access；`TestUsersPasswordChangeRevokesAccessToken` |
| 5 | P1-1 前端异常捕获 | **pass** | `runRequest` / `runBatchRequest` / recordSource prefill 三处 `constructRequest` try/catch；`invokeAction` / `invokeBatchAction` 补 `.catch`；fetcher 网络错误亦返回 ActionResult |
| 6 | P1-2 URL 校验统一 | **pass** | `PROTOCOL_URL_RE` 拒 `//`/反斜杠/空白；`buildPageTriggerRequest`/`Navigate`/`buildOutcomeNavigate` 均走 `isRelativeProtocolUrl`；`request-construction.test.ts` 覆盖 `//`、`/\\`、绝对 URL |
| 7 | P2-1 web 启动加固 | **pass** | `initTheme` / `applySystemDefaultTheme` / `setTheme` localStorage + matchMedia 全 try/catch |
| 8 | P2-2 登录文案 + autoComplete | **pass** | `AuthError.status` → `t(..., {status})`；`LoginPage.test` 断言无字面 `{status}`；`BaseInput` password → `autoComplete="new-password"`；登录页仍 `current-password` |

**本轮复跑（independent，2026-08-11）**

- `go test ./internal/handler/ -run "TestLoginRateLimiter|TestUploadRequires|TestUploadPerUserQuota|TestUsersPasswordChange"` → ok
- `go test ./internal/auth/` → ok
- `go test ./internal/store/ -run "TestMigrate|TestApplied|TestRestart"` → ok
- `npx vitest run src/protocol/conformance/request-construction.test.ts src/app/LoginPage.test.tsx` → 12 passed

（未复跑完整 `go test ./...` / 全量 vitest；self A-001 已记录全绿，本轮抽样覆盖八项关键回归。）

## 绕过面核对（prompt 重点）

| 面 | 结论 |
|----|------|
| 限流 `allow` 预建 key | 已消除；spray 路径由 `record` 驱逐 |
| 上传权限绕过 | POST 唯一挂载点经 gate；无第二上传入口 |
| meta 篡改绕过配额 | 进程内仅服务端写 meta；字节用 Stat；不可读 meta 保守计入（偏向拒绝） |
| GET 下载权限 | 仍 owner-only，不依赖 `files.write`（正确） |
| 改密后 access 仍可用 | middleware 版本比对 + 同事务吊销 refresh → 即时失效 |
| Refresh 换新 access | 改密后 refresh 已 revoked，无法换带新 `tv` 的 token |
| constructRequest 白屏 | 三 UI 调用点均捕获；协议层 builder 本身返回 `ok:false` 不抛 |
| pageTrigger 反斜杠/绝对 URL | 统一正则拒绝 |

## Findings

### N-001 · recommended · low · 配额扫描 O(files)（与 self A-001 同向）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| severity | low |
| evidence | `apps/api/internal/handler/upload.go` `quotaReached` 每次上传 `ReadDir` + 逐文件 Stat |

在 `files.write` 默认 admin-only + 默认 1000 文件上限下可接受。量级上升后应持久化计数或定期清理。**不阻断关门。**

### N-002 · recommended · low · 迁移计数硬编码三处 10→11（与 self A-001 同向）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| severity | low |
| evidence | `store/migrate_test.go`、`operations_test.go`、`restart_test.go`；`access_token_revocation` checksum `c3ea720a…` |

维护面已知；新增迁移须同步。**不阻断。**

### N-003 · recommended · low · 配额 check-then-save 并发 TOCTOU

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| severity | low |
| evidence | `upload.go`：`quotaReached` 与 `save` 非原子；并发两次均可在计数下通过后双写 |

需已持有 `files.write` 且并发打满边界。单实例 admin 工具场景风险有限；多 worker / 高并发上传时可能短暂超配额。**不阻断本波声明的 best-effort 配额。**

### N-004 · recommended · low · D-001 P2「低危批量」未交付且未进成功标准

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| severity | low |
| evidence | D-001「范围内」P2 第 8 条（page 上限 / WithTx defer / 413 / health 指纹）；`00-meta` 八项成功标准未列入；E-001 residual 标 recommended |

实现与 8 项成功标准一致，但决策「范围内」与 residual 表述略不一致。建议下波立项或正式 residual 接受，避免范围漂移。**不阻断按 00-meta 关门。**

### N-005 · recommended · low · workspace.md 波次表未登记 W4

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| severity | low |
| evidence | `workspace.md` 波次表仅 W1–W3；`goal-tree.md` 已含 GOAL-005 [active] |

工作区上下文指针滞后于目标树。建议 `/govern` 响应时补一行 W4。**不阻断 GOAL-005 关门。**

## 必改项汇总

- 开放 **required**：**0**
- recommended：N-001～N-005（均可带入下波或 residual，不挡本波关门）

## 与既有意见的异同

| 点 | A-001 (self) | A-002 (independent) |
|----|--------------|---------------------|
| 八项实现主张 | pass | 一致 pass（抽样复跑关键回归） |
| F-001 allow 剪枝语义 | fixed / 无缺陷 | 同意：拒绝分支写回 kept 正确 |
| N-001 配额扫描 | recommended | 同意并保留 |
| N-002 迁移计数 | recommended | 同意并保留 |
| 新增 | — | N-003 配额 TOCTOU；N-004 D-001 低危范围一致性；N-005 workspace 波次表 |
| verdict | conditional（无 open required） | **pass**（无 open required / 无 med required 缺口） |

self 的 conditional 主要因 recommended 提示；按 prompts 尺度（无未关闭 high required、无到期 required 信息项）本独立审给 **pass**。

## 结论 + 建议给编排器/用户

1. **verdict：pass** — W4 八项成功标准在代码与关键回归上可核对成立；开放 required = 0。
2. 可按 P-003 推进 **关门**（用户确认）；N-001～N-005 记入下波或 accepted-residual，不挡 `done`。
3. 建议 `/govern` 响应时：闭合/记录 recommended、同步 `workspace.md` 波次表、再改 GOAL-005 status + goal-tree。

## 声明

本意见 **source: independent**；**不修改** 目标 `status` / `progress` / goal-tree 状态列。响应与放行由 **`/govern`** 处理。
