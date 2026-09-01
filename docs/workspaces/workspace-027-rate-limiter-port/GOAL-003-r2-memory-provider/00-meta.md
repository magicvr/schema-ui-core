---
id: GOAL-003-r2-memory-provider
title: R2 内存供应商 + 7 处使用点迁移
status: done
parent: GOAL-001-rate-limiter-port
created: 2026-09-01
updated: 2026-09-01
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-027-rate-limiter-port
primary_plan: VP-027-rate-limiter-port
serves_summary: 承载 VP-027 R2 阶段（判据 #2/#3）：internal/ratelimit 内存供应商（演进 loginRateLimiter）+ 7 处使用点全量接入 + 回归不迁移语义（W12 D-002 常量 / D-001 P1）。
---

# GOAL-003 · R2 内存供应商 + 使用点迁移

## 概述

执行 Root 纲领 **R2**：在 D-002 合同（v0.1.1 冻结分母）之上，交付 **内存供应商**（`internal/ratelimit`，演进既有 `loginRateLimiter` 语义——allow 不注册 / 容量驱逐 / Retry-After / 惰性剪枝）并由**组合根单一持有**（fx.Provide）；随后 **7 处使用点全部接入** `kernel.RateLimiterProvider`（登录 / 验证码生成 / 密码修改 / 自助恢复 / MFA verify 独立桶 / MFA step-up / 邀请接受），删除旧 `handler/rate_limit.go`（client-IP 工具迁至 `client_ip.go` 保留在 handler 层）；回归 = 各迁入点既有 handler 测试套件全量通过 + 供应商单元语义（含迁移自 auth_test 的 allow-不注册/驱逐测试）+ trusted-proxy/`loginClientIP` 语义保持。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **迁移策略裁决**：I-027-002 用户裁决（P-004；required） | **已关门**（2026-09-01 用户裁决 **方案 A：演进为内存供应商 + 全量注入**（组合根 fx 单一持有；key 维度按 I-027-004 不扩展；多实例语义随 W12 D-002 单实例边界不变量保持）——D-001） |
| C2 | **供应商 + 迁移落地**：`internal/ratelimit`（Provider/Memory）+ 7 处构造点全部经注入接入；`handler/rate_limit.go` 删除（IP 助手迁移 `client_ip.go`）；既有 limiter 单元测试迁入供应商包；组合根/模块/测试装配更新 | **已关门**（2026-09-01：internal/ratelimit 落地（allow 不注册 / 驱逐 / RetryAfter 委托 kernel 谓词）；7 处注入点全量接入；`newLoginRateLimiter` 全仓 0 残留；`go build` / `go vet` / **`go test ./...` 全量绿（exit 0）**） |
| C3 | **审视与关门**：self A-001 + independent（grok build grok-4.6 · high A-002）合并响应；R2 关门、Root 信息台账回写（I-027-002 verified） | **已关门**（A-001 self `pass` + A-002 grok independent `pass`（0 required · F-001～F-005 全部处置（fixed ×4 · fixed-recording ×1））；R2 关门 3/3 · Root progress 2/4 · I-027-002 verified 回写；2026-09-01） |

`progress` = 已关门检查点数 / 3。当前 **3/3**（R2 已关门）。

## 成功标准（方向级）

1. 内存供应商可用（判据 #2）：滑动窗口 + 容量边界 + 驱逐语义实现并有测试（并发、窗口边界、驱逐、RetryAfter 计算；`-race`）。
2. 使用点迁移不回归（判据 #3 · 完整分母 V-F099）：7 处构造点全部接入端口；各迁入点既有 handler 测试套件全量通过（登录 429 语义 / captcha 10/分钟 / 密码 5 / 恢复 20 / MFA verify 10 独立桶 / MFA step-up 5 / 邀请 10）；W12 D-002 窗口常量（15min/20/`IP|identifier`/Retry-After）保持；GOAL-014 分层锁定显式排除。
3. 供应商语义与合同逐位一致：`Allow` 不注册 key（D-001 P1）、容量驱逐、`RetryAfterSeconds` 经 `kernel.RateLimiterRetryAfterSeconds`、剪枝仅 `Allow`（D-002 v0.1.1 §3）。
4. 层边界：handler 不接触供应商类型（只消费 `kernel.RateLimiterProvider`）；client-IP 工具（trusted-proxy/`loginClientIP`）保留在 handler 层。
5. 未越界：不改 Profile 默认集 / 模块矩阵 / Manifest 装配；不引入 Redis 客户端依赖；不改 Charter。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-027-001 | required | 端口 API 形态 | 判据 #1 | R1 | 用户裁决 | **verified** | — | R1 已裁决（GOAL-002 D-001） |
| I-027-002 | required | 既有 `loginRateLimiter` 迁移策略：演进为内存供应商（推荐） vs 保留并存（双轨） | 判据 #3 | R2（C1） | 用户裁决（P-004） | **verified** | — | 2026-09-01 用户裁决：**方案 A 演进 + 全量注入**（本目标 D-001 accepted；证据 = 迁移实施 + 回归） |
| I-027-003 | non-blocking | 窗口语义默认 | 判据 #2 | R1 | 用户确认 | **verified** | — | R1 已确认（GOAL-002 D-001） |
| I-027-004 | non-blocking | key 维度扩展 | 判据 #1 | R1 | 用户确认 | **verified** | — | R1 已确认（GOAL-002 D-001） |

## 父目标

- `GOAL-001-rate-limiter-port`（Root · 纲领 R2）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式（Root D-001 已定）：阶段关门 default self；本目标动 handler 公共面（登录/恢复/邀请/MFA 安全面）+ 7 个迁入点（security 门禁）→ **C3 走 cross**：A-001 self + A-002 本地 grok build（grok-4.6 · high）independent（项目级默认执行路径 `docs/architecture/independent-audit-execution.md`）。
- 迁移基线：`handler/rate_limit.go`（删除前）为语义逐位对照物；迁移测试 = auth_test 内联 limiter 单元测试迁至 `internal/ratelimit/memory_test.go`（包内直查 attempts/order 保持）。