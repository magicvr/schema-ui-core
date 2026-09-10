---
doc_type: goal-execution
record_id: E-001
id: E-001-r4-draft-and-handoff
doc: execution-entry
status: recorded
parent: GOAL-005-r4-roadmap-draft-and-close
date: 2026-09-10
created: 2026-09-10
updated: 2026-09-10
version: 0.1.0
---

# E-001 · R4 路线图草案落盘与交 `/vision`

## 已发生事实（2026-09-10）

1. **C1 完成**：[D-001](01-decision/D-001-r4-execution-boundary.md) 冻结 R4 边界（草案落工作区、四项文档卫生清单与顺序、关门证据、红线）。
2. **C2 草案落盘**：[roadmap-restatement-draft.md](attachments/roadmap-restatement-draft.md) —— 含
   - §1 现状锚点修正版（`MaxOpenConns=1` → 小连接池默认 4；补齐 A2/A4/A5/A6/A7 已交付事实）；
   - §2 A 序列现状（A1/A2/A4/A5/A6/A7 delivered；**仅 A3 仍 trigger-gated**）；
   - §3 三分支「已交付 vs 下一拍 vs 仍 gated vs 明确不做」；
   - §4 R3 冻结的 18 条 residual 总账；
   - §5 **10 项**路线图正文改动清单（逐项给锚点、现行文本、建议文本与依据），交 `/vision` editorial 逐项裁决。
3. **交 `/vision`**：草案已就绪等待 editorial；`docs/vision/**` 的正文改动（G-002/G-005 与投影一致性）**尚未执行**，须经用户 editorial 确认后与草案同一事务或紧随冻结执行（R4 D-001 §4）。
4. 本轮未改 `docs/architecture/**`（G-001/G-003 属 C3）、未改 `docs/vision/**`、未改 `apps/**`。

## 证据

| 主张 | 路径 |
|------|------|
| 草案正文 | `attachments/roadmap-restatement-draft.md` |
| 现状锚点依据 | `../GOAL-003-r2-as-built-matrix/attachments/as-built-matrix.md` 行 04（`store.go:29,104–113`） |
| A 序列与 residual 依据 | `../GOAL-004-r3-industry-comparison/attachments/{industry-comparison.md,r3-gap-classification.md}` |
| 改动清单锚点 | `docs/vision/roadmap.md:104,138,209,290–304,16,223,353,392`；`docs/vision/charter.md:73`；`docs/vision/workspaces.md:30,58`；`docs/vision/plans/VP-016-key-rotation-and-backup.md:115` |

## 边界

未改 Charter、未新增/消耗 trigger 行、未实现任何「现在修」代码项、未重开 closed VP；草案明确标注「不是已冻结权威」。
