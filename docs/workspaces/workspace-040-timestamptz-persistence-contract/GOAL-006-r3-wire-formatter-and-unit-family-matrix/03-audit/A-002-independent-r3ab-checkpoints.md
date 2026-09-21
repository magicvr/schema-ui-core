---
id: A-002-independent-r3ab-checkpoints
doc: audit-entry
status: active
parent: GOAL-006-r3-wire-formatter-and-unit-family-matrix
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# A-002 · independent 审计：R3-A/R3-B 检查点（formatter、fixture、单位族矩阵、VP-020 round-trip）

- **source**: `independent`
- **auditor**: grok build · grok-4.6 · reasoning high
- **日期**: 2026-09-20
- **scope**: `GOAL-006` 检查点 A（R3-A）与 B（R3-B）；边界依据 Root `D-003`/`D-005`/`D-008`/`D-009`/`D-018`/`D-019`、GOAL-003 `D-001` §3、`GOAL-002/attachments/r1-public-wire-inventory-v0.1.md`；不含 R3-C/D。基线 HEAD `acc51aab`（工作树 clean）。
- **类型**: execution-facts + finding-closure（对照检查点 A/B 与 `I-041-008`/`I-040-004`）
- **verdict**: `fail`

> 本意见不修改 `status`/`progress`/方案正文/goal-tree。响应由 `/govern` 处理。本会话按用户指令**不落盘**；编排器应原样写入 `03-audit/A-002-independent-r3ab-checkpoints.md` 并更新 `03-audit.md` 索引。

## 范围与区间

只审 GOAL-006 检查点 A/B。不审 R3-C/R3-D，不重开 R2 checksum/descriptor。共享资料目录为 `none`。未读取其他工作区作为状态源；仅在 `I-041-008` 反证时核对仓内 `docs/` 是否仍承诺 3 位小数输出。

## 成果（独立核实）

