---
doc_type: goal-audit
id: A-001-self-implementation-review
parent: GOAL-003-r2-offer-purchase-wallet
date: 2026-09-05
status: closed
version: 1.0.0
---

# A-001 · R2 实施 self 审计（execution-facts）

## A-001 · R2 实施自审（2026-09-05）

- **source**：self
- **auditor**：编排器（/govern 会话内自审）
- **类型 / scope**：execution-facts · GOAL-003 C2/C3 交付（`apps/api/modules/digitaloffer/` + handler/composition/errorcatalog 接线）对照 D-002 v1.1.0 合同与 VP-031 判据 1/2
- **verdict**：**pass**（2 项 recommended；无 required）

### 范围与区间

审本目标 C2/C3 已实施事实；不审 R3（核验/消耗 API、Telegram 命令——E-001 已标注 R3 归属）。工作区绑定与 canonical 一致。

### 成果（有证据）

- 判据 2 路径：单事务购买（freeze→deduct_frozen→凭证→权益）经 `service/purchase_test.go` 在 SQLite 与真 PostgreSQL 双库验证；并发双发恰一凭证/一次扣款/一份权益；余额恰够一次的并发竞争恰一成功且无冻结残留；幂等重放与跨 offer 冲突（`ErrRequestIdConflict`）可测。
- 判据 1 路径：Admin 路由 + 权限键 + schema 页 + 审计 fail-closed（注入审计失败 → 域写回滚，测试证明）；公开目录仅 on_sale 且不泄漏内部字段。
- 合同符合性：§4.4 协议逐条落地（attempt 起点回读、互异幂等键、ref 反链、自表唯一竞争经 `kernel.IsUniqueViolation`、不依赖 wallet 未导出 sentinel）；§7 事件名与 record id 冻结；§8 公开目录桶 AllowRecord、无 Clear；§10 红线（默认 Profile 排除、subject-only、占位符 `?`）。
- 全仓 `go build ./...` + `go test ./...` 回归绿（commit `0e65f392`）。

### 对照成功标准（GOAL-003）

| 标准 | 状态 | 证据 |
|------|------|------|
| 1 · Offer CRUD + Admin 面 + 审计 + C 端列表 | 达成 | provider/handler + `digitaloffer_test.go` |
| 2 · 购买路径 + 同事务 fail-closed + 并发测试 | 达成 | `purchase_test.go`（双数据库） |
| 3 · 不越界（红线/Profile/subject-only） | 达成 | composition `plan.HasModule` 门控；store 无 admin.users 关联 |
| 4 · 关门审计闭合 | 未开始 | 本条即 self 审；independent 待跑 |

### Findings

- **F-001 · 购买/价目限流桶的接线点在 R3**
  - 严重度：low；建议：recommended；状态：open → 转 GOAL-004 承接
  - 描述：§8 冻结的 `bizoffer|purchase|<subject_id>` / `bizoffer|price|<subject_id>` 桶，其生产入口（Telegram buy / price 命令）属 R3；R2 的购买入口仅为进程内 Service API（测试面），无 HTTP/通道调用方，桶无处可挂。Provider 已注入 `kernel.RateLimiterProvider`（构造参数），R3 接线即可。
  - 处置：登记为 R3（GOAL-004）实施项；不阻断 R2 关门（判据 2/1 不依赖限流，V-F119 的冻结语义已落 §8 合同与公开目录桶）。

- **F-002 · Admin purchases 列表未做分页过滤字段校验以外的范围控制**
  - 严重度：low；建议：recommended；状态：open → **closed（fixed）**
  - 描述：purchases/entitlements 列表为 `digitaloffer.read` 门控的跨主体审计视图（合同 §7 定位），但 Q/subjectId/offerId 过滤值未做 LIKE 转义以外的长度上限。复核确认 store 层走参数化查询 + `escapeLike`，无注入面；长度上限由 handler `INVALID_PAGE*` 与分页上限覆盖，Q 值超长仅影响查询效率。补一条长度防御即可。
  - 修复：`store` 层新增 `searchQ`（trim + 100 字符截断），三个列表过滤统一走该入口（commit 附于 A-002 审计轮）。

### 必改项汇总

- 无 required。

### 结论 + 建议下一步

- 实施与合同一致，测试分母覆盖判据 1/2 的可验证面。建议：F-002 补丁 → codex independent 审计（A-002，gpt-5.6-sol · medium）→ 意见响应 → GOAL-003 关门。
