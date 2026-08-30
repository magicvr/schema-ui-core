---
id: GOAL-016-w15-api-web-audit-remediation
title: W15 api/web 独立审计问题修正准备
status: draft
parent: GOAL-001-production-hardening
created: 2026-08-30
updated: 2026-08-30
version: 0.1.0
progress: 1/6
---

# GOAL-016 · W15 api/web 独立审计问题修正准备

## 意图

承接 2026-08-30 对 `apps/api` 与 `apps/web` 的独立代码审计。本目标先建立修正波次的治理上下文、范围和可追踪路线图；**当前只完成立项记录，不推进方案冻结、代码实施、回归验收或关门**。

本目标挂在 `workspace-009-production-hardening` 的长期安全程序 Root 下，不改变 VP-009、Root 或 VP-008 的状态与 `go` 语义。

## 审计发现分母（待修正）

正式意见见 [A-001](03-audit/A-001-w15-independent-intake.md)。

| 编号 | 级别 | 主题 | 处理意向 |
|------|------|------|----------|
| F-001 | P1（条件性） | `schema-ui serve` 默认开发配置、已知 `admin/admin` 与固定 JWT 密钥可被网络访问 | required |
| F-002 | P1（条件性） | 公共 `server` 配置链对非 development JWT secret 缺少强度门禁 | required |
| F-003 | P1 | 生产 bootstrap 初始管理员密码只检查非空，未执行密码策略 | required |
| F-004 | P2 | MFA 自助 disable / recovery rotate 成功 TOTP 未推进 `LastUsedStep`，可在时间窗内重放 | required |
| F-005 | P2 | 邀请 token 保留在 URL 与历史记录，成功后未 `replaceState` 清理 | required |
| F-006 | P1（质量门禁） | Web Vitest 多处硬编码旧的 `api/internal/modules` fixture 根，导致 13 个 suite / 76 个测试失败 | required |
| F-007 | P3（主机条件） | Local object store 使用 `0755` 目录和 `0644` 文件，可能被同机其他 OS 账号读取 | recommended，纳入本波评估 |
| R-001 | 设计残余 | refresh token 使用 `localStorage` 是既有书面接受的 XSS 权衡 | 不在本波重开 |

## 成功标准（路线图检查点，progress 依此派生）

- [x] S1 立项与审计意见落盘：创建本目标、记录 A-001、登记范围与暂不推进约束
- [ ] S2 配置/部署范围与修正方案冻结：确认 `schema-ui serve` 支持边界、密钥/密码策略、MFA 一次性语义、邀请 URL 处理与测试 fixture 权威路径
- [ ] S3 API 修正与回归：完成 F-001～F-004 的实现、单元/集成测试和配置启动负例
- [ ] S4 Web 修正与回归：完成 F-005～F-006 的实现/fixture 修复，并通过 TypeScript、构建和完整 Vitest
- [ ] S5 主机存储评估与全量验证：对 F-007 形成 fixed 或有界 residual 决策；完成 API/Web 回归与必要部署检查
- [ ] S6 分层审计与关门：self + 项目约定 independent 复核，按 P-003 闭合 required findings；用户书面授权后再决定是否 `done`

## 边界

- 只覆盖本次审计列出的 API、Web、测试 fixture 与 LocalStore 权限问题；不承载新的业务模块。
- 不修改 Charter、VP-009、Root `GOAL-001-production-hardening` 或 VP-008 `go` 状态。
- 不重开既有 `localStorage` refresh-token residual；若安全模型改变，另行走 P-004 决策。
- 当前目标保持 `draft`；S2 及以后均未开始。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | workspace-009 与 Root 当前记录的 independent provider 版本不一致（workspace 页写 `grok-4.6`，Root 旧记录写 `grok-4.5`）；本波 cross 审计采用哪一声明？ | S2 / S6 审计门禁 | S2 前 | 核对当前 provider 可用性与用户/项目约定，并在 D-002 或审计记录留痕 | open | 未延期 | 待确认 |
| I-002 | required | `schema-ui serve` 是仅本地开发入口，还是受支持的生产/局域网入口？默认监听地址与反向代理拓扑是什么？ | F-001/F-002 方案与验收 | S2 前 | 读取 CLI/部署文档并确认实际暴露边界；必要时做回环/非回环启动检查 | open | 未延期 | 待确认 |
| I-003 | required | 初始管理员密码是否必须复用现有 8–72 字节密码策略；对已有空库/升级库的兼容迁移如何处理？ | F-003 实施 | S2 前 | 对照 `authsession` policy、bootstrap migration 与部署合同，补启动负例 | open | 未延期 | 待确认 |
| I-004 | required | Web 测试 fixture 的 canonical 根是否统一为 `apps/api/modules`，缺失的 schema 是否需要补齐或删去过时断言？ | F-006 S4 门禁 | S2 前 | 盘点测试引用与实际文件，形成一次性路径/fixture 映射并在 CI 验证 | open | 未延期 | 待确认 |
| I-005 | non-blocking | Unix 主机是否存在同机不同 OS 用户读取 LocalStore 的威胁模型？ | F-007 S5 范围 | S5 前 | 结合 Docker 非 root/volume 拓扑确认；无该场景时记录有界 residual | open | 未延期 | 待确认 |

## 愿景与工作区对齐

- 工作区：`workspace-009-production-hardening`（显式、`delivery`、唯一 lead）。
- Root：`GOAL-001-production-hardening`（`parent: null`、`active`）。
- `plan_refs` / `primary_plan`：`VP-009-production-hardening`（`active`，`vision_ref` 对齐当前 Charter `schema-ui-core-admin-foundation@0.3.0`）。
- 本目标 parent 为同区 Root 的完整 ID；不跨工作区建 parent。

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger 目录与 `attachments/`。当前没有实施事实或关门审计。
