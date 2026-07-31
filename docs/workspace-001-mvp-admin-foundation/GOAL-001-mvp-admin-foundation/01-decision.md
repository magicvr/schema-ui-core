---
id: GOAL-001-mvp-admin-foundation
doc: decision
status: active
parent: null
created: 2026-07-31
updated: 2026-07-31
version: 0.4.3
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

权威表见 [00-meta.md](00-meta.md)「信息就绪与未知项」。摘要：

| ID | 级别 | 最晚阶段 | 状态 | 阻断 |
|----|------|----------|------|------|
| I-PROTO-001 | required | R2 方案冻结前 | **verified**（D-009；v0.1.3 冻结基线） | R2 已冻结；不得把该子集表述为完整协议支持，R4/R5 仍受各自信息门禁约束 |
| I-PROTO-002 | required | R4 实施前 | open | 未 verified 前不得宣称账号权限链路完成 |
| I-PROTO-003 | required | R5 验收前 | open | 未 verified 前不得验收/对照 VP 退出判据关门 |
| I-PROTO-004 | non-blocking | 实施前为宜 | open | 不阻断开区；影响校验工程策略 |
| I-STACK-001 | required | R1 实施前 | **verified**（D-004） | 已确认 monorepo 布局与包管理；可启动 R1 子目标实施 |
| I-STACK-002 | non-blocking | R1 内 | **verified**（D-004） | monorepo `apps/web`+`apps/api`；端口/env 细节随 GOAL-002/003 |

## D-001 · 开区并挂接 VP-001 为 primary

**决定**：

1. 创建显式工作区 `workspace-001-mvp-admin-foundation`（`vision_role: primary`）。
2. Root = `GOAL-001-mvp-admin-foundation`（`parent: null`，`status: active`）。
3. `plan_refs` / `primary_plan` = `VP-001-mvp-admin-foundation`（`vision_ref` 已匹配 `schema-ui-core-admin-foundation@0.1.0`）。
4. `shared_materials_catalog: none`（暂无共享资料）。
5. 将 VP-001 自 `planned` 升为 `active`，`lead_workspace` = 本工作区；同步 Charter `primary_workspace` 与 `workspaces.md`。

**为什么**：

- 冷启动上环（Charter → VP）已完成；VRev required 为 0 open；法定下一步是工作区 + Root。
- 用户明确指令：首区 + Root，挂 `primary_plan = VP-001-mvp-admin-foundation`。
- slug / 角色 / 状态经用户确认，禁止静默默认。

**未选方案**：

- **继续只停在愿景层**：无法承载实施证据与 P-001 路线图。
- **legacy `docs/goals/`**：新项目禁止；须显式工作区。
- **`vision_role: delivery` 作首区**：与「主交付 / primary_workspace」不一致。

## D-002 · 大目标先纲领路线图，不批量建细子目标

**决定**：

本回合只建立 Root 五件套与六段纲领路线图（R1–R6）及信息项；**不**在开区当下批量创建细粒度子目标。进入 R1 实施前先闭合 `I-STACK-001`（及按需 `I-PROTO-004`）；进入方案冻结前必须闭合 `I-PROTO-001`。

**为什么**：

- P-001：范围大、步骤多 → 先高层阶段再按阶段立项。
- P-005：覆盖子集与脚手架仍为开放 required 信息，不得假装已知。

**未选方案**：

- **立刻拆一堆前后端子目标**：在覆盖与脚手架未定时易返工，且易跳过信息门禁。

> **后续**：2026-07-31 在 `I-STACK-001` verified（D-004）后，按阶段创建 R1 三子目标（GOAL-002/003/004）；**不**撤回本条「开区当下不批量建细目标」的历史决定。

## D-003 · 协议边界与禁止主张

**决定**：

- 协议固定源与实施清单以 [protocol-inventory-v2.7.0.md](../../vision/protocol-inventory-v2.7.0.md) 为准。
- **禁止**在 `I-PROTO-001` verified 前主张“支持全部协议功能”或把 `mvp_candidate` 列当作已冻结覆盖集。
- `mvp_candidate` 仅作 R2 决策输入。

