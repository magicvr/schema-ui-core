---
doc_type: staging-evidence
title: R3 as-built 证据锚点（精确行号）
status: staged
created: 2026-09-09
updated: 2026-09-09
parent: GOAL-004-r3-industry-comparison
role: repo-input
version: 0.1.0
---

# R3 as-built 证据锚点（精确行号）

本文件是 R3 的**只读取证输入**：为四类参照集、gated/缺席能力、12 条入册 residual、G-001～G-004 提供可逐行核对的 `file:line` 锚点。**不是** finding、**不是**审计意见、**不关闭任何门禁、不改分类、不改 R2 矩阵**。

读数约定：每个锚点给出 `file:line` + 该行**原文**（源码行按去缩进后的字面内容引用；超长文档行标注「摘录片段」）。产品代码本次零修改；本次仅产出本文件与 `r3-anchor-commands.txt`。

配套：R4 文档卫生锚点见同目录 `r4-doc-hygiene-anchors.md`（G-001/G-002/G-003 的文档侧已由该文件登记，本文件 §D 独立复核并补齐**代码侧**锚点）。

## A · 四类参照集的本仓现状锚点

基线 HEAD：`5230097570e4f4c1f245359f94c63e94e8c1a0bd`（`git log -1 --format=%H`；与 R2 基线 `ebe6013c` 之间 `apps/**` **零差异**，见 §E.1）。

### A.1 模块化单体 + 组合根（6 条代表锚点）

1. `apps/api/kernel/module.go:11` → `const KernelAPIVersion = "2.0.0"`（内核 API 版本常量）
2. `apps/api/kernel/module.go:265` → `func (r *Registry) Resolve(enabled []string) (Plan, error) {`（显式集 → 校验依赖/重复/环 → Plan）
3. `apps/api/kernel/provider.go:19` → `type Provider interface {`（Descriptor / CompiledPersistence / Register 三项）
4. `apps/api/kernel/provider.go:28` → `type Registrar interface {`（六面贡献：29 `HTTP(RouteContribution) error` … 34 `Configuration(ConfigurationContribution) error`；**刻意无 Persistence 方法**）
5. `apps/api/internal/composition/composition.go:121` → `func newAppWithOptions(cfg *config.Config, secretValue, seedHash string, logger *slog.Logger, extra ...fx.Option) (*fx.App, error) {`（Fx 组合根唯一入口；`composition.go:130` `fx.Provide(`、`composition.go:158` `return fx.New(opts...), nil`）
6. `apps/api/assembly/assembly.go:26` → `func OpenStore(ctx context.Context, dialect kernel.Dialect, path, dsn string, catalog []kernel.MigrationContribution) (kernel.Store, error) {`（B+ 公开装配面）

补充锚点（覆盖要求项「provider/registrar 契约的落地」与「静态汇集」）：`apps/api/kernel/provider.go:70` → `func RegisterContributions(ctx context.Context, plan Plan, providers []Provider) (ContributionSet, error) {`；`apps/api/internal/composition/composition.go:509` → `var providers []kernel.Provider`；`apps/api/internal/composition/composition.go:678` → `set, err := kernel.RegisterContributions(context.Background(), plan, providers)`；`apps/api/internal/composition/composition.go:80` → `func ResolvePlan(cfg *config.Config) (kernel.Plan, error) {`（先显式 `:87` `modules := append([]string(nil), cfg.ModulesEnabled...)`、后 Profile `:89` `resolved, err := kernel.ResolveProfile(cfg.ProfileName, nil)`）；`apps/api/kernel/profile.go:25` → `var profileDefaults = map[ProfileName][]string{`（mvp 11 / admin 23 / demo 12 项）。

### A.2 基础设施端口 + 内存/本地默认实现

| 端口 | 端口接口锚点（原文） | 默认实现锚点（原文） | 组合根接线锚点（原文） |
|---|---|---|---|
| Store | `apps/api/kernel/store.go:30` `type Store interface {` | `apps/api/internal/store/store.go:145` `func (s *Store) Run(ctx context.Context, fn func(kernel.Tx) error) error {`（SQLite）／`apps/api/internal/store/postgres.go:30` `func openPostgres(ctx context.Context, opts OpenOptions, catalog []kernel.MigrationContribution) (*postgres, error) {` | `apps/api/internal/store/open.go:18` `func Open(ctx context.Context, opts OpenOptions, catalog []kernel.MigrationContribution) (kernel.Store, error) {` |
| ObjectStore | `apps/api/kernel/objectstore.go:108` `type ObjectStore interface {` | `apps/api/internal/objectstore/local.go:41` `func NewLocal(root string) *LocalStore {`（本地盘默认）／`apps/api/internal/objectstore/s3.go:50` `func NewS3(endpoint, region, bucket, accessKeyID, secretAccessKey string, usePathStyle bool) (*S3Store, error) {` | `apps/api/internal/composition/composition.go:816` `func newObjectStore(cfg *config.Config) (kernel.ObjectStore, func(context.Context) error, error) {`；`:836` `return objectstore.NewLocal(root), nil, nil` |
| Cache | `apps/api/kernel/cache.go:138` `type Cache interface {` | `apps/api/internal/cache/memory.go:64` `func NewMemory(maxEntries int) (*Memory, error) {` | `apps/api/internal/composition/composition.go:854` `func newCache(cfg *config.Config) (kernel.Cache, error) {` |
| RateLimiter | `apps/api/kernel/ratelimit.go:36` `type RateLimiter interface {`；`:86` `type RateLimiterProvider interface {` | `apps/api/internal/ratelimit/memory.go:37` `func NewProvider() *Provider { return &Provider{} }` | `apps/api/internal/composition/composition.go:844` `func newRateLimiters() kernel.RateLimiterProvider {` |
| EventBus | `apps/api/kernel/eventbus.go:87` `type EventBus interface {` | `apps/api/internal/eventbus/memory.go:53` `func NewMemory(buffer int, logger *slog.Logger) *Memory {` | `apps/api/internal/composition/composition.go:871` `func newEventBus(cfg *config.Config, logger *slog.Logger) kernel.EventBus {` |
| Mail | `apps/api/kernel/mail.go:73` `type MailSender interface {` | `apps/api/internal/mail/runtime.go:137` `func NewSwitcher(store kernel.Store, masterKey []byte, seed SeedConfig, logger *slog.Logger) (*Switcher, error) {`；`:237` `return NewOutboxSink(s.store, cfg.MockRetention), nil`（mock 渠道） | `apps/api/internal/composition/composition.go:1005` `func newMailRuntime(cfg *config.Config, st kernel.Store, logger *slog.Logger) (*mail.Switcher, func(context.Context) error, error) {` |
| Telegram | `apps/api/kernel/telegram.go:103` `type TelegramSender interface {`；`:123` `type TelegramDispatcher interface {` | `apps/api/internal/channel/telegram/disabled.go:21` `return kernel.ErrTelegramDisabled`（disabled sender）；`:36` `return nil`（disabled 注册 no-op） | `apps/api/internal/composition/composition.go:909` `func newTelegramRuntimeForFx(params telegramRuntimeParams) (*TelegramRuntime, error) {`；`:926` `func buildTelegramRuntime(plan kernel.Plan, cfg *config.Config, st kernel.Store, rateLimiters kernel.RateLimiterProvider, options telegramRuntimeOptions) (*TelegramRuntime, error) {` |
| Lifecycle | `apps/api/kernel/lifecycle.go:14` `func NewRuntime(plan Plan) *Runtime {`；`:18` `func (r *Runtime) Start(ctx context.Context) error {` | 逆序停止：`apps/api/kernel/lifecycle.go:70` `for i := len(modules) - 1; i >= 0; i-- {` | `apps/api/internal/composition/composition.go:1073` `func registerLifecycle(lc fx.Lifecycle, …)`（长签名） |
| Observability | 无独立 kernel 端口（进程内 Observer） | `apps/api/internal/obs/observer.go:70` `func NewObserver(build BuildInfo) *Observer {` | `apps/api/internal/composition/composition.go:264` `func newTracing(cfg *config.Config) *obs.Tracing {`；`:291` `func newMetricsServer(cfg *config.Config, observer *obs.Observer, logger *slog.Logger) *obs.Server {` |

