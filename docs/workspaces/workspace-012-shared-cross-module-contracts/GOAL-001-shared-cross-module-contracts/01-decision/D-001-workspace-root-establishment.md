---
id: D-001
doc: decision-entry
goal: GOAL-001-shared-cross-module-contracts
status: accepted
created: 2026-08-18
updated: 2026-08-18
version: 1.0.0
---

# D-001 · 开区 scaffold 与横切契约波路线图

## 背景

用户确认：不是“全部都不做”，而是需要分析哪些现在做、哪些等业务需求。经分析后选择启动**横切契约波**（P0 四项 + P1 两项），并确认新建独立 VP-012 / workspace-012 承载，避免塞进 VP-011 的 S/B 模块编号。

## 决策

1. 新建 `VP-012-shared-cross-module-contracts`（active）与 `workspace-012-shared-cross-module-contracts`（delivery lead）。
2. Root `GOAL-001-shared-cross-module-contracts` 建立，纲领路线图 R1～R6：
   - R1 correlation / request-id / 错误恢复
   - R2 审计事件模型增强
   - R3 乐观并发 + 幂等
   - R4 异步 Job / 长操作
   - R5 maintenance / degraded / read-only
   - R6 API Token / Service Credential
3. 首个子目标 = `GOAL-002-r1-correlation-error-contract`（R1）。
4. 不承载 Tier D 业务域；安全/符合性 gap 分流 VP-009/VP-010。

## 审计模式

开区为文档 scaffold（低风险、可逆）：**none**；后续契约实现按 security/data 门禁补 self / independent。

## 未选方案

- 扩展 VP-011/workspace-011 承接：会稀释“标准 Admin 功能模块”边界，且与 D-002（workspace-011）刚写的“横切不自动进入 VP-011”冲突。
- 挂到 VP-009/VP-010：这两个是持续程序，不是契约交付波；范围不匹配。
- 只登记不实施：用户已确认要启动横切契约波。
