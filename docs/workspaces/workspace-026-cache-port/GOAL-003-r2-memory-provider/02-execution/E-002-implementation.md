---
doc_type: goal-execution
id: E-002-implementation
parent: GOAL-003-r2-memory-provider
date: 2026-09-01
status: done
version: 0.1.0
---

# E-002 · R2 实施（C2 关门）

## 事实时间线

- 2026-09-01：`apps/api/internal/cache/` 落地——
  - `policy.go`：`AbsoluteExpiry` / `SlidingExpiry`（无状态值；TTL<=0 = 永不过期；滑动命中刷新）；compile-time 断言。
  - `memory.go`：`Memory` 实现 `kernel.Cache`（总条目有界 ≤ maxEntries · FIFO 驱逐 · 惰性清理仅读写路径 · 全局互斥 · `ValidateCacheSet` 先于存储触达 · `CacheEntryExpired` 谓词 · 拷贝边界）；`MaxEntries()` 访问器供 wiring 测试。
  - `typed.go`：`Typed[T]`（默认 `JSONCodec` + 可注入 `Codec[T]`；解码错误不伪装 miss）。
  - `memory_test.go` 14 测试 + `typed_test.go` 4 测试（含 `nextMidnightPolicy` 自定义策略样例 = 判据 #2 可插拔样本；`-race` 并发）。
- 2026-09-01：实现期发现并修正两处关键缺陷——
  1. **timeline 算术错误（测试自身）**：滑动过期用例第二命中后刷新至 t+13，用例只推进到 t+10 → 修正推进量。
  2. **`append([]byte(nil), empty...)` 返回 nil**：空值命中语义被破坏（stored `[]byte{}` 在 Get 拷贝时坍缩为 nil）→ 增加 `copyBytes`（`make+copy` 保留非 nil 空语义，合同 §1/§4 空值命中）。
- 2026-09-01：配置面——
  - `internal/config`：`CacheMaxEntries` 字段 + `DefaultCacheMaxEntries = 10000` 常量；yamlFile `cache.max_entries`；严格 env 解析 `CACHE_MAX_ENTRIES`（非法值 fail-closed，MAIL_SMTP_PORT 先例）；Load 期 `<=0` fail-closed；ValidateProd `<0` 门禁。
  - `config.default.yaml` / `configs/config.yaml`：`cache.max_entries: 10000` 块；`configs/.env.example`：`CACHE_MAX_ENTRIES` 文档（canonical-env 测试强制）。
- 2026-09-01：组合根 `newCache(cfg)`（内存供应商；0 → 默认 10000（零值 Config 兼容，db/objects 先例）；负值 fail-closed）+ 单一实例 holder + `cache_wiring_test.go` 3 用例。**组合测试首跑暴露零值 Config fx 路径失败 → 修正为零值回落默认语义。**
- 2026-09-01：验证——`go build ./...` 通过；`go vet` 0；`internal/cache`（含 `-race`）、`internal/composition`、`internal/config` 全绿；全模块回归见 E-003 记录。

## 产物（证据）

- `apps/api/internal/cache/{policy,memory,typed}.go` + `{memory,typed}_test.go`
- `apps/api/internal/config/config.go`（CacheMaxEntries/DefaultCacheMaxEntries/严格 env/两处校验）
- `apps/api/internal/config/config.default.yaml`、`apps/api/configs/config.yaml`、`apps/api/configs/.env.example`
- `apps/api/internal/composition/composition.go`（newCache + holder）、`cache_wiring_test.go`

## 下一步

- C3：A-001 self → A-002 本地 grok build（grok-4.6 · high）independent → A-003 合并响应 → GOAL-003 `done`（R2 关门）→ Root 进度 2/4。