---
doc_type: vision-plan
id: VP-035-foundation-architecture-health
title: 基架架构健康评估与路线图重述
status: active
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-035-foundation-architecture-health
created: 2026-09-09
updated: 2026-09-09
version: 0.2.0
parent: null
---

# VP-035 · 基架架构健康评估与路线图重述

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-09-09 · v0.2.0 · 用户指令激活 · lead `workspace-035-foundation-architecture-health`） |
| lead_workspace | `workspace-035-foundation-architecture-health`（2026-09-09 开区） |
| Vision required | 计划阶段 = [VRev-086](../reviews/VRev-086-vp035-foundation-architecture-health-planned.md) self `pass`；激活就绪 = [VRev-087](../reviews/VRev-087-vp035-foundation-architecture-health-activation.md) self `pass`（0 required · 架构类 freshness PASS `f2044cf3`→`5c341ec7`） |
| 组合位置 | **架构分支** · 基架交付波收口后的有界评估。产出 = as-built 对照 + 有界业界对照分类 + 下一版总路线图草案。不替代 VP-009/VP-010 持续程序 |

## 意图

基架交付波（VP-013～034 所覆盖的端口、渠道、分发与 Admin 增量）已经收口。当前组合层仍停在过期的「现状锚点」和未冻结的 A0–A7 建议顺序上；各 VP residual 也没有一张总账。

本意图做一次**有界的架构健康评估**，再把评估结果交给 `/vision` 做总路线图 editorial 重述。评估有两路输入，缺一不可：

1. **内部 as-built**：对照 [module-architecture.md](../../architecture/module-architecture.md)、内核端口、组合根、Profile 与已交付 RT-*，盘点实现与文档是否一致。
2. **有界业界对照**：用事先冻死的参照集给缺口分类（现在修 / 仍 gated / 接受残余 / 明确不做）。业界实践是分类输入，**不是**决策源头，不得用来推翻 Charter 非目标。

本 VP **不是**无限优化程序，**不是**符合性长期程序（那是 VP-010），**不是**安全扫描（那是 VP-009）。

## 用户已裁决（2026-09-09 · `/vision` 会话书面）

| 项 | 裁决 |
|----|------|
| 结构 | **新 VP + 新 delivery 工作区**；不塞进 VP-010 当普通波次；不改 Charter |
| 业界对照 | **要**，但是有界输入；每条必须落成「业界常见做法 → 本仓现状 → 保持 / 现在修 / 仍 gated / 明确不做」 |
| 参照集 | 只盯四类：模块化单体与组合根；基础设施端口（内存默认、外部实现 gated）；Schema 驱动 Admin 贡献；同进程基座（H-002） |
| 路线图 | 先评估再 `/vision` editorial 重写；禁止先按「主流实践」改路线图再回头找代码 |
| 实现 | 无证据的优化不开工；Redis / MQ / 搜索引擎 / 多实例 / K8s 仍 gated，评估里最多登记 |

## 首波冻结（退出分母）

| 项 | 本 VP 交付 | 不进本 VP |
|----|-----------|-----------|
| 对照矩阵 | 内核、组合根、Provider 六项、已交付端口（Store / Object / Cache / RateLimiter / EventBus / Mail / Observability / Shutdown）、Profile 默认集、architecture 文档 vs `apps/api` + `apps/web` | 业务域模块产品化、Admin 体验增强实现（Command Palette 等） |
| 缺口登记 | 每条缺口：证据路径、分类（现在修 / 仍 gated / 接受残余 / 明确不做）、影响的路线图行 | 把分类写成已经开工的优化 |
| 业界对照 | 下表四类参照；公开文档/稳定框架惯例即可，不要求付费调研或现场访谈 | 无限扩充参照集；用某家公司博客推翻 Charter |
| 路线图草案 | 现状锚点、A 序列、Admin/架构/业务下一拍、residual 总账；交 `/vision` editorial | 本 VP 直接改 Charter；本 VP 在用户确认前把草案写成已冻结路线图 |
| 文档卫生 | 至少刷新 `docs/architecture/overview.md` 现时节，以及 roadmap 过期锚点/A 序列（与草案同一事务或紧随 `/vision` 冻结） | 重写全部历史 VP 正文 |
| 代码整改 | 仅当分类为「现在修」且范围小、可在本 VP 证据闭环内完成；默认倾向**登记后另立 VP/波次** | Redis/MQ/K8s/ORM/第三库；重开已 closed VP |

### 有界业界参照集（I-035-002 冻结）

| # | 类别 | 允许对照 | 禁止据此推翻 |
|---|------|----------|--------------|
| 1 | 模块化单体 + 组合根 | Go + Fx / .NET Generic Host / Spring 启动装配 | 拆微服务、运行时插件市场 |
| 2 | 基础设施端口 | 内存默认缓存/限流/进程内总线，以及它们何时才升到 Redis/MQ | 「主流都用 Redis/Kafka，所以现在就上」 |
| 3 | Schema 驱动 Admin | 协议优先的页面/权限/导航贡献模型 | 把本仓做成某个 CMS 克隆 |
| 4 | 同进程基座 | 同进程扩 C 端模块、共享端口 | 默认拆 Admin / API / Worker 三进程 |

## 非目标

