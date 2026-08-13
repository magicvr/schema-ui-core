---
id: A-001
goal: GOAL-007-w6-container-smoke-reproducibility
source: self
date: 2026-08-14
scope: W6 关门（F-1a/b/c + 容器构建 + V-007/V-008 复跑）
verdict: pass
parent: GOAL-007-w6-container-smoke-reproducibility
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# A-001 · W6 关门自审（self）

## 结论

**pass**：F-1（a/b/c）全部修复并以**可执行验证**闭合（claim 双路径自检、nginx -t、容器构建、V-007 exit 8、V-008 exit 0 完整绿含 SM-006 种子可重复性）。`go` 消费重验证的「冻结命令可执行性」恢复满足。

## 成果

- F-1a claim GIT_COMMIT 接线（4 文件）；F-1b nginx upstream 作用域修正；F-1c smoke.sh SM-007 按 profile 页面集断言。
- 修复均为配置/脚本层，不改变协议、Profile 默认集、模块矩阵、运行时语义。
- 残余观察（fixtures 陈旧）非阻断，见 E-002。

## Findings

| finding | 级别 | 主张 | 处置 |
|---------|------|------|------|
| — | — | 无开放 | — |

## 偏差

无。修复过程中新发现 F-1b（nginx）与 F-1c（SM-007 陈旧断言），均已并入本波范围（D-001 增补）。
