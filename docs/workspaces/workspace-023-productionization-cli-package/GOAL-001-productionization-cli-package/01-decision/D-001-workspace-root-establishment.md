---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-001-productionization-cli-package
version: 0.1.0
---

# D-001 · 开区事务（R1 边界 + freshness 留痕）

## 决策

1. **绑定**：lead workspace = `workspace-023-productionization-cli-package`（delivery 单 lead）；`root_goal` = `GOAL-001-productionization-cli-package`（parent null）；`primary_plan` = VP-023（active v0.2.0）；不改 Charter primary（workspace-001）。
2. **审计模式**：阶段关门 default `self`；**R5 与 Root 关门 = `independent`**（grok build · grok-4.6 · high 先例；provider 失败不静默降级）。
3. **golden-field 初始化（R1 前置）**：最小骨架（Go 组合根 + Web 包消费骨架 + 探针 + README）；**`replace`/`file:` 为 R1 占位**——判据 #1 实证（真实 tag/registry）完成后移除并留痕；server 壳与 config 属 R4 运维交付。
4. **R1 信息门禁**：I-023-001（Go 通道：origin tag + 公共/私有 proxy 形态）、I-023-002（npm registry 目标与凭据）、I-023-003（PG 实例，最晚 R4）——均 open 登记，R1 生效前须闭合 #001/#002。

## freshness 三字段（VP-022 先例 · 2026-08-29）

| 字段 | 值 |
|------|-----|
| `consumer_vp` | VP-008 `go` 消费基线：候选链 `ed99e88` → … → `5c168070`（VP-022 激活）→ **`041744b3`（本次）**；不暂挂 |
| `last_freshness_review_at` | 2026-08-29（架构类轻量复核 · VRev-051：依赖锁/迁移/Profile/部署基线无外部变更；0063 与 pin 2.9 = VP-022 自产已验） |
| `next_freshness_review_trigger` | 下一激活/关门事务前复核；或 pin/依赖锁/迁移/Profile/部署基线/门禁语义任一变更时立即复核 |

## 未选方案

- 不改 Charter（fork 并存维持——用户指令）。
- 不做 G2 多模块细 tag、运行时镜像、业务域。
- 不在 R1 前置阶段推送真实发布（凭据/网络授权以用户为界）。