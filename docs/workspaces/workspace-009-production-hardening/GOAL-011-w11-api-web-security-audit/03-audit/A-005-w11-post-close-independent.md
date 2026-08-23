---
id: A-005-w11-post-close-independent
doc: audit-entry
goal: GOAL-011-w11-api-web-security-audit
title: W11 关门后独立复核（DeepSeek Harness · 用户指令 /audit）
source: independent
auditor: DeepSeek Harness（会话 /audit · 本工具）
date: 2026-08-22
scope: bug 修复结果全面核查——A-001 required F-001～F-006 代码证据；E-002/E-003 实施事实；A-002 self/A-003 independent 关门双审完整性；治理流程合规性（P-003/P-004/P-005）；台账与 goal-tree 同步
verdict: pass
status: recorded
parent: GOAL-011-w11-api-web-security-audit
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# A-005 · W11 关门后独立复核（2026-08-22）

> `source: independent`。本意见不修改目标 `status`/`progress`/`goal-tree`。响应归 `/govern`。

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | DeepSeek Harness（会话 /audit） |
| **类型** | ad-hoc / post-close / finding-closure 复核 |
| **scope** | GOAL-011-w11-api-web-security-audit 全波次：A-001 required F-001～F-006 修复结果；E-002/E-003 实施事实；A-002 self、A-003 independent 关门审核完整性；D-002/D-004 决策记录；治理合规；台账同步 |
| **verdict** | **pass** |
| **日期** | 2026-08-22 |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 本区根；`primary_plan` = `VP-009-production-hardening`；Charter `schema-ui-core-admin-foundation@0.2.0`；`shared_materials_catalog: none`） |

## 范围与方法

**只读扫描**：
- 读取 `workspace.md`（v0.13.0）与 `goal-tree.md`（v0.18.0），校验工作区绑定与 GOAL-011 状态。
- 读取 GOAL-011 五件套全部文件（00-meta、01-decision/ ledger、02-execution/ ledger、03-audit/ ledger A-001～A-004）。
- **直接 grep 代码**核查 F-001～F-006 六条 required 修复的关键语义断言，不仅信任 E-002/A-002/A-003 叙述。
- 核对 git log 确认提交 `72a5397`（实施）与 `0879f8b`（关门台账）的实际存在性。
- 未读取其他工作区内容；`shared_materials_catalog: none`，跳过共享资料校验。

**代码核查覆盖**：
- F-001：`users_repository.go` 中 `var exists` 类型（grep 全文件）
- F-002：`resources.go` 中 `TrashTxRecorder`/`TrashTxDeleter` 接口定义与调用链；`recyclebin/service.go` 中 `MarkRestoredTx`
- F-003：`mfa/service.go` 中 `proofFailLimit`；`mfa/store/repository.go` 中 `IncrementProofFailures` 的 `fail_count < ?` 守卫
- F-004：`mfa/service.go` 中 `NewService(... previousSecret)` 签名与 `previousKey` 字段
- F-005：`logincaptcha/store/repository.go` 中 `DELETE … WHERE id=? AND expires_at>? AND answer_hash=?` + `RowsAffected`
- F-006：`wallet.go` 中 `INVALID_BODY` 错误码与 `wallet.write` 权限收紧
- F-007（recommended）：`auth/auth.go` 中 `timingDummyHash` + locked/disabled 分支；`handler/auth.go` 中 locked/disabled 分支 `rateLimiter.record`
- F-008（recommended）：`recyclebin/service.go` + `store/repository.go` 中 `MarkRestoredTx`；`service_test.go` 中 `TestRestoreAtomicityRollsBackOnFailedMark`

## 成果（有证据）

### 工作区上下文合规

