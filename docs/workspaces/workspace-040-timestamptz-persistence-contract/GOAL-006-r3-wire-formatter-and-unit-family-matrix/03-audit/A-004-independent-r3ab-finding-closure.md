---
id: A-004-independent-r3ab-finding-closure
doc: audit-entry
status: active
parent: GOAL-006-r3-wire-formatter-and-unit-family-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# A-004 · independent 定向复审：A-002 required 闭合（F-I-001～F-I-003）

- **source**: `independent`
- **auditor**: grok build · grok-4.6 · reasoning high
- **provider**: grok build
- **日期**: 2026-09-21
- **scope**: `GOAL-006` 检查点 A 的定向复审（A-002 的 3 条 required + 2 条 recommended 处置；守卫不可绕过性；`ParseWireTime` 反例；PUT 回归决定性；`enrolledAt` 契约；治理一致性）。不含 R3-C/D，不重开 R2。基线：上一轮 `acc51aab` → 本轮 HEAD `28b9db62`（`fix(w040-r3): close the independent A-002 required findings and guard the class`）；工作树 clean。
- **类型**: finding-closure
- **verdict**: `pass`

> 本意见不修改 `status`/`progress`/方案正文/goal-tree。响应由 `/govern` 处理。本会话按用户指令**不落盘**；编排器应原样写入 `03-audit/A-004-independent-r3ab-finding-closure.md` 并更新 `03-audit.md` 索引。

## 范围与区间

只复核 A-002 三条 required 是否**合法闭合**，并按用户清单主动攻击守卫、PUT 回归、offset 拼写、MFA 契约与治理。共享资料目录为 `none`。未读取其他工作区作为状态源。

## 逐条闭合判定表

| ID | 级别 | 判定 | 依据 |
|----|------|------|------|
| `F-I-001` | required/high | **fixed** | 生产路径已改：`mail_admin.go:157-159` PUT 成功走 `mailConfigWire(view)`，与 GET `mail_admin.go:105` 共用投影。`mailConfigWire`（`mail_admin.go:86-96`）对非空 `UpdatedAt` 调 `FormatWireTime`，nil 保持 JSON `null`。回归 `TestMailConfigPutUpdatedAtIsCanonicalWire`（`w040_r3a_public_wire_fields_test.go:285-328`）存在且本轮 `go test ./internal/handler/` 绿。**测试并非 A-003 所声称的「巧合可忽略」**——见新 recommended `F-I-101`；缺陷本身可从源码直接核对，不阻断 `fixed`。 |
| `F-I-002` | required/high | **fixed** | `mfa.go:205-209`：缺席 `enrolledAt: nil`，在场 `FormatWireTime(enrolledAt)`。`TestMFAStatusEnrolledAtIsCanonicalWire`（`w040_r3a_public_wire_fields_test.go:335-373`）锁未注册 `null` 与 trailing-zero `.900000Z`。Web 类型仍是 `string \| null`（`mfa-manager.tsx:29`），UI 只读 `status?.enabled`（L255-257），测试 mock 已用 `null`。`user_mfa.created_at` 为 `NOT NULL` 且 insert 打戳（`modules/mfa/store/repository.go:88-97`、`vp040_temporal.go:100`），不存在「已注册但 `CreatedAt` 为零值」的真实路径。常驻守卫可抓住**本条那一类**漏网（见反例记录）；闭集主张过满 → `F-I-102`（recommended），不把本条打回未闭合。 |
| `F-I-003` | required/high | **fixed** | `ParseWireTime`（`rfc3339.go:62-74`）先 `temporal.Parse`，再在原始串上 `carriedOffsetSeconds` 要求 offset==0；逗号分隔符单独拒绝。`+08:00` 已从接受集移入拒绝集（`rfc3339_test.go:86-92`）。一次性探针（见下）覆盖用户点名的拼写：**没有**「codec 接受 + RFC3339Nano 读不到非零 offset」的漏网；`Z`/`+00:00`/`-00:00` 接受；非零 offset 拒绝。三个入站调用点仍只走该函数（`jobs.go:389`、`operations.go:111`、`service_credentials.go:170`）。`operations.go:105-110` 的 `YYYY-MM-DD` 先于 `ParseWireTime`，activity 筛选仍是 `datePicker`（`activity.json:110-117`）。Jobs UI 无 from/to 日期字段。凭据测试发 `time.RFC3339` 的 UTC `Z`（`service_credentials_test.go:22`）。仓内调用方未被误伤。`-00:00` 归零等价符合 A-002 关闭要求原文（保留 `Z`/`+00:00`/`-00:00`）与 `D-005`「`+00:00` 等等价 offset」。按已冻结决策实现，**不是**新决策，无需 P-004。 |
| `F-I-004` | recommended/med | **fixed** | `TestUnitFamilyProvenanceIsMachineChecked`（`wire_bypass_guard_test.go:246-295`）：分母 `temporalcontract.Columns()` unit + 生成描述符 `{Table: "mail_config", Column: "updated_at", NonNull: false}` + DDL 行无 `NOT NULL`。未改冻结 `Column` 类型，旁证可执行。不阻断 B（A-002 已判 B 达成）。 |
| `F-I-005` | recommended/low | **fixed** | `utc-roundtrip.test.tsx:199-211` 对 UTC/Shanghai/New_York/Kathmandu 把生产 `formatDate(WIRE, "zh-CN", {timeZone})` 与 `wallClockIn` 的 `HH:MM` 对拍（含 Kathmandu `18:42`）。自实现 helper 仍在，但已与生产 renderer 交叉约束。本轮 vitest 该文件 12 tests 绿。 |

