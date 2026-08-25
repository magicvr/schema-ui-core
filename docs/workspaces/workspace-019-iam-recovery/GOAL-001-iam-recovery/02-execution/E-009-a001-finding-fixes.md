---
id: E-009
doc: execution-entry
goal: GOAL-001-iam-recovery
status: recorded
created: 2026-08-26
updated: 2026-08-26
version: 1.0.0
---

# E-009 · 响应 Root A-001：F-001 代码修正 + F-002 部署拓扑注意项登记

> Root `done 4/4`（2026-08-25）之后的审计响应维护工作；处置取舍见 D-003，闭合核验见 A-002。未改动任何目标 status / progress。

## F-001 · fixed（代码证据）

| 文件 | 变更 |
|------|------|
| `apps/api/internal/modules/authsession/password_policy.go` | 新增导出 sentinel `ErrPasswordPolicyNotSeeded`；`UpdatePasswordPolicy` 改用 `kernel.Result.RowsAffected()` 检查，单例行缺失时返回该 sentinel（fail closed，不再静默 no-op） |
| `apps/api/internal/handler/password_policy_settings.go` | 删除 L136 死导入保持行 `var _ = errors.Is`；PATCH 失败路径 sentinel 细分：`ErrPasswordPolicyNotSeeded` → `404 SETTINGS_NOT_FOUND`（既有冻结码），其余维持 `500 INTERNAL` |
| `apps/api/internal/modules/authsession/password_policy_test.go` | 新增：unseeded 行 → `errors.Is(err, ErrPasswordPolicyNotSeeded)` 断言；seeded 行 → 写入持久化断言 |
| `apps/api/internal/handler/password_policy_settings_test.go` | 新增：sentinel → 404 + `error=SETTINGS_NOT_FOUND` 信封断言；成功路径仍 200 且值落库 |

验证（2026-08-26 本会话实测）：`gofmt` 干净；定向测试 `go test ./internal/modules/authsession/ ./internal/handler/ -run "PasswordPolicy"` 全绿（authsession 1.8s ok / handler 2.2s ok）；全量回归见 A-002 记录。错误契约漂移护栏（`error_contract_test.go`）未新增任何字面量——`SETTINGS_NOT_FOUND` 本就在冻结集内且仍被 settings.go 持续发射。

## F-002 · 部署拓扑注意项登记（按审计处方）

**登记内容（持久）**：恢复面与登录面共用的 `loginRateLimiter` 为**进程内内存滑动窗口桶**（15 min / 20 次 / `IP|identifier`，容量 65,536 键；定义 `apps/api/internal/handler/rate_limit.go` L12–39，接线 `recovery.go` L58 与 login 面）。语义要点：

1. **单实例内有效**：阻止单节点在线爆破；进程重启即清零。
2. **多实例部署时限流预算按节点各自计算**：N 节点实际放行上限 ≈ N × 单节点预算；跨节点无共享状态。
3. **本区未引入新模式**：workspace-019 R2 复用既有 login 面同型限流器（W4 P0-1 / D-001 P1 已有实现），非本区新增强化项。
4. **边界归属**：分布式/共享存储限流属生产化部署拓扑决策，不在 workspace-019 边界内（A-001 F-002 审计原文同判）。后续生产化波次评估的**自然归属位为 VP-009 production-hardening 程序**（长期共享基架安全与健壮性容器）；是否立项由该程序波次规划决定，本次仅完成登记，不代为立项。

## 治理同步

- 索引更新：`01-decision.md`（+D-003）、`02-execution.md`（+E-009）、`03-audit.md`（+A-002 响应条目）、`goal-tree.md` 关后维护指针。
- A-001 独立意见正文保持原样未改（响应走独立 A-002 条目，不污染 `source: independent` 台账）。
