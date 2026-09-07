# A-002 · grok build 独立审计原始输出（2026-09-01 · 原样收录）

> 来源：本地 grok build（grok-4.6 · reasoning high · headless 单轮），`source: independent`。编排器按 P-003 誊入，未修改意见内容；报告正文自「# A-002」起为 grok 原样输出。grok 当场独立复跑：`go vet ./...`（0）、`go test ./internal/cache/... -count=1 -race`（ok）、`go test ./internal/config/... ./internal/composition/...`（ok）、`git status` / `git diff`（越界核账）。grok 未创建或修改任何文件。

---

# A-002 · R2 内存供应商关门独立交叉审计（2026-09-01）

- **source**：`independent`
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：`close-out` / execution-facts
- **scope**：`workspace-026-cache-port` · `GOAL-003-r2-memory-provider` 全量——C1 方案冻结、C2 实施（`internal/cache` + 配置键 + 组合根）、判据 #2/#3/#6、合同 D-002 v0.1.1 逐条一致性、信息门禁、越界核账；对照 A-001 self
- **verdict**：**conditional**
- **开放 required 计数**：**1**
- **工作区**：`id=workspace-026-cache-port` · `root_goal=GOAL-001-cache-port` · `canonical_scope=docs/workspaces/workspace-026-cache-port/` · `shared_materials_catalog=none`（无共享资料引用可误当事实）
- **落盘**：本会话按用户硬约束**未写入任何文件**。编排器应将本意见落盘为 `03-audit/A-002-r2-impl-closeout-independent.md` 并更新 `03-audit.md` 索引。

## 范围与区间

独立复核（本会话实测，非转述 A-001）：

| 命令 | 结果 |
|------|------|
| `cd apps/api` → `go vet ./...` | 退出码 0，无输出 |
| `go test ./internal/cache/... -count=1 -race` | `ok` · 2.268s（无 data race） |
| `go test ./internal/config/... ./internal/composition/... -count=1` | `ok` config 2.921s · composition 27.378s |
| 仓库根 `git status --short` / `git diff --stat` | 见「越界核账」 |

只读扫描：workspace.md、Root 00-meta、GOAL-003 五件套、D-001、E-002、A-001、R1 D-002 v0.1.1、`internal/cache/{policy,memory,typed}.go` + 测试、`config.go` Cache 段、`composition.go` `newCache`、`cache_wiring_test.go`、`kernel/cache.go` 合同面。未读取其他工作区。

## 信息门禁核验（P-005）

| ID | 级别 | 最晚阶段 | 状态 | 对本 scope |
|----|------|----------|------|------------|
| I-026-001 / 002 / 003 | required（Root） | R1 | verified（GOAL-002 D-001） | 不阻断 R2 |
| I-026-004 | **non-blocking** | **R3** | 待确认 | 不阻断 R2 关门 |

R2 **无新 required 信息项**。无到期且影响本 scope 的 required 信息项；无 `accepted-residual` 缺口。共享资料目录 `none`。

## 方案冻结（C1 / P-004）

| 项 | 判定 | 证据 |
|----|------|------|
| 驱逐策略用户裁决 FIFO | **pass** | E-001：「用户**裁决 FIFO**（P-004）」；D-001 `status: accepted` 冻结 FIFO（container/list + map、覆盖写保位、过期 Set 重插） |
| 未选方案留痕 | **pass** | D-001：近似 LRU / 逐 ns 互斥锁 / Set 全量清扫 / 驱逐仅计活动条目 / 配置容错回落默认 —— 均有理由 |
| maxEntries 义务 / Typed / 配置键 / 组合根 / 审计模式 cross | **pass（有一处未锁定，见 F-001）** | D-001 表已写；C3 cross = A-001 self + 本 A-002，与 D-001 一致 |

本会话未见原始用户对话，以 `accepted` 决策台账为 P-004 书面证据。足够。

## 成果（有证据）

