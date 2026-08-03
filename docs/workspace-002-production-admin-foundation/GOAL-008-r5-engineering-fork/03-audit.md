---
title: 审计台账 · R5 · 工程化、fork 体验与集成关门
status: active
created: 2026-08-02
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.15.0
---

# 审计台账 · GOAL-008

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | independent | 2026-08-02 | goal-definition + design-plan · GOAL-008 立项信息与 Root I-005/D-013 合理性 | conditional | responded：F-001 **fixed**（D-002 + D-001 修订 + S2 对齐）；R-001/R-002 handled |
| A-002 | independent | 2026-08-02 | finding-closure · A-001 F-001 与 R-001/R-002 响应证据 | pass | responded：F-001 `fixed` 关闭成立；R-001/R-002 handled；R-003 handled（投影清理） |
| A-003 | self | 2026-08-02 | goal-definition + design-plan 复核 · GOAL-008 立项与 I-005/D-013 方案边界 + A-001/A-002 响应证据 | pass | —；补同 scope `source: self` 覆盖（P-004 §3.1） |
| A-004 | independent | 2026-08-02 | execution-facts · S1/S2 实施向审计（C-001～C-007） | pass | responded：pass 采纳；R-001 handled（由 F-002 修复承载）；R-002 fixed（README `.env` 注记） |
| A-005 | independent | 2026-08-03 | execution-facts · S1/S2 实施复核（C-001～C-007） | fail | responded：F-002 **fixed**（config.ValidateProd 生产守卫 + 回归 + 运行时反例复验）；F-003 **fixed**（00-meta 进度投影同步） |
| A-006 | self | 2026-08-03 | execution-facts 复核 · S1/S2 实施 + A-004/A-005 响应证据（含 F-002/F-003 fixed） | pass | —；补同 scope `source: self` 覆盖（P-004 §3.1） |
| A-007 | independent | 2026-08-03 | finding-closure · F-002/F-003 关闭证据复审 | pass | responded：pass 采纳；F-002/F-003 `fixed` 维持闭合；本 scope 无开放 required |
| A-008 | independent | 2026-08-03 | design-plan · I-008-002 fork 复现与 smoke 协议 v0.1.0 合理性 | conditional | responded：F-004 **fixed**（协议 v0.1.1：S4 强制 disposable SM-006 证据）；R-008-001～004 absorbed |
| A-009 | self | 2026-08-03 | design-plan 复核 · I-008-002 协议 v0.1.1 修订与 F-004 关闭证据 | pass | —；补同 scope `source: self` 覆盖（P-004 §3.1 用户裁决「补 self」） |
| A-010 | self | 2026-08-03 | execution-facts · S3/S4 实施（QUICKSTART + 独立复现 + smoke.sh + CI 接入） | pass | —；S3/S4 实施向自审（P-002 阶段审计） |
| A-011 | independent | 2026-08-03 | execution-facts · S3/S4 实施与协议/运行证据交叉复核 | fail | open：F-005～F-009 required；与 A-010 同 scope verdict 冲突，待 `/govern` 按 P-004 处理 |

## 当前审计边界

- 本目标于 2026-08-02 立项，`active / 2/5`。A-002（independent · finding-closure · pass）确认 A-001 **F-001 的 `fixed` 关闭成立**；R-001/R-002 handled；R-003 handled（投影清理）。**A-003（self · goal-definition + design-plan 复核 · pass）**：立项与 I-005/D-013 方案边界、A-001/A-002 响应闭合证据经 self 复核成立（P-004 §3.1）。**S1/S2 已实施（2026-08-02）**：env 清单 + health/启动验证 + dev/prod 区分文档 + Dockerfile × 2 + compose.yaml + nginx 反代 + CI `container-smoke`；契约 C-001～C-007 本机验证通过（`docker compose up` → healthz/登录/`/me`/SPA fallback/重启与 down-up DB 持久化）。**建议对 S1/S2 做一次实施向审计（self 或 `/audit`）**。
- 信息门禁：**`I-008-001` 已 verified（D-003 + [I-008-001-engineering-contract.md](attachments/I-008-001-engineering-contract.md) v1.0.0）**；**`I-008-002` 已 verified（D-004 + [I-008-002-fork-reproduction-protocol.md](attachments/I-008-002-fork-reproduction-protocol.md) v0.1.0）**，解除 S3/S4 的协议层信息门禁但不构成实施或验收；`I-008-003` 仍 `open`（required，仅当 S6 实施时阻断）；A-001 不把 Root `I-005: verified` 当成实现或验收证据。
- **A-004（independent · execution-facts · S1/S2 实施向审计 · pass）**：S1/S2 实施主张与仓库事实一致且可复现——C-001/C-002 静态与代码核对通过、C-003/C-005/C-006 本机 Docker 复跑通过、C-004 镜像构建成功、C-007 CI job 结构与本地 smoke 通过；无开放 required；`0/5 → 2/5` 勾选有据。
- **A-005（independent · execution-facts · S1/S2 实施复核 · fail）**：独立重跑 C-002～C-007 均通过，但 C-001/§5 “生产禁止启用 `AUTH_DEV_SESSION_ENABLED`”与实现不符：生产进程可由该 flag 启动并对未认证请求注入静态高权限开发身份（F-002）。同时 `00-meta.md` 仍写 `progress: 0/5`，而勾选成功标准和 `goal-tree.md` 均为 `2/5`（F-003）。这与 A-004 将同一安全事实列为非阻断 R-001 的分类存在相关分歧；不得跳过 P-004 裁决。
- **A-006（self · execution-facts 复核 · pass）**：补 S1/S2 实施 scope 的 `source: self` 覆盖（P-004 §3.1 用户裁决「需要补 self」）。**A-004/A-005 统一响应（2026-08-03）**：分类分歧按用户裁决采用 **F-002 `required/high`** 口径——F-002 → **fixed**（`config.ValidateProd()` 生产守卫：非 `development` + `AUTH_DEV_SESSION_ENABLED=true` → 启动报错退出；`config_test.go` 4 用例回归；运行时复验 `production + flag=true` → exit 1、`development + flag=true` → 正常启动；`go build`/`go vet`/`go test ./...` 全绿）；F-003 → **fixed**（`00-meta` frontmatter `progress: 0/5 → 2/5`，与正文/goal-tree 一致）；A-004 R-001 → handled（由 F-002 修复承载同一建议）、R-002 → fixed（根 README 补 `.env` 免重复 export 注记）。S1 的 C-001「生产禁止启用 dev session」现由运行时硬门禁成立；C-002～C-007 维持 pass。**建议对 F-002/F-003 关闭证据做一次 finding-closure 复审（`/audit`）**。
- **A-007（independent · finding-closure · pass）**：独立复核 F-002/F-003 的代码、回归、运行时与进度投影证据，确认两项 `fixed` 关闭成立；`I-008-002`/`I-008-003` 及 Root R5 门禁不在本次复审中改变。
- **A-007 响应（/govern · 2026-08-03）**：采纳 `pass`——F-002/F-003 `fixed` 维持闭合；同 scope（F-002/F-003 关闭证据）已有 A-006（self）覆盖，P-004 §3.1 无需再补自审；本 scope 无开放 required。当时 `I-008-002` 仍为 S3/S4 前置 required 门禁。
- **D-004 冻结后的当前边界（/govern · 2026-08-03）**：`I-008-002` 已由 [复现与 smoke 协议 v0.1.0](attachments/I-008-002-fork-reproduction-protocol.md) 置为 `verified`，其协议层信息门禁解除；协议没有形成新的审计 verdict，也没有改变 A-001～A-007 的历史结论。S3/S4 的文档、脚本、独立复现与运行/CI 证据仍未发生，检查点维持 `2/5`；`I-008-003` 仍在 S6 实施时适用。
- **A-008（independent · design-plan · conditional · 2026-08-03）**：独立审计 [I-008-002-fork-reproduction-protocol.md](attachments/I-008-002-fork-reproduction-protocol.md) v0.1.0 合理性。协议主体与 Root D-013 / A-001 R-002 / 仓库事实对齐，且未把协议冻结合并成 S3/S4 已完成；但 **F-004 required/medium 开放**——S4「种子可重复」未强制至少一次 disposable SM-006 通过证据，默认非 disposable 路径可在未验证种子幂等时 exit 0。R-008-001～R-008-004 为 recommended。本意见**不**改 `I-008-002` verified 状态、不改 progress、不实施 S3/S4。
- **A-009（self · design-plan 复核 · pass）**：按 P-004 §3.1 用户裁决「补 self」补齐 A-008 同 scope 的 `source: self` 覆盖——复核 [I-008-002 协议 v0.1.1](attachments/I-008-002-fork-reproduction-protocol.md) 修订（F-004 fixed + 吸收 R-008-001～004）关闭证据成立；`I-008-002` 维持 `verified`（v0.1.1 权威）。
- **A-008 响应（/govern · 2026-08-03）**：采纳 `conditional`；**F-004 → fixed**（D-005 + 协议 v0.1.1：S4 验收强制 ≥1 次 disposable/隔离运行且 `SM-006=PASS`，非 disposable exit 0 不得单独作为「种子可重复」关闭证据；不安全 destructive → exit 2）；**R-008-001～004 absorbed**（默认 URL 按路径区分、终点 4 = `list-edit-lifecycle`、SM-005 判据 `id="root"`、退出码 2/6 分离）。`I-008-002` 维持 `verified`（v0.1.1）；S3/S4 实施仍未发生，检查点维持 `2/5`；Root R5 未勾选（`4/5`）。
- 后续意见从 A-010 起。

## A-008 · I-008-002 fork 复现与 smoke 协议合理性独立审计（2026-08-03）

- **source**：independent
- **auditor**：GitHub Copilot（Grok 4.5）
- **类型 / scope**：design-plan；审计 [I-008-002-fork-reproduction-protocol.md](attachments/I-008-002-fork-reproduction-protocol.md) v0.1.0 是否足以合理回答 `I-008-002`（15 分钟计时复现协议 + smoke 验收判据），对照 Root D-013、I-005、A-001 R-002 最低清单、GOAL-008 S3/S4 成功标准、I-008-001 契约边界与当前仓库事实。不审计尚未发生的 S3/S4 实施、独立复现实测、CI run 或 Root R5 关门；不复判 S1/S2 或 F-002/F-003。
- **verdict**：conditional

### 范围与区间

- 工作区：`workspace-002-production-admin-foundation`；canonical root `docs/workspace-002-production-admin-foundation/`；Root `GOAL-001-production-admin-foundation` 与本目标 `parent` 一致；`shared_materials_catalog: none`（未把共享资料或其他工作区内容当事实）。
- 已审阅：本目标 `00-meta`（S3/S4、`I-008-002`）、`01-decision` D-004、`02-execution`「冻结 I-008-002」、本文件 A-001～A-007；协议附件 v0.1.0；Root D-013 与 I-005 §3/§4；I-008-001 §2/§8（计时/smoke 归属边界）；`compose.yaml` 端口、`apps/web` 5173/`WEB_PORT`、`app-manifest.json`（`homePageRef: overview`、page title `List + edit lifecycle`）、`seed_records.go`（8 条、`rec-1`/`Acme Console`）、records list `total`/`pageSize` cap、login/`/me` 字段、`index.html` `#root`；仓库 HEAD `5e27019482eb8d0695c402b784860233bbc90c39` 与协议 `frozen_revision` 一致；工作树含协议与 GOAL-008 文档未提交变更；`scripts/smoke.sh` **不存在**（与协议「未实施」声明一致）。
- 本轮为协议/设计合理性静态审计；**未运行**应用、Docker、计时复现或未来 smoke 脚本；不把静态核对写成 S3/S4 通过。

### 成果（有证据）

| 审计项 | 结论与证据 |
|--------|------------|
| 信息问题覆盖（A-001 R-002） | **总体满足**：协议冻结工具/平台基线（§2.1）、依赖缓存前提与排除项（§3.1）、计时起止与四终点（§3.2）、失败/重试与独立性（§3.3）、独立复现最小字段（§4）、SM-001～SM-006 + 输入/退出码（§5）。 |
| 与 Root D-013 方向对齐 | **总体对齐**：终点含登录 + 可交互列表；不含依赖下载；≥1 次独立复现；smoke 清单形状 `/healthz → login → /me → 代表页 → 种子可重复`；`scripts/smoke.sh` 位置与机器可判定退出码明确。 |
| 事实/门禁边界 | **成立**：§1/§6 明确本协议不是 S3/S4 实施或验收证据；D-004/`02-execution`/`00-meta` 同步「verified 只解除信息门禁」；未发现把协议冻结合并为检查点勾选。 |
| 双路径与 Compose 义务 | **成立**：§2.2 保留 `local-dual-process` 与 `compose`；S3 最少一条合格独立复现、未测路径不得宣称达标——不削弱 S2 Compose 必须交付。 |
| 安全与 disposable 边界 | **方向正确**：禁止输出 token/password/secret；缺工具/密码 fail-closed；普通开发库禁止 reset；CI 应用隔离 project/volume。 |
| 与仓库静态事实一致性 | **多数可核对**：默认种子 8/`rec-1`/`Acme Console` 与 `seed_records.go` 一致；page title `List + edit lifecycle` 与 manifest 一致；`accessToken`、`/healthz` `status:ok`、`/api/records?pageSize=100` 与 handler 形状一致；Compose web `8081:80` 与协议 compose 默认期望一致；freeze revision 与当前 HEAD 一致；smoke 脚本尚不存在。 |
| S3 浏览器 vs S4 HTTP 分层 | **合理**：终点 4 要求真实浏览器列表加载；SM-005 仅 HTTP 200 + SPA root，并写明浏览器可交互归 S3——避免把 curl 冒充 UI 验收。 |

