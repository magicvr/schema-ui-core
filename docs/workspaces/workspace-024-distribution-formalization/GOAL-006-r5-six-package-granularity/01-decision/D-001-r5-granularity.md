---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-006-r5-six-package-granularity
version: 0.1.0
---

# D-001 · 形态细化定档（2026-08-29）

## 决策

1. **renderer 依赖图 external 化（判据 #5）**：renderer 包由自包含 bundle 改为**显式依赖面**——重建时把内部面导入重写为包子路径并 external（消费端解析）：

   | src 引用 | 重写为 | external 声明 |
   |----------|--------|---------------|
   | `@/i18n/*`（16 处） | `@schema-ui/lib/i18n/*` | peerDependencies: `@schema-ui/lib` |
   | `@/protocol/*`（13 处） | `@schema-ui/protocol/*` | peerDependencies: `@schema-ui/protocol` |
   | `@/lib/*`（3 处） | `@schema-ui/lib/lib/*` | 同上（lib） |
   | `@/components/ui/*`（9 处） | `@schema-ui/ui/components/ui/*` | peerDependencies: `@schema-ui/ui` |
   | `@/theme/*`（0 处） | —（无依赖） | — |

   renderer peer 矩阵 = `react ^19` + `react-dom ^19` + `@schema-ui/{protocol,lib,ui}`（peer 声明，消费端解析；不重复 bundle）。
2. **六包 exports 子路径 + files 收窄**：各包 `exports: { ".": …, "./*": "./*" }`（子模块可解析）；`files` 收窄为各自实际目录（lib = index.js + lib/ + i18n/；protocol = index.js + protocol/ + lib/；ui = index.js + components/（ui）+ index + lib/（依赖工具若打包）——按产物实际目录校验后落）。
3. **版本推进**（npmjs 已发布旧版 → 新版本发布；changelog 注记「renderer 0.2.0 自包含 → 0.3.0 external 化（消费契约变化）」）：
   renderer **0.3.0** · protocol **0.2.1** · lib **0.1.1** · ui **0.1.1** · theme **0.1.1** · shell **0.1.1**。
4. **ui 纯原子拆分（判据 #6）**：断言 `components/ui` 无反向业务依赖（禁止 import `@/renderer` / `@/protocol` / `@/i18n`）+ ui 包独立消费（probe-six 重验证）；业务组件（form-controls 等 schema 驱动面）归属 renderer，**不**入 ui 包。
5. **冻结面 v1.4.0**：六包导出面（exports）+ peer 矩阵在附件定稿；V1.3.0 → v1.4.0 升格随本波。

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 全量子路径白名单 | 通配 `"./*"` | 包目录即产物面（冻结面 v1.4.0 逐项列白表以供审计）；通配降低维护 |
| renderer 双轨（dev=src/build=包） | B 路径单轨 | VP-023 教训：双轨漂移；构建期重写为包子路径即消费态 |
| 新建独立 forms 包 | 业务组件留 renderer | 判据 #6 只要求「业务组件出 ui 包」，非新建包 |