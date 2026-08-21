---
id: GOAL-002-metrics-export-contract
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
---

## E-002 · observability.metrics 配置面实施与验证（checkpoint）

### 事实

按 D-001 §1–§2 实施配置面（仅配置加载与校验，listener 归 R2）：

- `internal/config/config.go`：
  - `Config` 新增 `MetricsEnabled` / `MetricsAddr`（默认 `127.0.0.1:25081`）/ `MetricsAuthToken`；
  - `yamlFile` 新增 `observability.metrics.{enabled,addr,auth_token}`（KnownFields 严格解析）；
  - env 映射 `OBSERVABILITY_METRICS_ENABLED` / `OBSERVABILITY_METRICS_ADDR` / `OBSERVABILITY_METRICS_AUTH_TOKEN`；
  - `validateMetrics` + `validateObservability`：死 token 拒绝、addr 必须 host:port 且数字端口 1-65535、token ≥16 字符、**任何环境**非 loopback 绑定必须 token、错误只引键名不回显值；`Load` 与 `ValidateProd` 双侧生效，零值测试 Config 全缺省时跳过。
- `internal/config/config.default.yaml` 与 `configs/config.yaml`：新增 `observability:` 段（默认关、loopback、注释契约说明）。
- `internal/config/config_observability_test.go`：13 个子测试（缺省全关 / YAML 启用 loopback / 通配与非 loopback 绑定拒绝 / 插值 token 通过 / 死 token / 短 token / 非法 addr / 非数字端口 / env 覆盖 / 零值 ValidateProd 容忍 / 手工 Config 非 loopback 拒绝 / loopback 默认免 token）。

### 验证

- `gofmt -l internal/config` → 仅 `config_test.go` 命中（**既有漂移，本切片未触碰该文件**）；本切片文件全部干净。
- `go vet ./internal/config/` → exit 0。
- `go build ./...` → 通过；`go test ./...`（apps/api 全仓）→ 无 FAIL。

### Git checkpoint

| hash | scope |
|------|-------|
| `45489f4` | `apps/api/internal/config/config.go`、`config.default.yaml`、`config_observability_test.go`、`apps/api/configs/config.yaml`（4 files，+280） |

### 备注

- traces 配置键按 D-001 §9 未加（I-002 保持 open，归 R3）。
- listener 接线（fx lifecycle、registry、instrumentation）= R2 子目标范围。
