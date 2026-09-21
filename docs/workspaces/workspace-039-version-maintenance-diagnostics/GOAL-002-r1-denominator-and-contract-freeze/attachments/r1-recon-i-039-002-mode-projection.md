---
doc_type: recon
id: r1-recon-i-039-002-mode-projection
parent: GOAL-002-r1-denominator-and-contract-freeze
status: draft
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 侦察 · I-039-002 四模式投影

> 只读事实，不是冻结决策。

## 现行投影（代码）

| `runtime.mode` | 写门禁 HTTP | 错误码 | Host bootstrap `availability.mode` | Host 结果 | Shell 能否加载 | status `availabilityMode` |
|----------------|-------------|--------|--------------------------------------|-----------|----------------|---------------------------|
| `normal` | 放行 | — | `normal` | `READY` | 是 | `normal` |
| `maintenance` | 503（登录/恢复/邀请白名单除外） | `SERVICE_MAINTENANCE` | `maintenance` | **`MAINTENANCE` 终态** | **否**（HostFailureScreen） | `maintenance`（监控模块若仍编译） |
| `degraded` | 503 | `SERVICE_DEGRADED` | `degraded` | `READY_DEGRADED` | 是 | `degraded` |
| `read-only` | 503 | `SERVICE_READ_ONLY` | **`degraded`**（故意折叠） | `READY_DEGRADED` | 是 | **`read-only`**（status 保留精确值） |

锚点：`handler/operational.go`、`handler/bootstrap.go`、`host/bootstrap.ts:273-318`、`host/boot.ts:195-200`、`systemmonitoring/provider.go`。

前端 `feedback-policy.ts` 已把 `SERVICE_MAINTENANCE` 分类为 maintenance 反馈；无持久横幅。`apps/web` 无 `runtime.mode` 命中。

## 关键推论

`I-039-004`「Shell 横幅」**不能**覆盖 `maintenance`：availability-gate 在拉 manifest / 进 Shell **之前**终态。若坚持 maintenance 也出 Shell 横幅，必须改 Host 投影（把 maintenance 改成 `READY_DEGRADED`），这会改变已交付 VP-012 行为。

## 待 P-004（maintenance 呈现）

见编排器提问。
