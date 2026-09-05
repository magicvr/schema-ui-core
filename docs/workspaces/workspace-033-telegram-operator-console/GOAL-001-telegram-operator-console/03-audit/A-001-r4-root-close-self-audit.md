---
doc_type: goal-audit
id: A-001-r4-root-close-self-audit
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: self
auditor: Codex govern
audit_type: root-closeout
scope: GOAL-001 Root R4 关门核验（VP-033 方向级退出判据 1～8；R1～R3 证据、红线、required finding 与当前 API/Web 验证）
verdict: fail
open_required: 1
version: 0.1.0
---

# A-001 · Root R4 关门 self audit（2026-09-05）

## 自审结论

基于当前 `HEAD e026c1b7` 重新核对 Root 成功标准、VP-033 退出判据、R1～R3
子目标状态、源码与测试；不把历史目标状态或上一阶段的通过直接当成 Root 关门依据。
R1、R2、R3 的阶段证据与 required finding 处理链可核对，当前 API/Web 回归和 Web
构建均通过；但发现一项仍未实现的 Root required 合同，因此本次 `verdict: fail`、
`open_required: 1`，Root 保持 `active · 3/4`。

## 证据矩阵

| 退出判据 | 当前核验 | 证据 |
|----------|----------|------|
| 1. 连接状态与显式 URL | pass | `apps/api/internal/channel/telegram/fake_bot_api_test.go:92-176` 核对 polling 的 `getMe → deleteWebhook`、webhook 的 `getMe → setWebhook`、显式 URL/secret；`apps/api/internal/composition/composition_telegram_test.go:250-472` 核对持久化与 authoritative empty settings。 |
| 2. 模式互斥热切换与 fail-closed | pass | `apps/api/internal/channel/telegram/connection_manager_test.go:19-65,467-534,561-635,772-817` 覆盖建立、失败切换 drain、缺 secret/URL 与 settings hot switch；`apps/api/internal/composition/composition_telegram_lifecycle_test.go:20-80` 覆盖 Fx shutdown drain。 |
| 3. 轮询启停、heartbeat、占用位 | pass | `apps/api/internal/channel/telegram/connection_manager_test.go:67-187,400-465,637-714`；`apps/api/internal/composition/composition_telegram_test.go:528-685`；Web `apps/web/src/components/telegram-admin-tab.test.tsx:341-427,478-549`。 |
| 4. 会话落盘、私聊/群分栏、人工发送与权限 | pass | R3 子目标 `GOAL-004-r3-session-operator-console/00-meta.md:31-35` 与 `03-audit.md:65-66`；`apps/api/internal/handler/telegram_operator_test.go:141-321,325-540`；Web `apps/web/src/components/telegram-admin-tab.test.tsx:120-287`。 |
| 5. 首波边界与默认 Profile | pass | VP-033 首波冻结/非目标 `docs/vision/plans/VP-033-telegram-operator-console.md:50-74`；真实 composition 的 enabled/disabled surface 与 default `mvp` 404 由 `apps/api/internal/composition/composition_telegram_test.go:789-885` 覆盖。 |
| 6. 证据矩阵与 required finding 归零 | fail | 当前本条发现 F-001，不能声称 Root required 已归零。R1～R3 子目标均为 done，GOAL-004 的 A-039/A-040 已确认 C4 `open_required: 0`，但 Root 自身尚有本条。 |
| 7. polling 单实例声明 | fail | VP-033 明确要求 UI 明示多副本会丢 Update（`docs/vision/plans/VP-033-telegram-operator-console.md:92-101`）；当前连接状态 UI 仅渲染状态/接收器/租约（`apps/web/src/components/telegram-admin-tab.tsx:674-684`），组件与双语 catalog 未提供该警示。 |
| 8. 审计闭合 | blocked | 本次 self audit 已落盘，但在 F-001 修复、targeted test 与 Root independent re-audit 前不得关门。 |

## Finding

### F-001 · required · polling 单实例风险未在 UI 明示

- **事实**：VP-033 方向级退出判据 7 要求 polling 模式在文档与 UI 明示「多副本会丢
  Update」，并且不得把 polling 标成 HA 生产路径。文档已有该边界；当前
  `TelegramAdminTab` 的连接状态区只显示 connection、receiver、bot 和 lease，没有
  polling 单实例警示；当前英文/中文 catalog 也没有对应文案。
- **影响**：管理员在 UI 选择或使用 polling 时无法看到已冻结的单实例数据丢失边界，
  因此 Root 的边界退出判据尚未满足。该 finding 是 required，不是可静默保留的
  recommended 项。
- **修复路径**：在 polling 模式的连接状态区域加入双语、可定位的 warning 文案，
  明确多副本会丢 Update 且 polling 不是 HA 路径；补组件测试验证 polling 显示、
  webhook 不显示。修复后重新运行 Web targeted/full test、build，并由独立审计复核。
- **状态**：open；未使用 `accepted-residual` 或 `user-overruled`。

## 验证事实

- `go test ./... -count=1`（cwd `apps/api`）：通过。
- `npm test -- --run`（cwd `apps/web`）：92 个 test files、1213 个 tests 全部通过。
- `npm run build`（cwd `apps/web`）：通过；仅有 chunk size warning，无 TypeScript/build
  error。`da9d955e` 修复的 `form-controls.tsx:946-947` 基线构建错误未再出现；构建生成
  的 conformance projection 已恢复到当前 Git checkpoint，未把生成物变化冒充源码修复。
- `git diff --check`：通过。

## 门禁结论

Root 不得在本条 finding 未合法闭合前标为 `done`。本意见只记录当前缺口，不修改
Root 状态；修复与复审完成后再由 `/govern` 汇总 self/independent 意见并关门。
