# 契约冻结面清单 v1.0.1（2026-08-29 · R2 S2 路径勘误）

- 来源：workspace-022 R1 扫描（2026-08-29；GOAL-002 E-001 / A-001）
- 锚点：`KernelAPIVersion = "2.0.0"`（`apps/api/kernel/module.go`）
- 范围：**Go 侧**；Web 侧 npm 包组冻结面 = I-002（R3 另行落盘）
- 状态：~~v0.1.0 草案~~ → **v1.0.0 生效**（2026-08-29 用户确认「确认冻结面，关门 R1，开 R2」；GOAL-002 D-002）→ **v1.0.1**（2026-08-29 · R2 S2 外移重构路径勘误：`apps/api/internal/kernel` → `apps/api/kernel`、`apps/api/internal/modules` → `apps/api/modules`；契约内容不变，editorial）；B 层全量符号随 R2 回填（A-001 F-002）

## 0. 分界规则（冻结面 vs 内部自由演进面）

| 层 | 范围 | 下游 import | 变更纪律 |
|----|------|-------------|----------|
| **A · 内核冻结面** | `apps/api/kernel` 全部导出（见 §1） | ✅ 允许 | semver 约束；breaking = major + changelog 迁移说明 |
| **B · 模块契约面** | `apps/api/modules/*` 装配面（ModuleID / Provider / New* / Descriptor / CompiledPersistence / Register） | ✅ 允许（含模块 store 附件包） | 与 A 同纪律；先 deprecate 再删 |
| **C · 内部实现面** | composition / handler / server / auth / config / jobs / mail / objectstore / obs / manifest / migration / store 方言 | ❌ 默认禁止（张力项见 §5） | 内部自由演进，无版本承诺 |

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

## 3. semver 锚点与版本语义

- kernel（A 层）：**major = KernelAPIVersion 主号**（当前 2）；minor/patch 随主线节奏；`KernelAPIRange` 由各模块声明兼容窗。
- 模块（B 层）：版本 = `Descriptor.Version`；对 kernel 的兼容窗 = `KernelAPIRange`。
- 判定与流程：`semver-breaking-policy-v0.1.0.md`；发布记录：`changelog-template-v0.1.0.md` 起用。

## 4. 待完成（后续阶段回填）

- [ ] R2：逐包枚举 B 层导出并回填全量符号清单（A-001 F-002）；验证模块构造装配面、收敛 C 层泄漏（F-001 / D-001 §4）
- [ ] R3：Web 侧 npm 包组冻结面（protocol / renderer / shell / ui 导出边界 + peer 版本耦合矩阵）（I-002）
- [ ] R5：发布通道与 Go tag / npm artifact 版本绑定（I-003）

## 5. 边界张力与决策点

- **C 层泄漏（A-001 F-001）**：模块 `New*` 签名引用 C 层类型（例 `users.New(*auth.Authenticator, *authsession.Repository, ...)`）。收敛候选：① `auth.Authenticator` 等上移 A 层契约；② 公开装配工厂/适配器包（如 `foundation/assembly`）；③ 有限 C 层白名单。**R2 试点裁定，不预先决定。**
- **`kernel.Registrar` 调用面**：模块 Register 体使用 handler 工厂（C 层）构造路由——B 层内部实现可引用 C 层；此属模块自身实现自由，不要求模块内部纯净隔离。