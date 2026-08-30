---
status: recorded
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-003-r2-export-diff
version: 0.1.0
---

# E-002 · export + diff 实现与测试（2026-08-30）

1. **代码**：`apps/api/cmd/schema-ui/configpkg.go`（新 · ~540 行）——cliError 退出码、cfgTree/package 结构（yaml/json 双 tag）、export 管线（buildExportTree：默认∪显式 · 敏感登记表 + 宽规则不变量 · env 名提取）、buildPackage（装载校验 + 元数据）、diff（loadConfigLeaf 扁平化 · diffLeafMaps · 手写参数解析 · --against）；`apps/api/server/config.go` 新增只读 `DefaultConfigYAML()`；`main.go` 注册 `config` 子命令 + `cliError` 处理 + usage。
2. **测试**：`apps/api/cmd/schema-ui/configpkg_test.go`（10 用例：默认形态/env 保留/敏感剔除与无明文/覆盖合并/JSON 双面/diff 一致·modify·增删·--against·错误码·往返一致）——`go test ./cmd/schema-ui/` PASS。
3. **回归**：`go test ./...`（apps/api 全量 49 包）PASS；`server` 3.2s PASS。
4. **CLI 冒烟（实证）**：`config export -o` exit 0（产物 package 元数据 + `${VAR}` 保留 + 无敏感键）；`config diff` 两包（exported_at 不同）→ `[]` exit 0；`config diff <pkg> --against <src>`（改 http.addr）→ 单条 modify exit 1。
5. **验证覆盖判据**：VP-025 判据 #1（导出闭环 · 快测 + CLI 实证）/ #2（diff 可机器断言 · 0/1/2 退出码）本阶段全部满足（往返「导出→再导出→diff 无差」已由 TestRoundtrip + 冒烟覆盖；跨实例导入往返属 R3）。