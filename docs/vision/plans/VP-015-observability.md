---
doc_type: vision-plan
id: VP-015-observability
title: 可观测性（指标导出 + OpenTelemetry）
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-015-observability
created: 2026-08-21
updated: 2026-08-21
version: 0.2.0
parent: null
---

# VP-015 · 可观测性（指标导出 + OpenTelemetry）

## 状态与门闩（2026-08-21 · active）

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-08-21 用户确认激活并开区；VRev-033 `pass`，V-F064/V-F065 → `fixed`） |
| **lead_workspace** | **`workspace-015-observability`**（Root `GOAL-001-observability`；唯一 delivery；**不**重开 workspace-014） |
| **Vision required** | **已满足**：VRev-033 `pass`，open required = 0；V-F064/V-F065 recommended 由激活 + Root scaffold 闭合 |
| **激活门闩（现行）** | 已激活；实现证据在 lead 区。改变 Profile / 模块矩阵 / Manifest / 共同门禁时按 VP-008 `go` 消费有效性暂挂 |
| **组合位置** | 架构分支 A4；前提 = VP-013 A1 与 VP-014 A2 均已有界 `closed`；roadmap **RT-O01**/**RT-O02** 已 delivered，**RT-O03**/**RT-O04** 本 VP 冻结退出分母 |
| **完整 ≠ 架构清单无限扩张** | 本 VP 只承接 A4。A3 多实例/Redis/队列、A5 密钥轮换、Sentry、连续剖析、Admin 可观测页、业务域不进退出分母 |

## 意图

在 VP-003 单主线模块化内核、已交付的结构化日志 / correlation（RT-O01）与 `/healthz` `/readyz`（RT-O02）之上，把「指标 = 按需、当前无贡献契约」（workspace-003 Root D-011）收成**可导出的内核可观测合同**：

1. **指标导出（RT-O03）**：Prometheus 类 scrape（或同等 pull 面）。系列至少携带 `module_id`（[module-architecture.md](../../architecture/module-architecture.md) §7）。不是 Grafana 产品、不是 Admin 仪表盘。
2. **分布式 tracing（RT-O04）**：OpenTelemetry 导出（默认 OTLP）。请求级 correlation / request-id 可与 trace 关联。不是替换已有日志与错误包络。
3. **内嵌默认**：本地双进程与 Compose **不**要求 Prometheus / collector / Jaeger。未配置导出面时进程仍能启动与快测。不得把「没收集器就不能启动」做成 mvp/dev 默认。
4. **生产向验收**：显式配置后可核对至少一条 scrape/export 路径与至少一条 trace 导出。

本 VP 属**架构分支**，不承载 Admin 功能页或业务域。不重开 VP-012。

## 配置面

可观测导出由配置选择，**不是**改 Profile、也不是改模块矩阵：

- **缺省**：无外部收集器、无强制 `/metrics` 对公网暴露义务。本地与 Compose 默认不变。
- **生产 / 本 VP 验收**：显式打开指标 scrape 与/或 OTLP endpoint（具体键名由 lead Root 方案冻结）。凭证与 endpoint 走 YAML + env 插值、密钥 fail-closed，不把 secret 写入仓库。
- 未配置 tracing backend 时不得 fail-closed 挡住 mvp/dev；指标面若默认绑定本地 loopback，也不得成为启动硬依赖。

## 首波冻结（退出分母 = 架构 A4）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| 指标导出 | Prometheus 类 pull 面；内核 + 已启用模块路径可观察；标签含 `module_id` | Grafana/产品仪表盘；Admin 监控页；push gateway 作为唯一面；无限自定义业务指标 |
| Tracing | OTLP 导出；HTTP 请求至少可出 span；与现有 request-id / correlation 可关联 | 替换结构化日志；Sentry（RT-O06）；连续剖析（RT-O05）；服务网格 sidecar |
| 默认路径 | 无收集器仍能开发与快测 | 强制 Compose 常驻 Prometheus/Jaeger/Tempo |
| 既有可观测 | 保留日志 `module_id`、`/healthz` `/readyz`、request-id | 重开 VP-012 correlation/错误包络 |

## 非目标

