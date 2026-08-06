---
title: 文档体系说明（消费方精简入口）
status: active
created: 2026-07-18
updated: 2026-08-06
parent: null
version: 0.13.0
---

# docs/ · 文档体系（消费方）

本目录是 **Goal Governance** 在目标仓库中的核心规范入口：方法论、文档协议与模板。  
目标实例的状态真相只存在于各自 `docs/workspace-<NNN>-<slug>/` 根。  
仓库级愿景规则见 `docs/vision/alignment.md`（**不是**第二套目标状态）。

> **完整安装**：`docs/architecture/` 与 Skills **同级必备**。缺 architecture = 不完整安装。  
> monorepo 维护者长文与 dogfood 过程树**不**随本精简入口分发。

## 最小目录

```text
docs/
├── README.md                 # 本文件（消费方精简入口）
├── architecture/             # 治理原则与协议（必备）
│   ├── principles.md         # P-001～P-006
│   ├── workspace-protocol.md
│   ├── overview.md
│   └── directory-layout.md
├── templates/                # 五件套 + workspace-context + vision 冷启动模板
├── vision/
│   ├── alignment.md          # 愿景对齐契约（必备规则面）
│   ├── README.md             # 消费方愿景说明
│   ├── charter.md            # 冷启动后由 /vision 创建（实例）
│   └── plans/VP-*.md         # 冷启动后创建
└── workspace-<NNN>-<slug>/   # 工作区：goal-tree + GOAL-* 五件套
```

**不**随 Skills 包安装：`tech-stack.md`、monorepo dogfood 目标树、`web/`、`artifacts/`。

## 核心规则（摘要）

1. **工作区内目标平铺**：禁止用嵌套文件夹表达层级；层级只写在 `00-meta.md` 的 `parent`。
2. **GOAL-001 为 Root**：`parent: null`；编号单调不复用。
3. **五件套**：`00-meta` / `01-decision` / `02-execution` / `03-audit` / `attachments/`。
4. **总览同步**：新建/改状态/改 parent 后更新 `goal-tree.md`。
5. **P-001**：尚不可直接执行 → 先纲领路线图，再按阶段立项。
6. **P-002～P-004**：阶段质量意识；独立审计出意见、编排器响应；finding 三路径闭合；P-004 问用户。
7. **P-005**：可带未知立项；I-00N 与阶段门禁可追踪。
8. **P-006**：**单愿景**；Charter → VP → 工作区；对齐递归；Vision Review 使用 `reviews.md` 稳定索引 + `reviews/VRev-NNN-*.md` 平铺报告。

全文见 [architecture/principles.md](architecture/principles.md) 与 [architecture/workspace-protocol.md](architecture/workspace-protocol.md)。

## 与 Skills

| 层级 | 路径 | 说明 |
|------|------|------|
| 核心方法论 | 本 `docs/`（由 install 从包内 `core/docs` 安装） | 与 Skills 同级必备 |
| AI 适配器 | `skills/`（或改名包目录） | `/govern` · `/audit` · `/vision` · `/vision-audit` |
| 模板源（包内） | `skills/core/docs/templates/` | monorepo 编辑只改上游 `docs/templates/` 再 stage |

## 冷启动

1. 确认 `docs/architecture/principles.md` 存在。  
2. **`/vision`**：Charter → 首个 VP（+ Vision Review）。  
3. 建立工作区（`/govern` 或 install `--init-workspace`，slug **显式**）。  
4. **`/govern`** 创建 Root 五件套并推进。  
5. Goal 交叉审计用 **`/audit`**；独立 Vision Review 用 **`/vision-audit`**。

## 推荐阅读

1. [architecture/principles.md](architecture/principles.md)  
2. [architecture/workspace-protocol.md](architecture/workspace-protocol.md)  
3. [vision/alignment.md](vision/alignment.md)  
4. [templates/README.md](templates/README.md)  
5. 仓库根 `AGENTS.md`（install 安装）
