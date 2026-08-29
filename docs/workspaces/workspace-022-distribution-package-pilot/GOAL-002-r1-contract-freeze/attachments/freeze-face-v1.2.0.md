# 契约冻结面清单 v1.2.0（2026-08-29 · R5 定稿：Web 面 + 发布形态）

- 来源：workspace-022 R1 扫描（2026-08-29；GOAL-002 E-001 / A-001）+ R2 收口（GOAL-003）+ R5 定稿（GOAL-006）
- 锚点：`KernelAPIVersion = "2.0.0"`（`apps/api/kernel/module.go`）
- 范围：**Go 侧 + Web 侧**（R3 起）
- 状态：~~v0.1.0 草案~~ → **v1.0.0 生效**（GOAL-002 D-002）→ v1.0.1（R2 S2 路径勘误）→ **v1.1.0**（R2 收口：B+ 层）→ **v1.2.0**（R5 定稿：**Web 六包边界 + peer 矩阵**（§2c）+ **发布形态注记**（§6）；B 层盘点引用 = `GOAL-003 attachments/modules-export-inventory-v0.1.md`）

## 0. 分界规则（冻结面 vs 内部自由演进面）

| 层 | 范围 | 下游 import | 变更纪律 |
|----|------|-------------|----------|
| **A · 内核冻结面** | `apps/api/kernel` 全部导出（见 §1） | ✅ 允许 | semver 约束；breaking = major + changelog 迁移说明 |
| **B · 模块契约面** | `apps/api/modules/*` 装配面（ModuleID / Provider / New* / Descriptor / CompiledPersistence / Register；符号盘点见 §2） | ✅ 允许（含模块 store 附件包） | 与 A 同纪律；先 deprecate 再删 |
| **B+ · 装配工厂面** | `apps/api/assembly` 导出签名（见 §2b） | ✅ 允许 | 与 A 同纪律；内部实现自由 |
| **C · 内部实现面** | composition / handler / server / auth / config / jobs / mail / objectstore / obs / manifest / migration / store 方言 | ❌ 默认禁止（下游经 B+ 消费；张力项见 §5） | 内部自由演进，无版本承诺 |

## 1. A 层清单（apps/api/kernel · 11 非测试文件 · ~45 导出符号）

| 文件 | 导出符号 |
|------|----------|
| module.go | `KernelAPIVersion` · `Capability` · `ContributionKeys` · `Hooks` · `Module` · `ErrorCode` · `Error` · `Registry` · `NewRegistry` · `Plan` |
| provider.go | `Provider` · `Registrar` · `ContributionSet` · `RegisterContributions` · `DefaultNavigationOrder`（var）· `NormalizeNavigationOrder` |
| contribution.go | `ContributionIdentity` · `RouteContribution` · `PageContribution` · `PermissionContribution` · `NavigationContribution` · `FragmentContribution` · `ConfigurationContribution` · `MigrationContribution` · `ContributionKind` · `RouteKey` |
| profile.go | `ProfileName` · `ProfileResolution` · `ParseModuleList` · `ResolveProfile` · `BuiltinModules` · `StandardAdminCapabilities` · `SortedModuleIDs` |
| store.go | `Dialect` · `ErrNoRows` · `Store` · `Tx` · `Result` · `Rows` · `Row` |
| persistence.go | `MigrationChecksum` · `CollectPersistence` |
| lifecycle.go | `Runtime` · `NewRuntime` |
| mail.go | `MailMessage` · `MailSender` |
| objectstore.go | `ObjectNamespace` · `ValidObjectNamespace` · `ValidObjectID` · `ErrObjectNotFound` · `ObjectMeta` · `ObjectInfo` · `ObjectStore` |
| duplicate_object.go | `IsDuplicateObject` |
| unique_violation.go | `IsUniqueViolation` |

## 2. B 层规则（模块契约面）

每模块包（`apps/api/modules/<id>`）暴露：

- `const ModuleID`（如 `"admin.users"`）
- `type Provider struct` + `New*(...)` 构造（签名属冻结面，见 §5 张力）
- `Descriptor() kernel.Module`（六字段：ID / Version / KernelAPIRange / DependsOn / Requires / Contributions）
- `CompiledPersistence() ([]kernel.MigrationContribution, error)`
- `Register(ctx context.Context, reg kernel.Registrar) error`（经 reg.HTTP / Schema / Authorization / Navigation / Manifest 贡献六能力）

候选 B 层包：account · activity · authsession · compiled · dashboard · datadictionary · datapermission · datatransfer · dev/examples · filelibrary · logincaptcha · mfa · notifications · operationlog · recyclebin · roles · scheduledtasks · schemarender · settings · systemmonitoring · users · wallet（22）；corepersistence / jobs 为内核支撑；各 `*/store`、`*/migration`、`*/schema`、`*/manifest`、`*/systemdata` 为模块附件包。

