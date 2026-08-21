---
id: GOAL-022-w15-rectification-batch-b
title: W15 整改批 B · API 规范（F03 时间格式 / F11 GET 只读 / F10 429 配额 / F12 分页默认）
status: done
parent: GOAL-020-w15-user-perspective-findings
created: 2026-08-17
updated: 2026-08-17
version: 0.2.0
progress: 4/4
---

# GOAL-022 · W15 整改批 B

[GOAL-020](../GOAL-020-w15-user-perspective-findings/00-meta.md) 下级。D-002 批 B。

## 成功标准

- [x] **S1 · 冻结**：F03 只统一 RFC3339 不改字段名；F11 GET 404；F10 Retry-After + 配额 413；F12 DefaultPageSize=20
- [x] **S2 · 实施**：代码与测试
- [x] **S3 · 回归**
- [x] **S4 · 关门**

progress: **4/4**
