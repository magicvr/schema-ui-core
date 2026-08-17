---
id: GOAL-016-w14-rectification-batch-a
doc: audit
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# A-001 · W14 整改批 A 独立审计（F-01～F-04，重点 F-04）

- source: independent
- auditor: independent subagent
- date: 2026-08-17
- scope: GOAL-016 S1-S3 (F-01..F-04), F-04 focus
- verdict: conditional
- 成果（有证据）:
  - F-01 定时任务 handler 目录端点 `GET /api/scheduled-tasks/handlers`（`tasks.read`，返回 `{items:[{key,label}],total,page,pageSize}`）已实现：`apps/api/internal/handler/scheduledtasks.go` L230-249；create/edit schema `handler` select 接 `optionsSource` 且 create 默认 `system.noop`：`apps/api/internal/modules/scheduledtasks/schema/scheduled-tasks.json` L106-117 / L165-175；写入沿用 `validateHandler`（未知 handler 拒绝）。
  - F-02 数据权限范围设置入口 `data-permission-scopes` 已实现：选用户 → `GET /api/data-permission/scopes?userId=` → 逐 policy 提供 all/self（默认取 `defaultScope`）→ `PATCH /api/data-permission/scopes`：`apps/web/src/components/data-permission-scopes.tsx`。
  - F-03 审计日志结构化过滤与导出已实现：`operations.go` `ExtraQuery`（event/actorName/from/to）+ `parseOperationTime`（YYYY-MM-DD 按天边界与 RFC3339，非法返回 `INVALID_DATE_FILTER`）；`operationlog/repository.go` `OperationFilter` 与 `operationsWhere` 精确/日期范围匹配；`operations_export.go` 新增 `GET /api/operations/export`（`operations.read`，UTF-8 BOM、RFC 4180、公式注入防护、`maxExportRows` 10000）。
  - F-04 通知本地化 messageKey 已实现且自洽：迁移新增可空 `title_key`/`body_key`（`migration.go` version 37）；repository 读写字端、handler 列表返回 `titleKey`/`bodyKey`（非空时）、`NotifyAccountEvent` 新通知写 key 而 title/body 置空串；前端优先 `t(titleKey/bodyKey)` 回退 `title/body`：`notifications.go`、`notifications_repository.go`、`notification-center.tsx`；en/zh 均存在 `notification.account.locked/disabled/unlocked/passwordChanged.*` 键。

- Findings:

  - **F-001（required，阻断关门）** — S1 冻结方案与 as-built 迁移号不一致。`01-decision/D-001-s1-freeze.md` L41 写明「迁移 0018」，但实现 `migration.go` L60-65 使用 Version 37（注释「0037 · W14 F-04」），`02-execution.md` L27 亦记录 0037。四件套内部事实矛盾，须在关门审计前修正 D-001（改为 0037）或经用户书面裁决接受差异并留痕。证据：`01-decision/D-001-s1-freeze.md` vs `apps/api/internal/modules/notifications/migration/migration.go` vs `02-execution.md`。

  - **F-002（recommended，不阻断）** — messageKey 命名文本与实现用名不一致。D-001 L44 写 `notification.account.password-changed.*`，实现统一用 `notification.account.passwordChanged.*`（`notifications.go` L288、`en-US.json` L182-183、`zh-CN.json` L182-183）。运行时前后端一致、功能正确，仅建议把决策文本对齐到实际键名（或反向重命名），避免后继批误读。

  - **F-003（recommended，不阻断）** — S3 回归为断言事实，当前审阅给的文件集无法独立复核测试证据。`02-execution.md` 声明 Go 全量 / Web 全量 / tsc / build 通过，但 02-execution/E-003 不在本次允许读取范围；建议编排器在 S4 关门复核时附上 E-003 测试产物/日志作为证据留痕。

- 必改项汇总
  - F-001：修正 D-001 冻结文档中 notification messageKey 迁移号（0018 → 0037）或取得用户书面裁决记录差异，二者须留痕，闭合后方可关门。

- 结论: F-01～F-04 的功能/契约与前端接线在可读文件范围内基本与 D-001 冻结方案一致，F-04（重点）消息键迁移 0037、repository/handler/前端/i18n 四层自洽，旧记录回退语义正确；存在一处 required 文档一致性偏差（迁移号 0018 与 0037 冲突），故 verdict 为 conditional，闭门前须闭合 F-001，其余为 recommended 不阻断。

- 声明: This opinion does not modify status/progress; response handled by /govern.