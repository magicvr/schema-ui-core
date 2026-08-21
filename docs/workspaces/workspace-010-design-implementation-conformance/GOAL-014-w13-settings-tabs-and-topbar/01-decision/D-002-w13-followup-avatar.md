---
id: D-002
doc: decision
status: accepted
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-002 · 追加：移动端汉堡靠左修正 + T-05 个人中心头像上传

## 背景

2026-08-16 用户在本波关门后追加两项：① 移动端打开左侧导航的汉堡按键应靠在工具栏最左（用户直觉）；② 给目标加任务——个人中心【资料】选项卡增加头像上传（用户明示「极其附属功能」）。GOAL-014 重开承接（status active，S2～S4 重新开放）。

## 决策

### ① T-02 修正 · 移动端汉堡按键靠工具栏最左

- App.tsx：汉堡按钮从右对齐的功能区容器（`ml-auto`）中移出，置于主行最左（<lg 时品牌链接隐藏，汉堡即首元素）；功能区其余控件保持右对齐。
- ≥lg 布局不变（品牌链接占据左侧，汉堡隐藏）。
- 与既有断点约定一致（<lg = 移动端）。

### ② T-05 · 个人中心【资料】选项卡头像上传

**产品语义**：自服务附属功能——登录用户可在个人中心【资料】选项卡上传自己的头像；头像展示在顶栏用户菜单触发器上。

**后端（Go）**：
- 复用 W9 的品牌资产存储模式，抽象为共享 `RasterAssetStore`（`handler/raster_assets.go`）：品牌（logo/favicon）与头像（avatar）同走「服务端重编码 + 专用目录 + 公开 GET（nosniff/sandbox/immutable）」安全模型；头像最长边 256px（`BrandingAssetsOptions.AvatarDim` 默认 256）。
- 新端点：`POST /api/account/avatar`（仅认证，自服务无权限键；上传成功后 best-effort 删除旧头像文件）+ 公开 `GET /api/account/avatars/{id}`。
- `users` 表新增 `avatar_url` 列（admin.account 迁移 **0035**，全局台账末尾追加）；`operation_log` 事件 CHECK 增加 `account.avatar-change`（core.operationlog 迁移 **0036**，重建式，同 0032/0034 模式）。
- `authsession.User`/`UserPatch` 增加 `AvatarURL`；`GET/PATCH /api/account/profile` 支持 `avatarUrl`（PATCH 校验：必须为空或头像存储 URL，其他 URL 一律 400；替换/清空时 best-effort 删除旧文件）；`/me` 用户快照携带 `avatarUrl`（`account.User.AvatarURL` + `accountFromUser`）。
- 无启动 GC：替换/清空路径已覆盖常见场景；崩溃遗留最多一个孤儿文件（可接受，见 E-005 记录）。

**前端（Web）**：
- `account.json`：`actions.uploadAvatar`（type upload → /api/account/avatar，accept 图片、4 MiB 上限）+ 资料表单增加 `avatarUrl` upload 字段（`recordSource`/`saveProfile` 映射）+ meta 能力 `actions.upload`。
- i18n：`schema.account.field.avatarUrl`（en: Avatar / zh: 头像）。
- `parseAuthUser` 解析 `avatarUrl`；UserMenu 触发器在头像存在时渲染圆形头像图（否则保留首字母）。

**未选方案**：
- 不新增独立头像存储实现（复用共享存储，避免重复安全敏感代码）；
- 不做头像裁剪/圆角服务端处理（客户端展示层圆角即可；附属功能最小化）；
- 不把头像纳入管理端用户列表/详情（用户只要求个人中心自服务）。

## 影响

- **go 判定**：新增两条 admin.account 路由 + users 表追加列 + 事件 CHECK 扩展——均为能力追加，不改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义 → **VP-008 go 无影响、不暂挂**。
- **测试影响**：Go 新增 avatar 处理器测试（上传/公开读取/拒绝/替换/清空/非法 URL）+ 迁移台账（35/36）；web 增加 UserMenu 头像渲染与 e2e 头像全链路；既有 branding 测试作为共享存储重构回归护栏。
