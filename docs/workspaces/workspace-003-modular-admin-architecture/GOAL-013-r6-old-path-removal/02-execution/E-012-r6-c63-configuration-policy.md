---
id: E-012-r6-c63-configuration-policy
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-012 · R6 C6.3 Configuration 与 Policy 校验

## 已发生事实

- 提交 `2548e42` 新增框架无关 `ConfigurationContribution`、
  `Registrar.Configuration` 与 `ContributionSet.Configurations`。
- kernel 对 `Key == Namespace`、ASCII 小写点分 namespace、JSON object defaults、
  validator 必填/defaults 必须通过、bytes copy、全局冲突与确定性排序执行 fail-closed
  双检。
- `admin.settings` descriptor/runtime 同时声明并注册 `settings.branding`；模块-owned
  configuration 固定现有 `siteTitle`/`logoUrl` defaults 与严格 validator。Settings PATCH
  的 config-change header 从 contribution namespace 传入，生产 handler 不再持有该模块
  namespace 字面量。
- PolicyID 与 Visibility v1 现在都校验单 policy reference grammar；未版本化布尔表达式、
  大写、空段、非法连字符、下划线与非 ASCII 均在 Registrar 阶段拒绝。
- well-formed 但未知的 policy（测试为 `system.custom`）通过 kernel 语法层后，由
  `core.auth-session` `rolesForPolicy` allowlist 在 system-data reconcile 前分别拒绝
  permission 与 navigation，保持薄内核分层。

## 验证

- 定向 kernel/settings/auth-session/handler/composition 测试 → exit 0。
- `go test ./...`（`apps/api`）→ exit 0。
- `go vet ./...`（`apps/api`）→ exit 0。
- 生产 handler `settings.branding` 字面量零命中；旧 trim-only policy 校验文案零命中。
- Git staged diff check 通过；三份既存 handler 测试换行噪音未暂存。

## 事实边界

- Configuration/Policy 子切片已实现。Lifecycle matrix 与 C6.3 self + Grok independent
  尚未完成；R6-I003 继续为 `collecting`，`progress` 继续为 `2/4`。
