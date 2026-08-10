---
doc_type: vision-plan
id: VP-004-module-contribution-readiness
title: 一方模块贡献就绪（作者与 AI 操作契约）
status: closed
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-004-module-contribution-readiness
closed_under_vision_ref: schema-ui-core-admin-foundation@0.2.0
created: 2026-08-06
updated: 2026-08-06
version: 0.2.0
parent: null
---

# VP-004 · 一方模块贡献就绪（作者与 AI 操作契约）

## 意图

在 [VP-003](VP-003-modular-admin-architecture.md) 已关闭的单主线模块化终态之上，交付**可被合作者与 AI 工具共同遵循的操作契约**：新增一方功能模块时**必须完成什么、明确不必/禁止做什么**，以及后续演进应归入**薄内核 / 组合根**还是**独立（或既有）模块**的判定方法。

本 VP **不**再交付架构迁移或具体业务领域模块。架构边界权威仍为 [module-architecture.md](../../architecture/module-architecture.md)；本波次把该边界**操作化**为可发现、可审计的 playbook 与归属方法论，并为实现层提供可挂接的 `primary_plan`（治理过程落盘位）。

### 交付形态定名（主交付）

本 VP 的**主交付形态**是 `docs/architecture/` 下**产品模块贡献方法论 / 操作 playbook 的新增或修订**（可扩展 [module-architecture.md](../../architecture/module-architecture.md)，或新建并从该文与 overview 链出的 authoring 文）。它是**产品架构作者指南**，不是运行时功能交付波次。

明确划界：

| 是 | 不是 |
|----|------|
| `docs/architecture/` 产品模块贡献方法论 / playbook 新增或修订 | 默认交付代码脚手架、参考模块骨架、自动化检查脚本（可选加分，默认不进退出分母） |
| 把 `module-architecture.md` 终态**操作化**为 must / must-not / 归属判定 | 重开架构迁移或交付业务领域模块 |
| 过程台账挂本 VP 的 delivery 工作区（激活后） | 修订 Goal Governance 核心方法论（[principles.md](../../architecture/principles.md) P-001～P-006、workspace-protocol 等治理 MUST） |

词义说明：本仓 `overview.md` 将 `docs/architecture/` 放在「核心方法论与文档协议」框图内；`principles.md` 自称「Goal Governance 核心方法论」元规则。本 VP 所说「方法论文档」**仅指产品模块贡献侧**，**不**授权修订治理元规则或安装 MUST。

### AI 发现路径充分条件

标题中的「作者与 AI 操作契约」在退出层按下列充分条件理解（闭合 VRev-010 `F-V017`）：

1. **默认充分路径**：从 [overview.md](../../architecture/overview.md) 与根 [QUICKSTART.md](../../../QUICKSTART.md)（或与模块扩展相关的 README/QUICKSTART 等价入口）之一，能到达上述 playbook 权威文，且 playbook 正文本身对人类与 AI 可读、可遵循。
2. **不默认改写**根 `AGENTS.md`、Skills 发现路径或其他 AI 适配器入口；此类接线**不是**本 VP 默认退出分母。
3. 若激活后用户书面将指定 AI 入口接线纳入 Root 路线图，可作为可选加分或显式检查点，**不得**回溯把「未改 AGENTS/Skills」读成 exit 4 失败。

**内容权威**与**过程台账**分工：

| 层 | 写什么 | 落点 |
|----|--------|------|
| 内容（产品/架构真相） | 清单、判定树、反例、与代码路径对齐 | `docs/architecture/`（扩展既有架构文或新建 authoring 文）；再链到 overview / QUICKSTART 等发现路径 |
| 过程（治理台账） | 立项取舍、实施事实、审计 | 本 VP 绑定工作区的 Goal 五件套 + `goal-tree.md`（激活后由 `/govern` 建立） |

## 继承边界

| 来源 | 本 VP 继承 |
|------|------------|
| Charter `@0.2.0` | 单主线、薄内核、可 fork；业务领域非愿景成功条件 |
| VP-003 / `module-architecture.md` | 模块契约（核心六项 / 按需）、组合根静态接入、Profile、迁移与 Manifest 聚合规则；**不重开**架构退出判据 |
| I-PROTO-001 `v0.1.3` | 协议覆盖不扩张；本 VP 不改协议子集 |
| QUICKSTART §5 | 既有「加页面」最小步骤可被 playbook 引用或升级，不得与终态架构矛盾 |

## 方向级退出判据

在同时满足下列方向、且均有工作区 Q2 证据时，本 VP **可以**提议 `closed`：

