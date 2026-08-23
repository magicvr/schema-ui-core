---
title: E-002 · S2/S4 前批事实：A3 边界测试完成；H1/H2 核实为零写入
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-033-w22-residual-closeout
version: 0.1.0
---

# E-002 · S2/S4 前批事实（2026-08-23）

## C4 · A3 导航重复/大小写边界单测 —— 完成

- 实现核对文件：`apps/api/internal/kernel/provider.go`（`NormalizeNavigationOrder` / `sortNavigation`，未改动）。
- 新增测试：`apps/api/internal/kernel/navigation_order_test.go` 三条——`TestSortNavigationDuplicateLegacyOrderIsStable`（重复 order 值按 NodeID 字典序稳定排序）、`TestNormalizeNavigationOrderCaseSensitiveExactMatch`（覆盖键大小写精确匹配，未知键整体回退默认序）、`TestNormalizeNavigationOrderUnknownIDReturnDefault`（未知 ID 返回 DefaultNavigationOrder 且不 mutate）。
- 结果：kernel 包 `go test` PASS（0.731s）；代码行为与 W11·GOAL-013 F-004 残余描述的契约**完全一致**，无生产代码改动。
- 残余处置：W11·GOAL-013 F-004 由「无测」转「有测」，待 S4 回写源台账转 fixed。

## C14 · H1（W11 GOAL-017 F-004 MFA UI residual 提案）—— 核实为零写入

证据链已完整存在：
- [W11·GOAL-017 D-005](../../../../workspace-011-admin-functional-modules/GOAL-017-r3-s10-mfa-2fa/01-decision/D-005-a007-response.md) 第 3 条：用户 2026-08-15 书面裁决——不走 residual 接受路径，**新建下级子目标 GOAL-018-mfa-manager-ui 承接交付**；
- GOAL-018 `done 5/5`（2026-08-15）、GOAL-017 随后回归关门 `done 5/5`（00-meta status: done）。
结论：「建议 accepted-residual」的提案从未被采用，实际以 successor-goal 路径兑现且留痕齐全，无需任何补写。

## C15 · H2（W17 outbound-mail N-001 用词纠偏）—— 核实为零写入

纠偏已在既有关门链完成：
- [W17·GOAL-005-r4-readyz-evidence A-001](../../../../workspace-017-outbound-mail/GOAL-005-r4-readyz-evidence/03-audit/A-001-self-r4-readyz.md) 行 N-001 已改判「分母外 note，非 residual——独立审计 A-002 复核认可」，closed（note）；
- [W17·Root 03-audit](../../../../workspace-017-outbound-mail/GOAL-001-outbound-mail/03-audit.md) 明确「N-001 定性按独立意见更正为分母外 note」。
结论：本项在扫描时命中的只是历史描述文本，非未处理错标。

## 信息项状态复核（同日）

- I-001 **verified**：`netsh interface ipv4 show excludedportrange protocol=tcp` 输出中已无 8011–8110 区间，8080 无监听占用 → A1 补跑放行。
- I-002 **verified**（维持）：VP-003 closed 2026-08-06、W3 区 R5/R6 均 done、R4-I004 review date 2026-08-05 早于今日。

## 进度

C1 ✓ C4 ✓ C14 ✓ C15 ✓ → **4/18**。其余检查点对应代理进行中（A1/A2/A4/A5+A6/B1–B4/B5+B6/H3）。
