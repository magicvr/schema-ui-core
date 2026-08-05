---
id: E-013-r6-c63-lifecycle-matrix
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-013 · R6 C6.3 双 Profile 生命周期矩阵

## 已发生事实

- 提交 `9896a02` 完成生命周期切片。`kernel.Runtime` 在 Start 失败时逆序清理此前
  成功 Start 的模块，在 Ready 失败时逆序清理全部已 Start 模块；两条失败路径均清空
  started 集合，后续重复 Stop 为 no-op。
- 生命周期失败保留 hook 返回的结构化 stable code/module；若 hook 返回非结构化错误，
  使用阶段对应的 lifecycle code 与失败模块。cleanup error 追加到 detail，不覆盖原错误。
- `Stop` 即使遇到错误仍继续逆序清理全部模块，返回反向顺序中的首个错误并清空
  started 集合。
- composition 在 Ready 失败时不再重复调用 Runtime Stop；kernel 负责模块 hook 清理，
  composition 只关闭 listener/store。readiness gate 仍只在全部 Start+Ready 成功后设置。
- 新增以真实 `mvp` / `admin` resolved Plan 参数化的 Runtime 矩阵，覆盖成功
  Start+Ready+Stop、Start 失败清理、Ready 失败清理、Stop error continuation 与重复 Stop。
- composition 测试新增两种 Profile、`127.0.0.1:0` 的 Fx app 成功 Start/Stop，并将端口
  占用的 stable lifecycle failure 测试参数化为两种 Profile。

## 验证

- `go test ./internal/kernel ./internal/composition`（`apps/api`）→ exit 0。
- `go test ./...`（`apps/api`）→ exit 0。
- `go vet ./...`（`apps/api`）→ exit 0。
- owned staged diff check 通过；提交仅包含四个 lifecycle/composition 路径，三份既存
  handler 测试换行噪音未暂存。

## 事实边界

- D-003 的四个代码切片均已实现并形成 C6.3 self/cross 审计候选证据；本记录不自行
  关闭 Root A-010 F-003b。
- C6.3 self + Grok independent 尚未完成；R6-I003 保持 `collecting`，GOAL-013 保持
  `active / 2/4`。

## 下一步（计划）

- 写入 C6.3 self audit，再由 Grok Build 对同 scope 执行独立 `/audit`；由 `/govern`
  响应全部意见后，才判断 R6-I003、C6.3 与 Root F-003b 是否可按 `fixed` 闭合。
