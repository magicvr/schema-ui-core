---
doc_type: vision-plan
id: VP-006-full-protocol-contract-v2-7-0
title: schema-ui-docs@v2.7.0 整份契约可验证兼容
status: closed
vision_ref: schema-ui-core-admin-foundation@0.3.0
lead_workspace: workspace-005-full-protocol-contract-v2-7-0
created: 2026-08-08
updated: 2026-08-10
version: 0.3.1
parent: null
---

# VP-006 · schema-ui-docs@v2.7.0 整份契约可验证兼容

## 意图

**纠正组合焦点**：本仓库对 `magicvr/schema-ui-docs` **`v2.7.0`**（pinned commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`）的目标是 **整份契约的可验证兼容实现**，而不是长期停留在 MVP 覆盖子集 `I-PROTO-001 v0.1.3` 上把「钉死子集」误当成终态协议成功条件。

`I-PROTO-001 v0.1.3` 仅记录 **MVP 阶段**为验证结构而冻结的覆盖切片；它 **不**解除、也 **不**替代 Charter 成功边界第 1 条所要求的「对 `schema-ui-docs` `v2.7.0` 协议能力形成可验证的兼容实现与示例路径」。

本 VP 在 [VP-003](VP-003-modular-admin-architecture.md) 单主线模块架构与 [VP-004](VP-004-module-contribution-readiness.md) 贡献契约已关闭的底座上，交付：

1. 以 [protocol-inventory-v2.7.0.md](../protocol-inventory-v2.7.0.md) 与上游 pin 为权威输入的 **全量（整份契约）覆盖表**升版与决策留痕；
2. 前后端与 Renderer 对 **registry / 能力域 / conformance 纳入面** 的实现与 fail-closed 边界；
3. 每一纳入能力的 **范例路径 + 可复核验证**（结构 / fixture / 集成按域约定）。

### 用户裁决（2026-08-08）

- **必须**支持 `schema-ui-docs@v2.7.0` **整份契约**（本 VP 主交付）。
- **在此之前**，**不得**启动 [VP-005](VP-005-design-system-and-ui-experience.md) 的视觉优化 / 设计系统实施（含激活开区、Root 视觉阶段开工）；VP-005 仅可保持 `planned` 文档修订，不得作为 `primary_plan` 推进实现。
- 历史 MVP 子集与已关闭 VP 的交付事实 **不重写**；本波次是 **协议覆盖的纠正与扩张**，不是宣称过去已完成全量协议。

### 交付形态定名（主交付）

| 是 | 不是 |
|----|------|
| 覆盖表从 `I-PROTO-001 v0.1.3` **升版**到「整份 v2.7.0 契约」的可审计范围定义 + 实现 + 范例 + 验证 | 默认交付 Design Token / Shell 换皮 / Linear 级视觉产品化（归属 **VP-005**，且本 VP **closed 前禁止启动**） |
| Renderer / 协议 runtime / 必要后端 API 与模块贡献面对齐上游语义 | 订单、钱包、类目、通知等 **业务领域产品**（后续独立业务 VP） |
| 对未实现或有意延后的能力 **fail closed** 并在覆盖表中显式 disposition | 静默降级、用「MVP 够用」宣称全量兼容、私有 Schema 扩展冒充上游 |
| 过程台账挂本 VP 的 delivery 工作区（激活后） | 重开 VP-003 架构迁移；修订 Goal Governance 元规则（principles / 安装 MUST） |

### 覆盖表权威落点（F-V022）

| 角色 | id / 路径 | 规则 |
|------|-----------|------|
| **历史 MVP 基线（只读）** | `I-PROTO-001` v0.1.3 · [I-PROTO-001-coverage-draft.md](../../workspaces/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) | **禁止**就地改写语义或版本号冒充升版；仅作升版起点与回归对照 |
| **整份契约现行权威（本 VP）** | 信息项 id：`I-PROTO-FULL-001`；文件名建议：`I-PROTO-FULL-001-coverage-v2-7-0.md` | 激活后落在 **本 VP lead 工作区** Root `attachments/`；由 Root **新决策**冻结；**新版本号**、新决策，不覆盖 v0.1.3 文件 |
| **发现入口** | 意图期：本节 + [protocol-inventory-v2.7.0.md](../protocol-inventory-v2.7.0.md) 指针；S1 冻结后：Root 决策链 + `docs/architecture/overview.md` / QUICKSTART 可发现「现行覆盖表 = I-PROTO-FULL-001」 | 不得把 workspace-001 的 v0.1.3 继续标成「现行协议覆盖」 |

激活前：本表为意图权威；**尚无** `I-PROTO-FULL-001` 实体文件 = 覆盖未冻结，禁止宣称全量兼容。

## 继承边界

| 来源 | 本 VP 继承 |
|------|------------|
| Charter `@0.2.0` | 成功边界 1（v2.7.0 可验证兼容与示例路径）；协议来源 pin；非目标（不重定义上游语义、不做业务终端产品、无热插拔插件市场）。 |
| 协议 pin | `schema-ui-docs@v2.7.0` / commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`；清单 [protocol-inventory-v2.7.0.md](../protocol-inventory-v2.7.0.md)。 |
| `I-PROTO-001 v0.1.3` | **历史 MVP 基线**，作为升版起点与回归对照；**不是**本 VP 的退出上界。升版须新决策 + **`I-PROTO-FULL-001` 新文件**，禁止静默改写 v0.1.3 文件语义。 |
| VP-003 / `module-architecture.md` | 单主线、薄内核、Fx 组合根、Profile、后端聚合 Manifest；不破坏、不重开迁移。 |
| VP-004 / playbook | 一方模块贡献 must / must-not；新增协议相关模块或能力须符合 playbook。 |

