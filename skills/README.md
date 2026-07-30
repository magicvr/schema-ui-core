---
title: Skills · 目标治理可复用包
status: active
created: 2026-07-18
updated: 2026-07-31
parent: null
version: 1.8.0
---

# Skills

本目录提供可复制到**其他项目**的目标治理约定与模板。  
本仓库运行中的强制规则仍以根目录 [AGENTS.md](../AGENTS.md) 为准；此处是提炼后的**可复用交付物**。

Skills 是核心方法论的 **AI 消费适配器**。**核心方法论与 Skills 同级必备**（GOAL-019 D-003）：包内 [`core/`](core/) 为消费分发镜像，`install` **默认**安装到目标仓 `docs/architecture/`、`docs/templates/` 与精简 `docs/README.md`。缺 core = **不完整安装**。

在 monorepo 中，规范模板位于 [`docs/templates/`](../docs/templates/)；包内 **`core/docs/templates/`** 为 stage 生成的分发镜像（GOAL-022）。`skills/templates/` 仅保留指针 README，**不是**第三真相。机读契约以 [`docs/contracts/`](../docs/contracts/) 为 canonical，本包 `contracts/` 由 stage 逐字节生成。改 docs 后运行：`python scripts/stage_skills_mirrors.py`。

**发布与候选证据边界**：

| 身份 | 状态 |
|------|------|
| **`v0.9.0`** / **`v0.9.1`** / **`v0.9.2`** / **`v0.10.0`** | 已发布 annotated tag / Release 基线。 |
| **`v0.11.0`（本冻结）** | 矩阵 **`candidateRevision: v0.11.0`**；GOAL-002 Codex 安装面（`--codex` / `-All`）+ 文档 pin；四入口 × 三宿主 **runtime-verified**（2026-07-30 证据）。**不**宣称 Codex 矩阵 `committed` / `runtime-verified`。正式 GitHub Release 以 annotated tag + release evidence + Environment `release` 为准。 |

Claude Code / Grok Build / Copilot CLI 为 `committed` + `runtime-verified`；Web parser 为 `automated-verified`。权威字段见 [`docs/contracts/skills-consumer-contract.json`](../docs/contracts/skills-consumer-contract.json) 与 [`docs/contracts/skills-consumer-compatibility-matrix.json`](../docs/contracts/skills-consumer-compatibility-matrix.json)。

## 产品模型（必读）

| 层级 | 是什么 | 用户怎么用 |
|------|--------|------------|
| **核心方法论** | `docs/architecture` + `docs/templates` + 精简 `docs/README` | install 从 `core/` **默认**安装；与 Skills **同级必备** |
| **实现主入口（primary）** | 编排器：扫描 / 意见台账 / 分类 / P-004 裁决 / 确认 / 原语 | **`/govern`** |
| **决策入口** | 愿景与组合：Charter / VP / Review / re-align / 结构选型 | **`/vision`** |
| **Goal 交叉入口** | 独立 Goal 审计：只出意见（`source: independent`） | **`/audit`** |
| **愿景交叉入口** | 独立 Vision Review：只写 `reviews.md`（`source: independent`） | **`/vision-audit`** |
| **原语（primitives）** | 创建目标、记决策、更执行、写审计 | 由编排器调用；Copilot advanced 可选 |
| **规则** | AGENTS / copilot-instructions | 结构、编号、操作细则摘要 |

生命周期（P-006）：**愿景/意图 → 工作区+Root → 纲领路线图 → 阶段计划 → 子目标 → 审计/整改 → 关门**。  
决策层与 Vision finding 响应由 `/vision` 负责；Goal 交叉意见由 `/audit` 写入，独立 Vision Review 由 `/vision-audit` 写入；Goal 的**响应与放行**由 `/govern` 处理。

工作区协议：`/govern` 和 `/audit` 先定位当前 `docs/workspace-<NNN>-<slug>/workspace.md`，校验其 Root Goal、canonical 范围和共享资料固定引用；不匹配或多个工作区未指定焦点时 fail closed。冷启动缺 Charter 时先 **`/vision`**。没有显式工作区根的旧项目才按 `docs/goals/` 的 legacy 隐式单工作区工作。

