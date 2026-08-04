---
id: D-001-r4-c1-stage-scope
doc: decision-entry
goal: GOAL-006-r4-c1-freeze-decision
source: govern
date: 2026-08-05
status: accepted
---

# D-001 · 建立 R4-C1 冻结裁决子目标与路线图

## 决定

在工作区 3 中建立 `GOAL-006-r4-c1-freeze-decision`，作为
`GOAL-005-r4-full-module-migration` 的 C1 子目标，承接 Provider、Records 和
operationlog 的信息门禁、候选方案、用户裁决、最终复审和 C1 close-out。

## 理由

父目标已经形成多轮跨审计候选材料，但 C1 的 P-004 决策仍未落盘。将裁决与
最终复审独立成子目标可以保持 R4 implementation scope 与 C1 decision scope 分离，
让 C2 只接收已验证的 contract、scope 和 residual 边界。

## 边界

- 本目标不实施 C2/C3/C4 代码。
- 父目标 `GOAL-005` 仍保持 `active 0/5`；Root progress 保持 `3/6`。
- 只有完成 C1.2、C1.3、C1.4，才可向父目标传递 C2 entry context。
