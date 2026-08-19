---
id: A-001
goal: GOAL-007-w7-api-web-security-audit
title: W7 api/web 独立代码审计（bug 与安全漏洞）
source: independent
auditor: grok-4.6 · thinking high · 本会话独立代码审计（用户指令不加载 skills；P-003 代贴）
date: 2026-08-19
verdict: fail
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# A-001 · W7 api/web 独立代码审计（independent）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-4.6 · thinking high · 本会话独立代码审计（用户明确不加载 skills；意见按 P-003 代贴，非 `/audit` skill 入口冒充编排器） |
| **类型** | ad-hoc / execution-facts（当前实现静态审查） |
| **scope** | `apps/api` + `apps/web` 当前实现的 bug 与安全漏洞（鉴权、会话、MFA、验证码、上传/头像/品牌资源、RBAC/IDOR、钱包、SQL、XSS/CSP、CORS） |
| **verdict** | **fail** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 已校验；`shared_materials_catalog: none`） |

## 范围与区间

- **覆盖**：`apps/api/internal/{auth,handler,composition,config,server,modules/*}` 安全敏感路径；`apps/web/src/{account,renderer,host,app}`；`apps/web/nginx.conf`；`compose.yaml`。
- **方法**：源码通读 + 并行核对鉴权/上传/SQL·RBAC/前端 XSS；关键 finding 回读文件与行号。
- **不覆盖**：未跑动态 exploit / 渗透；未改代码；未把本意见当作修复完成。
- **排除**：已在 GOAL-002/004 书面接受的残余（见「已知残余」），不重开为 required。

## 成果（有证据 · 基线已加固）

下列主张来自本次核对，**不是**本波修复成果。

| 主张 | 证据 |
|------|------|
| JWT 拒绝非 HMAC；v5 默认不容 `alg=none` | `apps/api/internal/auth/auth.go` `ParseAccessToken` keyfunc |
| Refresh 轮换原子：`UPDATE … WHERE revoked_at IS NULL` | `apps/api/internal/auth/auth.go` Refresh；`authsession/accounts.go` |
| 密码变更 bump `token_version` 并吊销 refresh | users 仓库 UpdateUser 路径 |
| 通用上传：服务端嗅探 + 主动内容标记 + `attachment` + sandbox；`GET /api/files/{id}` 属主校验 + hex id | `apps/api/internal/handler/upload.go` |
| 品牌/头像栅格服务端重编码 JPEG/PNG | `apps/api/internal/handler/raster_assets.go` |
| 角色委派：`roles.assign` + 权限子集 + 仅 admin 赋 admin | `apps/api/internal/handler/users.go` `authorizeRoleAssignment` |
| 钱包 apply 表拒绝负余额；自助接口用会话身份 | `wallet/store/repository.go` Apply；`handler/wallet_self.go` |
| 反应表达式白名单，无 `eval` | `apps/web/src/renderer/reaction-expression.ts` |
| 导航 / branding / dataSource 拒 `javascript:` 与 `//host` | `branding.ts`、`resource.ts`、`render.tsx` navigate 正则 |
| 无 Cookie 会话；CORS 默认关闭且无 `Allow-Credentials` | `server/server.go` |
| nginx 基线：`X-Frame-Options DENY`、`frame-ancestors 'none'`、`script-src 'self'` | `apps/web/nginx.conf` |
| 生产启动：`APP_ENV` 必填；非 development 禁 dev-session、要求强 JWT | `config.go` ValidateProd；`cmd/server/main.go` |
| 2026-08-10 高危项（上传 XSS、refresh 双发、默认 development）在当前代码中已对上修复 | 对照 `raw/audit-20260810-api-web-bug-review.md` C1–C3 与现行实现 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 独立意见落盘 | 本条即为 | 本文件 + 索引 |
| S2 用户确认修复范围 | 未开始 | P-004 待确认 |
| S3 实施与回归 | 未开始 | — |
| S4 开放 required = 0 | 未达成 | 本条 12 条 required 均 open |

## Findings

### F-001 · MFA 登录门在存储出错时 fail-open

