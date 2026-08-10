---
id: r4-c1-freeze-package-draft
doc: proposal
goal: GOAL-005-r4-full-module-migration
date: 2026-08-05
status: accepted
decision_state: user_accepted
source: candidate-response-to-A-003
accepted_by: magicvr
accepted_date: 2026-08-05
accepted_via: D-003-r4-c1-decisions (GOAL-005 与 GOAL-006)
---

# R4 C1 Provider 冻结包（契约正文）

本附件是 A-003 六项 required finding 的候选响应，**已由用户整包接受**，作为
GOAL-005/GOAL-006 `D-003` 的精确契约正文。正文含 Contribution 字段、双检规则、
注册与发布顺序、`CompiledPersistence()` 契约、compiled-global Persistence 规则、
Authorization/seed/security owner matrix 与 operationlog Option A 边界；C2 实施
不得在未记录的情况下改变身份、冲突键、安全语义或顺序。文件保持该路径与文件名以
维持既有引用可追溯；`status: accepted` 为权威状态，文件名中的 draft 仅为历史遗留。

## 1. 待裁决轴

| 轴 | 当前候选 | 用户必须确认的内容 |
|---|---|---|
| Provider | framework-agnostic `Provider` + Plan-owned `Registrar` | 是否采用该公共形状，以及是否接受下述字段、双检和 composition-owned construction |
| Records | historical-only，维持 `0006 records_retire` 后的现状 | 是否恢复 Records 产品 CRUD；恢复时必须使用新迁移版本，不得改写 `0006` |
| operationlog | Option A：保留 best-effort，R4 不自动 purge/archive/delete | 是否接受日志失败可能造成审计缺口的 residual；Option B/C 需要单独扩大行为和数据生命周期范围 |

Records 与 operationlog 是 P-004 用户裁决项。Provider 方案虽有工程推荐，精确
字段也必须以用户确认后的 D-003 或其修订版为准。

## 2. Provider 与 Registrar 候选契约

### 2.1 公共形状

模块公共 API 只使用标准库、内核契约和稳定应用接口，不导入 Fx。候选形状：

```go
type Provider interface {
    Descriptor() Module
    CompiledPersistence() ([]MigrationContribution, error)
    Register(context.Context, Registrar) error
}

type Registrar interface {
    HTTP(RouteContribution) error
    Schema(PageContribution) error
    Authorization(PermissionContribution) error
    Navigation(NavigationContribution) error
    Manifest(FragmentContribution) error
}
```

Provider 必须是编译期候选集的一部分。`Descriptor()` 返回的 module ID、版本、
Kernel API range、依赖、capabilities 和声明的 contribution keys 是注册前元数据；
`Register` 只能为同一个 module ID 写入临时 surface 集合。Persistence 仍是六类
标准贡献之一，但通过 `CompiledPersistence()` 进入 compiled-global catalog，刻意
不通过 enablement-gated `Registrar` 注册。

### 2.2 Contribution 最小字段

每类结构化 contribution 都必须带不可变身份：`ModuleID`、`Kind`、`Key`。`Key`
在同一 `Kind` 内全局唯一；模块版本更新不能复用其他模块的 key。

| 类型 | 冻结前至少需要的字段 | 全局校验 |
|---|---|---|
| HTTP | method、pattern、handler、中间件顺序、公开/认证边界 | method+pattern 唯一；禁用模块不得发布 route |
| Schema | page ID、resource/action IDs、data source、owner module | page/resource/action ID 唯一；owner 必须来自 contribution |
| Authorization | permission key、resource、action、policy reference、secret sensitivity | key 唯一；缺少授权或 disabled permission fail closed |
| Navigation | node ID、parent、order、label、visibility expression、permission key | node ID 唯一；引用的 page/permission 必须存在 |
| Manifest | fragment ID、protocol/API version、required capabilities、确定性 payload | fragment/page/navigation/app identity 冲突 fail closed；登录前无秘密 |
| Persistence | global version、name、module ID、checksum、transaction callback、tombstone、reconcile version | global version/name/module/checksum 唯一；缺口、漂移、未知记录 fail closed |

`Configuration` 在 R4 不新增独立 Registrar 方法；若模块需要配置，必须继续在
`Module.Contributions.ConfigNamespaces` 声明稳定命名空间，并由后续明确的配置
contract 处理。R4 不宣称已有配置 namespace 已经完成 module-owned runtime
configuration 迁移。

候选字段名和类型如下。它们是冻结包的规范候选，C2 不得在未记录的情况下改变
身份、冲突键或安全语义；可在实现时增加不改变协议的内部字段。

