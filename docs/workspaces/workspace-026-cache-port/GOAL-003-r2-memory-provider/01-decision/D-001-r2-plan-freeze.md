---
doc_type: goal-decision
id: D-001-r2-plan-freeze
parent: GOAL-003-r2-memory-provider
date: 2026-09-01
status: accepted
version: 0.1.1
---

# D-001 · R2 方案冻结（2026-09-01 · 勘误 v0.1.1 见文末）

## 上下文

R2 按 R1 冻结合同（D-002 v0.1.1）实施内存供应商与双策略；驱逐算法为合同 §6 留给 R2 的决策点（「驱逐语义 R2 定（近似 LRU 或 FIFO 其一，随供应商冻结）」），2026-09-01 经用户裁决。

## 裁决与冻结内容

| 项 | 冻结 |
|----|------|
| **驱逐策略（用户裁决）** | **FIFO**：插入序队首逐出（container/list + map 双向结构，全操作 O(1)）；覆盖写不改变插入序；对已过期 key 的 Set 视同新条目重插。理由：O(1) 全操作 / 读路径零写锁维护 / 确定性可测 / 无消费者前 LRU 命中率优势属猜测；未来可换策略（供应商内部实现，合同不变）。 |
| **maxEntries 义务诠释** | 合同 §6「任一 Set 后条目数 ≤ maxEntries」按**总条目数**（含未清扫过期项）执行——内存有界最严格；过期项只在读/写路径被惰性清除时释放位置。 |
| **惰性清理** | Get 命中前查 `kernel.CacheEntryExpired`，过期即删（map+list）；Set 对过期 key 视同未命中重插并刷新插入序。 |
| **滑动过期刷新** | 条目持有政策实例（无状态）；Get 命中调 `policy.Refresh(now, expiresAt)`，返回刷新即写回 `expiresAt`。 |
| **拷贝边界** | Set 复制入参；Get 返回新拷贝；绝不出借内部数组（合同 §1）。 |
| **并发** | 供应商级 `sync.Mutex` 全局串行化（单机内存端口，正确性优先）；`-race` 并发测试必过。 |
| **校验顺序** | 每条路径先 `kernel.ValidateCacheSet`（Set）/ key 校验（Delete）/ `ValidCacheNamespace`（Namespace）再触达存储（合同 §8）；过期判定用 `kernel.CacheEntryExpired`（合同 §5）。 |
| **Typed[T] 封装** | `internal/cache.NewTyped[T](view, codec...)`；默认 `JSONCodec`（json.Marshal/Unmarshal）+ 可注入 `Codec[T]` 接口；Get 解码错误返回 error（不伪装 miss/命中）。合同 §1 承诺的 R2 交付。 |
| **配置键** | `cache.max_entries`（YAML）/ `CACHE_MAX_ENTRIES`（env）/ 默认 **10000**；`<= 0` 或解析失败 → fail-closed（LoadError；ValidateProd 再核）。落 `internal/config` + `config.default.yaml` + `configs/.env.example`。 |
| **组合根** | `newCache(cfg) (kernel.Cache, error)` 构建单一实例（内存供应商）；无消费者（R3 mail 迁移评估接入）；wiring 测试锁定构造与 maxEntries 传播。 |
| **审计模式** | C3 **cross**：A-001 self → A-002 本地 grok build（grok-4.6 · high）independent（生产路径先例 workspace-016/021）。 |

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 近似 LRU 驱逐 | 未选 | 每次 Get 需写锁移动节点，读路径成本上升；实现/测试复杂度高；无消费热点可验证收益；可后续替换（供应商内部实现） |
| 逐命名空间互斥锁 | 未选 | 全局锁更简单且正确性等价；细化留到有性能证据时 |
| Set 时全量清扫过期项 | 未选 | O(n) 每写；FIFO 总条目有界已保证内存上限，清扫只需在读/写触碰路径完成 |
| 驱逐仅计活动条目 | 未选 | 总条目计数才是「有界」的最严格兑现；活动计数会让过期洪峰撑破内存上限 |
| 配置键容错回落默认 | 未选 | 静默掩盖配置错误（合同 §6 fail-closed 已定） |

## 影响

- C2 实施分母 = 本合同 + D-002 v0.1.1；实施与验收（R4）以两者为准。

## 勘误（v0.1.1 · 2026-09-01 · A-002 grok build independent 响应）

A-002（independent · **conditional** · required F-001）三条语义级响应：

1. **maxEntries 计数域 = 进程总预算（F-001 · 用户裁决 2026-09-01）**：冻结为**跨命名空间共享的进程级预算**——任一 Set 后**进程总条目**（含过期未清扫、跨 ns）≤ `maxEntries`；驱逐 = **全局** FIFO 最旧条目（热 ns 可能挤掉冷 ns 旧条目，即共享预算的代价）；实现为全局顺序链 + 每 ns map + 总计数。原 v0.1.0 表格「maxEntries 义务诠释」按本裁决解释；`Memory` / `Config` / YAML / `.env.example` 注释与之一致（此前实现按每 ns 计数，与注释矛盾——已按裁决重构）。
2. **活条目覆盖写保位（F-004 → fixed）**：活条目覆盖写**无论是否更换策略实例**均保留全局插入位；仅**过期**条目被 Set 时重插（丢弃旧位）。供应商**不做**策略接口比较（避免不可比较策略类型触发 panic）。
3. **holder 义务（F-002 → 跟踪 R3）**：组合根 `_ = cachePort` 不保活实例（Go 事实）；R3（mail 迁移评估）必须把单一实例挂到长生命周期结构并移除 blank assign。