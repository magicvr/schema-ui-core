---
id: GOAL-004-r3-dryrun-import
title: R3 dry-run + 导入（配置包预检 / 安全应用）
status: active
parent: GOAL-001-config-export-diff-dryrun-import
created: 2026-08-30
updated: 2026-08-30
version: 0.2.0
progress: 1/3
plan_refs:
  - VP-025-config-export-diff-dryrun-import
primary_plan: VP-025-config-export-diff-dryrun-import
serves_summary: 承载 VP-025 R3 阶段：config dry-run（只读预检：结构校验 + secrets.exclude env fail-closed + 影响报告）与 config import（预检通过后安全应用；失败路径不破坏既有配置 per I-025-004）。
---

# GOAL-004 · R3 dry-run + 导入

## 概述

执行 Root 纲领 **R3**（VP-025 判据 #3/#4），在配置包合同 v0.1.0（GOAL-002 D-002）§2.3/§2.4 之上实现：

1. **`schema-ui config dry-run <pkg>`**（C1 · 不依赖 I-025-004 裁决）：只读预检三件套——① 结构校验（包按 `configPackage` 严格解码：未知键/多文档/格式或版本不符 → 预检失败）；② `secrets.exclude` 逐项检查所需 env（未设置 → 预检失败，fail-closed，不回落默认）；③ 影响报告 = 相对目标配置的将变更键列表（add/modify/remove）。退出码 `0` 通过 / `1` 预检失败 / `2` 错误；yaml/json 双格式。
2. **`schema-ui config import <pkg>`**（C2 · **前置 = I-025-004 用户裁决**，2026-08-30 已呈报待裁决）：预检通过后安全应用到目标配置文件；失败路径不破坏既有配置（快照/回滚语义按裁决冻结）。
3. **边界（红线）**：不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go`）；不改迁移台账；密钥 fail-closed；热加载不进分母；管理面不做；不引入配置中心/远程分发。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **dry-run 实现**：三件套预检 + 退出码 0/1/2 + yaml/json + 快测（通过 / env 缺失 / 影响列表 / 坏包 / IO 错误） | **已关门**（2026-08-30 · parsePackageStrict/dryRun/cmdConfigDryRun · TestDryRun* 4 用例全绿 · CLI 冒烟 exit 0/1 + fail-closed 实证；无裁决依赖） |
| C2 | **import 实现**：预检接线 + 安全应用（I-025-004 裁决后冻结写入语义）+ 快测（往返 / 失败路径不破坏） | 待 I-025-004 用户裁决 |
| C3 | **审视与关门**：self 审计（合同↔实现 · 红线 · 证据）+ R3 关门、Root 信息台账回写 | 待 C1/C2 |

`progress` = 已关门检查点数 / 3。

## 成功标准（方向级 · 对应 VP-025 判据 #3/#4）

1. dry-run 前后目标配置零变更（快测快照对比）；零写副作用。
2. env 缺失场景 → 预检失败（fail-closed 快测）；坏包（未知键/格式不符）→ 预检失败。
3. 影响报告可机器断言（复用 diff 语义：包→目标方向的 add/modify/remove）。
4. import 预检通过后应用；导入前后实例可启动、回归通过；失败路径不破坏既有配置（语义 per I-025-004）。
5. 未越界：红线清单未触碰；`go test ./apps/api/...` 全绿。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-025-004 | required | 导入失败语义：预检失败即止 vs 应用期失败快照回滚；与既有升级前快照的关系。 | 判据 #4（C2 方案冻结） | R3 | 用户裁决（2026-08-30 已呈报 A/B/C 建议，待裁决） | open（待裁决） | 已到最晚阶段 | 合同 §2.4/§7 边界预告 |

C1（dry-run）不依赖 I-025-004，先行推进；C2 冻结前必须裁决。

## 父目标

- `GOAL-001-config-export-diff-dryrun-import`（Root · 纲领 R3）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。