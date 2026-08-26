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

### E-003 · S3 API P3 与健壮性批 + S5 部署批实施完成（2026-08-26，checkpoint `b7954235`）

1. **用户三项裁决（D-002 落盘）**：F-007=fixed（承载于新建子目标 GOAL-014）；F-013=accepted-residual（复审触发：首个 self-scope 生产角色上线前必须谓词化修复）；F-020=保留 img-src https:（远程品牌图为受支持功能，normalizeLogoURL 证据）+ 实施 HSTS。
2. **API 修复**：F-005 TOTP 比较改 `subtle.ConstantTimeCompare`；F-008 `EvaluateRecoveryCode` 过期分类仅在哈希匹配时给出 EXPIRED（其余统一 INVALID，消除"近期申请过重置"探测预言机）；F-009 `BindEmail` 对**任意**地址派发执行每账号 60s 冷却（原换址即时发信=邮件炸弹原语；frozen 重绑语义不变）；B-1 `registerFailedAttempt` 用 UPDATE RowsAffected 区分"无挑战行"与"已作废"（并发下不再误报 EXPIRED）；F-010 `/api/schema/{pageId}` 挂认证中间件（匿名统一 401，禁枚举页面文档）；F-011 wallet×2 + recyclebin×1 的 q LIKE 通配转义（ESCAPE '\\'）；F-012 `WalletRoutes` 新增 OwnerExistsFunc 门（composition 以 authRepository.UserByID 接线），显式创建/adjust 前校验 owner 存在；B-2 filelibrary Get 改 Stat 单遍（不再整读 body）；B-3 新增 `keyedMutex`，上传配额与头像配额从全局锁改为分 owner 锁（W7 F-012/W11 F-018 不变量保持）；B-4 mail 新增 `ErrInvalidRetention` sentinel，handler 改 errors.Is。
3. **S5 部署/运维**：F-017 新增 `MAIL_MASTER_KEY_PATH`（yaml/env）支持密钥文件与数据目录分置（默认行为不变，备份同泄风险文档化）；F-018 头像/品牌图 Cache-Control 从一年 immutable 收敛为 max-age=300（删除分钟级生效；内容寻址正确性不依赖 immutable）；F-020 nginx 增加 HSTS（RFC 6797：HTTP 上被忽略、TLS 接入自动生效；I-001 拓扑仍属运维侧）。
4. **附带稳定化（非 A-001 分母，如实记录）**：全量回归第 2 轮暴露 `TestRolesOperationLogEvents` 顺序翻转——根因为 `newOperationID` 纯随机 id 使 `created_at DESC, id DESC` 同粒度写入顺序随机；按 GOAL-037 F-008 先例改毫秒前缀+单调计数器+随机尾（`resources.go`）。该 flake 同时解释了 E-002 第 5 点的第 1 轮未定位失败。
5. **验证**：go vet 0；go test ./... -count=1 全绿 46 包（修复后复跑）；checkpoint `b7954235`（33 files，+388/-117）。

### E-004 · S4 Web 前端批实施完成（2026-08-26，checkpoint `e93f7228`）

1. **F-014**：`authFetch.withAuth` 增加 `isSameOrigin` 守卫——Authorization 与 X-Refresh-Token 仅附同源目标（原仅按 pathname 判定，绝对跨源 URL 会带凭据）；Accept-Language 不受限。
2. **F-015**：`executeBootRecovery("support")` 增加 `isSafeSupportUrl` scheme 守卫（仅 http/https，相对路径按 origin 解析）——关闭潜伏 javascript: 执行点。
3. **F-016**：`validateClaim` 依赖游走补 visited 集——registry dependsOn 成环不再死循环，共享 DAG 依赖指数膨胀消除；INCOMPLETE 判定语义不变。
4. **回归锁测试**：`claim-dependency-walk.test.ts`（环图终止 + 未列依赖仍拒）、`boot-recovery-url.test.ts`（危险 scheme 全拒 + 不构建 anchor）、`auth-client.test.ts` 跨源不带凭据用例。
5. **验证**：vitest 全量 1128/1128（83 文件）+ `npm run build` 全绿；checkpoint `e93f7228`。

**路线图状态**：S1 ✅ S2 ✅ S3 ✅ S4 ✅ S5 ✅；下一阶段 S6（审计闭合与关门）。

## 执行记录目录

| 编号 | 文件 | 内容 | 状态 |
|------|------|------|------|
| E-001 | （本文件时间线第 1 节） | 审查执行 + 立项 + A-001/D-001 落盘 | done 2026-08-26 |
| E-002 | （本文件时间线第 2 节） | S2 API 必修批（F-001～F-004）+ 回归锁 + checkpoint `9da0084e` | done 2026-08-26 |
| E-003 | （本文件时间线第 3 节） | S3/S5 API P3+健壮性+部署批（F-005～F-012/F-017～F-019/B-1～B-4）+ D-002 三裁决 + checkpoint `b7954235` | done 2026-08-26 |
| E-004 | （本文件时间线第 4 节） | S4 Web 批（F-014～F-016）+ 回归锁 + checkpoint `e93f7228` | done 2026-08-26 |
