# 冻结面 v1.4.0 · 六包导出面与 peer 矩阵（2026-08-29 · GOAL-006 定稿）

> 继承 freeze-face v1.3.0（六包 registry 发布 · d.ts 管线）。v1.4.0 = renderer 依赖图 external 化后的**消费契约面定稿**（VP-024 判据 #5）。

## 1. 导出面（exports）

| 包 | 版本 | 入口 | 子路径 | files |
|----|------|------|--------|-------|
| `@schema-ui/protocol` | 0.2.11 | `"types": ./protocol/index.d.ts · "import": ./protocol/index.js` | `"./*": "./*"`（protocol/ · lib/ · schemas/） | 全量（无 files 白名单） |
| `@schema-ui/lib` | 0.1.10 | ./lib/index.d.ts · ./lib/index.js | `"./*": "./*"`（lib/ · i18n/ · messages 资产） | 全量 |
| `@schema-ui/theme` | 0.1.4 | ./theme/index.d.ts · ./index.js | `"./*": "./*"`（theme/） | 全量 |
| `@schema-ui/ui` | 0.1.8 | ./components/ui/index.d.ts · ./components/ui/index.js | `"./*": "./*"`（components/ · lib/ · i18n/） | 全量 |
| `@schema-ui/renderer` | **0.3.8** | ./renderer/index.d.ts · ./index.js | `"./*": "./*"`（renderer/ · components/） | 全量 |
| `@schema-ui/shell` | 0.1.4 | ./app/index.d.ts · ./index.js | `"./*": "./*"`（app/ host/ account/ …） | 全量 |

> 终值版本 = npmjs 实发（命名/元数据/运行链全绿后定稿）；早期中间版（0.3.0~0.3.6 等修正链）为历史记录，消费指引以终值为准。

## 2. peer 矩阵（消费端解析 · R5 定稿）

| 包 | peerDependencies | dependencies | 说明 |
|----|------------------|---------------|------|
| protocol | — | ajv ^8.20 | 协议面；schema JSON 资产入包（import attributes） |
| lib | react ^19（i18n React 面） | clsx · tailwind-merge | i18n/messages 资产入包 |
| theme | — | — | Token CSS/TS |
| ui | react ^19 · react-dom ^19 · **@magicvr/schema-ui-lib ^0.1.10（breadcrumbs i18n 面经 lib peer）** | clsx · tailwind-merge · class-variance-authority · lucide-react · @radix-ui/react-slot | 原子组件 + DataTable 设计系统面（data-table 归 ui · 用户 P-004 裁决）；独立消费实证 = UI-ONLY（仅装 ui+peer · exports=18 PASS） |
| renderer | react ^19 · react-dom ^19 · **@magicvr/schema-ui-protocol ^0.2.11 · @magicvr/schema-ui-lib ^0.1.10 · @magicvr/schema-ui-ui ^0.1.8** | — | 依赖图 external 化：产物 import 包子路径（17 处），消费端解析；不再自包含 protocol/lib/ui/i18n 面 |
| shell | react ^19 · react-dom ^19 | — | App/Host 壳 bundle 自包含 |

## 3. 契约记录

- **renderer 0.2.0（自包含）→ 0.3.0（external 化）为消费契约变化**（changelog 迁移说明：消费方安装 renderer 新版时需同时具备 protocol/lib/ui peer；npm 8+/pnpm 自动解析 peer）。
- **d.ts alias 重写**：六包产物 d.ts 的 `@/` 引用已按映射表重写为包子路径（protocol 1 · ui 1 · shell 12 处）；**shell 类型面残余 4 文件 7 处引用（`@/account/*`、`@/host/*` 无对应包）→ 消费端 tsc 类型面未验证，登记 R7 复核（GOAL-008 D-001 残余 2）**（JS 运行时自包含不受影响；`components/data-table` 不在此列——属 ui 设计系统面）。
- 主仓源码（`apps/web/src`）不受影响（重写仅作用于 dist-lib 产物）；`tsc -b` 回归绿。

## 4. 判定

- 判据 #5（renderer 依赖图 external 化 · 六包 peer 矩阵定稿）：**满足**。
- 判据 #6（ui 纯原子拆分 · 独立消费）：ui 包 = 设计系统面（原子组件 12 + DataTable 通用数据表——归 ui 为用户 P-004 裁决；breadcrumbs 的 i18n 经 `@magicvr/schema-ui-lib` peer）；业务组件（form-controls 等 schema 驱动面）留在 renderer ——**满足**（UI-ONLY 独立消费实证）。