```go
type ContributionIdentity struct {
    ModuleID string
    Key      string
}

type RouteContribution struct {
    ContributionIdentity
    Method     string
    Pattern    string
    Handler    http.Handler
    Middleware []string
    Public     bool
}

type PageContribution struct {
    ContributionIdentity
    PageID     string
    Resources  []string
    Actions    []string
    DataSource string
    Owner      string
}

type PermissionContribution struct {
    ContributionIdentity
    Permission        string
    Resource          string
    Action            string
    PolicyID          string
    SecretSensitivity string
}

type NavigationContribution struct {
    ContributionIdentity
    NodeID     string
    Parent     string
    Order      int
    Label      string
    Visibility string
    Permission string
}

type FragmentContribution struct {
    ContributionIdentity
    FragmentID             string
    ProtocolVersion       string
    ModuleAPIVersion      string
    RequiredCapabilities  []string
    JSON                   []byte
}

type MigrationContribution struct {
    ContributionIdentity
    Version          int
    Name             string
    Checksum         string
    Apply            func(*sql.Tx) error
    Tombstone        bool
    ReconcileVersion int
    ReconcileChecksum string
    Reconcile        func(*sql.Tx) error
}
```

`Handler` 使用标准 `net/http`，`Apply`/`Reconcile` 使用标准 `database/sql`；上述
类型不引用 Fx。`PolicyID`、`Visibility` 和 `JSON` 必须有校验器和确定性编码规则，
不能由实现者默认为任意字符串或任意 JSON。`CompiledPersistence` 对每个 compiled
provider 调用一次，返回的 migration contribution 不受 Profile 过滤。

`ContributionIdentity.Key` 是规范语义 ID，不是另一个可任意发明的别名：HTTP 使用
大写 method 加空格加 pattern；Schema 使用 `PageID`；Authorization 使用
`Permission`；Navigation 使用 `NodeID`；Manifest 使用 `FragmentID`；Persistence
使用稳定 `Name`。实现必须校验对应字段与 Key 相等。Persistence 另外对 global
`Version`、Name、ModuleID、Checksum 和 `ReconcileChecksum` 做全局校验。

### 2.3 Plan 元数据与运行时双检

现有 `Module.Contributions` 继续承担注册前的快速、可诊断声明；结构化 Registrar
承担运行时的真实贡献。两者不是互相绕过的两套 owner：

1. composition 先从全部 compiled providers 建立 registry，再解析 Profile 和
   `modules.enabled` 得到唯一 immutable `Plan`。
2. 每个启用 provider 的 `Descriptor` 必须与 Plan 中同 ID 的 Module 完全匹配。
3. `Register` 只能写入 descriptor 已声明的 `Kind + Key`；未声明、重复或字段身份
   不匹配立即返回结构化错误。
4. 临时 `ContributionSet` finalize 前执行全局冲突、引用完整性、确定性排序和
   security checks。失败时整个集合丢弃，不发布部分 route、Schema 或 Manifest。
5. `CompiledPersistence()` 不从启用 Plan 收集，详见第 4 节；其 catalog 与
   HTTP/Schema 等 surface contribution 的 enablement 过滤严格分离。`Registrar`
   没有 Persistence 写入口，避免实现者误把 migration 放入 Plan-gated path。

### 2.4 依赖注入和 Fx 边界

Provider 由 composition-owned factory 构造。factory 可以将 `*store.Store`、认证/
授权接口、日志和配置等框架无关依赖传入模块构造函数；模块包不得接收或导出
`fx.In`、`fx.Out`、`fx.Option`，也不得使用 service locator。Fx 只在 composition
中把已构造的 provider、Registrar、HTTP server 和 lifecycle 连接起来。

R4 C2 必须增加静态检查，确认 `apps/api/internal/modules/**` 和 `apps/api/internal/kernel`
不 import `go.uber.org/fx`；同时测试 provider register error、duplicate contribution
和 dependency mismatch 的结构化错误。

## 3. 注册、发布和生命周期顺序

候选顺序固定为：

1. 收集全部 compiled provider descriptors，并对每个 provider 调用
   `CompiledPersistence()`，形成唯一 compiled-global persistence catalog。
2. 解析 Profile/`modules.enabled`，得到 Plan；不满足依赖或 capability 时停止。
3. 仅对 enabled standard Admin providers 调用 `Register`，写入临时 surface set；
   persistence 不在此步骤注册。
4. 对 surface set 做全局冲突、引用、权限、Manifest secrecy 和确定性排序校验。
5. 对第 1 步的 compiled-global persistence set 做 ledger、tombstone、reconcile
   metadata 校验；该校验不读取 Plan 的 enabled module 列表。
6. finalize 后才由 composition 建立 mux、Manifest response 和 Fx graph。
7. 所有注册和校验成功后才允许 listen；Start、Ready 成功后才调用 `Serve`。
8. 任一 register/conflict/Start/Ready 失败，都不得留下可用的部分 surface；按已
   启动顺序反向 Stop，并关闭 listener、store 和其他外部资源。

