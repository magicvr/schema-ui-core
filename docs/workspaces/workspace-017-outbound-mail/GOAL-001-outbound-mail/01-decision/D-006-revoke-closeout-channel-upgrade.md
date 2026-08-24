---
id: D-006
doc: decision-entry
goal: GOAL-001-outbound-mail
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-006 · 否决 Root/VP 关门并升级渠道分母

## 背景

用户书面（2026-08-24）：

- 认为 017 交付偏移预期；**不**作废目标。
- **显式否决**工作区 017 与 VP-017 的关门：只回退关门状态，**不**回退、**不**改变实施历史。
- 用本工作区对应 VP 与 Root 承接讨论方案（可切换渠道；默认 mock 站内虚拟渠道；生产 Resend；设置热切换与配置；试发控制台）。
- 后继开设新子目标实现真实预期。
- VP-018 及对应工作区暂时冻结，在 017 重新关门之前不允许推进。

历史：D-005 / GOAL-005 / Root A-001+A-002 / VP-017 v0.3.0 / VRev-039 曾把 SMTP 专用 A6 收成组合层 `closed` 与 Root `done 4/4`。该关门相对**当时分母**的实施核验不在此否认；用户否决的是把它当作出站邮件的成功交付。

## 决定

1. **否决关门，不作废。** Root `status: done → active`。GOAL-002～005 保持 `done`。不改写 D-002～D-005、E-002～E-005、A-001/A-002 原文与 verdict。不删除 `apps/api` 已落地的 `MailSender` / SMTP / `CaptureSink`。
2. **升级现行分母**（投影 VP-017 v0.4.0）：
   - `MailSender` 仍是唯一发送合同。
   - 第一期渠道 = **mock（默认）** + **Resend（生产）**；SMTP 适配器**保留不删**，不再是唯一生产权威。
   - mock = 管理员可检视的站内出站记录，**不是**用户站内通知。
   - 设置「邮件」tab：选渠道、填配置、热切换、试发走同一端口；mock 记录用表，不整包塞进一张卡片，不 PATCH `/api/settings/default`。
3. **纲领扩展**：R1～R4 保持已完成；新增 R5（渠道合同）→ R6（mock+Resend）→ R7（设置/热切换/试发）→ R8（证据/`readyz`）。`progress` = 4/8。
4. **信息项**：I-001～I-006 历史 verified 不动。I-007 / I-008 / I-012 **verified**（本决策）。I-009 / I-010 / I-011 **collecting**，分别最晚 R7 / R6 / R5。
5. **本回合创建 R5 子目标** `GOAL-006-channel-provider-contract`；**不**创建 R6～R8 子目标（P-001 按阶段）；**不**改应用代码。
6. **018 冻结**不在本区五件套执行，由 VP-018 / workspace-018 Root `blocked` 落盘。本区再关门前不得声称 018 可推进。
7. D-005 中「Root/VP 可按 SMTP 专用分母关门」的组合结论由本条 **supersede**。D-005 对 R4 技术事实（HTML 不进分母、R4 启动单例、`readyz` 仅显式 SMTP 扩依赖）**仍成立**，不 superseded。

## 为什么

- 用户要保留 017 作为出站邮件主线，同时拒绝「SMTP 专用波次 = 出站邮件已成功」。
- 回退关门而不回退实施，避免洗白审计、也避免丢掉已付过钱的端口与 SMTP 适配器。
- 渠道模型与已有 `MailSender` 端口同构；Resend 是用户点名的生产渠道；mock 站内记录解决 capture 容量 1、无法给人测的问题。
- 冻结 018 避免身份波次把历史 capture/SMTP 锁成验收权威。

## 未选方案

- **作废 VP-017 / 取消 Root / 新开 workspace-019**：用户明确否决。
- **原地把 GOAL-001 当新意图却继续 `done`**：关门效力仍在，脏成功叙事还在。
- **回退 SMTP/CaptureSink 代码**：违反「不回退实施历史」。
- **第一期删除 SMTP、只留 mock+Resend**：会改实施史；SMTP 降为非唯一生产权威即可。
- **mock 发布到用户通知**：与 Notification Transport 非目标冲突。
- **本回合创建 R6～R8 全部子目标**：违反 P-001。
- **本回合实施 Resend/设置页**：R5 合同（尤其 I-011）未冻结。

## 影响

- VP-017 `closed → active`；RT-M01 `delivered → in-progress`。
- 再次关门必须走现行分母 + VRev-041 V-F075（recommended：independent）。
- 历史 Root 关门审计保留为当时分母下的意见，不构成现行 `done`。

## 后续

- GOAL-006 冻结 I-011（mock 持久化）及渠道注册形状。
- R6 前关闭 I-010；R7 前关闭 I-009。
- 不在本回合改 `apps/api` / `apps/web`。
