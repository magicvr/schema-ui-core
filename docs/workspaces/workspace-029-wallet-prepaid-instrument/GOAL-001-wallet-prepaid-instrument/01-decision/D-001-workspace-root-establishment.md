---
doc_type: goal-decision
id: D-001-workspace-root-establishment
parent: GOAL-001-wallet-prepaid-instrument
date: 2026-09-02
status: accepted
version: 0.1.0
---

# D-001 · 工作区 / Root 建立与开区决策

## 触发

用户 2026-09-02 指令「/vision-audit 审视 VP-029，没问题的话激活 VP-029，然后交编排器开设工作区」。激活门禁已满足：VRev-066 independent `pass`（0 required）+ Admin 类轻量 freshness PASS（`29727510` → `b5c39dfb` 五域零变更，不暂挂 `go`）。slug 按 VP-013～028 惯例 = `workspace-029-wallet-prepaid-instrument`（VP 正文已预登记）。

## 决定

| 项 | 决定 |
|----|------|
| 工作区 | `workspace-029-wallet-prepaid-instrument`（canonical `docs/workspaces/workspace-029-wallet-prepaid-instrument/`） |
| Root | `GOAL-001-wallet-prepaid-instrument`（`parent: null`；primary_plan = `VP-029-wallet-prepaid-instrument`） |
| 愿景角色 | `delivery`（不改变 Charter primary workspace） |
| 纲领阶段 | R1 合同冻结 → R2 主体接缝+账本入金 → R3 Admin 批次面+导出 → R4 证据与关门（串行；阶段内可并行子目标） |
| 审计模式 | 阶段关门 default **self**；资金路径 / 并发核销 / R4 证据与关门实证门禁可按需 **independent**（grok build 先例 · 项目级默认执行路径） |
| 红线 | 不重开 VP-011；不把 C 端用户做成 `admin.users`；不引入支付网关或 Telegram；不改 Profile 默认集 / 模块矩阵 / Manifest；不消耗 RT-Q03/Q05 trigger |

## freshness 三字段（VRev-066 · V-F113）

| 字段 | 值 |
|------|-----|
| consumer_vp | `VP-029-wallet-prepaid-instrument`（vision_ref `schema-ui-core-admin-foundation@0.4.0`） |
| last_freshness_review_at | 2026-09-02（`29727510` → `b5c39dfb` · Admin 类轻量 PASS · 五域零变更） |
| next_freshness_review_trigger | 实施改变共同门禁 / Profile 默认集 / 模块矩阵 / Manifest 装配 / 协议 pin，或激活 VP-031（业务域 freshness + H-002） |
| 原 `go` 候选 | `ed99e88`（clean）；本 VP **不是** Tier D，不消费业务模块解锁 scope |
| VP-008 go 现行 | **active**（workspace-009 Root `vp008_go_status: active` · 2026-09-01 W16 恢复） |
| 消费候选基线 | HEAD `b5c39dfb` |
| 结果 | **PASS（Admin 激活）**；不暂挂 `go` |

## 未选方案

- 不同时激活 VP-030/031（硬前置方向：030 要主体接缝、031 要资金原语；建议激活序 029→030→031）。
- 不把本 VP 写成支付业务域或重开 VP-011。
- 不把 Bot 用户写入 `admin.users`。
