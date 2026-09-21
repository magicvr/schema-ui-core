---
id: D-012-voucher-invalid-value-policy
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-012 · Voucher 时间异常值策略裁决

用户选择：`vouchers.expires_at` / `redeemed_at` 的 legacy `0` 转为 `NULL`；负值视为数据损坏并 fail closed；正值按 Unix seconds 转换。迁移预检必须分别统计 0、负值、正值；runtime 读取不能继续用 `>0` 把负值静默当 absence。
