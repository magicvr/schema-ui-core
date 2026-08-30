---
id: GOAL-001-config-export-diff-dryrun-import
title: 配置包导出 / diff / dry-run / 导入
status: done
parent: null
created: 2026-08-30
updated: 2026-08-30
version: 0.5.0
progress: 4/4
plan_refs:
  - VP-025-config-export-diff-dryrun-import
primary_plan: VP-025-config-export-diff-dryrun-import
serves_summary: 配置包导出 / diff / dry-run / 导入（Admin 功能分支 · 基架能力剩余 #3）：可移植配置包 + 键级差量 + 只读预检 + 安全导入
---

# GOAL-001 · 配置包导出 / diff / dry-run / 导入

## 概述

承接 [VP-025-config-export-diff-dryrun-import](../../vision/plans/VP-025-config-export-diff-dryrun-import.md)（active v0.2.0 · [VRev-054](../../vision/reviews/VRev-054-vp025-activation.md) self `pass` · Admin 类 freshness PASS `c9122478`→`055da2fd`）：把「配置包导出 / diff / dry-run / 导入」收成可核对的 Admin 合同。**对象面 = serve 壳配置树**（`apps/api/server/config.default.yaml` 内嵌默认 · env 插值 `$VAR` fail-closed / `$VAR:-default` · 骨架模板 `config.yaml.tmpl`）。红线（激活即生效）：不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；密钥 fail-closed；热加载不进分母。

## 成功标准（对应 VP-025 六条方向级退出判据）

- [x] 判据 #1（导出闭环）：当前生效配置可导出为可移植配置包；往返（导出 → 干净实例导入 → 再核对）一致，密钥/敏感值按冻结规则排除或脱敏（快测 + 至少一条 harness/CLI 实证）——R2（导出面交付；跨实例导入往返实证归 R3）
- [x] 判据 #2（diff 可核对）：两包 / 包 vs 运行配置的差量输出可机器读并可断言（快测覆盖一致、仅差、冲突场景）——R2（TestDiff* · CLI 冒烟 exit 0/1）
- [x] 判据 #3（dry-run 无副作用）：预检覆盖校验与影响报告，成功/失败路径均有快测，不产生写副作用——R3（TestDryRunPass 快照对比 · 冒烟 exit 0/1）
- [x] 判据 #4（导入不破坏）：预检通过后应用；导入前后实例可启动、回归快测通过；失败路径不破坏既有配置（快照/回滚语义按 I-025-004 冻结）——R3（用户裁决方案 A：备份+tmp+校验+rename · TestImport* 失败路径目标原样）
- [x] 判据 #5（边界保持）：未改 Charter；未改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；热加载不进分母；密钥 fail-closed 保持——全程（R1/R2 已核对）
- [x] 判据 #6（审计闭合）：开放 required finding = 0（或已合法闭合）——R4（A-001 self `pass` · A-002 grok build independent → F-001～F-008 全部 fixed · 开放 required=0 · VRev-055 `pass` · **2026-08-30 用户书面确认 VP-025 closed v0.3.0**）

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 合同冻结（判据 #5/6 边界）：配置包内容边界（I-025-001）· 落地形态（I-025-002）· diff/dry-run 语义基线 | **已关门**（2026-08-30 · GOAL-002 done 3/3 · A-001 self `pass` · 配置包合同 v0.1.0 冻结（D-002）；I-025-001/002/003 `verified`；I-025-004 预告 R3） |
| R2 | 导出 + diff（判据 #1/2） | 依赖 R1 | **已关门**（2026-08-30 · GOAL-003 done 3/3 · A-001 self `pass`（0 required）· configpkg.go 实现 · 10 用例 + 全量 49 包 + CLI 冒烟 exit 0/1 实证；判据 #1/#2 本阶段交付） |
| R3 | dry-run + 导入（判据 #3/4；I-025-004 前置裁决） | 依赖 R2 | **已关门**（2026-08-30 · GOAL-004 done 3/3 · A-001 self `pass`（0 required）· dry-run 三件套 + import 方案 A（用户裁决）· 18 用例 + 全量 49 包 + 往返冒烟 diff `[]`；判据 #3/#4 交付） |
| R4 | 证据与关门（判据 #6） | 依赖 R1–R3 | **已关门**（2026-08-30 · GOAL-005 done 3/3 · 证据矩阵 verified · A-001 self `pass` + A-002 grok build independent（F-001～F-008 全部 fixed · 开放 required=0）· VRev-055 `pass` · **用户书面确认 VP-025 `closed` v0.3.0**） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-025-001 | required | 配置包内容边界：包 = 当前生效配置树的哪些键；env 引用（不解析 vs 解析后值）；密钥/敏感值处理（排除 / 脱敏 / 占位 + fail-closed）。 | 方案冻结 + 退出判据 #1 | R1 | 用户裁决（R1 合同冻结前置；VRev-054 V-F090 登记） | **verified** | — | 2026-08-30 用户裁决：非敏感结构键全集 + env 保留 `$VAR` 形态 + 敏感键占位/`secrets.exclude` 清单 + 导入 fail-closed（GOAL-002 D-001 accepted · 合同 §1/§3） |
| I-025-002 | required | 落地形态：CLI（`schema-ui config *`）vs 管理面 vs 两者；与 VP-007 Settings 面的关系。 | 方案冻结 | R1 | 用户裁决（R1 合同冻结前置） | **verified** | — | 2026-08-30 用户裁决：**CLI 主路径**四子命令（export/diff/dry-run/import）；管理面不做；yaml/json 双格式；diff 机器可读（GOAL-002 D-001 accepted · 合同 §2） |
| I-025-003 | non-blocking | diff 语义与输出：键级规范化/排序/类型；输出格式（text / yaml / json 合一或分面）。 | 退出判据 #2 | R2 | lead 建议 + 用户确认 | **verified** | — | 随裁决确认（用户采纳建议）：规范化键级差量 + yaml/json 双输出 + 退出码 0/1/2（GOAL-002 D-001 accepted · 合同 §2.2） |
| I-025-004 | required | 导入失败语义：预检失败即止 vs 应用期失败快照回滚；与既有升级前快照（VP-013 方言级 / VACUUM INTO）的关系。 | 退出判据 #4 | R3 | 用户裁决（R3 前置） | **verified** | — | 2026-08-30 用户 GUI 裁决**方案 A**：原子替换 + 应用前备份 `.pre-import.bak`；失败路径清 tmp、原文件不被触碰；不引入 DB 级快照（GOAL-004 D-002） |
| I-025-005 | required | 是否触及 Profile 默认集 / 模块矩阵 / Manifest 装配？**本 VP 冻结为不进**（VP-008 `go` 红线）。台账投影。 | 退出分母 | R1 | 投影（V-F090 / D-001） | **registered**（冻结不进） | — | VP-025 §边界 |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001 已首条落盘，后续按编号递增。

## 备注

- 审计模式（D-001 已定）：阶段关门 default self；实证门禁（R4 证据 / 关门）可按需 independent（grok build 先例，项目级默认执行路径）。
- freshness 三字段（V-F089 执行）与现状锚点（V-F091 执行）：见 `01-decision/D-001-workspace-root-establishment.md`。
- I-025-001 / I-025-002 已于 R1 合同冻结（2026-08-30）经用户裁决 `verified`（GOAL-002 D-001）；I-025-004 为 R3 前置裁决点。