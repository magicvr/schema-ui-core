---
id: GOAL-003-dual-key-jwt
title: R2 JWT 双密钥实现（重叠窗验签）
status: done
parent: GOAL-001-key-rotation-and-backup
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
plan_refs:
  - VP-016-key-rotation-and-backup
primary_plan: VP-016-key-rotation-and-backup
serves_summary: 承接 Root R2：Authenticator 双密钥——签发只用 current；校验 current 失败再试 previous；refresh（opaque）不受影响。关门按 independent（grok build /audit）。
---

# GOAL-003 · R2 JWT 双密钥实现（重叠窗验签）

## 概述

承接 [workspace-016] GOAL-001 纲领阶段 **R2**，消费 R1 已冻结的配置面（Root D-002）：

- 签发永远只用 current（`AUTH_JWT_SECRET`）。
- 校验先试 current，签名不匹配再试 previous（`AUTH_JWT_SECRET_PREVIOUS`）；两者都强制过期与方法检查。
- previous 未配置 = 单密钥行为逐字节不变。

## 边界

- 不改 token 格式（不加 `kid`、不改 claims）；不改 refresh 模型（opaque SHA-256 不动）。
- 不做热加载；不做 KMS/多密钥环；不做 Admin 密钥页。
- 安全 finding → VP-009；符合性 gap → VP-010。

## 检查点（progress 来源）

| # | 检查点 | 状态 |
|---|--------|------|
| 1 | D-001 决策关闭 I-003（重叠窗语义 / kid / refresh） | done（D-001；Root I-003 verified） |
| 2 | Authenticator 双密钥实现 + composition 接线 + 单测（重叠窗 / 窗关闭 / 签发只用 current / 过期不可延长 / 单密钥不变式） | done（E-001，4/4 PASS） |
| 3 | 全仓验证（vet + go test ./...）+ E 记录 | done（E-002 v1.1：vet 0 + JWT 相关包 ok，双独立复现；整包措辞按 F-003 收窄） |
| 4 | self 审计 → independent 审计（grok build /audit）→ 意见合并响应 → goal-tree 同步 | done（A-001 self pass + A-002 independent pass；F-001/F-002/F-003 全部 fixed；03-audit 响应节） |

`progress` = 已完成检查点 / 4 = **4/4**。

## 信息就绪

I-003 由本目标 D-001 关闭（required，最晚需要阶段 = R2 接入前）。无其他新增未知。

## 父目标

- GOAL-001-key-rotation-and-backup（[Q2](../GOAL-001-key-rotation-and-backup/00-meta.md)）

## 台账布局

三台账目录 `01-decision/`、`02-execution/`、`03-audit/` 平铺追加；索引文件保留 frontmatter 与条目索引。
