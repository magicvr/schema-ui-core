---
doc_type: vision-review
id: VRev-012
status: active
source: independent
created: 2026-08-08
updated: 2026-08-08
version: 0.1.1
parent: null
---

# VRev-012 · VP-006 整份 v2.7.0 契约意图 / 退出边界 / 组合焦点（2026-08-08）

| 字段 | 值 |
|------|-----|
| source | independent |
| auditor | Grok 4.5 · `/vision-audit` |
| scope | `VP-006-full-protocol-contract-v2-7-0`（`planned` v0.1.0）；用户关注：`vp-006` |
| audit_type | vision-plan |
| verdict | conditional |
| 建议 class | editorial |

### 范围与结论

只读核对：`docs/architecture/principles.md` **P-006**、`docs/vision/alignment.md`、`charter.md` `@0.2.0`（含 2026-08-08 协议目标语义与 H-001 分列）、`plans/VP-006-full-protocol-contract-v2-7-0.md`（v0.1.0）、closed `VP-003`/`VP-004`、planned `VP-005`（仅组合门闩关系）、`roadmap.md`、`workspaces.md`、`revisions.md`（至 VR-007）、既有 `reviews.md`（至 VRev-011；`F-V018` 仍 open 且仅作用于 VP-005）、`protocol-inventory-v2.7.0.md`、`I-PROTO-001` 覆盖表 `v0.1.3`（路径与 disposition 结构）。未读 Goal 正文替代愿景证据；未把 `planned` 读成已交付；**未改** Charter / VP / Goal status。

**总判：conditional（1 open required · 2 open recommended）。**

VP-006 作为「同愿景下纠正组合焦点的新纲领波次」结构选型**合法且方向正确**：Charter 成功边界 1 + 2026-08-08 用户「整份契约」裁决与 H-001 分列 ③ 一致；`I-PROTO-001 v0.1.3` 被明确降为历史 MVP 升版起点而非终态上界；硬阻塞 VP-005 实施与 roadmap / Charter / VRev-011 响应一致；机读 `vision_ref` 精确；`planned` 零区合法；退出含过程可关门与回归诚实声明，整体优于曾被 F-V018 卡住的 VP-005 范围措辞。

但方向级退出 **#1** 将 `include-partial` 与 `include` 并列为「默认覆盖路径」，且未把「大面积明确排除子面」抬到 residual 级用户裁决，存在在 S1 再钉一个更大 partial 子集、仍宣称「整份契约覆盖表已决策」的可解释空间——这足以阻断「方向已稳、可无修订激活」的宣称。其余为覆盖表权威落点与 inventory 入口卫生问题（recommended）。

### 核对事实

| 核对项 | 结论 | 证据 |
|--------|------|------|
| 单愿景 | **pass** | 唯一 `status: active` Charter；`schema-ui-core-admin-foundation@0.2.0` |
| VP→Charter 机读 | **pass** | `vision_ref: schema-ui-core-admin-foundation@0.2.0` 精确匹配 |
| 语义对齐（抽样） | **pass（方向）** | 意图收口成功边界 1 + pin 整份契约；继承边界 4–5 / 非目标（业务产品、协议重定义、热插拔）；不收缩仍生效的边界 |
| VP 最小完备（P-006 §6.5） | **pass（骨架）** | 意图、方向级退出 1–6、`vision_ref`、工作区绑定表（空）、关门占位、规划短史、交付形态定名表均在 |
| planned 零区 | **pass** | alignment §5 允许；`lead_workspace: null`；roadmap 一致 |
| 结构选型 | **pass** | 新可关门主题波次 → 新 VP（P-006 §6.6）；禁止吸收进 closed VP-003/004 工作区 |
| 前置关闭 | **pass** | VP-001～004 `closed`；架构/playbook 可继承 |
| 组合编排同步 | **pass** | roadmap 顺序 5=VP-006 焦点、6=VP-005 实施冻结；Charter「无 active 交付 VP」+ 焦点指针一致；VR-007 editorial 与用户裁决一致 |
| 与 I-PROTO / inventory 关系 | **pass（意图）** | inventory + pin 为权威输入；v0.1.3 仅升版起点；禁止静默改写 v0.1.3 语义；禁止「MVP 曾排除」默许 |
| 与 VP-005 / F-V018 | **pass（门闩）** | VP-006 硬阻塞 VP-005 至 closed；`F-V018` **不**阻断本 VP 的 planned/激活讨论（reviews 投影已写明） |
| 过早交付主张 | **无** | `planned`、0 工作区；未把全量覆盖/实现写成已完成事实 |
| 退出 #1 纪律 | **fail → 见 F-V021** | `include-partial` 与默认覆盖路径并列表述；大面积 partial 子面排除未强制 residual 级裁决 |

