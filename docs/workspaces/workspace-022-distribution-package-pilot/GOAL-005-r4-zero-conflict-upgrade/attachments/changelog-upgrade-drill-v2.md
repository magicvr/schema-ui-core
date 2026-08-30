# Changelog · R4 升级演练（V1 → V2 · 2026-08-29）

> 演练产物：演示「上游演进 → changelog 迁移说明 → 下游仅 bump + 按说明执行」的闭环。
> 基线 V1 = commit `8686b3fd`（62 条迁移）；V2 = 本演进 commit。

## [0.2.0] - 2026-08-29（演练版本）

### Added（additive · 无 breaking）

- **kernel（A 层）**：`JoinKeys(parts ...string) string`——贡献键拼接工具（module-id 风格）。
- **protocol（Web）**：`normalizePageID(id: string) string`——页面标识规范化（trim + lowercase）。
- **迁移 0063**：`site_settings_updated_at_index`（admin.settings）——`CREATE INDEX IF NOT EXISTS idx_site_settings_updated_at ON site_settings (updated_at)`；**双方言同一语句**（SQLite / PostgreSQL 均支持），无新表对象。

### Migration 说明（下游升级必读）

1. **已有数据库升级**：启动时全局迁移 runner 自动应用 0063（checksum 校验通过后）；无需人工 SQL。**唯一副作用** = `site_settings.updated_at` 上的索引（查询加速，无行为变化）。
2. **全新数据库**：fresh 初始化自动包含 0062+0063 全链。
3. **兼容性**：A 层/协议面均为 additive——不支持旧 API 破坏；下游既有调用零迁移成本。
4. **无配置键或依赖变更**（本演练样本覆盖 A 层 + 协议面 + 迁移三类，配置键/依赖留给 R5 发布回归补测）。

### 关联闸门

- 契约冻结面 v1.1.0 不变（additive 未触 breaking 流程）；本演练验证 semver-breaking-policy §3.3 的下游路径。