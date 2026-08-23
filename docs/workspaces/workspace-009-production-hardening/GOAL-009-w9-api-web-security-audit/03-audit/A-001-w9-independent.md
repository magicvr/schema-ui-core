---
id: A-001-w9-independent
goal: GOAL-009-w9-api-web-security-audit
status: final
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-009-w9-api-web-security-audit
version: 0.2.0
---

# A-001 · W9 api/web 独立安全审计（2026-08-21）

> **消费勘误（编排器 · 2026-08-21 · 非原文）**：S2/S3 冻结输入以 [D-002](../01-decision/D-002-w9-finding-inventory.md) 为准。F-003 作废（空洞，不复用）；全文 P2-7 = **F-025**；required **12** = F-001/F-002 + F-004～F-012 + F-025（无 F-003）。索引曾写「22」为抄写错误。下列独立意见原文不改写。

- **source**：independent
- **auditor**：ox-alpha（DSH 会话模型；3 并行子代理广度审计 + 主线深读安全关键路径 + 全部 P1/P2 逐条源码交叉复核 + `go vet`）。**非工作区默认 grok provider，偏差见 D-001/I-003。**
- **类型 / scope**：ad-hoc（用户直接指令的独立审计）· `apps/api`（Go，~20k 行非测试代码）与 `apps/web`（React/TS）当前实现的 bug 与安全漏洞
- **verdict**：fail

## 范围与区间

2026-08-21 工作树快照（未含未提交假设）。方法与完整证据见全文附件：[attachments/audit-A-001-w9-full-report.md](../attachments/audit-A-001-w9-full-report.md)。

## 结论摘要

**P0=0 · P1=2 · P2=10 · P3=20+**。整体评价：代码库经多轮审计波次（W1–W8）后基础面扎实——SQL 注入零发现、路径穿越防护完备、认证/会话/令牌轮换严谨、钱包整数账务+乐观锁正确、上传与图片处理纵深防御充分；本轮问题集中在 **Postgres 方言的错误映射与并发正确性**（消费 ID：F-001/F-004/F-011/F-020 族；原文曾写未定义的 F-003）、**两处部署/门禁 fail-open**（F-002/F-008/F-009/F-010）与 **MFA 并发语义**（F-005/F-006）。

## Findings（required）

### F-001 · 钱包唯一冲突检测仅匹配 SQLite 文案，PG 幂等/去重契约失效
- 严重度：high ｜ 建议：required ｜ 状态：open
- `wallet/store/repository.go:845-847` `isUniqueViolation` 只匹配 modernc/sqlite 文案；pgx 报 SQLSTATE 23505 永不匹配。后果：CreateAccount 重复 owner 裸 500；GetOrCreateUserAccount（:294-315）并发回退死代码（PG 事务失败后中止，同事务重读亦不可行）；Mutate（:489-512）`errIdempotencyRace` 永不触发 → **幂等键并发提交得到 500 而非幂等重放**。余额仍受乐观版本 CAS 保护，无资金错账。

### F-002 · 生产 nginx 未代理 host-bootstrap，compose 生产栈必然启动失败
- 严重度：high ｜ 建议：required ｜ 状态：open
- `apps/web/nginx.conf:33-62` 仅代理 app-manifest；`/.well-known/schema-ui/host-bootstrap.json` 落入 SPA fallback 返回 200 text/html。`host/bootstrap.ts:172-188` 仅 404/410 视为未提供、非 JSON 判 protocol 失败 → boot 终态 HOST_PROTOCOL_REJECTED（retry:none）。已核实 Dockerfile:27 将该 conf 打入生产镜像、compose web 使用之、API 侧确有该路由（composition.go:532）。

### F-012 · scheduledtasks 列表 WHERE OR/AND 优先级绕过过滤器
- 严重度：med ｜ 建议：required ｜ 状态：open
- `scheduledtasks/store/repository.go:86-97, 299-311`：`key LIKE ? OR name LIKE ?` 后追加 ` AND enabled = ?` 未加括号 → `q+enabled=true` 返回 disabled 任务；ListAllRuns 同型绕过 status 过滤。对照 datadictionary repository.go:251 正确加括号。

### F-004 · RecordLoginFailure 读改写竞态削弱锁定阈值（PG）
- 严重度：med ｜ 建议：required ｜ 状态：open
- `authsession/accounts.go:173-196`：SELECT 计数后 UPDATE 计算值，无原子自增/FOR UPDATE；PG 并发失败登录丢失计数 → 锁定推迟到达，助长暴力破解（SQLite 单连接串行安全）。

### F-005 · MFA TOTP 同步重放拒绝跨事务 check-then-act
- 严重度：med ｜ 建议：required ｜ 状态：open
- `mfa/service.go:123-137`（+`store/repository.go:160-171`）：GetState(tx1)→Validate(candidate>LastUsedStep)→SetLastUsedStep(tx2)；两个并发 Verify（每口令登录各发 proof）可让同一验证码被接受两次，SQLite 上竞态同样存在。

### F-006 · 恢复码消费丢失更新/双花竞态
- 严重度：med ｜ 建议：required ｜ 状态：open
- `mfa/service.go:300-319`：读全表→bcrypt→整表回写跨事务；并发兑换两枚不同码后写覆盖前写（已消费码复活）；同码并发双成功。

