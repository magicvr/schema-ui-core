---
id: A-004
goal: GOAL-007-w7-api-web-security-audit
title: W7 close-out independent 复核（A-003 F-006 修正后）
source: independent
auditor: grok-4.6 · thinking high · /audit skill
date: 2026-08-19
verdict: pass
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# A-004 · W7 close-out independent 复核（independent）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-4.6 · thinking high · `/audit` skill（本会话；A-003 后续 close-out，非编排器） |
| **类型** | close-out / finding-closure |
| **scope** | A-001 F-001～F-012 required 闭合证据；重点复核 A-003 F-001（A-001 F-006 captcha generate limiter 未 `record()`）经 E-003 后是否 genuine fixed；对其余 11 条做现行代码抽查 |
| **verdict** | **pass** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 已校验；`shared_materials_catalog: none`） |

## 范围与区间

- **覆盖**：A-001 十二条 required（F-001～F-012）在现行代码中是否闭合。A-003 判定 11/12 fixed、F-006 关闭声明不实；本条核对 E-003 修正（`captcha.go` `record()` + `TestCaptchaPreflightRateLimited`）并回读其余闭合路径。
- **方法**：源码通读 + 关键测试阅读 + 定向回归（见下）。未做动态 exploit / 渗透。
- **不覆盖**：不改 `status` / `progress` / 方案正文 / goal-tree。不把 A-001 recommended F-014～F-016 升格为本波 required。不把本意见当作已关门。
- **排除**：GOAL-002/004 已书面接受残余（refresh localStorage、匿名 schema/manifest、Compose 无 TLS、bcrypt cost、data-permission v1 未接线、development JWT、单会话吊销不 bump `token_version`）。

## P-005 / 工作区核对

| 核对项 | 结论 |
|--------|------|
| I-001（finding 清单） | verified；本条不重开 |
| I-002（high 是否暂挂 VP-008 `go`） | verified（D-002：F-001/F-002 闭合前暂挂）。本条确认 **A-001 F-001/F-002 仍 genuine fixed**。恢复对外 go 宣称仍须 `/govern` 书面复核，不得由本意见直接改宣称 |
| 到期 required 信息项 | 无到期未关闭项阻断本 close-out scope（I-002 绑的是宣称门禁，不是本波代码闭合） |
| 共享资料 | 无（`shared_materials_catalog: none`） |
| 工作区绑定 | Root / canonical / `plan_refs`+`primary_plan` 与 `workspace.md` 一致 |

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| A-003 F-001 / A-001 F-006：生成路径在 `allow()` 通过后调用 `record()`，滑动窗口真正计数 | `apps/api/internal/handler/captcha.go` 57–65：`allow()` 拒绝 → 429 `RATE_LIMITED`；通过后 `captchaGenerateLimiter.record(loginClientIP(r), …)` 再 `Generate()` |
| 第 11 次匿名生成 → 429 已锁回归 | `apps/api/internal/handler/captcha_test.go` `TestCaptchaPreflightRateLimited`：独立 limiter、10 次 200 后第 11 次 429 + `error=RATE_LIMITED` |
| 其余 11 条 required 现行实现仍与 A-003 闭合判定一致 | 见下表；抽查路径未回退 |
| 定向 API 测试通过 | `cd apps/api && go test ./internal/modules/mfa ./internal/handler -count=1 -timeout 180s -run "TestServiceAdminReset\|TestServiceLifecycle\|TestMFASelfService\|TestAccountAvatarProfileRejectsAnotherUsersAsset\|TestAccountAvatarPerUserQuota\|TestAccountLockLifecycle\|TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer\|TestCaptchaPreflightRateLimited"` → ok |
| 定向 web 测试通过 | `apps/web` vitest：`download-behavior.test.tsx` / `mfa-manager.test.tsx` / `auth-client.test.ts` → 3 files / 37 tests passed |