| 字段 | 值 |
|------|-----|
| 严重度 | high |
| 建议 | required |
| 状态 | open |
| evidence | `apps/api/internal/modules/mfa/service.go` 75–86；`apps/api/internal/auth/auth.go` 196–201 |

`Required()` 把任何 `GetState` 错误（含 SQLite busy / 读失败）当成「不需要 MFA」。密码通过后 `Login` 直接 `issue()`。`ErrNotFound` 跳过 MFA 正确；存储异常不是。已知密码的攻击者在 `user_mfa` 短暂不可读时可拿到 access/refresh。

### F-002 · `users.mfa-reset` 无管理员目标边界，可拆管理员 2FA 并踢全会话

| 字段 | 值 |
|------|-----|
| 严重度 | high |
| 建议 | required |
| 状态 | open |
| evidence | `apps/api/internal/handler/mfa.go` 201–217；`modules/mfa/service.go` AdminReset；对比 `handler/users.go` 336–362 `authorizeAdminTargetBoundary` |

密码重置禁止非 admin 改 admin。MFA 重置无对等检查：`DELETE FROM user_mfa` 后无条件 `BumpTokenVersionAndRevokeAll`。持 `users.mfa-reset`、无 `admin` 角色的委派账号可对管理员拆 2FA + 强制全会话下线。对未开通 MFA 的用户，该接口仍是通用强制下线。

### F-003 · 头像 URL IDOR：绑定他人资源后可删磁盘文件

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `handler/account_self.go` 150–188；`handler/account_avatar.go` 72–84；公开 GET `account_avatar.go` 37 |

`PATCH /api/account/profile` 只校验 `/api/account/avatars/{32hex}` 形状，不校验上传者。随后清空或重传会 `Delete`/`dropPreviousAvatar` 该 id。`GET` 公开可读。用户列表 API 当前不返回 `avatarUrl`，利用取决于 URL 泄漏（UI / 审计 / 分享）。

### F-004 · 任意登录用户可无配额刷头像（孤儿文件无限涨）

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `handler/account_avatar.go` 47–68；`composition.go` 头像「无启动 GC」注释 |

`POST /api/account/avatar` 无 `files.write`、无配额、无限流。每次新 id；清理只删 profile 上旧 URL。循环上传且不 PATCH 留下孤儿。品牌资源有启动 GC，头像没有。

### F-005 · 栅格解码允许约 256 MiB/请求

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `handler/raster_assets.go` `maxRasterInputDimension = 8192`；decode 后 `image.NewRGBA` |

8192×8192×4 ≈ 256 MiB，再缩到 256/512px。头像对每个登录用户开放。小体积、大 IHDR 的 PNG 可压进 4 MiB 上限。并发可打内存。

### F-006 · 登录验证码可被机器秒解，生成接口无计量

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `modules/logincaptcha/challenge.go` 42–66；`handler/captcha.go` 39–59；`handler/auth.go` 113–121 |

默认关闭。一旦打开：`GET /api/auth/captcha` 匿名无限流；题目为 1–50 加减明文；验证失败不计入登录失败桶。机器人可解。还可灌满 `captcha_challenges`。

### F-007 · MFA enroll/confirm 无 step-up

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `handler/mfa.go` 108–145；disable 路径 147–174 需要第二因子 |

仅 Bearer 即可给**尚未开通 MFA** 的账号绑定攻击者 TOTP。会话窃取放大：受害者下次登录必须过攻击者第二因子。

### F-008 · `X-Real-IP` 信任全部 RFC1918，且 Compose 映射 API 端口

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `handler/rate_limit.go` 125–149；`compose.yaml` 44–45 发布 `:25080` |

W3 设计为「私网 peer 才信 X-Real-IP」。本次确认：全部 RFC1918 都算可信反代，且 API 端口对主机/局域网暴露。Docker 网 / 另一容器 / 局域网可轮换 `X-Real-IP` 打散 `ip\|username` 限流。账号锁（5 次）对已存在用户仍有效。

### F-009 · 登录错误码泄露账号状态

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `auth.go` 165–173；`handler/auth.go` 124–141 |

未知用户/错密：401 + dummy bcrypt。禁用：403 `ACCOUNT_DISABLED`（不验密）。锁定：423 `ACCOUNT_LOCKED`（不验密）。5 次错密可区分「存在且锁定」vs「不存在」。产品上 423 服务 `HOST_ACCOUNT_LOCKED`，仍构成枚举口。

