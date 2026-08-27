---
id: D-001
title: 激活与开区：lead 绑定、纲领 R1～R3、信息门禁登记、架构类 freshness 留痕
date: 2026-08-27
status: accepted
---

# D-001 · 激活与开区（2026-08-27）

## 决定

1. **lead 绑定**：`workspace-021-graceful-shutdown-and-connection-drain` 为 VP-021 唯一 lead delivery 工作区（`workspace.md` 已建；`vision_role: delivery`；`plan_refs`/`primary_plan` = `VP-021-graceful-shutdown-and-connection-drain`）。
2. **纲领路线图**（P-001）：R1 合同冻结 → R2 实现与测试 → R3 证据与关门；按阶段逐项立项（R1 阶段先立 GOAL-002）。
3. **信息门禁**：VP-021 I-021-001～004 同号镜像登记进 Root 台账（I-001/I-002 required · 最晚 R1；I-003 required · 最晚 R2；I-004 non-blocking · R3）；R1 关闭前禁止直接改进程生命周期实现 / 迁移台账相关 DDL。
4. **架构类 freshness 留痕**（VP-008 `go` 消费有效性）：

| 字段 | 值 |
|------|-----|
| 原 `go` 候选 | `ed99e88`（VP-008 S5，2026-08-10） |
| 上次架构类复核 | `250cb9c`（VP-017 激活，VRev-038，2026-08-22） |
| 本次 HEAD | `fddaf638`（2026-08-27，clean） |
| pin / 部署基线 | `provenance-v2.8.json`、`compose.yaml`、`config.yaml` 无变更 |
| 依赖锁 | `go.mod`/`go.sum`/web lockfile 无变更（`package.json` 仅 +1 script 行） |
| 区间变更归属 | VP-018/019/020 交付（Admin 面，均双审计关门）+ VP-009 W13（用户书面批准 + grok 独立闭环）+ VP-010 W27（done）；Profile 默认集 / 模块矩阵 / Manifest 装配语义未变 |
| VP-009 / VP-010 | 无现行暂挂 |
| 本 VP 意图 | 不改 Profile 默认集 / 模块矩阵 / Manifest / 协议 pin |
| 复核结果 | **PASS（架构激活）**，不消费业务解锁 scope，不暂挂 `go` |

## 未选方案

- 不加 independent 层激活审视：VP-021 激活门闩明文「VRev intent-activation：self；可加 independent」——self 即满足（P-004 规则可唯一判定时不询问）。后续如需强审可在关门时走 `/vision-audit` / `/audit`。
- 不把 A3 余项（多实例 / 就绪探针扩依赖 / PG 锁 vs Redis vs 队列）拉进本波：仍 trigger-gated。

## 证据

- VRev-046（self · `pass` · 0 required）：`docs/vision/reviews/VRev-046-vp021-intent-activation.md`
- VP-021 v0.2.0 激活 + 修订短史；roadmap 行 21 / 架构分支 A7 / 组合焦点；revisions VR-048
- 本区 `workspace.md`、`goal-tree.md`、Root 五件套（本文件 + `00-meta.md` 纲领/信息台账）