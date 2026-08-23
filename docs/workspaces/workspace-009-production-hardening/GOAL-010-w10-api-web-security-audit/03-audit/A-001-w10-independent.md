---
id: A-001-w10-independent
goal: GOAL-010-w10-api-web-security-audit
status: final
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-010-w10-api-web-security-audit
version: 0.1.0
---

# A-001 · W10 api/web 独立安全审计（2026-08-21）

> **消费调和（编排器 · 2026-08-21 · 非原文）**：S2 采纳后经 S3 逐条源码核实，7 条中 4 条 recommended 不成立作废（F-003 不成立 / F-004、F-005、F-006 误报），实际实施 3 条（F-001/F-002/F-007，全部 fixed）。调和依据与证据：[D-003](../01-decision/D-003-w10-scope-reconciliation.md) / [E-002](../02-execution/E-002-w10-s3-implementation.md)。下列原文按条目加状态注记，正文不改写。

- **source**：independent
- **auditor**：DSH 会话模型 + 2 并行子代理（主线深读安全关键路径 + 子代理广度审计 api/web 全部文件 + 逐条 P1/P2 源码交叉复核）。**非工作区默认 grok provider，偏差见 `00-meta.md` I-003。**
- **类型 / scope**：ad-hoc（用户直接指令的独立审计）· `apps/api`（Go，~20k 行非测试代码）与 `apps/web`（React/TS）当前实现的 bug 与安全漏洞
- **verdict**：conditional（1 条 HIGH required + 6 条 MEDIUM recommended + 5 条 informational）

## 范围与区间

2026-08-21 工作树快照。方法与完整证据见全文附件：[attachments/audit-A-001-w10-full-report.md](../attachments/audit-A-001-w10-full-report.md)。

## 结论摘要

**P0=0 · P1=1 · P2=6 · P3=5**。整体评价：代码库经多轮审计波次（W1–W9）后基础面扎实——SQL 注入零发现、认证/会话/令牌轮换严谨、密码 bcrypt + token_version 吊销、文件上传服务端重编码、CSP 头严格、操作日志完备、运行时模式门控有效。本轮问题集中在 **配置安全卫生**（env.example 真实凭据）与 **前端纵深防御不足**（无超时、noopener 缺失等），均已有先例或设计缓解。

## Findings（required）

### F-001 · env.example 含硬编码真实数据库凭据
- 严重度：**high** ｜ 建议：**required** ｜ 状态：open → fixed（E-002，2026-08-21）
- **文件**：`apps/api/configs/env.example:18, 25, 34-35`；另经全仓扫描发现 `apps/api/internal/config/config_test.go:827-847` 测试夹具含同一凭据（同波修复）
- **问题**：示例文件与测试文件包含真实数据库密码（本记录脱敏，下称 `<leaked-password>`）、内网 IP `192.168.31.213`、完整 PostgreSQL DSN（含超级用户凭据）。〔2026-08-21 脱敏注记：修复时本记录不再复现明文值；已按 D-002 从版本控制文件中移除〕
- **风险**：若 env.example 被提交到公开仓库或泄露，攻击者可直接获得数据库连接信息。虽 `.env` 已 gitignored，但 `env.example` 本身是模板文件，通常会被提交。
- **建议**：将所有真实凭据替换为占位符（如 `DB_PASSWORD=your_password_here`），确认 IP 地址不外泄，并建议相关环境改密。

## Findings（recommended）

### F-002 · Web 认证请求无超时机制
- 严重度：**med** ｜ 建议：**recommended** ｜ 状态：open → fixed（E-002，2026-08-21：`lib/fetch-timeout.ts` 30s 包装接线 auth-client/load-page/form-controls 默认路径）
- **文件**：`apps/web/src/account/auth-client.ts`（authFetch）、`apps/web/src/renderer/render.tsx`（runRequest）、`apps/web/src/protocol/load-page.ts`（loadPageDocument）
- **问题**：所有 `fetch` 调用（包括认证刷新、页面 schema 加载、资源请求）均未设置 `AbortController` 超时。网络缓慢或服务端挂起时请求可能无限等待，影响用户体验和故障恢复。
- **建议**：添加统一的超时包装（如 30s），超时后 abort 并返回适当错误。

### F-003 · window.open 文件预览未设置 noopener
- 严重度：**med** ｜ 建议：**recommended** ｜ 状态：open → **作废（不成立，D-003）**：预览窗口仅写入静态模板，无不可信内容插值（iframe `sandbox=""` 禁脚本）；`noopener` 特性使 `window.open` 返回 null，功能依赖该引用写入预览内容
- **文件**：`apps/web/src/renderer/render.tsx:353`
- **问题**：`window.open("about:blank", "_blank")` 打开预览窗口后写入 HTML 内容，但窗口引用未设置 `noopener`。虽后续写入的 iframe 使用了 `sandbox=""`，但顶层窗口本身缺少此防护。
- **建议**：在 previewDocument 的 HTML 中添加 `<meta name="referrer" content="no-referrer">`，或将 `window.open` 第三个参数设为 `"noopener,noreferrer"`。

