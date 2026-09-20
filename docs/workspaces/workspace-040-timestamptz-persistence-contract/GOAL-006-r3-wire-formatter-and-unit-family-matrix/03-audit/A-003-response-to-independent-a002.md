---
id: A-003-response-to-independent-a002
doc: audit-entry
status: active
parent: GOAL-006-r3-wire-formatter-and-unit-family-matrix
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# A-003 · 编排器响应：independent A-002（`fail`，3 条 required）

- **source**: `self`（编排器响应，非独立意见）
- **日期**: 2026-09-20
- **scope**: 响应 `A-002`（grok build · grok-4.6 · reasoning high，verdict `fail`，开放 required = 3）的**全部** findings，并说明检查点 A/B 状态变更
- **verdict**: `pass`（A-002 的 3 条 required 全部 `fixed`，2 条 recommended 亦已 `fixed`；开放 required = 0）

## 响应总表

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| `F-I-001` PUT `/api/mail/config` 仍走默认时间编解码 | required/high | **fixed** | `handler/mail_admin.go:157` → `mailConfigWire(view)`；`TestMailConfigPutUpdatedAtIsCanonicalWire` |
| `F-I-002` `GET /api/mfa/status` `enrolledAt` 是 map 里的裸 `time.Time` | required/high | **fixed** | `handler/mfa.go:200-209`；`TestMFAStatusEnrolledAtIsCanonicalWire`；两条常驻守卫 |
| `F-I-003` `ParseWireTime` 接受非零 offset，违反 `D-005` | required/high | **fixed** | `handler/rfc3339.go:51-76`；`TestParseWireTimeCompatibility`（`+08:00` 移入拒绝集）+ `TestParseWireTimeRejectsNonZeroOffsetsOnEveryInboundParser` |
| `F-I-004` 可空/sentinel 族身份不可从分母核对 | recommended/med | **fixed** | `TestUnitFamilyProvenanceIsMachineChecked`（分母 unit + 生成描述符 `NonNull: false` + DDL 无 `NOT NULL`） |
| `F-I-005` Web round-trip 自实现表面积 | recommended/low | **fixed** | `utc-roundtrip.test.tsx` 新增「helper ↔ 生产 `formatDate`」4 时区对拍表 |

## 1. `F-I-001` · **fixed**

审计事实成立：GET 已投影而 PUT 仍 `writeJSON(w, http.StatusOK, view)`，同一资源的两半会对同一瞬时发出不同宽度。这是编排器第一轮修复的**不完整**，不是审计误判。

- **改动**：`mailConfigPut` 成功路径改为 `writeJSON(w, http.StatusOK, mailConfigWire(view))`（与 GET 共用同一投影）。
- **回归**：`TestMailConfigPutUpdatedAtIsCanonicalWire`。断言不是「形状看起来对」而是**确定性等价**：`mail.Switcher.Update` 用 `time.Now()` 打戳并把同一个值写入 `mail_config`，因此「响应的规范形」必须逐字节等于「库里保存的规范值」。默认编解码器只有在**一位小数都不丢**时才可能碰巧相等，用例连做两次切换（第二次 `mockRetention` 不同，保证是新的 `now`），把巧合概率压到可忽略。
- 其余键（`channel`/`mockRetention`/`resend`/`smtp`/`secrets`）在该路径也未丢。

## 2. `F-I-002` · **fixed**

审计事实成立且属**编排器扫描方法的缺陷**：上一轮按「时间语义键白名单」扫 map，键表里没有 `enrolledAt`；而 `enrolledAt` 又是函数返回值，不是 `time.Time` 字段声明，因此两类扫描都漏。

