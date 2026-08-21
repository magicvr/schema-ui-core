---
id: GOAL-020-w15-user-perspective-findings
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.2.0
---

# D-001 · W15 开波：真实用户视角审视结论与改进项台账

## 1. 审视范围与方法

- **覆盖面**：
  - `apps/web`：认证客户端（`auth-client.ts`）、会话上下文（`AuthContext.tsx`）、登录页（`LoginPage.tsx`）、壳层（`App.tsx`）、Schema 驱动渲染器（`schema-table.tsx`、`form-controls.tsx`、`render.tsx`）、通用组件（`data-table.tsx`、`notification-center.tsx`、`locale-switcher.tsx`）、多语言运行时与字典。
  - `apps/api`：通用资源工厂（`resources.go`）、各业务 Handler（`auth.go`、`account_self.go`、`wallet*.go`、`scheduledtasks.go`、`dictionary.go`、`notifications.go`、`recyclebin.go`、`settings.go`、`upload.go`、`filelibrary.go`）、错误字典（`errorcatalog.go`）、中间件与服务组合根（`composition.go`、`server.go`）。
  - 运维与部署配置：`compose.yaml`、`nginx.conf`、`config.yaml`。
- **方法**：
  - 模拟真实终端用户交互走查（网络中断恢复、改密会话续期、长表格操作重试、多语言一致性、错误提示辨识度）；
  - 模拟前端接入开发者走查（接口字段命名统一性、时间戳格式解析、错误 JSON 信封完备性、CORS 与部署拓扑）；
  - 所有问题均基于实际代码进行 `file:line` 验证，杜绝空泛推测。

---

## 2. 改进项台账（W15-F01 ～ W15-F14）

### 一、高优先级（High · 核心体验阻断与数据契约破坏）

| ID | 模块 | 问题描述 | 代码证据 | 真实影响与改进建议 |
|:---|:---|:---|:---|:---|
| **W15-F01** | Web / Auth | **网络偶发抖动时刷新 Token 会永久抹除本地会话** | `apps/web/src/account/auth-client.ts:129-136`, `382-385` | **影响**：`doRefresh()` 在 `fetch` 抛出网络异常时直接 `clearTokens()`。用户遇弱网、离线或服务端重启抖动时打开网页，本地登录态被误当成凭证失效而强行清空。<br>**建议**：仅在响应状态为 401/403 时执行 `clearTokens()`；网络/5xx 异常保留 Token 并允许重试。 |
| **W15-F02** | Web / Table | **表格数据加载失败是死胡同：展示原始堆栈且无「重试」机制** | `apps/web/src/renderer/schema-table.tsx:517-521`, `apps/web/src/components/data-table.tsx:259-264` | **影响**：网络异常或 500 时，表格仅显示红色文本 `resource fetch failed: HTTP 500...`，无重试按钮，用户无法原地恢复。<br>**建议**：错误文案走多语言映射，并在错误提示下方提供「点击重试」操作按钮。 |
| **W15-F03** | API / Model | **时间字段在不同模块间格式割裂（RFC3339 字符串 vs Unix 秒整数）** | `apps/api/internal/handler/users.go:107` vs `apps/api/internal/handler/scheduledtasks.go:329` vs `apps/api/internal/handler/dictionary.go:279` vs `apps/api/internal/handler/filelibrary.go:177` | **影响**：users 返回毫秒 RFC3339 字符串；scheduled-tasks / dictionary 返回秒级整数；filelibrary 字段名叫 `created`。首方 `formatDisplayTime` 对非 ISO 返回 null，表格显示原始数字，不是 `Invalid Date`。第三方/通用解析仍无法当日期用。<br>**建议**：全局后端统一输出 RFC3339 字符串（`2006-01-02T15:04:05.000Z07:00`）。字段改名 `createdAt`/`updatedAt` 为契约破坏，须单独 go-impact（默认不改字段名）。 |
| **W15-F04** | API / Router | **路由层 404 / 405 绕过统一 JSON 错误信封，返回 Go 原生纯文本** | `apps/api/internal/composition/composition.go:207`（`http.NewServeMux()`，无自定义 NotFound/405） | **影响**：未注册路径返回 Go 原生 `text/plain` "404 page not found"。首方 `readEnvelope` 对非 JSON 降级为 `resource fetch failed: HTTP 404`，**不会** `response.json()` 崩溃。缺口是契约/第三方客户端与错误辨识度，不是首方运行时崩溃。<br>**建议**：顶层 Mux 自定义 NotFound / MethodNotAllowed，统一 JSON 错误信封。 |
| **W15-F05** | API / Server | **纯 API 进程缺少全局安全响应头与灵活的 CORS 支持** | `apps/api/internal/server/server.go:21-29` | **影响**：API 完全依赖外层 Nginx 反代；若开发者采用前后端分离独立域名部署（如前端在 Vercel/CDN，后端在独立主机），OPTIONS 跨域预检会直接失败。<br>**建议**：增加基础中间件层，支持可配置的 CORS 白名单与全局安全响应头（`nosniff` 等）。 |

---

### 二、中优先级（Medium · 交互逻辑欠严谨与文案误导）