| 工具 / 表面 | 安装位置 | 斜杠 | 当前契约层级 |
|------|----------|------|--------------|
| Claude Code CLI `2.1.220` | `.claude/skills/{govern,audit,vision,vision-audit}/` | `/govern` · `/audit` · `/vision` · `/vision-audit` | govern/audit/vision/vision-audit **`runtime-verified (2026-07-30)`** |
| Grok Build CLI `0.2.114` | `.grok/skills/{govern,audit,vision,vision-audit}/` | `/govern` · `/audit` · `/vision` · `/vision-audit` | govern/audit/vision/vision-audit **`runtime-verified (2026-07-30)`** |
| GitHub Copilot CLI `1.0.75` | `.github/…` + prompts | `/govern` · `/audit` · `/vision` · `/vision-audit` | 四个入口均 `runtime-verified` via BYOK（2026-07-30） |
| OpenAI Codex | `.agents/skills/{govern,audit,vision,vision-audit}/` | `$govern` · `$audit` · `$vision` · `$vision-audit` | **install surface shipped**（GOAL-002）；runtime 探针见目标证据链（非矩阵 committed） |

核心行为：

> Contract manifest 的 `verificationStatus` 仍是历史有界事实，不能替代候选矩阵。Claude Code 机读证据保存脱敏 stream transcript；Grok Build 机读证据保留辅助 session-title `grok-build` alias 的 502 警告，但主 `grok-4.5` 调用 exit `0` 且输出实际 dispatch marker。完整发行验收仍以全部 matrix 单元、coverage、CI 与 release 证据为准。

- 编排：[`prompts/00-govern-orchestrator.md`](prompts/00-govern-orchestrator.md)
- 交叉：[`prompts/05-independent-audit.md`](prompts/05-independent-audit.md)
- 愿景：[`prompts/06-vision-orchestrator.md`](prompts/06-vision-orchestrator.md)
- 独立愿景审视：[`prompts/07-independent-vision-review.md`](prompts/07-independent-vision-review.md)

## 目录结构

```text
skills/
├── README.md
├── AGENTS.template.md
├── install.sh / install.ps1
├── core/                               # GOAL-019：方法论镜像 → install 默认装到 ./docs/
│   └── docs/
│       ├── README.md                   # 精简文档入口
│       ├── architecture/               # principles, workspace-protocol, overview, layout
│       └── templates/                  # 五件套 + workspace-context
├── install/
│   ├── claude/
│   │   ├── AGENTS.md
│   │   └── skills/{govern,audit,vision,vision-audit}/SKILL.md
│   ├── grok/
│   │   └── skills/{govern,audit,vision,vision-audit}/SKILL.md
│   ├── codex/
│   │   └── skills/{govern,audit,vision,vision-audit}/SKILL.md  # → .agents/skills/
│   └── copilot/
│       ├── copilot-instructions.md
│       └── prompts/
│           ├── govern.md               # impl primary
│           ├── audit.md                # Goal cross-audit (default)
│           ├── vision.md               # decision layer (default)
│           ├── vision-audit.md         # independent Vision Review (default)
│           └── new-goal.md …           # advanced only
├── prompts/
│   ├── 00-govern-orchestrator.md       # PRIMARY impl core
│   ├── 01–04 …                         # primitives
│   ├── 05-independent-audit.md         # Goal cross-audit core
│   ├── 06-vision-orchestrator.md       # vision decision core
│   └── 07-independent-vision-review.md # independent Vision Review core
├── templates/README.md                 # 指针：模板在 core/docs/templates（GOAL-022）
├── contracts/                          # docs/contracts 的 stage 镜像
└── tests/
```

## 安装（双入口 · GOAL-023）

| 入口 | 是什么 | 网络 |
|------|--------|------|
| **1 · Bootstrap** | `scripts/bootstrap/install-online.ps1` / `.sh`（Release 亦挂同名脚本） | 在线：下 **skills zip**；或离线：本地 zip + `.sha256` |
| **2 · 包内 install** | 解压后的 `skills/install.ps1` / `install.sh` | **离线**（不访问网络） |

