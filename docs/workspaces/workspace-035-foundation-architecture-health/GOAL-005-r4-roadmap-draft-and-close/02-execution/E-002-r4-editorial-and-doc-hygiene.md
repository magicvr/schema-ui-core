---
doc_type: goal-execution
record_id: E-002
id: E-002-r4-editorial-and-doc-hygiene
doc: execution-entry
status: recorded
parent: GOAL-005-r4-roadmap-draft-and-close
date: 2026-09-10
created: 2026-09-10
updated: 2026-09-10
version: 0.1.0
---

# E-002 · R4 editorial 执行与文档卫生

## 已发生事实（2026-09-10）

1. **用户 editorial 裁决**：书面采纳草案第 5 节**全部 10 项** `docs/vision/**` 改动。
2. **`/vision` editorial 落盘**（提交 `2b511aa0`）：
   - `docs/vision/roadmap.md` → v0.80.0：现状锚点改为「SQLite 小连接池（默认 4）/PG 双方言 + S3 + 进程内 Job/EventBus + 内存限流 + 可选指标/traces + JWT 轮换」；RT-P04 现状锚点同步（池化已交付、读写分离/replica 仍 gated）；RT-D02 自相矛盾表述删除（保持 delivered）；RT-K03 与 VP-016 行的 mfa-wrap 表述改为现行事实；**A0–A7 重述为「已交付序列 + 唯一未触发项 A3」**并新增未立项候选 **C1**（`timestamptz`）；组合投影同步（workspace-035 3/4、Admin 最近一拍、当前组合焦点）。
   - `docs/vision/charter.md:73`、`docs/vision/workspaces.md`（VP-016 行与 workspace-016 说明、workspace-035 投影）、`docs/vision/plans/VP-016-key-rotation-and-backup.md`（追加 2026-09-10 注记，**历史原文不回改**）、`docs/vision/plans/VP-035-foundation-architecture-health.md`（工作区绑定 notes + 规划修订短史 3 行）。
   - [VR-075](../../../vision/revisions.md) 登记；[VRev-088](../../../vision/reviews/VRev-088-vp035-roadmap-restatement-editorial.md) self `pass`（open required = 0），分类判定 = **editorial**。
3. **G-001 文档卫生执行**：`docs/architecture/overview.md` → v0.11.0（仓库布局表拆为 `apps/api` + `apps/web` 并加历史说明；「当前阶段」改为带日期的指针式表述；Charter `@0.2.0`→`@0.4.0`；当前交付 VP 由 VP-005 改为 VP-035；删除只到 `workspace-006` 的陈旧清单）。
4. **G-003 文档卫生执行**：`docs/architecture/cache-redis-seam-and-track.md` → v1.2.0（§2.6.1 补现行 7 方法完整面 + `RateLimiterProvider` + 精确锚点；§2.6.2 补「原子三元组在同一原语上复合」，`Reserve`/`Cancel` 的 token 结构仍留触发后裁决）。
5. **判据取证**：[exit-criteria-matrix.md](attachments/exit-criteria-matrix.md) 逐条覆盖 VP-035 六条方向级判据——判据 1～5 达成，判据 6 待 R4 审计。

## 证据

| 主张 | 路径 / 提交 |
|------|-------------|
| editorial 改动 | `docs/vision/**`（提交 `2b511aa0`）；VR-075；VRev-088 |
| G-001 / G-003 执行 | [doc-hygiene-record.md](attachments/doc-hygiene-record.md)；`docs/architecture/overview.md` v0.11.0；`docs/architecture/cache-redis-seam-and-track.md` v1.2.0 |
| 判据矩阵 | [exit-criteria-matrix.md](attachments/exit-criteria-matrix.md) |
| 边界未越 | 本阶段 `git diff --name-only -- apps` 为空；Charter `@0.4.0` 未变；无 trigger 行变化 |

## 边界

未改生产代码/端口/Profile；未新增或消耗 trigger 行；未改 Charter `vision_id@version`；未把 C1（`timestamptz`）立项；VP-016 历史原文保留。