1. `apps/api/internal/cache/` 新包：`policy.go`（`AbsoluteExpiry` / `SlidingExpiry`）· `memory.go`（`Memory` + FIFO + 惰性清理 + 全局 `sync.Mutex`）· `typed.go`（`Typed[T]` + `JSONCodec` + 可注入 `Codec[T]`）+ 18 个测试函数（14 memory + 4 typed）。
2. 配置键：`DefaultCacheMaxEntries = 10000`；YAML `cache.max_entries`；env `CACHE_MAX_ENTRIES` 严格 `Atoi`（非法/`<=0` → `LoadError`）；Load 网 `CacheMaxEntries <= 0`；`ValidateProd` 先包装 `LoadError` 再拦 `<0`；`config.default.yaml` / `configs/config.yaml` / `.env.example` 已文档化（canonical-env 扫描 `os.Getenv("CACHE_MAX_ENTRIES")`）。
3. 组合根：`newCache` 零值回落默认、负值 fail-closed、`cache_wiring_test.go` 3 子例；`newMuxWithExtraProviders` 调用 `newCache`。
4. 实现期两处修正已在 E-002 留痕（滑动用例时间线；`copyBytes` 保留空值非 nil）。代码侧可复核。
5. 本会话：`go vet` 0；cache `-race` 绿；config + composition 绿。

## 合同-实施逐条一致性（D-002 v0.1.1 ↔ 实现）

| # | 合同 / D-001 义务 | 实现 | 测试 | 判定 |
|---|-------------------|------|------|------|
| 1 | `Memory` 实现 `kernel.Cache` | `var _ kernel.Cache = (*Memory)(nil)` | compile-time 断言 | **一致** |
| 2 | Set 先 `ValidateCacheSet`（key→value→policy）再触达存储 | `Set` 在 `mu.Lock` 之前调用 | `TestMemoryFailClosedValidation` | **一致** |
| 3 | 过期判定用 `CacheEntryExpired` | Get/Set 均调用 | 绝对/滑动/惰性清扫专测 | **一致** |
| 4 | 拷贝边界 | `copyBytes`（避免 `append(nil, empty)` 坍缩） | `TestMemoryCopySemantics` + 空值命中非 nil | **一致** |
| 5 | 空值命中 / nil fail-closed | 匹配 | basics + fail-closed | **一致** |
| 6 | Delete 幂等；Get 非法 key 当 miss | 匹配 | fail-closed + basics | **一致** |
| 7 | Namespace 校验 fail-closed | `ValidCacheNamespace` 先于建 space | `TestMemoryNamespaceValidation` | **一致** |
| 8 | 绝对不刷新 / 滑动命中刷新 | 匹配 | 两专测 | **一致** |
| 9 | TTL<=0 永不过期（零值 `time.Time`） | `ExpireAt` 返回 `time.Time{}` | `TestMemoryZeroTTLNeverExpires`（仅 Absolute） | **一致**（滑动零窗无专测，F-003） |
| 10 | 策略无状态、并发安全 | 值类型、无可变字段 | `-race` 并发测 | **一致** |
| 11 | 自定义策略可插拔（判据 #2） | `nextMidnightPolicy` 经接口注入 | `TestMemoryCustomPolicyPluggable` | **一致** |
| 12 | 惰性清理仅读写路径；无 goroutine | 生产代码无 `go` 关键字 | `TestMemoryLazyCleanupFreesCapacity` | **一致** |
| 13 | FIFO 覆盖写保位 / 过期 Set 重插 | 活条目 in-place；过期重插 | `TestMemoryEvictionOverwriteKeepsPosition` | **一致**（policy 换实例丢位，F-004） |
| 14 | 任一 Set 后**总条目** ≤ `maxEntries`（含过期未清扫） | **`for len(v.space.entries) >= maxEntries`** —— 计数域是**当前命名空间**，不是 `Memory` 全量 | 三则有界测试均单 ns | **不一致** → **F-001 required** |
| 15 | `Typed[T]` JSON 默认 + 注入 codec；解码错误不伪装 miss | 匹配 | 4 则 typed 测试 | **一致** |
| 16 | 配置键 + 非法值 fail-closed | 实现齐全 | wiring 3 例 + canonical-env；**无** LoadError 专测 | **实现一致 / 测试偏薄** → F-003 |
| 17 | `newCache` 零值回落 · 负值 fail-closed · 单一实例 | 行为匹配；`_ = cachePort` 并不能阻止 GC | wiring 3 例 | **构造一致**；holder 名不副实 → F-002 |
| 18 | `-race` 并发 Get/Set/Delete | 8 goroutine × 200 op，双 ns | 绿；未断言事后条目数 | **无竞争**；有界断言缺口 → F-003 |
| 19 | 不改端口合同 / 无 Redis / 不改 Charter | 无 content diff；无 redis 引用 | git + grep | **合同面未越界**；gofmt 脏文件 → F-005 |

## 对照成功标准

