# Changelog · v0.3.0（breaking 演练 · VP-023 R5 F-008）

## [0.3.0] - 2026-08-29

### Breaking

- **`kernel.JoinKeys` 已删除**，替换为 **`kernel.JoinIdentifiers`**（行为等价：以 `.` 拼接贡献键段）。

### 迁移说明（下游升级必读）

1. **定位**：搜索代码中 `kernel.JoinKeys(` 的所有调用。
2. **改写**：`kernel.JoinKeys("a", "b")` → `kernel.JoinIdentifiers("a", "b")`（签名与返回类型不变）。
3. **重新构建**：`go build ./...`；若仍有 `undefined: JoinKeys` 报错 = 遗漏调用（`JoinKeys` 已从契约面移除）。
4. **回归**：组合根冒烟 + 探针。

### 语义说明

0.x 阶段 minor 版本可承载 breaking（changelog Breaking 节必写迁移说明）；`KernelAPIVersion` 主号（2）不变——契约兼容窗校验仍以 KernelAPIRange 为准（下游模块 `>=2.0 <3.0` 不受影响）。