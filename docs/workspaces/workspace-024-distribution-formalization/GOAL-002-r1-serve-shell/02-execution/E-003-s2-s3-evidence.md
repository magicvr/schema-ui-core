---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-serve-shell
version: 0.1.0
---

# E-003 · S2/S3 实现与验收证据（2026-08-29）

## 交付物（commit 随本目标核销）

- `apps/api/server/config.go` + `config.default.yaml`：公开 config 装载（RT-K01 子集；`${VAR}`/`${VAR:-default}` fail-closed；KnownFields；env 定向覆盖；validate fail-closed）。
- `apps/api/server/serve.go`：`Options` / `Serve`（信号版）/ `Run(ctx, opts, signals)`（可测版）；标准组合 7 模块（server-registration · auth-session · manifest-route · navigation-capability · schema-render · operationlog · users）；中央面（healthz/readyz/登录、RegisterSchemas、RegisterManifest + Bootstrap）；timeouts + request-id + nosniff/CORS；RT-D02 §1 全序停机。
- `apps/api/server/config_test.go` + `serve_test.go`：13 项单测（fail-closed、插值、方言配对、env 覆盖、Run 干净排空 + healthz/readyz/manifest/登录 200）。
- `apps/api/cmd/schema-ui/main.go`：新增 `serve` 子命令（-config/-dialect/-dsn/-addr）；`apiVersion` 钉 `v0.4.0`（serve 面所在下一发布）；模板新增 `config.yaml.tmpl`；`main.go.tmpl` 重写为薄封装（-config/-dialect/-dsn 兼容语义）；README 模板更新。

## 验收证据

| 级别 | 内容 | 结果 |
|------|------|------|
| 单元 | `go test ./server/ -count=1`（13 项：默认/显式文件/env 覆盖/插值 fail-closed/非法 shutdown_timeout/dialect 配对/非 dev 密钥 fail-closed/未知键拒绝/Run 排空） | **PASS** |
| 回归 | `go test ./...`（apps/api 全量，含新 `server` 包） | **全绿**（exit 0） |
| E2E-L1 | `schema-ui serve -addr 127.0.0.1:25107 -dsn <tmpdb>`（SQLite）：/healthz 200 · /readyz 200 · POST /api/auth/login 200（dev 种子）· manifest 200 | **PASS** |
| E2E-L2 | `schema-ui create demo-e2e`（14 文件：thin main.go + config.example.yaml + …）→ replace 本地点缀下 `go mod tidy` + `go build` + `go run ./cmd/server` → healthz/readyz/login/manifest 200 | **PASS**（骨架可直接 serve 启动） |
| E2E-L3 | docker postgres:16（15433→5432）→ `schema-ui serve -dialect postgres -dsn …` → /healthz 200（迁移台账 apply） | **PASS**（双方言） |

## 残余与登记（有界口径）

1. **信号级 drain harness**（SIGTERM → exit 0/1 断言）：Windows 本环境不可执行（同 VP-021 先例「进程级 harness（!windows）CI 核销登记」）——由 `Run(ctx 取消)` 干净排空单测 + 主仓 RT-D02 契约（composition 既有 harness A/B/C）覆盖语义；**登记 = `compose CI 实跑`（R3）时在 linux runner 补齐**。
2. **registry 级骨架消费**：模板 `go.mod` 钉 `v0.4.0`（含 serve 面）；全 registry 语义（无 replace）消费实证随 **R2 公开发布**（tag apps/api/v0.4.0 → golden-field 升级）核销——本波以本地点缀运行级验证兜底。

## 验收判据映射（GOAL-002 成功标准）

| 判据 | 证据 | 结论 |
|------|------|------|
| C1（serve 子命令 + 骨架可直接启动） | E2E-L1/L2 | ✅ |
| C2（薄封装 + -config/-dialect/-dsn 保留） | E2E-L2（编译+运行）· 单测 | ✅ |
| C3（RT-D02 全序停机 / 预算/退出码） | 单测（ctx 干净排空）+ config fail-closed；信号/退出码 → CI 登记 | ✅（有界） |
| C4（双方言） | E2E-L1（sqlite）+ E2E-L3（postgres） | ✅ |
| C5（config fail-closed + 探针登录路径） | config_test 十三项 + E2E-L1/L2 login 200 | ✅ |