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
- **秒/毫秒两族的 PG 表达式（2026-09-20 用户裁决 B 收口，替代本决策原「sec/ms 均 `date_trunc` + 整数 interval」措辞）**：
  - **毫秒族**：`date_trunc('microseconds', TIMESTAMPTZ 'epoch' + <col> * INTERVAL '1 millisecond')` —— 整数 interval，全程不经二进制浮点（A-016 F-I-002.3）。
  - **秒族**：`date_trunc('microseconds', to_timestamp(<col>::double precision))` —— **保留**三份 C2 载体自 A-020 起已被 independent 接受的既有形式；不引入第二套秒列 SQL 家族（A-027 约束）。整数秒在 `double precision` 的精确整数范围内，不产生可观察精度损失。
- C2 禁止把 raw fractional SQL/driver value 直接 cast 到 `timestamptz(6)` 并宣称已满足向零规则；raw fractional 输入必须先经过 Go codec。
- 负 epoch 不是 sentinel；仅 **Root** D-012（voucher 异常值政策）明确的 0 sentinel 映射为 NULL。**child** `D-012-v73-allocation-negative-truncation.md` 是 v73 allocation / 负瞬间承接，与本条的 sentinel 表述无关（A-029 F-I-018 编号限定要求）。
