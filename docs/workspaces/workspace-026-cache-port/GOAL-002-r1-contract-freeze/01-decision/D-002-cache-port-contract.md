---
doc_type: goal-decision
id: D-002-cache-port-contract
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
status: accepted
version: 0.1.1
---

# D-002 · Cache 端口合同 v0.1.1（2026-09-01 冻结 · 勘误 v0.1.1 见 §11）

> **责任文件（frozen）**。实施（R2）与验收（R4）以本合同为分母。本波只落端口本体 + 合同级快测；内存供应商、双策略实装与配置键归 R2；Redis 接缝与共享约定归 R3。不改 Profile 默认集 / 模块矩阵 / Manifest 装配；不引入 Redis 客户端依赖；不改 Charter。

## 0. 适用与验收基线

- **契约面**：`apps/api/kernel` 公共面（Go 1.26）；模块与 C 端业务域经 `kernel.Cache` 消费，绝不接触供应商类型（VP-003 薄内核）。
- **先例对齐**：`kernel.Store` / `ObjectStore` / `MailSender` —— 非泛型接口 · ctx 首位 · fail-closed 校验 · sentinel errors（`errors.Is`）。
- **范围外（trigger-gated / 属其它阶段）**：Redis 供应商实现（RT-Q03 触发后才评估）；分布式锁（RT-Q04）；限流 / 消息语义（VP-027/028）；LRU / 分层 / 标签失效 / 批量策略（策略接口预留，实现待消费者触发）；Key 值的业务语义。

## 1. 端口形状（I-026-001 · 冻结：非泛型 `[]byte` 负载 + 类型化封装）

```go
// kernel.Cache：供应商无关缓存端口。消费方先取命名空间视图，再按 key 操作。
type Cache interface {
    Namespace(ns CacheNamespace) (CacheView, error)
}

// kernel.CacheView：单命名空间作用域内的 Get/Set/Delete。
// 全部方法并发安全（§7）；key 在进入供应商前校验（§3）。
type CacheView interface {
    Get(ctx context.Context, key string) ([]byte, bool)
    Set(ctx context.Context, key string, value []byte, policy ExpiryPolicy) error
    Delete(ctx context.Context, key string) error
}
```

- **负载**：`[]byte`，不做任何序列化整形——Key 值的业务语义不进端口（VP-026 首波冻结）。
- **未命中与零值**：`Get` 返回 `(nil, false)` 表示未命中；显式写入的空值（非 nil 零长 slice）命中返回 `(空 slice, true)`。**nil 值不可写入**（§4 fail-closed）。零值与未命中由 `ok` 区分，与 `[]byte` 负载组合后语义无歧义。
- **拷贝边界**：`Set` 复制入参（调用方后续可复用/变更自己的 buffer）；`Get` 返回新拷贝（调用方可自由持有/变更）。供应商不得共享内部底层数组。
- **类型化封装**（R2 交付承诺，不进端口本体）：`internal/cache` 提供 `Typed[T]`（默认 JSON codec，可注入自定义 codec），模块可 `Typed[UserProfile]` 消费，底层仍是 `kernel.CacheView`。
- **为何不选泛型端口**：Go 接口方法不可泛型；参数化接口使「同一端口」成为类型族，未来 Redis 适配需构造期注入 encode/decode，且与仓库全部端口先例不一致（D-001 未选方案）。

## 2. 命名空间（I-026-003 · 冻结：显式 scoped 视图）

- `Cache.Namespace(ns)` 一次校验、返回作用域视图；模块持视图后 `Get/Set/Delete(key)` 免重复传参。
- **开放集合 + 形状校验（fail-closed）**：与 ObjectStore 的封闭集合（三个已命名族）不同——缓存消费方是未来业务域模块（逐个出现的 C 端 / Admin 模块），封闭集合会迫使每新增模块改 kernel。形状规则：
  - 一个或多个小写字母/数字段，以**单中划线**连接（不允许首/尾中划线、不允许连续双中划线）：`^[a-z0-9]+(-[a-z0-9]+)*$`，≤ 64 字节。
  - `ValidCacheNamespace` 为唯一入口；不符合 → `ErrInvalidCacheNamespace`，`Namespace` 拒绝（fail-closed，不回落默认命名空间）。