**为什么**：

## D-004 · R1 脚手架与平行仓复用策略（闭合 I-STACK-001 / I-STACK-002）

**日期**：2026-07-31  
**状态**：accepted  
**关联信息项**：`I-STACK-001` → `verified`；`I-STACK-002` → `verified`

**决定**：

1. **仓库形态（I-STACK-002）**：本仓 **monorepo**  
   - 前端：`apps/web/`  
   - 后端：`apps/api/`  
   - 治理与分发保持根级：`docs/`、`skills/`、`AGENTS.md` 等  
2. **前端包管理与工具链（I-STACK-001）**：  
   - **npm** + `package-lock.json`（与平行仓 `allinme.web-client` 一致）  
   - Vite + React + TypeScript  
   - R1 即接入 **Tailwind CSS + shadcn/ui 基线** + 浅/深色最小占位（完整 Admin 外壳仍属 R3）  
3. **后端包管理与布局（I-STACK-001）**：  
   - Go modules，module 根在 `apps/api/`  
   - 分层取向：`cmd/server`、`internal/`、`pkg/`（参考平行仓，不照搬 module path）  
4. **平行仓复用（本地已克隆，主看 `dev`）**：  
   - 来源：`../allinme.core-api` @ `dev`（观察提交 `387896d`）、`../allinme.web-client` @ `dev`（`57616d4`）  
   - 策略：**结构 + 通用层参考 / 择优移植，禁止整仓拷贝**  
   - **可参考**：Go cmd/internal/pkg、Makefile/env/health 模式、JWT/SQLite/response 等通用模式（R1 不强制完整 auth）；Web 的 host/protocol/renderer 边界与 Vite 习惯  
   - **不搬入 MVP 默认树**：订单 / 钱包 / 通知等业务域与对应 mock；平行仓 page schema 中曾出现的协议 **2.4** 声明（本仓目标为 **schema-ui-docs@v2.7.0**）  
5. **R1 子目标拆分**：  
   - `GOAL-002-r1-repo-layout-conventions` — 布局与包管理约定落盘  
   - `GOAL-003-r1-api-go-scaffold` — Go 可运行骨架  
   - `GOAL-004-r1-web-react-scaffold` — React 可运行骨架（含 Tailwind/shadcn 基线）  
6. **仍开放**：`I-PROTO-004`（vendor vs pin）non-blocking；默认端口具体数字在 GOAL-002/003 细化（建议沿用 API `:8080` 作起点，未写死为门禁）。

**为什么**：

- R1 实施前 required `I-STACK-001` 必须闭合，否则不得批量生成代码树。  
- 单仓 monorepo 符合 Charter「可 fork 的 React+Go 基架」。  
- 平行仓已验证可运行分层，但含业务非目标与旧协议痕迹，整拷风险高于收益。  
- 用户 2026-07-31 在 `/govern` 中确认全部推荐项。

**未选方案**：

- **根下短名 `web/` + `api/`**：亦可，但与后续可能的 `apps/docs-site` 等扩展相比，`apps/*` 更清晰。  
- **仅本仓写约定、实现仍在外部平行仓**：不利 VP 证据落在本工作区，不作 MVP 主路径。  
- **pnpm**：省磁盘但与现成 `package-lock` 平行仓摩擦更大。  
- **R1 不做 Tailwind/shadcn、R3 再加**：增加产品化债；用户未选。  
- **整树拷贝再删业务**：噪声与协议版本污染风险高。  
- **单一大 R1 子目标**：并行与验收边界糊。

**影响**：

- 放行 R1 子目标创建与骨架实施；**不**放行 R2 协议覆盖冻结（`I-PROTO-001` 仍 open）。  
- 不自动创建应用代码（本决策只定布局与子目标）。

**后续**：

1. 执行 GOAL-002 → 约定文档 / 占位。  
2. 并行 GOAL-003 / GOAL-004 骨架。  
3. R1 可运行后 self 阶段审；再进入 R2。

