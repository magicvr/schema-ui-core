---
id: GOAL-006-r3-wire-formatter-and-unit-family-matrix
doc: audit
status: active
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# 审计台账 · GOAL-006-r3-wire-formatter-and-unit-family-matrix（R3-A/B）

> 本文件是唯一正式审计台账索引：`self` 与 `independent` **共用** `A-NNN` 序列。
> 每条意见正文在 `03-audit/A-NNN-<slug>.md`；本文件登记条目头（`source`/日期/scope/`verdict`）。
> 独立审计默认只写意见，不修改 `status`/`progress`/方案正文；响应归编排器。

## 意见索引

| A-ID | source | 日期 | scope | verdict | 摘要 | 文件 |
|------|--------|------|-------|---------|------|------|
| A-001 | self | 2026-09-20 | GOAL-006 检查点 A/B（formatter、fixture、单位族矩阵、VP-020 round-trip） | conditional | 8 项成果可核对；开放 required = 0；`F-S-001`/`F-S-002` 已 fixed；`F-S-003`～`F-S-006` 提交 independent 复核；`N-001`～`N-005` 已核对为非问题 | `03-audit/A-001-self-r3ab-checkpoints.md` |
| A-002 | independent（grok build · grok-4.6 · high） | 2026-09-20 | 同上（独立复核） | **fail** | 3 条 required：`F-I-001` PUT `/api/mail/config` 仍走默认编解码；`F-I-002` `GET /api/mfa/status` `enrolledAt` 为 map 内裸 `time.Time`（含零值伪造年 1）；`F-I-003` `ParseWireTime` 接受非零 offset 违反 `D-005`。另有 recommended `F-I-004`（族身份不可核对）、`F-I-005`（Web 自实现表面积）。**同意** `I-041-008`、`I-040-004` 与检查点 B 成立；判定检查点 A 未完成 | `03-audit/A-002-independent-r3ab-checkpoints.md` |
| A-003 | self（编排器响应） | 2026-09-20 | 响应 A-002 全部 findings | **pass** | 3 条 required 全部 **fixed**（PUT 共用 `mailConfigWire`；MFA `enrolledAt` 走 formatter 且缺席为 `null`；`ParseWireTime` 拒绝非零 offset——按 `D-005`/GOAL-003 `D-001` 已冻结决策实现，非新决策）；两条 recommended 亦 **fixed**；新增两条常驻守卫（含合成违规文件的决定性验证）。开放 required = 0；请求 `A-004` 定向复审 | `03-audit/A-003-response-to-independent-a002.md` |
| A-004 | independent（grok build · grok-4.6 · high） | 2026-09-21 | 定向复审：A-002 的 3 条 required 是否合法闭合（含守卫不可绕过性、offset 拼写反例、PUT 回归决定性、`enrolledAt` 契约、治理一致性） | **pass** | 开放 required = **0**：`F-I-001`/`002`/`003` 全部 `fixed`；同意检查点 A 成立、`I-041-008`/`I-040-004` 维持 verified、可进入检查点 C 关门；无 P-004 待裁项。新增 2 条 recommended：`F-I-101`（PUT 回归的巧合率实测约 10%/次，A-003 措辞过满）、`F-I-102`（守卫非闭集证明，给出 A–F 六类绕过反例） | `03-audit/A-004-independent-r3ab-finding-closure.md` |
| A-005 | self（编排器响应 / 关门记录） | 2026-09-21 | 响应 A-004 + 检查点 C 关门 | **pass** | `F-I-101` **fixed**（`Switcher.SetClock` 冻结瞬时，PUT 回归变为确定性）；`F-I-102` **fixed**（覆盖 B/D/E/F 并有决定性验证；收回「闭集」措辞；A/C 两类记录为书面缺口，A-004 已确认无现行漏网）。开放 required = 0 → 按用户裁决**静默关闭 `GOAL-006`**（`done · 3/3`）；Root 仍 `active · 2/3`，R3-C/D 待立项 | `03-audit/A-005-response-to-a004-and-checkpoint-c-closure.md` |

## 待复审事项（编排器登记，供独立审计取证）

| # | 事项 | 证据位置 | 说明 |
|--:|------|----------|------|
| 1 | ixed-6 输出替换范围是否覆盖 inventory 全部 Go 面、是否有遗漏的 inline 布局 | `internal/handler`、`cmd` | 逐项对照 `r1-public-wire-inventory-v0.1.md` §「当前 inline fixed-3 outputs」与「Current RFC3339 outputs」 |
| 2 | 输入 parser 兼容边界（`+00:00` 接受、非零 offset 拒绝、无时区拒绝）与 `D-005` 一致 | `internal/handler` parser + 测试 | 反例优先 |
| 3 | 单位族矩阵是否真的覆盖「秒/毫秒/可空/sentinel」四族各至少一个 endpoint | 测试证据 | 避免用同一族多次充数 |
| 4 | VP-020 会话时区展示 round-trip 是否证明「展示随会话时区、存储恒 UTC」 | Web 组件/单测 + Go 侧 | 载体按 `I-041-007` |
| 5 | 是否误把非 DB 文本纳入固定 6 位输出（`D-009` 例外） | `D-009` + fixture 对照 | 例外范围 |
| 6 | **E-003 的三处漏网字段修复是否完整、是否引入语义/形状变化** | `handler/health.go`、`handler/mail_outbox.go`、`handler/mail_admin.go` | 重点：`mailConfigResponse` 内嵌 + 同名字段 shadowing 的实际 JSON 行为；outbox 列表/详情 body 形状是否与修复前**逐键**一致；`mail.PublicView` 的 `*time.Time` + JSON null 是否仍满足 `GOAL-004/A-003` 的 `user-overruled` 裁决 |
| 7 | **「时间到达 wire 的路径」闭集扫描是否真的完备** | `E-003` §2 完整性论证 | 是否存在第三类路径（如 `any`/`json.RawMessage` 间接承载、模块层自行 `writeJSON`、`MarshalJSON` 自定义类型）未被两类扫描覆盖 |
| 8 | 单位族「族身份」是否可独立核对（分母 unit）而非仅凭测试命名 | `w040_r3b_unit_family_matrix_test.go` `requireDenominatorColumn` | 反例优先：能否构造出「实际不属于该族却被记为已覆盖」的情形 |
| 9 | Web 52 处 `.NNNZ` 与 pinned upstream JSON 不改写的判定是否站得住 | `E-003` §1、`D-019` §1 证据表 | 关键反证：是否有测试把「3 位小数**输出**」当作断言而非输入 |
| 10 | VP-020 round-trip 的精度声明（展示层秒粒度 vs 存储微秒）是否诚实且足以关闭 `I-040-004` | `w040_r3b_timezone_roundtrip_test.go`、`i18n/utc-roundtrip.test.tsx` | 是否存在「展示丢精度即等于存证不足」的反驳 |
| 11 | 基线是否仍然全绿且无回归 | `go test -count=1 ./...`（64/64）、`vitest run`（124 文件/1504 用例） | 审计可自行复跑；不采信本文叙述 |