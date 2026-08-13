---
id: E-004
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: S4 生产实现 — 上游 2.8 Host/App 互操作 pin、Host 模块与浏览器级证据（正式 2.8.0）
status: recorded
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# E-004 · S4 生产实现 — 上游 2.8 Host/App 互操作 pin、Host 模块与浏览器级证据（正式 2.8.0）

## 背景

上游 `schema-ui-docs` 的 ADR-0034～0037 于 2026-08-13 accepted（H1 变更集 commit `3936cf9`），
H2 机器契约同日落地（commit `453008d`：B0/F0/C0 Schema、`capability-registry.json`、B1/F1/C1
validator、host-bootstrap 23 / host-failure 43 / host-conformance-claim 30 正反 fixtures、
JS/Python 双 reference）。本仓 I-003（上游协议已发布/固定并进入本仓）按维护者裁定以
**固定上游 commit 的候选机器契约**满足 S4 开始条件：正式 2.8.0 发布前，本仓实现与证据均为
候选绑定，不构成生产支持声明。

## 已完成事实

### 1. 上游制品 pin（`provenance-v2.8.json`）

- 新增 vendored 制品：`docs/schemas/{host-bootstrap,host-failure,host-conformance-claim}.schema.json`、
  `docs/schemas/capability-registry.json`、三个 2.8 候选 fixture suite、重 pin 的
  `app-manifest.schema.json` 与 `app-manifest.cases.json`（returnIntentQueryKeys + capability ID
  连字符语法，上游 `453008d`，全部含 canonical sha256 记录）。
- 2.7 pin（`provenance.json`）以 note 记录两项重 pin 事实，其余条目保持 2.7.0。

### 2. 生产 Host 模块（`apps/web/src/host/*`，真实入口消费）

- `bootstrap.ts`：`discoverBootstrapDocument()`（`credentials: omit`、404/410 fallback、
  Content-Type/解析 fail-closed）+ `evaluateBootstrap()`（确定性生命周期，逐字段对齐上游 B1）；
- `failure.ts`：`classifyHostFetch`（§2.7 分类优先级）、`mapBootstrapResult`（0035 D7）、
  `validateFailure`（kind/hostCode 配对、scope/kind 兼容、retry/recovery 过滤）、
  `validateReturnIntent`（敏感 key 永久拒绝、allowlist 收窄）；
- `claim.ts`：`canonicalize`/`claimDigest`（D1a，跨语言一致）、`validateRegistry`、`validateClaim`
  （§4.8 顺序，vendored registry）；
- `boot.ts`：生产阶段编排——availability/auth 终态在 manifest 获取前裁定；manifest 装载失败保持
  ADR-0025 语义；integrity 以真实响应字节核验（Go API 同字节组装，声明的 `manifest.sha256`
  与真实响应一致）；degraded 收窄按 D5 执行。

### 3. 生产 UI 集成

- `main.tsx`：bootstrap discovery → `HostBootGate` → 终态 `HostFailureScreen` /
  `AuthGate`/`ManifestFailure`（各自保持原有语义）；`RenderFailureBoundary` → `HOST_RENDER_FAILED`
  （manual retry，无自动循环）；
- `HostFailureScreen`：`main` landmark 唯一错误标题、首次终态 focus 标题、assertive/polite
  live-region 按 kind、同 `failureId` 不重复播报、键盘可达 recovery actions；
- route 404 → `HOST_ROUTE_NOT_FOUND`（App shell 内 section 投影）；恢复动作后 focus 落到恢复
  页面主标题（§3.8 规则 5）；
- `app-manifest.ts`：Host 支持 `2.7 + 2.8`（strict 协商）、`returnIntentQueryKeys` M1 门控
  （`PROTOCOL_VERSION_TOO_LOW` / `MISSING_REQUIRED_CAPABILITY` + detail）、上游 M1 envelope
  （`CAPABILITY_REQUIRED` + detail）——vendored app-manifest suite 从此**零排除**执行。

### 4. Go API 生产 App 侧（`apps/api`）

- `handler.RegisterBootstrap`：`GET /.well-known/schema-ui/host-bootstrap.json`，与 manifest
  handler 同字节组装，`manifest.sha256` 始终等于真实 manifest 响应 bytes 的 SHA-256；
  注册于 composition（`core.manifest-route` 变更集内）。

### 5. Conformance 消费与回归

- `upstream-host-fixtures.test.ts`：上游三 suite **96 fixtures 零排除**（host-bootstrap 23 + host-failure 43 + host-conformance-claim 30）逐字段核对生产模块
  （+ sha256 pin 校验）；
