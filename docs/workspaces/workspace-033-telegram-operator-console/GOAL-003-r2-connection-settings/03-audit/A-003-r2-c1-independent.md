---
doc_type: goal-audit
id: A-003-r2-c1-independent
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: design-plan
scope: R2 C1 方案与信息门禁（D-001 用户三项书面裁决及未选方案、I-033-014～016 verified 证据、R2 五阶段路线、D-002+D-003 实施边界、A-001 r2 entry self、A-002 r2 C1 self response、现有 Telegram/config/migration/UI/Fx 基线是否支持进入 C2/C3；DB authoritative + YAML/env 首次 seed + Admin PATCH、引用计数/每会话 TTL 至少 20 秒、getUpdates 30 秒请求与 40 秒独立 client、遗漏 required 信息冲突）
verdict: pass
version: 0.1.0
---

# A-003 · R2 C1 方案与信息门禁独立交叉审计（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan · `[workspace-033-telegram-operator-console]` `GOAL-003-r2-connection-settings` 的 R2 C1 方案与信息门禁（含 D-001 / I-033-014～016 / 五阶段路线 / GOAL-002 D-002+D-003 / A-001 / A-002 / Telegram 基线）
- **verdict**：pass
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。A-001、A-002 原文均未改写。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **covered**：D-001 三项用户书面裁决与未选方案；I-033-014～016 的 C1 门禁、证据与是否 residual；R2 C1～C5 路线；GOAL-002 D-002+D-003 实施边界是否被 D-001 改写；A-001 required F-001～F-003 是否可由 D-001 合法闭合；A-002 self response 是否越权改状态或把未发生实现写成事实；现有 Telegram runtime / YAML-env / migration / settings API / Admin UI / Fx `OnStop` 是否足以**进入** C2/C3；DB authoritative + 首次 YAML/env seed + Admin PATCH；引用计数与每会话 TTL≥20s；`getUpdates` 30s 请求 / 40s 独立 client；遗漏 required 信息冲突
- **excluded**：C2/C3/C4/C5 生产代码实施；真实 Telegram / 公网 webhook 运行态；把 `progress: 1/5` 当作完成或放行证据；其他工作区台账；愿景层 VRev 改写

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51；Root 在本区且 `parent: null` |
| D-001 记录用户对 I-033-014～016 的三项书面选择、未选方案与理由；属领域事实而非 residual/overrule | D-001 L12–22 |
| I-033-014～016 在 meta/decision 投影为 `verified`，无延期、无 residual | GOAL-003 `00-meta.md` L53–55、L60；`01-decision.md` L17–19；E-002 L14–16 |
| R2 五阶段串行路线存在；C1 为参数/信息，C2 配置 API，C3 client/manager，写集不重叠才可并行 | GOAL-003 `00-meta.md` L38–48 |
| D-001 声明不改变 D-002+D-003 产品边界，并要求实施同时引用三者 | D-001 L14、L22、L37、L40 |
| A-001 原文仍为 `conditional` / open required 3；A-002 为 self response `pass`，未改写 A-001 | A-001 L6–10、L32–51；A-002 L6–18、L39–41 |
| 现有基线仍无 mode/URL 列、无 `getMe`/`setWebhook`/`deleteWebhook`/`getUpdates`、无 connection manager；`OnStop` 未挂 Telegram drain | `migration.go` L11–16；`runtime.go` L15–22、L178–196；`http_sender.go` L19–20、L69–70；`composition.go` L828–834、L857、L1039–1058；telegram 包内无 `GetMe`/`SetWebhook`/`GetUpdates`/`Start`/`Stop` receiver |
| 本条未把 C2/C3 代码或 Fake Bot API 测试写成已完成 | E-002 L17；A-002 L37；GOAL-003 `02-execution.md` L22 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| I-033-014：DB authoritative + YAML/env 首次 seed + Admin PATCH；未选「每次 YAML 覆盖」与「DB 只读」 | **已达成（方案层）** | D-001 L18、L26；与现有 token/secret 行权威模型同构（`runtime.go` L78–80、L89–118） |
| I-033-015：未绑定 polling 用活跃控制台会话引用计数；每会话独立 TTL，基线 20s（覆盖两个 10s 心跳）；归零/全部失效后 Stop+drain | **已达成（方案层）** | D-001 L19、L29–30；对齐 VP-033「heartbeat 引用计数」（VP-033 L41、L61）；不收缩为单共享 lease |
| I-033-016：`getUpdates` 请求 timeout 30s；独立 polling HTTP client 40s；40>30；禁止复用 10s `sendMessage` client | **已达成（方案层）** | D-001 L20、L30；满足 D-003 L31「client timeout 严格大于请求 timeout」；现有 10s 预算在 `http_sender.go` L19–20、L33–35、L117 |
| D-002+D-003 边界未被 D-001 改写：缺省 polling、显式 URL、webhook 非空 secret、`idle`+`receiver=none`、唯一 manager、模式建立与 loop 分离 | **已达成** | D-001 L14、L22、L27–30；D-002 L16–20、L24–39；D-003 L16、L20–38 |
| R2 五阶段路线可指向 C2/C3；C1 信息门禁无到期未关闭 required | **已达成** | `00-meta.md` L38–48、L53–60；I-033-017～018 为 non-blocking open（最晚 C3） |
| 现有基线支持**进入** C2/C3（扩展点存在，实现尚未发生） | **已达成** | 见下「基线可进入性」；缺口属于 C2/C3 实施范围，不是 C1 信息缺失 |
| A-001 F-001～F-003 在 C1 方案层可闭合 | **closed/fixed（合同/信息层）** | 见关闭证据表；R2 代码证据仍待 C2/C3 |
| C2/C3 生产实现已完成 | **不在本条 scope**；事实为未开始 | E-002 L17 |

