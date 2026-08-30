---
status: recorded
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-004-r3-dryrun-import
version: 0.1.0
---

# E-002 · C1 dry-run 实现与测试（2026-08-30）

1. **代码**：`apps/api/cmd/schema-ui/configpkg.go` 追加 —— `parsePackageStrict`（KnownFields 严格解码 + 多文档拒绝 + format/version 检查）、`flattenTree`、`dryRun`（三件套：结构校验 → `secrets.exclude` env fail-closed → 包→目标方向影响差量）、`cmdConfigDryRun`（手写参数解析 · 退出码 0/1/2 · yaml/json）；`main.go` usage 行更新（dry-run 转正，import 仍标注 R3）。
2. **测试**：`configpkg_test.go` 追加 4 用例（TestDryRunPass：checks 全 ok + 零变更 + 目标文件快照零副作用；TestDryRunEnvMissingFailsClosed：env 缺失 → fail check + 命令层 exit 1；TestDryRunChanges：`-config` 目标 modify http.addr + old/new 方向正确；TestDryRunInvalidPackage：坏包 exit 1 + 缺文件 exit 2）——`go test ./cmd/schema-ui/` PASS（14 用例）。
3. **回归**：`go test ./...`（apps/api 全量 49 包）PASS。
4. **CLI 冒烟（实证）**：`config export` → `config dry-run <pkg>` exit 0（checks ok · changes []）；`-f json` 双面一致；env 缺失 → exit 1（`2 check(s) failed (fail-closed)`）。
5. **验证覆盖判据**：VP-025 判据 #3（dry-run 无副作用）交付面完成（快照对比 + 零写断言）；判据 #4 待 C2（import · I-025-004 裁决）。