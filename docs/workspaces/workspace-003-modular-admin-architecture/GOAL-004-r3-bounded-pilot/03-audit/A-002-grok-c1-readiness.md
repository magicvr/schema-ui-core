---
id: A-002-grok-c1-readiness
doc: audit
goal: GOAL-004-r3-bounded-pilot
source: independent
auditor: Grok Build / grok-4.5
reasoning_effort: high
date: 2026-08-05
scope: R3 C1 / I-006 readiness
verdict: conditional
---

# A-002 · Grok Build R3 C1 独立审计

## 审计边界

本意见由本地 Grok Build CLI 以 `grok-4.5`、`high`、只读计划模式返回，
审查 workspace-003 canonical `GOAL-004-r3-bounded-pilot` 与当前代码快照。
范围为静态/嵌入式 Manifest、集中式 handler 注册、Profile 启停、Schema
所有权、Web Host/Shell 兼容与告警、回滚/数据保留，以及 C1/R3 关闭所需证据。
本意见不改 status/progress，不构成第三方鉴证，也不把源码或本地测试升级为
clean revision、CI、容器、生产或发布验收。

## 结论

`conditional`。C1 的源码盘点和 D-003 候选边界可追踪且与当前快照基本一致，
但 I-006 三项仍为 `collecting`，不能关闭 C1、冻结 R3 方案、推进 Root progress
或建立 R4 实施子目标。

## 已核对事实

- API 组合根按 Profile 生成 Manifest：`apps/api/internal/composition/composition.go:87-98`；公开端点在 `apps/api/internal/handler/manifest.go:16-42`。
- Web Docker 构建删除 `dist` 静态 Manifest：`apps/web/Dockerfile:17-18`；Nginx 精确代理 API：`apps/web/nginx.conf:10-21`；Vite 开发代理：`apps/web/vite.config.ts:20-28`。
- Web 静态 Manifest 仍作为源码 fixture 存在：`apps/web/public/.well-known/schema-ui/app-manifest.json:1`；API schema 测试仍依赖该 fixture 与全局 embed：`apps/api/internal/handler/schema_test.go:140-170`。
- `handler.Register` 仍无条件挂载试点及其他业务路由：`apps/api/internal/handler/health.go:27-37`；组合根调用时没有把 Profile plan 传入该业务注册边界：`apps/api/internal/composition/composition.go:87-89`。
- Profile 定义 MVP 不含 `admin.settings`/`admin.activity`、Admin 含两者：`apps/api/internal/kernel/profile.go:24-47,91-102`，但路由与 Schema 暴露未按同一计划过滤。
- Schema 仍由全局 fixtures 提供：`apps/api/internal/handler/schema.go:15-16,24-52`；Settings/Activity fixture 仍在 `apps/api/internal/handler/fixtures/schema/`。
- Settings branding 仍由 Host 特例触发：`apps/web/src/main.tsx:29-52`，由 App/LoginPage 消费：`apps/web/src/app/App.tsx:391-412`、`apps/web/src/app/LoginPage.tsx:35-50`；未见独立事件自动化测试。
- GOAL-004 的 D-003、A-001 与 Root I-006 台账诚实记录了政策和缺口，但未提供运行/演练闭环。

## Required findings

### F-IND-001 · I-006 不能只凭盘点和草案政策关闭

- impact: C1 close、R3 方案冻结、Root I-006 verified、R4 子目标创建
- status: open
- finding: 当前已有入口清单和候选保留/移除、兼容、回滚政策，但尚未验证删除清单、告警/兼容行为和回滚触发条件。
- closure: 提供具体关闭证据并由编排器响应；在此之前保持 I-006 open/collecting。

### F-IND-002 · Profile 启停没有贯穿 HTTP 注册边界

- impact: C1 禁用语义冻结、C2 A.2/B 中央 Register 病灶、V-2
- status: open
- finding: Manifest 可按 plan 过滤，但 `handler.Register` 仍始终挂载 Settings/Activity/Operations 等路由；当前 always-on 只能标为 gap，不能当作模块启停行为。
- closure: 发布 MVP/Admin（含禁用 Settings/Activity）的 route/nav/schema/Manifest 矩阵，并在 C2 实现后提供真实 HTTP 和重启验证。

### F-IND-003 · 兼容/告警契约只有政策文字

- impact: R3-I006-02、C1 close、R6 静态路径移除
- status: open
- finding: D-003 写了兼容窗口和告警要求，但未找到静态 fixture 命中时的可测试告警、生产无静态兜底的构建产物或运行证明。
- closure: 增加命名告警 hook 与自动化/可复现日志证据，并检查最终镜像无静态 Manifest、API/Nginx 为唯一生产来源。

### F-IND-004 · 回滚和数据保留边界没有演练证据

- impact: R3-I006-03、C1 close、R6 旧路径移除
- status: open
- finding: D-003 只规定触发器与保留 SQLite、operation_log、用户字段；没有失败演练、计数/字段核对或恢复后 readyz/Manifest/路由链证据。
- closure: 记录至少一次触发、回滚、数据保留和恢复后验证的演练；不得用政策文本替代。

### F-IND-005 · 信息门禁的证据时序未明确

- impact: C1 是否可关闭、C2 是否可开始
- status: open
- finding: 子目标信息表要求 C1 关闭前提供 I-006 证据，但执行记录将告警、禁用和回滚演练推迟至 C2/C3；不能静默把两者视为一致。
- closure: 在后续决策记录中显式处理该时序，并保持严格门禁；本轮不接受未获用户书面确认的 residual 或有界实验替代。

## Recommended findings

### F-IND-006 · Schema 所有权和 Host branding 应保持为 C2 病灶

- level: recommended / non-blocking for C1 wording
- impact: C2、V-3、V-4
- finding: 当前全局 Schema fixture 与 Host branding 特例盘点正确，但不得改写为已经模块化。

### F-IND-007 · 证据行号和 revision identity

- level: recommended / non-blocking
- finding: 修改后重新固定 `file:line`；使用 commit SHA 标注快照。Grok 返回中的工作树身份描述不作为本地当前 revision 证据，当前事实以本仓库实际 Git 校验为准。

## 与 A-001 的关系

A-002 与 A-001 的 `conditional` 同向，不构成结论冲突。F-C1-001/F-C1-002
分别由 F-IND-001/F-IND-003/F-IND-004 加强；A-002 另增 F-IND-002 的 Profile
路由矩阵要求和 F-IND-005 的证据时序要求。五项 required finding 均保持 open，
没有被本次审计自动关闭。

## 关闭所需证据包

1. R3-I006-01：MVP/Admin route/nav/schema/Manifest keep-remove 矩阵，含静态
   Manifest、API embed、Vite/Nginx/Docker、中心 Register、Host hooks、
   Settings/Activity schema fixture 的绝对行号。
2. R3-I006-02：兼容结束条件、静态命中告警 hook 与自动化/日志证据，以及最终
   构建产物无静态 Manifest、精确 API 代理仍为唯一生产路径的验证。
3. R3-I006-03：触发一次 D-003 条件，回退应用构建/Profile，保留 DB、
   operation_log、用户字段，核对恢复后的 `readyz`、Manifest、启停路由和日志读写。
4. P-004：在决策台账中显式写出 F-IND-005 的处理，不把继续推进理解成 silent
   residual；之后由 self 与 independent 复审，再作 C1 close-out。

## 审计边界声明

本意见只给出 `source: independent` 审计意见；GOAL-004 的 status/progress、
Root status/progress、goal-tree 和 finding closure 仍由 `/govern` 按 P-003/P-004
响应。R3 仍处于 C1 未完成状态。
