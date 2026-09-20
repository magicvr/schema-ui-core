---
id: D-012-v73-allocation-negative-truncation
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-012 · v73 allocation / negative truncation 承接

承接 **Root** `D-014`/`D-015`：v73–v87 是 R2 当前 baseline，未发布前可有记录地拆分/调整；发布后 append-only。负 fractional `time.Time` 由 Go codec 向零截断；legacy integer sec/ms 的 PG conversion 一律 `date_trunc` 包裹，不依赖 typmod rounding：**毫秒族**用 `TIMESTAMPTZ 'epoch' + <col> * INTERVAL '1 millisecond'`（整数 interval），**秒族保留** `to_timestamp(<col>::double precision)`（2026-09-20 用户裁决 B；见 **Root** `D-015`）。

> **编号限定（A-029 F-I-018）**：本条编号 `D-012` 属 **child**（v73 allocation / 负瞬间承接），与 **Root** `D-012-voucher-invalid-value-policy.md`（voucher 异常值政策）同号不同义；child 没有 D-013/D-014/D-015，这三个编号只存在于 Root。引用 voucher 政策时必须写 **Root D-012** + child `D-011`。
