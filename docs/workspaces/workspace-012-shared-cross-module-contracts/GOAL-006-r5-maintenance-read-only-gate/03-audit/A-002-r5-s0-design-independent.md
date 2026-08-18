---
id: A-002-r5-s0-design-independent
goal: GOAL-006-r5-maintenance-read-only-gate
doc: audit-entry
record_id: A-002
source: independent
scope: R5 S0 设计：运行态四模式、统一写门禁路由匹配/认证优先级、错误语义、bootstrap/status 投影、Profile/Manifest/readiness 不变式、I-002/I-003/I-004
verdict: conditional
status: recorded
parent: GOAL-006-r5-maintenance-read-only-gate
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# A-002 · R5 S0 设计独立审计

## 审计头

| 项 | 值 |
|----|----|
| source | independent |
| auditor | grok-build（grok-4.6 · reasoning high） |
| type | design-plan |
| 日期 | 2026-08-18 |
| scope | R5 S0：四模式、写门禁路由/认证优先级、错误语义、bootstrap/status 投影、Profile/Manifest/readiness 不变式、I-002/I-003/I-004 |
| verdict | **conditional** |
| required findings | 2（F-001 high / F-002 med） |

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`；Root `GOAL-001-shared-cross-module-contracts`；canonical `docs/workspaces/workspace-012-shared-cross-module-contracts/`；`shared_materials_catalog: none`；`primary_plan: VP-012-shared-cross-module-contracts`。
- **covered**：E-001 扫描事实、D-001/D-002 设计、I-001～I-005 相对 S0 门禁、以及用户点名的实现只读核对。
- **excluded**：S1/S2/S3 实施、关门审计、其他工作区上下文、D-002/00-meta/`status`/`progress`/业务代码修改。
- **P-005**：I-001/I-005 已 `verified`；I-002～I-004 为 required、最晚 S0 结束前、状态仍为 `collecting`。本意见不关闭信息项。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；无共享资料引用被当成事实 | `workspace.md` L8–16、L55–60；`goal-tree.md` L20–26 |
| E-001 的「无统一运行态 / 无热加载」成立 | `apps/api/internal/config/config.go` 无 `runtime`/`mode` 字段；`config.default.yaml` L15–38 无 runtime 节；全仓 `RUNTIME_MODE`/`SERVICE_*`/`availabilityMode` 在实现中不存在 |
| 核心 + Provider 路由汇入同一 mux，最终 handler 链是统一覆盖边界 | `composition.go` L423–425、L505–506；`server.go` L24、L35–54 |
| `RouteContribution.Public` 不参与装配，不能当豁免依据 | `kernel/contribution.go` L32；`composition.go` L423–425 只 `mux.Handle`；registrar `HTTP()` 原样追加（`provider.go` L211–222） |
| request-id 在最外层；CORS 先处理 OPTIONS | `server.go` L24、L48–51；`requestid/requestid.go` L59–69 |
| `WithJSONRouteErrors` 可在未注册路径/方法上产出 404/405 | `route_envelope.go` L10–28 |
| `/healthz` 不碰库；`/readyz` 只表达存储 + 模块图 readiness | `health.go` L59–77、L78–107 |
| public bootstrap 现固定 `normal`；Host 只接受四模式，无 `read-only` | `bootstrap.go` L36–45；`bootstrap.ts` L21、L141–143；`host-bootstrap.schema.json` L48–50 |
| system-monitoring status 是单行 list envelope，现无运行态字段 | `systemmonitoring.go` L27–36、L124–134、L140–150 |
| 既有错误包络可带 catalog code + `correlation_id` | `localize.go` L13–32 |
| I-003 公开写路径集合与 D-002 allowlist 对齐 | 无中间件的写路径仅 `POST /api/auth/login|refresh|logout`（`auth.go` L67–69）与 `POST /api/auth/mfa/verify`（`mfa.go` L67）；`POST /api/account/password` 是强制改密路径（`account_self.go` L60；`auth.go` L422–431） |
| 验证码预检是 GET，不受写门禁影响 | `captcha.go` L35–60 |
| Profile/Manifest/readiness 不变式在 D-002 中写明，且本设计未要求改装配算法 | D-002 §5；`health.go` L71–77；bootstrap 只改 availability，不改 Manifest bytes |

## 对照成功标准（S0 设计范围）

| 标准 | 状态 | 证据 |
|------|------|------|
| 四种运行态有单一来源、启动校验、确定性投影；非法配置 fail closed | 部分 | D-002 §1 已写来源/缺省/fail closed；实现尚未存在。空 env 与「空的显式值」若误用 `envOr` 会变成 unset 而非 fail closed（`config.go` L670–674） |
| maintenance/read-only 写拒绝覆盖核心及模块路由；GET/HEAD/探针与错误优先级符合冻结契约 | 部分 | D-002 §2 矩阵与门禁位置可实施；I-003 优先级可接受。I-002 的重试/消费语义对 degraded/read-only 不成立（见 F-002） |
| bootstrap 与 system-monitoring 消费同一状态；正常态兼容；degraded 不产生未知 capability | **未通过** | 未知 capability 可避免，但所选 `disabledCapabilities` 语义反转（F-001）；`host.readOnly` 在 READY_DEGRADED 路径不被消费 |
| 定向/全量测试；Profile/Manifest/pin/readiness 不变 | S0 不适用实施 | 设计未要求改 Profile/模块矩阵/Manifest 聚合；测试矩阵仍待 S1/S2（同意 A-001 F-001/F-002） |

## Findings

### F-001 · I-004 把 `form.controls.readonly` 放进 `disabledCapabilities` 是语义反转

| 字段 | 值 |
|------|-----|
| level | **required** |
| 严重度 | high |
| status | open |
| 影响门禁 | S0 方案冻结；I-004；S1 bootstrap 投影 |
| evidence | D-002 L31；`docs/schemas/capability-registry.json` L134–140；`bootstrap.ts` L273–284、L298–300、L317–318；`boot.ts` L66–77、L201–206、L256；`form-controls.ts` L336–349；`wallet.json` L6–16；`dictionary-entries.json` L6–16；`host-bootstrap.schema.json` L48–50、L67–72 |

`form.controls.readonly` 是协议能力：允许页面声明字段 `readOnly: true`（值仍参与提交）。Host `disabledCapabilities` 会从 `effectiveCapabilities` **剔除**该能力（`bootstrap.ts` L298–300）。因此 D-002 的 `degraded/read-only → degraded + disabledCapabilities:["form.controls.readonly"]` 并不表示「界面只读」，而是「Host 不再支持只读字段声明」。

当前生产路径里 `READY_DEGRADED` 直接进应用且 `effectiveCapabilities` 被丢弃（`boot.ts` L201–206、L256；`main.tsx` 只看 `failure === null`），所以该字段今天是空操作。若 S3 把 bootstrap 收窄后的能力集接到页面门禁，`wallet` / `dictionary-entries` 等已在 page meta 要求该能力的页面会 `FORM_CAPABILITY_REQUIRED`（`form-controls.ts` L344–349）。

I-004 要求「投影到既有 Host availability，且不引入未知 capability」。用已登记 ID 只满足后半句；前半句的投影对象选错了。正确方向：`degraded` / `read-only` 都映射为既有 `availability.mode=degraded`，**省略** `disabledCapabilities`；`read-only` 用 `messageKey` + status `availabilityMode` 区分；写阻断仍由 HTTP 门禁执行。不要用能力裁剪冒充运行态。

### F-002 · I-002 把 degraded/read-only 写拒绝的重试语义指到 bootstrap recovery，与现网 Host 行为不符

| 字段 | 值 |
|------|-----|
| level | **required** |
| 严重度 | med |
| status | open |
| 影响门禁 | S0 方案冻结；I-002；S2 错误契约 / S3 消费 |
| evidence | D-002 L22–24、L36；`bootstrap.ts` L273–276、L317–318；`boot.ts` L93–114、L201–206；`failure.ts` L124–126、L137–138；`resource.ts` L45–75 |

四模式的写阻断与 HTTP/code 表本身可执行：`maintenance/degraded → 503`，`read-only → 423`。但 §4 写「客户端使用 bootstrap 的 maintenance/manual 或 degraded recovery 语义」只对 **maintenance** 成立——Host 在 availability-gate 直接 `MAINTENANCE` 终态（`bootstrap.ts` L275）。`degraded` / 被投影的 `read-only` 是 `READY_DEGRADED`，应用继续运行，写请求会打到已挂载的 Resource 层；现有 `ResourceApiError` 按包络 `error`/`messageKey`/`correlation_id` 分类（`resource.ts` L45–75），**不会**走 Host failure recovery。

5xx 在 Host HTTP 映射里是 `unavailable`（`failure.ts` L124–126），与「降级仍可读」冲突。I-002 在 S0 结束前必须写清：受控态写拒绝是 API 包络错误、不承诺 `Retry-After`、不触发账号锁定/维护终态；bootstrap recovery 仅适用于 `mode=maintenance`。

### F-003 · `423 SERVICE_READ_ONLY` 与现有 `423 ACCOUNT_LOCKED` 共享状态码

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | med |
| status | open |
| 影响门禁 | S2 错误语义；S3 Web 消费 |
| evidence | D-002 L24；`auth.go` L126–128；`auth-client.ts` L328–331；`AuthContext.tsx` L147–152 |

登录路径把 **任意 423** 映射为 `ACCOUNT_LOCKED`（`auth-client.ts` L330–331）。login 在 allowlist 内，首版不会把 `SERVICE_READ_ONLY` 打到 login。但契约若只强调 HTTP 423，后续通用客户端容易复用该映射。S2 必须规定客户端以 catalog `error` 为准；S3 不得把业务写的 423 当成锁号终态。可考虑改用 `403`/`409`/`503` 之一，或在 D-002 明确「禁止按状态码分流」。

### F-004 · `messageKey: "host.readOnly"` 目前无法区分 read-only 与 generic degraded

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| 影响门禁 | S3 Host 消费（VP-012 允许 UI 后置） |
| evidence | D-002 L31；`boot.ts` L93–114、L201–206；`HostFailureScreen.tsx` L7–21、L92；`en-US.json` / `zh-CN.json` 无 `host.readOnly` |

`messageKey` 只进入 `terminalFailure`。`READY_DEGRADED` 的 `failure` 为 `null`，该 key 不会被读到；HostFailure 标题还按 `kind` 查表，不读 document `messageKey`。词典也没有 `host.readOnly`。S0 可用 status `availabilityMode` 作为权威区分；S3 若要横幅/文案，需单独消费 document，不能假定现网会显示该 key。

### F-005 · S1 配置加载必须证明「空的显式值 fail closed」，不能复用 `envOr`

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | med |
| status | open |
| 影响门禁 | S1 实施 |
| evidence | D-002 L15；`config.go` L132–136、L280–281、L670–674；A-001 F-001 |

现有环境覆盖把空字符串当未设置（`envOr`）。若 `RUNTIME_MODE` 用同一助手，空值会静默回落 `normal`，违反 D-002「空的显式值 fail closed」。同意 A-001 F-001，并补充：S1 测试必须覆盖 YAML 缺省、YAML 合法值、`RUNTIME_MODE` 覆盖、空字符串、未知值、解析错误。

### F-006 · S2 黑盒矩阵须钉死「已注册方法」判定与 allowlist 精确匹配

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| 影响门禁 | S2 实施 |
| evidence | D-002 L19–26；`route_envelope.go` L10–56；A-001 F-002 |

门禁放在 CORS 之后、`WithJSONRouteErrors` 之前，可用现有 `mux.Handler` 探测实现「只拦已注册的当前写方法」。allowlist 在「仅对 POST/PUT/PATCH/DELETE 生效」的前提下用精确 path 即可；实现不得前缀匹配。同意 A-001 F-002：核心与 Provider 各取一条真实写路径，并断言未知路径仍为 404/405。

## 必改项汇总

1. **F-001（required / high）**：重写 I-004 投影。禁止把 `form.controls.readonly`（或任何协议能力）当作运行态开关放进 `disabledCapabilities`。`degraded`/`read-only` 使用既有 `degraded` mode 且省略该字段；原始 mode 走 status `availabilityMode`。
2. **F-002（required / med）**：重写 I-002 重试/消费句。maintenance 可用 Host `MAINTENANCE` + manual retry；degraded/read-only 的写拒绝是运行中 API 包络错误，不走 bootstrap recovery，也不把 5xx 解释成 Host unavailable 终态。

未合法闭合前不得把 I-002/I-004 标为 `verified`，不得放行 S1。

## 与既有意见的异同

| 点 | A-001 self | 本意见 |
|----|------------|--------|
| 统一 mux / 404–405 / request-id+CORS 顺序 | 同意 | 同意 |
| allowlist 覆盖认证生命周期与强制改密 | 同意 | 同意；公开写路径已穷尽核对 |
| health/ready 与 status 追加字段 | 同意 | 同意；现有 statCard 按 `valueField` 取列，额外 JSON 字段不破坏页面 |
| Profile/Manifest/readiness 不变 | 同意 | 同意 |
| `form.controls.readonly` 投影 | 当作正确设计（pass，required=0） | **required F-001**：能力语义反转，不能关闭 I-004 |
| I-002 重试语义 | 未提出问题 | **required F-002** |
| S1/S2 测试门禁 | F-001/F-002 recommended | 保留并补充 F-005/F-006 |

A-001 为 `pass`，本意见为 `conditional`。不是 P-004.2 的 pass/fail 对撞；但对 I-004 投影主张构成「一侧视为已可接受、一侧标 required」。编排器应按 P-004 展示差异，不得静默取 self 的乐观侧。

## 信息项核对

| ID | 级别 | 最晚阶段 | 登记状态 | 本审计结论 |
|----|------|----------|----------|------------|
| I-001 | required | S0 结束前 | verified | 维持；E-001 与实现一致 |
| I-002 | required | S0 结束前 | collecting | **不能 verified**：写矩阵可用，重试/消费句错误（F-002） |
| I-003 | required | S0 结束前 | collecting | 设计充分：门禁先于 handler 内认证；业务写不泄露 401/403；allowlist 覆盖全部公开写路径；GET/HEAD/探针放行。闭合仍待编排器在响应 F-001/F-002 后一并处理 |
| I-004 | required | S0 结束前 | collecting | **不能 verified**：投影选错能力（F-001） |
| I-005 | required | S1 实施前 | verified | 本条即 cross independent；不改其状态 |

## 结论 + 建议给编排器/用户的下一步

S0 的运行态来源、统一写边界位置、认证优先级和 Profile/Manifest/readiness 不变式可以冻结。**不能**按 D-002 原文放行 S1：I-004 的 Host 投影与 I-002 的重试语义必须先改正。

建议 `/govern` 下一句：响应 A-002 F-001/F-002，修订 D-002 §3/§4（不要用 `disabledCapabilities` 表达只读；写清 in-app 错误消费），再决定是否关闭 I-002～I-004 并进入 S1。

## 声明

本意见不修改 `status` / `progress` / D-002 / 00-meta / goal-tree / 业务代码。响应与信息项闭合由 `/govern` 处理。
