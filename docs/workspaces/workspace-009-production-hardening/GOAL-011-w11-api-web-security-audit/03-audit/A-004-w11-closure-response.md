---
id: A-004-w11-closure-response
doc: audit-entry
goal: GOAL-011-w11-api-web-security-audit
title: W11 闭合记录与意见响应（A-001 + A-002 + A-003 合并）
source: self
auditor: 编排器（govern · 本会话）
date: 2026-08-22
scope: 全部相关意见（A-001/A-002/A-003）的 finding 闭合；required 三路径闭合记录；A-003 recommended/informational 响应；关门条件核对
verdict: pass
status: recorded
parent: GOAL-011-w11-api-web-security-audit
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# A-004 · W11 闭合记录与意见响应（2026-08-22）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | self（编排器响应记录；独立意见 A-001/A-003 各自的 source 不变） |
| **类型** | closure-response / finding-closure |
| **scope** | 相关意见 = A-001（independent · fail）、A-002（self · pass）、A-003（independent · pass）；所需关闭的 required findings = A-001 F-001～F-006 |
| **verdict** | **pass**（合并响应；代码闭合条件已由 A-003 独立确认） |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 本区根；`primary_plan` = `VP-009-production-hardening`） |

## 意见台账与响应汇总

| A-ID | source | verdict | 响应 |
|------|--------|---------|------|
| A-001 | independent | fail（6 required + 13 recommended + 6 informational） | 整单采纳 required（D-002）+ recommended 处置（E-003） |
| A-002 | self | pass | 与 A-003 合并；残余移交见下 |
| A-003 | independent（grok-build · 真实 PG 复跑） | **pass**，开放 required = 0 | 全部采纳；R-001/R-002 为 non-blocking 残余，移交后续（见下） |

## Required closing（三路径：fixed）

| F-ID | 闭合路径 | 证据 |
|------|----------|------|
| F-001 | fixed | E-002 + A-003（PG 真实创建/重名） |
| F-002 | fixed | E-002 + A-003（同事务回滚双用例） |
| F-003 | fixed | E-002 + A-003（限流 + 懒清理 + 原子计数） |
| F-004 | fixed | E-002 + A-003（旋转窗口双钥 + 重封） |
| F-005 | fixed | E-002 + A-003（PG 并发 1 胜） |
| F-006 | fixed | E-002 + A-003（400 + wallet.write） |

开放 required = **0**（fixed ×6；无 accepted-residual / user-overruled 用于 required）。

## Recommended closing

| F-ID | 处置 | 依据 |
|------|------|------|
| F-007 / 008 / 010 / 011 / 012 / 013 / 014 / 016 / 017 / 018 | fixed | E-003 代码证据 + A-003 有据核对 |
| F-009 | fixed（unknown-handler 部分）+ residual（lastRun 单实例文档化） | E-003 + A-003 I-A |
| F-015 | user-overruled（有据） | E-003：客户端跨标签重试协议依赖「重放无连带」；家族吊销即双标签互踢。复审触发：客户端协议变更或威胁模型升级 |
| F-019 | user-overruled（有据） | E-003：custom action 白名单映射到已鉴权 API；schema 无 permission target；硬门禁在服务端。A-001 原文已注 |

## A-003 recommended 残余响应（non-blocking，移交后续波次）

- **R-001**：`/api/auth/mfa/verify` 无独立 HTTP 限流——proof 签发已按登录桶限流（约 100 次 TOTP 尝试/15min），6 位 TOTP 不可实用穷举；`fail_count < 5` SQL 守卫封顶并发窗口。**接受 residual**：复审触发 = TOTP 位数缩短 / captcha 关闭且第二因子成为唯一远程秘密 / verify 层 DoS 威胁模型升级。
- **R-002**：wallet 写口测试用无 wallet 键的 editor 证 403——生产代码三写口为 `wallet.write`（A-003 源码核对确认）；测试缝不构成缺陷。**接受 residual**：下波若建 wallet.read-only 角色夹具可补强。
  - **2026-08-23 兑现复核：fixed**——verify 独立限流桶（15m/10/IP，429+Retry-After）落地，3 条测试 PASS；independent 复核 pass。[workspace-010 GOAL-033 E-005](../../../../workspace-010-design-implementation-conformance/GOAL-033-w22-residual-closeout/02-execution/E-005-s2-completion-facts.md)
- A-003 I-B（重封无 CAS）/ I-C（entries/tasks 无独立 HTTP 回滚用例）：记录在案，非阻断；I-B 窗口极窄且失败不阻断已通过的第二因子。
- A-003 I-D：E-003 在 A-003 之前处置 recommended 的顺序偏差——用户轮次指令「推进…直到顺利闭门」隐含全波闭环，E-003 内容经 A-003 独立核对全部有据，偏差不构成缺陷，已记录。

## Informational（A-001 F-020～F-025）

保持记录项（产品选择 F-020、文档缺口 F-021/F-024、死代码路径 F-023、低危 F-022/F-025）；本波不升格、不实施，A-003 I-F 确认。F-025（roles JSON 与 user_roles 不一致时登录 500）在 A-001 中即为 informational，若后续波次处理须先由用户裁决范围。

## 信息台账关闭

| I-ID | 状态 | 依据 |
|------|------|------|
| I-001 | verified | A-001 清单 + E-002/E-003 逐条回应 |
| I-002 | verified | D-002：整单 6 条 + 波内暂挂 VP-008 go；恢复条件已满足（A-003 pass + 回归全绿）→ [D-004](01-decision/D-004-w11-go-restore.md) |
| I-003 | **verified** | 工作区惯例 grok independent 腿 = A-003（grok-4.6 · reasoning high · 真实 PG 复跑）已满足；用户轮次指令书面关闭 |

## 关门条件核对

- [x] 相关意见无未合法闭合的 required（fixed ×6）
- [x] 相关信息项无到期开放 required（I-001/I-002/I-003 verified）
- [x] 至少一次阶段/关门向审计：A-002（self pass）+ A-003（independent pass）
- [x] 成功标准对照：S1–S4 检查点全勾（见 00-meta）
- [x] 回归证据：go vet 0 / go test 全绿（含真实 PG 2 例）/ web 1085/1085 / tsc 0；checkpoint `72a5397`
- [x] 恢复 VP-008 go 消费有效性宣称：[D-004](01-decision/D-004-w11-go-restore.md)

## 残余移交（记录在案，非本波范围）

1. 数据库密码轮换（用户侧残余，继承 W10 移交；不可因本波闭合关闭）。
2. F-009 `lastRun` DB 持久化（多副本调度启用时；A-002 R-001 / A-003 I-A）。
3. A-003 R-001 / R-002（触发条件见上）。
4. A-001 informational F-021/F-024/F-025（文档与数据卫生，另行裁决）。

## 结论

本波 required 全部按 A-001 原文处方修复并经 self + independent 双复核确认 genuine fixed；开放 required = **0**；关门条件全部满足。**目标状态变更与 VP-008 go 宣称恢复由 [D-004](01-decision/D-004-w11-go-restore.md) 落盘。**
