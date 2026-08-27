---
date: 2026-08-23
scope: GOAL-036 S5（防复发机制实施）
---

# E-002 · 防复发机制实施（S5）

## 实施（文件级证据）

**1. store 连接面回归测试** — `apps/api/internal/store/store_wal_test.go`（新增）：
- `TestSQLiteDSNPragmas`：`:memory:`/`file::memory:` 保持原样（单连接前提）；文件 DSN 必须携带 `_busy_timeout=5000` / `_journal_mode=WAL` / `_synchronous=NORMAL` 三元组；已有 query 的 DSN 以 `&` 追加；
- `TestFileStoreWALPoolAndPragma`：真实文件库打开后白盒断言 `MaxOpenConnections == sqlitePoolDefault`、`journal_mode == wal`、`busy_timeout == 5000`、`synchronous == 1`（NORMAL）；
- `TestMemoryStoreStaysSingleConnection`：`:memory:` 必须保持 1 连接；
- `TestFileStorePoolOverride`：`OpenOptions.PoolMaxOpenConns` 覆盖生效（7）。

**2. 渲染层回归测试扩展** — `apps/web/src/renderer/render.test.tsx`：
- 新增「statCard + chart 共享同一 dataSource → 1 次网络请求」（合并机制覆盖全部展示节点形态）；
- 新增「`refreshList` 只重拉目标 dataSource」：注册测试控制器组件（`perf-refresh-controller-test`，对应 monitoring-auto-refresh 行为），断言 `/api/status` 1→2 次、`/api/other` 保持 1 次。

**3. 渲染层定向刷新机制** — `apps/web/src/renderer/render.tsx`：
- `SchemaCrudValue` 新增 `refreshList(dataSource)` + `listRefreshToken(dataSource)`；provider 持 per-URL 刷新 token Map；
- `useDisplayData` 增加 `targetedRefreshToken` 依赖——token 变化只重跑该 URL 的取数 effect。

**4. monitoring 定向刷新** — `apps/web/src/components/monitoring-auto-refresh.tsx`：
- 轮询 tick 由 `crud.reloadList()`（整页重拉波，6+1+1 请求/次）改为 `crud.refreshList("/api/system-monitoring/status")`（6 张 statCard 合并后 1 请求/次）；事件表随手动刷新更新（D-002 留痕）。

**5. schema 组件注册校验** — `apps/web/src/renderer/custom-components.schema.test.ts`（新增）：
- 遍历 `apps/api/internal/modules/*/schema/*.json`，收集所有渲染层 `{type:"custom", component:"…"}` 节点（动作层 `handler` 型 custom 排除），断言在注册表中存在；测试内副作用导入 9 个组件模块（镜像 `main.tsx`）填充注册表；当前 7 个引用（wallet-ensure / mfa-manager / account-session-toolbar / activity-export / data-permission-scopes / notification-center / monitoring-auto-refresh）全部命中。

**6. 开发规范** — `docs/architecture/module-contribution-playbook.md`：
- 新增 **§6 页面数据面性能规范（Schema 页面 MUST · W25）**：展示节点同源复用允许（合并机制兜底）/ 定时刷新用定向刷新 / 自定义组件禁止挂载即写 + 整页 reloadList（探活后写契约）/ 新组件必须注册并有行为回归测试；修订表 +1（1.1.0）。

## 回归结果（事实）

- `go test ./internal/store/ -count=1`：通过（含 4 个新用例；`-count=1` 实测两轮：首轮暴露 `Open` 缺显式 `Dialect` 与 `synchronous` 数值化读回，已修正后绿）。
- 定向 vitest：`render.test.tsx` / `custom-components.schema.test.ts` / `wallet-ensure` / `schema-table` / `schema-crud` 全绿（98+1+29 用例，首轮注册校验因测试未导入组件模块而红——镜像 main.tsx 后绿）。
- 全量回归与 `tsc -b` 于下轮补跑（E-003）。

## 剩余（S6，not done）

- I-001：Playwright e2e（admin/mvp × sqlite）未跑；**监控页定向刷新为新增行为，e2e 断言面需过一遍**。
- I-002：活栈体感/请求计时复核未做。
- 自审 A-001 与 goal-tree/workspace.md 终态同步未做；按用户指示不闭门。