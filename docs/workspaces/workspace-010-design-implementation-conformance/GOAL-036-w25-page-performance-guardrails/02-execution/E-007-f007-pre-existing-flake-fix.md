---
date: 2026-08-23
scope: GOAL-036 A-001 F-007（预存 flake）根因定位与修复
---

# E-007 · F-007 修复：`TestScheduledTaskRunsPagination` 预存 flake 根因与处置

## 根因（诊断事实链，2026-08-23）

1. **症状**：`POST /api/scheduled-tasks/{id}/run` 偶发 500（`task execution failed`）；失败总落在相邻 POST（第 2/3 次）。基线（stash @HEAD，单连接、无 W25 改动）`-count=20` 复现 → **预存 flake 确认**（非本波引入）。
2. **错误定位**：临时将 500 响应携带 err 文本 → `insert task run <id>: constraint failed: UNIQUE constraint failed: task_runs.id`。`(1555)` 为 SQLite 错误码后缀（非冲突值）。
3. **候选排除**：产品 `newRunID` 的 crypto/rand 失败回落路径（时间串）→ 已先修（`runIDSeq` monotonic + UnixNano，rand 仅扰动），**仍复现** → 不是产品路径。
4. **真因**：handler 测试替身 `testTaskRunner.Execute`（`testhelpers_test.go`）的 run id = `"run-test-" + now.UnixNano()`。时钟 dump（POST 前打印 `time.Now().UnixNano()`）证明 **Windows 系统时钟量化下相邻调用返回完全相同纳秒值**（如 `NOW 0 == NOW 1 = 1787482946873605800`）→ 相邻 POST 生成相同 id → 主键冲突。冲突 id 每次不同（`run-test-1787482975726705500` 等）与"时钟量化偶发同值"自洽。
5. **修复**：`testTaskRunner.Execute` id 改为 `"run-test-%x-%x"`（UnixNano + 包级 `atomic.Int64` 单调计数）——进程内绝对唯一，与时钟量化无关。

## 保留的产品/测试改进（防御纵深，均有测试）

- `scheduler.go newRunID`：**保留**本次加固（`(UnixNano, monotonic seq)` 唯一性主体 + crypto/rand 仅扰动）——若熵源在受限环境退化为常量，随机 id 路径同样会撞；新增 `TestNewRunIDUniqueUnderEntropyFailure`（rand 失败 / rand 恒常量两场景 ×1000 次唯一）。
- `store/repository.go RecordRun`：失败错误文本携带 `run.ID`（可诊断性改进，保留）。

## 验证

- `go test ./internal/handler/ -run TestScheduledTaskRunsPagination -count=60`：**60/60 ok**（修复前该窗口高频失败）。
- `go test ./internal/modules/scheduledtasks/`：ok（含新熵失败唯一性测试）。
- 全量 `go test ./...`：见下条核对记录（除本 flake 外此前各包全绿的历史不再成立——本条目修复后全量应整绿）。

## 归因更正（相对 A-002/E-006 的记录）

A-002/E-006 曾将 F-007 候选机制记为「newRunID 回落路径同毫秒撞 ID」——该候选经论证被排除（修复后仍复现），**真因为测试替身 id 生成器依赖时钟量化精度**。产品 `newRunID` 加固仍为有效防御并保留。本条目为更正后的权威记录。

## 连环发现：wallet reconcile 竞态（新 finding F-008，移交）

F-007 修复后的高频回归暴露**第二个间歇点**：`TestWalletLifecycleAndAdjustFlow`（wallet 全生命周期）在池化+FK 时代偶发：

1. **SQLITE_BUSY（已修）**：`update wallet balances: database is locked`——WAL 下读后写事务被并发 checkpoint 淘汰读快照（SQLITE_BUSY_SNAPSHOT 不等待 busy handler）。修复：DSN `_txlock=immediate`（写事务 BEGIN 即取写锁并固定快照；busy_timeout 让写者排队）。A/B 验证：无 txlock 时 freeze/adjust 500 高发，有 txlock 时全部消失。
2. **reconcile result = inconsistent（开放，F-008）**：A/B 证明与 txlock 无关（两配置下均偶发）；E-003/E-004 时代（池化+WAL、FK=OFF、无深挖）该用例曾绿，现确认在池化+FK 配置下不稳定窗口存在（曾被 BUSY 提前失败掩盖）。机制未定论（候选：异步 job 对账事务与主链写入的可见性窗口）；**已记录为独立 finding 移交后续波次**（不在本轮强制闭合——样本频率低、与本轮已修的确定性缺陷正交）。

## 验证（F-007 本体）

- `go test ./internal/handler/ -run TestScheduledTaskRunsPagination -count=60`：**60/60 ok**（修复前该窗口高频失败）。
- `go test ./internal/modules/scheduledtasks/`：ok（含 `TestNewRunIDUniqueUnderEntropyFailure`：rand 失败/恒常量 ×1000 唯一）。
- 全量 `go test ./...`：除 wallet 偶发（F-008）外全绿。