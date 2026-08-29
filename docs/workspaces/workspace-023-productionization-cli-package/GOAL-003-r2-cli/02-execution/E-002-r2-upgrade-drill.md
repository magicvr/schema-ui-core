---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-cli
version: 0.1.0
---

# E-002 · registry 升级演练（2026-08-29 · F-001 核销）

## 执行（CLI 生成仓 demo-admin · `schema-ui upgrade`）

| 环节 | 命令 | 结果 |
|------|------|------|
| Go 升级 | `go get github.com/magicvr/schema-ui-core/apps/api@latest` | ✅ `v0.1.0 => v0.2.0`（sumdb 本轮直接收录，未绕行） |
| npm 升级 | `pnpm add @magicvr/schema-ui-protocol@latest @magicvr/schema-ui-renderer@latest` | ✅ 0.2.0/0.1.0 registry 重拉 · lockfile 过供应链策略 |
| 探针回归 | `node web/probe.mjs` | ✅ PASS · 2.9 |
| 冒烟回归 | `go run ./cmd/server ./data-upgrade.db` | ✅ kernel=2.0.0 · fresh=true · contrib 与基线一致 |

**冲突计数 = 0 · 无 merge · 全程 CLI 驱动**（用户侧仅一条 `schema-ui upgrade`）。

## 知识/观察

- sumdb 对新 tag（apps/api/v0.2.0）本轮即时可用（首次 404 为偶发时延，按 E-001-R1 知识处理）。
- minReleaseAge 策略提示（pnpm 新特性）不影响零冲突结论，登记观察。