---
id: E-001
doc: execution-entry
goal: GOAL-005-default-single-key
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-001 · R4 证据实跑核对（2026-08-22）

按 GOAL-005 D-001 映射表逐行实跑，全部成立：

| # | 判据面 | 实跑结果 |
|---|--------|----------|
| 1 | config 缺省 previous = 空 | `TestJWTSecretPreviousConfig` 8/8 PASS（含 `absent_previous_keeps_single-key_default`） |
| 2 | 生产缺 previous 不 fail-closed | `TestValidateProd` 9/9 PASS（含 `production_without_dev_session_passes`） |
| 3 | dev 低门槛不受影响 | `development_keeps_the_low_bar_even_with_a_weak_previous` PASS |
| 4 | composition 空 previous = 单密钥语义 | `TestNewAuthenticatorWiresPreviousSecret` PASS（空 previous 拒绝 old-key token） |
| 5 | Compose 未设 PREVIOUS 合法且为空 | `docker compose config`（设 dummy 必填项、不设 PREVIOUS）→ 解析成功，输出 `AUTH_JWT_SECRET_PREVIOUS: ""` |
| 6 | dev/快测启动路径不依赖 previous | `go test ./cmd/server/ -count=1` ok（11.8s，含真实启动/重启循环） |

**判定**：VP 方向级判据 2 达成——未配置 previous 时本地/Compose 默认仍能开发与快测；轮换不是任何环境的启动硬依赖；缺省单密钥行为与 R1 之前一致。

环境注记：`docker compose config` 为纯解析校验（不构建、不起容器）；必填项以 dummy 值满足 fail-closed 插值，仅用于证明 PREVIOUS 缺省合法。
