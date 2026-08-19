---
id: GOAL-007-r6-api-token-service-credential
title: R6 · API Token / Service Credential
status: done
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
progress: 100
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
serves_summary: 建立与用户会话分离、具备作用域、吊销和审计的机器凭据契约，并保持现有 Profile、Manifest 与协议装配语义不变。
---

# GOAL-007 · R6 · API Token / Service Credential

## 概述

R6 交付机器凭据管理与认证基架：管理员可创建、查看元数据和吊销凭据；调用方使用一次性返回的 secret 作为 Bearer credential，并只获得冻结作用域内的权限。机器凭据拥有独立生命周期，不复用用户 JWT/refresh session 的 `Subject`、`token_version` 或轮换语义。

## 范围

- 定义 service credential 的身份、secret 格式、hash-only 持久化、过期与吊销语义。
- 提供受权限保护的管理 API，以及凭据 Bearer 认证与作用域投影。
- 记录创建、调用和吊销审计事实，复用 correlation 与敏感字段脱敏契约。
- 验证正常调用、scope deny、过期/吊销、一次性 secret、运行态门禁和用户会话非回归。

## 非目标

- 不把机器凭据伪装为用户 JWT 或 refresh session，不提供浏览器自动登录。
- 不提供外部 IdP、OAuth client credentials、密钥轮换编排或 HSM/KMS 集成。
- 不新增 Profile/module ID，不改变默认模块集、Manifest bytes/聚合算法或协议 pin。
- 不在审计、日志、错误体或公开 Manifest 中暴露 secret/hash/Authorization header。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 现状扫描、信息门禁、精确身份/生命周期/权限/审计契约与 cross 设计审计 | ✅ 已完成（A-001 self / A-002 conditional / A-003 response / A-004 independent pass / A-005 close） |
| S1 | 数据模型、hash-only repository 与管理生命周期实现 | ✅ 已完成（E-005；A-006 self pass） |
| S2 | 独立 Bearer 认证、scope enforcement、错误与 operation log 集成 | ✅ 已完成（E-005；A-006 self pass） |
| S3 | 组合黑盒/全量回归、独立关门审计与治理收口 | ✅ 已完成（A-007 conditional；A-008 response；A-009 independent pass；A-010 F-010；A-011 Root response） |

## 成功标准

1. 机器凭据与用户 JWT/refresh session 分离；secret 只在创建时返回一次，持久层只保存不可逆 hash。
2. 有效、过期和吊销状态可确定性判定；scope 只从既有 permission key 解析，越权请求 fail closed。
3. 管理 mutation 受权限与 R5 operational gate 保护；创建/使用/吊销审计含 actor/credential/correlation 且不泄露 secret。
4. 定向与全量测试通过；Profile 默认集、模块矩阵、Manifest bytes/装配语义和协议 pin 不变。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 现有用户认证、opaque token、权限和审计的可复用边界 | S0 方案 | S0 结束前 | 扫描 auth/authsession/resources/operationlog/composition | verified | 2026-08-19 | E-001：用户 JWT 仅含 user subject/token_version；refresh 已有 256-bit opaque/hash-only 模式；权限为 persisted permission keys；operation log 可承载 correlation |
| I-002 | required | 机器 principal 形状、scope 权限上限与用户会话隔离方式 | S1/S2 实施 | S0 结束前 | D-003 修订 principal/context/permission ceiling 与 deny 语义并 A-004 cross 审计 | verified | 2026-08-19 A-004 pass；实施仍须按 D-003 测试 | D-003 §§3–4、§6；A-004 F-004/F-006 fixed |
| I-003 | required | secret 前缀/熵/hash、过期、吊销、一次性展示与并发语义 | S1 实施 | S0 结束前 | D-003 冻结 0044/0045 生命周期和数据约束；A-004 closure | verified | 2026-08-19 A-004 pass；运行证据留在 S1/S3 | D-003 §§1–3、§7–8；A-004 F-001/F-002/F-003/F-007 fixed |
| I-004 | required | 管理 API、权限键、错误码、审计事件与 operational gate 组合 | S1/S2 实施 | S0 结束前 | D-003 路由/错误/审计矩阵；A-004 closure | verified | 2026-08-19 A-004 pass；运行证据留在 S2/S3 | D-003 §§1、§4–6；A-004 F-001/F-004/F-005 fixed |
| I-005 | required | Profile/Manifest/protocol/readiness 不变边界 | S3 关门 | S0 结束前 | 核对 profile/manifest/Host 与 VP-012 边界；全量回归 | verified | 2026-08-19；A-007/A-009/A-010 关门闭合 | E-005/A-006：Profile/Manifest/kernel 定向与 API 全量通过；Web build 成功且生成 claim 已恢复，协议资产无交付 diff；A-007 非 finding 主张成立，整改提交不触及装配/协议资产 |
| I-006 | required | 审计模式与 independent provider | S0/S3 审计 | S0 结束前 | 按 security/data/migration/cross-module 风险分级 | verified | 2026-08-19 | 模式 cross；self + 项目级 grok-build（grok-4.6 reasoning high）independent |

## 父目标

- `GOAL-001-shared-cross-module-contracts`（Root；依赖已关闭 R2 审计模型）

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