每类失败必须保留稳定分类：register error、contribution conflict、dependency/
capability mismatch、Start failure、Ready failure 分别映射结构化 error code，且
测试能证明 register/conflict 时未 listen、Start/Ready 失败时已启动资源被清理。
`mvp` 与 `admin` 两个 Profile 必须各运行一次该矩阵；一个 Profile 的通过不能替代
另一个 Profile 的证据。

`/readyz` 在 R4 不能被宣称为完整模块图 readiness；当前只证明 store ping。若 R4
实施触及 readyz，必须同时实现迁移、reconcile、必需依赖和 module lifecycle 的
真实 readiness，不能以当前 handler 的 200 响应代替。

## 4. Compiled-global Persistence 规则

### 4.1 迁移集合

迁移 catalog 的输入是静态编译的一方 provider 候选集加上历史 tombstone，而不是
enabled Plan。实现上唯一允许的收集入口是 `Provider.CompiledPersistence()`；
`Registrar` 没有 Persistence 方法，任何从 `Register` 或 Plan 追加 migration 的实现
均违反本冻结候选。启用过滤只决定 HTTP、Schema、Navigation、Manifest 和可见授权 surface；
禁用模块不回滚、不删除其表或数据。

迁移 descriptor 至少保留：

- 全局递增 version、稳定 name、原始 module ID 和 checksum；
- 迁移执行函数或明确的 tombstone 状态；
- 版本化 system-data reconcile 标识和 checksum；
- 不能覆盖用户字段、不能删除用户数据的 ownership 说明。

现有 `0001` 至 `0008` 不重编号、不改名、不改 checksum。未来 module-owned migration
只能追加全局版本；退役模块继续以 descriptor/tombstone 进入 collector。Store 的
硬编码 `ModuleID: "core.persistence"` 只能作为当前历史事实，不能作为 R4 终态契约。

### 4.2 Fresh、upgrade 和 reconcile

- migration 与 ledger row 在同一个事务内完成；失败不得留下半个 schema 或 ledger row。
- fresh bootstrap 只建立初始管理员和基础系统数据；不将用户字段当作 seed 所有权。
- reconcile 是独立的、版本化、幂等路径；只补齐明确归属的 system keys，不覆盖用户
  字段、不删除用户数据。
- 既有数据库在 Settings/Activity 或其他模块 disabled 时仍执行 compiled migration，
  并保留已存在的数据；关闭 UI 不等于回滚 persistence。
- 必须测试 compiled-global collector、tombstone、unknown applied migration、
  checksum drift、缺口、事务回滚、重启幂等、disabled profile 数据保留和 user-owned
  field preservation。

## 5. Authorization、seed 与敏感信息

Authorization contribution 是 permission key、handler policy 和 system-data reconcile
的唯一候选来源。composition 不再按 module ID 分散维护 permission owner；中心
reconcile 可以执行，但只消费已验证的 module contributions：

- seed、Manifest、handler authorization 和 Schema action 使用同一稳定 permission key；
- duplicate、missing reference、disabled route with enabled permission mismatch 都
  fail closed；
- disabled UI 不删除已存在的系统权限、菜单或用户数据，新增 surface 仍按 Plan 过滤；
- public Manifest 登录前可读但不得包含 secret、token 或用户个性化信息；
- operationlog detail 不得包含 password、token 或 secret；
- `core.auth-session`、RBAC 和 `core.operationlog` 的 cross-cutting owner 在 module
  matrix 中单独登记，Activity 只是可选查询/UI，不得关闭 operationlog writer。

C2/C3 必须覆盖 permission/seed/Manifest/handler 一致性、公开 Manifest secrecy、
用户响应中的 password 字段隔离，以及 Activity disabled 后 Users/Roles/Auth/Settings
写入仍产生 operationlog 的闭环测试。

横切 owner matrix 也必须随 D-003 附件落盘：

| owner | 永久职责 | 可选 surface | R4 验证 |
|---|---|---|---|
| `core.auth-session` | session、认证 API、授权边界 | auth routes | auth failure、password 不出响应、依赖闭包 |
| RBAC reconcile | permission/menu/grant system data | navigation projection | key 一致、幂等、不覆盖用户字段 |
| `core.operationlog` | 所有关键写操作的 append writer | 无 | Activity disabled 仍写、secret exclusion |
| `admin.activity` | operationlog 查询与只读 UI | Activity routes/Schema/Nav | disabled 时 UI 消失但 writer 保持 |

