---
id: D-002-r1-contract-freeze-user-decisions
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-002 · R1 合同选型与用户裁决

## 触发

R1 的四项关键选型会改变物理 schema、编解码、存量库转换与回归范围。根据 P-004，本条记录用户在本轮对 `/vision` 提问的书面选择；实现前仍须由 R1 子目标形成逐表证据与自审/独立审计意见。

## 用户裁决

1. **PostgreSQL 物理类型**：改为字面 `timestamptz`，按 `timestamptz(6)` 的微秒语义实施；不沿用 VP-013 的 PG `BIGINT` 作为 VP-040 目标形态。
2. **SQLite 等价表示**：使用 `TEXT`，规范化为 UTC RFC3339 固定 6 位小数：`YYYY-MM-DDTHH:MM:SS.ffffffZ`。固定宽度用于保证文本排序与时间排序一致。
3. **首波分母**：纳入所有表达绝对时刻的持久化列，包括 `*_at`、`locked_until`、`applied_at` 等语义时间列；排除 ID 前缀、`duration_seconds`、TOTP step、version、计数器、金额与 flag。
4. **存量升级**：SQLite 与 PostgreSQL 各自提供原地转换；不提供产品级 SQLite→PostgreSQL 跨引擎搬运器。跨引擎继续承接 VP-013 已书面接受的 fresh bootstrap + 运维自备导出/回放 residual。
5. **NULL/零值**：有“未发生/无期限”语义的 sentinel `0` 转为 `NULL`；读取层映射为领域 absence/zero；真正非空的绝对时刻列禁止零值。R1 必须逐列标明 nullable 化、默认值与迁移回填。

## 边界与未选方案

- 本裁决只改变 C1 时间列的物理合同；不重开 VP-013 的 Store 端口、PostgreSQL 接入、checksum 台账、非时间整数宽度或公共契约决策。
- PostgreSQL `BIGINT` epoch 与 SQLite `INTEGER` epoch 作为 VP-013 的历史合同保留为迁移来源与证据基线，不作为 VP-040 新目标形态。
- SQLite→PostgreSQL 自动搬运器、ORM、第二数据库、Redis/MQ/A3 与 PITR 仍不在本 VP。
- `I-040-004` 的展示/输入回归矩阵与最终备份剧本细节仍待后续 R1/R3 证据，不在本条静默补齐。

## 后续门禁

- R1 子目标必须完成全仓 compiled catalog + 运行时读写的逐列 inventory，并给出旧单位（秒/毫秒）、新形态、NULL/默认值、转换函数与回滚/失败策略。
- 方案冻结前执行 self 审计，随后按项目默认路径调用本地 grok build（grok 4.6 · high）进行 independent 审计；未合法闭合 required finding 不得进入 R2。
