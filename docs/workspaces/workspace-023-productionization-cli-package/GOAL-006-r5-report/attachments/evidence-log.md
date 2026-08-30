# 可复验证据日志（R5 · 独立审 F-006 闭合）

## 1. 从零走查计时（2026-08-29 重录 · 依赖缓存预置口径）

```
> schema-ui create walk-admin --module github.com/acme/walk-admin --dir …/walk-admin
create: 0.5s
> go mod tidy && go run ./cmd/server -dialect sqlite -dsn ./w.db
go assembly: 6.9s
> pnpm install && node probe.mjs && node probe-render.mjs && node token-check.mjs
web+probes: 1.1s
TOTAL: 8.4s
```

输出（组合根）：`walk-admin kernel=2.0.0 profile=admin dialect=sqlite fresh=true contrib{r=10 p=2 perm=3 nav=1 frag=1}`（probe 三 PASS）。

## 2. PG 迁移台账计数（docker gf-pg · postgres:16）

```
psql: SELECT count(*) FROM schema_migrations → 63
（0013 会计：全局台账 1..63 连续；0063 = site_settings_updated_at_index）
组合根：dialect=postgres fresh=true（首跑）→ fresh=false（幂等重入）
```

> 注：文档早前「64 迁移」为笔误，已统一修正为 63。