1. **新增模块 playbook 落盘**  
   存在现行架构文档（或明确从 `module-architecture.md` 链出的权威子文），写清一方**标准 Admin 功能模块**从零接入时的**必须完成**项（至少覆盖：模块 id/版本/依赖声明、核心六项贡献、组合根静态候选注册、Profile/`modules.enabled` 成员关系、全局迁移台账参与规则、验证/回归最小集），并与仓库真实路径一致。

2. **明确「不必做 / 禁止做」**  
   同一权威入口含对等的 **DO NOT** 表（至少覆盖：不为接模块改 Renderer/Shell 中央业务注册；不在生产路径静默静态 Manifest；不私建平行认证/授权/DB；不把「按需能力」当作核心六项可永久缺省；不做运行时插件/热插拔幻想），避免合作者/AI 过度施工。

3. **Core vs 模块责任方法论**  
   同一权威入口含可执行的**归属判定**（薄内核 / 组合根 / 新模块 / 既有模块 / 模块内 util），含横切基础设施 vs 标准 Admin 功能模块的差异，以及若干**正反例**；判定结果不得与 `module-architecture.md` §1 / §6 冲突。

4. **可发现性**  
   合作者与 AI 从约定发现路径（至少：`docs/architecture/overview.md` 或等价索引、与模块扩展相关的 README/QUICKSTART 入口之一）能到达上述权威文，而无需阅读已关闭 VP-003 的全过程治理树。  
   **充分条件**（方向级）：上述 overview + QUICKSTART（或等价）接线即为 AI 侧默认充分；**不**要求默认修订 `AGENTS.md` / Skills 发现路径（见意图节「AI 发现路径充分条件」）。

5. **过程可关门**  
   lead 工作区 Root（或等价交付目标）完成约定范围、开放 required findings = 0，且 Vision Review 无阻断本 VP 关门的开放 required；用户确认关门。

可选加分（**不**默认计入退出分母，除非激活后 Root 路线图显式纳入）：参考模块骨架/脚手架、自动化「新模块契约」检查脚本、指定 AI 入口（如 `AGENTS.md` / Skills）接线。未纳入则不得以「缺脚手架」或「未改 AGENTS/Skills」阻断关门。

## 建议实现形状（非退出判据正文）

激活并开区后，Root 可用极薄纲领（示例，由 `/govern` 裁剪）：

| 阶段 | 目的 |
|------|------|
| S1 | 盘点现状缺口 vs `module-architecture` / 代码路径；冻结 playbook 大纲与权威文件路径 |
| S2 | 落盘 must / must-not / Core-vs-module；对齐修订 |
| S3 | 发现路径接线 + 与现有一方模块对照抽检 |
| S4 | 阶段/关门审计与 VP 关门提案 |

具体子目标编号与检查点以实现层为准。

## 工作区绑定

用户已于 2026-08-06 确认将本 VP **激活**（`planned` → `active`），并指定唯一 lead / delivery 工作区 **`workspace-004-module-contribution-readiness`**（slug 用户书面确认）。同日 `/govern` 完成物理 scaffold：Root [GOAL-001-module-contribution-readiness](../../workspaces/workspace-004-module-contribution-readiness/GOAL-001-module-contribution-readiness/00-meta.md)，`primary_plan` / `plan_refs` 均为本 VP。

建区**不**构成 playbook 已交付或任一退出判据已满足。禁止在 closed VP-003 工作区（`workspace-003-modular-admin-architecture`）吸收本意图。

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-004-module-contribution-readiness | GOAL-001-module-contribution-readiness | delivery / lead | 2026-08-06 | 历史绑定保留；Root `done / 4/4`；playbook 已交付；不改变 Charter `primary_workspace` |

## Non-goals

- 不交付订单、钱包、类目、通知等**业务领域模块**（仍属 roadmap 后续独立 VP 候选）。
- 不重开 VP-003 架构迁移、不恢复长期双线、不引入运行时插件市场 / `.so` / 热插拔。
- 不扩张 `I-PROTO-001 v0.1.3` 协议覆盖；不改写上游 Schema 语义。
- **不**修订 Goal Governance 核心方法论（`principles.md` P-001～P-006、workspace-protocol 等治理 MUST / 安装契约）。
- **不**以代码脚手架或默认改写 `AGENTS.md` / Skills 作为本 VP 退出必要条件。
- 不把「仅会话约定」或「仅 git message」当作本 VP 退出证据；过程必须在绑定工作区目标内。
- 不为 VP 建立 Goal 五件套替代品于 `docs/vision/`。

## 与前后 VP 的关系