| ID | 模块 | 问题描述 | 代码证据 | 真实影响与改进建议 |
|:---|:---|:---|:---|:---|
| **W15-F06** | Auth / Profile | **修改自身密码后会话立刻作废，成功反馈未说明须重新登录** | `apps/api/internal/handler/account_self.go:249-257` | **影响**（A-002 F-001 改写）：改密走 `UpdateUser`，同一事务 `token_version+1` 并撤销全部 refresh；中间件立刻拒绝旧 access（不是「再活约 15 分钟」）。账户页 `onSuccess: reload`。真实路径：204 成功 → reload/下一请求 401 → refresh 失败 → 立刻回登录。15m 只是默认 `AuthAccessTTL`，不是改密残存窗口。<br>**建议**：成功后明确提示「密码已修改，所有会话已注销，请重新登录」并引导重新认证，或即时下发新 Token 续期当前会话。 |
| **W15-F07** | Auth / Error | **Token 刷新失败错误码复用 `UNAUTHORIZED`，文案误报为「用户名或密码错误」** | `apps/api/internal/handler/auth.go:182-184`, `apps/api/internal/errorcatalog/errorcatalog.go:27` | **影响**：API 对过期/撤销 refresh 返回 401 + `UNAUTHORIZED`（目录文案「用户名或密码错误」）。首方 `doRefresh` 失败只 `clearTokens()`，**不**把该文案铺到 Toast/登录页；第三方/错误信封消费者仍会看到密码错误文案。<br>**建议**：独立错误码 `REFRESH_TOKEN_EXPIRED`，翻译为「登录已过期，请重新登录」。 |
| **W15-F08** | API / L10n | **字段级校验失败的 Reason 为纯英文，中英混杂** | `apps/api/internal/errorcatalog/errorcatalog.go:221-224`, `apps/api/internal/handler/resources.go:464` | **影响**：中文环境下出现 `"创建字段无效: name must not be empty"`，字段 reason 无法被多语言覆盖。<br>**建议**：下发结构化规则码（如 `required`、`min_length`）由前端本地化，或在后端扩充通用校验字典。 |
| **W15-F09** | Web / UI | **操作失败类 Toast 4 秒自动消失** | `apps/web/src/renderer/render.tsx:1144-1153`, `1192` | **影响**：error/success 一律 4s 定时消失。显式 `×` 关闭钮已在 `1179-1186`。关键错误视线稍离即遗漏原因。<br>**建议**：`error` 类型 Toast 保持常驻（Sticky），保留显式关闭。 |
| **W15-F10** | API / RateLimit | **429 响应缺少 `Retry-After`，存储配额超限误用 429** | `apps/api/internal/handler/auth.go:104-106`, `apps/api/internal/handler/upload.go:325-327` | **影响**：登录限流 429 `RATE_LIMITED` 无 `Retry-After`（文案是「登录失败次数过多…」，不是「请求过于频繁」）。配额超限也走 429 + `UPLOAD_QUOTA_EXCEEDED`（「超出每用户配额」），配额已满是持久状态，不宜 429。<br>**建议**：限流响应加 `Retry-After`；配额超限改用 413/403 并返回当前容量与上限。 |
| **W15-F11** | API / Wallet | **`GET` 请求产生数据库写副作用（钱包被动自动创建）** | `apps/api/internal/handler/wallet.go:82-103`, `wallet_self.go:36-51` | **影响**：`GET /api/wallet/by-owner/{id}` 会自动向数据库插入钱包并写入审计日志，破坏 HTTP GET 幂等与只读语义，爬虫/预检可被动触发写操作。<br>**建议**：GET 请求仅做只读查询，缺失时返回 404 或 `exists: false`；创建操作显式由 POST 触发。 |
| **W15-F12** | API / List | **各业务模块默认 `pageSize` 不一致（10 / 20 / 50 混用）** | `apps/api/internal/handler/resources.go:393`, `apps/api/internal/handler/wallet.go:62`, `apps/api/internal/handler/scheduledtasks.go:305` | **影响**：不同页面默认展现条数不一致，分页器体验跳变。<br>**建议**：定义全局常量 `DefaultPageSize = 20` 并在全仓保持一致。 |
| **W15-F13** | Web / Account | **个人中心「在线会话列表」无法识别「当前设备」** | `apps/api/internal/handler/account_self.go:312-328` | **影响**：会话列表仅有时间戳，无 User-Agent / IP，用户无法分辨哪条是当前正在使用的设备，撤销操作容易误伤自己。<br>**建议**：记录并在响应中返回 UA/IP，且为当前使用的 Token 注入 `current: true` 标记。 |

---

### 三、低优先级（Low · 细节抛光与轻微瑕疵）

| ID | 模块 | 问题描述 | 代码证据 | 真实影响与改进建议 |
|:---|:---|:---|:---|:---|
| **W15-F14** | Multi | **交互与工程细节小缺口** | 多处文件 | 1. 钱包金额缺少精度（`decimals`）元数据，前端需硬编码换算；<br>2. 顶栏用户/通知下拉菜单（`role="menu"`）缺少方向键键盘导航；<br>3. 表格空状态文案不区分「未搜索到数据」与「列表本身为空」；<br>4. 登录 MFA 阶段无取消按钮；<br>5. 遗留空目录 `apps/web/nginx.conf;C/`（拼写失误产物，不是空文件）。 |

---

## 3. 下一步规划（待用户裁决）

按指示，本阶段**仅建立治理台账，不开展代码修复**。待后续通过独立审计并由用户裁决（P-004）后，建议可按以下方式分批推进：
- **批次 A（核心安全与高危体验阻断）**：W15-F01（会话网络容灾）、W15-F02（表格重试）、W15-F04（404/405 JSON 信封）、W15-F07（401 文案纠偏）；
- **批次 B（API 规范与数据一致性）**：W15-F03（时间格式统一）、W15-F11（GET 幂等性）、W15-F10（429 与配额状态码）、W15-F12（分页大小归一）；
- **批次 C（界面交互与细节抛光）**：W15-F06（改密后会话处理）、W15-F08（校验原因本地化）、W15-F09（Toast 留存）、W15-F13（当前会话标记）、W15-F14（细节缺口）。