## 方向级退出判据

在同时满足下列方向、且均有工作区 Q2 证据时，本 VP **可以**提议 `closed`：

1. **整份契约覆盖表落盘且已决策**  
   存在现行覆盖表 **`I-PROTO-FULL-001`**（新文件 + 新版本号，继承 inventory / registry / fixture suites；**不**就地改写 `I-PROTO-001 v0.1.3`），对 `v2.7.0` pin 下的能力域与 component registry **逐项**给出 disposition。  
   **默认 disposition** = 对 inventory / registry / fixture **承诺面**的可验证 **`include`**。  
   **`include-partial` 收窄定义**：仅允许在「能力已纳入」前提下，对保真边角、未测 fixture 子集或可列明的次要语义缺口写显式边界；**不得**用 `include-partial` 表达「整域或主要子面不打算做」（与历史 I-PROTO 对 D-ACT/D-TABLE/D-COMP/D-FORM/D-UPLOAD 等范围 partial 模式同构的写法 **禁止**作为默认覆盖路径）。  
   **范围收缩**（示例：整域 upload、批量 selection / ADR-0022 全排除、registry 大支整支排除）必须二选一：**(a)** `exclude` + 用户书面有界残余（范围 + 复审触发），或 **(b)** 用户书面接受的有界范围收缩（范围 + 复审触发）；**不得**仅靠 S1 规划便利落成 `include-partial`。  
   任何 `exclude` 必须是用户书面接受的有界残余（范围 + 复审触发），**不得**以「MVP 曾经排除」作为默许。  
   **禁止**以「又钉一个更大/更小的 partial 子集 + 表上逐项有行」替代「整份契约覆盖表已决策」。

2. **Renderer / 协议 runtime 对齐纳入面**  
   对覆盖表 `include` /（合法）`include-partial` 的 node type、表单控件、table/form/action/upload 等能力，前端具备可观察实现或明确 fail-closed；禁止钉死白名单外静默忽略。保真度达到「契约语义可验证」，不要求 VP-005 级视觉产品化。

3. **后端与模块贡献面对齐**  
   纳入能力所需的列表/详情/动作/上传/权限等服务端契约在基架或一方模块路径上可验证；与 Manifest / Schema 贡献模型一致，不引入中央业务硬编码回潮。

