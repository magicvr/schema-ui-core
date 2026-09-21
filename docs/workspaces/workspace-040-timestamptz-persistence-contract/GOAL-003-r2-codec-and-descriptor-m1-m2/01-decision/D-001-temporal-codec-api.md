---
id: D-001-temporal-codec-api
doc: decision-entry
status: accepted
parent: GOAL-003-r2-codec-and-descriptor-m1-m2
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-001 · 共享时间 codec 的公共 API 形态（I-041-001 定稿）

## 决定的来源

- 本目标检查点 **A**（M1）要求「共享 codec 落码 + 单测全绿（含 `D-018` 的 Go 对拍用例）」；信息项 **`I-041-001`**（required，影响门禁 A/M1）要求定稿 codec 的公共 API 形态。
- 权威来源（**不得改写**）：
  - `GOAL-002/attachments/r1-c2-column-contract-draft-v0.1.md` §1（canonical 值合同）；
  - 同 §3「Codec boundary」（领域面只暴露 Go `time.Time`；**禁止** `pgtype`/driver 类型/raw storage string 泄漏）；
  - `GOAL-002/attachments/r1-c2-c3-guardrails-v0.1.md` §1；
  - **Root** `D-015`（写入统一 `UTC().Truncate(time.Microsecond)`；毫秒族已更正）；
  - **Root** `D-008`（精度与 0/D0 政策）；
  - `GOAL-002/attachments/r1-public-wire-inventory-v0.1.md` L17（**输入**兼容合法 RFC3339 的 0/3/6/9 位小数与 `+00:00` 等等价 offset；**输出**固定 6 位）；
  - `GOAL-002/D-018`（行拷贝用纯整数 SQL；Go codec 作**对拍验收**基线）。

## 决定

### 1. 包位置与依赖边界

- 包路径：**`apps/api/internal/temporal`**（`internal`，模块私有；不对外发布）。
- **只依赖标准库**（`time`、`errors`、`fmt`）。**禁止**导入任何驱动包、`database/sql` 或 `pgtype`。
- 包内**不**出现任何驱动/SQL 类型；导出面只用 `time.Time`、`string` 与自定义错误。

### 2. 导出面（最小）

```go
package temporal

// 单位显式，禁止按数值量级推断（合同 §1「never inferred from magnitude」）
func FromUnix(sec int64) time.Time          // 语义 = Go time.Unix(sec, 0).UTC()
func FromUnixMilli(ms int64) time.Time      // 语义 = Go time.UnixMilli(ms).UTC()

const Layout = "2006-01-02T15:04:05.000000Z" // 固定 6 位；27 字符

func Format(t time.Time) (string, error)    // 严格输出：固定 6 位 UTC，非法即 error
func Parse(s string) (time.Time, error)     // 兼容输入：0/3/6/9 位小数、offset、拒绝无时区

func Truncate(t time.Time) time.Time        // t.UTC().Truncate(time.Microsecond)，向零
func NewValue(t time.Time) Value            // 写路径：Truncate

type Value struct{ ... }                    // 非空时间值
func (v Value) Time() time.Time
func (v Value) String() string              // 已截断后的固定 6 位

type NullValue struct{ ... }                // 可空时间值（symmetric with sql.NullTime）
func NewNullValue(t time.Time) NullValue
func (n NullValue) Time() time.Time
func (n NullValue) String() string          // Valid=false → "NULL"
```

### 3. 错误模型

```go
var ErrRange    = errors.New("temporal: value out of supported range")
var ErrFormat   = errors.New("temporal: value is not canonical fixed-6 UTC RFC3339")
var ErrNoZone   = errors.New("temporal: timestamp has no timezone offset")
var ErrFraction = errors.New("temporal: unsupported fractional-second precision")
```

