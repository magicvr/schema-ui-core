---
id: GOAL-010-w9-branding-asset-upload
doc: decision
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# D-001 · 品牌图标上传方案冻结

> 决策日期：2026-08-15；用户书面确认（workspace-010 新增子目标讨论轮：先讨论、确认后落盘）。

## 背景

设置页【品牌】四个图标字段现行以「填写 URL 链接」的 textarea 方式配置（历史遗留，schema 文案自述 "Branding assets are URLs (upload not available)"）。用户裁定改为**上传控件**，并自动处理上传图像（用户图像可能不规范、也可能过大影响网页性能）。

## 决定（六项用户裁决，全部 accepted）

1. **URL 输入彻底移除**：settings schema 品牌字段不再提供 textarea/URL 输入；改为上传控件 + 单字段「移除图片」。已存旧 URL 值**仅兼容读取**（仍可显示/使用），**不迁移、不转换**。
2. **处理参数**：logo 上限 512px（等比限幅）、favicon 64px、JPEG 质量 82、单文件 ≤4 MiB；输入支持 PNG/JPEG/WebP（新增纯 Go 依赖 `golang.org/x/image` 解码 WebP）；输出 PNG（含透明）/ JPEG（不透明）。
3. **SVG 不支持**：延续上传仓安全基线（active content / stored XSS 拒收），SVG 上传返回 UNSUPPORTED_FILE_TYPE。
4. **清理语义**：字段替换时删除旧资产；恢复默认（reset）时清空全部品牌资产；启动时 GC 未被 site_settings 引用的孤儿文件（兜底取消上传/崩溃）。
5. **配置位置**：处理参数进 W7 `config.yaml` 体系（branding 节），env 可覆盖。
6. **审计模式**：关门 cross 审计（self + independent；provider 历史惯例 grok，S6 执行时确认）。

## 关键取舍（含未选方案）

| 取舍 | 选择 | 未选 | 理由 |
|------|------|------|------|
| 存储位置 | 独立 brand-assets 仓 + 公开读端点 | admin.file-library；通用 /api/upload | 品牌图为公开资源（登录页/外壳启动前加载）；文件库与通用仓按 owner 隔离/模块语义不符；独立仓可强控处理与清理 |
| 回传内容 | 仅服务端重编码产物 | 原始字节直存直回 | 公开服务用户派生内容，重编码把 XSS/内容注入面降到最低 |
| 字段值形态 | 仍为 URL 字符串（`/api/branding/assets/{id}`） | 引入新 id 列/表 | DB 结构、/api/branding 契约、前端消费逻辑零改动；旧 URL 兼容读取自然成立 |
| 上传交互 | schema 驱动 `type: upload` + actionRef（现有协议能力） | 自研表单控件 | users/file-library/dev.examples 模块已用同模式，渲染器零改动 |
| Favicon 输出 | 64px PNG（现代浏览器支持 PNG favicon） | 生成 .ico | 避免引入 ICO 编码依赖 |
| 参数配置 | config.yaml（env 覆盖） | 硬编码常量 | 符合 W7 单一配置权威；部署可调 |

## 附录 A · 消费点清单（I-007 证据，2026-08-15 创建时扫描）

- API：`/api/branding` 公共投影（handler/settings.go）；`site_settings` 表 4 列（modules/settings/repository）；settings schema（schema/settings.json 品牌表单 + actions.updateBranding）
- 前端：`apps/web/src/app/branding.ts`（favicon 应用）；`App.tsx` 外壳 Logo（浅色/深色/默认）；`LoginPage.tsx`
- i18n：zh-CN.json / en-US（品牌字段与描述文案「不支持上传」需改）
- 测试：settings_test.go、provider_test.go、repository_test.go、startup-config.test.tsx、schema-keys.structural.test.ts、s5-denominator-render.test.tsx、e2e shell.spec.ts / localization.spec.ts
