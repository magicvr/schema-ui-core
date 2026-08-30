---
status: accepted
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-001-config-export-diff-dryrun-import
version: 0.1.0
---

# D-005 · R3 关门投影（2026-08-30）

1. **R3 已关门**：GOAL-004-r3-dryrun-import `done 3/3`——C1 dry-run（结构校验 + env fail-closed + 影响报告）+ C2 import（**I-025-004 用户裁决方案 A**：备份 `.pre-import.bak` → tmp 写 → `LoadConfig` 装载校验 → `os.Rename` 原子替换；失败清 tmp、原文件不被触碰）+ C3 A-001 self `pass`（0 required）。
2. **测试证据**：cmd/schema-ui 18 用例全绿；`go test ./...`（apps/api 全量 49 包）PASS；CLI 端到端冒烟（import exit 0 · 往返 re-export diff `[]` · 二次导入 `.bak` 生成 · fail-closed 拒绝 exit 1）。
3. **Root progress**：R3 检查点关闭 → `3/4`（R4 证据与关门待启动）。
4. **VP-025 判据映射**：判据 #3/#4 交付面满足；判据 #1 跨实例往返闭环（冒烟实证）；判据 #2/#5 维持；判据 #6 归 R4。
5. **信息项**：I-025-004 → `verified`（用户 2026-08-30 GUI 裁决方案 A · GOAL-004 D-002）；I-025-005 `registered` 维持。
6. **红线**：延续 E-004 核账；R3 代码面仍在 `cmd/schema-ui`（configpkg.go）——store/profile/upstream/迁移零触碰；密钥 fail-closed 未改。