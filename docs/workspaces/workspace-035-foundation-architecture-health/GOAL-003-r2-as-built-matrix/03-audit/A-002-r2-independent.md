---
id: A-002-r2-independent
doc: audit-entry
parent_goal: GOAL-003-r2-as-built-matrix
source: independent
auditor: grok-build/grok-4.6 high
type: stage
audit_type: execution-facts
scope: R2 as-built 对照矩阵及验证证据（R1 17 面 + Web；Provider Persistence 全局路径；assembly B+；Profile 数量；文档差异；限定测试）
verdict: pass
status: recorded
parent: GOAL-001-foundation-architecture-health
created: 2026-09-09
updated: 2026-09-09
version: 0.1.0
---

# A-002 · R2 as-built 矩阵独立审计（2026-09-09）

- **source**：independent
- **auditor**：grok-build/grok-4.6 high
- **类型**：stage / execution-facts
- **scope**：R2 as-built 对照矩阵及验证证据（R1 17 面 + Web；Provider Persistence 全局路径；assembly B+；Profile 数量；文档差异；限定测试）
- **verdict**：pass
- **候选 / 代码基线**：候选 `ca5caa7e`；代码基线 `ebe6013c`（相对 `5c341ec7` 仅 VP-035 治理文档；相对 HEAD 仅本区 GOAL-003 产物）

## 范围与区间

- **工作区**：`workspace-035-foundation-architecture-health`；Root `GOAL-001-foundation-architecture-health`；canonical `docs/workspaces/workspace-035-foundation-architecture-health/`；`shared_materials_catalog: none`。未读取其他工作区上下文。
- **覆盖**：GOAL-003 矩阵、验证记录、R1 分母包含行、对应 `apps/api` / `apps/web` 代码与本审重跑的限定测试。
- **排除**：R3 业界对照分类冻结、R4 路线图 editorial、生产整改、live S3/Telegram/Resend 验收、全仓 `-race` / E2E、Charter/VP 状态变更。
- **立场**：治理文档只约束范围；不以 A-001 self `pass` 为依据。未跑 live 供应商或 gated 未实现 ≠ 本有界评估自动失败；无证据也不宣称它们通过。

## 工作区与信息门禁

| 项 | 结论 | 证据 |
|----|------|------|
| 工作区绑定 | 匹配 | `workspace.md`：id / Root / canonical / `plan_refs`+`primary_plan`=VP-035 |
| 共享资料 | 无引用被当事实 | catalog `none`；矩阵未引用跨区资料关闭证据 |
| I-035-001/002/004/005 | verified，R2 可用 | Root `00-meta.md`；GOAL-002 冻结附件 |
| I-035-003 | collecting，最晚 R3 | 未用于本阶段分类放行；矩阵未接受 Charter 残余 |
| 红线 | 未见越界整改 | `git diff ebe6013c ca5caa7e` 仅本区治理文件；无 Redis/MQ/K8s/ORM 实现、无 Profile 默认集改动 |

## 成果（有证据）

| 主张 | 本审证据 |
|------|----------|
| 代码基线未改生产源码 | `git diff --stat ebe6013c ca5caa7e`：15 个文件，全部在 workspace-035 GOAL-003 / Root E-003 / goal-tree |
| 矩阵覆盖 R1 17 面 + Web | 附件 `as-built-matrix.md` 行 01–17 + W1；与 GOAL-002 `r1-denominator-freeze.md` §1.1/§1.2 对齐 |
| Persistence 为编译期全局入口 | `kernel/persistence.go:49` `CollectPersistence`；`kernel/provider.go:25-27` Registrar 无 Persistence；`modules/compiled/persistence.go:26-51` 全候选（含 telegram/digital-offer）经 `CollectPersistence`，`composition.go:203-206` open 前收集 |
| Profile 数量 mvp 11 / admin 23 / demo 12 | `kernel/profile.go:26-113` 逐项计数；`kernel/kernel_test.go:21,45` 断言 mvp/demo 列表 |
| assembly B+ 限制 | `assembly/assembly.go:26` `OpenStore`→`kernel.Store`；`:37-38` `NewAuthenticator` 返回 `*auth.Authenticator` 且调用单密钥 `auth.New`；组合根双密钥在 `composition.go:254-259` |
| 文档差异未改代码 | G-001/G-002/G-003 可在 docs 复现；G-004 为工厂限制候选；表头写明 R3 前不授权整改 |
| 限定测试 | 本审重跑：Go 15 包 `ok` + `assembly` `[no test files]`；Web 4 files / 48 tests passed。命令与附件 `validation.md` 一致 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 冻结分母逐行覆盖，区分实现事实与文档差异 | 达到 | 下表 01–17+W1；G-001～G-004 单列候选，未并入 residual 接受 |
| C2 针对性验证完成，明确未覆盖项 | 达到 | `validation.md` 证据限制；本审重跑短输出一致；未把 fake S3 / 未专属 live 服务写成验收 |
| C3 self 后 independent、无开放 required | 本条正在落盘 | A-001 self 已在；本意见 required=0。关门仍须 `/govern` 响应，本条不改 status |
| 无越界分类 / 整改 | 达到 | 未改端口/Profile；未把排除项（Admin 体验、gated 缺席、业务规则）写成缺口；R1 12 入册/3 排除未改写 |

