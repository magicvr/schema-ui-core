---
doc_type: vision-review
id: VRev-034
status: active
source: self
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
parent: null
---

# VRev-034 · VP-015 关门就绪度审视（2026-08-22）

| 字段 | 值 |
|------|-----|
| source | self |
| auditor | `/vision` · 会话编排（grok-4.6）；本轮对代码 / 测试 / live 双路径**独立复验**，不以 Goal 台账为关门充分条件 |
| scope | `VP-015-observability` 组合层关门就绪 · 退出判据 1–6 · 代码实现 · 本会话测试与 live scrape/OTLP · Vision required · 有界 residual · 组合索引 |
| audit_type | vision-plan（关门就绪度 · 代码成果独立核验） |
| verdict | pass |
| 建议 class | editorial（组合层关门 + 索引同步 + residual 点名；不改 Charter 方向） |
| open required | 0 |

## 范围与结论

只读核对：`docs/architecture/principles.md` P-006、`docs/vision/alignment.md` §6/§7/§9、`charter.md` `@0.2.0`、[VP-015-observability](../plans/VP-015-observability.md)（审视时 `active` v0.2.0）、`roadmap.md`、`workspaces.md`、`revisions.md`（至 VR-035）、`reviews.md` 与 `reviews/VRev-001`～`VRev-033`、lead 工作区 `workspace-015-observability`（绑定与 Root 状态只作**指针**，不替代代码核验）。

**本轮独立核验（非转录治理记录）**：

1. **源码**：`apps/api/internal/obs/{observer,instrumented_mux,server,tracing}.go`；`internal/composition` 接线（`newTracing` / `newObserver` / `newMetricsServer` / `metrics.Start` 在 `gate.setReady()` 之前且不进 `readyz`）；`internal/config` 校验；两份 YAML 缺省 `enabled: false`；`compose.yaml` 仅 api+web、无 Prometheus/Jaeger/Tempo；`go.mod` 一等依赖含 `prometheus/client_golang` 与 `otel` / `sdk` / `otlptracehttp` / `trace`。
2. **测试（本会话复跑）**：`go test ./internal/obs ./internal/config ./internal/composition -count=1` 全绿（obs 2.322s / config 1.172s / composition 15.706s）；`go vet ./internal/obs ./internal/config ./cmd/otlp-sink` 通过。
3. **live（本会话，独立端口，不复用 GOAL-006 叙述）**：缺省路径 `APP_ENV=development` 无任何 `OBSERVABILITY_*` → `GET /healthz` 与 `/readyz` **200**，`127.0.0.1:25199` 与默认 `25081` 均不可 scrape，启动日志 observability/metrics 命中 **0**。显式同一次运行：`OBSERVABILITY_METRICS_ENABLED=true`（`127.0.0.1:25199`）+ `OBSERVABILITY_TRACES_ENABLED=true`（OTLP `http://127.0.0.1:24318`）→ scrape **200**（11771 字节）含 `suc_build_info`、`suc_http_requests_total{method="GET",module_id="core",route="/healthz",status="200"}`、`suc_kernel_modules_enabled{module_id="admin.users"} 1`；请求头 `X-Request-ID=vrev034-corr-0001` 响应回显同 id；capture sink 收到 `POST / bytes=1285`；live 载荷 UTF-8 可检出 `correlation.request_id` **与** 该 request-id 字节；指标标签未泄漏该 id。

未把 Goal `progress=5/5` 或 A-001/A-002 正文当作退出判据的充分证据。治理记录仅用于定位实现与对照开放 required。

**总判：pass（0 open required · 1 open recommended）。**

**关门的实质证据已齐备**（代码 + 本会话测试 + 本会话 live），可按 alignment §7 做**有界 closed**。Vision Review open required = 0；对齐链成立；激活后 Charter 仅 editorial（VR-034/VR-035），无 strategic 宽阻断。组合索引仍写 VP-015 `active`，且 `workspaces.md` 仍错误投影「Root active 0/5；未改代码」——这是待用户书面确认后的投影同步，**不是**实现缺口。本轮用户意图为「复核是否满足关门条件；独立审计代码成果；是的话关门」。