### 合理性总评（独立立场）

| 维度 | 立场 |
|------|------|
| 为何现在做 | **同意**：组合曾长期把 MVP 子集当协议成功代理；用户书面纠正后，在架构/playbook 已闭、视觉波次之前先收口整份契约，符合协议优先与可 fork 基架承诺。 |
| 波次粒度 | **方向合适、门闩须收紧**：S0–S5 与「可分批但关门须满足方向退出」合理；允许有界 residual exclude 是诚实治理。但若不收紧 partial 纪律，波次可能再次变成「新版 partial 表 + 部分实现」。 |
| 是否应用 Charter strategic | **否（对现行文本）**：成功边界 1 原文未改编号；VR-007 按 editorial 澄清「子集非终态」可接受。本审查不重开 strategic。 |
| 是否可保持 planned | **是**。**不建议**在 `F-V021` 未 editorial 闭合前宣称方向已稳或直接激活开区。 |
| 可行性 | **高不确定、非方向非法**：全量 registry/fixtures/前后端缺口相对 v0.1.3 很大，属实现风险，应由激活后 S0 差距盘点与用户 residual 裁决消化；不因此 fail 意图本身。 |

### Findings

#### F-V021 · 退出 #1 将 `include-partial` 并入默认覆盖路径，且未强制「大面积子面排除」走 residual 级裁决

- level: `required`
- status: `fixed`
- severity: high
- impact: 激活后 S1 可产出「逐项有 disposition」的覆盖表，却对多域保留与 I-PROTO 同构的大块 `include-partial`+明确排除子面（批量、upload、半个 registry 等），在字面上仍宣称「整份契约覆盖表已决策 / 可关门推进」，实质再钉一个更大子集；与标题「整份契约」、用户 2026-08-08 裁决及 exit 文「不得以又钉更小子集替代」的意图冲突，且关门时难以客观反驳。
- finding: |
  1. Exit 1 写：对能力域与 component registry **逐项** disposition（`include` / `include-partial`+边界 / `exclude`+书面理由）；**默认目标**是可验证 `include`（**或**带显式边界的 `include-partial`）覆盖上游契约承诺面。
  2. 对 `exclude` 明确要求用户书面有界残余（范围 + 复审触发），并禁止以「MVP 曾经排除」默许——**正确**。
  3. 但对 `include-partial` **未**要求 residual 级用户裁决，也未区分两类 partial：
     - **保真/边角 partial**（能力已纳入，边角语义或未测 fixture 子集有界）；
     - **范围 partial**（整域或主要子面长期排除，与历史 I-PROTO 对 D-ACT/D-TABLE/D-COMP/D-FORM/D-UPLOAD 的模式同构）。
  4. 历史 `I-PROTO-001 v0.1.3` 证明：大量能力可在「表已冻结」下以 include-partial + 明确排除子面长期停留；若 VP-006 允许同一模式作为默认覆盖路径，则「整份契约」可被读成「全表有行」而非「契约承诺面默认可验证 include」。
  5. Exit 后段禁止「又钉更小子集替代整份契约目标」——与 exit 1 前段 partial 并列表述之间，缺少可操作的判定规则（何时 partial 仍算整份契约、何时必须升格 exclude+residual 或用户书面接受范围收缩）。
- evidence:
  - `docs/vision/plans/VP-006-full-protocol-contract-v2-7-0.md` §方向级退出判据 1、§建议实现阶段 S1、§意图用户裁决
  - `docs/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md` v0.1.3（include-partial / exclude 先例）
  - `docs/vision/charter.md` 成功边界 1；H-001 分列 ③ open → VP-006
  - `docs/vision/roadmap.md` 组合门闩 2（MVP 子集不是终态成功条件）
