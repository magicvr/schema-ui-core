---
title: D-002 · W15 方案冻结与放行（S2）
status: active
created: 2026-08-30
updated: 2026-08-30
parent: null
version: 0.1.0
---

# D-002 · W15 方案冻结与放行（S2）

日期：2026-08-30 · scope：`[workspace-009-production-hardening] GOAL-016-w15-api-web-audit-remediation` · 决策性质：方案冻结（S2 检查点）与 S3～S6 放行

## §1 用户放行（P-004 书面确认）

用户对以下四项裁决全部选择推荐项（2026-08-30 会话内书面确认）：

1. **F-001/F-002**：默认监听回环 `127.0.0.1:25080`（显式 `-addr`/config/`HTTP_ADDR` 才能对外暴露）；自定义 YAML 未声明 `APP_ENV` 启动 fail-closed；非 development 的 JWT secret 复用主仓强度规则（≥32 字符 + 字母数字混合，单一来源）；显式 development 保持文档化 dev 链（`admin/admin` + 固定 dev 密钥）不变。
2. **F-003**：仅非 development 在 bootstrap 前强制既有冻结密码策略（8–72 字节、非空），不满足启动 fail-closed；development 保留 `admin` 回退。
3. **F-007**：**fixed**——LocalStore 目录 `0700`、文件 `0600`。
4. **放行**：S2 冻结后顺序推进 S3（API 修正+回归）→ S4（Web 修正+回归）→ S5（F-007+全量验证）→ S6（self → grok build independent → 用户关门授权）。

## §2 信息就绪结论（P-005，I-001～I-005 全部关闭）

| ID | 结论 | 证据 / 留痕 |
|----|------|-------------|
| I-001 | independent provider = **grok build · grok-4.6 · reasoning high · `/audit`**；Root 旧记录中 grok-4.5 系历史措辞，以用户目标指令 + workspace-009 页眉（2026-08-30）与 `docs/architecture/independent-audit-execution.md`（v1.0.0）为准 | 用户目标指令原文；workspace.md §愿景对齐；architecture 项目级决策 |
| I-002 | `schema-ui serve` 是**受支持的下游运行入口**（VP-024 R1 交付：HTTP 壳 + config 装载 + assembly 服务器面 + RT-D02 接线；双方言 E2E 实证；`create` 骨架可直接 serve 启动）。默认配置即 development 链（`:25080` 全网监听 + `admin/admin` + 固定 dev 密钥），不作为生产默认；生产/局域网暴露需显式配置，且必须满足非 development 门禁 | `apps/api/cmd/schema-ui/main.go:175-203`；`server/config.default.yaml`；`cmd/schema-ui/templates/{main.go,config.yaml}.tmpl`；workspace-024 GOAL-002 E-003（E2E-L1/L3） |
| I-003 | 初始管理员密码必须满足现有 authsession 冻结策略下限（8–72 字节、非空、含已配置类别）；fresh DB 的 policy 行（migration 0057）即播种默认 minLength=8，故 bootstrap 静态检查与 frozen default 等价；已有库无影响（bootstrap 仅在零用户 fresh 库执行） | `modules/authsession/password_policy.go`（floor 8 / ceiling 72）；`migration/migration.go:251-257`（seed min_length 8）；`systemdata/bootstrap.go`（needs-bootstrap 语义） |
| I-004 | canonical fixture 根 = **`apps/api/modules`**（`apps/api/internal/modules` 不存在，本轮实测 13 suite / 76 测试失败与此一致）；被引用 fixture 全部实际存在（含 `dev/examples/schema`、`wallet/schema`） | vitest 基线 `13 failed | 75 passed (88)` / `76 failed | 1081 passed (1157)`；目录盘点 |
| I-005 | 文档化部署拓扑 = Docker Compose **非 root `app` 用户 + 命名卷**（I-008-001 契约），对象经 HTTP 面校验读取、无静态直出；多 OS 账号共机非文档化拓扑。加固到 `0700`/`0600` 廉价无损（同一进程创建/读取）→ 用户裁决 **fixed** | `internal/objectstore/local.go:113-145`；workspace-002 I-008-001 契约（Dockerfile 非 root、`db-data` 卷）；handler owner 校验先例 |

## §3 冻结方案（A-001 分母逐条）

### F-001 · 公共 serve 默认暴露面收紧

