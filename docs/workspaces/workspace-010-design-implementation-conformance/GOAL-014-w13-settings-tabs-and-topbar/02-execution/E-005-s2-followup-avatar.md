---
id: E-005
doc: execution
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-005 · S2 追加实施：移动端汉堡靠左 + T-05 头像上传

## ① 移动端汉堡靠左（T-02 修正）

- `apps/web/src/app/App.tsx`：汉堡按钮移出右对齐功能区容器，置于主行最左（<lg 时品牌链接隐藏，汉堡即首元素）；功能区其余控件（亮暗/语种/铃/用户）保持 `ml-auto` 右对齐；≥lg 布局不变。

## ② T-05 · 个人中心头像上传

### 后端

- **共享存储抽象**：`apps/api/internal/handler/raster_assets.go`（新）——把 W9 品牌资产存储泛化为 `RasterAssetStore`（dir + opts + urlPrefix + kinds→目标尺寸），`NewBrandingAssetStore`/`NewAvatarAssetStore` 两个构造器；`branding_assets.go` 精简为选项/路由/门禁（settings.write 门禁移到路由包装层）。既有 branding 测试原样全绿（重构回归护栏）。
- **新端点**：`handler/account_avatar.go`——`POST /api/account/avatar`（仅认证，自服务；成功上传后 best-effort 删除旧头像文件；记录 `account.avatar-change` 审计事件）+ 公开 `GET /api/account/avatars/{id}`（nosniff/sandbox/immutable 同品牌资产）。头像最长边 256px（`AvatarDim` 默认），PNG/JPEG/GIF/WebP 输入服务端重编码。
- **迁移**：admin.account **0035** `account_avatar_url`（`ALTER TABLE users ADD COLUMN avatar_url TEXT NOT NULL DEFAULT ''`，checksum `a2872e…`）；core.operationlog **0036** `operation_log_avatar_events`（事件 CHECK 重建追加 `account.avatar-change`，checksum `3f4a67…`）。
- **持久化与快照**：`authsession.User`/`UserPatch` 增加 `AvatarURL`；4 处用户行 SELECT/scan 与 UPDATE 全部接入 `avatar_url`；`GET/PATCH /api/account/profile` 支持 `avatarUrl`（PATCH 校验仅接受空值或头像存储 URL，其他 URL 400；替换/清空时 best-effort 删除旧文件）；`/me` 用户快照携带 `avatarUrl`（`account.User.AvatarURL` + `accountFromUser`）。
- **组合根**：`composition.go` 构造 `avatars/` 目录的 avatar store 并注入 admin.account provider；`kernel/profile.go` 与 account provider 声明新增两条路由。

### 前端

- `account.json`：`actions.uploadAvatar`（upload → /api/account/avatar；accept 图片；4 MiB）+ 资料表单 `avatarUrl` upload 字段（recordSource/saveProfile 映射）+ meta `actions.upload` 能力。
- i18n：`schema.account.field.avatarUrl`（en Avatar / zh 头像）。
- `auth-client.ts` `parseAuthUser` 解析 `avatarUrl`（此前 /me 快照被丢弃——e2e 首跑发现）；`App.tsx` UserMenu 触发器在头像存在时渲染圆形头像图，否则保留首字母。

### 测试

- `handler/account_avatar_test.go`（新，7 例）：上传/公开读取（256 限幅、PNG 透明保留、nosniff 头）、拒绝（匿名 401/SVG 415/空 400/文本 415）、提交-替换-清空（旧文件删除）、非法 URL 400、缺失 404。
- `store` 迁移台账测试更新（35/36；fresh/restart/reopen 计数）。
- `shell.spec.ts`：真实 API 头像全链路（上传 → PATCH profile → reload → 用户菜单头像 img 可见）。
