---
status: done
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-001-design-implementation-conformance
version: 0.4.0
---

# 审计索引 · GOAL-036（W25）

> 本文件是稳定索引。正式意见在 `03-audit/A-NNN-*.md`。独立意见不改 `status` / `progress`。

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required | 文件 |
|------|------|--------|-------|---------|---------------|------|
| A-001 | 2026-08-23 | independent | 修正情况 + 修改方式（S1–S5 / I-001 / I-002；非关门） | conditional（原文保留） | **2 → 0**（F-001/F-002 已按 A-002 `fixed` 闭合） | [A-001-correction-review-independent.md](03-audit/A-001-correction-review-independent.md) |
| A-002 | 2026-08-23 | self | 响应 A-001：F-001～F-005 fixed；self 新增 F-006 fixed / F-007（后处理见 E-007） | 响应结论：required 全闭 | 0 | [A-002-response-a001-self.md](03-audit/A-002-response-a001-self.md) |
| A-003 | 2026-08-23 | self | 关门自审（全目标 C1–C6；F-008 由 GOAL-037 闭环后回归关门） | pass | 0 | [A-003-closeout-self.md](03-audit/A-003-closeout-self.md) |
| A-004 | 2026-08-23 | independent | 修正结果复核（F-001～F-008 关闭证据独立复跑；关门后复审） | pass（原文保留） | 0（新增 recommended：F-009/F-010） | [A-004-correction-recheck-independent.md](03-audit/A-004-correction-recheck-independent.md) |
| A-005 | 2026-08-23 | self | 响应 A-004：F-009 frontmatter 统一 done、F-010 DSN 测试补钉 `_txlock`（均 fixed） | 响应结论：recommended 全闭 | 0 | [A-005-response-a004-self.md](03-audit/A-005-response-a004-self.md) |

## 信息就绪核对

| 核对项 | 状态 | 备注 |
|--------|------|------|
| I-001（required · C6 e2e） | closed | 叙事已由 A-002 补充机制归因（FK/CASCADE）；e2e admin 9/9 复跑为 F-001 关闭证据 |
| I-002（non-blocking · C6 活栈） | closed | F-004 台账卫生已修 |
| F-007（预存 flake） | **fixed**（E-007） | 60/60 连跑全绿；真因 = 测试替身 run id 依赖时钟量化精度（非产品 newRunID 回落）；产品侧防御加固保留 |
| **F-008（wallet reconcile 竞态，新）** | open（**移交 GOAL-037-w25-f008-wallet-reconcile-race 承接**） | 池化+FK 时代偶发 `reconcile result = inconsistent`；与 `_txlock` 无关（A/B 实证）；曾受 BUSY 掩盖（BUSY 已由 `_txlock=immediate` 修复）；机制待定，E-007；**用户书面：GOAL-037 关门后再回归关门 GOAL-036** |
| 资料引用 | 无 | 工作区 `shared_materials_catalog: none` |

实施与验证事实见 E-001～E-007；**关门完成（A-003 pass · done 6/6 · 2026-08-23）**；A-004（independent）复核 pass → A-005 响应（F-009/F-010 fixed，2026-08-23）；下级 GOAL-037（F-008 承接）done 4/4。