兼容机制（已存在，继续为门禁）：`Descriptor.KernelAPIRange` 与 `kernel.KernelAPIVersion` 的 registry 校验（fail-closed；`kernel_test.go` `TestRegistryValidatesKernelAPIRanges`）。

**B 层符号盘点**（v1.1.0 回填引用）：`GOAL-003-r2-go-library-consumption/attachments/modules-export-inventory-v0.1.md`（22 包顶层导出扫描；构造参数含 internal 类型者以 ⚠️ 标注——S3 后仅 `assembly` 收敛包内保留）。

## 2b. B+ 层（装配工厂面 · v1.1.0 增列）

`apps/api/assembly`（v0.1.0 experimental · D-003 方案 β）：

| 工厂 | 签名 | 说明 |
|------|------|------|
| `OpenStore` | `(ctx, dialect kernel.Dialect, path, dsn string, catalog []kernel.MigrationContribution) (kernel.Store, error)` | 双方言打开 + 迁移随 Open apply |
| `NewAuthenticator` | `(secret []byte, accessTTL, refreshTTL time.Duration, runner authsession.TxRunner) *auth.Authenticator` | 返回 internal 类型，调用方**类型推断消费**（不命名） |
| `NewMailSender` | `(st kernel.Store, retentionCap int) kernel.MailSender` | mock 站内出站记录渠道 |

规则：签名 semver 约束（与 A/B 同纪律）；参数位不导出 internal 类型（仅返回值可，配合类型推断）；内部实现自由演进。机制：`kernel.Store.Run` 与 `authsession.TxRunner` 结构同构，Store 可直接作 runner。

## 3. semver 锚点与版本语义

- kernel（A 层）：**major = KernelAPIVersion 主号**（当前 2）；minor/patch 随主线节奏；`KernelAPIRange` 由各模块声明兼容窗。
- 模块（B 层）：版本 = `Descriptor.Version`；对 kernel 的兼容窗 = `KernelAPIRange`。
- 装配工厂（B+ 层）：版本 = package 语义版本；签名变更 = A/B 同纪律。
- 判定与流程：`semver-breaking-policy-v0.1.0.md`；发布记录：`changelog-template-v0.1.0.md` 起用。

## 4. 待完成（后续阶段回填）

- [x] R2：B 层符号盘点回填（F-002 fixed · inventory v0.1）；C 层泄漏收敛（F-001 fixed · 方案 β assembly 实证）
- [x] R3：Web 侧 npm 包组冻结面（I-002 闭合 · §2c）
- [x] R5：发布通道与 Go tag / npm artifact 版本绑定（I-003 定案 · §6）

## 5. 边界张力与决策点

- **C 层命名面（已闭合）**：R1 记录的三候选经 R2 实测收敛 = **方案 β**（公开装配工厂 `apps/api/assembly`，用户 2026-08-29 裁决）；① kernel 接口化（β'）与③ 白名单外移（α：auth/store 提升）留 **go/no-go 后评估**。
- **`kernel.Registrar` 调用面**：模块 Register 体使用 handler 工厂（C 层）构造路由——B 层内部实现可引用 C 层；此属模块自身实现自由，不要求模块内部纯净隔离。
## 2c. Web 侧包面（v1.2.0 定稿 · I-002 闭合）

| 包 | 来源 | 导出面（v0.1/v0.2 产物） | peer | 状态 |
|----|------|--------------------------|------|------|
| @schema-ui/protocol | src/protocol | app-manifest 协商/校验/路由 + page 文档加载 + conformance | ajv（bundle 自包含） | 0.2.0 产物实证 |
| @schema-ui/renderer（粗粒度单包含 ui/i18n 面） | src/renderer + components + i18n | RenderPage / I18nProvider / registerCustomComponent / 类型 | react ^19 / react-dom ^19（external） | 0.1.0 产物实证 |
| @schema-ui/ui / 	heme / lib / shell（六包细化目标） | 待六包化专项 | 按边界设计 v0.1 | React/Tailwind | **go 后正式化**（F-006 随 monorepo 化） |

Tailwind 契约：包组件仅 className，零 CSS 产物；样式由下游构建编译（Tailwind 4 @source 扫描包产物——golden-web README 指引）。

## 6. 发布形态注记（v1.2.0 · R5）

- Go：单模块粗粒度 tag（v0.0.x 试点 → 语义版本随 semver-breaking-policy）；发布载荷 = go get + replace/registry；装配经 apps/api/assembly（B+ 层）。
- Web：scripts/pack-npm-packages.mjs 一键 tgz（registry 载荷）；golden-web tarball 安装实证（GOAL-006 E-001）。
- 版本绑定：changelog 三向一致（kernel 主号 / 模块 KernelAPIRange / npm peer 窗）。