## 逐行核对（R1 分母）

| # | 矩阵主张 | 独立结论 | file:line / 命令 |
|---|----------|----------|------------------|
| 01 | 薄内核、Registrar 六项、Persistence 独立、fail closed | **属实**。`KernelAPIVersion="2.0.0"`；Registrar 六方法无 Persistence；Resolve 缺依赖 fail closed 不补齐 | `kernel/module.go:11,265-298`；`kernel/provider.go:19-35,70`；`kernel/persistence.go:49`；`modules/users/provider.go:65-67` |
| 02 | mvp11/admin23/demo12；compiled 不静默进默认 | **属实**。telegram / digital-offer 在 `BuiltinModules` 不在 `profileDefaults`；custom 无显式集失败 | `kernel/profile.go:25-113,141-148,208-227`；`kernel/kernel_test.go:11-50` |
| 03 | Fx 在组合根；B+ 不能宣称全部工厂返回 kernel 接口 | **属实**。Fx `Provide` 在组合根；`NewAuthenticator` 单密钥 + 内部类型 | `composition.go:111-158,509-678`；`server.go:17-30`；`assembly.go:26-38` |
| 04 | Store 双方言、无 `*sql.Tx` 公共句柄；池默认 4/内存 1；ErrNoRows 为 sentinel 别名 | **属实**。公开 `Store`/`Tx` 无驱动句柄；`database/sql` 仅 sentinel 别名。G-002 池化 vs 文档 `MaxOpenConns=1` 成立 | `kernel/store.go:20-47`；`internal/store/store.go:29,104-114,141-158`；`internal/store/postgres.go` 存在 |
| 05 | ObjectStore 本地/S3；公共面无 `os.File`；fake≠live | **属实**。端口用 `[]byte`/`io` 语义；S3 适配在 internal；`TestS3TransportErrorsPropagate` 用 fake client | `kernel/objectstore.go:44-63,108-114`；`internal/objectstore/local.go:119`；`s3.go:47-67,121-128`；`s3_test.go:270-276` |
| 06 | Cache 内存默认；go.mod 无 Redis；接缝未实现 | **属实**。未把缺席写成缺陷 | `kernel/cache.go:35-54,102-138`；`internal/cache/memory.go:64,227-235`；`composition.go:854-862`；`apps/api/go.mod` 无 redis |
| 07 | RateLimiter 同时保留旧/新原子 API；§2.6 仍写 Allow/Record/Clear | **属实**。`Cancel` 只删匹配 token。G-003 文档滞后成立 | `kernel/ratelimit.go:37-69`；`internal/ratelimit/memory.go:83,174,187-199`；`docs/architecture/cache-redis-seam-and-track.md:62-73` |
| 08 | EventBus 进程内；无跨进程可靠消息 | **属实**。未据缺席判缺陷 | `kernel/eventbus.go:43-58`；`internal/eventbus/memory.go:72-88,243-260`；`composition.go:871-876` |
| 09 | Mail 端口/渠道隔离；assembly 默认 OutboxSink | **属实**。SMTP 失败保留旧渠道；CaptureSink 非现行默认 | `kernel/mail.go:35-74`；`internal/mail/runtime.go:112-114,362-370`；`assembly.go:41-43`；`composition.go:1022-1032` |
| 10 | metrics/traces 默认关；exporter 失败退 noop | **属实**。专用 listener 在 `internal/obs/server.go`，不是 HTTP `internal/server/server.go` | `internal/config/config.go:495-499`；`internal/obs/tracing.go:48-62`；`internal/obs/server.go:21-47`；`composition.go:254-296` |
| 11 | 停机顺序 HTTP drain→retention/metrics→jobs→eventbus→runtime→store→tracing；Windows 可跑 HTTP/Job 单测 | **属实**。`shutdown_drain_test.go` 无 `!windows`；PG 子测门控 | `kernel/lifecycle.go:18-29,68-78`；`composition.go:1151-1170`；`shutdown_drain_test.go:134-191` |
| 12 | Job 六态 ≠ EventBus | **属实** | `internal/jobs/model.go:19-25`；`runner.go:125-148,185-211` |
| 13 | Telegram 端口无 SDK；disabled 为显式形态 | **属实**。kernel 仅本仓消息类型 | `kernel/telegram.go:45-80`；`internal/channel/telegram/disabled.go:9-36`；`runtime.go:91-99`；`composition.go:607-641` |
| 14 | 组合根注入 current+previous；B+ 仅 current | **属实**。未重做 dump | `composition.go:254-259`；`internal/auth/auth.go:151-168,424-436`；`assembly.go:34-38` |
| 15 | Manifest 聚合；group 可选；top/user 不参与 sidebar 归一化 | **属实** | `internal/manifest/manifest.go:40-43,311+`；`composition.go:678,750-771`；`manifest_test.go:389` |
| 16 | 六份 architecture 文档：主线大体一致；overview/roadmap/§2.6 滞后 | **属实**。`module-architecture.md` `vision_ref @0.2.0` 不覆盖现行 Charter@0.4.0 链 | overview「当前阶段」仍 Charter@0.2.0 / VP-005 / FastAPI `web/`；roadmap:104,138 `MaxOpenConns=1` 且 RT-D02 现状锚点仍写「无明确 drain 合同」但状态 delivered；directory-layout / monorepo-layout 指向 `apps/` |
| 17 | 泄漏抽检：modules 生产签名无目标供应商客户端 | **属实**。本审补查 `internal/handler`：`*sql.Tx` 仅测试。digital-offer 注入 kernel Telegram 端口 | `modules/digitaloffer/provider.go:26-37,69-70`；modules 生产 `provider.go` 无 `*sql.Tx`/redis/s3/pgx 客户端 |
| W1 | Web Shell 消费 Manifest；无中央业务导航表 | **属实**。自定义渲染组件自注册 ≠ 导航注册 | `apps/web/src/main.tsx:10-12,76-79`；`App.tsx:1049`；`navigation.ts:224-238`；`protocol/load-page.ts:81-90` |

