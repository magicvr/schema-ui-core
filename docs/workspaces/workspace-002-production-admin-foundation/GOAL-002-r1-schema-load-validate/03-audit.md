---
title: 审计台账 · R1 · Schema 加载、校验与统一错误面
status: active
created: 2026-08-01
updated: 2026-08-01
parent: null
version: 0.4.0
---

# 审计台账 · GOAL-002

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | independent | 2026-08-01 | 执行事实、成功标准、P-005 信息项与工作区边界 | pass | 已出具；无开放 required finding |
| A-002 | self | 2026-08-01 | 关门自审：成功标准、D-VAL 串联、统一错误面、测试证据、I-002-001/002 | pass | 已出具；无开放 required；F-002/F-003 recommended → accepted-residual（用户确认 follow-up） |

## A-001 · 加载、校验与统一错误面执行事实独立审计（2026-08-01）

- **source**：independent
- **auditor**：GitHub Copilot
- **类型 / scope**：execution-facts；核对 GOAL-002 的四项成功标准、`I-002-001` / `I-002-002`、当前工作区绑定与共享资料边界。
- **verdict**：pass

### 范围与区间

- 当前工作区为 `workspace-002-production-admin-foundation`，canonical root 为 `docs/workspaces/workspace-002-production-admin-foundation/`，Root 为 `GOAL-001-production-admin-foundation`；目标 parent、goal-tree 与工作区绑定一致。
- `shared_materials_catalog: none`，本次没有将共享资料作为事实或 finding 关闭证据。
- 审计仅覆盖本目标的加载、校验与错误面，不审默认渲染分支（GOAL-003）或代表性页面全集（GOAL-004）。

### 成果（有证据）

- 前端加载器 `apps/web/src/protocol/load-page.ts` 通过 `resolveSchemaUrl` 展开参数、请求页面文档、解析 JSON、调用 `validatePageDocument`，并以 `PageSchemaError` 统一暴露 `PAGE_LOAD_FAILED`、`PAGE_NOT_FOUND`、`PAGE_PARSE_FAILED`、`PAGE_SCHEMA_INVALID` 与 `PAGE_ID_MISMATCH`。
- 浏览器侧校验器 `apps/web/src/protocol/conformance/runtime-schema-validate.ts` 在构建期导入 pinned `docs/schemas` 的 page/node/action/reaction schema，并在校验失败时返回可观察 issues。
- API `apps/api/internal/handler/schema.go` 已注册 `GET /api/schema/{pageId}`，服务嵌入的 `overview` / `catalog` 页面 JSON，未知页面返回 `404 SCHEMA_NOT_FOUND`。
- 2026-08-01 的独立复核命令均成功：`apps/web` 中 `npm test -- --run src/protocol/load-page.test.ts`（10/10）；`npm run build`；`apps/api` 中 `go test ./internal/handler -run TestSchemaEndpoint -count=1` 与 `go test ./...`。

### 对照成功标准

1. 加载入口与参数模板展开：实现和定向测试均存在。
2. 加载路径结构校验与 fail-closed：校验在返回文档前执行，非法文档测试断言 `PAGE_SCHEMA_INVALID` 与 issues。
3. 统一可观察错误：错误类含 code、url、issues，定向测试覆盖 HTTP、网络、JSON、校验和 pageId 不一致路径。
4. 自动化证据：前端 10 项定向测试、API 端点测试、Web 生产构建与 API 全量测试均由本次复核通过。

### Findings

- **F-001** · recommended · low · 成功加载测试使用手工维护的 `VALID_DOCUMENT`，而非 `GET /api/schema/{pageId}` 实际返回的原始 fixture；`schema_test.go` 当前只检查 fixture 可解析及少数字段。现有证据足以支撑本目标的成功标准，但后续 fixture 或接口演进可能造成端到端契约漂移而不被该组测试捕获。建议在后续维护中增加实际端点 fixture 经加载器校验的跨边界用例。**不阻断当前目标的实施事实审计。**

### 必改项汇总

- 无开放 required finding。
- `I-002-002` 保持 open 且为 non-blocking；其错误包络对齐应在目标验收前由 `/govern` 复核，不构成当前实施证据的阻断项。

### 与既有意见的异同

- 本目标此前无正式 self 或 independent Goal Audit 意见，无冲突需要裁决。

### 结论 + 建议给编排器 / 用户的下一步

