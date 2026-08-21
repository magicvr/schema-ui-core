---
id: GOAL-029-w18-preview-copy-and-import-modal
title: W18 · 预览弹窗/复制链接与导入模态模板
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-18
updated: 2026-08-18
version: 0.2.0
progress: 4/4
---

# GOAL-029 · W18 · 预览弹窗/复制链接与导入模态模板

VP-010 / workspace-010 的**第十八波**：承接 [GOAL-024](../GOAL-024-w16-user-perspective-improvements/00-meta.md) A-007 F-001 / F-002（recommended open）。不重开 GOAL-024。

## 当前边界

- **范围**：文件预览改为手势内同步开窗再灌入 blob，并 `revokeObjectURL`；复制链接改为源站绝对 URL（非 `blob:`）；导入 CSV 模板入口移入导入模态；模板下载失败可见；导入 200 `fieldErrors` 定向测试。
- **非范围**：图片 Lightbox；带签名的对外可分享下载链；改下载端点鉴权模型；A-007 已闭合项。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：D-001。
- [x] **S2 · 实施**：预览/复制 + 模态模板 + 失败提示（E-002）。
- [x] **S3 · 定向验证**：Web 73/73 + `tsc`（E-003）。
- [x] **S4 · 自审与关门**：A-001 self pass；goal-tree / workspace 同步（E-004）。

progress: 四个等权检查点；当前 **4/4**。

## 审计策略

S4 关门 `self`（可逆 UX，无安全门禁语义变化）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required | 复制链接在无签名 URL 时复制什么 | S1 | S1 | 对照下载端点仍 Bearer 门禁 | **verified** | D-001：源站绝对 download 路径；不宣称对外免登 |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 溯源

- GOAL-024 A-007 F-001 / F-002（经本波 A-001 / A-010 记 **fixed**）
- GOAL-024 D-002 W16-F02 / W16-F03 原冻结（Lightbox / 对外直链本波不做）
