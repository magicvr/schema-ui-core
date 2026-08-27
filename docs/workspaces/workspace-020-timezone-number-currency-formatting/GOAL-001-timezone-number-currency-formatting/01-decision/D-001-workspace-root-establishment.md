---
id: GOAL-001-timezone-number-currency-formatting
title: 决策 · 激活与开区（lead 绑定 / 纲领 / 信息门禁）
status: accepted
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# D-001 · 激活与开区：lead 绑定、纲领 R1～R4、信息门禁登记

## 决策（2026-08-26 · 用户指令「对 VP-020 做激活审查并推进开区」）

1. **激活绑定**：VP-020 `planned → active`；lead = `workspace-020-timezone-number-currency-formatting`（唯一 delivery）；`primary_plan` = VP-020；不改变 Charter `primary_workspace`。
2. **激活门禁证据**：VRev-044（self）`pass`，V-F079/V-F080 → 激活事务内 fixed；Admin 类 freshness PASS（`66f5fd1f` → `c6fda691`，不暂挂 `go`）；Vision open required = 0。见 `docs/vision/reviews/VRev-044-vp020-intent-activation.md`。
3. **纲领阶段**（P-001）：R1 合同冻结 → R2 时区语义 → R3 数字/货币语义 → R4 证据与关门；同一纲领阶段内允许并行子目标。
4. **信息门禁**：I-001（时区来源）与 I-002（数字/货币落点）为 required，**R1 方案冻结前须用户裁决**；I-003/I-004 保持 VP 冻结投影（不进退出分母）；I-005 non-blocking 至 R2。R1 关闭前禁止直接改时区/格式相关 DDL 或迁移台账。

## 未选方案

- 不等 freshness 直接激活：否决——VP-008 `go` 消费有效性要求每激活前 freshness（候选自 `66f5fd1f` 起代码与迁移有变更，须核对）。
- 开区推迟到 independent 审查：否决——本 VP 为常规、可逆、非安全/数据/迁移类（P-003 审计模式 `self` 可唯一判定）。

## 影响

- 建立 Root 五件套与 ledger 目录；goal-tree、workspace.md、vision 索引原子同步（roadmap / README / workspaces）。
- 实现推进与阶段审计由 `/govern` 编排（P-002 原语），本 Root 不预写任何实施事实。