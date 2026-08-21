---
id: GOAL-001-observability
doc: audit-entry
record_id: A-002
source: independent
verdict: conditional
scope: Root 关门（GOAL-001-observability 全范围——R1～R5 + 00-meta 成功标准 5 条 + 信息门禁 + 愿景对齐）
status: recorded
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# A-002 · Root 关门独立审计（source: independent）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **日期**：2026-08-22
- **类型**：close-out
- **scope**：workspace-015-observability / GOAL-001-observability 全部交付（R1–R5 = GOAL-002/003/004/005/006）+ 00-meta 成功标准 1–5 + I-001～I-005 信息门禁 + Charter/VP-015 愿景对齐
- **verdict**：**conditional**（开放 required findings = 2，均为 med；无 high）

## 范围与区间

- covered：
  - 工作区绑定：`workspace.md`（id / root_goal / canonical_scope / plan_refs / primary_plan / shared_materials）
  - Root `00-meta` 纲领 R1–R5、成功标准 5 条、I-001～I-005
  - 五子目标 `00-meta` / `01-decision`（含 D-001）/ `02-execution`（含 E-001/E-002）/ `03-audit`（A-001 self）
  - Root `01-decision`、`02-execution`、`03-audit` A-001 self
  - 实现核对：`apps/api/internal/obs`、`internal/config` observability 面、`internal/composition` 接线、`cmd/otlp-sink`、两份 YAML、`compose.yaml`
  - 本会话复跑：`go test ./internal/obs ./internal/config` 与 `./internal/composition -run TestMetrics|TestTracing`（全绿）
  - 愿景对齐：Charter `schema-ui-core-admin-foundation@0.2.0`、VP-015、VRev-033、alignment 机读链（不写 `docs/vision/reviews.md`）
- excluded：
  - 不改任何目标 `status` / `progress` / 方案正文 / goal-tree
  - 不关闭 VP-015（愿景层结项走 `/vision`）
  - 不复跑 E-002 live 双路径冒烟（端口/进程；live 以落盘叙述 + 单测佐证）
  - 不审其他工作区

## 工作区与资料核对

| 项 | 结论 |
|----|------|
| workspace_id | `workspace-015-observability`，与 canonical 路径一致 |
| root_goal | `GOAL-001-observability`（`parent: null`） |
| canonical_scope | `docs/workspaces/workspace-015-observability/` |
| plan_refs / primary_plan | `VP-015-observability`（`active`） |
| vision_role | `delivery`；Charter `primary_workspace` 仍为 workspace-001（未改） |
| shared_materials_catalog | `none`；固定引用表空。无资料引用被当成事实或关闭证据 |
| 跨区 GOAL | 未用裸 `GOAL-*` 跨区取证 |

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| R1 导出合同冻结 + 配置面 | GOAL-002 D-001 §1–§9；E-001/E-002；checkpoint `499f97d` / `45489f4`；`observability.metrics.{enabled,addr,auth_token}` 在两份 YAML + `config.go` Load/ValidateProd；`config_observability_test.go` |
| R2 scrape 接入：独立 listener、`suc_*`、`module_id`、route=注册 pattern | GOAL-003 D-001 / E-002；`internal/obs/{observer,instrumented_mux,server}.go`；composition `newObserver`/`newMetricsServer` + OnStart `metrics.Start` 在 `gate.setReady()` **之前**且不进 readyz；checkpoint `ef33b40` / `5ba04c5` |
| R3 OTLP/HTTP traces + no-op 缺省 | GOAL-004 D-001（闭合 I-002）；`tracing.go`（`otlptracehttp` + ParentBased+ratio + BSP）；缺省 `enabled: false`；checkpoint `0470307` / `2ab4ec4` |
| R4 `correlation.request_id` + baggage `request-id` | GOAL-005 D-001（闭合 I-005）；`TestServerSpanCarriesCorrelationRequestID` 经真实 `requestid.Middleware`；checkpoint `8b52f2d` / `bc5e196` |
| R5 双路径证据方案 + otlp-sink + 落盘 live 叙述 | GOAL-006 D-001 / E-002；工具 `apps/api/cmd/otlp-sink/main.go`（`cf9df6c`）；GOAL-006 `00-meta` `done 4/4` |
| 缺省无收集器（Compose 未加 Prometheus/Jaeger/Tempo） | `compose.yaml` 仅 api+web；YAML `observability.*.enabled: false`；`TestTracingWiringNoopByDefault` / `TestMetricsServerWiringDisabledIsInert` |
| 五子目标产品交付有 D→E→A self 链 | 各 GOAL-002～006 A-001 self `pass`，开放 required = 0 |
| 本会话测试复跑绿 | `go test ./internal/obs/` ok；`./internal/config/` ok；`./internal/composition -run TestMetrics\|TestTracing` ok（2026-08-22） |
| 未进入 A3/A5/Admin 监控页/业务域；feat 切片未改 web | `git show --stat`：`45489f4`/`5ba04c5`/`2ab4ec4`/`bc5e196`/`cf9df6c` 均在 `apps/api`；Charter 目的/非目标/`vision_id@version` 未改（`@0.2.0`） |

