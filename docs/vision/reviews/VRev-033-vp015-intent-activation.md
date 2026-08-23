---
doc_type: vision-review
id: VRev-033
status: active
source: self
created: 2026-08-21
updated: 2026-08-21
version: 0.2.0
parent: null
---

# VRev-033 · VP-015 意图完备 / 可行性 / 激活就绪（2026-08-21）

| 字段 | 值 |
|------|-----|
| source | self |
| auditor | `/vision` · 会话编排（grok-4.6） |
| scope | `VP-015-observability`（审视时 `planned` v0.1.0）意图完备、Charter 对齐、退出分母、与 roadmap A4 / RT-O03–O04 一致性、激活与开区就绪、VP-008 `go` 消费前新鲜度（架构类） |
| audit_type | vision-plan（意图 / 激活就绪） |
| verdict | pass |
| 建议 class | editorial |
| open required | 0 |

## 范围与结论

只读核对：`docs/architecture/principles.md` P-006、`docs/vision/alignment.md`、Charter `@0.2.0`、[VP-015-observability](../plans/VP-015-observability.md) v0.1.0、roadmap v0.32.0 A4 / RT-O01–O06、VR-034、`module-architecture.md` §2.2/§7、workspace-003 Root D-011（指标 = 按需、当时无贡献契约）、现行代码（`apps/api/internal/requestid` 的 `X-Request-ID` / `correlation_id`；`apps/api` 无 prometheus/otel/expvar 命中）。未把 `planned` 读成已交付；本报告落盘时尚未改 VP status（激活与开区是用户本轮「通过后」的后续原子动作）。

**总判：pass（0 open required）。** 单愿景与 `vision_ref` 精确匹配；新 VP 承接架构 A4 的结构选型合法；退出分母与用户书面确认同构（指标导出 + OTel traces；无收集器为默认）；方向足以激活并开新 delivery 工作区。两条 recommended 约束 Root 纲领/信息项、退出 4 措辞与 `go` 新鲜度留痕，不阻断激活。

## 核对事实

| 核对项 | 结论 | 证据 |
|--------|------|------|
| 单愿景 | **pass** | 唯一 active Charter `schema-ui-core-admin-foundation@0.2.0` |
| VP→Charter 机读 | **pass** | `vision_ref` 精确匹配 |
| 语义对齐 | **pass** | 可 fork 的 Go 后端内核能力（可观测导出）；不把业务领域当成功条件；不改 Charter 非目标 |
| 最小完备 | **pass** | 意图、配置面、首波冻结、非目标、退出 1–6、邻接 VP、I-015-001～004、工作区表、短史均在 |
| 结构选型 | **pass** | 同愿景新纲领波次 → 新 VP；不重开 VP-012/013/014；不改 Charter；新 delivery 区（惯例 slug 见 D-001） |
| 与 A4 / RT-O03–O04 | **pass** | Prometheus 类 pull + OTLP traces；Sentry / 剖析 / A3 / A5 / Admin 页不进分母 |
| 退出分母有界 | **pass** | 明确排除 Grafana 产品、强制 Compose 收集器、把模块 Observability 改成 MUST |
| 配置面 | **pass** | 缺省无收集器；导出为显式配置；未配置不 fail-closed 挡住 mvp/dev |
| 可行性 | **pass（工作量大、边界清）** | 接入点已知：`requestid.Middleware` 已有请求级 ID 与 context；HTTP 边界可挂 span / 指标。现状无 metrics/otel 基础设施（D-011 仍成立）。工作是补内核导出合同 + HTTP 必埋点，不是换产品叙事 |
| 开放 VRev required | **pass** | 本报告前 open required = 0 |
| 过早交付主张 | **无** | `planned`、0 区；未主张驱动已写 |

## VP-008 `go` 消费前新鲜度（架构类）

VP-008 正文强制 freshness 的对象是**后续业务 VP**。VP-015 是架构分支，且自行把门闩写成激活前须复核。本表按该自设门闩做轻量复核，**不**把本 VP 误读成业务域解锁。

| 项 | 结论 |
|----|------|
| 原 `go` 候选 | `ed99e88`（2026-08-10，clean）；解锁 scope = 标准业务模块框架能力，不是可观测导出 |
| 现行 HEAD | `323c00a`（`docs(workspace-014): 登记 E-003 关门后跟进记录（postgres 方言扫描修复）`） |
| VP-009 | W1–W4、W6 与后续波次已 done；W5 扫描 0 中高危未开子目标；无新的共享基架暂挂宣称 |
| VP-010 | W1–W13 done；`go` 无新暂挂 |
| Vision open required | 0 |
| F-007 residual | 上传授权深度仍 **deferred**（owner=VP-008 lead）。本 VP 不得借导出面扩张授权 scope |
| 本 VP 是否改 Profile / 模块矩阵 / Manifest / 协议 pin | **意图否**。纯内核导出面；若实施时证据显示改变，按消费有效性暂挂 |
| 复核结果 | **PASS（架构激活）**。不消费业务解锁 scope；不暂挂 `go` |

## 不构成 fail / 不新开 required 的诚实边界