| VP | 关系 |
|----|------|
| VP-001 / VP-002 | 历史基线；不重开 |
| VP-003 | **前置已关闭**；本 VP 消费其终态与 `module-architecture.md`，不修改其关门事实 |
| 未来业务 VP | 应**引用**本 VP 交付的 playbook 为默认贡献约束；不得在业务 VP 内重新发明平台归属法作为唯一来源 |

## 关门记录

仅在 `closed` 或 `abandoned` 时填写。

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-06 | **closed** | 五条方向级退出判据均以 lead 工作区 Q2 证据满足：① **playbook MUST**（`module-contribution-playbook.md` §1 M1–M6：id/版本/依赖、核心六项、组合根、Profile、全局迁移、验证最小集；路径对齐现网）；② **DO NOT**（同文 §2 D1–D5，含 Renderer/Shell、静默静态 Manifest、平行认证授权 DB、按需误读、热插拔）；③ **Core vs 模块归属**（同文 §3 判定树+正反例+横切 vs 标准 Admin；与 `module-architecture.md` §1/§6 同向）；④ **可发现性**（`overview.md`「一方模块扩展」+ 根 `QUICKSTART.md` §5 + architecture §9；AI 默认充分路径 a，未默认改 AGENTS/Skills）；⑤ **过程可关门**（lead Root `GOAL-001-module-contribution-readiness` `done / 4/4`；A-001 self close-out `pass` + A-002 independent `pass` + A-003 `/govern` 响应采纳 pass 且 F-001 recommended `fixed`；Root 03-audit **开放 required=0**；Vision Review **0 open required**（VRev-010 及 F-V016/F-V017 已 `fixed`）；用户指令确认关门）。主交付：`docs/architecture/module-contribution-playbook.md`。 | [Root 00-meta](../../workspaces/workspace-004-module-contribution-readiness/GOAL-001-module-contribution-readiness/00-meta.md)；[goal-tree](../../workspaces/workspace-004-module-contribution-readiness/goal-tree.md)；[Root 03-audit](../../workspaces/workspace-004-module-contribution-readiness/GOAL-001-module-contribution-readiness/03-audit.md)；[A-001](../../workspaces/workspace-004-module-contribution-readiness/GOAL-001-module-contribution-readiness/03-audit/A-001-root-closeout-self.md)；[A-002](../../workspaces/workspace-004-module-contribution-readiness/GOAL-001-module-contribution-readiness/03-audit/A-002-root-closeout-and-vp004-alignment-independent.md)；[A-003](../../workspaces/workspace-004-module-contribution-readiness/GOAL-001-module-contribution-readiness/03-audit/A-003-response-a002.md)；[vp004-close-proposal](../../workspaces/workspace-004-module-contribution-readiness/GOAL-001-module-contribution-readiness/attachments/vp004-close-proposal.md)；[playbook](../../architecture/module-contribution-playbook.md)；[module-architecture §9](../../architecture/module-architecture.md) | 无有界 residual。可选加分（脚手架 / AGENTS·Skills 接线）从未纳入退出分母，不构成残余。Non-goals（业务领域模块、重开 VP-003 迁移、热插拔、principles 修订等）保持排除。 |

## 规划修订短史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2026-08-06 | `0.1.0` | 用户确认结构选型路径 B：为模块贡献方法论提供可挂接意图与治理过程落盘位；`planned`，未激活、未绑工作区。 |
| 2026-08-06 | `0.1.1` | 响应 VRev-010：意图节增加交付形态定名（产品模块贡献方法论/playbook；非脚手架默认交付；非 principles/治理 MUST）与 AI 发现路径充分条件（默认 overview+QUICKSTART；不默认改 AGENTS/Skills）；exit 4 / Non-goals / 可选加分同步。仍 `planned`，未激活、未绑工作区。 |
| 2026-08-06 | `0.1.2` | 用户确认激活：`planned` → `active`；绑定唯一 lead / delivery `workspace-004-module-contribution-readiness`（建议 Root `GOAL-001-module-contribution-readiness`）。物理 scaffold 交 `/govern`；未将激活写成 playbook 已交付。 |
| 2026-08-06 | `0.1.3` | `/govern` scaffold 完成：区目录 + Root 五件套 + goal-tree；绑定表 root 从「待 scaffold」改为已建立。未勾选 S1–S4、未交付 playbook。 |
| 2026-08-06 | `0.2.0` | 关门：五条方向级退出判据经 lead 工作区 Q2 证据满足（Root `done / 4/4` + A-001/A-002/A-003；playbook + 发现路径；Vision Review 0 open required），用户指令确认 → `status` `active` → `closed`；关门记录 + roadmap/workspaces/charter 关系节同步；`closed_under_vision_ref = @0.2.0`。 |
