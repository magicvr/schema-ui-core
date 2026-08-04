---
title: Skills 模板指针（非真相源）
status: active
created: 2026-07-19
updated: 2026-08-04
parent: null
version: 0.7.0
---

# Skills 模板指针

本目录**不是**模板真相源（GOAL-022）。

| 角色 | 路径 |
|------|------|
| **canonical** | 仓库 `docs/templates/` |
| **包内分发镜像** | `skills/core/docs/templates/`（由 `python scripts/stage_skills_mirrors.py` 从 docs stage） |

- 新建目标、工作区、愿景模板请读 **`core/docs/templates/`**（或安装后的 `docs/templates/`）。
- **禁止**在本目录手维 `goal-folder/`、`ledger-entry/` 等第三份副本。
- 改 canonical 后：运行 stage 脚本并提交 `skills/core` / `skills/contracts` 镜像 diff。
