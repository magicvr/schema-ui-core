# R4 证据矩阵 · VP-026 八条方向级退出判据 ↔ 阶段证据（2026-09-01）

> 责任人：workspace-026 R4（GOAL-005）；区间：激活规划 `54fb57e7` → R3 关门 `c4284450`。逐条 verified 项均给出可核对证据；红线核账覆盖区间内全部 82 个路径。

## 判据逐条映射

| # | 判据 | verified | 证据（阶段 / 产物） |
|---|------|----------|---------------------|
| 1 | **端口契约冻结**：Get/Set/Delete + TTL + 命名空间 + 并发安全；供应商无关、快测可断言 | **✓** | R1（GOAL-002）：合同 `D-002` v0.1.1；`apps/api/kernel/cache.go`（Cache/CacheView/ExpiryPolicy + Valid* + 4 sentinels + ValidateCacheSet/CacheEntryExpired）；快测 = `kernel/cache_test.go` 5 父测试 / 40 表驱动子例 + sentinel `%w` 链 + 编译期端口面断言；A-001 self + A-002 grok independent 双审 pass |
| 2 | **双策略 + 可插拔**：绝对/滑动 + 策略接口（含自定义策略测试样例） | **✓** | R2（GOAL-003）：`internal/cache/policy.go`（Absolute/Sliding 语义专测：绝对不刷新 / 滑动刷新 / 零窗永不过期）；自定义 `nextMidnightPolicy` 样例（`memory_test.go`）；A-002 review 判据 #2 「达成」 |
| 3 | **内存供应商可用**：有界 + TTL 清理 + 驱逐 + 并发边界测试 | **✓** | R2（GOAL-003）：`internal/cache/memory.go`（**进程总预算**（用户裁决）· 全局 FIFO · 惰性清理 · 拷贝边界 · 全局互斥）；23 父测试含 `-race`（TestMemoryConcurrentAccess / ConcurrentBudgetBound）+ 跨 ns 预算（GlobalBudgetAcrossNamespaces / GlobalFIFOInterleave）+ 驱逐（FIFOEvictionBound / EvictionOverwriteKeepsPosition / ExpiredEntriesStillBoundTheTotal / LazyCleanupFreesCapacity） |
| 4 | **Redis 接缝声明落盘**：端口不变 / 连接管理约定 / key 前缀与命名空间约定；`go.mod` 无 Redis 客户端 | **✓** | R3（GOAL-004）：架构短文 `docs/architecture/cache-redis-seam-and-track.md` §2（端口不变 / key `<ns>:<key>` / TTL 映射 / 连接管理 + PING fail-closed / 无客户端依赖）；`go.mod`+`go.sum` redis **0 命中**（grok 独立复核确认） |
| 5 | **共享约定登记**：Redis 轨道约定（VP-026/027）单一所有者文档落地（本区 owner） | **✓** | R3（GOAL-004）：短文 §3（单一所有者 VP-026 · VP-027 激活继承 · VP-028 排除 · 命名空间登记表 + owner 义务 · 变更流程 + 修订史 1.0.0） |
| 6 | **停机语义**：惰性清理避开新生命周期（VP-021 义务不触发） | **✓** | R1/R2：I-026-002 用户裁决惰性清理（GOAL-002 D-001 ①）；合同 §5；`memory.go` 无 goroutine / 无 Hooks（grok A-002 判据 #6 复核） |
| 7 | **边界保持**：未改 Charter / Profile 默认集 / Manifest；未预制 Redis（不消耗 RT-Q03 trigger）；未重开历史 VP | **✓** | 红线核账（下表）：Charter / `go.mod`+`go.sum` / Profile / Manifest / 迁移台账零触碰；区间代码 = 仅 kernel 端口 + internal/cache + config 键 + composition 接线；RT-Q03 保持 gated（VP-026 修订史 + 短文 §1） |
| 8 | **审计闭合**：开放 required finding = 0（或已合法闭合） | **✓**（基线） | 阶段审计：R1 A-003（9 findings 全处置）· R2 A-003（F-001 用户裁决 → fixed；8 findings 全处置）· R3 A-003（5 findings 全处置）；各阶段开放 required = 0；R4 关门双审（Root A-001 self + A-002 grok independent）见 GOAL-005 C2/C3 |

## 红线越界核账（`git diff --name-only 54fb57e7..HEAD` · 82 路径）

| 面 | 路径 | 结论 |
|----|------|------|
| Charter | `docs/vision/charter.md` | **零触碰** |
| 依赖锁 | `apps/api/go.mod` / `go.sum` | **零触碰**（redis 0 命中） |
| Profile 默认集 / 模块矩阵 | `apps/api/kernel/profile.go` 等 kernel 面（除 cache 端口） | **零触碰**（kernel 区仅 `cache.go` + `cache_test.go`） |
| Manifest 装配语义 | `internal/manifest/`、`modules/*/manifest/` | **零触碰** |
| 迁移台账 | `internal/store/migrate*`、`internal/migration/` | **零触碰** |
| 各异构模块 | `modules/**`（除 workspace 文档） | **零触碰** |
| mail（I-026-004 承诺零漂移） | `internal/mail/` | **零触碰**（git 空 diff · grok 确认） |
| 波次触碰面 | kernel/cache.* · internal/cache/ · internal/config（Cache 键）+ config.default.yaml · internal/composition（newCache 接线）+ 测试 · configs YAML ×2 · docs/workspaces/workspace-026-cache-port/** · docs/architecture/cache-redis-seam-and-track.md · docs/vision/**（VP-026 激活/关门台账）· docs/README.md | 全部属允许集 |

## 回归证据基线（R3 波 2026-09-01 实测 + 阶段复跑）

| 面 | 结果 |
|----|------|
| `go vet ./...`（apps/api） | 0（阶段复跑 ×4） |
| `go test ./... -count=1` 全模块 | **exit 0**（R3 波：无 FAIL，50 包；R2 波同样 exit 0） |
| `go test ./internal/cache/... -count=1 -race` | ok（含并发/预算/跨 ns 测试） |
| `go test ./internal/composition/... ./internal/config/...` | ok（fx 图真实拉起：composition_test / shutdown_drain / postgres_startup 等） |
| 配置键 fail-closed | `config_cache_test.go` 6 子例（默认 / YAML / env / 非法 env / 非法 YAML / ValidateProd）+ canonical-env 门禁 |
| 无 Redis 依赖 | `go.mod`+`go.sum` 大小写不敏感搜索 0 命中 |

## 信息门禁

I-026-001/002/003（R1 用户裁决）· I-026-004（R3 用户确认不迁移）——全部 **verified**；无 deferred required；无 open required 信息项。

## 结论

八条判据证据齐备；红线零触碰；阶段审计全部闭合（开放 required = 0）。**待**：Root A-001 self → A-002 grok build independent → 合并响应 → VRev-061 → 用户书面确认关门（GOAL-005 C2/C3）。