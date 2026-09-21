---
doc_type: goal-decision
record_id: D-001
id: D-001-r4-execution-boundary
doc: decision-entry
parent: GOAL-005-r4-roadmap-draft-and-close
date: 2026-09-10
status: accepted
version: 0.1.0
created: 2026-09-10
updated: 2026-09-10
---

# D-001 · R4 执行边界（草案范围 / 文档卫生清单 / 关门证据 / 顺序）

## 触发

Root R1～R3 已关门（R3 经 independent A-003 `fail` → A-004 `fail` → A-005 `pass` → A-006 响应，开放 required = 0）。VP-035 判据 1～3 已由 R2/R3 交付；剩判据 4（路线图草案交 `/vision`）、判据 5（边界保持）、判据 6（开放 required = 0）需在 R4 收口。R3 用户裁决 A 冻结「本 VP 只做文档卫生」。

## 决定

| # | 项 | 决定 |
|---|----|------|
| 1 | 草案产物位置 | `GOAL-005-r4-roadmap-draft-and-close/attachments/roadmap-restatement-draft.md`（工作区证据），**不直接改写** `docs/vision/roadmap.md` 作为已冻结权威 |
| 2 | 草案范围 | 现行锚点修正版；A0–A7 序列现状（A1/A2/A4/A5/A6/A7 已 delivered，A3 仍 gated）；已交付 vs 下一拍；18 条 residual 总账（R3 分类结论）；架构 / Admin 功能 / 业务域三分支建议下一拍 |
| 3 | 文档卫生清单 | G-001 `docs/architecture/overview.md`（现时节）；G-003 `docs/architecture/cache-redis-seam-and-track.md` §2.6（补 `AllowRecord`/`Reserve`/`Cancel`）；G-002 `docs/vision/roadmap.md`（现状锚点、RT-P04、RT-D02）；G-005 `docs/vision/{charter,roadmap,workspaces}.md` 与 `plans/VP-016-*.md` 的 mfa-wrap 旧表述 |
| 4 | 顺序 | ① 草案落盘 → ② 交 `/vision` editorial（VRev 记录）→ ③ `docs/vision/**` 的 G-002/G-005 与草案**同一事务或紧随冻结**执行 → ④ `docs/architecture/**` 的 G-001/G-003 可先行 |
| 5 | 关门证据 | `exit-criteria-matrix.md` 逐条对应 VP-035 六条判据；`doc-hygiene-record.md` 记录四项卫生的执行 diff 与核对；R4 自审 + independent（provider = 本地 codex `gpt-5.6-sol` · high，I-035-006 已裁决） |
| 6 | 红线 | 不改 Charter 目的/边界/非目标；不实现 gated 基础设施；不消耗 trigger 行；不改 Profile 默认集；不重开 closed VP；不把草案写成已冻结权威；不把「现在修」实现留在本 VP |
| 7 | residual 处置 | 沿用 R3 分类（18 条）：6 现在修 / 3 仍 gated / 4 接受残余 / 5 明确不做；本阶段只执行其中「文档」部分，代码类一律另立 |

## 为什么

- 判据 4 要求「已交 `/vision` 等待用户 editorial 确认」，因此草案必须有独立落点，权威路线图仅在 editorial 后更新。
- 判据 5 要求边界保持：把文档卫生与草案放在同一冻结事务，可避免 roadmap 出现「草案已写但陈旧锚点未修」的中间态（VP-035 首波表原文要求）。
- R3 裁决 A 明确本 VP 不实现任何代码整改，故 R4 的「现在修」只能是文档。

## 未选方案

- **R4 直接改写 `docs/vision/roadmap.md`**：越过 `/vision` editorial，违反判据 4 与「不得把未确认草案写成已冻结权威」。
- **R4 顺手实现 G-003 的 Redis 接缝或 RES-T03-tz**：违反裁决 A 与本 VP 红线。
- **把文档卫生拆成独立子目标**：四项卫生同属判据 4 的同一事务，无独立范围或依赖价值，拆分会增加门禁成本。
