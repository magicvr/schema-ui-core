---
id: GOAL-013-w12-product-surface-intent
doc: audit-entry
record_id: A-001
source: self
scope: S3 实施 ～ S4 验证/关门
verdict: pass
status: recorded
auditor: self（编排器）
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# A-001 · S3/S4 自审（2026-08-16）

- **source**：self
- **auditor**：编排器（workspace-010 波次流程）
- **类型** / **scope**：close-out · S3 实施（T-05/T-01/T-03/T-02/T-06）～ S4 验证/关门
- **verdict**：**pass**（scope 内无 required findings）

## 范围与区间

- covered：D-002～D-008 冻结方案 vs 实现逐条核对；E-002～E-004 主张 vs 代码证据；回归证据（Go 全量 + Web 全量 + tsc）。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| T-05：deletedAt/restoredAt 为 UTC ISO-8601 | `recyclebin.go` `recycleItemToMap` 格式化；`recyclebin_test.go` 新增 `time.Parse` 断言；存储层不变 |
| T-01：全断点单一用户下拉；菜单序 = projection.user 声明序 + 退出末项 | `App.tsx` `UserMenu`；`user-menu.test.tsx` 断言菜单项 ["Account","Settings","Sign out"]；抽屉不再渲染 user 链 |
| T-03：account.json 三档 Tabs；TabsView 支持 labelKey | `account.json` body.tabs 三 section；`render.tsx` TabsView `resolveTextProp("labelKey","label")`；D-VAL 全量结构回归绿 |
| T-02：搜索表单非 q 字段绑定 filters；清空即移除；表级 filters 保留 | `render.tsx` `searchFormSubmit`（owned 字段先删后写）；`search-form-filters.test.tsx` 断言 `q=ali&enabled=true` |
| T-02：12 页 schema 矩阵落实（含降级项） | 11 页改造 + notifications 新增；wallet-entries 按 D-003 降级为仅关键词（list 无 entryType 谓词，E-003 留痕） |
| T-02：后端筛选接线 | users/roles ExtraQuery + *bool filter + where；scheduled-tasks enabled；task-runs status；notifications q/read；wallet/recycle-bin 既有参数 |
| T-06：模块启用只认 YAML；preset/list 互斥；env 选择器废除 | `config.go` `resolveModulesFromYAML`/`loadPresetFile`；`config_test.go` 7 子用例（含旧 env 被忽略）；compose/dev.cmd/文档同步 |
| T-06：默认三档模块集未变（go 判定依据） | `kernel/profile.go` `profileDefaults` 无改动；config_test 内置预设解析结果与既有 `TestBuiltinProfilesResolveDeterministically` 一致 |
| 回归 | `go test ./... -count=1` 0 FAIL；`npx vitest run` 1027/1027；`tsc -b` 0 |

## 对照成功标准

- S3 检查点：P0（T-05+T-01）✓、P1（T-03+T-02）✓、P2（T-06）✓ —— E-002/E-003/E-004 逐项留痕。
- S4 检查点：回归绿 ✓；自审 A-001 ✓；独立审计 A-002（grok）另行核销；T-06 go 判定见 E-005（部署契约变化、默认集不变 → 不暂挂）。

## Findings

无 required findings。非阻断观察：
- R-001（recommended）：wallet-entries 的 entryType 筛选在账本 list 增加谓词后可再挂（D-003 既有降级条款，无需重开目标）。

## 结论

S3/S4 主张均可核对，回归全绿；等待 A-002 independent 核销后关门。
