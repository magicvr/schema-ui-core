---
id: A-003-w11-s4-independent
doc: audit-entry
goal: GOAL-011-w11-api-web-security-audit
title: W11 S4 关门前 independent 复核（grok-build）
source: independent
auditor: grok-build (grok-4.6 · reasoning high · `/audit`)
date: 2026-08-22
scope: A-001 required F-001～F-006 在 72a5397 是否 genuine fixed；修复是否改错既有语义；E-003 recommended 处置（fixed 11 + overruled 2）是否有据；API/Web 回归复跑
verdict: pass
status: recorded
parent: GOAL-011-w11-api-web-security-audit
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# A-003 · W11 S4 关门前 independent 复核（2026-08-22）

> 本文件为 grok-build 独立的审计意见（`source: independent`），原始输出见 [attachments/audit-A-003-w11-grok-output.txt](attachments/audit-A-003-w11-grok-output.txt)；意见内容按原样收录。不改 status/progress/goal-tree（P-003）。

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build（grok-4.6 · reasoning high · `/audit`） |
| **类型** | finding-closure / close-out（S4 复核腿） |
| **scope** | `GOAL-011-w11-api-web-security-audit`：A-001 required F-001～F-006 在 `72a5397` 是否 genuine fixed；修复是否改错既有语义；E-003 recommended 处置（fixed 11 + overruled 2）是否有据；API/Web 回归复跑 |
| **verdict** | **pass** |
| **日期** | 2026-08-22 |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 本区根；`plan_refs`/`primary_plan` = `VP-009-production-hardening`；Charter `schema-ui-core-admin-foundation@0.2.0`；`shared_materials_catalog: none`） |

## 范围与方法

- 只读：工作区 `workspace.md`、Charter/alignment/VP-009、GOAL-011 meta、D-002、E-002、E-003、A-001、A-002。未读取其他工作区过程树。
- 对照 A-001 原文逐条核 `72a5397` 源码（非只信 E-002/A-002 叙述）。
- 本会话复跑：`apps/api` `go vet ./...`；`go test ./...`（全包）；**PG 门控两例在本机真实 Postgres 上 PASS**（非 skip）；`apps/web` `npx vitest run` 76/1085；`npx tsc -b` 0。
- 未做：动态 exploit、浏览器真实交互、多副本调度演练。

P-005：I-001/I-002 verified；I-003 为本 `/audit` 腿（关闭仍归编排器）。共享资料：none。

## 逐条闭合判定（A-001 required ×6）

| F-ID | 原主张 | 代码证据 | 回归锁（本会话） | 判定 |
|------|--------|----------|------------------|------|
| **F-001** | Postgres `EXISTS` 扫进 `int`，创建/导入用户 500 | `users_repository.go` `CreateUserManagement`：`var exists bool`；与 `DeleteUser` / batch 路径一致 | `TestPostgresCreateUserManagementExistsScan` **PASS 1.43s（真实 PG）**：创建成功 + 重复用户名 `ErrUsernameTaken` | **fixed** |
| **F-002** | 删除提交后快照失败仍 204；字典类型/条目、定时任务不可恢复丢失 | `TrashTxRecorder`/`TrashTxDeleter` 双接口；`delete`/`batchDelete` 双方 opt-in 走同事务；`dictType`/`dictEntry`/`task` 实现 `DeleteTrashTx`；store `Delete*Tx` 内 `record` 失败回滚；断言失败 fail-closed。Trash 仅接这三类实体 | `TestRecycleFactoryDeleteAndSnapshotSameTx`（204 + 行消失 + `recycle_items=1`）；`TestRecycleFactorySnapshotFailureRollsBackDelete`（500 + 行仍在 + 快照 0） | **fixed** |
| **F-003** | 无限 proof + `/mfa/verify` 无限流 + `fail_count` 非原子 | login MFA 分支同桶 `allow`+`record`（15min/20）；`CreateProof` 同事务懒清过期；`IncrementProofFailures` `AND fail_count < 5`。TOTP 顺序预算约 100 次/15min，6 位 TOTP 不可实用穷举 | `TestCreateProofPurgesExpired`；`TestIncrementProofFailuresCapped`（8 次递增封顶 5）；既有 exhaustion 用例 | **fixed**（残余见 recommended R-001） |
| **F-004** | JWT 轮换后 MFA 密文只用当前 secret，用户全锁 | `NewService(..., previousSecret)` HKDF 双钥；`decryptSecret` 当前失败再试 previous；成功后 `maybeRewrap`（verify CAS 赢家 / confirm / disable / recovery）；`composition.go` 传入 `AuthJWTSecretPrevious`；previous 空 = 旧单钥 | `TestServiceRotationWindow`：A 注册 → B/A 窗口可登录 → 重封后仅 B 可解 → 仅 A → `ErrMFAInvalid` | **fixed** |
| **F-005** | captcha `SELECT` 再 `DELETE`，PG READ COMMITTED 双胜 | 单语句 `DELETE … WHERE id=? AND expires_at>? AND answer_hash=?` + `RowsAffected==1`；同事务无守卫 DELETE 保持「任意尝试即消费」 | `TestPostgresCaptchaConsumeConcurrent` **PASS 1.77s（真实 PG）**：并发恰好 1 胜；sqlite lifecycle 仍绿 | **fixed** |
| **F-006** | 坏 JSON 静默全库对账；submit/cancel/retry 只需 `wallet.read` | Decode 失败 400 `INVALID_BODY`；`io.EOF` 空 body 保持全库哨兵；POST reconcile/cancel/retry → `wallet.write`；GET 仍 `wallet.read`。空 `accountId` 哨兵按 D-002 有据保留 | `TestWalletReconcileBadBodyAndWriteGate`（垃圾 body 400；editor 三写口 403）；既有空 body 202 仍绿 | **fixed** |

