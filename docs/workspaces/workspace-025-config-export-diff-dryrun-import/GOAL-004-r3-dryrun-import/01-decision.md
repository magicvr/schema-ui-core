---
id: GOAL-004-r3-dryrun-import
doc: decision
status: active
parent: GOAL-001-config-export-diff-dryrun-import
created: 2026-08-30
updated: 2026-08-30
version: 0.1.0
---

# 决策记录 · GOAL-004

## 信息需求与阶段门禁

I-025-004（required · 导入失败语义）= C2 方案冻结前置，2026-08-30 已呈报用户裁决（建议 A：原子替换 + 应用前备份；未裁决前 C2 冻结）。C1（dry-run）无裁决依赖。

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-30 | C1 dry-run 实现口径：严格包解码 + env fail-closed + 影响报告（lead · 合同 §2.3 派生） | accepted | `01-decision/D-001-dryrun-posture.md` |
| D-002 | 2026-08-30 | import 写入语义冻结（I-025-004 → verified · 用户裁决方案 A：原子替换 + 应用前备份 .bak） | accepted | `01-decision/D-002-import-write-semantics.md` |

> 编号在本目标内单调不复用。