---
id: GOAL-001-wallet-prepaid-instrument
doc: execution-entry
record_id: E-004
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-004 · A-002 修复复审与 VP-029 关门（2026-09-02）

## 事实

1. 用户指令：复审 A-002 修复；没问题则关闭对应 VP。
2. 独立复审 **A-004 pass**（不以 A-003 声明为证据）：
   - F-001：协议页生成成功后同手势 CSV 下载。Vitest `automatically triggers CSV download when voucher generation returns plaintext codes` PASS。
   - F-002：非 CNY 生成/核销 fail-closed。`TestRedeemCurrencyMismatchFailClosed` + HTTP USD 400 PASS。
   - F-003：`GetOrCreateSubjectAccountInTx` 使用 `ON CONFLICT DO NOTHING`。`TestPostgresWalletVoucherAndSubject0064` **PASS 1.87s（未 skip）**。
3. `/vision` 关闭 VP-029：`active → closed` v0.3.0；VRev-067 self `pass`；VR-061 editorial。组合焦点回到无 active 交付 VP（持续程序 VP-009/010；VP-030/031 仍 planned）。

## 产物

- `03-audit/A-004-a002-finding-closure-independent.md`
- `docs/vision/reviews/VRev-067-vp029-wallet-prepaid-instrument-close-out.md`
- `docs/vision/plans/VP-029-wallet-prepaid-instrument.md` v0.3.0 `closed`
