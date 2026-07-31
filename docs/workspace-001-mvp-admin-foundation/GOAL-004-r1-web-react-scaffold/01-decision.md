---
id: GOAL-004-r1-web-react-scaffold
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.0
---

# 决策记录 · GOAL-004

## 信息需求与阶段门禁

| ID | 级别 | 最晚阶段 | 状态 | 阻断 |
|----|------|----------|------|------|
| I-004-001 | non-blocking | shadcn init 前 | **verified**（D-002） | new-york 预设 |
| I-004-002 | **required** | 骨架目录冻结前 | **verified**（D-002） | 方案 B 预建分层 |

服从父目标 [D-004](../GOAL-001-mvp-admin-foundation/01-decision.md)；目录交界服从 [GOAL-002 D-002](../GOAL-002-r1-repo-layout-conventions/01-decision.md)。

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

## D-002 · 闭合 A-001 F-001：I-004-002 分层策略 + shadcn 验收硬度

**日期**：2026-07-31  
**状态**：accepted  
**响应**：independent A-001 · F-001（required）；顺带收紧 F-002 recommended 的 shadcn 措辞  
**用户意图**：`/govern` 明确要求闭合 I-004-002 分层策略后推进 R1

**决定**：

1. **I-004-002** 级别升为 **required** @ **骨架目录冻结前**（与首次目录落盘同一门禁）。
2. **分层策略 = 方案 (B)**：R1 预建空目录  
   - `src/host/`  
   - `src/protocol/`  
   - `src/renderer/`  
   各目录仅放边界说明（如 `README.md` 或 `.gitkeep` + 根级一句），**不**实现 runtime/协议解析/业务页。
3. 另保留 `src/app/`（应用壳入口/页面占位）与 `src/components/ui/`（shadcn）。
4. **I-004-002 → verified**（本决策 + 后续目录证据）。
5. **I-004-001**：shadcn 预设 **new-york**，CSS variables 策略（shadcn 默认）→ **verified**。
6. **shadcn 验收**（A-001 F-002 recommended）：必须可指回初始化痕迹（`components.json` 与/或文档记载的 init 命令 + `components/ui`）；禁止用无工具链的手写 div 冒充 shadcn 基线。
7. **R1 vs R3**（A-001 F-004）：R1 通过条件**不含** App manifest 导航壳、多业务路由；单页/占位页 + 主题切换即可。
8. **目录交界**：若 `apps/web` 已有空壳，原地填充，不改路径（GOAL-002 D-002）。

**为什么**：

- 预建空分层防 R3/R5 引入 protocol/renderer 时大挪（meta 原建议）；比扁平-only 更贴平行仓边界思想。
- 「或等价」过宽会与 Charter/D-004 摩擦。

**未选方案**：

- **(A) R1 仅 `src/app`（+components/ui）**：短期更干净，但后续分层迁移成本高；用户要求闭合分层且 meta 倾向 B。
- **过深实现 host runtime 于 R1**：越 R1 范围，吞并后续阶段。

**影响**：

- 放行 `apps/web` 骨架实施；F-001 → `fixed`；I-004-002 verified。
- 不放行 Admin 外壳（R3）或协议 Renderer 全量（R5）。