- **改动**：`mfa.go` 的 status 投影改为：`row["enrolledAt"] = nil`（缺席时）或 `FormatWireTime(enrolledAt)`（在场时）。未注册**不再**输出伪造的 `"0001-01-01T00:00:00Z"`，与 Web 端 `enrolledAt: string | null` 一致。
- **测试夹具**：`fakeMFAService` 增加 `statusEnrolledAt` seam（可配置在场瞬时），使「缺席 → `null`」与「在场 → 逐字节 `.900000Z`」两个方向都可断言。`TestMFAStatusEnrolledAtIsCanonicalWire` 覆盖两者。
- **补扫描（按审计要求改为「方法 × 载体」闭集，且落成常驻守卫）**：新增 `internal/w040contracttest/wire_bypass_guard_test.go`：
  - `TestNoJSONTaggedWireTimeFieldOutsideProjections`（守卫 1）：扫描**全部非测试 Go 源**，任何 `time.Time`/`sql.NullTime` 字段带 `json:` tag 都必须出现在**带理由的 allowlist** 中，否则测试失败；allowlist 条目失效也会失败（防过期）。这直接封住「结构体直接交给 `writeJSON`」这一类（即 `F-I-001`）。
  - `TestNoUnformattedResponseTimeValue`（守卫 2）：扫描**全部非测试 Go 源**中形如 `*_at` / `*At` / `timestamp` 的键，其值必须包含 `Format`（或为 `nil`/布尔/字符串/数字字面量），否则失败。键模式比上一轮手写表宽（含 `enrolledAt`），且不再依赖「只用 GET」——它按源码扫描，方法无关。
  - **决定性验证**：临时放入一个合成违规文件（`EnrolledAt time.Time \\`json:"enrolledAt"\\`` + `map[string]any{"enrolledAt": instant}`）→ 两条守卫**各自**失败并点名该文件；随后删除（`Test-Path` = False）。不是空洞断言。
- **其余载体的处置（本轮重新逐类核对）**：
  - 非测试 `writeJSON` 调用点 120 处；其中真正可能承载时间的**必然**是「时间键 map」或「含时间字段的 struct/DTO」两类，二者现在都由上述守卫变成**构建期失败**。
  - 非测试 `json.Marshal` 调用点（`internal/handler`）：`bootstrap.go:63`（manifest 文档，无时间字段）、`export.go:141,147`（CSV 单元格编码的是**已投影**的字符串）、`resources.go:101`（资源描述符 JSON，无时间）。`MarshalJSON` 自定义实现全仓 = 0。
  - `internal/store/recovery.go:40` marker、artifact 文件名时间戳、`modules/recyclebin/service.go` 的 payload→domain 结构体初始化：分别为内部产物 / 文件名 / `D-009` 例外，已在守卫中**显式**标注理由（守卫 2 的 allowlist 只有 `modules/recyclebin/service.go` 一条，且注明是 6 处 `timeField(...)` 结构体字段而非响应 map）。

## 3. `F-I-003` · **fixed**（按已冻结决策实现，**不**构成新决策）

审计事实成立。核对原文后确认这**不是**「两种都说得通」的选型，而是编排器上一轮实现偏离了已冻结决策：

- Root `D-005` 原文第 3 条：「非法时间、**非零 offset 解析失败**；不接受模糊本地时间字符串。」
- `internal/temporal/temporal.go:140-141` 的包文档本身就写明：「Rejecting non-zero offsets is a transport-layer policy that belongs to the API input parser (**R3**, GOAL-003 D-001 §3), not to this package.」
- inventory C2 同义（`rejects non-zero offsets`）。

因此按**已冻结决策**实现拒绝（审计给出的 preferred 方案），**没有**静默改写 `D-005`/`GOAL-003 D-001`/inventory：

- `ParseWireTime` 先由 `temporal.Parse` 做权威格式校验（0/3/6/9 位、必须有 zone），再在**原始串**上取 offset 并要求为 0（`time.Parse(time.RFC3339Nano, ...)` 只用于读 offset，因为 `temporal.Parse` 会归一 UTC 而掩盖它）。`Z`、`+00:00`、`-00:00` 仍然接受。
- 测试：`+08:00` 从接受集移入拒绝集，新增 `-00:00` 接受项与 `TestParseWireTimeRejectsNonZeroOffsetsOnEveryInboundParser`（含「Go 本身能解析，所以拒绝是我们的策略」这一前提断言）。
- 影响面核对：三个入站解析器（`jobs.go`、`operations.go`、`service_credentials.go`）共用该函数；`operations.go` 的 `YYYY-MM-DD` 日期过滤走**独立分支**，不受影响（UI 的日期区间过滤仍可用）；仓内客户端只发 `Z` 形态或日期串。codec 层（`internal/temporal`）继续接受非零 offset 供存量 payload 解析，其自身测试未改。

