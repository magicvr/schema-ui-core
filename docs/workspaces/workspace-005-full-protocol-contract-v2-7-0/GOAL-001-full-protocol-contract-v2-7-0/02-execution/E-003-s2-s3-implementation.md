---
id: GOAL-001-full-protocol-contract-v2-7-0
doc: execution-entry
record_id: E-003
status: recorded
parent: null
created: 2026-08-08
updated: 2026-08-10
version: 0.1.1
---

# E-003 · S2/S3 纳入面实现（B1–B5）+ S4 验证登记

## 2026-08-08 · 实施批次 B1–B5

### 已发生事实

1. **B1（D-FORM/D-COMP 控件与展示）**：
   - `form-controls.ts`：新增 `inputNumber`（wire number）、`datePicker`（ISO string）、`dateRangePicker`（startField/endField 对）到白名单与门禁（DATE_RANGE_* 错误码、defaultValue 校验）；`form-controls.tsx` 新增 NumberField/DateField/DateRangeField；submit 投影展开 range 对。
   - `render.ts`/`render.tsx`：新增 `statCard`/`chart` 节点（registry props；dataSource 同 table F-001 单斜杠不变量；format 枚举 fail-closed；SVG line/bar/pie 无新依赖）；`SchemaCrudValue` 暴露 fetcher，`RenderPage` 接受 `dataFetcher`。
   - 范例页：`form-controls.json` 扩展（inputNumber/datePicker/dateRangePicker + advanced capability）；新增 `data-display.json`（statCard ×2 + chart，dataSource=/api/roles）；manifest/PageIDs/profile 注册。
2. **B2（D-EXPR 全引擎）**：
   - 新增 `reaction-expression.ts`：白名单语法 tokenizer/parser（`$deps`/`$self`/`$context.user|features`、`== != > >= < <= contains && || ! ( )`、字面量 string/number/bool/null）、严格类型语义（ADR-0016）、Unicode 码点字符串比较、contains 严格元素相等、表达式依赖提取、FORBIDDEN_VARIABLE/SYNTAX fail-closed。
   - 新增 `reaction-engine.ts`：快照多轮引擎（Snapshot→Evaluate→Commit→Next tick）、fulfill/otherwise 状态键、observers、baselines + resetMissingOtherwise、externalUpdates、深等提交检测、REACTION_LOOP_LIMIT（默认 10 轮）、MULTIPLE_VALUE_WRITES 警告；`runReactionEngineDetailed` 供 Renderer 取控件状态。
   - Renderer 集成：`resolveFullFormReactions`（上游 per-field reactions 形状）+ FormView 挂载基线 + 值提交合并 + 状态应用；`reactions.ts` $context 引擎保留为回退。
   - **16/16 reactions fixtures 全绿**（stage3）。
3. **B3（D-TABLE/D-ACT 批量）**：
   - `request-construction.ts`：`normalizeSelection`（D3：标量键、去重保序、count=keys.length）+ `buildBatchRequest`（D5：batchMapping path/query/body、$selection.keys 仅 body、$selection.count 标量、SELECTION_KEYS_BODY_ONLY、MISSING/EXTRA_PATH_BINDING、GET 拒绝、selectionAfterSuccessReload）；**11/11 batchRequest fixtures 全绿**（75/75 request-construction）。
   - `schema-table.tsx`：多选列 + 全选本页 + 查询变化清选（ADR-0022 D2）+ toolbar requiresSelection 禁用 + 批量触发分发；`render.tsx`：crud selection 状态、`invokeBatchAction`、`runBatchRequest`（权限门仅对已声明入口生效）、confirm 批量路径、reload 清空全部选中。
   - Go：通用 `POST {path}/batch-delete`（权限门、ids 标量校验、整批语义、`{deleted:n}` 响应）+ 测试；模块 descriptors/plan 同步。
   - 范例页：新增 `admin-list-batch.json`（selection.mode=multiple + 批量删除 toolbar）。
