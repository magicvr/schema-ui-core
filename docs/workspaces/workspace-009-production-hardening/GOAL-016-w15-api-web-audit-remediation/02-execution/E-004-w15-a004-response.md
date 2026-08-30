---
title: E-004 · W15 A-004 响应实施（F-008/F-009 fixed + N-002 注释同步）
status: active
created: 2026-08-30
updated: 2026-08-30
parent: null
version: 0.1.0
---

# E-004 · W15 A-004 响应实施（F-008/F-009 fixed + N-002 注释同步）

日期：2026-08-30

## 实施事实（响应 A-003 新发现）

1. **F-008 · 邀请页 replaceState 回归锁 → fixed**：新增 `apps/web/src/components/invite-accept.test.tsx`（jsdom ×3：token 从 query 清理且保留其它参数 / 无 token 不触碰 history / 二次挂载幂等）。首次实现时 spy 顺序导致 setup 调用被计数，修正为「先设 URL 再 spy」后 3/3 pass。
2. **F-009 · server 配置负例假绿修正 → fixed**：`apps/api/server/config_test.go` 的 `TestLoadConfigInvalidShutdownTimeout`（`0s`/`-1s`）与 `TestLoadConfigDialectPairing` 自定义 YAML 显式 `app.env: development`，用例重新咬到目标分支（超时 ≤0 / 方言配对）。
3. **N-002 · 注释同步**：`server/config.go` `AppEnv` 行内注释由「"" = development（缺省）」改为「"" = 未声明（validate 拒绝，refusing to guess；W15 F-001）」。

## 回归验证

- API：`go test ./server/ -count=1` pass。
- Web：`vitest src/components/invite-accept.test.tsx` 3/3 pass；全量 vitest **90 files / 1186 tests pass**（含新增 3 例）。
- 说明：本条目改动仅测试与注释，不触及 S3～S5 实现切片 `609cd6d6` 的代码语义。