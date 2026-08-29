---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-006-r5-report
version: 0.1.0
---

# E-002 · breaking 实演（2026-08-29 · F-008 fixed · 用户裁决实发）

## 变革

- `kernel.JoinKeys` → **`kernel.JoinIdentifiers`**（等价行为 · 新名）——真实契约面 breaking；零内部使用（演练友好）；changelog 迁移说明成文（`changelog-breaking-v0.3.0.md`：定位/改写/回归四步）。
- 发布：`apps/api/v0.3.0`（0.x 阶段 minor 承载 breaking + Breaking 节迁移说明——语义登记于 changelog）。

## 下游迁移路径实证（golden-field）

| 步骤 | 结果 |
|------|------|
| 旧绑定：`kernel.JoinKeys("a","b")` 存在于组合根（模拟 v0.2.0 时代下游） | — |
| `go get …@latest` | ⚠️ 初次仍解析 v0.2.0——**proxy/sumdb 收录时延**（知识项再次应验）；等收录后 `go get @v0.3.0` 成功（v0.2.0 → v0.3.0） |
| `go build` | 💥 **`undefined: kernel.JoinKeys`**（breaking 断裂实证） |
| 按 changelog 第 2 步改写 `JoinIdentifiers` | ✅ build 0 + 冒烟全绿（kernel=2.0.0 · fresh=true） |

**结论**：breaking 实演完成——迁移说明可执行、断裂可预期、修复分钟级。F-008 闭合。