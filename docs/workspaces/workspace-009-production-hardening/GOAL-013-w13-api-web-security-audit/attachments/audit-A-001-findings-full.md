---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-013-w13-api-web-security-audit
version: 0.1.0
---

# 附件 · audit-A-001 全文发现与证据（W13 api/web 安全审查）

> 权威摘要与 verdict 见 [03-audit/A-001](../03-audit/A-001-w13-security-review-findings.md)。本附件承载逐条证据摘录、攻击场景与建议修复。行号以 2026-08-26 工作树为准。

## Required

### F-001（P1）公开 invite/accept：bcrypt 前置 + 无限流 → 未认证 CPU DoS

位置：`apps/api/internal/handler/invites.go:293-307`；挂载 `:276-277`。

```go
if err := repo.ValidateNewPassword("", body.Password); err != nil { ... }   // L298 先查密码策略（DB 事务）
hash, herr := auth.HashPassword(body.Password, passwordHashCost)            // L302 再做 bcrypt-10（~60–100ms）
if _, aerr := repo.AcceptInvite(strings.TrimSpace(body.Token), ...)         // L307 最后才校验 token
```

`passwordHashCost = 10`（handler/users.go:29）。该路由直接挂中央 mux，与 login/recovery 不同——无任何 rate limiter。攻击者以垃圾 token + 合法格式密码打 `POST /api/auth/invite/accept`：每请求强制 ≥1 次存储事务 + 一次 bcrypt 后才被（单条索引 SELECT 的）token 查询拒绝；少量线程即可把带宽放大为 ~15× CPU 占用并与真实登录 bcrypt 抢算力。

**修复**：先按 `hashInviteToken(token)` 查邀请，未知/过期/已用/已吊销统一回 `INVITE_INVALID` 再进入密码校验；并照 recovery.go:58 挂 `newLoginRateLimiter(15*time.Minute, 20, 1<<16)`（key = IP|token 或 IP）。

### F-002（P2）/api/mfa/disable、/api/mfa/recovery/rotate 无第二因子失败限流

位置：`apps/api/internal/handler/mfa.go:214-239`（disable）、`:242-263`（rotate）；`modules/mfa/service.go:312-334`（requireActiveSecondFactor）。

```go
plain, fromPrevious := s.decryptSecret(st.SecretCiphertext)
if _, ok := ValidateTotp(plain, code, now, totpWindow, st.LastUsedStep); !ok {
    return ErrMFAInvalid
}
```

无失败计数、无锁定、路由无限流（对照 `/api/auth/mfa/verify` 有 proof 5 次上限 + IP 限流 10/15min，mfa.go:43-47,87-91）。持被盗 access token（默认 15min TTL）者可对 6 位码高频猜测（±1 窗口 ⇒ 瞬时命中 p≈3×10⁻⁶），数百 req/s 约 23 万次猜测即半数概率永久摘除 MFA；rotate 同理换出恢复码集。

**修复**：为两端点加持久化 per-user 失败计数（N 次失败冷却/作废，镜像 proof 上限），并套 IP|userID 的 loginRateLimiter（参照 changePassword）。

### F-003（P2）/api/mfa/enroll currentPassword 无限流密码预言机

位置：`apps/api/internal/handler/mfa.go:166-177`。

```go
if !auth.VerifyPassword(current.PasswordHash, body.CurrentPassword) {
    writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "current password is incorrect")
```

同类面在自助改密已明确加固并注明理由（"wrong currentPassword attempts are brute-force surface with a live access token"，account_self.go:82-86,271-288，limiter 定义 :51 = 5 次/15min）；enroll 无对应限流 → 被盗 token 可无限在线猜密码并烧 bcrypt。

**修复**：复用 passwordLimiter 模型（IP|userID，5/15min），置于 bcrypt 比对之前。

### F-004（P2·bug）MFA Confirm 持久化墙上时钟步进

位置：`apps/api/internal/modules/mfa/service.go:247-253`。

```go
if _, ok := ValidateTotp(plain, code, now, totpWindow, 0); !ok {   // 匹配步进被丢弃
    return ErrMFAInvalid
}
...
if err := s.repo.SetLastUsedStep(userID, now.Unix()/totpPeriodSeconds, now); err != nil {   // 写的是当前墙上步进
```

用户以 offset −1 窗口码完成确认（时钟漂移常见）时，last_used_step 领先其验证器显示码一步；随后首次登录所有候选码 `candidate <= lastUsedStep` 判重放拒绝，需等 30–60s。ValidateTotp 本就返回匹配步进（totp.go:58-74）。

**修复**：接收匹配步进并持久化之。

## Recommended（P3）

### F-005 TOTP 非常数时间比较

`modules/mfa/totp.go:66`：`if want == strings.TrimSpace(code)`。recovery/email 码均用 `subtle.ConstantTimeCompare`（authsession/recovery.go:229、email_identity.go:321）。修复：统一 `subtle.ConstantTimeCompare`。

