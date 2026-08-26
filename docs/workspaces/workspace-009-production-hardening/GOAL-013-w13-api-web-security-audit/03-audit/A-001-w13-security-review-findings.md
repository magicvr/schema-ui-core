---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-013-w13-api-web-security-audit
version: 0.1.0
---

# A-001 · W13 api/web 全量安全审查（发现台账）

- **source**: independent（4 个隔离上下文并行审查：认证会话/MFA、数据层/钱包/SQL、文件上传/对象存储、前端 XSS/令牌；另经编排会话对核心链路逐项复核。方法说明：非项目默认 grok-build 路径，属用户指令的多上下文审查；关门 independent 审计仍按项目默认在 S6 另行执行）
- **日期**: 2026-08-26
- **scope**: apps/api（auth/session/MFA/captcha/限流/邀请/恢复、store/kernel/wallet/recyclebin/settings 持久层、upload/objectstore/import/mail 密钥、composition/config/server）+ apps/web（account 令牌与传输层、host boot/failure/claim、renderer 渲染与反应引擎、protocol 构造、nginx/CSP）
- **verdict**: conditional（无 P0、无可直接利用高危；1×P1 + 3×P2 required 待修复；P3 与健壮性发现已由 D-001 纳入本波分母）

## 发现清单（编号为本波分母；全文证据见 [附件](../attachments/audit-A-001-findings-full.md)）

### Required（必修）

| 编号 | 级别 | 标题 | 位置 |
|------|------|------|------|
| F-001 | P1 | 公开 invite/accept 先 bcrypt 后验 token 且无限流 → 未认证 CPU DoS | handler/invites.go:276,293-307 |
| F-002 | P2 | /api/mfa/disable、/api/mfa/recovery/rotate 无第二因子失败限流 → TOTP 猜测预言机 | handler/mfa.go:214-263; modules/mfa/service.go:312-334 |
| F-003 | P2 | /api/mfa/enroll currentPassword 无限流密码预言机（对照 account_self passwordLimiter） | handler/mfa.go:166-177 vs handler/account_self.go:51 |
| F-004 | P2·bug | MFA Confirm 持久化墙上时钟步进而非匹配步进 → 确认后首登被误拒 30–60s | modules/mfa/service.go:247-251 |

### Recommended（P3 加固，D-001 已裁全量纳入）

| 编号 | 标题 | 位置 |
|------|------|------|
| F-005 | TOTP 比较非常数时间（recovery/email 已用 subtle） | modules/mfa/totp.go:66 |
| F-006 | 邀请链接由 Host/X-Forwarded-Proto 头拼出 → 链接注入钓鱼面 | handler/invites.go:156-162 |
| F-007 | 账号锁定可作定向 DoS（5 连败锁 15min 仅按账号计）——处置待裁 | internal/auth/auth.go:54-57,197-212 |
| F-008 | recovery 过期码响应可区分"近期申请过重置"状态 | handler/recovery.go:169-172 |
| F-009 | 邮箱绑定新地址即时发信无步长限制 + pending 位抢占 | modules/authsession/email_identity.go:166-177 |
| F-010 | GET /api/schema/{pageId} 未挂认证（匿名可枚举管理页 schema 元数据） | handler/schema.go:24（挂载 composition.go:596） |
| F-011 | 列表 q 的 LIKE 通配符未转义（wallet×2 / recyclebin×1）→ 扫描放大+逐字符探测 | wallet/store/repository.go:158,584-588; recyclebin/store/repository.go:98-105 |
| F-012 | 钱包账户可为不存在的 owner 创建（无 FK/存在性检查）→ 孤儿账本 | wallet/store/repository.go:271-308 |
| F-013 | 自助行级 scope 为 Go 侧预检 + 无谓词 UPDATE（TOCTOU，当前休眠）——处置待裁 | handler/resources.go:703-761 |
| F-014 | authFetch 仅按 pathname 决定附带令牌头，缺同源守卫（当前调用点安全） | web/src/account/auth-client.ts:181-211 |
| F-015 | boot recovery "support" 动作点击未校验 scheme 的 url（潜伏 javascript: 执行点） | web/src/host/boot.ts:328-334 |
| F-016 | validateClaim 依赖闭包遍历遇环注册表死循环 | web/src/protocol/conformance/claim.ts:273-284 |
| F-017 | mail 主密钥默认落在被加密数据同目录（备份/快照同泄） | composition/composition.go:713; mail/secrets.go:31-52 |
| F-018 | 头像/品牌图删除后因 immutable 一年缓存仍可取 | handler/raster_assets.go:264 |
| F-019 | 用户导入 CSV 含明文密码且导入后永久留存可再下载 | handler/import.go:140-193; filelibrary.go:224-258 |
| F-020 | nginx 无 HSTS + img-src https: 过宽（受 I-001 拓扑约束按可行范围实施） | web/nginx.conf:18,29 |

### 健壮性/性能 bug（一并纳入）

| 编号 | 标题 | 位置 |
|------|------|------|
| B-1 | registerFailedAttempt 以"行不存在"误判"已作废"，并发下提示错码 | authsession/email_identity.go:353-382 |
| B-2 | 文件库详情读整份 body 仅为取 meta，且多跑一次 Stat | handler/filelibrary.go:72-93 |
| B-3 | 全局互斥锁串行化所有用户的上传配额扫描与头像图像处理 | handler/upload.go:154-158; account_avatar.go:33 |
| B-4 | mail_admin 以错误字符串含 "retention" 分类错误类型（应 sentinel + errors.Is） | handler/mail_admin.go:110-113 |

### Informational / 既有残余（记录，不在本波修复分母内）

- refresh token 存 localStorage：tokens.ts 头注明的既有书面接受残余（access 仅内存 + 每次轮换 + 服务端吊销）；维持 D-002 接受。
- 算术 captcha 无真实机器人抵抗价值：定位为摩擦，不得作为 F-007 的对冲依赖。
- GET /healthz 暴露 version/commit：行业常规，接受。
- `apps/api/configs/.env` 存有本地真实 PG 密码与 Resend API key：已核实被 gitignore 且未入库（非仓库泄露）；建议用户轮换该 Resend key（用户侧动作，登记移交 Root）。

## 结论

整体工程质量显著高于同类内部平台基线：认证核心（bcrypt+防枚举时序、刷新令牌原子轮换、token_version 吊销）、持久层（全参数化 + 排序白名单 + 事务内乐观锁 + CHECK/幂等约束）、上传管线（服务端嗅探 + 硬拒绝活动内容 + owner-only + 穿越 白名单）、前端（零 HTML 注入汇点、白名单表达式引擎、开放重定向三重校验、原型污染防护）均经多轮既往审计迭代（W5–W11）。**开放 required = 4（F-001～F-004），全部可在 S2 批内以小改动闭合**；P3 分母见上表。
