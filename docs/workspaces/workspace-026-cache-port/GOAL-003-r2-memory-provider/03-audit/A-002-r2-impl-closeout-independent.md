---
id: GOAL-003-r2-memory-provider
doc: audit-entry
record_id: A-002
status: recorded
parent: GOAL-003-r2-memory-provider
created: 2026-09-01
updated: 2026-09-01
version: 0.1.0
---

# A-002 · R2 内存供应商关门独立交叉审计（grok build · independent）

> 誊入说明：本条由编排器自本地 grok build（grok-4.6 · reasoning high）headless 会话原样誊入（2026-09-01；grok 按指令只出报告文本、未落盘）。grok 当场独立复跑：`go vet ./...` 0、`go test ./internal/cache/... -count=1 -race` ok、`go test ./internal/config/... ./internal/composition/...` ok、`git status` / `git diff` 越界核账。原始输出见 [attachments/audit-A-002-grok-output.md](attachments/audit-A-002-grok-output.md)。

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · headless 单轮）
- **date**：2026-09-01
- **scope**：GOAL-003-r2-memory-provider 全量（C1 方案冻结 · C2 实施 · 判据 #2/#3/#6 · 合同 D-002 v0.1.1 逐条 · 越界核账）
- **verdict**：**conditional**（开放 required = **1**：F-001）
- **状态**：F-001 用户裁决 + F-002～F-006 处置后按 A-003 闭合

## 摘要

- **信息门禁**：R2 无新 required 信息项；I-026-004（R3）不阻断。C1 FIFO 用户裁决 pass；未选方案留痕 pass。
- **合同-实施逐条**：19 条核查中 18 条**一致**（ValidateCacheSet 先于触达 / CacheEntryExpired / 拷贝与空值语义 / Delete 幂等 / ns fail-closed / 绝对不刷新 / 滑动刷新 / TTL<=0 永不过期 / 自定义策略可插拔 / 惰性清理无 goroutine / FIFO 保位 / Typed 解码非 miss / 配置键 fail-closed 实现 / 零值回落 + 负值 fail-closed / `-race` 无竞争 / 合同面未越界）。
- **不一致 1 条 → F-001 required**：maxEntries **计数域未冻结**——实现按每命名空间计数（`len(v.space.entries)`），N 个 ns 最多 N×maxEntries；而 `Memory`/`Config`/YAML 注释均写「provider TOTAL ≤ maxEntries」；有界测试全为单 ns。判据 #3 在计数域选定并被测试锁定前不得无条件关门。
- **越界**：`kernel/cache.go`、`go.mod`、Charter、Profile 无 diff；无 Redis；红线未破。**工作树卫生问题**：约 60 个非允许集文件为 gofmt 噪音（编排器 `gofmt -w internal` 误扫），A-001「git status 仅允许集」陈述过时（F-005）——已由编排器恢复并纳入 A-003 处置。
- **独立复跑**：vet / cache `-race` / config / composition 全绿。

## Findings（6 条 · required 1 + recommended 4 + informational 1）

| # | 级别 | 内容摘要 | 处置（见 A-003） |
|---|------|----------|------------------|
| F-001 | **required** | maxEntries 计数域未冻结：每 ns 计数（N×maxEntries）vs 注释「provider 总条目」；跨 ns 证据缺失；判据 #3 关门受阻。两条闭合路径：**进程总预算**（审计建议）或每 ns 预算（勘误注释） | **fixed**（2026-09-01 **用户裁决：进程总预算**；`Memory` 重构为全局计数 + 全局 FIFO + 跨 ns 驱逐；新测试 `TestMemoryGlobalBudgetAcrossNamespaces` / `TestMemoryGlobalFIFOInterleave` / `TestMemoryConcurrentBudgetBound`；D-001 勘误 v0.1.1） |
| F-002 | recommended | 组合根 `_ = cachePort` 不能保活实例（`_ = x` 只是消未使用变量）；注释「keeps the reference live」不实 | **fixed-recording**（注释改为如实描述 + R3 义务：挂长生命周期结构并移除 blank assign；D-001 勘误 §3 跟踪 R3） |
| F-003 | recommended | 测试缺口：CACHE_MAX_ENTRIES / cache.max_entries 非法值 LoadError 专测；SlidingExpiry{Window:0} 永不过期；`-race` 后条目数 ≤ maxEntries | **fixed**（`config_cache_test.go` 6 用例；`TestMemoryZeroTTLNeverExpires` 扩展 Sliding 零窗；`TestMemoryConcurrentBudgetBound` Len≤25） |
| F-004 | recommended | 活条目 `e.policy != policy` 重插致 FIFO 位丢失，且接口比较对不可比较类型有 panic 风险；D-001 未冻结该语义 | **fixed**（实现改为：活条目覆盖写无论策略是否更换均保位，仅过期条目重插；**供应商不做策略接口比较**；D-001 勘误 §2；测试 `TestMemoryLiveOverwriteKeepsPositionAcrossPolicyChange`） |
| F-005 | recommended | 工作树 gofmt 噪音（约 60 文件）；A-001 越界陈述过时；`copyBytes` 注释称 nil 保持 nil 不实（make 产出非 nil 空） | **fixed**（编排器恢复全部非允许集文件，工作树仅剩 owned paths；`copyBytes` 注释勘误） |
| F-006 | informational | 台账卫生：GOAL-003 frontmatter progress 0/3 vs 正文 1/3；02-execution 索引 E-002 状态漂移 | **fixed**（A-003 响应后统一台账） |

## 关键结论（grok 原话要点）

- 「不能无条件 pass 的原因是判据 #3 的有界语义在『进程总预算 vs 每命名空间预算』上未冻结，且实现与多处『TOTAL』注释矛盾。」
- 「此条属 P-004 裁决点（语义分叉），独立审计给建议但不代裁。」
- 「F-001 `fixed` 前不得将 GOAL-003 标 `done`、不得把 Root R2 标完成。」
- 「C1 方案冻结与 C2 主体实施（双策略、可插拔样例、拷贝、校验顺序、惰性清理、FIFO 保位、Typed、配置键实现、`-race`、合同面未越界）经独立复跑与读码核对，大体属实。」

## 链接

- 原始输出全文：[attachments/audit-A-002-grok-output.md](attachments/audit-A-002-grok-output.md)
- 编排器合并响应：[A-003-response-to-a002.md](A-003-response-to-a002.md)
- 对照 self：[A-001-r2-impl-closeout-self.md](A-001-r2-impl-closeout-self.md)