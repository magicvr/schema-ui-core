---
id: GOAL-001-design-system-and-ui-experience
doc: decision
status: active
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.4
---

# 决策记录 · GOAL-001-design-system-and-ui-experience

## 信息需求与阶段门禁

> 权威信息项表见 `00-meta.md`；本索引同步关键门禁项。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 现行 UI 基线清单 | S1 冻结 | S1 前 | 盘点 apps/web | **closed** | — | E-002 + `attachments/I-S1-001-*` |
| I-002 | required | Token 语义分层与命名 | S1 冻结 | S1 决策时 | Root 决策 | **closed** | — | D-002 **accepted** |
| I-003 | non-blocking | 主范例页验收路径 | S4 | S4 前 | 对照 schemarender | **closed** | — | GOAL-004 |
| I-004 | non-blocking | 对比度是否进退出分母 | 可选 | 任意 | 用户裁决 | **open** | 默认否 | F-V019 路径 b |
| I-005 | required | 目标态视觉方向是否冻结 | S1 对照；S2/S3 呈现 | S1 实施前宜齐 | Stitch + Root 决策 | **closed** | — | D-004 **accepted** + E-004 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-09 | 开区与纲领路线图采纳 | accepted | `01-decision/D-001-workspace-open-and-roadmap.md` |
| D-002 | 2026-08-09 | S1 Token 语义分层与命名约定 | **accepted**（§3/§5 经 D-003 修订） | `01-decision/D-002-s1-token-naming-proposal.md` |
| D-003 | 2026-08-09 | 合并响应 A-001/A-002/A-003（Shadow 映射等） | **accepted** | `01-decision/D-003-audit-response-s1-token-mapping.md` |
| D-004 | 2026-08-09 | 视觉方向冻结（Stitch 定稿） | **accepted** | `01-decision/D-004-visual-direction-freeze.md` |
| D-005 | 2026-08-09 | Root 关门 — 用户书面确认 status: done | **superseded**（D-006） | `01-decision/D-005-root-closeout-user-confirmed.md` |
| D-006 | 2026-08-09 | 废止 D-005 · 回退 Root/工作区完成状态 | **accepted** | `01-decision/D-006-reopen-after-visual-fidelity-audit.md` |
