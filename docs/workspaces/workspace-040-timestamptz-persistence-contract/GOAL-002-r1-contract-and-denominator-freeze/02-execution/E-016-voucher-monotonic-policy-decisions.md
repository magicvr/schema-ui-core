---
id: E-016-voucher-monotonic-policy-decisions
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-016 · voucher / monotonic policy 用户裁决

C2 policy facts：voucher legacy 0→NULL、负值 fail closed、正值 seconds；users/roles updated_at 微秒单调 `max(truncatedNow, old+1µs)`。实现/测试未开始。
