---
id: GOAL-001-distribution-package-pilot
title: 分发形态包化试点（构建期包消费路径）
status: active
parent: null
created: 2026-08-29
updated: 2026-08-29
version: 0.2.0
progress: 1/5
plan_refs:
  - VP-022-distribution-package-pilot
primary_plan: VP-022-distribution-package-pilot
serves_summary: 以证据驱动验证「构建期包消费」分发路径（Go 库 + npm 包组 + 零冲突升级演练），产出 go/no-go 报告
---

# GOAL-001 · 分发形态包化试点

## 概述

承接 [VP-022-distribution-package-pilot](../../vision/plans/VP-022-distribution-package-pilot.md)（active v0.3.0）：在既有 fork 消费路径之外，验证「构建期包消费」最小闭环——本仓以 Go 库模块 + npm 包组形态发布 kernel / 标准模块 / Renderer / Shell，空下游仓以 `go get` / `pnpm add` 组装并完成一次零冲突升级演练；产出实测对比与 go/no-go 报告（是否推进 Charter strategic 修订由报告结论再议）。**试点不改 Charter、不弃 fork**。

## 成功标准（对应 VP-022 六条退出判据）

- [ ] 退出判据 #1：Go 库消费闭环（`go get` + 自建组合根 + kernel ≥1 标准模块 + 功能基线等价）
- [ ] 退出判据 #2：Web 包消费闭环（npm 包组组装 + 同一 schema 页面集 + Token 覆盖定制）
- [ ] 退出判据 #3：零冲突升级演练 PASS（真实演进 → 仅 bump 版本 + changelog 迁移说明 → 回归全绿、冲突 0、无 merge）
- [ ] 退出判据 #4：契约冻结面落盘（kernel/模块/npm 包 semver + breaking 流程 + changelog 模板）
- [ ] 退出判据 #5：发布可复现（一键 Go tag + npm 包组 + golden consumer 消费回归）
- [ ] 退出判据 #6：go/no-go 报告（实测对比 + Charter strategic 修订建议，按 VP 触发框架判向）

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标（R2/R3 各一个实施子目标；R5 可拆证据与报告）。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 契约冻结面落盘：kernel 公共 API / 模块契约 / npm 包组 semver + breaking 流程 + changelog 模板；「冻结面 vs 内部自由演进面」分界；发布通道初选 | **已关门**（2026-08-29 · GOAL-002 done 4/4 · 用户确认清单 v1.0.0 生效 · D-002；F-001/F-002 交接 R2） |
| R2 | Go 库包闭环：空下游仓 `go get` + 组合根装配（kernel + ≥1 标准模块），功能基线等价证据 | 依赖 R1 | **进行中**（GOAL-003 active：internal 阻断实验 done（E-001）· 外移方案待用户裁决（D-001）） |
| R3 | Web 包闭环：npm 包组组装 + schema 页面渲染 + Token 覆盖定制证据 | 依赖 R1（可评估与 R2 并行） |
| R4 | 零冲突升级演练：上游真实演进样本 → 下游仅 bump + 迁移说明 → 回归全绿、冲突 0、无 merge | 依赖 R2/R3 |
| R5 | 证据与 go/no-go：发布可复现（脚本/CI + golden consumer）+ 实测对比报告 + Charter 修订建议 | 依赖 R1–R4 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | kernel 公共 API 冻结面清单（哪些包/符号成为对外契约；与模块契约六贡献的导出边界） | R1 冻结 / R2/R3 实施 | R1 | 扫描 `apps/api/kernel` 导出面 + module-architecture 契约 → 形成冻结面清单 | open | — | 待确认 |
| I-002 | required | npm 包拆分与 peer 依赖策略（protocol/renderer/shell/ui 边界；React/Tailwind 版本耦合矩阵） | R3 实施 / R5 发布 | R3 | 包边界设计 + peerDependencies 草案 + 拆包可行性验证 | open | — | 待确认 |
| I-003 | required | 发布通道（npm 私有 registry vs GitHub Packages；Go 版本 tag 策略与 module proxy） | R5 发布可复现 | R5 | 通道调研 + 脚本试点 | open | — | 待确认 |
| I-004 | non-blocking | 零冲突演练的「真实演进样本」选择（哪类上游变更最能代表冲突压力：配置键 / 迁移 / 依赖 / 公共 API） | R4 演练质量 | R4 | 样本设计（≥3 类变更） | open | — | 待确认 |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001 已首条落盘，后续按编号递增。

## 备注

- 审计模式（D-001 已定）：阶段关门 default self；R5 发布/兼容门禁与 Root 关门 = independent（grok build 先例，项目级默认执行路径）。
- freshness 三字段（V-F084 执行）：见 `01-decision/D-001-workspace-root-establishment.md`。