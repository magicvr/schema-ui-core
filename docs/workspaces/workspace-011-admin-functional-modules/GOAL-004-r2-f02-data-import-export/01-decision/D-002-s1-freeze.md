---
id: D-002
goal: GOAL-004-r2-f02-data-import-export
title: S1 · 方案冻结 — admin.data-transfer 模块（导出 CSV + 导入校验/错误报告 + 权限键 + 必办-1 协议对照）
date: 2026-08-14
status: accepted
parent: GOAL-004-r2-f02-data-import-export
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-002 · S1 方案冻结（F-02 数据导入/导出 · 共享能力）

> 依据：I-011-001 `3 F-02、`8 必办-1；GOAL-004 00-meta 边界与 S1 门禁。

## 1. 模块设计：`admin.data-transfer`（共享能力模块）

| 项 | 冻结值 |
|----|--------|
| 模块 ID | `admin.data-transfer` |
| 依赖 | core.auth-session / core.navigation-capability / core.schema-render / core.operationlog |
| 能力 | StandardAdminCapabilities() |
| 路由 | `GET /api/export/{resource}`、`POST /api/import/{resource}` |
| 权限键 | `data.export`（PolicyAdminEditor：读类数据面，editor 可用）、`data.import`（PolicyAdmin：写类数据面，仅 admin） |
| 页面 | 无独立页面（共享能力；入口在 users/roles 页工具栏，跨模块动作引用） |
| 持久化 | 无新迁移（复用 uploads 目录 + operationlog） |

**R2 范围**：导出支持 `users` / `roles` 两个代表资源；导入支持 `users`（roles 导入涉权限引用风险，归 R3 再议）。Excel 格式（xlsx 库依赖）不在 R2（CSV 优先；I-003 关闭）。

## 2. 必办-1 · 独立协议对照（v2.8.0 面，A-002 F-004 口径）

