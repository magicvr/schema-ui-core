---
title: D-002 · W22 收尾 P-004 裁决：B1 续期 / B3 续期 / W13 追认
status: accepted
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-033-w22-residual-closeout
version: 0.1.0
---

# D-002 · 收尾 P-004 书面裁决（2026-08-23，ask_user_question）

## 1. B1 append 半边 —— 用户书面接受续期

R4-I004 operationlog 的 **append best-effort 半边维持 accepted-residual**：

- 范围：`recordAudit` fire-and-forget（业务写入与 operationlog 不保证一致）；TransactionalRecorder 已为敏感路径提供可选 fail-closed，但不全局强制。
- 缓解：下载/审计面既有控制 + retention 已兑现（W12-GOAL-008）+ 敏感路径可选用事务内记录。
- owner：`magicvr`
- 新复审日期：**2027-02-01**
- 复审触发：合规要求零日志丢失 / 引入幂等重试机制 / 安全审计发现审计缺口。
- 留痕位置：[W3·GOAL-005 03-audit](../../../../workspace-003-modular-admin-architecture/GOAL-005-r4-full-module-migration/03-audit.md) R4-I004 行；[W3·GOAL-006 01-decision](../../../../workspace-003-modular-admin-architecture/GOAL-006-r4-c1-freeze-decision/01-decision.md) C1-I003 行。

## 2. B3 allowlist 深化 —— 用户书面接受续期

C4-004 PolicyID/Visibility 校验器深化维持 residual：

- 范围：保持 v1 最小 dotted-identifier 语法；不实现 allowlist 与表达式语法（表达式由测试明确拒绝）。
- owner：`magicvr`
- 新复审日期：**2026-12-01**
- 复审触发：R6 清理完成后首个需要多条件 Visibility 表达式的业务模块 / 安全审计要求策略表达能力。
- 留痕位置：[W3·GOAL-010 A-003](../../../../workspace-003-modular-admin-architecture/GOAL-010-r4-c4-schema-other-migration/03-audit/A-003-r4-c4-schema-migration-response.md) C4-004 行；W3·GOAL-012 `02-execution/E-004-r5-c51-residuals-and-closeout.md` 对应行。

## 3. W13 SQLite→PG 搬运器 residual —— 用户书面追认

用户追认 [W13·GOAL-006-r5 D-002](../../../../workspace-013-store-dialects/GOAL-006-r5-dual-path-acceptance/01-decision/D-002-upgrade-backup-contract.md) 的有界 residual 条款：「本 VP 不提供自动化 SQLite→PG 数据搬运器；存量路径 = fresh bootstrap + 运维自备导出/回放」。自此该 residual 构成 P-003 意义上的合法 `accepted-residual`（范围=VP-013 退出判据 2；复审触发=出现 in-place 或内置搬运需求时另立目标）。追认节已追加至该 D-002 文件。

## 关联

本裁决覆盖 E-004/E-005 中标记「待用户书面确认」的全部开放项；无其他待裁冲突。
