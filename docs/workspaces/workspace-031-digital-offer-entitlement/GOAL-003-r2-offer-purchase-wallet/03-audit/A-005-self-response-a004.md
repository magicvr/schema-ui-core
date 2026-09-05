---
doc_type: goal-audit
id: A-005-self-response-a004
parent: GOAL-003-r2-offer-purchase-wallet
date: 2026-09-05
status: closed
version: 1.0.0
---

# A-005 · 响应 A-004（self · response）

## A-005 · 响应 A-004 closure 复审（2026-09-05）

- **source**：self（编排器响应记录，非独立审）
- **模式**：response · 响应 A-004（independent · codex gpt-5.6-sol · verdict **fail** · 1 open required + 2 项证据/投影缺口）
- **verdict**：**pass**（作为响应记录；放行以 A-006 independent closure 复审为准）

### 前提承认

A-004 的核心判定**成立并接受**：A-003 所称「验收矩阵双数据库执行」不实——`runPurchaseMatrix` 内部硬编码 `mustSQLite`，PG 测试虽通过但矩阵子测试实际仍落在 SQLite。这是关闭声明与测试构造的直接矛盾，属关闭声明不实，必须整改后重审。

### 关闭证据表

| Finding | 级别 | 闭合 | 证据 |
|---------|------|------|------|
| A-004（A-002 F-002 未闭合）· 矩阵未真实双库执行 | high required | fixed | `runPurchaseMatrix(t, newEnv)` 改为**环境注入**：矩阵子测试经工厂获取 env，不再硬编码方言。`TestPurchaseMatrixSQLite` 注入每场景独立 SQLite store；`TestPurchasePostgresAcceptance` 注入共享真 PG store（每子测试唯一 salt 命名空间隔离主体/请求守卫）。矩阵子测试（单事务、失败路径、幂等/冲突、并发收敛、余额竞争、币种、重试耗尽、终态不重试、故障恢复、审计 fail-closed、多字节 Q）现于两方言各执行一遍；PG 另保留 4 路同 request 并发收敛专项 |
| A-004（A-002 F-003 证据缺口）· emoji Q 与 rune 边界断言缺失 | med recommended | fixed | search 子测试新增：60 个 emoji Q（240 字节/60 rune）不产生非法 UTF-8 或错误；`🎫` 单 emoji Q 命中包含该 emoji 的 offer；100→101 rune 边界断言（101 rune Q 干净截断到 100 rune） |
| A-004（A-002 F-004 未闭合）· GOAL-003 00-meta 与 ASCII 树投影滞后 | low recommended | fixed | `00-meta.md` 的 serves_summary/概述/对齐/信息就绪四处分母引用统一为 D-002 v1.2.0；goal-tree ASCII 树 GOAL-003 进度改为 3/4；表头更新行同步 R1 合同现行为 v1.2.0 |

### 版本字段口径澄清（A-004 尾注）

`D-002-*.md` frontmatter 的 `version` 字段是**文档修订版本**（ledger 条目自身的版本）；合同语义版本在标题与修订说明（现 v1.2.0）。此口径与仓库五件套惯例一致（frontmatter version 描述文档而非被引用对象），不另行改动。

### 仍开放项

- 无（A-004 全部缺口已处理；GOAL-003 保持 active，待 A-006 independent closure 复审通过后关门）。

### 声明

本响应记录不修改 A-002/A-004 原文；`fixed` 证据以现行代码、测试与 git 历史可核对。
