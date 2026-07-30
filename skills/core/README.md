---
title: Skills 包内 Core 方法论镜像
status: active
created: 2026-07-24
updated: 2026-07-30
parent: null
version: 0.3.0
---

# skills/core · 消费方核心方法论镜像

本目录是 **GOAL-019 D-003 / D-004** 与 **GOAL-001 D-025（A-018 回流）** 规定的、随 Skills 包分发的**核心方法论子集**。

| 包内路径 | `install` 默认落到消费仓 |
|----------|---------------------------|
| `docs/README.md` | `docs/README.md`（精简入口） |
| `docs/architecture/*.md` | `docs/architecture/` |
| `docs/templates/**` | `docs/templates/` |
| `docs/vision/alignment.md`（+ 本目录 README） | `docs/vision/`（规则权威；**非** dogfood 实例） |

## 包含

- `architecture/principles.md` — P-001～**P-006**
- `architecture/workspace-protocol.md` — 工作区与共享资料协议
- `architecture/overview.md` — 逻辑架构（消费方；无 monorepo dogfood）
- `architecture/directory-layout.md` — 消费方最小目录树
- `templates/` — 五件套 + `workspace-context.md` + `vision/` 冷启动模板
- `vision/alignment.md` — 愿景对齐契约与门禁（完整安装必备规则面）
- `docs/README.md` — 精简文档入口

## 不包含

- `tech-stack.md`（实现栈）
- monorepo dogfood 目标树、现行 Charter/VP **实例**、`web/`、`artifacts/`
- 维护者-only 的 standalone 测试与 releases 长文

Charter / VP **实例**由冷启动从 `templates/vision/` 创建；规则以 `vision/alignment.md` + principles P-006 为准。

## 与 monorepo canonical（GOAL-022）

上游规范在仓库 `docs/architecture/`、`docs/templates/`、`docs/vision/alignment.md`。  
**日常只改 `docs/`**，再运行：

```bash
python scripts/stage_skills_mirrors.py
```

本目录中 architecture / templates / `vision/alignment.md` 为 stage 产物（可提交 git，禁止手改后不 stage）。  
**手维例外**：`docs/README.md`（消费方精简入口）、`vision/README.md`、本文件。

**缺 core = 不完整安装**（与 Skills 同级必备）。缺 `docs/vision/alignment.md` 或 active Charter = 愿景层不完整（仅引导补齐）。