未发现「假修复」或偷换 A-001 范围。D-002 未选方案（独立 MFA 密钥、显式全库哨兵字段）与代码一致。

## recommended 处置核对（E-003）

| F-ID | E-003 | 本轮核对 | 结论 |
|------|-------|----------|------|
| F-007 | fixed | locked/disabled 在返回前 `VerifyPassword(timingDummyHash, …)`；handler 对该两分支 `rateLimiter.record` | **有据 fixed** |
| F-008 | fixed | `Restore` 同一 `runner.Run` 内 `restoreRowTx` + `MarkRestoredTx`；`TestRestoreAtomicityRollsBackOnFailedMark` | **有据 fixed** |
| F-009 | 部分 fixed + residual | 未知 handler 记 `failed` 并返回 error，不再降级 `system.noop`/`ran`；`lastRun` 仍进程内存（文件头/compose 单实例语义） | **有据**（与 A-002 R-001 移交一致） |
| F-010 | fixed | `restoreSession` 刷新失败再试一次；401/403 清 token 语义不变；vitest「瞬时抖动后恢复」 | **有据 fixed** |
| F-011 | fixed | `coerceToKind("number")`：`""`/`undefined`/`null` → `undefined`；垃圾串仍 0 | **有据 fixed** |
| F-012 | fixed | effect 依赖 `JSON.stringify(crud?.route ?? null)`，同页 query 变化会重跑 prefill | **有据 fixed** |
| F-013 | fixed | `urlEscape` 对非 RFC3986 unreserved 全部百分号编码 | **有据 fixed** |
| F-014 | fixed | `errors.As(*http.MaxBytesError)` → 413 `FILE_TOO_LARGE` | **有据 fixed** |
| F-015 | overruled | 服务端仍单 token 吊销；`doRefresh` 非 2xx 且 localStorage 已有更新 token 时重试（跨标签旋转 A-002 F-003）。家族吊销会把该协议变成双标签互踢。复审触发已写 | **overruled 有据** |
| F-016 | fixed | `dictEntryFromPayload` 写回 `BadgeStyle` | **有据 fixed** |
| F-017 | fixed | `formulaSafe` 前缀含 `\t`/`\r` | **有据 fixed** |
| F-018 | fixed | 配额检查 + `storeUploadForOwner` 持 `avatarQuotaMu` | **有据 fixed** |
| F-019 | overruled | `runCustomAction` 白名单映射到已鉴权 API（export / library）；schema 无 permission target；硬门禁在服务端。审计原文已注「硬门禁仍在 API」 | **overruled 有据** |

## 回归复跑（本会话独立执行）