### A.3 Schema 驱动 Admin（模块贡献 → Manifest → sidebar 投影链）

同一条链的 6 个代表锚点：

1. `apps/api/kernel/contribution.go:38` → `type PageContribution struct {`（Schema 页面贡献；文档字节随贡献冻结）
2. `apps/api/kernel/contribution.go:76` → `type NavigationContribution struct {`（导航贡献；group 可选）
3. `apps/api/internal/composition/composition.go:678` → `set, err := kernel.RegisterContributions(context.Background(), plan, providers)`（六面贡献一次性校验后整体发布）
4. `apps/api/internal/composition/composition.go:746` → `data, err := manifest.ForModulesWithFragmentsAndGroups(`（manifest 聚合入口）
5. `apps/api/internal/manifest/manifest.go:213` → `func NormalizeSidebarGroups(data []byte, presentations []NavigationGrouping) ([]byte, error) {`（sidebar 分组归一化）
6. `apps/web/src/app/navigation.ts:224` → `export function projectNavigation(`（Web 只以 manifest 为输入做投影）

补充锚点（链上其余环节）：`apps/api/internal/composition/composition.go:729` `handler.RegisterSchemas(mux, a, set.Pages)`；`apps/api/internal/handler/schema.go:22` `func RegisterSchemas(mux routeRegistrar, a *auth.Authenticator, pages []kernel.PageContribution) {`；`apps/api/internal/composition/composition.go:763` `if err := handler.RegisterManifest(mux, data); err != nil {`；`apps/api/internal/handler/manifest.go:16` `func RegisterManifest(mux routeRegistrar, supplied ...[]byte) error {`；`apps/web/src/main.tsx:79` `manifestLoader: async () => loadAppManifestBytes(),`；`apps/web/src/protocol/app-manifest.ts:894` `export async function loadAppManifestBytes(options: {`；`apps/web/src/app/App.tsx:1049` `() => projectNavigation(manifest, path, navigationContext, t),`；`apps/web/src/protocol/load-page.ts:81` `export async function loadPageDocument(`（schemaUrl 解析在 `:90` `const url = resolveSchemaUrl(baseURL, page.schemaUrl, params);`）。

### A.4 同进程基座（装配与入口）

1. `apps/api/cmd/server/main.go:22` → `func main() {`（单进程入口）
2. `apps/api/cmd/server/main.go:50` → `app, err := composition.NewApp(cfg, secret, seedHash, logger)`
3. `apps/api/internal/composition/composition.go:1151` → `OnStop: func(ctx context.Context) error {`（停机全序起点；`:1152` `shutdownErr := srv.Shutdown(ctx)`、`:1164` `eventBusErr := eventBusPort.Stop(ctx)`、`:1170` `return errors.Join(shutdownErr, metricsErr, jobsErr, eventBusErr, runtimeErr, closeErr, tracingErr)`）
4. `apps/api/internal/jobs/runner.go:121` → `go r.loop()`（Job 运行时在 API 进程内）
5. `apps/api/server/serve.go:85` → `func Run(ctx context.Context, opts Options, signals <-chan os.Signal) (string, error) {`（第二条同进程装配/生命周期面；`serve.go:62` `func Serve(opts Options) error {`、`serve.go:215` `runtime := kernel.NewRuntime(plan)`）
6. `apps/api/internal/server/server.go:17` → `func New(cfg *config.Config, handler http.Handler, logger *slog.Logger) *http.Server {`（HTTP 超时/中间件面）

补充锚点：`apps/api/internal/composition/composition.go:87`（Profile/显式集解析）、`apps/api/cmd/schema-ui/main.go:76` `switch flag.Arg(0) {`（消费侧 CLI 子命令面：create / serve / add / upgrade / migrate-fork / config）；`compose.yaml:17` `# Non-goal: full production ops / CI-CD pipeline, TLS termination, multi-instance`（单进程基线）。

## B · gated/缺席能力的缺席证明

**口径**：缺席是**文档事实**，不是缺陷。以下均为「当前无该依赖 / 无该适配器 / 仅剩接缝声明」的可核对证据；分类仍由 R3 判据表决定，本文件不预判。

### B.1 Redis

| 证据 | 锚点与原文 | 结果 |
|---|---|---|
| 依赖 | `apps/api/go.mod:5` `require (` … `apps/api/go.mod:21` `modernc.org/sqlite v1.55.0` / `:22` `)`（全 require 块，直接 + 间接） | 无任何 redis 依赖；`git grep -i redis -- apps/api/go.mod` 0 命中，`go.sum` 0 命中 |
| 目录 | 无 `apps/api/**/*redis*` 目录（`Get-ChildItem -Recurse -Directory -Filter '*redis*'` 空） | 无适配器目录 |
| 代码内出现 | 仅注释级接缝：`apps/api/internal/composition/composition.go:426` `// dependency (Redis stays trigger-gated; the seam declaration lives in`；`apps/api/kernel/cache.go:9` `// view, and the active provider (in-memory default R2, Redis seam R3) resolves`；`apps/api/internal/ratelimit/memory.go:32` `// single instance (fx.Provide); a future Redis-tier provider implements the` | 无实现代码 |
| 配置键 | `git grep -i redis -- apps/api/internal/config apps/api/configs` **0 命中**；`apps/api/internal/config/config.default.yaml:215` `cache:` 段只有 `max_entries`（`:220` 注释 `# unparsable values fail closed at startup. Env override: CACHE_MAX_ENTRIES.`） | 无 Redis 配置键（连「承诺键」都未预制） |
| 接缝声明 | `docs/architecture/cache-redis-seam-and-track.md:59` `- **不引入 Redis 客户端依赖**（判据 #4 验证面：`go.mod` 无 redis；本波及触发前始终成立）。` | 文档侧与代码一致 |
| 路线图状态 | `docs/vision/roadmap.md:150` `| RT-Q03 | 缓存（Redis 等） | 无 | **trigger-gated** | …`（摘录片段）；`docs/vision/roadmap.md:152` RT-Q05 行「| RT-Q05 | 登录/API 限流跨实例 | 进程内滑动窗口 | **trigger-gated** |」（摘录片段） | 两行仍 gated，未被消耗 |

### B.2 外部消息 broker / 事务 outbox

