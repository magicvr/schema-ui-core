---
id: E-001
doc: execution
title: S6 · 设置页表单/详情页实现（recordSource + schema + 测试）
status: recorded
parent: GOAL-007-s6-settings-form-page
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-001 · S6 设置页表单/详情页实现（2026-08-09）

## 事实

- **Renderer —— 实现既有协议表面 `form.props.recordSource`（ADR-0021）**：
  - `render.ts`：`RenderFormNode.props` 增加 `title`/`titleKey`/`recordSource`；`parseRenderNode` form 分支放行三者；`RenderActionButtonNode.props` 放行 `labelKey`/`key`/`confirm`/`confirmKey`/`permissionIntent`；新增纯函数 `resolveResponsePath`（点号路径）。
  - `form-controls.ts`：新增 `FORM_RECORD_LOAD_CAPABILITY = "form.record.load"`。
  - `render.tsx`：`FormView` 拆为外层（recordSource GET 预填：capability 门禁 / search 拒绝 / `constructRequest` 构 URL / loading·error·ready 状态机 / `reloadToken` 重拉）+ 内层 `FormInner`（初值 `modalRow → prefill → 空`，`${formId}:submit` 只读门禁禁用字段与提交、form `title`/`titleKey` 渲染 `<h2>` 标题）；`invokeAction` 支持 `actionId` 回退；`ActionButtonView` 按 `props.key` 权限目标禁用。
- **Settings schema（`apps/api/internal/modules/settings/schema/settings.json`）**：
  - 4 个 recordSource 预填内联表单（General/Branding/Localization/Appearance，`titleKey` 复用 `schema.settings.toolbar.*`，`submitAction` 复用现有 update* PATCH，`responseMapping` 全 identity）；Restore defaults 改 `actionButton`（`actionId:"resetSettings"`、`permissionIntent:"edit"`、confirm）；body section 挂 `permissionCascade.keys:["edit"]` + `permissions.edit`。
  - 删除 `settings-table`/`settings-detail`/4 个 `open*` modal action/行内 Edit；保留 5 个 request action；meta 增 `form.record.load`、移除 `table.sort`/`actions.row.request`；全文本复用现有 catalog 键（无新增键）。
- **测试**：vitest **727/727**（40 文件；新增 renderer 用例：recordSource 预填/reload 重预填/capability fail-closed/search 拒绝/form title/actionButton actionId+confirm+权限禁用、`resolveResponsePath` 3 例；`startup-config` C4/C5 4 用例改写为内联表单形态）；`schema-keys.structural`（12 页 *Key 完备）通过；`npm run build`（tsc + vite）exit 0；`go test ./apps/api/...` exit 0。
- **e2e 环境降级**：`localization.spec.ts` M3 已改写为新流程（内联表单填写→保存→投影 + 四类 heading + 恢复默认 button）；本机 8080 落入 Windows 端口排除区间（8011–8110，`netsh` 快照见 `attachments/s6-e2e-env-block.log`）致 Go API 无法绑定，playwright 无法启动。按 S5 D-001 §3 诚实降级：捕获失败输出 + 单元覆盖 M3 逻辑为验收线（S5 同机曾跑通，排除区间后移所致，非代码问题）。

## 产物

| 路径 | 说明 |
|------|------|
| `apps/web/src/renderer/{render.ts, render.tsx, form-controls.ts}` | recordSource 预填 / title / actionButton 接线 / 只读门禁 / capability 常量 |
| `apps/api/internal/modules/settings/schema/settings.json` | 四类内联表单 + reset actionButton（C2） |
| `apps/web/src/{app/startup-config.test.tsx, app/representative-pages.integration.test.tsx, i18n/s5-denominator-render.test.tsx, renderer/render.test.ts, renderer/render.test.tsx}` | 测试更新与新增（C3） |
| `apps/web/e2e/localization.spec.ts` | M3 新流程（环境降级，C3） |
| `attachments/s6-e2e-env-block.log` | e2e 端口排除证据 |

## 里程碑 checkpoint

- commit：`ebd0288`（2026-08-09，S6 实现；owned paths = 上表全部路径，显式 `git add` 无 `-A`）。
- 治理前置 checkpoint：`a49f04c`（GOAL-007 五件套 + Root 回退，见 GOAL-001 `D-003`）。
