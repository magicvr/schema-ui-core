---
title: 执行记录 · R1 · Schema 加载、校验与统一错误面
status: active
created: 2026-08-01
updated: 2026-08-01
parent: null
version: 0.3.0
---

# 执行记录 · GOAL-002

## 2026-08-01 · 立项

- 用户确认 `/govern`：按 Root D-004 创建 R1 子目标（加载 + 主路径 + 代表性 Node 页）。
- 建立本目标五件套；`parent` = `GOAL-001-production-admin-foundation`；`status` = `active`；`progress` = `0/4`。
- 范围依据：I-001 矩阵 §4 候选「Schema 加载 + 校验 + 错误面」与 Root D-005。

> 本节仅记录立项；尚未实施加载器或测试。

## 2026-08-01 · 记录决策 D-003（I-002-001 裁决）

- 用户对 `I-002-001` 裁决：页面文档由 **Go `GET /api/schema/{pageId}` 端点**提供，manifest `schemaUrl=/api/schema/*` 不变；前端加载器运行时 fetch + 浏览器 Ajv D-VAL + 统一 `PageSchemaError`。
- 记录决策 `D-003`；`I-002-001` 由 open → **verified**（实施路径选择门禁解除）。
- 本轮进入实施（见后续条目）；`progress` 仍 `0/4`，成功标准未勾选。

## 2026-08-01 · 实施：Go 端点 + 前端加载器 + 校验 + 统一错误面

**Go 侧（`apps/api`）**
- 新增 `internal/handler/schema.go`：`GET /api/schema/{pageId}`，从 `//go:embed fixtures/schema/*.json` 提供页面文档；未知 pageId → `404 SCHEMA_NOT_FOUND`；已挂入 `Register(mux)`（`health.go`）。
- 新增 fixtures：`internal/handler/fixtures/schema/overview.json`、`catalog.json`（均满足 page/node schema 结构）。
- 新增 `internal/handler/schema_test.go`：服务已知 pageId（200 + 可解码 + meta 字段）、遍历全部 seeded 文档、未知 pageId 404、非 GET 不匹配。`go vet ./...` 干净，`go test ./...` 全绿。

**前端（`apps/web`）**
- 新增 `src/protocol/conformance/runtime-schema-validate.ts`：**浏览器安全** Ajv 校验（构建期导入 pinned `docs/schemas/*`，经新增 `@schemas` 别名；与 `schema-validate.ts` 同 schema 集，不重写语义）。`validatePageDocument` 提供加载路径 D-VAL。
- 新增 `src/protocol/load-page.ts`：`loadPageDocument(page, params, { baseURL, fetcher })` → 解析 schemaUrl（含 `{param}` 展开）→ fetch（可注入）→ JSON 解析 → 结构校验 → meta.pageId 与 manifest 核对；失败抛统一 `PageSchemaError`（`code`/`url`/`issues`）。**不**切换 App 默认渲染分支（属 GOAL-003）。
- 新增 `src/protocol/load-page.test.ts`：成功加载、参数展开、404→`PAGE_NOT_FOUND`、5xx→`PAGE_LOAD_FAILED`、网络失败、非 JSON body→`PAGE_PARSE_FAILED`、结构非法→`PAGE_SCHEMA_INVALID`、pageId 不符→`PAGE_ID_MISMATCH`；`validatePageDocument` 直接用例。
- 配置：`vite.config.ts`/`vitest.config.ts` 增 `@schemas` 别名；`tsconfig.app.json` 增 `resolveJsonModule` + `@schemas/*` paths。
- 验证：`npm test` 408 项全绿（含新增 10 项）；`tsc -b` 干净；`npm run build`（vite build）成功。

**成功标准勾选**：四条均以测试证据勾选；`progress` `0/4` → `4/4`（见 00-meta）。`I-002-002`（non-blocking）保持 open，验收前对齐 ManifestError/RenderError 复核。

**边界与诚实表述**
- Go 仅 seeded `overview`/`catalog` 两个最小页面文档；manifest 其余 pageId 返回 404（合法错误路径）。代表性页面全集属 GOAL-004。
- 加载器错误结构可被宿主层捕获；**UI 渲染错误面**尚未接线，属 GOAL-003（默认主路径）。
- 本目标不勾选 Root R1 检查点；未做阶段/关门审计（下一步提议）。

## 2026-08-01 · 关门自审 A-002 与合并意见响应

- 用户 `/govern`：自审，然后合并响应自身与独立审计意见，准备关门。
- 复核证据（全绿）：`npm test` 408/408（含 `load-page.test.ts` 10/10）；`npm run build`（vite 生产构建，含 `@schemas` 构建期导入）成功；`go test ./...` 全绿；`go vet ./...` 干净；manifest `schemaUrl=/api/schema/<pageId>` 与端点契约一致。
- 追加 `03-audit.md` **A-002**（source=self，close-out，verdict=pass）：四项成功标准全部有证据；F-002（与独立审 F-001 收敛的跨边界 fixture 漂移）与 F-003（防御性 `fetcher` 分支未测）为 recommended/low open，不阻断关门。
- 追加「编排响应」节：合并响应 A-001（independent）+ A-002（self）；`I-002-002` 验收前复核完成 → **concluded**（错误包络与 ManifestError/RenderError 约定一致，`PAGE_*` 域内专用码不冲突）。
- 状态仍 `active`（`progress` `4/4` 不变）；关门条件已满足，**待用户确认**后置 `done` 并同步 goal-tree。

## 2026-08-01 · 用户确认置 done

- 用户 `/govern` 确认「OK 采用（置 done + 推荐项处置）」。
- `GOAL-002` 状态 `active` → `done`（2026-08-01）；`progress` 仍 `4/4`。
- F-001 / F-002 / F-003 按 **accepted-residual** 合法闭合（用户书面确认）：F-001/F-002 → follow-up，随 R1 集成 / GOAL-004 补跨边界用例；F-003 → 接受防御性残余（浏览器不可达）。
- 同步 `goal-tree.md`（树 + 状态表）。Root R1 检查点**不**由本目标 `done` 自动勾选（R1 需 002+003+004）。