### F-008 · Web 权限门禁在动作缺 key 时整体失效（fail-open）
- 严重度：med ｜ 建议：required ｜ 状态：open
- `renderer/permissions.ts:353-356` targetId 取 `action.key`（缺省为空串），`render.tsx:914` 门禁查找用 actionRef 兜底 → 永不匹配 → `render.tsx:464-467` 完全跳过门禁；`schema-table.tsx:770,866` 对未注册 target 默认放行同病。服务端鉴权仍为硬边界，但协议承诺的 UI 门禁被 schema 笔误击穿。

### F-009 · cascade 缺权限源时 fail-open，L2 校验器生产路径未接线
- 严重度：med ｜ 建议：required ｜ 状态：open
- `renderer/permissions.ts:512-520`：cascade 键缺 permissions 源被跳过并最终 return true；能捕获该结构的 `validatePermissions`（:97）仅测试调用（grep 证实唯一调用方为 permissions-inheritance.test.ts:58）。

### F-010 · 资源工厂 delete() 归属预检 fail-open（潜伏越权模式）
- 严重度：med ｜ 建议：required ｜ 状态：open
- `handler/resources.go:716-723`：self-scope 预检 `if row, gerr := Entity.Get(id); gerr == nil` —— Get 出错即跳过归属检查并继续 Delete（对照 update() :680-683 fail-closed）。当前无生产资源接线 Scoper（仅测试），属潜伏缺陷，建议立即对齐。

### F-011 · 服务凭据重名判定方言 bug（PG 返回 500）
- 严重度：med ｜ 建议：required ｜ 状态：open
- `authsession/service_credentials.go:59` 匹配 `"service_credentials.name"`（SQLite 索引名）；PG 约束名为 `service_credentials_name_key`（migration.go:254 CITEXT UNIQUE）→ PG 重名返回 500 而非 400 ErrCredentialNameTaken。

### F-007 · job/task 处理器无 panic recover，panic 击穿进程
- 严重度：med ｜ 建议：required ｜ 状态：open
- `jobs/runner.go:278-281` 处理器 goroutine 与 `scheduledtasks/scheduler.go:64-75,131` tick/loop 均无 recover；单个 panic 处理器令整个服务下线。

## Findings（recommended / informational，摘录）

| F-ID | 严重度 | 建议 | 摘要（详见附件） |
|------|--------|------|------------------|
| F-013 | low | recommended | 头像每用户配额 CountOwner→save TOCTOU 可并发击穿（account_avatar.go:62-71；upload.go 已有 quotaMu 先例） |
| F-014 | low | recommended | POST /api/scheduled-tasks/{id}/run 手动执行无操作日志、无行为者归属 |
| F-015 | low | recommended | MFA proof 签发无每用户频控（持口令者可无限铸造 proof，每 proof 5 次 TOTP 猜测）+ /api/mfa/enroll 步进验证无频控（changePassword 已有先例） |
| F-016 | low | recommended | otpauthURL 转义仅空格/冒号，?&#% 可重构 URI（mfa/totp.go:77-84） |
| F-017 | low | recommended | LIKE 过滤未转义 %/_（9 处；通配注入+全表扫，非注入） |
| F-018 | low | recommended | cron 按服务器本地时区求值（site_timezone 未生效）+ 无匹配表达式 5 年分钟扫描性能 |
| F-019 | low | recommended | scheduler 无多实例抢占（双副本双跑；对照 jobs Runner guarded Claim） |
| F-020 | low | recommended | PG 并发错误映射族：RevokeServiceCredential 忽略 RowsAffected、RequestCancel/transitionJobs 忽略 RowsAffected、用户名并发重名裸 500、strings.Contains(err,"unique") 判重（4 处） |
| F-021 | low | recommended | recyclebin restore 与 MarkRestored 两事务，崩溃留永久冲突快照 |
| F-022 | low | recommended | rebindPostgres 改写字面量 '?'、store.WithTx 无 recover/rollback（单连接死锁）、CreateUser/linkUserRole 自动建任意角色（当前仅测试调用） |
| F-023 | low | recommended | Web 批量：App.tsx 未拦 //host、boot support URL 未校验 scheme、auth-client 子串匹配、maxSize 负值关闭上限、schema 正则 ReDoS、claim 依赖环死循环、API https URL 进 img src（追踪像素）、invokeAction 漏 t 依赖、LoginPage MFA resolver 泄漏、CSP style-src unsafe-inline |
| F-024 | info | — | 环境/卫生：configs/.env 本地真实凭据（已核实 gitignore 生效未入库，建议改密）；PG 默认 sslmode disable；refresh token localStorage 为文档化已接受取舍 |

## 必改项汇总（required）

原文自称 F-001～F-012 共 12 条。**消费以 D-002 为准**：F-003 作废；required 为 F-001、F-002 与 F-004～F-012、F-025（全文 P2-7）共 12 条，全部 open。recommended：F-013～F-023。informational：F-024。

## 结论 + 建议下一步

- verdict **fail**：2 条 high + 10 条 med required 开放。
- 建议下一步（编排器/用户裁决，对齐 W7/W8 先例）：① 用户裁决 required 修复范围与对 VP-008 go 消费有效性宣称的影响（I-002，D-002 模式）；② 按裁决实施（PG SQLSTATE 23505 判定统一化 + 原子 UPDATE 族 + nginx 一行代理 + 权限门禁 fail-closed 化）；③ self + independent 复核后关门。
- 本条目为独立意见，不改动目标 status/progress（per P-003）。