- **Redis key 前缀映射（预留）**：`<ns>:<key>`（实际前缀与连接约定由 R3 接缝文档落盘，属 VP-026/027 共享轨道 owner 义务）。
- 命名空间不承担 key 排序/枚举语义（无 List/Iterate——不在本波冻结）。

## 3. key 规则（fail-closed）

- 非空；≤ 256 字节；不含控制字符（byte < 0x20 或 0x7f 拒绝；UTF-8 多字节字符允许）。
- `ValidCacheKey` 为唯一入口；Set/Delete 对非法 key 返回 `ErrInvalidCacheKey`；**Get 无错误通道**，对非法 key 按未命中处理 `(nil, false)`（合同明示；调用方在 Set/Delete 侧已能发现编程错误）。

## 4. 值语义

| 情形 | 行为 |
|------|------|
| Set nil value | `ErrInvalidCacheValue`（fail-closed；nil 无法与未命中区分） |
| Set 空值（`[]byte{}`） | 允许；Get 命中返回 `(空 slice, true)` |
| Set nil policy | `ErrInvalidCachePolicy`（fail-closed；静默视为永不过期会掩盖编程错误） |
| Delete 不存在的 key | 幂等成功（nil） |
| Get 过期/未写入 | `(nil, false)`，无错误 |

## 5. TTL 与过期语义（I-026-002 · 冻结：惰性清理）

- **策略接口（可插拔 · 判据 #2 预留）**：

```go
// kernel.ExpiryPolicy：条目过期策略。实现必须无状态、并发安全。
type ExpiryPolicy interface {
    // ExpireAt 返回 now 时刻 Set 的过期时刻；零值 time.Time = 永不过期。
    ExpireAt(now time.Time) time.Time
    // Refresh 在 Get 命中（现过期时刻 previous）后调用；返回新过期时刻
    // 与是否刷新。滑动策略返回 (now+window, true)；绝对策略返回 (previous, false)。
    Refresh(now time.Time, previous time.Time) (time.Time, bool)
}
```

- **两种基础策略**（R2 实装，`internal/cache`）：`AbsoluteExpiry(ttl)`（过期 = 写入时刻 + ttl，命中不刷新）；`SlidingExpiry(window)`（写入与每次命中均刷新为 now + window）。策略实现无状态 → 条目可共享策略实例，Get 命中的刷新由供应商写回条目。
- **清理 = 惰性**：过期条目只在读（Get 发现过期即移除）与写（Set/Delete 路径顺带清扫）时消除；**无后台协程、无新生命周期**（不新增 kernel Hooks Start/Stop、无需 SIGTERM 排空声明）→ **判据 #6 自动满足**（VP-021 停机义务不触发）。
- 过期判断：`expiresAt != zero && !now.Before(expiresAt)` → 过期。**可执行谓词 = `kernel.CacheEntryExpired(expiresAt, now)`**（各供应商必须使用该谓词，见 §11 勘误）。
- 时钟：供应商统一使用单调时钟语义（`time.Now()`；R2 实现时以 `time.Time` 比较，测试注入可控时钟）。

## 6. 容量与驱逐（有界 · 判据 #3 义务）

- 内存供应商**必须有界**：构造参数 `maxEntries`（正整数）；容量来源 = 配置键 `cache.max_entries`（R2 落 `internal/config` + YAML/env 映射，默认 **10000**；非法值 fail-closed）。
- **驱逐语义 R2 定**（近似 LRU 或 FIFO 其一，随供应商冻结）；本合同只冻结义务：任一 Set 后条目数 ≤ maxEntries；驱逐须可测试（R2 并发/边界测试）。
- 端口本身不感知容量（供应商问题）；Redis 轨道将来以服务端 TTL 语义对齐（接缝文档 R3）。

## 7. 并发安全（判据 #1）

- 所有 `CacheView` 方法**并发安全**：多 goroutine 并行 Get/Set/Delete 同一视图无数据竞争、无不一致读取；命中与驱逐竞态最终收敛（允许短窗口内的过期残留，见 §5 惰性语义）。
- R2 供应商测试以 `-race` 覆盖并发边界。

