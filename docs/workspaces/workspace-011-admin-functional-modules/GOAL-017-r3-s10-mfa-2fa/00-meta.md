---
id: GOAL-017-r3-s10-mfa-2fa
title: R3-S10 · MFA / 2FA（TOTP 双因素认证）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-15
updated: 2026-08-15
version: 0.2.0
progress: 1/5
---

# GOAL-017-r3-s10-mfa-2fa · MFA / 2FA（TOTP 双因素认证）

## 概述

常用档 S-10（I-011-001 §4；R3 第三批次，2026-08-15 立项）：为登录流程叠加**第二因素**——TOTP（RFC 6238）候选、绑定/解绑与恢复码、登录挑战集成、安全存储与审计。安全敏感后台高频且增长中；基架已有失败自动锁定/限流（C-11）与改密吊销（C-10），无 2FA。

## 当前边界（立项；S1 方案冻结细化）

- 因子选型：**TOTP 优先**（不引入短信/邮件通道依赖；邮件通道依赖 B-09 模板基础设施，S1 冻结时复核）。
- 登录流程集成：既有 auth.login 语义上叠加第二因素挑战，不改变既有锁定/限流/会话吊销语义。
- 安全存储：TOTP secret 加密存储 + 权限收敛；恢复码策略（生成/吊销/轮换）。
- 审计事件：绑定/解绑/挑战成功与失败入 operationlog。
- 管理面：个人中心（自助绑定/解绑）+ 管理员视角（强制启用候选，S1 冻结）。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：因子与恢复流程（TOTP 基准 / 恢复码策略 / 解绑）、登录挑战集成点与协议面（auth.login 扩展动作键、会话面、失效语义）、安全存储、权限键与 Profile 归属；方案级 self 审视 + **grok build independent（security 门禁，grok-4.6 · high）**（D-002，2026-08-15）
- [ ] **S2 · 实现**：模块/内核扩展 + 挑战端点 + schema 页 + 测试
- [ ] **S3 · 验证**：单元/集成 + 全量回归（go 全绿 / web 全量 / e2e 双 profile）
- [ ] **S4 · go 影响判定 + 自审**
- [ ] **S5 · 关门**：独立审计（grok build）+ 关门 + goal-tree 同步

progress: 1/5 由五个等权检查点派生（S1 完成后更新）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 因子与恢复流程：TOTP 基准/漂移窗、恢复码策略（数量/一次性/吊销）、解绑流程 | S1 方案 | S1 | 业界惯例（RFC 6238）+ 既有 auth 会话面 | **verified** | — | D-002 §1（2026-08-15） |
| I-002 | required | 登录挑战集成点与协议面：auth.login 扩展动作键、会话态、失效/吊销语义 | S1 方案 | S1 | 对照 auth 模块 + protocol-inventory（登录/会话面） | **verified** | — | D-002 §3/§5（2026-08-15；AUTH-004 explicitly-out、AUTH-006 reserve-extension） |
| I-003 | non-blocking | 管理面与 Profile 归属：自助入口（个人中心）+ 管理员强制启用候选 | S1 方案 | S1 | F-03 个人中心先例 + 权限键清单 | **verified** | — | D-002 §4（2026-08-15） |
| I-004 | non-blocking | 与已关门 S-11 登录验证码的 login 链路合成裁定：先后/并存、失败计数分轨 | S1 方案 | S1 | A-002 017-F-003 登记；对照 GOAL-011 登录挑战先例 | **verified** | — | D-002 §3（2026-08-15：串行并存 + 分轨计数） |

## 审计策略

MFA/2FA 属 **security 门禁**（P-003 independent）：S1 方案冻结与 S5 关门必须 grok build independent（用户书面偏好沿用：grok-4.6 · reasoning high）。

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 01-decision/、02-execution/、03-audit/ 平铺 ledger；索引与目录条目共同构成正式记录。
