---
id: E-001
doc: execution-entry
goal: GOAL-006-r5-dual-path-acceptance
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-001 · R5 立项与方案

## 2026-08-20 · R5 建立

### 已发生事实

- Root R1–R4 全部 done（GOAL-002..GOAL-005）。创建 GOAL-006 五件套；R5 方案写入 D-001（U0 基线固化 → U1 升级策略 I-001 → U2 备份合同 I-004 → U3 共事务/readyz → U4 关门 → Root 5/5）。
- 既有双路径证据可作 U0 基线：sqlite `go test ./...` 0 FAIL；live PG `TestFullCatalogPostgresBootstrapIntegration` + `TestCompositionPostgresStartup` 全绿。

### 证据

| 主张 | 路径 / commit |
|------|---------------|
| GOAL-006 建立 + 方案 | `docs/workspaces/workspace-013-store-dialects/GOAL-006-r5-dual-path-acceptance/`（D-001） |
| 前提 R1–R4 done | Root `00-meta.md` progress 4/5；GOAL-002..GOAL-005 `00-meta.md` done |
| U0 基线 | `apps/api`: `go test ./...` 0 FAIL；store `TestFullCatalogPostgresBootstrapIntegration`、composition `TestCompositionPostgresStartup` |
