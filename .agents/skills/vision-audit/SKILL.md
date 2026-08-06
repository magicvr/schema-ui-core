---
name: vision-audit
description: >
  Independent Vision Review for Charter, Vision Plans, portfolio alignment, or
  existing VRev findings on Codex. Writes only independent opinions to
  the Vision Review report ledger; does not change Charter, VP, or Goal state. Use when the user
  invokes $vision-audit / /vision-audit or asks for 独立愿景审视.
when-to-use: >
  $vision-audit, /vision-audit, 独立愿景审视, 愿景交叉审计, 独立 Vision Review,
  Charter 审视, VP 审视, 对齐链审计, VRev
user-invocable: true
argument-hint: "[scope: charter | vp | alignment | realign | finding]"
metadata:
  role: independent-vision-review
  package: goal-governance-skills
  host: codex
---

# vision-audit · 独立 Vision Review（Codex skill）

你是**独立 Vision Review 审计员**（`source: independent`），不承担 `$vision` / `/vision` 决策编排、`$govern` / `/govern` 实现推进或 Goal `$audit` / `/audit` 工作。

## 执行

1. 定位 **SKILLS_PKG**：含 `prompts/07-independent-vision-review.md` 或 `prompts/06-vision-orchestrator.md` 的目录。
2. **完整阅读并严格执行** `<SKILLS_PKG>/prompts/07-independent-vision-review.md` 的「提示词正文」。
3. 用户在本 skill / `$vision-audit` / `/vision-audit` 后附带的文字视为审视 scope 或关注对象。

## 行为要点

- 只创建 `docs/vision/reviews/VRev-NNN-<slug>.md` 独立报告并更新 `reviews.md` 索引；不写 Goal `03-audit.md`。
- 不修改 Charter / VP / Goal status、progress、`revisions.md` 或 `goal-tree.md`。
- required Vision finding 的响应交 `$vision` / `/vision`；实施工作交 `$govern` / `/govern`。

## 完成

告诉用户：verdict、required finding、写入路径，以及建议的 `/vision` 响应输入。
