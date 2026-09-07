# A-002 · grok build 独立审计原始输出（2026-09-01 · 原样收录）

> 来源：本地 grok build（grok-4.6 · reasoning high · headless 单轮），`source: independent`。编排器按 P-003 誊入，未修改意见内容；报告正文自「# A-002」起为 grok 原样输出。grok 当场独立复跑：`go vet ./kernel/...`（0）、`go test ./kernel/... -count=1`（PASS）、`git status --short` / `git diff --stat`（越界核账）。grok 未创建或修改任何文件。

---

# A-002 · R1 合同冻结关门独立交叉审计（2026-09-01）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **date**：2026-09-01
- **类型**：close-out（C3 independent；对照 A-001 self）
- **scope**：`workspace-026-cache-port` / `GOAL-002-r1-contract-freeze` 全量——C1 信息裁决（I-026-001/002/003 · P-004）、C2 合同 D-002 §1～§8 与 `apps/api/kernel/cache.go` 逐节一致性、合同级快测 `kernel/cache_test.go`、未命中/零值/拷贝边界/Set 校验顺序、R1 越界核账（Profile 默认集 / 模块矩阵 / Manifest / `go.mod` / Charter / §0 范围外清单）
- **verdict**：**pass**
- **开放 required 计数**：**0**

本意见不修改 `status` / `progress` / 方案正文；落盘由编排器完成。响应请用 `/govern`。

## 范围与区间

| 项 | 值 |
|----|----|
| 工作区 | `workspace-026-cache-port`（`root_goal` = `GOAL-001-cache-port`；`canonical_scope` = 本区；`shared_materials_catalog` = `none`） |
| 被审目标 | `GOAL-002-r1-contract-freeze`（`parent` = `GOAL-001-cache-port`；`status: active`；C1 已关门 / C2 进行中 / C3 计划） |
| 冻结分母 | `01-decision/D-002-cache-port-contract.md` v0.1.0 `accepted` |
| 被审实现 | `apps/api/kernel/cache.go`、`apps/api/kernel/cache_test.go`（接口 + 校验 helper + sentinel；无供应商实装） |
| 先例 | `kernel.Store` / `ObjectStore` / `MailSender`：非泛型 · ctx 首位 · fail-closed helper · `errors.Is` sentinel |
| 共享资料 | 无固定引用（`none`）；本意见未把资料目录当作事实或关闭证据 |
| 不在本意见关闭 | I-026-004（R3）；判据 #2/#3 实装（R2）；Redis 接缝正文（R3） |

未读取或比较其他工作区上下文。

## 独立复核命令（本会话实测）

在 `apps/api`：

| 命令 | 结果 |
|------|------|
| `go vet ./kernel/...` | **0**（`VET_EXIT=0`） |
| `go test ./kernel/... -count=1` | **PASS**（`ok github.com/magicvr/schema-ui-core/apps/api/kernel` 0.880s） |
| `go test ./kernel -count=1 -run "TestValidCacheNamespace\|TestValidCacheKey\|TestCacheSentinelErrors" -v` | 三函数全绿；子测见下节 |

仓库根：

| 命令 | 结果 |
|------|------|
| `git status --short` | 见越界核账 |
| `git diff --stat` | 已跟踪 5 文件、+21/−17；未暂存。未跟踪：`cache.go` / `cache_test.go` / `GOAL-002-r1-contract-freeze/` / Root `E-002-r1-goal-established.md` |

红线路径（`apps/api/go.mod` / `go.sum`、`docs/vision/charter.md`、Profile / Manifest / composition / config）**无 diff**。`apps/api/go.mod` 无 `redis` 字符串。

## 信息门禁核验（P-005 / P-004）

C1 三条均有 **用户书面裁决**（`D-001` `status: accepted`，2026-09-01），目标 / Root / 决策索引均标 **verified**。独立审计不重放原对话，只核台账是否可追踪、是否与合同一致。

