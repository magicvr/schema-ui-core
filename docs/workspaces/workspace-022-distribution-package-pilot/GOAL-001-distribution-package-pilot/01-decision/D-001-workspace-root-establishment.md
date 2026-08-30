---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-001-distribution-package-pilot
version: 0.1.0
---

# D-001 · 工作区/Root 成立（开区事务）

## 决策

1. **绑定**：lead workspace = `workspace-022-distribution-package-pilot`（`vision_role: delivery`，单一 lead）；`root_goal` = `GOAL-001-distribution-package-pilot`（`parent: null`）；`primary_plan` = `VP-022-distribution-package-pilot`（active v0.3.0）；不改变 Charter `primary_workspace`。
2. **纲领路线图冻结**：R1 契约冻结面 → R2 Go 库包闭环 → R3 Web 包闭环 → R4 零冲突升级演练 → R5 证据与 go/no-go（串行；R2/R3 同依赖 R1 可评估并行；R5 可拆证据与报告子目标）。
3. **审计模式**（P-002 §3.1 按风险确定）：阶段关门 default `self`；**R5 发布/兼容门禁（release/compatibility 高影响）与 Root 关门 = `independent`**（项目级默认执行路径：自审后本地 grok build 独立审计；provider 失败不静默降级）。R1 契约冻结面属 migration/兼容前瞻，R1 关门采用 `self` + 标记 R5 independent 预留。
4. **信息项**：I-001（kernel 冻结面清单，required · R1）、I-002（npm 包拆分/peer 策略，required · R3）、I-003（发布通道，required · R5）、I-004（演练样本设计，non-blocking · R4）——登记于 00-meta，按 P-005 门禁推进。

## freshness 三字段（VRev-048 V-F084 开区义务执行 · 2026-08-29）

| 字段 | 值 |
|------|-----|
| `consumer_vp` | VP-008 `go` 消费基线：候选 `ed99e88`（VP-008 关门）→ 历次复核链 `fddaf638`（VP-021 激活）→ **`5c168070`（本次，2026-08-29）**；`go` 未暂挂 |
| `last_freshness_review_at` | 2026-08-29（平台/架构类轻量复核 · VRev-049：协议 pin / 依赖锁 / 迁移 / Profile 默认集无变更；compose.yaml +5 = VP-021 交付内 `stop_grace_period`，已随其关门核销） |
| `next_freshness_review_trigger` | 下一激活/关门事务前复核；或协议 pin、依赖锁、迁移台账、Profile 默认集、部署基线、认证授权门禁语义任一变更时立即复核并留痕 |

## 未选方案

- 不采用多模块细版本（G2）起步：单主线统一节奏下粗粒度（G1）成本最低，G2 留待 go 后评估。
- 不做 CLI 交付（仅评估 CLI 形态）；不做运行时基线镜像/模块下载（与 Charter 非目标边界保持距离）。
- 不把「包消费」写入 Charter 成功边界 #1（试点不改 Charter；结论由 R5 go/no-go 报告再议）。