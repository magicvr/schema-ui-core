---
id: GOAL-016-w14-rectification-batch-a
doc: execution
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-003 · S2/S3 实施与回归

## 事实

- **2026-08-17**：F-01 实施——`GET /api/scheduled-tasks/handlers` 返回 `{items:[{key,label}]}`；scheduled-tasks schema create/edit 增加 `handler` select + `optionsSource`；provider/kernel profile 路由清单同步。
- **2026-08-17**：F-02 实施——新增 `apps/web/src/components/data-permission-scopes.tsx` 自定义组件并注册到 `main.tsx`；data-permission schema 增加 custom 节点；复用 scopes GET/PATCH。
- **2026-08-17**：F-03 实施——operationlog `OperationFilter` 增加 Event/ActorName/From/To；operations 列表 ExtraQuery `event/actorName/from/to`；`GET /api/operations/export` CSV；activity schema 增加过滤字段 + `activity-export` custom 组件。
- **2026-08-17**：F-04 实施——notifications 迁移 0037 增加 `title_key`/`body_key`；repository/handler 写入与返回 messageKey；notification-center/bell 优先 i18n key 渲染；i18n en/zh 新增通知文案键。
- **2026-08-17**：S3 回归——Go 全量 `go test ./...` 通过；Web 全量 `npm test` 1041/1041 通过；`npx tsc -b` 通过；`npm run build` 通过（含 conformance claim 重新生成）。
