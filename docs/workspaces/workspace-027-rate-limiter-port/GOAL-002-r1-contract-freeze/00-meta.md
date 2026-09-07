---
id: GOAL-002-r1-contract-freeze
title: R1 合同冻结（RateLimiter 端口契约 / key 语义 / 窗口语义）
status: done
parent: GOAL-001-rate-limiter-port
created: 2026-09-01
updated: 2026-09-01
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-027-rate-limiter-port
primary_plan: VP-027-rate-limiter-port
serves_summary: 承载 VP-027 R1 阶段（判据 #1）：冻结 RateLimiter 端口契约——API 形态（I-027-001）、窗口语义（I-027-003）、key 维度（I-027-004）；端口本体落 kernel/ratelimit.go + 合同级快测。
---

# GOAL-002 · R1 合同冻结

## 概述

执行 Root 纲领 **R1**：在仓库既有内核端口先例（`kernel.Cache` / `Store` / `ObjectStore` / `MailSender`）与既有 `loginRateLimiter`（`internal/handler/rate_limit.go`：滑动窗口 · allow 不注册 key · 容量驱逐 · trusted-proxy）之上，冻结 VP-027 **RateLimiter 端口契约** — API 形态（拆分 vs 内聚由用户裁决）、窗口语义（滑动保持）、key 维度（不扩展）。**合同正文 = GOAL-002 D-002 产物**；端口本体（`apps/api/kernel/ratelimit.go`）+ 合同级快测（`kernel/ratelimit_test.go`）在本目标落地并验证；内存供应商与 7 处使用点迁移（判据 #2/#3）归 R2（GOAL-003）。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **信息裁决**：I-027-001（API 形态）、I-027-003（窗口语义）、I-027-004（key 维度）用户裁决（P-004；I-027-001 required，I-027-003/004 non-blocking lead 建议确认） | **已关门**（2026-09-01 用户裁决全部采纳建议：语义拆分保持 / 滑动窗口保持、策略接口独立 / 本波不新增复合 key——D-001） |
| C2 | **合同正文 + 端口落地**：D-002 合同冻结；`kernel/ratelimit.go`（RateLimiter / RateLimiterProvider + InWindow/RetryAfter 纯函数）实现；合同级快测绿 | **已关门**（D-002 v0.1.1 冻结（含 A-002 F-006 勘误 §3：剪枝仅 Allow）；端口 + helper 落地；`go vet` 0 / `go test ./kernel/...` 绿（15 子例）/ `go build ./...` 通过 / `gofmt -l` 空；2026-09-01） |
| C3 | **审视与关门**：合同自审（self A-001）+ independent（本地 grok build grok-4.6 · high A-002）合并响应；R1 关门、Root 信息台账回写 | **已关门**（A-001 self `pass` 0 required + A-002 grok build independent `pass` 0 required；A-003 合并响应 F-001～F-007 全处置（fixed ×6 · fixed-recording ×1）；R1 关门 3/3 · Root progress 1/4 · VP-027/workspace 台账回写；2026-09-01） |

`progress` = 已关门检查点数 / 3。当前 **0/3**。

## 成功标准（方向级）

1. 端口契约冻结：RateLimiter（Allow/Record/RetryAfterSeconds/Clear + key 寻址 + 供应商无关）按语义拆分保持冻结、快测可断言（判据 #1）。
2. API 形态按用户裁决落盘：Allow 不注册 key（D-001 P1 防喷洒）、Record 失败才计数、Retry-After 语义保持（W12 D-002）；`now` 注入保证确定性。
3. 窗口语义按用户裁决落盘：滑动窗口保持 + 访问路径内剪枝（惰性清理，无后台协程 → 无新生命周期 → 停机义务不触发）；不与 Cache ExpiryPolicy 共用策略形态。
4. key 按用户裁决落盘：不透明字符串、不新增复合维度；`IP|identifier` / `op|IP|user` / 纯 IP 形态保持（R2 迁移分母）。
5. 容量与驱逐义务冻结：capacity ≤ 0 → 默认 `1<<16`；distinct key 上限驱逐最老 key（D-001 P1）。
6. 未越界：不改 Profile 默认集 / 模块矩阵 / Manifest 装配；不引入 Redis 客户端依赖；不改 Charter；7 处使用点不动（R2 迁移）。

## 信息就绪与未知项

与 Root / VP-027 同号镜像（I-027-00x）；I-027-002（`loginRateLimiter` 迁移策略）最晚阶段 R2，不在本目标关闭。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-027-001 | required | RateLimiter 端口 API 形态：语义拆分保持 vs 内聚 Allow vs 回调式；RetryAfter 语义 | 方案冻结 + 判据 #1 | C1 | 用户裁决 | **verified** | — | 2026-09-01 用户裁决：**语义拆分保持**（Allow 不注册 + Record 失败计数 + RetryAfterSeconds + Clear；now 注入）（D-001 accepted；合同 §1） |
| I-027-002 | required | 既有 `loginRateLimiter` 迁移策略：演进为内存供应商 vs 保留并存（双轨）；key 维度是否扩展 | 判据 #3 | R2 | 用户裁决（R2 前置） | 待裁决 | — | — |
| I-027-003 | non-blocking | 窗口语义默认：滑动窗口 vs 固定/混合；策略接口是否与缓存 VP-026 共用形态 | 判据 #2 | C1 | lead 建议 + 用户确认 | **verified** | — | 2026-09-01 用户确认：**滑动窗口保持 + 策略接口独立**（D-001 accepted；合同 §3） |
| I-027-004 | non-blocking | 限流 key 维度扩展：是否新增"路由+用户"复合 key | 判据 #1 | C1 | lead 建议 + 用户确认 | **verified** | — | 2026-09-01 用户确认：**本波不新增复合 key**（D-001 accepted；合同 §2） |

## 父目标

- `GOAL-001-rate-limiter-port`（Root · 纲领 R1）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式（Root D-001 已定）：阶段关门 default self；本目标落 kernel 公共面（兼容性门禁）→ **C3 走 cross**：A-001 self + A-002 本地 grok build（grok-4.6 · high）independent（项目级默认执行路径 `docs/architecture/independent-audit-execution.md`）。
- R1 合同为本 VP 首波冻结分母；D-002 冻结后实施（R2）与验收（R4）以本合同为准。既有 `internal/handler/rate_limit.go` 语义（allow 不注册 / 驱逐 / Retry-After / trusted-proxy）为迁移不回归基线（判据 #3）。