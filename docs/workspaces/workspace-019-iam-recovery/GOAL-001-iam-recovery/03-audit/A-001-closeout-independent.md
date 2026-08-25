---
id: A-001
doc: audit-entry
goal: GOAL-001-iam-recovery
source: independent
auditor: ox-alpha（DSH harness 独立交叉审计会话）
audit_type: close-out
verdict: pass
created: 2026-08-26
updated: 2026-08-26
version: 1.0.0
---

# A-001 · workspace-019 全链关门独立复核（close-out · 代码证据制）

## 范围与区间

- **被审对象**：workspace-019-iam-recovery 完成情况整体（Root GOAL-001 及 R1～R4 四个子目标的关门主张）。
- **审计模式**：`close-out`（关门复审）。
- **证据规则**：按用户本轮明确要求，**完成证据只采信代码实现与实际运行结果**；治理文档仅用于提取「待核对主张清单」，不作为任何完成事实的依据。
- **区间**：R2 `299f8f52`/`9628ca8f`、R3 响应 `2f088d55`、关后维护 `b9ae434f`～HEAD `9ced003d`（四个提交哈希均经 `git log` 核实存在且信息吻合）。
- **工作区校验**：`workspace.md` id/root_goal/canonical 一致，`shared_materials_catalog: none`，无跨区资料引用问题。

## 成果（有证据）

以下每条均由本审计员**独立读取代码或实际运行命令**取得，非转抄治理记录：

### E-1 · 测试全量复跑（决定性证据）

| 运行 | 结果 |
|------|------|
| `go test ./...`（apps/api 全量复跑） | **47 包全部 ok，exit 0，零 FAIL**。关键包实测：handler 47.8s、store 40.8s、composition 26.3s、authsession 18.4s、mfa 17.9s |
| `npx tsc -b`（apps/web 复跑） | **exit 0，无输出** |
| `npx vitest run`（apps/web 复跑） | **81 文件 / 1113 tests 全部通过，exit 0**（关门时主张 1105；差额 +8 来自关后维护新增组件测试，与 E-008 时间线一致，非矛盾） |

### E-2 · R1 冻结合同 ↔ 代码常量逐项一致

- 6 位码 + sha256 落库 + 恒时比较：`authsession/recovery.go` L4/L27-31（`crypto/subtle`）、invites.go L65（`sha256.Sum256`）。
- TTL 10min / 重发冷却 60s / ≤5 次作废：recovery.go L28-30 三常量逐字吻合。
- 策略边界 8–72 / 0–4 / 0–10：password_policy.go L24-27 与 handler 校验 L117-119 双侧一致。

### E-3 · R2 自助恢复全链

- 迁移 0056 挑战表存在于 `modules/authsession/migration/migration.go` L534（0057/0058/0059 同文件 L542-558）。
- `handler/recovery.go`：start 对无自助路径账号返回同形 202 且不投递（L100-116 防枚举）；complete 门序 = 码匹配→MFA 第二因子（L185-197）→密码基线→策略校验（L208）；错码与第二因子失败共用挑战预算（failAttempt/failAttemptMFA → ConsumeRecoveryAttempt）；IP|identifier 限流桶贯穿全部失败路径。
- `mfa.Service.VerifySecondFactor` 存在且有专项测试 `mfa/recovery_gate_test.go`（nil fail-closed、TOTP/恢复码双路、重放拒绝）。
- operational allowlist 含 `/api/auth/recovery/start|complete`（operational.go L75-76）。
- Web：LoginPage 两步恢复流含第二因子字段（LoginPage.tsx L292/L632）；i18n `login.recovery.*` **恰 19 键** zh/en 对称。

### E-4 · R3 密码策略四口强制 + 邀请全链