- closure: |
  `/vision` **editorial**（建议激活前完成）：
  1. 改写 exit 1：**默认 disposition** = 对 inventory/registry/fixture 承诺面的可验证 `include`；`exclude` 仅用户书面 residual（已有）。
  2. 将 `include-partial` **收窄定义**为：能力已纳入前提下的保真/边角/未测子集边界，**不得**用 include-partial 表达「整域或主要子面不打算做」。
  3. 凡相当于范围收缩的 partial（示例：整域 upload、批量 selection/ADR0022 全排除、registry 大支整支排除），必须二选一写清：**(a)** `exclude` + residual 字段，或 **(b)** 用户书面接受的有界范围收缩（范围 + 复审触发），**不得**仅靠 S1 规划便利落 partial。
  4. 可选：在 S1 检查点增加「与 v0.1.3 差集中，转为 include 的计数 / 仍 residual 的清单」可审计摘要，防止静默子集化。
  不改 `vision_ref`、不要求 Charter strategic。
- resolution: |
  **fixed**（2026-08-08 `/vision`）：VP-006 → `0.1.1` exit 1 重写为默认 `include`；`include-partial` 仅保真/边角；范围收缩强制 exclude 或用户书面 residual；S1 增加 v0.1.3 差集可审计摘要。未改 `vision_ref`、未激活。
- 建议 class: `editorial`

#### F-V022 · 升版覆盖表的权威落点与版本身份未一锤定音

- level: `recommended`
- status: `fixed`
- severity: medium
- impact: 激活后可能在 workspace-001 附件上「升版改写」历史冻结、或散落多份覆盖表导致发现路径与关门证据分母不清。
- finding: |
  Exit 1 要求「新版本号」覆盖表 + 新决策，并禁止静默改写 `I-PROTO-001 v0.1.3` 文件语义——方向正确。
  但正文**未**定名：权威文件路径/id（例如新 `I-PROTO-00N-…` 附件 vs 独立 architecture/vision 索引）、与 v0.1.3 的只读继承关系（链接 vs 复制）、以及「现行覆盖表」的单一发现入口。
  交付形态定名表说明了「是什么/不是什么」，未给出可引用的权威落点句。对 planned 不阻断；激活前 editorial 或 S1 方案冻结应补齐。
- evidence:
  - `docs/vision/plans/VP-006-full-protocol-contract-v2-7-0.md` §交付形态定名、exit 1、继承边界 `I-PROTO-001`
  - 对照：`I-PROTO-001-coverage-draft.md` 现权威路径在 workspace-001
- closure: |
  `/vision` editorial 或激活后 Root S1：写清权威路径 + id + 「v0.1.3 只读回归对照、不就地改语义」；并保证 overview/QUICKSTART 或 Root 决策可发现。
- 建议 class: `editorial`
- resolution: |
  **fixed**（2026-08-08 `/vision`）：VP-006 → `0.1.1` 新增「覆盖表权威落点」：历史只读 `I-PROTO-001 v0.1.3`；现行权威 id `I-PROTO-FULL-001`（lead 区 Root attachments + 新决策）；发现入口意图期=本 VP+inventory，S1 后=Root 决策+overview/QUICKSTART。

#### F-V023 · inventory 入口仍以 VP-001 / MVP 为服务对象，易误导 S0 差集默认值

- level: `recommended`
- status: `fixed`
- severity: low
- impact: S0 盘点时把 `mvp_candidate` 或 inventory §3.1「纳入 MVP 的范例」读成 VP-006 上界；或继续把 inventory `serves: VP-001` 当唯一消费方。
- finding: |
  VP-006 正确声明 inventory + pin 为权威输入，且不得以 MVP 排除默许。
  但 `protocol-inventory-v2.7.0.md` frontmatter 仍 `serves: VP-001-mvp-admin-foundation`；§3.1 范例职责写「每个**纳入 MVP 的**能力域」；§4/§5 仍以 VP-001 覆盖冻结与 I-PROTO-001 信息项为中心，**未**指针到 VP-006 为整份契约收口意图。
  这不是 VP-006 正文自相矛盾，而是**组合入口卫生**缺口：实现层易从 inventory 单独读出错误默认。
