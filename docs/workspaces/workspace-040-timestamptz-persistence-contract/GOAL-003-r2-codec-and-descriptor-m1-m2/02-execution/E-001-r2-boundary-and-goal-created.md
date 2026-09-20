---
id: E-001-r2-boundary-and-goal-created
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-003-r2-codec-and-descriptor-m1-m2
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-001 · R2 边界冻结与 GOAL-003 立项

## 事实

1. **用户 2026-09-20 裁决三项**（承接上一轮「进入 R2」）：
   - **首个子目标编号与 slug = `GOAL-003-r2-codec-and-descriptor-m1-m2`**（AGENTS §11：slug 必须用户确认，禁止静默默认）；
   - **公共 wire formatter 的实施归 R3**（R2 只做 Store/持久化层，**不改 handler 的响应编码**）；
   - **Backup Port 实施归 R2**（`D-016` 第 9 项，与 `D-021` residual 复审触发同批）。
2. **Root `D-016-r2-boundary.md` 已按上述裁决修正**：范围表第 9 项标注「归 R2、M4 前完成」；新增范围修正说明（wire formatter 归 R3，故第 1–8、10 项即 R2 全部范围）；信息表 `I-041-002` 由 `open` → **`verified`**（用户裁决），`I-041-001` 由 `open` → **`collecting`**（在本子目标内定稿并由审计复审）。
3. **`GOAL-003` 立项**：五件套齐备（`00-meta.md` / `01-decision.md` / `02-execution.md` / `03-audit.md`）+ 三个 ledger 目录（`01-decision/`、`02-execution/`、`03-audit/`）+ `attachments/`。
   - **范围**：共享 codec（M1）、codec 单测（M1）、SQLite 15 descriptor（M2）、PG 显式 DDL（M2）、canonical SQL + 真实 `MigrationChecksum` 记录（M2，**即 `D-021` residual 复审触发点**）、`rebuildOperationLog` fail-closed 断言（M2）、v73+ 目录断言且 v1–v72 不变（M2）。
   - **非目标**：仓储/谓词改造与双方言回归（M3）、`postgres_test` 金额列拆分与 leftover 21 名（M3）、wire formatter（**R3**）、Backup Port（**M4 前**）、边界测试重定向（M3）、任何 v1–v72 改动、ORM/第三库/Redis/MQ。
   - **检查点 A～D**（`progress: 0/4`）：A codec + 单测；B SQLite 15 descriptor + v1–v72 不变；C PG 显式 DDL；D 真实 checksum 落盘 + **发起 independent 复审**。
4. **`goal-tree.md` 已同步**：树新增 `GOAL-003` 分支；状态表新增一行（`active · 0/4`）；说明节更新为「M3/M4 待 M1/M2 完成后立项，不预创建」。
5. **本轮 `apps/**` 未修改**——本目标只立项，尚未开始实现。

## 证据

- Root 决策：`GOAL-001/01-decision/D-016-r2-boundary.md`（含范围修正与 `I-041-002` verified）。
- 本目标五件套：`GOAL-003-r2-codec-and-descriptor-m1-m2/00-meta.md` 等四文件 + 四目录。
- 工作区树：`goal-tree.md`（树 / 状态表 / 说明三处）。
- 继承权威：Root `D-004`（模块归属）、`D-014`（未发布 baseline）、`D-015`（毫秒族已更正）、`D-016`（R2 边界）；`GOAL-002` 的 `D-017`（checksum 约定）、`D-018`（行拷贝/无过渡期）、`D-019`（F-5 + 断言 + PG 显式 DDL）、`D-020`（测试载体）、`D-021`（residual 范围与复审触发）。

## 状态评估

- **Root `progress: 1/3`**（R1 completed；R2 `active` 进行中；R3 pending）。
- **`GOAL-003` = `active · 0/4`**；无实现事实，无审计意见。
- **R2 未放行生产 schema 变更**：`I-041-001` 须在检查点 A 前定稿；`I-041-003`（PG 可执行验证环境）须在 M3 前关闭——本目标内若触及，将按 P-004 询问用户。
- **下一步**：进入检查点 A——先定稿 `I-041-001`（Go codec 公共 API 形态）并落盘为本目标首条决策，再实现 codec 与单测。