## 对照成功标准

| # | 00-meta 成功标准 | 独立核对 | 证据 |
|---|------------------|----------|------|
| 1 | 指标导出面落地；≥1 内核或启用模块路径可 scrape；系列携带 `module_id` | **满足** | 代码：`suc_http_requests_total` 标签含 `module_id`；`TestMetricsCompositionTagsModuleOwnership`（probe=`admin.probe`，health=`core`）；GOAL-003/006 E-002 live 叙述含 `suc_http_requests_total{module_id="core",route="/healthz"}` 与 `suc_kernel_modules_enabled{module_id="admin.users"}` |
| 2 | OTLP traces 可导出；HTTP ≥1 span；可与 request-id / correlation 关联 | **满足**（live 关联为间接证据，见 F-003） | `TestTracingExporterDeliversToOTLPSink`；`TestServerSpanShapeAndStatusError`；`TestServerSpanCarriesCorrelationRequestID`；GOAL-004 E-002 sink 负路径；GOAL-006 E-002 sink `POST … bytes=1037` |
| 3 | 未配置收集器时本地/Compose 默认仍能开发与快测 | **满足** | 缺省 YAML `enabled: false`；Compose 无收集器服务；disabled no-op 单测；GOAL-006 E-002 缺省路径（healthz/readyz 200、25081/4318 无监听、日志零提及） |
| 4 | 显式配置后 metrics scrape **与** ≥1 trace 导出都有可核对证据 | **满足** | 同一次运行叙述见 GOAL-006 E-002；单测分别覆盖 scrape 与 OTLP POST；命令序列可重复 |
| 5 | 未进入 A3/A5/Admin 功能/业务域；未改 Charter；未假装交付 Sentry/剖析/Grafana | **满足** | 各 D 边界；feat 切片无 web/Sentry/pprof/Grafana；Charter 仍 `@0.2.0`，`primary_workspace` 未改 |

## 信息门禁核对（P-005）

| ID | 00-meta | Root `01-decision.md` 索引 | 子目标闭合证据 | 独立结论 |
|----|---------|----------------------------|----------------|----------|
| I-001 | verified（2026-08-21） | **仍标 open** | GOAL-002 D-001 §1–§6 | 实质已闭合；登记不一致 → F-002 |
| I-002 | verified（2026-08-22） | **仍标 open** | GOAL-004 D-001 | 同上 |
| I-003 | verified（出局） | **仍标 open** | GOAL-002 D-001 §7 | 同上 |
| I-004 | verified（不进 readyz） | **仍标 open** | GOAL-002 D-001 §8；composition 接线证实 metrics 不进 gate | 同上 |
| I-005 | verified（2026-08-22） | **仍标 open** | GOAL-005 D-001 | 同上 |

