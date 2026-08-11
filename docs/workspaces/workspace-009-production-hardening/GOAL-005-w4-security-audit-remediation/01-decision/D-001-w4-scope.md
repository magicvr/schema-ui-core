---
id: D-001
goal: GOAL-005-w4-security-audit-remediation
title: W4 修复范围与技术取舍
date: 2026-08-11
status: accepted
---

# D-001 · W4 修复范围与技术取舍

## 背景

2026-08-11 新一轮 api/web 全量审计（四路并行 agent + 主路径第一手核对）确认多项 bug/安全发现。workspace-009 Root 保持 active；本决策冻结 W4 波次范围。W3（GOAL-004）8 项已完成，F-001 等已闭合，本波不重复。

## 审计发现汇总（确认 bug / 确认行为）

### api 侧

1. **限流器容量驱逐死代码**（`rate_limit.go:43-77`）：`allow()` 无条件 `attempts[key] = kept` 预建 key；`record()` 仅当 key 不存在才触发驱逐 → `exists` 恒真 → 驱逐分支永不执行。攻击者喷洒不同用户名可无界增长 map → OOM。**medium**
2. **上传无权限门/配额**（`upload.go:116-201`，挂载于 `composition.go:176`）：仅要求已登录，不检查任何 permission key；无每用户/全局配额、无清理。低权限账户可无限上传 8MiB 文件灌满磁盘。**medium**
3. **改密不吊销已签发 access token**（`users_repository.go:187-193`）：改密只吊销 refresh token；已签发 access token（15min TTL）仍有效。**medium**
4. **`/readyz` 饿死风险**（`store.go:49,68-81`）：单 SQLite 连接 + 全部 `withTx` 用 `context.Background()`；慢查询持连接时 `/readyz` 1s 超时 503 → 容器健康门误判。**medium**
5. **开发模式密钥暴露**（`main.go:82` `config.go:44`）：dev 回退 JWT 公开常量 + 默认 `admin/admin` + `HTTP_ADDR` 默认 `:25080`（全接口）。dev 实例可达即 admin 沦陷。**medium**
6. **其他 low**：`page` 参数无上限 → 巨大 OFFSET（`resources.go:314`）；非 admin users-writer 可改 admin `name`（`users.go:302-322`，W3 已标注 recommended）；`/healthz`/`/readyz` 泄露 version/commit（`health.go:50`）；谎报 Content-Length 超限上传返回 400/500 而非 413（`upload.go:137-154`）；`WithTx` 无 defer 回滚（`store.go:68-81`）。

### web 侧

7. **recordSource prefill 异常未捕获 → 白屏**（`render.tsx:930` + `request-construction.ts:80-86`）：`useEffect` 内 `constructRequest` 无 try/catch；`serializeQueryValue` 对非标量 throw → 整树卸载（无 ErrorBoundary）。**medium**
8. **动作静默失败**（`render.tsx:320,419` + 调用点 600/634/697）：`runRequest`/`runBatchRequest` 内 `constructRequest` 无 try/catch；调用方 `void`/无 `.catch` → unhandled rejection，点击无反馈。**medium**
9. **pageTrigger/outcomeNavigate URL 校验弱**（`request-construction.ts:532-551,576-584`）：`buildPageTriggerRequest` 只查 `startsWith("/") && !startsWith("//")`，不拒反斜杠；`buildOutcomeNavigate` 完全不校验，绝对 URL 透传。当前被 action.schema 结构性校验屏蔽（应用内不可达），属防御纵深缺口。**low**
10. **启动路径白屏**（`main.tsx:72` `theme.ts:71-76`）：`initTheme` 顶层调用无 try/catch；存储禁用时模块求值失败 → React 不挂载。**low-medium**
11. **登录失败文案字面 `{status}`**（`LoginPage.tsx:90` + 消息 `{status}` 占位符）：`t()` 未传 params → 渲染字面占位符。**low**
12. **密码字段缺 autoComplete**（`form-controls.tsx:81-88`）：改密表单可能被浏览器误填。**low**
13. **行操作/批量无 in-flight 去重 + invocationId 未接通**（`render.tsx:551,610` + `request-construction.ts:238-256`）：双击并发重复提交；声明 `retryPolicy: idempotent` 的 action 因缺 invocationId 一律 fail-closed。**medium**（当前无 schema 声明 idempotent → 潜在）

