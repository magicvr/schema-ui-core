---
doc_type: goal-execution
id: E-001-goal-opened
parent: GOAL-004-r3-seam-and-shared-conventions
date: 2026-09-01
status: done
version: 0.1.0
---

# E-001 · 目标开启（R3 裁决）

## 事实时间线

- 2026-09-01：编排器完成 R3 事实调查——`internal/mail/runtime.go` 全读（cachedAdapter 版本戳失效语义：`updatedAt == cfg.UpdatedAt` 命中否则重建；失效率零延迟热切）；组合根 Fx 结构（`fx.Provide` + `newMux` → `newMuxWithExtraProviders` seam；4 个测试调用点）；RT-Q03 触发条件与 VRev-059 V-F100 语义（轨道收窄 VP-026/027；架构短文或 owner 决策落盘）。
- 2026-09-01：向用户提交 I-026-004（三选项带建议）与 F-002 挂载（两选项带建议）；用户裁决 **不迁移，评估留痕** / **fx 容器持有 + newMux 注入点**（P-004 · D-001 落盘）。
- 2026-09-01：scaffold `GOAL-004-r3-seam-and-shared-conventions` 五件套。

## 产物

- `GOAL-004-r3-seam-and-shared-conventions/` 五件套；`01-decision/D-001-r3-adjudication.md`。

## 下一步

- C2：架构短文 `docs/architecture/cache-redis-seam-and-track.md`（§2 接缝声明 / §3 轨道约定）→ mail 评估附件 → 组合根 fx 改造（fx.Provide(newCache) + newMux 注入 + 4 调用点）→ `go vet`/`go test` 全绿 + `go.mod` 无 redis 复核 → E-002。