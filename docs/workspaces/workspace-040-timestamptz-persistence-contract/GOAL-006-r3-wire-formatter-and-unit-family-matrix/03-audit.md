---
id: GOAL-006-r3-wire-formatter-and-unit-family-matrix
doc: audit
status: active
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# 审计台账 · GOAL-006-r3-wire-formatter-and-unit-family-matrix（R3-A/B）

> 本文件是唯一正式审计台账索引：`self` 与 `independent` **共用** `A-NNN` 序列。
> 每条意见正文在 `03-audit/A-NNN-<slug>.md`；本文件登记条目头（`source`/日期/scope/`verdict`）。
> 独立审计默认只写意见，不修改 `status`/`progress`/方案正文；响应归编排器。

## 意见索引

| A-ID | source | 日期 | scope | verdict | 摘要 | 文件 |
|------|--------|------|-------|---------|------|------|
| — | — | — | — | — | 尚无意见（本目标刚立项） | — |

## 待复审事项（编排器登记，供独立审计取证）

| # | 事项 | 证据位置 | 说明 |
|--:|------|----------|------|
| 1 | ixed-6 输出替换范围是否覆盖 inventory 全部 Go 面、是否有遗漏的 inline 布局 | `internal/handler`、`cmd` | 逐项对照 `r1-public-wire-inventory-v0.1.md` §「当前 inline fixed-3 outputs」与「Current RFC3339 outputs」 |
| 2 | 输入 parser 兼容边界（`+00:00` 接受、非零 offset 拒绝、无时区拒绝）与 `D-005` 一致 | `internal/handler` parser + 测试 | 反例优先 |
| 3 | 单位族矩阵是否真的覆盖「秒/毫秒/可空/sentinel」四族各至少一个 endpoint | 测试证据 | 避免用同一族多次充数 |
| 4 | VP-020 会话时区展示 round-trip 是否证明「展示随会话时区、存储恒 UTC」 | Web 组件/单测 + Go 侧 | 载体按 `I-041-007` |
| 5 | 是否误把非 DB 文本纳入固定 6 位输出（`D-009` 例外） | `D-009` + fixture 对照 | 例外范围 |