## 新 findings

### F-I-101 · recommended · med

**问题**：`TestMailConfigPutUpdatedAtIsCanonicalWire` **不能**兑现 A-003「默认编解码器只有在一位小数都不丢时才可能碰巧相等 / 巧合概率可忽略」的主张。A-002 关闭要求写的是「trailing-zero 瞬时逐字节固定 6 位」；实现改成「响应规范形 == `SELECT updated_at` 扫成 string」。

**证据**：

- `Switcher.Update`（`runtime.go:409-428`）用**未截断**的 `time.Now().UTC()` 写入，并把同一 `now` 放进内存 view（不重读库）。SQLite bind（`store.go:242-245`）经 `temporal.NewValue` 写成 canonical TEXT；测试 `Scan` 进 `string` 不走 time adapter（`scan.go:86-89` 只拦截 `*time.Time`/`NullTime`）。
- 因此：**正确路径**下 `FormatWireTime(now)` 与 stored 恒等（探针 200/200）。
- **错误路径**（`json.Marshal(now)` vs stored）：本机一次性探针 **20/200 = 10%** 相等。两次循环的联合假通过约 1%，不是「可忽略」。
- GET 回归 `TestMailConfigUpdatedAtIsCanonicalWire` 已用 `trailingZeroInstant`；PUT 没有。成功 PUT 不可能返回 `updated_at: null`（Update 必打戳），故 PUT-NULL 不必补。

**关闭要求（可选，不阻断 A/C）**：PUT 回归改为注入/冻结一个 trailing-zero 瞬时（或断言响应**不是** `json.Marshal` 的短形，例如 `.9Z` / 无小数 `Z`），不要只靠 `got == stored`。

生产路径仍是 `mailConfigWire`，本条不把 `F-I-001` 打回未闭合。

### F-I-102 · recommended · med

**问题**：两条常驻守卫**各自能抓住 A-002 那一类漏网**（本审用自己的合成文件复现：两条都失败并点名文件），但**不是** A-003 所说「时间到公共 wire 的路径已变成构建期闭集」。`formattedValue` 含 `Format` 子串即通过，过宽；键模式与「必须有 json tag」过窄。

**证据（合成文件，已删除，工作树 clean）**：下列六类同时存在时，守卫 1 **和** 守卫 2 **都 PASS**（即漏过）：

