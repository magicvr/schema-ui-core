---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-013-w13-api-web-security-audit
version: 0.1.0
---

# A-002 · W13 S6 关门前 self 审计（编排器自审）

- **source**: self（编排器，`04-write-audit` 路径；独立腿按 D-001/S6 由 grok build 另行执行）
- **日期**: 2026-08-26
- **scope**: GOAL-013 全分母——A-001 required（F-001～F-004）、P3 加固（F-005～F-020）、健壮性（B-1～B-4）在 checkpoint `9da0084e` / `b7954235` / `e93f7228` 的处置真实性核对；D-001/D-002 用户裁决留痕齐备性；回归证据复核
- **verdict**: **pass**（开放 required = 0）

## 逐条核对

| 编号 | 处置 | 证据 | 自核结论 |
|------|------|------|----------|
| F-001 | fixed | `handler/invites.go`（PeekInviteToken 先于 bcrypt + 每 IP 滑窗）；`authsession/invites.go::PeekInviteToken`；测试 `TestInviteAcceptDeadTokenShortCircuitsBeforePasswordWork`、`TestInviteAcceptRateLimitBoundsUnauthenticatedSpray` | ✓ 死 token 零哈希开销；限流与登录模型一致 |
| F-002/F-003 | fixed | `handler/mfa.go` mfaStepUp* 常量 + `guardMFAStepUp`；enroll/disable/recovery-rotate 接线；测试 `TestMFAStepUpDisableAndRotateRateLimited`、`TestMFAEnrollWrongPasswordRateLimited` | ✓ 仅 MFAInvalid 记账、成功清桶、超限 429+Retry-After |
| F-004 | fixed | `mfa/service.go::Confirm` 持久化匹配步进；测试 `TestServiceConfirmPersistsMatchedStep`（旧代码必败） | ✓ 见"自审备注 1"残余说明 |
| F-005 | fixed | `mfa/totp.go` subtle.ConstantTimeCompare | ✓ 与 recovery/email 同强度 |
| F-006 | fixed | `config.go` AUTH_PUBLIC_BASE_URL（yaml/env/启动校验）→ users provider → `inviteLink`；`.env.example` 登记 | ✓ 配置后头派生路径不可达；未配置保持旧行为（文档化） |
| F-007 | fixed（承载子目标） | D-002 决策 1；GOAL-014 五件套（parent=GOAL-013） | ✓ 非 required，不阻断本波关门；关门顺序待用户裁决 |
| F-008 | fixed | `authsession/recovery.go` 过期分类仅哈希匹配时给出 | ✓ 统一无效面，无状态探测 |
| F-009 | fixed | `email_identity.go::BindEmail` 全地址派发冷却；测试更新为新契约 | ✓ 见"自审备注 2"语义变更说明 |
| F-010 | fixed | `handler/schema.go` 认证挂载；`TestSchemaEndpointRejectsAnonymous`；composition 投影测试改双 mux | ✓ 匿名统一 401；启用/禁用区分保留于认证后 |
| F-011 | fixed | wallet store ×2 + recyclebin store ×1 escapeLikePattern + ESCAPE '\' | ✓ 两方言可移植 |
| F-012 | fixed | `handler/wallet.go` OwnerExistsFunc；composition 接线 UserByID | ✓ 显式创建与 adjust 双门；nil 仅测试环境 |
| F-013 | accepted-residual | D-002 决策 2（用户书面），复审触发硬登记 | ✓ 合法闭合路径三之一 |
| F-014～F-016 | fixed | `auth-client.ts` isSameOrigin；`boot.ts` isSafeSupportUrl；`host/claim.ts` visited 集；三组回归锁测试 | ✓ vitest 1128/1128 |
| F-017 | fixed | `MAIL_MASTER_KEY_PATH`（yaml/env/.env.example）；composition 接线 | ✓ 默认行为不变，分置可选 |
| F-018 | fixed | `raster_assets.go` Cache-Control max-age=300；两处断言同步 | ✓ 删除分钟级生效 |
| F-019 | fixed | `import.go` load 后立即 best-effort Delete（uploads namespace） | ✓ 明文密码不再驻留可下载面 |
| F-020 | 部分 fixed（经裁决） | nginx HSTS；img-src https: 保留依据 normalizeLogoURL 功能证据 + D-002 决策 3 | ✓ 用户书面裁决留痕 |
| B-1～B-4 | fixed | email_identity RowsAffected 分支；filelibrary Stat 单遍；keyedMutex（upload/avatar）；mail.ErrInvalidRetention sentinel | ✓ 各带说明注释 |

## 回归证据

- API：`go vet ./...` 0 输出；`go test ./... -count=1` 46 包全绿（修复 op-id 顺序 flake 后复跑，日志存查）。
- Web：vitest 1128/1128（83 文件）+ `npm run build` 成功。
- Checkpoints：S2 `9da0084e`；S3/S5 `b7954235`；S4 `e93f7228`；docs `9337fd84`/`e5fe5829`。

## 自审备注（提交 independent 审计重点复核）

1. **F-004 残余语义**：确认时若匹配"下一步"码（设备时钟超前 >0<30s），高水位将保守拒绝该码之后的同窗码 ≤30s——与登录路径 AdvanceLastUsedStep 既有语义一致，属 TOTP 高水位模型的固有权衡而非回归。
2. **F-009 有意行为变更**：GOAL-002 时代"换址重绑即时发信"被审计判定为邮件炸弹原语；现所有派发共享每账号 60s 冷却。状态机覆盖语义（overwrite→pending、旧槽释放）不变。此为审计驱动变更，非静默偏离。
3. **operationID 时间有序化**：非 A-001 分母内的附带稳定化，先例 GOAL-037/F-008（wallet），因全量回归暴露同粒度写入顺序随机而实施。

## 结论

A-001 全部 required（F-001～F-004）已 genuine fixed 且有缺陷形状回归锁；P3/健壮性分母按用户裁决全部处置并留痕；开放 required = **0**。具备进入 independent 审计腿条件。
