---
id: GOAL-003-r2-memory-provider
title: R2 内存供应商 + 双策略 + 容量配置键（判据 #2/#3）
status: done
parent: GOAL-001-cache-port
created: 2026-09-01
updated: 2026-09-01
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-026-cache-port
primary_plan: VP-026-cache-port
serves_summary: 承载 VP-026 R2 阶段（判据 #2/#3）：内存供应商（有界 + TTL 惰性清理 + FIFO 驱逐 + 并发安全）+ 绝对/滑动双策略 + 可插拔自定义策略样例 + Typed[T] 封装 + cache.max_entries 配置键（fail-closed）+ 组合根接线。
---

# GOAL-003 · R2 内存供应商 + 双策略

## 概述

执行 Root 纲领 **R2**：按 [D-002 合同](../GOAL-002-r1-contract-freeze/01-decision/D-002-cache-port-contract.md)（v0.1.1）实施内存供应商与双过期策略——`internal/cache` 包（`Memory` 供应商实现 `kernel.Cache`；`AbsoluteExpiry` / `SlidingExpiry` 实现 `kernel.ExpiryPolicy`；`Typed[T]` 类型化封装）；配置键 `cache.max_entries`（`CACHE_MAX_ENTRIES`，默认 10000，非法值 fail-closed）落 `internal/config`；组合根构建单一 `kernel.Cache` 实例（无消费者，R3 mail 迁移评估接入）。**判据 #6 惰性清理语义已随 R1 冻结**（无后台协程）。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **方案冻结**：驱逐策略（用户裁决 FIFO）+ maxEntries 义务诠释 + Typed 封装形态 + 配置键与校验 + 审计模式 | **已关门**（2026-09-01 用户裁决 FIFO；D-001） |
| C2 | **实施**：`internal/cache`（Memory + 双策略 + Typed + 自定义策略样例测试 + `-race` 并发/边界/驱逐测试）；`config` 键 + fail-closed；组合根 `newCache` 接线 + wiring 测试；全量回归 | **已关门**（2026-09-01：实现落地；F-001 用户裁决进程总预算重构 + F-003/F-004 补齐；`go vet` 0 / 全模块 50 包 `go test` 全绿） |
| C3 | **审视与关门**：A-001 self + A-002 grok build（grok-4.6 · high）independent 合并响应；R2 关门、Root 台账回写 | **已关门**（A-001 self `pass` + A-002 grok build independent `conditional`（required F-001 → **用户裁决进程总预算** → fixed）；A-003 合并响应 6+2 findings 全处置；开放 required = 0；2026-09-01） |

`progress` = 已关门检查点数 / 3。当前 **3/3**（R2 已关门）。

## 成功标准（方向级）

1. `Memory` 实现 `kernel.Cache`（compile-time 断言）；`Get/Set/Delete` 行为符合合同 §1～§8（拷贝边界 / 未命中 vs 空值 / nil fail-closed / Delete 幂等 / `ValidateCacheSet` 先于存储触达 / `CacheEntryExpired` 谓词）。
2. 绝对过期与滑动过期策略按合同 §5 语义实现（命中刷新）并有边界测试；自定义策略样例证明接口可插拔（判据 #2）。
3. 有界容量：任一 Set 后总条目数 ≤ `maxEntries`；FIFO 驱逐可测试（判据 #3）；惰性清理仅在读/写路径（判据 #6 维持）。
4. 并发安全：`-race` 下多 goroutine 并发 Get/Set/Delete 无数据竞争（判据 #1/#3）。
5. `Typed[T]` 封装（默认 JSON codec + 可注入 codec）交付（合同 §1 承诺）。
6. 配置键 `cache.max_entries` / `CACHE_MAX_ENTRIES` 生效；默认 10000；非法值 fail-closed（LoadError / ValidateProd）。
7. 组合根单一 `kernel.Cache` 实例接线 + wiring 测试。
8. 未越界：不改端口合同；不改 Profile 默认集 / 模块矩阵 / Manifest；不引入 Redis 依赖；不改 Charter。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-026-004 | non-blocking | 既有 mail runtime `cachedAdapter` 是否迁移到端口（评估，不强制；版本戳失效语义可能不匹配通用 TTL） | 判据 #2（评估面） | R3 | lead 评估 + 用户确认 | 待确认 | R3 目标承载 | — |

（R2 无新 required 信息项；全部前期裁决项已 verified。）

## 父目标

- `GOAL-001-cache-port`（Root · 纲领 R2）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式（Root D-001）：阶段关门 default self；R2 为生产路径实施（配置面 + 并发代码）→ **C3 走 cross**：A-001 self + A-002 本地 grok build（grok-4.6 · high）independent（workspace-016/021 生产路径先例）。
- A-002（R1）F-002 供应商义务（`ValidateCacheSet` 先于触达 / `CacheEntryExpired` / `-race`）列为 C2 必做输入。