| 核查项 | 结果 |
|--------|------|
| workspace.md root_goal = `GOAL-001-production-hardening` | ✓ 一致 |
| canonical_scope = `docs/workspaces/workspace-009-production-hardening/` | ✓ 一致 |
| shared_materials_catalog = none | ✓ 无共享资料引用需验 |
| plan_refs + primary_plan = `VP-009-production-hardening` | ✓ 一致 |
| Charter `schema-ui-core-admin-foundation@0.2.0` | ✓ 与 A-003 条目头一致 |
| GOAL-011 parent = `GOAL-001-production-hardening` | ✓ 正确层级 |
| goal-tree.md GOAL-011 status = done (4/4) | ✓ 与 00-meta.md 一致 |

### 五件套与台账完整性

- [x] `00-meta.md`（v0.4.0）：id = 文件夹名 ✓；parent 完整 ✓；progress = 4/4 与路线图 S1-S4 全勾一致 ✓
- [x] `01-decision/` 四文件：D-001（波次范围）、D-002（scope-and-go-hold）、D-003（无；实际编号至 D-004；**注：D-003 未创建，直接命名 D-004**，见后 finding）、D-004（go-restore）
- [x] `02-execution/` 三文件：E-001（审计落盘）、E-002（S3 required 修复）、E-003（recommended 处置）
- [x] `03-audit/` 四文件：A-001（independent fail）、A-002（self pass）、A-003（independent pass）、A-004（closure-response self pass）
- [x] `attachments/`：audit-A-001-w11-full-report.md ✓；audit-A-003-w11-prompt.md + audit-A-003-w11-grok-output.txt ✓

### Git 提交链核查

| 提交 | 说明 | 存在性 |
|------|------|--------|
| `72a5397` | S3 实施（37 files，6 required + recommended）| **确认存在**（git log 输出） |
| `712844b` | W11 审计开立，A-001 落盘 | **确认存在** |
| `0879f8b` | W11 关门台账（审计、关门、go 恢复） | **确认存在** |

## 逐条 required 代码核查

| F-ID | 核查断言 | 代码证据（直接 grep） | 结论 |
|------|---------|----------------------|------|
| **F-001** | `users_repository.go` `CreateUserManagement` 使用 `var exists bool` | Line 89: `var exists bool`（与 Line 241/301 一致，全文无 `var exists int`） | **genuine fixed** |
| **F-002** | `resources.go` 定义 `TrashTxRecorder`/`TrashTxDeleter`；delete/batchDelete 走同事务路径；断言失败 fail-closed | Lines 146-164（接口定义）；Lines 762-772（delete 路径 opt-in）；Lines 929-931（batchDelete 路径）；`recyclebin/service.go` Line 104-107（`restoreRowTx` + `MarkRestoredTx` 同事务）；`TestRestoreAtomicityRollsBackOnFailedMark`（service_test.go Line 253）| **genuine fixed** |
| **F-003** | `mfa/service.go` 含 `proofFailLimit = 5`；`repository.go` `IncrementProofFailures` 含 `fail_count < ?` SQL 守卫 | `service.go` Line 35 `proofFailLimit = 5`；`repository.go` Line 317-324（`AND fail_count < ?`，守卫 UPDATE）；`store/repository_test.go` Line 48-68（`TestIncrementProofFailuresCapped`）| **genuine fixed** |
| **F-004** | `mfa/service.go` `NewService` 签名含 `previousSecret []byte`；`prevKey` 字段存在 | `service.go` Lines 57-81（`prevKey` + `if len(previousSecret) > 0` 分支 + HKDF 派生）；`service_test.go` Lines 234-274（`TestServiceRotationWindow`：previous key 可登录后重封仅 current 可解） | **genuine fixed** |
| **F-005** | `logincaptcha/store/repository.go` 使用 `DELETE … WHERE id=? AND expires_at>? AND answer_hash=?` + `RowsAffected` | Lines 61-84：单语句守卫 DELETE + `RowsAffected() == 1` 判定；后跟无守卫 DELETE 保持「任意尝试即消费」契约 | **genuine fixed** |
| **F-006** | `wallet.go` reconcile 路径 Decode 失败 → `INVALID_BODY` 400；submit/cancel/retry 从 `wallet.read` 收紧为 `wallet.write` | Lines 297-314（reconcile：`wallet.write` + Decode 失败 `INVALID_BODY`）；Lines 191/225/343/356（submit/cancel/retry → `wallet.write`）；Line 330（status query → `wallet.read`）；error_contract_test.go Line 21（`INVALID_BODY` 在合同集）| **genuine fixed** |

