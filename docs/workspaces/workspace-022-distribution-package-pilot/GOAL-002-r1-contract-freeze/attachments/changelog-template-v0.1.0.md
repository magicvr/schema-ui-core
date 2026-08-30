# Changelog 模板 v0.1.0

适用：契约冻结面（kernel / 模块 / npm 包组）发布记录。Each release section 由 R5 发布流水线生成骨架，人工补内容。

```markdown
# Changelog

与契约冻结面发布绑定（semver 语义见 semver-breaking-policy）；历史未发布期内容归入主线笔记，不追溯补造。

## [Unreleased]

### Added
- （新增能力 / 导出符号）

### Changed
- （行为扩展；标注影响面）

### Fixed
- （缺陷修复；无行为契约变化）

### Deprecated
- （冻结面符号进入弃用宽限期；注明移除版本计划）

### Removed（仅 major；必须逐条附迁移说明）
- （删除的符号/签名）

### Security
- （安全相关；与 VP-009 波次交叉引用）

## [2.<minor>.<patch>] - YYYY-MM-DD

### Breaking（major 版本专属节；每条必须含迁移说明）
- **迁移说明**：旧调用 → 新调用改写示例；受影响包/版本窗

### Added / Changed / Fixed / Deprecated / Security
（同上节模板）

## 版本窗声明

- kernel：`KernelAPIVersion` 主号 ↔ changelog major 对应；
- 模块：`Descriptor.Version`/`KernelAPIRange` 须与 changelog 一致；
- npm 包组（R3 后）：peer 依赖窗变更列于 Changed。
```

## 必填项检查（发布前）

- [ ] Unreleased 清空、版本节生成
- [ ] major 发布含 Breaking 节且每条有迁移说明
- [ ] kernel 主号 / 模块 KernelAPIRange / npm peer 窗三者一致
- [ ] golden consumer 升级演练结果回填（R4 后）