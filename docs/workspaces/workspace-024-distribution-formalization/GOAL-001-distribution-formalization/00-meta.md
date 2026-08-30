---
id: GOAL-001-distribution-formalization
title: 分发形态正式化（cli+包 对外服务化收口）
status: done
parent: null
created: 2026-08-29
updated: 2026-08-29
version: 0.4.0
progress: 7/7
plan_refs:
  - VP-024-distribution-formalization
primary_plan: VP-024-distribution-formalization
serves_summary: 收口 VP-022/023 go 后残余：serve 壳 / npmjs 公开发布 / compose CI 实跑 / fork 对照计时 / 六包形态细化 / 迁移工具化 / 方法 B 置顶
---

# GOAL-001 · 分发形态正式化

## 概述

承接 [VP-024-distribution-formalization](../../vision/plans/VP-024-distribution-formalization.md)（active v0.2.0 · VRev-052 pass）：把「cli+包 分发路径」从机制已验证升级为对外正式化——8 条方向级退出判据（= go 后合并残余 7 项 + 方法 B 置顶与收口报告）逐条落地；不改 Charter（fork 与包消费并存维持）；npmjs 正式 scope/凭据属外部动作（R2 前置门禁 I-024-001，用户授权为界）。实验下游仓 = `github.com/magicvr/golden-field`（registry 语义消费）。

## 成功标准（对应 VP-024 八条退出判据）

- [x] 判据 #1：serve 壳闭环（`schema-ui serve` · HTTP 壳 + config 装载 + assembly 服务器面 · RT-D02 合同接线）——R1 已关门（GOAL-002 done 5/5 · A-001 self `pass`）；registry 级骨架消费随 R2 发布核销
- [x] 判据 #2：公开发布通道闭环（npmjs.com 六包 + CLI · golden-field 免凭据消费实证 · 发布流程成文）——R2 已关门（GOAL-003 done 4/4 · @magicvr/schema-ui-* 实发 · 五探针消费实证 · scope 定稿（用户裁决））
- [x] 判据 #3：compose CI 实跑（compose/Dockerfile + consumer-regression workflow 真实 CI 或等价实跑 PASS）——R3 已关门（GOAL-004 done 4/4 · grok independent `pass`（0 required）· harness A/B linux 容器 · hosted 触发登记 R7）
- [x] 判据 #4：fork 对照计时实验（同一演进集实测对比：耗时/冲突计数/契约迁移成本）——R4 已关门（GOAL-005 done 4/4 · A-002 grok `pass`（0 required）· v0.4.0 定量实证：冲突 1 vs 0 · 改写点 2 vs 0 · ≈13.2s vs ≈4.8s）
- [x] 判据 #5：renderer 依赖图 external 化（ui 包可消费 renderer）→ 冻结面 v1.4.0——R5 已关门（GOAL-006 done 5/5 · 187.5kB · 17 处包子路径 · peer 矩阵实发 · 五探针 + UI-ONLY 全绿）
- [x] 判据 #6：纯原子拆分（业务组件出 ui 包）——data-table 属 ui 设计系统面（用户 P-004 裁决 · 2026-08-29）；ui 独立消费实证（仅装 ui+peer）
- [x] 判据 #7：fork→包迁移工具化（`schema-ui migrate-fork` 或等价）——R6 已关门（GOAL-007 done 4/4 · A-002 grok `pass`（0 required）· 9510023 旧态实测：v0.3.0→v0.4.0 · npmjs 钉死 · build 绿）
- [x] 判据 #8：默认主路径宣告与收口报告（方法 B 置顶 · 核销表 · 残余清零）——R7 已关门（GOAL-008 done 4/4 · QUICKSTART 首段置顶 · closure-report 判据 8 条核销 · 残余四项登记/评述/候选 · Root 关门审计 A-002 grok `pass`）· **用户 2026-08-29 书面确认关门 → Root done 7/7 · VP-024 closed**

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | serve 壳闭环（判据 #1） | **已关门**（2026-08-29 · GOAL-002 done 5/5 · A-001 self `pass` · E2E-L1~L3 · 有界登记：信号 harness→R3 CI · registry 消费→R2） |
| R2 | 公开发布通道（判据 #2；前置门禁 I-024-001） | 依赖 R1 | **已关门**（2026-08-29 · GOAL-003 done 4/4 · A-002 grok fail → F-001~F-004 fixed · scope 定稿 @magicvr） |
| R3 | compose/CI 实跑（判据 #3） | 依赖 R1/R2 | **已关门**（2026-08-29 · GOAL-004 done 4/4 · A-002 grok `pass`（0 required）· harness A/B · hosted 触发登记 R7） |
| R4 | fork 对照计时（判据 #4） | 依赖 R2/R3 | **已关门**（2026-08-29 · GOAL-005 done 4/4 · A-002 grok `pass`（0 required）· 附件 fork-comparison-report） |
| R5 | 六包形态细化（判据 #5/#6） | 依赖 R1 | **已关门**（2026-08-29 · GOAL-006 done 5/5 · A-002 conditional → F-001~F-004 fixed · 判据 #5/#6 核销） |
| R6 | 迁移工具化（判据 #7） | 依赖 R2 | **已关门**（2026-08-29 · GOAL-007 done 4/4 · A-002 grok `pass`（0 required）· F-001~F-003 fixed） |
| R7 | 置顶与收口报告（判据 #8） | 依赖 R1–R6 | **已关门**（2026-08-29 · GOAL-008 done 4/4 · Root 关门审计 A-002 grok `pass`（0 required · F-001~F-006 fixed）· 用户书面确认 · VP-024 closed（VRev-053）） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-024-001 | required | npmjs 公开发布授权（正式 scope 名 + 凭据方式）；未授权时裁决降级方案 | 判据 #2 | R2 | 用户书面授权或裁决（V-F088 → 激活事务内登记） | **verified**（2026-08-29 用户裁决：@magicvr 定稿 · 真实发布 · token 注入点） | — | D-001-r2-publication §7 |
| I-024-002 | required | CI 槽位环境（真实 runner / 用户环境等价 + 凭据） | 判据 #3 | R3 | workflow 实跑验证 | **verified**（2026-08-29 · 本地等价 + linux 容器实跑；hosted 触发登记 R7 复核） | — | GOAL-004 E-002 |
| I-024-003 | required | fork 对照的同一演进集样本（V 演进选择） | 判据 #4 | R4 | 样本设计与 fork 基线 | **verified**（2026-08-29 用户裁决：v0.3.0→v0.4.0 真实演进集） | — | D-001 · GOAL-005 E-002 |
| I-024-004 | non-blocking | serve 壳与 assembly 服务器面接线方式（新壳 vs 扩展 assembly） | 判据 #1 设计 | R1 | 设计核对（复用 RT-D02 合同） | **verified**（2026-08-29 · D-001 定案：公开 `server` 面 = assembly 服务器面扩展；RT-D02 全序接线） | — | GOAL-002 D-001 |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001 已首条落盘，后续按编号递增。

## 备注

- 审计模式（D-001 已定）：阶段关门 default self；R2（release）与 R4（对照实验）实证门禁与 Root 关门 = independent（grok build 先例，项目级默认执行路径）。
- freshness 三字段（V-F087 执行）：见 `01-decision/D-001-workspace-root-establishment.md`。