---
title: 目录布局
status: active
created: 2026-07-18
updated: 2026-07-31
parent: null
version: 0.6.3
---

# 目录布局

```text
goal-governance/
├── AGENTS.md                 # AI 助手强制规则
├── README.md                 # 项目入口说明
├── docs/
│   ├── README.md             # 文档体系规范
│   ├── vision/               # 仓库级愿景体系（非 goal-tree）
│   │   ├── README.md
│   │   ├── charter.md        # 现行愿景（不可 done）
│   │   ├── roadmap.md        # VP 索引
│   │   ├── plans/VP-*.md     # 愿景规划（可关门）
│   │   ├── revisions.md
│   │   ├── workspaces.md
│   │   ├── alignment.md
│   │   └── consumer-checklist.md
│   ├── workspace-001-example/ # 显式工作区根
│   │   ├── workspace.md       # Root/范围/资料/规划对齐
│   │   ├── goal-tree.md       # 本工作区目标树与状态总览
│   │   ├── GOAL-001-.../      # 目标（平铺，无嵌套）
│   │   └── GOAL-00N-.../
│   ├── shared-materials/      # 工作区外的资料候选库存
│   ├── templates/
│   │   ├── README.md            # 核心模板层说明
│   │   ├── goal-folder/         # canonical 五件套模板
│   │   ├── vision/              # Charter / VP 冷启动模板
│   │   └── workspace-context.md # workspace-<NNN>-<slug>/workspace.md 模板
│   ├── contracts/               # canonical 机读协议/模板版本与兼容声明
│   │   ├── skills-consumer-contract.schema.json
│   │   └── skills-consumer-contract.json
│   ├── architecture/
│   │   ├── overview.md
│   │   ├── principles.md     # 治理原则（元规则）
│   │   ├── workspace-protocol.md
│   │   ├── tech-stack.md
│   │   └── directory-layout.md
│   └── _index/               # 预留
├── skills/                    # AI/Agent 消费适配器与分发包
│   ├── prompts/
│   ├── templates/             # docs/templates 的同步镜像
│   │   └── workspace-context.md
│   ├── contracts/             # docs/contracts 的同步镜像
│   └── install.*
└── web/
    ├── main.py
    ├── requirements.txt
    ├── README.md
    ├── static/
    └── templates/
```

## 约束

- `docs/workspace-<NNN>-<slug>/GOAL-*` 之间**不得**再嵌套目标目录。
- 新目标只新增当前工作区根内的同级文件夹，并改 `parent` + 该工作区 `goal-tree.md`。
- `docs/templates/goal-folder/` 是核心 canonical 模板；包内分发镜像为 `skills/core/docs/templates/goal-folder/`（由 `scripts/stage_skills_mirrors.py` 从 docs stage；**不**再手维 `skills/templates/` 第三份）。
- `docs/workspace-<NNN>-<slug>/workspace.md` 是显式工作区上下文，绑定一个 Root Goal 与该工作区根范围；`docs/templates/workspace-context.md` 与 core 镜像必须经 stage 一致。没有显式工作区根时：**仅当**存在 `docs/goals/` 才按 legacy 隐式单工作区处理；否则不得猜测工作区根。
- `docs/vision/` 是仓库级愿景与规划对齐层；**不是**目标状态库，不得写入 progress% 或替代各区 goal-tree。Primary 冲突与 VP 空转规则见 `docs/vision/alignment.md`。
- 共享资料只以版本/哈希固定引用出现在工作区上下文或受控记录中，不能成为跨工作区目标状态或第二真相源。
- `GOAL-*` id 仅工作区内唯一，**形状不嵌工作区编号**；跨区引用见 [workspace-protocol.md](workspace-protocol.md) §2.6（文档默认 **Q2** 路径，对话默认 **Q3** 标签）。
- `docs/contracts/` 是消费适配器版本与兼容声明的 canonical；`skills/contracts/` 由 stage 从 docs 生成，必须逐字节一致且不得另立版本真相。
- **stage 门禁（本 monorepo）**：改 `docs/architecture` 白名单、`docs/templates/**`、`docs/vision/alignment.md` 或 `docs/contracts/**` 后，必须运行 `python scripts/stage_skills_mirrors.py` 并**提交**生成的 `skills/core` / `skills/contracts` 变更；禁止手改镜像、禁止只交 docs。AI 操作入口见根 `AGENTS.md` §8c；说明见 [../README.md](../README.md)。
