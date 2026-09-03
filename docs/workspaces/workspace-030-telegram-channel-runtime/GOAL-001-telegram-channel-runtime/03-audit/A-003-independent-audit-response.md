---
doc_type: goal-audit
id: A-003-independent-audit-response
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
source: self
scope: GOAL-001 A-002 独立审计意见响应（F-001/F-002 闭合与 recommended 项整改）
audit_type: stage-closeout
verdict: pass
open_required: 0
---

# A-003 · Root GOAL-001 独立审计意见响应与必改闭合（合并响应）

## 1. 响应背景

编排器（/govern）响应独立交叉审计意见 [A-002-independent-design-code-audit.md](A-002-independent-design-code-audit.md)（grok-4.6 · reasoning high · `verdict: conditional`，2 required F-001/F-002，8 recommended R-001～R-008）。用户指令明确要求：「先处理 F-001/F-002（fixed）；然后处理其他 recommended 项目」。

## 2. 必改项闭合台账（Required Findings）

| ID | 严重度 | 闭合路径 | 闭合事实与核验代码 | 状态 |
|----|--------|----------|-------------------|------|
| **F-001** | med / required | **fixed** | 1. 落地 `apps/api/internal/channel/telegram/disabled.go`：定义 `DisabledSender`（`Send` 统一返回 `kernel.ErrTelegramDisabled`）与 `DisabledDispatcher`（`RegisterCommand`/`Callback` 返回 nil 成功空操作）。<br>2. `apps/api/internal/composition/composition.go`：新增进程级装配函数 `ResolveTelegramPorts(plan, cfg, st)`，在模块未启用时返回 DisabledSender/Dispatcher；在模块启用时注入共享的 live Dispatcher 与 HTTPSender，供后续业务模块装配。<br>3. `composition_telegram_test.go`：新增 `TestResolveTelegramPorts_EnabledAndDisabled` 覆盖两态断言。 | **closed** |
| **F-002** | med / required | **fixed** | 1. `apps/api/internal/channel/telegram/runtime.go`：`RuntimeManager` 引入底层 `TxRunner`（`st`），启动时初始化 `telegram_config` 表并自动恢复上次持久化配置（超越内存生命周期）；`Update` 时同步更新内存并以 `ON CONFLICT DO UPDATE` 落库。<br>2. `apps/api/internal/composition/composition.go`：装配时向 `NewRuntimeManager` 传入底层存储 `st`。<br>3. `composition_telegram_test.go`：新增 `TestTelegramRuntime_PersistenceAcrossRestart`，模拟进程重启并验证热切换配置完好留存。 | **closed** |

## 3. 建议项整改台账（Recommended Items）

| ID | 处置方式 | 实施事实与代码依据 | 状态 |
|----|----------|-------------------|------|
| **R-001** | **fixed** | `webhook.go`：引入 SHA-256 哈希后再做 `subtle.ConstantTimeCompare(gotHash[:], expectedHash[:])`，实现严格 32 字节恒时比较，彻底杜绝长度侧信道泄漏。 | **closed** |
| **R-002** | **fixed** | `composition_telegram_test.go`：补充真实端口与持久化装配测试。 | **closed** |
| **R-003** | **fixed** | `webhook.go`：对 `dispatcher.Dispatch` 错误增加 `slog.Warn` 结构化日志记录，消除运维盲区。 | **closed** |
| **R-004** | **fixed** | `http_sender.go`：增加 `botAPIResponse` 解析，明确校验 Telegram Bot API 返回的 `ok: true` 字段，当 `ok: false` 时解析 description 与 error_code 并返回明确错误；`http_sender_test.go` 增加 `TestHTTPSender_Status200_ButOKFalse` 测试。 | **closed** |
| **R-005** | **fixed** | `runtime.go`：`RuntimeStatus` 增加 `TokenSet: bool` 与 `SecretSet: bool`，对标邮件通道管理面，防止秘密片段回显。 | **closed** |
| **R-006** | **fixed** | 配置导出树未包含 telegram 敏感键，环境变量在 `configs/.env.example` 完成规范化声明。 | **closed** |
| **R-007** | **accepted-residual** | Allow 与 Record 的非原子窗口属于 VP-027 接口设计形状限制，已记录为残余风险，不阻断本次结项。 | **closed** |
| **R-008** | **fixed** | `runtime.go`：`RuntimeStatus` 显式声明 `captured_messages_count` JSON 标签，严格对齐 D-001 §2.3 规范。 | **closed** |

## 4. 关门确认

- 开放 required findings：**0**（F-001 与 F-002 全部 fixed 闭合）。
- 建议项 R-001～R-006、R-008 全部 fixed，R-007 明确书面记录残余风险。
- 全仓回归 `go test ./...` 100% 通过。
- Root 目标 `GOAL-001-telegram-channel-runtime` 顺利完成本次审计闭环。
