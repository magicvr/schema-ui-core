---
doc_type: goal-attachment
id: r4-profile-route-evidence
status: recorded
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# R4 · Profile × Permission × Route 证据矩阵

## 1. Profile 分母与 provider oracle

来源：当前仓库 `apps/api/kernel/profile.go`、`apps/api/configs/config.yaml`、模块 Schema corpus，以及 `apps/web/src/app/searchable-profile-matrix.test.ts`。

| profile | 目标页数（含 notifications Shell） | 可收录 action 数 | 测试结果 | 备注 |
|---|---:|---:|---|---|
| `mvp` | 5 | 6 | pass | users/roles/account/dashboard + notifications；users-invites inner page 不进页面分母 |
| `admin` | 18 | 16 | pass | admin 17 Manifest leaves + notifications；dictionary-entries/wallet-entries/task-runs 等 inner/row scope 排除 |
| `demo` | 13 | 6 | pass | mvp + dev.examples；admin-list-batch selection trigger 排除 |
| 当前 custom（admin + channel.telegram + biz.digital-offer） | 22 | 18 | pass | admin + Telegram entry + digital-offer create；不含 dev.examples |

`searchable-profile-matrix.test.ts` 使用当前模块 Schema 文件的 `meta.pageId`/`actions`，四个 profile Manifest 构造严格对齐 R1 §2.1，并断言动态/inner page、batch selection、rowless 约束与 provider error = 0。

## 2. Permission / feature / route 矩阵

| 场景 | 预期 | 证据 |
|---|---|---|
| `menu_*` feature 为 true + 页面 Manifest leaf | page item 可见 | `searchable.ts` 的 `projectNavigation` provider；`searchable-profile-matrix.test.ts` |
| `menu_*` 缺失/false | 页面和其 action 不进入 provider | `searchable.test.ts` denied/visible projection；既有 `navigation.test.ts` feature fail-closed |
| 页面 action 的 `permissionIntent` / 局部 `permissions` 通过 | action 可收录，执行复核同一 target/cascade | `programmatic-action-gate.test.tsx` allowed modal/custom；App integration New user modal |
| 页面 action permission denied | 页面可保留，但 action 不收录；程序化 invoke 不发请求 | `searchable.test.ts`、`programmatic-action-gate.test.tsx` |
| page selection | 只走已存在的 Manifest route 与 App `onNavigate` History API | `command-palette-app.test.tsx`、E2E Users route |
| active grouped route | 既有 `item.active` 逻辑自动展开父组 | `command-palette-app.test.tsx` Admin group `aria-expanded=true` |
| navigate action 含未绑定 `{…}` | 不调用 host | `programmatic-action-gate.test.tsx` |
| row navigation mapping 非标量/缺失 | feedback，绝不调用 host | `programmatic-action-gate.test.tsx` |
| 未知/malformed context expression | fail-closed | `app-manifest.test.ts`、`apps/api/internal/account/permission_test.go` |
| stale disabled-profile DB/menu rows | 不作为 denominator | provider 只读当前运行时 Manifest；Profile 事实见 R1 matrix |

## 3. 浏览器与自动化证据

新建 `apps/web/e2e/command-palette.spec.ts`，每次 2 tests（visible Users search + New user modal；mobile trigger + Escape）：

| profile | dialect | result |
|---|---|---|
| mvp | SQLite | 2 passed |
| admin | SQLite | 2 passed |
| mvp | Postgres scratch | 2 passed |
| admin | Postgres scratch | 2 passed |

真实 API + Vite proxy、登录/forced-password seed、authenticated Schema transport、SQLite isolation 与 Postgres scratch teardown 均由既有 Playwright harness 处理；没有复用开发库或输出凭据。

## 4. 非目标/红线核对

- `SearchableProvider` 只读取 Manifest projection 与认证 page Schema，未读取 API 实体行或数据库表；无 entity full-text / RT-X01 / RT-X02 / Redis / MQ / multi-instance。
- pinned `AppManifest` envelope 未增加 provider/action/profile 字段；没有修改协议版本或 Manifest schema。
- 未增加 Saved Views、recent/pinned persistence、batch results、unsaved protection、Toast global rewrite 或第二业务域。
- action selection 只传 page owner + mount trigger 给既有 executor；后端 401/403 仍是最终授权边界。

## 5. 验证命令结果

- `node node_modules/typescript/bin/tsc -b --pretty false`：exit 0。
- `node node_modules/vitest/vitest.mjs run --reporter=dot`：**104 files / 1366 tests passed**。
- `node node_modules/vite/bin/vite.js build`：exit 0；仅有既有 chunk size warning。
- `go test ./...`（`apps/api`）：exit 0。
- `pnpm test` 包装命令因本机 `ERR_PNPM_IGNORED_BUILDS`（esbuild build script 未批准）在 install 前退出；直接 Vitest/tsc/Vite/Go 结果作为独立证据，限制已在 execution/audit 留痕。