本意见原文**不**把组合索引改写成 `closed`。

### 核对事实

| 核对项 | 结论 | 证据（本会话独立核验，除非标明指针） |
|--------|------|------|
| 单愿景 / `vision_ref` | **pass** | 唯一 active Charter `schema-ui-core-admin-foundation@0.2.0`；VP-015 `vision_ref` 精确匹配 |
| 工作区绑定 | **pass** | `workspace-015` 唯一 lead / delivery；`plan_refs` / `primary_plan` / `vision_role: delivery` 合规；Charter `primary_workspace` 仍为 workspace-001 |
| 区证据指针（§7.1） | **pass（指针）** | goal-tree Root `done 5/5`、GOAL-002～006 `done`。**不**据此放行；放行依据是下表退出 1–5 的代码/测试/live |
| 实现层开放 required | **pass（指针）** | Root A-002 F-001/F-002 → A-003 `fixed`；开放 required = 0。本轮未改 Goal finding |
| 退出 1 · 指标 scrape + `module_id` | **pass** | 私有 registry；`suc_http_*` 标签含 `module_id`；贡献路由 `mux.Own(pattern, route.ModuleID)`；主 mux **无** `GET /metrics`。live scrape 见上 |
| 退出 2 · OTLP + HTTP span + request-id 关联 | **pass** | `otlptracehttp`；SERVER span；属性 `correlation.request_id`；baggage 键 `request-id`。单测 `TestServerSpanCarriesCorrelationRequestID` / `TestTracingExporterDeliversToOTLPSink` 本轮绿。live 载荷含属性名与 id 字节 |
| 退出 3 · 无收集器默认 | **pass** | YAML `observability.*.enabled: false`；Compose 无收集器服务；live 缺省无额外端口、无 observability 日志 |
| 退出 4 · 显式双路径 | **pass** | 本会话同一次显式运行同时成立 scrape **与** OTLP POST |
| 退出 5 · 未进 A3/A5/Admin/业务；未改 Charter；未假装 Sentry/剖析/Grafana | **pass** | `apps/api` 无 sentry/grafana/pprof 命中；`go.mod` 无可观测产品依赖；Charter 仍 `@0.2.0` |
| 退出 6 · required = 0 | **pass** | Vision Review 与实现层开放 required 均为 0 |
| Vision required（§6 门禁 8） | **pass** | `reviews.md` open required = 0；VRev-033 为激活审视，本条为关门就绪首份 |
| Charter strategic 后 re-align | **pass** | 激活后仅 VR-034/VR-035（editorial）；无宽阻断 |
| 组合索引当前陈述 | **pass（待同步）** | Charter 关系节 / `roadmap.md` 第 15 行与 RT-O03/O04「实现未做」/ `reviews.md` 摘要 / `workspaces.md`（仍写 0/5 未改代码）/ 区 `workspace.md` 仍写 VP-015 `active` |

## Findings

#### V-F066 · 组合层关门须同步索引，并显式映射 exit 1–6 ↔ 本轮独立证据、点名有界 residual

- level: `recommended`
- status: `open`
- severity: low
- impact: alignment §7.2 允许有界 closed，但 residual 必须点名到具体 workspace / goal id。若只让 Root `done` 而组合索引仍称 `active`、且 `workspaces.md` 仍写「未改代码」，后续读者会把 A4 读成未交付或与代码事实相反。
- finding: |
  1. 用户确认组合层关门时一次写清 exit 1–6 ↔ **本轮独立**证据（源码路径、本会话 `go test`/`go vet`、本会话 live 缺省 + 显式双路径；治理 A 条目仅作指针）。
  2. residual 至少点名：`workspace-015` / `GOAL-001-observability` / Root A-002 **F-003**：仓库内 `cmd/otlp-sink` 仍不解析 OTLP protobuf（证据工具，不是 collector）。本轮关门取证用一次性 capture sink 核对 live 载荷含 `correlation.request_id` 与请求 id 字节；引入可解析 collector 时再补语义断言。另点名 **I-015-003**（Root I-003）：Store / 对象存储 / Job 指标不进本波分母（已 verified 出局，不是未交付）。
  3. 同步 `roadmap.md`（VP 行 + **RT-O03/O04 → delivered**）/ `workspaces.md`（纠正 0/5、未改代码）/ Charter 关系节 / `reviews.md` 摘要 / 区 `workspace.md`：VP-015 → `closed`；当前无 active **交付** VP（持续程序仍为 VP-009/VP-010）。将 VP 信息表 I-015-001～005 与 Root verified 对齐。Root `done` 不能冒充 VP 层用户确认。
