---
id: GOAL-010-w9-branding-asset-upload
doc: audit
source: self
date: 2026-08-15
scope: S1 方案冻结～S5 验证 + go 影响判定（S6 关门前置）
verdict: pass
status: closed
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# A-001 · S6 自审（S1～S5 + go 判定）

- **source**：self ｜ **date**：2026-08-15 ｜ **scope**：S1 方案冻结～S5 验证 + go 影响判定
- **verdict**：**pass**（自审必改项 F-001 已 fixed；无未闭合 required）

## 范围与区间

目标五件套 + ledger（D-001、E-001～E-003）+ 全部实现文件与测试（commit `9b751b4`）。

## 成果（有证据）

- S1：D-001 六项用户裁决（2026-08-15 书面）全部 accepted；I-001～I-007 closed。
- S2：专用 brand-assets 仓（独立于 uploads 仓与 file-library）+ POST（settings.write）+ 公开 GET（nosniff/sandbox/immutable）；multipart 限量读取（413 在读取载荷前拒绝）；替换/清空/重置清理 + 启动 GC。
- S3：PNG/JPEG/GIF/WebP 解码（golang.org/x/image v0.45.0 纯 Go）；logo ≤512 / favicon ≤64 等比限幅；透明→PNG、不透明→JPEG(q82)；原始字节永不明文存储/回传；SVG/HTML/脚本三重拒收。参数进 config.yaml branding 节 + BRANDING_* env。
- S4：settings schema 四字段 textarea→upload 控件（actionRef + accept + maxSize，meta 增 actions.upload）；UploadField 移除按钮（I-008 closed）；i18n zh/en + form.upload.remove + 新错误码登记（D-002 契约 + errorcatalog 双语）；旧 URL 兼容读取；/api/branding 契约未变，App.tsx/LoginPage.tsx/branding.ts 零改动。
- S5：Go 全量 `go test ./...` PASS；Web vitest 967/967 PASS；settings.json 通过 runtime AJV（page/action/node schema）校验；活栈 HTTP 级点验（上传→512 JPEG→公开 GET 头→PATCH→/api/branding 投影→替换/清空/重置清理→404）。

## Findings

### F-001（self · required → **fixed**）

替换清理未考虑**多字段共享同一资产 id** 的边界：logoUrl 与 faviconUrl 同时引用同一资产时，仅替换 logoUrl 会删除共享资产，导致 faviconUrl 引用失效（仅 API 直连可构造，UI 每次上传均产生新 id）。

- 修复：`cleanupReplacedBrandAssets` 仅当**没有任何**更新后字段仍引用该 id 时才删除；
- 回归测试：`TestBrandingAssetSharedReferenceSurvivesReplace`（替换后共享资产保留，最后引用清除后删除）；
- 证据：handler 测试全绿（`go test ./internal/handler/ -run TestBrandingAsset` ok）。

### N-001（recommended）

UploadField 移除按钮为共享组件扩展：对所有单文件上传字段生效（users CSV 导入、file-library、dev examples）。属预期延伸（I-008 用户裁决「+ 单字段移除图片」），影响面已在 E-002 留痕。

### N-002（recommended）

活栈点验为 HTTP 级（未做浏览器自动化上传交互）。浏览器级品牌投影由既有 e2e（shell.spec / localization.spec）覆盖；上传控件浏览器交互可留待后续 e2e 增补。

### N-003（信息）

`tsc -b` 的 6 处既有类型错误在 HEAD 即存在（git stash 复核一致，行号随改动位移），与本波无关；vitest 构建链不受影响。

## go 影响判定

**不 held、不暂挂**：未改 Profile 默认集、模块矩阵、Manifest/装配语义或协议 pin；新增路由属 admin.settings 模块贡献面（mvp 仅新增公开读面 GET /api/branding/assets/{id}，且 /api/branding 契约形状未变）。

## 结论

S1～S5 事实与证据可核对；F-001 已修复并回归；无未闭合 required。放行进入 independent 复审（A-002）与关门。