**skills zip 内嵌 core**（GOAL-019）：默认安装路径**不**再从网上拉 core。  
并行 **core-only** 资产 `goal-governance-core-vX.Y.Z.zip` 供 standalone / 无 Skills 场景；**不是**本表默认 Skills 路径（见 [standalone-bootstrap](../docs/standalone-bootstrap.md)、`pack_core_release.py`）。

### 入口 1 · Bootstrap（推荐 · 其他项目）

**当前示例 pin 最新正式 tag `v0.11.0`**（每次发版更新本节与根 README；固定 tag URL，**禁止** `main`/branch raw 作权威入口；**不是**无 pin 的 always-latest 安装）。也可 clone monorepo 用 `scripts/bootstrap/`。

在目标项目根执行（先落盘 bootstrap，再跑；默认**不**管道直跑）：

```powershell
Invoke-WebRequest -Uri "https://github.com/magicvr/goal-governance/releases/download/v0.11.0/install-online.ps1" `
  -OutFile .\install-online.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\install-online.ps1 -Version 0.11.0 -Force

# 离线（本地 skills zip + .sha256）：
powershell -NoProfile -ExecutionPolicy Bypass -File .\install-online.ps1 `
  -Version 0.11.0 -ZipPath .\goal-governance-skills-v0.11.0.zip -Force
```

```bash
curl -fsSL -o install-online.sh \
  "https://github.com/magicvr/goal-governance/releases/download/v0.11.0/install-online.sh"
chmod +x install-online.sh
bash ./install-online.sh --version 0.11.0 --force

# 离线：
bash ./install-online.sh --version 0.11.0 --zip-path ./goal-governance-skills-v0.11.0.zip --force
```

Bootstrap 会：校验 SHA-256 → 落到 `./skills` → 调用包内 install **默认 `-All` / `--all`**（四入口 + core → `docs/`）。digest 失败 **fail closed**。详见 [scripts/bootstrap/README.md](../scripts/bootstrap/README.md)。

### 入口 2 · 包内 install（解压后）

