---
id: A-001-w8-independent
doc: audit-entry
goal: GOAL-008-w8-api-web-security-audit
source: independent
date: 2026-08-20
scope: apps/api + apps/web current implementation
verdict: fail
---

# A-001 · W8 api/web 独立安全审计

## 审计范围与方法

独立检查 `apps/api` 与 `apps/web` 的认证/会话、授权边界、输入处理、分页、上传下载、配置、前端 token、URL 导航、CSP 与危险脚本面。审计代理未加载任何 skills，未修改文件。主代理对关键行号进行了抽查，并运行 API/Web 回归验证。

## Findings

### F-001 · required · High · 分页整数溢出可触发未处理 panic/DoS

**证据**：

- [resources.go:899](../../../../../apps/api/internal/handler/resources.go:899) 的 `intParam` 只检查 `value < 1`，没有 page 上限或乘法溢出保护。
- [datapermission.go:64](../../../../../apps/api/internal/handler/datapermission.go:64) 计算 `(page - 1) * pageSize` 后直接执行 `policies[start:end]`。
- [account_self.go:305](../../../../../apps/api/internal/handler/account_self.go:305) 计算分页边界后直接切 `all[start:end]`。
- [filelibrary.go:62](../../../../../apps/api/internal/handler/filelibrary.go:62) 对通用资源分页同样执行乘法并切片。

**触发条件**：已认证用户提交极大的正整数 `page`，使 `(page-1)*pageSize` 在 64 位平台溢出为负数；例如 `page=92233720368547760&pageSize=100`。

**影响**：切片边界可能触发 `slice bounds out of range`。即使标准 `net/http` 回收单次请求 panic，攻击者仍可重复制造连接失败、日志与资源压力；若部署层缺少 recovery，可能导致进程级退出。

**建议**：在通用分页入口统一限制 page，使用 checked multiplication/overflow guard；对所有内存切片和 SQL `OFFSET` 计算复用同一安全函数，并补充极大值、溢出与边界回归测试。

**状态**：open / required。

### F-002 · required · Low · CSP 阻止 inline 首屏主题 bootstrap

**证据**：

- [index.html:13](../../../../../apps/web/index.html:13) 包含 inline FOUC 主题初始化脚本。
- [nginx.conf:29](../../../../../apps/web/nginx.conf:29) 设置 `script-src 'self'`，没有 nonce 或 hash。

**影响**：生产浏览器会阻止该 inline 脚本，首屏可能出现错误主题或主题闪烁；这属于可观察的功能/安全策略不一致。

**建议**：将 bootstrap 移到外部静态脚本，或为固定脚本配置 CSP hash/nonce，并增加实际响应头下的浏览器回归检查。

**状态**：open / required。

### F-003 · recommended / known residual · refresh token 存储于 localStorage

**证据**：[tokens.ts:2](../../../../../apps/web/src/account/tokens.ts:2) 与 [tokens.ts:30](../../../../../apps/web/src/account/tokens.ts:30) 将长期 refresh token 明文存入 `localStorage`。

**影响**：任何未来同源 XSS 可读取 token 并跨页面刷新会话。

**处置**：代码已明确记录这是用户接受的 XSS trade-off；本报告未发现独立 XSS 注入点，也未将其列为本波 required。是否升级为独立安全工作，登记为 I-003，待后续范围评估。

### F-004 · conditional / recommended · development fallback 误用风险

**证据**：[main.go:80](../../../../../apps/api/cmd/server/main.go:80) 在 `APP_ENV=development` 缺少 JWT secret 时使用固定开发密钥；[main.go:91](../../../../../apps/api/cmd/server/main.go:91) 使用固定初始密码 `admin`。

**影响**：若把开发配置暴露到外部环境，可能导致令牌伪造或默认账户接管。

**现状**：非 development 环境的 `ValidateProd` 要求强 JWT secret 并拒绝 dev session；Compose 明确设置 `APP_ENV=production`、注入 secrets 并关闭 dev session。因此这是部署误配置条件风险，不是默认生产路径的 confirmed required finding。

## 已观察到的防护与未发现项

- API 已设置请求超时、CORS allow-list、上传 ID/owner 校验及下载响应安全头。
- 未发现生产代码使用 `dangerouslySetInnerHTML`、`eval`、`new Function`，也未发现把 refresh token 放入 URL/query/hash。
- API/Web 依赖与回归验证在 E-001 中记录；验证通过不抵消 F-001/F-002 的开放状态。

## 覆盖缺口

未覆盖外部生产网络边界、真实浏览器代理部署结果、数据库权限、容器运行时文件系统权限，以及完整 fuzz/竞态测试；这些不影响本 A-001 对上述代码路径的 finding 判定。