- evidence:
  - 本会话 live：缺省 healthz/readyz 200 且无 metrics 监听；显式 scrape 11771 字节 + OTLP POST 1285 字节且载荷含关联字段
  - 本会话测试：`internal/obs` `internal/config` `internal/composition` 全绿
  - `docs/vision/roadmap.md` RT-O03/O04 仍写「实现未做」；`workspaces.md` 仍写 Root 0/5
  - alignment.md §7.2
- closure: |
  `/vision` 在用户书面确认组合层关门时按上列三项一并完成。本 finding 不阻断「就绪」结论，只约束关门落盘形状。
- 建议 class: `editorial`

### 不构成 fail / 不新开 required 的诚实边界

1. 本 `pass` **不是**「组合索引已 closed」：用户书面确认与索引原子同步仍待发生（本轮用户已给出「满足则关门」）。
2. `/vision` 本入口只写 `source: self`。用户要求「独立审计代码成果」已在本报告用源码+测试+live 兑现；不是 `/vision-audit` 的 independent VRev。无独立 Vision Review 不是 alignment 强制项（强制时机仅为 Charter 初建与 strategic）。
3. 本轮 live 关联核验是 **protobuf 字节中的 UTF-8 字符串存在性**，不是 otlp collector 语义解码。它强于 GOAL-006 E-002 的「只计 POST 字节」；仍不等于完整 span 模型断言。单测已锁属性名与 baggage。不新开 required。
4. Compose `web` 宿主端口默认 `${WEB_HOST_PORT:-25081}` 与 metrics 缺省绑定位同号（A-002 N-003）。缺省 metrics 关闭，Compose 默认路径无冲突。不进本 VP residual。
5. 不把 progress=`5/5` 当作关门权威。
6. 架构 A3（多实例/Redis/队列）、A5 密钥轮换、Sentry、连续剖析、Admin 监控页本就不在退出分母。

### 声明

本意见不修改 Charter / VP / Goal status 或 progress；required/recommended finding 的响应由 `/vision` 追加在本报告中；实现层执行仍交 `/govern`。原 verdict 与 finding 原文不得改写。本入口不关闭 Goal finding。

### 门禁含义

- Vision Review **open required = 0**。
- **允许**：用户确认后，`/vision` 按 V-F066 执行 VP-015 有界组合层关门与索引同步。
- **禁止**：在无用户书面确认时把组合索引改成 VP-015 `closed`；把 Root `done` 或 Goal 审计原文冒充本轮代码核验；把 Grafana / Sentry / 剖析 / Admin 监控页写成已交付。

### 响应（对 self 意见 · VRev-034 findings 闭合 · 2026-08-22）

| date | actor | summary |
|------|-------|---------|
| 2026-08-22 | `/vision` · 用户书面「复核 VP-015 是否满足关门条件（请独立审计代码成果，不能单纯以治理记录作为关门依据），是的话关门」 | **不回溯改写**原 verdict `pass` 与 finding 正文。**V-F066 → `fixed`**：VP-015 组合层确认 **有界 `closed`**（架构 A4）。关门记录含 exit 1–6 ↔ 本轮独立证据；residual 点名 `workspace-015` / `GOAL-001` / A-002 F-003 与 I-015-003。`roadmap.md` / `workspaces.md` / Charter 关系节 / `reviews.md` / 区 `workspace.md` 原子同步（VR-036）。本 scope **0 open required、0 open recommended**。 |
