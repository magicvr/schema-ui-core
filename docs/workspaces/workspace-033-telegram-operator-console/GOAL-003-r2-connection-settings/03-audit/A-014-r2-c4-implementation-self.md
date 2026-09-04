---
doc_type: goal-audit
id: A-014-r2-c4-implementation-self
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
source: self
auditor: Codex govern
audit_type: execution-facts
scope: R2 C4 Admin settings UI、polling lease HTTP 接缝、认证/会话隔离、composition wiring、disabled profile gating 与相关 API/Web 测试
verdict: pass
open_required: 0
version: 0.1.0
---

# A-014 · R2 C4 Admin 设置页与 polling lease 自审（2026-09-04）

## 核对结论

`verdict: pass`，`open_required: 0`。C4 scope 的实现事实可由提交 `d95f7544`、源码与定向测试核对；未把本条 self 结论当作 independent 证据，C4 检查点仍待 Grok independent review。

## 核对项

| 项 | 结果 | 证据 |
|----|------|------|
| lease route contract | pass | `lease_handler.go` 三个 POST action；provider descriptor/registration 与 `kernel/profile.go` 同步声明；均 `Public: false` |
| authentication / session isolation | pass | `auth.IdentityFrom` 后只取非空 `SessionID`；缺身份、权限或 session fail closed；不以 user id fallback；`lease_handler_test.go` 覆盖多 session |
| manager接缝 | pass | composition 为真实 `TelegramRuntime.Connection` 构造 lease handler；mux 测试验证 acquire/release 影响同一 manager；`ConnectionManager` 继续保持唯一 receiver owner |
| Admin UI | pass | custom component 支持 mode、显式 URL、write-only secrets、非 secret connection 状态；5 个定向 Vitest 用例通过 |
| browser lease lifecycle | pass | polling 模式 acquire + 10 秒 heartbeat + cleanup release；promise queue 保证 heartbeat/release 顺序；模式切换/卸载清理不重复创建租约 |
| profile / i18n boundary | pass | enabled composition route 可达；disabled profile settings、lease、webhook、schema 均 404；en-US/zh-CN 新增 key 可解析 |
| required findings | pass | C4 scope 未发现 required finding；C5 矩阵缺口仍保留为后续范围，不在本条静默关闭 |

## 验证事实

- `go test ./internal/channel/telegram ./modules/channel/telegram ./internal/composition -count=1 -timeout=180s`：通过。
- `go test -race ./internal/channel/telegram -count=1 -timeout=180s`：通过。
- `npm test -- --run src/components/telegram-admin-tab.test.tsx`：5 tests 通过。
- `git diff --check`：通过；本条不改写 A-010/A-012 等历史意见。

## 边界与后续

- 用户 lease route/permission 选择交互未返回答案；实现采用 AI 推荐默认，已在 E-010 记录，未伪造为用户 accepted decision。该选择如被用户后续改写，应通过 `/govern` 留痕。
- A-010 F-004～F-005、A-012/A-013 转入的 C5 recommended 项（真实 30 秒等待、Stop timeout/时序、多 lease 完整过期矩阵、迁移/导出/并发 PATCH）仍开放；本条不关闭它们。
- 本条仅为 C4 self opinion；需要 independent opinion 后，编排器才可关闭 C4 并将 progress 更新为 `4/5`。
