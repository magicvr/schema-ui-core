---
id: A-002-w9-a001-reasonableness
doc: audit-entry
goal: GOAL-009-w9-api-web-security-audit
title: W9 A-001 独立意见复审（合理性）
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-08-21
scope: A-001 审计意见是否合理（finding 技术主张、严重度/required 分级、清单一致性）
verdict: conditional
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# A-002 · W9 A-001 独立意见复审（2026-08-21）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`）
- **类型 / scope**：ad-hoc · 复审 A-001（及全文附件）对 `apps/api` + `apps/web` 的 bug/安全主张是否合理，能否作为 S2 整单冻结输入
- **verdict**：conditional
- **完整意见**：[attachments/audit-A-002-a001-reasonableness.md](../attachments/audit-A-002-a001-reasonableness.md)

## 范围与区间

- **覆盖**：A-001 全部 required 主张（F-001～F-012 及条目中实际列出的 11 条）、全文附件 P1/P2（含未编号的 P2-7 cron DOM/DOW）、P3 抽样、A-001 台账计数/编号一致性。
- **方法**：工作区绑定核对 → GOAL-009 五件套通读 → 对 A-001 每条 required 回读现行源码（不依赖 A-001 的「已亲验」声明）→ P3 抽样。未做动态 exploit、未跑 `go test`/`npm test`、未起 compose。
- **不覆盖**：不改 `status` / `progress` / 检查点 / 方案正文 / goal-tree；不闭合 A-001 代码 findings；不替用户做 I-002 范围/go 宣称裁决；不是 S4 关门审。
- **工作区**：`workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 已校验；`shared_materials_catalog: none`；`primary_plan` = `VP-009-production-hardening`）。未读取其他工作区作为证据。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| A-001 已以 `source: independent` 落盘，S1 产物存在 | [A-001](A-001-w9-independent.md)、[全文](../attachments/audit-A-001-w9-full-report.md)、[E-001](../02-execution/E-001-w9-audit-performed.md) |
| 两条 high 的代码事实成立 | 钱包 `isUniqueViolation` 仅 SQLite 文案；生产 nginx 只代理 manifest，host-bootstrap 落入 SPA fallback |
| 多数 med 代码事实成立 | 见附件逐条表；MFA 跨事务 check-then-act、恢复码整表回写、凭据名方言匹配、scheduledtasks OR/AND 均复核到源码 |
| A-001 作为「12 条 required 清单」不自洽 | 索引写「22 required」；正文称 12 条但 required 节只有 11 个标题且 **F-003 从未定义**；全文 P2-7 无 F-ID |

## 对照成功标准

| 标准 | 本 scope | 状态 |
|------|----------|------|
| S1 独立意见落盘 | 只核 A-001 作为产物是否可消费 | 部分：落盘形式合格，清单不能无条件冻结 |
| S2 用户确认 required 范围 | 本条不裁决 | 未开始；本意见要求先调和清单 |
| S3 / S4 | 不在本 scope | 未开始 |

## P-005 / 工作区核对

| 核对项 | 结论 |
|--------|------|
| I-001（finding 清单，方案前 required） | 元数据标 `verified`，但所指清单不自洽（本条 F-001）。**不把 I-001 当作已可冻结的方案输入**；是否改回 `collecting` 由 `/govern` 处理 |
| I-002（修复范围与 go 宣称，实施前 required） | 仍 open；最晚阶段未到。本条不替用户裁 |
| I-003（provider 偏差，关门前 non-blocking） | 仍 open。本条是工作区默认 grok provider 对 ox-alpha A-001 的交叉复审，**不是** S4 关门 |
| 到期 required 信息项 | 无（相对本复审 scope） |
| 共享资料 | 无 |

## Findings

### F-001 · A-001 required 清单自相矛盾，不能作为 S2 整单冻结输入

- 严重度：med ｜ 建议：required ｜ 状态：open
- A-001 自称 **F-001～F-012 共 12 条 required 全部 open**。独立回读后：
  1. `03-audit.md` 索引写「开放 required **22**」——与 2 high + 10 med = 12 不符，属抄写错误。
  2. A-001 required 节实际标题为 F-001、F-002、**F-012**、F-004、F-005、F-006、F-008、F-009、F-010、F-011、F-007（**11 条**）。**F-003 从未定义**，却出现在结论摘要「F-001/F-003/F-004/F-011/F-020 族」和 I-001 的「F-003～F-012 med」。
  3. 必改项写「F-003～F-012 对应 … F-012/F-004/…/F-011」——9 个编号，凑不上 10 条 med。
  4. 全文附件 P2 第 7 项（cron DOM/DOW 用 AND、非 POSIX「两者受限取 OR」，`scheduledtasks/store/cron.go:99-107`）**没有对应 F-ID**，也未进入 A-001 required 节。
- 后果：I-001 不能支撑「方案前 finding 清单已冻结」。`/govern` 若按「整单采纳 F-001～F-012」推进 S2，会采纳一条不存在的 F-003，并漏掉全文已确认的 P2-7。

### F-002 · A-001 F-002 技术成立，但「compose 生产栈必然启动失败」过述

- 严重度：low ｜ 建议：recommended ｜ 状态：open
- 代码事实成立：`apps/web/nginx.conf` 仅精确代理 `app-manifest.json`；`host-bootstrap.json` 落入 `try_files … /index.html` → 200 `text/html`；`bootstrap.ts:172-188` 仅 404/410 视为未提供；非 JSON → protocol 失败 → `boot.ts:169-180` `BOOTSTRAP_DOCUMENT_FAILED`，映射 `HOST_PROTOCOL_REJECTED`、`retry: none`。Dockerfile:27 打入该 conf；compose `web` 使用该镜像；API 在 `composition.go:532` 注册 bootstrap。Vite/preview 代理了该路径，故缺陷限于生产 nginx。
- 过述：compose **容器仍会 healthy**（web healthcheck 只 wget `/`）。失败的是**浏览器 boot 终态**，不是栈进程起不来。high / required 仍建议保留（生产 compose 路径应用不可用）。

### F-003 · 若干 A-001 required 应在 S2 按「安全必改 / 可靠性 / 潜伏 / 协议默认」拆开，不宜整单等同

- 严重度：med ｜ 建议：recommended ｜ 状态：open
- 独立源码复核同意「代码现象存在」，但 **required 作为安全必改** 过宽或需用户明示：
  - **A-001 F-007**（job/task 无 recover）：路径应为 `apps/api/internal/jobs/runner.go:278-281`（A-001 写成 `jobs/runner.go`）。handler goroutine 无 recover，Go 未捕获 panic 会击穿进程。这是**可用性**，不是鉴权/注入。
  - **A-001 F-008**：缺 `action.key` 且声明了 `permissionIntent` 时，targetId 为空、查找用 `actionRef` 兜底 → 门禁被跳过。属实。但 `render.tsx:753-754` 写明「Absent target = engine default is allow」——**未声明 intent 的默认放行是协议行为**，不能与「声明了 intent 却漏 key」混成「协议门禁被笔误整体击穿」。`schema-table.tsx` 在 `apps/web/src/renderer/`（非 `components/`）。
  - **A-001 F-009**：cascade 缺 source 时跳过并 `return true`。`validatePermissions` 能报 `PERMISSION_CASCADE_SOURCE_MISSING`，生产渲染路径确实未调用。若生产者 schema 不经 L2 即下发，运行时 fail-open 成立；若 L2 是发布门禁，则本条是「门禁未接线」而非已爆发越权。
  - **A-001 F-010**：`delete()` Get 出错跳过归属预检。对比 `update()` 对非 not-found 的 Get 错误 fail-closed。属实。全仓 `Scoper:` 仅测试夹具，A-001 已承认潜伏。标 required 偏严。
  - **A-001 F-012**（scheduledtasks OR/AND）：过滤器绕过属实，且与 datadictionary 已加括号对照成立。路由为 admin `tasks.read`，是**列表过滤正确性**，不是未授权读。
  - **全文 P2-7 cron AND**：调度语义 bug 属实；A-001 未纳入 F 清单。不宜在用户未看见的情况下被 S2 漏掉或被「12 条」暗示已覆盖。

### F-004 · A-001 F-004 的「助长暴力破解」过述（锁定竞态本身成立）

- 严重度：low ｜ 建议：recommended ｜ 状态：open
- `RecordLoginFailure`（`authsession/accounts.go:173-196`）确为 SELECT 计数再 UPDATE 写回，无原子 `+1` / `FOR UPDATE`。PG 并发可丢计数，账户锁定阈值推迟。SQLite 单写手串行可掩盖。
- 过述：登录路径另有进程内滑动窗口限流 `newLoginRateLimiter(15*time.Minute, 20, …)`（`auth.go:60,102-111`），以及可选 captcha。主暴力破解防线不是账户锁定计数。本条作为 **PG 锁定语义正确性** 仍合理，不宜写成已突破限流的 P2 攻击面。

## 对 A-001 各条的独立判定（摘要）

完整表见附件。一句话：

| A-001 | 独立判定 | 建议 S2 如何用 |
|-------|----------|----------------|
| F-001 high 钱包 unique 仅 SQLite | **成立**（无资金错账；幂等/去重契约在 PG 失效；同事务重读在 PG 亦不可行） | 保留 required high |
| F-002 high nginx host-bootstrap | **成立**（措辞见本条 F-002） | 保留 required high |
| F-003 | **空洞**（未定义） | 不得整单点名采纳 |
| F-004 登录失败 RMW | 代码成立；影响过述 | 保留或降为 recommended：用户裁 |
| F-005 TOTP 跨事务重放 | **成立**（GetState→Validate→SetLastUsedStep；DeleteProof 亦不核 RowsAffected；proof 按 user 不唯一） | 保留 required med |
| F-006 恢复码丢失更新/双花 | **成立** | 保留 required med |
| F-007 panic recover | 代码成立；属可用性 | 用户裁 required vs recommended |
| F-008 UI 门禁缺 key | 部分成立（intent+缺 key 失配是 bug；未标记默认 allow 是协议） | 修失配；勿扩大成「门禁整体无效」 |
| F-009 cascade 缺源 fail-open | 运行时行为成立；L2 仅测试 | 用户裁 |
| F-010 delete 归属预检 fail-open | 潜伏成立；无生产 Scoper | 建议 recommended，对齐 update 仍值得做 |
| F-011 凭据重名方言 | **成立**（匹配 `service_credentials.name`；PG 约束名 `service_credentials_name_key`） | 保留 required med |
| F-012 scheduledtasks WHERE | **成立**（编号错位，内容是全文 P2-1） | 保留为过滤器 bug；可与 P2-7 一并裁 |
| 全文 P2-7 cron DOM/DOW | **成立但无 F-ID** | S2 必须显式纳入或书面排除 |
| F-013～F-023 recommended | 抽样成立（头像 TOCTOU、otpauth 转义、`Contains("unique")`×4、`//host` 过 startsWith("/")、负 maxSize、recycle 两事务）。F-022 把无 recover 算到生产 `WithTx` **过严**：生产走 `Store.Run`，有 recover+rollback 后再 panic | 保持 recommended |
| F-024 info `.env` | **成立**：`git check-ignore` → `apps/api/.gitignore:9`；从未 tracked。本意见不摘录凭据 | 本地卫生；非本波 required |

