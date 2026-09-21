---
id: A-001-self-r3ab-checkpoints
doc: audit-entry
status: active
parent: GOAL-006-r3-wire-formatter-and-unit-family-matrix
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# A-001 · self 审计：R3-A/R3-B 检查点（formatter、fixture、单位族矩阵、VP-020 round-trip）

- **source**: `self`
- **日期**: 2026-09-20
- **scope**: `GOAL-006` 检查点 A（R3-A）与 B（R3-B）；边界依据 Root `D-003`/`D-005`/`D-009`/`D-018` 与 `GOAL-002/attachments/r1-public-wire-inventory-v0.1.md`；不含 R3-C/D。
- **verdict**: `conditional`（**开放 required = 0**；2 条 recommended 已在本轮修正，4 条 recommended/note 提交 independent 复核）

## 成果（逐项可核对）

| # | 判据 | 结论 | 证据 |
|--:|------|------|------|
| 1 | 单一 shared fixed-6 formatter 落码并替换 inventory 的 Go 面 | **达成** | `internal/temporal/temporal.go` `FormatWire`；`internal/handler/rfc3339.go`；E-002 逐文件清单（15 文件约 40 处 inline + `time.RFC3339` 输出） |
| 2 | 输入端兼容矩阵（0/3/6/9 位 + `+00:00` + 拒绝无时区）有可执行测试 | **达成** | `rfc3339_test.go` `TestParseWireTimeCompatibility`（6 接受 + 5 拒绝）；`TestFormatWireTime`（4 子例，含「截断不舍入」与「非 UTC 归一化」） |
| 3 | Go/Web fixture 同步 | **达成（有界）** | Go：`rfc3339_test.go` 重写 + 4 个测试文件 3 位断言改 6 位；Web：`lib/datetime.test.ts`、`host/return-intent.test.ts` 正向固定 6 位用例；E-003 §1 |
| 4 | `I-041-008` 判定 | **达成** | Root `D-019` §1（用户 P-004：无破坏性、不需兼容期），证据表 6 条 |
| 5 | 单位族矩阵：秒/毫秒/可空/sentinel 各 ≥1 endpoint | **达成** | `w040_r3b_unit_family_matrix_test.go` 四子测试，四表四端点；族身份由 `requireDenominatorColumn` 从冻结分母读回 |
| 6 | VP-020 会话时区展示 round-trip（Go + Web） | **达成** | `w040_r3b_timezone_roundtrip_test.go`（4 时区 × 6 瞬时，含 DST 边界、+05:45、负 epoch、`time.Local` 独立性）；`apps/web/src/i18n/utc-roundtrip.test.tsx`（8 用例，L1/L2/L3 层解析 + 秒粒度恒等）；`I-040-004` → verified |
| 7 | 非 DB 文本未被纳入固定 6 位（`D-009`） | **达成** | 邮件正文、audit `detail`、recycle payload 解析器保留原状（E-002 明列）；本轮未新增例外 |
| 8 | v1–v87 canonical SQL/checksum 不可变 | **未触碰** | 本轮只改 `internal/handler`（3 个文件）与 Web 测试；`git show --stat 34e9941c` 可核对无 migration/descriptor 文件 |

## 偏差（本轮发现并已修正）

| ID | 级别 | 偏差 | 处理 |
|----|------|------|------|
| `F-S-001` | recommended | `mailConfigWire` 首版把 `nil` 视图投影为零值对象，而修复前的 `writeJSON(view)` 对 nil 指针编码为 JSON `null` —— 属行为漂移 | **fixed**：投影改为返回 `*mailConfigResponse`（nil 保持 nil），并加 `TestMailConfigWireProjectsExactlyOneUpdatedAt` 断言 `null` 保真 + 单 `updated_at` 键 |
| `F-S-002` | recommended | `apps/web/src/lib/datetime.ts` 文档注释仍以 `2026-08-01T18:02:44.000Z`（3 位）举例，与本轮 wire 合同变更后的事实不符 | **fixed**：注释示例改为 `...44.000000Z` |

## 提交 independent 复核的事项（不自行放行）

