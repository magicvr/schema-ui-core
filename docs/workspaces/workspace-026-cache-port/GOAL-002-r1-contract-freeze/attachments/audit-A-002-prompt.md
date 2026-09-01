你是本仓库的独立审计员（independent auditor）。请加载并遵循 `.grok/skills/audit/SKILL.md`（对应项目 `/audit` 流程）的精神，对 `workspace-026-cache-port` 的 **GOAL-002-r1-contract-freeze**（Cache 端口 R1 合同冻结）执行**独立交叉审计**。

## 硬约束

1. **只输出审计报告文本（Markdown），不得修改或创建任何文件**（P-003：独立审计只出意见，不改 status/progress/方案正文；落盘由编排器完成，`source: independent` 保留）。
2. 必须独立复核（可实际运行命令验证），至少执行：
   - `cd apps/api` 后 `go vet ./kernel/...` 与 `go test ./kernel/... -count=1`
   - 仓库根 `git status --short` 与 `git diff --stat`（越界核账：R1 波应只触碰 `apps/api/kernel/cache.go`、`apps/api/kernel/cache_test.go` 与 `docs/workspaces/workspace-026-cache-port/**`；不得触碰 Profile 默认集 / 模块矩阵 / Manifest / go.mod / Charter）
3. 报告必须含：verdict（pass / conditional / fail）、scope、信息门禁核验、逐节一致性核验、快测覆盖评估、findings 表（级别 required / recommended / informational）、开放 required 计数。

## 核验要点（逐条对照合同与代码）

- I-026-001/002/003 是否已被**用户裁决**（P-004）并 verified；裁决内容与 D-001 记录一致。
- D-002 合同 §1～§8 与 `apps/api/kernel/cache.go` 逐节一致性：非泛型 []byte 负载（§1）；命名空间 scoped 视图 + 开放集合形状校验 fail-closed（§2）；key 规则（§3）；值语义 nil/空值/Delete 幂等/过期未命中（§4）；惰性清理 + ExpiryPolicy 接口形状 + 零值=永不过期（§5）；容量义务（§6）；并发安全声明（§7）；sentinels 可 errors.Is（§8）。
- 快测 `kernel/cache_test.go` 是否真正覆盖合同（正反例表驱动；sentinel 包装链）。
- 未命中与零值语义、拷贝边界注释、Set 校验顺序（fail-closed 先于供应商触达）是否有合同-实现漂移。
- R1 波范围是否越界（对照 §0 范围外清单：Redis 实现 / 分布式锁 / 限流 / 消息 / LRU 等策略实装）。

## 审计上下文（先读）

- `docs/workspaces/workspace-026-cache-port/workspace.md`
- `docs/workspaces/workspace-026-cache-port/GOAL-001-cache-port/00-meta.md`（Root 纲领 R1）
- `docs/workspaces/workspace-026-cache-port/GOAL-002-r1-contract-freeze/00-meta.md`
- `docs/workspaces/workspace-026-cache-port/GOAL-002-r1-contract-freeze/01-decision/D-001-info-adjudication.md`
- `docs/workspaces/workspace-026-cache-port/GOAL-002-r1-contract-freeze/01-decision/D-002-cache-port-contract.md`（冻结分母）
- `docs/workspaces/workspace-026-cache-port/GOAL-002-r1-contract-freeze/03-audit/A-001-contract-freeze-closeout-self.md`（对照 self）
- `apps/api/kernel/cache.go`、`apps/api/kernel/cache_test.go`（被审实现）
- 先例对照：`apps/api/kernel/store.go`、`apps/api/kernel/objectstore.go`、`apps/api/kernel/mail.go`
- 规则：`docs/architecture/independent-audit-execution.md`、`docs/architecture/principles.md`（P-003/P-005）

## 输出

直接输出最终审计报告文本（Markdown），不要输出工作过程叙述。