---
id: E-004-r4-c2-child-closeout
doc: execution-entry
goal: GOAL-008-r4-c2-module-contract-extension
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-004 · R4-C2 模块契约扩展子目标关门

## 已发生事实

- kernel 契约层实现：`contribution.go`（六类 Contribution + 身份校验）、
  `provider.go`（Provider/Registrar + `RegisterContributions` 双检/fail-closed +
  `descriptorsMatch` 全规范匹配 + `contributionConflictError` 稳定 code + navigation
  两遍扫描）、`persistence.go`（`CollectPersistence` compiled-global catalog）、
  `module.go`（`ContributionKeys.Fragments` + Plan 声明冲突）。
- Grok A-003 `conditional` 三条 required（F-IND-C2-001/002/003）由 A-004 以 `fixed`
  闭合；recommended C2-007/008 `fixed`，C2-004/005/006 显式延至 C3/C5 并登记门禁。
- 验证：kernel 测试 + 全量 API 测试 + vet 通过；fx 静态检查通过；双 Profile 契约
  矩阵通过。
- C2.1-C2.4 勾选；meta `progress: 4/4`；goal-tree 同步为 `done 4/4`。

## 向 GOAL-005 传递的 C2 context（供 C3 使用）

- Provider/Registrar/Contribution 精确契约在 `apps/api/internal/kernel/` 可用；
  `RegisterContributions` + `CollectPersistence` 为 C3 业务模块 provider 化的入口。
- C3 门禁（冻结 §3/§4.2/§2.2）：Manifest secrecy 扫描、Ready 失败反向清理、PolicyID/
  Visibility/JSON 校验器、ledger 侧 drift/unknown 运行时 fail-closed、真实双 Profile
  register/conflict/Start/Ready 矩阵。
- 业务模块（Users/Roles/Settings/Activity）仍为中心注册，C3 迁移未开始。

## 提交

本目标 close checkpoint 已 git 提交，提交标题 `docs(workspace-003): close GOAL-008
R4-C2 module contract extension`（exact hash 见 git log）。
