---
id: GOAL-010-w9-branding-asset-upload
doc: audit-entry
record_id: A-002
source: independent
scope: S1 方案冻结 ～ S5 验证（S6 close-out；security 面必审）
verdict: pass
status: recorded
auditor: grok build（grok-4.6 · high）
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# A-002 · 独立交叉审计 · S6 关门前（2026-08-15）

- **source**：independent
- **auditor**：grok build（grok-4.6 · high）
- **类型** / **scope**：close-out · S1 方案冻结 ～ S5 验证；设置页【品牌】图标 URL→上传的完整实现（专用存储 / 公开 GET / 图像处理 / 清理 / 配置 / schema / UploadField / i18n / 错误码 / 使用点 / 测试）
- **verdict**：**pass**（scope 内无未关闭 high required；无到期未闭合 required 信息项；S1～S5 主张可核对）
- **工作区**：`workspace-010-design-implementation-conformance`（Root `GOAL-001-design-implementation-conformance`；`canonical_scope` 已校验；`shared_materials_catalog: none`；`primary_plan` = VP-010）

## 范围与区间

- **covered**：GOAL-010 五件套 + D-001 + E-001～E-003 + A-001 self；`branding_assets.go` / `_test.go`；`settings.go` 清理挂钩；`composition.go` 接线与启动 GC；`config.go` + `configs/config.yaml` + `config.default.yaml` + `TestBrandingConfig`；`settings/provider.go` + `kernel/profile.go`；`settings.json`；`form-controls.tsx` UploadField + `form-controls.upload.test.tsx`；zh-CN / en-US i18n；`errorcatalog.go` + `error_contract_test.go`；`upload.go` 安全基线对照；`branding.ts` / `App.tsx` / `LoginPage.tsx` 消费；`action.schema.json` UploadAction。
- **独立复跑**：`apps/api` `go test ./... -count=1` **exit 0**；`apps/web` `npx vitest run` **967/967 PASS**；W9 切片 `TestBrandingAsset*` / `TestBrandingConfig` / `TestErrorCodeContract*` / `TestSettingsProviderRegistersSurfaces` / `form-controls.upload.test.tsx` + `schema-keys.structural.test.ts` 均绿。
- **excluded**：未复现 E-003 活栈 HTTP 进程点验（采信 E-003 叙述 + 同源单测路径）；未跑浏览器 e2e / Playwright；未改 status / progress / goal-tree / 方案正文 / 业务代码。
- **P-005**：I-001～I-007 在 D-001 关闭；I-008 / I-009 在 00-meta 与 03-audit 信息表关闭（E-002/E-003）；01-decision.md 信息表仍写 I-008/I-009 `open`（见 F-004，非门禁阻断）。I-006 provider 本轮确认为 grok-4.6。
- **共享资料**：`none`，无引用可核。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 上传 fail-closed：无身份 401、无 `settings.write` 403 | `branding_assets.go` `upload()` → `requirePermission(..., "settings.write")`；`auth.Middleware` 无 token → 401；`resources.go` `requirePermission`；`TestBrandingAssetUploadRejections` |
| 超大文件 413；声明大小不能绕过上限 | `LimitReader(MaxBytes+1)` + `len(body) > MaxBytes` → `FILE_TOO_LARGE` 413；`TestBrandingAssetOversizeRejected`（1 KiB 策略） |
| SVG / 文本 / 不可解码 → 415；原始字节不入库、不回传 | `dangerousInlineTypes` + `containsActiveContent`（`<svg`/`<script`/`<?xml`）+ `decodeBrandingImage` 失败三重门；`save()` 只收 `processBrandingImage` 产物；公开 GET 解出 JPEG/PNG |
| 公开 GET：id 仅 `^[0-9a-f]{32}$`；穿越不 200；nosniff + CSP sandbox + immutable | `load()` / `BrandAssetIDFromURL`；`file()` 头；`TestBrandingAssetGetMissingAndInvalidID`；CSP 头在实现中存在、单测未钉扎（F-002） |
| 限幅不放大；透明→PNG；不透明→JPEG q82；WebP 纯 Go | `processBrandingImage` `scale` 仅当 `maxSide > target`；`imageIsOpaque`；`jpeg.Options{Quality: opts.JPEGQuality}`；`golang.org/x/image/webp` 无 `import "C"` / `#cgo`（模块缓存 v0.45.0） |
| 清理与 I-004 一致：替换/清空即删；重置 DeleteAll；启动 GC 孤儿；共享引用不误删 | `cleanupReplacedBrandAssets` `stillReferenced`；`settingsReset` `DeleteAll`；`composition.go` 启动 `GC` 四列；`TestBrandingAssetCleanupOnReplaceAndClear` / `CleanupOnReset` / `SharedReferenceSurvivesReplace` / `StartupGC` |
| W7 配置：YAML `branding` + `BRANDING_*` env；越界 jpeg_quality 回退 82 | `config.yaml` / `config.default.yaml`；`Load` YAML 层 + `positiveIntEnv`；`NewBrandingAssetStore` 质量 1..100；`TestBrandingConfig` 4 子例 |
| schema：`type: upload` + actionRef；accept/maxSize=4MiB；`actions.upload` | `settings.json` `uploadBrandingLogo/Favicon`；`action.schema.json` UploadAction；`meta.requiredCapabilities` |
| `/api/branding` 契约未变；旧 URL 兼容读取 | `brandingResponse` 字段集未增删；`normalizeLogoURL` 接受同源 `/api/branding/assets/{id}` 与 http(s)；`App.tsx` / `LoginPage.tsx` / `branding.ts` 无本波改动 |
| i18n + 错误码契约 | zh-CN/en-US：`form.upload.remove`、`error.assetNotFound`、`error.invalidKind`、描述文案已去「不支持上传」；`ASSET_NOT_FOUND` / `INVALID_KIND` 入 `frozenLiteralCodes` + catalog |
| S5 回归本轮独立复跑绿 | `go test ./...` exit 0；vitest **967/967** |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结（六项裁决 + 信息表） | 达成 | D-001 accepted；I-001～I-007 closed |
| S2 专用仓 + POST settings.write + 公开 GET | 达成 | `brand-assets` 目录；非 uploads / file-library；mvp 仅 `RegisterPublicBrandingAssets` |
| S3 处理参数 + x/image + config.yaml | 达成 | 实现 + `TestBrandingConfig`；WebP 单测缺口见 F-002 |
| S4 schema/i18n/移除按钮/旧 URL 兼容 | 达成 | settings.json；UploadField；I-008 实现已在 E-002 |
| S5 单测 + 全量回归 + 活栈 | 单测/全量本轮复跑达成；活栈采信 E-003 | 本轮 go/vitest；E-003 |
| I-004 清理闭环 | 达成（含 A-001 F-001 共享引用修复） | settings.go + 4 个清理测试 |
| go 不 held | 同意 A-001：未改 Profile 默认集 / 模块矩阵 / Manifest 装配 / 协议 pin | `profile.go` 仅 admin.settings 路由键增量；mvp 公开读不改 `/api/branding` 形状 |

