---
id: GOAL-010-w9-branding-asset-upload
title: W9 · 品牌图标上传（专用资产存储 + 自动图像处理）
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
progress: 6/6
---

# GOAL-010 · W9 · 品牌图标上传（专用资产存储 + 自动图像处理）

## 概述

VP-010 / workspace-010 的**第九波**（用户 2026-08-15 裁决立项）：设置页【品牌】四个图标字段（logoUrl / logoUrlLight / logoUrlDark / faviconUrl）从「填写 URL 链接」的**历史遗留方式**改为**上传控件**：

1. **专用资产存储**：不走常规文件库（admin.file-library），也不走通用 /api/upload（owner 隔离、非公开）；新增独立品牌资产仓（`<data>/brand-assets`）+ 鉴权上传端点 + 公开读端点。
2. **自动图像处理**：服务端重编码（PNG/JPEG/WebP → PNG/JPEG，限幅 logo ≤512px / favicon 64px，JPEG 质量 82，单文件 ≤4 MiB），永不回传原始字节。
3. **全量使用点更新**：settings schema、i18n、相关测试；外壳/登录页/favicon 消费逻辑不变（URL 形态不变，旧 URL 值仅兼容读取）。

## 当前边界

- 范围：apps/api settings 模块（schema / 端点 / 存储 / 处理 / 配置）、apps/web 品牌消费与测试、e2e、i18n 文案。
- 非范围：不改 /api/branding 契约形状与前端消费逻辑；不迁移旧 URL 值；不支持 SVG 上传；不引入文件库模块语义。

## 成功标准与路线图（P-001）

- [ ] **S1 · 方案冻结**：六项用户裁决 + 资产存储/端点/处理/清理/配置设计落盘（D-001），信息表 closed
- [ ] **S2 · 专用资产存储与端点**：brand-assets 目录、`POST /api/branding/assets`（settings.write）、`GET /api/branding/assets/{id}`（公开、nosniff/缓存头）
- [ ] **S3 · 自动图像处理**：PNG/JPEG/WebP 解码（新增 golang.org/x/image）、限幅重编码、参数进 config.yaml（env 可覆盖）
- [ ] **S4 · 前端与消费点**：settings schema 品牌字段改 upload 控件（+ 单字段移除）、文案/i18n、旧 URL 兼容读取、相关测试更新
- [x] **S5 · 验证与回归**：单测（处理/端点/安全）+ 全量回归 + 活栈点验（E-003）
- [x] **S6 · go 判定 + cross 审计 + 关门**：go 不 held；A-001 self + A-002 independent（grok-4.6 high）pass；E-004 响应全 closed；goal-tree 同步

progress: 0/6 由六个等权检查点派生（S1～S6 全勾后 6/6）。

## 审计策略

本目标含 **security 面**（上传校验、公开静态服务、图像处理、XSS 防护）→ 按 P-003 采用 **cross** 审计模式（S6：self + 至少一个指定 provider 的 independent；provider 历史惯例 grok，执行时确认）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | URL 输入彻底移除：表单仅上传 + 单字段移除；旧 URL 值仅兼容读取不迁移 | 方案 | S1 | 用户 2026-08-15 书面确认 | **closed** |
| I-002 | required | 处理参数：logo ≤512px、favicon 64px、JPEG 质量 82、单文件 ≤4 MiB；WebP 输入新增 golang.org/x/image | 方案 | S1 | 用户 2026-08-15 书面确认 | **closed** |
| I-003 | required | SVG 继续不支持（安全：active content 拒收） | 方案 | S1 | 用户 2026-08-15 书面确认 | **closed** |
| I-004 | required | 清理语义：替换/重置即删旧资产 + 启动 GC 孤儿文件 | 方案 | S1 | 用户 2026-08-15 书面确认 | **closed** |
| I-005 | required | 处理参数配置位置：config.yaml（W7 体系，env 可覆盖） | 方案 | S1 | 用户 2026-08-15 书面确认 | **closed** |
| I-006 | required | 审计模式：cross（self + independent） | 关门 | S6 | 用户 2026-08-15 书面确认；provider 执行时确认（历史惯例 grok） | **closed**（provider 待 S6 确认） |
| I-007 | required | 消费点完整清单（/api/branding、shell、login、favicon、结构测试、s5 分母、e2e） | 实施 | S4 | 创建时代码扫描（D-001 附录 A） | **closed** |
| I-008 | required | UploadField 是否支持「移除已设图片」交互（或需扩展） | 实施 | S4 | 检查 form-controls.tsx UploadField 能力；不支持则扩展 | **closed**（扩展 UploadField：单文件上传值非空时渲染「移除图片」按钮，E-002） |
| I-009 | non-blocking | 公开 GET 缓存与安全头细节（Cache-Control / nosniff / Content-Type 存储） | 实施 | S2 | D-001 设计内定（沿用上传仓安全基线） | **closed**（nosniff + CSP sandbox + immutable 缓存，E-002/E-003 实测） |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。