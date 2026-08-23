# W10 独立审计全文报告 · attachments

> 本附件为 A-001 的完整审计报告，包含详细分析、代码引用与逐条评估。

---

## 一、审计范围与执行方法

**审计对象**：`apps/api`（Go 后端）与 `apps/web`（React/TypeScript 前端）

**执行方式**：
- 主线（DSH 会话模型）深读所有安全关键路径
- 2 个并行子代理广度覆盖：① api（handlers / store / modules / auth / config / composition 全部文件）；② web（account / host / protocol / renderer / app / components + nginx / vite / Dockerfile）
- 全部 P1/P2 结论由主线逐条重读源码交叉复核

**覆盖的关键领域**：
- API：认证（JWT + bcrypt + refresh 轮换）、会话管理、权限门控、文件上传（品牌资源/头像/通用上传）、CSV 导入导出、服务凭据、MFA/TOTP、验证码、速率限制、操作日志、运行模式门控、配置加载、数据库迁移、PostgreSQL 方言
- Web：令牌存储（内存 + localStorage）、认证客户端（并发刷新去重、跨 tab 竞争）、Host 引导协议、Schema 页面加载与验证、渲染器（表单/表格/权限门控/反应引擎）、文件下载、品牌化、nginx 配置、CSP 策略、Vite 构建配置

---

## 二、整体评价

代码库经 W1–W9 九轮审计波次后基础面扎实：

- **SQL 注入**：零发现。全部使用参数化查询（`?` 占位符），无字符串拼接 SQL。
- **认证体系**：JWT 短 TTL（15min 访问 + 720h 刷新）；refresh token SHA-256 哈希存储；bcrypt cost=10；token_version 使密码变更后立即吊销旧令牌；防用户枚举（不存在用户与密码错误均返回相同 401 + bcrypt 恒时比较）。
- **文件上传**：图片强制服务端重编码（PNG/JPEG/GIF/WebP → 新 PNG/JPEG，维度限制），绝不存储原始字节；品牌资源/头像/通用上传分三个独立目录；上传所有权校验。
- **运行时模式**：维护/降级/只读模式阻止业务写操作，保留登录/刷新/MFA/密码修改通路。
- **CSP / 安全头**：nginx 配置了 `default-src 'self'; script-src 'self'; frame-ancestors 'none'; object-src 'none'` 等严格策略。
- **速率限制**：登录按 IP+用户名维度滑动窗口；密码修改按用户维度限流；验证码生成按客户端 IP 限流。
- **信任代理**：`X-Real-IP` 仅从显式配置的 CIDR 白名单信任（默认仅回环），支持 compose 网络内代理。
- **操作日志**：所有关键操作（登录/刷新/注销/CRUD/设置变更/数据导入导出）可审计。
- **配置安全**：生产环境强制 JWT secret ≥32 字符且含字母+数字；`AUTH_DEV_SESSION_ENABLED` 非开发环境强制关闭。

---

## 三、Finding 详细评估

### F-001 · HIGH · env.example 硬编码真实数据库凭据

**文件**：`apps/api/configs/env.example:18, 25, 34-35`

**源码**（2026-08-21 修复时脱敏：不再复现明文凭据，值以下称 `<leaked-password>`）：
```
# Line 18: DB_PASSWORD=<leaked-password>
# Line 25: DB_DSN=postgres://sa:<leaked-password>@192.168.31.213:5432/schema_ui?sslmode=disable
# Line 34-35: PG_TEST_PASSWORD=<leaked-password>, PG_TEST_HOST=192.168.31.213
```

**分析**：`env.example` 是开发者复制为 `.env` 的模板文件，通常会被提交到版本控制。该文件含：
1. 真实数据库密码（`<leaked-password>`）
2. 内网 IP `192.168.31.213`
3. 完整 PostgreSQL 连接字符串（含超级用户 `sa` 的凭据）

另经修复期全仓扫描：`apps/api/internal/config/config_test.go:827-847` 的测试夹具硬编码了同一真实凭据（测试文件同样入库），已同波修复。

**风险**：若仓库公开或被未授权访问，攻击者可直接获取数据库连接信息。虽然 `.env` 本身已 gitignored，但 `env.example` 作为模板文件通常会被提交。此外，该密码泄露可能影响同一密码在其他环境的使用。

