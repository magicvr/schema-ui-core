---
id: A-003
goal: GOAL-007-w7-api-web-security-audit
title: W7 close-out independent 复核（A-001 F-001～F-012 闭合）
source: independent
auditor: grok-4.6 · thinking high · /audit skill
date: 2026-08-19
verdict: conditional
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# A-003 · W7 close-out independent 复核（independent）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-4.6 · thinking high · `/audit` skill（本会话；非 A-001 代贴） |
| **类型** | close-out / finding-closure |
| **scope** | A-001 F-001～F-012 required 闭合证据（`apps/api` + `apps/web` 现行实现与相关测试）；对照 E-002 / A-002 |
| **verdict** | **conditional** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 已校验；`shared_materials_catalog: none`） |

## 范围与区间

- **覆盖**：A-001 十二条 required（F-001～F-012）在现行代码中是否真正闭合，而非仅被 E-002 / A-002 声称闭合。核对路径以 E-002 证据表为准，并回读实现与测试。
- **方法**：源码通读 + 关键测试阅读 + 定向回归（见下）。未做动态 exploit / 渗透。
- **不覆盖**：不改 `status` / `progress` / 方案正文 / goal-tree。不把 recommended A-001 F-014～F-016 升格为本波 required。不把本意见当作关门。
- **排除**：GOAL-002/004 已书面接受残余（refresh localStorage、匿名 schema/manifest、Compose 无 TLS、bcrypt cost、data-permission v1 未接线、development JWT、单会话吊销不 bump `token_version`）。

## P-005 / 工作区核对

| 核对项 | 结论 |
|--------|------|
| I-001（finding 清单） | verified；本条不重开 |
| I-002（high 是否暂挂 VP-008 `go`） | verified（D-002：F-001/F-002 闭合前暂挂）。本条确认 **A-001 F-001/F-002 已 genuine fixed**；恢复对外 go 宣称仍须 `/govern` 复核，不得由本意见直接改宣称 |
| 到期 required 信息项 | 无到期未关闭项阻断本 close-out scope（I-002 绑的是宣称门禁，不是本波代码闭合） |
| 共享资料 | 无（`shared_materials_catalog: none`） |
| 工作区绑定 | Root / canonical / `plan_refs`+`primary_plan` 与 `workspace.md` 一致 |

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 两条 high（MFA fail-closed、mfa-reset 管理员边界 + 仅 active 踢会话）已落地 | `modules/mfa/service.go` `Required` / `AdminReset`；`handler/mfa.go` admin reset；`handler/auth.go` BeginChallenge 失败 → 500 不发 token |
| 头像所有权、每用户配额、栅格 2048 边长、enroll step-up、显式反代 CIDR、登录 401 对齐、preview sandbox iframe、refresh 头仅会话列表、上传 quota 串行化 | 见闭合表 |
| 定向 API 测试通过 | `cd apps/api && go test ./internal/modules/mfa ./internal/handler -count=1 -timeout 180s -run "TestServiceAdminReset\|TestServiceLifecycle\|TestMFASelfService\|TestAccountAvatarProfileRejectsAnotherUsersAsset\|TestAccountAvatarPerUserQuota\|TestAccountLockLifecycle\|TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer"` → ok |
| 定向 web 测试通过 | `apps/web` vitest：`download-behavior.test.tsx` / `mfa-manager.test.tsx` / `auth-client.test.ts` → 3 files / 37 tests passed |
| F-013 顺手修复（非本波 required） | `handler/upload.go` `save`：meta 写失败删除对象、返回 error |

未核对 E-002 声称的全量 `go test ./...` / `npm test` 1069；本条不以全量绿作为闭合依据。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 独立意见落盘 | 已达成（前序） | A-001 |
| S2 用户确认修复范围 | 已达成（前序） | D-002 整单采纳 F-001～F-012 |
| S3 按确认范围实施并回归 | **部分** | 11/12 required 可核对闭合；A-001 F-006 声称已修，实现为 no-op |
| S4 独立/cross 复核后关门 | **未达成** | 本条 conditional；存在 1 条 med required 未闭合。本意见不改 status |

## 对照 A-001 required 闭合

