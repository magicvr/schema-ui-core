---
doc_type: vision-plan
id: VP-015-observability
title: 可观测性（指标导出 + OpenTelemetry）
status: closed
vision_ref: schema-ui-core-admin-foundation@0.3.0
lead_workspace: workspace-015-observability
created: 2026-08-21
updated: 2026-08-22
version: 0.3.0
parent: null
---

# VP-015 · 可观测性（指标导出 + OpenTelemetry）

## 状态与门闩（2026-08-22 · 已关门）

| 项 | 值 |
|----|-----|
| status | **`closed`**（2026-08-22 用户书面确认有界组合层关门；VRev-034 `V-F066` → `fixed`；关门依据 = 本轮独立代码/测试/live，不以 Goal 台账为充分条件） |
| **lead_workspace** | **`workspace-015-observability`**（Root `GOAL-001-observability` `done 5/5`；唯一 delivery；**不**重开 workspace-014） |
| **Vision required** | **已满足**：VRev-033 / VRev-034 均为 `pass`，open required = 0；`V-F064`/`V-F065`/`V-F066` recommended 已闭合 |
| **关门门闩（现行）** | 已 `closed`；保留 workspace-015 历史绑定，默认不接新区；reopen 须用户确认 |
| **组合位置** | 架构分支 A4；前提 = VP-013 A1 与 VP-014 A2 均已有界 `closed`；roadmap **RT-O01**/**RT-O02** 已 delivered，**RT-O03**/**RT-O04** 本 VP 交付 |
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
| I-015-001 | 指标面：Prometheus scrape 路径与端口、绑定/鉴权、基数上限、内核 vs 模块贡献最小集合、标签不得含秘密。禁止「先堆业务指标再补合同」。 | required | 方案冻结 / 实施 | R1 合同冻结 | verified（GOAL-002 D-001；本轮 live scrape 核对） |
| I-015-002 | Tracing：OTLP HTTP vs gRPC、采样默认、未配置 endpoint 时的 no-op 语义。 | required | 方案冻结 / 实施 | R3 接入前 | verified（GOAL-004 D-001：OTLP/HTTP、ParentBased+ratio 缺省 1.0、disabled 纯 no-op；本轮 live POST） |
| I-015-003 | Store / 对象存储 / Job 是否进本 VP 分母还是后续增量。（HTTP 请求 span 已由退出 2 冻结，不再作为本项问题。） | required | 方案冻结 | R1 合同冻结 | verified（出局 · GOAL-002 D-001 §7；本 VP residual 点名） |
| I-015-004 | `/metrics` 或 OTLP 是否进入 `readyz` 依赖？默认建议：未显式配置则不扩 ready。 | required | 方案冻结 | R2/R3 接入前 | verified（均不进 readyz；composition `metrics.Start` 在 `gate.setReady()` 之前且 health.go 无 metrics probe） |
| I-015-005 | request-id / correlation 如何写入 span（属性名、是否 baggage）才能满足退出 2 的关联判据。 | required | 方案冻结 | R4 关联前 | verified（属性 `correlation.request_id` + baggage `request-id`；本轮 live 载荷含二者） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-015-observability | GOAL-001-observability | lead | 2026-08-21 | 2026-08-21 用户确认激活并开区；2026-08-22 VP 组合层 `closed`；Root done 5/5；默认不接新区 |

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-22 | **closed**（有界 · 架构 A4） | 用户确认组合层关门，且要求以**独立代码核验**而非治理记录为充分条件。exit 1：专用 Prometheus listener + `suc_*` 系列携带 `module_id`；本轮 live scrape 核到 `module_id="core"` `/healthz` 与 `admin.users` gauge。exit 2：OTLP/HTTP SERVER span；属性 `correlation.request_id` + baggage `request-id`；本轮 live 载荷含属性名与请求 id 字节。exit 3：缺省 `enabled: false`；Compose 无收集器；本轮缺省 healthz/readyz 200 且无 metrics 监听。exit 4：同一次显式运行 scrape **与** OTLP POST 同时成立。exit 5：未进 A3/A5/Admin/业务；Charter 仍 `@0.2.0`；无 Sentry/剖析/Grafana。exit 6：实现层与 VRev required = 0。V-F066 按本表闭合。 | 本轮 `/vision` 独立核验：`apps/api/internal/obs` + composition 接线 + 两份 YAML + `compose.yaml` + `go.mod`；`go test ./internal/obs ./internal/config ./internal/composition -count=1` 全绿；live 缺省路径 + 显式 `25199` scrape / `24318` OTLP capture。指针：[VRev-034](../reviews/VRev-034-vp015-closeout-readiness.md)；lead `workspace-015` goal-tree（Root done 5/5）；Root A-002/A-003（required 已 fixed，不作充分条件） | **`workspace-015` / `GOAL-001-observability` / A-002 F-003**：仓库内 `cmd/otlp-sink` 不解析 OTLP protobuf（证据工具，不是 collector）；本轮关门取证用一次性 capture sink 核对 live 载荷含 `correlation.request_id` 与请求 id 字节。**I-015-003**（Root I-003）：Store / 对象存储 / Job 指标不进本波分母（已 verified 出局，不是未交付）。 |

### 退出判据 ↔ 证据（本轮独立核验）

| 退出 | 结论 | 证据 |
|------|------|------|
| 1 指标导出 + `module_id` | 满足 | `internal/obs` 私有 registry；`suc_http_*` 标签含 `module_id`；贡献路由 `mux.Own`；主 mux 无 `/metrics`。live：`suc_http_requests_total{module_id="core",route="/healthz"}` 与 `suc_kernel_modules_enabled{module_id="admin.users"}` |
| 2 OTLP traces + request-id 关联 | 满足 | `otlptracehttp`；SERVER span；`correlation.request_id` + baggage `request-id`。本轮测试绿；live dump 1285 字节含属性名与 `vrev034-corr-0001` |
| 3 无收集器默认 | 满足 | YAML 缺省 `enabled: false`；Compose 无收集器；live 缺省 healthz/readyz 200、25081/25199 不可 scrape、日志零提及 |
| 4 显式 scrape **与** trace | 满足 | 本会话同一次显式运行：scrape 200（11771 字节）**与** OTLP POST |
| 5 未进 A3/A5/Admin/业务 / 未改 Charter | 满足 | 无 sentry/grafana/pprof；Charter 仍 `@0.2.0`；无 Admin 监控页 |
| 6 required = 0 | 满足 | VRev-033/VRev-034；实现层 A-003 后开放 required = 0 |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-21 | 初创 `planned`：用户确认新建本 VP 承接架构 A4；退出分母 = Prometheus 类指标导出 + OpenTelemetry traces；A3 / A5 / Sentry / 剖析 / Admin 页 / 业务域不进分母。未激活、未开区 |
| 2026-08-21 | VRev-033 self `pass`（0 required）；用户确认激活并开区。v0.2.0 `planned → active`；lead = `workspace-015-observability`；退出 4 editorial「或→与」；I-015-003 收窄；I-015-001 补绑定/鉴权。Root 承接 P-001 与 I-00N（V-F064）及架构类 freshness（V-F065） |
| 2026-08-22 | VRev-034 self `pass`：组合层关门就绪；核验 = 本轮独立源码/测试/live，不以 Goal 台账为充分条件。用户确认有界组合层关门（v0.3.0）：关门记录含 exit↔证据映射 + F-003 / I-015-003 residual 点名；组合索引原子同步（VR-036） |
