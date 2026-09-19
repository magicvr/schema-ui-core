---
id: E-001-w34-implementation
doc: execution
parent: GOAL-046-w34-batch-export-availability
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
status: recorded
---

# E-001 · W34 实施与证据

## 事实（2026-09-19）

### 1. 诊断与复现（`I-046-001`）

以 `configs/config.yaml` **原样**启动 API（`CONFIG_FILE` 指向该文件副本，仅把 `app.env` 设为 `development` 并把 DB 指向临时 SQLite，避免触碰开发者库），逐项观测：

- `POST /api/jobs/batch-export`（admin Bearer）→ `404` + `{"error":"NOT_FOUND","message":"未找到","messageKey":"error.notFound"}`。
- `GET /api/jobs` → `404 NOT_FOUND`（同一 token）。
- `GET /api/export/users`（同步导出对照）→ 正常，证明缺的只有 `admin.jobs`。
- 改以 `profile: admin` 启动：同一请求 → `403 MUST_CHANGE_PASSWORD`（**路由存在**，先撞强制改密门）；完成改密后 → `202 queued`，轮询至 `succeeded`，`GET /api/jobs/{id}/result` 返回 200 与可下载 CSV。
- 以 `profile: mvp` 启动：`POST /api/jobs/batch-export → 404`、`GET /api/export/users → 404`，而 `GET /api/schema/users` **仍含** `jobs-batch-export` 节点 —— 证明触发面在无模块 profile 下照常发布。

结论与证据链见 `D-001` §1。

### 2. 实施产物（C2）

| 文件 | 改动 |
|------|------|
| `apps/api/configs/config.yaml` | 内联模块列表补 `- admin.jobs`（附成因注释：该列表镜像 admin preset，缺它则 `POST /api/jobs/batch-export` 未挂载、按钮必 404） |
| `apps/web/src/components/jobs-batch-export.tsx` | 新增 `canBatchExport(context)`：同时要求 `jobs.write` 与 `data.export`；不可用时按钮 `disabled` + `data-jobs-batch-export-unavailable` + 说明文案（`-unavailable-note`）+ `title`；`submit()` 同步短路；上下文未知 → fail-open；提交命中 404 → 报 `schema.jobs.batchExport.unavailable` 而非回显裸 `NOT_FOUND` |
| `apps/web/src/i18n/messages/en-US.json`、`zh-CN.json` | 新增 `schema.jobs.batchExport.unavailable`（en: "Batch export is not available in this deployment (the admin.jobs module is not enabled)"；zh: "当前部署未启用批量导出（未装配 admin.jobs 模块）"） |
| `apps/api/internal/config/operator_config_coverage_test.go`（新增） | `TestOperatorConfigCoversAdminPreset`：解析 operator config 的 `app.modules.list`，断言其为 admin preset 的超集且含 `admin.jobs`；失败信息直接指名缺失模块与失效形态 |
| `apps/web/src/components/jobs-batch-export.test.tsx` | 新增 4 例可用性用例（见 §4） |

### 3. 修复后真实服务端到端（operator config）

以补齐后的 `configs/config.yaml` 启动，admin 完成强制改密后：

- 权限面：`jobs.write = True`、`data.export = True`（两者齐备 → 触发面在用户环境**保持可用**，不被新门禁误禁）。
- `POST /api/jobs/batch-export`（users，`["user-admin"]`）→ `202 queued` → `succeeded` → `downloadable=True` → 结果 `resource=users rowCount=1 fileName=users-selection.csv`。
- `POST /api/jobs/batch-export`（roles，`["role-admin"]`）→ `202 queued` → `succeeded`；结果 CSV 2155 字节（含 BOM 与冻结列序）。

### 4. 回归与变异验证（C3）

**变异验证（每条都真实跑过并观察到预期的红）**

| 变异 | 观察 |
|------|------|
| 从 `configs/config.yaml` 删除/注释 `- admin.jobs` | `TestOperatorConfigCoversAdminPreset` 失败：`missing 1 admin-preset module(s): admin.jobs` |
| 把组件判据改成**只查 `jobs.write`** | 用例「only jobs.write is granted」失败（判别性证明，A-002 F-001 后补） |
| 把组件判据改成**只查 `data.export`** | 用例「only data.export is granted」失败（对称方向） |
| 移除组件的可用性判断（`return true`） | 可用性用例「不可用时禁用并说明」失败 |
| 把 404 分支改回 `body?.message ?? …`（回显裸 `NOT_FOUND`） | 用例「404 报部署未启用而非裸 NOT_FOUND」失败 |
| e2e：把组件判据改成 `return true` 后跑 `batch-export-availability`（mvp） | `toBeDisabled` 失败（`Received: enabled`） |
| e2e：`APP_PROFILE=admin` 下跑 config 守卫 | 守卫 **SKIP**（形态不匹配）—— 登记为 recommended `F-003`，非静默通过 |

**全量回归（含 A-002 响应后的最终状态）**