- 本 scope 可作为 GOAL-002 实施事实的正向审计证据；本意见不等同于目标关门或 Root R1 放行。
- 使用 `/govern` 汇总 A-001；在目标验收或关门前复核 `I-002-002`，并决定是否处理 F-001 推荐项。

### 声明

本意见不修改 status/progress；响应、finding 关闭和生命周期推进由 `/govern` 处理。

## A-002 · GOAL-002 关门自审（2026-08-01）

- **source**：self
- **auditor**：Claude Code · `/govern`（关门自审）
- **类型 / scope**：close-out；核对 GOAL-002 四项成功标准、D-VAL 串联、统一错误面、测试证据、`I-002-001` / `I-002-002`、工作区与共享资料边界
- **verdict**：pass
- **audit_type**：close-out

### 范围与区间

- 当前工作区 `workspace-002-production-admin-foundation`；canonical root 与 Root `GOAL-001-production-admin-foundation` 绑定一致；`parent` / goal-tree / workspace.md 核对无冲突。
- `shared_materials_catalog: none`；未把共享资料作为事实或关闭证据。
- 仅审本目标加载 / 校验 / 错误面；默认渲染分支（GOAL-003）与代表性页面全集（GOAL-004）不在 scope。
- 证据为 2026-08-01 实测：`npm test` 408/408（含 `load-page.test.ts` 10/10）；`npm run build`（vite，含 `@schemas` 构建期导入）成功；`go test ./...` 全绿；`go vet ./...` 干净。

### 成果（有证据）

- **加载入口**：`apps/web/src/protocol/load-page.ts` `loadPageDocument` 经 `resolveSchemaUrl` 展开 `{param}` 占位符，注入式 fetch → JSON 解析 → `validatePageDocument` D-VAL → `meta.pageId` 与 manifest 核对 → 返回文档或抛 `PageSchemaError`。未切换 App 默认渲染分支（属 GOAL-003）。
- **D-VAL**：`apps/web/src/protocol/conformance/runtime-schema-validate.ts` 构建期导入 pinned `docs/schemas/*`（`@schemas` 别名在 `vite.config.ts` / `vitest.config.ts` / `tsconfig.app.json` 三处一致）；Ajv `allErrors`；校验失败返回可观察 `issues`。
- **端点**：`apps/api/internal/handler/schema.go` `GET /api/schema/{pageId}` 经 `//go:embed fixtures/schema/*.json` 服务 `overview` / `catalog`，未知 pageId → `404 SCHEMA_NOT_FOUND`；`health.go` `Register` 已挂载。
- **契约保持**：manifest `schemaUrl=/api/schema/<pageId>`（`apps/web/public/.well-known/schema-ui/app-manifest.json`）未改写，符合 D-003。
- **错误面**：`PageSchemaError` 覆盖 `PAGE_LOAD_FAILED` / `PAGE_NOT_FOUND` / `PAGE_PARSE_FAILED` / `PAGE_SCHEMA_INVALID` / `PAGE_ID_MISMATCH`，含 `code` / `url` / `issues` / `message`。

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1 · 加载入口 + 路由参数模板展开 | ✅ | `load-page.ts`；`load-page.test.ts`「expands route-parameter placeholders」；`schema.go` |
| 2 · 加载路径强制结构校验 fail-closed | ✅ | `validatePageDocument` 在返回文档前执行；「fails closed … `PAGE_SCHEMA_INVALID` + issues」测试；非法文档加载器不返回 |
| 3 · 统一可观察错误结构 | ✅ | `PageSchemaError` 5 code；测试覆盖 404 / 5xx / 网络 / 非 JSON / 结构非法 / pageId 不符 |
| 4 · 自动化测试覆盖 | ✅ | `load-page.test.ts` 10/10；`schema_test.go`（已知 / 遍历 seeded / 404 / 方法不匹配）；全量 408 绿；生产构建通过 |

### Findings

- **F-002 · 成功加载测试未直接消费 Go 实际 fixture（与独立审 F-001 收敛）**
  - 严重度：low；建议：recommended
  - `load-page.test.ts` 的 `VALID_DOCUMENT` 为手工镜像 `overview.json`，`schema_test.go` 仅浅查 meta 字段与 `body.type`。fixture 或接口演进可能造成跨边界契约漂移而不被本组测试捕获。
  - 状态：closed（accepted-residual；2026-08-01 用户书面确认 follow-up——随 R1 集成 / GOAL-004 建代表性页面时补「实际 fixture → 加载器」跨边界用例；复审触发 = 新增页面 fixture 或 schemaUrl 契约变更）