### A-001 F-001～F-003 关闭证据核对

| A-001 finding | A-002 声明 | 本条状态 | 可重复核对的修正 |
|---------------|------------|----------|------------------|
| F-001 · 配置来源优先级（I-033-014） | 信息项 verified | **closed/fixed** | D-001 L18、L26：行存在后 DB authoritative（空值也算）；仅首次无行 YAML/env seed；PATCH 部分更新 mode/URL |
| F-002 · heartbeat 引用计数/TTL（I-033-015） | 信息项 verified | **closed/fixed** | D-001 L19、L29：按会话引用计数；独立 TTL 基线 20s；单会话到期只移除自己；归零/失效 → idle + drain |
| F-003 · getUpdates timeout（I-033-016） | 信息项 verified | **closed/fixed** | D-001 L20、L30：请求 30s、独立 client 40s；正常空结果继续 loop；取消退出；Bot API/传输错误 fail closed |

上述 `closed/fixed` 仅覆盖 **C1 信息/方案缺口**。`00-meta.md` L53–55 的「补 … 测试」仍是 C2/C3 实施证据，不得把 I-033-014～016 `verified` 读成测试已发生。

### 三项数值与优先级逐项核对

| 用户关注点 | D-001 | 与 D-002/D-003 / 基线 |
|------------|-------|------------------------|
| DB authoritative | L18、L26：行存在后 mode/URL 空值也权威；不回落 YAML | 同构于 `runtime.go` L78–80「including empty ones」；D-002 L27 空 mode 按 polling 解释并在后续写回，二者相容（空 = 不 overlay YAML，再规范化为 polling） |
| YAML/env 首次 seed | L18、L26：仅 `COUNT(*)=0` | 现有 `initPersistence` 在无行时 INSERT token/secret（`runtime.go` L89–101）。C2 必须把 mode/URL 纳入**同一次**无行 seed，而不是把「已有行 + 新列为空」当成首次 seed（本条 F-001 recommended） |
| Admin PATCH | L18、L26：可更新 mode/URL；未提供字段保留；token/secret 仍加密 | 现有 PATCH 已是指针部分更新（`settings_handler.go` L22–26、L74–85），但只接受 token/secret。C2 扩展字段即可，不必改权限模型（`settings.read`/`settings.write`，L37–46） |
| 引用计数 + TTL≥20s | L19：至少覆盖两个 10s 心跳，实现基线 20s | 否决单共享 lease；已绑定 polling 仍常驻、不受 heartbeat 约束（L30）；webhook 不启动 polling。I-033-009 仍 non-blocking open（10s 默认在 D-002 L46），不重开 I-033-015 |
| getUpdates 30s / client 40s | L20、L30 | 40>30 满足 D-003 L31；明确禁止复用 `OutboundHTTPTimeout = 10s`（`http_sender.go` L19–20） |

