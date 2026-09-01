---
doc_type: goal-execution
id: E-002-seam-and-track-landed
parent: GOAL-004-r3-seam-and-shared-conventions
date: 2026-09-01
status: done
version: 0.1.0
---

# E-002 · 落盘 + 实现（C2 关门）

## 事实时间线

- 2026-09-01：**架构短文** `docs/architecture/cache-redis-seam-and-track.md` v1.0.0 落盘（owner 文档）——
  - §2 Redis 接缝声明（判据 #4）：端口不变（同一 `kernel.Cache`）/ key = `<ns>:<key>` / TTL 映射表（PX + 滑动 EXPIRE 续期募集）/ 连接管理约定（组合根单一持有 + 启动 PING fail-closed）/ 不引入客户端依赖 / 测试 harness。
  - §3 轨道共享约定（判据 #5）：key 前缀 / 命名空间形状（引用 kernel 段式规则）/ **命名空间登记表（owner 义务，当前为空）** / 连接与 harness / 变更流程（owner 决策；修订史）。
  - §4 触发后专项 / §5 边界复核（`go.mod` 无 redis 实测 0 命中 ✓）。
- 2026-09-01：**I-026-004 评估附件** `attachments/mail-cached-adapter-evaluation-2026-09-01.md`（版本戳 vs TTL 语义对照；三候选方案否决论证；用户确认不迁移 → 判据 #2 评估面闭合）。
- 2026-09-01：**组合根 fx 改造（F-002 兑现）**——`fx.Provide(newCache)` 注册单一实例进 Fx 容器（进程级长生命周期持有 + 沿 newMux→newServer→lifecycle 依赖链 eager 构造 + fail-closed）；`newMux` / `newMuxWithExtraProviders` 增加 `cachePort kernel.Cache` 注入参数（= 首个消费者显式接入点）；seam 内移除局部构建与旧 holder 注释，改为注入参数 + 诚实注释；4 个测试调用点（s2_access_drill / metrics_composition / r5_operational_gate / composition_test）同步补参。
- 2026-09-01：验证——`gofmt`（仅 owned 文件）0；`go build ./...` 通过；`go vet ./internal/{composition,cache,config}/...` 0；`go test ./internal/composition/... ./internal/cache/...` 全绿 + config Cache 测试绿；`go.mod`/`go.sum` redis 命中 0（判据 #4 验证面）。

## 产物（证据）

- `docs/architecture/cache-redis-seam-and-track.md`（v1.0.0 · owner 文档）
- `GOAL-004/attachments/mail-cached-adapter-evaluation-2026-09-01.md`（I-026-004 评估）
- `apps/api/internal/composition/composition.go`（fx.Provide(newCache) + 双函数注入参数）+ 4 测试文件
- `docs/README.md`（architecture 索引行）

## 下一步

- C3：A-001 self → A-002 本地 grok build（grok-4.6 · high）independent → A-003 合并响应 → GOAL-004 `done`（R3 关门）→ Root 进度 3/4。