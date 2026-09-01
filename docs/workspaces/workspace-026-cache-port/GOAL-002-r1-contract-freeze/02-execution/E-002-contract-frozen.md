---
doc_type: goal-execution
id: E-002-contract-frozen
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
status: done
version: 0.1.0
---

# E-002 · 合同冻结 + 端口落地（C2 关门）

## 事实时间线

- 2026-09-01：D-002 合同正文 v0.1.0 冻结（§0～§10：端口形状 / 命名空间 / key / 值语义 / TTL 惰性清理 / 容量义务 / 并发安全 / 错误面 / 验收 / 未选方案）。
- 2026-09-01：`apps/api/kernel/cache.go` 落地（Cache / CacheView / ExpiryPolicy + ValidCacheNamespace / ValidCacheKey + 4 sentinels，文档引用 D-002）。
- 2026-09-01：`apps/api/kernel/cache_test.go` 快测落地（命名空间 16 + key 11 表驱动子例 + 1 sentinel 测试（4 条 `%w` 包装链））。
- 2026-09-01：**首次跑测发现合同实现偏差并当场修正**：原命名空间正则 `^[a-z0-9][a-z0-9-]{0,63}$` 允许尾/连续中划线，与戒严测试期望矛盾 → 收紧为段式规则 `^[a-z0-9]+(-[a-z0-9]+)*$`（首尾禁中划线、禁连续双中划线），D-002 §2 同步修订；sentinel `errors.Is` 测试改为真实 `%w` 包装（原 `errors.New("wrap: "+...)` 不构成包装链，属测试自身缺陷）。
- 2026-09-01：验证——`go vet ./kernel/...` 0；`go test ./kernel/... -count=1` 全绿（含全 kernel 包）；`go build ./...` 全模块编译通过。
- 2026-09-01（A-002 响应补充）：`kernel.ValidateCacheSet` + `kernel.CacheEntryExpired` 入 kernel（F-002 可执行化，D-002 §5/§8 + §11 勘误）；快测增 `ValidateCacheSet` 8 + `CacheEntryExpired` 5 子例 + 编译期端口面断言（stub）；Get godoc 补空值命中（F-005）；`gofmt` 后复跑绿。

## 产物（证据）

- 合同：`01-decision/D-002-cache-port-contract.md`（frozen · v0.1.0）
- 端口：`apps/api/kernel/cache.go`（+ 修订记录见上文）
- 快测：`apps/api/kernel/cache_test.go`
- 回归证据：`go vet` / `go test ./kernel/...` / `go build ./...` 均绿（2026-09-01 现场执行）

## 下一步

- C3：A-001 self 审计 → A-002 本地 grok build（grok-4.6 · high）independent → 合并响应 → GOAL-002 `done`（R1 关门）→ Root 进度 1/4 + 台账同步。