| 证据 | 锚点与原文 | 结果 |
|---|---|---|
| 依赖 | `apps/api/go.mod:5`–`:22` require 块：`git grep -E 'kafka\|rabbit\|amqp\|nats\|nsq\|pubsub\|sarama' -- apps/api/go.mod` 0 命中 | 无 broker 客户端依赖 |
| 实现 | `apps/api/internal/eventbus/` 只有 `memory.go`（`apps/api/internal/eventbus/memory.go:53` `func NewMemory(buffer int, logger *slog.Logger) *Memory {`） | 仅进程内 channel 实现 |
| 端口语义 | `apps/api/kernel/eventbus.go:8` `// drain (VP-021). Public types carry neither provider handles nor broker`（摘录片段）；`:9` `// clients. Domain code consumes EventBus, never a channel implementation.` | 端口刻意不含 broker 类型 |
| 名称歧义（重要） | 代码中 `outbox` 全部指**站内出站邮件记录**，非事务 outbox：`apps/api/internal/mail/outbox.go`、`apps/api/assembly/assembly.go:43` `return mail.NewOutboxSink(st, retentionCap)`、`apps/api/internal/composition/composition.go:437` `handler.RegisterMailOutbox(mux, a, mail.NewOutboxSink(st, mail.DefaultOutboxCap))` | `outbox` 一词在本仓**不是**事件 outbox 的实现证据 |
| 路线图状态 | `docs/vision/roadmap.md:149` RT-Q02 行「| RT-Q02 | 外部消息队列 / Job broker | 无 | **trigger-gated** |」（摘录片段）；`docs/vision/roadmap.md:153` `| RT-Q06 | 事务 outbox / inbox | 无 | **trigger-gated** | 先 DB outbox，后可选 broker |` | 两行仍 gated |

### B.3 搜索引擎

| 证据 | 锚点与原文 | 结果 |
|---|---|---|
| 依赖 | `apps/api/go.mod`：`git grep -E 'elastic\|meilisearch\|typesense\|bleve\|opensearch\|solr\|tantivy' -- apps/api/go.mod` 0 命中 | 无搜索引擎依赖 |
| 端口/适配器 | `apps/api/kernel/` 无 search 端口文件（目录列表见 `r3-anchor-commands.txt`）；`apps/api/internal/` 无 search 目录 | 无端口、无适配器 |
| 路线图 | `docs/vision/roadmap.md:306` `**刻意后置**：MongoDB、ORM、Redis、消息队列、搜索引擎、K8s、SMS。它们是部署或产品触发的后果，或已否决的技术选型。` | 明示后置 |

### B.4 多实例 / K8s

| 证据 | 锚点与原文 | 结果 |
|---|---|---|
| K8s 依赖/清单 | `apps/api/go.mod`：`git grep -E 'k8s.io\|client-go\|operator-framework'` 0 命中；仓库内无 `helm` / `k8s` / `kubernetes` / `charts` / `manifests` 目录（递归目录扫描为空） | 无 K8s 交付物 |
| 编排形态 | `compose.yaml:17` `# Non-goal: full production ops / CI-CD pipeline, TLS termination, multi-instance`；`:18` `# scaling (see I-008-001 §8).`；`:67` `    stop_grace_period: 15s`；`:68` `    restart: on-failure`（单容器，无 `deploy.replicas`） | 单实例 Compose |
| 进程内假设 | `apps/api/internal/ratelimit/memory.go:32`（进程内单实例注释，摘录片段）；`apps/api/internal/composition/composition.go:435` `_ = eventBusPort // accepted, not yet consumed (intentional until a consumer lands)` | 单进程假设显式 |
| 路线图 | `docs/vision/roadmap.md:211` `| RT-D04 | 多实例 / 水平扩展 | Compose 非目标 | **trigger-gated** | 拉动 RT-P02/P04、RT-Q\*、RT-S02、RT-Q05 |` | 仍 gated |

### B.5 ORM

| 证据 | 锚点与原文 | 结果 |
|---|---|---|
| 依赖 | `apps/api/go.mod`：`git grep -E 'gorm\|ent.io\|jmoiron\|uptrace/bun\|sqlboiler\|xorm'` 0 命中 | 无 ORM 依赖 |
| 数据访问面 | 模块 SQL 直接走内核 `Tx`：`apps/api/kernel/store.go:43` `type Tx interface {`；`:44` `Exec(ctx context.Context, query string, args ...any) (Result, error)`；占位符重绑在 `apps/api/internal/store/rebind.go:13` `func rebindPostgres(query string) string {` | 手写 SQL + 双方言端口 |
| 路线图 | `docs/vision/roadmap.md:137` RT-P03 行「| RT-P03 | Store 双方言端口（无 ORM） | 内核持久化端口 + SQLite/PG 双方言实现（VP-013 交付） | **delivered** |」（摘录片段）；`docs/vision/roadmap.md:306`（同 B.3） | 「无 ORM」是已交付的**设计选择**，不是缺口 |

## C · 12 条入册 residual 的现状核对

行号基线与 §A 同（HEAD `5230097570e4f4c1f245359f94c63e94e8c1a0bd`）。「一致」= 现行代码/文档与 R1 冻结表 §2 的记载**在事实层面**吻合；「不一致」= 记载已被现行代码推翻或范围需收窄。