### F-006 邀请链接 Host/XFP 头注入

`handler/invites.go:156-162`：

```go
scheme := "http"
if proto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); proto != "" { scheme = proto }
return scheme + "://" + r.Host + "/invite/accept?token=" + rawToken
```

代理转发任意 Host 时，邮件中的活体邀请 token 可指向攻击者域名的仿冒接受页。修复：从服务端配置取 canonical base URL（其他外联邮件面的既有口径）。触发需 users.invite 权限，故 P3。

### F-007 账号锁定定向 DoS（处置待裁）

`internal/auth/auth.go:54-57,197-212`；`handler/rate_limit.go:60`。5 连败开 15min 锁仅按账号计（无 IP 维度），锁开还吊销全部会话；登录 limiter 允许同 IP|username 20 次/15min，captcha 默认关闭。已知用户名者可单 IP 无限续锁。可选：提高阈值+指数退避 / 首锁自动启用 captcha / 显式 accepted-residual。

### F-008 recovery EXPIRED 响应状态区分

`handler/recovery.go:169-172`：identifier 能解析且有陈旧挑战时回 `RECOVERY_CODE_EXPIRED`，其余统一 `RECOVERY_CODE_INVALID` —— "近期申请过重置"的部分枚举 oracle（受 20/15min limiter 约束）。另：外部可对已知目标投 5 次垃圾作废其新鲜挑战（void-on-5）。修复：expired 也回 INVALID 并静默删行；void 行为保留但评估配对语义。

### F-009 邮箱绑定新地址无步长限制

`modules/authsession/email_identity.go:166-177`：cooldown 仅覆盖"同一 pending 地址重复请求"；绑新地址即时发信，端点仅 identity 门。滥用面：任意账号把绑定当免费发信原语打任意第三方地址；抢占受害者邮箱 pending 位阻塞真实主人绑定。修复：per-user 发信步长 + 新地址短冷却（或 currentPassword step-up），过期 pending 自动释放。

### F-010 GET /api/schema/{pageId} 未挂认证

`handler/schema.go:24`（挂载点 composition.go:596）：`mux.Handle("GET /api/schema/{pageId}", h.schema())` 无 a.Middleware。全部模块路由工厂均自包 a.Middleware；manifest 公开有注释声明而 schema 无 → 读作遗漏。匿名可枚举管理页 dataSource/action/字段形状/权限键名。修复：包 a.Middleware（或按"预登录页需要"最小公开集），若刻意公开则补注释声明。

### F-011 LIKE 通配符未转义 ×3

`wallet/store/repository.go:153-159`（ListAccounts）、`:584-588`（ListEntries）、`recyclebin/store/repository.go:98-105`。值均参数化（无注入），但 `%`/`_` 保留通配义：`q=%a%a%a…` 强制病态 LIKE 全表扫描（wallet_ledger_entries 为无界增长表）、`q=_` 击穿精确语义逐字符探测 owner/memo。修复：共享 `likePattern(q)` 转义 `\ % _` 并在各 LIKE 后加 `ESCAPE '\'`（SQLite/PG 可移植）。

### F-012 钱包孤儿 owner 账户

`wallet/store/repository.go:271-308`（GetOrCreateUserAccount）；DDL 仅 `UNIQUE(owner_type, owner_id, currency)` 无 FK；`handler/wallet.go:113-187` 直接采信路径 ownerId。一次手滑 typo 即产生永久对账健康却无主用户的账本行。修复：GetOrCreateUserAccount 内（OwnerUser 时）先 SELECT users 存在性，或加 FK。

### F-013 自助 scope TOCTOU（休眠，处置待裁）

`handler/resources.go:703-714`（update）、`:748-761,:787`（delete）、`:872-884,:903`（batch delete）：Go 侧 Get 预检 ownership 后执行裸 id UPDATE/DELETE。ownership/scope 分派在两步之间变更则越权写。v1 未注册生产 scoped resource（resources.go:224-228）故休眠。修复方向：scoped resource 合同落地时把 owner 判据推进语句 WHERE（RowsAffected 判定）。

### F-014 authFetch 缺同源守卫

`web/src/account/auth-client.ts:181-211`：Bearer 无条件附加；X-Refresh-Token 按 `new URL(String(input), origin).pathname === "/api/account/sessions"` 判定——绝对跨源 URL 路径撞名即携带。当前所有构造点强制单斜杠同源相对路径（resource.ts DATASOURCE_URL_PATTERN、request-construction.ts PROTOCOL_URL_RE、app-manifest.ts PATH_PATTERN），不可利用；但 authFetch 是通用传输层，缺最后一道闸。修复：withAuth/authFetch 内 resolve 后比对 `resolved.origin === window.location.origin`，不同源剥离 Authorization/X-Refresh-Token。

### F-015 boot support 动作 url 未校验 scheme

`web/src/host/boot.ts:328-334`；类型允许 url（host/failure.ts:59）；validateFailure（failure.ts:218-229）不验 url：

