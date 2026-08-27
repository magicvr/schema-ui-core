---
title: 目标树 · workspace-009-production-hardening
status: active
created: 2026-08-10
updated: 2026-08-26
parent: null
version: 0.24.0
workspace_id: workspace-009-production-hardening
---

# 目标树 · 生产加固（共享基架持续安全与健壮性）

> 工作区：`workspace-009-production-hardening`
> canonical：`docs/workspaces/workspace-009-production-hardening/`
> Root：`GOAL-001-production-hardening`，**长期程序容器 · active**
> primary_plan：`VP-009-production-hardening`，**active**
## 树
```text
GOAL-001-production-hardening [active]  · 持续安全程序
├── GOAL-002-audit-findings-remediation [done] (16/16)    · W1
├── GOAL-003-upload-ownership-hardening [done] (4/4)      · W2
├── GOAL-004-w3-security-audit-remediation [done] (8/8)   · W3
├── GOAL-005-w4-security-audit-remediation [done] (8/8)   · W4
├── GOAL-006-w6-scan-findings-remediation [done] (4/4)    · W6
├── GOAL-007-w7-api-web-security-audit [done] (4/4)      · W7
├── GOAL-008-w8-api-web-security-audit [done] (4/4)      · W8
├── GOAL-009-w9-api-web-security-audit [done] (4/4)     · W9
├── GOAL-010-w10-api-web-security-audit [done] (4/4)    · W10
├── GOAL-011-w11-api-web-security-audit [done] (4/4)  · W11
├── GOAL-012-w12-multi-instance-rate-limiting [done] (4/4) · W12 · 多实例限流拓扑评估（单实例边界维持 · 载体预登记 Redis）
├── GOAL-013-w13-api-web-security-audit [done] (6/6)      · W13 · api/web 全量安全审查发现修复（P1×1+P2×3 必修 + P3 全量）
│   └── GOAL-014-w13-account-lockout-redesign [done] (6/6) · W13-F007 账号锁定模型重设计（fixed · 用户裁决承载子目标）
└── GOAL-015-w14-schema-auth-wiring-lock [done] (4/4)     · W14 · 页面 Schema 鉴权装配修复与生产装配回归锁
```

