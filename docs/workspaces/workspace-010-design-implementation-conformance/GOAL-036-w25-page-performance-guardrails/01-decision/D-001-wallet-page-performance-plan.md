---
source: self
date: 2026-08-23
scope: GOAL-036 方案冻结（S2）
verdict: pass
---

# D-001 · 我的钱包性能诊断结论与修复方案取舍

## 诊断结论（S1，事实证据见 E-001）

一次访问 `/my-wallet` 的完整数据链路为 **10 个请求**（1× schema + 1× POST /me + 8× GET，其中后 4 个来自 wallet-ensure 成功后的 `reloadList()` 整页重拉波），每个请求经认证中间件 DB 查用户 + handler 事务，全部在 **SQLite 单连接**（`MaxOpenConns=1`）上串行排队，每次 COMMIT 伴随文件 fsync。四个放大因素：

1. **后端串行 + fsync**：`apps/api/internal/store/store.go` `SetMaxOpenConns(1)` → 池空转排队；无 WAL、`synchronous=FULL` → Windows/杀软下逐提交 fsync，单次页面 ≈ 20 个事务 × fsync。
2. **前端同 URL 重复请求**：三个 statCard 的 `dataSource` 均为 `/api/wallet/me`，`useDisplayData` 各自独立 fetch，无共享合并（`render.tsx`），一次页面 3× 同一 GET。
3. **wallet-ensure 无条件 POST + 整页重拉**：`wallet-ensure.tsx` 每次挂载无条件 `POST /api/wallet/me` 并 `reloadList()` → 追加一整波 4 个 GET；钱包已存在时纯属浪费。
4. **schema 每次导航重取 + 全量 D-VAL**：`load-page.ts` 无缓存，每次访问重新 fetch `/api/schema/my-wallet` 并跑结构校验。

索引与查询本身健康（`idx_wallet_ledger_account(account_id, created_at DESC)`、账户唯一约束），数据量不是当前瓶颈。

## 方案（已实施）

| 因子 | 修复 | 落点 |
|------|------|------|
| A | 文件库连接池 4（可配 `DB_POOL_MAX_OPEN` 覆盖）+ DSN `_busy_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL`；`:memory:` 保持单连接 | `store.go` / `open.go` / `options.go` / `config.go` |
| C | CRUD provider 新增 `fetchList`：**仅 in-flight 合并**（同 URL 并发消费者共用一个请求），reload/换 transport 清空 | `render.tsx`（statCard/chart/table 共用） |
| B | wallet-ensure 经 `fetchList` 与 statCards **合并同一个 GET /me 探活**；探到 `WALLET_NOT_FOUND`（404）才 POST 并刷新，钱包存在时零 POST、零重拉 | `wallet-ensure.tsx` |
| D | schema 文档按 resolved URL 内存缓存（App shell 实例级、上限 64、失败不缓存、opt-in） | `load-page.ts` / `App.tsx` |

## 取舍与未选方案

- **不做长期 URL 结果缓存（拒）**：只合并 in-flight。理由：查询变更/重置/重载必须拿到新数据（既有 `search-form-filters` T-02 语义：reset 后仍须重新拉取）；长期 memo 会命中过期结果。性能收益主体（首屏突发 3→1）来自 in-flight 合并，不损失。
- **provider 状态优先运输器（拒）**：`fetchList` 增加显式 `transport` 参数，调用方传入自己的当前运输器。理由：测试经 `tableRenderer` 注入 fixture、生产注入 authFetch；若改由 provider state 优先，首帧（注册未完成）会用错运输器，破坏既有注入语义（schema-crud 全组暴露）。缓存键只按 URL，transport 只决定谁执行请求。
- **schema 缓存为显式 opt-in**：`loadPageDocument` 仅在有 `cache` Map 时启用，测试/直调路径不受模块级状态污染。
- 未改动 statCard 的 `pageSize=100` 请求形态（RPC 契约冻结面），仅新增共享查询常量 `DISPLAY_LIST_QUERY` 供探活复用同一 URL。