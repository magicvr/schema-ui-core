---
id: D-001
goal: GOAL-030-w19-my-wallet-lazy-open-empty-state
title: W19 方案冻结：我的钱包惰性开通与未开户空态
status: accepted
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: GOAL-001-design-implementation-conformance
---

# D-001 · W19 方案冻结

## 1. 触发

用户确认：新用户必须手动开通、开通前「我的钱包」报错，两者都不符合预期。W15-F11 把 GET 改为只读是对的；当时用 toolbar「开通钱包」补首方面，没有空态，也没有进页自动 POST。

## 2. 决定

1. **惰性开通（恢复产品意图，遵守 GET 只读）**：`my-wallet` 页挂 custom component `wallet-ensure`，挂载时 `POST /api/wallet/me`（幂等 get-or-create），成功后 `reloadList` 刷新余额卡与流水。
2. **空态（竞态/失败兜底）**：`WALLET_NOT_FOUND` 在 statCard / 流水表不当硬错误——余额卡显示 0，表格空列表。
3. **CTA**：去掉表格上常驻「开通钱包」。仅 `wallet-ensure` 在 POST 失败时显示错误 + 重试。
4. **不改 API**：`GET /me` 与 `GET /me/entries` 仍只读 404；创建仍只走 POST。

## 3. 未选方案

| 方案 | 未选理由 |
|------|----------|
| GET 再改回 get-or-create | 违反 W15-F11 |
| 只做空态、仍要人手点开通 | 违背 GOAL-020/022 惰性开户 |
| 只做进页 POST、404 仍当报错 | 开通完成前三张卡+表会闪错误 |

## 4. 后续

S2 实施；S3 Web 定向 + tsc。go 不暂挂。
