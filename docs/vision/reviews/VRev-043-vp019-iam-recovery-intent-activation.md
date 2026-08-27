---
doc_type: vision-review
id: VRev-043
status: active
source: independent
created: 2026-08-25
updated: 2026-08-25
version: 0.1.0
parent: null
---

# VRev-043 · VP-019 IAM 意图/对齐链/可行性/激活就绪（2026-08-25）

> 代贴说明：本条由编排器自本地 grok build（grok-4.6 · reasoning high）headless 会话原样誊入（2026-08-25 08:43–08:50 运行约 7m19s；会话提示词见本会话记录，.grok sessions 本地存档）。grok 按指令**只出报告文本、未写入** `reviews/`——落盘与索引由 `/vision`（编排器）完成，`source: independent` 保持不变。

| 字段 | 值 |
|------|-----|
| source | independent |
| auditor | grok-build (grok-4.6 · reasoning high) |
| scope | `VP-019-iam-recovery` 意图完备、Charter `@0.2.0` 对齐、硬前置 VP-017/VP-018 成立性、退出分母与三件同波、P-005 信息项、VP-008 `go` 消费有效性（Admin 类 freshness） |
| audit_type | vision-plan（意图 / 激活就绪） |
| verdict | pass |
| 建议 class | editorial |
| open required | 0 |

## 范围与结论

只读核对（未改任何文件）：P-006 / `alignment.md`、Charter `@0.2.0`、[VP-019-iam-recovery](../plans/VP-019-iam-recovery.md)（`planned` v0.1.0）、[VP-017](../plans/VP-017-outbound-mail.md) v0.5.0 `closed`、[VP-018](../plans/VP-018-account-email-identity.md) v1.0.0 `closed`、roadmap Admin 功能分支、`reviews.md`（open required = 0）、workspaces.md，以及代码：`kernel.MailSender`、`core.auth-session` 账号模型、迁移 0054/0055、`must_change_password`、设置面 tab 形状、Profile 默认集。Freshness 比对 `092bf37` → HEAD `66f5fd1f`。

**总判：pass（0 open required）。** 单愿景与 `vision_ref` 精确匹配；硬前置两条均已按现行分母关门；新 VP + 单 lead delivery 区的结构选型成立；退出 1–6 方向可判定；Admin 类 freshness **PASS**，不暂挂 `go`。本 `pass` 允许**用户确认后激活 VP-019 并开 `workspace-019-iam-recovery`**。它不构成「忘记密码已交付 / 密码策略已可配置 / 邀请已上线」的任何宣称。

### 退出分母判定（判据 1–6）

| 判据 | 判定 |
|------|------|
| 1 自助恢复全链 | **可判定**：登录页发起 → verified 邮箱投递 → 冷却/过期 → 设新密 → 登录。今日登录页无入口（属交付面，非意图缺陷） |
| 2 密码策略强制 + 渐进生效 | **可判定**，依赖 I-019-003/007 在 R2 冻结默认参数与既有账号边界 |
| 3 邀请生成/投递/激活/撤销 | **可判定**，依赖 I-019-004/005 在 R3 冻结形态 |
| 4 只经现行 `MailSender` | **可判定**（端口唯一性已在代码） |
| 5 保留 `must_change_password`；非目标未混入成功条件 | **可判定** |
| 6 开放 required = 0 | **过程可判定** |

非目标充分：SMS/推送、模板中心、多邮箱、组织权限、OIDC/SSO/SCIM、业务域、outbox、重开 012–018、改 Charter，均已点名。枚举/重放/邀请滥用显式留给 VP-009/010，不把持续程序塞进本波成功条件。

**三件同波：不越界。** roadmap 已把密码策略/邀请/自助恢复写成同一 IAM 截；三者共享账号模型、`MailSender` 与新密码校验。内部串行（R2 恢复 → R3 策略+邀请）避免单 Goal 一把做完。未遗漏第四件：管理员重置已是既有特权路径。邀请「管理员出示链接」仍可用，故邀请不把生产邮件渠道做成硬依赖。

### P-005 允许带未知立项

I-019-001～008 字段完整（问题/级别/门禁/最晚阶段/collecting）。级别总体合理：001/003/004/006/007 绑方案冻结；005 绑 R3；008 non-blocking 合理。缺陷（不升 required、须在 Root 纠正）：

1. **I-019-002** 影响「方案/实施」但最晚阶段写成 **R4 接入前**（VP-019:93）。恢复全链是 R2，TTL/冷却必须在 R2 方案冻结前关闭，不能拖到关门证据阶段。
2. **缺一条 required 未知：已登记 MFA 的账号如何走自助恢复**（绕过第二因子 / 要求 TOTP / 拒绝自助只走管理员重置）。`admin.mfa` 已在 admin 默认集；登录页已有 MFA 错误码。不登记则 R2 可能做成 MFA 旁路。最晚 = R1 合同冻结。
3. **I-019-003** 应顺带冻结：策略配置 UI 仅 `admin.settings`（与邮件 tab 同形），强制面在 `core.auth-session` 对 mvp/admin/demo 生效；禁止为策略把 `admin.settings` 加入 mvp 默认集。
4. **I-019-006** 产品事实已于 2026-08-22 确认（无邮箱不自助、走管理员重置）；R1 应落成 verified，不必重裁。
5. **I-019-001** 应把 VP-018 已冻结的 6 位码（`email_identity.go:31-36`）当作一致性默认候选，而不是空白选择题。