| # | 判据 | 结论 | 证据 |
|--:|------|------|------|
| 1 | 单一 shared fixed-6 formatter 落码 | **达成** | `apps/api/internal/temporal/temporal.go:120-128` `FormatWire`（经 `NewValue` → `Truncate`，D-008 向零）；`apps/api/internal/handler/rfc3339.go:25-32` `FormatWireTime` 委托之；`TestFormatWireTime` 含截断不舍入与非 UTC 归一化（`rfc3339_test.go:11-40`） |
| 2 | inventory 列出的 Go **inline / `time.RFC3339` 输出**已替换 | **达成（inventory 面）** | 全仓 `2006-01-02T15:04:05.000Z` 仅剩 `modules/recyclebin/service.go:277`（payload 解析，D-009）；`.Format(time.RFC3339)` 仅剩邮件正文（`invites.go:182-184`、`email_identity.go:110`、`recovery.go:53`）与 audit detail（`service_credentials.go:218`）——均在 D-009 例外内。`filelibrary.go:86,121` 与 `cmd/schema-ui/configpkg.go:308,325` 走 `FormatWireTime`/`temporal.FormatWire` |
| 3 | 输入 0/3/6/9 位与 `+00:00`、拒绝无时区 | **部分达成** | `TestParseWireTimeCompatibility`（`rfc3339_test.go:58-90`）接受 0/3/6/9 与 `+00:00`，拒绝无时区/非法/2 位小数。**但同一测试把 `+08:00` 列为接受项**（L68），与 D-005 第二句及 GOAL-003 `D-001` §3 冲突，见 F-I-003 |
| 4 | Go/Web fixture 同步（有界） | **达成** | Web：`datetime.test.ts:36-44` 固定 6 位与 legacy 3 位同瞬时；`return-intent.test.ts:154-185` 接受 6/3/0 位、拒绝无时区。Go：`rfc3339_test.go` 输出 27 字符。未改写 pinned upstream JSON，与 D-019 §1.3 一致 |
| 5 | 单位族矩阵四端点、非同族充数 | **达成** | `w040_r3b_unit_family_matrix_test.go`：秒=`digital_offers`/`GET /api/digitaloffer/offers`（L48-98）；毫秒=`jobs`/`GET /api/jobs`（L100-139）；可空=`service_credentials.revoked_at/last_used_at`（分母 L39-40，unit `sec`，SQL NULL→JSON `null`，L173-197）；sentinel=`mail_config.updated_at`（分母 L44，unit `ms`，D-008 D0）。四表四端点。`requireDenominatorColumn` 从 `temporalcontract.Columns()` 读回 unit（L22-33） |
| 6 | VP-020 展示 round-trip（Go + Web） | **达成（声明的粒度诚实）** | Go：`w040_r3b_timezone_roundtrip_test.go` 4 区 × 6 瞬时，wire 与 zone 无关、parse 恢复 `temporal.Truncate`、`time.Local` 不渗入（L83-102）。Web：`utc-roundtrip.test.tsx` 用生产 `formatDate` 锁 L1 层（L146-171）；Kathmandu 墙钟硬编码 18:42（L201-203）；DST 期望值为硬编码 UTC 秒串（L207-211）。展示层秒粒度、微秒由 Go 侧证明 |
| 7 | `mailConfigResponse` GET 投影：单键 + nil→`null` | **达成（仅 GET）** | 已有 `TestMailConfigWireProjectsExactlyOneUpdatedAt`（`w040_r3a_public_wire_fields_test.go:164-193`）。独立在仓库外一次性程序复现：内嵌+浅层 shadow → 恰好 1 个 `updated_at`；nil 指针 → JSON `null`；非空视图且 `UpdatedAt==nil` → `"updated_at":null`（无 omitempty 丢键） |
| 8 | 三处修复的 JSON **键名**与 `mail.PublicView` NULL 语义 | **GET/health/outbox 无键名漂移** | health：`json:"timestamp"` 仍为字符串（`health.go:20-24`）。outbox：map 键与 `OutboxRecord` 标签一致，`keepEmptyBody` 区分列表恒带 body / 详情 `omitempty`（`mail_outbox.go:111-124`）。GET config：模型仍为 `*time.Time`（`runtime.go:81`），NULL→JSON `null`（测试 L252-267），符合 GOAL-004 `A-003` 对 F-I-002 的 `user-overruled` |
| 9 | 无自定义 `MarshalJSON`；CSV/bootstrap/configpkg 非第三类漏洞 | **达成（这些类别）** | 全仓 `func (.*) MarshalJSON` = 0。CSV `exportRow` 取已格式化 map 的 `createdAt`/`updatedAt` 字符串（`export.go:118-135`）。bootstrap 无时间字段（`bootstrap.go:21-38`）。`configpkg` 已走 `FormatWire` |
| 10 | `I-041-008` 证据表「仓内无 3 位**输出**断言」 | **同意判定** | 见下「I-041-008」。定向 vitest 30/30 绿；未发现把 3 位小数当作 API **输出**契约的现行测试 |
| 11 | 本轮建议复跑 | **绿** | `cd apps/api && go test -count=1 ./internal/handler/ ./internal/temporal/ ./internal/mail/` → ok（99.452s / 0.423s / 1.698s）。`cd apps/web && vitest run` 三个指定文件 → 3 files / 30 tests ok。未复跑全仓 64 包 / 124 文件（不采信 self 该叙述） |
| 12 | v1–v87 未触碰 | **未在本轮 diff 范围核到 migration 文件** | 工作树 clean；本审只读。inventory 替换与三处 handler 修复不涉及 descriptor/SQL |

## 对照成功标准

| 检查点 | 编排器状态 | 本审 |
|--------|------------|------|
| **A** | completed | **未达成**。formatter 与 inventory 替换成立，但 D-003 公共 wire 仍有漏网（F-I-001/002），且 D-005/GOAL-003 D-001 的「API 拒绝非零 offset」未落地（F-I-003）。`progress: 2/3` 把 A 计为完成，过早 |
| **B** | completed | **达成**。四族端点可核对；`I-040-004` 证据足以按 `I-041-007` 载体关闭。不因 A 的漏网自动否定 B |
| **C** | pending | 仍 pending；本意见为 C 的 independent 输入，**不能**在 required 未闭合时关门 |

