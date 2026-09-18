---
id: E-006-r5-closeout-evidence-digest
doc: execution-entry
status: recorded
goal_id: GOAL-006-r5-composition-acceptance
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-006-r5-composition-acceptance
version: 1.0.0
---

# E-006 · 关门证据一页与投影滞后修正

## 事实

2026-09-18，你要求「先给一页证据再定」（未选择立即关门），据此执行两项**不涉及关门**的工作：

1. 产出关门证据一页：`attachments/r5-closeout-evidence-digest.md`（v1.0.0，`status: provided`）。内容含：R5-I-004 门禁定位、十项台账快照、VP-037 七条方向级退出判据逐条证据、本次复跑结果、关门后仍开放/仍 gated 清单、确认后的 C4 投影清单与确认模板。
2. 裁决前复跑证据刷新（不改代码，仅运行与核对）：

| 检查 | 结果 |
|------|------|
| `npm run typecheck`（`tsc -b` + `tsc -p e2e/tsconfig.json`） | exit 0 |
| `npm test`（Vitest） | 112 文件 / 1426 测试通过 |
| `npm run test:e2e -- list-visual-surface`（`APP_PROFILE=admin`） | 2 passed（32.9s） |
| 同上（`APP_PROFILE=mvp`） | 2 passed（33.3s） |
| `git diff --check` | 通过 |
| `git diff --name-only 89666e5c..HEAD -- apps/api` | 空（API 零漂移） |
| 五件套结构扫描 | Root 与 R1～R6、GOAL-008、GOAL-009 全部齐全 |

对照 R5 C3（E-004）的 `110/1408`：增量来自 GOAL-008/GOAL-009 新增测试文件，无既有用例删除或跳过。

3. 投影滞后修正（**目标内当前态**，非历史时间线）：

- `00-meta.md` 父目标行：Root 当前 `active · 4/5` → **`active · 5/6`**（R6 完成后 E-021 已投影）。
- `GOAL-008.../00-meta.md` 备注：`progress: 4/6` → **`progress`（开设时 4/6，R6 关闭投影后为 5/6）**。

其余 `4/5`/`4/6` 出现处（`01-decision.md` 时间线、`E-011`～`E-021`、GOAL-006 `E-001`/`E-003`、`A-002`/`A-003` 意见原文）均为**当时事实的历史记录**，按 P-003「审计意见不代改」原则保留原样。

4. Vision 层滞后**未在本轮修改**：`docs/vision/roadmap.md`（VP-037 行与三处叙述仍写 plan `v1.1.0`、R6 `done · 4/4`）与 Charter「现行组合投影」快照按 A-002 F-002 的处置计划，列入 C4 投影清单、Charter 部分由 `/vision` 同步。本轮不在确认前把愿景投影写成已关门。

## 结论

R5-I-004 仍为 `collecting`；`GOAL-006` 保持 `active · 3/4`，C4 未勾选；Root/VP/workspace 保持 `active`。本条目只提供裁决证据并修正目标内当前态表述，不放行、不关门、不闭合任何 required finding。
