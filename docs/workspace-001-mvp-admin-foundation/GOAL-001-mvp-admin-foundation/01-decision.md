---
id: GOAL-001-mvp-admin-foundation
doc: decision
status: active
parent: null
created: 2026-07-31
updated: 2026-07-31
version: 0.4.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

权威表见 [00-meta.md](00-meta.md)「信息就绪与未知项」。摘要：

| ID | 级别 | 最晚阶段 | 状态 | 阻断 |
|----|------|----------|------|------|
| I-PROTO-001 | required | R2 方案冻结前 | **collecting**（草案已落盘，未冻结） | 未 verified 前不得冻结实施范围、不得主张完整协议支持 |
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
