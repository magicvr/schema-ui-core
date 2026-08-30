---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-contract-freeze
version: 0.1.0
---

# E-001 · S1 扫描事实（2026-08-29）

扫描 `apps/api`（HEAD `6008c16d` 前后同源；非测试文件口径）：

1. **包结构**（顶层）：`cmd/server`（main 3 文件）、`internal/kernel`（11）、`internal/composition`（11）、`internal/modules`（166，24 子目录）、`internal/store`（8 非测试）、`internal/handler`（98）、`internal/auth`（3）、`internal/config`（7）、`internal/jobs`（8）、`internal/mail`（13）、`internal/objectstore`（5）、`internal/obs`（7）、`internal/manifest`（2）、`internal/migration`（2）、`internal/server`（2）等；`pkg/version`（Version/Commit/BuiltAt，ldflags 注入）。
2. **kernel 导出面**（A 层候选，11 文件）：`module.go`（`KernelAPIVersion="2.0.0"` · `Capability` · `ContributionKeys` · `Hooks` · `Module` · `ErrorCode` · `Error` · `Registry` · `NewRegistry` · `Plan`）、`provider.go`（`Provider` · `Registrar` · `ContributionSet` · `RegisterContributions` · `DefaultNavigationOrder` · `NormalizeNavigationOrder`）、`contribution.go`（九类 Contribution + `ContributionKind` · `RouteKey`）、`profile.go`（`ProfileName` · `ProfileResolution` · `ParseModuleList` · `ResolveProfile` · `BuiltinModules` · `StandardAdminCapabilities` · `SortedModuleIDs`）、`store.go`（`Dialect` · `ErrNoRows` · `Store` · `Tx` · `Result` · `Rows` · `Row`）、`persistence.go`（`MigrationChecksum` · `CollectPersistence`）、`lifecycle.go`（`Runtime` · `NewRuntime`）、`mail.go`（`MailMessage` · `MailSender`）、`objectstore.go`（`ObjectNamespace` · `ValidObjectNamespace` · `ValidObjectID` · `ErrObjectNotFound` · `ObjectMeta` · `ObjectInfo` · `ObjectStore`）、`duplicate_object.go`（`IsDuplicateObject`）、`unique_violation.go`（`IsUniqueViolation`）——合计约 45 个导出符号。兼容窗校验先例：`kernel_test.go` `TestRegistryValidatesKernelAPIRanges`。
3. **模块契约模式**（B 层候选，抽样 `internal/modules/users/provider.go`）：`const ModuleID = "admin.users"` · `type Provider struct` · `func New(a *auth.Authenticator, repository *authsession.Repository, operations operationlog.Recorder, mailSender kernel.MailSender, publicBaseURL string) *Provider` · 方法 `Descriptor() kernel.Module`（`Version:"2.0.0"` · `KernelAPIRange:">=2.0 <3.0"` · `DependsOn` · `Requires: kernel.StandardAdminCapabilities()` · `Contributions` 六键）· `CompiledPersistence()` · `Register(ctx, kernel.Registrar)`（经 `reg.HTTP/Schema/Authorization/Navigation/Manifest` 注册）。全仓 24 个模块子目录中 16 个标准模块有 provider.go；core 支撑（corepersistence/jobs）与 store 子包为装配依赖。
4. **组合根**（`internal/composition/composition.go`）：静态 import 全部 B 层模块 + 各 store + C 层（auth/handler/jobs/mail/manifest/server/objectstore/obs…），Fx 装配——**这正是包化后下游自建组合根的模板**（从「import 全集」变「import 所选集」）。
5. **张力发现**：模块 `New*` 构造签名引用 C 层类型（如 `*auth.Authenticator`）——下游组合根将被迫 import C 层白名单；其余 `NewService(repo *X) *Service` 类构造引用 B 层 store 子包（`modules/*/store`），属 B 层内部附件（见 D-001 §4 候选收敛路径，R2 定）。

产物：`attachments/freeze-face-v0.1.0.md`（清单）、`semver-breaking-policy-v0.1.0.md`、`changelog-template-v0.1.0.md`。