- 闭合 VRev `F-V001` 只证明清单提取，不证明覆盖冻结（Charter H-001 分列）。

## D-005 · 进入 R2 并起草 I-PROTO-001 纳入/排除表（草案 · 未冻结）

**日期**：2026-07-31  
**状态**：proposed（草案已落盘；**待用户确认后**升 `accepted` 并 `I-PROTO-001` → verified）  
**关联信息项**：`I-PROTO-001` → `collecting`（**非** verified）

**决定**：

1. 纲领阶段 **进入 R2**（R1 已完成；R2 标记为进行中）。
2. 对照 [protocol-inventory-v2.7.0.md](../../vision/protocol-inventory-v2.7.0.md) §3，落盘纳入/排除**草案**：  
   [attachments/I-PROTO-001-coverage-draft.md](attachments/I-PROTO-001-coverage-draft.md)
3. **草案默认取向**（可被用户改写，确认前不构成冻结）：
   - **include（8）**：`D-NODE`、`D-EXPR`、`D-DATA`、`D-ACT`、`D-PERM`、`D-APP`、`D-VER`、`D-VAL`
   - **include-partial（3）**：`D-COMP`（组件 type 白名单）、`D-TABLE`（排序+搜索表；**不含**多选批量）、`D-FORM`（基础表单；**不含** 2.6/2.7 扩展进阶控件）
   - **exclude（1）**：`D-UPLOAD`
   - fixture：`uploads` exclude；`scenarios` 仅 support-only；`component-format` partial；其余与上表 include 域对齐为 include
4. **明确不做的主张**：本决策**不**将 `I-PROTO-001` 标为 verified；**不**放行 R4/R5 实施范围冻结；**不**主张“支持全部协议功能”。
5. **本回合不新建 R2 子目标**：覆盖冻结以 Root 决策 + 附件表为主交付；确认冻结后再按需拆 R3+ 子目标。

**为什么**：

- 用户 `/govern` 指令：进入 R2 并起草 `I-PROTO-001` 纳入/排除表草案。
- P-005：允许带未知推进信息收集；草案 = `collecting`，冻结才 `verified`。
- Charter 要求核心账号权限 + 每纳入域可验证；`mvp_candidate` 不可直接当覆盖集（D-003）。
- R1 已关门，R2 为下一串行纲领阶段。

**未选方案**：

- **直接把 inventory 全 domain 标 include**：违反 MVP 最小与 Charter 非目标，且无法在本波次举证。
- **静默 verified**：无用户确认的纳入边界，属伪闭合。
- **立刻建 R2 细子目标再写表**：表尚未确认，拆目标易空转；P-001 允许先决策再按阶段立项。
- **把 D-UPLOAD / 进阶表单强行纳入 MVP**：扩大 R4/R5 面，偏离“核心账号权限 + 基架范例”。

**影响**：

- R2 进行中；信息项 `I-PROTO-001` = collecting。
- **仍阻断**：方案冻结完成宣称、R4 权限实施范围定稿、R5 验收对照、任何全协议支持表述。
- 不修改子目标 status；不改 `progress`（仍 1/6，R2 检查点未完成）。

**后续**：

1. 用户确认或修订草案开放点 Q1–Q5（见附件 §5）。  
2. 确认后追加正式冻结决策（见后续 D 号），并将 `I-PROTO-001` → `verified`。  
3. 再规划 R3 外壳子目标；R4 依赖已冻结的 `D-PERM` 边界（`I-PROTO-002`）。

## D-006 · 响应 A-001：接受 R1 独立复核 pass，继续 R2 收集且不冻结覆盖

**日期**：2026-07-31  
**状态**：accepted  
**关联意见**：Root `03-audit` **A-001**（independent，verdict=pass）→ 响应节 **A-002**  
**关联信息项**：`I-PROTO-001` 保持 **collecting**（**不** verified）

**决定**：