| 对照对象 | 证据 | 结论 |
|----------|------|------|
| `docs/schemas/node.schema.json` `Permissions` 扩展动作键 | 行 222：`扩展动作键（如 approve/export）；协议层不限制键名，建议由接入方通过项目级 CI/Review 约束` | **协议显式允许 export 作扩展键**；无导出契约定义 |
| grid-dashboard（protocol-inventory `2.5 信息性场景） | 非语义权威，范例候选 | 与导出无关；呈现自由 |
| user-profile-*/order-list-* 上游样例（Manifest 点名） | 非语义权威（`_samples/`） | 呈现自由，不构成导出契约 |
| 上游 actions/request 契约（request-construction fixtures） | 本地实现按 fixture 校验 | 动作方法/URL 契约复用；无导出语义 |

**处置**：协议对「导出/导入」**无契约** → **本地契约**（`/api/export/*`、`/api/import/*` 为本模块定义的本地端点，非协议面）；**fail-open 留痕**：模块未启用时页面按钮因权限键缺失而禁用（视觉 fail-open），直打 API 404/403（服务端 fail-closed）。文档引用本对照。

## 3. 导出设计（GET /api/export/{resource}）

- 参数：`q` / `sort` / `order`（与列表端点同口径）；`pageSize` 上限 **10000**（超出即 INVALID_PAGE_SIZE；全量导出由管理员自行加筛选或分页导出——文档化上界）。
- 响应：`text/csv; charset=utf-8` + UTF-8 BOM + `Content-Disposition: attachment; filename="<resource>.csv"`；表头 = 列表行字段（稳定顺序：users = id,username,name,roles,enabled,locked,createdAt,updatedAt；roles = id,key,name,system,permissions,menuItems,assignedUsers,editable,deletable,createdAt,updatedAt）。
- CSV 转义：含 `,` / `"` / 换行的字段加引号并双写引号（RFC 4180）；数组字段（roles/permissions/menuItems）JSON 序列化。
- 权限：`data.export`；匿名 401 / 无键 403（requirePermission）。
- 审计：operationlog `data.export`，detail `{"resource":"...","rows":n}`。
- 实现：复用资源仓库（users/roles 的 List 过滤逻辑，PageSize=cap）直接读库，不经 resource 工厂（避免 JSON 双编码）。

## 4. 导入设计（POST /api/import/users）

- 请求：`{"fileId":"<upload id>"}`（复用既有上传基建 C-09：multipart 上传 → `/api/upload` 返回 id；**owner 校验**——fileId 必须属于当前 actor，缺失/外籍 → `FILE_NOT_FOUND` 404；大小上限 2 MiB → `FILE_TOO_LARGE`）。
- 解析：CSV（RFC 4180），首行表头；字段白名单 `username` / `name` / `roles` / `password`（可选）；未知表头列 → 忽略并计入告警（不失败）。
- **不回滚语义（按方案）**：逐行校验 + 逐行应用（每行一个事务内 insert）；错误行收集，**不回滚已成功行**；响应 200 `{applied, failed, total, errors:[{row, message}]}`（errors 仅含失败行，1-based 行号）。
- 行级校验：username/name 非空；password 8–72 字节（缺失 → 行错误 `password required`；**安全留痕：CSV 携带密码属敏感通道，导入为 admin-only 操作**）；roles 逗号分隔且必须存在（INVALID_ROLE_REF 语义）；username 冲突 → 行错误 `username taken`（不中断其它行）。
- 权限：`data.import`。
- 审计：operationlog `data.import`，detail `{"resource":"users","applied":n,"failed":m}`。
- 前端反馈：成功 toast（静态文案）+ reload；详细错误报告在响应与操作日志 detail 中留痕（**错误报告 UI 归 R3**——与 S-02 文件库/独立传输页一起设计，本方案明确此边界）。

## 5. 前端设计

- **Renderer 本地扩展（download 行为）**：请求动作 `onSuccess.behavior: "download"` → 成功后以 blob 触发浏览器下载（文件名取 `Content-Disposition`）。本地契约：上游行为集（toast/reload/navigate/closeModal）不变、未知行为 fail-open（现状 emitBehavior default 忽略）；本扩展仅在有显式 `behavior:"download"` 时生效。需 conformance 文档留痕（本地扩展，非协议新增）。
- **users 页**：工具栏 `Export`（`GET /api/export/users` + download；`permissions.edit = contains "data.export"`）+ `Import`（modal：upload 控件 + `submitImport` POST `/api/import/users` `{fileId}`；`permissions.edit = contains "data.import"`）。
- **roles 页**：工具栏 `Export`（`GET /api/export/roles` + download；`data.export`）。
- **i18n**：`schema.users.toolbar.export/import`、`schema.roles.toolbar.export`、`schema.users.import.*`（modal 标题/字段/提交/反馈）、`feedback.downloadStarted` 等键（en/zh 同步）。
- **fail-open**：admin.data-transfer 未启用 → 权限键不存在 → 按钮禁用（视觉）+ 404/403（服务端）；与 F-03 同一模式。

## 6. Profile 声明

- `admin.data-transfer` 加入 **admin** 默认集（导出/导入是管理面能力，mvp 不引入——与 F-03 的账号安全基线不同；F-02 属功能扩展）。**Profile 内容扩展**声明同 D-002 `6 模式（F-03）：不改装配语义。
- `adminFunctionalOrder` **不追加**（无页面贡献，不影响 home 推导）。
- smoke.sh SM-007 admin 页面集不变化（无新页面）。

## 7. 必办核对（I-011-001 `8）

| 必办 | 适用 | 处置 |
|------|------|------|
| **必办-1（独立协议对照）** | **适用** | **✅ 本方案 `2**（node.schema.json 扩展键 + grid-dashboard + 上游样例 → 本地契约 + fail-open 留痕） |
| 必办-2/3/4/5 | 其它目标 | 不适用 |

## 8. 未选方案（留痕）

- 不引入 xlsx 二进制依赖（R2 CSV；Excel 归 R3 评估）。
- 不做导入预览/两段式（importId + 错误表）——schema 引擎数据源静态，动态 importId 无法表达；R2 采用单段提交 + 结构化错误报告（响应 + operationlog），完整报告 UI 归 R3。
- 不做「整文件失败回滚」选项（R2 固定为不回滚语义；需求出现再按 P-004 加）。
- 不做服务端文件清理任务（上传文件生命周期归 R3 文件库 S-02）。
- 导入密码列保留（admin-only 通道）但不做密码哈希列（不存在该字段）。

## 9. 实现范围（S2 清单）

1. `handler/export.go`：GET /api/export/{resource}（users/roles CSV；cap 10000；转义；审计）。
2. `handler/import.go`：POST /api/import/users（fileId + owner 校验 → CSV 解析 → 逐行校验/应用 → 错误报告；审计）。
3. upload store 暴露 load 能力给 import（同包 helper）。
4. 模块 provider `admin.data-transfer` + 权限贡献（data.export/import）+ 无 fragment。
5. profileDefaults admin += `admin.data-transfer`；composition 装配。
6. users.json/roles.json 工具栏动作；renderer download 行为 + 测试；i18n。
7. 测试：Go（导出形状/转义/cap/权限/审计；导入校验/部分成功/错误报告/owner 校验/审计）+ Web（download 行为、i18n）。