4. **范例与验证路径闭合**  
   每一纳入能力域有可发现的范例（页面或场景）与验证入口（schema / conformance fixture / 集成测试按域登记）；上游 pin 的 conformance 纳入面有执行证据，exclude 面有表可查。

5. **回归与兼容声明诚实**  
   既有 `I-PROTO-001 v0.1.3` 主路径与 E2E/Smoke **不回退**；对外不得在覆盖表未闭合时宣称「已完整支持 v2.7.0」；文档（inventory 指针、overview/QUICKSTART 发现路径）反映全量覆盖权威（`I-PROTO-FULL-001`）而非仅 MVP 子集。

6. **过程可关门**  
   lead 工作区 Root（或等价交付目标）完成约定范围、开放 required findings = 0，Vision Review 无阻断本 VP 关门的开放 required；用户确认关门。

## 建议实现阶段（非退出判据正文，供后续 `/govern` 参考）

| 阶段 | 阶段目的 | 建议产物 / 检查点 |
|------|----------|-------------------|
| S0 | 差距盘点 | 覆盖表 v0.1.3 vs inventory/registry/fixtures 差集；前端保真债 vs 未纳入 type 分列 |
| S1 | 覆盖表升版冻结 | 落盘 `I-PROTO-FULL-001` + Root 决策；默认 `include`；`include-partial` 仅保真/边角；范围收缩 → exclude 或用户书面 residual；**可审计摘要**：相对 v0.1.3 差集中「转为 include 的计数 / 仍 residual 的清单」 |
| S2 | 核心缺口实现 | 未实现 registry type / 批量 selection / upload 等按表纳入批次交付 |
| S3 | 保真与 runtime | 钉死内降级控件（如 cascader/richText 等）提升到契约语义；表达式/权限边角 |
| S4 | 范例 + conformance | 范例路径、fixture/集成门禁、失败路径 fail-closed |
| S5 | 文档与关门 | 发现路径、兼容声明、回归、close-out 审计 |

具体子目标与批次以实现层为准；允许按域分批，但 **VP 关门**须满足上述方向级退出，不得以「又钉了一个更小子集」或「大面积 include-partial 伪装整份契约」替代目标，除非用户对明确 exclude / 范围收缩做 residual。

## Non-goals（非目标）

- **不**启动或吸收 VP-005 视觉/设计系统/Shell 产品化（本 VP closed 前硬禁止）。
- **不**交付订单、钱包、类目、通知等业务领域模块。
- **不**重开 VP-003 架构迁移、不恢复长期双线、不引入运行时插件/热插拔。
- **不**在本项目内重新定义或替代上游协议语义；分歧回上游或书面兼容决策。
- **不**修订 Goal Governance 核心方法论（`principles.md` P-001～P-006 等）。
- **不**把历史 `I-PROTO-001 v0.1.3` 或已关闭 VP 的「子集完成」改写成「全量协议已完成」。
- 不为 VP 在 `docs/vision/` 建立 Goal 五件套或 progress% 权威。

## 与前后 VP 的关系

| VP | 关系 |
|----|------|
| VP-001 / VP-002 | 历史 MVP/生产基线与子集验证；已关闭。本 VP **扩张覆盖**，不重开其过程树。 |
| VP-003 / VP-004 | **前置已关闭**；本 VP 在其架构与贡献契约上补齐全量协议。 |
| **VP-005** | **后置且硬门闩**：本 VP **未 `closed` 前**，禁止激活/开区/实施 VP-005 视觉优化。VP-005 将来继承本 VP 交付的全量契约面再做产品化视觉。 |
| 后续业务 VP | 默认继承本 VP 的协议覆盖与 VP-003/004 架构/playbook；建 VP 前仍须 `/vision` 复核。 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-005-full-protocol-contract-v2-7-0 | GOAL-001-full-protocol-contract-v2-7-0 | lead | 2026-08-08 | 用户确认激活后唯一 lead / delivery；`/govern` 已 scaffold；**禁止**在 closed VP-003/004 工作区吸收本意图 |

