---
id: GOAL-002-r1-contract-freeze
title: R1 合同冻结（Cache 端口契约 / 命名空间 / TTL 与策略接口）
status: done
parent: GOAL-001-cache-port
created: 2026-09-01
updated: 2026-09-01
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-026-cache-port
primary_plan: VP-026-cache-port
serves_summary: 承载 VP-026 R1 阶段（判据 #1/#6）：冻结 Cache 端口契约——API 形态（I-026-001）、TTL/清理语义（I-026-002）、命名空间形态（I-026-003）、可插拔策略接口形态；端口本体落 kernel/cache.go + 合同级快测。
---

# GOAL-002 · R1 合同冻结

## 概述

执行 Root 纲领 **R1**：在仓库既有内核端口先例（`kernel.Store` / `ObjectStore` / `MailSender`：非泛型接口 · ctx 首位 · fail-closed 校验 · sentinel errors）之上，冻结 VP-026 **Cache 端口契约** — API 形态（三选一由用户裁决）、TTL/清理语义（惰性 vs 后台协程）、命名空间形态、可插拔策略接口。**合同正文 = GOAL-002 D-002 产物**；端口本体（`apps/api/kernel/cache.go`）+ 合同级快测（`kernel/cache_test.go`）在本目标落地并验证；内存供应商与双策略实装（判据 #2/#3）归 R2（GOAL-003）。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **信息裁决**：I-026-001（API 形态）、I-026-002（TTL 清理）、I-026-003（命名空间）用户裁决（P-004；I-026-001/002 required，003 non-blocking lead 建议确认） | **已关门**（2026-09-01 用户裁决全部采纳建议：[]byte 负载 + 类型化封装 / 惰性清理 + 配置化容量驱逐 / 显式命名空间 scoped 视图——D-001） |
| C2 | **合同正文 + 端口落地**：D-002 合同冻结；`kernel/cache.go`（Cache / CacheView / ExpiryPolicy + 命名空间与 key 校验 + sentinels）实现；合同级快测绿 | **已关门**（D-002 v0.1.1 冻结（含 A-002 响应勘误 §11）；端口 + helper 落地；`go vet` 0 / `go test ./kernel/...` 绿（40 表驱动子例 + 1 sentinel 测试 + 编译期断言）/ `go build ./...` 通过；2026-09-01） |
| C3 | **审视与关门**：合同自审（self A-001）+ independent（本地 grok build grok-4.6 · high A-002）合并响应；R1 关门、Root 信息台账回写 | **已关门**（A-001 self `pass` 0 required + A-002 grok build independent `pass` 0 required；A-003 合并响应 9 条 findings 全处置（fixed ×8 · fixed-recording ×1）；开放 required = 0；2026-09-01） |

`progress` = 已关门检查点数 / 3。当前 **3/3**（R1 已关门）。

## 成功标准（方向级）

1. 端口契约冻结：Cache（Get/Set/Delete + TTL + 命名空间 + 并发安全）供应商无关、快测可断言（判据 #1）。
2. API 形态按用户裁决落盘：非泛型 []byte 负载 + 类型化封装承诺（R2 交付 `internal/cache` Typed 封装）。
3. TTL/清理语义按用户裁决落盘：惰性清理（读/写时清扫），无后台协程 → 无新生命周期 → 停机语义（判据 #6）自动满足。
4. 命名空间按用户裁决落盘：显式 scoped 视图 + fail-closed 校验；Redis key 前缀映射约定预留（R3 接缝文档落盘）。
5. 可插拔策略接口（ExpiryPolicy）形状冻结：绝对/滑动两基础策略在 R2 实装（判据 #2 预留）。
6. 未越界：不改 Profile 默认集 / 模块矩阵 / Manifest 装配；不引入 Redis 客户端依赖；不改 Charter。

## 信息就绪与未知项

与 Root / VP-026 同号镜像（I-026-00x）；I-026-004（mail cachedAdapter 迁移评估）最晚阶段 R3，不在本目标关闭。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-026-001 | required | Cache 端口 API 形态：Go 泛型 vs `[]byte` vs 结构化值；零值/未命中语义 | 方案冻结 + 判据 #1 | C1 | 用户裁决 | **verified** | — | 2026-09-01 用户裁决：**`[]byte` 负载 + 非泛型端口 + 类型化封装**（D-001 accepted；合同 §1） |
| I-026-002 | required | TTL 清理语义：惰性（读时清理） vs 后台协程清理；边界与容量来源 | 判据 #3/#6 | C1 | 用户裁决 | **verified** | — | 2026-09-01 用户裁决：**惰性清理 + 配置化容量驱逐**（D-001 accepted；合同 §5/§6） |
| I-026-003 | non-blocking | 命名空间 / key 前缀约定：模块 ID 前缀 vs 独立命名空间参数 | 判据 #1/#4 | C1 | lead 建议 + 用户确认 | **verified** | — | 2026-09-01 用户确认：**显式命名空间 scoped 视图**（D-001 accepted；合同 §2） |

## 父目标

- `GOAL-001-cache-port`（Root · 纲领 R1）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式（Root D-001 已定）：阶段关门 default self；本目标落 kernel 公共面（兼容性门禁）→ **C3 走 cross**：A-001 self + A-002 本地 grok build（grok-4.6 · high）independent（项目级默认执行路径 `docs/architecture/independent-audit-execution.md`）。
- R1 合同为本 VP 首波冻结分母；D-002 冻结后实施（R2）与验收（R4）以本合同为准。