| # | 标准 | 状态 | 证据 |
|---|------|------|------|
| 1 | `Memory` 实现端口 + §1–§8 行为 | **达成** | 上表 1–7 |
| 2 | 双策略 + 可插拔（判据 #2） | **达成** | 上表 8–11 |
| 3 | 有界 + FIFO + 惰性清理（判据 #3/#6） | **部分** | FIFO/惰性达成；**有界计数域未锁定且与注释矛盾（F-001）** |
| 4 | `-race` 并发安全 | **达成** | cache 包 `-race` ok |
| 5 | `Typed[T]` | **达成** | typed_test 4 例 |
| 6 | 配置键 fail-closed | **达成（实现）** | 代码路径完整；Load 专测缺失（F-003） |
| 7 | 组合根接线 | **达成（构造）** | wiring 3 例；holder 不保活（F-002） |
| 8 | 未越界（合同面） | **达成** | kernel/go.mod/Charter 无 diff；Redis 未引入 |

## 测试覆盖评估

**充分**：Set/Get/Delete 基础与 sentinel 顺序、拷贝与空值、ns 隔离与校验、绝对/滑动过期、自定义策略、FIFO 保位、惰性清扫释容、单 ns 过期项仍占预算、Typed JSON/自定义 codec/解码非 miss、`NewMemory` 非正 fail-closed、wiring 三路径、`-race`。

**缺口（影响关门的只有 F-001）**：

- 跨命名空间有界（进程总预算 vs 每 ns 预算）——**无测试，语义未冻结**。
- `CACHE_MAX_ENTRIES` / YAML `max_entries` 非法值 LoadError（仓库对 `MAIL_SMTP_PORT` 有对等专测）。
- `SlidingExpiry{Window: 0}` 永不过期。
- 并发结束后 Σ entries ≤ maxEntries。
- 同一 ns 两个 `CacheView` 共享底层 space。

## 越界核账

**R2 语义改动（允许集）**：

- `?? apps/api/internal/cache/`
- `M` `config.go` + `config.default.yaml` + `configs/config.yaml` + `.env.example`（Cache 键；另有 gofmt 对齐）
- `M` `composition.go`（`newCache` + import；gofmt 重排 `kernel` 组）
- `?? cache_wiring_test.go`
- `docs/workspaces/workspace-026-cache-port/**`（GOAL-003 五件套 + Root/workspace/goal-tree 进度回写）

**禁止面**：

| 面 | 结果 |
|----|------|
| `apps/api/kernel/cache.go`（及整个 `kernel/`） | **无 content diff** |
| `apps/api/go.mod` / `go.sum` | **无 diff**；全树无 redis 客户端 |
| Charter | **无 diff** |
| Profile 默认集 | **无 kernel profile 改动** |
| Manifest 装配语义 | `manifest.go` status=M 但 **空 content diff**（CRLF）；`manifest_test.go` 仅删空行 |

**工作树噪音**：`git diff --name-only` 另有约 60 个非允许集文件（handler/auth/mail/store/composition 测试等），抽检为 **gofmt import 分组 / 复合字面量对齐 / 空行**，非 Redis、非端口合同、非 Profile。A-001「git status 仅允许集」**当前不成立**。R2 提交不得 `git add -A`（P-002）。

## Findings

| # | 级别 | 严重度 | 内容 | 证据 | 状态 |
|---|------|--------|------|------|------|
| **F-001** | **required** | **med** | **maxEntries 计数域未冻结，且实现与「供应商总条目」注释不一致。** `Set` 以 `len(v.space.entries)` 做 FIFO 驱逐，故 N 个命名空间最多约 **N × maxEntries**。`Memory` 注释、`Config.CacheMaxEntries` 注释、YAML/`.env.example` 均写「provider TOTAL ≤ maxEntries」；D-001「总条目数（含未清扫过期项）」字面可作「含过期」也可作「全实例」。有界三测均为单 ns。A-001 将成功标准 3 标「达成」——跨 ns 证据不足。判据 #3 在计数域选定并被测试锁定前，不得无条件关门。 | `memory.go` 驱逐循环；`config.go` Cache 注释；有界测试；A-001 标准 3 | **open** |
| **F-002** | recommended | low | 组合根 `_ = cachePort` **不能**保活对象（Go 中 `_ = x` 只消未使用变量）。注释「blank assign keeps the reference live」不实。构造 fail-closed 仍有效；进程内无长期引用。同意 A-001 F-002：R3 必须把实例挂到长生命周期结构并去掉 holder。 | `composition.go` | **open**（跟踪 R3） |
| **F-003** | recommended | low | 测试缺口：(1) `CACHE_MAX_ENTRIES`/`cache.max_entries` 非法值 LoadError（对标 `config_mail_test.go` 的 `MAIL_SMTP_PORT`）；(2) `SlidingExpiry{Window:0}`；(3) `-race` 后条目数 ≤ maxEntries。实现已读过，属覆盖而非行为缺失。 | `config.go`；`policy.go`；`TestMemoryConcurrentAccess` | **open** |
| **F-004** | recommended | low | 活条目 `e.policy != policy`（接口值比较）时按重插处理，FIFO 位丢失。D-001 只冻结「覆盖写保位 / 过期 Set 重插」，未写 policy 更换。应补测或在 D-001 勘误一句，避免 R4 验收歧义。 | `memory.go` 活条目分支 | **open** |
| **F-005** | recommended | low | 工作树大量非 R2 语义的 gofmt 脏文件；A-001 越界陈述过时。提交只暂存允许路径。另：`copyBytes` 注释称 nil 入参保持 nil，实际 `make([]byte, 0)` 为非 nil（Set 路径不会传入 nil，无合同破坏）。 | `git diff --name-only`；`memory.go` `copyBytes` | **open** |
| **F-006** | informational | low | 台账卫生：`00-meta.md` frontmatter `progress: 0/3` 与正文/goal-tree `1/3` 不一致；`02-execution.md` 索引 E-002 为「进行中」、文件 `status: done`。不阻断。A-001 F-001（实现期两处修正）已闭合，不重开。 | 台账文件 | **open**（卫生） |