- **F-003 · `fetcher` 非函数防御分支未直接测试**
  - 严重度：low；建议：recommended
  - `loadPageDocument` 的 `typeof fetcher !== "function"` 分支在浏览器中不可达（`globalThis.fetch` 恒存在），为防御性代码；当前无定向测试。
  - 状态：closed（accepted-residual；2026-08-01 用户书面确认接受防御性残余；浏览器中不可达，无运行时风险；可后续补定向测试）

### 必改项汇总

- 无开放 required finding。
- `I-002-001` = verified（D-003 裁决，证据齐全）。
- `I-002-002`（non-blocking，验收前复核）本轮完成复核，结论见「编排响应」节。

### 结论 + 建议下一步

- 四项成功标准均有可核对测试证据；无开放 required；无到期 required 信息项。
- 本 close-out 自审 verdict = **pass**。建议编排器合并 A-001（independent）与本条意见，响应 F-001 / F-002 / F-003 与 I-002-002 后，经用户确认将目标置为 `done` 并同步 goal-tree。

## 编排响应 · 合并 A-001（independent）+ A-002（self）（2026-08-01 · `/govern`）

**响应对象**：A-001（F-001）与 A-002（F-002、F-003、I-002-002）。

**对 A-001（independent）**：
- verdict `pass` 采纳；其 scope（执行事实、成功标准、P-005 信息项、工作区边界）与本次自审一致，无冲突。
- F-001（recommended / low）：与自审 F-002 收敛为同一风险（跨边界 fixture 漂移）。**处置（2026-08-01 用户书面确认）**：accepted-residual——不阻断关门；记录为 follow-up，随 R1 集成 / GOAL-004 建代表性页面时补「实际端点 fixture 经加载器校验」跨边界用例；复审触发 = 新增页面 fixture 或 schemaUrl 契约变更。

**对 A-002（self）**：
- F-002：与 F-001 合并处置（同上，accepted-residual）。
- F-003（recommended / low）：**处置（2026-08-01 用户书面确认）**：accepted-residual——防御性代码残余；浏览器中不可达，无运行时风险；可后续补定向测试。

**I-002-002 复核结论（验收前复核完成）**：
- 错误包络与 `ManifestError` / `RenderError` 约定一致（code + message + 位置锚点）；`PageSchemaError` 以 `url`（资源锚点）替代 `path`（文档内 JSON 路径）并增 `issues`（Ajv 校验明细），属 D-003 已明示设计；`PAGE_*` 为域内专用码，与 `MANIFEST_*` / `RENDER_*` 不冲突。
- 证据：`load-page.ts`、`render.ts`、`app-manifest.ts`。**`I-002-002` → concluded（closed）**。

**关闭证据表**

| F / I | source | 级别 | 状态 | 处置 / 证据 |
|-------|--------|------|------|-------------|
| F-001 | independent | recommended | **closed（accepted-residual）** | 与 F-002 合并；用户书面确认 follow-up：R1 集成补跨边界用例；复审触发 = fixture/schemaUrl 变更 |
| F-002 | self | recommended | **closed（accepted-residual）** | 同上 |
| F-003 | self | recommended | **closed（accepted-residual）** | 用户书面确认：防御性不可达分支，接受残余 |
| I-002-002 | self | non-blocking | **closed** | 复核：错误包络对齐；证据 `render.ts` / `app-manifest.ts` |

**仍开放项**：无开放 required；F-001 / F-002 / F-003 已按 accepted-residual 合法闭合（用户书面确认）。

**关门条件检查**（已通过）：
- 无未合法闭合 required finding ✅
- 无到期未处理 required 信息项 ✅（`I-002-001` verified；`I-002-002` non-blocking closed）
- 至少一次阶段 / 关门向审计 ✅（A-001 independent + A-002 self）
- 成功标准可核对 ✅
- 2026-08-01 用户书面确认 → `GOAL-002` 置 `done` 并同步 goal-tree。

## 当前审计边界

- 已有 A-001（independent，实施事实）与 A-002（self，关门自审）意见；两意见无冲突，verdict 均为 pass。
- F-001 / F-002 / F-003（recommended）均已按 **accepted-residual** 合法闭合（2026-08-01 用户书面确认），不阻断关门。
- 编排响应（合并 A-001 + A-002）见「编排响应」节；`I-002-002` 已复核关闭。
- 本目标已置 `done`（2026-08-01）。后续意见从 `A-003` 起共用序列。
