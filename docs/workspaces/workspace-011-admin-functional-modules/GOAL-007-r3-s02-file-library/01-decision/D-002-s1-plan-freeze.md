---
id: D-002
goal: GOAL-007-r3-s02-file-library
title: 方案冻结：文件/附件库设计（S1）
date: 2026-08-14
status: accepted
parent: GOAL-007-r3-s02-file-library
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · 方案冻结（S-02 文件/附件库）

## 1. 协议对照（I-001 闭合）

- 协议面（v2.8.0，I-PROTO-FULL-001）：上传为 `type: "upload"` action（action.schema.json，POST /api/upload 契约，D-UPLOAD include）；页面呈现无新语义契约——**呈现自由 + fail-open** 处置（沿用 F-01/F-02 先例，本文件留痕）。
- 行操作下载 = 协议 `CustomAction` 扩展点（action.schema.json，sanctioned）+ 本仓 F-02 白名单机制（render.tsx CUSTOM_HANDLER_URLS）；删除 = 标准 request action + 资源 DELETE。
- 上传/下载/删除**不属于** schema-ui-docs 既有动作契约 → 本地契约 + 白名单 + 文档化，无协议 pin 变更。

## 2. 存储与数据模型

- 复用中心上传存储：`<DB dir>/uploads`（composition.go:183），对象文件 + `<id>.meta.json`（name/type/owner）。
- 列表 = 扫描目录（同 quotaReached 模式）：id / name / type / owner（meta），size（stat 对象文件），created（mtime 近似——meta 无时间戳，v1 文档化）。
- 排序/筛选：name/type/owner（meta 字段）、size/created（stat/mtime）；服务端分页 + total，遵循资源列表契约（items/total）。
- 校验/加固：id 必须匹配 `uploadFileIDPattern`（32 hex）；文件名 sanitize（下载 Content-Disposition 仅安全字符或固定值）。

## 3. 权限与导航

| 键 | Policy | 用途 |
|----|--------|------|
| `files.read` | PolicyAdmin | 库列表 + 任意文件下载（管理面） |
| `files.write` | PolicyAdmin（已有，core.server-registration） | 上传（复用中心端点） |
| `files.delete` | PolicyAdmin | 删除（破坏性操作仅 admin） |

- 导航 `menu_files`（visibility PolicyAdmin，permission files.read）→ features.menu_files 投影。
- 既有 `GET /api/files/{id}` owner-only 契约**不变**（表单控件字段用）；库下载走独立端点。

## 4. 端点（路由贡献）

| 端点 | 门禁 | 说明 |
|------|------|------|
| GET /api/library/files | files.read | 列表（资源工厂 Listable；q/sort/order/page/pageSize/total） |
| GET /api/library/files/{id}/download | files.read | 下载（nosniff + attachment + CSP sandbox，同中心下载加固） |
| DELETE /api/library/files/{id} | files.delete | 硬删对象 + meta；不存在 → 404 FILE_NOT_FOUND；审计 files.delete |
| POST /api/library/files/upload | files.write | 上传**确认**端点：JSON `{file: id|url}`，校验 id 存在 + 上传者一致；幂等；审计 files.upload |

- 说明：真正存储由中心 `POST /api/upload` 完成（复用全部加固：8 MiB 上限、类型嗅探、活跃内容拦截、每用户配额）；本确认端点为 schema 表单提交流程（modal 上传字段 → submit）的落点 + 审计挂点，不做二次存储。
- 下载审计 files.download 在下载端点记录（best-effort）。

## 5. 页面与 Schema

- 页面 `file-library`（manifest fragment + page schema，标准 schema 驱动，不改 Renderer 主路径）：
  - 列表列：name / type / size / owner / created
  - 行操作：**下载** = CustomAction `library.download`；**删除** = request DELETE（confirm）
  - 工具栏：**上传文件** = modal 表单（upload 字段 `actionRef uploadFile → /api/upload` + submit `POST /api/library/files/upload` → reload）
  - 权限级联：行/工具栏操作按 files.read / files.write / files.delete 表达式
- i18n：zh/en 增 `manifest.title.fileLibrary`、`manifest.nav.fileLibrary`、`schema.fileLibrary.*` 键。
- 页面 titleKey/关键 label 走既有 i18n 机制（F-01 先例）。

## 6. 引用与清理策略（I-002 闭合）

- v1 **不做引用注册表**：删除即硬删；表单字段值 / 进行中导入对已删文件的引用在渲染/读取时 fail-open（404 呈现为缺失，导入报 FILE_NOT_FOUND）。
- 清理 = 已有配额（UPLOAD_MAX_FILES_PER_USER / UPLOAD_MAX_BYTES_PER_USER）+ 库内手动删除；孤儿策略 v1 不引入（文档化残余）。
- 风险记录：删除可能破坏客户端引用 —— v1 接受（admin 管理面 + 审计留痕），S-12 回收站（R3 后续）再评估软删。

## 7. 测试与验证路径

- API：模块 provider 测试（Descriptor/权限贡献）、列表契约（分页/排序/筛选）、下载门禁（files.read 403/200、id 格式、不存在 404）、删除门禁 + 审计事件、上传确认（幂等/所有者校验/不存在 404）
- 组合根：mvp 不含文件库页面、admin 含页面/导航/权限（composition_test 计数更新）
- 全量回归：go test ./...、web vitest（fixture 同步：admin fixture + STATIC_MANIFEST_SHA256 重钉）、e2e（admin 冒烟手测按需）
- 冒烟：SM-007 admin 页面集加 file-library

## 8. 未选方案

- DB 镜像表（含 ref_count 预留列）：双写一致性与死代码风险，v1 否决（D-001）。
- 库级下载权限放宽到 owner-only：与「统一管理面」语义冲突，否决。
- 扩展中心 /api/upload 加审计：中心端点属 core 契约路径，避免横切改动，审计挂在模块确认端点。