1. 退出 4 同时写「scrape **或** trace」与「二者须都有证据」。冻结表与意图 §4 已要求两项都交付。属措辞，不是方向空洞；激活时 editorial 收口。
2. I-015-003 问「HTTP 请求是否必须」，但退出 2 已冻结 HTTP 至少可出 span。激活后该信息项应收窄为 Store / 对象存储 / Job 是否进本波。
3. `/metrics` 默认绑定、是否鉴权、标签是否含秘密仍 open，落在 I-015-001 / R1，不阻断激活。架构 §7 已禁止诊断面泄露秘密。
4. 本 `pass` 允许激活与开区，**不是** R1 导出合同已冻结，也不是可以开始无设计地引入 OTel SDK。
5. 无独立 Vision Review 不是 alignment 强制项。若用户要求交叉审视，另走 `/vision-audit`。
6. 用户本轮确认开区但未点名 slug。不把惯例 slug 写成用户已口头点名；开区记录须留痕推导。

## Findings

### V-F064 · recommended · Root 须写出纲领阶段，收口退出 4 / I-015-003，并登记 I-00N

- level: `recommended`
- status: `open`
- severity: medium
- impact: 若不开区就引 SDK，或把 HTTP 埋点、Store/Job 埋点、scrape 暴露面混成未登记未知，A4 会在 R1/R2 才爆。退出 4 的「或」会让只交 metrics 的人误读成可关门。
- finding: |
  VP 已建议 R1→R5 与 I-015-001～004，但 `planned` 无工作区。激活后 lead Root **必须**：
  1. 写出串行纲领（导出合同冻结 → 指标 scrape → OTel traces → 与 request-id 关联 → 默认无收集器 + 显式导出双路径证据）。
  2. 把 I-015-001/002/004 登记为 required（最晚 R1 或 R2/R3）；I-015-003 收窄为「Store / 对象存储 / Job 是否进本波」（HTTP span 已由退出 2 冻结）；另登记 request-id 写入 span 的属性名（required，最晚 R4）。
  3. I-015-001 须覆盖 scrape 路径/端口/**绑定与鉴权**、基数、内核 vs 模块最小集合、标签不得含秘密。
  4. 激活时把退出 4 的「或」改为与冻结表同构：metrics **与** traces 都要有可核对证据，允许分路径验收。
- evidence: VP-015 退出 1–4；I-015-001～004；`module-architecture.md` §7；`requestid.HeaderName`。
- closure: |
  Root `00-meta` 含 P-001 阶段表 + 上述 I-00N；VP 退出 4 与 I-015-003 作 editorial 修订。不要求本 Review 落盘时已经有答案。
- 建议 class: `editorial`

### V-F065 · recommended · 激活记录须留下架构类 freshness 结论，避免被读成已消费业务 `go`

- level: `recommended`
- status: `open`
- severity: low
- impact: VP-015 把门闩写成「激活前 freshness review」。若激活记录只写「开区」而不点名：本 VP 非业务域、不改 Profile 意图、F-007 不升格、现行无 `go` 暂挂，后续读者会把架构 A4 误读成已走 VP-011 那种业务消费。
- finding: 激活时在 VP 短史或 lead Root D-001 写入上表复核结论与候选/HEAD 指针。
- evidence: VP-008 §`go` 消费有效性（业务 VP）；VP-015 激活门闩；GOAL-007 D-001 原候选 `ed99e88`。
- close requirement: D-001 或 VP 激活短史含 freshness 表；不要求重开 VP-008。
- 建议 class: `editorial`

### 声明

本意见不直接修改 Charter / VP / Goal status。required finding 的响应由 `/vision` 追加在本报告中；原 verdict 与 finding 原文不得改写。

### 门禁含义

- Vision Review **open required = 0**。
- **允许**：用户确认后激活 VP-015、开新 delivery 工作区、按 V-F064/V-F065 写 Root 纲领与 freshness。
- **禁止**：把本 `pass` 写成 R1 已冻结或 SDK 已引入；重开 workspace-014；把 VP-015 当业务 VP 消费 `go` 解锁 scope。

### 响应（2026-08-21 · `/vision` 激活与开区）

不回溯改写原 verdict `pass` 与 finding 正文。

| finding | 闭合 | 证据 |
|---------|------|------|
| V-F064 | **fixed** | Root `GOAL-001-observability` 纲领 R1～R5；I-001～I-004 required open（对应 I-015-001/002/004 与收窄后的 I-015-003）；I-005 request-id↔span required（最晚 R4）。VP-015 v0.2.0 退出 4 改为 metrics **与** traces；I-015-003 收窄。答案仍 open，登记已满足 close requirement。 |
| V-F065 | **fixed** | D-001 架构类 freshness：原 `go` `ed99e88`；HEAD `323c00a`；非业务解锁；不暂挂 `go`；F-007 不升格。 |

用户书面「对 VP-015 做 self Review；通过后激活并交 /govern 开区」已执行：VP `active`；lead `workspace-015-observability`（惯例 slug，D-001 留痕）；Root scaffold。本 scope **0 open required、0 open recommended**。
