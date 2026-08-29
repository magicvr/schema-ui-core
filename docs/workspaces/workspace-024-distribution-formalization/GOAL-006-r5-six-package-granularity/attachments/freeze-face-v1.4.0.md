# 冻结面 v1.4.0 · 六包导出面与 peer 矩阵（2026-08-29 · GOAL-006 定稿）

> 继承 freeze-face v1.3.0（六包 registry 发布 · d.ts 管线）。v1.4.0 = renderer 依赖图 external 化后的**消费契约面定稿**（VP-024 判据 #5）。

## 1. 导出面（exports）

| 包 | 版本 | 入口 | 子路径 | files |
|----|------|------|--------|-------|
| `@schema-ui/protocol` | 0.2.1 | `"types": ./protocol/index.d.ts · "import": ./index.js` | `"./*": "./*"`（protocol/ · lib/） | index.js + protocol/ + lib/ |
| `@schema-ui/lib` | 0.1.1 | ./lib/index.d.ts · ./index.js | `"./*": "./*"`（lib/ · i18n/） | index.js + lib/ + i18n/ |
| `@schema-ui/theme` | 0.1.1 | ./theme/index.d.ts · ./index.js | `"./*": "./*"`（theme/） | index.js + theme/ |
| `@schema-ui/ui` | 0.1.1 | ./components/ui/index.d.ts · ./index.js | `"./*": "./*"`（components/ · lib/ · i18n/） | index.js + components/ + lib/ + i18n/ |
| `@schema-ui/renderer` | **0.3.0** | ./renderer/index.d.ts · ./index.js | `"./*": "./*"`（renderer/） | index.js + renderer/ |
| `@schema-ui/shell` | 0.1.1 | ./app/index.d.ts · ./index.js | `"./*": "./*"`（app/ host/ account/ …） | index.js + 全部目录 |

## 2. peer 矩阵（消费端解析 · R5 定稿）

| 包 | peerDependencies | 说明 |
|----|------------------|------|
| protocol | — | 自包含 bundle（ajv 打包） |
| lib | — | 自包含 |
| theme | — | Token CSS/TS，无运行时依赖 |
| ui | react ^19 · react-dom ^19 | 原子组件；lib/i18n 面 bundle 自含（无包级依赖） |
| renderer | react ^19 · react-dom ^19 · **@schema-ui/protocol ^0.2.1 · @schema-ui/lib ^0.1.1 · @schema-ui/ui ^0.1.1** | 依赖图 external 化：产物 import 包子路径（17 处），消费端解析；不再自包含 protocol/lib/ui/i18n 面 |
| shell | react ^19 · react-dom ^19 | App/Host 壳 bundle 自包含 |

## 3. 契约记录

- **renderer 0.2.0（自包含）→ 0.3.0（external 化）为消费契约变化**（changelog 迁移说明：消费方安装 renderer 新版时需同时具备 protocol/lib/ui peer；npm 8+/pnpm 自动解析 peer）。
- **d.ts alias 重写**：六包产物 d.ts 的 `@/` 引用已按映射表重写为包子路径（protocol 1 · ui 1 · shell 12 处）；**shell 类型面残余 7 处引用（`@/account/*`、`@/host/*`、`@/components/data-table` 无对应包）→ 消费端 tsc 类型面未验证，登记 R7 复核**（JS 运行时自包含不受影响）。
- 主仓源码（`apps/web/src`）不受影响（重写仅作用于 dist-lib 产物）；`tsc -b` 回归绿。

## 4. 判定

- 判据 #5（renderer 依赖图 external 化 · 六包 peer 矩阵定稿）：**满足**。
- 判据 #6（ui 纯原子拆分 · 独立消费）：ui 包 = 12 个原子组件（无 renderer/protocol/i18n 反向依赖）+ 独立消费实证（见 GOAL-006 E-002）；业务组件（form-controls 等 schema 驱动面）留在 renderer ——**满足**。