### 对照成功标准 / 信息项 / 上游口径

| 项 | 审计结论 |
|----|----------|
| `I-008-002` 问题「协议与判据是什么」 | **有条件可回答**：主体判据可实施、可核对；但 S4「种子可重复」的**验收强制力**不足（见 F-004）。 |
| S3（fork 文档 + ≤15 分钟 + 独立复现） | 计时、字段、独立性、四终点总体可指导实施；路径默认 URL 与「后台首页」操作化需文档消歧（R-008-001/002）。 |
| S4（`scripts/smoke.sh` + 本地/CI 全绿） | SM-001～SM-005 + 退出码可用于脚本骨架；**SM-006 未成为 S4 绿的硬条件**，与成功标准字面「种子可重复」及 D-013 smoke 清单不完全同构（F-004）。 |
| I-008-001 边界 | **未混叠**：协议不重开容器/env 契约；health/登录终态与 I-008-001 §2 启动验证口径衔接。 |
| 状态主张 | 协议与 D-004 **未**主张 S3/S4 已完成；本意见不改变 `active / 2/5` 或 Root `4/5`。 |

### Findings

#### F-004 · S4 验收未强制至少一次 disposable SM-006（种子可重复）通过

- **级别 / 严重度**：required / medium
- **状态**：open
- **关联门禁**：`I-008-002` 作为 S4 验收判据的完整性；S4 实施与勾选；Root D-013 smoke 清单中的「种子可重复」。
- **证据**：
  1. GOAL-008 S4 与 Root D-013 均要求 smoke 覆盖「种子可重复」。
  2. 协议 §5.2 将 SM-006 标为「仅 disposable 模式」；§5.3 退出码 `0` =「SM-001～SM-005 通过，**且要求的** disposable SM-006 也通过」——未要求 disposable 时，SM-006 可整段跳过仍 exit 0。
  3. §5.1 仅写「CI **应**使用隔离 Compose project/volume」，未把「S4 验收证据必须包含一次 SM-006=PASS 的 disposable/隔离运行」写成硬判据。
- **风险**：实现者可交付默认非破坏 smoke、本地/CI 全绿并勾选 S4，却从未验证空库种子数量/幂等/不重复播种；与 D-013 与 S4 字面成功标准漂移。
- **要求（建议修复路径）**：由 `/govern` 修订协议（建议 v0.1.1）并留痕，至少明确其一并写进验收：
  1. **S4 验收必须**包含 ≥1 次 disposable/隔离运行且 `SM-006=PASS`（推荐：CI `container-smoke` 或等价 job 默认 disposable）；且/或
  2. 非 disposable 默认路径的 exit 0 **不得单独**作为 S4「种子可重复」关闭证据。
  - 保持「拒绝对普通开发库 destructive reset」的安全边界。
  - 本 finding **不**要求已存在的脚本（当前无 `scripts/smoke.sh`）；闭合证据 = 协议修订 + D-004 补记/新决策，而非运行结果。

#### R-008-001 · 默认 `WEB_BASE_URL=http://localhost:8081` 未按路径区分

- **级别 / 严重度**：recommended / medium
- **证据**：协议 §3.2 写当前默认期望 API `:8080`、Web `:8081`；`compose.yaml` web 映射 `8081:80` 成立，但 `local-dual-process` 默认 Vite 为 `5173`（`WEB_PORT` 可覆盖，见 `apps/web/README.md` / `vite.config.ts`）。协议虽要求记录实际 URL，仍给出无条件 8081 期望。
- **影响**：不推翻双路径设计；但 QUICKSTART/复现记录易误用 compose 端口跑双进程路径，造成假失败或文档歧义。
- **建议**：按 `path` 给出默认期望表（compose → 8081；local-dual-process → `http://localhost:${WEB_PORT:-5173}`），并重申不得以默认覆盖实测端口。

#### R-008-002 · 「后台首页可交互」操作化为 `/list-edit-lifecycle` 未显式对齐 manifest 首页

- **级别 / 严重度**：recommended / medium
- **证据**：Root D-013 / S3 用语为「后台首页可交互（列表加载）」；manifest `homePageRef` 为 `overview`，而协议终点 4 / SM-005 固定 `/list-edit-lifecycle`（title `List + edit lifecycle` + `Acme Console`）。该操作化**强于**仅打开 overview，且与 I-005「至少一个可交互代表页」一致，但是否等于「首页」未在协议中写明。
- **影响**：不构成终点过弱；但 S3 文档若只引导到 overview，会与协议终点 4 不一致。
- **建议**：在协议 §3.2 增加一句：D-013「后台首页可交互（列表加载）」在本目标操作化为 R4 代表页 `list-edit-lifecycle`（非 `overview`），QUICKSTART 必须覆盖该路由。

#### R-008-003 · SM-005「含 SPA root」机器判据偏虚

- **级别 / 严重度**：recommended / low
- **证据**：§5.2 SM-005 仅写「HTTP 200 且含 SPA root」；`apps/web/index.html` 实际标记为 `<div id="root"></div>`。生产 nginx 构建产物仍应保留该挂载点，但协议未钉死可 grep 的稳定子串。
- **建议**：将通过条件写为响应体包含 `id="root"`（或等价稳定标记），避免实现时主观解释「SPA root」。

#### R-008-004 · 退出码 6 混用「不安全 destructive」与「种子断言失败」

- **级别 / 严重度**：recommended / low
- **证据**：§5.3 退出码 `2` = 参数/工具/安全前提不满足；`6` = 种子重复性失败，**或**请求了不安全的 destructive 模式。不安全模式更接近安全前提（应偏 `2`），与 SM-006 数据断言失败同码会降低 CI 分类信号。
- **建议**：不安全 destructive → `2`；SM-006 断言失败单独 `6`。

### 必改项汇总

- **F-004（required / medium）· 开放**：修订 I-008-002 协议，使 S4 验收**强制**至少一次 disposable/隔离路径下 SM-006 通过；在此之前，不得把「仅 SM-001～SM-005 全绿」当作 S4「种子可重复」已满足。
- F-004 **不**否定协议其余计时/复现字段的可用性；**不**自动重开 S1/S2；**不**由本 independent 意见改写 `I-008-002` 的 verified 机读状态（是否补丁后维持 verified 或短暂回退 collecting 由 `/govern` + 用户裁决）。
- R-008-001～R-008-004 为 recommended，建议同版协议补丁一并处理，非本意见放行阻断项。

### 与既有意见的异同

- 同意 A-001 R-002：关闭 `I-008-002` 前必须把高层方向落成可执行契约；本意见确认 v0.1.0 **大部分**做到了，并指出 S4 种子项仍有强制力缺口。
- 同意 D-004 / A-007 边界：协议 verified ≠ S3/S4 实施完成；本意见不把协议审计扩成执行事实通过。
- 与 A-004/A-005/A-006/A-007 无同 scope 冲突（那些是 S1/S2 或 F-002/F-003）；F-002/F-003 保持已闭合，不在本 scope 重开。
- 无 self 同 scope 意见；按 P-004 §3.1，未来若以本意见推进协议修订或 S3/S4，应询问是否补 self。

### 结论 + 建议给编排器/用户的下一步

- **conditional**：I-008-002 v0.1.0 **总体合理**——覆盖 R-002 最低清单、对齐 D-013、安全 disposable 边界清晰、未伪造 S3/S4 完成；但因 **F-004**，不能把现行协议无条件当作完整的 S4 验收契约。
- 建议 `/govern`：
  1. 响应 A-008：展示 F-004；建议 `fixed` = 协议 v0.1.1（强制 S4 disposable SM-006 证据）+ 可选吸收 R-008-001～004；
  2. 询问是否补同 scope self；
  3. 补丁合并前，S3 文档/计时可按现协议推进草稿，但 **S4 脚本验收判据应以补丁后文本为准**；
  4. 勿勾选 S3/S4，勿改 Root R5。

### 声明

本意见仅追加独立审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文、`I-008-002` 状态或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。

## A-009 · I-008-002 协议 v0.1.1 修订与 F-004 关闭 self 复核（2026-08-03）

- **source**：self
- **auditor**：/govern（self）
- **类型 / scope**：design-plan 复核；按用户 P-004 §3.1 裁决「补 self」补齐 A-008（independent）同 scope 的 `source: self` 覆盖——复核 [I-008-002-fork-reproduction-protocol.md](attachments/I-008-002-fork-reproduction-protocol.md) **v0.1.1** 修订（F-004 fixed + 吸收 R-008-001～004）及其关闭证据；不审 S3/S4 实施、独立复现实测或 Root R5 关门。
- **verdict**：pass

### 范围与依据

- 工作区：`workspace-002-production-admin-foundation`；canonical root 一致；Root `GOAL-001-production-admin-foundation`；`shared_materials_catalog: none`。
- 已复核：协议附件 v0.1.1（§3.2 默认 URL 表、终点 4 操作化、§5.1 disposable 强制、§5.2 SM-005/SM-006、§5.3 退出码、§7 修订记录）；D-005；A-008 F-004 要求与 R-008-001～004；Root D-013 与 I-005 §3/§4；A-001 R-002 最低清单；I-008-001 §2/§8 边界。
- 本 self 为静态核对；**未运行**应用、Docker、脚本或计时复现；不把协议修订写成 S3/S4 实施或验收。

### Findings

- **F-004（required/medium）→ `fixed`**：协议 v0.1.1 §5.1 明确「S4 验收必须包含至少一次 disposable/隔离运行且 `SM-006=PASS`；非 disposable 默认路径的 exit 0 不得单独作为 S4「种子可重复」关闭证据」；§5.2 SM-006 标为「disposable 模式 · S4 必检」；§5.3 退出码 `0` =「SM-001～SM-005 通过，且 SM-006（disposable）通过——S4 完整绿」；安全边界「拒绝对普通开发库 destructive reset」保留（不安全 destructive → exit 2）。与 A-008 要求同构，`fixed` 成立。
- **R-008-001 → absorbed**：§3.2 默认期望按路径区分（compose → API `:8080`/Web `:8081`；local-dual-process → API `:8080`/Web `${WEB_PORT:-5173}`），并重申不得以默认覆盖实测端口。
- **R-008-002 → absorbed**：§3.2 增加操作化说明——D-013「后台首页可交互（列表加载）」= R4 代表页 `list-edit-lifecycle`（非 manifest `homePageRef: overview`），QUICKSTART/README 必须覆盖该路由。
- **R-008-003 → absorbed**：§5.2 SM-005 通过条件写为「响应体含 `id="root"`（或等价稳定 SPA 挂载标记）」。
- **R-008-004 → absorbed**：§5.3 退出码 `2` 并入「不安全 destructive 模式」；`6` 仅指 SM-006 种子断言失败。
- **一致性**：v0.1.1 与 Root D-013、A-001 R-002 最低清单、I-008-001 边界及仓库静态事实（`id="root"`、Vite 5173/`WEB_PORT`、compose 8081、manifest overview）一致；未把协议修订扩写为 S3/S4 完成。

### 对照成功标准与信息门禁

- `I-008-002` 问题「协议与判据是什么」：**可回答**——v0.1.1 使 S4「种子可重复」具备硬性验收强制力。
- S3/S4：计时、字段、独立性、四终点、SM-001～SM-006 与退出码总体可指导实施；实施与独立复现实测仍未发生。
- `I-008-002` 维持 `verified`（v0.1.1 权威）；`I-008-003` 仍 open（仅当 S6）。本 self 不改动任何 status/progress。

### 声明

本自审仅追加审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。

## 统一响应 · A-008（2026-08-03）

按 P-004 §3.1 用户裁决「补 self」先行 A-009 自审，再统一响应 A-008。