| A-001 | 严重度 | 本条判定 | 证据 | 说明 |
|-------|--------|----------|------|------|
| F-001 | high | **fixed** | `apps/api/internal/modules/mfa/service.go` `Required()`：仅 `ErrNotFound` 视为不需要 MFA，其它读错误 `return true`；`auth.go` Login 不发 token；`handler/auth.go` `BeginChallenge` 失败 → 500 `LOGIN_FAILED` | 存储异常不再 fail-open。无专门「GetState 出错」单测（见本条 F-004） |
| F-002 | high | **fixed** | `handler/mfa.go` 非 admin 对 admin 目标 → 403 `ADMIN_ACCOUNT_FORBIDDEN`；`AdminReset` 返回 `removedActive`，仅 true 时 `BumpTokenVersionAndRevokeAll`；无 enrollment → `removed=false`。`service_test.go` `TestServiceAdminReset`；`mfa_test.go` active 路径仍踢会话 | 管理员目标边界与「未开通勿当通用踢会话」均落地。缺「委派 mfa-reset 打 admin」单测（见本条 F-004） |
| F-003 | med | **fixed** | `raster_assets.go` `save` 写 `owner`；`storeUploadForOwner`；`account_self.go` PATCH 校验 `meta["owner"]==user.ID`；`account_avatar.go` `dropPreviousAvatar` 仅删自有。`TestAccountAvatarProfileRejectsAnotherUsersAsset` | 原攻击链（绑他人 URL 再删盘）已断。PATCH `DeleteOrphan` 仍不查 owner，见本条 F-003 recommended |
| F-004 | med | **fixed** | `maxAvatarPerUser=10` + `CountOwner`（不可读 meta 保守计入）；`composition.go` 启动 GC 只保留 users.avatar_url | 「无限孤儿」已有上界 + 启动回收。CountOwner 与 save 未串行化，并发可略超 10，不恢复为 unbounded |
| F-005 | med | **fixed** | `maxRasterInputDimension = 2048`（约 16 MiB 解码预算）；`decodeRasterImage` 先 `DecodeConfig` 再解码 | 相对原 8192（~256 MiB）为实质性降预算 |
| F-006 | med | **仍 open（声称 fixed 不实）** | `handler/captcha.go` 只调用 `captchaGenerateLimiter.allow()`，**从不 `record()`**；`rate_limit.go` `allow()` 在无记录时恒 true，且「never creates a new map entry — only record() registers a key」 | 生成限流是空操作。算术题仍明文 1–50。`captcha_test.go` 无 429 覆盖。见本条 Findings F-001 |
| F-007 | med | **fixed** | `handler/mfa.go` enroll 强制非空 `currentPassword` 并 `VerifyPassword`；已 active 仍 `ErrActive`。`mfa-manager.tsx` 发送 `currentPassword`。`TestMFASelfService` 带密码 enroll；`mfa-manager.test.tsx` 有当前密码步骤 | Bearer  alone 不能绑 TOTP |
| F-008 | med | **fixed** | `SetTrustedProxyCIDRs` 显式 CIDR，默认 `127.0.0.1/8`；`config.default.yaml` `trusted_proxies: ""`；`composition.go` 启动安装；`compose.yaml` API **无** `ports:`。`TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer`：loopback 信 X-Real-IP，公网 peer 不信 | 不再默认信任全部 RFC1918；Compose 不再把 API 端口映射到宿主。测试注释仍写「private peer」、未断言 10/8·172.16/12 不信任。compose 注释示例 CIDR 过宽，见本条 F-002 |
| F-009 | med | **fixed** | `handler/auth.go`：`ErrAccountLocked` / `ErrAccountDisabled` → 401 `UNAUTHORIZED` 同未知/错密。`TestAccountLockLifecycle`、`account_self_test.go` disabled/locked login 401 | 错误码枚举口已关。锁定/禁用仍跳过 bcrypt（未知用户走 dummy hash），时序侧信道见本条 F-005 |
| F-010 | med | **fixed** | `apps/web/src/renderer/render.tsx`：preview 不再 `location.replace(blob)`；`iframe sandbox=""`（无 allow-scripts / allow-same-origin）。`download-behavior.test.tsx` 断言 sandbox iframe 且 `location.replace` 未被调用 | 预览页不再把 blob 当顶层同源文档。父 tab 仍是同源 `about:blank` 且 blob URL 存活 60s，属残余，不重开 required |
| F-011 | med | **fixed** | `auth-client.ts` `withAuth`：仅 `pathname === "/api/account/sessions"` 时设 `X-Refresh-Token`；refresh 走 POST body | 常规 `authFetch` 不再附带长寿命 refresh。测试只断言 Bearer，未断言 refresh 头缺席（见本条 F-004） |
| F-012 | med | **fixed** | `upload.go` `quotaMu` 包住 `quotaReached` + `save` | 单进程 TOCTOU 已关。无并发单测 |

## Findings

### F-001 · A-001 F-006 验证码生成限流未生效（allow 不记账）

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `apps/api/internal/handler/captcha.go` 36–60；`apps/api/internal/handler/rate_limit.go` 41–47、90–107；`apps/api/internal/modules/logincaptcha/challenge.go` 42–65；`apps/api/internal/handler/captcha_test.go`（无生成 429 用例） |

