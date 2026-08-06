---
doc_type: vision-alignment
title: 愿景对齐契约与门禁
status: active
created: 2026-07-28
updated: 2026-08-06
version: 0.7.0
parent: null
---

# 对齐契约 · Vision → VP → Workspace

本文件是愿景体系的**规则权威**。与实例声明冲突时，以本文件与 [charter.md](charter.md) 为准；`consumer-checklist.md` 不得宽于或严于本文件。  
元原则见 [../architecture/principles.md](../architecture/principles.md) **P-006**。

## 0. 不变量与完整安装

### 0.1 单愿景制

1. **每一个项目**（一次完整治理安装）**有且仅有一个**现行愿景：`docs/vision/charter.md` 且 `status: active`。
2. **禁止多愿景**（并行多个 active 北极星或多个争用的 `vision_id`）。
3. 不同项目可以有不同愿景——那是**独立治理实例**，不是本项目内多愿景。
4. **换代**：旧 Charter 仅 `superseded`；现行对齐只认唯一 active；历史 `vision_ref` 可只读。
5. 多工作区、多 VP **共一愿景**。

### 0.2 完整安装与冷启动

1. 完整安装文件集以本节 **Minimal Complete Install（MUST）** 为**唯一权威**；[consumer-checklist.md](consumer-checklist.md)、[standalone-bootstrap.md](../standalone-bootstrap.md)、[principles P-006 §6.2](../architecture/principles.md#62-完整安装与冷启动) **必须同表**，禁止三处各写「建议/必含」分裂定义。缺任一 MUST 行 = **不完整安装**。
2. **法定顺序（严格串行）**：最小完备 Charter → 首个 VP 落盘 → 工作区 + Root（`plan_refs` + `primary_plan`）→ 区内纲领路线图与子目标。
3. **缺 Charter**：仅允许**引导补齐**；拒绝非引导开区、推进、放行、关门；须报告不完整安装。
4. **legacy** 无 vision：不完整；补齐前不得新开显式区或放行/关门（只读与引导除外）。

#### Minimal Complete Install（MUST） vs Recommended

> **权威表**。更新本表后须同步 checklist 勾选与 standalone 核对项。
> 「完整独立启用」= 下表全部 MUST 已存在且冷启动顺序合法；**不是**「仅有 architecture + templates」。

| 层级 | 路径 / 条件 | 级别 | 说明 |
|------|-------------|------|------|
| 规则入口 | 根 `AGENTS.md`（或项目声明的等价 AI 规则） | **MUST** | 操作摘要；全文原则仍以 architecture 为准 |
| 文档入口 | `docs/README.md` | **MUST** | 核心文档索引 |
| 方法论 | `docs/architecture/principles.md` | **MUST** | P-001～P-006 全文 |
| 方法论 | `docs/architecture/workspace-protocol.md` | **MUST** | 工作区/资料协议 |
| 模板 | `docs/templates/goal-folder/`（五件套 + `attachments/`） | **MUST** | canonical 目标模板 |
| 模板 | `docs/templates/workspace-context.md` | **MUST** | 工作区页模板 |
| 模板 | `docs/templates/vision/charter.md`、`vision-plan.md`、`reviews-index.md`、`review.md` | **MUST** | 冷启动与 Vision Review 复制源 |
| 契约（若分发消费适配器） | `docs/contracts/` 消费契约文件 | **MUST**（有 Skills/Web 分发时） | 纯文档-only 仓可无，但不得假装已装适配器契约 |
| 愿景规则 | `docs/vision/alignment.md` | **MUST** | 本文件；规则权威 |
| 愿景入口 | `docs/vision/README.md` | **MUST** | 目录地图与硬边界 |
| 愿景实例 | `docs/vision/charter.md`（`status: active`） | **MUST** | 单愿景；缺 = 不完整 |
| 愿景树 | `docs/vision/roadmap.md` | **MUST** | 组合编排索引（可极简，但文件必须存在） |
| 愿景树 | `docs/vision/revisions.md` | **MUST** | Charter 修订台账（可极简） |
| 愿景树 | `docs/vision/reviews.md` | **MUST** | Vision Review 稳定索引；报告在 `reviews/VRev-NNN-*.md`（有条目时创建目录） |
| 愿景树 | `docs/vision/workspaces.md` | **MUST** | 工作区贡献图（可极简） |
| 愿景树 | `docs/vision/consumer-checklist.md` | **MUST** | 与本表一致的操作勾选 |
| 意图 | 至少一个 `docs/vision/plans/VP-*.md` | **MUST**（开区前） | `vision_ref` 精确匹配 Charter |
| 工作区 | 显式 `docs/workspace-<NNN>-<slug>/workspace.md` | **MUST**（开区后） | 含必填 `plan_refs` / `primary_plan` |
| 目标 | 工作区根 `goal-tree.md` + Root 五件套 | **MUST**（开区后） | Root `parent: null` |
| 方法论（可选扩展） | `docs/architecture/overview.md`、`directory-layout.md` 等 | Recommended | 增强可读性，不替代 MUST |
| 实例 dogfood | 他仓过程树、本仓历史 GOAL 附件 | 勿复制 | 不是完整安装条件 |

**半安装**：仅有 architecture/templates、无 Charter 或无上表愿景树 MUST 文件 → 只可读原则，**不得**记为完整独立启用通过，也不得非引导推进/放行/关门。

### 0.3 命名

| 名称 | 含义 |
|------|------|
| **组合编排** | 愿景级 VP 索引与波次（[roadmap.md](roadmap.md)），非 progress% |
| **意图** | 已落盘 [plans/VP-*.md](plans/)；草案不可作 `primary_plan` |
| **纲领路线图** | Root/大目标 P-001 阶段 |
| **阶段计划** | 目标内方案/实施安排（非树节点） |
| **Vision Review** | 愿景层审视；≠ Goal Audit |

## 1. 类型与禁止

| 类型 | 允许 status | 禁止 |
|------|-------------|------|
| Charter (`doc_type: vision-charter`) | `active` \| `superseded` | Goal 的 `done` / `draft` / `blocked` / `cancelled`；progress%；goal-tree；并行第二 active |
| Vision Plan (`doc_type: vision-plan`) | `planned` \| `active` \| `closed` \| `abandoned` | Goal 的 `done` 作 VP status；完整五件套；progress% 权威 |
| 工作区目标 | 既有 Goal status | 把 vision/VP 目录当目标父节点；跨区 `parent` |

Charter **没有 canonical `draft` 状态**：尚不满足最小完备或尚未获用户确认的草案只能留在会话/提案中，不得占用现行 `docs/vision/charter.md`。`active` Charter 可以显式登记战略假设/未知；若某项影响“方向已稳”，在 verified 或合规 residual 前不得作该宣称（见 P-006 §6.5）。

愿景体系**不是**第二套目标状态源；不汇总各区 progress，不关闭 Goal finding。

## 2. 对齐递归与三层链

```text
Charter: vision_id@version          ← 链源头（单愿景）
    ↑  VP.vision_ref 精确匹配（不做 semver 范围）
Vision Plan: VP-NNN-slug            ← 意图
    ↑  workspace/Root.plan_refs 含该 VP；primary_plan 必填且 ∈ plan_refs
Workspace + Root Goal
    ↑  子目标 parent = 父目标完整 id
子目标 …
```

- Root **不强制**再抄 `vision_ref`；经 VP 间接对齐 Charter。
- 子目标默认继承 Root 的规划语境；**有界偏离**须 P-004 留痕；改边界则升级修订 VP/Charter。
- **残差不自动继承**：父级 `accepted-residual` / `overruled` 不自动扩到子树。
- **复述**：允许短 `serves_summary` + 链接；禁止目标内第二套愿景边界。

### 2.1 最小可检查对齐

| 层 | 机读 | 语义 |
|----|------|------|
| 子→父 | `parent` 完整 id 存在且同区 | 不与父边界/非目标明显冲突 |
| Root/区→VP | `plan_refs`、`primary_plan` 必填；VP 文件存在 | `serves_summary` 服务该意图 |
| VP→Charter | `vision_ref` = `{vision_id}@{version}` | 意图在 Charter 边界内 |

失败 → **fail closed**（引导 re-align 除外）。

## 3. 工作区声明（`workspace.md` frontmatter）

| 字段 | 要求 |
|------|------|
| `vision_role` | `primary` \| `delivery` |
| `plan_refs` | **必填**，至少一个 VP id；多个用逗号分隔 |
| `primary_plan` | **必填**，且必须 ∈ `plan_refs`；对应 `docs/vision/plans/<id>.md` |

- **无 plan opt-out**：任何工作区都不得省略 `plan_refs` / `primary_plan`；`vision_role` 仅允许 `primary` / `delivery`。
- 至多一个工作区 `vision_role: primary`（与 [workspaces.md](workspaces.md) 一致）。

### 3.1 Primary 声明冲突裁决

Primary 可能出现在三处：`workspace.md` 的 `vision_role: primary`、[workspaces.md](workspaces.md) 的 `role: primary`、Charter 的 `primary_workspace`。

| 情形 | 行为 |
|------|------|
| 三处一致 | 通过 |
| 仅一处声称 primary，其它未声明或为空 | 以**已声明**处为准，并应在下一维护回合补齐另外两处 |
| 两处或以上**互相矛盾**（不同 `workspace_id`） | **fail closed**：不得推进受影响的新建 Root/放行/关门；展示冲突；按 P-004 等用户裁决后留痕再改 |
| Charter `primary_workspace` 指向不存在的工作区 | fail closed，直至修正 Charter 或创建/声明该区 |

权威顺序（仅用于**修复建议**，不能静默覆盖用户已确认的矛盾）：`alignment` 本文件规则 → 用户书面裁决 → 再改 Charter / workspaces.md / workspace.md。

## 4. Root Goal 声明

Root `00-meta.md` 应含与 workspace 一致的 `plan_refs`、`primary_plan`，以及简短 `serves_summary`（frontmatter 或「愿景对齐」节）。  
`primary_plan` 必须能解析为 `docs/vision/plans/<id>.md`（id 与文件名一致）。

## 5. VP 与工作区绑定

| VP status | 工作区绑定 |
|-----------|------------|
| `planned` | 允许 0 个工作区 |
| `active` | 期望 ≥1；若为 0，见下「空转」规则 |
| `closed` | 保留历史绑定；默认不接新区，除非 reopen + 用户确认 |
| `abandoned` | 不要求绑定 |

- 一规划 : 0..N 工作区；一工作区 : 1..N 规划（须标 `primary_plan` 焦点）。
- **多于一个**工作区绑定同一 VP 时 **`lead_workspace` 必填**；VP 关门提案默认由 **lead** 侧发起，support 证据须链接，经**用户确认**。单区时可省略 lead 或等于该区。

### 5.1 `active` VP 零工作区（空转）

`status: active` 且绑定工作区数为 0 时：

1. 编排器**必须告警**，并询问用户：挂接工作区 / 改回 `planned` / 接受有时限的空转。
2. **空转宽限**：自 VP 标为 `active` 或自上次「零区复核」起 **14 个日历日**（以 `updated` 或决策留痕日期较晚者为准）。宽限内可扫描与规划，但**不得**把该 VP 当作已在推进的交付证据。
3. **超过宽限**仍无工作区、且无用户书面「继续空转」记录（须含下一复核日 ≤ 再 14 日）：对该 VP 相关的新建挂接以外的**放行/关门** fail closed，直到挂区、降为 `planned`/`abandoned`，或留下新的有界空转接受。
4. 「长期空转」即指超过上述宽限且无合规留痕。

## 6. 门禁时机（fail closed）

在下列时机，Skills / 编排器 / 适配器应校验对齐链；失败则不得假装推进或关门：

1. 完整安装判定 / 冷启动（缺 Charter → 仅引导）  
2. 新建工作区或新建 Root（须已有 Charter + 可挂接 VP）  
3. 新建子目标（继承与语义对齐检查）  
4. 推进影响成功边界/非目标的阶段  
5. 目标 close-out 前  
6. Charter **strategic** 修订后：受影响 VP 与挂接工作区须 re-align（见 §8 宽阻断）  
7. VP 关门前（区证据 + lead 规则）  
8. 相关 **Vision Review** 的 required 意见未合法闭合时（开区 / VP 关门 / 宣称方向已稳）

失败模式：缺 Charter；缺 `plan_refs` / `primary_plan`；`primary_plan` 不在列表中；VP 文件缺失；`vision_ref` 与 charter 版本不一致；非法 `vision_role`；primary 无规划；Charter/VP 非法 status；待 re-align 宽阻断生效中。

## 7. 规划关门（轻量）

1. 退出判据方向满足；**证据链接**指向工作区目标的 done / 有界结项路径。  
2. 允许**有界 closed**：residual 必须点名到具体 workspace / goal id。  
3. 无区证据不得将 VP 标为 `closed`。  
4. 多区时由 **lead** 发起关门提案 + 用户确认。  
5. **禁止**为 VP 建立 Goal 五件套或独立目标 `03-audit` 台账替代区内审计。

## 8. Charter 修订与宽阻断

| class | 含义 | version | 工作区 |
|-------|------|---------|--------|
| `editorial` | 措辞、链接，不改方向/边界/非目标 | 可补丁级 | 不强制 re-align |
| `strategic` | 目的、边界、非目标、原则优先级 | 至少 minor | impact 所列 VP/区必须 re-align |

流程：更新 charter → 追加 [revisions.md](revisions.md) → 更新 VP `vision_ref` → 刷新工作区/Root 声明；Charter 初建或 strategic 后完成 [Vision Review](reviews.md)（可为 self）。

**strategic 后宽阻断**（受影响范围 re-align 完成前）：

- **禁止**：新建子目标、放行、关门、非引导开区  
- **允许**：只读扫描、引导 re-align、补写 revisions/Review 响应  

## 9. Vision Review

权威台账：[reviews.md](reviews.md) 稳定索引 + `reviews/VRev-NNN-<slug>.md` 平铺报告目录。编号 **`VRev-00N`**（与 revisions 的 `VR-` 修订号区分）。

| 项 | 约定 |
|----|------|
| 条目头 | `source`（`self` \| `independent`）、日期、scope、`verdict`（pass \| conditional \| fail）；建议 auditor |
| 正文 | findings；required/建议；建议修订 class |
| 默认效力 | **不**直接改 Charter/VP status |
| required 闭合 | 与 P-003 同构：`fixed` / `accepted-residual` / `user-overruled` + 留痕 |
| 强制时机 | Charter 初建；每次 `strategic` 后（P-006） |
| 长文 | 报告保留摘要/verdict/findings；更长证据可链愿景层附件 |

### 9.1 可扩展台账布局

1. `reviews.md` 保留 frontmatter、使用约定、当前 `open required` 投影与条目链接，不内联新 VRev 正文。
2. 一条正式意见一个 `reviews/VRev-NNN-<slug>.md`；目录单层平铺，文件 frontmatter 的 `id` 必须与文件名前缀一致，self / independent 共用编号序列。
3. 索引与报告共同构成唯一正式台账；写入必须同时创建报告和更新索引。编号扫描合并 legacy inline 与目录报告后取最大值 +1。
4. finding 响应由 `/vision` 追加在原 VRev 报告中，保留原 verdict 与 finding 原文；索引 `open required` 随合法闭合证据更新。
5. legacy inline VRev 继续有效。兼容 reader 合并读取；达到 32 KiB、800 行、12 条记录任一阈值后，下一条必须写入目录。迁移不得重编号、改变历史语义或丢失响应。
6. 全新安装从第一条 VRev 起使用目录报告；模板见 `docs/templates/vision/reviews-index.md` 与 `review.md`。

### 9.2 工具入口（与 Goal 台账分界）

| 入口 | 写什么 | 不写什么 |
|------|--------|----------|
| **`/vision`** | self Vision Review → VRev 报告 + 索引；在原报告追加 finding 响应；Charter / VP / 组合编排 / re-align 决策 | Goal `03-audit`；不推进子目标执行；不改写原审计结论 |
| **`/vision-audit`** | independent Vision Review → VRev 报告 + `reviews.md` 索引（`source: independent`） | 不改 Charter / VP / Goal status；不自行闭合 finding |
| **`/audit`** | Goal `03-audit` independent | **禁止**写入 `docs/vision/reviews.md` |
| **`/govern`** | 实现层推进与 Goal finding 响应 | 无 Charter 时不得假装完整推进 |

独立愿景审视**必须**走 `/vision-audit`；Goal 交叉审计走 `/audit`。两台账不得混写。

## 10. 结构选型（摘要）

详见 P-006 §6.6。要点：改源头 → Charter；新波次 → VP；独立树/隔离 → 新工作区（仍挂 VP）；同 Root → 子目标；高不确定探索 → P-005 有界信息收集阶段/目标。

## 11. 引用格式

- 愿景：`{vision_id}@{version}`，例 `vision-goal-governance@0.1.0`
- 规划：`VP-NNN-slug`，路径 `docs/vision/plans/VP-NNN-slug.md`
- Review：`VRev-00N`

## 12. 非目标（本契约 v1）

- 多愿景；双段目标编号；pillars；`vision_ref` 的 semver 范围匹配  
- Web UI / CT 专为愿景的写入流程（只读发现与校验即可）  
- 将执行流水或 progress 写入 vision 目录  
- 自动跳过 Vision Review 的算法；硬编码目标嵌套层数  
- 与 VP 并列的独立「Intent」文档类型  
- 任何工作区省略 `plan_refs` / `primary_plan` 的 opt-out
