---
id: GOAL-005-r3-admin-shell-navigation
doc: audit
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.3.0
---

# 审计 · GOAL-005

> 本文件是目标的唯一正式意见台账（P-003）。self / independent 意见共用 `A-00N` 序列。

## 信息就绪核对

| 核对项 | 状态 | 备注 |
|--------|------|------|
| 影响本 scope 的 I-00N | `I-005-001` 至 `I-005-005` 均 open/required | 见 [00-meta.md](00-meta.md) |
| 到期 required 是否已 verified / residual | 当前未进入方案冻结或实施门禁；不能放行 | 未有用户书面 residual；五项均 open；A-001 F-001 已由 A-002 `fixed` |
| 固定资料引用 | 部分固定（愿景/清单级 pin） | artifact `2.7.0` + commit `ca9e5fe…` 已登记；`shared_materials_catalog: none`，无工作区共享资料引用；本地 schema/fixture 接入仍待确认 |
| 当前实现证据 | 未开始 | `apps/web` 仍为 R1 单页占位；无 router / manifest loader / navigation shell 实现 |

## 意见台账索引

| ID | 日期 | source | scope | verdict | 开放 required |
|----|------|--------|-------|---------|---------------|
| A-001 | 2026-07-31 | independent | R3 目标定义 / 规划立项与信息门禁就绪性 | conditional | 1（F-001） |
| A-002 | 2026-07-31 | self（编排响应） | A-001 F-001 修正与 recommended 响应 | conditional | 0（F-001 fixed；`I-005-*` 仍 open） |

## 当前开放门禁

- `I-005-001` 至 `I-005-005` 在其最晚阶段前必须完成验证，或按 P-004 留下有范围和复审触发的用户书面 residual。
- A-001 **F-001** 已由 A-002 以 `fixed` 闭合：I 表与 D-002 现均要求在方案冻结前处理 `I-005-001` 至 `I-005-005`；这不等同于任何信息项已 verified，也不放行方案冻结。
- 在相关 required 信息项和审计意见处理完成前，不得把 R3 方案冻结、实现或 `status: done` 写成已放行事实。
- 本目标不处理父目标的 `I-PROTO-002` / `I-PROTO-003`；它们仍分别阻断 R4 实施和 R5 验收/关门。
- 父目标 `I-PROTO-004`（vendor vs pin）仍 open（non-blocking 于 Root），现已通过 `I-005-001` 显式关联；关闭该项时仍须记录校验方式和失败边界。

---

## A-001 · R3 目标定义与规划立项独立审计（2026-07-31）

- **source**：independent
- **auditor**：GitHub Copilot（`/audit` · independent）
- **类型**：goal-definition
- **scope**：`workspace-001-mvp-admin-foundation` / `GOAL-005-r3-admin-shell-navigation` 的目标定义、范围边界、P-001 路线图、P-005 信息台账、D-001～D-003、立项执行事实；对照 Root R3、I-PROTO-001（D-APP）、protocol-inventory、Charter/VP 与 `apps/web` 现状。**不审**实施完成、不审 R4/R5、不写 Vision Review。
- **verdict**：**conditional**
- **完整意见**：本节即全文（未另附 attachments）

### 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-001-mvp-admin-foundation`（`root_goal=GOAL-001-mvp-admin-foundation`，`canonical_scope` 已校验，`shared_materials_catalog: none`） |
| 被审目标 | `GOAL-005-r3-admin-shell-navigation`（`parent=GOAL-001-mvp-admin-foundation`，`status=active`） |
| 审计区间 | 立项/规划文档 + 可复核代码现状（2026-07-31） |
| 明确不在 scope | 方案冻结正文、R3 实现、验收关门、父目标 R4/R5 门禁改写、`docs/vision/reviews.md` |

### 成果（有证据）

| 主张 | 证据 | 核对 |
|------|------|------|
| 五件套与 `attachments/` 已建齐 | 目标目录：`00-meta`/`01-decision`/`02-execution`/`03-audit`/`attachments/` | 通过 |
| 挂接 Root、对齐 R3 纲领名 | [goal-tree.md](../goal-tree.md)；Root 路线图 R3「Admin 外壳与导航」 | 通过 |
| 范围纳入 D-APP 外壳/导航，排除 R4 权限、R5 Renderer/范例、完整协议支持 | [00-meta.md](00-meta.md) 范围；[I-PROTO-001 v0.1.3](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) `D-APP=include`（fixtures `app-manifest`/`app-navigation`）；inventory `D-APP` primary=装载与导航壳 | 通过 |
| P-001 四段路线图已写（信息→方案冻结→实施→验证关门） | [00-meta.md](00-meta.md) 高层路线图；[01-decision.md](01-decision.md) D-002 | 通过 |
| P-005 已登记 `I-005-001`…`005` 为 required/open，未伪造成 verified | [00-meta.md](00-meta.md) 信息表；[02-execution.md](02-execution.md) | 通过 |
| 上游 pin 与证据边界已声明，且未把 excluded runner 当兼容证明 | D-003；inventory §2.3 excluded 说明；覆盖表 excluded 提示 | 通过 |
| 立项未改 `apps/web`、未产生 shell/router 实现事实 | [02-execution.md](02-execution.md)；`apps/web/src/main.tsx` 无 router；`package.json` 无 react-router；`App.tsx` 仍标 R1 单页占位；`src/host`/`protocol` 仅为预留 README | 通过 |
| Root `progress: 2/6` 未被本目标立项抬升；R3 未写成完成 | [goal-tree.md](../goal-tree.md)；GOAL-005 成功标准全未勾选 | 通过 |
| 工作区共享资料未误用 | `workspace.md` catalog=`none`；GOAL-005 无伪造 material/sha256 引用 | 通过 |

