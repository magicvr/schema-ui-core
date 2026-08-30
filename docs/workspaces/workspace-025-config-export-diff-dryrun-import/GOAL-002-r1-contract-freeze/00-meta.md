---
id: GOAL-002-r1-contract-freeze
title: R1 合同冻结（配置包格式 / 敏感键规则 / CLI 命令面 / diff·dry-run 语义基线）
status: done
parent: GOAL-001-config-export-diff-dryrun-import
created: 2026-08-30
updated: 2026-08-30
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-025-config-export-diff-dryrun-import
primary_plan: VP-025-config-export-diff-dryrun-import
serves_summary: 承载 VP-025 R1 阶段：把 serve 壳配置树收成可核对的配置包合同——包格式 v1、敏感键排除与 fail-closed、CLI 四子命令面、diff/dry-run 语义基线、导入边界预告。
---

# GOAL-002 · R1 合同冻结

## 概述

本目标执行 Root 纲领 **R1**：在既有配置事实（`apps/api/server/config.default.yaml` 内嵌默认（`profile: admin`）· `server/config.go` 装载（env 插值 `$VAR` fail-closed / `$VAR:-default`）· 骨架模板 `config.yaml.tmpl`）之上，冻结 VP-025 合同正文——**配置包格式 v1 / 敏感键排除与 fail-closed / CLI 命令面（export·diff·dry-run·import）/ diff·dry-run 语义基线 / 导入边界预告（I-025-004 于 R3 冻结）**。合同正文 = GOAL-002 D-002 产物；不在此目标内实现 CLI 或改配置装载（红线约束）。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **信息裁决**：I-025-001（包键边界/密钥处理）、I-025-002（落地形态）两条 required 由用户裁决；I-025-003（diff 语义/输出）随 lead 提案确认 | **已关门**（2026-08-30 用户界面裁决全部采纳建议：非敏感结构键全集 + env 保留 `$VAR` + 敏感键占位/清单 + 导入 fail-closed；CLI 主路径四子命令 + 管理面不做 + yaml/json 双格式 + diff 机器可读——D-001 accepted） |
| C2 | **合同正文**：配置包格式 v1、敏感键规则、CLI 命令面（含退出码）、diff/dry-run 语义基线、导入边界预告落盘（D-002） | **已关门**（D-002 合同 v0.1.0 冻结，2026-08-30） |
| C3 | **审视与关门**：合同自审（self）+ R1 关门、Root 信息台账回写 | **已关门**（A-001 self `pass` 0 required；2026-08-30） |

`progress` = 已关门检查点数 / 3。当前 **3/3**（R1 已关门）。

## 成功标准（方向级）

1. 配置包格式 v1 明确：非敏感结构键全集 + env 引用保留形态 + 敏感键占位与 `secrets.exclude` 清单（永不出现明文）。
2. CLI 命令面四子命令（export / diff / dry-run / import）输入输出与退出码可核对（0/1/2 语义）。
3. diff 语义基线 = 键级规范化差量，yaml/json 双输出，机器可断言。
4. dry-run 语义基线 = 只读预检（结构校验 + 敏感键 env fail-closed 检查 + 影响报告），无写副作用。
5. 未越界：管理面本波不做；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；密钥 fail-closed；热加载不进分母；不引入配置中心/远程分发。

## 信息就绪与未知项

与 Root / VP-025 同号（I-025-001 ↔ 001，…）。C1 已关闭；I-025-004（导入失败语义）最晚 R3 冻结，本合同只作边界预告。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-025-001 | required | 配置包内容边界与密钥处理 | 方案冻结 + 判据 #1 | R1 | 用户裁决 | **verified** | — | 2026-08-30 用户裁决：非敏感结构键全集；env 引用保留 `$VAR` 形态不解析；敏感键导出为占位 + `secrets.exclude` 清单；导入 fail-closed（缺 env 即拒绝）（D-001 accepted；合同 §1/§3） |
| I-025-002 | required | 落地形态：CLI vs 管理面 | 方案冻结 | R1 | 用户裁决 | **verified** | — | 2026-08-30 用户裁决：**CLI 主路径**（`schema-ui config export/diff/dry-run/import`）；管理面本波不做（VP-007 Settings 不重开）；yaml/json 双格式 + diff 机器可读（D-001 accepted；合同 §2） |
| I-025-003 | non-blocking | diff 语义与输出格式 | 判据 #2 | R2 | lead 提案 + 用户确认 | **verified** | — | 随裁决确认（用户采纳建议）：规范化键级差量 + yaml/json 双输出 + 退出码 0/1/2（D-001；合同 §2.2） |
| I-025-004 | required | 导入失败快照/回滚语义 | 判据 #4 | R3 | 用户裁决（R3 前置） | **verified** | — | 2026-08-30 用户 GUI 裁决**方案 A**（GOAL-004 D-002 冻结 · 原子替换 + 应用前备份）；本合同 §7 当时仅边界预告 |

## 父目标

- `GOAL-001-config-export-diff-dryrun-import`（Root · 纲领 R1）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。