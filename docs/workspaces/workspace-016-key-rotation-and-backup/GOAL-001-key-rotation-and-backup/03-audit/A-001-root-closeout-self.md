---
id: A-001
doc: audit-entry
goal: GOAL-001-key-rotation-and-backup
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# A-001 · Root 关门自审（2026-08-22）

- **source**：self
- **auditor**：编排器（/govern）
- **类型 / scope**：close-out · Root GOAL-001 整体（R1～R5 五阶段；VP-016 方向级判据 1–6）
- **verdict**：pass（待 independent 复核）

## 范围与区间

基线 = 开区提交 `5195104`；区间内四个子目标（GOAL-002～GOAL-005）全部 done + GOAL-006 证据登记完成。审计材料：各目标五件套、四份 git checkpoint（`c96e963`/`8116565`/`1b8e9b0`/`1dc6975`）、GOAL-006 E-001 新鲜实跑记录。

## 成功标准逐条对照（VP-016 方向级判据）

| # | 判据 | 状态 | 证据链 |
|---|------|------|--------|
| 1 | JWT 轮换合同落地：可配置 current+previous；签发只用 current；重叠窗 previous 可验 | 达成 | R1 配置面（D-002 + config 测试）；R2 双密钥验签（verifyAccess + TestDualKeyRotationOverlapWindow 4/4） |
| 2 | 未配置 previous 时本地/Compose 默认仍能开发与快测；轮换非启动硬依赖 | 达成 | R4 六层证据实跑（GOAL-005 E-001，含 compose 解析实证） |
| 3 | 轮换后恢复：两方言既有备份路径上轮换后启动且鉴权可核对 | 达成 | R3 双循环（SQLite VACUUM INTO / PG pg_dump→restore；GOAL-004 E-001 + A-001） |
| 4 | 显式双密钥下轮换路径与恢复路径都有可核对证据 | 达成 | GOAL-006 E-001 四项新鲜实跑全 PASS |
| 5 | 未进入 A3/KMS/PITR/Admin 功能/业务域；未改 Charter；无热加载/第二套 dump 宣称 | 达成 | GOAL-006 E-001 越界核对单（diff 范围 10 文件；charter 零改动） |
| 6 | 开放 required finding = 0（或已合法闭合） | 达成 | 各目标意见台账：GOAL-002 A-001、GOAL-003 A-001+A-002（F-001/2/3 recommended 已 fixed）、GOAL-004 A-001、GOAL-005 A-001 —— 开放 required 均为 0 |

## Findings

无 required。备注（非 finding）：I-005 non-blocking 保持 collecting——其默认措辞（previous 可验）已在 VRev-035 冻结进退出判据 1，本波按默认交付；"旧 access 立即失效残余"若未来被用户选择，属新决策而非本 Root 缺口。

## 必改项汇总（required 列表）

空。

## 结论 + 建议下一步

五阶段全部关门、六条判据证据齐备、0 开放 required finding、信息门禁 I-001～I-004 全部 verified。建议下一步：independent 关门审计（grok build `/audit` 流程，落盘 Root A-002）→ 编排器合并响应 → 通过后 Root `done` 5/5。