- 四口接线实证：users Create（handler/users.go:167）、users Update password（:203）、account_self changePassword（account_self.go:260）、recovery complete（recovery.go:208），统一 `INVALID_PASSWORD`。
- 配置 MinLength 权威生效（password_policy.go L86-89 clamp 至 8 下限）——GOAL-004 A-001 F-001 的修复确实在码；专项测试 `TestValidateNewPasswordConfiguredMinLengthBites` 在 authsession 包绿测中。
- 历史捕获/裁剪随 UpdateUser 同事务（capturePasswordHistory + bcrypt CompareHashAndPassword 比对）。
- settings 面 GET=settings.read / PATCH=settings.write 分权（password_policy_settings.go L38-39）——F-004 修复在码。
- 邀请域：单事务激活（withTx "accept invite"）、角色 revalidate → INVITE_ROLE_GONE、用户名冲突 fail-closed、resend 撤旧发新 + 60s 冷却、激活不签发会话；公开 accept 中央注册（invites.go L272）并入 allowlist（operational.go L79）。
- 管理 四 路 `users.invite` 权限校验 ×4（handler/invites.go L111/196/223/236）；kernel/profile.go admin.users 声明同步（L166）。
- 错误码 INVALID_INVITE_BODY / INVITE_INVALID / INVITE_ROLE_GONE + 恢复面 4 码均入契约冻结集（error_contract_test.go L114-122）且 errorcatalog 有中英文案。
- composition 黄金计数 mvp wantPermissions=11 / admin=33（composition_test.go L462/486）与声明吻合；mvp +1 来源为**已在本集内**的 admin.users 模块内部新增 users.invite 权限——属模块内贡献增量，未改 Profile 默认集/模块矩阵，无越界。
- Web：password-policy-tab 自注册（L137 registerCustomComponent）、invite-issue-card/resend-dialog 组件及测试、`/invite/accept` 未认证分支（main.tsx L96）、users.json create-user roles checkboxGroup（optionsSource /api/roles）+ `createUser.bodyMapping.roles`（E-008 缺陷修复在码，users.json L27-31/L195-200）。

### E-5 · R4 关门证据独立可重复

- `r4_evidence_test.go` 三链均为真实 HTTP mux + mock 渠道出站取码（恢复 bind→verify→start→取码→complete→新密码登录 200；邀请建邀→取链接→accept 204→viewer 登录→一次性回放 INVITE_INVALID；策略 PATCH minLength=12→弱码 400→强码成功），并带 viewer 无权限 403 断言（A-001 F-001/F-002 闭环）。该测试在全量复跑中通过。
- store 黄金断言（identity/migrate/catalog/restart）随全量复跑绿，迁移目录一致性获独立验证。

## 对照成功标准（VP-019 三件）

| 交付 | 判定 | 依据 |
|------|------|------|
| 密码策略（配置面+强制面） | **达成** | E-2/E-4：配置 API + 四口强制 + 渐进边界（无存量扫描） |
| 邀请入职 | **达成** | E-4：双形态投递 + 即建号 + 完整生命周期 + Web 面 |
| 自助恢复状态机 | **达成** | E-3：全链 + MFA 门 + 会话撤销语义（CompleteRecovery→UpdateUser token_version） |

## Findings（F-00N；均 recommended，无 required）

| ID | 级别 | 严重度 | 描述 | 证据路径 |
|----|------|--------|------|----------|
| F-001 | recommended | low | `password_policy_settings.go` 末行 `var _ = errors.Is // keep errors imported for future sentinel mapping` 为死导入保持行；PATCH 失败路径一律映射 INTERNAL，无 sentinel 细分。纯卫生问题，不影响行为 | apps/api/internal/handler/password_policy_settings.go L136 |
| F-002 | recommended | info | loginRateLimiter 为进程内内存桶（15min/20 次/IP\|identifier），多实例部署时限流预算按节点各自计算。与既有 login 面同型（本区未引入新模式），不在 workspace-019 边界内；登记为部署拓扑注意项，供后续生产化波次评估 | apps/api/internal/handler/recovery.go L58 |

两条均不构成关门障碍，不阻断任何门禁。

## 必改项汇总

无 required 必改项。

## 与既有意见的异同（若有 self/independent 历史）

- 各子目标既有 independent/self 意见（GOAL-003 A-001/A-002、GOAL-004 A-001 conditional→F-001~F-004 fixed→A-002 pass、GOAL-005 A-001→F-001/F-002 fixed→A-002 pass）与本审结论**方向一致**；本审进一步把「测试绿」从文档主张升级为**本会话实际复跑结果**，并把 F-001/F-002 的修复逐一指认到代码行级证据。
- 本审补充发现 F-001/F-002（recommended），此前各轮意见未覆盖，属增量卫生/部署观察，非冲突。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass。** workspace-019 的 done 4/4 关门主张在代码证据制下成立：三条交付链（密码策略/邀请入职/自助恢复状态机）均有真实实现、真实测试且本审独立复跑全绿；冻结合同数值与代码常量一致；无越界；开放 required = 0。两条 recommended findings 可留待后续维护波次顺手处理，无需重开本工作区。

建议：如需响应 F-001/F-002 或推进后续波次，使用 **`/govern`** 编排处理；本意见不改任何 status/progress。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
