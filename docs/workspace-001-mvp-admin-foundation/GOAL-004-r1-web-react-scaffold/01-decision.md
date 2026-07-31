---
id: GOAL-004-r1-web-react-scaffold
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# 决策记录 · GOAL-004

## 信息需求与阶段门禁

| ID | 级别 | 最晚阶段 | 状态 | 阻断 |
|----|------|----------|------|------|
| I-004-001 | non-blocking | shadcn init 前 | open | 不阻断 Vite 初始化 |
| I-004-002 | non-blocking | 骨架冻结前 | open | 建议预建分层 |

服从父目标 [D-004](../GOAL-001-mvp-admin-foundation/01-decision.md)。

## D-001 · 前端骨架范围与平行仓复用边界

**日期**：2026-07-31  
**状态**：accepted

**决定**：

1. 目标路径：`apps/web/`；包管理 **npm + package-lock.json**。
2. 工具链：Vite + React + TypeScript；测试/lint 可沿用 Vitest + oxlint 或等价，R1 不强制完整测试矩阵。
3. UI：R1 **接入** Tailwind + shadcn/ui 基线 + 浅/深色最小占位（Charter 产品化方向）；Admin 外壳视觉完整度留给 R3。
4. 复用：参考 `../allinme.web-client`（`dev`）的 host runtime / protocol / renderer **边界思想**与 Vite 配置习惯；可择优移植无业务工具函数。
5. **明确不搬**：`mock-app` 业务演示、与特定业务 page 绑定的路由、生产 mock 排除脚本中的业务假设（除非改写为通用）。

**为什么**：

- 平行仓无 Tailwind/shadcn，与 Charter 不一致；R1 补基线可降低 R3 返工。
- 整仓拷贝会带入业务 demo 与自定义 CSS 双轨。

**未选方案**：

- **R1 仅 Vite+React，UI 库全放 R3**：短期更像平行仓，但 Charter 方向债后置。
- **沿用平行仓自研 tokens、不做 shadcn**：偏离 Charter，需 residual（用户未选）。
