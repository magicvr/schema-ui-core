---
id: A-005
goal: GOAL-017-r3-s10-mfa-2fa
title: S1 方案 A-004 required 闭合复审
date: 2026-08-15
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-004 F-001/F-002 required 闭合验证（D-002 §2/§3/§4/§6 + D-003；状态机 / fail_count / last_used_step / admin reset / 组合根）
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# A-005 · 独立复审意见（S1 · A-004 required 闭合）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：finding-closure · A-004 F-001、F-002（required med）闭合验证；对照 D-002 §2/§3/§4/§6、D-003 与现网 auth/组合根
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（绑定与 `plan_refs`/`primary_plan` 已校验；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。同区 GOAL-011 仅作登录挑战先例只读对照。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001、D-002、D-003、E-002、`03-audit.md`、A-001～A-004。
- **代码核对**：`internal/auth/auth.go` `Login` L105–145（验密后 `issue`；L59 `RevokeAllRefreshTokensForUser`）；`handler/account_self.go` L5、L198（改密吊销先例）；`composition/composition_test.go` L465（admin 24/12）；`operationlog/migration/migration.go` L225–226（max Version 26）；`logincaptcha/store/repository.go` L50–55（一次性挑战 ≠ 5 次计数，对照用）。
- **covered**：A-004 F-001/F-002 关闭证据是否真实、充分、可实施；`status` 是否消自锁；`fail_count`/`last_used_step` 与声明控制同构；admin reset 是否与 disable 同级；组合根 24→25；方案自洽。
- **excluded**：S2 实现、S3～S5；不改 `status` / `progress` / goal-tree / D-002 正文 / `00-meta`。
- **保证等级**：L0。

## 成果（有证据）

| A-004 主张 / 本复审核对项 | 闭合证据 |
|---------------------------|----------|
| `user_mfa.status` pending/active；仅 active 触发登录 MFA | D-002 §2 L27；§4 enroll「行状态 pending 至 confirm」。原「行存在即启用」已删除 |
| enroll 后 pending 不触发 → 会话过期可再登录，消除自锁 | D-002 §2 L27「pending 不触发——防 enroll 后会话过期自锁」；`Login` 插在 L145 `issue` 前，`Required` 只读 active 可实施 |
| `mfa_proofs.fail_count`；达 5 立即失效 | D-002 §2 L28、§3 L37；与 captcha `ConsumeChallenge` 一次即删不同构，已显式落库计数 |
| `user_mfa.last_used_step`；同窗重放拒绝 | D-002 §2 L27、§6 L64（proof 单次 + lastUsedStep） |
| admin reset = 解绑 + `token_version+1` + 吊销全部 refresh | D-002 §4 表 L50；与 §1/§4 disable 同级；先例 `auth.go` L415–426、`account_self.go` L5 |
| 组合根 admin 权限 24→25、导航 12→12 | D-002 §6 L67；`composition_test.go` L465 = 24/12 |
| 迁移 0029/0030 在 016 的 0027/0028 之后、相对 max 26 不碰撞 | D-002 §2/§6；`operationlog` Version 26 |
| D-003 声明路径 | D-003 L24–27：状态机一致化 + reset 强化；闭合路径 **fixed** |

## 对照 A-004 required 闭合标准

| 标准 | 状态 | 证据 |
|------|------|------|
| §2 表结构与 §3/§4/§6 声明控制同构（status / fail_count / last_used_step） | **闭合** | D-002 §2 L27–28 与 §3 L37、§6 L64 一致 |
| 仅 active 触发，消除 enroll 后自锁 | **闭合** | D-002 §2 L27 |
| admin reset ≥ 自助 disable（token_version+1 + 吊销 refresh） | **闭合** | D-002 §4 L50 |
| 组合根 24→25 与 live snapshot 一致 | **闭合** | L465 = 24 权限；+`users.mfa-reset` → 25 |
| I-001/I-002 required 信息项未重新打开 | **满足** | 设计层仍成立 |

## Findings

### F-001 · pending 行再次 enroll 覆盖未写死（丢失一次性 secret 时的恢复）

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | A-004 F-001 最低清单含「过期 pending 可被新 enroll 覆盖」。D-002 §2/`user_mfa.user_id` PK（L27）+ §4 enroll 未规定第二轮 enroll 覆盖 pending。自锁已解除（可登录），但 secret/恢复码「仅此一次可见」（§4 L46）：丢失后无法 confirm，也不能再 enroll，只能走 admin reset。 |
| closure | — |
| 影响门禁 | 不阻断 S2 |

建议：pending 允许新 enroll 覆盖整行（或等价「放弃并重绑」）；`status=active` 时 enroll 必须先 disable（防已绑定会话被静默换 secret）。

### F-002 · A-004 F-003/F-004 仍开放（强制启用未裁定；proof 回传路径未写死）

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | ① I-003 仍标 verified、`00-meta` L24 仍写「强制启用候选」；D-002 §4 只有 admin **reset**，§7 未点名否决强制启用。I-003 为 non-blocking。② `Login` 签名仍为 `(access, refresh, user, err)`（`auth.go` L105）；错误时 user 为空（L108–138）。D-002 §3 L36 要 200 `{mfaRequired, mfaProof}`，仍未写 proof 由哨兵携带 / 改签名 / handler 在已知 userID 后签发。S2 清单第 1 条仍未补这一句。 |
| closure | — |

不恢复为 required。S2 前各补一句即可：本波不纳入强制启用；handler/mfa 在验密成功后签发 proof，`Login` 返回可识别哨兵 + userID。

### F-003 · 组合根「26→30」仍为两目标合计；jwt 轮换残余未登记

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | D-002 §6 L67「迁移 26→30」是 016+017 合计；017 单独为 26→28 或 28→30。权限 24→25 已正确。A-004 F-005：HKDF(`jwt` 同源 secret, `"mfa/totp"`) 可接受，但 jwt_secret 轮换会使全部 TOTP secret 不可解密——D-003 未登记该残余。 |
| closure | — |

S2 按 live snapshot 改断言；轮换程序或书面 residual 择一留痕。

## 必改项汇总

无 required。A-004 F-001、F-002 已合法闭合（fixed）。

## 与既有意见的异同

- A-004 independent conditional：开放 F-001（模型/控制不同构 + 自锁）、F-002（reset 弱于 disable）。本意见核对 D-002 修正后：**两条 required 均闭合**。
- A-004 F-003～F-005（recommended）部分仍开放（本意见 F-002/F-003），不升级、不阻断 S2。
- A-003 self pass 与本复审不冲突。
- D-003 声称全 fixed：就 A-004 **required** 范围成立。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass**。A-004 F-001 / F-002（required med）关闭证据充分、可重复核对：pending/active 消自锁，`fail_count`/`last_used_step` 与声明控制同构，admin reset 已与 disable 同级，组合根 24→25 与 L465 一致。无 high/med required；无到期 required 信息项。

**可放行 S2 实施**（含写 0029）。本意见 F-001～F-003 为 recommended，随 S2 一并处理。

建议 `/govern`：响应本意见（记录 A-004 F-001/F-002 → fixed，本 A-005 recommended 带入 S2）后开 S2。勿用 `progress: 1/5` 作为放行依据。

## 声明

本意见不修改 `status` / `progress` / goal-tree / 方案正文。响应由 `/govern` 处理。保证等级 L0。