### 对照成功标准（目标定义阶段）

| 成功标准（00-meta） | 定义阶段评价 |
|--------------------|--------------|
| 信息项 verified 或有界 residual | **未满足（预期）**：五项均 open；台账存在且可追踪 |
| manifest 装载用冻结 schema/版本/子集 + 失败边界 | **未满足（预期）**：属方案/实施；I-005-001 open |
| shell/导航/默认·fallback·active-route 可复核 | **未满足（预期）**：无 router/shell；I-005-002…004 open |
| app-manifest / app-navigation 验证路径已执行 | **未满足（预期）**：仅有上游路径锚点，无本地执行证据 |
| R4/R5/完整协议仍边界外；无开放 required finding 才关门 | **定义层边界清楚**；关门条件尚未适用 |

定义阶段**不要求**勾选实施类成功标准；要求的是：边界诚实、门禁可执行、无假完成。前三点已大体满足；门禁表内部一致性见 F-001。

### Findings

#### F-001 · required · 中 · I 表最晚阶段与 D-002 冻结范围不一致

- **问题**：D-002 要求在方案冻结前处理 `I-005-001`～`I-005-005`，并把「默认/fallback/active-route 语义」写入冻结内容；但 [00-meta](00-meta.md) 将 `I-005-003`、`I-005-004` 的**最晚需要阶段**标为「实施前」，仅 `I-005-001/002/005` 为「方案冻结前」。编排器若只读 I 表，可能在 003/004 仍 open 时放行方案冻结；若只读 D-002，则冻结门槛高于 I 表。
- **证据**：
  - [01-decision.md](01-decision.md) D-002：「先处理 I-005-001 至 I-005-005，再冻结 … 默认/fallback/active-route 语义」
  - [00-meta.md](00-meta.md) 信息表：`I-005-003`/`I-005-004` → 最晚「实施前」
- **影响门禁**：方案冻结（P-005）
- **建议闭合**：
  1. **fixed（推荐）**：将 `I-005-003`/`I-005-004` 最晚阶段改为「方案冻结前」，与 D-002 对齐；或
  2. **fixed（替代）**：改写 D-002，明确方案冻结只强制 001/002/005，003/004 可延至实施前并写清可接受的冻结粒度；或
  3. 用户书面 residual（须范围 + 复审触发）——不推荐作为默认。

#### F-002 · recommended · 中 · 未显式挂钩父级 `I-PROTO-004`（vendor vs pin）

- **问题**：R3 验收依赖 `app-manifest` schema 与 navigation fixtures 的可核对接入；覆盖表 §4 写明 `I-PROTO-004`「建议 R3 前定」。GOAL-005 的 `I-005-001` 仅写「本地接入仍待确认」，未登记与 Root `I-PROTO-004` 的依赖/联动，易在方案中默认“路径已知=可校验”。
- **证据**：
  - [I-PROTO-001 覆盖表 §4](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)
  - Root [00-meta](../GOAL-001-mvp-admin-foundation/00-meta.md) `I-PROTO-004=open`
  - GOAL-005 `I-005-001` 证据列
- **建议**：在 meta/decision 增加对 `I-PROTO-004` 的显式依赖说明（仍可由 Root 决策）；关闭 `I-005-001` 时写明 vendor 路径或 pin 远程校验方式及失败边界。不把 `I-PROTO-004` 误升为本目标独立 required，除非用户决定下沉。

#### F-003 · recommended · 低 · Root 纲领 R3 状态文案仍为「未开始」

- **问题**：子目标已 `active` 并完成规划立项，Root 路线图表仍写 R3「未开始」，读者可能理解为“尚未进入 R3”，与 GOAL-005 备注「已进入规划阶段」并列表述摩擦。
- **证据**：Root [00-meta 纲领路线图](../GOAL-001-mvp-admin-foundation/00-meta.md) R3 行；GOAL-005 [00-meta 备注](00-meta.md)
- **建议**：由 `/govern` 在 Root 将 R3 标为「规划中 / 信息收集」或等价文案，**不要**标完成，**不要**改 `progress: 2/6`。本 independent 意见不改 Root。

### 必改项汇总

