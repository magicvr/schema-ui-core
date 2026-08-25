---
id: D-001
doc: decision-entry
goal: GOAL-001-iam-recovery
status: accepted
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# D-001 · 开区 scaffold 与 IAM 纲领路线图

## 背景

用户 2026-08-25 书面指令：「立项新vp（调用grok build做独立审计），并开设对应的工作区，同时顺手修正对齐债务」。按 `/vision` 建议立项 **VP-019-iam-recovery**（Admin 功能 · IAM：密码策略 / 邀请入职 / 自助恢复状态机）。本地 **grok build**（grok-4.6 · reasoning high）独立 Vision Review **VRev-043** verdict `pass`（0 required；V-F076/077/078 recommended → 本事务内 fixed）。未另写工作区 slug。

工作区 slug 按 VP-013～018 惯例取 `workspace-019-iam-recovery`（`workspace-NNN` + VP slug）；Root slug 取 `GOAL-001-iam-recovery`。用户确认的是「开区」动作；slug 为惯例推导并在此留痕。

## 决策

1. 激活 `VP-019-iam-recovery`（v0.2.0 `planned → active`；VR-047）。
2. lead 工作区 slug = `workspace-019-iam-recovery`；Root = `GOAL-001-iam-recovery`。
3. 纲领路线图 R1～R4：合同冻结 → 自助恢复全链 → 密码策略 + 邀请入职 → 证据/关门（V-F076 落地）。
4. 配置/模块面：归属既有 `core.auth-session`；不改 Profile 默认集 / 模块矩阵 / Manifest 装配；策略配置 UI 仅 `admin.settings`（与邮件 tab 同形），强制面在 `core.auth-session` 对全 Profile 生效，禁止为策略把 `admin.settings` 加入 mvp 默认集。
5. 运输面：只消费 VP-017 `kernel.MailSender` 现行渠道（mock 默认 / Resend 生产）；无生产渠道时 mock 出站记录可端到端取信。
6. 管理员重置：保持既有 `must_change_password` 特权路径，不冒充自助恢复；无邮箱账号不自助（2026-08-22 产品事实，I-006 registered）。
7. 开区审计模式 **none**（可逆文档 scaffold）。R1 合同冻结起按身份/数据/security 门禁走 **self**；恢复全链实施按 **independent**（项目默认 grok build · grok-4.6 · `/audit`）。
8. 本回合**不**创建 R1 子目标、**不**改 `apps/api` / `apps/web` 代码、**不**直接改恢复/邀请/策略 DDL。
9. 登记 I-001/I-002/I-009 required collecting（最晚 R1 / R2 方案冻结前）；I-003/004/005/007 collecting（R2/R3）；I-006 registered（2026-08-22 产品事实）；I-008 non-blocking。

## Admin 类 freshness（V-F077）

VP-008 强制 freshness 的对象是后续**业务** VP。本 VP 是 Admin 功能 IAM 面，不是 Tier D（VRev-043 独立复核）：

| 项 | 值 |
|----|-----|
| 原 `go` 候选 | `ed99e88`（2026-08-10，clean）；解锁 scope = 标准业务模块框架能力，不是 IAM 三件本身 |
| 本轮基线 | `092bf37`（VP-018 激活时 HEAD；VRev-040 Admin 类 PASS） |
| 现行 HEAD | `66f5fd1f`（`chore(web): 再生成协议符合性声明（buildId→76bd21aa）`） |
| 工作树代码 | HEAD 相对干净；脏路径仅为愿景索引 + VP-019 草案，不含 Profile/模块矩阵/协议 pin |
| 比对区间 | VP-017 渠道模型 + VP-018 邮箱身份（迁移 0054/0055）；`kernel/profile.go` 默认集未改、`admin.account` 仅增 mail 路由 |
| 共同门禁 | 认证/授权/fail-closed 语义未改成新的暂挂条件；W16 改密路径仍 bump `token_version` |
| VP-009 | W1–W4、W6–W11 done；W5 扫描 0 中高危；W11 D-004 恢复 `go`；无现行暂挂 |
| VP-010 | W1–W25（含 GOAL-037）done；无现行 `go` 暂挂 |
| Vision open required | 0 |
| F-007 residual | 上传授权深度仍 deferred（owner = VP-008 lead）；本 VP 不借 IAM 面扩张上传授权 scope |
| 数据库意图 | 会变（恢复令牌表 / 邀请表 / 策略存储）——本波分母内工作，属实现期 freshness 复核点 |
| **结果** | **PASS（Admin 激活）**。不消费 Tier D 解锁；不暂挂 `go` |

`consumer_vp` = VP-019；`last_freshness_review_at` = 2026-08-25；`next_freshness_review_trigger` = 若实施改变共同门禁 / Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 协议 pin 则重做。

## 为什么

- 新纲领波次独立工作区，避免写入已 closed 的 workspace-017/018。
- slug 与 VP id 对齐，便于组合索引。
- I-001～I-009 满足 V-F076；freshness 表满足 V-F077；roadmap 第 19 行 + workspaces.md 登记满足 V-F078。
- 运输与身份面已由 VP-017/018 交付；本区只消费端口，不重做渠道/SMTP/邮箱身份。

## 未选方案

- 继续 `planned` 只写 VP：用户已指令开区。
- 重开 workspace-017/018：均已 closed 且默认不接新区。
- 一开区就改恢复/邀请/策略 DDL：R1 合同未冻结（I-001/I-002/I-009 仍 collecting）。
- 把 SMS/模板/组织权限/OIDC 打进本 Root：混层，VP 非目标已排除。
- 把本 VP 当 Tier D 业务 VP 走订单/CMS 解锁：解锁 scope 不匹配（VRev-043 门禁含义）。
- 为密码策略改 mvp Profile 默认集：VRev-043 明确禁止。
- 等待用户另写 slug 再开区：用户本轮已指令「开区」；惯例与 VP-013～018 同构，记录推导以免静默。