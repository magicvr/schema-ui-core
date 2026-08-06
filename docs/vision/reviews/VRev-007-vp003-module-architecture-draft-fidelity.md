---
doc_type: vision-review
id: VRev-007
status: active
source: independent
created: 2026-08-04
updated: 2026-08-04
version: 0.1.0
parent: null
---

# VRev-007 · VP-003 相对 MODULE-ARCHITECTURE-DRAFT 的意图保真独立审视（2026-08-04）

| 字段 | 值 |
|------|-----|
| source | independent |
| auditor | Grok 4.5 · `/vision-audit` |
| scope | `VP-003-modular-admin-architecture`；根目录 `MODULE-ARCHITECTURE-DRAFT.md` 预设意图；现行 `docs/architecture/module-architecture.md`；Charter `schema-ui-core-admin-foundation@0.2.0`；组合编排与「中间形态≠终态」漂移风险 |
| audit_type | vision-plan |
| verdict | pass |
| 建议 class | editorial |

### 范围与结论

用户要求独立核对：**新 VP-003 是否符合 `MODULE-ARCHITECTURE-DRAFT.md` 的预设意图**，并专门防「意图曲解 / 把中间形态当终态 / 目标飘逸」。

只读证据：`MODULE-ARCHITECTURE-DRAFT.md`、`docs/vision/plans/VP-003-modular-admin-architecture.md`、`docs/architecture/module-architecture.md`、`charter.md` `@0.2.0`、`roadmap.md`、`workspaces.md`、`revisions.md`（VR-004）、`reviews.md`（至 VRev-006）、`dual-track-contract.md` 现行性说明。未用 Goal 正文替代愿景证据；未把 planned VP 读成已交付。

**总判：pass（0 open required）**。VP-003 的**终态意图**与 draft 的目标架构愿景一致，且在「禁止试点充当终态」上**不低于** draft；未发现把 R3 试点、文档落盘或局部模块化写成退出边界的曲解。下列 recommended findings 不推翻方向，但应在 R1 契约冻结 / R3 试点门禁前由 `/vision` 或 `/govern` 收紧措辞，避免实施期软化 draft 已拍板的门闩。

### 核对：draft 预设意图 → VP-003 / 正式架构

| draft 预设意图 | VP-003 / `module-architecture.md` 落点 | 保真判定 |
|----------------|----------------------------------------|----------|
| §1 用「单主线 + 模块注册 + 启动配置」替代双线代码树 | 退出判据 1；已接受决策「以单主线与 Profile 替代双线」；roadmap 将 dual-track 标为历史 | **保真**；且诚实声明不伪造不存在的分支删除 |
| §2.1 / P-1～P-6 薄内核、模块、Profile=装配列表、前端随协议、Manifest 后端驱动 | 退出判据 2/4/7；Charter 成功边界 4–5；正式架构 §1–§5 | **保真** |
| D-1 轻量 DI（建议 Fx/Wire） | Uber Fx；明确不采用已归档 Wire；Fx 类型不进入模块公共 API | **可接受细化**（仍在建议 A 内） |
| D-2 标准管理界面红线内前端零改动 | 同一前端 build + 后端聚合；增减标准模块不改 Renderer/Shell 中央路径 | **保真** |
| D-3 模块 9 项贡献点（6 项必须 + 3 项可选） | 正式架构能力表写「按需实现」；VP 写「统一模块描述与可选能力契约」 | **措辞弱化风险** → `F-V010` |
| D-4 禁用不删表、模块自带迁移 | 退出判据 3；禁用/退役不删表；已编译一方模块迁移独立于启用状态 | **保真且加强**（全局台账 / tombstone / fail-closed） |
| D-5 运行时配置 + 启动固定 Profile | 静态编译候选 + 启动时 Profile；无热插拔 | **保真** |
| D-6 运行时后端聚合 Manifest | `/.well-known/schema-ui/app-manifest.json`；生产禁静默静态兜底 | **保真**（路径命名细化） |
| D-7 / §5 试点先行，再全量迁移 | 路线图 R3 有界试点 → R4 全量一方模块 → R6 旧路径移除 | **保真**；退出边界**未**缩到试点 |
| §4.5 / §5.1 试点必须切除 4 个旧架构病灶 + V-1～V-4 | R3 写启停/依赖/冲突/双 Profile 等，**未完整复述** 5 项交付物与 4 病灶手术门闩 | **门闩变薄** → `F-V011` |
| §6 非目标：无热插拔、无 `.so`、不把试点定义成「只写新模块不碰旧代码」 | VP Non-goals 对齐；且明文禁止把试点/能启动/局部模块化当终态证据 | **保真且加强** |
| §8 正式 ADR → 试点 → 指南 → 存量迁移 | `module-architecture.md` 已固化；R1–R6 覆盖后续；VP `planned` 未绑工作区 | **保真**；实施仍属未来 `/govern` |

