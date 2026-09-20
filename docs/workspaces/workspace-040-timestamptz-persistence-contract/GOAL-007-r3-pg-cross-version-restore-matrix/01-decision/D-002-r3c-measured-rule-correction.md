---
id: D-002-r3c-measured-rule-correction
doc: decision-entry
status: accepted
parent: GOAL-007-r3-pg-cross-version-restore-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# D-002 · 矩阵实测对 `D-001` §3 预期规则的外推更正

## 决定的来源

- `D-001` §3 在冻结判定口径时，把前置探测（平凡 schema、只测了 client 15 与 17 两档）得到的观察**外推**成预期规则：`dumper ≤ client ≤ target_server`，并写明「矩阵本体若出现与上述规则不符的格，按 `unexpected-failure` 处理并逐格记明」。
- 矩阵本体（真实 VP-040 迁移链、54 个 restore 格，证据 `attachments/r3c-pg-cross-version-matrix-v0.1.md`）**出现反例**：由 16 client 生成的归档用 **16 client 恢复到 15 server** 成功且形状校验通过（`src15_by_16`→16→15、`src16_by_16`→16→15 两格）。
- 该反例不是失败格，因此不触发 `D-001` §3 的 `unexpected-failure` 分支；被更正的是**规则陈述**，不是判定分类或组合枚举。

## 1. 更正后的实测规则

| # | 规则 | 证据 |
|--:|------|------|
| 1 | `pg_dump`：client major **≥** server major；低于则 `pg_dump: error: aborting because of server version mismatch`（3/9 格） | dump 表（9 格，6 supported） |
| 2 | `pg_restore` 读归档：client major **≥** 生成归档的 client major；否则 `unsupported version (1.1x) in file header` | 归档格式版本随**生成归档的 client** 递增（1.14/1.15/1.16）；restore 表 24 格 toolgate |
| 3 | `pg_restore` 写目标：**PG 17 的 client 只能恢复到 PG 17 server**（17 client 的恢复会话设置 PG 17 引入的 `transaction_timeout`，15/16 server 不识别；12 格 serverguc） | restore 表 serverguc 行 |
| 4 | 目标 server major **不要求** ≥ client major：client ≤ 16 时未观测到 server 侧门（16 client → 15 server supported） | 2 个反例格 |
| 5 | 归档由哪个 server 生成，不构成额外的目标 server 约束（16 server 的内容经 16 client 可恢复到 15 server） | `src16_by_16`→16→15 |

**表述形式（供后续引用）**：

```text
dump:     client_major >= server_major
restore:  client_major >= dumper_client_major_default:  client_major == 17  ⇒  target_server_major == 17
```

## 2. 对 `D-001` 的处置

- **`D-001` 不重写、不撤销**：其 §1（版本取值）、§2（组合枚举 9+27/实测 54）、§3 的判定分类与 fail-closed 要求、§4（形状校验四项）、§5（驱动/清理）**全部有效**并被本次运行遵守。
- 被更正的是 `D-001` §3 末段对探测结果的**外推句**；本节即其修订记录，`D-001` 保持历史原文以便追溯（AGENTS：不重写已落盘决策，用新条目更正）。
- 分类结果**未**受影响：矩阵没有出现 `unexpected-failure`；54 格全部落在 `supported`（18）/`unsupported-toolgate`（24）/`unsupported-serverguc`（12）之内。

## 3. 对 `I-041-004` 的收口

信息项 `I-041-004`（non-blocking，继承）要求「PG 15/16/17 跨版本 `pg_restore` 兼容矩阵」。本轮以**逐组合实测记录**满足：

- 9 个 dump 格与 54 个 restore 格全部落盘（含 unsupported 的原因与退出码），无跳过、无概括；
- 18 个 supported 格的形状校验全部通过（ledger head 87 / `timestamptz(6)` / 微秒尾零逐字节 / sentinel NULL）；
- 附带正向结论：VP-040 迁移链在 PG 16.15 与 17.11 上可完整应用至 v87。

→ `I-041-004` = **verified**（本目标检查点 B 关闭），**不**记 residual：矩阵是**可逐格复核的事实记录**，不存在「未收集的必需信息」。

## 4. 边界（不得越读）

- 容器结果仍是 **CI/reproducibility** 证据，不是生产就绪证据（`D-017` §3 约束②，用户裁决）。
- 本节**不**作出产品级「支持哪些组合」的承诺；若 R3-D 退出矩阵需要该表述（例如「发布说明应写明 17 client 不可恢复到 15/16」），按 P-004 由用户裁决并留痕。
- 结论绑定 `attachments/r3c-pg-cross-version-matrix-v0.1.md` 记录的版本串（15.19 / 16.15 / 17.11）与命令形态（`pg_dump -F c --no-owner`、`pg_restore --exit-on-error --no-owner`）；镜像 tag 可变，复跑须重新记录。

## 未选方案

- **直接改 `D-001` 正文**：会抹掉「先冻结口径、后被实测更正」的可追溯链，且违反本项目既有做法（用新条目更正），未采用。
- **把反例格改判为 `unexpected-failure`**：反例格退出码 0 且形状校验通过，是真的 supported；改判等于用规则否定实测，未采用。
- **只记录 supported 格、把 unsupported 归为「工具问题」不落盘**：违反 `D-018` §2 第 4 项「逐组合记录 supported / unsupported」，未采用。
