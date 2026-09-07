---
doc_type: goal-decision
id: D-001-inherit-r1-contract-and-migration-scope
parent: GOAL-003-r2-handler-migration
date: 2026-09-03
status: accepted
version: 0.1.0
---

# D-001 · 继承 R1 合同与 14 处迁移范围冻结

## 背景

Root 纲领 R1（`GOAL-002-r1-contract-freeze`）已圆满关门，正式冻结 `AllowRecord` 端口合同（D-002 v0.1.0）并在 `kernel.RateLimiter` 与内存供应商 `Memory` 中落地单锁原子实现。
本目标（`GOAL-003-r2-handler-migration`）作为纲领 R2，承担将 14 处生产调用点完全迁移至 `AllowRecord` 并确保回归全绿的职责。

## 决策内容

1. **继承 D-002 合同与分母**：
   - 生产迁移范围严格限定在 D-002 §5 列出的 14 处调用点，不扩大范围，不侵入领域逻辑。
2. **两类迁移口径落地**：
   - **立即消费模式**（4 处：验证码生成、Telegram IP/Chat/User 3 桶）：
     - 将原 `if !limiter.Allow(...) { 429 }` 后跟 `limiter.Record(...)`，重构为单一原子调用：`if !limiter.AllowRecord(...) { 429 }`。
   - **失败预算模式**（10 处：登录、密码修改、自助恢复 2 处、MFA 4 处、邀请接受、钱包核销）：
     - 在操作入口处将 `limiter.Allow(...)` 改为原子乐观占槽：`if !limiter.AllowRecord(...) { 429 }`。
     - 移除失败分支中多余的二次 `limiter.Record(...)`（占槽已在入口完成）。
     - 成功分支保持既有的 `limiter.Clear(...)`。在并发场景下，此模式比原有 TOCTOU 更保守安全，且单请求成功或失败净状态与原设计完全一致。
3. **红线保持**：
   - 不修改 Profile 默认集、不实现 Redis、不改动其它内核端口、不删除兼容接口 `Allow`/`Record`。
4. **审计模式**：
   - 阶段关门 default **self**（A-001）；R3 证据关门时再按需安排 independent 审计。

## 影响门禁与验证

- C1 迁移完成后，复跑 `go test ./internal/handler/...` 及 `go test ./internal/channel/telegram/...`。
- 静态核账：14 处位置无遗漏，生产代码中不再保留 Allow→Record 配对。
