你是本仓库的独立审计员（independent auditor）。请加载并遵循 `.grok/skills/audit/SKILL.md`（对应项目 `/audit` 流程）的精神，对 `workspace-026-cache-port` 的 **Root `GOAL-001-cache-port`（VP-026 通用缓存端口）R4 关门**执行**独立交叉审计**（close-out 全量）。

## 硬约束

1. **只输出审计报告文本（Markdown），不得修改或创建任何文件**（P-003；落盘由编排器完成，`source: independent` 保留）。
2. 必须独立复核（可实际运行命令验证），至少执行：
   - `cd apps/api`：`go vet ./...`；`go test ./internal/cache/... ./internal/config/... ./internal/composition/... ./kernel/... -count=1 -race`（cache 带 -race）；如时间允许 `go test ./... -count=1`（全模块，可后台）
   - 仓库根：`git diff --name-only 54fb57e7..HEAD`（越界核账区间）；`Select-String -Path apps/api/go.mod,apps/api/go.sum -Pattern redis`（判据 #4）
   - 抽查：`git status --short`（R4 波工作树仅 owned paths）
3. 报告必须含：verdict（pass / conditional / fail）、scope、八条判据逐条核验、信息门禁核验、阶段审计链核验、越界核账、findings 表（required / recommended / informational）、开放 required 计数、对「可否呈报用户书面关门」的独立意见。

## 核验要点（逐条对照）

- **判据 #1～#8**（对照 `GOAL-005/attachments/r4-evidence-matrix.md` 的映射是否真实，证据是否可核对）：
  1. 端口契约（kernel/cache.go ↔ 合同 D-002 v0.1.1：Cache/CacheView/ExpiryPolicy + Valid* + 4 sentinels + ValidateCacheSet/CacheEntryExpired + 快测）
  2. 双策略 + 可插拔（Absolute/Sliding + 自定义策略样例）
  3. 内存供应商（进程总预算（F-001 用户裁决）/ 全局 FIFO / 惰性清理 / 拷贝边界 / -race 并发 + 跨 ns 预算测试）
  4. Redis 接缝声明（架构短文 §2；go.mod+go.sum redis 0 命中）
  5. 共享约定登记（短文 §3：单一所有者 VP-026 / VP-027 继承 / VP-028 排除 / 登记表 + 变更流程）
  6. 停机语义（惰性清理无新生命周期；无 goroutine）
  7. 边界保持（红线核账：Charter / go.mod / Profile / Manifest / 迁移台账 / mail 零触碰；RT-Q03 保持 gated）
  8. 审计闭合（阶段审计链：GOAL-002 A-002 pass · GOAL-003 A-002 conditional→F-001 用户裁决 fixed · GOAL-004 A-002 pass；各阶段开放 required = 0）
- **信息门禁**：I-026-001/002/003/004 全部 verified（用户书面裁决/确认留痕）。
- **契约面一致性**：D-002 v0.1.1 ↔ kernel/cache.go ↔ internal/cache 实现 ↔ 架构短文四点一致。
- **回归证据**：独立复跑核对 `go vet` / 测试（含 `-race`）/ redis 搜索。
- **R4 波工作树**：仅允许集（apps/api/internal/composition/* 不应再变——R4 无代码改动；docs/workspaces/workspace-026-cache-port/**、docs/architecture/cache-redis-seam-and-track.md、docs/vision/**、docs/README.md）。

## 审计上下文（先读）

- `docs/workspaces/workspace-026-cache-port/workspace.md`
- `docs/workspaces/workspace-026-cache-port/GOAL-001-cache-port/00-meta.md`（Root 纲领 + 成功标准 + 信息台账）
- `docs/workspaces/workspace-026-cache-port/GOAL-001-cache-port/03-audit/A-001-root-closeout-self.md`（对照 self）
- `docs/workspaces/workspace-026-cache-port/GOAL-005-r4-evidence-closeout/00-meta.md` + `01-decision/D-001-r4-closeout-design.md` + `attachments/r4-evidence-matrix.md`
- 阶段产物：GOAL-002 `D-002`（合同）+ `A-002`（独立）；GOAL-003 `D-001`（含 F-001 勘误）+ `A-002`（conditional 原文）+ `A-003`；GOAL-004 `D-001` + `A-002` + attachments（mail 评估）
- 代码：`apps/api/kernel/cache.go`、`apps/api/internal/cache/{memory,policy,typed}.go`、`apps/api/internal/composition/composition.go`（fx 段）、`apps/api/internal/mail/runtime.go`（L116～L229）
- 规则：`docs/architecture/independent-audit-execution.md`、`docs/architecture/principles.md`（P-003/P-005）、`docs/vision/alignment.md`

## 输出

直接输出最终审计报告文本（Markdown），不要输出工作过程叙述。