### F-010 · `library.preview` 用 `blob:` 打开文件，拆掉服务端 attachment / CSP sandbox

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `apps/web/src/renderer/render.tsx` 349–371；对比 `handler/filelibrary.go` 261–267、`upload.go` 87–114 |

API 下载加 `Content-Disposition: attachment`、`nosniff`、`CSP: sandbox`。前端 `createObjectURL` + `location.replace` 后这些头不再生效，`blob:` 也不继承 nginx CSP。当前上传拒 HTML/SVG/`<script>`，完整利用还需过滤绕过；过滤一旦放宽即为同源脚本读 refresh token。

### F-011 · 每个 `authFetch` 都带长寿命 `X-Refresh-Token`

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `apps/web/src/account/auth-client.ts` 174–187；API 仅 `account_self.go` 338–347 用于标记当前会话 |

HAR、反代、APM、扩展可在常规 GET 上看到 refresh。CORS 在配置 origins 时亦允许该头。

### F-012 · 通用上传配额 TOCTOU

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `handler/upload.go` 323–329 `quotaReached` 后 `save`，无按 owner 串行化 |

持 `files.write` 并发上传可按并发度突破 1000 文件 / 256 MiB。默认该键仅 admin。

### F-013 · 对象写入成功、meta 失败仍 200

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `handler/upload.go` 130–138 |

无 meta 的对象不计入配额、GET fail-closed。磁盘残留。

### F-014 · 定时任务列表 `q` + `enabled` SQL 优先级

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `modules/scheduledtasks/store/repository.go` 83–94 |

`WHERE instr… OR instr… AND enabled = ?` 中 AND 绑定更紧，key 匹配可绕过 enabled 过滤。仍需 `tasks.read`。

### F-015 · `otpauth://` URL 只转义空格和冒号

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `modules/mfa/totp.go` 76–84 |

显示名中的 `&` / `?` 可能破坏 URI。

### F-016 · Refresh 重放不杀整个 token family

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `auth.go` Refresh 原子吊销单条；不吊销该用户其余 live refresh |

轮换本身正确。偷到的 refresh 先用一次，攻击者获 30 天新会话。

## 已知残余（本波不重开 required）

| 项 | 出处 |
|----|------|
| refresh 在 localStorage | GOAL-004 范围外 / 产品 D-002 |
| schema/manifest 匿名可读 | GOAL-002 D3 accepted-residual |
| Compose 无 TLS | VP-009 / compose 非目标 |
| bcrypt cost 10 | 残余 |
| data-permission v1 未接生产资源 | 产品 D-002 |
| development 空 JWT → 公开开发密钥 | 生产 ValidateProd fail-closed |
| 单条会话吊销不 bump `token_version` | JWT 15m TTL 设计 |

## 必改项汇总（required · 均 open）

1. F-001 MFA `Required()` 存储错误 fail-closed  
2. F-002 MFA admin reset 管理员目标边界 + 未开通勿无条件踢会话  
3. F-003 头像 URL 所有权 / 引用计数删除  
4. F-004 头像配额或覆盖写 + 未引用 GC  
5. F-005 降低栅格解码边长 / 像素预算  
6. F-006 验证码不可机读或 generate 限流  
7. F-007 MFA enroll 要求当前密码（已开通则要 TOTP）  
8. F-008 `X-Real-IP` 改为显式反代 CIDR；勿把 API 端口暴露给非反代网络  
9. F-009 锁定/禁用登录响应与 401 对齐，或接受枚举并书面 residual  
10. F-010 预览勿裸 `blob:` 导航  
11. F-011 `X-Refresh-Token` 仅挂会话列表请求  
12. F-012 上传配额与 save 串行化  

## 结论 + 建议下一步

**verdict: fail** — scope 内 2 条 high required 未闭合，不得将本波标 `done`，不得假装共享基架无开放高危。

建议 `/govern` 下一步：用户确认 S2 修复范围（整单采纳 F-001～F-012，或逐条 residual/overruled），并裁决 I-002（是否暂挂 VP-008 `go`）。确认前不实施。
