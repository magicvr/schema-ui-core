---
id: E-005-a002-required-fixes
doc: execution-entry
status: active
parent: GOAL-006-r3-wire-formatter-and-unit-family-matrix
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-005 · 响应 independent A-002：3 条 required 修复 + 两条常驻守卫

## 事实

- **A-002（grok build · grok-4.6 · high）判 `fail`**，3 条 required：
  1. `F-I-001`：`PUT /api/mail/config` 仍 `writeJSON(view)`（编排器只修了 GET，同一资源两半宽度不一致）；
  2. `F-I-002`：`GET /api/mfa/status` 的 `enrolledAt` 是响应 map 里的裸 `time.Time`（键不在手写扫描表内，且未注册时输出伪造的 `0001-01-01T00:00:00Z`，而 Web 类型是 `string | null`）；
  3. `F-I-003`：`ParseWireTime` 接受非零 offset，与 `D-005`「非零 offset 解析失败」及 `internal/temporal` 包文档「该拒绝属 R3 的 API 入站解析器」相反；此前的测试甚至把 `+08:00` 写成接受项。
  - 同时判 **B 成立**、**同意** `I-041-008`（无破坏性/不需兼容期）与 `I-040-004`（verified），并明确「A 不能算完成、`progress: 2/3` 的 A=completed 过早」。
- **修复**（详见 `03-audit/A-003-response-to-independent-a002.md`）：
  - `handler/mail_admin.go`：PUT 成功路径改用 `mailConfigWire(view)`。
  - `handler/mfa.go`：`enrolledAt` 缺席 → `null`，在场 → `FormatWireTime`。
  - `handler/rfc3339.go`：`ParseWireTime` 在原始串上校验 offset 必须为 0（`Z` / `+00:00` / `-00:00` 仍接受）；codec 层 `internal/temporal` 不变。
  - `handler/mfa_test.go`：`fakeMFAService` 增加 `statusEnrolledAt` seam，使「在场」方向也可断言。
- **新增常驻守卫**：`internal/w040contracttest/wire_bypass_guard_test.go`
  - 守卫 1 `TestNoJSONTaggedWireTimeFieldOutsideProjections`：非测试 Go 源中任何 `json:` tag 的 `time.Time`/`sql.NullTime` 字段都必须在**带理由的 allowlist** 内，且 allowlist 不得过期。
  - 守卫 2 `TestNoUnformattedResponseTimeValue`：非测试 Go 源中 `*_at` / `*At` / `timestamp` 键的值必须含 `Format` 或为字面量；键模式比上一轮手写表宽（含 `enrolledAt`），与 HTTP 方法无关。
  - 守卫 3 `TestUnitFamilyProvenanceIsMachineChecked`（响应 recommended `F-I-004`）：族身份由「冻结分母 unit + 生成描述符 `NonNull: false` + DDL 无 `NOT NULL`」三处旁证核对。
  - **决定性验证**：临时放入合成违规文件后两条守卫各自失败并点名该文件，随后删除；不是空洞断言。
- **Web**：`utc-roundtrip.test.tsx` 增加 4 时区「测试内墙钟 helper ↔ 生产 `formatDate`」对拍表（响应 `F-I-005`）。

## 进度评估

- 修复后 **开放 required = 0**；检查点 A 重新成立（此前被 A-002 判未完成），B 未受影响；`progress` 保持 `2/3`，检查点 C 待 `A-004` 定向复审与关门。
- **未关闭 `GOAL-006`**：按 `GOAL-005` 先例，`fail` 意见的 3 条 required 修复后需 independent **定向复审**（`A-004`）确认合法闭合。
