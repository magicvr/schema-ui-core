---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-go-library-consumption
version: 0.1.0
---

# E-004 · 方案 β 装配工厂实证（2026-08-29）

## 产出

1. **`apps/api/assembly` 包（v0.1.0 experimental · B+ 层）**：三个导出工厂，签名全部公开类型：
   - `OpenStore(ctx, dialect kernel.Dialect, path, dsn string, catalog []kernel.MigrationContribution) (kernel.Store, error)` — 内部 `store.Open`（迁移随 Open apply + checksum 校验）
   - `NewAuthenticator(secret []byte, accessTTL, refreshTTL, runner authsession.TxRunner) *auth.Authenticator` — 返回 internal 类型但调用方类型推断消费
   - `NewMailSender(st kernel.Store, retentionCap int) kernel.MailSender`
2. 机制确认：`kernel.Store` 的 `Run` 与 `authsession.TxRunner` 结构同构 → Store 直接满足 TxRunner（无需包装）。

## 黄金下游仓全链验证（`golden.local/consumer`，SQLite 临时库）

```
kernel=2.0.0 profile=admin dialect=sqlite fresh=true contrib{routes=10 pages=2 perms=3 nav=1 frag=1} users_module=admin.users
```

- 迁移台账：`compiled.PersistenceCatalog()`（B 层）→ `assembly.OpenStore`（迁移从零 apply，`WasFresh()=true`）
- 装配链：Authenticator（类型推断，无 internal 命名）→ `authsession.NewRepository` / `operationlog.NewRepository`（B 层构造，Store 直作 TxRunner）→ `users.New(...)` → `RegisterContributions`
- 贡献计数与 `users.Provider.Descriptor` 声明逐项一致（10 routes / 2 pages / 3 perms / 1 nav / 1 frag）→ 功能基线与 fork 形态等价（贡献面）
- 主仓回归：`go build ./...` exit 0（新增 assembly 包无破坏）

## 残余（有界）

- **PG 方言的 external 消费**：本环境无本地 PostgreSQL 服务，golden-consumer 未跑 PG 路径；PG 方言实现与 SQLite 走同一 `store.Open` 端口语义，主仓 PG 测试全绿（`internal/store` postgres_test 等）。**复审触发** = R4 零冲突演练（可选用 PG 侧）或 R5 发布回归。
- assembly 仍 experimental：签名冻结于冻结面 v1.1.0；工厂集最小化（store/auth/mail），jobs/objectstore/obs 等后续按需增补。