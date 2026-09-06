---
id: GOAL-040-w28-admin-passwd-convention
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 1.0.0
---

# E-002 · 端到端真实验证（连现有库 + ADMIN_PASSWD 回退）

## 2026-09-06 · 实际验证效果

### 已发生事实

用户在本地 `apps/api/configs/.env` 配置了 `ADMIN_PASSWD`（键存在，值不落日志）。编排器拉起真实栈验证：

1. **启动**：`go run ./cmd/server`（`APP_ENV=development`，其余从 `configs/.env` 读取；进程 env 优先）——API 连现有 **postgres 库**启动成功，`readyz` 200，custom profile 25 模块 Start+Ready。Web（vite dev `:25173`）就绪。
2. **核心登录验证（AI 助手视角）**：从 `configs/.env` 读取 `ADMIN_PASSWD`（不打印值）→ `POST /api/auth/login {admin, <ADMIN_PASSWD>}` → **HTTP 200 · accessToken 非空 · mustChangePassword=False**；`GET /api/accounts/me` → roles=`admin,editor`。现有库 admin 现密码声明有效，不再撞墙。
3. **smoke 回退验证（自动化测试视角）**：**不设 `SMOKE_PASSWORD`**，仅从 `configs/.env` 导出 `ADMIN_PASSWD`，`bash scripts/smoke.sh`（非 disposable）：
   - SM-001=PASS（参数/安全前提：缺省 SMOKE_PASSWORD 时回退 ADMIN_PASSWD，不再报"缺少 SMOKE_PASSWORD"）
   - SM-002=PASS（readiness）· SM-003=PASS（代理登录）· SM-004=PASS（admin 身份，mustChangePassword=False）· SM-005=PASS（路由）
   - SM-006/SM-007/SM-008=SKIP（预期：非 disposable / 未设 profile / 未设 CSP）
4. **清理**：停止 API 与 Web 后台进程，25080/25173 端口已释放。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 登录 200 + token + mustChangePassword=False | `curl POST :25080/api/auth/login`（密码从 configs/.env 读取，值未输出） |
| admin 身份有效 | `GET /api/accounts/me` → roles=`admin,editor` |
| smoke 回退生效（SM-001~005 全 PASS） | `SMOKE_PASSWORD` 未设 + `ADMIN_PASSWD` 导出 → `bash scripts/smoke.sh` |

### 评估

两个成功检查点（约定可用性 + smoke 回退）在真实现有库上验证通过；无开放问题。关门自审 A-001 证据链补齐。
