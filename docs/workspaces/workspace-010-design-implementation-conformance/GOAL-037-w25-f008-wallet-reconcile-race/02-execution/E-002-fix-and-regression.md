---
date: 2026-08-23
scope: GOAL-037 S3（修复与回归验证）
---

# E-002 · F-008 修复实施与回归（S3）

## 修复实施（文件级证据）

1. **产品 `newID`**（`apps/api/internal/modules/wallet/provider.go`）：id 改为 `UnixMilli(16hex) + 进程内同毫秒单调计数(8hex) + 随机后缀(24hex)`——D-002 §1 回放排序契约（created_at ASC, id ASC）在**同一毫秒**内由计数恢复创建序；注释同步修正（原"毫秒前缀保证同秒创建序"在"同毫秒"处失效的假设）。`entryIDSeq atomic.Uint64`（& 0xFFFFFFFF 防溢出）。
2. **测试替身 `walletServiceStub`**（`apps/api/internal/handler/wallet_test.go`）：Mutate/Reconcile 的 id 生成同样加同毫秒单调计数（`walletStubEntryID`）——handler 测试实际暴露点（stub 生成 `{ms}+{random}` 才是回放乱序的直接来源）。
3. **测试断言同步化**：`wallet.reconcile` 审计事件为 job 提交后的 **best-effort 异步副作用**（`RegisterWithTerminalHook` 后置写 operationlog，产品与测试 stub 均吞错）→ 断言改为**轮询等待**（2s deadline，5ms 间隔，六事件全集命中）。
4. 附带：`wallet_test.go` 失败输出携带 reconcile `details`（mismatch reason 可诊断性，E-001 记录）。

## 机制结论（E-001 补充）

- 根因链：同毫秒多笔流水 → 旧 id 随机后缀决定回放字典序 ≠ 写入序 → 回放把 `freeze` 排在 `adjust` 前 → `Apply` 反负（`replay apply failed: insufficient balance`）→ reconcile inconsistent。
- 与池化/FK/txlock 无直接因果；池化+`synchronous=NORMAL` 提高同毫秒连发概率而暴露（基线同窗低频存在）。
- 第二暴露：terminal 审计事件与测试查询的异步时序竞态（独立于 id 排序，本目标一并修复断言侧）。

## 回归（事实）

- `go test ./internal/handler/ -run TestWalletLifecycleAndAdjustFlow -count=100`：**100/100 ok**（修复前该窗口混合失败：inconsistent / 审计缺失）。
- `go test ./internal/handler/ -count=30`（同用例）：ok。
- `go test ./... -count=1` **两轮全绿**（0 FAIL）；`go test ./internal/modules/wallet/`（含新 `TestNewIDSameMillisecondOrdering` / `TestNewIDCrossMillisecondOrdering`）：ok。
- vitest **1097/1097**；`tsc -b` 0（前端未动，伴行确认）。
- 提交：`be62582`（F-008 代码修复）；治理文档随关门提交。

## 残余（D-001 已记录）

既有库中"同毫秒随机序"旧流水在重对账时仍可能 inconsistent（触发苛刻、量级可忽略）；复审触发已记录于 D-001。