| ID | 级别 | 最晚阶段（GOAL-002） | 台账状态 | 用户裁决（D-001） | 合同落点 | 判定 |
|----|------|----------------------|----------|-------------------|----------|------|
| I-026-001 | required | C1（方案冻结 + 判据 #1） | verified | 采纳①：非泛型 `[]byte` + `(value, ok)` 区分未命中/零值；`Typed[T]` 为 R2 便利层，不进端口 | D-002 §1；`Cache`/`CacheView` 无类型参数，负载 `[]byte` | **通过**。非静默代替 |
| I-026-002 | required | C1（判据 #3/#6） | verified | 采纳①：惰性清理 + 配置化容量驱逐；无后台协程 → VP-021 不停机义务；容量键 R2、默认 10000 | D-002 §5/§6；`ExpiryPolicy` 冻结，无 goroutine / Hooks | **通过** |
| I-026-003 | non-blocking | C1（判据 #1/#4） | verified | 采纳①：`Cache.Namespace(ns) → CacheView` + fail-closed 形状校验 | D-002 §2；开放集合 + 段式正则 | **通过**（主裁决） |
| I-026-004 | non-blocking | R3 | 待确认 | — | 不在本目标关闭 | **不阻断 R1** |

无 `deferred required`、无 `accepted-residual` 未写范围。C2/C3 无到期未关闭 required 信息项。

**D-001 与 D-002 的一处预留细节不一致**（不改变已裁决的端口形态，见 F-001）：

- D-001 I-026-003 接受行写 Redis 前缀 **`cache:<ns>:<key>`**
- D-002 §2 与 `cache.go` 文件头写 **`<ns>:<key>`**，并声明实际前缀由 R3 接缝文档落盘

主冻结项（scoped 视图 + 形状校验）一致。前缀属 R3 预留，**不把 R1 合同形态打回重裁**，但冻结分母与裁决行必须在 R3 前对齐。

VP-026 信息表 I-026-002 最晚阶段仍写 **R2**、影响门禁仅「退出判据 3」；Goal 台账写 **R1 / #3/#6**。该项已 verified，不构成到期 required；台账应对齐（F-007）。

## 逐节一致性核验（D-002 §1～§8 ↔ `cache.go`）

R1 交付是 **接口 + 校验函数 + sentinel**，与 D-002 §9 C2 范围及 ObjectStore R1 先例同构。行为语义（命中/拷贝/惰性删除/容量）以注释义务约束 R2 供应商，本波无 adapter 可执行那些路径。

| 节 | 合同要点 | 实现 | 判定 |
|----|----------|------|------|
| **§1 端口形状** | 非泛型；`Namespace(ns) (CacheView, error)`；`Get(ctx, key) ([]byte, bool)`；`Set(ctx, key, []byte, ExpiryPolicy) error`；`Delete(ctx, key) error` | 签名与合同示例一致；无 `Cache[T]` / `any`；ctx 在 I/O 方法首位（`Namespace` 无 ctx，合同如此） | **一致** |
| **§1 未命中/零值** | miss = `(nil, false)`；空值 `[]byte{}` 命中 = `(空 slice, true)`；nil 不可写 | 文件头：ok 区分 miss 与空值、nil fail-closed。`Get` godoc 写明 miss `(nil, false)`，**未在方法注释复述空值命中 `(空, true)`**。`Set` godoc 写 non-nil / `ErrInvalidCacheValue` | **语义无矛盾**；Get 空值命中只在文件头，方法 godoc 不完整（F-005） |
| **§1 拷贝边界** | Set 复制入参；Get 返回新拷贝；不得共享内部底层数组 | `Set`：「input slice is copied at the boundary」；`Get`：「fresh copy… never share internal buffer」 | **注释与合同一致**；本波无拷贝代码（接口），R2 必须兑现 |
| **§2 命名空间** | 开放集合；`^[a-z0-9]+(-[a-z0-9]+)*$`；≤64 字节；`ValidCacheNamespace` 唯一入口；非法 → `ErrInvalidCacheNamespace`，不回落默认 | 同正则；`CacheNamespaceMaxLen = 64`；`len` 按字节；注释写明开放集合、非法 fail-closed | **一致**（含 A-001 F-001 已收紧的段式规则：禁首/尾/连续中划线） |
| **§3 key** | 非空；≤256；byte < 0x20 或 0x7f 拒绝；UTF-8 允；Set/Delete → `ErrInvalidCacheKey`；**Get 无错误通道**，非法 key 当 miss | `CacheKeyMaxLen = 256`；字节扫描；`CacheView` 包注释写明 Get 非法 key → `(nil, false)` | **一致** |
| **§4 值语义** | nil value / nil policy fail-closed；空 slice 允许；Delete 缺 key 幂等；过期/未写入 Get `(nil, false)` 无错误 | 四条均在接口注释：`ErrInvalidCacheValue` / `ErrInvalidCachePolicy`；Delete idempotent；Get miss 含 expired | **注释一致**；无执行体（R2） |
| **§5 TTL** | `ExpiryPolicy{ExpireAt, Refresh}`；无状态、并发安全；零值 `time.Time` = 永不过期；惰性清理、无后台协程；过期谓词 `expiresAt != zero && !now.Before(expiresAt)` | 接口形状与零值注释一致；文件头写 lazy、无新生命周期。**谓词未编码为 kernel helper** | **接口冻结一致**；谓词仅合同正文（F-002） |
| **§6 容量** | 端口不感知容量；有界义务归 R2 | `cache.go` 无 `maxEntries` / LRU / FIFO 类型 | **一致（未越界实装）** |
| **§7 并发** | 全部 `CacheView` 方法并发安全 | 接口注释声明；无 `-race` 供应商测试（合同归 R2） | **声明一致** |
| **§8 错误面** | 四枚 sentinel，`errors.Is`；未命中非错误；校验先于供应商 | 四枚 `errors.New("kernel: …")`；无 `ErrCacheMiss`。**无 `ValidateCacheSet` / 包装 adapter**；nil value/policy 只在注释 | **sentinel 面一致**；「先于供应商」是义务声明而非可调用端口函数（F-002） |

