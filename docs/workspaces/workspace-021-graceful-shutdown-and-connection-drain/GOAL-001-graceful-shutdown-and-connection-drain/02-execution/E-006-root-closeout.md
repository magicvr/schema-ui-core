---
id: E-006
title: Root 关门（双审闭合 · 结项）· VP-021 关门提案就绪
date: 2026-08-27
status: done
---

# E-006 · Root 关门（2026-08-27）

## 事实

1. **关门双审闭合**：A-001 self `pass`（0 required；F-001/F-002 台账口径）→ A-002 grok independent `conditional`（F-001/F-002 required / F-003～F-006 recommended）→ **编排器全部响应**：F-001、F-002 `fixed`（台账一次对齐 + 超前声明收回）；F-004 `fixed`（新连接拒绝断言）；F-005 `fixed`（`shutdown.starting` 携带 signal 字段）；F-006 `fixed`（合同 v0.1.1 勘误 §9）；F-003 recommended 登记残差。**开放 required = 0**。
2. **独立复跑**（grok 会话）：`go vet ./...` + config/jobs/composition/store `-count=1` + 全量 `go test ./...` 均 exit 0；**PG drain 实测 PASS**（PG_TEST_* 可用环境）——双方言证据链闭合（不再视为 skip 残余）。
3. **结项残余（登记，不阻断）**：
   - F-003（recommended）：进程级 A′/B′（`!windows` 构建排除）与 `docker compose stop` 以 **linux CI 实跑**核销；复审触发 = CI 失败或下一架构 VP 激活前；闭合 = CI 证据 → fixed，或用户书面 accepted-residual。
   - `bye` 打印保留（legacy，不改，A-001 F-002 记录）。
4. **Root `status: done` · progress 3/3**（R1 合同冻结 → R2 实现与测试 → R3 证据与关门全部关门）；goal-tree / workspace.md / vision workspaces.md 同步。
5. **范围纪律**：A3 余项、RT-D03/Q04/Q02、K8s、TLS 终止均未进波（仍 trigger-gated / default-non-goal）；未改 Charter、Profile 默认集、模块矩阵、Manifest、迁移账本。

## VP-021 关门提案（下一步，决策层）

- 本工作区结项证据齐备（退出判据 1～5 对照见 A-001）；交 **`/vision`**：VRev 关门就绪审查（self；可 grok independent）+ **用户书面确认**后 VP-021 `closed`，vision workspaces.md / roadmap 同步；A3 余项与三分支并行规则不变。