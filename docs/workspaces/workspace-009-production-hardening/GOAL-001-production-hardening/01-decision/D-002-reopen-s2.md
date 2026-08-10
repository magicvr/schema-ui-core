---
id: D-002
goal: GOAL-001-production-hardening
title: 重开 Root 并增加 S2（上传所有权加固）
date: 2026-08-10
status: recorded
---

# D-002 · 重开 Root 并增加 S2

## 决策

用户于 2026-08-10 指示：在工作区 9 添加子目标承载本轮安全审视治理上下文，并修复确认问题。

据此：

1. Root `GOAL-001-production-hardening` 从 `done` **重开为 `active`**，纲领检查点由 2 项扩展为 **S0–S2（3 项）**；`progress` 2/3。
2. 新建子目标 `GOAL-003-upload-ownership-hardening`（parent = Root），承接上传 IDOR / owner 绑定 / `ReadHeaderTimeout`。
3. GOAL-002 **保持 `done`**，不重开、不改写其 16 项清单。
4. VP-009 愿景层状态不在本决策内修改；实现层先推进修复。若需 VP 重开/修订，走 `/vision`。

## 原因

S1 关门后的新一轮代码审视发现访问控制缺口（认证后文件 IDOR），属于同一生产加固工作区的延续波次，而非新业务域。