未发现被 A-001 漏报的 P0。未发现 A-001 两条 high 为误报。A-001 `fail` 对**代码基线**仍成立。

## 必改项汇总（本条 required）

- **F-001**：S2 前调和 A-001 finding 清单（消除 F-003 空洞、索引「22」、全文 P2-7 无编号），**禁止**把「F-001～F-012 整单」当作已冻结范围。

本条 recommended：F-002（措辞）、F-003（分级拆开）、F-004（暴力破解过述）。

A-001 的代码 required **全部仍 open**（本条不闭合、不降级裁定——降级属 P-004，须用户书面）。

## 与既有意见的异同

- 与 A-001：同意 fail 的方向与两条 high 的核心事实；不同意「12 条 required 清单已可整单消费」；不同意把协议默认 allow、潜伏 Scoper、进程 panic 与 PG 方言/MFA 竞态放在同一 required 权重上不经用户拆分。
- 无 self 意见。本条是默认 grok provider 对 ox-alpha A-001 的交叉复审（I-003 语境），不是 S4。

## 结论 + 建议给编排器/用户的下一步

- **verdict conditional**：A-001 作为安全/缺陷审计，技术主干合理；作为 S2 冻结清单，**不合格**。
- 建议 `/govern`：① 响应本条 F-001：产出一份调和后的 finding 表（含 P2-7 去留、F-003 作废或改挂、索引 22→实际条数）；② 再问用户 I-002（哪些 required 修复、是否暂挂 VP-008 go）；③ 不要在清单未调和时整单开工。
- 建议保留为修复主干的最少集（本审计员建议，**非裁决**）：A-001 F-001、F-002、F-005、F-006、F-011，外加 scheduledtasks WHERE（现 F-012）与显式处理 P2-7。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
