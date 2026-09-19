---
id: A-003-a001-a002-response
doc: audit-entry
parent: GOAL-005-r4-result-center-experience
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-003 · R4 审计响应（A-001 self + A-002 independent）

## A-003 · R4 C1～C3 意见响应（2026-09-19）

- **source**：orchestrator（`/govern` 编排器对 A-001 self + A-002 independent 的合并响应）
- **类型** / **scope**：`stage` · GOAL-005 C1～C3
- **两腿结论**：A-001 self `pass`（0 required + 4 recommended）；A-002 independent `pass`（0 required + 4 recommended）。**无冲突**：两腿在合同、隔离、门禁、派生字段、前端打点、边界上同向；independent 未升级任何 self 项为 required，未提出与 self 相反的必改项，故**不触发 P-004 冲突裁决**。
- **响应后状态**：**开放 required = 0**；8 条 recommended 中 **5 条 `fixed`**、**3 条 `pending-user`**（拟 `accepted-residual`，等用户书面接受）。

## 1. 闭合清单（逐条）

| # | 来源 | finding | 处置 | 证据 |
|---|------|---------|------|------|
| 1 | self F-001 **+** indep F-001 | 写门禁测试无法区分 `jobs.write` 与 `jobs.read` | **fixed** | `TestJobWriteGateRequiresJobsWriteNotJobsRead`（`2ae1da37`） |
| 2 | self F-001（连带） | self 曾判定「判别性主体不可构造」——**该判断有误** | **fixed（纠正记录）** | 独立腿指出 `CreateRoleWithGrants` 路径；本响应据实纠正 §3 |
| 3 | indep F-002 | HEAD 前端目录缺 `error.job*` 键（写面错误会回退服务端串并记 missing-translation） | **fixed** | 双目录各 +7 键（1227/1227）+ 判别例（`2ae1da37`） |
| 4 | indep F-003 | 前端夹具 `retryable` 忽略 attempt 预算，`disabledWhen` 若改用 status 判据不会变红 | **fixed** | `ROWS` 增 exhausted-failed 判别行 + 变异验证（`2ae1da37`） |
| 5 | indep F-004 | 文档索引漂移（01-decision 索引、00-meta 审计状态/计数、E-001 checkpoint、D-001 口径） | **fixed** | `58c5614f` 逐条纠正，docscheck 全绿 |
| 6 | self F-002 | 状态列显示原始状态码，未逐值本地化 | **pending-user**（拟 `accepted-residual`） | 见 §2.1 |
| 7 | self F-003 | `jobs-auto-refresh` 依赖 `reloadList()` 清空选择的语义 | **pending-user**（拟 `accepted-residual`） | 见 §2.2 |
| 8 | self F-004 | 自动刷新在无进行中作业时仍按档位请求 | **pending-user**（拟 `accepted-residual`） | 见 §2.3 |

## 2. 拟 accepted-residual 的三条（等用户书面接受）

> 三条均为 low 严重度、已记录范围与复核触发条件的**有界残余**，不阻断任何门禁（开放 required = 0）。按 P-003，`accepted-residual` 需用户书面接受，故**不静默接受**。

### 2.1 self F-002 · 状态/错误文本不做逐值本地化

- **范围**：结果中心**状态列**显示冻结状态码（`queued`/`running`/…）+ 颜色徽标；作业行的**错误详情**显示作业行内存储的原始诊断串（`errorMessage`）。
- **已本地化的部分**：列头、筛选器六态选项（`schema.jobs.status.*`）、写面错误反馈（`error.job*` 七键，本轮补齐）、详情抽屉字段名。
- **为何不修**：逐值本地化需要协议级扩展（`valueLabels`/`tagMap` 类）→ 触碰 pinned 工件（`docs/schemas/**`）或新增渲染器能力，属 VP-038 明确非目标；且与既有产品约定一致（`wallet.status`、`scheduledTasks.enabled`、`users.mfaEnabled` 同样显示原始值）。
- **风险**：中文操作员在状态列读到英文码。
- **复核触发**：用户提出逐值本地化需求，或后续 VP/协议波次引入 `valueLabels`。

### 2.2 self F-003 · 自动刷新依赖 `reloadList()` 清空选择的语义

