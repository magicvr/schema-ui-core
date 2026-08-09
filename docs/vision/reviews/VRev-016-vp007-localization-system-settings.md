---
doc_type: vision-review
id: VRev-016
status: active
source: independent
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
parent: null
---

# VRev-016 · VP-007 多语种与系统设置产品化 · 意图 / 退出边界 / 对齐链（2026-08-09）

| 字段 | 值 |
|------|-----|
| source | independent |
| auditor | Grok · `/vision-audit` |
| scope | `VP-007-localization-and-system-settings`（`planned` v0.1.0）；用户关注：`vp-007` |
| audit_type | vision-plan |
| verdict | conditional |
| 建议 class | editorial |

## 范围与结论

只读核对：`docs/architecture/principles.md` **P-006**、`docs/vision/alignment.md`、`charter.md` `@0.2.0`、`plans/VP-007-localization-and-system-settings.md`（v0.1.0）、closed `VP-003`～`VP-006`、`roadmap.md`（v0.14.0）、`workspaces.md`、`revisions.md`（至 VR-011）、既有 `reviews.md`（至 VRev-015；仓库级 open required = 0）、以及可观察现状（`apps/web` 对 `titleKey`/`labelKey` 解析为字面串、`/api/branding` 仅 `siteTitle`+`logoUrl`、composition 测试确认 `mvp` 不启用 `admin.settings`）。未读 Goal 正文替代愿景证据；未把 `planned` 读成已交付；**未改** Charter / VP / Goal status。

**总判：conditional（1 open required · 2 open recommended）。**

VP-007 作为「同愿景下新纲领波次」的结构选型**合法且方向正确**：用户已书面确认将多语种与系统设置产品化合并为同一波次；不改 Charter 目的/边界/非目标；不重开 closed VP/区；`vision_ref` 精确匹配；`planned` 零区合法；组合索引与 Non-goals / Profile 边界（多语种双 Profile、Settings 编辑面仅 admin）与实现层可观察事实一致。信息门禁 `I-L10N-001`～`005` 与退出 #2 的协议不越权声明，整体优于「先做业务模块再补基架」的反模式。

但方向级退出 **#2** 将 Schema/表单覆盖钉在未枚举的「代表性」上，同时用「主流程不得依赖硬编码英文」作硬门禁却未给出可核对主流程分母——激活后与关门时存在主观放行空间，**足以阻断「方向已稳、可无修订激活」的宣称**。其余为 exit #5 服务端条件句残余路径、以及 Charter 协议来源段仍写 VP-005 `active` 的组合指针卫生（recommended）。

## 核对事实

| 核对项 | 结论 | 证据 |
|--------|------|------|
| 单愿景 | **pass** | 唯一 `status: active` Charter；`schema-ui-core-admin-foundation@0.2.0` |
| VP→Charter 机读 | **pass** | `vision_ref: schema-ui-core-admin-foundation@0.2.0` 精确匹配 |
| 语义对齐（抽样） | **pass（方向）** | 服务「可 fork 基架」目的 + 成功边界 3（产品化/主题）与 4–5（Profile/模块边界）的延伸；Non-goals 对齐 Charter（不重定义协议、不做业务产品、不热插拔）；未收缩已 closed 的协议/架构边界 |
| VP 最小完备（P-006 §6.5） | **pass（骨架）** | 意图、方向级退出 1–6、`vision_ref`、工作区绑定表（空）、关门占位、规划短史、信息门禁、建议阶段均在 |
| planned 零区 | **pass** | alignment §5 允许；`lead_workspace` 空；roadmap「0 区，待激活」一致 |
| 结构选型 | **pass** | 新可关门主题波次 → 新 VP（P-006 §6.6）；用户确认不修订 Charter、不吸收进 closed 区 |
| 前置关闭 | **pass** | VP-001～006 均 `closed`；roadmap 前置写 VP-003/004/005/006 |
| 组合编排同步 | **pass（索引）** | roadmap 顺序 7 = VP-007 `planned`；「无 active 交付 VP」与 VP status 一致；后续业务方向明示不得预支本 VP 能力 |
| Charter 指针卫生 | **fail → 见 F-V031** | 同文件「协议来源」仍写 VP-005 **active** 为现行焦点；「关系」节与 VR-011 已写无 active / VP-005 closed |
| 过早交付主张 | **无** | `planned`、0 工作区；信息项全 `open`；未把双语/四类设置写成已完成事实 |
| Profile / Settings 边界 | **pass（意图）** | 与 composition 测试「mvp disables optional settings」一致；禁止静默把 `admin.settings` 塞进默认 mvp |
| 协议字段锚点 | **pass（素材）** | 本地 `component-registry` / Manifest 校验已有 `titleKey`/`labelKey`；前端目前多将 key 当字面展示——exit #2「真解析」有可核对起点，非杜撰字段 |
| 既有 Vision required | **pass** | VRev-001～015 open required = 0；不继承到本 VP |
| 退出 #2 分母 | **fail → 见 F-V029** | 「代表性」+ 未枚举「主流程」使可验证覆盖分母不可客观钉死 |
| 退出 #5 条件句 | **weak → 见 F-V030** | 服务端 locale「有界可行」无「不可行」时的 residual / 退出降级写法 |