`progress: 2/3` 确由 `00-meta.md` A～C 等权派生（不是手填百分比），但派生输入「A=completed」本审不接受。

## Findings

### F-I-001 · required · high

**问题**：`GET /api/mail/config` 已投影为 `mailConfigWire`，**同资源的 `PUT /api/mail/config` 仍 `writeJSON(view)`**，把 `mail.PublicView.UpdatedAt *time.Time` 交给 `encoding/json` 默认编解码。这是编排器已识别的字段，只修了 GET。GET 与 PUT 现在会对同一瞬时发出不同宽度。

**证据**：

- `mail_admin.go:105` GET → `mailConfigWire(view)`；`mail_admin.go:157` PUT → `writeJSON(w, http.StatusOK, view)`。
- `mail.PublicView.UpdatedAt *time.Time \`json:"updated_at"\``（`runtime.go:81`）。
- 仓库外一次性 `json.Marshal(&PublicView{UpdatedAt: &instant})` 对 `2026-09-20T12:57:15.900000Z` 得到 `"updated_at":"2026-09-20T12:57:15.9Z"`（与 `TestDefaultTimeMarshallerBreaksTheWireContract` 前提一致）。
- R3-A 回归只覆盖 GET（`TestMailConfigUpdatedAtIsCanonicalWire`），无 PUT 断言。

**关闭要求**：PUT 成功路径改为 `writeJSON(..., mailConfigWire(view))`；增加 PUT 回归：trailing-zero 瞬时逐字节固定 6 位、NULL 仍为 JSON `null`、非时间键不丢。不得改 `PublicView` 的 `*time.Time` 类型（GOAL-004 A-003 `user-overruled`）。

### F-I-002 · required · high

**问题**：第三类「时间到达公共 wire」路径存在且可构造：`time.Time` 作为 `map[string]any` 的值、键名不在 E-003 的两类扫描词表里。`GET /api/mfa/status` 的 `enrolledAt` 即是。未登录/未注册时零值 `time.Time` 还会编成 `"0001-01-01T00:00:00Z"`，而 Web 类型是 `string | null`。

**证据**：

- `handler/mfa.go:200`：`writeJSON(..., map[string]any{"enabled": enabled, "enrolledAt": enrolledAt})`，`enrolledAt` 为 `time.Time`。
- `modules/mfa/service.go:184-192`：未找到记录时返回 `time.Time{}`。
- E-003 §2 扫描 2 的键表为 `created_at|updated_at|createdAt|updatedAt|timestamp|expiresAt|expires_at|finishedAt|lastMessageAt|occurredAt|verifiedAt|deletedAt|restoredAt`——**不含 `enrolledAt`**。无 json tag，也躲过扫描 1。
- 仓库外 marshal：非零瞬时 → `"enrolledAt":"2026-09-20T12:57:15.9Z"`；零值 → `"enrolledAt":"0001-01-01T00:00:00Z"`。
- `mfa_test.go:53-54` fake `Status` 恒返回 `time.Time{}`，status 断言只查 `"enabled":true`（L221），宽度/缺席从未锁住。
- Web：`mfa-manager.tsx:27-29` `enrolledAt: string | null`；测试 mock 用 `"2026-08-15T00:00:00Z"` / `null`（输入，不是 3 位输出契约）。

**关闭要求**：该字段走 `FormatWireTime` / `FormatWireTimePtr`；未注册不得伪造年 1 瞬时（`null` 或省略，与现有 Web 类型一致）；补 handler 测试。并再扫一遍：所有 `writeJSON`/`json.Marshal` 的 `map[string]any` 与无 tag/`[]any` 时间值，按 HTTP 方法（含 PUT/POST 响应）闭集，而不是只扫 GET 与固定键名。