- 不改 Charter 目的、成功边界、非目标或 `vision_id@version`
- 不重开 VP-003/004/008/009/010 或任何已 closed 交付 VP
- 不实现 Redis、外部队列、搜索引擎、多实例、K8s、ORM、第三数据库
- 不把 Command Palette、文件扫描策略、组织/部门当成本 VP 的实现范围
- 不把业界对照写成 Charter strategic 修订；若对照结论要动非目标，本 VP **停住**并交 `/vision` strategic
- 不替代 VP-009 安全扫描或 VP-010 例行符合性波次

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003 / VP-004** | 只读权威（module-architecture / playbook）。本 VP 可建议有界文档修订，不重开架构迁移史 |
| **VP-008 `go`** | 评估若发现改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义的 gap，按 freshness 规则暂挂并登记；不重开 VP-008 |
| **VP-009** | 正交。安全漏洞归 009；本 VP 发现的安全面只转交，不在本 VP 开加固波 |
| **VP-010** | 正交。010 继续做长期 as-designed vs as-built 程序。本 VP 是一次组合层评估 + 路线图重述，不并入 010 波次 |
| **RT-\*** | 评估可重述状态（delivered / registered / gated）；**不消耗**任何 trigger-gated 行 |
| **Admin 功能 / 业务域** | 路线图草案可建议下一拍；本 VP 不实现那些产品项 |

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界关门（证据必须在工作区目标内）：

1. **对照矩阵**：内核 / 组合根 / 模块契约 / 已交付端口 / Profile 与 architecture 文档相对 as-built 有可核对矩阵（覆盖分母在 R1 冻结）。
2. **缺口分类**：每条缺口都有证据，并落入「现在修 / 仍 gated / 接受残余 / 明确不做」之一；无「感觉该优化」而无证据的条目进入分母。
3. **业界对照**：四类参照集每类至少一条对照行，且每行含四格（业界常见 → 本仓现状 → 分类 → 不推翻项）；对照未导致静默改 Charter。
4. **路线图草案**：含现行锚点、已交付 vs 下一拍、residual 总账、三分支建议下一拍；已交 `/vision` 等待用户 editorial 确认（本 VP 不把未确认草案写成已冻结权威）。
5. **边界保持**：未实现 gated 基础设施；未重开已关闭 VP；未改 Charter；未把 Admin 体验增强/新业务域打进本 VP 实现。
6. **审计闭合**：开放 required finding = 0（或已合法闭合）。

建议 Root 纲领（激活后由 `/govern` 写入，本文件只给意图级顺序）：

```text
R1 分母冻结：对照范围、closed-VP residual 清单、业界参照集（已冻结）、「现在修 vs 另立」规则
 → R2 as-built 对照矩阵
 → R3 业界对照 + 缺口分类（若触及 Charter 非目标则停住）
 → R4 路线图草案 + 文档卫生 + 证据与关门
```

同一阶段内的小范围「现在修」代码整改，仅在 R1 规则允许且有独立证据时才做；默认不在本 VP 内做大重构。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 |
|----|--------------|------|----------|----------|------------------|------|
| I-035-001 | 对照分母精确覆盖哪些包、端口、Profile 与文档？ | required | 判据 1 / R2 | R1 | 列出 `apps/api/kernel`、composition/serve、各端口包、`docs/architecture/*` 与 RT-* 行的包含/排除表 | **verified**（2026-09-09 · GOAL-002 D-001 + r1-denominator-freeze.md §1） |
| I-035-002 | 业界参照集是否冻死为会话已确认的四类？ | required | 判据 3 / R3 | R1 | 用户 2026-09-09 书面确认四类；本 VP 正文已写入 | **verified**（2026-09-09 用户书面） |
| I-035-003 | 业界对照是否产生「必须改 Charter 非目标」的结论？ | required | 判据 3/4 / R3→R4 | R3 | 对照表逐行检查；若是则停住交 `/vision` strategic，不得在本 VP 内改 Charter | collecting |
| I-035-004 | 各 closed VP 的 named residual 哪些进入本登记册，哪些保持原 VP 点名？ | required | 判据 2 / R1 | R1 | 扫描 VP-013～034 residual（搬运器、Redis、timestamptz、指标分母等）并分类 | **verified**（2026-09-09 · GOAL-002 D-001 + 附件 §2） |
| I-035-005 | 「现在修」的代码整改是本 VP 内完成，还是只登记、另立 VP/波次？ | required | 判据 5 / 任何代码整改前 | R1 | 用户在 R1 冻结：默认另立；仅小范围、可闭环、不改端口公开语义的项可留在本 VP | **verified**（2026-09-09 · 默认另立；本 VP 仅文档卫生与只读断言） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-035-foundation-architecture-health | GOAL-001-foundation-architecture-health | lead delivery | 2026-09-09 | `/govern` scaffold；Root active 1/4；R1 分母冻结 done（GOAL-002） |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-09 | 初创（v0.1.0）；用户确认结构选型（新 VP + 新区，不作为 VP-010 波次）+ 有界业界对照四类参照集；I-035-002 verified；其余 I 项 R1 冻结；0 区 `planned` |
| 2026-09-09 | v0.2.0 · 激活（VRev-087 self `pass`；架构类 freshness PASS `f2044cf3`→`5c341ec7`；五域 pin/锁/迁移/Profile 默认集零变更；区间 = VP-034 已审结目 + 发布包装；不暂挂 `go`）；lead `workspace-035-foundation-architecture-health`；交 `/govern` 开区 |
| 2026-09-09 | R1 分母冻结（GOAL-002 done）：I-035-001/004/005 verified；默认另立；未改生产代码 |
