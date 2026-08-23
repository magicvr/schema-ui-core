---
id: GOAL-004-otel-traces
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
---

## D-001 · Tracing 合同与 no-op 语义（闭合 Root I-002 / VP I-015-002）

### 触发

R2 已落地 InstrumentedMux 拦截点与 metrics 面；I-002 为 required，最晚需要阶段 = R3 接入前。已核对：otel SDK v1.45.0 可达（proxy 实测）；OTLP/HTTP 约定端口 4318。

### 决定

**§1 协议 —— OTLP/HTTP（protobuf），不做 gRPC**

`go.opentelemetry.io/otel/sdk` + `exporters/otlp/otlptrace/otlptracehttp`。endpoint 走 `observability.traces.endpoint`（如 `http://localhost:4318`；SDK 自动追加 `/v1/traces`）。

**§2 配置面（键名冻结，语义本波实现）**

| 键 | env | 默认 | 语义 |
|----|-----|------|------|
| `observability.traces.enabled` | `OBSERVABILITY_TRACES_ENABLED` | `false` | 显式打开；缺省完全 no-op |
| `observability.traces.endpoint` | `OBSERVABILITY_TRACES_ENDPOINT` | 空 | OTLP/HTTP base endpoint；enabled 时必填 |
| `observability.traces.sample_ratio` | `OBSERVABILITY_TRACES_SAMPLE_RATIO` | `1.0` | 采样比 ∈ (0,1] |

fail-closed 校验（对齐 metrics 先例）：enabled + endpoint 空 → 拒绝；endpoint 非 http(s) 绝对 URL → 拒绝；disabled + endpoint 有值 → 拒绝（死配置）；sample_ratio 出 (0,1] → 拒绝。校验进 `Load` + `ValidateProd` 双侧。collector 鉴权 headers 本波不加（后续增量须新 D 记录）。

**§3 no-op 语义（VP 意图 3 的 tracing 半边）**

`enabled=false`（含全缺省）：不创建 exporter/provider 后台设施，`Wrap` 不创建 span——零行为变化、mvp/dev 缺省路径不受影响。`enabled=true`：BatchSpanProcessor 异步导出，导出失败只经自定义 error handler 记 warn 日志，**绝不 crash 进程、不影响请求路径**（tracing 是旁路观测面，同 R2 listener 语义）。

**§4 采样**

`ParentBased(TraceIDRatioBased(sample_ratio))`：尊重上游 traceparent 的采样决定（分布式一致），无上游时按比例根采样。缺省 ratio=1.0（显式打开即全采，运维可调低）。

**§5 span 面（复用 R2 拦截点）**

Observer.Wrap 扩展：tracing 启用时每个注册路由包 SERVER span——name = `<METHOD> <route>`（注册 pattern），属性 `http.request.method` / `http.route` / `url.path`（不含 query，防凭证入 span）/ `http.response.status_code`；status ≥500 → SpanStatus Error。传播：W3C TraceContext + Baggage 提取（join 上游 trace；baggage 提取为 R4 关联预留输入）。注入响应头不做。

**§6 Resource 与生命周期**

Resource：`service.name`=app.name、`service.version`=版本、`deployment.environment`=APP_ENV。fx lifecycle：provider 组装即生效（BSS 自启），OnStop `Shutdown(ctx)` flush（join 停机链）。obs 包不 import config，参数以原始值传入（沿用 R2 边界）。

### 为什么

- OTLP/HTTP 无需额外 gRPC 依赖栈、穿透反向代理/负载均衡最简单，且是 OTel 默认推荐起点。
- ParentBased+ratio 是社区默认采样链；no-op 路径连 provider 都不建，保证「未配置零开销」可测试。
- 复用 Wrap 拦截点让 trace/metrics 共享同一 route 标签源，避免两套 instrumentation 漂移。

### 未选方案

- **gRPC OTLP**：多引入 grpc-go 全家桶，代理穿透配置更重；无当前需求。
- **全局 otel.SetTracerProvider + otelhttp 中间件**：全局状态伤测试隔离；otelhttp 外层中间件拿不到 route pattern（同 R2 D-001 分析），route 标签会退化。
- **每请求同步导出（SimpleSpanProcessor）**：把导出延迟放进请求路径，违背旁路原则。