- **verdict 采纳**：`conditional` 成立——I-008-002 协议 v0.1.0 总体合理（覆盖 A-001 R-002 最低清单、对齐 D-013、安全 disposable 边界清晰、未伪造 S3/S4 完成），但 F-004 使 v0.1.0 不足以单独作为完整 S4 验收契约。
- **F-004（required/medium）→ fixed**：协议补丁 **v0.1.1**（D-005）——S4 验收强制 ≥1 次 disposable/隔离运行且 `SM-006=PASS`；非 disposable 默认路径 exit 0 不得单独作为 S4「种子可重复」关闭证据；保持「拒绝对普通开发库 destructive reset」安全边界（不安全 destructive → exit 2）。`I-008-002` 维持 `verified`，权威版本 v0.1.1。
- **R-008-001～004 → absorbed**：同版 v0.1.1 一并处理——默认 URL 按路径区分（compose 8081 / 双进程 `${WEB_PORT:-5173}`）；终点 4 操作化对齐 `list-edit-lifecycle`（非 overview）；SM-005 判据钉死 `id="root"`；退出码 2/6 语义分离。
- **P-004 §3.1 处置（用户裁决）**：A-008 为 `source: independent` 且协议设计 scope 无 self 审计；用户裁决「**补 self**」——已落盘 **A-009（self · design-plan 复核 · pass）** 作为同 scope self 覆盖（复核 v0.1.1 补丁与 F-004 关闭证据）。
- **仍开放**：S3/S4 实施与验收仍未发生（检查点维持 `2/5`）；`I-008-003` open（仅当 S6 实施）；Root R5 未勾选（Root `4/5`）。
- **下一步**：实施 S3（QUICKSTART/fork 文档 + ≥1 次独立复现记录，终点 4 = `list-edit-lifecycle` 列表加载）与 S4（`scripts/smoke.sh` + 本地/CI 全绿，含 ≥1 次 disposable `SM-006=PASS` 证据）；S4 完成后建议实施向审计（self 或 `/audit`）。
- **证据路径**：本响应节；`01-decision` D-005；协议附件 v0.1.1（§3.2/§5.1/§5.2/§5.3/§7）；本文件 A-008/A-009；`00-meta` 信息表 `I-008-002`；`02-execution` 2026-08-03「响应 A-008」节。

## A-001 · GOAL-008 立项信息与 I-005 工程化 / fork 报告独立审计（2026-08-02）

- **source**：independent
- **auditor**：Codex（GPT-5）
- **类型 / scope**：goal-definition + design-plan；审计 GOAL-008 立项目的、成功标准、信息门禁与父级 Root `I-005` / D-012 / D-013 的合理性，重点核对 [I-005 工程化 / fork 信息收集报告](../GOAL-001-production-admin-foundation/attachments/I-005-engineering-fork-collection.md)。不审计尚未发生的 R5 实施或关门，不复判 R1～R4。
- **verdict**：conditional

### 范围与依据

- 工作区：`workspace-002-production-admin-foundation`；canonical root `docs/workspace-002-production-admin-foundation/`；Root `GOAL-001-production-admin-foundation` 与本目标 `parent` 一致；工作区与 Root 均绑定 `VP-002-production-admin-foundation`，`vision_ref` 精确匹配现行 Charter `schema-ui-core-admin-foundation@0.1.0`。
- 共享资料：`shared_materials_catalog: none`；I-005 与本意见均只使用当前仓库代码、README、CI 和治理记录，未把外部或其他工作区资料当成本目标事实。
- 已审阅：本目标 `00-meta.md`、`01-decision.md`、`02-execution.md`、本文件；Root `00-meta.md`、`01-decision.md` D-012/D-013、I-005 附件；工作区页、`goal-tree.md`、VP-002、Charter、alignment、P-001～P-006 与 workspace protocol；并静态点验 API/Web 配置、README、health、CI、Docker/Compose/smoke 现状。
- 本 scope 是立项与信息质量审计；未运行应用、测试、Docker 或 15 分钟计时。没有把静态扫描扩写成运行时通过结论。

### 成果（有证据）

| 审计项 | 结论与证据 |
|--------|------------|
| 工作区与愿景链 | **通过**：workspace / Root / VP / Charter 机读链完整；本目标 `parent` 合法，未混入其他工作区状态。 |
| 立项目的与边界 | **总体合理**：S1～S5 承接 VP-002 #6/#7 的环境配置、Docker、≤15 分钟 fork、smoke 与阶段审计；S6 操作日志维持非阻断加分项。`progress: 0/5` 与零实施事实一致。 |
| I-005 当前工程事实 | **核心事实成立**：本地双进程、API env/fail-closed、公开 `/healthz`、Vite dev `/api` 代理、现有 CI 与运行时 DB gitignore 均有仓库证据；当前工作树未发现 Dockerfile、Compose、生产静态托管/反代或持久化 operation log。报告也明确本轮仅静态核对、未运行应用/测试。 |
| I-005 `verified` 的含义 | **边界合理**：D-013 的用户书面裁决足以关闭 Root 层“选哪条 R5 路线、采用什么计时与复现方向”的立项门禁；它不证明 Docker、smoke、CI 或 ≤15 分钟体验已实现。精确 env/Compose/探针/反代契约与计时/smoke 判据已由 `I-008-001/002` 继续 `open required`，未被静默放行。 |
| 实施与审计主张 | **通过**：`02-execution` 明记未改代码/配置/文档/容器/脚本、未运行测试；S1～S5 未勾选，Root R5 未勾选。本意见不产生实施或验收事实。 |

### 对照成功标准与信息门禁

| 项 | 审计结论 |
|----|----------|
| 目标意图、父级与最小可验证方向 | **满足**，足以立项并继续信息收集。 |
| 高层检查点与派生进度 | **满足**，S1～S5 可枚举且 `0/5` 可确定计算。 |
| `I-008-001` / `I-008-002` | **正确保持 open required**；当前只允许收集和冻结精确契约，不得实施或验收受影响范围。 |
| `I-008-003` | **边界合理**：仅在用户决定实施 S6 时阻断 S6，不阻断 S1～S5。 |
| I-005 报告可追踪性 | **有条件满足**：D-013 与当前状态可定位，但附件仍混有裁决前的当前时态，见 R-001。 |

### Findings

#### F-001 · Docker Compose 是核心交付还是可选加分路径，现有表述互相冲突

- **级别 / 严重度**：required / medium
- **关联门禁**：`I-008-001` 方案冻结；S1/S2 实施与验收；R5 / Root 关门。
- **证据**：VP-002 #7 把“Docker 一键启动”列入基础工程化成功标准；本目标 S2 是 S1～S5 五个核心检查点之一，要求交付 Compose，且 D-001 未选方案明确“仅做文档不落地容器与 smoke”会留下 VP #7 未交付。与此同时，D-001 边界又把“Docker 一键启动”写成“可选加分路径”；这与真正可选且不进分母的 S6 表述相同语义层级。
- **风险**：若按“加分项”解释，S2 可被跳过却仍宣称 `5/5` / R5 完成；若按核心检查点解释，D-001 的范围边界失真。该歧义会直接污染 `I-008-001` 的契约与关门判据。
- **要求**：由 `/govern` 在 GOAL-008 决策/审计响应中明确：建议采用“**Compose 是 R5 必须交付和验收的第二启动路径；对 fork 使用者而言可选择本地双进程或 Compose；完整生产拓扑 / CI-CD 仍为非目标**”。若用户确实要把 Compose 降为加分项，则须显式裁决并同步 S2、进度分母、Root D-013 与 VP-002 对齐边界，不能只保留现有含混措辞。

#### R-001 · I-005 v0.2.0 混用裁决前与裁决后的当前时态

- **级别 / 严重度**：recommended / medium
- **证据**：附件顶部与 §6 首段已写 D-013 完成裁决、`I-005: verified`；但 §2/§3 标题仍为“待用户裁决”，§6 随后仍写“保持 collecting”“裁决后再判断 verified”，frontmatter 仅列 `related_decision: D-012`。
- **影响**：Root `00-meta` / D-013 的 canonical 状态本身清楚，因此本项不重开 I-005，也不阻断立项；但附件作为 I-005 证据会让读者误判当前是否仍有 Root 层未决门禁。
- **建议**：把 §2～§4 与 §6 未决清单显式标为“D-013 前历史候选”，把过程性状态改为过去时，并在元数据/正文关联 D-013；保留候选比较，不要删除历史。

#### R-002 · `I-008-001/002` 关闭时应把高层方向变成可重复执行的契约

- **级别 / 严重度**：recommended / medium（现有 required 信息门禁的审计提示，不新增一扇门）
- **证据**：当前 `/healthz` 只证明 API 进程响应；现有 Web 代理仅为 Vite dev；CI 没有 Compose/smoke job；“不含依赖下载”尚未界定镜像 pull/build、Go/npm cache、Playwright 浏览器与首次 migration/seed。
- **建议**：`I-008-001` 至少冻结 production env/secrets、DB volume、Web SPA fallback 与 `/api` 反代、API/Web readiness、服务依赖/超时/失败行为和 CI 入口；`I-008-002` 至少冻结工具/平台基线、依赖缓存前提、计时起止、失败/重试规则、证据字段与 `scripts/smoke.sh` 的机器可判定退出码。未完成前维持现有 required 阻断。

### 必改项汇总

- **F-001（required / medium）**：明确 Compose 的交付义务与“可选”的对象。F-001 闭合前，不得关闭 `I-008-001`、实施/验收 S1/S2，亦不得据此勾选 R5 或关门。
- F-001 **不阻断** GOAL-008 保持 `active / 0/5`，也不阻断围绕 `I-008-001/002` 的信息收集。

### 与既有意见的异同

- 本目标此前无正式 self / independent 意见，因此不存在同 scope verdict 或 required finding 冲突。
- 同意 Root D-013：部署基线 A 与 15 分钟 / smoke 方向足以支持建立 R5 子目标；不同之处是本意见要求把 D-001 中“可选加分路径”与 S2 核心交付之间的语义冲突先闭合。
- 同意 I-005 对“候选/决策 ≠ 实施/验收”的边界；本意见不因 R-001 的历史措辞重开 Root `I-005`。

### 结论 + 建议给编排器/用户的下一步

- **conditional**：GOAL-008 的立项、S1～S5 结构、`I-008-001/002/003` 分层门禁和 I-005 的主要工程事实总体合理；当前可继续信息收集，但不可无条件冻结/实施工程化方案。
- 建议 `/govern` 先响应 A-001：修正 F-001；同步清理 I-005 的历史/当前时态（R-001）；随后以 R-002 作为 `I-008-001/002` 的最低收集清单。未来基于本 independent 意见推进门禁时，按 P-004 询问用户是否补同 scope self 审计。

### 声明

本意见仅追加独立审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。

### 响应 · A-001（/govern · 2026-08-02）

- **verdict 采纳**：`conditional` 成立——GOAL-008 立项、S1～S5 结构、`I-008-001/002/003` 分层门禁与 I-005 主要工程事实总体合理；F-001 使 D-001 的 Compose 交付义务表述冲突，需先闭合再冻结 `I-008-001`。本轮以 **fixed** 路径修正，未走 overruled/residual。
- **F-001 → fixed（Compose 交付义务澄清）**：
  - **明确口径**：**Docker Compose 是 R5 必须交付和验收的第二启动路径**（核心检查点 S2，计入 `0/5` 进度分母），**不是**像 S6 那样的可选加分项；对 fork 使用者可选择本地双进程或 Compose 两条启动路径；**完整生产拓扑 / CI-CD 部署流水线仍为非目标**。
  - **落盘**：GOAL-008 `01-decision` 新增 **D-002**；D-001 边界原「Docker 一键启动为可选加分路径，非强制生产拓扑」修订为「不创建完整生产运维 / CI-CD 部署流水线；Compose 为 R5 必须交付和验收的第二启动路径」；`00-meta` S2 同步为「交付 Docker Compose（必须的第二启动路径）」；Root I-005 附件 v0.2.1 与 GOAL-008 `00-meta` 概述对齐。
  - **证据路径**：本响应节；GOAL-008 `01-decision` D-002 / D-001 修订；`00-meta` S2 行；I-005 附件 v0.2.1 §2 F-001 澄清。
- **R-001 → handled（I-005 时态清理）**：I-005 附件 v0.2.1——§2/§3/§4 标题标注「D-013 前历史候选 · 已裁决」；§5 补裁决注记；§6 移除「保持 collecting」「裁决后再判断 verified」残余过程时态，改为历史候选清单 + 已裁决结论；frontmatter `related_decision: D-012` → `related_decisions: D-012, D-013`。
- **R-002 → handled（I-008-001/002 最低收集清单）**：GOAL-008 `00-meta` 信息表 `I-008-001`（production env/secrets、DB volume、Web SPA fallback 与 `/api` 反代、API/Web readiness、服务依赖/超时/失败行为、CI 入口）与 `I-008-002`（工具/平台基线、依赖缓存前提、计时起止、失败/重试规则、证据字段、`scripts/smoke.sh` 机器可判定退出码）已把 R-002 列为收集最低清单；后续关闭该两项时以此为准。
- **P-004 §3.1 处置**：A-001 为 `source: independent` 且本 scope 无 `source: self` 审计。本轮**不**放行下一阶段（`I-008-001` 仍 open），仅闭合 F-001 并处理 recommended——属整改闭环，不触发「仅用独立意见推进」门禁。**未来冻结 `I-008-001`（方案冻结门禁）或进入 S1/S2 实施前**，按 P-004 §3.1 询问用户是否补同 scope self 审计（A-001 亦建议如此）。
- **仍开放**：`I-008-001` / `I-008-002` / `I-008-003`（required · 阻断 S1/S2、S3/S4、S6 若实施）；Root R5 未勾选（Root 保持 `4/5`）；本目标 `active / 0/5`。
- **证据路径**：本响应节；GOAL-008 `01-decision` D-002/D-001；`00-meta`（S2、信息表）；I-005 附件 v0.2.1；02-execution 2026-08-02「响应 A-001」节。