A-001 必改项为「验证码不可机读**或** generate 限流」。E-002 / A-002 选择后者并标 `fixed`。实现新建 `captchaGenerateLimiter = newLoginRateLimiter(time.Minute, 10, 1<<16)`，但登录限流器的契约是：**`allow()` 只读、不建条目；只有 `record()` 才计数**。验证码生成路径在 `allow()` 为 true 后直接 `Generate()`，从不 `record()`。因此任意客户端可无限 `GET /api/auth/captcha`：题目仍是明文 1–50 加减，表 `captcha_challenges` 仍可被灌满。这是关闭声明不实，不是「限流偏松」。

闭合要求：生成路径必须对每次成功/尝试 `record()`（或换用按 allow 即计数的限流器），并补「第 11 次 → 429 RATE_LIMITED」测试。不可机读题目可作为替代闭合路径（A-001 原文 OR）。

### F-002 · Compose 文档示例把信任面扩回整段 Docker 网

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `compose.yaml` 44–48；对比 `rate_limit.go` 默认仅 loopback |

默认实现正确（不发布 API 端口 + 空 `trusted_proxies` → 只信 loopback）。注释却建议 `HTTP_TRUSTED_PROXIES=172.16.0.0/12` 以「启用 per-client 限流」。该 CIDR 覆盖默认 Docker 网桥，其它容器可再伪造 `X-Real-IP`，部分撤回 A-001 F-008。应改为 compose 网络的**具体**网段，或 nginx 容器 IP `/32`。

### F-003 · 头像 PATCH 清理仍 `DeleteOrphan` 不校验 owner

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `handler/account_self.go` `DeleteOrphan(oldAvatar)`；对比 `dropPreviousAvatar` 已查 owner |

新绑定已强制 owner，原 IDOR 链已断。若历史 profile 已写入他人资产 URL，清空/替换仍会删盘。防御深度：`DeleteOrphan` 应对齐 `dropPreviousAvatar`。

### F-004 · 若干闭合路径缺少对准原 finding 的回归测试

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | 无 `Required()` 存储错误测试；`mfa_test.go` 无委派账号重置 admin；无 captcha 生成 429；`TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer` 未覆盖 RFC1918；`auth-client.test.ts` 未断言非 sessions 请求不含 `X-Refresh-Token` |

现有测试覆盖了多数修复的 happy/主路径，但未锁住本波最容易回退的失败模式。F-006 正是因此漏过。

### F-005 · 锁定/禁用登录仍跳过密码哈希（时序枚举残余）

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `apps/api/internal/auth/auth.go` 158–173 vs 未知用户 dummy bcrypt |

A-001 F-009 的错误码枚举已闭合。锁定/禁用仍在验密前返回（handler 映射为同一 401），未知用户仍烧 dummy bcrypt。远程时序仍可能区分「存在且锁定/禁用」与「不存在」。非本波必改。

## 必改项汇总（required · 仍 open）

1. **F-001**：使 `GET /api/auth/captcha` 生成限流真正计数（`record()` 或等价实现）并补回归；或改为不可机读题目。A-001 F-006 不得标 `fixed` 直至可重复核对。

无 high required 未闭合。

## 仍开放但不在本波 required 范围

A-001 F-014 / F-015 / F-016（SQL 优先级、otpauth 转义、refresh family 不杀）保持 recommended；D-002 未纳入 S3/S4。F-013 已顺手修。

## 与既有意见的异同

| 来源 | 异同 |
|------|------|
| A-001 independent fail | 十二条 required 当时均 open。本条确认其中 11 条已 genuine fixed；**驳回** A-002 对 F-006 的 `fixed` |
| A-002 self pass | 实施范围声称开放 required = 0。本条不同意：F-006 关闭证据不可复核（限流器未记账）。高危两条同意闭合 |
| E-002 | 证据路径大体可指回真实改动；F-006 行「按真实客户端 IP 限流（10/min）」与代码行为不符 |

无 P-004 冲突需用户在「是否 required」上二选一：F-006 仍是 A-001/D-002 范围内的 required，本条未降级。

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional** — 两条 high 已闭合；11/12 required 可核对；A-001 F-006 关闭声明不实，仍为 med required。不得将 GOAL-007 标 `done`，不得把 A-002 的「开放 required = 0」当作 S4 依据。

建议 `/govern`：

1. 响应本条 F-001：修补验证码生成限流（或改不可机读）并补 429 测试，再请独立复核该单条。
2. I-002：A-001 F-001/F-002 已闭合，可就是否恢复 VP-008 `go` 宣称做书面复核；恢复前仍应披露本条 F-006 未闭合。
3. recommended F-002～F-005 可同波顺手或明确 residual。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