与先例：身份校验走 ObjectStore 式自由函数（`Valid*` +「实现必须调用」），而非 Mail 的 `MailMessage.Validate()` 值方法。合同 §8 引用 Mail 先例指 **fail-closed 原则**，未强制 Set 也做成值方法。不因此判合同失败。

## 快测覆盖评估（`kernel/cache_test.go`）

D-002 §9 C2 明文范围：**`ValidCacheNamespace` / `ValidCacheKey` 正反例表驱动 + sentinel 存在性**。Get/Set/Delete 行为、过期、驱逐、`-race` 明确归 R2。

本会话 `-v` 实测（非 A-001 自称）：

| 测试 | 子例 | 正/反 | 覆盖合同点 |
|------|------|-------|------------|
| `TestValidCacheNamespace` | **16** | 正 6 / 反 10 | 单字母/数字、段式中划线、满 64、空、大写、首/尾/双中划线、`_` `.` 空白、Unicode、65 字节 |
| `TestValidCacheKey` | **11** | 正 5 / 反 6 | 单字符、`:`/`-`、UTF-8、满 256、空格；空、257、tab、newline、DEL、NUL |
| `TestCacheSentinelErrors` | **1 函数 × 4 sentinel** | — | 非 nil；`fmt.Errorf("wrap: %w", …)` 后 `errors.Is`；消息两两不同 |

**合计可数子例：16 + 11 = 27 表驱动 + 1 个 sentinel 测试（内含 4 条包装链断言）。** A-001 / E-002 / Root 笔记所称「命名空间 17 + key 12 + sentinel 4 = **33 例**」与仓库事实不符（各多计 1 条 ns/key）。测试本身绿，属证据算术错误（F-004），不是「声称有测试但文件没有」。

对照合同戒严点：

- 段式规则负例（尾中划线、连续中划线、首中划线）**有**，锁定 A-001 F-001 的当场收紧。
- sentinel **真实 `%w` 包装链**（A-001 F-002 已改；本会话复核成立）。
- 未测（且 §9 不要求本波）：空值命中 vs miss、拷贝边界、Set 三条件校验顺序、Get 非法 key 当 miss、过期谓词、并发。无供应商则无法做行为断言。

相对 Mail R1：`mail_test.go` 有 `var _ MailSender = …` 编译期端口面守卫；cache 快测无 `var _ Cache` / `CacheView` / `ExpiryPolicy`（F-005）。签名已由接口声明固定，缺口为推荐而非必改。

## 未命中 / 零值 / 拷贝 / Set 校验顺序（专项）

