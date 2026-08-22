---
id: A-002-w11-self
doc: audit-entry
goal: GOAL-011-w11-api-web-security-audit
title: W11 S3/S3b 实施 self 审计（required 6 + recommended 处置）
source: self
auditor: 编排器（govern · 本会话）
date: 2026-08-22
scope: E-002（F-001～F-006 required）+ E-003（F-007～F-019 recommended）实施事实与回归证据；required 合法闭合判定（S3/S4 门禁）
verdict: pass
status: recorded
parent: GOAL-011-w11-api-web-security-audit
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# A-002 · W11 S3 self 审计（2026-08-22）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | self |
| **auditor** | 编排器（govern） |
| **类型** | execution-facts / finding-closure |
| **scope** | E-002（6 required）+ E-003（recommended 处置）的代码改动、回归锁与全量回归；required 闭合三路径核对 |
| **verdict** | **pass**（self 视图；独立复核见 A-003） |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 本区根；`plan_refs`/`primary_plan` = `VP-009-production-hardening`；Charter `schema-ui-core-admin-foundation@0.2.0`；`shared_materials_catalog: none`） |

## 范围与方法

- 通读 E-002/E-003 与逐文件 diff（`72a5397`，37 files）；逐条 required 对 A-001 原文核对修复路径与边界（未偷换范围、未「假修复」）；recommended 逐条核对 fixed 的代码证据与 overruled 的依据。
- 回归复跑：`go build ./...` / `go vet ./...` / `go test ./...`（含新增回归锁；PG 门控用例在无 PG_TEST_* 环境下按设计 skip）；web `npx vitest run`（76 files / 1085 tests）、`npx tsc -b`。
- 未做：动态 exploit、live Postgres 集成、浏览器端真实交互（web 层由 vitest 覆盖）。

## 逐条闭合判定（required ×6）

| F-ID | 修复证据 | 回归锁 | 判定 |
|------|----------|--------|------|
| F-001 | `users_repository.go` `CreateUserManagement` EXISTS → `bool`（与 DeleteUser/scanUserListRow 一致） | `TestPostgresCreateUserManagementExistsScan`（PG 门控） | **fixed** |
| F-002 | `TrashTxRecorder`/`TrashTxDeleter` 双接口 + 同事务 delete+snapshot；失败整体回滚；legacy 路径字节不变 | `TestRecycleFactoryDeleteAndSnapshotSameTx` / `TestRecycleFactorySnapshotFailureRollsBackDelete`（sqlite 实库回滚断言） | **fixed** |
| F-003 | MFA 分支同桶限流（allow+record）+ proof 懒清理 + `fail_count < 5` 守卫 UPDATE | `TestCreateProofPurgesExpired` / `TestIncrementProofFailuresCapped`；既有 exhaustion 用例保持 | **fixed** |
| F-004 | `NewService(..., previousSecret)` 双密钥 HKDF + 解密回退 + 成功后重封（CAS 赢家） | `TestServiceRotationWindow`（A 注册 → B/A 窗口可登录 → 重封后仅 B 可解） | **fixed** |
| F-005 | 单语句守卫生成 DELETE + RowsAffected 判定胜者；保留任意尝试即消费 | `TestPostgresCaptchaConsumeConcurrent`（PG 门控）+ sqlite lifecycle 用例不变 | **fixed** |
| F-006 | reconcile Decode 错误 400（`io.EOF` 空 body 保持全库哨兵）；submit/cancel/retry → `wallet.write` | `TestWalletReconcileBadBodyAndWriteGate`；既有 202 流程（空 body）保持；gates 测试 editor 403 不变 | **fixed** |

## recommended 处置核对

- **fixed（11）**：F-007（dummy bcrypt 燃烧 + limiter record）、F-008（restore+mark 同事务）、F-009（未知 handler 记 failed；lastRun 单实例为文档化设计 → residual 移交）、F-010（restoreSession 重试一次）、F-011（空数字 → undefined）、F-012（routeKey 内容级依赖）、F-013（otpauth 全量百分号编码）、F-014（MaxBytesError → 413）、F-016（badgeStyle 恢复）、F-017（formulaSafe + \t\r）、F-018（avatarQuotaMu 串行化）。逐条有源码证据与回归锁（F-010：auth-client.test 新增瞬时恢复用例；F-011/F-012：既有单元/渲染测试保持）。
- **overruled（2，有据）**：F-015（客户端跨标签刷新重试 A-002 F-003 依赖「重放无连带」，家族吊销会杀死正常双标签旋转流；服务端旋转原子 + 单 token 吊销语义保留）、F-019（custom action 后端硬门禁不变；前端 executeAction 无 permission target 语义可挂；纯 UI 对齐项）。两案均给出源码依据与复审触发，符合 P-003「作废需有据」先例（W9 F-003 / W10 D-003）。
- **无 required 升格**：recommended 处置不改变 A-001 的 F-020～F-025 informational 状态（属记录项，非缺陷）。

## 回归证据（本会话复跑）

- API：`go build ./...` 0；`go vet ./...` 0；`go test ./...` 全绿（含 handler/composition/store/auth/authsession/mfa/logincaptcha/recyclebin/datadictionary/scheduledtasks/wallet 等全部包；新增 8 个回归锁；PG 门控 2 例按环境 skip）。
- Web：`npx vitest run` **76 files / 1085 tests 全绿**（基线 1084 + 新增 1）；`npx tsc -b` **0 错误**。
- Git checkpoint：`72a5397`（37 files；owned paths 显式暂存，无 `git add -A`）。

## 信息台账核对

| I-ID | 状态 | 依据 |
|------|------|------|
| I-001 | verified | A-001 清单 + E-002/E-003 逐条回应 |
| I-002 | verified | D-002（整单 6 条 + 波内暂挂 go；复核通过后恢复） |
| I-003 | open（deferred · 本次触发） | 工作区惯例 grok independent 腿 = A-003；其输出后由关门记录关闭 |

## Findings

无 required findings。记录以下非阻断项（供关闭记录与审计台账备查）：

- **R-001（residual · 移交）**：F-009 `lastRun` 仍为进程内存去重；调度器 best-effort 单实例为文档化语义（文件头 + compose 文档）。若未来启用多副本调度，需 DB 持久化 lastRun + 租约（特性变更，留待专波）。
- **R-002（residual · 移交）**：F-015/F-019 的 overruled 记录（E-003）含复审触发：客户端跨标签重试协议变更或威胁模型升级时重审。

## 必改项汇总

开放 required = **0**（self 视图：6/6 fixed，三路径：fixed ×6）。

## 结论 + 建议

- **verdict = pass**：S3/S3b 实施事实与回归证据一致，6 条 required 已按 A-001 原意修复且未改错周边语义；recommended 处置有据。
- 建议：① 追加工作区惯例 grok independent 复核（A-003）后由编排器合并响应；② 合并通过后写闭合记录（A-004/响应节）并恢复 VP-008 go 宣称（D-004，关门）。