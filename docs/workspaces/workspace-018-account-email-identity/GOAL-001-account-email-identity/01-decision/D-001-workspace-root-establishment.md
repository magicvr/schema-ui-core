---
id: D-001
doc: decision-entry
goal: GOAL-001-account-email-identity
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-001 · 开区 scaffold 与 Admin 邮箱身份纲领路线图

## 背景

用户确认：按 `/vision` 建议的 1234 顺序全部执行——有界关闭 VP-017 → 落盘 VP-018 → Admin 类 freshness → 激活并交 `/govern` 开区。未另写工作区 slug。VRev-040 self `pass`（0 required；V-F073/V-F074 recommended）。VP-018 退出分母已冻结：账号邮箱绑定+校验；IAM 恢复 / 邀请 / 密码策略 / SMS 不进。

工作区 slug 按 VP-013～017 惯例取 `workspace-018-account-email-identity`（`workspace-NNN` + VP slug）。Root slug 取 `GOAL-001-account-email-identity`，与 VP / 工作区后缀对齐。用户确认的是「开区」动作；slug 为惯例推导并在此留痕，不是另一次口头点名。

## 决策

1. 激活 `VP-018-account-email-identity`（v0.2.0 `planned → active`）。
2. lead 工作区 slug = `workspace-018-account-email-identity`；Root = `GOAL-001-account-email-identity`。
3. 纲领路线图 R1～R4：身份合同冻结 → 双方言 schema + 唯一性 → 绑定/校验消费 `MailSender` → 证据（capture + 唯一性 + 无 IAM 恢复）。
4. 配置/模块面：不改 Profile 默认集；归属既有 `core.auth-session`。缺省无 SMTP（capture sink）；校验信走 VP-017 端口。
5. VP 已冻结：账号可空邮箱（I-003）；换绑进本波（I-004）。R1 仍须冻结唯一性细则（I-001）与校验形态（I-002）。
6. 开区审计模式 **none**（可逆文档 scaffold）。R1 合同冻结起按身份/数据门禁走 **self**；schema 迁移与绑定实施按 **independent**（项目默认 grok build · grok-4.6 · `/audit`）。
7. 本回合**不**创建 R1 子目标、**不**改 `apps/api` / `apps/web` 代码。
8. 登记 I-001/I-002 required collecting；I-003/I-004 registered（VP 冻结投影）；I-005 required collecting（最晚 R3）；I-006 non-blocking collecting。

## Admin 类 freshness（V-F074）

VP-008 强制 freshness 的对象是后续**业务** VP。本 VP 是 Admin 功能身份面，不是 Tier D：

| 项 | 值 |
|----|-----|
| 原 `go` 候选 | `ed99e88`（2026-08-10） |
| 现行 HEAD | `092bf37` |
| 工作树 | clean |
| VP-009 | W11 D-004 已恢复 `go`；无现行暂挂 |
| VP-010 | W1–W25 done；无现行暂挂 |
| Vision open required | 0 |
| F-007 | 上传授权深度仍 deferred；本 VP 不扩张授权 scope |
| 是否消费 Tier D 解锁 | **否**。不改 Profile / 模块矩阵 / Manifest 为意图。R2 将改 `users` 迁移台账（分母内） |
| 结果 | **PASS**。不暂挂 `go` |

`consumer_vp` = VP-018；`last_freshness_review_at` = 2026-08-24；`next_freshness_review_trigger` = 若实施改变共同门禁 / Profile / 模块矩阵 / Manifest / 协议 pin 则重做。

## 为什么

- 新纲领波次独立工作区，避免写入已 closed 的 workspace-017。
- slug 与 VP id 对齐，便于组合索引。
- I-001～I-006 满足 V-F073；freshness 表满足 V-F074。
- 运输面已由 VP-017 有界交付；本区只消费端口，不重做 SMTP。

## 未选方案

- 继续 `planned` 只写 VP：用户已要求 1234 全部执行含开区。
- 重开 workspace-017：VP-017 已关门且默认不接新区；该区明确不改 `users`。
- 一开区就改 `users` DDL：R1 合同未冻结（I-001/I-002 仍 collecting）。
- 把忘记密码 / 邀请 / 密码策略打进本 Root：混层，VP 非目标已排除。
- 把本 VP 当 Tier D 业务 VP 走订单/CMS 解锁：解锁 scope 不匹配。
- 等待用户另写 slug 再开区：用户本轮已指令「开区」；惯例与 VP-013～017 同构，记录推导以免静默。
