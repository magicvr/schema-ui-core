---
doc_type: goal-audit
id: A-002-r2-closeout-independent
parent: GOAL-003-r2-memory-provider
date: 2026-09-01
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: close-out
scope: GOAL-003 R2 全量（D-002 v0.1.1 ↔ internal/ratelimit 逐节一致性 / 7 处构造点迁移 / 回归不迁移 / 层边界 / 越界核账 / I-027-002 信息门禁）
verdict: pass
open_required: 0
status: active
version: 0.1.0
---

# A-002 · R2 内存供应商与使用点迁移独立交叉审计（independent）

> 编排器代贴（本地 grok build · grok-4.6 · 思考强度 high · headless 单轮输出），全文证据见 `attachments/audit-A-002-grok-output.md`；`source: independent` 保留，正文未改写。

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：close-out（C3 交叉审计 · 实施事实 + 合同一致性 + 迁移完整性 + 回归 + 层边界 + 越界）
- **scope**：`workspace-027-rate-limiter-port` / `GOAL-003-r2-memory-provider`（R2 全量）
- **verdict**：**pass**
- **开放 required 计数**：**0**
- **对照 self**：`03-audit/A-001-r2-closeout-self.md`（`verdict: pass` · `open_required: 0`）

## 1. 范围与区间

| 项 | 值 |
|----|----|
| 工作区 | `workspace-027-rate-limiter-port`（`canonical_scope` 已校验 · `primary_plan` = VP-027） |
| 被审目标 | `GOAL-003-r2-memory-provider`（`status: active` · C1/C2 已关门 · C3 本轮） |
| 冻结分母 | GOAL-002 `D-002-rate-limiter-port-contract.md` **v0.1.1** |
| 迁移策略裁决 | GOAL-003 `D-001-migration-strategy-adjudication.md`（方案 A） |
| 被审实现 | `apps/api/internal/ratelimit/memory.go` + `memory_test.go`；7 处 handler 构造点；`handler/client_ip.go`；`composition` / `server/serve.go` 装配 |
| 迁移对照物 | `git show e53de690:apps/api/internal/handler/rate_limit.go` |
| 排除 | R3 Redis 接缝实现；R4 全仓证据矩阵关门；GOAL-014 DB 行锁；Charter / Profile / Manifest 改写 |

## 2. 独立复跑（2026-09-01）

| 命令 | 结果 |
|------|------|
| `go build ./...` | **exit 0** |
| `go vet ./...` | **exit 0** |
| `go test ./... -count=1` | **exit 0**；handler ok 56.470s · ratelimit ok 2.430s · composition ok 37.233s · kernel ok 2.638s · server ok 5.402s；无 FAIL/PANIC |
| `go test ./internal/ratelimit -count=1 -race` | **ok 3.934s**（RACE_EXIT=0） |
| 仓库根 `git status --short` / `git diff --stat` | 见 §8；禁区空 diff |
| `go.mod` + `go.sum` 检索 `redis` | **0 命中** |

## 3. 信息门禁核验（P-005 / P-004）

| ID | 级别 | 结论 |
|----|------|------|
| I-027-002 | required | **verified**：D-001 `accepted`（方案 A 演进 + 全量注入）；实施 = `internal/ratelimit` + 7 处注入 + 旧类型删除；无双轨残留构造函数 |
| I-027-001/003/004 | required/non-blocking（R1） | **verified**（R1 已关闭；不阻断 R2） |

无到期未关闭 required；无 deferred required；共享资料 `none`。

## 4. 供应商 ↔ 合同逐节一致性（D-002 v0.1.1 ↔ memory.go）

§1 形状（编译期断言 + 签名一致）**通过**；§1 Allow 不注册（直查 attempts==0）**通过**；§2 key 不透明 **通过**；§3 剪枝仅 Allow + RetryAfter 不剪枝 + 窗外→1（含测试锁定）**通过**；§4 capacity≤0 → 1<<16 + FIFO 驱逐（65537-key 测试）**通过**；§5 Retry-After 经 kernel 谓词且与旧实现逐位同构 **通过**；§6 互斥 + `-race` **通过**；§7 无生命周期 **通过**；§8 红线 **通过**。

与旧 `loginRateLimiter`（e53de690）逐位对照：Allow 剪枝 / Record FIFO / Clear 同步删 order / RetryAfter 不剪枝 / capacity 回落 **同构**；可执行谓词改为 kernel 谓词 = 合同强制，非语义漂移。

## 5. 迁移完整性核验（判据 #3 · V-F099 分母 = 7）

| # | 使用点 | 注入面 | 窗口/阈值/容量 |
|---|--------|--------|----------------|
| 1 | 登录 | 中央 `RegisterWithMFAProbes` → `authsHandler` | 15min / 20 / 1<<16 |
| 2 | 验证码生成 | `logincaptcha.New` → `CaptchaRoutes`（包级 var 已删除） | 1min / 10 / 1<<16 |
| 3 | 密码修改 | `account.New` → `AccountSelfRoutes` | 15min / 5 / 1<<16 |
| 4 | 自助恢复 | 中央 `RegisterRecovery` | 15min / 20 / 1<<16 |
| 5 | MFA verify 独立桶 | `mfa.New` → `MFARoutes` 独立实例 | 15min / 10 / 1<<16 |
| 6 | MFA step-up | 同函数第二个实例（不共用桶） | 15min / 5 / 1<<16 |
| 7 | 邀请接受 | 中央 `RegisterInviteAccept` | 15min / 10 / 1<<16 |