不得把验证码 vs 链接、邀请即建号、策略默认参数在 VP 层假装已选。

### VP-008 `go` 消费有效性 · Admin 类 freshness

VP-008 强制 freshness 的对象是**后续业务 VP**。VP-019 是 Admin 功能 IAM 面，不是 Tier D。按激活门闩做 Admin 类复核：

| 项 | 结论 |
|----|------|
| 原 `go` 候选 | `ed99e88`（2026-08-10，clean）；解锁 scope = 标准业务模块框架能力，不是 IAM 三件本身 |
| 本轮基线 | `092bf37`（VP-018 激活时 HEAD；VRev-040 Admin 类 PASS） |
| 现行 HEAD | `66f5fd1f`（`chore(web): 再生成协议符合性声明（buildId→76bd21aa）`） |
| 工作树代码 | HEAD 相对干净；脏路径仅为愿景索引 + 本 VP 草案，不含 Profile/模块矩阵/协议 pin |
| 比对区间 `092bf37`→`66f5fd1f` | VP-017 渠道模型（邮件 tab / Switcher / Resend / mock outbox）+ VP-018 邮箱身份（迁移 0054/0055 + 绑定/校验路由）；`kernel/profile.go` 默认集未改、`admin.account` 仅增 mail 路由（`profile.go:175-177`） |
| 共同门禁 | 认证/授权/fail-closed 语义未改成新的暂挂条件；W16 改密路径仍 bump `token_version` |
| VP-009 | W1–W4、W6–W11 **done**；W5 扫描 0 中高危未开子目标；W11 D-004 **恢复** `go`；无现行暂挂 |
| VP-010 | W1–W25（含 GOAL-037）**done**；无现行 `go` 暂挂 |
| Vision open required | 0 |
| F-007 residual | 上传授权深度仍 **deferred**（owner = VP-008 lead）。本 VP 不得借 IAM 面扩张上传授权 scope |
| 本 VP 意图是否改 Profile / 模块矩阵 / Manifest 装配 / 协议 pin | **意图否**。实施若新增模块或改默认集 → 迁移台账 |
| 数据库意图 | **意图会变**（恢复令牌表 / 邀请表 / 策略存储）。这是本波分母内工作，不是激活前已存在的失效；须复核 |
| **结果** | **PASS（Admin 激活）**。不消费 Tier D 解锁；不暂挂 `go` |

`consumer_vp` = VP-019；`last_freshness_review_at` = 2026-08-25；`next_freshness_review_trigger` = 若实施改变共同门禁 / Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 协议 pin 则重做。

### 不构成 fail / 不新开 required 的诚实边界

- 恢复/邀请/策略的**迁移台账会变**——属本波分母内工作，属实现期 freshness 复核点，不构成激活前失效。
- README 前向链接（workspace-019 指针）不得读成已开区；`workspaces.md` 仅在 `/govern` 开区后登记。

## Findings 清单

| id | level | 状态 | 一句话 |
|----|-------|------|--------|
| V-F076 | recommended | open | Root P-001 + 纠正 I-019-002 最晚阶段 + 补 I-019-009（MFA）+ 策略面/Profile 边界 |
| V-F077 | recommended | open | 激活记录留下 Admin 类 freshness（`092bf37` → `66f5fd1f`，不暂挂 `go`） |
| V-F078 | recommended | open | roadmap 已落盘表补第 19 行 VP-019；README 前向链接 ≠ 已开区 |

**required findings：无。**

### V-F076 · recommended · Root P-001 与 I-00N 纠正

- level: `recommended`；status: open；severity: low
- impact: 若只把 VP 落盘而不给 Root 阶段目标，R1 合同冻结没有承接地；I-019-002 的「R4 接入前」会让 TTL/冷却拖延至关门证据阶段；缺 MFA 未知项则 R2 可能做成 MFA 旁路。
- finding: lead Root `00-meta` 写 P-001 纲领阶段表（R1 合同冻结 → R2 自助恢复全链 → R3 密码策略 + 邀请入职 → R4 证据/关门），登记 I-019-001～008；**把 I-019-002 最晚阶段改为「R2 方案冻结前」**（不得复制 VP 表的「R4 接入前」）；**增补 I-019-009 required（最晚 R1）**：已登记 MFA 的账号，自助恢复是旁路第二因子、要求 TOTP、还是拒绝自助只走管理员重置；I-019-003 冻结：策略 UI 仅既有 settings tab（admin profile），强制面在 `core.auth-session` 对全 Profile 生效，禁止把 `admin.settings` 加入 mvp 默认集；I-019-001 对照 VP-018 已落地的 6 位码（`email_identity.go:31-36`）；I-019-006 按 2026-08-22 产品事实在 R1 标 verified；不得把验证码 vs 链接、邀请即建号、策略默认参数在 VP 层假装已选。
- evidence: VP-019:90-101；`email_identity.go:31-36`；`profile.go:26-93`；`LoginPage.tsx` MFA 错误码；P-005 门禁须与最晚阶段同构。
- closure: Root `00-meta` 含 P-001 阶段表 + 上述 I-00N（含 I-019-009 与 I-019-002 纠正）。不要求本 Review 落盘时已经有 MFA/TTL 答案。
- 建议 class: `editorial`