| residual id | 现状（代码/文档事实） | 精确锚点 | 与 R1/R2 记载是否一致 | 备注 |
|---|---|---|---|---|
| RES-013-migrator | 仍无产品 SQLite→PG **数据**搬运器：双方言只共享**同一** compiled 迁移目录 + 各自应用器（schema 层），数据层仅备份/恢复（SQLite `VACUUM INTO`、PG `pg_dump`/`pg_restore`）。CLI 子命令面只有 `create/serve/add/upgrade/migrate-fork/config`，无 data-migrate。 | `apps/api/internal/store/open.go:18` `func Open(ctx context.Context, opts OpenOptions, catalog []kernel.MigrationContribution) (kernel.Store, error) {`；`apps/api/internal/store/migrate.go:294` `if _, err := s.db.Exec("VACUUM INTO '" + strings.ReplaceAll(target, "'", "''") + "'"); err != nil {`；`apps/api/cmd/schema-ui/main.go:76` `switch flag.Arg(0) {`；文档侧 `docs/vision/charter.md:73`（长行，摘录「VP-013 … residual = 无产品 SQLite→PG 搬运器」） | **一致** | R2 矩阵未单列本行（不在 17 面分母内），本轮为新增核对 |
| RES-014-migrator | 仍无产品本地盘→对象存储搬运器：`internal/objectstore/` 生产文件仅 `local.go` / `s3.go`，无 copy/transfer 函数；组合根按 `driver` 二选一构造。 | `apps/api/internal/objectstore/local.go:41` `func NewLocal(root string) *LocalStore {`；`apps/api/internal/objectstore/s3.go:50` `func NewS3(endpoint, region, bucket, accessKeyID, secretAccessKey string, usePathStyle bool) (*S3Store, error) {`；`apps/api/internal/composition/composition.go:816`（二选一入口）；文档侧 `docs/vision/charter.md:73`（摘录「VP-014 … residual = 无产品本地盘→对象存储搬运器」） | **一致** | 同上 |
| RES-015-otlp-sink | in-repo `cmd/otlp-sink` 明确不解析：只 `io.Copy(io.Discard, r.Body)` + 计数/打印，文件头自述「It is NOT a collector: no parsing」。 | `apps/api/cmd/otlp-sink/main.go:4` `// bytes and prints one line per request. It is NOT a collector: no parsing,`；`:32` `n, err := io.Copy(io.Discard, r.Body)`；`:5` `// no semantics, no storage. Do not use in production.` | **一致** | R2 行 10 亦按「显式 endpoint 才导出」计 |
| RES-015-metrics | 指标分母仍只有 4 个 family：build_info / http_requests_total / http_request_duration_seconds / kernel_modules_enabled；**无** Store / 对象 / Job 指标。 | `apps/api/internal/obs/observer.go:74` `Name: "suc_build_info",`；`:87` `Name: "suc_http_requests_total",`；`:94` `Name:    "suc_http_request_duration_seconds",`；`:102` `Name: "suc_kernel_modules_enabled",` | **一致** | R2 行 10 的「仅当前指标分母」可复现 |
| RES-016-revoke | 轮换重叠窗内 **previous 可验**（= 不选「立即失效」）：`verifyAccess` 当前钥失败且存在 previous 时回退 previous。另一条**不同轴**的机制：access token 携带 `token_version`，改密/禁用后中间件立即拒绝旧 token（W4 P0-3）——不改变本 residual 的指向。 | `apps/api/internal/auth/auth.go:431` `func (a *Authenticator) verifyAccess(raw string) (ParsedAccessToken, error) {`；`:433` `if err == nil || len(a.previousSecret) == 0 {`；`:436` `return ParseAccessToken(a.previousSecret, raw)`；`:630` `if parsed.TokenVersion != u.TokenVersion {`；文档侧 `docs/vision/plans/VP-016-key-rotation-and-backup.md:103`（I-016-005 行，摘录「默认「previous 可验」已随 VRev-035 冻结进退出判据 1 并按此交付」） | **一致** | 备注：本 residual 只覆盖**密钥轮换重叠窗**；不得把 `token_version` 的批量失效读成「立即失效已交付」。R1 归类「接受残余」的指向未变 |
| RES-016-mfa-wrap | 记载的「不随 JWT previous 重包」**已被现行代码推翻**：MFA 服务在配置 previous 时派生 `prevKey`，密文在 previous 钥下可解，且**成功 TOTP 校验后**由 CAS 赢家惰性重包为当前钥。仍存的范围：恢复码路径不重包、无启动期批量重包、无主动轮换重包。 | `apps/api/modules/mfa/service.go:57` `// prevKey is the HKDF key derived from the previous JWT secret during a`；`:60` `prevKey []byte`；`:72` `func NewService(repo *store.Repository, serverSecret, previousSecret []byte) *Service {`；`:165` `s.maybeRewrap(st, plain, fromPrevious, now)`；`:450` `func (s *Service) decryptSecret(encoded string) (plain string, fromPrevious bool) {`；`:466` `func (s *Service) maybeRewrap(state *store.State, plain string, fromPrevious bool, now time.Time) {`；接线 `apps/api/internal/composition/composition.go:395` `mfaService = mfamodule.NewService(mfastore.NewRepository(st), []byte(secret), []byte(cfg.AuthJWTSecretPrevious))` | **不一致**（范围需收窄，非「已修复=闭合」） | R1 引用的原文来自 VP-016 关门时点（`docs/vision/plans/VP-016-key-rotation-and-backup.md:115`：摘录「本 VP 未改 `internal/modules/mfa`」）——其后 W11 F-004 改了。R3 若据此收窄表述，属**分类输入**；若扩大或新增残余须按 D-001 裁决 B 走用户书面接受 |
| RES-021-harness | 进程级停机 harness 仍带 `!windows` 构建约束；Linux CI（`ubuntu-latest`）跑 `go test ./...`，即以 linux 侧核销；Compose 侧有 15s `stop_grace_period`。 | `apps/api/cmd/server/shutdown_harness_test.go:1` `//go:build !windows`；`.github/workflows/r6-basic-matrix.yml:28` `name: api (Linux, Go 1.26)`；`:29` `runs-on: ubuntu-latest`；`:74` `- run: go test ./...`；`compose.yaml:67` `stop_grace_period: 15s` | **一致** | 复审触发仍 = CI 失败（R1 原文） |
| RES-026-redis | 仍 trigger-gated：`go.mod`/`go.sum` 无 redis、无适配器目录、无配置键；仅注释级接缝 + owner 文档 §2.5/§2.6.5 的「不引入客户端依赖」。 | `apps/api/go.mod:5`–`:22`（require 块，0 命中）；`docs/architecture/cache-redis-seam-and-track.md:87` `- **不引入 Redis 客户端依赖**（判据 #4 验证面：`go.mod` 无 redis；触发前始终成立）。`；`docs/vision/roadmap.md:150`（RT-Q03 行）；`docs/vision/roadmap.md:152`（RT-Q05 行） | **一致** | 与 §B.1 同证 |
| RES-028-broker | 仍 gated：无 broker 依赖；EventBus 只有 `internal/eventbus/memory.go`；outbox/MQ 仅接缝声明。 | `apps/api/internal/eventbus/memory.go:53` `func NewMemory(buffer int, logger *slog.Logger) *Memory {`；`docs/vision/roadmap.md:149`（RT-Q02 行）；`docs/vision/roadmap.md:153` `| RT-Q06 | 事务 outbox / inbox | 无 | **trigger-gated** | 先 DB outbox，后可选 broker |` | **一致** | 注意 `mail.OutboxSink` 不是事件 outbox（§B.2） |
| RES-030-keyfile | bot master key 默认文件与 DB **同目录**（`telegram-master.key`），邮件 master key 同形（`mail-master.key`）；两者均可由配置改路径。 | `apps/api/internal/composition/composition.go:930` `masterKeyPath = filepath.Join(filepath.Dir(cfg.DBPath), "telegram-master.key")`；`apps/api/internal/composition/composition.go:1016` `masterKeyPath = filepath.Join(filepath.Dir(cfg.DBPath), "mail-master.key")` | **一致** | R1 记「已 accepted-residual；复审 = KMS 波」——现状未变，无需重开 |
| RES-T03-tz ⏱️**时敏** | DB `timestamptz` 未做：`apps/api` 全仓 `timestamptz` **0 命中**；时间列仍为 SQLite 兼容的 `INTEGER`（epoch）形状；roadmap RT-T03 保持 `registered`。 | 代码侧反向证据 `git grep -i timestamptz -- apps/api` = 0 命中；时间列形状 `apps/api/modules/authsession/migration/migration.go:41` `created_at INTEGER NOT NULL`；文档侧 `docs/vision/roadmap.md:252` `| RT-T03 | 时区在持久化层的合同 | locale 在 Admin 功能面；DB `timestamptz` 未做 | **registered** |`（摘录片段）；`docs/vision/plans/VP-020-timezone-number-currency-formatting.md:25`（摘录「DB `timestamptz` 持久化时区合同仍归架构分支 RT-T03（`registered`，不进本波）」） | **一致** | **时敏（R1/R2 明确点名）**：R3 须裁定「是否建议下一拍」。本轮未发现任何实现动作或迁移草案；RT-T03 的状态值仍是 `registered`（未被消耗、也未被降级） |
| RES-P04-pool ⏱️**时敏** | **文档锚点过期、代码已池化**：代码文件库连接池默认 `4`（`sqlitePoolDefault`），内存库固定 `1`，可经 `PoolMaxOpenConns` 覆盖；文档现状锚点仍写 `MaxOpenConns=1`。RT-P04 行本身（连接池/读写分离/replica）仍 trigger-gated——池化已交付不等于读写分离/replica 已交付。 | 代码侧 `apps/api/internal/store/store.go:29` `const sqlitePoolDefault = 4`；`:107` `db.SetMaxOpenConns(1)`（memory 分支）；`:113` `db.SetMaxOpenConns(pool)`（文件库分支）；文档侧 `docs/vision/roadmap.md:104` `> 现状锚点：单进程 + SQLite（`MaxOpenConns=1`）+ 本地盘上传 + 进程内 Job + 内存限流。…`（摘录片段）；`docs/vision/roadmap.md:138` `| RT-P04 | 连接池 / 读写分离 / replica | `MaxOpenConns=1` | **trigger-gated** | 多实例或 PG 之后才有意义 |` | **一致**（R1 要求核实「文档锚点是否过期」→ 结论：**已过期**；gated 判定不变） | **时敏（R1/R2 明确点名）**：与 R2 G-002、同目录 `r4-doc-hygiene-anchors.md` §G-002 第 1/2 行结论一致；修文档归 R4（判据 4），不得把「已池化」扩写成读写分离/replica 已交付 |

