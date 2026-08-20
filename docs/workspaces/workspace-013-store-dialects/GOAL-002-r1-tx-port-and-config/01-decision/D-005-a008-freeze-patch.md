---
id: D-005-a008-freeze-patch
doc: decision-entry
status: accepted
created: 2026-08-20
updated: 2026-08-20
parent: GOAL-002-r1-tx-port-and-config
version: 1.0.0
---

# D-005 · 响应 A-008：冻结合同 v1.4（全部 recommended `fixed`）

- **日期**：2026-08-20
- **状态**：accepted
- **工作区**：`workspace-013-store-dialects`
- **用户确认**：本轮 `/govern` 输入「响应 GOAL-002 A-008」。A-008 independent **pass**，开放 required = 0。F-001～F-004 均为 recommended。与 A-002 / A-004 / A-006 响应同口径：recommended 一并 `fixed`（不走 residual / overruled）。A-008 允许写入合同补丁或登记到 R2/R3 方案边界；本条选前者，避免 R2 配置校验与 R3 对写再猜。

## 决定

修订 [r1-tx-port-and-config-freeze.md](../attachments/r1-tx-port-and-config-freeze.md) 至 **v1.4.0**。D-001～D-004 的 Tx 形状、占位符 `?`、upsert、时间单位与宽度、配置键名、path 文件根、`WasFresh`、R2 postgres `Open` 不 apply catalog、无 ORM、缺省 sqlite 仍成立。本条只补 A-008 指出的 recommended 边界。GOAL-002 **保持 `done`**（不改检查点 2/2）。本回合不改 `apps/api`。

1. **F-001**：postgres `db.path` 增加谓词 `filepath.Ext(filepath.Base(path))` 必须非空。拦住尚不存在的 `./data` / `.\data`（谓词 1–3 只拦尾部分隔符与已存在目录）。缺省 `./data/schema-ui.db`、cwd 文件 `schema-ui.db`、`./schema-ui.db` 必须通过。不要单靠 `Dir == "."` 拒绝。
2. **F-002**：点名 `COLLATE NOCASE` 为 R3 成对改写债：DDL `service_credentials.name … COLLATE NOCASE UNIQUE`；查询 `users_repository.go` / `roles_repository.go` 的 `ORDER BY … COLLATE NOCASE`。禁止按字面抄进 postgres。
3. **F-003**：checksum **算法不变**；digest 输入 = 该 version 的 sqlite/canonical 历史 stmts + 既有 transform id。postgres 成对 SQL 不进入该 digest。可选同一 catalog（checksum 绑 sqlite）或按方言分列 catalog。禁止改 hash 函数，禁止把「算法不变」读成两方言 DDL byte-identical，禁止改写 sqlite 历史文本去对齐 `BIGINT`。
4. **F-004**：嵌套 `Run` 按当前回调检测（`ctx` / 调用栈 / goroutine 局部）。禁止 Store 级或进程级 `inRun` 互斥（会误杀 postgres 并发 `Run`）。

不另立 GOAL-002 I-00N：缺口是合同文本，不是新的未知信息。A-008 无 required，本条不打开新的冻结门禁。

## 理由

A-008 independent `pass` 给出可核对边角：path 谓词与自身 `./data` 禁例不完全同构；`COLLATE NOCASE` 会在 postgres apply/排序失败；成对 DDL 与唯一 checksum 冲突若到 R3 才猜会分叉；Store 级嵌套标志与并发 `Run` 互斥。用户要求响应本条；与前三次响应同口径写入附件，R2/R3 实施者不必另翻审计原文。

F-001 选扩展名谓词（A-008 建议）：与缺省 `.db` 同构，且能区分 `schema-ui.db`（cwd 文件，`Dir == "."` 合法）与 `./data`（目录分量、无扩展名）。

F-003 选「checksum 绑 sqlite 历史文本、postgres SQL 不进 digest」为主路径，并允许方言分列 catalog：满足「算法不变」且不漂移现网 sqlite 库。

## 未选方案

- 只把 F-001～F-004 登记到未来 R2/R3 方案、不改本附件：A-008 允许，但与本目标前三次「recommended 一并写入合同」口径不一致；R2 配置校验仍会在 `./data` 上分叉。
- F-001 单靠 `Dir == "."` 拒绝：误杀 cwd 文件 `schema-ui.db`。
- F-003 改写 sqlite 历史 stmts 为 `BIGINT` 以保持一份 SQL：现网 checksum 漂移 fail-closed。
- F-003 改 hash 函数：违反「算法不变」。
- `accepted-residual` / `user-overruled`：用户未选。
- 重开 GOAL-002 或改 `status`：A-008 为 pass，无必改门禁；补丁记在本条 + 附件版本。GOAL-002 `done` 仍不证明冻结质量。

## 影响范围

- R2：`db.path` 扩展名谓词；`Run` 嵌套检测不得用 Store 级互斥。Open/Ping 对照 **v1.4**。
- R3：点名 `COLLATE NOCASE`；checksum 输入绑 sqlite 历史 SQL。时间列 postgres 仍为 `BIGINT`（v1.3）。
- 不改 Profile、Manifest、模块矩阵、Compose 默认。
- Root I-002（驱动）仍 open，最晚 R2 方案冻结。

## 关联信息项

- GOAL-002 I-001：仍 verified（打开 + `WithTx`）；A-008 所指 recommended 由本条写入 v1.4，不另立 I-00N。
- GOAL-002 I-002：仍 verified。本条补上 A-008 响应，不回溯改模式记录。
- Root I-002 / I-003 / I-001 / I-004：不变。

## 后续

R2 立项前先冻结 Root I-002（驱动）。实施合同 = 本附件 **v1.4**。A-008 为 pass，关闭证据复审可选。R2 方案冻结后再 `/audit`（independent）。