未发现与 I-033-011～013、D-003 secret/`idle`/长轮询分流相反的 required 信息项。I-033-017～018 保持 non-blocking open，D-001 L34–35 已给出实施方向（不扩展 kernel 端口；不卸载 HTTP 路由），不构成 C1 必改。

### 基线可进入性（C2/C3）

| 扩展点 | 现状 | C2/C3 含义 |
|--------|------|------------|
| 单行 `telegram_config` + 加密 token/secret | `migration.go` L11–16；`runtime.go` L81–118 | C2 增加非敏感 `mode` / `webhook_public_base_url` 列与回读测试 |
| YAML/env 仅 token/secret | `config.go` L244–245、L417–421、L627–628、L749–750；`.env.example` TELEGRAM_BOT_TOKEN / TELEGRAM_WEBHOOK_SECRET | C2 增加 `telegram.mode` 与 `telegram.webhook_public_base_url`；禁止复用 `runtime.mode`（`config.go` L422–424、L807–808） |
| Settings GET/PATCH + write-only secrets | `settings_handler.go`；`provider.go` L71–94；UI `telegram-admin-tab.tsx` L14–20、L162–180 | C2 扩展 API；C4 才改 UI。现有 UI 无 mode/URL 不是 C1 缺口 |
| Fake Bot API 注入点 | `http_sender.go` L16、L33–38 `apiBaseURL` | C3 可为管理 API 另建 client；长轮询不得共用该 10s client |
| Fx `OnStop` | `composition.go` L1039–1058 | 接缝存在，尚未调用 connection manager `Stop`；C3 接入 |
| Dispatcher 占用位 | `kernel/telegram.go` L121–127 仅 Register/Unregister；`dispatcher.go` L10–17、L29–85 无只读探测 | 与 D-001 L34 / I-033-018 一致：C3 在具体 `*Dispatcher` 或内部 adapter 上实现 `HasBusinessHandlers()` |
| Webhook HTTP 路由常挂 | `provider.go` L96–108；secret 空则 401（`webhook.go` L113–121） | 与 D-001 L35 / D-003 F-001 相容；「不保留 receiver」= Telegram 侧 `deleteWebhook` + 不跑 polling loop |
| Profile gating | `disabled.go` L24–47；`composition.go` L841、L873–876 | I-033-017 可沿用；不重开默认 Profile |

`composition.go` L828–834 的 `TelegramRuntime.Manager` 仍是密钥热切换 `RuntimeManager`，不是 D-002 的 connection manager。这是已知命名陷阱（GOAL-002 A-002 F-008），C3 必须新增独立 owner，而不是给现有 `Manager` 加 `getUpdates`。

## Findings

### F-001 · C2 必须把「已有 telegram_config 行 + 新列为空」当成 authoritative 空值，而不是第二次 YAML seed

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| status | open |
| 关联 | I-033-014（方案已 verified）；`runtime.go` L78–80、L89–101 |
| evidence | 见下 |

