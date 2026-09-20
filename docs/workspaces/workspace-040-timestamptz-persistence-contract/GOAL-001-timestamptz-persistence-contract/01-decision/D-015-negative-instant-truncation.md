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
- 旧 DB seconds/milliseconds 是整数来源：秒值与毫秒值本身没有微秒以下余数，PG conversion 使用 `date_trunc` + 整数 interval 只负责形成 `timestamptz(6)`，不会以 typmod cast 对任意 fractional input 做 round。
- C2 禁止把 raw fractional SQL/driver value 直接 cast 到 `timestamptz(6)` 并宣称已满足向零规则；raw fractional 输入必须先经过 Go codec。
- 负 epoch 不是 sentinel；仅 D-012 明确的 0 sentinel 映射为 NULL。
