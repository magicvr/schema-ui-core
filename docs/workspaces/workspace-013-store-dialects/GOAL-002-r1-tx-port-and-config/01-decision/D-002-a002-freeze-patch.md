---
id: D-002-a002-freeze-patch
doc: decision-entry
status: accepted
created: 2026-08-20
updated: 2026-08-20
parent: GOAL-002-r1-tx-port-and-config
version: 1.0.0
---

# D-002 · 响应 A-002：冻结合同 v1.1（全部 required `fixed`）

- **日期**：2026-08-20
- **状态**：accepted
- **工作区**：`workspace-013-store-dialects`
- **用户确认**：`ok 全部 fixed`（不走 residual / overruled）

## 决定

修订 [r1-tx-port-and-config-freeze.md](../attachments/r1-tx-port-and-config-freeze.md) 至 **v1.1.0**。D-001 的 Tx 形状、占位符 `?`、配置键名、无 ORM、缺省 sqlite 仍成立；下列补丁覆盖 A-002 F-001～F-006。F-007 由本响应过程闭合。GOAL-002 **保持 `done`**（不改检查点 2/2）。本回合不改 `apps/api`。

1. **F-001**：把 upsert 与时间规则写入本冻结（不 defer 到 R3 I-00N；不新增 Go upsert API）。允许 `ON CONFLICT(cols) DO UPDATE/NOTHING`；禁止 `INSERT OR IGNORE` / `INSERT OR REPLACE` / `REPLACE INTO`。时间由 Go UTC 绑参（默认 INTEGER Unix 秒）；禁止 `datetime('now')` / `now()` / 不等价 `CURRENT_TIMESTAMP`。现行 `operationlog/retention.go` 的 `INSERT OR IGNORE` 留 R3/R4 改写。
2. **F-002**：postgres 下 `db.path` **不作 SQL**，**仍作数据目录根**（`filepath.Dir` → uploads / brand-assets / avatars）。省略时仍缺省 `./data/schema-ui.db`。不新增 `db.data_dir`。监控 `os.Stat` 仅 sqlite 报库文件大小；postgres 下 `DBSizeBytes=0`。
3. **F-003**：`WasFresh()` 留在 `kernel.Store`。方言中立：Open 后、迁移 apply 前用户表数为 0（sqlite = 现行 `sqlite_master` 探测；postgres = `search_path` 首个用户 schema 内非系统表为 0）。快照只读。`MarkSystemDataReady` / `SystemDataReady` 保持进程内 atomic。
4. **F-004**：`OpenOptions` 可在不改 Tx 的前提下增字段；`Tx` 仅在 `Run` 回调内有效；sqlite 串行 `Run`，PG 池允许并发 `Run`，嵌套仍禁。
5. **F-005**：`errors.Is` 对 `kernel.ErrNoRows` 与 `sql.ErrNoRows` 同真。
6. **F-006**：模块仓库禁止按 `Dialect()` 分支 SQL。

降级 A-001「合同与 VP-013 / RT-P03 同构」：v1.1 覆盖连接/事务/占位符/upsert/时间 + 配置面；迁移 runner 对写与备份仍属 R3/R5。R2 **必须按 v1.1** 实施，不得按 v1.0 猜测。

## 理由

A-002 independent `conditional` 给出可核对反证：附件未冻 upsert/时间却要求模块 SQL 可移植；`db.path`「不使用」与 composition 文件根冲突；`WasFresh` 无 PG 语义。用户确认三条必改与四条 recommended 全部 `fixed`。规则现在可写清，无需新信息项。不拆 `WasFresh`、不新增配置键，以免 R2 换一套类型或键名。

## 未选方案

- F-001 登记 required I-00N、最晚 R3：用户否决；拖会使 R3 无对照合同。
- F-002 新键 `db.data_dir`，或字面「postgres 完全不用 path」：拆现网文件根或增加无必要配置。
- F-003 把 bootstrap 就绪移出 `kernel.Store`：composition/测试仍依赖该方法。
- `accepted-residual` / `user-overruled`：用户未选。
- 重开 GOAL-002 或改 `status`：F-007 不要求因审计模式改 status；补丁记在本条 + 附件版本。

## 影响范围

- R2：`store.Open`、配置校验、`WasFresh` 实现、文件根 mkdir、监控 Stat。
- R3/R4：迁移/模块 SQL 必须遵守 upsert 与时间规则。
- 不改 Profile、Manifest、模块矩阵、Compose 默认。
- Root I-002（驱动）仍 open，最晚 R2 方案冻结。

## 关联信息项

- GOAL-002 I-001：仍 verified（打开 + WithTx）；A-002 所指未登记缺口由本条写入合同，不另立 I-00N。
- GOAL-002 I-002：仍 verified（当时 self）；本条补上冻结层 independent 响应，不回溯改模式记录。
- Root I-002 / I-003 / I-001 / I-004：不变。

## 后续

R2 立项前先冻结 Root I-002（驱动）。实施合同 = 本附件 v1.1。建议 R2 方案冻结后再 `/audit`（independent）。