## 8. 错误面

| sentinel | errors.Is | 触发 |
|----------|-----------|------|
| `kernel.ErrInvalidCacheNamespace` | ✓ | `Namespace` 非法 ns |
| `kernel.ErrInvalidCacheKey` | ✓ | Set/Delete 非法 key |
| `kernel.ErrInvalidCacheValue` | ✓ | Set nil value |
| `kernel.ErrInvalidCachePolicy` | ✓ | Set nil policy |

- 未命中不是错误；供应商内部故障（未来 Redis）以供应商错误返回，不得伪装为未命中。
- 端口校验先于供应商触达（fail-closed，与 `MailMessage.Validate` 先例一致）。**可执行入口 = `kernel.ValidateCacheSet(key, value, policy)`**（顺序：key → value → policy；供应商必须先调用再触达存储）（§11 勘误）。

## 9. 验收方式（R2/R4 预告）

- **合同级快测（本目标，C2）**：`kernel/cache_test.go` —— `ValidCacheNamespace` / `ValidCacheKey` 正反例表驱动 + sentinel 存在性 + `ValidateCacheSet` 顺序/三条件 + `CacheEntryExpired` 边界 + 编译期端口面断言（stub 实现 `CacheView` / `ExpiryPolicy`；`%w` 包装后 `errors.Is`）。
- **R2 供应商测试**：compile-time 断言 `Memory *kernel.Cache` 实现；Get/Set/Delete 行为、绝对/滑动过期、惰性清理、驱逐、容量边界、`-race` 并发。
- **R4 证据矩阵**：判据 #1～#8 逐条映射 + 越界核账 + 无 Redis 依赖核验（`go.mod` diff）。

## 10. 未选方案（除 D-001 已记录外）

| 项 | 未选 | 理由 |
|----|------|------|
| Get 带 error 通道（三返回值） | 未选 | 未命中非错误；key 校验在 Set/Delete 侧已 fail-closed；三返回值破坏轻量端口形态 |
| 端口级 List/Iterate/Stats | 未选 | 无消费者需求；供应商内部实现需要时再评估（不冻结） |
| 条目内嵌策略接口实例引用 | 未选（合同不禁止） | 供应商实现细节；合同只冻结策略接口形状与每条目过期时刻语义 |

---

**引用链**：证据 → `D-001`（信息裁决）；端口本体 → `apps/api/kernel/cache.go`（同 commit 落地）；实施责任 → R2（GOAL-003）；验收 → R4（GOAL-005）。

## 11. 勘误（v0.1.1 · 2026-09-01 · A-002 grok build independent 响应）

A-002（independent · pass · 0 required）7 条 recommended/informational 中有 4 条涉及合同正文，全部按下列修正；语义不变，均为表述精确化与可执行化：

1. **Redis 前缀一致性（F-001 → fixed）**：`D-001` I-026-003 裁决行与 §2 / `cache.go` 文件头统一为 **`<ns>:<key>`**；具体前缀与连接约定仍由 R3 接缝文档落盘。
2. **可执行校验入口（F-002 → fixed）**：§8 增补 `kernel.ValidateCacheSet(key, value, policy)`（顺序 key → value → policy），供应商必须先调用再触达存储；§5 增补 `kernel.CacheEntryExpired(expiresAt, now)` 谓词（过期即达即过期：`!expiresAt.IsZero() && !now.Before(expiresAt)`），各供应商必须使用。
3. **Get 空值命中注释（F-005a → fixed）**：`CacheView.Get` godoc 补「存储空值命中 = (空 slice, true)」；端口面编译期断言（stub）补入快测（F-005b）。
4. **台账勘误（F-004/F-007）**：快测计数改为实测（27 表驱动 + sentinel 包装链 → 勘误后含新增 helper 测试共 **40 表驱动子例 + 1 sentinel 测试 + 编译期端口面断言**）；VP-026 I-026-002 行最晚阶段对齐为 R1（语义随合同冻结；容量键 R2）。

原 §1～§10 除上述修订外视为按 v0.1.0 理解；本勘误不改变任何已裁决选项。