**建议**：立即将所有真实凭据替换为占位符（如 `DB_PASSWORD=your_password_here`），并建议相关环境改密。同时确认 IP 地址 `192.168.31.213` 是否仍为活跃地址。

---

### F-002 · MED · Web 认证请求无超时机制

**文件**：`apps/web/src/account/auth-client.ts`（authFetch）、`apps/web/src/renderer/render.tsx`（runRequest）、`apps/web/src/protocol/load-page.ts`（loadPageDocument）

**分析**：所有 `fetch` 调用（包括认证刷新 `/api/auth/refresh`、页面 schema 加载 `/api/schema/{pageId}`、资源 CRUD 请求）均未设置 `AbortController` 超时。虽然 HTTP 协议本身有 TCP 超时，但应用层无显式超时意味着：
- 网络缓慢或服务端挂起时，请求可能无限等待。
- 用户体验下降（页面无响应、无错误提示）。
- 故障恢复延迟（需等待浏览器默认超时，通常 300s+）。

**建议**：添加统一的超时包装（如 30s），超时后 abort 并返回适当错误。可利用 `AbortSignal.timeout()` 或 `setTimeout` + `AbortController`。

---

### F-003 · MED · window.open 文件预览未设置 noopener

**文件**：`apps/web/src/renderer/render.tsx:353`

**源码**：
```typescript
const previewWindow = window.open("about:blank", "_blank");
```

**分析**：`window.open` 返回的窗口引用默认保留 `window.opener` 关系，允许子窗口通过 `window.opener.location` 重定向父窗口（钓鱼攻击）。虽后续写入的 HTML 中 iframe 使用了 `sandbox=""`，但顶层窗口本身缺少此防护。

**建议**：在 previewDocument 的 HTML 中添加 `<meta name="referrer" content="no-referrer">`，或使用 `window.open("about:blank", "_blank", "noopener,noreferrer")`。

---

### F-004 · MED · 刷新令牌并发旋转原子性（PostgreSQL 路径）

**文件**：`apps/api/internal/auth/auth.go:259-264`

**分析**：`RevokeRefreshToken` 使用单条 UPDATE 设置 `revoked_at`。SQLite 用 `SetMaxOpenConns(1)` 使所有请求串行执行，天然安全。但在 PostgreSQL 部署中，两个并发请求可能同时读到同一令牌未吊销，各自执行 UPDATE，导致令牌被"双重使用"——虽然后续的 UPDATE 会成功（两行都设了 revoked_at），但第一个请求已用旧令牌签发了新令牌对。

**当前缓解**：Web 端 auth-client.ts 有 `inflightRefresh` 去重 + 跨 tab 生成号竞争 + 服务器端公钥对比。三重防护使实际利用概率极低。

**建议**：在 PostgreSQL 路径使用 `UPDATE ... WHERE revoked_at IS NULL` 并检查 `RowsAffected == 1`，确保同一令牌最多被一个请求成功吊销。这是低成本的纵深防御。

---

### F-005 · LOW · 文件下载文件名消毒允许点前缀

**文件**：`apps/web/src/renderer/render.tsx:417-425`

**分析**：`sanitizeClientFilename` 正确移除了路径分隔符、控制字符和大多数特殊字符，但保留了 `.`。某些操作系统（macOS）将 `._` 前缀解释为资源 fork，Windows 文件管理器隐藏以 `.` 开头的文件。`...` 等纯点文件名在 Windows 上有效但可能引起混淆。

**建议**：可添加对起始 `.` 的截断逻辑（如 `replace(/^\.+/, '')`），确保文件名不以点开头。

---

### F-006 · LOW · 服务凭据作用域数组无长度上限

**文件**：`apps/api/internal/handler/service_credentials.go:65`

**分析**：`Scopes []string` 无长度限制。需要管理员权限才能创建凭据，且每个作用域仍需对应有效权限键。实际风险低，但作为纵深防御，建议添加上限。

**建议**：添加合理上限（如 100 个作用域），超过时返回 400。

---

### F-007 · LOW · 动态选项源 URL 正则验证边界