| # | 级别 | 摘要 | 建议 |
|---|------|------|------|
| F-001 | **required** | D-002 与 I-005-003/004 最晚阶段不一致 | 对齐 I 表或改写 D-002 后再谈方案冻结 |
| F-002 | recommended | 显式挂钩 `I-PROTO-004` / 本地校验接入策略 | 方案冻结前澄清 |
| F-003 | recommended | Root R3 状态文案与子目标规划态对齐 | `/govern` 更新 Root 路线图备注 |

### 与既有意见的异同

- 本目标此前**无** self/independent `A-00N`；文件头「信息就绪核对」为立项脚手架，不构成正式审计意见。
- 与 Root R2 意见传统一致：区分清单/pin、覆盖冻结与实现完成；本审确认 GOAL-005 **未**把 `I-PROTO-001=verified` 误写成 R3 已交付。

### 结论 + 建议给编排器/用户的下一步

**结论**：GOAL-005 作为 R3 **目标定义/规划立项**整体诚实、边界清楚、与 D-APP/R3 纲领和代码现状一致，**无假完成**；但信息门禁表与 D-002 存在 **required** 一致性缺口，故 verdict = **conditional**。不可无条件进入方案冻结或实施。

**建议 `/govern` 输入**（示例）：

```text
/govern 响应 workspace-001 GOAL-005 A-001：
1) fixed F-001：将 I-005-003/004 最晚阶段改为「方案冻结前」（或改写 D-002 明确分流）
2) 采纳 F-002：在 00-meta/01-decision 挂钩 Root I-PROTO-004
3) 可选 F-003：Root 纲领 R3 改为「规划中」
4) 然后推进 I-005-001…005 信息收集（勿实施代码）
```

### 声明

本意见 `source: independent`，**不修改**目标 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree 状态列。响应、finding 闭合与阶段推进由 **`/govern`** 处理。

---

## A-002 · A-001 编排响应（2026-07-31）

- **source**：self（编排响应；不是 independent 审计）
- **auditor**：Codex（`/govern`）
- **类型**：response
- **scope**：响应 A-001 的 F-001 required 一致性缺口，以及 F-002/F-003 recommended；不审 R3 实现、方案冻结或关门。
- **verdict**：**conditional**
- **完整意见**：本节即全文（未另附 attachments）

### 响应哪些意见 / Findings

| 原意见 | 处理 | 说明 |
|--------|------|------|
| A-001 F-001 · required | `fixed` | 采用审计推荐路径：`I-005-003` / `I-005-004` 的最晚需要阶段统一为「方案冻结前」，D-002 无需改写。 |
| A-001 F-002 · recommended | `fixed` | `I-005-001` 与 Root `I-PROTO-004` 的依赖和关闭要求已显式登记；后者仍为 Root 的 open/non-blocking 项。 |
| A-001 F-003 · recommended | `fixed` | Root 路线图 R3 由「未开始」改为「规划中」；未改 Root `progress: 2/6`。 |

### 关闭证据

| Finding / 信息项 | 状态 | 证据路径 |
|------------------|------|----------|
| A-001 F-001 | `fixed` | [00-meta.md](00-meta.md) 的 `I-005-003` / `I-005-004` 最晚需要阶段；[01-decision.md](01-decision.md) D-004；[02-execution.md](02-execution.md) 本次时间线。 |
| A-001 F-002 | `fixed` | [00-meta.md](00-meta.md) 的 `I-005-001` 与表后依赖说明；[01-decision.md](01-decision.md) D-004。 |
| A-001 F-003 | `fixed` | [Root 00-meta.md](../GOAL-001-mvp-admin-foundation/00-meta.md) R3 路线图行；[02-execution.md](02-execution.md) 本次时间线。 |

### 仍开放项与 P-004

- `I-005-001` 至 `I-005-005` 仍均为 `required` / `open`；本次未收集或验证任何契约事实，故不放行方案冻结、实施、验收或关门。
- Root `I-PROTO-004` 仍为 `open` / `non-blocking`；它只作为 `I-005-001` 的显式工程策略依赖，不构成已验证事实。
- A-001 是同 scope 的 independent 意见，当前没有同 scope 的 self 审计。本 A-002 只是编排响应，不能被当作 self 审计或静默跳过自审；在提议基于独立意见放行方案冻结前，必须由用户按 P-004.1 决定是否需要自审。

### 结论 + 建议下一步

**结论**：A-001 的唯一 required finding F-001 已合法以 `fixed` 闭合，信息表与 D-002 的方案冻结门禁现已一致。R3 仍处于信息收集前的规划阶段，不能因本响应而写成已冻结、已实施或已完成。

**建议下一步**：先收集并记录 `I-005-001` 至 `I-005-005` 的可核对事实；在任何方案冻结提议前，请用户明确决定是否需要同 scope 自审。

---

## 当前结论与下一步

- 一句话结论：R3 的信息表与 D-002 现已一致，A-001 F-001 已闭合；五项 required 信息仍 open，方案冻结、实施和关门均未放行。
- 建议下一步：信息收集 `I-005-*`（勿实施代码）；在提出方案冻结前先完成 P-004.1 的同 scope 自审选择。
- **声明**：上表索引以外的聊天讨论不构成正式意见。
