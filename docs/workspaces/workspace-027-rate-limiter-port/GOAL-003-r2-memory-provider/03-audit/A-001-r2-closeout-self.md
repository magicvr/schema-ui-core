---
doc_type: goal-audit
id: A-001-r2-closeout-self
parent: GOAL-003-r2-memory-provider
date: 2026-09-01
source: self
scope: GOAL-003 R2 全量（D-002 ↔ internal/ratelimit 逐节一致性 / 7 处迁移完整性 / 回归不迁移 / 层边界 / 越界核账 / 信息门禁）
verdict: pass
open_required: 0
status: active
version: 0.1.0
---

# A-001 · R2 内存供应商与迁移关门自审（self）

## 1. 信息门禁（P-005）

| ID | 级别 | 状态 | 证据 |
|----|------|------|------|
| I-027-002 | required | **verified**（本目标 C1） | 2026-09-01 用户裁决方案 A（D-001 accepted）；实施证据 = 本目标 E-002 + 全量回归绿 |
| I-027-001/003/004 | 已 verified（R1） | 不影响本目标 | — |

## 2. 合同 ↔ 供应商一致性（D-002 v0.1.1）

- **§1 形状**：`Provider.Memory` 编译期断言实现 `kernel.RateLimiterProvider` + `kernel.RateLimiter`；接口方法签名一致；`now` 注入保持。
- **§2 key**：内存供应商不解析 key；handler 侧 key 约定（`IP|identifier` / `op|IP|user` / 纯 IP）零改动；`loginClientIP`/trusted-proxy 语义原样保留（client_ip.go）。
- **§3 窗口**：剪枝仅 `Allow`（`kernel.RateLimiterInWindow`）；`RetryAfterSeconds` 不剪枝（D-002 v0.1.1 勘误逐字落实；测试锁定：窗外键返回 1、Allow 剪枝后归 0）。
- **§4 容量**：capacity≤0 → `DefaultRateLimiterCapacity`（测试：65537-key 驱逐最老键，map 恒 ≤ 1<<16）；构造点容量常量全部保持 1<<16。
- **§5 Retry-After**：`kernel.RateLimiterRetryAfterSeconds` 单一权威，与旧实现在 remain≤0→1 / Round 语义逐位一致。
- **§6 并发**：`sync.Mutex` 保持；`-race` 并发用例绿。
- **§7 停机**：无后台协程、无新生命周期；无 Start/Stop。
- **§8 红线**：无 Redis 依赖（go.mod 零变更）；Profile/模块矩阵/Manifest/Charter 零触碰；GOAL-014 未纳入；7 处分母全部接入（见 §3）。

## 3. 迁移完整性（判据 #3 · V-F099 分母）

代码检索确认：`newLoginRateLimiter` 与 `loginRateLimiter` 全仓 **0 残留**；`rate_limit.go` 已删除；7 处构造点全部经 `kernel.RateLimiterProvider` 注入（handler 只依赖 kernel 接口，无供应商类型 import 落入 handler/模块公共面；`composition` 与 `server/serve.go` 为组合根合法装配面）。

## 4. 回归证据

- 全量 `go test ./... -count=1` **exit 0**（含：登录 429 + Retry-After（TestLoginRateLimit）、captcha 10/分钟 429、密码修改限流、恢复 20 次预算、MFA verify/step-up 限流、邀请接受防刷、trusted-proxy X-Real-IP 信任矩阵）。
- 迁移测试：`internal/ratelimit`（allow 不注册 / 窗口剪枝 / 容量驱逐 / RetryAfter / provider 默认容量 / `-race` 并发）。
- W12 D-002 常量表逐项核对（15min/20/1min/10/15min/5/15min/20/15min/10/15min/5/15min/10 全部保持）。

## 5. 越界核账

变更面 = `apps/api/internal/ratelimit/**`、`apps/api/internal/handler/{auth,health,captcha,account_self,recovery,mfa,invites,client_ip}.go`（+测试）、`apps/api/modules/{account,mfa,logincaptcha}/`、`apps/api/internal/composition/`、`apps/api/server/serve.go`、`docs/workspaces/workspace-027-rate-limiter-port/**`。禁区（go.mod/go.sum、kernel/profile.go、manifest、config.default.yaml、charter）零触碰。

## 6. 验证复跑（2026-09-01）

`go build ./...` ✓ · `go vet ./...` ✓ · `go test ./... -count=1` ✓（exit 0）。

## Verdict

**pass**（0 required）。R2 满足关门条件；建议 A-002（grok build · grok-4.6 · high）independent 复核后合并响应关门。

## Findings

- required：无。
- recommended：无。