- 默认监听地址 `:25080` → `127.0.0.1:25080`：`server/config.go` 代码默认、`server/config.default.yaml`、`cmd/schema-ui/templates/config.yaml.tmpl`（create 骨架）三处同步；对外暴露必须显式 `-addr`/`http.addr`/`HTTP_ADDR`。
- 空 `APP_ENV` fail-closed：自定义 YAML 未声明 `app.env` 时启动报错（内嵌默认仍显式 `development`，dev 流不受影响）。
- 显式 `development` 保持现状（`admin/admin` + 固定 dev 密钥），作为文档化 dev 链。

### F-002 · 非 development JWT secret 强度门禁

- `internal/config` 导出 `ValidateJWTSecretStrength`（≥32 字符 + 字母 + 数字，错误不携带密钥值），`ValidateProd` 内部复用同一实现（单一来源）。
- `server.Config.validate()` 非 development 分支调用同一校验；短/纯字母/纯数字 secret 启动 fail-closed；补负例测试。

### F-003 · 生产 bootstrap 密码策略

- `modules/authsession/password_policy.go` 导出 `ValidateSeedPassword(plain)`：复用策略边界常量（8–72 字节、非空）。
- `cmd/server/main.go` `resolveSeedHash`：非 development 先过策略再 bcrypt，不满足 fail-closed。
- `server/serve.go` `bootstrapAdmin`：非 development 先过策略再 bcrypt；config 层已保证非空。
- development 回退 `admin` 不变。补启动负例测试（短密码 fail-closed）。

### F-004 · MFA step-up 一次性语义

- `Service.requireActiveSecondFactor`：TOTP 校验成功后以匹配 step 执行 `AdvanceLastUsedStep`（CAS，登录路径同语义）；CAS 失败（同窗重放/并发）→ `ErrMFAInvalid`；`maybeRewrap` 仅 CAS 赢者执行。
- 覆盖 disable / recovery rotate / VerifySecondFactor 共用路径；补「同码二次提交拒绝」与并发守卫测试。

### F-005 · 邀请 token URL 清理

- `invite-accept.tsx`：挂载时读取 token 后立即 `history.replaceState` 删除 query 中的 `token`（地址栏/历史不再携带一次性 bearer）；成功跳转路径不变。
- `no-referrer` 与短 TTL 维持既有行为（A-001 建议的 referrer 策略属 nginx 面，不在本波）。

### F-006 · Web 测试 fixture 根统一

- 13 个失败 suite 的 fixture 根统一切换 `apps/api/internal/modules` → `apps/api/modules`（含 `schema-keys.structural` 的 fragment 清单与 `row-action-bindings` 的深度路径）。
- 新增 fixture 根 guard 测试（`src/protocol/fixture-root.guard.test.ts`）：断言 canonical 根存在 + `src` 内无任何文件再引用已退役的 `internal/modules` 根（A-001「CI 路径存在性检查」落地为 vitest 锁）。

### F-007 · LocalStore Unix 权限（S5 范围，fixed）

- `Put`：目录 `0o700`、临时文件与最终文件 `0o600`、sidecar `0o600`；既有文件不强制改写（新写对象即收紧）；补权限断言测试。

## §4 未选路径

| 路径 | 不采用原因 |
|------|------------|
| 保留 `:25080` 全网默认、只靠文档提醒 | 无法满足 F-001 required；默认即安全、显式才放开 |
| 空 `APP_ENV` 继续按 development 启动 | 与主仓 `ValidateProd`「refusing to guess」语义不一致，静默 dev 启动正是 F-001 根因 |
| `internal/config` 规则复制到 server 包 | 双份密钥规则会漂移；导出共享校验函数保持单一来源 |
| dev 也强制 8+ 字节 bootstrap | 破坏文档化 `admin/admin` dev 链与既有 E2E/smoke 流；F-003 只针对生产 bootstrap |
| F-007 有界 residual | 用户已书面选择 fixed（廉价无损加固） |
| F-005 顺带改 nginx referrer 策略 | 属部署面配置、超出本波代码分母；A-001 建议项中仅 replaceState 属代码修正 |

## §5 放行语义

- 本决策放行 S3～S6 的执行与回归，不构成任何 finding 已 fixed 的宣称；闭合证据（实现 + 测试 + 独立审计）在 E 条目与 03-audit 响应节落盘。
- 审计模式：security/production 影响范围按工作区约定走 `cross`（S6 self → grok build independent），provider = grok-4.6 · high（I-001 已关闭）。
- checkpoint：冻结（本决策）、S3 实现切片、S4 实现切片、required 闭合（独立审计后）、关门前最终验证后各提交一次，只 `git add` owned paths。