## A-002 · A-001 F-001 与 R-001/R-002 响应证据独立复核（2026-08-02）

- **source**：independent
- **auditor**：Codex（GPT-5）
- **类型 / scope**：finding-closure；仅复核 A-001 **F-001**（Compose 交付义务）、**R-001**（I-005 时态与决策关联）、**R-002**（`I-008-001/002` 最低收集清单）的响应证据。不审计 Compose/smoke 实施，不冻结信息项，不放行 S1～S4，不评估 R5 / Root 关门。
- **verdict**：pass

### 范围与依据

- 工作区：`workspace-002-production-admin-foundation`；canonical root `docs/workspace-002-production-admin-foundation/`；Root 与 GOAL-008 `parent`、VP-002 / Charter 绑定未改变；`shared_materials_catalog: none`。
- 修订证据：提交 `f12973f6ae52e4f059c0153dda16e12e79448445`；复核时工作树 clean。已核对该提交差异、GOAL-008 `00-meta` / `01-decision` / `02-execution` / A-001 响应、Root `00-meta` / `02-execution`、I-005 v0.2.1 与 `goal-tree.md`。
- 本轮只核对文档修正与状态边界；未运行应用、测试、Docker、Compose、CI smoke 或 15 分钟计时，且没有把 finding closure 扩写成实现验收。

### 成果（有证据）

| A-001 项目 | 复核结果与证据 |
|------------|----------------|
| **F-001 required** · Compose 交付义务冲突 | **`fixed` 成立**：D-002 明确 Compose 是 R5 必须交付和验收的第二启动路径、对应 S2 且进入 `0/5` 分母；D-001 边界与 S2 已同步；fork 使用者可在双进程与 Compose 间选择，完整生产拓扑 / CI-CD 仍为非目标。VP-002 #7、Root I-005 投影与 I-005 v0.2.1 均同向。 |
| **R-001 recommended** · I-005 时态与关联 | **handled**：附件升级 v0.2.1，frontmatter 关联 D-012/D-013；§2～§4 标为 D-013 前历史候选，§6 明确已裁决并移除 `collecting` / “再判断 verified”的当前状态主张。候选比较被保留，没有改写历史。 |
| **R-002 recommended** · 信息收集最低清单 | **handled**：`I-008-001` 已纳入 production env/secrets、DB volume、SPA fallback/反代、readiness、依赖/超时/失败行为与 CI 入口；`I-008-002` 已纳入工具/平台、缓存、计时边界、失败/重试、证据字段与 smoke 退出码。两项仍为 `open required`，没有被文字补充冒充为 verified。 |
| 状态与进度边界 | **保持**：GOAL-008 `active / 0/5`，Root `active / 4/5`；S1～S5、Root R5 均未勾选，`I-008-001/002/003` 仍 open。A-001 响应只闭合 finding/处理建议，没有越过实施或验收门禁。 |

### Findings

- **无新 required**。

#### R-003 · 三处当前投影/历史短句仍可进一步消歧

- **级别 / 严重度**：recommended / medium
- **证据**：
  1. GOAL-008 `00-meta` 概述仍写“文档双进程为默认 + **可选 Docker Compose**”，没有像同文件 S2 那样点明“使用者可选、交付必需”；
  2. Root `00-meta` 当前进度说明仍写“R5 待立项”，但上一行与 goal-tree 已记 GOAL-008 立项；
  3. I-005 v0.2.1 §2 已标为历史候选且已说明 F-001，但段末仍以当前时态写“最终形态、镜像方案与是否纳入 R5 由用户裁决”，其中“是否纳入 R5”已由 D-013/D-002 决定。
- **影响**：D-002、S2、Root I-005 行与 §6 已足以确定当前权威口径，因此这些短句不推翻 F-001 `fixed`，也不重开 Root I-005；但全文检索或只读概述时仍可能误报“Compose 可跳过”或“R5 尚未立项”。
- **建议**：由 `/govern` 做一次窄幅投影清理：概述改为“Compose 必须交付、fork 使用者可选”；Root 改为“R5 已立项、待实施”；I-005 历史段末改为过去时，并明确精确镜像/Compose 契约仍由 `I-008-001` 冻结。

### 必改项汇总

- **无开放 required**（本 finding-closure scope）。A-001 F-001 的 `fixed` 关闭证据充分，可维持闭合。
- R-003 为 recommended，不阻断 GOAL-008 信息收集、`I-008-001/002` 后续方案冻结或目标状态；是否修正由 `/govern` 记录响应。

### 与既有意见的异同

- 同意 A-001 的 `conditional` 原判断，并确认其 required F-001 已按推荐路径修正；本意见不改写 A-001 历史 verdict。
- 同意 A-001 对 I-005“方向裁决 ≠ 实施/验收”的边界；I-005 `verified` 不关闭 `I-008-001/002`。
- 无 self / independent verdict 或 required finding 冲突；R-003 是关闭复核中新识别的非阻断投影卫生项。

### 结论 + 建议给编排器/用户的下一步

- **pass**：F-001 `fixed` 关闭成立，R-001/R-002 handled；本 scope 无开放 required。GOAL-008 可继续 `I-008-001/002` 信息收集，但仍不得把当前文档修正当作 S1～S4 实施或验收完成。
- 建议 `/govern` 记录采纳 A-002 `pass`，可选处理 R-003；在冻结 `I-008-001` 或进入 S1/S2 实施前，继续按 P-004 §4.1 询问用户是否补同 scope self 审计。

### 声明

本意见仅追加独立审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。

### 响应 · A-002（/govern · 2026-08-02）

- **verdict 采纳**：`pass` 成立——A-002 独立复核确认 A-001 **F-001 的 `fixed` 关闭成立**（D-002 + D-001 边界修订 + `00-meta` S2 + Root/I-005 投影对齐）；R-001/R-002 handled；本 scope 无开放 required、无新 required finding。
- **F-001 关闭复核确认**：Compose 为 R5 必须交付和验收的第二启动路径、对应 S2 且计入进度分母；D-001/S2/Root I-005/I-005 v0.2.1 同向；`fixed` 维持闭合。
- **R-003 → handled（投影/历史短句消歧）**：
  1. `GOAL-008 00-meta` 概述改为「文档双进程为默认；**Docker Compose 为 R5 必须交付的第二启动路径**，fork 使用者可选本地双进程或 Compose」；
  2. Root `00-meta` 进度说明由「R5 待立项」改为「R5 已立项 `GOAL-008-r5-engineering-fork`，待实施」；
  3. I-005 附件 v0.2.2 §2 末句改为过去时「最终形态与镜像方案已由 D-013（部署基线 A）决定；精确镜像 / Compose 契约由 `I-008-001` 在 GOAL-008 冻结」。
  - 均为文档投影清理；不改写历史候选、不重开 Root `I-005`。
- **P-004 §3.1 处置**：A-002 亦为 `source: independent`，本 scope 仍无 `source: self` 审计。本轮仅记录 finding-closure 采纳与 R-003 投影处理，**不**放行下一阶段。**未来冻结 `I-008-001`（方案冻结门禁）或进入 S1/S2 实施前**，按 P-004 §3.1 询问用户是否补同 scope self 审计。
- **仍开放**：`I-008-001` / `I-008-002` / `I-008-003`（required · 阻断 S1/S2、S3/S4、S6 若实施）；Root R5 未勾选（Root 保持 `4/5`）；本目标 `active / 0/5`。
- **证据路径**：本响应节；GOAL-008 `00-meta`（概述、S2、信息表）；Root `00-meta`（进度说明、I-005 行）；I-005 附件 v0.2.2 §2；02-execution 2026-08-02「响应 A-002」节。

## A-003 · GOAL-008 立项与 R5 方案边界 self 复核（2026-08-02）

- **source**：self
- **auditor**：/govern（self）
- **类型 / scope**：goal-definition + design-plan 复核；对 GOAL-008 立项信息与 Root **I-005 / D-012 / D-013** 方案边界，以及 **A-001 F-001 / R-001/R-002** 与 **A-002 R-003** 的响应闭合证据做 self 复核（P-004 §3.1 · 用户裁决「进行自审计」；同 scope 现有 A-001/A-002 independent，本自审补 `source: self` 覆盖）。不审计尚未发生的 R5 实施或关门，不复判 R1～R4。
- **verdict**：pass

### 范围与依据

- 工作区：`workspace-002-production-admin-foundation`；canonical root `docs/workspace-002-production-admin-foundation/`；Root `GOAL-001-production-admin-foundation` 与本目标 `parent` 一致；workspace/Root 绑定 `VP-002-production-admin-foundation`（`vision_role: delivery`、`primary_plan` 合法），`vision_ref` 匹配 Charter `schema-ui-core-admin-foundation@0.1.0`；`shared_materials_catalog: none`。
- 已复核：本目标五件套与 `01-decision` D-001/D-002、`00-meta` 成功标准/信息表；Root `00-meta`、`01-decision` D-012/D-013、I-005 附件 v0.2.2；A-001/A-002 与各自响应；工作区 `goal-tree.md`；代码 `apps/api/internal/config/config.go`、`.env.example`、`handler/health.go`、`apps/web/vite.config.ts`、`apps/web/README.md`、`.github/workflows/r6-basic-matrix.yml`。
- 本自审为文档/静态核对；**未运行**应用、测试、Docker、Compose、CI smoke 或 15 分钟计时；不把静态核对扩写成运行时通过结论。

### 成果（有证据）

| 审计项 | self 复核结论与证据 |
|--------|---------------------|
| 立项与边界 | **通过**：S1～S5 承接 VP-002 #6/#7（环境配置、Docker、≤15 分钟 fork、smoke、阶段审计），S6 操作日志维持非阻断加分；`progress: 0/5` 与零实施一致（`02-execution` 明记未改产品代码）。 |
| I-005 / D-013 方案边界 | **通过**：部署基线 A（Compose 必须交付、fork 用户可选）、建议计时口径、复现方法、I-006 方案甲均有用户书面裁决留痕（D-013）；I-005 附件 v0.2.2 与 Root `00-meta` 同向。 |
| F-001 `fixed` | **成立**：D-002 + D-001 边界修订 + `00-meta` S2 对齐——Compose 为 R5 必须交付和验收的第二启动路径（S2 核心检查点、计入进度分母），非可选加分项；完整生产拓扑/CI-CD 仍非目标。与 A-001/A-002 独立复核一致。 |
| R-001/R-002/R-003 handled | **成立**：I-005 附件 v0.2.2 时态清理（§2～§6「D-013 前历史候选 · 已裁决」、`related_decisions`）；`I-008-001/002` 信息表含 A-001 R-002 最低收集清单；GOAL-008 概述 / Root 进度说明 / I-005 §2 末句三处投影消歧已落实。 |
| 信息门禁边界 | **成立**：`I-008-001/002/003` 仍 open required，未被静默放行；Root `I-005: verified` 只解除立项/方案方向门禁，不充当 Docker、smoke、CI 或 ≤15 分钟体验已实现的证据。 |
| 状态与进度边界 | **成立**：GOAL-008 `active / 0/5`，Root `active / 4/5`；S1～S5、Root R5 均未勾选；A-001/A-002 响应只闭合 finding/处理建议，未越权放行实施或验收。 |

### 对照成功标准与信息门禁

| 项 | self 复核结论 |
|----|---------------|
| 立项意图、父级与最小可验证方向 | **满足**，可继续信息收集。 |
| 高层检查点与派生进度 | **满足**，S1～S5 可枚举且 `0/5` 可确定计算。 |
| `I-008-001` / `I-008-002` | **正确保持 open required**（自审时）；本轮随 D-003 冻结 `I-008-001`，`I-008-002` 继续阻断 S3/S4。 |
| `I-008-003` | **边界合理**：仅在用户决定实施 S6 时阻断 S6。 |
| 进入 I-008-001 方案冻结 | **可进入**：本 self 复核 + A-001/A-002（independent）同向；方案冻结本身由 D-003 决策记录，不由本自审代替。 |

### Findings

