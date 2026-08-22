---
id: D-001
doc: decision-entry
goal: GOAL-005-default-single-key
status: accepted
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# D-001 · R4 判据 2 证据映射（默认单密钥仍可用）

## 证据映射表

| # | 判据面 | 证据（既有，R1～R3 已落盘） | 补充核对（本目标） |
|---|--------|------------------------------|---------------------|
| 1 | config 层：缺省 previous = 空字符串（单密钥模式） | `TestJWTSecretPreviousConfig/absent_previous_keeps_single-key_default`（R1） | 重跑确认 |
| 2 | config 层：生产缺 previous 不 fail-closed | `TestValidateProd/production_without_dev_session_passes`（仅 current，无 previous，R1 前已有 + R1 回归） | 重跑确认 |
| 3 | config 层：dev 低门槛不受 previous 影响 | `TestJWTSecretPreviousConfig/development_keeps_the_low_bar_even_with_a_weak_previous`（R1） | 重跑确认 |
| 4 | composition 层：生产装配路径空 previous = 单密钥语义 | `TestNewAuthenticatorWiresPreviousSecret` 第二段（空 previous 拒绝 old-key token，R2 A-002 F-001） | 重跑确认 |
| 5 | Compose 层：`AUTH_JWT_SECRET_PREVIOUS` 可选透传默认空；未设时配置合法 | R1 compose.yaml diff（结构证据） | **新增实证**：`docker compose config` 解析（设 dummy 必填项、不设 PREVIOUS）→ 输出中该键为空串 |
| 6 | dev 启动：本地快测路径不依赖 previous | `cmd/server` 既有启动/重启测试 + `resolveJWTSecret` dev fallback（[workspace-001] GOAL-005 D-004 既有合同，见 `cmd/server/main.go:74` 注释） | 重跑 server 包测试确认 |

## 判定标准

上表 6 行全部成立 ⇒ VP 方向级判据 2 达成（轮换不是 mvp/dev/production 启动硬依赖；缺省单密钥行为与 R1 之前逐字节一致）。

## 未选方案

- 为 R4 新写端到端"dev 双进程起服"脚本：`server` 包既有测试已覆盖启动路径，重复建设无增量价值。
- 修改 compose 默认注入 PREVIOUS：违反判据 2 本身。
