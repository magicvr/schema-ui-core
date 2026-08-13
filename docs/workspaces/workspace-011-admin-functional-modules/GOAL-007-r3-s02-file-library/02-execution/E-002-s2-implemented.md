---
id: E-002
goal: GOAL-007-r3-s02-file-library
date: 2026-08-14
status: recorded
parent: GOAL-007-r3-s02-file-library
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-002 · S2 实现完成

## 事实

- 2026-08-14：S2 实现完成，覆盖 D-002 §2–§5：
  - **API 模块** apps/api/internal/modules/filelibrary/：provider.go（admin.file-library，五面贡献：HTTP/Schema/Authorization/Navigation/Manifest）、schema/file-library.json（列表 + 行下载 CustomAction library.download + 行删除 + 上传 modal）、manifest/fragment.json（file-library 页 + menu_files 导航）。
  - **handler** apps/api/internal/handler/filelibrary.go：GET /api/library/files（列表，磁盘扫描单一事实源）、GET /{id}（detail）、GET /{id}/download（files.read + nosniff/attachment/CSP/sanitize 文件名）、DELETE /{id}（files.delete 硬删）、POST /upload（files.write 确认端点：id 形状/存在性/所有者校验 + 审计挂点）；审计事件 files.upload/files.download/files.delete。
  - **migration 0018**：operation_log CHECK 扩展三个 file 事件（rebuild 模式同 0014/0015）；checksum 3351b6e6…（Go 权威计算，与台账 0014/0015 复算校验一致）。
  - **装配**：kernel/profile.go（admin 默认集 + BuiltinModules）、composition.go（plan.HasModule 装配）、testsupport/store.go 镜像。
  - **web**：render.tsx CUSTOM_HANDLER_URLS + runCustomAction 行上下文 {id} 解析 + 行 name 作下载文件名（客户端 blob，无 header 注入面）；i18n zh/en 键（manifest.title/nav.fileLibrary + schema.fileLibrary.* + error.invalidUploadBody/invalidFileId）。
  - **测试**：handler filelibrary_test.go（lifecycle/权限门禁/ack 校验）、provider_test.go（注册面 + 端到端）；错误码目录 + 钉住集 + 迁移计数（17→18）断言全量更新。
- 页面 schema 经 AJV 对照 docs/schemas/page.schema.json 校验通过（meta 不含 titleKey——titleKey 只在 fragment）。
