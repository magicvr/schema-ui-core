---
id: D-001
doc: decision-entry
goal: GOAL-001-outbound-mail
status: accepted
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# D-001 · 开区 scaffold 与 A6 纲领路线图

## 背景

用户确认：响应独立审计意见，然后激活 VP-017，然后交 `/govern` 开区。未另写工作区 slug。VRev-037 independent `pass`（0 required；V-F070/V-F071 recommended）。VRev-038 self `pass`（0 required）。A6 退出分母已由 VP-017 v0.1.0 / VR-040 冻结。

工作区 slug 按 VP-013 / VP-014 / VP-015 / VP-016 惯例取 `workspace-017-outbound-mail`（`workspace-NNN` + VP slug）。Root slug 取 `GOAL-001-outbound-mail`，与 VP / 工作区后缀对齐。用户确认的是「开区」动作；slug 为惯例推导并在此留痕，不是另一次口头点名。

## 决策

1. 激活 `VP-017-outbound-mail`（v0.2.0 `planned → active`）。
2. lead 工作区 slug = `workspace-017-outbound-mail`；Root = `GOAL-001-outbound-mail`。
3. 纲领路线图 R1～R4：端口与发送合同冻结 → SMTP 接入与配置面 → 默认 sink + 公共面去客户端类型 → 显式 SMTP 路径证据 **与** `readyz`。
4. 配置面：缺省无 SMTP（capture/log sink）；SMTP 为显式配置的生产/验收路径。不改 Compose 默认依赖。重启生效，热加载不进本波（I-006 / I-017-006）。
5. 开区审计模式 **none**（可逆文档 scaffold）。R1 合同冻结起按内核门禁走 **self**；SMTP 生产路径实施按 **independent**（项目默认 grok build · grok-4.6 · `/audit`）。
6. 本回合**不**创建 R1 子目标、**不**改 `apps/api` 代码。
7. 登记 I-001～I-004 required collecting（对应 I-017-003/004/001/002）；I-005 non-blocking（I-017-005）；I-006 non-blocking registered（V-F071）。

## 架构类 freshness（V-F070）

VP-008 强制 freshness 的对象是后续**业务** VP。本 VP 是架构 A6，按自身激活门闩做轻量复核：

| 项 | 值 |
|----|-----|
| 原 `go` 候选 | `ed99e88`（2026-08-10） |
| 现行 HEAD | `250cb9c` |
| VP-009 / VP-010 | 无开放中高危暂挂宣称 |
| Vision open required | 0 |
| F-007 | 上传授权深度仍 deferred；本 VP 不扩张授权 scope |
| 是否消费业务解锁 | **否**。不改 Profile / 模块矩阵 / Manifest 为意图 |
| 结果 | **PASS**。不暂挂 `go` |

`consumer_vp` = VP-017；`last_freshness_review_at` = 2026-08-22；`next_freshness_review_trigger` = 若实施改变共同门禁 / Profile / 模块矩阵则重做。

## 为什么

- 新纲领波次独立工作区，避免写入已 closed 的 workspace-016。
- slug 与 VP id 对齐，便于组合索引。
- I-001～I-006 满足 V-F071；freshness 表满足 V-F070。
- 开区时账号表仍无 email 列：本区不补该列（后续 Admin VP）。

## 未选方案

- 继续 `planned` 只写 VP：用户已要求响应独立意见后激活并开区。
- 重开 workspace-016：VP-016 已关门且默认不接新区。
- 一开区就接 SMTP 客户端：R1 合同未冻结。
- 默认改为必须有 SMTP：与 A6 内嵌默认冲突。
- 把账号 email / 自助恢复打进本 Root：混层，VP 非目标已排除。
- 把本 VP 当业务 VP 走 VP-011 式完整 freshness 矩阵：解锁 scope 不匹配。
- 等待用户另写 slug 再开区：用户本轮已指令「开区」；惯例与 VP-013/014/015/016 同构，记录推导以免静默。
