---
id: E-003-r3a-web-fixtures-and-wire-field-gap
doc: execution-entry
status: active
parent: GOAL-006-r3-wire-formatter-and-unit-family-matrix
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-003 · R3-A 第二段：Web fixture 同步 + 三个漏网的公共 wire 时间字段

## 1. Web fixture 同步（inventory §Web consumer）

- `apps/web/src/lib/datetime.test.ts`：新增固定 6 位与 legacy 3 位正向用例（同一瞬时两种小数宽度渲染一致）、非零 offset 的固定 6 位等价用例、**无时区**固定 6 位必须 `null`（不猜本地时区）。
- `apps/web/src/host/return-intent.test.ts`：新增 `validateReturnIntent` 的 wire 形状兼容组——接受固定 6 位、legacy 3 位、整秒 `Z`；拒绝无时区；固定 6 位的过期瞬时仍判 expired。
- **未改写** `apps/web/src/protocol/upstream/*.json`：这些是**带 sha256 pin 的上游 provenance fixture**（`HOST_FAILURE_FIXTURE_SHA256` 等），属外部 producer 制品，改写会破坏 pin 且超出本目标范围。
- **未改写** 52 处 `.NNNZ` 字面量：全部是测试**输入**（Web 解析器容忍变宽小数），不是输出断言；改写它们不增加任何覆盖，反而扩大 diff。真正的固定 6 位正向覆盖由上面两个文件承担。

## 2. 发现的缺口：三个公共 wire 时间字段绕过共享 formatter

R3-A 第一段（E-002）按 inventory 的逐行列清单替换了 inline 布局与 `time.RFC3339` 输出，但 inventory 存在**类别级遗漏**：字段本身是 `time.Time` 结构体/DTO 成员时，输出由 `encoding/json` 的默认时间编解码器完成，源码里没有任何 `Format(...)` 可替换，因此逐行列清单搜不到。

**证据（可复现）**：对固定 6 位存储值，`time.Time.MarshalJSON` 会去掉小数尾零：

```text
ns=0            -> {"updated_at":"2026-09-20T12:57:15Z"}
ns=100000000    -> {"updated_at":"2026-09-20T12:57:15.1Z"}
ns=900000000    -> {"updated_at":"2026-09-20T12:57:15.9Z"}
ns=914522000    -> {"updated_at":"2026-09-20T12:57:15.914522Z"}
ns=914520000    -> {"updated_at":"2026-09-20T12:57:15.91452Z"}
nil             -> {"updated_at":null}
```

宽度可变（0–6 位），直接违反 `D-003` 固定 6 位输出合同。命中三处：

| # | 字段 | 端点 | 修复 |
|--:|------|------|------|
| 1 | `handler.healthResponse.Timestamp time.Time` | `GET /healthz`、`GET /readyz`（4 条分支） | 改为 `string` + `healthBody(status)` 统一经 `FormatWireTime` |
| 2 | `handler.toMapItems` 的 `"created_at": rec.CreatedAt` 与 `outboxDetail` 直接 `writeJSON(rec)` | `GET /api/mail/outbox`、`GET /api/mail/outbox/{id}` | 新增单一投影 `outboxWireRecord(rec, keepEmptyBody)`，两条路径共用（`keepEmptyBody` 保留各自已冻结的 body 形状：列表恒带 body、详情保留 `omitempty`） |
| 3 | `mail.PublicView.UpdatedAt *time.Time` | `GET /api/mail/config` | 新增处理器侧 `mailConfigResponse`（**内嵌** `mail.PublicView` 以免字段漂移，浅层 `UpdatedAt *string` 覆盖同名字段）+ `mailConfigWire` |

**为什么第 3 处不直接改模型类型**：`mail.PublicView.UpdatedAt *time.Time` + JSON null 已由用户 2026-09-20 P-004 书面裁决为 `user-overruled`（`GOAL-004/03-audit/A-003`，F-I-002：NULL = 未配置，不再伪造瞬时）。本次**不动模型**（Go 类型与 NULL 语义原样保留），只在处理器边界把**非空值**改走共享 formatter：两条用户裁决（`*time.Time`+null、固定 6 位输出）同时成立。

### 完整性论证（为什么这份清单是闭集）

两个互补的扫描，均在非测试 Go 源码上执行：

1. `\*?time\.Time\s+`json:` 全仓扫描 → 8 个字段：`internal/mail/outbox.go:65`、`internal/mail/runtime.go:81`、`internal/store/recovery.go:40`、`modules/wallet/subject/subject.go:29`、`modules/wallet/voucher/voucher.go:52,54,55,56`。
   - 命中 1/2 已按上表修复；
   - `internal/store/recovery.go:40` 是**内部 marker 文件**（sidecar JSON / PG comment），非 HTTP wire，且属 R2 冻结产物，**不改**；
   - `modules/wallet/voucher` 四个字段已由 `voucherJSON` 经 `FormatWireTime` 投影（R3-A 第一段）；
   - `modules/wallet/subject` 不在任何 HTTP 投影路径（`Subject` 仅模块内使用，handler 无引用），非公共 wire。
2. 时间语义键的 map 字面量扫描（`created_at|updated_at|createdAt|updatedAt|timestamp|expiresAt|expires_at|finishedAt|lastMessageAt|occurredAt|verifiedAt|deletedAt|restoredAt`）→ 唯一非 `FormatWireTime` 命中即 `internal/handler/mail_outbox.go:114`（已修复）。

## 3. 回归测试

`apps/api/internal/handler/w040_r3a_public_wire_fields_test.go`：

- `TestDefaultTimeMarshallerBreaksTheWireContract`：让「默认编解码器确实会把这两个瞬时印短」成为**可执行前提**（断言 `"2026-09-20T12:57:15.9Z"` / `"2026-09-20T12:57:15Z"`），并断言 `FormatWireTime` 给出规范形——防止上面的断言空转。
- `TestHealthProbeTimestampIsCanonicalWire`：`/healthz` + `/readyz`，`timestamp` 必须是字符串且匹配 `...\.\d{6}Z`。
- `TestMailOutboxCreatedAtIsCanonicalWire`：直接写入 `.900000Z` 行，列表与详情都必须逐字节等于 `2026-09-20T12:57:15.900000Z`。
- `TestMailConfigUpdatedAtIsCanonicalWire`：`.900000Z` 与整秒两个瞬时精确断言；NULL → JSON `null`；并做**字段漂移守卫**（`channel`/`mockRetention`/`resend`/`smtp`/`secrets` 仍在内嵌投影里）。

## 证据

- `go build ./...` exit 0。
- `go test -count=1 ./internal/handler/ ./internal/mail/`：ok（99.279s / 1.638s）。
- 全仓 `go test -count=1 ./...`：64/64 包 ok。
- `apps/web` `vitest run`：124 文件 / 1504 用例全绿。
- 上述新增 4 个名称的定向复跑：`go test -count=1 ./internal/handler/ -run "TestDefaultTimeMarshallerBreaksTheWireContract|TestHealthProbeTimestampIsCanonicalWire|TestMailOutboxCreatedAtIsCanonicalWire|TestMailConfigUpdatedAtIsCanonicalWire"` → ok。

## 未完成 / 诚实边界

- 检查点 A 的其余要件（formatter、输入兼容矩阵、Go fixture）已在 E-002 完成；**`I-041-008` 判定后 A 才成立**，判定见 Root `D-019`。
- 三处修复都是「让已冻结的 `D-003` 合同在漏网字段上成立」，**未**改变任何字段语义、JSON 键名或 NULL 语义；`mail.PublicView` 的 Go 类型保持 `*time.Time`。
