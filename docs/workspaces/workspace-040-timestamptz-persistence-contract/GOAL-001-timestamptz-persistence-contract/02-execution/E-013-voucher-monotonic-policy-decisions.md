---
id: E-013-voucher-monotonic-policy-decisions
doc: execution-entry
status: recorded
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-013 · Voucher 与微秒单调策略用户裁决

用户选择 voucher legacy 0→NULL、负值 fail closed、正值 Unix seconds；users/roles `updated_at` 采用 `max(truncatedNow, old+1µs)` 保持单调。C2 仍需把这两项转为逐列 migration/runtime tests。

证据：Root D-012 / D-013。