- **无新 required**。
- **注记（recommended · 非阻断）**：D-003 冻结 `I-008-001` 后，S1/S2 实施的精确 nginx.conf / Dockerfile / 镜像 tag 属实施细节，由实施留痕；建议 S2 完成后做一次实施向审计（self 或 `/audit`），并在 S3/S4 前关闭 `I-008-002`。

### 必改项汇总

- **无开放 required**（本 scope）。

### 与既有意见的异同

- 相对 A-001/A-002（independent）：结论一致——立项、S1～S5、`I-008-001/002/003` 分层门禁与 I-005 主要工程事实总体合理；F-001 `fixed`、R-001/R-002/R-003 handled 闭合证据经本自审复核成立。本自审补齐同 scope 的 `source: self` 覆盖，满足 P-004 §3.1。
- 无 self/independent verdict 或 required finding 冲突。

### 结论 + 建议下一步

- **pass**：GOAL-008 立项与 I-005/D-013 方案边界、A-001/A-002 响应闭合证据经 self 复核成立；当前 scope 无开放 required，可进入 `I-008-001` 方案冻结与 S1/S2 实施边界判断。
- 建议 `/govern`：记录 D-003 冻结 `I-008-001`（`verified`）；进入 S1 实施（环境/配置基线文档）前保持 `I-008-002` 对 S3/S4 的阻断；S2 完成后建议一次实施向审计。

### 声明

本自审仅追加审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；finding 响应、方案冻结与阶段推进由 `/govern` 处理。

## A-004 · S1/S2 实施向审计：环境/配置基线与容器一键启动（2026-08-02）

- **source**：independent
- **auditor**：Claude（Opus 4.8）
- **类型 / scope**：execution-facts；审计 GOAL-008 **S1**（环境与配置基线，契约 **C-001/C-002**）与 **S2**（容器与一键启动，契约 **C-003～C-007**）的实施事实，对照 [I-008-001-engineering-contract.md](attachments/I-008-001-engineering-contract.md) v1.0.0 验收清单；不评估 S3～S6、`I-008-002/003` 或 Root R5 关门。
- **verdict**：pass

### 范围与依据

- 工作区：`workspace-002-production-admin-foundation`；canonical root 一致；Root `GOAL-001-production-admin-foundation`；`shared_materials_catalog: none`。`goal-tree` 与 `00-meta` 一致：GOAL-008 `active / 2/5`。
- 已审阅：本目标 `00-meta`、`01-decision`（D-001～D-003）、`02-execution`（2026-08-02「实施 S1 + S2」节）、本文件；I-008-001 契约 v1.0.0（§1～§7、C-001～C-007）；实施物 `apps/api/.env.example`、`apps/api/internal/config/config.go`、`apps/api/internal/handler/health.go`、`apps/api/cmd/server/main.go`、`apps/api/Dockerfile`、`apps/web/Dockerfile`、根 `compose.yaml`、`apps/web/nginx.conf`、根 `.dockerignore`、`apps/api/.dockerignore`、`.github/workflows/r6-basic-matrix.yml`（`container-smoke`）；`apps/api` / `apps/web` / 根 README 的 S1 文档段；`apps/api/internal/handler/auth.go`、`records.go`（核对 CI smoke 端点形状）。
- 验证方式：静态核对契约↔实施物；本机复跑——API `go build`/`go vet`/`go test ./...` 全绿、web `npm run build`（tsc -b && vite build）成功；`docker compose config --quiet` 通过且缺 secret 时 fail-closed abort；`docker compose build` 两镜像成功（api 走 BuildKit 缓存）；`up -d` + 冒烟（healthz、nginx 反代登录、`/me`、SPA fallback、restart/down-up 持久化）；结束后 `down -v` 清理，git 工作树无新增改动。
- 本意见是实施事实审计：本机运行仅用于核对已记录的验证主张，不替代 `I-008-002` 的 smoke 判据，也不勾选任何检查点。

### 成果（有证据）

| 契约项 | 审计结论与证据 |
|--------|----------------|
| C-001 `.env.example` 与 config.go 键一致、dev/prod 注释齐全 | **通过**：13 键（APP_NAME / APP_ENV / HTTP_ADDR / HTTP_READ/WRITE/IDLE_TIMEOUT / LOG_LEVEL / AUTH_JWT_SECRET / AUTH_ACCESS_TTL / AUTH_REFRESH_TTL / DB_PATH / ADMIN_INITIAL_PASSWORD / AUTH_DEV_SESSION_ENABLED）逐一与 `config.Load()` 键一致；dev 默认 vs production 必填/fail-closed 注释齐全；fail-closed 由 `main.go` `resolveJWTSecret`/`resolveSeedHash` 真实执行（非 development 且缺失 → 报错退出 1）。 |
| C-002 `/healthz` 200 + `{"status":"ok"}`（本地与容器内） | **通过**：`healthz()` 返回 `200 {"status":"ok","timestamp":...,"version":...,"commit":...}`；容器内实测 `curl :8080/healthz` → `{"status":"ok","timestamp":"2026-08-02T15:55:07Z","version":"0.1.0","commit":"unknown"}`。 |
| C-003 `docker compose up` 后 api healthy、web 200、登录种子 admin | **通过（本机复跑）**：`up -d` 后 api `Healthy`、web 200；经 nginx 反代 `POST :8081/api/auth/login` 成功（accessToken 176 字符），`GET :8081/api/accounts/me` → `user: Admin`。 |
| C-004 Dockerfile ×2 多阶段 + 根 compose.yaml 与 §3 一致 | **通过**：api 多阶段（golang:1.26-alpine → alpine:3.20，`CGO_ENABLED=0` 静态、非 root `app` 用户、`DB_PATH=/app/data`）；web 多阶段（node:22 → nginx:1.27-alpine，仓库根 context + `COPY docs/schemas` 解决 `@schemas` 别名）；compose 服务 api/web、`db-data` 命名卷挂 `/app/data`、healthcheck 探针、`depends_on service_healthy`、`restart: on-failure`；`docker compose config` 通过，两镜像构建成功。 |
| C-005 nginx SPA fallback + `/api` 反代；刷新 `/list-edit-lifecycle` 可回退 | **通过（本机复跑）**：`nginx.conf` `location / { try_files $uri $uri/ /index.html; }` + `location /api { proxy_pass http://api:8080; }`；实测 `:8081/` 与 `:8081/list-edit-lifecycle` 均返回含 `id="root"` 的 index，登录/`/me` 经反代成功。 |
| C-006 重启/down-up 后 DB 数据保持 | **通过（本机复跑）**：创建 `rec-3bb6770ac47aff74` 后，`docker compose restart api` 与 `down`→`up` 两次均经反代 GET 该记录成功，`db-data` 命名卷保持数据。 |
| C-007 CI 增加容器/smoke 入口 job | **通过**：`r6-basic-matrix.yml` 新增 `container-smoke` job（build → up → 等 ready → 经 nginx 反代 login + `/me` → SPA/route fallback → restart 持久化 → `down -v` teardown），job env 提供 fail-closed secret。CI job 尚无 GitHub Actions 运行记录（`02-execution` 如实只记「本机验证」），C-007 判据（job 存在 + 本地 smoke 通过）满足。 |

**其它核验**：契约 §1 的 Web 相对路径 `/api/*`（`auth-client.ts` 硬编码 `LOGIN_URL="/api/auth/login"` 等）与 §4 同源反代一致；CI smoke 的端点形状与 handler 实测吻合（`accessToken` 字段、`POST /api/records` 校验 `name/status/owner`、detail 返回 `id`）。

### 对照成功标准与信息门禁

| 项 | 审计结论 |
|----|----------|
| S1 环境与配置基线 | **成立**：env 清单、health/启动验证、dev/prod 区分文档齐备，与 C-001/C-002 一致。 |
| S2 容器与一键启动 | **成立**：Compose 第二启动路径交付物完整且本机可复现，C-003～C-007 全通过。 |
| `I-008-001` 契约↔实施一致性 | **一致**：契约 §1～§7 的字面要求均落到实施物；实施细节（镜像 tag、BuildKit cache、`COPY docs/schemas`）留痕于 `02-execution`，未改动契约形状。 |
| 状态与进度边界 | **保持**：S1/S2 勾选（`0/5 → 2/5`）与已核实实施事实相符；`I-008-002/003` 仍 open required（阻断 S3/S4、S6 若实施）；Root R5 未勾选（`4/5`）。本意见不改动任何状态。 |

### Findings

- **无新 required**（C-001～C-007 全部通过，实施主张可复现）。

#### R-001 · 「AUTH_DEV_SESSION_ENABLED 生产禁止启用」靠部署约定而非运行时硬门禁

- **级别 / 严重度**：recommended / low
- **证据**：契约 §5 与 `.env.example`/README 写「生产禁止启用」（M9）；实际 enforce 是 compose + api Dockerfile 显式 `AUTH_DEV_SESSION_ENABLED=false`，而 `auth.New` 只按 flag 显式 opt-in，无「`APP_ENV=production` 且 flag=true → 拒绝启动」的运行时检查。
- **影响**：不在 C-001～C-007 验收内，且 S2 交付的 compose 路径与镜像均已置 false，不构成 S1/S2 缺口；属 R2 既有边界。非 compose 直接跑生产二进制并显式开启时不会被硬拦。
- **建议**：可于后续在 `main.go` 增加生产守卫（`AppEnv==production && AuthDevSessionEnabled` → 启动报错），或将契约措辞收敛为「所有交付的生产路径必须 false（由部署强制）」，避免被读成硬保证。

#### R-002 · `docker compose down` 也受 fail-closed 插值影响（新 shell 需重复 export secret）

- **级别 / 严重度**：recommended / low
- **证据**：`${AUTH_JWT_SECRET:?}` / `${ADMIN_INITIAL_PASSWORD:?}` 对整份 compose 生效，本审计 teardown 时在未带 secret 的 shell 里 `docker compose down` 直接 abort（`required variable ... missing a value`）。
- **影响**：README/compose 注释已说明可写入仓库根 `.env`（gitignored）或 export；但仅靠 export 的 fork 用户在新 shell `down`/`config` 会遇到 fail-closed 报错，属轻微 UX 脚枪。
- **建议**：在根 `README.md`「Docker Compose 一键启动」段加一句「将密钥写入仓库根 `.env`（gitignored）可避免每次 `down`/`config` 重复 export」。

#### 注记（非 finding）

- CI `container-smoke` job 尚无 GitHub Actions 运行记录；建议在 S4 回归或 Root R5 关门证据里留一次 CI 运行记录。
- api 镜像 `VERSION/COMMIT/BUILT_AT` 未由 compose 注入 build args → healthz 显示 `0.1.0` / `unknown`；纯装饰性，契约未要求 provenance。

### 必改项汇总

- **无开放 required**。S1/S2 的 C-001～C-007 已核对/复跑通过；R-001/R-002 为 recommended 非阻断，由 `/govern` 决定是否响应。

### 与既有意见的异同

- 同意 A-001/A-002/A-003 对「契约冻结 ≠ 实现完成」的边界：I-008-001 verified 只解除方案冻结门禁，S1/S2 是否完成以本实施向审计的 C-001～C-007 为准。
- 本意见是 GOAL-008 首个 execution-facts 实施向意见；与既有 self/independent 无 verdict 或 required finding 冲突。

### 结论 + 建议给编排器/用户的下一步

- **pass**：S1/S2 的实施主张与仓库事实一致且可复现——C-001/C-002 静态与代码核对通过，C-003/C-005/C-006 本机 Docker 复跑通过，C-004 镜像构建成功，C-007 CI job 结构与本地 smoke 通过；无开放 required。`0/5 → 2/5` 的进度勾选有据。
- 建议 `/govern` 响应采纳 A-004 `pass`（可选处理 R-001/R-002）；下一拍按 `02-execution` 计划收集并冻结 `I-008-002`（15 分钟计时复现协议 + smoke 判据），再实施 S3（fork 文档 + 独立复现记录）与 S4（`scripts/smoke.sh` 正式化）。

### 声明

本意见仅追加独立审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。

## A-005 · S1/S2 实施复核：生产开发会话门禁与派生进度一致性（2026-08-03）

- **source**：independent
- **auditor**：Codex（GPT-5）
- **类型 / scope**：execution-facts；复核 `workspace-002-production-admin-foundation` 中 `GOAL-008-r5-engineering-fork` 的 S1（C-001/C-002）与 S2（C-003～C-007），对照 `attachments/I-008-001-engineering-contract.md` v1.0.0；不审 S3～S6、`I-008-002/003` 或 Root R5 关门。
- **verdict**：fail

### 范围与依据

