---
id: E-006-a002-f001-spotcheck-d4-d5
goal_id: GOAL-001-module-contribution-readiness
status: recorded
created: 2026-08-06
updated: 2026-08-06
version: 0.1.0
parent: null
---

# E-006 · 响应 A-002 F-001：补抽 DO NOT D4/D5

## 事实

1. 用户经 `/govern` 书面确认：采纳 A-002 independent `pass`，Root 维持 `done`；对 recommended **F-001** 选择 **补抽**（非 residual / overruled）。
2. 更新抽检附件 `attachments/s3-users-spotcheck.md` → v0.2.0：在 DO NOT 表补齐 **D4**、**D5** 行与补抽说明；结论改为 D1–D5 均已逐行勾选。
3. 核对依据（只读代码/架构，未改 runtime 主路径）：
   - D4：`apps/api/internal/modules/users/provider.go` Register 覆盖 HTTP/Schema/Authorization/Navigation/Manifest；`CompiledPersistence` 空返回 + 依赖 `core.auth-session`（与 playbook §1.2 一致）；未用按需能力替代核心六项。
   - D5：`apps/api/internal/composition/composition.go` 对 `admin.users` 为静态候选 + `plan.HasModule` 装配；无模块侧热插拔/`.so`/远程下载路径。
4. 未改 Root `status`/`progress`（保持 `done` / `4/4`）；未改 VP-004 status。

## 阻塞

无。

## 下一步（计划）

- 编排响应落盘见 `03-audit/A-003-response-a002.md`（计划/并行写入）。
- VP-004 正式 `closed` 仍须用户 `/vision` 确认（非本 E 范围）。
