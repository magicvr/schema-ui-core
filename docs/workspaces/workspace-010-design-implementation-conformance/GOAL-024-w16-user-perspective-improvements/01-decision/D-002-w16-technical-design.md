---
id: D-002
goal: GOAL-024-w16-user-perspective-improvements
title: W16 技术方案与接口设计（W16-F01～W16-F10）
status: approved
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-002 · W16 技术方案与接口设计

## 1. 触发

在 D-001 台账冻结后，为 W16-F01～W16-F10 逐项明确 Schema 契约、API 路由、前端交互与存储迁移方案。同时关闭 `00-meta` I-001：现有 `apps/web/src/renderer` 已具备 custom-node 注册表（`registerCustomComponent`）与 schema 驱动表格/表单能力，因此各项实现路径可统一为「优先原生 schema 扩展，特殊交互使用 custom component」。

## 2. Renderer 兼容路径（I-001 结论）

- **原生 schema 扩展**：新增字段/列属性/表单控件属性即可承载的项，不改 Renderer 中央注册路径。
- **custom component**：文件预览/复制、Cron 解析预览、监控自动轮询、导入错误明细、MFA 密钥复制与恢复码下载等需要业务交互或剪贴板/定时器的项，注册 `{type:"custom", component:"<key>"}` 节点。
- **不修改 Host/App 协议**：所有新增前端能力均在既有 Renderer 白名单与 custom-node 机制内实现，不改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义。

## 3. 逐项技术方案

### W16-F01 · 首次登录强制修改初始密码

- **存储**：`users` 表新增 `must_change_password INTEGER NOT NULL DEFAULT 0`（0/1）；种子 `admin` 与新建/导入用户默认 `1`，自改密成功后置 `0`。
- **API**：
  - `POST /api/auth/login` 成功响应增加 `mustChangePassword: boolean`。
  - 现有 `POST /api/account/password` 改为在 `must_change_password=1` 时允许使用初始密码作为 `currentPassword` 完成改密；改密成功后清标记并 bump `token_version`。
  - 新增后端强制门禁：当用户 `must_change_password=1` 时，除 `POST /api/auth/login`、`POST /api/account/password`、`GET /api/account/profile` 等必要端点外，其余业务 API 返回 `403 MUST_CHANGE_PASSWORD`（或由前端拦截 + 后端兜底）。
- **前端**：`AuthContext` 识别 `mustChangePassword` 后进入强制改密页，复用现有改密表单（`currentPassword + newPassword`），成功后刷新会话并进入系统。

### W16-F02 · 文件库在线直接预览与一键复制直链

- **API**：文件库资源行已含可访问 URL 或可通过现有 GET 端点组装；无需新增存储字段。
- **前端**：
  - 文件库表格操作列增加 `preview` 与 `copyLink` 两个 row action。
  - `preview`：图片类型打开 Lightbox（custom component `file-preview-lightbox`）；PDF/文本使用 `window.open(url, "_blank")`。
  - `copyLink`：`navigator.clipboard.writeText(url)` + Toast 提示。
- **Renderer**：通过 schema 表格 `actions` 声明 + 自定义 action handler 或 custom component 实现；不新增协议。

### W16-F03 · 数据导入弹窗下载 CSV 示例模板与逐行错误定位反馈

- **API**：
  - 新增 `GET /api/import/template`（或按数据类型的 `GET /api/{resource}/import/template`）返回 CSV 模板文件。
  - 导入接口响应结构扩展为 `{ ok, imported, errors: [{ rowNumber, field, reason }] }`（现有 `errors` 保留兼容字段，新增 `fieldErrors` 或直接升级 `errors` 为明细对象，冻结时二选一）。
- **前端**：导入模态框增加“下载 CSV 模板”链接；校验失败后以表格/列表渲染 `rowNumber + field + reason`，便于逐行定位。

### W16-F04 · 钱包余额金额“分转元”标准格式化展示与调账警示

