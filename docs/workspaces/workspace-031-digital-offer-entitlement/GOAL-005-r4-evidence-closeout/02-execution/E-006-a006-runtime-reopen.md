---
doc_type: goal-execution
id: E-006-a006-runtime-reopen
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: done
version: 1.0.0
---

# E-006 · A-006 关门撤回与运行时集成整改

## 事实（时间线）

- 2026-09-05 · **A-006 independent runtime-integration 审计**（用户侧追加，codex local）：verdict **fail**、4 required——F-003（high）`biz.digital-offer` 未进入 `kernel.BuiltinModules()`，`ResolvePlan` 注册表无法解析该模块，生产配置路径不可达（composition 分支为不可达代码）；F-004（high）Offer Delete 缺失、「CRUD」表述与实现不一致；F-005（high）0070 全局迁移策略未裁决；F-006（med）缺组合根到 Manifest/schema/HTTP 的真实验收。
- 2026-09-05 · **关门撤回**：按 P-003（未合法闭合 required 存在时不得保持关门），GOAL-005 → active 1/2、Root → active 3/4、VP-031 → active（v0.3.1，revision history 记录撤回）；goal-tree/workspace 投影同步。
- 2026-09-05 · **D-003 裁决落盘**：F-004 按 D-002 §2 既有冻结收窄（无删除，生命周期管理）；F-005 裁决保留 compiled-global persistence（平台级既定语义，dormant schema 副作用显式接受并记录复审触发）。
- 2026-09-05 · **代码修复**：`kernel/profile.go` BuiltinModules 增补 `biz.digital-offer` descriptor（与 Provider.Descriptor 逐字段一致）；新增 `internal/composition/composition_digitaloffer_test.go`（plan 解析 enabled/disabled + 真实 NewApp/mux 组合根验收：Manifest 两页面、schema 认证/匿名/未知 pageId、Admin 401/200、公开目录 200、Telegram-disabled 装配语义）。发现并修复测试构造缺陷：手写 Config 缺 AuthAccessTTL/AuthRefreshTTL 默认值导致登录令牌 exp==nbf 即时过期——补 15min/30d 默认。
- 2026-09-05 · **A-007 响应**：四项 required 全部 closed（fixed ×3 + contract-conformant ×1）。`go test -count=1 ./modules/digitaloffer/... ./internal/composition/ ./internal/channel/telegram/` 全绿；全仓 `go test ./...` 回归绿（BuiltinModules 增补后）。

## 产物路径

- `01-decision/D-003-a006-runtime-response.md`
- `03-audit/A-007-self-response-a006.md`
- `apps/api/kernel/profile.go`、`apps/api/internal/composition/composition_digitaloffer_test.go`

## 进度评估

- C2 重开进行中：待 A-008 independent closure 复审确认 open required = 0 后重新关门（GOAL-005 2/2、Root 4/4、VP-031 closed）。
