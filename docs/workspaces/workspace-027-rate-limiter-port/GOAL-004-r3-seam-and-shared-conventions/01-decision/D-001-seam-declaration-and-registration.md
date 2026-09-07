---
doc_type: goal-decision
id: D-001-seam-declaration-and-registration
parent: GOAL-004-r3-seam-and-shared-conventions
date: 2026-09-01
status: accepted
version: 0.1.0
---

# D-001 · RateLimiter Redis 接缝声明与轨道登记（VP-027 owner 决策）

## 上下文

VP-027 R3（判据 #4/#5）。轨道条款（key 前缀 / 命名空间 / 连接管理 / harness / 变更流程）为 VRev-059 V-F100 收窄后的 owner 冻结约定（短文 v1.0.0），VP-027 激活即继承（短文 §1）；workspace-026 登记义务（短文 §3.3「命名空间登记义务跟踪至首个消费者 / VP-027 激活」）于本区履行。无需新用户裁决（P-004 不触发：无冲突、无 required、无 residual）。

## 决策

| 项 | 决定 |
|----|------|
| 接缝声明 | owner 短文 `cache-redis-seam-and-track.md` v1.0.0 → **v1.1.0**，新增 **§2.6 RateLimiter 供应商接缝声明**：① 端口不变（同一 `kernel.RateLimiter` / `RateLimiterProvider`，7 处使用点零感知）；② key 映射沿用 `<ns>:<key>`（限流命名空间段 = `rl`）；③ **原子窗口原语 = INCR + EXPIRE**（Record = INCR + 首次 EXPIRE；Allow = 只读；Clear = DEL），滑动窗口 Redis 表达（ZSET vs 双桶近似）**触发立项时裁决**（不预裁）；④ 连接管理 = 组合根单一持有 + 配置存在时构造 PING fail-closed；⑤ 测试 harness 与缓存轨道相同（端口契约测试双供应商）；⑥ 不引入 Redis 客户端依赖（`go.mod` 无 redis） |
| 轨道登记 | §3.3 登记表**首条登记**：`rl` = RateLimiter 桶（7 处使用点：登录 / 验证码 / 密码修改 / 恢复 / MFA verify / MFA step-up / 邀请接受；key 内携带使用点维度），归属 VP-027；§1 端口分母增列 `kernel/ratelimit.go` |
| 继承确认 | VP-028 不属 Redis 轨道保持（outbox/MQ 演化）；短文 §3.5 变更流程 → 修订史 v1.1.0 行 |
| 范围 | **零 Go 代码变更**；不改 Profile / Manifest / Charter；不消耗 RT-Q05 trigger |
| 审计模式 | cross：A-001 self + A-002 grok build（grok-4.6 · high）independent |

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 新建独立限流接缝文档 | 未选 | V-F100 收窄 = 单一所有者轨道（Cache+限流共享基建约定）；分文档重开轨道分歧，违反「单一所有者」 |
| Redis key 独立前缀（非 `<ns>:` 约定） | 未选 | 短文 §3.1 冻结条款（`<ns>:<key>`，ns 段登记）；双前缀 = 轨道内两套映射 |
| 预裁滑动窗口实现（ZSET / 双桶） | 未选 | RT-Q05 触发后才立项；预裁 = 预制 Redis 实现红线（Charter 0.4.0 成功边界 #6 / VP-027 非目标） |

## 影响

- 短文 v1.1.0 为 VP-027 判据 #4/#5 的责任载体；触发立项时以 §2.6 + §3 为实施分母。
- Root / VP-027 判定据映射随 R3 关门回写。