## 合理性总评（独立立场）

| 维度 | 立场 |
|------|------|
| 为何现在做 | **同意**：协议面（VP-006）与设计系统（VP-005）已闭；业务模块前统一 locale 运行时与系统设置产品面，可降低后续业务 VP 分叉成本；用户 2026-08-09 已确认合并波次。 |
| 波次粒度 | **方向合适、退出须收紧**：多语种 + Settings 四类在 Localization/Branding 上自然耦合；S0–S5 与 I-L10N-* 门禁合理。若 exit #2 分母不钉死，波次可能变成「译了几页 + Settings 能改标题」仍宣称产品化完成。 |
| 是否应用 Charter strategic | **否（对现行文本）**：不改目的/边界/非目标编号；属可 fork 产品化延伸。本审查不重开 strategic。 |
| 是否可保持 planned | **是**。**不建议**在 `F-V029` 未 editorial 闭合前宣称方向已稳或直接激活开区。 |
| 可行性 | **中等不确定、非方向非法**：Schema 文案策略（I-L10N-001）、公开 bootstrap（I-L10N-003）、错误 envelope（I-L10N-004）属实现/兼容风险，应由 S0 与用户 residual 消化；不因此 fail 意图本身。 |

## Findings

#### F-V029 · 退出 #2 将 Schema/表单钉在「代表性」，且「主流程」分母未枚举，关门不可客观判定

- level: `required`
- status: `open`
- severity: high
- impact: 激活后实施与关门可在「仅双语化登录/Shell + 任意 1～2 个示范页」与「fork 主路径全中英可完成」之间漂移；无法客观反驳「主流程已无硬编码英文依赖」；削弱本 VP「可直接用于 fork 的系统设置产品面 / 多语种基础」方向主张。
- finding: |
  1. 退出判据 2 固定面写清了登录、Shell、导航、通用 Renderer 状态/反馈、Settings——**正确且可核对**。
  2. 但对表单/表格与 Schema 驱动页面使用「**代表性**」限定，同时要求「不存在必须依赖硬编码英文才能完成的**主流程**」，却未给出：
     - 「主流程」最小清单（例如：登录 → 导航到代表性管理页 → 读/写代表性表单或表格 → 打开 Settings 改一项并看到生效）；
     - 「代表性」页面/Schema 的选取规则或冻结清单（可链到 I-PROTO-FULL-001 已 include 的范例页、或本 VP 书面钉死的 pageId/schema 集合）；
     - 缺失翻译可观察性与安全回退是否计入「主流程可完成」。
  3. Manifest `titleKey` / `labelKey`「真解析」是可机读子门闩，**不足以**单独证明 exit 2 的产品化覆盖。
  4. `I-L10N-001` 冻结的是 Schema 本地化**策略**（不越权协议），**不**自动冻结 exit 2 的覆盖分母；二者必须分开钉死，否则 S2 可在策略「已决策」下只译极少 Schema 面仍主张 exit 2。
  5. 与历史 F-V018 / F-V021 同类：方向正确但**覆盖分母可解释空间过大**时，不得宣称方向已稳。
- evidence:
  - `docs/vision/plans/VP-007-localization-and-system-settings.md` §方向级退出判据 2、§信息门禁 `I-L10N-001`、§建议实现阶段 S2、§交付范围「前端覆盖」
  - `docs/vision/plans/VP-007-localization-and-system-settings.md` §意图（「可直接用于 fork 项目的系统设置产品面」）
  - `apps/web/src/app/navigation.ts`、`App.tsx`（`titleKey`/`labelKey` 目前多作字面回退展示）
  - 对照：`docs/vision/reviews/VRev-011-vp005-design-system-ui-experience.md` F-V018；`VRev-012` F-V021（分母/partial 纪律先例）
- closure: |
  `/vision` **editorial**（建议激活前完成）：
  1. 为 exit 2 增加**可枚举最小证据面**（主流程步骤清单 + 固定 UI 面 + Schema/表单「代表性」选取规则或显式 page/schema 列表）。
  2. 明确：未列入最小证据面的文案面允许 residual，但须用户书面接受范围 + 复审触发；**不得**用「代表性」一词在无清单时关闭 exit 2。
  3. 与 exit 6 的「两语种 × 两 Profile × 登录前后…」矩阵交叉引用，使同一分母可被自动化或等价证据覆盖。
  4. 可选：把 `titleKey`/`labelKey` 真解析列为 exit 2 的显式子检查点（解析 + 文本 fallback + 缺失可观察）。
  不改 `vision_ref`、不要求 Charter strategic、不强制现在激活。