1. **接受** A-001 对 R1 关门证据的独立复核结论：**pass**（开放 required finding = 0）。  
2. **维持** R1 完成事实不变：纲领 R1 = 完成；GOAL-002/003/004 保持 `done`；Root `progress` 保持 **1/6**；不回退、不重开 R1。  
3. **不**因 A-001 另做 Root 级 R1 自审（P-004.1：用户本轮书面确认以独立意见推进；子目标侧已有各 A-003 self pass 作既有自审证据）。  
4. **继续** R2 的 `I-PROTO-001` 信息收集：草案 [attachments/I-PROTO-001-coverage-draft.md](attachments/I-PROTO-001-coverage-draft.md) 与 D-005（proposed）仍有效。  
5. **明确不做**：不将 `I-PROTO-001` 标为 verified；不冻结 MVP 协议覆盖范围；不主张“支持全部协议功能”；不进入 R3/R4 实施范围定稿；不改 Root `status`。

**为什么**：

- 用户 `/govern` 指令原文：确认 R1 关门证据复核为 pass；维持 R1 完成事实；继续 R2 `I-PROTO-001` 信息收集；不冻结协议覆盖范围。  
- A-001 独立复现了可执行证据，与三子目标 A-003 self pass 同向、无冲突 required。  
- P-005 / D-003：覆盖子集冻结须用户对纳入表书面确认后另决策；本回合仅响应审计，不越权闭合信息门禁。

**未选方案**：

- **因独立审 pass 而静默冻结 I-PROTO-001**：A-001 scope 仅 R1，不含覆盖冻结；且用户明确禁止。  
- **回退或重开 R1**：无 contradictory finding；与 pass 结论相反。  
- **未问用户即强制 Root 再自审一遍 R1**：用户已书面选择基于现有独立意见继续。  
- **跳过响应留痕直接推进 R3**：须先闭合 R2 信息门禁（`I-PROTO-001` verified）后再谈 R3 立项。

**影响**：

- A-001 无 required finding 需 `fixed`/`residual`/`overruled`；响应 = 接受 verdict + 维持阶段事实。  
- R2 仍进行中；`I-PROTO-001` 仍阻断「方案冻结完成」宣称。  
- progress / goal-tree 状态表不因本决策变化。

**后续**：

1. 用户确认/修订草案 Q1–Q5 → 正式冻结决策 → `I-PROTO-001` → `verified`。  
2. 冻结后再规划 R3 Admin 外壳子目标。  
3. 可选：对 R2 草案另跑 `/audit`（非本回合强制）。

## D-007 · 响应 A-003 的用户裁决：固定已定边界，保留 Q2 控件清单为未决

**日期**：2026-07-31
**状态**：accepted（已确认部分；**不是** R2 冻结决定）
**关联意见**：Root `03-audit` A-003 / A-004
**关联信息项**：`I-PROTO-001` 保持 **collecting**；A-003 `F-001` / `F-002` 保持 **open / required**

**决定**：

1. **Q1**：D-TABLE 的多选批量不进入 MVP。D-ACT 只纳入与非批量 page/row action 及 request lifecycle 对齐的语义；依赖 D-TABLE 多选批量的 action/request 语义不在本次覆盖范围。
2. **Q3**：D-UPLOAD 保持 `exclude`。
3. **Q4**：R2 冻结前必须写入 D-COMP 的边界原则与初始 type 表；R3 可以在不改变该范围的前提下实现，任何扩域须另行决策，不得静默发生。
4. **Q5**：`scenarios` 不进入自动化门禁，只作为范例或手工验证路径。
5. **P-004.1**：用户选择在草案具备可核对修订后，补一次 R2 `I-PROTO-001` 同 scope self audit。
6. **Q2**：用户确认方向为“是”，即允许纳入一个或多个 2.6/2.7 扩展控件；但尚未列明实际控件、组件 type、结构验证或 fixture 子集。该方向**不**等同于把全部扩展控件纳入，也不改变当前 `I-PROTO-001=collecting`、`F-001=open` 或 R2 未冻结的状态。

**为什么**：

