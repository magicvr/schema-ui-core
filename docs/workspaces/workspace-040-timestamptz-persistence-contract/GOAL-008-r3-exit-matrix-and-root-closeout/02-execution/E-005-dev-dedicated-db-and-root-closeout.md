---
id: E-005-dev-dedicated-db-and-root-closeout
doc: execution-entry
status: active
parent: GOAL-008-r3-exit-matrix-and-root-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-005 · 检查点 C：用户确认关门 + dev 环境改指专用库（实测）

## 事实

### 1. Root 关门（判据 6）

- 用户 **2026-09-21** 书面确认关闭 Root（`01-decision/D-001-root-closeout-user-confirmation.md`；Root `D-020`）。
- 处置：`GOAL-008` → `done · 3/3`；Root `GOAL-001` → `done · 3/3`；Root 六条成功标准**勾选**（勾选依据 = 用户确认 + 判据 1–5 的逐条证据矩阵 + 跨目标开放 required = 0）；`goal-tree.md`（树/表/纲领段/说明/版本 1.0.0）、Root `00-meta.md`、`docs/vision/roadmap.md` 投影、`workspace.md`（v0.2.0）同步。

### 2. dev 环境改指专用库（用户裁决）

- `apps/api/configs/.env`（**gitignored**）`DB_NAME`：`postgres` → **`schema_ui_dev`**；其余键（`PG_TEST_*`、`ADMIN_PASSWD`、mail/telegram 键）**不变**。
- 在常驻实例创建 `schema_ui_dev`。

### 3. 实测（用户命令、无环境覆盖）

```text
cmd /c "set API_PORT=25083&& set WEB_PORT=25175&& dev.cmd start --no-browser"
  ok   API listening on :25083
  ok   API ready (/readyz 200)
  ok   Web listening on :25175
  exit=0

GET http://127.0.0.1:25083/readyz → 200 {"status":"ok","timestamp":"2026-09-21T01:41:27.014980Z",...}
GET http://127.0.0.1:25175/       → 200
cmd /c "... dev.cmd stop"          → stopped services (:25083 / :25175)
```

- **专用库迁移结果**：`schema_ui_dev` 的 `schema_migrations` = **max 87 / 87 行**（迁移链在该库完成）。
- **旁证**：`/readyz` 的 `timestamp` 为 R3-A 的规范 fixed-6 形状（`…T01:41:27.014980Z`）。

### 4. 观察（不属本次改动，供用户知悉）

- 共享常驻实例的**维护库 `postgres`** 上存在 `public.schema_migrations`：这是 `PG_TEST_DB=postgres` 配置下测试路径（如 `internal/composition` 的 PG 启动测试）在维护库上跑迁移的既有结果，**不是**本次 dev 改动引入；本次已把 **dev** 路径从该库移开。若也要让测试不再落在维护库，需另行调整 `PG_TEST_DB`（属用户测试环境决策，未擅自更改）。

## 证据

- `01-decision/D-001-root-closeout-user-confirmation.md`；Root `D-020`。
- `goal-tree.md`（全部 8 目标 `done`）、Root `00-meta.md`（判据勾选 + `progress: 3/3`）、`workspace.md`、`docs/vision/roadmap.md`。

## 进度评估

检查点 C 完成 → `progress: 3/3`；Root **`done`**。工作区实现层工作全部收口；VP-040 波次关闭 / Vision Review 属决策层（`/vision`），不在本次范围。