**计数**：一致 **11** / 不一致 **1**（RES-016-mfa-wrap）/ 无法核对 **0**。

## D · G-001～G-004 的复核

以下均为独立复核：文档侧与代码侧各自给出**精确锚点**。**不给出分类建议**（分类由 R3 判据表 + 用户裁决定）。

### D-1 · G-001（`overview.md` 现时节过期）

| 侧 | 锚点与原文 | 复核结论 |
|---|---|---|
| 文档 | `docs/architecture/overview.md:71` `- **愿景**：[charter.md](../vision/charter.md) **`schema-ui-core-admin-foundation@0.2.0`**，且当前仅有一个 active Charter；未来方向已从双线长期维护改为单主线模块化。` | 与现行 Charter 版本不一致（G-001 主张成立） |
| 文档 | `docs/architecture/overview.md:64` `| `web/` | FastAPI Web 应用（有界受控写入，默认门闩关闭） |` | 与产品树（`apps/api` + `apps/web`）不符（G-001 主张成立） |
| 文档 | `docs/architecture/overview.md:73` `- **当前交付 VP**：[VP-005](../vision/plans/VP-005-design-system-and-ui-experience.md)（设计系统与 Schema 驱动 UI/UX）**`active`**（v0.4.1）；lead = `workspace-006-design-system-and-ui-experience`。…`（摘录片段） | 「当前交付 VP」停在 VP-005（G-001 主张成立） |
| 文档 | `docs/architecture/overview.md:74`–`:79`（工作区清单只到 `workspace-006`；`:79` `- **workspace-006**：状态 **active**，角色 **delivery**，Root `GOAL-001-design-system-and-ui-experience` **active / 0/5**，…`） | 与实际 `workspace-001`～`035` 不符（G-001 主张成立） |
| 代码侧对照 | 现行 Charter 版本见 `docs/architecture/cache-redis-seam-and-track.md:9` `vision_ref: schema-ui-core-admin-foundation@0.4.0`；产品树见 `apps/api/go.mod:1` `module github.com/magicvr/schema-ui-core/apps/api` 与 `apps/web/package.json` | 代码侧无「过期」问题——G-001 是**纯文档**项 |

### D-2 · G-002（roadmap 现状锚点 / RT-P04 / RT-D02）

| 侧 | 锚点与原文 | 复核结论 |
|---|---|---|
| 文档 | `docs/vision/roadmap.md:104` `> 现状锚点：单进程 + SQLite（`MaxOpenConns=1`）+ 本地盘上传 + 进程内 Job + 内存限流。Compose 已声明非目标含 TLS 终止与多实例（`compose.yaml`）。` | 现状锚点过期（见 §C RES-P04-pool） |
| 文档 | `docs/vision/roadmap.md:138` RT-P04 行（同 §C 引文） | 现状锚点过期；**状态值 `trigger-gated` 本身仍成立** |
| 文档 | `docs/vision/roadmap.md:209` `| RT-D02 | 优雅停机 / 连接排空 | 进程生命周期有，无明确 drain 合同 | **delivered**（VP-021 `closed` v0.3.0 · 2026-08-27） | 2026-08-26 立项 → 2026-08-27 交付：停机顺序 / HTTP drain / Job 语义（中断标记重跑）/ 双方言 Store 排空；单进程基线；与 Job 租约相关部分仍随 A3 |` | 同一行「无明确 drain 合同」与 `delivered` 并存，属自相矛盾表述（G-002 主张成立） |
| 代码 | `apps/api/internal/composition/composition.go:1151` `OnStop: func(ctx context.Context) error {`；`:1152` `shutdownErr := srv.Shutdown(ctx)`；`:1170` `return errors.Join(shutdownErr, metricsErr, jobsErr, eventBusErr, runtimeErr, closeErr, tracingErr)` | drain 合同**已实现**（HTTP drain → retention/metrics → jobs → eventbus → runtime → store → tracing） |
| 代码 | `apps/api/internal/store/store.go:29` `const sqlitePoolDefault = 4`；`:107` `db.SetMaxOpenConns(1)` | 池化已交付；内存库仍 1 |

### D-3 · G-003（cache-redis-seam 文档落后于 RateLimiter API）

**文档仍描述的 API 名**（`docs/architecture/cache-redis-seam-and-track.md`）：

| 锚点与原文 | 文档提到的 API |
|---|---|
| `:68` `- 端口语义保持：`Allow` **不注册**（只读检查）、失败才 `Record`、`Clear` 清桶、`RetryAfterSeconds` 语义分母 = `kernel.RateLimiterRetryAfterSeconds`（触发立项时按远端 TTL 细化）。` | `Allow` / `Record` / `Clear` / `RetryAfterSeconds`（+ 文末 `RateLimiterInWindow`） |
| `:73` `- **原子窗口原语（VP-027 冻结）**：`INCR` 计数 + `EXPIRE`（窗口）为最小原子原语——`Record` = `INCR` + 首次 `EXPIRE`；`Allow` = 读计数（`GET`，不写、不续期失败桶）；`Clear` = `DEL`。` | 同上四名；**未提** `AllowRecord` / `Reserve` / `Cancel` |
| `:66` `- Redis 级限流供应商实现**同一** `kernel.RateLimiter` / `kernel.RateLimiterProvider` 接口（workspace-027 GOAL-002 D-002 v0.1.1）；7 处使用点消费方零感知、零代码改动。` | 端口分母只点名 kernel 两接口，未列原子方法 |

**端口现行声明的 API 名**（`apps/api/kernel/ratelimit.go`）：

| 锚点与原文 | 端口 API |
|---|---|
| `:40` `Allow(key string, now time.Time) bool` | `Allow`（旧三元组，保留兼容） |
| `:44` `Record(key string, now time.Time)` | `Record` |
| `:52` `AllowRecord(key string, now time.Time) bool` | **`AllowRecord`（文档缺）** |
| `:65` `Reserve(key string, now time.Time) (token uint64, ok bool)` | **`Reserve`（文档缺）** |
| `:71` `Cancel(key string, token uint64)` | **`Cancel`（文档缺）** |
| `:76` `RetryAfterSeconds(key string, now time.Time) int` | `RetryAfterSeconds` |
| `:79` `Clear(key string)` | `Clear` |

