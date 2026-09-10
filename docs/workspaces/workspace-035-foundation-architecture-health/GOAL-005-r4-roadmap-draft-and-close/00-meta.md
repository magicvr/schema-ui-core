---
id: GOAL-005-r4-roadmap-draft-and-close
title: R4 路线图草案、文档卫生与关门
status: active
created: 2026-09-10
updated: 2026-09-10
parent: GOAL-001-foundation-architecture-health
version: 0.2.0
progress: 1/5
plan_refs:
  - VP-035-foundation-architecture-health
primary_plan: VP-035-foundation-architecture-health
---

# R4 路线图草案、文档卫生与关门

按 VP-035 方向级判据 4/5/6 收口：产出**下一版总路线图草案**（现状锚点、已交付 vs 下一拍、residual 总账、三分支建议）交 `/vision` editorial；执行 R3 分类为「现在修（文档）」的文档卫生项；完成证据矩阵与 Root 关门。不改 Charter，不实现任何 gated 基础设施，不在用户确认前把草案写成权威路线图。

## 成功标准与检查点

- [x] C1：执行边界冻结——草案范围、文档卫生清单（G-001/G-002/G-003/G-005）、关门证据清单写入 [D-001](01-decision/D-001-r4-execution-boundary.md)；草案与 `docs/vision/**` 的先后关系（草案 → `/vision` editorial → 同事务/紧随冻结）已冻结。
- [ ] C2：路线图草案落盘（现状锚点修正版、A 序列现状、已交付 vs 下一拍、18 条 residual 总账、架构/Admin/业务三分支建议下一拍），交 `/vision` 等待用户 editorial 确认。
- [ ] C3：文档卫生执行——`docs/architecture/overview.md`（G-001）、`docs/architecture/cache-redis-seam-and-track.md` §2.6（G-003）、`docs/vision/roadmap.md` 锚点/RT-P04/RT-D02（G-002，须与草案同一事务或紧随 `/vision` 冻结）、mfa-wrap 旧表述（G-005，`docs/vision/**` 部分随 G-002 一并处理）。
- [ ] C4：VP-035 六条方向级退出判据逐条取证（判据 1～3 已由 R2/R3 交付；判据 4 = 草案已交 `/vision`；判据 5 = 边界未越；判据 6 = 开放 required = 0）。
- [ ] C5：self 自审 + independent 交叉审计（provider = 本地 codex `gpt-5.6-sol`·high，I-035-006 已裁决）后无开放 required，Root `GOAL-001` 关门并在 `goal-tree.md` / `workspace.md` 同步。

进度由五项等权计算；当前 1/5。progress 只作展示，不放行阶段、不关闭 finding、不推导 `done`。

## 本阶段必须遵守的既有裁决

| 来源 | 内容 |
|------|------|
| R3 D-001 裁决 A（用户 2026-09-10） | 「现在修」一律另立；**本 VP 只做文档卫生**，端口/Profile/实现类改动一律出本 VP |
| R3 D-001 裁决 B | residual 处置继承原 VP 留痕；新增/扩大残余须用户逐条书面接受 |
| R3 D-001 裁决 C | independent provider = 本地 codex `gpt-5.6-sol` · 思考强度 high |
| Root D-001 | R4 关门与路线图 editorial 冻结前建议 independent（V-F122） |

## 信息门禁（P-005）

| ID | 级别 | 需要回答 | 影响门禁 | 最晚需要阶段 | 收集动作 | 状态 |
|----|------|----------|----------|--------------|----------|------|
| I-035-003 | required | 业界对照是否迫使改 Charter 非目标 | C4 | R3 | 已判定 = 否 | **verified** |
| I-035-006 | required | independent provider | C5 | C5 之前 | 用户裁决 | **verified (user decision)** |
| — | — | 本阶段暂无新增 required 信息项 | — | — | 如草案需要新数据（如 ROI/成本），先登记 I-035-007 并走 P-004 | — |

## 红线（继承）

不实现 Redis/MQ/K8s/ORM/第三库；不消耗 trigger-gated 行；不改 Profile 默认集；不重开已 closed VP；不把草案写成已冻结权威；不改 Charter。

## 台账布局

平铺五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。

## 父目标

- `parent: GOAL-001-foundation-architecture-health`（Root；R3 已于 2026-09-10 completed）
