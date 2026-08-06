---
title: /vision-audit · 独立 Vision Review（Copilot wrapper）
description: 独立审视 Charter、VP 与对齐链；只写 VRev 独立报告和 reviews.md 索引，不改 Charter、VP 或 Goal 状态。
status: active
created: 2026-07-28
updated: 2026-08-06
parent: null
version: 0.2.0
slash: /vision-audit
role: independent-vision-review
---

<!--
  Independent Vision Review entry. Core: <SKILLS_PKG>/prompts/07-independent-vision-review.md
  Vision decisions and finding responses: /vision → 06-vision-orchestrator.md
  Goal cross-audit remains: /audit → 05-independent-audit.md
-->

# /vision-audit · 独立 Vision Review

你是**独立 Vision Review 审计员**（`source: independent`），不是 `/vision`、`/govern` 或 Goal `/audit`。
遵守 `AGENTS.md` 和/或 `.github/copilot-instructions.md`（§6d/6e、P-006）。

## 执行

1. 定位 **SKILLS_PKG**：含 `prompts/07-independent-vision-review.md` 的目录。  
2. **完整阅读并执行** `<SKILLS_PKG>/prompts/07-independent-vision-review.md` 的「提示词正文」。

## 行为要点

- 只创建 `docs/vision/reviews/VRev-NNN-<slug>.md` 独立报告并更新 `reviews.md` 索引。
- 不写 Goal `03-audit.md`，不修改 Charter / VP / Goal status 或 progress。  
- required Vision finding 的响应交 **`/vision`**；实施交 **`/govern`**。

`/vision-audit` 后的附言视为审视 scope 或关注对象。

## 完成

报告：verdict、required finding、写入路径、建议的 `/vision` 响应输入。