- 用户对 A-004 的书面回复确认了 Q1、Q3-Q5 和同 scope self audit 的选择；P-004 要求将这些裁决留痕。
- A-003 F-002 要求 D-ACT/D-TABLE 对批量语义有一致的覆盖与验证边界，Q1 已足以先收敛该项。
- A-003 F-001 要求可逐项核对的 D-COMP/D-FORM 边界。仅回答 Q2=是不足以确定哪些扩展控件和验证子集进入 MVP，不能据此伪造白名单或冻结。

**未选方案**：

- **把 Q2=是解释为全部 2.6/2.7 扩展控件**：用户未给出该范围，且会越过 MVP 最小边界。
- **把 Q2=是解释为任取一个控件并自行挑选**：控件选择、type 与验证要求仍属用户未确认事实。
- **因 Q1/Q3-Q5 已定而直接冻结 R2**：F-001 尚无可核对白名单和表单子集，P-005 门禁仍开。

**影响**：

- 草案按 Q1、Q3-Q5 修订 D-ACT/D-TABLE、fixture 与待确认表；`F-002` 仍等待同 scope self audit 后才可评估 `fixed`。
- `F-001` 保持 open，直至用户列明 Q2 纳入控件，并在草案中给出 D-COMP/D-FORM 初始 type 表和结构验证子集。
- Root `status: active`、R2=进行中、`I-PROTO-001=collecting` 与 `progress: 1/6` 不变；不更新 `goal-tree.md`。

**后续**：

1. 用户列明 Q2 要纳入的具体 2.6/2.7 控件（或明确全部），以及适用的 type / 验证预期。
2. 补齐 D-COMP/D-FORM 初始 type 表与 `component-format` / 相关结构验证子集，完成草案修订。
3. 执行用户要求的同 scope self audit；根据可核对证据再评估 F-001/F-002 的闭合与 R2 冻结。

## D-008 · 澄清 Q2：纳入全部 2.6 / 2.7 表单控件及其版本能力

**日期**：2026-07-31
**状态**：accepted（范围澄清；**不是** R2 冻结决定）
**关联意见**：Root `03-audit` A-003 / A-004
**关联信息项**：`I-PROTO-001` 仍为 **collecting**，待正式冻结决定

**决定**：

1. 用户确认：MVP D-FORM 纳入固定上游 `schema-ui-docs@v2.7.0` 的全部 2.6 与 2.7 表单控件能力。
2. 2.6 范围为：`textarea`、`switch`、`checkbox`、`radio`，以及 `select.props.mode: multiple`；页面须声明 `protocolVersion >= "2.6"` 与 `form.controls.extended`。
3. 2.7 范围为：`cascader`、`checkboxGroup`、`richText`、`password`，以及任一纳入字段的 `props.defaultValue`；页面须声明 `protocolVersion >= "2.7"` 与 `form.controls.advanced`。`defaultValue` 是属性能力，不是独立 component type。
4. Q1 仅排除 D-TABLE 的多选批量及其 action/request 语义；它**不**排除表单内 `select.props.mode: multiple`，两者的 payload 和动作语义不同。
5. D-COMP 初始 type 表与验证边界按修订草案 §5/§5.1 固化：它是当前 MVP 的最小白名单，不宣称支持完整 registry；任何新增 type 或扩域仍须另行决策。

**为什么**：

- 用户明确将 Q2 从“纳入任一扩展控件”澄清为“纳入 2.6/2.7 全部控件”。
- 固定提交的 migration、registry、node/page schema 与 fixture 内容给出了可逐项核对的 type、版本、capability 与结构验证依据，可补足 A-003 F-001 所要求的初表和验证边界。
- 将 D-TABLE 批量语义与表单多值 select 区分，避免 Q1 被误解为缩小用户刚确认的 2.6 控件范围。

**未选方案**：

- **只选其中一部分 2.6/2.7 控件**：与用户“全部控件”的书面确认相反。
- **将 `props.defaultValue` 当成独立 type**：与上游 2.7 迁移定义不符。
- **把全量 registry 一并纳入**：超出本目标的 MVP 初始白名单，也会违反 D-COMP 的 partial 边界。

