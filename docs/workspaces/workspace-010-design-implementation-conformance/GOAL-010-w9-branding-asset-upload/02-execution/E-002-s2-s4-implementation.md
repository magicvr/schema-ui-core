---
id: GOAL-010-w9-branding-asset-upload
doc: execution
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-002 · S2～S4 实施（存储/端点/处理/配置/前端）

2026-08-15 完成 S2～S4 实施，证据如下：

## S2 · 专用资产存储与端点

- 新增 `apps/api/internal/handler/branding_assets.go`：`BrandingAssetStore`（专用 `<data>/brand-assets` 目录，独立于通用 uploads 仓与 admin.file-library）+ `POST /api/branding/assets?kind=logo|favicon`（settings.write 门禁）+ `GET /api/branding/assets/{id}`（公开、nosniff + CSP sandbox + immutable 缓存）。
- multipart 读取用 `MultipartReader/NextPart` + `LimitReader(MaxBytes+1)`：超大文件在读取载荷前即拒绝（413），声明大小说谎也无法绕过上限。
- 路由接线：settings 模块 provider（`BrandingAssetRoutes`）+ mvp 无 settings 模块时 `RegisterPublicBrandingAssets`（公开读不依赖编辑模块）；kernel profile admin.settings 路由键同步。
- 清理（I-004）：settingsPatch 替换/清空即删旧资产；settingsReset 清空全部资产；组合根启动 GC 删除未被 site_settings 引用的孤儿文件。

## S3 · 自动图像处理

- 解码 PNG/JPEG/GIF（stdlib）+ WebP（新增 `golang.org/x/image` v0.45.0，纯 Go 无 CGO）；SVG/HTML/script 经 sniff + active-content 标记 + 解码失败三重拒收。
- 服务端重编码：logo ≤512px、favicon ≤64px（等比、不放大，CatmullRom 高质量缩放）；含透明 → PNG，不透明 → JPEG（质量 82）；**永不回传原始字节**。
- 参数进 W7 config.yaml `branding:` 节（max_bytes / logo_max_dimension / favicon_dimension / jpeg_quality），env 覆盖（BRANDING_*）。

## S4 · 前端与消费点

- `settings.json`：品牌四字段 textarea → `type: upload` + actionRef（`uploadBrandingLogo` / `uploadBrandingFavicon`，accept image/png,image/jpeg,image/webp，maxSize 4 MiB）；meta 增 `actions.upload`；描述文案改为上传语义。旧 URL 值兼容读取（DB 列与 /api/branding 契约未变）。
- UploadField 扩展（I-008 closed）：单文件上传字段值非空时渲染 i18n 化「移除图片」按钮（multiple/readOnly/disabled 不显示）；表单提交仍为原有 PATCH。
- i18n：zh-CN/en-US 描述与字段文案 + `form.upload.remove` + 新增错误码文案（`error.assetNotFound` / `error.invalidKind`，同时登记 D-002 错误码契约）。
- 消费点核对：/api/branding 契约形状未变；App.tsx / LoginPage.tsx / branding.ts 零改动（URL 形态一致）。
