---
status: accepted
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-001-config-export-diff-dryrun-import
version: 0.1.0
---

# D-003 · R2 关门投影（2026-08-30）

1. **R2 已关门**：GOAL-003-r2-export-diff `done 3/3`——C1 export 实现（configpkg.go · buildExportTree/buildPackage）+ C2 diff 实现（diffLeafMaps/loadConfigLeaf · cliError 0/1/2 · 手写参数解析）+ C3 A-001 self `pass`（0 required）。
2. **测试证据**：10 用例全绿（export 形态/无明文/覆盖合并/JSON/diff 一致·modify·增删·--against/错误码/往返）；`go test ./...`（apps/api 全量 49 包）PASS；CLI 端到端冒烟（export exit 0 · diff `[]` exit 0 · `--against` modify exit 1）。
3. **Root progress**：R2 检查点关闭 → `2/4`（R3 dry-run+导入（I-025-004 前置裁决）→ R4 证据与关门）。
4. **VP-025 判据映射**：判据 #1（导出闭环）/ #2（diff 可核对）交付面完成（跨实例往返实证归 R3 判据 #4 验收）；判据 #5/#6 维持。
5. **红线**：Profile 默认集 / 模块矩阵 / Manifest 装配零触碰；迁移台账零变更；密钥 fail-closed 与装载语义未改；`server` 仅新增只读 `DefaultConfigYAML()`。
6. **I-025-004**：仍为 R3 前置裁决点（合同 §2.4/§7 边界预告；快照/回滚语义待用户裁决）。