> P-004 说明：审计给出的备选是「用户书面修订 `D-005`」。由于冻结决策本身要求拒绝、且实现该拒绝不改变任何用户裁决，本项按 `fixed` 闭合，无需再次询问用户。

## 4. `F-I-004` · **fixed**

新增 `TestUnitFamilyProvenanceIsMachineChecked`，把「族身份」从散文变成三条可执行检查：

1. 四族所用 8 个列必须存在于冻结分母且 unit 与矩阵声明一致（`temporalcontract.Columns()`）；
2. sentinel 族的身份来自**生成描述符**本身：`modules/corepersistence/migration/vp040_temporal.go` 必须含 `{Table: "mail_config", Column: "updated_at", NonNull: false}`（D0 转换把 legacy 0 变 NULL，故列可变空）；
3. 可空族的身份来自**表 DDL**：`service_credentials.revoked_at` / `last_used_at` 的 DDL 行不得含 `NOT NULL`。

（未给 `temporalcontract.Column` 增加标志位：该类型与 `Count = 90` 是从 `r1-time-column-inventory-v0.3.md` 机械生成的冻结分母，改它等于改冻结制品；改用「生成描述符 + DDL」作为旁证，同样可执行且不动冻结物。）

## 5. `F-I-005` · **fixed**

`utc-roundtrip.test.tsx` 新增 4 时区对拍表：对同一 wire 瞬时，**生产** `formatDate(value, "zh-CN", {timeZone})` 必须包含期望的 24 小时 `HH:MM`，且测试内 `wallClockIn` 给出同一 `HH:MM`（UTC 12:57 / Asia/Shanghai 20:57 / America/New_York 08:57 / Asia/Kathmandu 18:42）。自实现若与生产展示不一致会在此失败，不再只靠自洽。

## 6. 状态与门禁

- **开放 required = 0**（`F-I-001`/`002`/`003` 全部 `fixed`，`F-I-004`/`005` 亦已修）。
- **检查点 A 现在成立**：formatter（E-002）+ fixture 与三处漏网字段（E-003、本响应）+ 输入矩阵**按 `D-005` 真实语义**（Z / `+00:00` / `-00:00` 接受，非零 offset 与无时区拒绝）+ `I-041-008`（Root `D-019` §1）。`progress` 仍为 `2/3`，检查点 C 待复审与关门。
- **`I-041-008` / `I-040-004`**：A-002 的结论（均同意 verified）未受本轮修正影响，编排器接受。
- **B**：A-002 判 `达成`；本轮对 Web round-trip 的加强（`F-I-005`）只增证据。
- **未决冲突 / P-004**：`F-I-003` 已按冻结决策闭合，**无**需用户裁决的冲突项。
- **请求复审**：按 `GOAL-005` 先例（A-002 → A-003 修复 → A-004 independent 复审 → A-005 关门），下一轮请 grok build 对修复后的提交做**定向复审**（`A-004`），确认 3 条 required 已合法闭合且无新 required；在此之前不关 `GOAL-006`、不宣布检查点 C 完成。

## 证据（可复跑）

- `cd apps/api && go build ./...` → exit 0。
- `go test -count=1 ./internal/handler/ -run "TestFormatWireTime|TestParseWireTime|TestDefaultTimeMarshaller|TestHealthProbe|TestMailOutbox|TestMailConfig|TestMFAStatus|TestR3B"` → 全绿。
- `go test -count=1 ./internal/w040contracttest/` → 守卫 3 项 + 既有边界矩阵全绿；合成违规文件验证两条守卫各自失败后已删除。
- 全仓 `go test -count=1 ./...` 与 `apps/web` `vitest run` 见 `02-execution.md`。