未重跑 E-002 声称的全量 `go test ./...` / `npm test` 1069；本条不以全量绿作为闭合依据。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 独立意见落盘 | 已达成（前序） | A-001 |
| S2 用户确认修复范围 | 已达成（前序） | D-002 整单采纳 F-001～F-012 |
| S3 按确认范围实施并回归 | 已达成 | E-002 + E-003；12/12 required 可核对闭合 |
| S4 独立/cross 复核后关门 | **本条判定代码闭合条件已满足** | 本条 pass：A-001 F-001～F-012 无未闭合 required。本意见不改 status；关门由 `/govern` 响应后执行 |

## 对照 A-001 required 闭合（F-001～F-012）

| A-001 | 严重度 | 本条判定 | 证据 | 说明 |
|-------|--------|----------|------|------|
| F-001 | high | **fixed** | `apps/api/internal/modules/mfa/service.go` `Required()`：仅 `ErrNotFound` 视为不需要 MFA，其它读错误 `return true`；`handler/auth.go` `BeginChallenge` 失败 → 500 `LOGIN_FAILED` 不发 token | 存储异常不再 fail-open。无专门「GetState 出错」单测（A-003 recommended F-004 残余） |
| F-002 | high | **fixed** | `handler/mfa.go` 非 admin 对 admin 目标 → 403 `ADMIN_ACCOUNT_FORBIDDEN`；`AdminReset` 返回 `removedActive`，仅 true 时 `BumpTokenVersionAndRevokeAll`。`TestServiceAdminReset`：active → `removed=true`；无 enrollment → `removed=false` | 管理员目标边界与「未开通勿当通用踢会话」均仍落地 |
| F-003 | med | **fixed** | `raster_assets.go` 写 `owner`；`account_self.go` PATCH 校验 `meta["owner"]==user.ID`；`dropPreviousAvatar` 仅删自有。`TestAccountAvatarProfileRejectsAnotherUsersAsset` | 原 IDOR 链仍断。PATCH `DeleteOrphan` 仍不查 owner（A-003 recommended F-003） |
| F-004 | med | **fixed** | `maxAvatarPerUser=10` + `CountOwner`；`composition.go` 启动 GC 只保留 users.avatar_url。`TestAccountAvatarPerUserQuota` | 无限孤儿仍有上界 + 启动回收。CountOwner 与 save 未串行化，并发可略超 10，不恢复 unbounded |
| F-005 | med | **fixed** | `maxRasterInputDimension = 2048`（约 16 MiB）；`decodeRasterImage` 先 `DecodeConfig` | 相对原 8192（~256 MiB）预算仍压住 |
| F-006 | med | **fixed**（A-003 当时 open；E-003 后闭合） | `captcha.go`：`allow()` 后 **`record()`**；限流器契约仍是「`allow()` 只读、`record()` 才建条目」（`rate_limit.go` 41–47、90–107）。`TestCaptchaPreflightRateLimited` 本会话通过 | A-001 闭合路径为「不可机读 **或** generate 限流」。本波选后者且现已真正计数（10/min/IP）。题目仍为明文 1–50 加减，属该 OR 路径残余，不重开 required。`allow`/`record` 分锁，并发可略超 10，不恢复为无限流 |
| F-007 | med | **fixed** | `handler/mfa.go` enroll 强制非空 `currentPassword` 并 `VerifyPassword`；`mfa-manager.tsx` 发送 `currentPassword`。`TestMFASelfService`；`mfa-manager.test.tsx` | Bearer  alone 不能绑 TOTP |
| F-008 | med | **fixed** | `SetTrustedProxyCIDRs` 显式 CIDR，默认 `127.0.0.1/8`；`composition.go` 启动安装；`compose.yaml` API **无** `ports:`。`TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer` | 不再默认信任全部 RFC1918。compose 注释示例 `172.16.0.0/12` 仍过宽（A-003 recommended F-002） |
| F-009 | med | **fixed** | `handler/auth.go`：`ErrAccountLocked` / `ErrAccountDisabled` → 401 `UNAUTHORIZED`。`TestAccountLockLifecycle`、`account_self_test.go` | 错误码枚举口仍关。锁定/禁用仍跳过 bcrypt（A-003 recommended F-005 时序残余） |
| F-010 | med | **fixed** | `render.tsx` preview：`iframe sandbox=""`（无 allow-scripts / allow-same-origin），不再 `location.replace(blob)`。`download-behavior.test.tsx` | 预览页不再把 blob 当顶层同源文档 |
| F-011 | med | **fixed** | `auth-client.ts` `withAuth`：仅 `pathname === "/api/account/sessions"` 时设 `X-Refresh-Token` | 常规 `authFetch` 不再附带长寿命 refresh。测试仍未断言非 sessions 请求缺该头（A-003 recommended F-004） |
| F-012 | med | **fixed** | `upload.go` `quotaMu` 包住 `quotaReached` + `save`；meta 写失败删对象（F-013 顺手） | 单进程 TOCTOU 仍关 |

