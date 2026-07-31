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

npm test        # vitest run (398 tests)
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

## Example pages (R5)

Registered in `src/app/examples/registry.tsx` and rendered by manifest route:

- `data-table` — list surface
- `search-form-table` — search + table
- `list-edit-lifecycle` — row actions gate through the D-PERM engine against
  the boot session; Edit/Delete are disabled for non-admin `$context`
- `form-controls` — whitelisted form controls
- `form-with-reactions` — reactions + context snapshot

The schema-driven page Renderer remains a later protocol boundary; these are
direct React example surfaces.

## 鉴权边界（MVP 声明，非生产）

- 会话为静态开发会话（`/api/accounts/me`），无真实登录 / 令牌 / IAM。
- 权限求值引擎（D-PERM）提供**渲染层**门禁（按钮禁用、导航隐藏）；后端
  `/api/records` 写路由（PATCH/DELETE）挂 fail-closed 鉴权（需 admin 会话），GET 只读开放。

## 测试

```bash
npm test
```

- `src/protocol/` — manifest 契约与上游 fixture 一致性
- `src/renderer/` — D-PERM 权限求值、reactions、form controls、records 客户端
- `src/app/` — 外壳集成、导航投影、示例页渲染与权限拒绝路径
- `src/account/` — 会话快照加载（含失败 fail-closed 路径）
