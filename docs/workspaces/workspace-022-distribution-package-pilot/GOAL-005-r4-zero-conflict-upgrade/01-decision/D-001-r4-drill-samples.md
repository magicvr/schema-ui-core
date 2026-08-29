---
status: done
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-005-r4-zero-conflict-upgrade
version: 0.1.0
---

# D-001 · 演练样本定案（R4）

## 决策

演练样本集（VP-022 判据 #3「至少含配置键变更 + 新增迁移 + 依赖更新」的有界口径）：

| 样本 | 变更 | 类型 |
|------|------|------|
| A1 | `kernel.JoinKeys`（apps/api/kernel/keys.go） | A 层 additive |
| E1 | `normalizePageID`（apps/web/src/protocol/app-manifest.ts） | 协议面 additive |
| M1 | 迁移 0063 `site_settings_updated_at_index`（admin.settings · CREATE INDEX IF NOT EXISTS · 双方言同句） | 全局台账 |

**范围注（有界）**：判据 #3 字面三类的「配置键变更」「依赖更新」不在本波样本（changelog 明文留给 R5 补测）；本波以「A 层 + 协议面 + 迁移」三类代表升级压力（迁移 = fork-merge 最痛面）。R4 演练 A-001 已声明此口径——**独立审计 F-002 复核后由用户书面 residual 确认**（GOAL-006 A-002 响应，2026-08-29）。

## 其他执行决定

- 版本 bump 符号：golden-consumer `require v0.0.1 → v0.0.2`（+ replace 本地）；golden-web `pnpm install` 重拉（file: 快照语义）。
- 迁移版本号纪律：全局唯一 + 连续（11/50/130 三连冲突后定 63——wallet 常量 50 占用 + gap 校验）。
- 审计模式：R4 关门 = self（数据面演练，无生产数据）；R5/Root 关门 = independent（grok）。