**W14（2026-08-26 开波，2026-08-26 关门 · done 4/4）**：用户报障「所有页面都显示无法显示此页面」→ 定位为 GOAL-013 F-010（`/api/schema` 挂认证，checkpoint `b7954235`）之后**生产入口缺 `schemaFetcher` 认证传输**——页面文档请求全部匿名 401，全站渲染 PageSchemaErrorSurface；30+ 测试均显式注入 fetcher，测试全绿与生产断裂并存（「测试装配 ≠ 生产装配」）。hotfix `schemaFetcher={authFetch}` 经用户确认先行落地；用户追加指令补防回归锁并在本区立项承载。**S2** 冻结双措施：hotfix 追认 + AuthGate 提取为可测模块与生产装配回归锁（[D-001](GOAL-015-w14-schema-auth-wiring-lock/01-decision/D-001-w14-scope-and-freeze.md)；三项替代方案有据否决）。**S3** 实施落地：`AuthGate.tsx` 提取 + `auth-gate.wiring.test.tsx` ×2，全量 vitest **1130/1130**（84 文件）+ `tsc -b` 0。**S4 复核关门**：self [A-001](GOAL-015-w14-schema-auth-wiring-lock/03-audit/A-001-w14-self-closeout.md) `pass`（含变异验证红→绿）；用户指令「把 R-001 并入本波。处理完再关门」→ **R-001 fixed**：新增 `e2e/schema-auth-transport.spec.ts` 真实网络层 Bearer 冒烟（[E-003](GOAL-015-w14-schema-auth-wiring-lock/02-execution/E-003-w14-r001-e2e-bearer-smoke.md)）；关门验证附带修复 shell.spec 匿名 schema 探测陈旧契约（F-010 后匿名恒 401；定性非本波回归，[E-004](GOAL-015-w14-schema-auth-wiring-lock/02-execution/E-004-w14-shell-spec-contract-fix.md)；根因=e2e 自 F-010 起未再完整运行）；最终全量绿：vitest 1130/1130 + Playwright chromium **10 passed / 1 skipped · exit 0**。**[D-002](GOAL-015-w14-schema-auth-wiring-lock/01-decision/D-002-w14-closeout.md) 用户书面关门：done (4/4)**。Root 保持 active。见 [GOAL-015](GOAL-015-w14-schema-auth-wiring-lock/00-meta.md)。
**W13（2026-08-26 开波，2026-08-26 关门 · done 6/6）**：用户指令审视 api/web bug 与安全漏洞 → 会话内 4 路隔离上下文并行审查 + 核心链路复核 → A-001 `source: independent` **conditional**（required = F-001～F-004；P3 F-005～F-020 + B-1～B-4 全量分母）。用户裁决：落位 W13 + **范围=全部发现一次修完**（D-001）。**S2**：F-001 先验 token+每 IP 限流、F-002/F-003 MFA 三端点第二因子失败限流、F-004 Confirm 匹配步进；checkpoint `9da0084e`。**S3/S5**：D-002 三项用户裁决（F-007=fixed→子目标 GOAL-014；F-013=accepted-residual+Root E-008 复审硬门；F-020=保留 img-src https:+HSTS）；其余 P3/健壮性全量 fixed；附带 operationID 时间有序化；checkpoint `b7954235`。**S4**：F-014 同源守卫、F-015 scheme 守卫、F-016 环图死循环修复；vitest 1128/1128 + build 绿；checkpoint `e93f7228`。**S6 审计腿**：A-002 self pass → A-003 grok-build independent **pass**（required ×4 genuine fixed、回归独立复跑全绿）→ A-004 响应（R×3 闭合：Root E-008 / 叙事约束 / I-001 信息项）。**GOAL-014 承载 F-007 fixed**：分层锁定模型（来源锁 5/15min 仅拒施害来源 + 全局天花板 100/24h 滑动 + 失败零吊销），迁移 0061，真实 PG 方言复核，A-002 grok-build independent **pass**（F-007 genuine fixed）；checkpoint `26655b55`。**D-003/D-004 用户书面关门：两目标一并 done（2026-08-26）**。残余移交：I-001 TLS 拓扑（deferred non-blocking）、F-013 TOCTOU 复审硬门（Root E-008）、GOAL-014 A-003 R-F002 Refresh 残余。Root 保持 active。见 [GOAL-013](GOAL-013-w13-api-web-security-audit/00-meta.md) / [GOAL-014](GOAL-014-w13-account-lockout-redesign/00-meta.md) / [D-004](GOAL-013-w13-api-web-security-audit/01-decision/D-004-w13-closeout.md)。
**W11（2026-08-22，已关门）**：用户指令独立审计 api/web（禁止加载 skills）→ 开 GOAL-011；A-001 `source: independent` **fail**（3 HIGH required + 3 MEDIUM required = 6 开放；13 recommended；6 informational）。用户轮次指令「推进…直到顺利闭门」授权闭门路径：**D-002 整单采纳 6 条 + 波内暂挂 VP-008 go**（不因 fail 单点自动悬挂的原则保留）→ E-002 S3 实施 6/6 修复（F-001 PG EXISTS bool；F-002 删除+快照同事务；F-003 MFA proof 限流+原子计数；F-004 MFA 双密钥轮换窗口；F-005 验证码单语句消费；F-006 对账 400+写权限）→ E-003 recommended 处置（fixed 11：F-007/008/010/011/012/013/014/016/017/018 + F-009 部分；overruled 2 有据：F-015 客户端跨标签重试协议依赖、F-019 API 硬门禁不变）→ **A-002 self pass** → **A-003 grok-build（grok-4.6 · reasoning high）independent pass：6/6 genuine fixed + F-001/F-005 真实 Postgres 复跑**（回归全绿：go vet 0 / go test 全绿 / web 1085/1085 / tsc 0；checkpoint `72a5397`）→ A-004 闭合记录（fixed ×6 + overruled ×2 + recommended ×11；开放 required = 0；I-003 关闭）→ **D-004 关门 + 恢复 VP-008 go 宣称**。`status: done` (4/4)。**关门后复核：A-005（independent · DeepSeek Harness）post-close pass（代码直接核验 6/6 genuine fixed；无必改项）→ A-006 响应（R-002：03-audit.md frontmatter 转 done；R-001/R-003/I-A/I-B 有据记录）→ 正式关门确认，无重新打开条件。** Root 保持 active。残余移交：数据库密码轮换（用户侧）；R-001/R-002 与 F-009 lastRun（见 A-004）。见 [GOAL-011](GOAL-011-w11-api-web-security-audit/00-meta.md) / [D-002](GOAL-011-w11-api-web-security-audit/01-decision/D-002-w11-scope-and-go-hold.md) / [A-003](GOAL-011-w11-api-web-security-audit/03-audit/A-003-w11-s4-independent.md) / [A-004](GOAL-011-w11-api-web-security-audit/03-audit/A-004-w11-closure-response.md) / [D-004](GOAL-011-w11-api-web-security-audit/01-decision/D-004-w11-go-restore.md) / [A-005](GOAL-011-w11-api-web-security-audit/03-audit/A-005-w11-post-close-independent.md) / [A-006](GOAL-011-w11-api-web-security-audit/03-audit/A-006-w11-a005-response.md)。
**W9（2026-08-21，已关门）**：A-001 fail + A-002 conditional → D-002 调和 required **12**（F-003 作废，全文 P2-7 补号 F-025）→ D-003 整单采纳 + 暂挂 VP-008 go → E-004 S3 实施 12/12 修复 + 回归全绿（A-004 self pass）→ **A-005 independent/pass（grok-build · grok-4.6 · reasoning high · `/audit`）：12/12 genuine fixed、无新引入缺陷、回归复跑一致** → A-006 闭合记录（fixed ×12，开放 required = 0，I-003 verified）→ E-005 将 A-005 三条 recommended 全部实施并锁定（L2 校验器接线生产路径、恢复码 CAS 换值令牌、6 组原缺陷形状回归锁；API/Web/build 再全绿）→ **D-004 用户书面恢复 VP-008 go 宣称**。`status: done` (4/4)。Root 保持 active。见 [GOAL-009](GOAL-009-w9-api-web-security-audit/00-meta.md) / [D-004](GOAL-009-w9-api-web-security-audit/01-decision/D-004-w9-go-restore.md) / [E-005](GOAL-009-w9-api-web-security-audit/02-execution/E-005-w9-recommended-hardening.md) / [A-005](GOAL-009-w9-api-web-security-audit/03-audit/A-005-w9-s4-independent.md) / [A-006](GOAL-009-w9-api-web-security-audit/03-audit/A-006-w9-a005-response.md)。
**W10（2026-08-21，已关门）**：用户指令独立审计 → A-001 `source: independent` **conditional**（1 HIGH required = F-001 env.example 硬编码凭据；6 recommended；5 informational）→ D-002 整单采纳 7 条 + 暂挂 VP-008 go → S3 实施前逐条核实，**D-003 调和 4 条 recommended 作废**（F-003 不成立 / F-004、F-005、F-006 误报），实际修复 F-001/F-002/F-007 → E-002 回归全绿（go vet/test 全绿；web 1083/1083 + build 全绿）→ A-002 self pass → **A-003 independent/pass（grok-build · grok-4.6 · reasoning high · `/audit`）：3/3 genuine fixed、4 作废有据、无新引入缺陷；新增 recommended ×3** → 用户 `/govern` 书面闭合授权：E-003 将 A-003 三条 recommended 全部修正（索引同步 / fetch-timeout listener 清理 / 预览窗 opener 置空；web 复跑 1084/1084 + tsc 0）→ A-004 闭合记录（fixed ×6 + 作废 ×4，开放 required = 0，I-003 关闭）→ **D-004 关门 + 恢复 VP-008 go 宣称**。`status: done` (4/4)。Root 保持 active。残余移交：数据库密码轮换（用户侧）。见 [GOAL-010](GOAL-010-w10-api-web-security-audit/00-meta.md) / [A-003](GOAL-010-w10-api-web-security-audit/03-audit/A-003-w10-s4-independent.md) / [A-004](GOAL-010-w10-api-web-security-audit/03-audit/A-004-w10-closure-response.md) / [D-004](GOAL-010-w10-api-web-security-audit/01-decision/D-004-w10-go-restore.md)。
**W8（2026-08-20，已关门）**：独立代码审计 A-001（`source: independent`）判定 **fail**（F-001 分页整数溢出/切片 panic-DoS、F-002 生产 CSP 阻止 inline 主题脚本为 2 条 required；F-003/F-004 为非阻断/条件风险）。用户目标轮次指令授权修复并闭门（D-002 整单采纳 F-001/F-002 + 暂挂 VP-008 go 宣称）；实现 E-002 修复 + E-003 处置 recommended；self A-002 pass + independent A-003 pass（grok-4.6）确认 required 0 开放；D-003 恢复 VP-008 go 宣称。`go test ./...`、web `npm test`（1072）、`npm run build` 全绿。`status: done`。Root 保持 active。见 [GOAL-008](GOAL-008-w8-api-web-security-audit/00-meta.md) / [D-002](GOAL-008-w8-api-web-security-audit/01-decision/D-002-w8-scope-and-go-hold.md) / [A-003](GOAL-008-w8-api-web-security-audit/03-audit/A-003-w8-independent.md)。