自定义 `MarshalJSON`、HTML 模板、bootstrap、CSV 原始 `time.Time`、configpkg：**本审未找到反例**。`json.RawMessage` 的 jobs payload/result 属 D-009 任意 payload，不升 required。

### F-I-003 · required · high

**问题**：R3 的 API 入站解析器 `ParseWireTime` **接受非零 offset**，与已冻结决策相反。这不是「多兼容一点」的无害放宽：GOAL-003 已把「codec 接受、API 拒绝」明确划给 R3。

**证据**：

- Root `D-005`：「非法时间、**非零 offset 解析失败**」。
- GOAL-003 `D-001` §3（`01-decision/D-001-temporal-codec-api.md` L71）：`temporal.Parse` 接受并归一 UTC；**「API 入站解析器必须拒绝非零 offset —— 该拒绝属 R3 的 wire formatter 工作」**。
- inventory C2：「API input parser … **rejects non-zero offsets**」。
- `temporal.go:139-141` 仍把拒绝推给 transport 层。
- `ParseWireTime`（`rfc3339.go:51-56`）只调用 `temporal.Parse`，无 offset==0 检查。
- `rfc3339_test.go:68` 把 `"2026-08-17T20:00:00+08:00"` 列为**接受**并规范化。`jobs.go:389`、`operations.go:111`、`service_credentials.go:170` 均走该 parser。

**关闭要求**（二选一，须留痕）：

1. **preferred**：`ParseWireTime` 拒绝非零 offset（保留 `Z` / `+00:00` / `-00:00`）；测试从「接受 +08:00」改为「拒绝」；handler 入站路径复跑。
2. 用户书面 P-004 **修订** D-005 / GOAL-003 D-001 / inventory（明确「任意显式 offset 归一 UTC」），再改决策正文。口头或测试现状不算 overruling。

未决前不得声称「输入矩阵与 D-005 一致」，不得把检查点 A 标 completed。

### F-I-004 · recommended · med

**问题**：`requireDenominatorColumn` 只能断言分母 `Unit`（`sec`/`ms`）。`temporalcontract.Column`（`columns.go:104-108`）没有 nullable / D0 标志，因此「可空族 vs sentinel 族」的身份不是可执行的分母检查，只靠用例选列 + SQL NULL 断言。

**证据**：四列确实不同（`digital_offers.created_at` sec；`jobs.*` ms；`service_credentials.revoked_at` 设计可空 sec；`mail_config.updated_at` D-008 D0 ms）。**不是同族充数**。sentinel 子测试直接写 NULL，不重演 v73 `0→NULL`（与 A-001 `F-S-003` 一致）；R2 转换证据仍在 `w040contracttest` / `backup/verify.go`，不应读作本轮新证据。

**关闭要求**：在矩阵注释或分母旁路中写明 D0/nullable 的判定来源（D-008 + 列清单），或给 Column 增加可核对的可空/sentinel 标记。不阻断 B。

### F-I-005 · recommended · low

**问题**：Web round-trip 的 `zoneOffsetMs` / `instantFromWallClock` / `canonicalWire`（`utc-roundtrip.test.tsx:70-97`）是测试内自实现。`canonicalWire` 用 `toISOString()` 的 3 位再替换成 `.000000Z`，只证明秒粒度。自逆辅助函数不能单独证明生产输入路径。

**为何仍不阻断 `I-040-004`**：生产 `formatDate`（`format.ts:18-34`）被 L1/L2/L3 用例对拍；Kathmandu 18:42 与 DST 期望 UTC 串是硬编码，符号/整小时错误会失败；微秒由 Go 测试证明；Web 时间输入为 date-only（用例注释 L18-20）。展示层秒粒度的声明诚实。

**关闭要求**：可选——用生产 `formatDate` 墙钟或 `Intl` 直接断言，减少自实现表面积。非关门门禁。

## 必改项汇总

