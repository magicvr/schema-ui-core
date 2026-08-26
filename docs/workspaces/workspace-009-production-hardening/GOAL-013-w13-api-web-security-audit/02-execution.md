---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-013-w13-api-web-security-audit
version: 0.1.0
---

# 执行索引 · GOAL-013

## 时间线（事实）

### E-001 · 立项与审计落盘（2026-08-26，S1 完成）

1. **审查执行**：用户指令"审视 api 和 web 的代码实现是否存在 bug 和安全漏洞"。会话内派出 4 个隔离上下文并行深审（①认证会话/MFA/captcha/限流/邀请/恢复；②store/kernel 持久层/wallet 钱包/recyclebin/settings；③upload/objectstore/import/mail 密钥/nginx/composition/config；④web 前端令牌传输/host/renderer/protocol/CSP），每路全文件通读 + 交叉 grep；编排会话另对 server/registrar/auth/session/rate-limit/login/recovery/upload/localStore/schema/invites/mfa/totp/service.go 核心面逐项复核。`go vet ./...` 干净。P1/P2 与关键 bug 均经编排会话二次读源码确认。
2. **落位与范围裁决**：结构化提问获用户书面选择——落位「workspace-009 · W13 波次」、范围「全部发现一次修完」。记录于 D-001。
3. **产物**：
   - `00-meta.md`（意图 + S1–S6 路线图，progress 来源登记）
   - `01-decision.md` + `01-decision/D-001-w13-scope-and-placement.md`
   - `03-audit.md` + `03-audit/A-001-w13-security-review-findings.md`（verdict: conditional；required = F-001～F-004）
   - `attachments/audit-A-001-findings-full.md`（逐条证据/场景/修复建议）
4. **goal-tree 同步**：workspace-009 `goal-tree.md` 增 GOAL-013 行与树节点。

**路线图状态**：S1 ✅；S2～S6 待启动（下一阶段：S2 API 必修批）。

### E-002 · S2 API 必修批实施完成（2026-08-26，checkpoint `9da0084e`）

1. **F-001（P1）**：`handler/invites.go` 公开 accept 面重排——新增 `authsession.Repository.PeekInviteToken`（单次索引 SELECT 的廉价存活检查，未知/过期/已用/已撤销一律统一 `INVITE_INVALID`，无枚举预言机），在密码策略与 bcrypt 之前执行；同路由挂每客户端 IP 滑窗限流（15min / 10 次 / 容量 64K，`loginRateLimiter` 模型），peek 失败与 Accept 失败记账、成功清桶。死 token 不再触发任何哈希开销。
2. **F-002/F-003（P2）**：`handler/mfa.go` 新增共享 step-up 失败桶（15min / 5 次 / 64K，key=`op|clientIP|userID`）覆盖三端点——`/api/mfa/enroll` 错误 currentPassword 记账（对照 account_self changePassword 模型）；`/api/mfa/disable`、`/api/mfa/recovery/rotate` 仅对 `MFA_INVALID`（第二因子证明失败）记账，成功清桶，超限 429 RATE_LIMITED + Retry-After。TOTP 猜测预言机关闭。
3. **F-004（P2·bug）**：`modules/mfa/service.go Confirm` 改持久化 ValidateTotp 返回的**匹配步进**（原为墙上时钟步进）：确认时匹配 ±1 邻步码后首登当前步码不再被回放水位误拒（原缺陷窗口 30–60s）。已知残余（与登录路径 AdvanceLastUsedStep 一致的语义）：确认时匹配"下一步"码（设备时钟超前）则该码之后的同窗码按高水位语义保守拒绝 ≤30s——记录于 A-001 响应待审。
4. **回归锁测试**：`TestServiceConfirmPersistsMatchedStep`（旧代码下必败）、`TestInviteAcceptDeadTokenShortCircuitsBeforePasswordWork`、`TestInviteAcceptRateLimitBoundsUnauthenticatedSpray`、`TestMFAStepUpDisableAndRotateRateLimited`、`TestMFAEnrollWrongPasswordRateLimited`。
5. **验证**：`go vet ./...` 0 输出；`go test ./... -count=1` 全绿 46 包（第 1 轮出现 1 次未定位非稳定失败 exit 1，相同代码第 2 轮全量复跑 46/46 ok、零 FAIL——判定 flake，如实记录）；checkpoint 提交 `9da0084e`（7 files，+322/-5）。

**路线图状态**：S1 ✅ S2 ✅；下一阶段 S3（API P3 与健壮性批；F-007/F-013 按 D-001 须用户三路径裁决）。

## 执行记录目录

| 编号 | 文件 | 内容 | 状态 |
|------|------|------|------|
| E-001 | （本文件时间线第 1 节） | 审查执行 + 立项 + A-001/D-001 落盘 | done 2026-08-26 |
| E-002 | （本文件时间线第 2 节） | S2 API 必修批（F-001～F-004）+ 回归锁 + checkpoint `9da0084e` | done 2026-08-26 |