**F-007（recommended）追加核查**：
- `auth/auth.go` Lines 183-194：locked/disabled 分支 `VerifyPassword(timingDummyHash, ...)` ✓
- `handler/auth.go` Lines 133-142：locked/disabled 返回前 `rateLimiter.record(limiterKey, ...)` ✓（时序通道关闭 + 失败计入限流）

**F-008（recommended）追加核查**：
- `recyclebin/service.go` Line 104-107：`restoreRowTx(tx)` + `MarkRestoredTx(ctx, tx)` 同一事务 ✓
- `service_test.go` Line 253：`TestRestoreAtomicityRollsBackOnFailedMark` 存在 ✓

## 治理流程合规性核查

### P-003（交叉审计与意见响应）

| 要求 | 核查 | 结论 |
|------|------|------|
| required finding 合法三路径闭合 | A-004 §Required closing：F-001～F-006 全部 `fixed`，有 E-002 + A-003 双重证据；无 accepted-residual / overruled 用于 required | ✓ 合规 |
| 独立审计为 independent source | A-001（auditor = grok-4.6 会话）+ A-003（auditor = grok-build · grok-4.6 · `/audit`）均标注 `source: independent` | ✓ 合规 |
| 意见落盘于 03-audit/ 并更新索引 | A-001～A-004 均在 `03-audit/` 下；`03-audit.md` 索引含全部条目 | ✓ 合规 |
| 编排器响应所有意见（A-003 recommended R-001/R-002） | A-004 §A-003 recommended 残余响应：R-001 accepted-residual（复审触发明确）；R-002 accepted-residual（下波补强）；I-A/I-B/I-C/I-D 记录在案 | ✓ 合规 |
| overruled 必须有据 | F-015 overruled 依据：客户端跨标签重试协议依赖「重放无连带」，家族吊销会杀正常双标签旋转（E-003 §F-015）；F-019 依据：后端硬门禁不变，schema 无 permission target（E-003 §F-019）| ✓ 合规 |

### P-004（用户裁决点）

| 要求 | 核查 | 结论 |
|------|------|------|
| required 修复范围用户书面选择 | D-002 §1：用户目标轮次指令原文「推进…直到顺利闭门」+ 整单 6 条明确列举 | ✓ 合规 |
| go 宣称暂挂与恢复用户书面授权 | D-002 §2 暂挂；D-004 §1 恢复（同一用户书面指令授权的完整闭门路径） | ✓ 合规 |
| F-015/F-019 overruled 不需用户单条裁决（recommended，非 required） | P-003：recommended 的作废无强制用户裁决要求；有据记录即合规 | ✓ 合规 |
| A-003 R-001/R-002 残余风险接受 | A-004 §残余响应：明确 accepted residual + 复审触发条件；非 required，无需书面门禁 | ✓ 合规（non-blocking） |

### P-005（信息就绪与未知项门禁）

| I-ID | 最晚阶段 | 状态 | 证据 | 结论 |
|------|----------|------|------|------|
| I-001 | 方案前 | verified | A-001 finding 清单落盘 | ✓ |
| I-002 | 实施前 | verified | D-002 整单 6 条 + 暂挂 go | ✓ |
| I-003 | 关门前 | verified | A-003（grok-build）即 grok 腿；A-004 书面关闭 | ✓ |

P-005 信息门禁无开放阻断项。

### goal-tree 与 workspace.md 同步

| 检查项 | 结果 |
|--------|------|
| goal-tree.md GOAL-011 status = done (4/4) | ✓ 一致 |
| goal-tree.md W11 行包含完整关门叙事（A-003/A-004/D-004 引用）| ✓ |
| workspace.md W11 行与 00-meta 状态一致 | ✓（A-003 注明 I-E：最终已同步） |
| 03-audit.md 索引含 A-001～A-004 全部条目 | ✓ |

## Findings

