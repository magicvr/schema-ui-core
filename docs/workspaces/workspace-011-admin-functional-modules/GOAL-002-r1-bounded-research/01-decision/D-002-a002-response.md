---
id: D-002
doc: decision-entry
goal: GOAL-002-r1-bounded-research
status: accepted
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-002 · A-002（grok build independent）响应方案

## 背景

A-002（2026-08-14，grok build · grok-4.6 · xhigh，conditional）对 R1 调研结果提出 F-001～F-008。用户裁决 **A**：F-05/F-06 降档。

## 决策

1. **F-001（required high）→ fixed**：订单 / 钱包降为常用（S-13/S-14）；R2 第一波 = F-01 仪表盘、F-02 导入导出、F-03 个人中心与账户安全（含账号启停）、F-04 通知中心。
2. **F-002（required med）→ fixed**：S-12 回收站/软删除修正——需新持久化 + 管理 UI，不复用 0006 records_retire。
3. **F-003（required med）→ fixed**：C-01 修正为「不含产品态启停」；账号启停（启用/停用/手动解锁）并入 F-03。
4. **F-004（recommended）→ fixed（口径修正）**：`7 协议对照不再外推 S3 9/0；R2 方案做独立对照。
5. **F-005（recommended）**：登记 R2 方案必办（通知边界）。
6. **F-006（recommended）→ fixed（显式补入）**：组织/部门/岗位 → B-10；登录日志独立视图 → B-11。
7. **F-007（recommended）→ fixed**：Root I-001 verified（required 全部闭合后）；workspace.md 指针同步；附件路径用 Q2 全路径。
8. **F-008（recommended）**：登记 R2 方案必办（home 装配内容扩展声明、订单依赖/桩声明）。

## 门禁影响

- Root I-001（required）→ verified（证据 = I-011-001 v1.1.0）。
- R2 立项门禁解除；立项时逐项核对 R2 方案必办清单（I-011-001 `8）。

## 审计模式

响应为事实修正 + 用户裁决留痕：**self 记录**（A-002 响应节），不新开独立审计。
