---
title: I-008-002 · 15 分钟 fork 复现与 smoke 验收协议
status: active
doc_type: information-contract
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-008-r5-engineering-fork
version: 0.1.1
related_info: I-008-002
related_decisions: D-004, D-005
frozen_revision: 5e27019482eb8d0695c402b784860233bbc90c39
---

# I-008-002 · 15 分钟 fork 复现与 smoke 验收协议

## 1. 冻结范围

本协议回答 `I-008-002` 的信息问题：15 分钟 fork 体验如何计时、什么构成独立复现记录、以及 `scripts/smoke.sh` 应以什么机器可判定标准验收。

- **依据**：[Root D-013](../../GOAL-001-production-admin-foundation/01-decision.md) 已采纳的终点、依赖下载排除和独立复现方向；[I-005 收集报告](../../GOAL-001-production-admin-foundation/attachments/I-005-engineering-fork-collection.md) 的候选方法；[A-001 R-002](../03-audit.md) 的最低收集清单。
- **冻结内容**：工具/平台基线、依赖缓存前提、计时起止、失败与重试规则、独立复现记录字段、smoke 检查项、脚本输入和退出码。
- **不在本协议内**：S3 的 QUICKSTART/README 实现、任何独立复现实测、S4 的 `scripts/smoke.sh` 实现或运行结果、CI run、R5/Root/VP 关门结论。本协议不是上述事项的证据。

## 2. 基线与可移植性

### 2.1 冻结时观察到的工程基线

| 项 | 值 | 用途 |
|----|----|------|
| 仓库 revision | `5e27019482eb8d0695c402b784860233bbc90c39` | 协议形成时的实现参照；后续复现记录必须写明实际 ref，不能默认继承此值。 |
| 工作树 | clean | 仅说明协议形成时没有未提交变更；不替代后续复现记录。 |
| 记录主机 | Windows 11 Pro `10.0.26200` / x64 | 当前本地收集环境。 |
| Git | `2.47.0.windows.2` | 复现记录的版本字段基准。 |
| Go / Node / npm | `go1.26.0 windows/amd64` / `v22.17.0` / `10.9.2` | 当前 API/Web 工具链基准。 |
| Docker / Compose | `29.6.2` / `v5.3.1` | Compose 路径的当前工具基准。 |
| curl | `8.21.0 (Windows)` | HTTP smoke 工具基准。 |

当前 Windows 主机的 WSL `bash` 启动失败，故本次**没有**运行任何未来的 `scripts/smoke.sh`。S4 脚本必须在可用的 POSIX `bash` 环境执行：Linux CI，或本机已验证可用的 Git Bash/WSL。该约束本身不构成 S3/S4 已实施或未通过的结论。

### 2.2 可接受的复现路径

复现记录必须选择并声明一条路径：

1. `local-dual-process`：按文档启动 API 与 Web；或
2. `compose`：按文档从仓库根启动 API 与 Web 容器。

两条路径都仍是 R5 交付范围。S3 的最小要求是至少一条路径有一份合格的独立复现记录；未实测的另一条路径不得被写为已经满足 15 分钟指标。

## 3. 15 分钟计时协议

### 3.1 计时前提与排除项

开始计时前必须完成并记录下列事实：

- 已 `git clone` 并 checkout 到待测 ref；工作树 clean，或在记录中列出并哈希未提交差异。
- 所需运行时已经安装；外部依赖下载、镜像 pull 与语言依赖缓存已经完成。可记录的准备命令包括 `go mod download`、`npm ci` 与 Compose 所需镜像/层获取。
- 不得把项目自身的编译、配置、服务启动、数据库迁移、种子初始化或登录步骤预先完成；它们属于被测体验。
- 已取得但不在日志中输出 `AUTH_JWT_SECRET` 与 `ADMIN_INITIAL_PASSWORD` 的值；记录只写来源类别（例如本地 `.env` 或 CI secret），不写 secret 本身。

依赖下载、网络拉取、工具安装和 clone 的耗时不计入 15 分钟；其余按文档完成环境配置、构建和启动的耗时计入。

### 3.2 起点、终点与判定

