---
id: GOAL-006-dual-path-evidence
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-001 · R5 双路径证据方案（判据映射 + 收集步骤 + otlp-sink）

### 触发

R1–R4 全部落地（GOAL-002～005 done）；R5 需把 VP-015 退出判据 3/4 与 Root 成功标准 3/4 变为可核对证据，并支撑 Root 独立关门审计（项目决策 `docs/architecture/independent-audit-execution.md`：self 后直接调本地 grok build `/audit`）。

### 决定

**§1 判据映射**

| VP/Root 判据 | 证据形态 | 收集点 |
|--------------|----------|--------|
| 退出 3 / 成功 3：未配置收集器时本地/Compose 默认仍能开发与快测 | 无 observability 配置启动：`/healthz`+`/readyz` 200；**无** 25081（metrics 默认）与 4318（collector 默认）监听；启动日志无 observability 相关项 | 本目标 E-002 缺省路径 |
| 退出 4 / 成功 4：显式配置后 metrics scrape **与** ≥1 trace 导出可核对 | ①启用 metrics+tracing 启动；②`GET /metrics` 实测 `suc_build_info` / `suc_http_requests_total{module_id,route}` / `suc_kernel_modules_enabled`；③带 `X-Request-ID` 请求后真实 OTLP sink 收到 span 导出 POST | 本目标 E-002 显式路径 |
| 退出 2 关联判据（R4 已锁） | `X-Request-ID` == span 属性（单测锁定；live 侧以 sink 收包 + 既有单测佐证） | 引用 GOAL-005 测试 |

**§2 工具：`apps/api/cmd/otlp-sink`**

极简 OTLP/HTTP 接收缘（POST 任意路径、计字节、逐行日志），仅用于证据与运维排障；不实现采集语义。env `OTLP_SINK_ADDR`（默认 `:4318`）。

**§3 判定标准**

证据满足 §1 三行即判据成立；任一失败 → 标记为 finding 修复后重跑。命令序列以 E-002 记录为准（可重复）。

**§4 关门顺序**

GOAL-006 A-001（self）→ Root A-001（self 关门核对，成功标准逐条对照）→ grok build `/audit`（independent）→ 合并响应 → 无未闭合 required 后 Root `done`。

### 为什么

- 用 live 证据而非仅单测，满足「生产向验收以显式配置为准」的 VP 语义；命令序列落 E 条目，任何协作者/CI 可重跑（N-004 固化）。
- otlp-sink 独立于 obs 测试（httptest sink 只证明包内行为），避免「用被测对象的测试设施当证据」的循环取证。

### 未选方案

- **仅引用既有单测/live 冒烟**：R1（metrics）/R3（traces）的冒烟证据分散且未记录为双路径对照；R5 需要同一次显式配置运行里两条路径同时成立。
- **引入真实 collector（Prometheus/Jaeger/OTel Collector）**：重部署依赖，超出「显式配置后可核对」的最低证据要求；sink 已足够判定导出行为。