| 主题 | 合同 | 代码 | 漂移？ |
|------|------|------|--------|
| 未命中 | `(nil, false)`，不是 error | `Get` 无 error 返回值；注释 miss = `(nil, false)` | 无 |
| 空值命中 | 非 nil 零长 slice → `(空, true)` | 仅文件头用 ok 区分；`Get` godoc 未写空值命中 | **注释完整度缺口**，非反合同 |
| nil 禁止写入 | `ErrInvalidCacheValue` | sentinel + `Set` godoc | 无（R2 必须调用） |
| 拷贝 | Set 入参拷贝；Get 新拷贝；禁共享内部数组 | 两条方法注释均写明 | **无漂移**；无可执行拷贝 |
| Set 校验顺序 | 端口校验 **先于** 供应商触达；key / nil value / nil policy 各有 sentinel。合同 **未规定三者之间的先后** | 无 wrapper、无 recording stub、无 `ValidateCacheSet`。`CacheView` 注释：key 在进入供应商前校验 | **原则已声明、顺序未冻结、执行未编码**。不构成 §1～§8 签名失败；R2 漏检风险见 F-002 |

## 越界核账（§0 / 红线）

**允许触达（用户本轮硬清单）**：`apps/api/kernel/cache.go`、`apps/api/kernel/cache_test.go`、`docs/workspaces/workspace-026-cache-port/**`。

| 路径 | 状态 | 是否越界 |
|------|------|----------|
| `apps/api/kernel/cache.go`、`cache_test.go` | untracked，仅接口/helper/sentinel/快测 | 否 |
| `GOAL-002-r1-contract-freeze/` | untracked 五件套 + D/E/A | 否 |
| Root `00-meta` / `02-execution` / `E-002-r1-goal-established.md` / `goal-tree.md` / `workspace.md` | 台账同步 | 否（在 canonical 工作区内） |
| `docs/vision/plans/VP-026-cache-port.md` | **已修改**（I-026-001/002/003 → verified + 修订短史一行） | **超出本轮「只触碰 kernel cache + workspace-026」清单**；内容是 P-005 回写，**不是** Charter / Profile / Manifest / `go.mod`（F-006） |
| Charter / Profile 默认集 / 模块矩阵 / Manifest / `go.mod` / `go.sum` | 无变更 | 红线未触 |
| Redis 客户端 / 分布式锁 / 限流 / 消息 / LRU 实装 | `cache.go` 仅注释提及 Redis 接缝；无类型、无 import、无 goroutine | **未越 §0 范围外清单** |

`cache.go` 无 `func (` 方法体，无 Memory 供应商，无策略实装，无容量驱逐代码。R2/R3 责任未被本波偷运进 kernel。

**过程问题（不改 verdict）**：Root `00-meta` 已把判据 #1/#6 打 `[x]`，`goal-tree` 写 Root `1/4`，但 Root frontmatter 仍 `progress: 0/4`，GOAL-002 frontmatter `progress: 0/3` 而正文/goal-tree 为 `1/3`。C3 双审尚未闭合。属台账抢跑与派生进度不一致（F-003），不是端口越界。

## 成果（有证据）

