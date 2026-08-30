---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-serve-shell
version: 0.1.0
---

# E-002 · S1 设计冻结完成（2026-08-29）

1. **用户 P-004 裁决**：serve 面构成 = **方案 A · 标准下游组合闭环**（公开 `server` 包：config 装载 + 标准组合装配 + 中央面接线 + RT-D02 停机）；模板形态 = **方案 A · 薄封装单一形态**（`cmd/server` = flag 解析 + `server.Serve`；组合代码移入公开面）。
2. **D-001 落盘**：决策正文、未选方案（B 全量对等 / C 最小壳；模板双入口）、信息门禁 I-001/I-002 → verified。
3. **实现落位**（本阶段已随冻结同步落地，S2/S3 并作）：
   - `apps/api/server/config.go`：`LoadConfig(path)`（代码默认 → YAML `${VAR}`/`${VAR:-default}` 插值 fail-closed → env 定向覆盖 → validate fail-closed）；内嵌 `config.default.yaml`。
   - `apps/api/server/serve.go`：`Options{Config, Logger, Store, ExtraProviders}` · `Serve`（信号版）/`Run(ctx, opts, signals)`（可测版）——标准组合（7 模块）装配、中央面（RegisterWithMFAProbes / RegisterSchemas / RegisterManifest + Bootstrap）、timeouts + request-id + nosniff/CORS（镜像 internal/server）、RT-D02 §1 全序停机。
   - `apps/api/cmd/schema-ui/main.go`：新增 `serve` 子命令（-config/-dialect/-dsn/-addr）；`apiVersion` 对齐 `v0.3.0`；模板新增 `config.yaml.tmpl`；`main.go.tmpl` 重写为薄封装；README 模板更新。
4. **待验证**：`go build ./...` + `go vet ./server/` + `go test ./server/`（S2/S3 验收，见 E-003）。