---
doc_type: goal-audit
id: A-011-r2-c3-finding-remediation-self
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
source: self
auditor: Codex govern
audit_type: finding-remediation
scope: 响应 A-010 F-001～F-003 required findings 的 C3 修复复核
verdict: pass
open_required: 0
version: 0.1.0
---

# A-011 · R2 C3 必改项修复 self 复核（2026-09-04）

## 复核结论

A-010 `source: independent` 的 verdict=`fail`、open required=`3` 原文保留。对其 F-001～F-003 的代码修复已在 `4cc96b06` 落盘，并由本条 self 复核为可核对的 `fixed` 候选；因这些修复影响独立审计结论，仍须 Grok re-audit 后由 `/govern` 追加最终响应。本条不改写 A-010、不自行关闭 C3。

| finding | 修复证据 | self 结论 |
|----------|----------|-----------|
| A-010 F-001 · 30s/40s polling deadline | `bot_api.go` 保留 Telegram timeout=30s，新增 `PollingRequestContextTimeout=35s`，polling client timeout=40s；`bot_api_test.go` 核验严格顺序与 request context deadline grace | fixed 候选 |
| A-010 F-002 · transport error 泄漏 token | `bot_api.go` 对 create/execute request failure 使用不含 URL 的安全诊断；`bot_api_test.go` transport error 断言不含 bot token/secret；manager error status 使用该安全错误 | fixed 候选 |
| A-010 F-003 · Stop 后 watcher 可能重启 | `connection_manager.go` 的 watcher 与 `startPolling` 双重 started/lifecycle 门禁；`TestConnectionManager_WatchDemandDoesNotStartWhenStopped` 回归 | fixed 候选 |

## 保留的推荐项

- A-010 F-004（缺 URL 对等测试、Stop timeout、多 lease、真实长等待矩阵）仍为 recommended/open，转入 C5 证据矩阵。
- A-010 F-005 以及 A-006 F-002～F-004 仍为 recommended/open，转入 C4/C5；未因本次 required 修复而静默关闭。
- C4 Admin settings UI/heartbeat HTTP surface 与 R3 会话范围仍未实施。

## 放行边界

当前 self verdict=`pass`、open required=`0`，但 A-010 的独立 required finding 尚待复审；C3 checkpoint 与 GOAL-003 仍保持 `active · 2/5`。