**文件**：`apps/web/src/renderer/form-controls.tsx:78`

**分析**：正则 `/^\/(?!\/)[^\s\#]*$/` 正确拒绝了双斜杠、空格、`#`。但 URL 来自服务端 schema 文档，若 schema 被篡改可能注入恶意同源路径。由于 schema 由服务端签名且 fetch 仅同源，实际风险低。

**建议**：可增加额外的 URL 安全校验（如禁止 `..` 路径穿越字符）。

---

### F-008 · INFO · 刷新令牌 localStorage 存储

**分析**：刷新令牌存储在 `localStorage`，易受 XSS 攻击。代码注释明确承认这是"用户接受的 XSS 权衡"（D-002），并通过短 TTL 访问令牌（仅内存）、刷新令牌轮换（每次刷新旧令牌立即吊销）、服务端吊销等机制缓解。CSP 严格（`script-src 'self'`）进一步降低了 XSS 风险。**接受此设计决策。**

---

### F-009 · INFO · 无 CSRF 令牌

**分析**：所有 API 端点依赖 `Authorization: Bearer` 头进行认证。浏览器不会自动附加该头，因此标准 CSRF 攻击不适用。若未来引入 cookie 认证，需增加 CSRF 保护。**当前架构下 CSRF 风险低，已接受。**

---

### F-010 · INFO · SQLite .db 备份文件存在于仓库

**文件**：`apps/api/data/schema-ui.db` 及 48 个 `.pre-v*` 备份文件

**分析**：这些 SQLite 数据库可能包含开发环境中的用户数据（密码哈希等）。若仓库公开，这些文件会泄露。建议确认 `.gitignore` 已排除 `data/*.db` 和 `data/*.sqlite`。

---

### F-011 · INFO · 生产构建 sourcemap 未显式禁用

**文件**：`apps/web/vite.config.ts`

**分析**：未显式设置 `build.sourcemap`。默认情况下 Vite 生产构建可能生成 sourcemap，导致源码泄露。建议确认生产构建中 `sourcemap: false` 已设置。

---

### F-012 · INFO · Dockerfile go mod download 重试逻辑

**分析**：`RUN go mod download || go mod download` 是合理的容错设计（BuildKit 缓存挂载），若 `go.sum` 缺失或损坏，重试仍会失败。已接受。

---

## 四、验证排除项（未发现问题的领域）

以下领域经逐条检查未发现安全漏洞：

| 领域 | 检查结果 |
|------|----------|
| SQL 注入 | 全部参数化查询，0 发现 |
| 路径穿越 | 文件上传/下载路径均经严格校验 |
| 命令注入 | 无 `exec.Command` 等外部命令调用 |
| XSS（服务端） | JSON API 无 HTML 渲染 |
| XSS（客户端） | 无 `dangerouslySetInnerHTML`、无 `eval()` |
| 不安全的反序列化 | JSON 仅用于 API 通信 |
| 硬编码密钥（生产路径） | AUTH_JWT_SECRET 必须外部提供，开发环境有明确警告 |
| 认证绕过 | 所有业务端点经 auth.Middleware 保护 |
| IDOR（上传） | 所有权校验完善 |
| 越权（资源 CRUD） | 权限键门控完整 |
| 敏感信息泄露（错误消息） | 错误码目录化，无堆栈泄露 |
| 加密弱点 | bcrypt cost=10，JWT 签名验证完整 |
| 竞态条件（钱包） | 乐观锁版本 CAS 保护 |
| 日志注入 | 结构化日志（slog.JSON），无格式化字符串注入 |

---

## 五、与 W9 审计的对比

W9（2026-08-21）发现 12 条 required（2 HIGH + 10 MED），涉及 PostgreSQL 方言错误映射、nginx 配置缺失、MFA 并发语义、权限门禁 fail-open 等深层架构问题。W10 审计仅发现 1 条 HIGH（配置卫生），说明 W9 修复后代码库安全基线显著提升。本轮 finding 集中在配置管理和前端纵深防御，属于"最后一公里"的加固。

---

*报告生成：2026-08-21 · DSH 会话模型 · 主线 + 2 并行子代理*