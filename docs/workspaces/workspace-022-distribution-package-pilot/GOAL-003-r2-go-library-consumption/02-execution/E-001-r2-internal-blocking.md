---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-go-library-consumption
version: 0.1.0
---

# E-001 · internal 阻断实验（2026-08-29）

## 实验设置

- 黄金下游仓骨架：`attachments/golden-consumer/`（`module golden.local/consumer` · `go 1.26` · `require github.com/magicvr/schema-ui-core/apps/api v0.0.0` · `replace => ../../../../../../apps/api`——嵌套 go.mod，主仓 `./...` 自动排除）。
- 桩程序：`main.go` import `github.com/magicvr/schema-ui-core/apps/api/internal/kernel` 并读取 `kernel.KernelAPIVersion`（冻结面 A 层锚点）。

## 结果（工具链 go 1.26.0 · windows/amd64）

```
go build ./...  →  退出码 1
main.go:8:2: use of internal package github.com/magicvr/schema-ui-core/apps/api/internal/kernel not allowed
```

**结论**：Go `internal` 命名空间规则（仅允许同模块树导入）**阻断外部模块 import `apps/api/internal/*`**。kernel 与 modules 位于 `internal/` 下意味着下游组合根在现有目录结构下**无法装配任何契约面**——包化必须以「A/B 层移出 `internal/`」为先决条件（方案见 D-001，待用户裁决）。

## 附带观察

- replace 本地路径消费本身工作正常（依赖解析成功，错误只出在 internal 规则）→ 「单 go.mod + replace」的 G1 粗粒度路径无其他阻碍。
- golden-consumer 桩保留为失败证据；S3 将重写为真实组合根。

## 影响

- GOAL-003 S1 完成第一个实验检查点；I-005 收集到关键证据。
- 外移属结构性重构（数百文件 import 改写）→ **关键决策，交用户裁决**（D-001）。