### 核对：「中间形态当终态 / 意图曲解 / 目标飘逸」

**不成立的漂移主张（本轮未采纳）：**

1. **未把 R3 当退出判据**：VP 正文写「完整终态，不是 Activity/Settings 试点的妥协版本」「即使试点通过也不得据此关闭本 VP」；R3 列「通过仅允许继续扩迁」；Non-goals 第 5 条禁止试点/文档/局部模块化充当终态证据。roadmap 同步标注试点非退出边界。
2. **未把「架构文档落盘」当架构完成**：VP `status: planned`、零 workspace 绑定（`workspaces.md`）；正式架构写 draft 仅为评议输入；VRev-006 已声明非实施证据——本轮独立复认。
3. **未把生产运维闭环偷换成业务产品**：Non-goals 排除订单/钱包等业务域与协议扩张；roadmap 方向 4 仍要求单独建 VP。运维/数据/双 Profile 验收是 draft 模块化可 fork 前提的生产化补全，继承 VP-002 基线，**不构成**北极星换成业务产品。
4. **未逆向恢复双线终态**：Charter `@0.2.0` 与 dual-track 历史化一致；退出判据 1/6 要求平行代码线与旧中央注册路径退出。

**可接受的工程细化（不记为 finding）：**  
operationlog / activity 拆分、Settings 纳入试点（呼应 draft §4.5 host 特例）、迁移「已编译即参与全局顺序」、Lifecycle/Observability、Manifest 路径命名、仅采用 Fx。它们未改写「单主线 + Profile + 后端聚合 + 全量迁完才算完」的终点。

**仍须警惕的实施期软化（recommended，非 required）：**  
见 `F-V010`、`F-V011`。它们不会让当前「方向已稳 / planned 合法」失效，但若不在 R1/R3 门禁前闭合，**有可能在执行中**把 draft 已否决的「新模块自嗨、老病灶仍在」或「最小能力集模块」偷渡为阶段通过。

### Findings

#### F-V010 · 正式模块能力表相对 draft D-3「必须贡献点」措辞偏弱

- level: `recommended`
- status: `fixed`
- closed_at: `2026-08-04`
- closed_by: `/vision` · V6 响应（用户指令：接着用 `/vision` 做修正）
- severity: medium
- impact: R1「模块 API / 能力契约」冻结；一方模块验收是否允许缺 Schema/Manifest/迁移等核心贡献
- finding: |
  用户已按 draft 会前清单接受 D-3：每个模块通过 DI 注册 **9 项**贡献点，其中路由、Schema、权限、导航、Manifest 片段、迁移为 **必须**；DependsOn / Config / Seeds 为可选（DependsOn 强烈建议）。
  现行 `docs/architecture/module-architecture.md` §2 将能力一律写成「**按需实现，未实现即表示模块不贡献该能力**」；VP-003 退出判据 2 亦强调「可选能力契约」。
  若按字面执行，一方 Admin 功能模块可能在缺 Schema/Manifest/迁移等 draft **必须**项时仍被解释为「合法模块」。这属于对 draft 模块自治颗粒度的**可执行口径弱化**，不是终态缩成试点，但是契约漂移入口。
- evidence:
  - `MODULE-ARCHITECTURE-DRAFT.md` §3 D-3（必须/可选表）
  - `docs/architecture/module-architecture.md` §2
  - `docs/vision/plans/VP-003-modular-admin-architecture.md` 退出判据 2、已接受决策「能力可选」
- closure: |
  由 `/vision` editorial（或 R1 冻结决策）明确：**一方标准 Admin 功能模块**必须具备 draft 的 6 项核心贡献（HTTP/Schema/Authorization/Navigation/Manifest/Persistence）；DependsOn/Config/Seeds（及 Lifecycle/Observability）按需；「按需」不得被读成「核心六项亦可永久缺省」。可在 `module-architecture.md` 与/或 VP-003 增加对照表，无需改 Charter 战略。