- evidence:
  - `docs/vision/protocol-inventory-v2.7.0.md` frontmatter `serves`、§3.1、§4、§5
  - `docs/vision/plans/VP-006-full-protocol-contract-v2-7-0.md` 继承边界 / exit 1
- closure: |
  `/vision` editorial（可与激活同步）：inventory 增加「全量能力清单；MVP 子集见 I-PROTO v0.1.3；整份契约收口见 VP-006」指针；`serves` 可扩写或改为 multi-serve 说明；§3.1 范例句与 VP-006 对齐（纳入能力域而非仅 MVP）。不要求改 pin 或重提清单。
- 建议 class: `editorial`
- resolution: |
  **fixed**（2026-08-08 `/vision`）：`protocol-inventory-v2.7.0.md` → `0.1.2` multi-serve（VP-001 历史 + VP-006 整份契约）；用途分列 + §3.1 范例改为「纳入覆盖表」；§4/§5 指针 `I-PROTO-FULL-001` / VP-006。未改 pin、未重提清单。

### 不构成 fail / 不新开额外 required 的诚实边界

1. 本 `conditional` **不是**对「做整份 v2.7.0 契约」方向的否决；结构选型、用户裁决、组合门闩与 Charter 边界 1 对齐成立。
2. 实现体量大、保真债与未实现 type 并存，属 S0 后实现风险，不单独升格 required。
3. `F-V018`（VP-005）仍 open **不**计入本 VP 的 open required，也不阻断本 VP 的 planned 讨论；VP-005 解冻仍依赖 VP-006 closed + 其后对 F-V018 等的响应。
4. inventory 未 vendor 上游全文、conformance reference runner 被 manifest exclude——VP-006 exit 4 已要求按域登记验证入口，不强制上游 runner；不新开 finding。
5. Charter `primary_workspace` 仍为 workspace-001 / VP-001 历史 primary——与 closed 波次 + 未来 delivery 新区模式一致；不要求本 VP 改 primary。
6. 独立 Vision Review **不**激活 VP、**不**开区、**不**闭合 finding。

### 声明

本意见不修改 Charter / VP / Goal status 或 progress；required finding 的响应由 `/vision` 协调；实现层执行仍交 `/govern`。独立 Vision Review **不**自行闭合 finding。

### 门禁含义

- Vision Review 对本 scope（审计当时）：**open required = 1**（`F-V021`）；响应后见下节。
- recommended open（审计当时）：`F-V022`、`F-V023`。
- 仓库级既有 open required：`F-V018`（VRev-011 · 仅 VP-005）仍 open，与本条并行。
- **允许**：保持 `planned`；只读规划与 editorial 修订讨论；组合焦点继续指向 VP-006 而非 VP-005 实施。
- **禁止（在 F-V021 合法闭合前）**：以本审查宣称「方向已稳」；将本 VP 作为已就绪 `primary_plan` 直接开区/激活交付；将 exit 1 读成授权用大面积 `include-partial` 再钉子集并宣称整份契约已决策。
- 闭合 `F-V021` 后：仍须 `/vision` + 用户确认方可 `active` 与 `/govern` 开区；本 `conditional` 不自动变为实施放行。

### 响应（对独立意见 · VRev-012）

| date | actor | summary |
|------|-------|---------|
| 2026-08-08 | `/vision` | 采纳 VRev-012 `conditional` / `editorial`。**F-V021 → `fixed`**：VP-006 → `0.1.1` 重写 exit 1（默认 `include`；`include-partial` 仅保真/边角；范围收缩 → exclude 或用户 residual；S1 差集可审计摘要）。**F-V022 → `fixed`**：覆盖表权威落点 `I-PROTO-FULL-001` + v0.1.3 只读 + 发现入口。**F-V023 → `fixed`**：inventory → `0.1.2` multi-serve 与 VP-006 / 全量指针。未改 Charter、`vision_ref`、VP status；**未激活、未开区**。本 scope **0 open required、0 open recommended**。仓库级仍余 **F-V018**（仅阻断 VP-005）。激活与工作区 slug 仍须用户确认后 `/vision` + `/govern`。 |