1. **C1 用户裁决落盘**：`D-001-info-adjudication.md`（accepted）；三问均有选项、采纳项、未选方案与理由；非编排器静默代裁。
2. **C2 合同冻结**：`D-002-cache-port-contract.md` v0.1.0；§0 范围外清单与 VP-026 非目标对齐。
3. **端口本体**：`apps/api/kernel/cache.go`——`Cache` / `CacheView` / `ExpiryPolicy` + `ValidCacheNamespace` / `ValidCacheKey` + 四 sentinel；签名与 §1/§2/§3/§5/§8 一致。
4. **合同级快测绿（独立跑测）**：27 表驱动子例 + sentinel `%w` 包装链；`go vet ./kernel/...` 0；`go test ./kernel/... -count=1` PASS。
5. **停机语义（判据 #6 边界）**：I-026-002 选惰性清理；端口无 Start/Stop、无后台协程 → VP-021 排空义务不触发。策略 **形状** 已冻、实装仍归 R2（判据 #2 未在本波主张完成）。
6. **self 已闭环的测试缺陷**：命名空间段式规则与 `%w` 包装——本会话代码与测试均已体现，无需重开 A-001 F-001/F-002。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 判据 #1：端口契约冻结，供应商无关，快测可断言 | **本目标 C2 范围内达成**（C3 待编排器合并本意见后关门） | D-002 + `cache.go` 接口 + 快测绿；无供应商类型泄漏 |
| 判据 #6：惰性清理 → 无新生命周期 | **边界达成** | D-001/D-002 §5；代码无 goroutine |
| API 形态 = []byte + 非泛型 + 类型化封装承诺 R2 | 达成（封装本身不在 R1） | D-001 ①；接口无泛型 |
| 命名空间 scoped 视图 + fail-closed | 达成 | `Namespace` + 段式 `ValidCacheNamespace` |
| ExpiryPolicy 可插拔形状 | 达成（绝对/滑动实装 R2） | 接口 `ExpireAt` / `Refresh` |
| 未改 Profile / Manifest / Charter / 无 Redis 依赖 | 达成 | git 红线路径空；`go.mod` 无 redis |
| I-026-001/002 required 在 C1 前 verified | 达成 | D-001 + 三处台账 |

Root 上过早 `[x]` 判据 #1/#6 **不能**当作本独立审已经放行关门；那是编排器在合并意见后的动作。

## Findings

| # | 级别 | 严重度 | 内容 | 证据 | 状态 | 影响门禁 |
|---|------|--------|------|------|------|----------|
| **F-001** | recommended | med | Redis key 前缀预留在 **D-001 裁决行**（`cache:<ns>:<key>`）与 **D-002 §2 / `cache.go` 文件头**（`<ns>:<key>`）不一致。主裁决（scoped 视图）不受影响，但冻结分母与 P-004 记录不应各写一套。R3 接缝前必须择一写入并回写另一侧。 | `D-001` I-026-003 行；`D-002` §2；`cache.go` L18–19 | open | 不阻断 R1 C3；阻断 R3 接缝文档若未调和 |
| **F-002** | recommended | med | 「端口校验先于供应商」与过期谓词 `expiresAt != zero && !now.Before(expiresAt)` 仅合同/注释。不像 `MailMessage.Validate`，没有 `ValidateCacheSet(key, value, policy) error`（或等价）把 key/nil value/nil policy 收成可测入口，也没有过期谓词 helper。R2 各供应商可能漏检或把 `now.Equal(expiresAt)` 判错。建议 R2 开工前在 kernel 补 helper + 用假供应商证明 **fail-closed 先于存储触达**（三者相对顺序仍可在合同补一句）。 | D-002 §5/§8；`cache.go` 仅接口注释；对比 `mail.go` `Validate` | open | 不阻断 R1 接口冻结；R2 实施前应消化 |
| **F-003** | recommended | low | 台账抢跑与派生进度不一致：① Root 成功标准判据 #1/#6 已 `[x]`，但 GOAL-002 C3 仍待 independent；② `goal-tree` Root `1/4` vs Root `00-meta` frontmatter `0/4`；③ GOAL-002 frontmatter `progress: 0/3` vs 正文/goal-tree `1/3`；④ `02-execution.md` 索引写 E-002「进行中」，E-002 文件 `status: done`。progress 不得当放行证据；编排器响应时应按检查点重算并避免父级判据先于双审闭合。 | Root `00-meta`；GOAL-002 `00-meta`；`goal-tree.md`；`02-execution.md` vs `E-002-contract-frozen.md` | open | 不阻断合同技术面；阻断「宣称 R1 已关门 / Root 1/4」直到 C3 响应完成 |
| **F-004** | informational | low | 「快测 33 例」不成立。实测 16 个 namespace 子例 + 11 个 key 子例 + 1 个 sentinel 测试（4 条 `errors.Is`）。A-001 / E-002 / Root E-002 应改称实测数字。测试仍全绿。 | `go test -v` 本会话；`cache_test.go` | open | 无 |
| **F-005** | informational | low | （a）`Get` 方法 godoc 未写空值命中 `(空 slice, true)`；（b）无 Mail 式 `var _ Cache` 编译期端口面测试。不改变冻结签名。 | `cache.go` Get 注释 vs §1；`mail_test.go` `TestMailSenderPortSurface` | open | 无 |
| **F-006** | informational | low | 工作区外文件 `docs/vision/plans/VP-026-cache-port.md` 被改（信息项 verified + 短史）。属 VP 台账回写，不是 Charter/Profile/go.mod。本轮硬清单未包含该路径；若要坚持「R1 只碰 kernel cache + workspace-026」，应由编排器说明 VP 回写授权。 | `git diff --stat`；VP-026 hunk 仅信息表 + 修订短史 | open | 无（红线未破） |
| **F-007** | informational | low | I-026-002 最晚阶段：VP-026 写 R2 / 仅判据 #3；Goal 写 R1 / #3+#6。已 verified，无到期 required。建议 VP 表补 #6、或注明「R1 提前冻结清理语义」。 | VP-026 信息表；GOAL-001/002 信息表 | open | 无 |

