---
status: active
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-016-w15-api-web-audit-remediation
version: 0.1.0
---

# A-001 · W15 api/web 独立审计意见（立项收录）

- **source**：independent（本会话使用 API、Web、交叉核验三个隔离审查上下文，并由主会话回到源码核对；属于 L0 入口分离级独立意见。项目约定的 `grok build` provider-specific 复核留待后续 S6，不在本条冒充已完成。）
- **日期**：2026-08-30
- **scope**：`apps/api`（公共/主服务配置、认证/JWT、bootstrap、MFA、LocalStore）+ `apps/web`（邀请 token、认证令牌、测试 fixture、构建验证）
- **verdict**：conditional（发现 6 项 required、1 项 recommended；未发现主生产路径无条件认证绕过、IDOR、CSRF、SQL 注入、路径穿越或 React HTML 注入）

## Required findings

### F-001 · P1（条件性）公共 serve 默认暴露开发凭据

- **证据**：`cmd/schema-ui` 直接加载 `server.LoadConfig` 并调用 `server.Serve`；默认配置使用 development、`:25080`、空 secret/password；公共校验将空 `AppEnv` 视为 development；bootstrap 回退 `admin`，JWT 回退固定开发密钥。
- **位置**：`apps/api/cmd/schema-ui/main.go:185-202`；`apps/api/server/config.default.yaml:8-25`；`apps/api/server/config.go:289-295`；`apps/api/server/serve.go:270-301`。
- **影响**：在端口可被局域网/外部访问时，默认实例可使用已知 `admin/admin` 登录，且签名密钥公开可猜。
- **建议**：默认监听回环地址；要求显式 `APP_ENV`；生产禁止使用该默认链，或统一调用严格生产校验。

### F-002 · P1（条件性）公共 serve 生产 JWT secret 缺少强度门禁

- **证据**：`apps/api/server/config.go:289-295` 在非 development 仅要求 secret 非空；主服务 `ValidateProd` 才执行最小长度和字母/数字检查；JWT 使用 HS256。
- **影响**：`schema-ui serve` 以短、可猜 secret 启动时，可通过离线猜测伪造访问 token。
- **建议**：公共配置复用 `internal/config.ValidateProd` 的密钥规则，并增加短 secret 启动失败测试。

### F-003 · P1 生产 bootstrap 初始管理员密码无强度校验

- **证据**：主服务 `resolveSeedHash` 仅处理空值后直接 bcrypt；公共 serve 的 `bootstrapAdmin` 同样只回退空值，不调用密码策略；数据库的 `must_change_password=1` 发生在登录之后。
- **位置**：`apps/api/cmd/server/main.go:96-112`；`apps/api/server/serve.go:270-285`；`apps/api/modules/authsession/systemdata/bootstrap.go:52-55`。
- **影响**：部署者配置极短或可猜的初始密码时，攻击者可先登录，再完成强密码替换；强制改密不能消除首次凭据风险。
- **建议**：bootstrap 前调用统一 8–72 字节密码策略（及已配置复杂度），不满足时 fail closed。

### F-004 · P2 TOTP step-up 可在时间窗内重放

- **证据**：`requireActiveSecondFactor` 只用 `ValidateTotp` 比较 `LastUsedStep`，成功后没有 `AdvanceLastUsedStep`；登录路径却使用 CAS 推进水位。
- **位置**：`apps/api/modules/mfa/service.go:321-342`；对照 `apps/api/modules/mfa/service.go:152-165`；调用接口 `apps/api/internal/handler/mfa.go:302-332`。
- **影响**：同一 TOTP 可在有效窗口内重复触发 recovery-code rotate 等敏感 step-up 操作。该问题需要已有会话和一次有效 TOTP，不是初始认证绕过。
- **建议**：在公共 helper 中以匹配 step 做 CAS 更新，并与敏感操作保持原子性；补重复提交/并发测试。

### F-005 · P2 邀请 token 保留在 URL 与历史记录

- **证据**：前端从 `window.location.search` 读取 token，提交成功后仅 `window.location.href = "/"`，未调用 `history.replaceState`；邀请 token 是一次性 bearer，默认有效期可达 7 天。
- **位置**：`apps/web/src/components/invite-accept.tsx:45-50`、`:67-69`、`:84-88`；`apps/api/internal/handler/invites.go:156-160`。
- **影响**：消费前 token 可能留在地址栏、浏览器历史、截图或同源日志。当前 nginx 的跨站 Referer 策略可降低普通跨站泄露，但不能清理历史记录。
- **建议**：读取后立即 `history.replaceState` 清理 query；邀请页使用 `no-referrer`；保持短 TTL 与一次性消费。

### F-006 · P1（质量门禁）Web Vitest fixture 根路径陈旧

- **证据**：多个测试硬编码 `../../../api/internal/modules`，而当前 schema/fixture 实际位于 `apps/api/modules`。
- **位置**：`apps/web/src/protocol/all-module-schemas-dval.test.ts:17`；`apps/web/src/app/representative-pages.integration.test.tsx:27-46`；`apps/web/src/renderer/schema-crud.test.tsx:28-36`。
- **验证**：Vitest `13 failed | 75 passed (88)`，`76 failed | 1081 passed (1157)`；失败主要为上述路径或缺失 fixture 的 `ENOENT`。
- **影响**：大量集成与 schema 回归测试无法运行，不能作为发布质量证据。
- **建议**：统一 canonical fixture 根、补齐/删除过时 fixture，并在 CI 加入路径存在性检查。

## Recommended finding

### F-007 · P3 主机侧 LocalStore 文件权限偏宽

- **证据**：LocalStore 使用 `os.MkdirAll(..., 0o755)`，对象临时文件及 sidecar 最终为 `0o644`。
- **位置**：`apps/api/internal/objectstore/local.go:113-145`。
- **影响**：多用户 Unix 主机上，其他 OS 账号可绕过 HTTP owner 检查读取对象及元数据；Docker 非 root volume 场景会降低实际暴露。
- **建议**：评估威胁模型；若纳入修正，使用 `0700` 目录和 `0600` 文件，并处理既有文件权限。

## 已核验的安全基线

- 主服务 `cmd/server` 调用 `ValidateProd`；空 `APP_ENV`、生产 dev-session、短 JWT 会 fail closed。
- CORS 使用精确 allowlist，不反射任意 Origin，也不启用 credentials。
- 上传大小、对象 ID、下载 owner 检查和 SQL 参数化在本次覆盖范围内未发现直接漏洞。
- Web 生产源码未发现 `dangerouslySetInnerHTML`、`innerHTML`、`eval` 或 `new Function` 等直接注入汇点。

## 验证基线

- API：`go test ./...` 通过；`go vet ./...` 通过。
- Web：`tsc --noEmit` 通过；Vite production build 通过（仅 bundle 大小警告）。
- Web：Vitest 受 F-006 阻断，不能宣称全量回归通过。
- 工作区：审计开始前代码工作树 `git status` clean，`git diff --check` 通过；本目标治理文档随后按用户指令新增，属于本次开区记录，不是代码实施变更。
