# schema-ui-core Web

React 19 + Vite + TypeScript + Tailwind application shell for the MVP Admin
workspace. Manifest-driven navigation and route resolution (R3), R4 account
session + permission-engine gating, and R5 example pages.

## Requirements

- Node.js 20+
- npm with the committed `package-lock.json`

## Run

```bash
npm install
npm run dev
# http://localhost:5173  (Vite proxies /api to :8080)

npm test        # vitest run (458 tests)
npm run build   # tsc -b && vite build
```

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

Pages are Schema documents embedded by the Go API (`GET /api/schema/{pageId}`)
and rendered through the schema-driven default path
(`manifest route → loadPageDocument → RenderPage`). New/adjusted pages only
edit `apps/api/internal/handler/fixtures/schema/*.json` — the Renderer main path
stays generic (T-UI-10).

- `data-table` — list surface over `/api/records`
- `search-form-table` — search form bound to its target table query (R4)
- `list-edit-lifecycle` — R4 representative CRUD page: create/edit/delete via
  Schema actions (`actionRef` → modal forms + row DELETE with confirm); write
  affordances gate through the D-PERM engine against the boot session
  (`records.write` → edit/delete), so they are disabled for read-only `$context`
- `form-controls` — whitelisted form controls
- `form-with-reactions` — reactions + context snapshot

## 鉴权边界

- 会话为真实登录会话/Token（GOAL-005 R2），请求级身份中间件按会话权限求值；
  `/api/accounts/me` 返回 `{ user: { id, roles, permissions }, features }` 作为
  `$context` 快照。
- `/api/records` 读路由要求 `records.read`、写路由（POST/PATCH/DELETE）要求
  `records.write`；匿名 `401 UNAUTHENTICATED`、缺权限 `403 FORBIDDEN`（I-007-001）。
- 权限求值引擎（D-PERM）提供**渲染层**门禁（按钮禁用/隐藏）；**后端为权威**——
  前端隐藏不是安全边界，直接写请求仍由 API 403/401 拦截（T-UI-08/09）。

## 测试

```bash
npm test
```

- `src/protocol/` — manifest 契约与上游 fixture 一致性
- `src/renderer/` — D-PERM 权限求值、reactions、form controls、records 客户端
- `src/app/` — 外壳集成、导航投影、示例页渲染与权限拒绝路径
- `src/account/` — 会话快照加载（含失败 fail-closed 路径）
