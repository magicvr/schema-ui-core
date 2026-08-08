---
id: GOAL-002-s1-design-tokens-and-primitives
doc: decision
status: active
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# 决策记录 · GOAL-002-s1-design-tokens-and-primitives

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | S1 实施输入齐备（D-002/D-003/D-004 + I-S1-001） | C1–C6 | 实施前 | 读 Root 决策链与基线盘点 | **closed** | — | Root 00-meta I-001/I-002/I-005 closed（2026-08-09） |
| I-002 | required | 深浅色 headless 验证环境可用 | C6 | C6 前 | Playwright chromium 检查 | **closed** | — | chromium-1234 已装（2026-08-09） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| — | — | （本目标不产生新决策；实施按 Root D-002/D-003/D-004 accepted 方案执行；若实施中发现方案外取舍，以 D-001 落盘） | — | — |

## 方案要点（引用 Root 决策，非新决策）

1. **Token 权威**：`apps/web/src/index.css` 唯一运行时权威；`:root`/`.dark` 存原始语义值，`@theme inline` 仅 alias（F-005 双层纪律）。
2. **F-002**：原始语义 `--elevation-sm|md|lg` + `@theme` `--shadow-sm|md|lg: var(--elevation-*)`；禁止同名自引用；闭合条件见 D-003 §2（a–d）。
3. **Typography**：`--font-sans`/`--font-mono`；字阶依赖 Tailwind 默认 `text-*` scale（F-006），禁止硬编码 px。
4. **FOUC**：`index.html` 同步内联引导脚本 + 可单测主题单元（class + `localStorage.theme` + `color-scheme`）。
5. **消费闭环**（F-003）：success/chart/overlay/shadow 硬编码点迁移语义 Token；不进 Token 的写边界理由。
