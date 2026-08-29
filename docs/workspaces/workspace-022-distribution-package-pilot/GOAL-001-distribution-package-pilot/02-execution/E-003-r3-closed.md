---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-001-distribution-package-pilot
version: 0.1.0
---

# E-003 · R3 完成事实（2026-08-29）

1. **Web 包分发链路**（B 路径 · 用户裁决）：Vite lib + declaration 双链路（protocol / renderer 两个 lib 配置 + 两个 tsconfig）；`.gitignore` 排除产物。
2. **首包**：`@schema-ui/protocol` v0.1.0（306 kB 自包含 · 144 modules）——manifest 协商/校验/URL 解析/表达式。
3. **渲染包**：`@schema-ui/renderer` v0.1.0（436.7 kB · 1653 modules · React external 为 peer）——RenderPage + I18nProvider + registerCustomComponent。
4. **golden-web 三探针 PASS**：protocol 功能断言 · SSR 渲染真实形态 schema 文档（能力门控 fail-closed 可观测）· Token 覆盖纪律（brand ⊆ index）。
5. **契约**：Web 冻结面候选（六包边界 + peer 矩阵 v0.1 设计）→ go 后升格（v1.2）；F-006（d.ts TS5056）与 I-007（pin 漂移）登记。
6. **判据 #2 满足声明**；Root progress 2/5 → 3/5；GOAL-004 done 4/4。

下一步：**R4 立项（零冲突升级演练）**——上游真实演进样本 → golden-consumer/golden-web 仅 bump + 迁移说明 → 回归全绿、冲突 0；含 PG external 核销（F-005）、pin 漂移 `/vision` 处理（I-007）、A-001 F-003（PG drain 时序）复核。