| 命令 | 结果 |
|------|------|
| `apps/api` `go vet ./...` | **0** |
| `apps/api` `go test ./...` | **全绿**（handler 169.9s；store 含 PG 门控） |
| `TestPostgresCreateUserManagementExistsScan` | **PASS 1.43s**（真实 PG，非 skip） |
| `TestPostgresCaptchaConsumeConcurrent` | **PASS 1.77s**（真实 PG，非 skip） |
| `apps/web` `npx vitest run` | **76 files / 1085 tests 全绿** |
| `apps/web` `npx tsc -b` | **0 错误** |

与 A-002 主张一致；本轮额外确认 F-001/F-005 在 Postgres 上真实执行。

## Findings

### required

无。开放 required = **0**。

### recommended

- **R-001**（low–med，F-003 残余，不阻断闭合）：`/api/auth/mfa/verify` 仍无独立 HTTP 限流；`Verify` 先 `ValidateTotp` 再 `IncrementProofFailures`，并发窗口内同一 proof 可多做几次 TOTP 校验（计数仍被 SQL `fail_count < 5` 封顶）。proof 签发已被登录桶限制，6 位 TOTP 在约 100 次/15min 下仍不可实用穷举。复审触发：缩短 TOTP 位数、关闭验证码且把第二因子当唯一远程秘密、或威胁模型升级为「verify 层 DoS/并发放大」。
- **R-002**（low，测试缝，不阻断）：`TestWalletReconcileBadBodyAndWriteGate` 用无任何 wallet 键的 `editor` 证 403；生产代码三写口已是 `wallet.write`。更贴原主张的锁是「仅有 `wallet.read`、无 `wallet.write`」的角色。

### informational

- **I-A** F-009 `lastRun` 进程内存去重：与 A-002 R-001 一致，单实例文档化 residual；多副本需 DB 租约，非本波。
- **I-B** F-004 `UpdateSecretCiphertext` 无 CAS：重封窗口内若发生 admin reset，理论上可能把旧密文写回；窗口极窄，best-effort 重封失败不阻断已通过的第二因子（与注释一致）。
- **I-C** F-002 工厂回归锁覆盖 dict-types；entries/scheduled-tasks 走同一 `DeleteTrashTx` 接口，无独立 HTTP 回滚用例。
- **I-D** D-002 §4 曾写 recommended 在 A-003 **之后**处置；实际 E-003 已先做。属顺序偏差，内容有据，不构成缺陷。
- **I-E** `goal-tree.md` 仍写 GOAL-011 `1/4`、workspace.md W11 行仍写「S2 未开始」；与 meta `3/4`、D-002/E-002/A-002 不一致。编排器关门时同步（本意见不改）。
- **I-F** A-001 informational F-020～F-025 保持记录项，本波未升格。

## 必改项汇总

开放 required = **0**（fixed ×6）。无新 required。recommended 残余不阻断 S4 代码闭合。

## 与既有意见的异同

- 与 **A-001**（fail，6 required 开放）：6 条均按原文处方落地，主张不再成立。
- 与 **A-002**（self pass）：独立复跑一致；同意 6/6 fixed 与 recommended 处置。新增 R-001/R-002 为残余与测试缝，**不与 A-002 冲突**，不升格 required。
- F-015/F-019 overruled：与客户端跨标签重试协议、API 硬门禁一致，接受 E-003 依据。

## 结论 + 建议给编排器/用户

**verdict = pass。** 代码闭合条件已满足：A-001 六条 required 在 `72a5397` 为 genuine fixed；未改错既有 fail-closed 语义；E-003 11 fixed + 2 overruled 有据；本会话 API/Web 回归全绿，且 F-001/F-005 在真实 Postgres 上通过。

建议 `/govern`：① 将本意见落盘并更新索引；② 闭合记录（F-001～F-006、F-007～014/016～018 → fixed；F-015/F-019 → user-overruled 有据；I-003 → verified）；③ 按 D-002 写恢复 VP-008 go 宣称 + 关门（用户书面确认后）；④ 同步 `goal-tree.md` / `workspace.md`；⑤ R-001/R-002 移交后续波次，不阻塞关门。

## 声明

本意见 **source=independent**，**不修改** 目标 `status` / `progress` / 检查点 / 方案正文 / goal-tree。响应、finding 闭合与关门由 **`/govern`** 处理。最终关门仍归编排器与用户书面裁决。