- 上层（迁移/仓储）须用 `fmt.Errorf("...: %w", err)` 包装并附 **表/列/行** 证据（C3 边界的 `m0` 失败语义）；**迁移期非法值一律 fail closed，不静默转换**。
- `Parse` 的**兼容面**（0/3/6/9 位小数 + `±HH:MM`/`Z`）与 `Format` 的**严格面**（仅 `Z`、仅 6 位）**刻意不对称**——这是合同要求的：入站宽容、出站规范。
- **非零 offset 的处理**：`Parse` **接受并归一化到 UTC**（满足 Web 显示解析需求，见 wire inventory L71）。**API 入站解析器必须拒绝非零 offset** —— 该拒绝属 **R3 的 wire formatter 工作**，不在本包内实现（`D-016` 范围修正）。

### 4. 与 `D-018` 的关系（对拍验收）

- `FromUnix` / `FromUnixMilli` 是 `D-018` 要求的**Go 对拍基线**：R2 的迁移必须证明 SQL 转换结果与这两个函数逐行等价。
- 因此本包**不得**改变 `time.Unix`/`time.UnixMilli` 的语义（例如不得自行 round）；`FromUnixMilli(-1)` 必须等于 `time.UnixMilli(-1).UTC()`。

### 5. 截断语义

- `Truncate` = `t.UTC().Truncate(time.Microsecond)`：`time.Time.Truncate` 对负值**向零截断**（与 **Root** `D-008`/`D-015` 一致）。
- **写路径**（`NewValue`/`NewNullValue`）自动 `Truncate`；**读路径**（`FromUnix`/`FromUnixMilli`/`Parse`）不额外截断（源本身无亚微秒余数，或由 `Parse` 归一化）。

### 6. 边界与拒绝规则

| 情形 | 行为 |
|------|------|
| 年份 < 0 或 > 9999 | `Format` → `ErrRange`（固定 4 位年无法表示） |
| 非 UTC（`t.Location() != time.UTC`） | `Format` **不**拒绝；内部先 `UTC()` 归一（合同：存储只存 UTC，归一属实现义务） |
| 已有亚微秒余数且非 `Truncate` 来源 | `Value` 构造时已截断；`Format` 对裸 `time.Time` 仍按 6 位输出（`time.Format` 会**四舍五入**到 6 位，故文案要求：**经 `Value`/`Truncate` 之后再 `Format`**；直接 `Format` 裸值不保证向零语义，文档须写明） |
| 无时区 / 本地时间串 | `Parse` → `ErrNoZone` |
| 小数位非 0/3/6/9 | `Parse` → `ErrFraction` |
| 非法日历值（如 2 月 30 日） | `Parse` → 底层 `time.Parse` 错误（包装为 `ErrFormat`） |
| 非规范存储文本（`+00:00`、空格、变长小数） | 持久化边界拒绝（`ErrFormat`）——见合同 §1 |

## 未选方案

- **把 codec 放在 `kernel`**：会让内部实现细节进入平台公共面；合同要求领域面只暴露 `time.Time`。未采用。
- **在 codec 内实现 API 入站的「拒绝非零 offset」**：该策略属 wire 层（**R3**），放进 codec 会让存储层承担传输层策略。未采用。
- **隐式 `Truncate` 于 `Format`**：会让 `Format` 变成有损操作且掩盖调用方漏截断。未采用（改为文档强制「先 `Value` 再 `Format`」）。
- **用 `int64` 作为领域值**：违反 Root 红线（驱动/存储表示不得进入领域面）。未采用。

## 影响与边界

- 本决策定稿 **`I-041-001`**；其状态由 `collecting` → **`verified`**（本目标内实现并单测后）。
- **未经用户裁决的技术细节**（函数命名、错误哨兵具体措辞）由后续 independent 审计复审；若审计要求变更，按 P-004 回到用户。
- 本决策**不**改变任何 `GOAL-002` 决策；`D-017`（单 checksum）、`D-018`（行拷贝）、`D-019`（F-5）、`D-020`（测试载体）、`D-021`（residual）均不受影响。
