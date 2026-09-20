---
id: E-003-r3c-matrix-measured
doc: execution-entry
status: active
parent: GOAL-007-r3-pg-cross-version-restore-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-003 · 检查点 B：矩阵实测与 `I-041-004` 收口

## 事实

- **驱动落码**：`apps/api/internal/backup/pg_cross_version_matrix_test.go`（`VP040_PG_MATRIX=1` 门控，默认跳过；`VP040_PG_MATRIX_OUT` 可落盘报告）。
  - 3 个固定版本 server 容器（15-alpine / 16 / 17-alpine），各自发布随机宿主端口（宿主侧跑迁移与读取）并加入共享容器网络（client 侧跑 dump/restore）；归档经宿主临时目录 bind mount 交换；`pg_isready` 就绪门；`t.Cleanup` 删除容器与网络。
  - 每个 server 用**本仓真实迁移链**（`compiledmodules.PersistenceCatalog()` 全量 → head **87**）建源库并播种 3 行代表性数据（微秒尾零、sentinel 非空/空、毫秒族）。
  - 判定按 `D-001` §3 三类 + `unexpected-failure` fail closed：出现未知失败即整体失败。
  - 每个 `supported` 格执行 `D-001` §4 四项形状校验。
- **实测结果**（2026-09-21，退出码 0，`--- PASS (105.30s)`）：
  - 版本（实测）：server **15.19 / 16.15 / 17.11**；client 工具同版本；命令 `pg_dump -F c --no-owner`、`pg_restore --exit-on-error --no-owner`。
  - dump：**9 格**，6 supported / 3 `unsupported-toolgate`（client < server）。
  - 归档格式版本由**生成归档的 client** 决定：1.14 / 1.15 / 1.16。
  - restore：**54 格** = 18 supported（**形状校验全部通过**）/ 24 `unsupported-toolgate` / 12 `unsupported-serverguc`；`unexpected-failure` = 0。
  - **正向结论**：VP-040 迁移链在 **PG 16.15 与 17.11** 上可完整应用至 v87（三个源库分别建在三个 server 上，全部成功），且恢复后的库保持转换后合同形状（ledger head 87、`timestamptz(6)`、微秒尾零逐字节、sentinel NULL 保持）。
- **规则更正**：`D-001` §3 曾把探测结果外推为 `dumper ≤ client ≤ target_server`；矩阵给出反例（16 client → 15 server 为 supported 且形状通过）。实测规则更正为 `client ≥ server`（dump）、`client ≥ dumper_client`（归档门）、`client == 17 ⇒ target == 17`（server-GUC 门），见 `D-002` §1。**判定分类与组合枚举未被推翻**，`D-001` 保持原文以便追溯。
- **信息项**：`I-041-004`（non-blocking）→ **verified**（逐组合记录满足其信息需求；不记 residual）。
- **清理**：驱动自建容器/网络随测试清理；宿主临时目录由 `t.TempDir()` 回收；常驻 15.4 实例未被本矩阵使用或修改。

## 证据

- 实测记录（驱动原样输出，含 9+54 格逐格表与版本表）：`attachments/r3c-pg-cross-version-matrix-v0.1.md`。
- 前置探测（工具层规则与机制验证）：`attachments/r3c-pg-tool-compatibility-probe-v0.1.md`。
- 默认基线的隔离性：`go test -count=1 -run TestPGCrossVersionRestoreMatrix ./internal/backup/` → **SKIP**（未设门控），故本仓默认 `go test ./...` 仍不依赖 docker。

## 进度评估

检查点 B 完成 → `progress: 2/3`。**检查点 C 待办**：self 审计 + grok independent 审计、required 合法闭合 → 静默关门。