4. **B4（D-UPLOAD）**：
   - 新增 `upload-orchestration.ts`：约束前置（multiple/maxSize/accept token 匹配）、逐文件 multipart 请求构造、`runUploadBatch`（fixture 驱动）与 `uploadFilesWithFetch`（真实 fetch 传输）；**13/13 uploads fixtures 全绿**。
   - `form-controls`：`upload` 进白名单（wire string/string-array by multiple），门禁 UPLOAD_ACTION_REQUIRED/CONFLICT/CAPABILITY_REQUIRED/INVALID（registry oneOf + actions.upload）；`UploadField` 文件输入 + 实时上传 + 字段值回写；`SchemaCrudValue.uploadFiles`（actionRef → 页面 action / 直接 URL 双模式）。
   - Go：`POST /api/upload`（multipart、FILE_TOO_LARGE/UNSUPPORTED_FILE_TYPE/STORAGE_UNAVAILABLE、`{url,id,name,size}` 响应）+ `GET /api/files/{id}`（字节回读）+ 测试。
   - 范例页：新增 `form-with-upload.json`（actionRef=uploadAttachment → /api/upload）。
5. **B5（保真与 fail-closed）**：白名单外节点/控件拒绝测试、门禁错误渲染、批量空选禁用 + 不发请求（UI 集成）、上传字段缺 action 拒绝、循环引擎阻断提交、表达式非法 fail-closed（单元）。

### 测试证据

| 套件 | 结果 |
|------|------|
| `cd apps/web && npm run test`（vitest） | **25 文件 / 569 测试全绿** |
| `npx vitest run src/protocol/conformance/stage3-fixtures.test.ts` | **260 测试全绿**（其中 250 个为 15 套件 fixture case 执行） |
| `npx vitest run src/protocol/upstream-fixtures.test.ts` | 53 测试全绿（app-manifest 37 + app-navigation 16 case 执行） |
| `src/renderer/permissions-inheritance.test.ts` | 17 case 全绿 |
| `cd apps/api && go test ./...` | 全包 ok（含新增 batch-delete / upload 端点测试） |
| fixtures 执行面 | **16 行为套件 320/320 case 全部执行**（stage3 250 + upstream-fixtures 53 + permissions 17）：reactions 16/16、batchRequest 11/11、uploads 13/13 从「排除」转为「执行」 |

### 与覆盖表的关系

- S4 验证入口已在 `attachments/I-PROTO-FULL-001-coverage-v2-7-0.md` 登记（真实路径回填）；disposition 未变（12/12 include）。
- 未实现任何 exclude / include-partial；无用户 residual 引入。

### 证据路径

| 主张 | 路径 |
|------|------|
| 表达式引擎 | `apps/web/src/renderer/reaction-expression.ts` |
| 多轮引擎 | `apps/web/src/renderer/reaction-engine.ts` |
| 批量构造 | `apps/web/src/protocol/conformance/request-construction.ts`（buildBatchRequest） |
| 上传编排 | `apps/web/src/protocol/conformance/upload-orchestration.ts` |
| 渲染集成 | `apps/web/src/renderer/render.tsx`、`schema-table.tsx`、`form-controls.tsx` |
| Go 批量/上传端点 | `apps/api/internal/handler/resources.go`（batchDelete）、`upload.go` |
| 范例页 | `apps/api/internal/modules/schemarender/schema/{admin-list-batch,data-display,form-with-upload,form-controls}.json` |
| 覆盖表登记 | `attachments/I-PROTO-FULL-001-coverage-v2-7-0.md` §1/§3 验证入口列 |

## 2026-08-10 · 勘误注

本记录保留 2026-08-08 当时的执行声明。D-003 / E-007 复核后，`upstream-fixtures.test.ts` 的真实口径应为 app-manifest 35 executed + 2 excluded、app-navigation 16 executed；全体行为 fixture 为 **318 executed + 2 excluded**，不是 320 个 case 全部执行。当前权威见 `I-PROTO-FULL-001` v1.0.1。
