---
doc_type: goal-audit
id: A-005-self-response-a004
parent: GOAL-002-r1-contract-freeze
date: 2026-09-05
status: closed
version: 1.0.0
---

# A-005 · 响应 A-004（self · response）

## A-005 · 响应 A-004 closure 复审意见（2026-09-05）

- **source**：self（编排器响应记录，非独立审）
- **模式**：response · 响应 A-004（independent · codex gpt-5.6-sol · verdict **fail** · 2 required + 1 recommended）
- **verdict**：**pass**（作为响应记录：A-004 全部 3 项 finding 按 `fixed` 处理；放行以 A-006 independent closure 复审为准）

### 前提承认

A-004 指出 A-003「全部 5 项均 fixed」的关闭声明不实——**接受该判定**：A-003 对 A-002 F-001/F-002 的修复文本存在与代码事实的矛盾（未导出 sentinel 的跨包识别、事务边界），当时不应记为 closed。本响应修正流程：以 A-006 复审通过作为唯一放行依据。

### 关闭证据表

| Finding | 级别 | 闭合 | 证据 |
|---------|------|------|------|
| A-004 F-001 · wallet ledger 唯一竞争 sentinel 对域调用方不可见 | med required | fixed | D-002 §4.4 重写错误分类：**不要求跨包识别 wallet 内部未导出 sentinel**。同 request 竞争由「attempt 起点回读（第 1 步）+ 自表 INSERT 唯一违反经 `kernel.IsUniqueViolation` 识别 + 通用有界重试」收敛：失败方下一次 attempt 回读命中已提交胜者 → 幂等重放，不触碰 wallet ledger；wallet 互异幂等键保留为纵深防御。终态错误集合显式（ErrInsufficient/币种/on_sale/subject/ErrRequestIdConflict），其余非终态错误统一按竞争重试。新增收敛论证段 |
| A-004 F-002 · Consume 重试循环与 `Store.Run` 事务边界矛盾 | med required | fixed | D-002 §5.2 重写：attempt 循环移到 `store.Run` **之外**，每次 attempt 一个全新事务；callback 内无 COMMIT/ROLLBACK，以返回值表达 committed-OK / deterministic-insufficient / 竞争失败三种结局；隔离级别声明改为「平台不固定隔离级别（PG `BeginTx(ctx, nil)` 未设 TxOptions，隔离由部署默认决定；SQLite 写者串行），算法不依赖 serializable，正确性仅依赖谓词重检 + 全新事务重读 + 终态不足的确定性判定」；补充「不足响应的线性化点」论证（部分扣减随事务回滚，无消费发生） |
| A-004 F-003 · workspace R1 审计投影未含 A-003/A-004 状态 | low recommended | fixed | `workspace.md` 纲领阶段表 R1 行更新为完整审计循环状态（A-001 → A-002 → A-003 → A-004 fail → A-005 已修订 → 待 A-006） |

### 仍开放项

- 无（A-004 全部 finding 已处理；D-002 保持 draft，待 A-006 independent closure 复审通过后转 accepted v1.0.0 并关门 C3）。

### 声明

本响应记录不修改 A-002/A-004 原文；`fixed` 证据以 D-002 现行文本与 git 历史可核对。