1. PUT `/api/mail/config` 走 `mailConfigWire` + 测试（F-I-001）。
2. `GET /api/mfa/status` `enrolledAt` 走共享 formatter；零值不得上线；补测试；按「方法 × map/DTO」补扫描（F-I-002）。
3. `ParseWireTime` 拒绝非零 offset，或用户书面修订 D-005/D-001（F-I-003）。

## 逐条结论表

### `03-audit.md` 待复审 11 条

| # | 事项 | 结论 | 依据 |
|--:|------|------|------|
| 1 | inventory Go 面 inline 是否覆盖 | **inventory 面覆盖；闭集不成立** | 旧 milli 布局仅 D-009 解析器残留；见成果 2。漏网见 F-I-001/002 |
| 2 | parser 与 D-005（`+00:00` 接受、非零 offset 拒绝、无时区拒绝） | **不成立** | `+00:00`/无时区符合；非零 offset 被接受（F-I-003） |
| 3 | 四族各 ≥1 endpoint、非充数 | **成立** | 四表四端点；秒/毫秒 unit 可执行；可空与 sentinel 不是同一列的两种说法（F-I-004 仅方法弱点） |
| 4 | VP-020「展示随时区、存储恒 UTC」 | **成立** | Go wire 与 `time.Local` 无关；Web 生产 formatter + 层解析 |
| 5 | 未把非 DB 文本纳入固定 6 位 | **成立** | 邮件/audit detail/recycle payload 仍 RFC3339/milli 解析；符合 D-009 |
| 6 | E-003 三处修复完整性与形状 | **GET/health/outbox 完整；mail PUT 不完整** | 键名未改；outbox `keepEmptyBody` 保留列表/详情 body 形状；PublicView `*time.Time`+null 仍在。PUT 仍默认编解码（F-I-001）。独立 marshal：GET 单 `updated_at`、nil→`null` |
| 7 | 时间到 wire 路径是否闭集 | **不成立** | 反例：PUT 写裸 `PublicView`（已扫到类型、漏了方法）；`map[string]any`+`enrolledAt`（键表外）。无 `MarshalJSON`。CSV/bootstrap/HTML/configpkg 本审无反例 |
| 8 | 族身份可否独立核对 | **unit 可以；nullable/sentinel 不能从分母结构核对** | `Column` 只有 `Unit`。选列本身正确（F-I-004） |
| 9 | 52 处 `.NNNZ` / pinned JSON 不改写 | **成立，未推翻 I-041-008** | 现行测试里 3 位字面量是输入/mock。e2e、`.github`、`compose.yaml`、scripts 无 3 位输出契约。历史工作区文档仍写 3 位，已被 VP-040/`D-003`/`D-019` 取代，不是运行时消费者 |
| 10 | 展示秒 vs 存储微秒是否足以关 `I-040-004` | **足以关闭** | 粒度声明诚实；Go 微秒恒等；Web 硬编码墙钟/DST 期望（F-I-005 不阻断） |
| 11 | 基线全绿 | **指定包绿；未复跑全仓** | handler/temporal/mail 与 3 个 vitest 文件通过。不采信 self「64/64、124/1504」 |

### A-001 `F-S-003`～`F-S-006`、`N-001`～`N-005`