**影响**：

- 修订草案已经将 2.6/2.7 类型、`defaultValue`、版本/capability 与验证子集写为可核对边界；`component-format` 的现有五个 format cases 只作为格式非强制转换验证，不能替代 type 白名单检查。
- Q5 不变：form-controls 场景保持范例/手工路径，不成为自动化门禁。
- A-005 将在同 scope 自审中复核 F-001/F-002 的修订；即便通过，`I-PROTO-001` 仍须由用户确认后另行冻结并转为 `verified`。

**后续**：

1. 执行同 scope self audit，核对草案与固定提交的类型、版本/capability、结构和 fixture 边界。
2. 若自审无开放 required finding，向用户提出 R2 正式冻结决定；未确认前不得改 `I-PROTO-001`、status 或 progress。

## D-009 · 正式冻结 R2 / I-PROTO-001 的 v0.1.3 覆盖边界

**日期**：2026-07-31
**状态**：accepted
**关联信息项**：`I-PROTO-001` → **verified**
**关联意见**：A-003 `F-001` / `F-002` 已由 A-005 以 `fixed` 合法闭合

**决定**：

1. 按用户本轮书面确认，正式冻结 [I-PROTO-001 覆盖表](attachments/I-PROTO-001-coverage-draft.md) **v0.1.3**，作为 VP-001 MVP 的 R2 覆盖基线。
2. 冻结范围以该表为准：7 个 `include` 域、4 个 `include-partial` 域及 `D-UPLOAD` 的 `exclude`，包括 D-COMP/D-FORM 白名单、D-ACT/D-TABLE 的非批量边界和 fixture 映射；不得将其扩大为完整 registry、全协议支持或未列 type。
3. `I-PROTO-001` 自 `collecting` 转为 **`verified`**，R2 标记为完成，Root 派生进度按 R1-R6 路线图重算为 **2/6**；Root `status` 仍为 `active`。
4. 本冻结只固定后续实施与验证的范围，**不**证明 R3-R5 已实施、已通过 conformance，或 VP 已可关门。`I-PROTO-002` / `I-PROTO-003` 继续分别阻断 R4 / R5。
5. 对覆盖子集的任何扩大、缩小或语义变更，必须追加新的 D-00N、修订覆盖表版本，并重新评估 `I-PROTO-002` / `I-PROTO-003`；不得静默改写 v0.1.3 基线。

**为什么**：

- 覆盖表规定的闭合条件是用户书面确认 → Root 冻结决策 → `I-PROTO-001` → `verified`；本轮指令满足该确认。
- A-005 同 scope self audit 为 `pass`，A-003 的 required `F-001` / `F-002` 已有可核对的 `fixed` 留痕；当前 scope 没有开放 required finding。
- Charter / VP 的协议 pin、工作区绑定和 Vision required 门禁均已核对；本决定仅实例化已确认的 MVP 子集，不改变愿景或协议源。

**未选方案**：

- **冻结完整协议或完整 registry**：超出 v0.1.3 的明确 include/include-partial/exclude 边界，也会违反既有 MVP 非目标。
- **继续保持 `collecting` 或接受 residual**：用户已明确要求正式冻结，且所需修订与自审证据已经齐备。
- **同时启动 R3-R5 实施**：本轮只确认 R2 范围；后续阶段仍须分别规划并满足 `I-PROTO-002` / `I-PROTO-003` 门禁。

**影响**：

- 允许以 v0.1.3 作为 R3-R5 的范围基线进行后续规划；不自动创建子目标或放行具体实现。
- 更新 Root 信息表、路线图、执行/审计响应和工作区 `goal-tree.md` 的同一事实投影。

**后续**：

1. 规划 R3 Admin 外壳与导航子目标。
2. 在 R4 前闭合 `I-PROTO-002`，在 R5 验收前闭合 `I-PROTO-003`。
