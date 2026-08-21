---
id: GOAL-016-w14-rectification-batch-a
doc: decision
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · GOAL-016 S1 方案冻结（F-01～F-04）

## 背景

GOAL-016 承接 GOAL-015 D-003 用户书面裁决：批 A 实施 F-01～F-04。本文件冻结 S1 设计决策，作为 S2 实施输入。

## 决策

### F-01 · 定时任务可指定 handler

- **新增端点**：`GET /api/scheduled-tasks/handlers`（`tasks.read` 权限），返回统一列表信封 `{items:[{key,label}],total,page,pageSize}`。
- **后端**：由 `Scheduler.HandlerKeys()` 提供可用 handler；v1 只有 `system.noop`，label 暂与 key 相同（前端可用 i18n key 覆盖展示）。
- **Schema**：create/edit 表单新增 `handler` select 字段，`optionsSource` 指向新端点（`valueField=key`, `labelField=key`），create 默认 `system.noop`。
- **写入契约**：沿用现有 `validateHandler`，未知 handler 拒绝（`INVALID_HANDLER`）。

### F-02 · 数据权限范围设置入口

- **方案**：在 `data-permission` 页面增加自定义组件 `data-permission-scopes`（schema `type:"custom"`）。
- **交互**：选择用户（`/api/users` 动态选项）→ 加载该用户现有 assignments（`GET /api/data-permission/scopes?userId=`）→ 对每条已注册 policy 提供 `all`/`self` 选择（默认取 policy.defaultScope）→ 保存时 `PATCH /api/data-permission/scopes`。
- **API 契约**：不新增后端端点；复用现有 scopes GET/PATCH。组件仅做 UI 接线。

### F-03 · 审计日志结构化过滤与导出

- **过滤**：`/api/operations` 增加 `ExtraQuery`：`event`、`actorName`、`from`、`to`。`from`/`to` 接受 `YYYY-MM-DD`（按天边界）或 RFC3339；非法值返回 400 `INVALID_DATE_FILTER`（通过 List DomainError 通道）。
- **仓库**：`OperationFilter` 增加 `Event`、`ActorName`、`From`、`To`；`operationsWhere` 按需追加精确/前缀匹配与 created_at 范围。
- **导出**：新增 `GET /api/operations/export`（`operations.read`），按同一过滤条件返回 CSV（UTF-8 BOM、RFC 4180 转义、attachment、公式注入防护），行数上限 `maxExportRows`（10000）。
- **Schema/UI**：activity 搜索表单增加 `event`/`actorName` 输入与 `from`/`to` 日期选择；新增自定义组件 `activity-export` 按钮，读取当前表格查询并触发下载。

### F-04 · 通知本地化 messageKey

- **存储**：在 `notifications` 表新增可空 `title_key`、`body_key` 列（迁移 0037）。新通知写入 messageKey，`title`/`body` 保留为空字符串（NOT NULL 兼容）；旧数据保持原英文 title/body 作为回退。
- **API**：通知列表/详情返回 `titleKey`/`bodyKey`（非空时）以及原有 `title`/`body` 回退字段。
- **前端**：notification-center 渲染时优先 `titleKey`/`bodyKey` 经 i18n 翻译，缺失时回退 `title`/`body`。
- **i18n**：新增 `notification.account.locked.title/body`、`notification.account.disabled.*`、`notification.account.unlocked.*`、`notification.account.passwordChanged.*` 中英文键。

## 信息项更新

| ID | 状态 | 说明 |
|----|------|------|
| I-001 | **closed** | F-01 端点路径与权限已冻结：`GET /api/scheduled-tasks/handlers` + `tasks.read` |
| I-002 | **closed** | F-04 旧文案迁移已冻结：保留旧 title/body 回退，新增可空 key 列，不重发/迁移旧记录 |
