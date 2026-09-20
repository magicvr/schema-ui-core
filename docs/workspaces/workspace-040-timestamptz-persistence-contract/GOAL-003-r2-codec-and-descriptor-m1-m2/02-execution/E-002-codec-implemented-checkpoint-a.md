---
id: E-002-codec-implemented-checkpoint-a
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-003-r2-codec-and-descriptor-m1-m2
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-002 · 共享 codec 落码与检查点 A 完成

## 事实

1. **`I-041-001` 定稿并落盘**：本目标首条决策 `01-decision/D-001-temporal-codec-api.md`（`accepted`）——包路径 `apps/api/internal/temporal`、**只依赖标准库**、导出面最小化、错误模型四个哨兵、`Parse` 兼容面与 `Format` 严格面**刻意不对称**、写路径强制 `Truncate`、并明确「拒绝非零 offset 属 R3 wire 层、不在本包」。

2. **codec 落码**：`apps/api/internal/temporal/temporal.go`（新包）：
   - `FromUnix` / `FromUnixMilli`：语义**严格等于** `time.Unix(sec,0).UTC()` / `time.UnixMilli(ms).UTC()`（`D-018` 对拍基线，**不得偏离**）；
   - `Layout = "2006-01-02T15:04:05.000000Z"`（27 字符）、`CanonicalLen = 27`；
   - `Format`（严格：UTC + 固定 6 位 + `Z`；年份越界 → `ErrRange`）、`MustFormat`（限已校验调用点）；
   - `Parse`（兼容：0/3/6/9 位小数、`Z`/`±HH:MM` 偏移归一化 UTC；无时区 → `ErrNoZone`；小数位非法 → `ErrFraction`；非法日历值 → `ErrFormat`）；
   - `Truncate` = `t.UTC().Truncate(time.Microsecond)`（向零，含 epoch 前负值）；
   - `Value` / `NullValue` 写路径包装（自动 `Truncate`；`NullValue` 缺失时 `String()` 返回 `"NULL"`）。

3. **可执行单测**：`temporal_test.go`，**11 个测试全绿**（`go test ./internal/temporal/ -count=1` → `ok`）：
   `TestFromUnixMatchesGo`、`TestFromUnixMilliFloorSemantics`（含 `-1/-999/-1000/-1001/-86400000/-1758320000123/253402300799999` 逐项断言，并**同时**与 `time.UnixMilli` 比对）、`TestFormatCanonicalShape`、`TestFormatRejectsUnrepresentableYear`、`TestTruncateTowardZero`（含负 instant）、`TestTruncateNormalizesZone`、`TestParseAcceptsCompatibleInput`、`TestParseRejects`（4 子测试）、`TestRoundTripCanonical`、`TestValueAndNullValue`、`TestCanonicalConstants`。

4. **依赖边界实测**（不靠代码内自证，靠构建层查询）：
   ```text
   go list -f "{{range .Imports}}{{.}} {{end}}" ./internal/temporal/   =>   errors fmt time
   go list -deps ./internal/temporal/ | grep 'database/sql|pgx|sqlite|mattn|modernc'   =>   无匹配
   ```
   → **未引入任何驱动包或 `database/sql`**；Root 红线「不把驱动类型泄漏进公共面」在 codec 处成立。

5. **质量门**：`gofmt -l` 无输出、`go vet ./internal/temporal/` exit 0、`go build ./...` exit 0。

6. **本轮自纠两处自身错误**：
   - `temporal_test.go` 初稿有一个 `if got := ...; got != want` 的遮蔽赋值，**编译失败** → 已修；
   - `TestParseAcceptsCompatibleInput` 初稿把 4 种小数位形式当作同一时刻比对，**语义错误**（它们是**不同**时刻）→ 已改为逐例断言 + 另设「等价拼写归一后一致」子断言；
   - 初稿的 `TestNoDriverImports` **名不副实**（并未检查 import）→ 已改名 `TestCanonicalConstants`，并在注释中写明「无驱动依赖」由**构建层查询**验证、验证命令记入本条。

## 证据

- 代码：`apps/api/internal/temporal/temporal.go`；测试：`apps/api/internal/temporal/temporal_test.go`。
- 决策：`GOAL-003/01-decision/D-001-temporal-codec-api.md`。
- 命令与结果：见上文 §3–§5（`go test` / `go vet` / `go list` / `go build` 的实际输出）。
- 权威对位：`GOAL-002/attachments/r1-c2-column-contract-draft-v0.1.md` §1/§3；`r1-public-wire-inventory-v0.1.md` L17；**Root** `D-008`/`D-015`；`GOAL-002` `D-018`。

## 状态评估

- **检查点 A = completed**（codec 落码 + 单测全绿 + 依赖边界实测）；`I-041-001` = **verified**。
- `progress: 0/4 → 1/4`（A 完成；B/C/D pending）。
- **尚未触及 v1–v72**：本轮只新增一个 `internal` 包，**无迁移改动、无 schema 变更**。
- **下一步 = 检查点 B**：15 个 conversion descriptor 的 SQLite `Apply` 落码（按 `D-019` F-5 / 裸四步；逐表 DDL 取 `GOAL-002/attachments/r1-c2-per-table-rebuild-ddl-v1.0-fc.md`），并保持 v1–v72 canonical SQL/checksum 逐条不变。
- **未做的独立审计**：本检查点尚未送 `/audit`；按计划在检查点 D（首次记录哈希）发起，或按风险提前。
