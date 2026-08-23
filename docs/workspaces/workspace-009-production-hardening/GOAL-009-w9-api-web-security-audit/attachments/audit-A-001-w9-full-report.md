# 附件 · A-001 全文 · W9 api/web 独立安全审计报告（2026-08-21）

> auditor：ox-alpha（DSH 会话；3 并行子代理 + 主线深读 + P1/P2 逐条源码交叉复核 + go vet）。
> 本附件为 A-001 摘要条目的全文载体；结论以 A-001 为准。

**方法**：本次为独立审计，未加载任何 skill。三路并行子代理（handlers 广度 / store+modules 层 / web 前端）+ 主线深读全部安全关键路径（auth、session、permission、rate-limit、upload、service-credentials、MFA/TOTP、captcha、wallet 账务、config/composition/main、tokens、auth-client、reaction-expression、nginx、Dockerfile/compose），对子代理的 P1/P2 结论逐条亲自复核源码确认，另跑 `go vet`（0 告警）与 git 追踪核查。**未发现 P0；共确认 P1×2、P2×10、P3×20+。**

---

## P1（2 项）

**[P1-1] API·Postgres 方言下钱包唯一冲突检测失效** — `wallet/store/repository.go:845`
`isUniqueViolation` 只匹配 SQLite 错误文案（"UNIQUE constraint failed"），pgx 报 `SQLSTATE 23505 duplicate key value violates...` 永不匹配。后果：① `CreateAccount` 重复 owner 返回裸 500 而非 `ErrOwnerTaken`；② `GetOrCreateUserAccount`（:294-315）并发创建回退成死代码（且 PG 事务失败后即中止，同事务重读也不可行）；③ `Mutate` 的 `errIdempotencyRace`（:489-512）在 PG 上永不触发——**幂等键并发提交得到 500 而非幂等重放，钱包幂等保证在一级方言上失效**。余额本身仍受乐观版本 CAS 保护，无资金错账。修复：按 `pgconn.PgError.Code=="23505"` 做驱动无关判定，get-or-create 回退改新事务。

**[P1-2] Web·生产 nginx 未代理 host-bootstrap，compose 栈必然启动失败** — `apps/web/nginx.conf:33-62`
nginx 只代理了 `app-manifest.json`；`/.well-known/schema-ui/host-bootstrap.json` 落入 SPA fallback 返回 **200 text/html**（index.html）。而 `host/bootstrap.ts:172-188` 只把 404/410 视为"未提供"，非 JSON content-type 判为 `protocol` 失败 → `boot.ts:169-180` 终态 `BOOTSTRAP_DOCUMENT_FAILED`，retry:none。已核实 `Dockerfile:27` 将该 nginx.conf 打进生产镜像、compose web 使用它、API 侧确实注册了该路由（composition.go:532）。**修复**：为 host-bootstrap 增加与 manifest 相同的 proxy 块（或将 `/.well-known/` 排除出 SPA fallback，让 404 保持 404）。

## P2（10 项）

**API（7 项，均已亲验）**

1. `scheduledtasks/store/repository.go:86-97, 299-311` — `WHERE key LIKE ? OR name LIKE ?` 后追加 ` AND enabled = ?` 未加括号，AND 优先级更高 → `q + enabled=true` 会返回 **disabled** 任务；`ListAllRuns` 同型绕过 status 过滤。
2. `authsession/accounts.go:173-196` — `RecordLoginFailure` 读-改-写非原子（无 `+1`/`FOR UPDATE`），PG 并发失败登录丢失计数 → **锁定阈值推迟到达，助长暴力破解**（SQLite 因单连接串行而安全）。
3. `mfa/service.go:123-137` — TOTP 同步重放拒绝跨两个事务 check-then-act（GetState→Validate→SetLastUsedStep），两个并发 Verify（每口令登录各发一个 proof）可让**同一验证码被接受两次**。
4. `mfa/service.go:300-319` — 恢复码消费为"读全表→bcrypt→整表回写"，并发兑换两枚不同码时后写覆盖前写（已消费码复活）；同码并发则双成功。
5. `jobs/runner.go:278-281` + `scheduledtasks/scheduler.go:64-75,131` — job/task 处理器无 `recover`，**panic 直接击穿整个进程**。
6. `authsession/service_credentials.go:59` — 重名判定匹配 `"service_credentials.name"`（SQLite 索引名），PG 约束名为 `service_credentials_name_key`（migration.go:254）→ PG 上重名返回 500 而非 400。
7. `scheduledtasks/store/cron.go:99-107` — DOM/DOW 用 AND 组合，非 POSIX 的"两者受限取 OR"：`0 0 1 * 1` 只在"1 号恰逢周一"触发，标准表达式被静默少调度。

**Web（2 项，均已亲验）**

