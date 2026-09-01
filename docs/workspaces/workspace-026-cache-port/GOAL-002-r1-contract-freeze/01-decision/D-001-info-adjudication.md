---
doc_type: goal-decision
id: D-001-info-adjudication
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
status: accepted
version: 0.1.0
---

# D-001 · 信息裁决：I-026-001 / I-026-002 / I-026-003（2026-09-01 用户裁决）

## 上下文

R1 合同冻结前置三条信息项（P-005 / P-004）。编排器基于仓库事实提出带建议的选项（端口先例 `kernel.Store` / `ObjectStore` / `MailSender`、`internal/mail` `cachedAdapter` 版本戳缓存现状、Go 1.26、VP-021 停机合同、VP-026 判据 #1/#6），2026-09-01 经用户裁决**全部采纳建议项**。

## 裁决记录

| ID | 级别 | 选项（用户所见） | 裁决 |
|----|------|------------------|------|
| I-026-001 | required | ① `[]byte` 负载端口 + 类型化封装 ② Go 泛型端口 `Cache[T]` ③ `any` + 断言 | **采纳①**：`kernel.Cache` 非泛型；负载 `[]byte`；未命中/零值以 `(value, ok)` 区分；类型化封装（Typed[T]）作为 R2 便利层交付，不进入端口本身。理由：与仓库全部端口先例一致；Redis 适配天然；nil 语义清晰。 |
| I-026-002 | required | ① 惰性清理 + 配置化容量驱逐 ② 后台清理协程 + SIGTERM 排空 | **采纳①**：无后台协程 → 无新生命周期 → VP-021 停机义务（判据 #6）自动满足；容量来源 = 配置项（R2 落键，默认 10000 条）+ 驱逐（R2 定策略）。 |
| I-026-003 | non-blocking | ① 显式命名空间 scoped 视图 ② ObjectStore 式参数 ③ 模块 ID 前缀软约定 | **采纳①**：`Cache.Namespace(ns) → CacheView`（fail-closed 形状校验）；Redis key 前缀映射 = `<ns>:<key>`（具体前缀与连接约定由 R3 接缝文档落盘，属 VP-026/027 共享轨道）为预留。 |

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 泛型端口 `Cache[T]` | ② | 与仓库端口先例不一致（全部非泛型）；未来 Redis 适配需构造期注入编解码；契约面复杂化 |
| `any` 负载 | ③ | 无类型安全，调用方断言，违背 fail-closed 风格 |
| 后台协程清理 | ② | 引入协程生命周期（Start/Stop hook + SIGTERM 排空声明），风险面大于收益；VP-026 明确「否则选惰性清理避开新生命周期」 |
| 参数式命名空间 | ② | 与 ObjectStore 同构但调用面冗长；scoped 视图语义等价且更优 |
| 模块 ID 前缀软约定 | ③ | 隔离靠自觉、无法 fail-closed 校验，判据 #1「命名空间隔离」证据弱 |

## 影响

- 合同正文 D-002 按本裁决冻结（§1 端口形状 / §2 命名空间 / §5 TTL 与清理）。
- Root / VP-026 信息台账同步 `verified`（证据 = 本条目）。