**A-001 F-002**（holder → R3）：本意见以 F-002 强化复述，不另开号。

## 必改项汇总

**开放 required = 1**

1. **F-001**：用户/编排器必须二选一并留痕 + 补锁定测试后，方可把判据 #3 视为达成：
   - **建议（与「内存有界最严格」+ 进程级 `cache.max_entries` 一致）**：进程总条目（含过期未清扫、跨 ns）≤ `maxEntries`；跨 space 计数 + 可测的跨 ns 驱逐；测试：两 ns 合计不可超过预算。
   - **或者**勘误 D-001 / `Memory` / `Config` / YAML 注释为**每命名空间独立预算**，并加测试：两 ns 可**各**持 `maxEntries` 条。

未选、未测、不关门。此条属 P-004 裁决点（语义分叉），独立审计给建议但不代裁。

无 high required。无到期 required 信息项。

## 与既有意见的异同（A-001 self）

| 点 | A-001 self | 本意见 |
|----|------------|--------|
| C1 FIFO 用户裁决 | pass | **同意** |
| 双策略 / Typed / 拷贝 / ValidateCacheSet / CacheEntryExpired / `-race` / 无 Redis / 合同面未改 | pass | **同意**（本会话复跑测试） |
| 判据 #3 有界「总条目」达成 | pass | **不同意** → F-001 required |
| 越界「git status 仅允许集」 | pass | **当前不成立**（gofmt 脏树）→ F-005 recommended，合同面仍干净 |
| holder | recommended → R3 | **同意并强化**（`_` 不保活） |
| verdict | pass（0 required） | **conditional**（1 required） |

无「一要一否」式 finding 冲突需用户立刻在两条 required 之间裁断；差异是 self 漏掉计数域。F-001 闭合路径（实现全局 vs 勘误 per-ns）才是 P-004。

## 结论 + 建议给编排器 / 用户的下一步

C1 方案冻结与 C2 主体实施（双策略、可插拔样例、拷贝、校验顺序、惰性清理、FIFO 保位、Typed、配置键实现、`-race`、合同面未越界）**经独立复跑与读码核对，大体属实**。不能无条件 pass 的原因是 **判据 #3 的有界语义在「进程总预算 vs 每命名空间预算」上未冻结**，且实现与多处「TOTAL」注释矛盾。

**verdict = conditional。开放 required = 1（F-001）。**

建议 `/govern`：

1. 展示 F-001 两条闭合路径 + 本意见建议（进程总预算），等用户书面裁决并写入 D-00N 或 D-001 勘误。
2. 按裁决改代码或勘误注释，并补跨 ns 测试；可选顺手 F-003 LoadError 专测。
3. 合并响应 A-001 + A-002 全部 findings；F-001 `fixed` 前**不得**将 GOAL-003 标 `done`、不得把 Root R2 标完成。
4. 提交只暂存 R2 允许路径，禁止 `git add -A`。

## 声明

本意见 `source: independent`，不修改 `status` / `progress` / 方案正文 / goal-tree。响应与落盘由 `/govern` 处理。进度数字不得作为 finding 闭合或放行依据。