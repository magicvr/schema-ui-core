---
id: GOAL-003-metrics-scrape-endpoint
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
---

## D-001 · R2 scrape 接入实施接缝

### 触发

R1 合同已冻结（GOAL-002 D-001）；R2 开工需把合同映射到本仓代码结构。已核对：composition 唯一 mux 装配点 `newMuxWithExtraProviders`；全部中心注册走 `mux.Handle`/`mux.HandleFunc`（handler 包 31 处，含 `auth.go` HandleFunc）；模块路由经 `kernel.RegisterContributions` 返回 `RouteContribution`（内嵌 `ContributionIdentity.ModuleID`）。

### 决定

**§1 新包 `internal/obs`（不依赖 config/kernel/composition）**

- `Observer`：持有私有 `prometheus.Registry` 与固定系列（D-001-R1 §3 全集 + Go/process collectors）；`BuildInfo{Version,Commit,GoVersion,Profile}` 静态标签。
- `InstrumentedMux`：内嵌 `*http.ServeMux`，覆写 `Handle` **与** `HandleFunc`（auth.go 使用后者），按注册 pattern 包装 handler；`Own(pattern, moduleID)` 声明所有权，缺省 `core`。
- `Server`：可选专用 listener（enabled=false 时 Start 为 no-op、不创建任何资源）。
- 依赖方向：obs → prometheus/client_golang + pkg/version；composition → obs；obs 不 import config（参数用原始值传入），避免反向耦合。

**§2 instrumentation 标签规则**

- `method` / `route` 由注册 pattern 派生（首个空格切分；无 method 前缀时 method=""），永不取 `r.URL.Path`——通配 pattern（如 `/api/files/{id}`）保证基数有界。
- `status` 记录 WriteHeader 显式值，未调用时隐式 200。
- 所有权：contributed routes 在装配循环以 `route.ModuleID` 调 `Own`；中心注册（health/upload/branding/manifest/schema/auth…）默认 `core`。

**§3 listener 生命周期与失败语义**

- fx lifecycle OnStart 在 readiness gate 置位后启动 metrics listener；**显式 enabled=true 而 bind 失败 → 启动失败**（fail-closed：显式配置即硬要求；缺省关闭仍是「非硬依赖」的保证面）。
- Serve 运行期意外错误：error 日志但**不杀进程**（旁路观测面不得拖垮主服务；与主 listener 的 fail-closed 语义有意不同）。
- OnStop 与主 server 一并 Shutdown，错误 join 进停机链。
- 不影响 `/readyz`（R1 D-001 §8）：metrics listener 的健康不由 readyz 表达，也不进 readyz 探针。

**§4 Bearer 校验**

`Authorization: Bearer <token>`，`subtle.ConstantTimeCompare` 恒时比较；scheme 大小写不敏感；失败回 401 + `WWW-Authenticate`。

**§5 build info 来源**

`version.Version` / `version.Commit`（ldflags 可覆盖）+ `runtime.Version()` + `cfg.ProfileName`；`suc_kernel_modules_enabled` 在组装时按 `plan.IDs()` 逐个 Set(1)。

### 为什么

- 覆写式 InstrumentedMux 是唯一能同时覆盖中心注册与模块贡献的单一拦截点（否则要在 handler 包 8+ 个 Register* 函数里逐一穿参，侵入面大且易漏）。
- bind 失败 fail-closed 对齐仓库一贯姿态（配置错误应当响亮）；运行期降级仅日志则守住「旁路面不杀主服务」边界。
- obs 不 import config 保持包层级干净，Server/Observer 可独立单测（config 无关）。

### 未选方案

- **在 requestid/server 层做全局中间件**：拿不到匹配后的 route pattern（Go 1.22+ `r.Pattern` 在 mux 匹配后才填充），要么标签退化为 raw path（违反 R1 §4），要么需要二次匹配 hack。
- **每个 Register\* 函数显式传 Observer**：侵入 8+ 处签名，漏一处即静默丢指标。
- **复用 prometheus 默认全局 registry**：全局状态使测试隔离与多实例组装困难；私有 registry 是 client_golang 推荐做法。