## 专项核对

### Persistence 全局路径

编译期全候选，不以 Profile 启用过滤：`compiled.PersistenceProviders()` 含 `telegrammigration` / `digitaloffermigration`；`PersistenceCatalog()` 调用 `kernel.CollectPersistence`。与 playbook M5 / architecture §4.1 一致。Registrar 六项 ≠ 六个 `Register` 调用——矩阵纠正成立。

### assembly B+

公共三工厂中仅 `NewAuthenticator` 返回 internal 指针且只接 current。不能把抽检无 SDK 命中解释成 B+ 已全部 kernel 化。G-004 作为候选而非故障宣称，符合范围。

### Profile 数量

| Profile | 数量 | 本审计数来源 |
|---------|------|----------------|
| mvp | 11 | `profile.go` + `TestBuiltinProfilesResolveDeterministically` 精确列表 |
| admin | 23 | `profile.go:46-92` 手工枚举（该测试只断言 Resolve 成功，未列 23 项） |
| demo | 12 | mvp 能力面 + `dev.examples`；测试精确列表 |

`channel.telegram` / `biz.digital-offer` 不在三套默认集。`config.default.yaml:21` 默认 `profile: mvp`，modules 块注释掉。

### 文档差异与分类边界

| ID | 本审是否复现 | 是否越界 |
|----|----------------|----------|
| G-001 | 是：overview 现时节 Charter@0.2.0、VP-005 active、产品树 `web/` FastAPI | 标 R4 文档卫生；明确不能用来判断产品未实现 |
| G-002 | 是：roadmap 现状锚点/`RT-P04` 仍 `MaxOpenConns=1`；代码 `sqlitePoolDefault=4`；RT-D02 行内「无明确 drain 合同」与 delivered 并存 | 未把池化扩写成 replica/读写分离已交付 |
| G-003 | 是：§2.6.1/2.6.2 只写 Allow/Record/Clear 与 INCR/GET/DEL | 「建议现在修文档」是候选措辞，表头禁止本阶段整改；未冻结 Redis 实现 |
| G-004 | 是 | 留 R3/用户；无消费者失败证据 |