- app-manifest 41 + app-navigation 16 零排除（原 2 个 envelope 排除项已消除）；
- `claim-artifact.test.ts`：构建生成的 claim 通过 C0（Ajv）、C1（§4.8 → `CLAIM_OK`）、
  evidence sha256 绑定、D1a canonical digest 复现一致；
- 全量回归：apps/web vitest **857 通过**、tsc 无错误、Go `internal/handler` +
  `internal/composition` 测试通过。

### 6. 浏览器级证据（Playwright，真实双服务）

`e2e/host-failure.spec.ts`（4 tests，全过；全量 e2e 7 通过 + 1 既有 skip）：

- maintenance 终态：manifest **未**被获取（阶段顺序）、polite 播报、focus 标题、Retry 重建实例；
- protocol-rejected 终态：assertive 播报、无“继续渲染”动作；
- route not-found：`HOST_ROUTE_NOT_FOUND` surface、home 恢复、focus 落回页面主标题；
- 真实入口正常 boot：Go API 服务真实 bootstrap document（`bootstrapVersion 1.0` +
  `manifest.sha256` 十六进制），shell 正常渲染登录门。

### 7. Claim 与 evidence 绑定（正式 2.8.0 pin）

- `scripts/generate-claim.mjs`（`prebuild` 挂载）：`public/protocol/conformance-claim.json` +
  `conformance-local-report.json` + canonical digest 文件；
- 绑定值（上游正式 tag `v2.8.0` @ `521cff8`，上游审计 0080 V379 权威）：`protocolArtifact.contentSha256` = 正式制品
  contentDigest `4fae4605…`；`conformance.fixtureSha256` = 正式 fixture 树 digest
  `7aacf133…`；`host.buildId` = `git:<HEAD>`（claim/report/evidence 三处逐字一致，
  由 `claim-artifact.test.ts` 门禁；每次生产修正后重生成，绑定修正所在 commit）；
  2026-08-13 身份纠偏后重生成（E-005；此前 `593f625`/`40690917…` 为 H4 预备身份）。
- vendored host 三 suite 的文件 digest 已按正式 release 原样重 vendor（
  host-bootstrap `bfc71bbd…` / host-failure `59efc00a…` / host-conformance-claim
  `e3511ecc…`，语义内容与 H2 相同，字节级 pin 对齐正式 tag）；
- suites：app-manifest / app-navigation / host-bootstrap / host-failure /
  host-conformance-claim 全部 `pass`（CI 实际运行结果）；
- evidence kind `local-report`，报告 bytes 以 sha256 绑定且与 `subjectBuildId` 一致。

## 已登记 residual（不得冒充生产支持）

1. **（已关闭）候选绑定**：上游 2.8.0 正式发布后已按正式 artifact/fixture digest 重 pin 并重生成
   claim（见 §7）；当前 claim 绑定正式 2.8.0 digest。
   > **纠偏注记（E-005，2026-08-13）**：§7 原记录绑定的 `593f625` / content `40690917…`
   > 经上游审计 0080 V379 裁定为 H4 预备身份，非正式 pin；正式 v2.8.0 = tag `521cff8` /
   > content `4fae4605…`（fixture `7aacf133…` 不变）。已按正式身份重 pin 并重生成 claim
   > （buildId `git:fd641c6…`），vendored 工件与正式 tag 字节级一致。
2. **页面协议 2.7 mandatory behavior**：R5 已登记的 multi-round `$deps` reactions 子集未实现；
   闭环前 claim 的 `pageVersions` 条目视为候选绑定。
3. **返回意图消费链**：`validateReturnIntent` 已实现并消费上游 fixtures；登录流程接入
   return intent 的端到端恢复尚未上线（后续 S4 迭代）。
4. **304/ETag 缓存复用**：生产实现采用 200-only 装载（合规路径），conditional GET 复用为
   可选优化，未实现。
5. **session adapter 终态映射**：当前会话模型只有 loading/authenticated/unauthenticated；
   bootstrap 归一化的 `reauth-required` / `locked` 输入尚无生产来源（映射层已就位），后续
   认证迭代接入。
6. **hostOwnedPaths 封闭列表**：当前生产 shell 没有额外 Host 自有 path——未认证时登录页由
   `AuthGate` 在 shell 外拦截（不会误产生 `HOST_ROUTE_NOT_FOUND`），认证后未命中 manifest
   route 一律产生协议 404；显式 `hostOwnedPaths` 集合与对应 e2e 待未来出现宿主页时补。

## 下一步

- 上游后续 MINOR 发布时按同流程重 pin、重生成 claim 并重跑全部证据；
- S2 出口门禁（I-001/I-002/I-005/I-006）与 cross 方案审视按 GOAL-004 计划继续；
- 本仓 W3 阶段进度更新见 `02-execution.md` 执行索引与事实边界。
