---
id: GOAL-002-metrics-export-contract
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
---

## D-001 · 指标导出合同与配置面冻结（R1）

### 触发

Root GOAL-001-observability 纲领阶段 R1 开工。VP-015 I-015-001（指标面合同）与 I-015-003（分母）为 required，最晚需要阶段 = R1 合同冻结；I-015-004（readyz）最晚 R2/R3 接入前。已核对现状：主 HTTP 面 `:25080` 单 mux（`internal/server` + `requestid.Middleware` 外层包装）；配置体系 YAML + `${VAR}` 插值 + KnownFields fail-closed（`internal/config`）；`/healthz` `/readyz` 语义已由 GOAL-003/GOAL-009 冻结。

### 决定

以下各项即本波（架构 A4 / VP-015）指标导出合同，R2 起按此实施：

**§1 暴露面形态 —— 独立 listener，不挂主 mux**

新增独立 metrics listener，与主 HTTP server 完全分离：

| 配置键 | env 覆盖 | 默认 | 语义 |
|--------|----------|------|------|
| `observability.metrics.enabled` | `OBSERVABILITY_METRICS_ENABLED` | `false` | 显式打开；缺省完全不监听 |
| `observability.metrics.addr` | `OBSERVABILITY_METRICS_ADDR` | `127.0.0.1:25081` | 独立 listener 绑定地址 |
| `observability.metrics.auth_token` | `OBSERVABILITY_METRICS_AUTH_TOKEN` | 空 | 可选 Bearer token；设置后所有请求须带 `Authorization: Bearer <token>` |

`enabled=false`（含全缺省）时进程不创建该 listener、零行为变化——这是「无收集器仍能开发快测」（VP 意图 3）的合同基础。

**§2 鉴权与暴露守卫（fail-closed）**

- token 为 secret：只能经 `${VAR}` 插值 / env / `configs/.env` 提供，不得写 YAML 字面量（同 `DB_PASSWORD` 约定）。
- 设置时长度 ≥ 16 字符（去空白后），否则启动失败。
- **非 loopback 绑定必须配 token**：`enabled=true` 且 addr 的 host 不是 loopback（空 host / `0.0.0.0` / `::` / 非 loopback IP / 非 `localhost` 域名均视为非 loopback）且 token 为空 → 任何环境下启动失败。loopback 绑定允许无 token（本机 scrape / SSH 隧道场景合法）。
- 死凭证拒绝：`auth_token` 有值但 `enabled=false` → 启动失败（对齐 `storage.objects.s3.*` + `driver=local` 的既有先例，防止死 secret 滞留配置）。
- 校验在 `Load` 与 `ValidateProd` 双侧生效（零值测试 Config 全缺省时跳过），错误信息只引用键名，绝不回显 token 值。

**§3 内核最小系列（命名前缀 `suc_`）**

| 系列 | 类型 | 标签 | 说明 |
|------|------|------|------|
| `suc_build_info` | gauge=1 | `version`, `commit`, `go_version`, `profile` | 进程构建身份 |
| `suc_http_requests_total` | counter | `module_id`, `method`, `route`, `status` | 主 mux 全部注册路由的请求计数 |
| `suc_http_request_duration_seconds` | histogram | `module_id`, `method`, `route` | 固定默认 bucket 集 |
| `suc_kernel_modules_enabled` | gauge=1 | `module_id` | 每 enabled 模块一行 |
| （Go runtime / process collectors） | 标准集 | — | `go_*` / `process_*`，有界标准集 |

`suc_` 为 schema-ui-core 缩写前缀。

**§4 标签契约（秘密禁令）**

标签白名单封闭：`module_id`、`method`、`route`、`status`（HTTP 系列）+ §3 静态标签。禁止把 query、header、raw path、用户/租户 ID、任何凭证或 PII 放进标签。`route` 一律取 ServeMux 注册 pattern（含 `{id}` 通配段），永不取实际请求路径。新系列/新标签必须以新 D 记录修订本合同。

**§5 基数边界**

`route` ∈ 已注册 pattern 集合（有限）；`method`/`status` 天然有界；直方图 bucket 固定。基数上界 ≈ 注册路由数 × 方法 × 状态码类，随路由清单线性、可控。不存在无界标签源。

**§6 模块贡献面（本波）**

模块**不**新增自定义业务指标义务：已启用模块的路由经统一 instrumentation 自动携带 `module_id`（所有权来自 kernel.Provider `ContributionIdentity`）。架构 §2.2「Observability 按需」语义不变。「内核 vs 模块最小集合」= 内核提供 face + HTTP instrumentation + build_info/modules_enabled/runtime；模块本波零额外贡献。

**§7 本波分母（闭合 Root I-003 / VP I-015-003）**

Store 连接池统计、对象存储、Job runner 指标**不进本波退出分母**（后续增量，须新 D 记录扩展合同）。HTTP request span 已由 VP-015 退出判据 2 冻结进 tracing 分母（归 R3/R4）。理由：A4 退出判据仅要求 ≥1 条内核或模块路径可 scrape 且系列携带 `module_id`；收窄分母防范围膨胀。

**§8 readyz 边界（提前闭合 Root I-004 / VP I-015-004）**

metrics listener 与后续 OTLP 导出器**均不进 readyz**；`/readyz` 语义保持不变（store ping + module graph Start+Ready + 显式对象后端 probe）。导出面是旁路观测设施而非服务依赖；未显式配置不得影响启动、就绪或 mvp/dev 缺省路径。

**§9 实现载体**

引入官方 exposition 库 `github.com/prometheus/client_golang`（v1.23.x，Prometheus 文本格式 + registry 语义正确性）。备选「手写 text format」因转义/解析边缘正确性风险被拒。OpenTelemetry SDK 选型**不在本决定范围**（I-002 未闭合前不动 traces 面；本波不加 `observability.traces.*` 配置键）。

### 为什么

- 独立 listener 让「默认仅本地回环」成为结构保证而非运行时判断；主 addr 默认 `:25080` 绑全部接口，若 /metrics 挂主 mux，显式打开后即随主端口进入反向代理拓扑，暴露边界难审。
- 非 loopback 必 token 是可机械核对的暴露守卫；比 IP allowlist 简单，不引入与 W7 trusted_proxies 的耦合。
- 分母与 readyz 边界沿用 VP 文本自身的建议（未显式配置不扩 ready；A4 最小集合），避免过度交付。
- client_golang 为 Prometheus 官方客户端，缓存已有 v1.23.2，网络代理可达；依赖成本一次付清，正确性收益长期。

### 未选方案

- **主 mux 挂 `/metrics` + IP allowlist**：allowlist 需信任代理链判断，复杂且易与 CORS/trusted_proxies 语义纠缠；无法落实默认 loopback。
- **永远跟随 http.addr 复用端口**：同上，且无法在不开 token 时满足「不对公网暴露」。
- **手写 Prometheus 文本格式**：省一个依赖，但转义、escaping、TYPE/HELP 边缘案例的正确性风险自担；不值得。
- **本波同时冻结 traces 配置键名**：I-002 尚 open，先冻结名会诱导实现绕过 R3 门禁；traces 键留给 R3 决策。