### 审计重点逐项

1. **鉴权**：POST 经 `a.Middleware` + `requirePermission("settings.write")` 双门；缺身份 fail-closed 401；editor 无 settings.write → 403。公开 GET 故意无鉴权（登录页/外壳预认证加载），与 D-001 取舍一致。
2. **类型欺骗**：不信任客户端 MIME；`DetectContentType` 危险类型 + 全文 active-content 标记 + 解码失败。SVG（含 `<script>`）与纯文本已测 415。HTML 走 `text/html` 硬拒绝，无单独用例（F-002）。重编码后只可能是 PNG/JPEG 字节。
3. **公开 GET**：id 形状与通用上传仓同一正则，阻断 `..` / Windows 盘符。头：`nosniff` + `CSP: sandbox` + `immutable`。meta 缺失时 Content-Type 回落 `application/octet-stream`（偏安全）。
4. **图像处理**：等比缩小、不放大；透明 PNG / 不透明 JPEG(q82)。GIF 为额外输入（D-001 写 PNG/JPEG/WebP；用户 scope 含 GIF）。**解码前无 `DecodeConfig` 像素上限**（F-001）。
5. **清理**：替换/清空按字段 diff 且检查其余列是否仍引用；重置全清；启动 GC 忽略 legacy http(s) URL。运行期未提交的上传依赖下次启动 GC（D-001 已定）。
6. **配置 / schema**：单一 YAML 权威 + env 覆盖；越界质量回退。UploadAction accept/maxSize 合法。前端 accept 不含 GIF，服务端更宽，不构成安全缺口。
7. **消费点**：`/api/branding` 形状未变；`safeBrandingUrl` 接受 `/api/branding/assets/{id}`；错误码已钉 D-002 冻结集。