- 多实例 / Redis / 外部队列 / 分布式锁 / 进程分离（A3：RT-Q\* / RT-D03 / RT-D04）
- JWT 密钥轮换 / KMS / TLS 终止（A5 / RT-K\* / RT-D05）
- 错误汇聚 Sentry 类（RT-O06）、连续性能剖析（RT-O05）
- Admin 全局搜索、Command Palette、监控产品页；业务域页面
- 改变 Charter 边界；重开 VP-012 / VP-013 / VP-014；替代 VP-009 / VP-010
- 把模块 Observability 贡献从「按需」改成每个模块 MUST（架构 §2.2 仍为按需；**内核导出面**本 VP 必交）

## 与相邻 VP 的边界

| VP | 关系 |
|----|------|
| **VP-003** | 遵守薄内核。本 VP 是 D-011「指标若引入必须带 `module_id`」的落地波，不重开 VP-003 |
| **VP-012** | 已 closed 的 correlation / 错误包络 / Job 六态不重开；本 VP 在其上加 metrics + traces |
| **VP-013 / VP-014** | 已 closed 的 Store / 对象存储端口可被埋点消费，不改持久化或对象存储合同 |
| **VP-008 `go`** | 若实现改变 Profile 默认集 / 模块矩阵 / Manifest 装配 / 共同门禁，按消费有效性做 freshness review。纯导出面接入若证据显示未改上述语义，不自动暂挂 `go`。激活前仍须架构类 freshness |
| **VP-009 / VP-010** | 安全 finding 与设计符合性 gap 仍归持续程序；本 VP 不扩扫描范围，也不做 Sentry |
| **Admin 功能 / 业务域** | 只消费导出合同；不得在本 VP 加监控产品页或领域表 |

## 方向级退出判据

1. 指标导出面已落地；至少一条内核或已启用模块路径可 scrape/核对；系列携带 `module_id`。
2. OpenTelemetry traces 可经 OTLP（或书面冻结的同等协议）导出；HTTP 请求至少可出 span，并能与现有 request-id / correlation 关联（或书面记录不可关联的有界 residual）。
3. 未配置收集器时本地/Compose 默认仍能开发与快测；导出不是 mvp/dev 启动硬依赖。
4. 生产向验收以显式配置为准：指标 scrape **与** 至少一条 trace 导出都须有可核对证据（允许分路径验收）。
5. 未进入 A3 / A5 / Admin 功能 / 业务域；未改 Charter；未假装交付 Sentry / 剖析 / Grafana 产品。
6. 开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root `GOAL-001-observability`（P-001）书写：R1 导出合同冻结 → R2 指标 scrape → R3 OTel traces → R4 与 request-id 关联 → R5 默认无收集器 + 显式导出双路径证据。

## 信息需求（P-005）

允许带未知立项。下列不影响「本 VP 意图已冻结」，但必须在对应阶段前关闭或经用户接受残余。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-015-001 | 指标面：Prometheus scrape 路径与端口、绑定/鉴权、基数上限、内核 vs 模块贡献最小集合、标签不得含秘密。禁止「先堆业务指标再补合同」。 | required | 方案冻结 / 实施 | R1 合同冻结 | open |
| I-015-002 | Tracing：OTLP HTTP vs gRPC、采样默认、未配置 endpoint 时的 no-op 语义。 | required | 方案冻结 / 实施 | R3 接入前 | open |
| I-015-003 | Store / 对象存储 / Job 是否进本 VP 分母还是后续增量。（HTTP 请求 span 已由退出 2 冻结，不再作为本项问题。） | required | 方案冻结 | R1 合同冻结 | open |
| I-015-004 | `/metrics` 或 OTLP 是否进入 `readyz` 依赖？默认建议：未显式配置则不扩 ready。 | required | 方案冻结 | R2/R3 接入前 | open |
| I-015-005 | request-id / correlation 如何写入 span（属性名、是否 baggage）才能满足退出 2 的关联判据。 | required | 方案冻结 | R4 关联前 | open |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-015-observability | GOAL-001-observability | lead | 2026-08-21 | 2026-08-21 用户确认激活并开区；唯一 delivery；不重开 workspace-014 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-21 | 初创 `planned`：用户确认新建本 VP 承接架构 A4；退出分母 = Prometheus 类指标导出 + OpenTelemetry traces；A3 / A5 / Sentry / 剖析 / Admin 页 / 业务域不进分母。未激活、未开区 |
| 2026-08-21 | VRev-033 self `pass`（0 required）；用户确认激活并开区。v0.2.0 `planned → active`；lead = `workspace-015-observability`；退出 4 editorial「或→与」；I-015-003 收窄；I-015-001 补绑定/鉴权。Root 承接 P-001 与 I-00N（V-F064）及架构类 freshness（V-F065） |