| 类 | 构造 | 为何漏 |
|----|------|--------|
| A | 导出字段 `EnrolledAt time.Time`（无 json tag） | 守卫 1 正则要求 `json:"..."` |
| B | `map[string]any{"EnrolledAt": instant}` | 守卫 2 的 `*At` 要求小写开头 |
| C | `"created": instant` | 键不是 `*_at`/`*At`/`timestamp` |
| D | `"enrolledAt":\n\tinstant`（键后换行） | `formattedValue("")` 为 true |
| E | `"enrolledAt": notFormatted`（变量名含 `Format`） | `strings.Contains(value, "Format")` |
| F | `inner["enrolledAt"] = instant`（运行时赋值） | 字面量 `key:` 扫描看不到 |

全仓只有一处 `writeJSON`（`health.go:123-127`）。本轮**未找到**上述类别的现行 HTTP 漏网：`filelibrary` 的 `"created"` 已是 `FormatWireTime` 字符串（`filelibrary.go:86,121,157`）；`RecycleItem` 无 tag 但走 `recycleItemToMap`（`recyclebin.go:169-184`）；`[]any` 命中均为 SQL args。守卫只扫 `apps/api`：仓内 HTTP JSON 生产者只在该模块；`apps/web` 无 Go `time.Time`；`scripts/` 是打包脚本。

**关闭要求（可选）**：不要把「两类守卫绿」写成闭集证明。若要加强：扫无 tag 导出 `time.Time` 是否进入 `writeJSON`、禁止 `Contains("Format")` 当充分条件、覆盖 `"created"`/`"modified"` 等键、或对 `writeJSON` 实参做浅层类型检查。有现行漏网再升 required。

## 反例尝试记录

### 1. 守卫合成文件（自己的，不采信编排器叙述）

- **命中文件** `a004_guard_hit.go`（`EnrolledAt time.Time \`json:"enrolledAt"\`` + `map[string]any{"enrolledAt": instant}`），临时放入 `apps/api/internal/handler/`：
  - `TestNoJSONTaggedWireTimeFieldOutsideProjections` **FAIL**，点名 `apps/api/internal/handler/a004_guard_hit.go: EnrolledAt time.Time \`json:"enrolledAt"\``。
  - `TestNoUnformattedResponseTimeValue` **FAIL**，点名 `a004_guard_hit.go:12 (enrolledAt): return map[string]any{"enrolledAt": instant}`。
- **绕过文件** `a004_guard_bypass.go`（上表 A–F）：两条守卫均 **PASS**。
- 用后删除。`git status`：nothing to commit, working tree clean。

### 2. `ParseWireTime` 一次性探针（`TestA004OffsetProbe`，已删）

`codec` = `temporal.Parse`；`std` = `time.Parse(RFC3339Nano)`；`wire` = `ParseWireTime`。

| 输入 | codec | std | wire | 备注 |
|------|:-----:|:---:|:----:|------|
| `…Z` / `+00:00` / `-00:00` / `.123Z` / 6 位 Z / 9 位 Z / 9 位 `+00:00` | ✓ | ✓ | ✓ | 零等价 |
| `+08:00` / `-05:00` / 9 位 `+08:00` | ✓ | ✓ | ✗ | 非零 offset，拒绝正确 |
| 小写 `z`、`+0000`、`+0800`、`+08`、前后空白、`24:00:00`、`:60` 闰秒、小写 `t`、`+00:00:00`、`UTC` 后缀、无时区、空格分隔 | ✗ | ✗ | ✗ | 进不了 offset 探针 |
| `+00:60` / `+24:00` | ✓ | ✓ | ✗ | Go 把非法 offset **溢出**成 3600/86400 秒；`carriedOffsetSeconds` 读到非零，wire 拒绝。不是漏网 |
| 逗号小数 `,123Z` | ✓ | ✓ | ✗ | 与测试前提一致：codec 容忍，wire 拒绝 |
| `.12Z`（2 位） | ✗ | ✓ | ✗ | fraction 门先挡 |

**无**「`temporal.Parse` 接受且 RFC3339Nano 失败」的发散行，因此 `carriedOffsetSeconds` 在成功路径上「parse 失败返回 0」不可达。没有需要升 required 的 offset 漏网。

### 3. PUT 巧合探针（`TestA004PutCoincidenceProbe`，已删）

`n=200`，每次 `time.Now().UTC()` + 1ms：`FormatWireTime==stored` **200/200**；`json.Marshal==stored` **20/200（10%）**。见 `F-I-101`。