- workspace、Root、VP-002、Charter、alignment、principles 与 workspace protocol 链已核对；本目标五件套、`goal-tree.md` 与 `I-008-001` 均在当前 canonical root 内，`shared_materials_catalog: none`。
- `I-008-001` 为 `verified`；`I-008-002` 仍为 open required（阻断 S3/S4），`I-008-003` 仍为 open conditional（若实施则阻断 S6）。本意见不把它们静默关闭。
- 独立复跑：`apps/api` 的 `go test ./...` 通过；`apps/web` 的 `npm test -- --run` 通过（23 个文件、458 个测试）；唯一 Compose 项目的 `config --quiet`、镜像构建、启动和 teardown 通过。运行时依次验证 `/healthz`、Web 200、登录与 `/api/accounts/me`、SPA fallback、API 重启后的持久化、Compose down/up 后持久化，以及删除后 down/up 不复活；测试资源已清理。
- 另以 `APP_ENV=production` 且 `AUTH_DEV_SESSION_ENABLED=true` 启动 API，未认证请求 `GET /api/accounts/me` 返回 200 和静态开发身份 `dev-001`。`internal/auth/auth.go` 的 middleware 在无效/缺失 bearer 时注入 `StaticDevSession()`，而 `internal/account/session.go` 赋予该身份 `admin/editor`、`records.read` 和 `records.write`；`main.go` 只对生产 JWT/初始密码做 fail-closed，没有禁止该 flag 的生产运行时门禁。

### 成果（有证据）

| 契约项 | 独立复核结论 |
|--------|--------------|
| C-001 `.env.example` / `config.go` / dev-prod 配置行为 | **不通过**：键清单和缺 secret fail-closed 成立，但契约 §1、§5 要求生产禁止启用 `AUTH_DEV_SESSION_ENABLED`，上述运行时反例不满足。 |
| C-002 `/healthz` | **通过**：本机及容器内均返回 200 和 `status: ok`。 |
| C-003 Compose 启动、健康、登录、`/me` | **通过**：独立项目启动后 API healthy、Web 200，反代登录和 `/api/accounts/me` 成功。 |
| C-004 双 Dockerfile 与 Compose 结构 | **通过**：两镜像构建成功，服务、卷、healthcheck、依赖关系符合契约。 |
| C-005 nginx API 反代与 SPA fallback | **通过**：根路径及 `/list-edit-lifecycle` 返回 SPA，`/api` 反代可用。 |
| C-006 重启/down-up 持久化与删除不复活 | **通过**：创建记录跨 API 重启及 Compose down/up 保留，删除后不复活。 |
| C-007 CI `container-smoke` 入口 | **通过**：job 与本地 smoke 路径存在并可复跑；没有把缺少远端 Actions 运行记录误记为证据。 |

### 对照成功标准与信息门禁

- **S1 不能无条件判定通过**：C-001 的生产安全断言被运行时反例否定；C-002 通过不抵消该门禁缺口。
- **S2 的实现证据成立**：C-003～C-007 全部通过，且包含完整运行时序列；这不解除 `I-008-002` 对 S3/S4 的 required 阻断。
- **派生进度不一致**：`GOAL-008-r5-engineering-fork/00-meta.md` frontmatter 仍为 `progress: 0/5`，但其中 S1/S2 已勾选、正文写 `0/5 → 2/5`，`goal-tree.md` 树和状态表均为 `2/5`。这是治理投影错误，不应由本意见直接修复或改状态。

### Findings

#### F-002 · 生产环境未硬拒绝开发会话（required / high）

- **状态**：open；未按 `fixed`、`accepted-residual` 或 `user-overruled` 合法闭合。
- **证据**：`I-008-001` §1、§5 明确生产禁止 `AUTH_DEV_SESSION_ENABLED=true`；实际生产配置可启动该 flag，并使未认证请求获得开发静态高权限身份。
- **风险**：部署者或环境注入错误 flag 时，认证边界被绕过并暴露写权限，不能仅作为 Compose 默认值或文档约定处理。
- **建议修复**：在配置加载或 server 启动处对 `APP_ENV=production && AUTH_DEV_SESSION_ENABLED=true` 明确报错退出；补生产配置回归测试，并重新执行 C-001 及相关 smoke。

#### F-003 · `00-meta` 与 `goal-tree` 的进度投影不一致（required / medium）

- **状态**：open；未合法闭合。
- **证据**：`00-meta.md` frontmatter 为 `0/5`，而成功标准勾选、正文和 `goal-tree.md` 为 `2/5`；A-004 中“`goal-tree` 与 `00-meta` 一致”的事实陈述因此不成立。
- **风险**：下游 Root gate、审计索引和用户界面可能读取不同进度，削弱治理记录的可追溯性。
- **建议修复**：由 `/govern` 按已核实检查点同步 `00-meta` 与 `goal-tree`，然后复核所有引用；不得把同步本身当作 S1/S2 技术 finding 的关闭证据。

### 必改项汇总

- **F-002 required / high：开放，阻断 S1 无条件通过及相关放行。**
- **F-003 required / medium：开放，阻断本目标进度投影的可信闭合。**
- `I-008-002` 的既有 required 门禁仍开放，继续阻断 S3/S4。

### 与既有意见的异同

- A-004 对 C-002～C-007 的通过结论与本复核一致，且本意见没有删除或改写 A-004。
- A-004 将同一 `AUTH_DEV_SESSION_ENABLED` 生产安全事实列为 `recommended / low` R-001；本意见依据契约的“生产禁止启用”硬要求及运行时反例，将其列为 `required / high` F-002。该 required 与非阻断分类构成 P-004 §4.2 的相关意见分歧，不能静默裁决。
- A-004 关于“`goal-tree` 和 `00-meta` 一致”的陈述被 F-003 的当前文件事实否定；这不改变 A-004 记录的历史内容，只要求 `/govern` 响应时纠正投影。

### 结论与建议下一步

- **fail**：S2（C-003～C-007）可复现通过，但 S1 的生产开发会话门禁不成立，且进度投影不一致；因此不能把 GOAL-008 S1/S2 整体作为无条件通过。
- `/govern` 应先展示 A-004 R-001 与本意见 F-002 的分类分歧，按 P-004 等待用户裁决；建议选择 `fixed`，修复生产硬门禁并同步 F-003，再做针对性复审。`I-008-002` 仍需独立收集并冻结后，方可进入 S3/S4。

### 声明

本意见仅追加独立审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。

## A-006 · S1/S2 实施复核 self：生产守卫与进度投影（2026-08-03）

- **source**：self
- **auditor**：Claude（Opus 4.8）
- **类型 / scope**：execution-facts 复核；按用户 P-004 §3.1 裁决补齐 GOAL-008 **S1/S2 实施 scope** 的 `source: self` 覆盖（A-004/A-005 均为 independent；A-003 self 仅覆盖立项/方案边界）。复核 A-004/A-005 证据与 F-002/F-003 `fixed` 关闭证据；不审 S3～S6、`I-008-002/003` 或 Root R5 关门。
- **verdict**：pass

### 范围与依据

- 工作区/Root/VP-002/Charter/alignment 链已在 A-004/A-005 核对一致（`vision_role: delivery`、`primary_plan: VP-002`、`shared_materials_catalog: none`）；本 self 复核聚焦 S1/S2 实施事实与两项 `fixed` 关闭证据。
- 已审阅：`internal/config/config.go`（`ValidateProd`）、`internal/config/config_test.go`、`cmd/server/main.go` 接线、`compose.yaml` / `apps/api/Dockerfile` 的 `APP_ENV=production` + `AUTH_DEV_SESSION_ENABLED=false`、`00-meta.md` frontmatter、根 README「Docker Compose 一键启动」段。
- 验证方式：`go build` / `go vet` / `go test ./...`（apps/api）全绿（config 包新增 4 用例通过）；运行时复验——`APP_ENV=production AUTH_DEV_SESSION_ENABLED=true` 启动即 `ERROR startup failed`（exit 1）；`APP_ENV=development AUTH_DEV_SESSION_ENABLED=true` 正常启动（`dev_session: true`）。

### Findings

- **F-002（required/high）→ `fixed`**：`config.ValidateProd()` 在非 `development` 环境且 `AUTH_DEV_SESSION_ENABLED=true` 时返回启动错误，`main.go` 于日志初始化后立即校验并 `os.Exit(1)`；回归用例覆盖 development / production / staging 分支；运行时反例（A-005 举证）复验为拒绝启动。契约 §1/§5「生产禁止启用」现由硬门禁成立，不再仅靠部署约定。
- **F-003（required/medium）→ `fixed`**：`00-meta.md` frontmatter `progress: 0/5 → 2/5`，与成功标准勾选、派生进度段、`goal-tree.md`（`active / 2/5`）一致；未把投影同步当作 S1/S2 技术 finding 的关闭证据（F-002 关闭证据独立成立）。
- A-004 R-001（recommended/low，与 F-002 同事实）→ handled：分类分歧按用户 P-004 §3.2 裁决采用 F-002 `required/high` 口径并修复；R-001 自身建议的「main.go 生产守卫」正是本次落地。
- A-004 R-002（recommended/low）→ fixed：根 README 补「将密钥写入仓库根 `.env` 可避免新 shell 重复 export」注记。
- C-002～C-007 维持 A-004/A-005 独立复跑的通过结论；本 self 不复跑 Docker 冒烟（守卫仅作用于 `flag=true` 路径，compose/镜像均显式 false，不受影响）。

### 对照成功标准与信息门禁

- S1（C-001/C-002）：C-001 的「生产禁止启用 dev session」已由运行时守卫成立；C-002 维持通过。S1 判定恢复为**成立**。
- S2（C-003～C-007）：A-004/A-005 独立复跑通过，本 self 未改动 S2 实施物，维持**成立**。
- `I-008-001` verified 维持；`I-008-002`/`I-008-003` 仍 open required（阻断 S3/S4、S6 若实施）。本 self 不改动任何状态/progress。

### 声明

本自审仅追加审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。

## 统一响应 · A-004 / A-005（2026-08-03）

按 P-004 §3.1 用户裁决「需要补 self」先行 A-006 自审，再统一响应；§3.2 分类分歧按用户裁决采用 **F-002 `required/high`** 口径并修复。

### 响应 A-004（independent · execution-facts · pass）

- **采纳 `pass`**：S1/S2 的 C-001～C-007 实施主张与仓库事实一致且可复现。
- **R-001（recommended/low）→ handled**：与 A-005 F-002 同事实的分类分歧按用户裁决采用 `required/high` 口径，并由 F-002 的 `fixed` 修复承载（生产运行时守卫落地），不再单独处理。
- **R-002（recommended/low）→ fixed**：根 README「Docker Compose 一键启动」段补注记「将密钥写入仓库根 `.env`（gitignored）可避免新 shell 里 `docker compose config`/`down` 因 fail-closed 插值重复 export」。
- **注记**：CI `container-smoke` 尚无 GitHub Actions 运行记录，留待 S4 回归或 Root R5 关门证据（沿用 A-004 注记）。

### 响应 A-005（independent · execution-facts · fail）

- **采纳 `fail`**：F-002/F-003 均按 `fixed` 合法闭合。
- **F-002（required/high）→ fixed**：
  - 新增 `config.ValidateProd()`（`internal/config/config.go`）：`AppEnv != "development" && AuthDevSessionEnabled` → 返回启动错误；`main.go` 于 `config.Load()` 后立即校验，错误则 `logger.Error` + `os.Exit(1)`。
  - 回归测试 `internal/config/config_test.go`：4 用例覆盖 development 允许 / production+flag fail-closed / production 无 flag 通过 / 其他非生产环境（staging）fail-closed。
  - 运行时复验：`APP_ENV=production AUTH_DEV_SESSION_ENABLED=true` → `startup failed: AUTH_DEV_SESSION_ENABLED must be false when APP_ENV="production"`（exit 1）；`development + flag=true` → 正常启动（`dev_session: true`）。
  - 回归：`go build` / `go vet` / `go test ./...`（apps/api）全绿；compose 与 api 镜像均显式 `AUTH_DEV_SESSION_ENABLED=false`，守卫不影响第二启动路径。
- **F-003（required/medium）→ fixed**：`00-meta.md` frontmatter `progress: 0/5 → 2/5`，与勾选成功标准、派生进度段、`goal-tree.md` 一致；复核引用无其它残留 `0/5`。
- **A-004「goal-tree 与 00-meta 一致」陈述**：当前文件事实已由 F-003 `fixed` 恢复为一致；不改写 A-004 历史记录。

### 关闭证据与后续

- F-002/F-003 关闭证据路径：代码 + 回归测试 + 运行时复验（`02-execution` 同步留痕）；**建议做一次 `/audit` finding-closure 复审**确认关闭成立。
- `I-008-002`（计时复现协议 + smoke 判据）仍 open required，为进入 S3/S4 的前置门禁；下一拍收集并冻结后再实施 S3/S4。
- 本目标维持 `active / 2/5`；Root R5 未勾选（Root `4/5`）。

## A-007 · F-002/F-003 关闭证据独立复审（2026-08-03）

- **source**：independent
- **auditor**：Codex（GPT-5）
- **类型 / scope**：finding-closure；仅复审 A-005 的 **F-002**（生产环境开发会话硬门禁）与 **F-003**（GOAL-008 派生进度投影）`fixed` 关闭证据。不审 S3～S6、`I-008-002/003`、Root R5 或 VP-002 关门。
- **verdict**：pass