无 `deferred` / `accepted-residual`。无到期未处理 required **实质**项。VP-015 正文 I-015-001～005 仍标 `open`——属愿景层滞后，不构成本 Goal 关门 required（N-001）。

## Findings

### F-001 · Root 权威台账未同步 R5 完成态（阻断无条件关门）

| 字段 | 值 |
|------|-----|
| level | required |
| 严重度 | med |
| status | open |
| evidence | `goal-tree.md` 树：GOAL-006 `[active] …（证据收集中）`；表：Root `progress: 3/5`、GOAL-006 `active 0/4`；树头却写 Root `4/5`。`GOAL-001-observability/00-meta.md` 路线图 R5 =「未开始」「承载子目标：待立项」，`progress: 4/5`。`workspace.md` 纲领表 R5 =「待立项」。对照：`GOAL-006-dual-path-evidence/00-meta.md` 已 `status: done` `progress: 4/4`；A-001 self 主张「GOAL-002～006 均 done」 |
| closure | 须 `/govern` 同步 goal-tree（树+表一致）、Root 路线图 R5→已完成并链接 GOAL-006、workspace.md R5 状态；**不得**在台账仍写 R5 未开始时把 Root 标 `done` |

AGENTS §7：改子目标 status/progress 必须同步 `goal-tree.md`。P-001：纲领路线图须随进展更新。goal-tree **内部** 3/5 与 4/5 已自相矛盾。产品交付可以成立，但 Root 关门声明目前与本区权威台账不一致。

### F-002 · Root `01-decision.md` 信息表与 `00-meta` 矛盾

| 字段 | 值 |
|------|-----|
| level | required |
| 严重度 | med |
| status | open |
| evidence | `GOAL-001-observability/01-decision.md` 信息需求表 I-001～I-005 全部 `状态: open` /「待确认」；同目标 `00-meta.md` 同号项全部 `verified` 且指向子目标 D-001 |
| closure | 将 Root `01-decision.md` 信息表与 `00-meta` 对齐为 verified（或明确声明 00-meta 为唯一登记并删除矛盾行）；关门前 P-005 登记不得两套真相 |

### F-003 · R5 live 未解码 OTLP 载荷核对 `correlation.request_id`

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| evidence | GOAL-006 E-002：sink 只计 `POST / bytes=1037`；`cmd/otlp-sink` 明确「不解析」。关联判据 live 侧仅为「请求/响应头回显同 id」+ 引用 GOAL-005 单测。D-001 §1 事先接受该取证形态 |
| closure | 可选：解码一次 protobuf / 用可解析 sink 断言属性；或用户接受「单测锁定 + 收包」为残余。不阻断本波退出 2（单测已锁） |

### F-004 · 运维文档未暴露 observability 配置面

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| evidence | `apps/api/README.md` 无 observability/metrics/traces 命中；`apps/api/configs/env.example` 无 `OBSERVABILITY_*`。配置键已在 YAML + `config.go` env 映射落地。与 A-001 self N-002 同向 |
| closure | 补 README 配置说明 + env.example 键（secret 仍只给 env 名）；愿景层 VP 关门记录走 `/vision` |

### F-005 · 一等 tracing 依赖在 go.mod 标为 indirect

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| evidence | `apps/api/internal/obs/tracing.go` 直接 import `go.opentelemetry.io/otel` 及 `otel/sdk`、`otlptracehttp`；`apps/api/go.mod` 将这些模块标 `// indirect`（仅 `prometheus/client_golang` 为 direct） |
| closure | `go mod tidy`（或显式 `go get`）把一等 import 提升为 direct，降低后续 tidy 漂移 |

### N-001 · 愿景层索引仍写「实现未做」（不阻断 Goal 关门）

- level: note
- VP-015 信息表 I-015-001～005 仍 `open`；`docs/vision/roadmap.md` RT-O03/RT-O04 仍「实现未做」；`docs/vision/workspaces.md` 仍「Root active 0/5；未改代码」。与本区已落地实现不一致。属 `/vision` 结项范围，本 `/audit` 不改 vision 文件。

