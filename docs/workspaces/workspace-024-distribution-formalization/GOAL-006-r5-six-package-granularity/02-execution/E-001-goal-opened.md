---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-006-r5-six-package-granularity
version: 0.1.0
---

# E-001 · 目标建立（2026-08-29）

1. **立项**：承接 Root 纲领 R5（六包形态细化 · VP-024 判据 #5/#6 · go 后清单 renderer external + 纯原子拆分）；goal-tree 同步。
2. **勘察（设计依据）**：
   - renderer 内部面导入：`@/i18n` 16 · `@/protocol` 13 · `@/components/ui` 9 · `@/lib` 3 · `@/theme` 0（共 41 处子路径导入）；
   - `components/ui` = 12 文件纯原子组件（async-state/badge/button/card/input/label/skeleton/textarea/breadcrumbs）——无业务组件；
   - `dist-lib/@schema-ui/lib` 已含 `lib/` + `i18n/`（i18n 面归属 lib 包 ✓）；六包 files 声明为全目录（需收窄）。
3. **设计定档**：D-001（重写映射表 + exports/files + 版本推进 + ui 纯原子边界）。