`Module.Hooks` 继续是生命周期 owner；Registrar 不拥有或替换 `Hooks`。Provider
registration 只完成静态 surface/catalog 收集，composition 负责将 `Module.Hooks` 按
拓扑顺序适配到 lifecycle。Observability 事件若在 R4 需要新增，必须通过独立且
框架无关的 contribution 或明确留在 composition，不得隐含混入 Registrar。

## 6. operationlog 选项的冻结边界

### Option A 候选验收

若用户选择 A，冻结的只是 R4 运行行为：业务写入成功后 best-effort append；append
失败记录服务日志且不翻转业务成功；R4 不自动 purge、archive 或 delete
`operation_log` 行，迁移和重启保留现有行，Activity UI 开关不改变 writer。

这不是永久 retention policy。正式 D-003 必须另附 residual：

| 字段 | 当前状态 |
|---|---|
| residual | operationlog append 失败可能产生审计缺口；duration/archive 尚未定义 |
| scope | R4 Users/Roles/Auth/Settings 写入和现有历史 events |
| owner | `pending_user` |
| review trigger | 合规/运营 retention 要求、日志规模阈值、恢复演练发现缺口，或进入 R5 数据生命周期决策 |
| review date | `pending_user` |

### Option B/C 候选边界

Option B 必须让所有 writer 共用同一 `*sql.Tx`，并接受 API 错误语义、锁/重试和并发
行为变化；Option C 必须另行定义 retention duration、archive/restore/purge、失败
重试、查询一致性和数据恢复契约。任何一项都不能在 C1 被隐含选择。

无论选择哪项，必须覆盖 append/read/order、failure injection、迁移/重启、Auth/Settings
writer、Activity disabled 仍写入，以及 event CHECK 兼容性。Option A 的建议兼容矩阵
保留现有 HTTP success semantics、Schema 和 event set；改变这些契约必须新建 migration
或取得另行兼容决策。

## 7. Records 分叉和兼容性

### Historical-only 候选

维持 `0006 records_retire` 的现状：R4 迁移当前仍存在的 Schema-driven Admin 能力，
保留历史 `records.*` operation events，不声称存在 Records 产品 CRUD。R4 success
criteria 和 VP 对齐记录必须显式写明该收敛，不得只通过删除测试或静默排除。

### Restore CRUD 候选

恢复 Records 产品 CRUD 时必须另建高于 `0008` 的 migration、表/权限/menu/Schema/
handler/恢复测试；不得重写或复活 `0006` 的已发布语义，也不得以历史 operation events
代替现行产品数据。

### 中心特例切换顺序

不论 Records 分叉，中心特例的候选切换顺序为：

1. 先添加 provider metadata 和无发布的 contract tests。
2. 再由 typed provider 生成 surface，并在测试中与现有中心输出做兼容比较；不能在
   生产 mux 中永久注册两份相同 route。
3. 切换 composition 消费 provider finalize 结果，移除 Schema owner map、Manifest
   `adminModules` 和按 module ID 的 handler Register 分支。
4. 保留必要的历史 migration/event tombstone，证明旧 API status/payload、Schema ID、
   operationlog CHECK 和 migration ledger 未被无记录地改变。
5. 以静态扫描和运行测试证明中心业务 Register、全局 Schema fixture 占用和永久双路径
   已删除；tombstone 只保留数据兼容所需事实。

兼容性清单至少固定为：

- HTTP：现有 auth、users、roles、settings 写入的成功 status、错误分类和响应字段；
- Schema：`users`、`roles`、`settings`、`activity` page/resource/action IDs；
- operationlog event CHECK：`records.create`、`records.update`、`records.delete`、
  `auth.login`、`auth.logout`、`auth.refresh`、`users.create`、`users.update`、
  `users.delete`、`roles.create`、`roles.update`、`roles.delete`、`settings.update`；
- persistence：`0001` 至 `0008` 的 version/name/checksum、历史行保留和新版本追加规则。

每一项都要有切换前后可比较的测试或静态证据；“行为应保持兼容”不能单独作为
closure evidence。

## 8. C1 关闭条件

C1 只有在以下证据齐全后才可标记完成：

- 用户对 Provider 精确契约、Records 分叉和 operationlog 选项作出书面裁决并形成
  D-003；
- R4-I002/R4-I003/R4-I004 从 `collecting` 变为 `verified`，每项都有对应 evidence；
- A-003 的 F-IND-R4-OPT-001 至 006 逐项以 `fixed` 或合法用户 residual 响应，不能
  用推荐文字代替闭合；
- self + Grok independent 对最终冻结包复审，且没有 open required finding；
- C2 实施计划、迁移/seed/reconcile、compatibility 和 failure-matrix 的边界已写入
  execution ledger。

本草案只完成“候选响应材料”这一事实，当前不得推进 `GOAL-005` progress 或进入 C2。
