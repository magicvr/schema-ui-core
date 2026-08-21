---
id: D-002-bucket-model
title: 桶模型裁定：单桶 + 命名空间前缀（I-002 闭合）
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-001-object-storage
version: 0.1.0
---

# D-002 · 桶模型：单桶 + 命名空间前缀（闭合 I-002）

## 决定

对象存储采用**单桶 + 命名空间前缀**模型：

- S3 兼容实现：一个显式配置的桶；key 形如 `<namespace>/<object-id>`，namespace ∈ { `avatars`, `brand-assets`, `uploads` }。
- 本地盘实现：一个存储根目录（缺省仍由 `filepath.Dir(db.path)` 派生）；每个 namespace 对应根下同名子目录。**与现状磁盘布局完全一致**——现行 `${DB_DIR}/avatars`、`${DB_DIR}/brand-assets`、`${DB_DIR}/uploads` 三个目录即三个命名空间，零迁移。

## 理由（证据）

1. 现行三类落盘就是"一个本地根 + 三个子目录"（composition.go:308/318/337，均从 `filepath.Dir(cfg.DBPath)` 派生）。单桶+前缀使本地适配器与 S3 适配器的 key/路径映射一一对应，端口语义天然平等（VP-014 意图 3）。
2. VP-014 明文用"命名空间（或等价隔离）分开头像、品牌资源、通用上传"；多桶会把隔离做在桶级，凭证与建桶负担翻三倍，且与本地子目录模型不再同构。
3. 三类落盘的访问控制差异（公开读 brand-assets vs owner-only uploads）由 HTTP 面与调用方门禁承担（raster_assets.go 头注、upload.go owner-only GET），不依赖桶级策略；单桶不削弱该模型。
4. 禁止第三对象存储方言（VP I-014-001）不受影响：桶模型与驱动选择正交。

## 未选方案

- **多桶（每类一桶）**：隔离更强但配置面膨胀（3×endpoint/bucket/凭证组合或桶命名约定），本地盘无对应物，破坏"两实现合同平等"；且现有安全模型不在桶级。
- **按模块/业务分桶**：违反薄内核边界（VP-003）；对象端口只认三类第一方命名空间，不认业务域。

## 影响

- R1 端口冻结以 namespace 枚举为合同的一部分（见 workspace-014 GOAL-002 D-001）。
- 对象 id 规则 `^[0-9a-f]{32}$` 在端口层 fail-closed 校验（与现行 uploadFileIDPattern 一致），杜绝路径/key 注入。