### 4. 入站调用方 / Web 取值形态

- 三个 parser 均调用 `ParseWireTime`；无平行的 `time.Parse(RFC3339)` 入站残留（本轮 grep）。
- Activity `from`/`to`：`datePicker` → `YYYY-MM-DD` → `parseOperationTime` 日期分支，**不**进 offset 拒绝。
- Jobs 搜索 schema 无 from/to 日期控件；空 query 对 `parseJobTime` 为无界。
- 凭据创建：仓内测试发 UTC `RFC3339`（`Z`）；Web 无向该 API 发 `+08:00` 的表单。Wallet `expiresAt` 走独立 `parseVoucherExpiry`（Unix / `YYYY-MM-DD`），不是 `ParseWireTime`。
- Host `return-intent` 的 `toISOString()` 是浏览器 Host 协议，正则已要求 `Z`（`failure.ts:93`），不经 API parser。

### 5. 基线复跑（本审）

- `cd apps/api && go test -count=1 ./internal/handler/ ./internal/w040contracttest/ ./internal/temporal/` → ok（98.059s / 2.627s / 0.209s）。
- `cd apps/web && vitest run` 三个指定文件 → 3 files / **34** tests ok（utc-roundtrip 现 12 tests，含 F-I-005 对拍）。
- 未复跑全仓 64 包 / 全量 vitest。未跑迁移/Docker。

## 治理一致性

- `00-meta.md` 检查点 A=completed、B=completed、C=pending，`progress: 2/3` 由 A～C 等权派生。A-002 曾不接受 A=completed；**本轮接受**（3 条 required 已合法闭合）。C 仍 pending，与「未静默关门」一致。
- `A-003` 对三条 required 的**代码修复**有证据；「PUT 回归决定性」「守卫 = 闭集」两句过满，降为 `F-I-101`/`F-I-102`，不构成「required 未闭合却放行」。
- `03-audit.md` 索引含 A-001 / A-002 / A-003，待补 A-004。无「以测试通过代替审计结论」的关门（A-003 明确请求本轮复审，GOAL-006 仍 `active`）。
- goal-tree `GOAL-006 active · 2/3` 与 meta 一致。`workspace.md` 过时属已知，不升 required。

## 明确写出

- **开放 required 数量：0**（`F-I-001`/`002`/`003` 均为 `fixed`）。新 findings 两条均为 **recommended**，不阻断检查点 A，不阻断进入 C。
- **是否同意「检查点 A 成立」：同意。** formatter + inventory 替换 + 本轮三处漏网已修 + 输入矩阵按 `D-005`（`Z`/`+00:00`/`-00:00` 接受，非零 offset 与无时区拒绝）+ `I-041-008`。
- **是否同意 `I-041-008` / `I-040-004` 维持 verified：同意。** 本轮修正未引入 3 位输出消费者，也未削弱 VP-020 展示矩阵（还加了 production `formatDate` 对拍）。
- **是否同意可以进入检查点 C 的关门（GOAL-006 静默关门）：同意进入 C。** 本意见是 C 所需的 independent 输入；编排器落盘本条、响应 recommended（可修可 residual）、把 C 标 completed 并写关门记录后，可以静默关 `GOAL-006`。本意见本身不改 `status`。关门仍不等于 Root 关门（R3-C/D）。
- **P-004：无。** `F-I-003` 按已冻结 `D-005` / GOAL-003 `D-001` §3 实现拒绝，不是新裁决。`-00:00` 不与 `D-005` 冲突。recommended 的补强不需要用户先裁才能进 C。

## 最小补齐动作（不阻断 C）

1. （可选）PUT 回归改 trailing-zero 或显式排除默认短形（`F-I-101`）。
2. （可选）收敛守卫或收回「闭集 / 构建期失败」措辞（`F-I-102`）。
3. 编排器落盘本意见为 `A-004`，更新 `03-audit.md` 索引；然后走检查点 C / `/govern` 关门。

不要把「两类守卫绿」或「`got == stored`」单独当成闭集/决定性证明；生产三处 required 缺陷本身已经可核对闭合。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
