---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-005-r4-pg-ops
version: 0.1.0
---

# E-001 · PG external 消费实测（2026-08-29 · F-005 核销）

## 环境与执行

- docker `postgres:16` 容器 `gf-pg`（用户指令 · 127.0.0.1:15432 → 5432 · 库 `golden`）。
- golden-field 组合根双方言参数化（`-dialect sqlite|postgres` + `-dsn`；sqlite 仍为内嵌默认；CLI 模板同步待下轮）。

## 实测结果

```
首次：golden-field kernel=2.0.0 profile=admin dialect=postgres fresh=true contrib{r=10 p=2 perm=3 nav=1 frag=1}
重入：golden-field kernel=2.0.0 profile=admin dialect=postgres fresh=false contrib{…一致}
库内：site_settings 种子行存在（psql count=1）
```

- 迁移台账（含 0063）在 **PG 方言从零 apply**（fresh=true）+ **幂等重入**（fresh=false）——与 SQLite 行为一致（双方言契约平等实证）。
- **workspace-022 F-005（PG external 消费）→ 核销**（回填该 A-001）。

## 知识

- PG 生产权威方言的 external 消费与 SQLite 走同一 `assembly.OpenStore` 端口路径——方言差异（BIGINT updated_at 等）由模块迁移 ApplyPostgres 处理（主仓契约既有），下游无感知。