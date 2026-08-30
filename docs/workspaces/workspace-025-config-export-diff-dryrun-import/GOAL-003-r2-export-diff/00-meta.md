---
id: GOAL-003-r2-export-diff
title: R2 导出 + diff（CLI config export / config diff）
status: done
parent: GOAL-001-config-export-diff-dryrun-import
created: 2026-08-30
updated: 2026-08-30
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-025-config-export-diff-dryrun-import
primary_plan: VP-025-config-export-diff-dryrun-import
serves_summary: 承载 VP-025 R2 阶段：实现 schema-ui config export（配置包 v1 产物 · 敏感剔除 · 无明文）与 config diff（规范化键级差量 · 退出码 0/1/2 · yaml/json · --against）。
---

# GOAL-003 · R2 导出 + diff

## 概述

执行 Root 纲领 **R2**（VP-025 判据 #1/#2）：在配置包合同 v0.1.0（GOAL-002 D-002）为分母之上，实现（上一目标只冻结合同、未动代码）：

1. **`schema-ui config export`**：导出 serve 壳配置树（内嵌默认 ∪ 显式文件）为配置包 v1（package 元数据 + config 非敏感键 + `secrets.exclude`）；env 引用保留 `${VAR}` 形态；敏感键按清单+宽规则剔除并记录所需 env；`-f yaml|json`；`-o` 文件或 stdout。
2. **`schema-ui config diff`**：两包之间 / 包 vs 配置（`--against`）的键级差量（add/modify/remove + 路径 + old/new）；忽略信息性元数据；退出码 0 无差 / 1 有差 / 2 错误。
3. `dry-run` / `import` 仅注册占位（R3 实现）。

**边界（红线）**：不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go`）；不改迁移台账；不改既有配置装载语义（server.LoadConfig 只新增只读导出函数）；密钥 fail-closed；热加载不进分母；管理面不做。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **export 实现**：导出树管线（默认∪显式合并 · 敏感剔除 · env 名提取）+ 包 v1 输出（yaml/json）+ 快测（含无明文断言） | **已关门**（2026-08-30 · configpkg.go buildPackage/buildExportTree · 测试 TestExport* 全绿） |
| C2 | **diff 实现**：扁平化键级比较 + 三退出码 + `--against` + yaml/json + 快测（一致/仅差/增删/忽略元数据） | **已关门**（2026-08-30 · diffLeafMaps/loadConfigLeaf + cliError 0/1/2 · TestDiff* 全绿 · CLI 冒烟 exit 0/1 实证） |
| C3 | **审视与关门**：self 审计（合同↔实现一致 · 红线核对 · 快测证据）+ R2 关门、Root 信息台账回写 | **已关门**（2026-08-30 · A-001 self `pass`（0 required）· Root 2/4） |

`progress` = 已关门检查点数 / 3。

## 成功标准（方向级 · 对应 VP-025 判据 #1/#2）

1. export 产物通过包 v1 结构校验（快测 schema 断言）；无敏感明文（快测扫描）。
2. 往返一致性支撑：export 产物重新解析 → 与源树逐键一致（快测）。
3. diff 可机器断言：一致 → 空差 + 退出码 0；单键改 → modify + 退出码 1；增删 → add/remove；yaml/json 双格式等价。
4. 未越界：红线清单未触碰；`go test ./apps/api/...` 全绿（含既有回归）。

## 信息就绪与未知项

R1 已关闭全部 C1 阶段信息项（I-025-001/002/003 `verified`）；本目标无新 required 信息项（I-025-004 仍属 R3 前置）。实现细节（如规范化键序、占位显示）以合同 D-002 为分母，lead 口径 + self 审计。

## 父目标

- `GOAL-001-config-export-diff-dryrun-import`（Root · 纲领 R2）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。