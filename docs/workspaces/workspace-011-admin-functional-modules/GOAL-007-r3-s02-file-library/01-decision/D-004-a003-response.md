---
id: D-004
goal: GOAL-007-r3-s02-file-library
title: A-003 响应：4 条 recommended 全 fixed
date: 2026-08-14
status: accepted
parent: GOAL-007-r3-s02-file-library
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-004 · A-003 响应（关门）

## 结论

A-003（grok-build independent，security/data）verdict **pass**、0 required；4 条 recommended **全部 fixed**（无 residual / overruled）：

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| F-001 客户端 blob 文件名消毒弱于服务端 | recommended | **fixed** | render.tsx sanitizeClientFilename：allowlist [A-Za-z0-9._-] + 去首尾分隔符 + 100 上限 + 空回退；单测期望同步（download-behavior.test.tsx） |
| F-002 DELETE 对象缺失时残留孤儿 meta | recommended | **fixed** | filelibrary.go 删除逻辑重构：对象与 meta 均 best-effort 移除；两者皆缺才 404；移除任意一项即审计 files.delete |
| F-003 审计写错误静默丢弃 | recommended | **fixed** | filelibrary.go recordFileEvent：slog.Error 留痕（与 export/users 对等），仍 best-effort 不影响 HTTP |
| F-004 工具栏上传受 files.read（表 edit）门控而非 files.write | recommended | **fixed** | file-library.json 工具栏 upload 显式 permissions.edit = files.write；modal 表单本就 files.write，两层一致 |

## 验证

- go test ./internal/handler/ -run TestFileLibrary 全绿；vitest（download-behavior + schema-keys 结构）8/8；页面 schema 经 AJV 校验通过。
- 全量回归（go ./... / vitest / e2e 双 profile）在 E-004 记录。
