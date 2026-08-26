---
id: GOAL-013-w13-api-web-security-audit
title: W13 api/web 全量安全审查发现修复（P1×1 + P2×3 必修 + P3 加固全量）
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-001-production-hardening
version: 0.3.0
progress: 5/6
---

# GOAL-013 · W13 api/web 全量安全审查发现修复

## 意图

承接 2026-08-26 用户指令的 api/web 代码审查（bug + 安全漏洞）：4 路隔离上下文并行深审（认证会话/MFA、数据层/钱包/SQL、文件上传/对象存储、前端 XSS/令牌）+ 会话内核心链路独立复核。结论：无 P0、无可直接利用高危；**1×P1（公开路由 CPU DoS）+ 3×P2（MFA 猜测预言机 ×2、Confirm 步进 bug）必修**，另有 16 项 P3 加固/健壮性发现。用户裁决：**全部发现一次修完**。

审计意见权威台账：[03-audit/A-001](03-audit/A-001-w13-security-review-findings.md)（全文证据见 `attachments/audit-A-001-findings-full.md`）。范围与落位决策：[01-decision/D-001](01-decision/D-001-w13-scope-and-placement.md)。

## 路线图（progress 来源：以下 6 个检查点等权）

- [x] **S1 审计意见落盘与范围冻结** —— A-001 findings 台账（F-001～F-020 + B-1～B-4）写入 `03-audit/`；D-001 记录用户两项裁决（落位 workspace-009/W13；范围=全量）
- [x] **S2 API 必修批** —— F-001 invite-accept 先验 token + 限流；F-002/F-003 MFA 三端点限流；F-004 Confirm 匹配步进；回归全绿（checkpoint `9da0084e`）
- [x] **S3 API P3 与健壮性批** —— F-005～F-012 + B-1～B-4 全量 fixed；F-007=fixed 承载于子目标 GOAL-014、F-013=accepted-residual（用户裁决 D-002 留痕）；checkpoint `b7954235`
- [x] **S4 Web 前端批** —— F-014～F-016 fixed + 回归锁；vitest 1128/1128 + build 全绿（checkpoint `e93f7228`）
- [x] **S5 部署/运维批** —— F-017 密钥路径分置配置化；F-018 缓存窗口收敛；F-020 HSTS 实施、img-src https: 经用户裁决保留（D-002）
- [ ] **S6 审计闭合与关门** —— self 审计 pass → independent 审计（项目默认 grok build · grok-4.6 · reasoning high · `/audit`）确认全部 required 闭合 → 用户书面关门

## 边界

- 不改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 协议 pin（若实施中必须触及，按 VP-008 go 规则先暂挂并记录）。
- Root 保持 active；本波关门不推导 Root done。
- localStorage refresh token 为既有书面接受的残余（tokens.ts D-002），本波不重开。