- 旧类型残留：`newLoginRateLimiter` / `loginRateLimiter` / `captchaGenerateLimiter` 全仓 **0**；`rate_limit.go` 已删除。
- IP 助手迁移：`SetTrustedProxyCIDRs` / `mustCIDR` / `loginClientIP` / `trustedReverseProxy` / `clientIP` 与旧实现控制流一致（X-Real-IP 仅受信 CIDR）；现居 `handler/client_ip.go`。
- 密码修改 429 与验证码 429 不写 Retry-After = 与旧 handler 相同，非本波回归。

## 6. 装配与层边界

| 断言 | 独立核对 |
|------|----------|
| 组合根单一持有 | `fx.Provide(..., newRateLimiters, ...)` → `ratelimit.NewProvider()`；newMux/newMuxWithExtraProviders 透传 |
| 模块透传 | account / logincaptcha / mfa 生产 provider.go 只持有 `kernel.RateLimiterProvider` |
| `server/serve.go` | 注入 `ratelimit.NewProvider()`（合法第二装配面；非 7 点漏装） |
| handler / modules 生产 import `internal/ratelimit` | **0**；生产 import 仅 composition + serve.go |
| handler 接触供应商类型 | 仅 kernel 接口；无 `*ratelimit.Memory` |

测试装配用 `ratelimit.NewProvider()` 为测试合法缝，不破坏生产层边界。

## 7. 回归评估

全量回归 exit 0（登录 429+Retry-After / captcha 10→429 / 密码修改限流 / 恢复 20 预算 / 邀请防刷 / MFA verify+step-up 429 / trusted-proxy 矩阵）；供应商单元（allow 不注册 / 滑动窗口 / Clear / FIFO / RetryAfter 不剪枝 / 默认容量 65537-key / 并发）；`-race` ok；旧 limiter 单测迁移（auth_test 删除 ~93 行，对等断言在 memory_test）。GOAL-014 分层锁定未纳入端口 ✓。

## 8. 越界核账

变更面 ⊆ 允许集（internal/ratelimit/ · handler · modules/{account,mfa,logincaptcha} · composition · server/serve.go · workspace-027 docs）。禁区（go.mod / go.sum / kernel/profile.go / manifest / config.default.yaml / charter / kernel/ratelimit.go）**未触碰**；`git diff --stat` 已跟踪面 36 files +173/−439（限流本体迁移净删），无 Redis。

## 9. 对照成功标准

判据 #2（供应商可用）**满足**；判据 #3（7 处接入不回归）**满足**；供应商与合同逐位一致 **满足**；层边界 **满足**；未越界 **满足**。

## 10. Findings

| ID | 级别 | 严重度 | 标题 | 建议闭合 |
|----|------|--------|------|----------|
| F-001 | recommended | low | 新文件未过 gofmt：缺 EOF newline；`Memory` 结构体字段对齐 | `fixed`：gofmt 三新文件 |
| F-002 | recommended | low | `03-audit.md` 索引未登记已存在的 A-001 | `fixed`：索引补 A-001 行并代贴本 A-002 |
| F-003 | recommended | low | `00-meta.md` 正文进度句与 frontmatter/goal-tree 不一致（正文写 0/3，实际 2/3） | `fixed`：正文与检查点表对齐 |
| F-004 | informational | low | 注释仍出现历史名 `loginRateLimiter`（无 type/func 残留） | 可选改注释；不作为双轨未迁完 |
| F-005 | informational | low | Root `00-meta` 路线图与 `workspace.md` 仍标 R2「待启动」 | C3 关门时回写 R2 已完成 |

**required：无。开放 required = 0。**

## 11. 与 self（A-001）的异同

A-001 与 A-002 同向 **pass**，无 verdict 冲突、无「一要一否」必改项。独立意见补充 F-001～F-005（gofmt / 索引 / meta 正文 / 注释历史名 / R2 行回写），self 的「loginRateLimiter 0 残留」按构造函数/类型口径成立（注释历史名见 F-004）。

## 12. 结论

**verdict = pass**（0 required）。R2 实施在 scope 内成立：I-027-002 方案 A 已落地；供应商与 D-002 v0.1.1 逐节一致；V-F099 七处构造点全部注入；旧 `rate_limit.go` 删除且 IP 语义保持；生产层边界干净；全量回归与 `-race` 独立复跑绿；禁区未触碰。建议 `/govern`：落盘 A-002 + 索引登记；合并响应（F-001～F-003 关门前提修；F-004/F-005 记录/回写）；C3 关门（GOAL-003 done 3/3 · Root R2 完成 2/4 · workspace.md 同步）。

## 13. 声明

`source: independent`。不修改目标 status / progress / 方案正文 / goal-tree；响应与关门由 `/govern` 处理；保证等级 L0。