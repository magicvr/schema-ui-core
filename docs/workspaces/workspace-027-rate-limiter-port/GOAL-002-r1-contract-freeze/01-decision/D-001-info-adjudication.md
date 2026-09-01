---
doc_type: goal-decision
id: D-001-info-adjudication
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
status: accepted
version: 0.1.0
---

# D-001 · 信息裁决：I-027-001 / I-027-003 / I-027-004（2026-09-01 用户裁决）

## 上下文

R1 合同冻结前置三条信息项（P-005 / P-004）。编排器基于仓库事实提出带建议的选项（既有 `internal/handler/rate_limit.go` 语义——allow 不注册 key（D-001 P1）、容量驱逐、Retry-After 计算；W12 D-002「窗口常量与语义保持现状」；`kernel.Cache` 端口先例；VP-026 ExpiryPolicy 形态；VP-021 停机合同），2026-09-01 经用户裁决**全部采纳建议项**（选项 A ×3）。

## 裁决记录

| ID | 级别 | 选项（用户所见） | 裁决 |
|----|------|------------------|------|
| I-027-001 | required | ① **语义拆分保持**：`Allow(key)`（不注册）+ `Record(key)`（失败才计数）+ `RetryAfterSeconds(key)`（秒）+ `Clear(key)`（成功清零）；`now` 注入；供应商无关 kernel 接口 + 注入工厂（window/max/capacity，capacity≤0 默认 `1<<16`） ② 内聚 Allow（检查即记数） ③ 回调式 Allow | **采纳①**：完整保留 D-001 P1「allow 不注册 key → 喷洒不能撑爆 map」与 W12 D-002「Retry-After 语义保持」；对既有 `loginRateLimiter` 最小语义迁移（R2 演进基线）。 |
| I-027-003 | non-blocking | ① **滑动窗口保持** + 策略接口独立（不与 VP-026 Cache 的 ExpiryPolicy 共用） ② 固定/混合窗口 | **采纳①**：限流是计数窗口、缓存是过期策略，语义不同；共用反耦合两轨道；滑动窗口实现随 R2 供应商演进（含窗口内剪枝）。 |
| I-027-004 | non-blocking | ① **本波不新增复合 key**：现有 key 形态（`IP|identifier` / `op|IP|user` / 纯 IP）完整保持 ② 本波预留「路由+用户」复合 key | **采纳①**：C 端复合维度留给业务域 VP 自行定义（不预制、分母不扩）。 |

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 内聚 Allow（检查即记数） | ② | 取消「allow 不注册」安全保证（D-001 P1）；与既有调用点（allow 后再 record）不兼容，需改 7 处行为语义 |
| 回调式 Allow | ③ | 难测、易错（失败回调易漏/重复）；与既有工具函数形态断裂 |
| 固定/混合窗口 | ② | 与 W12 D-002「窗口常量与语义保持现状」冲突；登录限流从滑动改固定是退化 |
| 复合 key 预留 | ② | 无消费者、扩大分母；违反「不预制」红线（Charter 0.4.0 成功边界 #6 / VP-027 非目标） |
| 容量配置化（config 键） | 另拟 | W12 D-002 常量保持现状；容量随构造参数显式传（默认 `1<<16`），不引入新配置键——与 VP-026 `cache.max_entries` 不同（缓存有持久容量需求，限流容量是内存护栏） |

## 影响

- 合同正文 D-002 按本裁决冻结（§1 端口形状 / §2 key / §3 窗口）。
- Root / VP-027 信息台账同步 `verified`（证据 = 本条目）。