## 范围内（按优先级）

### P0

1. **限流器容量驱逐修复**：`allow()` 不预建新 key（仅读）；驱逐逻辑统一由 `record()` 在登记时执行；或 `allow()` 建 key 时同时登记 order。目标：map 有界。
2. **上传权限门 + 配额**：`POST /api/upload` 增加权限检查（默认 `files.write`，admin 默认持有；system 角色授予说明）。配额作为第二层：每用户文件数与总字节上限（best-effort 存储层统计）。
3. **改密吊销 access token**：递增用户 `token_version`，写入 access token claims；中间件校验版本比对 → 改密即即时吊销。新增迁移/列。

### P1

4. **前端异常统一捕获**：`constructRequest` 在 UI 调用点（recordSource prefill、runRequest、runBatchRequest）包 try/catch，失败走既有 error 反馈通道，绝不白屏/静默。
5. **pageTrigger/outcomeNavigate URL 校验统一**：统一用 `isRelativeProtocolUrl`（拒反斜杠/空白/`//`）；`buildOutcomeNavigate` 拒绝绝对 URL。

### P2

6. **web 启动路径加固**：`initTheme`/`setTheme` 包 try/catch（与 tokens.ts 一致）。
7. **登录失败文案 + 密码 autoComplete**：`t(loginErrorKey(code), { status })` 或在 catch 回退 `err.message`；密码字段加 `autoComplete="new-password"`。

> 原第 8 条「低危修复」（`page` 上限 / `WithTx` defer / 上传 413 / health 指纹）经 A-002 N-004 指出与成功标准不一致——本波 8 项成功标准不含该条，实施也未交付。**随 A-002 响应调整为「明确不做 · recommended · 下波」**（见下方），消除范围漂移。

## 明确不做（本波）

- **refresh token 存 localStorage**：沿用 D-002 文档化 XSS 权衡（不重复）。
- **schema/manifest 匿名可读**：沿用 GOAL-002 D3 accepted-residual（不重复）。
- **Compose 默认 TLS / 生产 CDN CSP / HSTS**：部署层职责，非本波。
- **`/readyz` 独立连接 / MaxOpenConns 调整 / 请求 ctx 穿透 withTx**：可用性加固，涉及存储层结构性改动，标记为 recommended（下波）。
- **dev 密钥强制随机 / HTTP_ADDR 默认 loopback**：配置策略，涉及开发体验取舍，标记为 recommended。
- **`page` 参数上限 / `WithTx` defer 回滚 / 上传谎报 Content-Length 413 / health 探针指纹**：低危加固，标记为 recommended（下波；A-002 N-004 一致性响应）。
- **操作日志/refresh token 清理、错误目录丢弃 message**：观察项，不做。

## 审计模式

security 高影响实施 → 完成后 **self** 审计（A-001）；独立审由编排器按 VP-009 provider 另开（不阻断本波落地）。审计模式沿用 W3：self + 可选 independent。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 限流 `allow()` 直接拒建 key（只读） | `allow()` 返回后调用方必然 `record()` 失败，若 `allow` 不建 key 则 record 每次触发驱逐，语义混乱；统一由 `record` 登记建 key + 驱逐 |
| 上传仅 admin-only（无委派） | 破坏现有 viewer 可上传的示例语义；用 `files.write` 权限门 + 默认 admin 持有更符合 RBAC 委派模型 |
| access token 完全无状态（不吊销） | 15min 窗口可被窃取 token 利用；token_version 成本低、即时生效 |
| 前端只改一个调用点 | recordSource/runRequest/runBatch 是三处独立路径，需各自捕获才能全覆盖 |
