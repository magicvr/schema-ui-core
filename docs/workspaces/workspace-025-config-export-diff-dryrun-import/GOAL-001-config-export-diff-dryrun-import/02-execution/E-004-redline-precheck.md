---
status: recorded
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-001-config-export-diff-dryrun-import
version: 0.1.0
---

# E-004 · 红线核账预检（2026-08-30 · R4 越界核账前置）

`git diff --name-only cf68c7ce..HEAD`（workspace-025 开区至 R3 C1 全量提交面）：

1. **代码变更面**（仅 4 文件）：`apps/api/cmd/schema-ui/configpkg.go`（新 · 功能本体）、`configpkg_test.go`（新 · 测试）、`main.go`（config 子命令注册 + cliError）、`apps/api/server/config.go`（仅新增只读导出 `DefaultConfigYAML()`）。
2. **红线域零触碰**：`apps/api/internal/store`（迁移台账）/ `kernel/profile.go`（Profile 默认集与装配）/ `apps/web/src/protocol/upstream`（provenance）/ 任何 `migration`/`migrate` 面 —— 均无变更。
3. **结论**：VP-008 `go` 红线维持（未改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义）；迁移台账 checksum 面零变更；密钥 fail-closed 装载语义未改。R4 越界核账的直接证据已预积累。