### F-004 · 刷新令牌并发旋转原子性（PostgreSQL 路径）
- 严重度：**med** ｜ 建议：**recommended** ｜ 状态：open → **作废（误报，D-003）**：accounts.go:337-359 `RevokeRefreshToken` 已是防护式 UPDATE（`WHERE id=? AND revoked_at IS NULL` + RowsAffected + ErrAlreadyRevoked）；审计只读了 auth.go 调用层未核实仓库层实现
- **文件**：`apps/api/internal/auth/auth.go:259-264`
- **问题**：`RevokeRefreshToken` 依赖单次 UPDATE 的原子性。SQLite（`SetMaxOpenConns(1)`）天然串行安全，但 PostgreSQL 部署中并发刷新可能触发竞态。auth-client.ts 已有 `inflightRefresh` 去重 + 跨 tab 重试，但服务端缺乏显式 WHERE 条件防护。
- **建议**：在 PostgreSQL 路径使用 `UPDATE ... WHERE revoked_at IS NULL` 并检查 `RowsAffected`，确保同一令牌最多被吊销一次。

### F-005 · 文件下载文件名消毒允许点前缀
- 严重度：**low** ｜ 建议：**recommended** ｜ 状态：open → **作废（误报，D-003）**：render.tsx:418 `/^[._-]+|[._-]+$/g` 已剥离前导点与尾部分隔符
- **文件**：`apps/web/src/renderer/render.tsx:417-425`
- **问题**：`sanitizeClientFilename` 正确移除了路径分隔符和控制字符，但保留了 `.` 字符。某些操作系统/文件管理器可能将 `._` 前缀解释为资源 fork 或隐藏文件。
- **建议**：可添加对起始 `.` 的截断逻辑（如 `replace(/^\.+/, '')`）。

### F-006 · 服务凭据作用域数组无长度上限
- 严重度：**low** ｜ 建议：**recommended** ｜ 状态：open → **作废（误报，D-003）**：service_credentials.go:151-155 去重后强制 1..64 唯一作用域上限
- **文件**：`apps/api/internal/handler/service_credentials.go:65`
- **问题**：创建服务凭据时 `Scopes []string` 无长度上限。恶意管理员可创建超大权限集合，消耗存储和序列化资源。
- **风险**：低——需要管理员权限才能创建凭据，且每个作用域仍需对应有效权限键。
- **建议**：添加合理上限（如 100 个作用域）作为纵深防御。

### F-007 · 动态选项源 URL 正则验证可被边界绕过
- 严重度：**low** ｜ 建议：**recommended** ｜ 状态：open → fixed（E-002，2026-08-21；核实后升级为真实缺口）：字符类缺 `\\` 排除——WHATWG URL 将特殊 scheme 的 `\` 规范化为 `/`，`/\host` 变协议相对 `//host` 逃逸同源；form-controls.tsx:78 补 `[^\s\\\#]` + 反斜杠回归用例；其余 10 处同型正则核查无同类缺口
- **文件**：`apps/web/src/renderer/form-controls.tsx:78`
- **问题**：正则 `/^\/(?!\/)[^\s\#]*$/` 正确拒绝了双斜杠、空格、`#`。但 URL 来自 schema 文档（服务端渲染），若 schema 被篡改可能注入恶意同源路径。
- **风险**：低——schema 由服务端签名，fetch 仅同源。
- **建议**：可增加额外的 URL 安全校验（如禁止 `..` 路径穿越字符）。

## Findings（informational）

| F-ID | 严重度 | 摘要 |
|------|--------|------|
| F-008 | info | 刷新令牌 localStorage 存储 — 为文档化已接受取舍（D-002），有短 TTL 访问令牌 + 服务端吊销 + 轮换缓解 |
| F-009 | info | 无 CSRF 令牌 — 当前基于 `Authorization: Bearer` 头认证，浏览器不会自动附加，标准 CSRF 不适用；若未来引入 cookie 认证需增加 CSRF 保护 |
| F-010 | info | SQLite .db 备份文件存在于 `apps/api/data/` — 48 个 `.pre-v*` 备份文件可能含开发环境数据哈希；建议确认 `.gitignore` 排除 |
| F-011 | info | 生产构建 sourcemap 未显式禁用 — `vite.config.ts` 未设置 `build.sourcemap`；建议确认生产构建中 `sourcemap: false` |
| F-012 | info | 服务端 Dockerfile `go mod download \|\| go mod download` 重试逻辑 — 合理容错，已接受 |

## 必改项汇总（required）

F-001（HIGH · env.example 硬编码真实凭据）共 1 条 → **fixed**（E-002）；开放 required = **0**。recommended 处置（D-003 调和）：F-002/F-007 fixed；F-003 作废（不成立）；F-004/F-005/F-006 作废（误报）。informational：F-008～F-012 维持原判。

## 结论 + 建议下一步

- ~~verdict **conditional**：1 条 HIGH required 开放。~~ **更新（2026-08-21）**：F-001 已 fixed，开放 required = 0；recommended 3 fixed + 4 作废（D-003）。verdict 消费状态转为「可关门候选」，S4 复核后由用户裁决关门与 go 恢复。
- 建议下一步：① S4 independent 复核（工作区惯例 grok provider；本会话不可用，I-003 待用户裁决）；② 用户书面恢复 VP-008 go 宣称（对齐 W9 D-004）；③ 密码轮换残余项跟踪。
- 本条目为独立意见，不改动目标 status/progress（per P-003）。