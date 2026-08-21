---
id: GOAL-003-s2-s3-renderer-and-shell
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

## D-001 · GOAL-003 成功标准对齐 D-004（重写确认）

### 已采纳

1. C1/C2 分母以 Root D-004 为准（见 `00-meta.md` 重写检查点）；作废 A-001 过窄 pass。
2. 实施不拆第二工作区；在本目标内一次交付 S2+S3 视觉 fidelity 重做（E-002）。
3. 不把 Stitch `code.html` 接入生产；生产仍为 React + Schema + Token + primitives。
4. 验收证据：代码结构标记 + 真实模块 vitest + build + Playwright；不强制像素对照 gitignored PNG。

### 为什么

A-006/A-002 证伪过窄完成声明。分母必须与 Root S2/S3 与 VP-005 exit 2–3 可复核对齐。

### 未选

| 方案 | 原因 |
|------|------|
| 仅补文档勾选 | 违反过程诚实 |
| 拆 GOAL-006/007 平行树 | 本目标已足够承载 C1/C2；避免空壳目标 |