| 项 | 冻结定义 |
|----|----------|
| 起点 | 依赖就绪后，执行 QUICKSTART/README 中第一条环境配置或启动命令的瞬间；复制/编辑 `.env`、构建和启动均在计时内。 |
| 终点 1 | `GET ${API_BASE_URL}/healthz` 返回 HTTP 200，JSON 含 `status: "ok"`。 |
| 终点 2 | 经 `${WEB_BASE_URL}/api/auth/login` 使用种子 admin 登录成功，HTTP 200 且响应含非空 `accessToken`。 |
| 终点 3 | 携带该 token 请求 `${WEB_BASE_URL}/api/accounts/me` 返回 HTTP 200，含 user 与 features 投影。 |
| 终点 4 | 浏览器登录后打开 `${WEB_BASE_URL}/list-edit-lifecycle`；页面显示 `List + edit lifecycle` 且代表性列表已加载 `Acme Console`。 |
| 达标 | 四个终点均满足，单次计时 `<= 900` 秒。 |

`API_BASE_URL` 与 `WEB_BASE_URL` 必须写入记录；默认期望按复现路径（§2.2）区分，且不得以默认值覆盖实际端口：

| 路径 | 默认期望 |
|------|----------|
| `compose` | API `http://localhost:8080`、Web `http://localhost:8081`（compose.yaml `web` 映射 `8081:80`） |
| `local-dual-process` | API `http://localhost:8080`、Web `http://localhost:${WEB_PORT:-5173}`（Vite dev 默认 5173，`WEB_PORT` 可覆盖） |

D-013 的「后台首页可交互（列表加载）」在本目标操作化为 R4 代表页路由 **`list-edit-lifecycle`**（page title `List + edit lifecycle`，含列表加载 `Acme Console`），而非 manifest `homePageRef: overview`；QUICKSTART/README 必须覆盖该路由。

### 3.3 独立性、失败与重试

- **独立复现**：在与实现/文档编辑会话隔离的 shell、checkout 与可控数据库状态中执行；可以由同一操作者完成，但记录必须标注 `same-operator-clean-session`。同一 shell、复用已启动服务或已存在数据库的直接重跑不算独立复现。
- **失败留痕**：每次开始计时的尝试都必须记录开始/结束、结果、失败检查项和是否修改了环境或文档。不得只保留最快的成功样本。
- **重试**：只有 readiness 等待可在同一尝试中每秒重试，最多 30 次；登录、`/me`、页面/数据断言失败应立即使该尝试失败。环境修复后的新尝试必须另起记录。
- **不得掩盖**：不能以预热后的服务、预先生成的 token、人工修改响应或省略失败尝试来满足时限。

## 4. 独立复现记录的最小字段

每份 S3 记录必须至少包含下表；可存于该目标 `attachments/`，并从 `02-execution.md` 链接。

| 字段 | 要求 |
|------|------|
| protocol | `I-008-002` 版本与本协议路径 |
| attempt | 唯一编号、日期、时区、操作者/执行器标识和独立性声明 |
| source | 仓库 URL、commit/ref、工作树状态及任何 diff 的摘要/哈希 |
| path | `local-dual-process` 或 `compose`；API/Web 实际 base URL |
| platform | OS/架构、Git、Go、Node、npm、Docker/Compose（若使用）与 bash/curl 版本 |
| cache precondition | 已完成的依赖/镜像准备命令、是否命中缓存，以及哪些耗时被排除 |
| timing | 起点/终点 UTC 时间、单调计时秒数和四个终点各自时间戳 |
| checks | health、login、`/me`、浏览器页面/列表与 smoke 输出的逐项结果 |
| secrets | 仅写配置来源与是否脱敏；禁止记录 token、password 或 secret 值 |
| result | pass/fail、失败原因、重试编号、日志/截图/命令输出路径 |

## 5. `scripts/smoke.sh` 验收契约

### 5.1 运行边界

S4 新建的脚本必须位于仓库根 `scripts/smoke.sh`，使用 `bash` 执行，默认对正在运行的实例进行非破坏性检查。脚本至少要求 `bash`、`curl` 与 `node`；缺任一工具、缺少密码或 URL 配置无效时必须失败而不是猜测默认 secret。

