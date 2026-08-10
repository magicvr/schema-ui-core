---
doc_type: vision-plan
id: VP-009-production-hardening
title: 生产加固（共享基架安全与健壮性缺陷整改）
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-009-production-hardening
created: 2026-08-10
updated: 2026-08-10
version: 0.2.0
parent: null
---

# VP-009 · 生产加固（共享基架安全与健壮性整改）

## 意图

在 VP-008 准入波次关门（2026-08-10，`go` 签发，候选 `ed99e88`）之后，对当前代码主线执行一次**共享基架安全与健壮性加固波次**：修复代码审查（2026-08-10）发现的共享基架安全/健壮性缺陷，使生产部署路径 fail-closed，并恢复 VP-008 `go` 的消费有效性。

本 VP 是 VP-008 §`go` 消费有效性规则的**新准入 VP** 分支：审查发现命中"影响共享基架或共同风险语义的问题"（VP-008 §173），按 roadmap 门闩第 9 条由 `/vision` 新建准入 VP；VP-008 的 `go` 在该等缺陷重验证前按规则暂挂。

本 VP **不**重开 VP-001～008，**不**修改 Charter 的目的、成功边界或非目标（Charter `@0.2.0` 保持）。

具体缺陷清单（2026-08-10 代码审查产物）属实现层，由工作区第一个子目标承接，不写入本 VP；输入见 `raw/audit-20260810-api-web-bug-review.md`（gitignored 临时记录）。

## 方向级范围

| 区域 | 方向级范围 |
|------|------------|
| 认证与会话安全 | 上传/下载边界的存储型 XSS 消除；refresh token 轮换原子化；生产部署 fail-closed（弱默认拒绝）；种子账户可恢复性 |
| 异步与交互健壮性 | 前端异步错误路径全覆盖（无卡死/静默失败）；查询语义（清空可清除）；权限缺省契约一致；路由 query 贯通 |
| 边界一致性与体验 | 用户/角色/设置各边界语义一致；迁移升级健壮性；表单与主题细节一致 |

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在工作区目标内）：

1. 共享基架安全/健壮性缺陷（工作区第一个子目标承接的审查发现清单）全部修复并回归；新增/更新测试覆盖，Go + vitest 全绿，基线不回归。
2. 生产部署路径 fail-closed（无公开弱默认）；认证/授权/会话边界无已知可利用缺陷。
3. 未改变 Charter 边界；未引入超出安全加固范围的新语义。
4. 共享基架重验证证据落盘，VP-008 `go` 消费有效性按规则恢复（或用户书面裁决继续暂挂）。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-009-production-hardening | GOAL-001-production-hardening | lead | — | 由 `/govern` 开区；第一个子目标 = 本批代码审查发现修正（清单见 `raw/audit-20260810-api-web-bug-review.md`） |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-10 | 初创（`planned`、0 区）；意图 = 共享基架安全与健壮性加固；VP-008 `go` 因共享基架缺陷按规则暂挂，本 VP 承担重验证 |
| 2026-08-10 | 按用户指示修订：VP 仅保留安全加固意图（决策层）；具体审查发现移出 VP，由工作区第一个子目标承接（输入 `raw/audit-20260810-api-web-bug-review.md`），不写入 vision 层 |