### 范围与依据

- 工作区：`workspace-002-production-admin-foundation`；canonical root `docs/workspace-002-production-admin-foundation/`；目标 `parent`、workspace Root 和 `goal-tree.md` 的当前投影一致；`shared_materials_catalog: none`，未使用共享资料作为关闭依据。
- 已核对：本目标 `00-meta.md`、`01-decision.md`、`02-execution.md`、本台账 A-004～A-006 与统一响应、`I-008-001-engineering-contract.md` §1/§5/C-001；当前 `goal-tree.md`；`apps/api/internal/config/config.go`、`config_test.go`、`cmd/server/main.go` 与根 `compose.yaml`。
- 独立验证：在 `apps/api` 运行 `go test ./...`、`go vet ./...`、`go build ./...` 均通过；`TestValidateProd` 的 development / production / staging 四个分支通过；以 `APP_ENV=production AUTH_DEV_SESSION_ENABLED=true go run ./cmd/server` 复验，进程输出 `startup failed` 并以预期 exit 1 拒绝启动。
- 本复审不把 F-002/F-003 的关闭扩写为 S3/S4 已实施、`I-008-002/003` 已关闭、Root R5 已勾选或目标可关门。

### 成果（有证据）

| Finding | 独立复核结果与证据 |
|---------|--------------------|
| **F-002 required/high** · 生产环境未硬拒绝开发会话 | **`fixed` 成立**：契约要求生产禁止 `AUTH_DEV_SESSION_ENABLED`；`Config.ValidateProd()` 对所有非 `development` 环境的该 flag 返回启动错误，`main.go` 在加载配置后、创建认证器前执行该校验并 `os.Exit(1)`。4 个回归用例覆盖 development 允许、production flag=true 拒绝、production flag=false 允许与 staging 拒绝；实际 production 反例亦得到预期拒绝。Compose 仍显式设置该 flag 为 `false`。 |
| **F-003 required/medium** · 派生进度投影不一致 | **`fixed` 成立**：`00-meta.md` frontmatter 为 `progress: 2/5`；S1/S2 是五个核心检查点中仅有的两项已勾选；派生进度段同样为 `2/5`；`goal-tree.md` 的 ASCII 树与状态表均投影 `active / 2/5`。检索到的 `0/5` 均为立项或修复前的历史时间线，不与当前投影冲突。 |

### Findings

- **无新 required 或 recommended finding**。

### 必改项汇总

- 本 finding-closure scope 无开放 required；F-002/F-003 可维持 `fixed`。
- `I-008-002` 仍为 S3/S4 的 open required 门禁，`I-008-003` 仍在 S6 实施时适用；它们不因本意见而关闭或降级。

### 与既有意见的异同

- A-005 的 `fail` 记录修复前的运行时反例与投影不一致，保留为历史事实；统一响应提供 `fixed` 路径，A-006 提供同 scope self 覆盖。
- 本意见独立复跑代码与测试证据后确认关闭成立，不以 A-006 的 self 结论代替交叉复审，也不改写 A-004/A-005 的历史 verdict。

### 结论 + 建议给编排器/用户的下一步

- **pass**：F-002/F-003 的关闭证据充分且可重复核对，两项可维持 `fixed`；本 scope 无开放 required。
- 建议 `/govern` 记录 A-007 的 `pass`，随后仍按现有信息门禁收集并冻结 `I-008-002` 后再进入 S3/S4；本意见不授权任何状态、进度或 Root R5 变更。

### 声明

本意见仅追加独立审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。

### 响应 · A-007（/govern · 2026-08-03）

- **verdict 采纳**：`pass` 成立——A-007（independent · finding-closure）独立复核 F-002/F-003 关闭证据（代码、回归、运行时与进度投影），确认两项 `fixed` 维持闭合，本 scope 无开放 required、无新 required/recommended finding。
- **F-002（required/high · 生产环境开发会话硬门禁）关闭复核确认**：`Config.ValidateProd()` 对所有非 `development` 环境的 `AUTH_DEV_SESSION_ENABLED=true` 返回启动错误，`main.go` 于 `config.Load()` 后校验并 `os.Exit(1)`；4 回归用例覆盖 development/production/staging 分支，production 反例复验拒绝启动；Compose 与 api 镜像显式 `false`。`fixed` 维持闭合。
- **F-003（required/medium · 派生进度投影）关闭复核确认**：`00-meta` frontmatter `progress: 2/5` 与成功标准勾选、派生进度段、`goal-tree.md`（`active / 2/5`）一致，检索无当前投影残留。`fixed` 维持闭合。
- **P-004 §3.1 处置**：A-007 为 `source: independent`；同 scope（F-002/F-003 fixed 关闭证据）已有 **A-006（self · execution-facts 复核）** 覆盖，无需再补自审。本 scope 无意见冲突、无 required finding，不触发 §3.2/§3.3。
- **仍开放**：`I-008-002`（required · 阻断 S3/S4）；`I-008-003`（required · 仅当 S6 实施）；Root R5 未勾选（Root 保持 `4/5`）；本目标 `active / 2/5`。
- **下一步**：收集并冻结 `I-008-002`（15 分钟计时复现协议 + smoke 判据），随后实施 S3（fork 文档 + 独立复现记录）与 S4（`scripts/smoke.sh` 正式化）；建议对 `I-008-002` 冻结做一次方案冻结审计（self 或 `/audit`）。
- **证据路径**：本响应节；`03-audit.md` A-005/A-006/A-007；`02-execution` 2026-08-03「响应 A-007」节。

## A-010 · S3/S4 实施向自审（execution-facts · 2026-08-03）

- **source**：self
- **auditor**：/govern（self）
- **类型 / scope**：execution-facts；审计 GOAL-008 **S3**（QUICKSTART fork 文档 + 独立复现记录）与 **S4**（`scripts/smoke.sh` + 本地/CI smoke + disposable `SM-006` 证据）的实施事实，对照 [I-008-002-fork-reproduction-protocol.md](attachments/I-008-002-fork-reproduction-protocol.md) **v0.1.1**（§2.2 路径、§3 计时、§4 复现字段、§5 smoke 契约）与成功标准 S3/S4 字面；不审 S5、S6、`I-008-003` 或 Root R5/VP-002 关门。
- **verdict**：pass

### 范围与依据

- 工作区：`workspace-002-production-admin-foundation`；canonical root 一致；Root `GOAL-001-production-admin-foundation`；`shared_materials_catalog: none`。`goal-tree` 与 `00-meta` 同步后一致：GOAL-008 `active / 4/5`。
- 已复核：根 `QUICKSTART.md`；[R5-S3-REPRO-001](attachments/R5-S3-REPRO-001.md)（含 attempt/source/path/platform/cache precondition/timing/checks/secrets/result 全字段）；[r5-repro-endpoint4.mjs](attachments/r5-repro-endpoint4.mjs) 与截图 [r5-repro-endpoint4.png](attachments/r5-repro-endpoint4.png)；`scripts/smoke.sh`；[r5-smoke-disposable-local.txt](attachments/r5-smoke-disposable-local.txt)；`.github/workflows/r6-basic-matrix.yml` `container-smoke` 修改；`00-meta` S3/S4 行与信息表 `I-008-002`；`01-decision` D-006；I-008-002 协议 v0.1.1；Root D-013 与 I-005 §3/§4。
- 验证方式：实施时**真实运行**——`docker compose down -v` 后全新 `up -d`（compose 路径）；四终点实测（healthz 200/`status:ok`、login 200 + accessToken 176 字符、`/me` 200 + user + features、Playwright Chromium 浏览器打开 `/list-edit-lifecycle` 标题 + cell `Acme Console`）；`bash scripts/smoke.sh --disposable` → SM-001～006 全 PASS、EXIT=0；非 disposable 默认路径 SM-001～005 PASS + SM-006 SKIP；`go test ./... -count=1`（apps/api）全绿、web `vitest run` 458/458；workflow YAML 解析通过。

### 成果（有证据）

| 审计项 | self 结论与证据 |
|--------|-----------------|
| QUICKSTART 覆盖（S3 文档） | **成立**：前置/双路径/四终点/命令行+完整 smoke/接业务齐备；终点 4 明确操作化为 `list-edit-lifecycle`（非 `overview`），对齐 I-008-002 §3.2 R-008-002。 |
| 独立复现记录（S3 证据） | **成立**：[R5-S3-REPRO-001](attachments/R5-S3-REPRO-001.md) 含协议 §4 最小字段；`compose` 路径；`same-operator-clean-session`（先 `down -v` 清卷、全新 `up -d`，未复用服务/DB）；四终点全 PASS；单次计时 `34.5s ≤ 900s`；浏览器终点由真实 Chromium 驱动（截图+脚本留痕）。 |
| `scripts/smoke.sh` 契约（S4） | **成立**：SM-001～SM-006 形状、退出码 `0/2/3/4/5/6/70`、`--disposable` 开关、不输出 secret、无开关不 reset，均对齐协议 §5.1/§5.2/§5.3；SM-006 经 `SMOKE_RESTART_CMD` 重启后二次断言「数量不变、同一记录仍在、不重复播种」。 |
| disposable SM-006 强制证据 | **成立**：本地 Git Bash + Docker 运行 `--disposable` → **SM-006=PASS**（[log](attachments/r5-smoke-disposable-local.txt)）；CI `container-smoke` 以 `--disposable` 调用（runner 每次隔离 project/volume，满足「CI 默认 disposable」）。非 disposable exit 0 未单独用作「种子可重复」关闭证据。 |
| 状态/进度边界 | **成立**：S3/S4 勾选（`2/5 → 4/5`）与已核实实施事实相符；`I-008-002` 维持 verified（v0.1.1）；`I-008-003` 仍 open（仅 S6）；Root R5 未勾选（`4/5`）；S5 待实施。本 self 不改动任何状态。 |

### Findings

- **无新 required**。
- **注记（recommended · 非阻断）**：本自审为 `same-operator-clean-session` 复现（协议 §3.3 允许）；建议未来由用户侧或独立执行者做一次真正隔离（不同 shell/checkout）的复现交叉核验；CI `container-smoke` 的 `--disposable` 证据需一次 GitHub Actions 运行记录留痕（当前仅本机 log）。

### 对照成功标准与信息门禁

| 项 | self 结论 |
|----|-----------|
| S3 成功标准 | **满足**：QUICKSTART 交付 + ≥1 次独立复现记录，终点 4 登录 + 后台可交互（列表加载），34.5s ≤ 15 分钟。 |
| S4 成功标准 | **满足**：`scripts/smoke.sh` 交付（healthz→login→`/me`→代表页→种子可重复），本地 smoke 全绿（SM-001～006）+ CI 接入。 |
| `I-008-002` | verified 维持；实施对照 v0.1.1 判据成立，未把协议修订混为实施证据。 |
| 关门边界 | S5 未实施；Root R5 与 VP-002 关门不因 S3/S4 自动达成。 |

### 声明

本自审仅追加审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。

### 响应 · A-010（/govern · 2026-08-03）

- **verdict 采纳**：`pass` 成立——S3/S4 实施主张与仓库事实一致且可复现（QUICKSTART + 独立复现 34.5s + smoke.sh SM-001～006 + 本地 disposable SM-006=PASS + CI 接入）；无开放 required。
- **注记处置**：same-operator 复现边界已在记录声明；建议在 S5 阶段或用户侧做一次真正隔离复现交叉核验；CI 运行记录留待 `container-smoke` 首次 GitHub Actions run 后补记。
- **仍开放**：S5（阶段审计 + Root R5 勾选/关门条件评估）待实施；`I-008-003` open（仅 S6）；Root R5 未勾选（Root 保持 `4/5`）。
- **下一步**：实施 S5——对 R5 工程化交付做阶段审计（self + 视需要 independent），评估并记录 Root R5 勾选与 Root / VP-002 关门证据口径；建议对 S3/S4 关闭证据做一次 `/audit` finding-closure 复审。
- **证据路径**：本响应节；`02-execution` 2026-08-03「实施 S3 + S4」节；`01-decision` D-006；`00-meta` S3/S4 行 + `I-008-002`；QUICKSTART.md；`scripts/smoke.sh`；附件 R5-S3-REPRO-001 / r5-repro-endpoint4.* / r5-smoke-disposable-local.txt。

## A-011 · S3/S4 实施与协议/运行证据独立交叉审计（2026-08-03）

- **source**：independent
- **auditor**：Codex（GPT-5）
- **类型 / scope**：execution-facts；交叉复核 GOAL-008 S3（QUICKSTART、R5-S3-REPRO-001、浏览器终点证据）与 S4（`scripts/smoke.sh`、本地 disposable 运行、CI 接入与当前 revision 运行证据），对照 `I-008-002` v0.1.1 §§3～6；不审 S5、S6 或 Root R5 关门。
- **verdict**：fail

