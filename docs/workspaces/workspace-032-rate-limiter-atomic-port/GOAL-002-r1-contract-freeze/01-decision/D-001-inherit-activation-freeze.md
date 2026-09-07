---
doc_type: goal-decision
id: D-001-inherit-activation-freeze
parent: GOAL-002-r1-contract-freeze
date: 2026-09-03
status: accepted
version: 0.1.0
---

# D-001 · 继承激活冻结（R1 无新 P-004 裁决）

## 上下文

用户 2026-09-03 指令「/govern 继续推进下一拍，冻结 r1 方案。有什么需要我决策的吗？」扫描结论：VP-032 激活时 [VRev-073](../../../../vision/reviews/VRev-073-vp032-rate-limiter-atomic-port-activation.md) 已冻结 I-032-001/002；Root D-001 已定审计模式与红线。R1 是合同落盘，不是再开信息门禁。

## 决定

| 项 | 决定 |
|----|------|
| I-032-001 | **沿用** VRev-073：`AllowRecord(key string, now time.Time) bool`；不返回剩余额度；`RetryAfterSeconds` 独立 |
| I-032-002 | **沿用** VRev-073：14 处生产 Allow→Record 全迁（**R2**）；Clear 无需原子变体；立即消费 vs 失败预算两口径 |
| 新 P-004 | **无**。签名、分母、乐观占槽、兼容方法去留均已唯一确定 |
| R1 审计 | 沿用 Root D-001：阶段关门 **self**（不因本次是 kernel 加法而升格 cross；VP-027 R1 的 cross 是新公共面从零冻结） |
| R1/R2 切分 | R1 落地接口 + Memory.AllowRecord（否则 Go 接口无法编译）+ 合同级测试；**不**迁 14 处调用点 |

## 理由

- 再问签名/分母 = 重开已 verified 的 required 信息项，违反 P-005。
- 失败预算入口改为 AllowRecord 的「更保守并发」是消除 TOCTOU 的方向本身，不是新产品行为，激活时已冻结。
- Memory 必须实现新方法，否则 `var _ kernel.RateLimiter = (*Memory)(nil)` 编译失败；这是语言约束，不是把 R2 提前。

## 未选方案

- 不把 Memory.AllowRecord 推迟到 R2（仓库无法编译）。
- 不在 R1 迁移生产调用点（Root 纲领 R2；本拍只冻合同）。
- 不新增剩余额度返回值 / 原子 Clear / Redis。