- **范围**：`jobs-auto-refresh` 每拍调用页面级 `reloadList()`，该 seam 会清空本页所有表选择（ADR-0022 D2）。
- **当前为何惰性**：`jobs` 表未声明 `props.selection`，无选择可清；交互测试断言页面无选择列（前提被钉住，一旦破坏即变红）。已登记为 `D-001` §7 R-1.2。
- **为何不修**：渲染器未提供「定向刷新表格而不清空选择」的 seam（`refreshList` 只覆盖 display 节点），新增该 seam 属渲染器能力扩展，超出本次已冻结范围。
- **复核触发**：`jobs` 表启用行选择时（届时必须改用定向刷新或显式排除）。

### 2.3 self F-004 · 无进行中作业时仍按档位请求

- **范围**：自动刷新 tick 不区分列表内是否仍有非终态作业。
- **为何不修**：组件读不到表格行（自定义组件无行数据 seam），要判断「是否还有进行中作业」需新增数据 seam 或额外请求；当前形态与既有 `monitoring-auto-refresh` 同形，且**默认 Off**、由操作员显式开启。
- **风险**：操作员开启后即使作业全部终结仍持续请求（有界的无谓流量）。
- **复核触发**：用户反馈后台流量问题，或渲染器提供行可见性 seam。

## 3. 对 self F-001 判断错误的纠正（据实留痕）

A-001 self 的 F-001 处置曾写「真正的 write-vs-read 反例在策略集分离前**不可构造**」，并以嵌套守卫替代。**该判断有误**：`env.authRepository.CreateRoleWithGrants(key, name, []string{"jobs.read"}, nil, now)` 可建只授一个权限的自定义角色（seeded roles 并非唯一授权来源），A-002 独立腿指出了这条路径。本响应据此落地真反例并取代嵌套守卫：

- **正向对照**：同一主体 `GET /api/jobs` = 200 —— 证明它确实持有 `jobs.read`；缺此对照，下面的 403 也可能由「什么权限都没有」的主体产生，断言将不可判别。
- **判别断言**：`POST .../cancel` 与 `.../retry` 均 403，且作业行未被改动。
- **变异验证**：把路由门禁改为 `"jobs.read"` 后，**本例变红**（错误信息打印出 200 + cancelled 投影），而旧 editor 用例**仍绿**——正好证明旧用例不可判别、新用例可判别。还原后 `jobs.go` diff 为空。
- **教训**：自审的「不可构造」结论属于**未穷尽验证路径的判断**，应表述为「未找到构造路径」而非「不可能」；独立腿的价值正在于此。

## 4. 独立腿采纳情况

- A-002 的 11 项「成果」独立核验与本响应无冲突，全部采纳；其独立复跑（Go 全量、`Any` 用例、点名前端测试）与自做变异（`retryable` 忽略预算→红→还原）一并留痕。
- A-002 明确未复跑 web 全量与构建；本响应以 `2ae1da37` / `58c5614f` 之后的最终全量复跑补足：`go build ./...` exit 0、`go test ./...` 全绿、`npm run typecheck` exit 0、`npm test` 117 files / 1458 tests 全绿、`unknown custom` 警告 0。
- A-002 的 F-001 建议「维持嵌套守卫即可 / 若要一次钉死键名则用自定义授予」：本响应选择**后者**（更强的真反例），超出其最低建议。

## 5. 门禁判定

- 相关意见（scope 覆盖 R4 C1～C3）：A-001、A-002、本条 —— 全部已汇总并响应。
- **开放 required = 0**（两腿均 0 required；5 条 recommended 已 `fixed`）。
- 3 条 `pending-user` 均为 low 且不阻断（无门禁语义）；按 P-003 需用户书面接受方能记为 `accepted-residual`，故本条目**不预先关闭**它们。
- 因此 **C4 检查点的判据已满足**（self + independent 落盘、开放 required = 0）；Root R4 投影与 `GOAL-005` 关门在用户就 §2 三条作出选择后执行，避免把未接受的残余写进关门记录。

## 6. 本条不修改

- 不修改 `status`（本条目自身为审计响应，不改目标 status/progress）。
- 不修改 A-001 / A-002 正文（独立意见保持原样，响应只落本条）。
- 不修改实现与方案正文（实现对应用户可见变更由 §1 的 commit 承载）。