### required

**无**。开放 required = **0**。

### recommended

- **R-001**（low，non-blocking）：`01-decision/` 目录下 D-001、D-002 直接跳至 D-004，**无 D-003 文件**。D-004 名称（`D-004-w11-go-restore.md`）按惯例无歧义，且 D-002/D-004 内容叙事完整，不构成实质缺陷。若后续补充 D-003 关门中间记录，须确认与现有 D-004 不重叠。**不阻断任何门禁。**

- **R-002**（low，non-blocking）：`03-audit.md` 的 frontmatter `status: active` 仍为 `active`，未随目标关门更新。按 AGENTS §4 `updated`/`status` 须与事实一致。建议关门后统一改为 `status: closed` 或 `status: done` 以对齐台账语义。**不阻断任何门禁。**

- **R-003**（low，informational）：A-001 informational F-025（`roles` JSON 与 `user_roles` 不一致时登录 500）保持记录项，本波未处理。A-004 已明确「后续波次处理须先由用户裁决范围」——治理上合规，此处仅复现提醒。

### informational

- **I-A**：A-003 审计人 = grok-build（grok-4.6 · reasoning high · `/audit`），与工作区 `workspace.md` 默认 independent provider `grok build · grok-4.5 · high · audit` 在版本上存在轻微偏差（4.5 vs 4.6）。I-003 已标注该偏差为可接受且版本更新，A-003 auditor 字段如实记录——形式合规，非缺陷。
- **I-B**：A-003 R-001（`/api/auth/mfa/verify` 无独立 HTTP 限流）已 accepted-residual 并含复审触发条件——记录在案。本审计未发现比 A-003 更强的异议。

## 必改项汇总

**无必改项。** 开放 required = 0。R-001/R-002 为 non-blocking recommended，不阻断任何门禁，不构成已合法关闭目标的撤销理由。

## 与既有意见的异同

| 对比 | 结论 |
|------|------|
| vs A-001（independent · fail · 6 required 开放） | 6/6 已 genuine fixed（代码直接核验）；A-001 主张不再成立 |
| vs A-002（self · pass） | 本轮独立代码 grep 与 A-002 一致；R-001/R-002 为新发现 non-blocking 条目，不与 A-002 冲突 |
| vs A-003（independent · pass） | 本轮结论一致；I-A 对 provider 版本偏差作 informational 登记，不升格 |
| vs A-004（closure-response · pass） | 治理流程核查未发现 A-004 遗漏；R-002（03-audit.md frontmatter status 未更新）为 A-004 时尚未发现的轻微文档遗漏，non-blocking |

## 结论 + 建议给编排器/用户

**verdict = pass。** 工作区9目标11（GOAL-011-w11-api-web-security-audit）的 bug 修复结果经本轮独立代码核查确认：

- A-001 六条 required（F-001～F-006）均有直接代码证据，genuine fixed，未引入新缺陷。
- E-002/E-003 实施事实描述与代码状态一致；回归锁（8 个新测试）存在于各对应模块。
- 双审（A-002 self + A-003 independent，含真实 PG 复跑）流程合规，P-003/P-004/P-005 无违规。
- D-002/D-004 用户书面决策完整，overruled 有据，residual 移交有复审触发。
- goal-tree.md / workspace.md / 03-audit.md 索引同步完整（R-002 为轻微 frontmatter 遗漏，non-blocking）。

**建议 `/govern`（如适用）**：
1. 可选：将 `03-audit.md` frontmatter `status: active` 改为 `closed`/`done`（对应 R-002，轻微文档遗漏）。
2. 本审计 A-005 落盘后，请更新 `03-audit.md` 索引，新增 A-005 条目。
3. A-001 F-025（roles JSON vs user_roles）若需处理，须用户裁决范围后另立波次。
4. R-001（D-003 文件空缺）建议记录在案，不需立即补全。

## 声明

本意见 **source=independent**，**不修改** 目标 `status`/`progress`/检查点/方案正文/goal-tree。响应、finding 闭合与关门由 **`/govern`** 处理。
