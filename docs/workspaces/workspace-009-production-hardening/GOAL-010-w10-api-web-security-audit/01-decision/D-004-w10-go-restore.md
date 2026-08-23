---
id: D-004-w10-go-restore
goal: GOAL-010-w10-api-web-security-audit
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-010-w10-api-web-security-audit
version: 0.1.0
---

# D-004 · W10 关门与 VP-008 go 宣称恢复（2026-08-21）

## 决定（用户书面指令）

用户 `/govern` 指令原文："……同步滞后索引。修正审计报告提出的 recommended 意见，**然后关门并恢复go宣称**。"

1. **GOAL-010 关门**：`status: done`（S1–S4 4/4）。闭合条件 = D-003 调和后的 **3 条 fixed**（F-001/F-002/F-007，A-003 independent pass 确认 genuine）+ **4 条作废合法闭合**（user-overruled，本指令书面确认）+ A-003 recommended ×3 fixed（E-003）+ 开放 required = 0（A-004 闭合记录）。
2. **VP-008 go 消费有效性宣称：恢复**。暂挂区间 2026-08-21（D-002 落盘）至本决定；恢复依据与 W7/W8/W9 D-00x 先例一致：本波 required/recommended 全部合法闭合、self（A-002）与 independent（A-003 · grok-build grok-4.6）双复核 pass、回归全绿（go vet/test 全绿；web 1084/1084 + tsc 0）。
3. I-003 随本关门关闭（grok 腿已满足 + 用户书面关闭，见 A-004）。

## 未选方案

- **维持暂挂至密码轮换完成**：密码轮换是环境侧残余动作（A-002/A-003 一致定性，非本波 required）；以其为 go 前置会把代码闭合与运维动作不当耦合。轮换责任保留在残余移交清单。

## 残余与移交（不阻断关门）

| 项 | 责任 | 说明 |
|----|------|------|
| `192.168.31.213` 数据库密码轮换 | 用户/运维 | 凭据曾入版本控制历史，视为已泄露；git 历史与本地 gitignored `.env` 不改写 |
| A-001 informational F-008～F-012 | 后续波次评估 | 维持原判（已接受取舍/环境卫生） |