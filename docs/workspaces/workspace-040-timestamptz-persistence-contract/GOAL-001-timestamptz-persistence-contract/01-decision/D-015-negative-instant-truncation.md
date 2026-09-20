---
id: D-015-negative-instant-truncation
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-015 · 负时间与微秒截断解释

C2 解释用户 D-008 的“向零截断到微秒”如下：

- 新领域 `time.Time` 写入统一先在 Go codec `UTC().Truncate(time.Microsecond)`，包括 epoch 前的负 instant。
- 旧 DB seconds/milliseconds 是整数来源：秒值与毫秒值本身没有微秒以下余数，PG conversion 只负责形成 `timestamptz(6)`，不会以 typmod cast 对任意 fractional input 做 round。
- **秒/毫秒两族的 PG 表达式（2026-09-20 用户裁决 B 收口；毫秒族同日在 PG 15/16/17 实测后更正）**：
  - **毫秒族（已更正）**：`TIMESTAMPTZ 'epoch' + ((<col> - CASE WHEN <col> >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + (((<col> % 1000) + 1000) % 1000) * INTERVAL '1 millisecond'` —— **整数拆分式，全程不经浮点**。
    > **⚠️ 原式已弃用**：本决策原写 `date_trunc('microseconds', TIMESTAMPTZ 'epoch' + <col> * INTERVAL '1 millisecond')`。A-044 在 **PG 15.19 / 16.15 / 17.11** 三版本独立复现：该式在**大数值**上产生**非零误差**（`253402300799999` → `.999008`，**+8 µs**），误差自约 **8×10¹³ ms** 起出现，且部分 remainder 偏差达 ±16 µs；`::numeric` 强制转换**无效**。根因是 `BIGINT * INTERVAL` 内部经浮点。**该式不得再被实现或引用。**
  - **秒族（实测精确，保留）**：`date_trunc('microseconds', to_timestamp(<col>::double precision))` —— **保留**三份 C2 载体自 A-020 起已被 independent 接受的既有形式；不引入第二套秒列 SQL 家族（A-027 约束）。A-044 在 `253402300799` 与负值上实测**精确**。
  > 两族的**完整更正记录与实测表**见 child `GOAL-002` 的 `attachments/r1-c2-per-table-pg-ddl-v1.0-fc.md` §0.1。
- C2 禁止把 raw fractional SQL/driver value 直接 cast 到 `timestamptz(6)` 并宣称已满足向零规则；raw fractional 输入必须先经过 Go codec。
- 负 epoch 不是 sentinel；仅 **Root** D-012（voucher 异常值政策）明确的 0 sentinel 映射为 NULL。**child** `D-012-v73-allocation-negative-truncation.md` 是 v73 allocation / 负瞬间承接，与本条的 sentinel 表述无关（A-029 F-I-018 编号限定要求）。
