---
doc_type: freeze-matrix
id: r1-mode-projection-matrix
parent: GOAL-002-r1-denominator-and-contract-freeze
status: frozen
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# 冻结矩阵 · 四模式投影（I-039-002）

| runtime.mode | 写门禁 | 错误码 | Host 文档 availability.mode（生产者 · 本波将改） | Host 消费者结果 | Shell | `/me.runtimeMode` | 横幅 |
|--------------|--------|--------|--------------------------------------------------|-----------------|-------|-------------------|------|
| normal | 放行 | — | `normal`（不变） | READY | 加载 | `normal` | 无 |
| maintenance | 503；登录/恢复/邀请白名单 | `SERVICE_MAINTENANCE` | **`degraded`（本波改；现行为 `maintenance`）** | READY_DEGRADED | **加载** | `maintenance` | 有（维护文案） |
| degraded | 503 | `SERVICE_DEGRADED` | `degraded`（不变） | READY_DEGRADED | 加载 | `degraded` | 有（降级文案） |
| read-only | 503 | `SERVICE_READ_ONLY` | `degraded`（不变） | READY_DEGRADED | 加载 | `read-only` | 有（只读文案） |

Host *消费者* 仍实现 `availability.mode=maintenance` 终态（协议）。本仓生产文档不再发出该值。