**差值结论**：文档 `:68`/`:73` 描述的端口面 = {`Allow`, `Record`, `Clear`, `RetryAfterSeconds`}；端口现行面 = 上述 7 个方法，**多出 `AllowRecord` / `Reserve` / `Cancel`**（VP-032 原子化；该 VP 的路线图行在 `docs/vision/roadmap.md:51`，行首为 `| 32 | [VP-032-rate-limiter-atomic-port](plans/VP-032-rate-limiter-atomic-port.md) |`）。**旁证**：`git grep -E 'AllowRecord|Reserve|Cancel' -- docs/architecture/cache-redis-seam-and-track.md` = **0 命中**（该文档全文从未出现这三个名字）；三个方法名在文档中的确只有 `:68`/`:73` 描述的旧面。G-003 主张成立，且落在**文档**侧（§2.6.1/§2.6.2 两处）。内存供应商对应实现：`apps/api/internal/ratelimit/memory.go:158` `func (l *Memory) AllowRecord(key string, now time.Time) bool {`、`:174` `func (l *Memory) Reserve(key string, now time.Time) (uint64, bool) {`、`:187` `func (l *Memory) Cancel(key string, token uint64) {`。

### D-4 · G-004（`assembly.NewAuthenticator` 返回 internal 指针且仅 current 单密钥）

**assembly 侧（公共工厂）**：

| 锚点与原文 |
|---|
| `apps/api/assembly/assembly.go:5` `// ——Go 类型推断对 *auth.Authenticator、kernel.Store 等返回值直接消费。` |
| `apps/api/assembly/assembly.go:34` `// NewAuthenticator 构造 JWT 认证器（current 单密钥形态）。` |
| `apps/api/assembly/assembly.go:37` `func NewAuthenticator(secret []byte, accessTTL, refreshTTL time.Duration, runner authsession.TxRunner) *auth.Authenticator {` |
| `apps/api/assembly/assembly.go:38` `return auth.New(secret, accessTTL, refreshTTL, runner, false)` |

**auth 侧（单密钥形态的证据）**：

| 锚点与原文 |
|---|
| `apps/api/internal/auth/auth.go:147` `// New builds an Authenticator. The caller is responsible for a non-empty secret`（`:151` `func New(secret []byte, accessTTL, refreshTTL time.Duration, runner authsession.TxRunner, devSession bool) *Authenticator {`） |
| `apps/api/internal/auth/auth.go:152` `return NewWithRepository(secret, accessTTL, refreshTTL, authsession.NewRepository(runner), devSession)` |
| `apps/api/internal/auth/auth.go:158` `return &Authenticator{secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL, repository: repository, devSession: devSession}`（**`previousSecret` 缺省零值 → 单密钥**） |
| `apps/api/internal/auth/auth.go:131` `previousSecret []byte`（字段声明；唯一置值路径为 `:167`–`:168` `func NewWithRepositoryAndPrevious(current, previous []byte, …` / `return &Authenticator{secret: current, previousSecret: previous, …}`） |
| `apps/api/internal/auth/auth.go:433` `if err == nil || len(a.previousSecret) == 0 {`（`len == 0` 即严格单密钥：`verifyAccess` 不回退） |

**复核结论**：G-004 的两项主张——① 公共工厂返回 `internal` 指针（`*auth.Authenticator`）；② 只注入 current、`previousSecret` 恒为空——在 HEAD 上**逐字成立**。补充事实：仓内并非没有「current + previous」的接线，只是走内部构造器 `apps/api/internal/composition/composition.go:259` `return auth.NewWithRepositoryAndPrevious([]byte(secret), []byte(cfg.AuthJWTSecretPrevious), cfg.AuthAccessTTL, cfg.AuthRefreshTTL, repository, cfg.AuthDevSessionEnabled)`——**B+ 公共面与内部接线不同形**，这是 G-004 的准确表述。

## E · 证据矛盾与不确定

### E.1 基线可比性（先行事实）

- `git diff --name-only ebe6013c..HEAD -- apps` → **空**（0 个路径）。R2 基线 `ebe6013c` 与本次 HEAD `5230097570e4f4c1f245359f94c63e94e8c1a0bd` 之间的 29 个改动路径**全部**在 `docs/workspaces/workspace-035-*` 内。
- 推论：R2 矩阵的**每一条代码锚点在 HEAD 上都应逐字复现**。§E.3 列出的错行因此是**写作/计数误差**，不是代码漂移；不得以「代码变了」解释。

### E.2 与 R2 矩阵结论矛盾或需限定的事实

1. **R2 行 10 `internal/config/config.go:493` 指向的语句与主张无关**（详见 §E.3 第 14 项）：`:493` 是 `CacheMaxEntries: DefaultCacheMaxEntries,`，而「metrics/traces 默认关闭」的语句是 `:495` `MetricsEnabled:           false,` / `:498` `TracesEnabled:            false,`。这是本轮发现的**唯一一处「锚点指向另一主题的语句」**（其余错行均指向紧邻的注释/空行/包裹函数）。
2. **R2 行 09 的「CaptureSink 是历史实现，不能当现行默认」——主张正确，但产品代码内注释仍与现行接线矛盾**：现行 mock 渠道由 `apps/api/internal/mail/runtime.go:237` `return NewOutboxSink(s.store, cfg.MockRetention), nil` 构成，而 `apps/api/internal/mail/capture.go:12` 仍自述 `// CaptureSink is the embedded default kernel.MailSender for unconfigured SMTP`。即：**接线已换到 OutboxSink，`capture.go` 的文档注释未同步**。本文件不改代码，仅登记（R3 可判为文档卫生输入；不在 R2 矩阵主张内）。
3. **「构造内存 provider」类锚点中，R2 行 06 的 `composition.go:853` 是注释行**，`func newCache` 在 `:854`；同类的 R2 行 08 `composition.go:871`（`func newEventBus`）**本身正确**。同一叙述族里正确与错行并存，故 §E.3 逐条列出而非整行否定。
4. **无法核对项（0 条硬缺失，2 条口径限制）**：
   - R2 行 05 的「本地与 fake S3 证据不等于 live S3 验收」、行 11 的「进程/Compose Linux 证据本轮未重跑」属**验证口径**声明，本轮为只读取证，**未重跑**任何 live/Compose 证据，故无法对这两句的强度做增量核对（不构成矛盾）。
   - R2 行 16 的「module-architecture §1–7 主线与当前实现一致」是**全局一致性主张**，本轮只抽检了其 frontmatter（`docs/architecture/module-architecture.md:9` `vision_ref: schema-ui-core-admin-foundation@0.2.0`，与 A-002「旧架构记录」判断一致），未逐节复核 §1–7 全文。
5. **本次未能确定的问题**：
   - RES-016-mfa-wrap 的**正确残余措辞**（应写「恢复码路径不重包 / 无启动批量重包 / 无主动轮换重包」中的哪一条或哪几条）——需 owner 或用户按 D-001 裁决 B 定稿；本文件只给可核对事实。
   - `docs/architecture/cache-redis-seam-and-track.md:66` 的「**7 处使用点**」是否仍等于现行调用点数量——本轮未逐一清点使用点（不在四类参照集取证要求内），仅登记该数字为**未复核**。
   - `.github/workflows/r6-basic-matrix.yml` 是否在最近一次 main 推送中实际运行成功（本轮只读工作树，未查 GitHub Actions 运行历史）。
6. **未触碰项**：未读 `apps/api/configs/.env`；未运行完整测试套件；未改任何 `apps/**`、`docs/vision/**`、`docs/architecture/**` 或治理文件；未执行 `git add/commit/checkout/stash`。