## Findings

### F-001 · 解码前无像素上限（压缩炸弹）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | high |
| status | open |
| evidence | `branding_assets.go` `decodeBrandingImage` / `processBrandingImage`：直接 `png/jpeg/gif/webp.Decode` 后再缩放，无 `DecodeConfig` 与最大宽高门禁 |
| 影响门禁 | 不阻断 S6 关门。需 `settings.write`；压缩体已 ≤4 MiB。残余：高权限账号可用小体积、超大声明尺寸的图像触发大内存分配 |
| closure | — |

建议：先 `DecodeConfig`，拒绝超过约定上限（例如边长数千）的输入，再全量解码。

### F-002 · 处理与安全头的单测钉扎不完整

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | med |
| status | open |
| evidence | `branding_assets_test.go` 仅 PNG 输入；公开 GET 断言 nosniff/Cache-Control，**未断言** `Content-Security-Policy: sandbox`；无 WebP/JPEG/GIF 输入、无不放大（小图保持原尺寸）、无 HTML 专用 415 例 |
| 影响门禁 | 不阻断。实现路径存在；本轮全量测试绿 |
| closure | — |

### F-003 · `golang.org/x/image` 在 go.mod 标为 indirect

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | `branding_assets.go` 直接 import `golang.org/x/image/draw` 与 `.../webp`；`apps/api/go.mod` 将其放在 indirect 块 |
| 影响门禁 | 不阻断。模块在 go.sum，构建/测试通过。`go mod tidy` 应升为 direct |
| closure | — |

### F-004 · 01-decision.md 信息表与 00-meta / 03-audit 不一致

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | `01-decision.md` I-008/I-009 仍 `open` /「待确认」「待实施」；`00-meta.md` 与 `03-audit.md` 信息表已 closed（E-002/E-003） |
| 影响门禁 | 不阻断 P-005（权威关闭证据在 00-meta + E-002/E-003）。编排响应时应把决策索引表与 meta 对齐 |
| closure | — |

## 必改项汇总

无。本意见 **0** 条开放 required finding。

## 与既有意见的异同

- **A-001 self · pass**：同意 S1～S5 事实、go 不 held、以及 F-001（共享引用误删）已 fixed——本轮核对 `stillReferenced` 与 `TestBrandingAssetSharedReferenceSurvivesReplace` 并通过 handler 复跑。
- **A-001 N-001**（UploadField 移除按钮作用于全部单文件上传）：同意，属 I-008 预期延伸，不升级。
- **A-001 N-002**（无浏览器自动化上传）：仍成立；本审亦未跑 Playwright。不升为 required。
- **本审新增** F-001（解码像素上限）、F-002（测试钉扎）、F-003（go.mod）、F-004（决策信息表漂移）。均 recommended。

## 结论与建议下一步

S1～S5 在 security 面上可核对：鉴权 fail-closed、类型欺骗与原始字节不回传成立、公开 GET 头与 id 形状成立、清理与 I-004（含共享引用）成立、配置/schema/i18n/错误码/消费点一致，全量回归本轮独立复跑绿。

**建议 `/govern`**：响应本意见；recommended 项可本波顺手修或记入后续卫生；无 required 阻断后，用户确认即可勾 S6 并走关门（勿用 progress% 代替检查点）。

## 声明

本意见不修改 status / progress / 检查点 / 方案正文 / goal-tree；响应由 `/govern` 处理。
