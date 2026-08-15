---
id: GOAL-011-w10-account-page-conformance
doc: execution
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-002 · 数据权限页修复（I-002 裁决执行）

用户 2026-08-15 书面裁决：data-permission 是明确设立的交付物（workspace-011 GOAL-016），**保留并修复**（页面报错 + 菜单无合理图标）。

## 根因链（页面报错共四层 + 图标一层）

| # | 根因 | 证据 | 修复 |
|---|------|------|------|
| L1 | schema 顶层用 `view` 而 page.schema.json `required: ["meta","body"]` + additionalProperties:false（无 view 属性）→ D-VAL 失败 → loadPageDocument fail closed | data-permission.json 顶层键 [meta,actions,view]；page.schema.json L39-40 | 顶层 `view` → `body`（`type: "column"` → 渲染器支持的 `"section"` + id） |
| L2 | table 节点的 columns/dataSource/toolbar 放在节点**顶层**而非 `props` → node.schema.json additionalProperties:false 拒绝（×3）+ 渲染器读 `props.columns` 为空 | D-VAL 报 /body/children/0 must NOT have additional properties ×3 | 三字段移入 `props` |
| L3 | 策略行无 `id` 字段，默认 rowKey="id" → F-002 rowKey 校验 fail closed | 后端返回 {resource, ownerColumn, defaultScope, enabled}；临时渲染验证报 row 0 no valid "id" key | table 加 `"rowKey": "resource"` |
| L4 | registerPolicy 用 `{resource}` path 槽，而 formAction 只绑定 `{id}` 槽（request-construction.ts buildFormAction 无 path 绑定；render.tsx 仅 `{id}` 特判）→ 表单提交打到字面 `/policies/resource` | registerPolicy.url + requestMapping；临时验证 | 后端 `PATCH /api/data-permission/policies`（resource 移入 body，`Resource` 字段必填）；schema 同步去占位符 + requestMapping 删除 |
| L5 | 菜单 icon `shield` 不在前端 `iconRegistry`（App.tsx 硬编码映射）→ iconFor 返回 null → 无图标 | fragment.json icon: "shield"；App.tsx iconRegistry 无 shield | App.tsx 注册 `shield: Shield`（lucide-react） |
| L6 | GET policies 返回 `{"items": ...}` 无统一信封（缺 total/page/pageSize）→ 前端 parseResourceList fail closed（"expected a finite number"）；且 Policy 结构体无 json tag，序列化字段为大写 `Resource` 等，与 schema 列 camelCase 不匹配 → 列值缺失 | 用户实测报错 `parseResourceList.total: expected a finite number`；Policy struct（repository.go:48）无 json tag | 返回 `resourceList` 完整信封（Total/Page/PageSize）+ camelCase 行投影（resource/ownerColumn/defaultScope/enabled/updatedAt）；pageSize = max(1, len(items)) 使列表恒单页 |

## 改动文件

- `apps/api/internal/modules/datapermission/schema/data-permission.json`：view→body/section；table props 化 + rowKey；registerPolicy url/bodyMapping
- `apps/api/internal/handler/datapermission.go`：PATCH 路由无占位符 + body.Resource 必填校验；GET policies 完整信封 + camelCase 行投影
- `apps/api/internal/handler/datapermission_test.go`：PATCH 用例改 body 形态
- `apps/api/internal/kernel/profile.go`、`apps/api/internal/modules/datapermission/provider.go`、`provider_test.go`：Routes 契约 `PATCH /api/data-permission/policies`
- `apps/web/src/app/App.tsx`：iconRegistry + Shield
- `apps/web/src/renderer/representative-pages.test.tsx`：data-permission 纳入跨边界回归（D-VAL + 渲染 policies 表格）

## 验证

- Go 全量 `go test ./...` 全绿（含 handler、composition、datapermission、kernel；L6 修复后复跑无 FAIL）
- Web 全量 985/985（新增 1 个 representative-pages 用例）
- tsc -b --noEmit 0 错误
- 临时验证（已删）：D-VAL pass + RenderPage 渲染出 orders 行与 Register policy 按钮

## 状态

- I-002 **closed**（裁决：保留并修复；本记录为执行证据）。
- I-001 / I-003（参考样式对齐）仍 open，S1 未完成；本目标不变门。
