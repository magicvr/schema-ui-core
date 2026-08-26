---
id: E-001
title: 开区与激活记录（VRev-046 · 架构类 freshness · 五件套建立）
date: 2026-08-27
status: done
---

# E-001 · 开区与激活（2026-08-27）

## 事实

1. **决策层（/vision）**：VRev-046（self）`pass` 落盘（`docs/vision/reviews/VRev-046-vp021-intent-activation.md`）；reviews.md 索引行 + open required 投影（0）同步；**VP-021 v0.2.0 `planned → active`**（lead 字段、状态表、修订短史更新）；roadmap 行 21 / A7 / 架构分支当前拍 / 组合焦点同步；revisions 追加 VR-048（editorial）。
2. **架构类 freshness**：VP-008 `go` 消费有效性复核 **PASS**——`ed99e88` →（VP-017 `250cb9c`）→ 现行 HEAD `fddaf638`（2026-08-27，clean）；`provenance-v2.8.json` / `compose.yaml` / `config.yaml` / `go.mod` / `go.sum` / web lockfile 无变更（`apps/web/package.json` 仅 +1 script 行，无依赖变化）；区间变更全可追溯至已审节目（VP-018/019/020、VP-009 W13、VP-010 W27）；Profile 默认集 / 模块矩阵 / Manifest 装配语义未变；VP-009/VP-010 无现行暂挂；不暂挂 `go`、不消费业务解锁 scope。
3. **实现层（/govern）**：scaffold `docs/workspaces/workspace-021-graceful-shutdown-and-connection-drain/`（`workspace.md` + `goal-tree.md` + `attachments/`）；Root `GOAL-001-graceful-shutdown-and-connection-drain` 五件套 + 三个 ledger 目录一次建齐（00-meta 含 P-001 纲领 R1～R3 + 信息台账 I-001～004；D-001 决策含 freshness 留痕；E-001 本记录；A 台账空索引）。
4. **决策层绑定投影**：`docs/vision/workspaces.md` 追加 workspace-021 行；VP-021 `lead_workspace` 表 0 区 → 1 区绑定。

## 验证 / 后续

- V-F081（Root 纲领 + 信息台账承接）→ fixed（本事务完成）。
- V-F082（架构类 freshness 激活留痕）→ fixed（D-001 + 本记录 + VRev-046）。
- 下一步建议：R1 立项 GOAL-002 前，先按 P-005 关闭 I-001（Job 停机语义）与 I-002（grace/超时默认与配置键）两个 required 信息项（用户裁决）。