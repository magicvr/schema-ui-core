---
id: GOAL-004-w3-schema-host-protocol-conformance
title: W3 · Schema-UI 语义对齐与 Host/App 协议增补
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-12
updated: 2026-08-13
version: 0.2.0
progress: 2/6
---

# GOAL-004 · W3 · Schema-UI 语义对齐与 Host/App 协议增补

## 概述

本子目标是 VP-010 / workspace-010 的**第三波**：审视 `apps/api` 与 `apps/web` 的 UI 语义来源，把发现分为两类并分别闭环：

1. **实现偏离**：现行 Schema-UI 协议已经定义语义，但 API/Web 仍使用私有字段、重复语义源、宽松校验或硬编码业务 UI；在新协议基线确定后修正实现。
2. **协议缺口**：业务确实需要、但现行协议没有覆盖的 Host/App 契约；先回到协议源补齐 auth、bootstrap、branding、shell、全局 error 等契约，再消费新协议。

本波不把所有手写 Host UI 都判为绕过。登录壳、启动壳、全局导航和错误边界可以由 Host 实现，但其数据形状、状态机、能力协商、安全边界和扩展点必须有协议依据。页面业务语义若已有 Schema-UI 表达能力，则不得另造本地 DSL。

## 当前边界

- 当前固定协议来源为 `schema-ui-docs`，Web provenance 显示 `artifactVersion: 2.7.0`。
- Renderer 对未知节点 fail closed；`form`、`table`、`recordView`、`actionButton` 等已进入白名单的节点不因使用 Host adapter 而自动构成绕过。
- 已发现的实现整改候选包括：Settings PATCH 未显式拒绝未知字段、Users 导航 `label` 在 provider 与 manifest 双重定义，以及后续语义 validator 盘点所得项目。
- 已发现的协议缺口包括：公共 branding 启动配置、认证/会话生命周期、bootstrap 顺序、Shell 与偏好、全局错误恢复等 Host/App 契约。
- **停止线（已解除，2026-08-13）**：上游 `schema-ui-protocol v2.8.0`（tag `v2.8.0` @ `593f625`）已发布并固定到本仓
  （`provenance-v2.8.json`，正式 content/fixture digest 绑定；E-004）。I-003 已满足，`apps/api` / `apps/web` 的
  正式问题修复不再被停止线阻断；S2 方案级 cross 审视与 S4–S6 残余门禁仍按各自定义执行。

## 成功标准与路线图（P-001）

- [x] **S1 · 基线与候选目录**：建立 Host/App 协议候选附件；区分现有覆盖、实现偏离、协议缺口与未来业务候选。（2026-08-12 · E-001）
- [ ] **S2 · 上游协议方案**：逐项给出 `adopt-now` / `reserve-extension` / `explicitly-out` 处置，形成 schema、能力声明、状态机、错误语义和 fixtures 提案；完成方案级 cross 审视。
- [x] **S3 · 新协议到手**：上游增补已合并并形成可固定的版本/commit（tag `v2.8.0` @ `593f625`）；本仓 provenance（`provenance-v2.8.json`）、schemas、registry、fixtures（正式 release 字节级 pin）已更新并通过结构/语义验证（E-004；vitest 862 + claim-artifact 门禁）。（2026-08-13）
- [ ] **S4 · 实现整改**：仅以 S3 固定的新协议为依据，修正 API/Web 的实现偏离并接入 Host/App 契约；不保留无治理的本地私有方言。
- [ ] **S5 · 符合性验证**：API/Web validator、协议 fixtures、代表性业务页面、auth/bootstrap/shell/error 流程与回归测试通过；旧协议兼容策略有证据。
- [ ] **S6 · 关门**：全部 required 信息项与 findings 合法闭合；完成实施事实审视、go 影响判定和关门审计。

`progress: 1/6` 由上述六个等权检查点中的一个完成项确定。该展示不放行 S2～S6，也不覆盖信息门禁或审计 finding。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 当前实现与 2.7.0 的逐项覆盖/偏离基线是否完整 | S2 方案 | S2 冻结前 | 对照 upstream schema/registry/fixtures 与 API/Web 语义，维护附件处置列 | collecting | — | `attachments/I-HOST-APP-001-protocol-gap-catalog.md` |
| I-002 | required | Host/App 候选哪些进入本次协议、哪些仅保留扩展点或明确排除 | S2 方案 | S2 冻结前 | 对每个候选记录 `adopt-now` / `reserve-extension` / `explicitly-out` 及理由 | verified | — | ADR-0034–0037 accepted（2026-08-13），D10 逐项处置 95/95 且 shape/state/security/fixtures 已随 2.8.0 发布齐备；本仓 S2 方案级 cross 审视仍待执行 |
| I-003 | required | 新协议是否已完成上游合并、发布和本仓固定引用 | **S4 实施** | S4 开始前 | 核对版本/commit、provenance、schema/registry/fixtures 与 capability matrix | verified | — | 上游 v2.8.0 发布（tag `v2.8.0` @ `593f625`，content `40690917…` / fixture 树 `7aacf133…`）；本仓 `provenance-v2.8.json` 固定 + release 字节级 fixture pin（E-004） |
| I-004 | required | 协议/跨边界变更的 independent provider | S2 审视、S6 关门 | 首次 cross 审视前 | 用户指定 provider；self + independent 分别落 A 台账 | verified | — | 用户指定 `grok build`（grok 4.5，reasoning high）；self=A-001（pass）；independent=A-002（conditional，BLOCKING_COUNT=0）已落盘 |
| I-005 | required | 2.7.0 消费方的兼容、迁移、弃用和 fail-closed 规则 | S3/S4 | S3 固定前 | 形成版本矩阵、迁移说明、正反 fixtures | collecting | — | 上游已交付：`migrations/2.7-to-2.8.md`（双轨兼容矩阵、零动作项）、registry `deprecatedSince`/`removedIn` 弃用机制、正反 fixtures（host 99 + app-manifest 41 + version-negotiation 2.8 向量）；本仓 S3 台账录入随审计轮收尾 |
| I-006 | required | `recordView` 行上下文、抽屉/详情交互等争议语义的最终归属 | S2/S4 | S2 冻结前 | 对照现行 registry；上游明确标准能力或明确 Host extension | verified | — | 上游已裁定：ADR-0034 D6 IMP-004（row selection → drawer/detail）保留独立 overlay ADR；`record.view.load` 复用既有；D7 明确 reserve 位置不冒充 capability |
| I-007 | non-blocking | P1/P2 业务候选的引入顺序 | 后续发布节奏 | S2 后 | 依据产品需求分批，但每项必须有协议处置 | collecting | 责任人=维护者；复核=每次协议 release proposal | 附件已给初始优先级，不等于最终承诺 |

## 审计策略

本波涉及协议、Host/App 边界、兼容性与 API/Web 跨层语义，风险模式为 `cross`。S2 方案冻结与 S6 关门至少各需要一条覆盖对应 scope 的 `self` 和一条用户指定 provider 的 `independent` 意见；未指定 provider 时保持门禁，不静默降级。

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
