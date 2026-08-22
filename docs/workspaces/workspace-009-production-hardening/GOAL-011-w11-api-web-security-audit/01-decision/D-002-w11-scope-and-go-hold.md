---
id: D-002-w11-scope-and-go-hold
goal: GOAL-011-w11-api-web-security-audit
status: accepted
created: 2026-08-22
updated: 2026-08-22
parent: GOAL-011-w11-api-web-security-audit
version: 0.1.0
---

# D-002 · W11 修复范围裁决与 go 宣称暂挂（2026-08-22）

## 决定（用户书面指令）

用户目标轮次指令原文："推进工作区9目标11，直到顺利闭门"（本波按 W7–W10 先例整单执行闭门路径）。

1. **required 修复范围 = 6 条**：整单采纳 A-001 全部 required——F-001（HIGH · Postgres 创建/导入用户 500）、F-002（HIGH · 删除与回收站快照非原子）、F-003（HIGH · MFA 第二因子可在线穷举）、F-004（MEDIUM · JWT 轮换使 MFA 密文不可解）、F-005（MEDIUM · 验证码并发一次性问题）、F-006（MEDIUM · 钱包对账坏 JSON 静默全库 + 写操作权限过松）。
2. **VP-008 go 消费有效性宣称：暂挂**——自本决定落盘起至 6 条 required 全部合法闭合、self（A-002）与 independent（A-003 · grok）双复核通过、回归全绿为止，不宣称 VP-008 go 消费有效性；复核通过后按 W9/W10 D-004 先例在关门记录中恢复（对齐 W7/W8/W9/W10 D-00x 先例）。S1 落盘时"不自动暂挂"仅指不因审计 `fail` 单一事实擅自悬挂；用户本轮指令已授权整波闭门路径，按先例进入暂挂-恢复周期。
3. **实施顺序**：S3 按 F-001 → F-002 → F-003 → F-004 → F-005 → F-006 实施 + API/Web 回归全绿 → A-002 self 审计 → S4 independent 复核（工作区惯例 grok build · grok-4.6 · high）→ 闭合记录 + go 恢复 + 关门。
4. **recommended（13 条）**：A-003 独立复核后按 W9 E-005 / W10 E-003 先例处置——真实缺口就地修正，纯设计取舍/误报逐条留痕（fixed/overruled 有据），不以推荐项阻塞关门。
5. I-002 状态：open → **verified**（证据 = 本文件）。

## 未选方案

- **仅修 3 条 HIGH 最小范围**：与"顺利闭门"指令及 W7–W10 整单采纳先例不符，不采用。
- **MFA 密钥改用独立专用 secret**（F-004 建议分支一）：需新增部署配置面与密钥生命周期管理；轮换窗口内尝试 previous 解密（分支二）与既有 VP-016 JWT 双密钥机制同构、零配置面变更，选分支二。
- **wallet 全库对账改显式哨兵字段**（F-006 建议第三点）：空字符串哨兵已被 jobs 测试与协议行为锁定，在线文档与 I-PROTO 测试套件 pin 该语义；解码失败 400 已消除"垃圾 JSON 静默全库"的利用面，选择保留 `accountId: ""` 语义并在代码注释与闭合记录中留痕（对齐 W9"F-003 作废需有据"的逐条核实精神——保留有测试锁定的行为，不因 audit 措辞改动协议）。
- **本波同时实施 13 条 recommended**：S3 先收敛 6 条 required；recommended 在独立复核后逐条处置，避免在关门路径内扩散回归面。

## 对 go 宣称的影响

- 暂挂起点：2026-08-22（本决定落盘时）。
- 恢复条件：6 条 required 全部 fixed，self（A-002）与 independent（A-003 · grok build）确认无新引入缺陷，API/Web 回归全绿后，在关门记录中书面恢复（对齐 W9 D-004 / W10 D-004 模式）。