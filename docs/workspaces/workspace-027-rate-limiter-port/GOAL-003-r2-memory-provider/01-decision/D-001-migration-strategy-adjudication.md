---
doc_type: goal-decision
id: D-001-migration-strategy-adjudication
parent: GOAL-003-r2-memory-provider
date: 2026-09-01
status: accepted
version: 0.1.0
---

# D-001 · 信息裁决：I-027-002 迁移策略（2026-09-01 用户裁决）

## 上下文

R2 前置 required 信息项（P-005 / P-004）。编排器基于仓库事实提出带建议的选项（既有 7 处构造点清单与装配链——authsHandler / CaptchaRoutes / AccountSelfRoutes / RegisterRecovery / MFARoutes / RegisterInviteAccept；`kernel.Cache` + `internal/cache` 供应商先例；W12 D-002 单实例边界与 I-027-004 裁决），2026-09-01 经用户裁决**采纳方案 A**。

## 裁决记录

| ID | 级别 | 选项（用户所见） | 裁决 |
|----|------|------------------|------|
| I-027-002 | required | ① **演进为内存供应商 + 全量注入**：`internal/ratelimit` 实现 `kernel.RateLimiter` + `Provider`（允许不注册/驱逐/RetryAfter 逐位保持）；7 处构造点全部改注入 `kernel.RateLimiterProvider`（组合根 fx 单一持有，Redis 级后续同接口替换）；`handler/rate_limit.go` 删除限流本体，trusted-proxy/client-IP 工具保留 handler 层；既有 limiter 单元测试迁移 `internal/ratelimit`；回归 = 既有 handler 测试套件全绿 ② 保留双轨（旧 `loginRateLimiter` 继续 + 新端口旁路） | **采纳①**：双轨 = 两套限流语义并存且无法满足判据 #3「7 处全部接入」。key 维度按 I-027-004 裁决不扩展（既有 `IP|identifier` / `op|IP|user` / 纯 IP 形态保持）；多实例 key 维度（实例前缀等）随 W12 D-002 单实例官方边界不变量保持（RT-Q05 触发才评估） |

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 保留双轨（并存） | ② | 两套限流语义并存、维护面翻倍；判据 #3「7 处使用点全部接入端口」无法闭环；迁移成本只是推迟非消除 |
| 包级全局 setter（SetRateLimiterProvider 式） | 另拟（未提出） | 组合根注入（fx.Provide）与 kernel.Cache / MailSender / ObjectStore 装配先例一致；显式传参让 7 个注入点与模块 Provider 可测试、可替换（Redis 级同接口） |

## 影响

- 实施（C2）与回归（C3）以本裁决与 D-002 合同 v0.1.1 为分母。
- Root / VP-027 信息台账 I-027-002 随 R2 关门回写 `verified`（证据 = 实施 + 双审）。