- resolution: |
  **editorial fixed**（非 strategic）：
  1. `docs/architecture/module-architecture.md` → `1.0.1`：§2 拆为 §2.1 核心六项**必须** + §2.2 按需能力；明确「按需不得覆盖核心六项」；横切 `operationlog` 可经显式说明豁免不适用 UI 项。
  2. `VP-003` → `0.1.1`：退出判据 2 与「已接受的架构决策 · API 边界」同步核心六项必须口径。
  **未改** Charter、`vision_ref`、VP `status: planned`、工作区绑定。

#### F-V011 · R3 试点门闩未完整继承 draft §4.5/§5 的「病灶切除 + V-1～V-4」

- level: `recommended`
- status: `fixed`
- closed_at: `2026-08-04`
- closed_by: `/vision` · V6 响应（用户指令：接着用 `/vision` 做修正）
- severity: medium
- impact: VP-003 激活后 R3 阶段放行；是否允许「写出 activity/settings 模块」在未证明 Kernel 切口时进入 R4
- finding: |
  draft §4.5 / §5.1 / §5.2 / 会前清单第 9 条明确：试点**不是**「新增 Activity 就完事」，必须切除 4 个旧架构病灶（中心化 Register、全局 Schema fixtures、静态 Manifest、Host 特例），完成 5 项交付物，并通过 V-1～V-4（启停、前端零改动、Schema 贡献等）；卡在 Kernel 改造则先加固，**不盲目全量迁移**。
  VP-003 R3 写的是 operationlog/activity 拆分与 settings、启停/依赖/冲突/配置事件/双 Profile/升级等——方向正确且含 Settings（呼应 host 特例），但**未点名**将 draft 的 4 病灶拆除与 V-1～V-4 列为 R3 **通过门闩**。
  这不会把 R3 抬成 VP 终态（退出判据 6 与 R6 仍要求旧路径删除），却增加「软 R3 通过 → 带着泥潭进 R4」的执行漂移风险；与 draft「试点必须动手术」的预设意图不完全同构。
- evidence:
  - `MODULE-ARCHITECTURE-DRAFT.md` §4.5、§5.1、§5.2、§5.3、§7 项 9
  - `docs/vision/plans/VP-003-modular-admin-architecture.md` 迭代路线图 R3、退出判据 6、Non-goals
  - `docs/architecture/module-architecture.md` §8（试点不是终态——正确，但未恢复 draft 试点硬门闩全文）
- closure: |
  在 VP-003 R3 行或未来工作区阶段计划中**显式继承** draft 试点硬门闩：4 病灶切除清单 + 5 交付物 + V-1～V-4；并保留 draft §5.3「未通过则加固 Kernel、不盲目 R4」。闭合不要求改退出判据终态本身。
- resolution: |
  **editorial fixed**：
  1. `VP-003` `0.1.1`：新增 **R3 通过门闩**（A 五项交付 / B 四病灶 / C V-1～V-4 / D 决策门）；路线图表 R3 行指向该节；已接受决策增加「试点」行。
  2. `module-architecture.md` `1.0.1` §8：摘要 R3 硬门闩并链回 VP-003。
  终态退出判据与「试点不关 VP」边界**未缩减**。

### 声明

本意见不修改 Charter / VP / Goal status；required finding 的响应由 `/vision` 协调，实施工作交 `/govern`。审计当时 **无 required** finding；`F-V010` / `F-V011` 为 recommended，现已由 `/vision` editorial 响应为 `fixed`。

### 门禁含义

- Vision Review **required 仍为 0 open**；VRev-007 recommended 亦已 **0 open**。
- 对用户关注点的独立结论：**无「终态被试点替换」类 required 缺陷**；契约/试点门闩措辞缺口已由 editorial 修正。
- 本 `pass` 与后续 `/vision` 响应 **均不是** VP 激活、开区、实施放行或关门证据。

### 响应

| date | actor | summary |
|------|-------|---------|
| 2026-08-04 | `/vision` | 用户指令「接着用 `/vision` 做修正」：采纳 VRev-007 `pass` / `editorial`。**F-V010 → `fixed`**：`module-architecture.md` v1.0.1 §2.1 核心六项必须 + §2.2 按需；VP-003 v0.1.1 退出判据 2 / API 边界同步。**F-V011 → `fixed`**：VP-003 写入 R3 通过门闩（5 交付 + 4 病灶 + V-1～V-4 + 失败则不盲目 R4）；架构 §8 摘要。未改 Charter、未激活 VP、未开区。 |

---

> **迁移说明（2026-08-07）**：本报告自 legacy inline `docs/vision/reviews.md` 原样拆出，编号与历史结论未改；相对链接已按 `reviews/` 目录深度调整。
