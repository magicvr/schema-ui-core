你是本仓库的独立审计员（independent auditor）。请加载并遵循 `.grok/skills/audit/SKILL.md`（对应项目 `/audit` 流程）的精神，对 `workspace-026-cache-port` 的 **GOAL-004-r3-seam-and-shared-conventions**（R3：Redis 接缝声明 + 轨道 owner 文档 + mail 迁移评估 + 组合根 fx 改造）执行**独立交叉审计**。

## 硬约束

1. **只输出审计报告文本（Markdown），不得修改或创建任何文件**（P-003；落盘由编排器完成，`source: independent` 保留）。
2. 必须独立复核（可实际运行命令验证），至少执行：
   - `cd apps/api` 后：`go vet ./internal/composition/... ./internal/cache/... ./internal/config/...`；`go test ./internal/composition/... ./internal/cache/... ./internal/config/... -count=1`
   - 仓库根：`git status --short` 与 `git diff --stat`（越界核账）；`Select-String -Path go.mod,go.sum -Pattern redis`（判据 #4 验证面）
   - 检查 `internal/mail/` 是否有任何改动（I-026-004 承诺 mail 零改动）
3. 报告必须含：verdict（pass / conditional / fail）、scope、信息门禁核验、交付物逐条核验、findings 表（required / recommended / informational）、开放 required 计数。

## 核验要点（逐条对照）

- **C1 裁决**：D-001 是否记录用户裁决（I-026-004 不迁移 + F-002 fx 挂载；2026-09-01）；未选方案留痕。
- **架构短文 `docs/architecture/cache-redis-seam-and-track.md`**：
  - 判据 #4（§2 接缝声明）：端口不变（同一 kernel.Cache 接口、供应商类型只在 internal/、[]byte 无需序列化）；key = `<ns>:<key>`；TTL 映射（PX / 滑动 EXPIRE 续期 / 永不过期 / 惰性由服务端承担）；连接管理（组合根单一持有 + 启动 fail-closed PING）；不引入客户端依赖。
  - 判据 #5（§3 轨道约定）：单一所有者 VP-026；VP-027 激活继承；VP-028 排除（outbox/MQ）；命名空间登记表 + owner 义务；变更流程（owner 决策 + 修订史）。
  - 边界：短文是否把任何 Redis **实现**当成本波交付（不得）；RT-Q03 trigger 是否保持 gated。
- **I-026-004**：评估附件（attachments/mail-cached-adapter-evaluation-2026-09-01.md）论证是否成立（版本戳 vs TTL 语义；三候选否决）；`internal/mail/runtime.go` 是否零改动。
- **F-002 兑现（组合根）**：`NewApp` fx.Provide 是否含 `newCache`；`newMux` / `newMuxWithExtraProviders` 注入参数 `cachePort kernel.Cache`；fx 容器 = 长生命周期持有者（kernel.Cache 单例）；eager 构造链（newMux→newServer→lifecycle）是否成立；seam 内是否还有旧「局部构建 + 谎言注释」残留；4 个测试调用点是否全部补参。
- **越界**：`internal/mail/`、`kernel/cache.go`、Profile / Manifest / Charter / go.mod 是否零改动；工作树是否只有 R3 owned paths（apps/api/internal/composition/*、docs/workspaces/workspace-026-cache-port/**、docs/architecture/cache-redis-seam-and-track.md、docs/README.md、docs/vision/plans/VP-026-cache-port.md）。

## 审计上下文（先读）

- `docs/workspaces/workspace-026-cache-port/workspace.md`
- `docs/workspaces/workspace-026-cache-port/GOAL-001-cache-port/00-meta.md`（Root 纲领 R3）
- `docs/workspaces/workspace-026-cache-port/GOAL-004-r3-seam-and-shared-conventions/00-meta.md`、`01-decision/D-001-r3-adjudication.md`、`02-execution/E-002-seam-and-track-landed.md`、`03-audit/A-001-r3-closeout-self.md`
- 被审交付：`docs/architecture/cache-redis-seam-and-track.md`、`attachments/mail-cached-adapter-evaluation-2026-09-01.md`
- 代码：`apps/api/internal/composition/composition.go`（NewApp / newMux / newMuxWithExtraProviders / newCache）+ 4 个被改测试文件；对照 `apps/api/internal/mail/runtime.go`（L116～L229）
- 分母链：workspace-026 `GOAL-002/01-decision/D-002-cache-port-contract.md`（v0.1.1 端口合同）、`GOAL-003/01-decision/D-001-r2-plan-freeze.md`、`03-audit/A-002-r2-impl-closeout-independent.md`（F-002 原文）
- 规则：`docs/architecture/independent-audit-execution.md`、`docs/architecture/principles.md`（P-003/P-005）

## 输出

直接输出最终审计报告文本（Markdown），不要输出工作过程叙述。