### N-002 · 子目标仅 self；independent 按 GOAL-006 D-001 §4 在 Root 执行

- level: note
- Root D-001 曾写「指标/tracing 生产路径实施按 independent」。后续 GOAL-006 D-001 §4 将 independent 收口到 Root 关门。本条意见即该路径。子波次缺 independent 不另开 required。

### N-003 · Compose web 宿主端口与 metrics 缺省端口同号

- level: note
- `compose.yaml` web `${WEB_HOST_PORT:-25081}`；metrics 缺省 `127.0.0.1:25081`。缺省 metrics 关闭，Compose 默认路径无冲突。本机同时开 Compose web 与本地 metrics 默认绑定时可能抢端口。

### N-004 · 既有子目标 note（N-001～N-009）均不阻断

- level: note
- gofmt/CRLF、BSP ~9s、baggage 字符集、headers 增量等均为 open-note。无未闭合 required。

## 必改项汇总

1. **F-001**：同步 `goal-tree.md`（消除 3/5 vs 4/5；GOAL-006 与其 `00-meta` 一致）+ 更新 Root `00-meta` R5 路线图（已完成、链接 GOAL-006）+ `workspace.md` R5 行。未同步前不得 Root `done`。
2. **F-002**：对齐 Root `01-decision.md` 与 `00-meta` 的 I-001～I-005 状态。

推荐（不阻断关门，建议同轮处理）：F-003 文档化残余或补一次载荷断言；F-004 README/env.example；F-005 go.mod tidy。愿景层 N-001 交 `/vision`。

## 与既有意见的异同

| 点 | A-001 self | 本条 independent |
|----|------------|------------------|
| 成功标准 1–5 产品证据 | pass，逐条 ✓ | **同意**（标准 2 的 live 关联为间接，F-003 recommended） |
| I-001～I-005 实质闭合 | 全部 verified | **同意实质闭合**；**不同意**登记已干净（F-002） |
| 五子目标交付 | 均 done | **同意**子目标 `00-meta` 为 done；**不同意**可忽略 goal-tree/Root 路线图滞后（F-001） |
| verdict | **pass** | **conditional** |
| 开放 required | 0 | 2（F-001、F-002） |

**冲突**：self `pass` vs independent `conditional`。对「Root 现在能否无条件关门」结论不同；独立审新增两条 required，self 未列。按 P-004：编排器须展示冲突、给出建议、等用户裁决后才能放行 Root `done`。独立审建议：采纳 conditional，先修 F-001/F-002，再关门；不要用 self pass 覆盖台账门禁。

## 结论 + 建议给编排器/用户的下一步

R1–R5 的**产品交付**有可核对证据：合同、实现、测试（本会话复跑绿）与 R5 落盘命令序列支撑 00-meta 成功标准 1–5 与 VP-015 退出 1–5 的 Goal 侧主张。信息项实质已闭合。愿景机读链（Charter `@0.2.0` → VP-015 `active` → workspace-015 → Root）成立，未越 A3/A5/Admin/Sentry 边界。

**不能 pass 的原因**：Root 关门所依赖的本区权威台账（goal-tree、Root 路线图、workspace 纲领表、Root 决策信息表）仍停留在 R4/R5 未完成叙事，且 goal-tree 自相矛盾。把 Root 标 `done` 会留下「官方树说 R5 进行中、meta 说未开始、self 说已关」三套状态。

建议 `/govern`：响应 A-002 → 闭合 F-001/F-002（只改台账，不重做实现）→ 再问用户是否 Root `done`。VP-015 `closed` 与 roadmap RT-O03/O04 更新另走 `/vision`。

## 声明

本意见 `source: independent`，不修改 `status` / `progress` / 检查点 / 方案正文 / goal-tree。响应与状态变更由 `/govern` 处理。