**W7（2026-08-19，已关门）**：独立代码审计落盘（A-001 fail，12 required）；用户确认整单采纳 F-001～F-012 并暂挂 go 宣称（D-002）。已实施全部 12 条 required（E-002/E-003），self A-002 pass，independent A-003 conditional 指出 F-006 限流未 record → 修正（E-003）后 independent A-004 **pass**（12/12 required 闭合）。**A-005 独立代码复核确认 F-001/F-002/F-006 genuine fixed + 恢复 VP-008 go 宣称（D-003）**；A-003 recommended F-002/F-003 修复 + F-004/F-005 处置（E-004）。`go test ./...` 与 web `npm test` 全绿。`status: done`。Root 保持 active。见 [GOAL-007](GOAL-007-w7-api-web-security-audit/00-meta.md) / [E-002](GOAL-007-w7-api-web-security-audit/02-execution/E-002-w7-implementation.md) / [A-002](GOAL-007-w7-api-web-security-audit/03-audit/A-002-w7-self.md) / [A-004](GOAL-007-w7-api-web-security-audit/03-audit/A-004-w7-independent.md) / [A-005](GOAL-007-w7-api-web-security-audit/03-audit/A-005-w7-independent-code-review.md)。
**W6（2026-08-15）**：承接本会话对 api/web 的代码审视——scheduler 未到期任务 5 年分钟空扫描改 O(1) Matches 快速路径（每日一次诊断保留）；回收站还原孤儿字典项 500 退化改 409 DICT_KEY_NOT_FOUND（快照保留可重试）；branding data:image 内联评估后**不采纳**（API normalizeLogoURL 与 errorcatalog 均拒绝，web 测试锁定，保持一致收紧）。`go test ./...` 全绿，self 审计 A-001 pass，开放 required = 0；**2026-08-17 补记用户授权关门（D-002）+ close-out self 审计 A-002 pass**，`status: done` 维持。
**W5（2026-08-14 扫描）**：全量审计 **0 中高危**（L-001～L-006 低危就地修补，见 Root [E-002](GOAL-001-production-hardening/02-execution/E-002-w5-scan-zero-midhigh.md)）；按程序约定未开子目标。**go 判定：无影响、不暂挂**（安全加固与已冻结 fail-closed 语义一致；未改 Profile 默认集 / 模块矩阵 / Manifest 装配 / 协议 pin）。
Root **保持 active**。W1–W4 为已关门波次档案；W4 承接 2026-08-11 新一批 api/web 全量审计修复（限流驱逐、上传权限门+配额、改密吊销 access token、前端异常捕获、URL 校验、启动加固、文案）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-production-hardening | 生产加固（共享基架持续安全与健壮性） | null | active | —（程序容器，不用 n/n→done） | 2026-08-10 |
| GOAL-002-audit-findings-remediation | 审查发现修复（W1） | GOAL-001-production-hardening | done | 16/16 | 2026-08-10 |
| GOAL-003-upload-ownership-hardening | 上传所有权与下载鉴权加固（W2） | GOAL-001-production-hardening | done | 4/4 | 2026-08-10 |
| GOAL-004-w3-security-audit-remediation | W3 安全审计发现修复（api/web） | GOAL-001-production-hardening | done | 8/8 | 2026-08-11 |
| GOAL-005-w4-security-audit-remediation | W4 安全审计发现修复（api/web） | GOAL-001-production-hardening | done | 8/8 | 2026-08-11 |
| GOAL-006-w6-scan-findings-remediation | W6 扫描审计发现修复（api/web） | GOAL-001-production-hardening | done | 4/4 | 2026-08-15 |
| GOAL-007-w7-api-web-security-audit | W7 api/web 独立安全审计（落盘） | GOAL-001-production-hardening | done | 4/4 | 2026-08-19 |
| GOAL-008-w8-api-web-security-audit | W8 api/web 独立安全审计（审计报告落盘） | GOAL-001-production-hardening | done | 4/4 | 2026-08-20 |
| GOAL-009-w9-api-web-security-audit | W9 api/web 独立安全审计（审计报告落盘） | GOAL-001-production-hardening | done | 4/4 | 2026-08-21 |
| GOAL-010-w10-api-web-security-audit | W10 api/web 独立安全审计（审计报告落盘） | GOAL-001-production-hardening | done | 4/4 | 2026-08-21 |
| GOAL-011-w11-api-web-security-audit | W11 api/web 独立安全审计（审计报告落盘） | GOAL-001-production-hardening | done | 4/4 | 2026-08-22 |
| GOAL-012-w12-multi-instance-rate-limiting | W12 多实例限流拓扑评估与加固（承接 workspace-019 E-009 §F-002） | GOAL-001-production-hardening | done | 4/4 | 2026-08-26 |
| GOAL-013-w13-api-web-security-audit | W13 api/web 全量安全审查发现修复（P1×1+P2×3 必修 + P3 全量） | GOAL-001-production-hardening | done | 6/6 | 2026-08-26 |
| GOAL-014-w13-account-lockout-redesign | W13-F007 账号锁定模型重设计（fixed · 承载自 GOAL-013） | GOAL-013-w13-api-web-security-audit | done | 6/6 | 2026-08-26 |
| GOAL-015-w14-schema-auth-wiring-lock | W14 页面 Schema 鉴权装配修复与生产装配回归锁 | GOAL-001-production-hardening | done | 4/4 | 2026-08-26 |
| — | W5 scan（0 中高危；低危就地修补，未开子目标） | GOAL-001-production-hardening | — | — | 2026-08-14 |

## 维护说明

- Root 是长期能力容器；`status: done` 仅在程序废弃或 `primary_plan` 迁移且用户确认时使用。
- 波次 progress 只写在子目标；不得用波次完成数推导 Root done。
- 层级唯一来源是目标 `00-meta.md` 的 `parent`。
