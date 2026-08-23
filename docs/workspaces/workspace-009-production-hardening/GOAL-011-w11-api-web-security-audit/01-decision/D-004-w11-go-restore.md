---
id: D-004-w11-go-restore
goal: GOAL-011-w11-api-web-security-audit
status: accepted
created: 2026-08-22
updated: 2026-08-22
parent: GOAL-011-w11-api-web-security-audit
version: 0.1.0
---

# D-004 · W11 关门 + 恢复 VP-008 go 消费有效性宣称（2026-08-22）

## 决定（用户书面指令）

用户目标轮次指令原文："推进工作区9目标11，直到顺利闭门"——授权本波完整闭门路径（对齐 W7–W10 先例：D-002 采纳修复 → S3 实施 → self/independent 双审 → 关门 + go 恢复）。

1. **VP-008 go 消费有效性宣称：恢复**。暂挂区间 = D-002 落盘（2026-08-22）至本决定。恢复依据与 W7/W8/W9/W10 D-00x 先例一致：6 条 required 全部合法闭合（fixed）、self（A-002）与 independent（A-003 · grok-build grok-4.6 · 真实 PG 复跑）双复核 pass、API/Web 回归全绿（go vet 0 / go test 全绿 / web 1085/1085 / tsc 0）。
2. **GOAL-011-w11-api-web-security-audit 关门**：S1–S4 检查点全勾（4/4），`status: active → done`，同步 goal-tree 树与状态表。
3. I-003 随本关门 closed（A-003 即工作区惯例 grok 腿，A-004 书面关闭）。
4. 残余移交（A-004 §残余移交）：不因关门失效；数据库密码轮换仍为用户侧动作。

## 未选方案

- **继续暂挂 go 宣称直至新 VP 激活前的额外 freshness review**：VP-008 §`go` 消费有效性机制独立存在（触发失效规则时门闩自动暂挂），本波双审通过即满足 D-002 恢复条件；是否在消费前再加 freshness review 属消费方 VP 的门禁流程，不由本波决定。
- **把 A-003 R-001/R-002 升格为本波必改**：grok 独立判定为 non-blocking 且给出复审触发；强制关闭会扩大回归面，违背"有界波次"语义。

## 影响

- Root `GOAL-001-production-hardening` 保持 active（长期程序容器，不随单波 done 关闭）。
- `workspace.md` / `goal-tree.md` W11 行同步为本波最终事实（A-004 I-E 同步项）。