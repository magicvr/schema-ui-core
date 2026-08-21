---
id: A-003
goal: GOAL-007-r3-s02-file-library
source: independent
date: 2026-08-14
scope: S-02 实现安全/数据门禁（admin.file-library vs D-002 冻结方案）
verdict: pass
auditor: grok-build
audit_type: execution-facts
status: recorded
parent: GOAL-007-r3-s02-file-library
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-003 · independent 安全/数据审计（S-02 实现）

## 范围与区间

- **auditor**：grok-build（independent cross-audit）
- **type**：execution-facts / security-data gate
- **covered**：
  - `apps/api/internal/handler/filelibrary.go`（list/detail/download/delete/upload-ack）
  - `apps/api/internal/modules/filelibrary/`（provider、schema、manifest）
  - `apps/api/internal/modules/operationlog/migration/migration.go`（0018）
  - `apps/api/internal/handler/upload.go`（中心上传基线）
  - `apps/web/src/renderer/render.tsx`（CUSTOM_HANDLER_URLS / runCustomAction / blob 文件名）
  - `apps/api/internal/kernel/profile.go`（admin 默认集）
  - 计划契约 `01-decision/D-002-s1-plan-freeze.md`（及 D-001/D-003 中 Profile 声明）
- **excluded**：端到端浏览器手测、生产部署配置、S-12 回收站软删设计、其他工作区上下文

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 库表面门禁 fail-closed（anon 401 / 无 key 403） | `requirePermission`（resources.go:231–241）；download/delete/upload-ack 各调 files.read/delete/write（filelibrary.go:204/232/269）；Resource list/detail 用 files.read（filelibrary.go:191）；`TestFileLibraryPermissionGates` |
| 中心 `GET /api/files/{id}` owner-only 未改 | upload.go:336–341 仍校验 `meta.owner == user.ID`；库下载走独立 `/api/library/files/{id}/download` |
| 上传确认校验所有权 | filelibrary.go:298–303；`TestFileLibraryUploadAckValidation` foreign → 403 |
| id 形状白名单防路径穿越 | `uploadFileIDPattern` `^[0-9a-f]{32}$`（upload.go:131）；Get/scan/delete/ack/load 均校验；`filepath.Join(dir, id)` 仅在合法 id 后 |
| 下载头加固 | filelibrary.go:251–257：nosniff + attachment + CSP sandbox + `sanitizeFileHeaderName` |
| 配额不可被库表面绕过 | 库无二次存储；真实写入仅 `POST /api/upload`（quotaReached + files.write，upload.go:299–304） |
| 审计事件与 0018 CHECK 一致 | 发射：`files.upload`/`files.download`/`files.delete`（repository.go:35–37 + filelibrary.go:222/246/304）；CHECK 列表含且仅扩这三项（migration.go:71） |
| 0018 重建保行 + 版本连续 | rebuildOperationLog rename/copy/drop（migration.go:174–193）；台账版本 1..18 无空洞（authsession/corepersistence/operationlog/settings/account/notifications） |
| Profile 内容扩展未改装配语义 | profileDefaults admin 仅追加模块 id（profile.go:65–67）；ResolveProfile 逻辑未变；D-001/D-003 声明一致 |
| Renderer `{id}` / 未知 handler fail-closed | render.tsx:287–298（CUSTOM_HANDLER_NOT_FOUND / CUSTOM_HANDLER_MISSING_ROW_ID）；`encodeURIComponent(rowId)` |

## 对照成功标准（D-002 安全/数据相关）

| 标准 | 结论 |
|------|------|
| files.read / files.write / files.delete 门禁 | **满足** |
| 中心 owner-only 契约保留 | **满足** |
| upload-ack 所有者校验 | **满足** |
| id 校验 + 安全 Join | **满足** |
| 下载 XSS/头注入加固 | **满足**（服务端）；客户端 blob 文件名见 F-001 |
| 配额不可绕过 | **满足** |
| 0018 CHECK ↔ 事件 | **满足** |
| 迁移安全 | **满足** |
| Profile 仅内容扩展 | **满足** |
| Renderer 扩展 fail-closed | **满足** |

## Findings

