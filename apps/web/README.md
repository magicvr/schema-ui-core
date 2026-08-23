# schema-ui-core Web

React 19 + Vite + TypeScript + Tailwind application shell for the MVP Admin
workspace. Manifest-driven navigation and route resolution (R3), R4 account
session + permission-engine gating, and R5 example pages.

## Requirements

- Node.js 22+
- npm with the committed `package-lock.json`

## Run

```bash
npm install
npm run dev
# http://localhost:25173  (Vite proxies /api to :25080)
# Override with WEB_PORT (dev web) and/or HTTP_ADDR (API) when another port is
# needed, e.g. $env:WEB_PORT=3000; npm run dev
npm test        # vitest run
npm run test:e2e        # Playwright Chromium, sqlite dialect (default)
npm run test:e2e:postgres  # same suite against a dedicated scratch PostgreSQL
# Bash: run the runtime profiles against the same Web code (demo = non-production)
# e2e profile follows apps/api/configs/config.yaml (T-06)
# Default ports (API :25080 / web :25173) are above Windows Hyper-V excluded
# ranges; override via HTTP_ADDR / WEB_PORT if a port is taken.
npm run build   # tsc -b && vite build
```

### 浏览器 E2E 双数据库方言（W24 / GOAL-035）

收尾阶段浏览器套件**按方言各跑一次**：sqlite（默认）与 postgres。挂具在
`playwright.config.ts` 显式声明 `DB_DIALECT` 契约，并在 `e2e/global-setup.ts`
启动后校验契约（sqlite：临时 DB 文件出现；postgres：scratch 库已迁移），
违反时立即失败并给出诊断——**本地 `apps/api/configs/.env` 无法再静默改道**
（W23 N-001 的根因正是它把挂具指向共享开发库）。

- `npm run test:e2e`：每轮全新临时 SQLite（`DB_PATH`），确定性种子。
- `npm run test:e2e:postgres`：经 `apps/api/cmd/e2e-pgset` 自动
  create → run → drop 专用 `schema_ui_e2e_*` 库；连接凭据取自进程 env 或
  `apps/api/configs/.env` 的 `DB_*`（与 API 自身加载一致），要求该用户有
  CREATEDB 权限。遗留库可用
  `go run ./cmd/e2e-pgset drop <name>`（cwd: `apps/api`）清理。
- CI（`r6-basic-matrix.yml` `browser-e2e`）按 `profile × dialect` 矩阵全量覆盖。

### 生产（compose / nginx · GOAL-008 S2 / I-008-001）

- 生产形态由仓库根 `compose.yaml` 提供第二启动路径：`apps/web/Dockerfile` 多阶段构建（`node:22` → `nginx:alpine`），`nginx.conf` 服务 `dist/` 并做 SPA fallback + `/api` 反代到 `api` 服务；同源免 CORS。
- SPA 使用相对路径 `/api/*`（`auth-client.ts`），故生产**必须**同源反代（`location /api → proxy_pass http://api:25080`），无需 `VITE_` API base。
- 完整契约见 GOAL-008 `attachments/I-008-001-engineering-contract.md`。

## Shell & session (R3/R4)

- The shell fetches the pinned app manifest from
  `/.well-known/schema-ui/app-manifest.json`, validates the 2.7 contract,
  projects navigation, and resolves History API routes.
- On boot, the app loads `GET /api/accounts/me` and attaches the `$context`
  snapshot (`{ user, features }`) to the shell. Navigation permission checks
  and row-action gates evaluate against this real identity. If the session
  load fails, the shell renders a non-blocking notice instead of silently
  dropping the error (fail-closed: navigation and actions stay restrictive).

## Example pages (schema-driven, R1/R4)

Pages are Schema documents owned by a Go module (`GET /api/schema/{pageId}`)
and rendered through the schema-driven default path
(`manifest route → loadPageDocument → RenderPage`). New/adjusted pages add a
document under the owning module schema package (core examples are in
`apps/api/internal/modules/schemarender/schema/`) and a Provider contribution;
the Renderer main path
stays generic (T-UI-10).

- `data-table` — list surface over `/api/users`
- `search-form-table` — search form bound to its target table query (R4)
- `users` — GOAL-011 representative CRUD page: create/edit/delete via Schema
  actions (`actionRef` → modal forms + row DELETE with confirm); write
  affordances gate through the D-PERM engine against the boot session
  (`users.write` → edit/delete), so they are disabled for read-only `$context`
- `roles` — GOAL-011 second semantic resource CRUD page (roles.write gate)
- `form-controls` — whitelisted form controls
- `form-with-reactions` — reactions + context snapshot

## 鉴权边界

- 会话为真实登录会话/Token（GOAL-005 R2），请求级身份中间件按会话权限求值；
  `/api/accounts/me` 返回 `{ user: { id, roles, permissions }, features }` 作为
  `$context` 快照。
- `/api/users`、`/api/roles` 读路由要求 `users.read`/`roles.read`、写路由要求
  `users.write`/`roles.write`；匿名 `401 UNAUTHENTICATED`、缺权限 `403 FORBIDDEN`
  （GOAL-011 · I-011-001 §6）。
- 权限求值引擎（D-PERM）提供**渲染层**门禁（按钮禁用/隐藏）；**后端为权威**——
  前端隐藏不是安全边界，直接写请求仍由 API 403/401 拦截（T-UI-08/09）。

## 测试

```bash
npm test
```

- `src/protocol/` — manifest 契约与上游 fixture 一致性
- `src/renderer/` — D-PERM 权限求值、reactions、form controls、资源 transport 客户端
- `src/app/` — 外壳集成、导航投影、示例页渲染与权限拒绝路径
- `src/account/` — 会话快照加载（含失败 fail-closed 路径）