| ID | 结论 | 依据 |
|----|------|------|
| F-S-001 | **同意 fixed（仅 GET）** | GET 投影 + 测试；独立 marshal 单键/`null`。PUT 未修 → 新开 F-I-001 |
| F-S-002 | **同意 fixed** | `datetime.ts:4` 示例已为 `.000000Z` |
| F-S-003 | **同意限定为 wire 输出；不是 R2 转换新证据** | 子测试写 NULL 后断言 JSON `null`。v73 `0→NULL` 仍以 R2 为准。不升 required |
| F-S-004 | **升级为 required（F-I-001/002）** | 闭集主张被反例推翻，不是「理论上可能」 |
| F-S-005 | **同意：未发现 3 位输出断言** | 见待复审 9。不推翻 `I-041-008` |
| F-S-006 | **同意有自实现风险；不阻断 `I-040-004`** | 见 F-I-005。DST/Kathmandu 硬编码 + 生产 `formatDate` 足够 |
| N-001 | **同意非 HTTP wire** | `store/recovery.go:40` sidecar/PG comment；改它会重开 R2 |
| N-002 | **同意文件名非 wire** | `migrate.go:337`、`postgres.go:243`、`recovery.go:397` |
| N-003 | **同意仓内无解析 health body 的消费者** | `compose.yaml:58` `wget -qO-` 只作探针退出码；改 `time.Time`→`string` 键名仍为 `timestamp`，JSON 类型本就是字符串 |
| N-004 | **同意 D-009 例外** | `recyclebin/service.go:277` |
| N-005 | **同意 dist-lib 非受控制品** | `.gitignore:113` `apps/web/dist-lib/`；`git ls-files` 0 行。构建产物仍含 3 位注释，不影响合同 |

## 与既有意见的异同

- **同意 self**：formatter 落码、inventory 替换、四族矩阵、VP-020 矩阵载体、D-009 例外、GET mail 投影、F-S-001/002 已修、N-001～N-005、`I-041-008` 无 3 位输出依赖。
- **不同意 self**：开放 required ≠ 0；检查点 A 未完成；E-003「闭集」名不副实；`ParseWireTime` 输入矩阵与 D-005/D-001 不一致。self 把 F-S-004 列为 recommended 交独立复核——本审给出反例并升为 required。

## 明确写出

- **开放 required 数量：3**（F-I-001、F-I-002、F-I-003）。无 `accepted-residual` / `user-overruled` 之前不得推进检查点 A 完成、不得 GOAL-006 关门（P-003）。
- **`I-041-008`（无破坏性 / 不需兼容期）**：**同意 verified**。前端正则变宽小数（`datetime.ts:12-13` `(?:\.\d+)?`；`failure.ts:93` `(?:\.\d+)?Z`）；`apps/web/src` 的 `split(".")`/`length ===` 不是时间宽度断言；52 处 `.NNNZ` 为输入；dist-lib gitignored；权威合同 `D-003`/`VR-091` 即固定 6 位。本审 required 是「6 位输出尚未盖全 + 入站 offset 政策」，不是「仍有 3 位输出消费者」。
- **`I-040-004`（已验证）**：**同意 verified**。载体符合 `I-041-007`（Go 单测 + Web 组件/单测，无浏览器 e2e）。F-I-005 不阻断。
- **意见冲突 / P-004**：
  1. **必须**：F-I-003 —— 实施「拒绝非零 offset」还是书面改写 D-005/D-001。本审建议实施拒绝（D-001 已把该工作指派给 R3）。
  2. **非冲突**：F-I-001/002 应直接修，无需在「是否算公共 wire」上再裁一次（D-003 + E-003 已把默认 `time.Time` 编解码当漏网）。
- **治理**：goal-tree 与 GOAL-006 `00-meta` 的 `active · 2/3` 自洽，但 A=completed 本审不接受。`workspace.md` 仍写 R2/R3 尚未创建、R1 active——与 goal-tree 过时，建议编排器顺手更新，**不**升 required。`shared_materials_catalog: none` 合格。

## 最小补齐动作（在此之前 verdict 保持 fail）

1. PUT mail config 与 GET 共用 `mailConfigWire` + 测试（F-I-001）。
2. MFA `enrolledAt` 纳入 formatter；零值→`null`/省略 + 测试；补 map/多方法扫描（F-I-002）。
3. `ParseWireTime` 拒绝非零 offset 并修正测试，**或**用户书面修订 D-005/D-001（F-I-003）。
4. 编排器响应本意见（A-003）；required 按三路径闭合后，才能把检查点 A 标 completed、再谈 C/关门。

不要为迎合「A/B 已完成、开放 required=0」而放行。B 与 `I-040-004`/`I-041-008` 可以保留；A 不能。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。

