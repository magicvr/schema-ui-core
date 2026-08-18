---
id: E-003-r5-s1-implementation
goal: GOAL-006-r5-maintenance-read-only-gate
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# E-003 · R5 S1 运行态配置与状态投影实现

## 已核对事实

- `c4856f2` 增加 `runtime.mode` 的 normal/maintenance/degraded/read-only 配置模型；YAML 缺省为 normal，`RUNTIME_MODE` 覆盖 YAML，显式空值和未知值写入 `LoadError`，`ValidateProd` 对非零非法 mode fail closed。
- bootstrap 保持上游四模式不变：read-only 映射为既有 degraded；没有把 `form.controls.readonly` 放入 `disabledCapabilities`，也没有新增 protocol mode。
- system-monitoring status 单行 envelope 追加 `availabilityMode`，原 `status`/`ready` 仍只表达存储与模块图 readiness；页面增加 Availability stat card。
- `go test ./internal/config ./internal/composition ./internal/modules/systemmonitoring` 通过；`go test ./internal/handler -count=1` 通过（208.310s）；`go test ./internal/docscheck` 通过。

## 边界

统一写请求门禁与核心/Provider 黑盒矩阵属于 S2；本条不将其写成 S1 已完成事实。