- 建议 class: `editorial`

#### F-V030 · 退出 #5 服务端 locale 以「成本有界」为条件，未规定不可行时的 residual / 退出降级

- level: `recommended`
- status: `open`
- severity: medium
- impact: S4 可因「成本无界」整段跳过服务端协商，却仍按字面关闭 exit 5（因前端保底已写）；或反过来在无 residual 留痕时被要求做完服务端才关门，范围争议。
- finding: |
  用户确认「有界可行时」后端按请求语种返回已编目消息——方向诚实。
  Exit 5 前半（稳定错误码 + 前端按码/key/参数本地化）是硬要求；后半服务端协商挂在「若成本有界」。
  正文**未**写明：当 I-L10N-004 结论为「成本无界 / 本波次不做服务端 locale」时，exit 5 是否仅以前端路径关闭，以及是否必须 `accepted-residual`（范围 + 复审触发）。条件句单独存在时，关闭纪律弱于本 VP 对 I-L10N-* 其他门禁的严肃程度。
- evidence:
  - `docs/vision/plans/VP-007-localization-and-system-settings.md` §退出 5、§信息门禁 `I-L10N-004`、§用户确认「有界可行时」
- closure: |
  `/vision` editorial：在 exit 5 或 I-L10N-004 关闭规则中写清二选一——(a) 服务端路径纳入本波次可验证证据；或 (b) 用户书面 residual（仅前端本地化 + 错误码契约保留），并声明不阻断 exit 5 的前端半边。
- 建议 class: `editorial`

#### F-V031 · Charter「协议来源」段仍宣称 VP-005 为 active 现行焦点，与关系节 / roadmap / VR-011 冲突

- level: `recommended`
- status: `open`
- severity: low
- impact: 读者（含后续激活 VP-007 时的对齐扫描）可能误判组合仍有 active 交付 VP，或忽略 planned VP-007 为下一意图；属组合指针卫生，**不是** VP-007 正文自相矛盾。
- finding: |
  `charter.md` 「与工作区 / VP 的关系」与 `revisions.md` VR-011、`roadmap.md` 一致写明：**无 active 交付 VP**；VP-005 **closed**。
  但同文件「协议来源」段末句仍写：现行组合交付焦点为 VP-005（**active**）。
  单文件内指针冲突；VP-007 已在 roadmap 落盘 planned，Charter 关系节亦未索引该 planned 波次（可选卫生，非门禁）。
- evidence:
  - `docs/vision/charter.md` 「协议来源」段 vs 「与工作区 / VP 的关系」段
  - `docs/vision/revisions.md` VR-011；`docs/vision/roadmap.md` 当前交付意图
- closure: |
  `/vision` editorial：删除或改写「协议来源」中 VP-005 active 句；可选在关系节增加 planned VP-007 指针（不改 `vision_id@version`、无 re-align）。
- 建议 class: `editorial`

## 声明

本意见不修改 Charter / VP / Goal status；required finding 的响应由 `/vision` 协调，实施工作交 `/govern`。原 verdict 与 finding 原文不得改写；闭合响应追加在本报告。

## 响应（对独立意见 · VRev-016）

| date | actor | summary |
|------|-------|---------|
| 2026-08-09 | `/vision` | 采纳 `conditional` / `editorial`，保留原 verdict 与 finding 原文。**F-V029 → `fixed`**：VP-007 `0.1.1` 新增“最小可枚举证据面”，以固定 UI、S0 双 Profile Runtime Manifest `pages[]`/schema 并集、M1～M4 主流程钉死 exit #2 分母；exit #6 强制复用 `zh-CN/en-US × mvp/admin × 匿名/认证` 同一矩阵，未覆盖面只能经用户书面 residual。**F-V030 → `fixed`**：exit #5 / I-L10N-004 明确二选一关闭规则——服务端 locale 实施证据，或用户书面 `accepted-residual`；前端本地化 + 稳定错误码始终为硬门禁，本轮未接受 residual。**F-V031 → `fixed`**：Charter 删除 VP-005 active 旧指针，明确当前无 active、planned VP-007，并以 VR-012 记录 editorial；`vision_id@version` 保持 `@0.2.0`，无 re-align。VP-007 保持 `planned`、0 区，未激活、未开工作区。本 scope **0 open required、0 open recommended**；激活仍是后续独立用户裁决。 |