D-001 L26 已写明：只有首次**无行**才 YAML/env seed；行存在后空 mode/URL 也是权威值。现有 `initPersistence` 在模块首次启动就会 INSERT `id=1`（即使 token/secret 为空）。因此升级路径上几乎总会「行已存在」。C2 迁移加列后若把新列 NULL/空当作「尚未 seed」并回写 YAML，会把用户否决的「每次启动 YAML 覆盖」从后门放进来，并与 I-033-012「已有配置缺省 polling」冲突。

C2 最小测试应覆盖：① 无行 → YAML/env seed mode/URL；② 已有 token 行 ALTER 后空列保持空/polling、忽略 YAML URL；③ PATCH 部分更新后重启仍读 DB；④ 空 mode 按 D-002 解释为 polling 并在后续写回显式值。本条不把该实施风险升级为 C1 required，因为优先级合同已经写清。

### F-002 · C3/C4 须钉死「会话身份」键；TTL 基线 20s 随心跳间隔耦合，不重开 I-033-015

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | I-033-015；I-033-009（non-blocking open）；D-001 L19、L29 |
| evidence | D-001 要求「活跃控制台会话引用计数」「单个会话到期只移除自己的引用」 |

C1 已否决单共享 lease，足以进入 C3。实现仍需选择计数键（建议：控制台连接/heartbeat client 身份，而不是 Admin 用户 id），否则同用户多标签仍会误停。TTL「至少覆盖两个 10s 心跳、基线 20s」与 D-002 L46 的 10s 默认一致；若 C4 调整 I-033-009 间隔，TTL 应保持 ≥2×interval，而不是冻结死 20s 却改间隔。不构成 C1 信息冲突。

### F-003 · C3 实施入口必须合并 D-002+D-003+D-001，并避开现有 10s client 与 Manager 命名

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| status | open |
| 关联 | GOAL-002 A-002 F-008；A-004 F-001/F-002；D-003 L16、L31 |
| evidence | 见下 |

D-001 L37 已要求三份合同同时引用。仍须在 C3 代码/测试中落实，否则会把已闭合的 D-003 分流再实现错：

- D-002 L36/L56/L66 原文仍把 polling 写成 `deleteWebhook` 后立刻 `getUpdates`、并把「超时」写入失败表；冲突以 D-003 L16 为准。
- `http_sender.go` L19「D-002 §4」与 `webhook.go` L22–29、L83「GOAL-002 D-002 §6」**不是本区 D-002**（本区 D-002 无 §4/§6、未冻三桶限流）。C3 不得把这些注释当成本区合同。
- 新 polling client timeout=40s，请求 `timeout=30`；`getUpdates` 正常空结果/正常等待继续同一 loop（D-003 L28）。
- HTTP 409/429 等未逐条枚举；默认走 D-003 L30 非正常传输错误 fail closed。若要对 429 做有界重试，必须另写决策。
- `getUpdates` offset 可在 C3 进程内内存持有；跨重启重投不在本 C1 信息门禁内，不得静默做成多实例抢 offset。

### F-004 · GOAL-003 下一决策号 `D-002` 会与 R1 实施源同名

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| evidence | GOAL-003 `01-decision.md` L29「后续决策从 `D-002-*.md` 起递增」；实施源是 GOAL-002 `D-002-r1-contract-freeze` / `D-003-r1-audit-correction` |

区内 `parent` 仍用短 id 合法，但 C2/C3 计划与 PR 若只写「按 D-002」会读错合同。建议引用时用 Q2 路径或 `GOAL-002` D-002/D-003 + 本目标 D-001。本条不要求改编号规则。

### F-005 · A-002 self 把 GOAL-002 A-002 的 F-004～F-009 误写成「A-004 F-004～F-009」

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| evidence | GOAL-003 A-002 L41；GOAL-002 A-004 仅有 recommended F-001/F-002；F-004～F-009 在 GOAL-002 A-002 L123–193 |

不改变 A-002 原文，也不把那些 recommended 升级为 required。C2～C5 计划应继续追踪：**GOAL-002 A-002 F-004/F-005** 已由本目标 D-001 在方案层回应，仍待代码/测试；**F-006/F-007** 已由 D-001 L34–35 定向，待 C3/C4；**F-008** 见本条 F-003；**F-009** 仍是「不要把未实现写成失败」。