### 范围与区间

- 工作区页眉：`workspace_id: workspace-002-production-admin-foundation`；canonical scope：`docs/workspace-002-production-admin-foundation/`；Root：`GOAL-001-production-admin-foundation`；目标：`GOAL-008-r5-engineering-fork`；`shared_materials_catalog: none`。本意见未读取或比较其他工作区。
- 已核对：本目标五件套与 `goal-tree.md`、`I-008-002` v0.1.1、`QUICKSTART.md`、`R5-S3-REPRO-001.md`、`r5-repro-endpoint4.mjs`/截图、`scripts/smoke.sh`、`.github/workflows/r6-basic-matrix.yml`、本地 smoke log、Root D-013/I-005 边界。
- 本轮独立运行：使用显式 Compose 项目 `audit-goal008-s34` 与新命名卷执行 `docker compose -p audit-goal008-s34 up -d --build`；Git Bash 执行 `bash scripts/smoke.sh --disposable`，SM-001～SM-006 全部 `PASS`、退出码 `0`；随后 `docker compose -p audit-goal008-s34 down -v --remove-orphans` 清理成功。该运行证明脚本在操作者已提供隔离环境时可执行，不证明脚本自身强制隔离，也不替代 CI 运行证据。
- 当前 `HEAD` 为 `f2f1acbbc126648c8f218145ec6df0ae027e06e8`；`gh run list --workflow r6-basic-matrix.yml` 最新成功运行 `30750997885` 的 `headSha` 为 `7a6eb5e6b48fd96cc2c8503c45ff2dbdefb76bf8`，早于 S4 脚本接入提交；该 run 的 jobs 没有 `container-smoke`。

### 成果（有证据）

| 审计项 | 独立结论 |
|--------|----------|
| S3 文档路径与终点 | **部分通过**：QUICKSTART 覆盖 Compose/双进程两路径；协议终点 4、QUICKSTART 与 Playwright 脚本均指向 `list-edit-lifecycle`、标题和 `Acme Console`。复现记录时间算术自洽（`03:04:45.596Z` 到 `03:05:20.141Z`，约 `34.545s`），但计时边界与记录字段存在 required 缺口（F-005/F-006）。 |
| 浏览器证据 | **通过**：`r5-repro-endpoint4.mjs` 使用真实 Chromium 登录、导航并断言标题/单元格；截图显示 `List + edit lifecycle`、`Acme Console` 与列表内容。 |
| S4 disposable 路径 | **运行通过但证据有界**：本轮显式隔离项目的 `SM-001`～`SM-006` 全 PASS；协议要求的隔离前提并未由脚本机器校验（F-008）。 |
| secret 输出边界 | **通过**：脚本将 token 保存在变量并抑制响应输出，未发现输出 password/JWT secret 的路径。 |
| 当前 CI 证据 | **不足**：工作流有 `--disposable` wiring，但最新成功 run 早于当前 S4 revision 且无 `container-smoke` job（F-009）。 |

### 对照成功标准与信息门禁

- **S3**：QUICKSTART 与真实浏览器终点内容成立；但 `R5-S3-REPRO-001` 将已构建镜像作为前置并只计 `up` 后耗时，违反协议“项目构建计入计时”的冻结口径，且未满足完整 source/checks/result 记录字段，不能无条件证明 `<=15` 分钟（F-005/F-006）。
- **S4**：本轮隔离运行的 disposable SM-006 通过；但是非 disposable 分支仍以 `SM-006=SKIP` 输出 `SMOKE RESULT: PASS` 并退出 0，和 v0.1.1 §5.3 的完整绿退出码定义不一致（F-007）；`--disposable`/`eval` 不验证隔离或普通开发库安全边界（F-008）；当前 revision 没有对应 CI run（F-009）。
- **`I-008-002`**：协议文件本身仍为 `verified`，本意见不改其状态；本意见判定的是实施/验收证据未满足协议，不把协议冻结当作 S3/S4 通过。
- **状态边界**：不修改目标 `status: active`、`progress: 4/5`、S3/S4 勾选、Root R5 或 `goal-tree.md`。

### Findings

#### F-005 · S3 计时记录排除了协议要求计入的项目构建

- **级别 / 严重度**：required / high
- **状态**：open
- **关联门禁**：S3 `<=900s` 体验验收与 `I-008-002` §3.1/§3.2。
- **证据**：协议要求项目自身编译、配置、启动、迁移、种子和登录不得预先完成，且“其余按文档完成环境配置、构建和启动的耗时计入”（`I-008-002-fork-reproduction-protocol.md:53-66`）。记录却写明 Compose 镜像已在先前 S2/本次 smoke 构建中完成，本轮 `up -d` 复用缓存且未重新 build（`R5-S3-REPRO-001.md:55-58`）；QUICKSTART 的 Compose 首条启动命令是 `docker compose up -d --build`（`QUICKSTART.md:32-39`）。
- **风险**：`34.5s <= 900s` 只证明预构建镜像后的启动/页面时间，不能证明按 fork 文档从依赖就绪到构建、启动、登录和页面加载的完整体验时限。
- **要求**：重新生成一份 clean-ref 记录，依赖下载/镜像 pull 可排除，但项目 build、配置、迁移、种子、登录和页面加载必须从计时起点开始；或由 `/govern` 按 P-004 留痕有界 residual（不等同于 verified）。

#### F-006 · S3 复现记录未满足协议规定的可核验字段

- **级别 / 严重度**：required / medium
- **状态**：open
- **关联门禁**：S3 独立复现证据完整性与 `I-008-002` §4。
- **证据**：协议要求 source 列出未提交 diff 摘要/哈希、checks 逐项包含 smoke 输出、result 提供失败原因/重试编号/日志、截图和命令输出路径（`I-008-002-fork-reproduction-protocol.md:93-104`）。记录只写了 3 个未提交新增但没有 diff hash（`R5-S3-REPRO-001.md:29-32`），checks 只有四个终点、没有 smoke 输出或“不适用”声明（`:72-79`），result 只有截图/脚本路径，没有命令输出或运行日志路径（`:85-89`）；“脚本已提交工作树”与同段“3 个未提交新增”也有语义歧义（`:22-31`）。
- **风险**：第三方无法从该记录重建确切 ref/diff、核对命令输出或区分未运行项，独立复现不能作为完整可追溯记录。
- **要求**：补齐 clean ref、工作树和 diff hash、完整命令/日志输出路径、smoke 各项（或明确标注 S4 不适用）及失败/重试记录；保持 secret 脱敏。

#### F-007 · 非 disposable smoke 仍报告 `PASS`/exit 0，违反 v0.1.1 退出码契约

- **级别 / 严重度**：required / medium
- **状态**：open
- **关联门禁**：S4 机器可判定验收与 `I-008-002` §5.1/§5.3。
- **证据**：协议定义 exit `0` 为 SM-001～SM-005 通过且 disposable SM-006 通过，并明确非 disposable exit 0 不得单独关闭种子可重复性（`I-008-002-fork-reproduction-protocol.md:112,123,129-139`）。脚本非 disposable 分支输出 `SM-006=SKIP`、`SMOKE RESULT: PASS`，随后无条件 `exit 0`（`scripts/smoke.sh:190-199`）。
- **风险**：调用方只看退出码或 `PASS` 摘要时，会把只完成 SM-001～SM-005 的部分检查误判成 S4 完整绿。
- **要求**：使非 disposable 路径与 v0.1.1 的完整绿语义一致（例如返回明确的 partial/non-success 结果并让 CI 不把它当完整绿），同时保留默认路径的非破坏性检查。

#### F-008 · `--disposable` 没有机器可判定的隔离与 destructive 安全守卫

- **级别 / 严重度**：required / high
- **状态**：open
- **关联门禁**：S4 disposable/隔离证据与 `I-008-002` §5.1/§5.2。
- **证据**：协议要求 disposable 运行只能在隔离环境验证空库种子，必须拒绝普通开发库 reset，CI 使用隔离 Compose project/volume 或等价临时 DB（`I-008-002-fork-reproduction-protocol.md:110-112,123,127-132`）。脚本只把 `--disposable` 映射为布尔开关（`scripts/smoke.sh:20-32,40`），接受任意 `SMOKE_RESTART_CMD` 并用 `eval` 执行（`:31,170-173`），不验证 Compose project、volume、DB_PATH、临时数据库或 disposable 标记；CI 也用默认 `docker compose up -d`，没有显式隔离 project/volume（`.github/workflows/r6-basic-matrix.yml:76-95`）。本轮能通过显式 `-p audit-goal008-s34` 新卷运行，证明的是外部隔离前提有效，不是脚本具备安全守卫。
- **风险**：操作者可在普通开发数据库上带 `--disposable` 并注入 destructive restart 命令，或在非隔离 CI 上把种子断言当作 S4 证据；这违反协议的 fail-closed 安全边界。
- **要求**：限制 restart/reset 命令并校验隔离项目、临时卷/DB 或等价标记；不满足时返回安全前提失败（exit 2），并在 CI/日志中记录隔离身份。

#### F-009 · 当前 S4 revision 没有对应的 CI smoke 运行证据

- **级别 / 严重度**：required / medium
- **状态**：open
- **关联门禁**：S4 “本地与 CI smoke 全绿”成功标准与协议 §6 的运行/CI 证据要求。
- **证据**：当前 `HEAD` 为 `f2f1acbbc126648c8f218145ec6df0ae027e06e8`；`gh run list --workflow r6-basic-matrix.yml` 最新成功 run `30750997885` 的 `headSha` 为 `7a6eb5e6b48fd96cc2c8503c45ff2dbdefb76bf8`，早于 S4 smoke 接入提交，`gh run view` jobs 仅有 browser E2E、api、web，没有 `container-smoke`。仓库工作流 wiring 存在（`.github/workflows/r6-basic-matrix.yml:76-95`），但 wiring 不是 CI smoke 全绿的运行事实；已提交的本地 log 只有 SM-001～SM-006 输出，没有 run/ref/provenance（`attachments/r5-smoke-disposable-local.txt:1-8`）。
- **风险**：无法证明当前 revision 的 CI runner 真的执行了 disposable SM-006、退出码传播和 teardown；S4 的 CI 部分仍是计划/接线证据而非验收事实。
- **要求**：对当前 revision 触发并保留 `container-smoke` 成功 run/job URL、head SHA、SM-001～SM-006 输出和 teardown 结果；若 CI 无法运行，按 `/govern` 记录有界 residual，不得写成“CI 全绿”。

### 推荐项（非阻断）

- **R-011 · 重启后的 readiness 未重新判定**：脚本复用初始 `ready=1`，重启循环结束后不检查 readiness 结果（`scripts/smoke.sh:174-178`）；建议重置并在失败时返回协议定义的 readiness 失败码，避免把后续登录失败误归类为其它检查项。

### 必改项汇总

- **F-005 required/high**：S3 计时须纳入项目 build。
- **F-006 required/medium**：补齐 S3 记录的 diff hash、命令/日志输出、smoke 字段和 result 追溯。
- **F-007 required/medium**：修正非 disposable `PASS`/exit 0 与 v0.1.1 退出码语义的冲突。
- **F-008 required/high**：实现 disposable/隔离与 destructive 安全的机器守卫，并在 CI/记录中固定隔离身份。
- **F-009 required/medium**：补当前 revision 的 CI `container-smoke` 运行证据。

### 与既有意见的异同

- A-010 是同一 S3/S4 execution-facts scope 的 `self · pass`，本 A-011 的 `fail` 与其在“可无条件判定 S3/S4 通过”上构成 P-004 §4.2 的 verdict/门禁冲突；A-010 中的浏览器与本地 disposable 主体事实经本轮点验/复跑基本成立，但本意见识别了其未覆盖的计时、字段、退出码、安全和 CI 证据缺口。
- A-008/A-009 仅审 `I-008-002` design-plan/finding-closure，不与本 execution-facts scope 冲突；`I-008-002` 的 verified 机读状态不因本意见自动改写。
- 本意见不否定 S3 浏览器脚本/截图或本轮显式隔离 disposable 运行；它们不足以单独闭合上述 required findings。

### 结论 + 建议给编排器/用户的下一步

- **fail**：S3 浏览器终点和本轮显式隔离的 S4 disposable smoke 可复现，但当前记录和实现不满足完整的协议/验收边界；F-005～F-009 均为未闭合 required。不能仅凭 A-010 self `pass` 将 S3/S4 作为无条件验收完成。
- 建议使用 `/govern`：展示 A-010/A-011 冲突并按 P-004 等用户裁决；建议优先修复 F-005/F-007/F-008，刷新 S3 record 与当前 revision CI 证据，再做同 scope self/independent finding-closure 复审。未完成前不要推进 S5、勾选 Root R5 或关门。

### 声明

本意见仅追加 `source: independent` 审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文、信息项状态或 `goal-tree.md`；finding 响应与后续推进由 `/govern` 处理。