### E.3 R2 矩阵锚点错行清单（**不静默修改矩阵**，仅列出 旧锚点 → 精确锚点）

判定口径：**旧锚点行不是其所声称的语句/声明**。凡只是「指向包裹该行为之函数声明」的，另列 §E.4（可再精确，非错行）。

| # | R2 矩阵位置 | 旧锚点 | 旧锚点行原文 | 精确锚点 | 精确锚点行原文 |
|---|---|---|---|---|---|
| 1 | 行 01 | `apps/api/kernel/provider.go:27` | `// into the compiled-global catalog (freeze package §4).` | `apps/api/kernel/provider.go:28` | `type Registrar interface {` |
| 2 | 行 04 | `apps/api/kernel/store.go:24` | `// errors.Is(err, kernel.ErrNoRows).` | `apps/api/kernel/store.go:25` | `var ErrNoRows = sql.ErrNoRows` |
| 3 | 行 04 | `apps/api/internal/store/store.go:141` | `// Nested Run is forbidden and detected per-callback via a goroutine-local` | `apps/api/internal/store/store.go:145` | `func (s *Store) Run(ctx context.Context, fn func(kernel.Tx) error) error {` |
| 4 | 行 04 | `apps/api/internal/store/store.go:104`（作「内存池 1」的锚点） | `if memory {` | `apps/api/internal/store/store.go:107` | `db.SetMaxOpenConns(1)` |
| 5 | 行 05 | `apps/api/internal/objectstore/s3.go:47`（作「供应商适配」的锚点） | `// Credentials are static-only (Root D-005): the SDK default chain (~/.aws,` | `apps/api/internal/objectstore/s3.go:50` | `func NewS3(endpoint, region, bucket, accessKeyID, secretAccessKey string, usePathStyle bool) (*S3Store, error) {` |
| 6 | 行 06 | `apps/api/kernel/cache.go:117`（作「key/expiry 约束」的锚点） | `return nil` | `apps/api/kernel/cache.go:107`（key/value/policy 顺序校验） | `func ValidateCacheSet(key string, value []byte, policy ExpiryPolicy) error {` |
| 6b | 行 06（同上，expiry 契约本体） | 同上 | 同上 | `apps/api/kernel/cache.go:124` | `type ExpiryPolicy interface {` |
| 7 | 行 06 | `apps/api/internal/cache/memory.go:227`（作「拷贝保持空 slice hit」的锚点） | `// copyBytes copies src preserving empty-value semantics: an empty-but-non-nil` | `apps/api/internal/cache/memory.go:232` | `func copyBytes(src []byte) []byte {` |
| 8 | 行 07 | `apps/api/kernel/ratelimit.go:35` | `// for deterministic tests; production callers pass time.Now().UTC().` | `apps/api/kernel/ratelimit.go:36` | `type RateLimiter interface {` |
| 9 | 行 07 | `apps/api/internal/ratelimit/memory.go:83`（作「Cancel 仅取消本次 token」的锚点） | `func (l *Memory) Allow(key string, now time.Time) bool {` | `apps/api/internal/ratelimit/memory.go:187` | `func (l *Memory) Cancel(key string, token uint64) {` |
| 10 | 行 08 | `apps/api/kernel/eventbus.go:43` | `// DefaultEventBusBuffer is the fallback per-subscription buffer size` | `apps/api/kernel/eventbus.go:87` | `type EventBus interface {` |
| 11 | 行 09 | `apps/api/kernel/mail.go:33` | `// recipient. There is deliberately no From field: the configured default` | `apps/api/kernel/mail.go:35` | `type MailMessage struct {` |
| 12 | 行 09 | `apps/api/internal/mail/runtime.go:105`（作「渠道切换」的锚点） | `var ErrUnknownChannel = errors.New("mail: unknown channel")` | `apps/api/internal/mail/runtime.go:316` | `func (s *Switcher) Update(ctx context.Context, req UpdateRequest) (*PublicView, error) {` |
| 13 | 行 09 | `apps/api/internal/mail/runtime.go:363`（作「失败保留旧渠道」的锚点） | `if err != nil {` | `apps/api/internal/mail/runtime.go:370` | `return nil, fmt.Errorf("mail: new SMTP endpoint is unreachable, keeping the previous channel: %w", err)` |
| **14** | **行 10** | **`apps/api/internal/config/config.go:493`** | **`CacheMaxEntries:          DefaultCacheMaxEntries,`** | **`apps/api/internal/config/config.go:495`（metrics）/ `:498`（traces）** | **`MetricsEnabled:           false,` / `TracesEnabled:            false,`** |
| 15 | 行 10 | `apps/api/internal/obs/server.go:21` | `// Server is the optional dedicated metrics listener (R1 D-001 s1, R2 D-001` | `apps/api/internal/obs/server.go:26` | `type Server struct {` |
| 16 | 行 11 | `apps/api/kernel/lifecycle.go:67`（作「逆序停止」的锚点） | （空行） | `apps/api/kernel/lifecycle.go:70` | `for i := len(modules) - 1; i >= 0; i-- {` |
| 17 | 行 13 | `apps/api/kernel/telegram.go:45`（作「消息/分发契约」的锚点） | `var (` | `apps/api/kernel/telegram.go:103` / `:123` | `type TelegramSender interface {` / `type TelegramDispatcher interface {` |
| 17b | 行 13（消息契约本体） | `apps/api/kernel/telegram.go:76` | `func (m TelegramMessage) Validate() error {` | `apps/api/kernel/telegram.go:65` | `type TelegramMessage struct {` |
| 18 | 行 13 | `apps/api/internal/channel/telegram/disabled.go:9` | `// DisabledSender is the fail-closed kernel.TelegramSender provided when channel.telegram is disabled.` | `apps/api/internal/channel/telegram/disabled.go:21`（sender 返回 disabled）/ `:36`（注册 no-op） | `return kernel.ErrTelegramDisabled` / `return nil` |
| 19 | 行 13 | `apps/api/internal/channel/telegram/runtime.go:90`（作「参数校验」的锚点） | `// "主密钥离开源码"). A nil/empty key is a construction error (fail-closed).` | `apps/api/internal/channel/telegram/runtime.go:99` | `return nil, fmt.Errorf("telegram: master key is required")` |
| 20 | 行 14 | `apps/api/internal/auth/auth.go:424`（作「current 失败再试 previous」的锚点） | `// verifyAccess verifies an access token against the current signing key and,` | `apps/api/internal/auth/auth.go:431`；回退语句 `:433` / `:436` | `func (a *Authenticator) verifyAccess(raw string) (ParsedAccessToken, error) {` / `if err == nil || len(a.previousSecret) == 0 {` / `return ParseAccessToken(a.previousSecret, raw)` |
| 21 | 行 15 | `apps/api/internal/manifest/manifest.go:40`（作「结构化贡献投影」的锚点） | `// presentation declaration that can affect the public manifest. Entries with a` | `apps/api/internal/manifest/manifest.go:43` | `func NavigationPresentationsFromContributions(contributions []kernel.NavigationContribution) []NavigationGrouping {` |
| 22 | 行 02 | `apps/api/internal/composition/composition.go:83`（作「先显式再 Profile」的锚点） | `}` | `apps/api/internal/composition/composition.go:87`–`:89` | `modules := append([]string(nil), cfg.ModulesEnabled...)` / `if len(modules) == 0 {` / `resolved, err := kernel.ResolveProfile(cfg.ProfileName, nil)` |
| 23 | 行 03 | `apps/api/internal/composition/composition.go:127`（作「Fx 在组合根」的锚点） | `logger = slog.Default()` | `apps/api/internal/composition/composition.go:130` / `:158` | `fx.Provide(` / `return fx.New(opts...), nil` |
| 24 | 行 14 | `apps/api/internal/composition/composition.go:255`（作「注入 current + previous」的锚点） | `// VP-016 R2 (workspace-016 GOAL-003 D-001): an empty previous keeps exact` | `apps/api/internal/composition/composition.go:259` | `return auth.NewWithRepositoryAndPrevious([]byte(secret), []byte(cfg.AuthJWTSecretPrevious), cfg.AuthAccessTTL, cfg.AuthRefreshTTL, repository, cfg.AuthDevSessionEnabled)` |
| 25 | 行 10 | `apps/api/internal/composition/composition.go:262` / `:289`（作「接线」的锚点） | `// newTracing maps the observability.traces config surface onto the tracer` / `// newMetricsServer maps the observability.metrics config surface onto the` | `apps/api/internal/composition/composition.go:264` / `:291` | `func newTracing(cfg *config.Config) *obs.Tracing {` / `func newMetricsServer(cfg *config.Config, observer *obs.Observer, logger *slog.Logger) *obs.Server {` |
| 26 | 行 03 / 行 15 | `apps/api/internal/composition/composition.go:677`（作「按 plan 静态装配/聚合贡献」的锚点） | `providers = append(providers, extra...)` | `apps/api/internal/composition/composition.go:678` | `set, err := kernel.RegisterContributions(context.Background(), plan, providers)` |
| 27 | 行 15 | `apps/api/internal/composition/composition.go:771`（作「聚合已启用贡献」的锚点） | `}` | `apps/api/internal/composition/composition.go:746` | `data, err := manifest.ForModulesWithFragmentsAndGroups(` |
| 28 | 行 06 | `apps/api/internal/composition/composition.go:853`（作「构造内存」的锚点） | `// handling); a negative budget is a programming error and fails closed.` | `apps/api/internal/composition/composition.go:854` | `func newCache(cfg *config.Config) (kernel.Cache, error) {` |
| 29 | 行 17 | `apps/api/modules/digitaloffer/provider.go:25`（作「注入 kernel.TelegramDispatcher/Sender」的锚点） | `// Provider implements kernel.Provider for biz.digital-offer.` | `apps/api/modules/digitaloffer/provider.go:26`；字段 `:30` / `:31` | `type Provider struct {` / `dispatcher kernel.TelegramDispatcher // nil = channel surface not assembled` / `sender     kernel.TelegramSender` |

