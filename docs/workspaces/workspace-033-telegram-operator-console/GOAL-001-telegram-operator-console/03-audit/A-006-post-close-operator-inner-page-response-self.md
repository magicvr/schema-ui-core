---
doc_type: goal-audit
id: A-006-post-close-operator-inner-page-response-self
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: self
auditor: Codex govern
audit_type: finding-response
scope: 响应 A-005 F-001/F-002；Telegram 设置页与人工会话内页分离
verdict: pass
open_required: 0
version: 0.1.0
---

# A-006 · A-005 finding response（2026-09-05）

## 响应矩阵

| finding | 等级 | 响应 | 证据 |
|---------|------|------|------|
| A-005 F-001 | required | **fixed**：将 `captured_messages_count` 从 settings JSX 移入 `isOperatorSurface` 的 operator section；设置 surface 测试断言不显示该统计，operator surface 测试断言其仍可显示 | `apps/web/src/components/telegram-admin-tab.tsx`；`apps/web/src/components/telegram-admin-tab.test.tsx` |
| A-005 F-002 | recommended | **covered**：补齐 provider page contribution 的 ModuleID、Key、Owner、Resources、Actions、datasource 断言，并保留 manifest/schema 与 settings negative-surface 测试 | `apps/api/modules/channel/telegram/provider_test.go`；`apps/api/modules/channel/telegram/manifest/manifest_test.go`；`apps/api/modules/channel/telegram/schema/schema_test.go` |

## 验证

- Web 定向测试 3 files/38 tests 通过；API Telegram/manifest/schema、composition、kernel
  定向测试通过。
- `git diff --check` 与新增文件尾随空白检查通过；代码 checkpoint 为
  `6a94ba28fad08de43d3b88a129c5dcbcd0b18ccb`。
- A-005 的 required F-001 已走 `fixed`，无 residual 或 user-overruled；本条不改变 Root
  status/progress，等待 A-007 最终 independent re-audit。