1. 打开本仓库 [Releases](https://github.com/magicvr/goal-governance/releases)，下载与 tag 对应的  
   `goal-governance-skills-vX.Y.Z.zip`（可对照同目录的 `.sha256` 校验）。  
   包内含 Skills 适配器 + **core 方法论镜像**（prompts、install、`core/`、模板/契约），**不含** monorepo dogfood 过程树、`web/` 或 `artifacts/`，**不含** `tech-stack.md`。  
   若目标 tag 的 Release 尚未出现，可从源码树复制 `skills/`，或使用维护者提供的预打包 zip（**不**等于正式 Release 身份）。
2. 在目标项目根目录解压，使包内容落在 `./skills/`（或你选择的目录名）：

```bash
# 示例：已下载 zip 到当前目录
unzip goal-governance-skills-vX.Y.Z.zip
# zip 根目录名为 goal-governance-skills-vX.Y.Z/ — 重命名为 skills 便于默认参数
mv goal-governance-skills-vX.Y.Z skills
```

```powershell
Expand-Archive .\goal-governance-skills-vX.Y.Z.zip -DestinationPath .
Rename-Item .\goal-governance-skills-vX.Y.Z skills
```

3. 安装宿主入口（默认 **`/govern` + `/audit` + `/vision` + `/vision-audit`**）**并默认安装 core → `./docs/`**：

```bash
bash ./skills/install.sh --all --skills-dir ./skills
# 或单宿主：--claude / --grok / --copilot / --codex（同样会装 core）
```

```powershell
.\skills\install.ps1 -All -SkillsDir .\skills
# 或：-Claude / -Grok / -Copilot / -Codex（同样会装 core）
```

4. 确认 `docs/architecture/principles.md` 等已存在；冷启动先 **`/vision`**（Charter→VP），再建立工作区并 **`/govern`**；Goal 交叉审计用 **`/audit`**；独立 Vision Review 用 **`/vision-audit`**。

> 维护者正式发布：推 **annotated** `v*` tag → CI pack → Environment **`release` 审批** → 硬 `release_evidence --mode release` 通过后自动 `gh release create` 并挂 **skills zip + core zip + bootstrap 脚本** / sha256 / evidence。详见 [docs/releases/README.md](../docs/releases/README.md)。  
> 本地调试 zip：  
> `python scripts/pack_skills_release.py --version X.Y.Z --output-dir dist/`  
> `python scripts/pack_core_release.py --version X.Y.Z --output-dir dist/`  
> 尚未对齐矩阵/`candidateRevision` 的工作树**不要**推正式 tag；门禁失败则**不会**创建 Release。
### 0. 从源码树复制包（开发者 / 无 Release 时）

```bash
cp -R /path/to/goal-governance/skills ./skills
```

```powershell
Copy-Item -Recurse path\to\goal-governance\skills .\skills
```

### 1. 手动安装

**默认安装面**（与脚本一致）：每个所列安装产物都装 **`/govern` + `/audit` + `/vision` + `/vision-audit`**（Codex 侧为 `$govern` 等）。填表类 advanced slash 仍为可选。四入口在 Claude / Grok / Copilot CLI 上均为 **`runtime-verified`（2026-07-28）**：`/vision-audit` 证据为只读 dispatch（wrapper 路由、核心提示词加载、愿景发现），**不是**写盘全路径 e2e；见兼容矩阵。Codex 已提供 install 面（GOAL-002），**尚未**写入矩阵 `committed` / runtime-verified。

#### Claude Code

```text
install/claude/AGENTS.md
  →  <repo>/AGENTS.md
install/claude/skills/{govern,audit,vision,vision-audit}/SKILL.md
  →  <repo>/.claude/skills/{govern,audit,vision,vision-audit}/SKILL.md
```

```bash
mkdir -p .claude/skills/govern .claude/skills/audit .claude/skills/vision .claude/skills/vision-audit
cp ./skills/install/claude/AGENTS.md ./AGENTS.md
cp ./skills/install/claude/skills/govern/SKILL.md .claude/skills/govern/SKILL.md
cp ./skills/install/claude/skills/audit/SKILL.md .claude/skills/audit/SKILL.md
cp ./skills/install/claude/skills/vision/SKILL.md .claude/skills/vision/SKILL.md
cp ./skills/install/claude/skills/vision-audit/SKILL.md .claude/skills/vision-audit/SKILL.md
```

#### Grok Build

```text
install/grok/skills/{govern,audit,vision,vision-audit}/SKILL.md
  →  <repo>/.grok/skills/{govern,audit,vision,vision-audit}/SKILL.md
```

（建议同时有根 `AGENTS.md` 作项目规则；可与 Claude 共用。）

```bash
mkdir -p .grok/skills/govern .grok/skills/audit .grok/skills/vision .grok/skills/vision-audit
cp ./skills/install/grok/skills/govern/SKILL.md .grok/skills/govern/SKILL.md
cp ./skills/install/grok/skills/audit/SKILL.md .grok/skills/audit/SKILL.md
cp ./skills/install/grok/skills/vision/SKILL.md .grok/skills/vision/SKILL.md
cp ./skills/install/grok/skills/vision-audit/SKILL.md .grok/skills/vision-audit/SKILL.md
```

#### OpenAI Codex

```text
install/claude/AGENTS.md
  →  <repo>/AGENTS.md
install/codex/skills/{govern,audit,vision,vision-audit}/SKILL.md
  →  <repo>/.agents/skills/{govern,audit,vision,vision-audit}/SKILL.md
```

（官方 REPO skill 根为 `.agents/skills`；显式调用 `$govern` / `$audit` / `$vision` / `$vision-audit`。）

```bash
mkdir -p .agents/skills/govern .agents/skills/audit .agents/skills/vision .agents/skills/vision-audit
cp ./skills/install/claude/AGENTS.md ./AGENTS.md
cp ./skills/install/codex/skills/govern/SKILL.md .agents/skills/govern/SKILL.md
cp ./skills/install/codex/skills/audit/SKILL.md .agents/skills/audit/SKILL.md
cp ./skills/install/codex/skills/vision/SKILL.md .agents/skills/vision/SKILL.md
cp ./skills/install/codex/skills/vision-audit/SKILL.md .agents/skills/vision-audit/SKILL.md
```

#### GitHub Copilot

```text
install/copilot/copilot-instructions.md
  →  .github/copilot-instructions.md
install/copilot/prompts/govern.md
  →  .github/prompts/govern.prompt.md
install/copilot/prompts/audit.md
  →  .github/prompts/audit.prompt.md
install/copilot/prompts/vision.md
  →  .github/prompts/vision.prompt.md
install/copilot/prompts/vision-audit.md
  →  .github/prompts/vision-audit.prompt.md
```

| Wrapper | 斜杠 | 何时安装 |
|---------|------|----------|
| govern.md | `/govern` | **默认**（实现主入口） |
| audit.md | `/audit` | **默认**（交叉审计） |
| vision.md | `/vision` | **默认**（决策层） |
| vision-audit.md | `/vision-audit` | **默认**（独立 Vision Review） |
| new-goal … write-audit | advanced | 仅 `--with-primitives` |

### 2. 脚本安装

| 参数 | 作用 |
|------|------|
| `--claude` / `-Claude` | `AGENTS.md` + `.claude/skills/{govern,audit,vision,vision-audit}` + **core → docs/** |
| `--grok` / `-Grok` | `.grok/skills/{govern,audit,vision,vision-audit}` + **core → docs/** |
| `--copilot` / `-Copilot` | copilot-instructions + `govern`/`audit`/`vision`/`vision-audit` prompts + **core → docs/** |
| `--codex` / `-Codex` | `AGENTS.md` + `.agents/skills/{govern,audit,vision,vision-audit}` + **core → docs/** |
| `--with-primitives` / `-WithPrimitives` | 可选：四个 advanced 填表 slash（new-goal 等） |
| `--all` / `-All` | Claude + Grok + Copilot + Codex + prompts/templates/contracts + **core** |
| `--init-workspace` / `-InitWorkspace` | 可选：scaffold `docs/workspace-NNN-slug/`（**须**同时给 slug） |
| `--workspace-slug` / `-WorkspaceSlug` | 与 init-workspace 联用；小写短横线；**禁止静默默认** |
| `--root-slug` / `-RootSlug` | 与 init-workspace 联用 → 计划中的 `GOAL-001-<slug>` |
| `--root-title` / `-RootTitle` | 可选；计划中 Root 标题 |
| `--workspace-nnn` / `-WorkspaceNnn` | 可选；默认 `001` |
| `--skills-dir` / `-SkillsDir` | 默认 `./skills` |

```bash
bash ./skills/install.sh --all --skills-dir ./skills
bash ./skills/install.sh --claude --skills-dir ./skills
bash ./skills/install.sh --grok --skills-dir ./skills
bash ./skills/install.sh --copilot --skills-dir ./skills
bash ./skills/install.sh --codex --skills-dir ./skills
```

```powershell
.\skills\install.ps1 -All -SkillsDir .\skills
.\skills\install.ps1 -Claude -SkillsDir .\skills
.\skills\install.ps1 -Grok -SkillsDir .\skills
.\skills\install.ps1 -Copilot -SkillsDir .\skills
.\skills\install.ps1 -Codex -SkillsDir .\skills
```

可选：安装同时 scaffold 工作区骨架（**不**创建 Root 五件套；slug 必须显式给出）：

```bash
bash ./skills/install.sh --all --skills-dir ./skills \
  --init-workspace --workspace-slug my-product --root-slug product-vision \
  --root-title "Product vision"
```

```powershell
.\skills\install.ps1 -All -SkillsDir .\skills `
  -InitWorkspace -WorkspaceSlug my-product -RootSlug product-vision `
  -RootTitle 'Product vision'
```

安装后：冷启动用 **`/vision`**（Charter→VP）；**`/govern`** 推进（若已 scaffold，则创建 Root 五件套）；Goal 交叉审计用 **`/audit`**；独立 Vision Review 用 **`/vision-audit`**。

## 最小可运行集（消费方）

| 必备 | 来源 |
|------|------|
| 根 `AGENTS.md`（或 copilot-instructions） | install |
| `/govern` + `/audit` + `/vision` + `/vision-audit` + `skills/prompts/*` | install + 包 |
| **`docs/architecture/`**（principles、workspace-protocol、overview、directory-layout） | install 从 `core/` |
| **`docs/templates/`** + 精简 **`docs/README.md`** | install 从 `core/` |
| 现行 Charter + 至少一 VP（完整治理） | `/vision` 冷启动 |
| `docs/workspace-…/workspace.md` + `goal-tree` | `/govern` S0，或 install `--init-workspace`（slug **显式**） |
| Root 五件套 | `/govern` / 原语 01 创建（init-workspace **不**代建） |

| 不要期望随包出现 | 原因 |
|------------------|------|
| monorepo dogfood `GOAL-*` 树 | 过程数据 |
| `tech-stack.md` | 实现栈，非方法论 |
| 完整 monorepo `docs/README` / standalone 测试 | 维护者路径 |

## 在其他项目中快速启用

1. 安装规则 + `/govern` + `/audit` + `/vision` + `/vision-audit`（**同时默认安装 core → `docs/`**）。
2. 核对 `docs/architecture/principles.md` 存在。  
3. 调用 `/vision`：Charter → 首个 VP（+ Vision Review）。  
4. 建立工作区（`/govern` S0 或 `--init-workspace`）并挂 `primary_plan`。  
5. 调用 `/govern`：Root / 子目标推进。  
6. 调用 `/audit`：目标独立审计意见（不改 status）。
7. 调用 `/vision-audit`：独立 Vision Review（只写 `reviews.md`）。

## 核心约定（摘要）

| 规则 | 说明 |
|------|------|
| 核心 + Skills | 同级必备；仅装适配器不算完整 |
| 扁平存储 | 目标平铺在当前 `docs/workspace-<NNN>-<slug>/` 根 |
| 编号 | `GOAL-001` 为 Root；slug 自定 |
| 层级 | 仅 `parent` 字段 |
| 总览 | 变更后更新 `goal-tree.md` |
| 五件套 | meta / decision / execution / audit / attachments |
| 工作区 | `workspace.md` 绑定 Root Goal 与 canonical 范围；legacy `docs/goals/` 仅旧仓兼容 |
| 共享资料 | 固定 `material_id` / source / version / SHA-256；非第二状态 |
| 信息就绪 | 可带未知立项；I-00N 与阶段门禁 |
| 包目录名 | 常为 `skills/`，可改名；按含 `prompts/00-…` 定位 |

## 测试

```bash
# 结构契约 +（Windows）PowerShell 隔离安装冒烟（F-018）
python skills/tests/test_skills_orchestrator.py
```

```powershell
# 仅跑隔离安装冒烟（不经过 unittest）
powershell -NoProfile -ExecutionPolicy Bypass -File .\skills\tests\test_install_ps1_isolated.ps1
```

Windows 上隔离安装冒烟断言默认**四入口**（`/govern`+`/audit`+`/vision`+`/vision-audit`）+ **core docs 落点**，且不含填表 advanced slash、不含 `tech-stack`。`install.sh` 的真实执行仍依赖 bash 环境（本仓库 Windows 主证据以 PS1 为准）。

## 尚未包含

- Marketplace 完整包  
- 编号 / parent 自动校验工具  
- 自动在无维护者授权时创建 GitHub Release（tag CI 仅 pack + 上传 artifact）  
- `/vision-audit` 写盘全路径 e2e（矩阵证据为只读 dispatch / probe，不是写盘 e2e）

当前交付：**core 方法论镜像（默认 install）+ Skills 适配 + 默认四入口（`/govern` `/audit` `/vision` `/vision-audit`）+ 原语 01～05 + 愿景 06 + 独立 Vision Review 07 + 多宿主安装 + 可选 `--init-workspace` + 模板/契约镜像 + pack zip**。monorepo `docs/` 仍为维护者 canonical 上游。