```ts
case "support": if (action.url) { const link = document.createElement("a"); link.href = action.url; … link.click(); }
```

现无生产者发送带 url 的 support 动作（HostFailure 均本地构造）；首个未来服务器驱动载荷将使其成为 javascript:-URL 执行原语。修复：点击前按 branding.ts isSafeBrandingUrl 同款规则校验（https: 或同源路径）。

### F-016 validateClaim 环死循环

`web/src/protocol/conformance/claim.ts:273-284`：`pending.push(...dependencyEntry.dependsOn)` 无 visited 集合，A→B→A 互依时队列永不清空冻结标签页；兄弟 validateRegistry（claim.ts:185-205）已有 visiting/visited 正确范本。当前固定注册表无环不可达。修复：照抄 visited 守卫。

### F-017 mail 主密钥与密文同目录

`composition/composition.go:713`：`mail.LoadOrCreateMasterKey(cfg.MailConfigMasterKey, filepath.Join(filepath.Dir(cfg.DBPath), "mail-master.key"))`。MAIL_CONFIG_MASTER_KEY 未设（默认，.env.example 标注可选）时 AES-256-GCM 密钥自动生成于 `<db目录>/mail-master.key`；数据目录任何部分泄露（备份/快照/误配挂载）= 密文+密钥同泄，"加密落盘"退化为混淆。密钥文件未入库（.gitignore /data/ 与 *.key 已核实）。修复：生产文档/env 引导独立 keyring 路径；未设 env 且 production 启动时告警。

### F-018 raster immutable 缓存 vs 删除

`handler/raster_assets.go:245-269`（L264 `Cache-Control: public, max-age=31536000, immutable`）；删除 :306-311；公共路由 branding_assets.go:65、account_avatar.go:44。替换/清除头像或 logo 后，见过该 URL 的浏览器/中间缓存继续供给旧资产一年（immutable 不再验证）。URL 以 128-bit 随机 id 为键且公共可取。修复：缩短 TTL（如 max-age=3600 可变），或确证 id 单次使用内容寻址后保留 immutable 并记录残余。

### F-019 导入 CSV 明文密码留存

`handler/import.go:140-193,225,292-301`：模板强制 password 列；CSV 文件经 uploads 通道永久留存（owner 可经 GET /api/files/{id} 再下载，files.read 持有者可经库下载面读取），无导入后清理。MustChangePassword 首登轮换缩小窗口但不消除"导入→首登"间的静态明文凭据。修复：成功导入后自动删除源 CSV（owner 知情）或文档要求手动即刻删除；长期支持预哈希列。

### F-020 nginx HSTS 缺失 + img-src 过宽（受 I-001 约束）

`web/nginx.conf:18`（listen 80）、`:29`（CSP `img-src 'self' data: https:`）；注释明确 HSTS/生产 CDN CSP 属部署层残余。TLS 层无 HSTS 时存在降级面；img-src https: 允许任意远端图（branding/avatarUrl 字段 scheme 已由 branding.ts:36-49 校验，非注入而是追踪/内网探测信道）。修复：TLS 终结层补 Strict-Transport-Security（或本文件起 443 后就地加）；评估收窄 img-src 或远端图同源代理。拓扑未知部分按 I-001 文档化。

## 健壮性/性能 bug

### B-1 registerFailedAttempt 行缺失误判

`authsession/email_identity.go:353-382`：以读回 attempts==0 判"已作废"，并发消费/UPDATE 出错（事务结果被弃 `_ =`）同样得 0 → VerifyEmail 对刚成功的挑战回 EMAIL_CODE_EXPIRED 而非 EMAIL_NOT_PENDING。修复：同一事务内由守卫 UPDATE/delete 判定（范本：ConsumeRecoveryAttempt）。

### B-2 filelibrary 详情整读 body 取 meta

`handler/filelibrary.go:72-93`：`objects.Get(...)` 读全量（至 8MiB）仅为 {name,type,owner}，随后又 Stat 一次。修复：改用 Stat（ObjectInfo 含 Meta/Size），Get 仅留 fallback。

### B-3 全局锁串行化昂贵操作

`upload.go:154-158,361-368`：quotaMu 包住 quotaReached 的 List + 全命名空间逐对象 Stat（O(总对象数) 在锁内）→ 所有用户上传串行且延迟随库存线性劣化；`account_avatar.go:33,70-83`：avatarQuotaMu 跨 multipart 读取 + 解码 + CatmullRom 重采样 + 编码全程持有 → 头像处理全局串行。修复：配额锁按 owner 分片/临界区内移 + 乐观复查；头像锁仅包 CountOwner + save。

### B-4 mail_admin 错误字符串分类

`handler/mail_admin.go:110-113`：`strings.Contains(err.Error(), "retention")` 决定 400/409 分类；文案演进即静默错分。修复：导出 sentinel error + errors.Is。