## 必改项汇总

无。本条 **open required = 0**。

I-033-014～016 保持 required **verified**（用户书面裁决 + D-001 实施合同）。无 `accepted-residual` / `user-overruled`。无到期且影响 C1 的 required 信息项。未发现与 D-002+D-003 相反、需要 P-004 重裁的 required 信息冲突。

A-001 历史条目仍保留其原始 `open`/`required` 表述；闭合效力在 A-002 self response + 本条复审，不靠改写 A-001。

## 仍开放但不构成本条必改

- 本条 F-001～F-005：recommended，供 C2/C3 实施与测试。
- I-033-017～018：non-blocking open，最晚 C3。
- I-033-009/010：Root non-blocking；009 不阻断 TTL 基线，010 属 R3。
- GOAL-002 A-002 F-004～F-009 与 A-004 F-001/F-002：recommended，方案层已部分回应，代码证据在 R2 对应阶段。

`workspace.md` L40–45 仍写 R1「进行中」、R2「待开始」，与本区 goal-tree / GOAL-002 `done` 不一致。工作区状态权威不在 `workspace.md`；本独立审不改该文件，建议 `/govern` 顺手刷新纲领表，不作为 required。

## 与既有意见的异同

| 项 | A-001 self | A-002 self response | 本条 independent |
|----|------------|---------------------|------------------|
| 原文是否保留 | 保留（仍 conditional / F-001～F-003 open） | 声明保留且未改写 A-001 | 核对属实 |
| I-033-014～016 | open required | verified | **同意 verified**（方案层）；测试仍待 C2/C3 |
| A-001 F-001～F-003 | required / open | 无开放 required（未逐条写 F-id 关闭表） | **closed/fixed**（合同/信息层） |
| verdict | conditional | pass | **pass**（本 C1 scope） |
| open required | 3 | 0 | **0** |
| 可进入 C2/C3 实现 | 否（先裁决） | 是，但先跑本 independent | 同意：本条通过后由 `/govern` 放行实现；本条不改检查点 |

本条与 A-002 在「三项用户裁决真实、非 residual、D-002+D-003 边界未改、代码未发生」上一致；补上 A-001 finding 的逐条闭合核对，以及 C2 升级路径/命名陷阱的 recommended。A-001 的 `conditional` 与 A-002 的 `pass` 是先后响应关系，不是并存的 P-004 结论冲突。

## 结论 + 建议给编排器/用户的下一步

R2 C1 方案与信息门禁可重复核对：**verdict = pass**，**open required = 0**。现有 Telegram/config/migration/UI/Fx 基线足以进入 C2/C3，但不能跳过它们。C2 可在 `/govern` 响应后开工；C3 的 Bot API client 可与 C2 并行，但 manager 的持久化/回读接线应等 C2 schema（`00-meta.md` L46「部分依赖 C2」）。

建议 `/govern`：

1. 响应本条：将 A-001 F-001～F-003 视为已合法闭合（`fixed`）；I-033-014～016 保持 `verified`；不要重开 I-033-011～013。
2. 按 D-001 + GOAL-002 D-002 + D-003 进入 C2（配置 schema、迁移、runtime 回读、settings API）。不要把本条 recommended 或 I-033-017～018 当成 C2 阻断。
3. C3 实施时显式处理本条 F-001（已有行不加 YAML overlay）、F-003（独立 40s polling client、新 connection manager、D-003 分流），并在 C3/C4 钉死会话身份键（F-002）。
4. 不要把 `progress: 1/5` 或 C1 检查点展示当作本条闭合或 C2 测试证据。

## 声明

本意见 `source: independent`，不修改 status/progress/检查点/goal-tree/decision 正文或生产代码。响应与是否进入 C2/C3 由 `/govern` 处理。