用户已于 2026-08-08 确认将本 VP **激活**（`planned` → `active`），并指定唯一 lead / delivery 工作区 **`workspace-005-full-protocol-contract-v2-7-0`**（slug 按 VP-006 id 与既有 workspace-00N 惯例，用户本轮书面授权开区）。同日 `/govern` 完成物理 scaffold：Root [GOAL-001-full-protocol-contract-v2-7-0](../../workspaces/workspace-005-full-protocol-contract-v2-7-0/GOAL-001-full-protocol-contract-v2-7-0/00-meta.md)，`primary_plan` / `plan_refs` 均为本 VP。激活与建区 **不**构成覆盖表冻结或协议全量兼容证据。

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-08 | **closed** | 用户书面确认关门（「确认关门」）。退出判据 1–6 全部闭合：覆盖表 `I-PROTO-FULL-001` v1.0.0 冻结（12/12 域、24/24 registry、16/16 行为套件 include；0 exclude）；Renderer/后端纳入面实现；8 范例页 + 320 case 全量分母；vitest 569/569 + `go test ./...` 全绿 + e2e 2/2 + 回归不回退；文档指针 = `I-PROTO-FULL-001`；S1/S5 独立审计 A-001/A-002，**开放 required = 0** | 历史关门证据保留：覆盖表 v1.0.0、`E-001～E-005`、`A-001/A-002` | 2026-08-10 执行分母勘误见下行；历史关门事实不变 |
| 2026-08-10 | **errata** | 响应 workspace-008 F-001：覆盖权威升为 `I-PROTO-FULL-001` v1.0.1；12/12 域、24/24 registry、16/16 suite include 保持；320 case 现行执行口径为 **318 executed + 2 local adapter excluded**，两项 exclusion 原因与复审触发见 D-003 / E-007。VP-006 与 Root 终态不重开。 | `docs/workspaces/workspace-005-full-protocol-contract-v2-7-0/GOAL-001-full-protocol-contract-v2-7-0/attachments/I-PROTO-FULL-001-coverage-v2-7-0.md`; `…/01-decision/D-003-i-proto-full-errata.md`; `…/02-execution/E-007-i-proto-full-errata.md` | 域/协议范围 residual 为空；local adapter exclusion 为有界执行口径，复审触发 = 协议 pin/disposition 或错误包络变化 |

## 规划修订短史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2026-08-08 | `0.1.0` | 用户书面裁决：目标为 `schema-ui-docs@v2.7.0` **整份契约**；创建本 VP 为当前组合主意图；**硬阻塞** VP-005 视觉实施直至本 VP closed。`planned`，未激活、未绑工作区。 |
| 2026-08-08 | `0.1.1` | 响应 [VRev-012](../reviews/VRev-012-vp006-full-protocol-contract.md)：exit 1 收紧 disposition（默认 `include`；`include-partial` 仅保真/边角；范围收缩 → exclude 或用户 residual）；覆盖表权威落点 `I-PROTO-FULL-001`（F-V021/F-V022）。仍 `planned`，未激活、未绑工作区。 |
| 2026-08-08 | `0.2.0` | 响应 [VRev-013](../reviews/VRev-013-vp006-post-closure-reaudit.md) pass；用户确认激活：`planned` → `active`；绑定唯一 lead / delivery `workspace-005-full-protocol-contract-v2-7-0`（Root `GOAL-001-full-protocol-contract-v2-7-0`）。物理 scaffold 交 `/govern`；未将激活写成全量兼容已交付。 |
| 2026-08-08 | `0.3.0` | **用户书面确认关门**（「确认关门」，2026-08-08）：`active` → `closed`；关门记录见上表。整份契约覆盖与实现证据链完整；VP-005 视觉实施冻结解除条件满足（本 VP closed），但 `F-V018` 仍为 open required，VP-005 不自动放行，解冻与否由用户另行决策。 |
| 2026-08-10 | `0.3.1` | editorial 勘误：`I-PROTO-FULL-001` 升至 v1.0.1；保留 12/12 域、24/24 registry、16/16 suite include，将执行分母修正为 318 executed + 2 local adapter excluded。VP status、退出判据和历史审计 verdict 不变。 |