### F-001 · 客户端 blob 下载文件名消毒弱于服务端 Content-Disposition

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `apps/web/src/renderer/render.tsx:314–318` 仅替换 `[\r\n"]`；服务端 `sanitizeFileHeaderName`（filelibrary.go:313–329）为字母数字/`.`/`_`/`-` 白名单 |
| severity | low（`download` 属性非响应头；现代浏览器会剥离路径分量；非同源 XSS 面） |

**说明**：行 `name` 来自上传者可控的存储文件名。服务端响应头已安全；客户端建议文件名消毒不一致，可能造成文件名混淆（路径分隔符、控制字符、超长名），不是库下载的 header-injection 面。

**建议修复**：与服务端同源规则对齐（allowlist 或至少剥离 `/\\`、控制字符并截断长度），或优先使用响应 `Content-Disposition` 解析结果。

### F-002 · 删除在对象已缺失时提前 404，不清理孤儿 meta

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `filelibrary.go:213–220`：`os.Remove(object)` 若 `ErrNotExist` 立即 404 返回，不执行 `os.Remove(id+".meta.json")`（221 行） |
| severity | low–med（数据完整性；配额扫描按 meta 计费，孤儿 meta 可持续占用 per-user 配额） |

**说明**：正常路径先删对象再删 meta。若进程在两步之间中断、或对象被外部移除而 meta 残留，列表仍显示该行（scan 读 meta），下载/detail 失败，且重复 DELETE 永远无法清 meta。

**建议修复**：对合法 id 同时 best-effort 删除对象与 meta；仅当**两者皆不存在**时返回 `FILE_NOT_FOUND`；若至少删除了 meta 或对象则记 `files.delete` 并 204。

### F-003 · 审计写入静默丢弃错误（弱于仓内其它 handler）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `filelibrary.go:337`：`_ = operations.RecordOperation(...)`；对比 `export.go:239` / `users.go:393–394` 使用 `slog.Error` |
| severity | low（D-002 对 download 写明 best-effort；upload/delete 同样静默） |

**说明**：事件字符串与 0018 CHECK **完全匹配**，正常路径可落盘（测试覆盖）。失败时（DB 故障、未迁移）操作仍成功且无日志，削弱安全审计可观测性。

**建议修复**：与 export/users 一致：`if err := ...; err != nil { slog.Error(...) }`；保持 best-effort 不阻断主路径。

### F-004 · 工具栏「上传」权限意图映射到 files.read 而非 files.write

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `schema/file-library.json:89` table `edit` = files.read；`126–133` toolbar upload `permissionIntent: "edit"`；表单自身 `edit` = files.write（43–44） |
| severity | low（服务端 POST /api/upload 与 POST /api/library/files/upload 均要求 files.write；无绕过） |

**说明**：若 RBAC 将来授予 `files.read` 而不授予 `files.write`，UI 仍显示上传入口（表单内再挡一层）。与 D-002 §5「按 files.read / files.write / files.delete 级联」的意图不完全一致。

**建议修复**：为 table 增加独立 `write`/`upload` 权限键映射 `files.write`，工具栏 `permissionIntent` 指向该键；或拆分 cascade keys。

## 必改项汇总

- **required / 必改**：无
- **recommended**：F-001、F-002、F-003、F-004（均可在后续维护波次处理，不阻断 S-02 安全门禁）

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-002 self（S2–S4，verdict pass） | **同意**核心安全结论；本意见补充 4 条 recommended hardening，不推翻 pass |
| A-001 self（S1） | 方案级安全意图与实现一致 |

## 结论 + 建议给编排器/用户的下一步

**verdict: pass** — 在 D-002 冻结的安全/数据门禁范围内，实现与计划一致：授权 fail-closed、owner-only 控制端点保留、确认端点校验所有者、id 白名单防穿越、下载头加固、配额不可绕过、0018 CHECK 与事件一致、迁移重建安全、Profile 仅为声明的内容扩展、Renderer 白名单与 `{id}` fail-closed。无 required findings。

建议 `/govern`：记录对本 independent 意见的响应；F-001～F-004 可选排入维护或 accepted-residual；不因本意见阻断 S-02 关门（若其它门禁已满足）。

### 声明

本意见 `source: independent`，**不修改**目标 `status` / `progress` / goal-tree / 方案正文；响应与状态变更由 `/govern` 与用户裁决处理。