| 面 | 结果 |
|----|------|
| Web 单测 | vitest **121 files / 1482 tests** 全绿（W33 为 1476；本波 +6：可用性用例 4 → 5 且其中 2 例为对称判别性用例，另有可用性 e2e 1 例不计入 vitest） |
| Web 类型/构建 | `npm run typecheck` exit 0；`npm run build` exit 0（构建后 `public/protocol` 已 `git checkout` 还原） |
| 浏览器 e2e（mvp） | **18 passed / 5 skipped / 0 failed**（含 always-on 可用性 spec） |
| 浏览器 e2e（admin） | **19 passed / 4 skipped / 0 failed** |
| 浏览器 e2e（custom，定向） | `batch-export-availability` **1 passed** |
| Go | `go build ./...` / `go vet ./...` exit 0；`go test ./internal/config/ ./modules/jobs/ ./internal/handler/ -count=1` 全绿（W34 前另跑过 `go test ./...` 全绿） |

> 数字更正说明：`E-001` 初稿（W33 相同数字 17/5/0 与 18/4/0）是**可用性 e2e 入仓前**的快照，被独立审计 `A-002` F-002 指出；入仓并补 harness 后复跑得到上表数字。

### 5. 实施中发现并修正的回归（真实、可核对）

第一次全量 e2e 出现 2 例失败，均由**本波第一版实现**（不可用时 `return null`）引起，非既有 flake：

- `list-visual-surface.spec.ts:82` → `page-action controls must share one height (got 0, 32, 32, 32)`：组件返回 `null` 后插槽宿主成为高度 0 的空 flex 子元素。
- `command-palette.spec.ts:37` → 同一页面结构变化的连带失败。

修正 = 改为 **disabled + 说明**（`D-001` §3），两个 spec 单独重跑 **6 passed**，随后双 profile 全量 e2e 全绿。

> **依据澄清（据独立审计 `A-002` 成果 4）**：硬约束是**空插槽宿主的高度契约**（宿主仍被注册时留下高度 0 的 flex 子元素，使 page-actions 行控件高度不一致）。W33 `D-001` §3 冻结的是**插槽布局** fail-open（schema 笔误不得让入口静默消失），把它直接当作「权限不足必须显示」是**过度延伸**；「可见 + disabled + 说明」的成立依据是「入口不撒谎 + 布局契约」，不是该冻结条文本身。

### 6. 边界与未做

- 未改服务端两道门禁语义、Job 六态合同、批量导出契约；未把 `admin.jobs` 加入 `mvp`/`demo` preset（Profile 内容决策，需用户 P-004，本波不动，仅让 UI 不再假装可用）。
- 未改 VP-038 `status`/判据与 workspace-038 台账正文；未触碰 `docs/schemas/**` 与 `apps/web/src/protocol/upstream/**`。
- Windows 端 `go test ./...` 的既有 PG drain flake（`TestShutdownDrainHarnessPostgres`）本波未触发；本次为 sqlite 路径全绿。

### 7. Git checkpoint

| commit | 内容 | 覆盖路径 |
|--------|------|----------|
| `e1013c8b` | 根因修复 + config 守卫 + 组件可用性用例 | `apps/api/configs/config.yaml`、`apps/api/internal/config/operator_config_coverage_test.go`、`apps/web/src/components/jobs-batch-export.tsx`、`apps/web/src/components/jobs-batch-export.test.tsx`、`apps/web/src/i18n/messages/{en-US,zh-CN}.json` |
| `A-002` 响应后（本波第二个 checkpoint） | e2e 可用性契约**入仓** + harness `custom` 列表对齐 + XOR 用例改为判别性 | `apps/web/e2e/batch-export-availability.spec.ts`（新增）、`apps/web/playwright.config.ts`、`apps/web/src/components/jobs-batch-export.test.tsx` |

### 8. 审计与响应（C4）

- **self `A-001`**：`pass`（0 required + 3 recommended）。
- **independent `A-002`**（grok build · grok-4.6 · high · `/audit`）：`conditional`（2 required + 3 recommended）。独立复验了根因（用 `784e9152` 的 YAML 直接核对「缺 `admin.jobs`」）、修复面（`git diff --stat` 仅 6 文件、未触禁止路径）、config 守卫变异、vitest 全量与可用性 e2e 双 profile；**推翻** self 的两处结论：① 「门禁缺一」单测不具判别性（标题与断言相反、未选行）；② 双 profile e2e 契约当时**未入仓**，不能作为覆盖缺口的关闭证据（并指出 harness 的 `customE2EModules` **当时就不含** `admin.jobs`，属现在时而非将来时）。
- **响应 `A-003`**：`F-001`/`F-002` 均 `fixed`（判别性用例双向变异 + spec 入仓 + harness 对齐 + 双 profile 复跑）；`F-003`/`F-004` recommended 保留并写明关闭要求；`F-005`（投影抢先写「当日关门」）接受并收回措辞。**开放 required = 0**。
- 被独立审计推翻的两处 self 结论已在 `A-003` 中 append-only 更正（不改写 `A-001` 原文）。