### V-F077 · recommended · 激活记录须留下 Admin 类 freshness 结论

- level: `recommended`；status: open；severity: low
- impact: 若激活只写「开区」而不点名：本 VP 非 Tier D、不改 Profile 意图、F-007 不升格、W11 已恢复 `go`、基线 `092bf37`、HEAD `66f5fd1f`，后续读者会把 IAM 误读成订单/CMS 解锁。
- finding: 激活时在 VP 短史或 lead Root D-001 写入上表复核结论与候选/HEAD 指针。
- evidence: VP-008 §`go` 消费有效性；VP-019 激活门闩；GOAL-007 D-001 原候选 `ed99e88`；workspace-009 GOAL-011 D-004 最近一次恢复；`kernel/profile.go` 默认集未改、`admin.account` 仅增 mail 路由（`profile.go:175-177`）。
- close requirement: D-001 或 VP 激活短史含 freshness 表；不要求重开 VP-008。
- 建议 class: `editorial`

### V-F078 · recommended · 组合索引补第 19 行；勿把 README 前向链接读成已开区

- level: `recommended`；status: open；severity: low
- impact: 「已落盘意图」表止于 VP-018，而组合焦点已指向 VP-019，读者会以为 IAM 尚未立项。
- finding: `/vision` 激活时在 roadmap 已落盘意图表追加第 19 行 VP-019（`planned`→`active` 随激活更新）；`workspaces.md` 仅在 `/govern` 开区后登记；README 的 `workspace-019-iam-recovery` 保持「开区后生效」直到目录真实存在。
- evidence: `roadmap.md` 已落盘表第 1–18 行；README 工作树前向链接。
- close requirement: 索引与 VP status 同拍；不要求本独立意见改文件。
- 建议 class: `editorial`

## 门禁含义

- 本 scope **open required = 0**。
- **允许**：用户确认后激活 VP-019、开新 delivery 工作区、按 V-F076/V-F077/V-F078 写 Root 纲领 / I-00N / freshness / 索引。
- **禁止**：把本 `pass` 写成忘记密码已交付、密码策略已可配置、或邀请已上线；重开 workspace-017/018；把 VP-019 当 Tier D 业务 VP 消费订单/CMS 解锁 scope；为密码策略改 mvp Profile 默认集。

## 结论

**VP-019 是否可以激活并开区？可以。** 用户确认后：`/vision` 落盘本 VRev、激活 VP、`/govern` 开 `workspace-019-iam-recovery`。开区后第一件事是 V-F076 的 R1 合同冻结，不要直接改恢复 DDL。

## 声明

本意见 `source: independent`，**不直接修改** Charter / VP / Goal status，不自行闭合任何 finding。required finding 的响应由 `/vision` 追加在正式 VRev 报告中；原 verdict 与 finding 原文不得改写。实施工作交 `/govern`。按用户本轮指令，本独立意见**只出报告文本、未写入** `docs/vision/reviews/VRev-043-*.md` 或 `reviews.md`。

## `/vision` 响应（2026-08-25 · V-F076/077/078 → fixed）

本 VRev 落盘与 VP-019 激活/开区同事务完成（用户本轮书面指令：「立项新vp（调用grok build做独立审计），并开设对应的工作区」）：

- **V-F076 → `fixed`**：lead Root `GOAL-001-iam-recovery` `00-meta` 写 P-001 纲领阶段表（R1 合同冻结 → R2 自助恢复全链 → R3 密码策略 + 邀请入职 → R4 证据/关门）；I-019-002 最晚阶段纠正为「R2 方案冻结前」；增补 **I-019-009 required（最晚 R1，MFA 与自助恢复）**；I-019-001 对照 VP-018 6 位码一致性默认候选、I-019-003 冻结策略面/Profile 边界、I-019-006 按 2026-08-22 产品事实标 verified（详见 VP-019 v0.2.0 与 Root 00-meta）。
- **V-F077 → `fixed`**：VP-019 激活短史 + Root `D-001-workspace-root-establishment` 写入 Admin 类 freshness 复核表（`ed99e88` → 基线 `092bf37` → HEAD `66f5fd1f`；PASS，不暂挂 `go`）。
- **V-F078 → `fixed`**：roadmap 已落盘意图表补第 19 行 VP-019（`planned → active`）；`workspaces.md` 在 `/govern` 开区后登记 workspace-019；README 指针随目录真实存在保持同拍。

原 verdict（pass）与 finding 原文未改写；本响应为 append-only 补充。