- **API**：金额传输保持最小货币单位（分）不变，不破坏契约。
- **前端**：
  - `SchemaTableColumnSpec` 增加可选 `format?: "currency"`（或 `amountUnit?: "cents"`），渲染时除以 100 并保留两位小数、千分位。
  - 调账/扣款模态框输入负数或超过阈值（如单笔 ≥ 100000 分）时显示高亮警示与二次确认。

### W16-F05 · 定时任务 Cron 表达式自然语言解析与未来运行时间预览

- **API**：新增 `POST /api/scheduled-tasks/cron/preview`，请求 `{ cron }`，响应 `{ description, nextRuns: ["2006-01-02T15:04:05Z", ...] }`（复用 Go 端 5-field cron 解析，计算未来 3 次）。
- **前端**：Cron 表单字段下方挂 custom component `cron-preview`，输入防抖调用该端点，展示中文描述与未来运行时间。

### W16-F06 · 系统监控页面支持定时自动轮询刷新

- **方案**：给监控页面的表格/指标节点增加可选 `refreshInterval?: number`（毫秒，如 5000/10000/30000），`SchemaTable` 在开启时按间隔静默 refetch。
- **前端**：监控页面顶部提供 `关闭 / 5秒 / 10秒 / 30秒` 下拉，写入页面局部状态并透传为节点 `refreshInterval`；不新增 API。

### W16-F07 · 个人中心活跃会话支持“一键下线其他所有设备”

- **API**：新增 `POST /api/account/sessions/revoke-others`，请求头携带当前 refresh token 或会话 ID；后端吊销除当前会话外的全部活跃 refresh tokens，并 bump `token_version`（注意不能把当前请求自己也吊销）。
- **前端**：个人中心会话卡片头部增加“下线其他设备”按钮 + 确认对话框，成功后刷新会话列表并 Toast。

### W16-F08 · 登录算术验证码主动刷新 + MFA 密钥复制与恢复码 txt 下载

- **前端**：
  - 登录页验证码区域增加“换一题”按钮，调用现有验证码生成接口（或前端重新加载 captchaId/图片），无需后端新契约。
  - MFA 绑定弹窗增加“复制密钥”按钮（`navigator.clipboard`）与“下载恢复码 (.txt)”按钮（生成 Blob 下载）。

### W16-F09 · 数据字典项支持 Badge 颜色/标签风格配置

- **存储/API**：`dict_entries` 表新增 `badge_style TEXT NOT NULL DEFAULT 'default'`；条目 CRUD 的 schema/API 增加 `badgeStyle` 字段，取值 `default|success|warning|destructive|info`。
- **前端**：`SchemaTableColumnSpec` 增加可选 `badgeStyleField`；当列声明该字段时，单元格渲染为对应颜色的 Badge（无自定义组件也可完成，但可提供 custom cell renderer 统一实现）。

### W16-F10 · 系统设置支持页脚版权文字与备案号

- **存储/API**：`settings` 资源增加 `copyrightText` 与 `icpNumber` 字段（通用设置 schema），`GET /api/settings` 与 `PATCH /api/settings/{id}` 自动覆盖。
- **前端**：Shell/`App.tsx` 底部统一读取设置并渲染页脚；未配置时隐藏或显示默认文案。

## 4. 风险与兼容

- F01 属安全门禁：实施需独立审计；涉及登录响应与业务 API 门禁，需回归登录/会话/用户管理全链路。
- F03 若调整导入错误响应结构，需与现有导入调用方（前端、测试）兼容；倾向新增字段而非破坏 `errors`。
- F04/F09 为展示层扩展，低风险；需补充列渲染快照/组件测试。
- F10 仅设置字段与页脚渲染，低风险。
- 所有 API 新增端点需在对应模块 `Descriptor().Contributions.Routes` 登记并补权限/操作日志。

## 5. 结论

本方案可作为 S2 技术方案冻结输入。I-001 已通过代码核验（`apps/web/src/renderer/custom-components.ts`、`schema-table.tsx`、`form-controls.ts(x)`）关闭：现有 Renderer 支持 custom-node + 原生 schema 扩展两条路径，无需新增协议或中央渲染框架。
