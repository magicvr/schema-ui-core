---
id: GOAL-001-design-system-and-ui-experience
doc: decision-entry
record_id: D-007
status: superseded
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
---

## D-007 · Root 关门用户书面确认（**superseded**）

### 原主张（已作废）

曾将多轮目标指令「继续推进工作区6，直至根目标顺利关门」直接视为 D-005 类书面确认并写 `status: done`。

### 为何 superseded

- 与 A-006 / D-006 教训及计划门禁冲突：**Root/`workspace` `done` 须用户对「可否关门」的显式书面确认**（或书面驳回并列出缺口）。
- 目标指令描述的是任务意图，**不等于**对当前 closeout 证据包的验收签字。
- 编排器/验证器要求：无显式确认时诚实停在 **closeout-ready**（Root `active`，开放 required = 0）。

### 现行状态（本文件 supersede 后）

1. Root `status: active`；S2/S3 可保持勾选（E-002 + A-008）；**S5 过程关门检查点取消勾选**直至新确认。
2. `workspace.md` `status: active`。
3. D-007 **superseded**；待用户显式确认后另立 **D-008**（或用户书面重申接受后恢复/修订本决策）。
4. 实施与 finding 闭合（F-VUI-001/002 fixed、GOAL-003 done）**不**回滚。
