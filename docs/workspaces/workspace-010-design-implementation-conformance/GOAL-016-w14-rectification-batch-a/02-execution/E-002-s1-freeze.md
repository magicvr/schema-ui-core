---
id: GOAL-016-w14-rectification-batch-a
doc: execution
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-002 · S1 方案冻结

## 事实

- 2026-08-17：`01-decision/D-001-s1-freeze.md` 落盘，冻结 F-01～F-04 设计：
  - F-01：新增 `GET /api/scheduled-tasks/handlers`（tasks.read），表单 `optionsSource` 动态加载。
  - F-02：data-permission 页面新增 `data-permission-scopes` 自定义组件，复用 scopes GET/PATCH。
  - F-03：operations 列表增加 event/actorName/from/to 过滤；新增 `GET /api/operations/export` CSV。
  - F-04：notifications 新增 `title_key`/`body_key` 可空列（迁移 0037），新通知存 messageKey，前端回退 title/body。
- I-001/I-002 由 open 转为 closed（非阻塞信息项完成设计）。