## Findings

本条 **无新 required**。A-003 的必改项（当时 F-001 = A-001 F-006）已按 `fixed` 可核对闭合。

A-003 recommended 仍开放（本波不升格 required，不阻断 S4 代码闭合）：

| A-003 | 严重度 | 建议 | 状态 | 备注 |
|-------|--------|------|------|------|
| F-002 · Compose 文档示例 CIDR 过宽 | low | recommended | open | `compose.yaml` 44–48 仍写 `HTTP_TRUSTED_PROXIES=172.16.0.0/12` |
| F-003 · PATCH 清理 `DeleteOrphan` 不校验 owner | low | recommended | open | `account_self.go` 196；新绑定已强制 owner，原 IDOR 链已断 |
| F-004 · 若干失败模式缺回归 | low | recommended | open（部分缓解） | captcha 429 已补；仍无 `Required()` 存储错误、委派 mfa-reset 打 admin、RFC1918 不信任、非 sessions 不含 `X-Refresh-Token` 断言 |
| F-005 · 锁定/禁用跳过 bcrypt 时序残余 | low | recommended | open | 错误码枚举已关 |

## 必改项汇总（required · 仍 open）

无。A-001 F-001～F-012 均为 `fixed`。A-003 required（captcha generate 未记账）已闭合。

无 high required 未闭合。

## 仍开放但不在本波 required 范围

A-001 F-014 / F-015 / F-016（SQL 优先级、otpauth 转义、refresh family 不杀）保持 recommended；D-002 未纳入 S3/S4。F-013 已顺手修。

## 与既有意见的异同

| 来源 | 异同 |
|------|------|
| A-001 independent fail | 十二条 required 当时均 open。本条确认 12/12 现已 genuine fixed |
| A-002 self pass | 当时声称开放 required = 0，A-003 驳回 F-006。E-003 后本条同意：12 条 required 现可复核闭合。不采纳 A-002 当时对 F-006 的提前 `fixed` 作为历史证据，以 E-003 + 现行代码为准 |
| A-003 independent conditional | 11/12 fixed、F-006 关闭声明不实。本条 **同意** 当时判定；**确认** E-003 已按 A-003 闭合要求补 `record()` 与 429 回归。recommended F-002～F-005 仍开放，不升格 |
| E-003 | 证据路径与代码一致：`allow()` 后 `record()`；`TestCaptchaPreflightRateLimited` 本会话通过 |

无 P-004 冲突：A-003 未降级 F-006；本条亦未降级，而是核对修复后标 `fixed`。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass** — A-001 十二条 required（含 2 条 high）均可核对闭合；A-003 唯一开放 required（captcha generate 限流未记账）已由 E-003 闭合。scope 内无未关闭 high/med required，无到期阻断本 close-out 的 required 信息项。

建议 `/govern`：

1. 响应本条：将 A-001 F-006 / A-003 F-001 标 `fixed`（证据：E-003 + 本条）。
2. 按 D-001/D-002 推进 S4 关门（本意见不改 `status`）。
3. I-002：F-001/F-002 仍闭合，可书面复核是否恢复 VP-008 `go` 宣称；F-006 已不再构成披露缺口。
4. A-003 recommended F-002～F-005 可同波 residual 或后续波次，不阻断本波关门。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
