---
title: 消费方愿景规则镜像说明
status: active
created: 2026-07-29
updated: 2026-07-29
parent: null
version: 0.1.0
---

# docs/vision · 消费方最小规则面

本目录随 **Skills core** 安装到消费仓，提供愿景层**规则权威**，**不是** dogfood 过程树。

| 文件 | 角色 |
|------|------|
| `alignment.md` | 愿景对齐契约与门禁（P-006 操作细则权威） |
| （本 README） | 说明边界 |

## 不随 core 预装的实例

- 现行 `charter.md`、具体 `plans/VP-*.md`、`reviews.md` 过程条目  
- monorepo dogfood 工作区绑定与 progress 叙述  

冷启动时从 `docs/templates/vision/` 复制并改写 Charter / VP，再按 alignment 校验 `plan_refs` / `primary_plan`。

## 完整安装

缺 active Charter 或本目录 `alignment.md` → **不完整安装**（仅允许引导补齐）。  
原则全文仍在 `docs/architecture/principles.md`（P-001～**P-006**）。