A-001 self 的 F-001（正则过宽）与 F-002（非 `%w` 包装）在当前树中 **已是 fixed 状态**；独立复核未发现回退，不重开。

## 必改项汇总

**开放 required = 0。**

无 high required，无到期且影响本 scope 的 required 信息项。按 P-003 verdict 尺度：**可以在响应 recommended 后无条件放行 GOAL-002 C3 关门**（recommended 不阻断 R1；F-001 跟踪到 R3，F-002 跟踪到 R2）。

## 与 A-001 self 的异同

| 点 | A-001 self | 本意见 independent |
|----|------------|---------------------|
| verdict | pass（0 required） | **pass**（0 required） |
| C1 用户裁决 | 认定成立 | 同意；补充 D-001/D-002 前缀预留不一致（F-001） |
| §1～§8 签名/helper | 「逐节一致」 | 同意签名与 helper；行为语义为注释义务，与 §9 C2 范围匹配 |
| 快测 | 「33 例绿」 | **实测 27 表驱动子例 + sentinel 包装链全绿**；否定 33 这个数字（F-004） |
| `go vet` / `go test ./kernel/...` | 自称绿 | **本会话复跑绿** |
| 越界 | 「仅 cache.* + 本区文档」 | 红线未破；**另见 VP-026 工作区外回写**（F-006） |
| 未写 | — | 缺 Set/过期 helper（F-002）；Root 判据/progress 抢跑（F-003） |

**无冲突必改项**（两边 required 均为 0）。不触发 P-004.2 意见冲突裁决。

## 结论 + 建议给编排器 / 用户的下一步

GOAL-002 R1 **合同冻结在技术面上成立**：用户已书面裁决三条信息项；D-002 将非泛型 `[]byte` 端口、scoped 命名空间、惰性过期接口与四枚 sentinel 冻成可编译的 `kernel` 公共面；独立跑测与越界核账支持「供应商无关契约 + 未偷运 Redis/LRU/锁/限流」。self 的 pass 不被推翻。

**建议 `/govern`：**

1. 将本意见落盘为 **A-002**（`source: independent`），索引 `open required = 0`。
2. 响应 F-001～F-003（recommended）：R3 前对齐 Redis 前缀；R2 前补校验/过期 helper（或在 D-002 写明「校验由每个供应商调用 Valid* + 三条件，kernel 不提供包装」）；纠正 progress / 父级判据勾选，使 GOAL-002 `done` 与 Root 纲领 R1 同步发生，而不是已经提前 `[x]`。
3. F-004 把「33 例」改成实测计数（执行记录勘误即可）。
4. **可以提议 GOAL-002 `status: done`（R1 关门）**——无开放 required、无到期 required 信息项。不得用 `1/3` 或 `1/4` 百分比代替本结论。
5. 下一阶段：按纲领开 **GOAL-003（R2 内存供应商 + 双策略）**；F-002 列为 R2 方案输入。

## 声明

- 本意见 `source: independent`；**不**修改目标 `status` / `progress` / 方案正文 / `goal-tree` 状态列。
- 本会话 **未创建或修改任何仓库文件**（用户硬约束；落盘交编排器）。
- 独立审计不是放行：合并响应与关门由 `/govern` 执行。
- 建议用户下一句：`/govern` 响应 GOAL-002 A-002（independent pass，0 required），闭合 recommended 跟踪项并评估 R1 关门。