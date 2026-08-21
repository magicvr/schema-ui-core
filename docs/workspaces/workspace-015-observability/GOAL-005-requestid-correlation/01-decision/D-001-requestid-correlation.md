---
id: GOAL-005-requestid-correlation
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-001 · request-id ↔ span 关联合同（闭合 Root I-005 / VP I-015-005）

### 触发

R3 已证明 HTTP server span 可导出；VP-015 退出判据 2 要求「能与现有 request-id / correlation 关联」。已核对：`requestid.Middleware` 包在 mux 外层（server.New 链 = requestid → security → operational gate → 路由），因此 Wrap 执行时上下文必已携带合法 `requestid`；现有 API 合同用 `correlation_id` 为 JSON body 字段名（requestid.BodyName）。

### 决定

**§1 关联判据（退出 2 可核对形式）**

开启 tracing 时，每个 wrap 处理的 HTTP 请求生成的 server span 必须携带属性 `correlation.request_id`，其值 = `requestid.FromContext(r.Context())`（= 传入的 `X-Request-ID`，或中间件生成的新 id）。判据核对方式：带 `X-Request-ID: <id>` 请求 → span 属性等于 `<id>`；不带 → span 属性等于响应头 `X-Request-ID` 值。

**§2 属性名与 baggage**

- 冻结属性名：`correlation.request_id`（不采用裸 `request_id`，用命名空间前缀避免与将来 correlation.* 属性冲突）。
- W3C baggage 键：`request-id`（小写短横线，与 header 名语义一致）。开启 tracing 时将该键写入请求的 baggage 上下文，使下游（未来）instrumented 调用无需额外配置即可继承关联 id。本仓当前无出站调用，注入为惰性就绪（由同 ctx 提取断言验证）。
- 不修改 `requestid` 包、不改变 HTTP 响应头行为、不把 request-id 加入 metrics 标签（R1 §4 标签白名单仍封闭）。

**§3 边界**

- 无效/空 id 不出现（requestid.Valid 已由中间件保证）；Wrap 层不做二次校验，直接读取。
- span 属性写入失败不影响请求（旁路语义，同 R3）。
- 不重开 VP-012 correlation/错误包络合同。

### 为什么

- span 属性是 tracing 侧关联的最短闭环；baggage 键与 header/字段既有命名惯例一致，为跨服务传播预留。
- 关联判据可机械核对（header 值 == 属性值），正好落在 R5 双路径证据的 trace 侧。

### 未选方案

- **仅 baggage 不写属性**：独立 span 在 collector 里看不到关联值，判据不可核对。
- **把 request-id 写进 metrics 标签**：违反 R1 §4 白名单（该值高基数且含可识别性），明确拒绝。
- **改 requestid 包注入 span**：requestid 属 transport 层不该依赖 otel；在 obs 拦截点读取上下文已足够。