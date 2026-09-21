---
id: D-003-r1-public-wire-contract-user-decision
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-003 · R1 公共时间输出合同裁决

## 用户裁决

用户选择：VP-040 不仅改数据库 persistence，还将公共 API 时间输出统一为 **6 位微秒 RFC3339**（UTC `Z` 形态），而不是保持现有的固定 3 位毫秒或模块间混用的 RFC3339 输出。

## 影响边界

- 读面由 repository/store codec 返回 UTC `time.Time`，handler/serializer 统一输出 `YYYY-MM-DDTHH:MM:SS.ffffffZ`。
- 需要同步 API response formatter、schema/contract fixtures、Web/Go 测试与 VP-020 时区展示回归；不能把公共 wire 变更伪装成纯内部迁移。
- `pgtype`、驱动时间类型与 SQLite codec 仍不得泄漏到 handler/模块公共契约；公共面使用领域 `time.Time` 或已冻结的 JSON 字符串形态。
- 此裁决不改变 UTC instant、NULL/absence 或分母规则；不把 JSON/TEXT payload 内嵌时间自动纳入列分母。

## 门禁

这是 R1/C2 的公共面范围决策；C2 必须列出受影响 endpoints、formatter、协议 fixtures 与兼容回归，R3 必须执行微秒 wire 与 VP-020 时区往返验证。