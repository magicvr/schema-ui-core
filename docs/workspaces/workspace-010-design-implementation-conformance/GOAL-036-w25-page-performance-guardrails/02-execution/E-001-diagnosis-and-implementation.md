---
date: 2026-08-23
scope: GOAL-036 S1–S4（诊断、实施、自动回归）
---

# E-001 · 诊断、实施与自动回归（S1–S4）

## S1 · 诊断（事实）

- 链路追查：schema 驱动渲染 `loadPageDocument → RenderPage → wallet-ensure / statCard(×3) / table`；后端 `GET /api/wallet/me`、`GET /api/wallet/me/entries`、`POST /api/wallet/me`（`wallet_self.go`）。
- 关键证据：`store.go:49 SetMaxOpenConns(1)`；`render.tsx` `useDisplayData` 每组件独立 fetch（无共享层）；`wallet-ensure.tsx` 无条件 POST + `reloadList()` 整页重拉；`load-page.ts` 每导航重取 schema + 全量 D-VAL；索引健康（`migration.go` `(account_id, created_at DESC)`）。
- 结论：请求数 10/次 ×（认证查用户 + handler 事务）× 单连接串行 × 逐提交 fsync = 慢的主链；见 D-001。

## S2/S3 · 实施（文件级证据）

**后端（Fix A）**
- `apps/api/internal/store/store.go`：新增 `sqlitePoolDefault=4`、`sqliteDSNParams`、`sqliteDSN()`（`:memory:` 原样返回）；`open()` 改收 `OpenOptions`，文件库按池 4（`DB_POOL_MAX_OPEN>0` 时覆盖）。
- `apps/api/internal/store/open.go`：sqlite 分支传 `opts`。
- `apps/api/internal/store/options.go`、`internal/config/config.go`：注释/语义同步（sqlite 现使用 `PoolMaxOpenConns`）。

**前端（Fix B/C/D）**
- `apps/web/src/renderer/resource.ts`：新增 `resourceListURL()`（缓存键与线上请求单一事实源）、`DISPLAY_LIST_QUERY`（statCard/chart 与探活共享）；`fetchResourceList` 复用 URL 构造。
- `apps/web/src/renderer/render.tsx`：`SchemaCrudValue.fetchList`（in-flight 合并，`transport` 显式参数；reloadList / transport 变更时清空 in-flight）；`useDisplayData` 走 `fetchList`。
- `apps/web/src/renderer/schema-table.tsx`：表格 fetch 走 `fetchList`。
- `apps/web/src/components/wallet-ensure.tsx`：改为 GET 探活（经 `fetchList` 与 statCards 合并）→ 仅 `WALLET_NOT_FOUND`（404，`wallet.go:511` 确认码）时 POST + `reloadList()`。
- `apps/web/src/protocol/load-page.ts`：`LoadPageOptions.cache`（opt-in Map，上限 64，失败不缓存）。
- `apps/web/src/app/App.tsx`：shell 持有 `schemaDocumentCache` 并下传。

**请求数对比（钱包已存在时，一次访问）**：10（1 schema + 1 POST + 3 GET /me + 1 GET entries + 4 重拉波）→ **3**（schema 缓存后 + 1 合并 GET /me + 1 GET entries）。首开用户额外 1 POST + 1 刷新波（必要）。

## S4 · 自动回归（事实）

- `go build ./...` ✓；`go test ./internal/store/... ./internal/modules/wallet/... ./internal/handler/... ./internal/composition/... ./internal/config/...` 全 ok。
- vitest 全量 **76 文件 / 1093 测试全绿**；新增回归测试：`render.test.tsx`（三 statCard 同 URL → 1 请求）、`load-page.test.ts`（缓存命中/失败不缓存 2 例）、`wallet-ensure.test.tsx`（重写 4 例：存在→无 POST/无 reload；404→POST+reload；创建失败/探活失败→CTA）。
- 过程中修复一处设计偏差：初版含长期结果 memo，被 `search-form-filters` T-02（reset 需重新拉取）否决 → 收敛为纯 in-flight 合并（D-001 取舍留痕）。
- `tsc -b` 0 错误。

## 剩余（S5，not done）

- I-001：Playwright e2e（admin/mvp × sqlite）未跑。
- I-002：活栈体感/请求计时复核未做。
- 自审与 goal-tree/workspace.md 终态同步未做；**按用户指示不闭门**。gofmt 检查说明：仓库文件为 CRLF 行尾，`gofmt -l` 全仓库命中，属既有现象，本次未引入格式化偏差。