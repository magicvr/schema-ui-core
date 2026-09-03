---
doc_type: goal-execution
id: E-002-r1-contract-freeze
parent: GOAL-001-rate-limiter-atomic-port
date: 2026-09-03
status: active
version: 0.1.0
---

# E-002 · R1 立项并冻结合同（C1/C2）

## 事实时间线

- 2026-09-03：用户指令冻结 R1 方案；扫描无 P-004 裁决（I-032-001/002 已 verified）。
- 2026-09-03：立项 `GOAL-002-r1-contract-freeze`；D-001 继承激活冻结；D-002 v0.1.0 冻结 AllowRecord 合同。
- 2026-09-03：端口落地 `kernel.RateLimiter.AllowRecord` + Memory 单锁实现 + 合同级测试。生产 14 处调用点未迁。
- 2026-09-03：Git checkpoint `bdfe925f`。

## 产物

- `docs/workspaces/workspace-032-rate-limiter-atomic-port/GOAL-002-r1-contract-freeze/`
- `apps/api/kernel/ratelimit.go` / `ratelimit_test.go`
- `apps/api/internal/ratelimit/memory.go` / `memory_test.go`

## 下一步（计划）

- GOAL-002 C3：R1 阶段关门 self（A-001）。通过后 Root progress 1/3，立项 R2。
