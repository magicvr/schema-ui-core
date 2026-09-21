---
id: E-002-r3a-go-wire-sweep
doc: execution-entry
status: active
parent: GOAL-006-r3-wire-formatter-and-unit-family-matrix
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-002 · R3-A 第一段：Go 侧 wire formatter 落码与全量替换

## 事实

- **单一 formatter 落码**：
  - `internal/temporal/temporal.go` 新增 `FormatWire(t) string`（= `NewValue(t).String()`，UTC + 固定 6 位 + `Z`，**向零截断**）；
  - `internal/handler/rfc3339.go` 重写为：`FormatWireTime` / `FormatWireTimePtr`（nil 不伪造瞬时）/ `ParseWireTime`（D-005 兼容矩阵：0/3/6/9 位与 offset → UTC；无时区拒绝）+ `ErrWireTimeInvalid`；
  - `modules/wallet/jobs.go` 改用 `temporal.FormatWire`（模块不得向上依赖 handler）。
- **全量替换（inventory §「当前 inline fixed-3 outputs」与「Current RFC3339 outputs」的 Go 面）**：
  - inline `2006-01-02T15:04:05.000Z07:00` 布局（15 个文件、约 40 处）→ `FormatWireTime(...)`；nullable 字段按 `*ptr` 解引用（7 处，均在既有 nil 守卫内）。
  - `time.RFC3339` **输出** → `FormatWireTime`：`invites.go`（2）、`scheduledtasks.go`（nextRuns）、`service_credentials.go`（5，含 nullable 解引用）、`telegram_operator.go`（4）、`cmd/schema-ui/configpkg.go`（`ExportedAt` + `# imported_at:` 头）。
  - **输入 parser 统一**：`jobs.go`、`operations.go`、`service_credentials.go` 改调 `ParseWireTime`（集中 D-005 规则）。
  - **按 `D-009` 明确保留在 wire 合同之外**（不改为固定 6 位）：邮件正文（`invites.go` 正文两处、`authsession/email_identity.go`、`authsession/recovery.go`）、审计 `detail` JSON（`service_credentials.go:218`）、recycle payload 解析（`modules/recyclebin/service.go:280`）。
- **Go fixture 同步**：`rfc3339_test.go` 重写为三组测试（输出契约含「截断不舍入」「非 UTC 归一化」；指针缺省语义；D-005 输入矩阵含拒绝项）；`recyclebin_test.go` / `wallet_voucher_test.go` / `server_restart_test.go` / `resources_test.go` 的 3 位断言改为固定 6 位。

## 证据

- `go build ./...` exit 0；`go test -count=1 ./internal/handler/ ./cmd/schema-ui/` 全绿。
- `TestFormatWireTime`（4 子例）/ `TestFormatWireTimePtr` / `TestParseWireTimeCompatibility`（6 接受 + 5 拒绝）全绿。
- 全仓 `go test -count=1 ./...`：见 `02-execution.md` 事实边界。

## 未完成（诚实边界）

- **Web 侧 fixture 同步**（inventory §Web consumer：新增 fixed-6 用例；现有 52 处 3 位字面量是**输入** fixture，解析器容忍变宽小数，故不强制改写）。
- **`I-041-008`（wire 输出破坏性/兼容期）尚未判定**：已实测前端 `apps/web/src/lib/datetime.ts` 的正则接受变宽小数（`D-005` 输入容忍），但需补固定 6 位的正向用例后才能判定无需兼容期。
- **R3-B 未开始**：单位族矩阵（秒/毫秒/可空/sentinel 各至少一个 endpoint）与 VP-020 会话时区展示 round-trip。
- 检查点 A 因此**尚未完成**，`progress` 保持 0/3。