R1 residual 12 入册 / 3 不入册未改写、未自动 `accepted-residual`。

### 限定测试（本审重跑，非转录 self）

Windows；`GOCACHE`=`%TEMP%\schema-ui-w035-audit-gocache`；未打开 `.env`、未输出密钥。

| 命令 | 结果 |
|------|------|
| `go test ./kernel ./internal/cache ./internal/ratelimit ./internal/eventbus ./internal/objectstore ./internal/mail ./internal/obs ./internal/auth ./internal/jobs ./internal/store ./internal/manifest ./internal/composition ./assembly ./internal/channel/telegram ./internal/config ./internal/server -count=1 -timeout=180s` | exit 0；15 包 `ok`；`assembly` `[no test files]` |
| `npm test -- src/app/navigation.test.ts src/protocol/app-manifest.test.ts src/protocol/load-page.test.ts src/app/App.integration.test.tsx` | 4 files / 48 tests passed |

矩阵列出的代表测试名均在源码中存在。短输出只证明包通过，不能统计每个子测是否执行——与 `validation.md` 自陈一致。

`pgtest` 会从 gitignored `apps/api/configs/.env` 读取 **仅** `PG_TEST_*` 键。本审另跑 `-run TestShutdownDrainHarnessPostgres` 时该用例 **实际执行并通过**（约 3.3s），不是 skip。这只说明本机在 DSN 可解析时门控 PG drain 能跑；**不**等于本任务专属 PG 验收，也 **不**外推 live S3/Telegram/Resend。矩阵未把它们写成已通过。

未实现 Redis/MQ 客户端：`apps/api/go.mod` 无 redis。按 R1 §1.2，缺席不是本 VP 缺口。

## Findings

### F-001 · 若干 file:line 是邻近锚点而非精确语句

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| 描述 | 主张本身成立，但个别行号需顺读上下文才能对上：`composition.go:601` 是 wallet 装配，Telegram 运行时共享在 `607-618`、channel 装配在 `641`；`composition.go:198` 为空行，catalog 在 `203-206`；`kernel/objectstore.go:30` 是类型声明，ID 校验在 `60-63`；行 10 的 `server.go:37` 未限定包——HTTP `internal/server/server.go:37` 是 CORS，metrics 专用 listener 在 `internal/obs/server.go:21-37`。不构成分母漏行或错误分类。 |
| 证据 | `as-built-matrix.md` 行 03/04/05/10/13；上述源文件 |

无 required finding。

## 必改项汇总

无。required = 0。

## 与既有意见的异同

- **A-001 self `pass`**：覆盖面与边界声明与本审一致。本审不继承其 verdict；独立读代码并重跑测试后同意 C1/C2 与无越界整改。
- **差异**：本审补了 PersistenceProviders 全候选清单、admin=23 的手工计数、handler 抽检、以及 PG 门控测在本机可能实跑而非 skip 的观察；增加 recommended F-001（行号精度）。
- **无冲突**：不否决 self 的「评估产物可交下一审」；也不把 self 或本 pass 当成各模块生产就绪证明。

## 结论 + 建议给编排器/用户的下一步

**verdict: `pass`。** R2 矩阵在冻结分母上真实、完整；实现事实与文档差异已分开；验证界限诚实；未见越界分类或生产整改。

建议用 `/govern` 响应本意见（及 A-001）：闭合 F-001（可 `accepted-residual` 或 R3 前顺手校正行号），然后在无开放 required 时关闭 GOAL-003 C3，将 Root R2 标 completed。不要在响应里把 G-001～G-004 冻结成 R3 分类，也不要开始改端口。

## 声明

本意见 `source: independent`。不修改 `status` / `progress` / goal-tree / 方案正文 / 生产代码。未 git commit。未调用其他审计 provider。响应由 `/govern` 处理。
