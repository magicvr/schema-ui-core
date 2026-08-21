---
id: GOAL-010-w9-branding-asset-upload
doc: execution
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-004 · S6 审计响应（A-001 F-001 + A-002 F-001～F-004 全部处置）

2026-08-15 完成 S6 cross 审计响应，全部 finding 处置：

## A-001（self）响应

- F-001（required）→ **fixed**：`cleanupReplacedBrandAssets` 仅在无其他字段仍引用时删除替换资产；回归测试 `TestBrandingAssetSharedReferenceSurvivesReplace`。

## A-002（grok build independent · grok-4.6 · high · pass）响应

| finding | level | 处置 | 证据 |
|--------|-------|------|------|
| F-001 解码前无像素上限（压缩炸弹） | recommended | **fixed**：DecodeConfig 头读取预检，最长边 >8192 直接 415（`maxBrandingInputDimension`）；测试 `TestBrandingAssetRejectsOversizedDimensions`（构造 30000×30000 IHDR + 合法 CRC） | `branding_assets.go` / `_test.go` |
| F-002 WebP/JPEG/GIF/CSP 单测未钉 | recommended | **fixed**：新增 JPEG 与 GIF 输入测试 + GET 断言 `Content-Security-Policy: sandbox`；WebP 解码路径由 x/image 覆盖（无纯 Go 编码器可造夹具，留作已知限制） | `TestBrandingAssetJpegAndGifInputs` / `TestBrandingAssetUploadAndPublicServe` |
| F-003 x/image 标为 indirect | recommended | **fixed**：`go mod tidy` 后升为 direct require | `go.mod` |
| F-004 01-decision 信息表 I-008/I-009 未同步 | recommended | **fixed**：索引表与 00-meta / 03-audit 对齐为 closed | `01-decision.md` |

## 复核

- handler 测试 11 例全绿；`go test ./...` 全量复跑（见 E-005/关门证据）；vitest 967/967 不受后端修复影响。
