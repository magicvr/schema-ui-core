---
id: GOAL-012-w12-multi-instance-rate-limiting
doc: execution-entry
record_id: E-001
status: recorded
goal: GOAL-012-w12-multi-instance-rate-limiting
created: 2026-08-26
updated: 2026-08-26
version: 1.0.0
---

# E-001 · 立项与登记项代码现状核验（S1）

## 已发生事实（2026-08-26）

1. 开波写入完成：五件套 + 三 ledger 目录 + attachments 建齐；`parent = GOAL-001-production-hardening`；编号 011→012 无冲突。
2. **来源核验**（本会话直接读取代码，非转抄登记文本）：
   - `apps/api/internal/handler/rate_limit.go` L12–39：`loginRateLimiter` = 进程内互斥锁保护的内存滑动窗口桶（window/max/capacity 参数化；capacity 缺省 `1<<16`）；类型注释明示 "It is process-local and best-effort — it does not protect against distributed attacks"。
   - `apps/api/internal/handler/recovery.go` L58：恢复面接线 `newLoginRateLimiter(15*time.Minute, 20, 1<<16)`；登录面同型复用（W4 P0-1 / D-001 P1 既有实现）。
   - 键空间 = `loginClientIP(r) + "|" + lower(account)`（recovery.go L74–76）；X-Real-IP 仅信任显式配置的反代 CIDR。
   - 上游登记：[workspace-019 GOAL-001 E-009 §F-002](../../../../workspace-019-iam-recovery/GOAL-001-iam-recovery/02-execution/E-009-a001-finding-fixes.md)（Q2 引用）于 2026-08-26 落盘，含语义要点四条与本波归属指认。
3. 阻断检查：W1–W11 全部 `done`（GOAL-011 关门后独立复核 A-005 pass / A-006 正式确认，开放 required = 0）；`docs/vision/reviews.md` VRev open required = 0 → 开波无阻断。

## 待办（非事实，仅指针）

- I-001/I-002 required 裁决 → S2 方案冻结（D-002）。在此之前不实施任何代码变更。
