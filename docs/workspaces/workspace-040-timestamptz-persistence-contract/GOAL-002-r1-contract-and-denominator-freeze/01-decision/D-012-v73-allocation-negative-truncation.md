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

承接 **Root** `D-014`/`D-015`：v73–v87 是 R2 当前 baseline，未发布前可有记录地拆分/调整；发布后 append-only。负 fractional `time.Time` 由 Go codec 向零截断；legacy integer sec/ms 的 PG conversion 不依赖 typmod rounding：

- **毫秒族（已更正，2026-09-20 PG 15/16/17 实测）**：`TIMESTAMPTZ 'epoch' + ((<col> - CASE WHEN <col> >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + (((<col> % 1000) + 1000) % 1000) * INTERVAL '1 millisecond'`（**整数拆分式，无浮点**）。
  > **⚠️ 原式已弃用**：本条原写 `TIMESTAMPTZ 'epoch' + <col> * INTERVAL '1 millisecond'`。A-044 在 PG 15.19 / 16.15 / 17.11 独立复现其在大数值上有 **+8 µs 误差**（自约 8×10¹³ ms 起），`::numeric` 无效。**不得再实现或引用原式。**
- **秒族（实测精确，保留）**：`date_trunc('microseconds', to_timestamp(<col>::double precision))` —— 2026-09-20 用户裁决 B；见 **Root** `D-015`。

> 完整更正记录（含实测表与受影响区间）见 child `GOAL-002` 的 `attachments/r1-c2-per-table-pg-ddl-v1.0-fc.md` §0.1（A-044 要求：Root `D-015` 与 child `D-012` 亦须标弃用，本轮已办）。

> **编号限定（A-029 F-I-018）**：本条编号 `D-012` 属 **child**（v73 allocation / 负瞬间承接），与 **Root** `D-012-voucher-invalid-value-policy.md`（voucher 异常值政策）同号不同义；child 没有 D-013/D-014/D-015，这三个编号只存在于 Root。引用 voucher 政策时必须写 **Root D-012** + child `D-011`。