| ID | 级别 | 事项 | 独立复核点 |
|----|------|------|-----------|
| `F-S-003` | recommended | sentinel 族子测试以直接写 NULL 表达「v73 已把 legacy 0 转成 NULL」的**转换后状态**，本身不重演 0→NULL 转换 | 该子测试的断言范围是否被正确限定为「wire 输出」；是否会把 R2 的转换证据误读为本轮新证据（指针：`internal/w040contracttest/migration_boundaries_test.go`、`internal/backup/verify.go:185-203`） |
| `F-S-004` | recommended | E-003 §2 的「闭集」是**两类扫描 + 推理**，不是不可能性证明：`time.Time` 字段若无 json tag、或时间值嵌在 `[]any`/嵌套 map 中，可同时躲过两类扫描 | 能否构造第三类「时间到达 wire」的路径反例；已核对 `json.Marshal` 调用点与全部 wire map 均走 `FormatWireTime` |
| `F-S-005` | note | 51–52 处 `.NNNZ` 与 pinned upstream JSON 未改写，理由为「输入 fixture + provenance pin」 | 是否存在把 3 位小数**输出**当断言的测试（若有，`I-041-008` 结论需推翻） |
| `F-S-006` | note | Web round-trip 的 `zoneOffsetMs`/`instantFromWallClock` 是测试内自实现，理论上可自洽地掩盖时区错误 | 该实现对 DST 边界与 +05:45 的处理是否真能失败（三个硬编码期望值 + 组件侧与生产 formatter 对拍是否足以兜底） |

## 已核对为「非问题」的项（供审计反证）

| ID | 项 | 依据 |
|----|----|------|
| `N-001` | `internal/store/recovery.go:40` `VerifiedAt time.Time json:"verifiedAt"` 仍走默认编解码 | 内部 marker 产物（SQLite sidecar / PG `COMMENT`），**非 HTTP wire**；形状由 R2 冻结（`D-017`），改它属重开 R2 |
| `N-002` | artifact 文件名时间戳 `20060102T150405.000Z`（`store/migrate.go:337`、`store/postgres.go:243`、`store/recovery.go:397`）保持 3 位 | **文件名**，非 wire；R2 冻结 |
| `N-003` | `/healthz` `timestamp` 由 `time.Time` 改 `string` 是否破坏消费者 | 仓内无解析 body 的消费者：`compose.yaml:58` 健康检查用 `wget` 只看退出码；`cmd/server/server_restart_test.go:210-221` 与 `shutdown_harness_test.go:104` 只看状态码；Web 无调用点 |
| `N-004` | `modules/recyclebin/service.go:277` 仍是 3 位布局解析器 | `D-009` 非 DB 例外（任意 payload 字段），R3 非目标 |
| `N-005` | `apps/web/dist-lib/**` 仍含 3 位示例文本 | `.gitignore:113` 忽略的构建产物（`git ls-files` = 0），非受控制品 |

## 自审可复跑证据

- `cd apps/api && go build ./...` → exit 0。
- `go test -count=1 ./...` → **64/64 包 ok**（含真实 PostgreSQL 15.4 路径的 `internal/backup`/`internal/composition`）。
- 定向：`go test -count=1 ./internal/handler/ -run "TestDefaultTimeMarshallerBreaksTheWireContract|TestHealthProbeTimestampIsCanonicalWire|TestMailOutboxCreatedAtIsCanonicalWire|TestMailConfigWireProjectsExactlyOneUpdatedAt|TestMailConfigUpdatedAtIsCanonicalWire|TestR3BUnitFamilyWireMatrix|TestR3BWireRoundTripIsZoneIndependent|TestR3BSessionTimezoneDoesNotMoveStoredOrSentInstants"` → PASS。
- `cd apps/web && vitest run` → **124 文件 / 1504 用例 ok**。

## 自审结论

R3-A/B 的功能性判据全部有可执行证据，开放 required = **0**；两条 recommended 已修正，四条事项（`F-S-003`～`F-S-006`）**不自审放行**，提交 grok independent 复核。检查点 C 的成立以独立意见（`A-002`）与编排器响应（`A-003`）为准。
