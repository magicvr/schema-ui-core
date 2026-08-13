---
id: E-003
goal: GOAL-004-r2-f02-data-import-export
title: S2 实现 + S3 验证（export/import 全量落地 + 回归 + 本地冒烟）
date: 2026-08-14
status: recorded
parent: GOAL-004-r2-f02-data-import-export
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-003 · S2 实现 + S3 验证

## S2 事实（checkpoint `39a1671`，24 files / +1272）

| 项 | 落地 | 证据 |
|----|------|------|
| 导出 | `GET /api/export/{resource}`（users/roles · RFC 4180 · UTF-8 BOM · cap 10000 · q/sort/order 过滤 · Content-Disposition） | handler/export.go |
| 导入 | `POST /api/import/users` `{fileId}`（owner 校验 · 2 MiB 上限 · 逐行校验+不回滚 · `{applied,failed,total,errors}` 报告） | handler/import.go |
| 权限键 | `data.export`（PolicyAdminEditor）/ `data.import`（PolicyAdmin）；匿名 401 / 无键 403 / 模块关闭 404 | modules/datatransfer/provider.go |
| 迁移 0015 | operation_log CHECK 扩展（data.export/data.import） | operationlog/migration |
| 前端 | 协议 `CustomAction` 扩展点白名单（`export.users` / `export.roles` → authed blob 下载；未知 handler fail-closed）；users 页 Export/Import（modal 含 upload 控件 + submitImport）；roles 页 Export | render.tsx、users.json、roles.json |
| 装配 | admin profile += `admin.data-transfer`（内容扩展；无 fragment → adminFunctionalOrder 不变） | kernel/profile.go、composition.go |
| i18n | en/zh 键（toolbar/import modal） | i18n/messages |

**方案偏差**：D-002 `5 原定「download behavior 本地扩展」改为 **CustomAction 白名单**——上游 action.schema.json 严格校验 `onSuccess.behavior` 枚举（toast/navigate/reload/closeModal）且 `additionalProperties:false`，download 行为无法过结构校验；`CustomAction` 是协议自带的扩展点（handler 白名单语义），比行为扩展更契合必办-1「本地契约 + fail-open」留痕。文件名由 handler 名推导（`users.csv` / `roles.csv`）。

## S3 验证事实（2026-08-14）

| 项 | 结果 |
|----|------|
| `go test ./... -count=1` | ✅ 全绿（导出形状/BOM/转义/权限/未知资源 404；导入部分成功 2/1/3、错误报告、owner 403、越权 403、fileId 404、审计落盘；错误契约新增 3 码；边界测试 users-manager 角色补 `data.export`——editor 角色新增该键后委托边界要求持有方一致） |
| `npm test` | ✅ 893/893（含 custom action 3 用例；schema 结构校验通过——users/roles 页含 export/import 动作后仍全绿） |
| `npx tsc -b` + `npm run build` | ✅ |
| 本地冒烟（admin profile） | ✅ 导出（BOM+转义+CT+附件头）→ viewer 403 → 上传 CSV → 导入（applied=2 failed=1，row 4 password required）→ 查询确认 2 用户 → oplog `data.export/data.import` 落盘 |

## 门禁结论

S2/S3 完成。进入 S4（go 判定 + self 审计）。
