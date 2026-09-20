---
id: D-001-r1-contract-freeze
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-001 · R1 合同冻结承接

本子目标承接 Root `D-002-r1-contract-freeze-user-decisions` 的用户裁决：

- PostgreSQL：`timestamptz(6)`。
- SQLite：UTC RFC3339 固定 6 位小数 `TEXT`，规范形如 `YYYY-MM-DDTHH:MM:SS.ffffffZ`。
- 分母：全部表达绝对时刻的持久化列；排除 ID、duration、TOTP step、version、计数、金额、flag。
- 升级：SQLite 与 PostgreSQL 各自原地转换；不提供 SQLite→PG 产品级搬运器，沿 VP-013 residual。
- sentinel：按语义把“未发生/无期限”的 `0` 转为 `NULL`；非空绝对时刻禁止零值。
- 公共 wire：用户选择统一输出 6 位微秒 RFC3339 UTC `Z`，R3 需同步验证 formatter、协议 fixtures 与 VP-020 展示回归。

## 仍需验证

上述是用户已选方向，不是逐列实现证据。C1 必须核对所有 compiled migration 与 runtime repository 的列、单位、默认值、扫描/写入转换；C2 必须把特殊列（例如 `locked_until`、`applied_at`、`updated_at DEFAULT 0`）逐列落盘。

## 方案边界

本合同只扩展 C1 时间列形状，不重开 VP-013 的 Store 端口、事务、checksum 算法、非时间整数宽度与跨引擎 residual。