8. `renderer/permissions.ts:353-356` + `render.tsx:464-467,914` — 声明了 `permissionIntent` 但漏写 `key` 的动作：targetId 记为空串，而门禁查找用 `actionRef` 兜底 → 永不匹配 → **客户端权限门禁整体失效（fail-open）**；`schema-table.tsx:770,866` 对未注册 target 默认放行同病。
9. `renderer/permissions.ts:512-520` — cascade 键缺 permissions 源时被跳过并最终 `return true`（fail-open）；能捕获此问题的 L2 校验器 `validatePermissions`（:97）**仅测试调用，生产渲染路径从未执行**。

**API（1 项潜伏）**

10. `handler/resources.go:716-723` — `delete()` 的 self-scope 归属预检 `if row, gerr := Entity.Get(id); gerr == nil`：**Get 出错时静默跳过归属检查并继续 Delete**（对比 `update()` :680-683 是 fail-closed）。当前无生产资源接线 `Scoper`（仅测试），属潜伏的越权 fail-open 模式，建议立即对齐 update 的错误处理。

## P3（择要）

**API**：头像每用户配额 CountOwner→save TOCTOU 可并发击穿（account_avatar.go:62-71，对照 upload.go 已有 quotaMu 修复同类问题）；POST /api/scheduled-tasks/{id}/run 手动执行任务无操作日志、无行为者归属；MFA proof 签发无每用户频控（持有口令者可无限铸造 proof，每 proof 5 次 TOTP 猜测，mfa_proofs 表亦可被灌爆——建议加每用户 proof 数/频率上限与全局尝试计数）；/api/mfa/enroll 的 currentPassword 步进验证无频控（changePassword 已有而此处没有）；otpauthURL 转义只处理空格/冒号（?&#% 可重构 URI）；LIKE 未转义 %/_（9 处，通配注入+全表扫）；cron 按服务器本地时区求值，site_timezone 未生效；scheduler 无多实例抢占（双副本双跑）；jobs/repository.go:131-141,406 忽略 RowsAffected（取消静默丢失/终态钩子漏发）；recyclebin restore 与 MarkRestored 两事务可留下永久冲突快照；用户名并发重名裸 500；strings.Contains(err,"unique") 判重（4 处）驱动文案依赖；rebindPostgres 会改写字符串字面量中的 ?；store.WithTx 无 recover/rollback（单连接下可死锁）；CreateUser/linkUserRole 可自动建任意角色（当前仅测试调用，提权脚枪）。

**Web**：App.tsx:830 未拦 //host 协议相对路径（pushState 抛未捕获异常）；boot.ts:328 "support" 恢复动作点击未校验 URL（当前无生产者，潜伏 javascript: 汇点）；auth-client.ts:215 子串匹配判定密码修改通知；upload-orchestration.ts:123 负 maxSize 关闭大小上限；form-controls.ts:511 schema 正则直接 new RegExp（ReDoS）；claim.ts:274 依赖环死循环（无 visited 集）；API 控制的 https:// URL 直接进 img src（头像/预览追踪像素，CSP img-src https: 放行）；render.tsx:930 useCallback 漏 t 依赖（切换语言后确认框文案滞留旧语言）；LoginPage MFA resolver 卸载后永不 settle；CSP style-src unsafe-inline。

**环境/卫生（info）**：apps/api/configs/.env 含真实 PG 凭据（sa / 192.168.31.213）——已核实 gitignore 生效且从未入库，属本地卫生问题，建议改密码并保持不入库；PG 默认 sslmode: disable；refresh token 存 localStorage 为文档化已接受的取舍（access 仅内存+轮换+generation 守卫缓解）。

## 已排除的疑点（验证为干净，避免误报）

SQL 注入零发现（排序字段全白名单、IN 列表占位符计数匹配、值全参数化）；无路径穿越（上传/资产 ID 强制 ^[0-9a-f]{32}$）；ServiceCredentialByHash 不会返回 (nil,nil)；refresh 轮换有 guarded UPDATE 原子防双花；钱包为纯整数运算+三余额恒等式+乐观锁 RowsAffected 校验；迁移账本 checksum+连续性校验、每迁移单事务；上传面（MIME 嗅探+活动内容标记+附件强制下载+CSP sandbox）、图片解码炸弹防护（DecodeConfig 预检 2048px）、CSV 公式注入防护（formulaSafe）、导入角色委派边界、CORS 无凭据回显、dev-session 生产 fail-closed（ValidateProd）、JWT 密钥生产长度门禁——均确认有效。Web 端无 dangerouslySetInnerHTML/innerHTML/eval；表达式引擎为手写白名单解析器；return-intent 单次使用+TTL+scheme 与双斜杠路径拒绝；401 并发刷新去重正确。

**修复优先级建议**：P1-2（nginx，一行配置即可恢复生产可用）→ P1-1 + P2-1/2/3/4（PG 方言与并发正确性，统一引入 SQLSTATE 23505 判定与原子 UPDATE）→ P2-8/9/10（权限 fail-open 三处，统一改 fail-closed）→ 其余 P3 按模块批量清理。
