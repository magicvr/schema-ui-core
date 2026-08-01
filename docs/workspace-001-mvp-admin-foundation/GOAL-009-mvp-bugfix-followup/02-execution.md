---
id: GOAL-009-mvp-bugfix-followup
doc: execution
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.3.0
---

# 执行记录 · GOAL-009

## 时间线

### 2026-08-01 · 立项与审视附件落盘

- 用户确认：在工作区 `workspace-001-mvp-admin-foundation` 为 Root `GOAL-001-mvp-admin-foundation` 新建子目标，承接代码审视发现的 bug 修正；审视内容作独立附件。
- 创建五件套：`GOAL-009-mvp-bugfix-followup/`。
- 写入附件 [attachments/audit-code-review-bugs-2026-08-01.md](attachments/audit-code-review-bugs-2026-08-01.md)（F-009-001～007 全表）。
- `03-audit.md` 登记 A-001（`source: independent`，索引附件）。
- 同步 [goal-tree.md](../goal-tree.md)；Root `00-meta` 备注轻量提及本子目标。
- **尚未**修改 `apps/api` / `apps/web` 业务代码。

### 2026-08-01 · 实施 F-009-001～005（required 五项）

- **F-009-001**：`apps/api/internal/handler/records.go` `update()` 成功 PATCH 后写 `rec.UpdatedAt = time.Now().UTC()`；`records_test.go` 新增 `TestRecordsUpdateRefreshesUpdatedAt`（断言 updatedAt 严格后移）。
- **F-009-004**：`apps/api/internal/handler/account.go` 注释改为「nil 即 fail-closed 无会话」；`me()` 在 `sessionProvider == nil` 时返回 `401 UNAUTHENTICATED`（不再 panic）；`account_test.go` 新增 `TestAccountsMeNilProviderFailsClosed`。
- **F-009-002**：`apps/web/src/app/examples/list-edit-lifecycle-page.tsx` 删除硬编码 `{ roles: ["admin"] }` context，改为接收真实 `navigationContext` prop（App→PageSurface→registry 透传）；`PAGE_DOCUMENT` 表节点挂 `permissionCascade: { keys: ["edit","delete"] }` + `permissions: { edit/delete: '$context.user.roles contains "admin"' }`，使 Edit/Delete 门禁可失败；`registry.tsx` 组件签名带 `context?: NavigationContext`；`row-action.ts` 的 `context` 类型改 `NavigationContext`；新增 `list-edit-lifecycle.test.tsx`（admin 会话按钮启用 / viewer 会话按钮 disabled 拒绝路径）。
- **F-009-003**：`apps/web/src/main.tsx` 保留 `loadAccountContext` 的 `error`、`console.error` 日志、并以 `accountError` 传入 App；`apps/web/src/app/App.tsx` 新增非阻断 alert 横幅（`role="alert"`）；`App.integration.test.tsx` 新增横幅断言（含健康路径无横幅）。
- **F-009-005**：重写 `apps/api/README.md` 与 `apps/web/README.md`（当前端点、session 模型、鉴权边界、测试命令）。
- **回归证据**：`go test ./...` 全绿；web vitest **398/398**（16 文件）；`tsc -b --noEmit` 通过；`npm run build` 通过。
- 遗留：F-009-006/007（recommended）按 I-009-001/002 待用户裁决；阶段/关门审计未做。

### 2026-08-01 · 实施 F-009-006/007（用户裁决纳入）

- 用户裁决（I-009-001/002）：「都纳入实施」。
- **F-009-006**：`apps/api/internal/handler/records.go` 新增 `writeGate()`——PATCH/DELETE 写路由 fail-closed 鉴权：无会话 → `401 UNAUTHENTICATED`；非 admin 角色 → `403 FORBIDDEN`（`account.Allow` `$context.user.roles contains "admin"`）；`recordHandler` 注入 `sessionProvider`。新增 `TestRecordsWriteRequiresSession`、`TestRecordsWriteDeniedWithoutAdminRole`。
- **F-009-007**：records.go `update()` 挂 `http.MaxBytesReader`（4 KiB）；`list()` `pageSize ≤ 100` 超限 `400 INVALID_PAGE_SIZE`。新增 `TestRecordsUpdateBodyTooLarge`、`TestRecordsListPageSizeCap`。
- 更新两 README 鉴权边界（写路由已鉴权；上限说明）。
- 回归证据：`go test ./...` 全绿（含 4 项新测试）；web vitest **398/398** 不变。

### 2026-08-01 · 关门（A-002 pass + 用户授权）

- 独立交叉审计 **A-002**（`/audit`，`source: independent`，close-out）`verdict: pass`：复核 A-001 七条 `fixed` 证据 + 成功标准 5/5 + I-009 resolved；复跑 `go test ./...` + 聚焦 vitest（7/7 pass）；open required = 0；给出 2 条 recommended（F-A002-001 鉴权措辞、F-A002-002 台账卫生）。
- 用户 P-004 裁决：不补 self 关门审，接受 A-002；关闭前修复两条 recommended。
- 修复 F-A002-001：`apps/api/README.md` / `apps/web/README.md` 鉴权边界各加「gate 绑定进程内会话提供者、非请求头身份、默认 admin 会话下无凭证仍可写」说明。
- 修复 F-A002-002：01-decision 文首 I-009 对齐 resolved；A-001「编排提示」加历史基线说明；附件顶加「实施后状态以 03-audit 为准」注。
- 用户授权关门 → GOAL-009 `status: done`；同步 goal-tree；不改 Root/VP。

## 待办

- 无（目标已 `done`）。浏览器手测保持 optional，未纳入必做。

## 进度评估

**5/5** required 检查点完成；A-001 七条 findings + A-002 两条 recommended 全部 `fixed`；A-002 关门复审 `pass`；用户书面授权 → **`status: done`**（2026-08-01）。
