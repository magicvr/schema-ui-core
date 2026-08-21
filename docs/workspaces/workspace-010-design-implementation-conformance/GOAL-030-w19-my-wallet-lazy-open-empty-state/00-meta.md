---
id: GOAL-030-w19-my-wallet-lazy-open-empty-state
title: W19 · 我的钱包惰性开通与未开户空态
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-18
updated: 2026-08-18
version: 0.2.0
progress: 4/4
---

# GOAL-030 · W19 · 我的钱包惰性开通与未开户空态

VP-010 / workspace-010 的**第十九波**：修正「我的钱包」新用户必须手动开通、开通前整页报错。不重开 workspace-011 GOAL-020/022，也不回退 W15-F11 的 GET 只读。

## 当前边界

- **范围**：进「我的钱包」由前端对 `POST /api/wallet/me` 惰性开通；`WALLET_NOT_FOUND` 在余额卡/流水表按空态而非硬错误；开通失败才露出重试 CTA。
- **非范围**：改回 GET 写库；签名/公开钱包 URL；管理端 `/api/wallet/accounts` 开户流；资金操作自服务。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：D-001。
- [x] **S2 · 实施**：`wallet-ensure` + 空态（E-002）。
- [x] **S3 · 定向验证**：Web 73/73 + `tsc`（E-003）。
- [x] **S4 · 自审与关门**：A-001 self pass；goal-tree / workspace 同步（E-004）。

progress: 四个等权检查点；当前 **4/4**。

## 审计策略

S4 关门 `self`（可逆 UX；不改 GET 写语义）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required | 惰性开通走 POST 是否仍满足 W15-F11 | S1 | S1 | 对照 GOAL-020 D-001 / W15-F11 | **verified** | D-001：GET 仍只读；开通仅 POST |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 溯源

- 用户 2026-08-18：新用户手动开通 + 开通前报错不符合预期
- W15-F11 / A-004 F-001：GET 只读后首方面用 toolbar「开通钱包」补洞，空态未做
- workspace-011 GOAL-020/022：用户钱包本应惰性自动开户