脚本不得输出 access token、password、JWT secret 或 `.env` 内容。若要验证空数据库的种子重复性，只能在明确标记为 disposable 的环境中运行；实现必须拒绝对普通开发数据库执行 reset。**S4 验收必须包含至少一次 disposable/隔离运行且 `SM-006=PASS`**；CI 必须使用隔离的 Compose project/volume 或等价临时 DB，并默认以 disposable 模式运行（或等价 job 覆盖该次证据）。非 disposable 默认路径的 exit 0 不得单独作为 S4「种子可重复」关闭证据。

### 5.2 机器可判定检查项

| ID | 检查 | 通过条件 |
|----|------|----------|
| SM-001 | 参数与工具 | `bash`/`curl`/`node` 可用，必需 URL 与脱敏密码输入有效。 |
| SM-002 | API readiness | `${API_BASE_URL}/healthz` 在 30 秒内返回 HTTP 200 且 `status=ok`。 |
| SM-003 | 代理登录 | `${WEB_BASE_URL}/api/auth/login` 返回 HTTP 200 和非空 `accessToken`。 |
| SM-004 | 当前身份 | Bearer 调用 `${WEB_BASE_URL}/api/accounts/me` 返回 HTTP 200，含 user 与 features。 |
| SM-005 | 代表页路由 | `${WEB_BASE_URL}/list-edit-lifecycle` 返回 HTTP 200 且响应体含 `id="root"`（或等价稳定 SPA 挂载标记）；S3 独立复现记录另行完成真实浏览器可交互/列表加载验证。 |
| SM-006 | 种子重复性（disposable 模式 · S4 必检） | 从空 DB 启动后，认证的 `GET ${WEB_BASE_URL}/api/records?pageSize=100` 返回 `SMOKE_EXPECTED_SEED_TOTAL` 条种子记录并含 `${SMOKE_RECORD_ID}`（默认 `rec-1`）/ `Acme Console`；重启 API 后再次断言数量和同一记录不变，不产生重复种子。**S4 验收必须至少一次以 disposable/隔离运行且 `SM-006=PASS`；无该次证据不得把 S4「种子可重复」判为满足。** |

### 5.3 输入与退出码

最小输入为 `API_BASE_URL`、`WEB_BASE_URL`、`SMOKE_USERNAME`、`SMOKE_PASSWORD`、`SMOKE_RECORD_ID`（默认 `rec-1`）和 `SMOKE_EXPECTED_SEED_TOTAL`（默认 `8`）。实现可增加明确文档化的 `--disposable` 或等价安全开关；无该开关时不得执行种子 reset。

| 退出码 | 含义 |
|--------|------|
| `0` | SM-001～SM-005 通过，且 SM-006（disposable）通过——S4 完整绿。 |
| `2` | 参数、工具或安全前提不满足；请求了不安全的 destructive 模式。 |
| `3` | readiness 在 30 秒内未达到。 |
| `4` | 登录或身份检查失败。 |
| `5` | 代表页路由或数据检查失败。 |
| `6` | 种子重复性（SM-006）断言失败。 |
| `70` | 未分类的脚本内部错误；实现必须输出失败检查项和脱敏诊断。 |

脚本输出应以稳定的 `SM-00N=PASS|FAIL` 形式逐项报告，方便 CI 归档；非零退出不得被 CI 忽略。

## 6. 门禁效果

本协议使 `I-008-002` 的“协议与判据是什么”信息问题可核对，因此解除 S3/S4 的**方案冻结/实施前信息门禁**。它不解除 S3/S4 的实施或验收责任：仍须交付文档、实际脚本、至少一份独立复现记录和相应运行/CI 证据；S1/S2 的既有事实、S5 审计和 Root R5 状态均不因本协议改变。

## 7. 修订记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v0.1.0 | 2026-08-03 | 初始冻结（D-004）。 |
| v0.1.1 | 2026-08-03 | 响应 A-008（independent · design-plan · conditional）：**F-004 → fixed**——S4 验收强制 ≥1 次 disposable/隔离运行且 `SM-006=PASS`，非 disposable exit 0 不得单独作为「种子可重复」关闭证据（§5.1/§5.2/§5.3）；吸收 **R-008-001**（默认 URL 按路径区分）、**R-008-002**（终点 4 操作化对齐 `list-edit-lifecycle`）、**R-008-003**（SM-005 判据钉死 `id="root"`）、**R-008-004**（不安全 destructive → 退出码 2，SM-006 失败单独 6）。D-005。 |