**合计**：**31 个错行锚点**（编号 1–29；其中 6b/17b 为第 6/17 行追加的第二子锚点）。性质分布：30 个为「指向紧邻注释/空行/大括号/邻行语句」，**1 个为主题性错锚**（第 14 条 `config.go:493`——指向 cache 默认值，而主张是 metrics/traces 默认关闭）。

### E.4 可再精确（旧锚点**不算错行**，但未指向其声称的具体语句）

| R2 位置 | 旧锚点（原文） | 更精确锚点（原文） |
|---|---|---|
| 行 07 | `apps/api/internal/ratelimit/memory.go:83` `func (l *Memory) Allow(key string, now time.Time) bool {`（「窗口/容量」部分成立） | `apps/api/internal/ratelimit/memory.go:89` `func (l *Memory) allowLocked(key string, now time.Time) bool {`（窗口裁剪/容量判定）+ `:132` `func (l *Memory) Record(key string, now time.Time) {`（容量驱逐入口） |
| 行 08 | `apps/api/internal/eventbus/memory.go:254` `func (s *subscription) run(payload []byte) {`（「panic 隔离」） | `apps/api/internal/eventbus/memory.go:256` `if rec := recover(); rec != nil {` |
| 行 09 | `apps/api/internal/mail/runtime.go:363` 之外的「渠道切换」主体 | `apps/api/internal/mail/runtime.go:234` `func (s *Switcher) buildAdapter(cfg RuntimeConfig) (kernel.MailSender, error) {`（候选配置构造成败即切换成败） |
| 行 12 | `apps/api/internal/jobs/runner.go:214` `func (r *Runner) loop() {`（「扫描恢复/过期」） | 需按具体语句另行定位（本轮未逐行定位扫描/过期语句，登记为**未定位**，不臆造锚点） |
| 行 05 | `apps/api/internal/objectstore/local.go:119` `func (s *LocalStore) Put(_ context.Context, ns kernel.ObjectNamespace, id string, body []byte, meta kernel.ObjectMeta) error {`（「文件/元数据写入」成立） | 元数据写入语句在 `Put` 内部（本轮未逐行定位，登记为**未定位**） |
| W1 | `apps/web/src/main.tsx:76` `bootHost({`（「bootHost + manifest loader」） | `apps/web/src/main.tsx:79` `manifestLoader: async () => loadAppManifestBytes(),` |
| W1 | `apps/web/src/protocol/load-page.ts:81` `export async function loadPageDocument(`（「schemaUrl 加载」） | `apps/web/src/protocol/load-page.ts:90` `const url = resolveSchemaUrl(baseURL, page.schemaUrl, params);` |

### E.5 已复核为**正确**的 R2 锚点（供 R3 免于重复核对）

`kernel/module.go:11`、`kernel/module.go:265`、`kernel/provider.go:19`、`kernel/provider.go:70`、`kernel/persistence.go:49`、`kernel/profile.go:25`、`kernel/store.go:30`、`kernel/store.go:43`、`kernel/objectstore.go:30`、`kernel/objectstore.go:60`、`kernel/cache.go:35`、`kernel/mail.go:73`、`kernel/lifecycle.go:18`、`kernel/contribution.go` 相关匿名引用；`internal/store/store.go:29`、`internal/composition/composition.go:199,203`（A-002 已校正）、`:509`、`:611,617,641`（A-002 已校正）、`:871`、`:1024`、`:1151`；`internal/server/server.go:17`、`:35`；`internal/obs/tracing.go:48`、`internal/obs/server.go:37`；`internal/objectstore/local.go:119`、`internal/objectstore/s3.go:117`；`internal/cache/memory.go:64`；`internal/jobs/model.go:19`、`internal/jobs/runner.go:125,185,214`；`internal/manifest/manifest.go:213,457`；`modules/users/provider.go:65`；`modules/digitaloffer/provider.go:68`；`assembly/assembly.go:26,37,43`；`apps/web/src/main.tsx:76`、`app/App.tsx:1049`、`app/navigation.ts:224`、`protocol/load-page.ts:81`。

## 边界声明

- 本文件是 R3 的**只读证据输入**：不构成 finding、不改变 `status`/`progress`、不关闭任何门禁、不修改 R2 矩阵 v0.2.0 正文、不重开任何 closed VP。
- 分类（现在修 / 仍 gated / 接受残余 / 明确不做）与去向（R4 文档卫生 / 另立波次）由 R3 判据表与用户按 D-001 裁决 A/B 决定；本文件只提供锚点与差值。
- 行号一律相对 HEAD `5230097570e4f4c1f245359f94c63e94e8c1a0bd`；R4/R5 引用前须按当时 HEAD 复核。
- 取证命令与观察结果见同目录 `r3-anchor-commands.txt`。
