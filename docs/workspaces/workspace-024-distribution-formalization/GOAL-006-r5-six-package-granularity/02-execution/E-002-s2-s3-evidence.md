---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-006-r5-six-package-granularity
version: 0.2.0
---

# E-002 · S2/S3 六包形态细化实现与消费实证（2026-08-29）

## 交付物

- **构建链改造**（scripts/）：`build-lib-packages.mjs`（renderer 内部面重写插件 + external 函数；protocol/lib/ui 改 tsc 全产物——js+d.ts 子模块树，子路径可解析；入口元数据 ./ 合法化）；`rewrite-lib-aliases.mjs`（`@/` → `@magicvr/schema-ui-*` 全名映射 + ESM 扩展名转写，js+d.ts）；`finalize-lib-dist.mjs`（scope 全名统一 + 扩展名规范化 + JSON 资产面：schema/message 资产入包 + `with { type: "json" }` import attributes）。
- **六包终版**（npmjs @magicvr）：protocol **0.2.11** · lib **0.1.9** · renderer **0.3.7** · ui **0.1.7** · shell **0.1.2** · theme **0.1.2**。契约面：exports（"." 与 "./*" 均以 "./" 前缀合法）· peer 矩阵（renderer: react 系 + protocol/lib/ui；ui/shell: react 系；lib: react（i18n React 面））· dependencies（protocol: ajv；lib: clsx/tailwind-merge；ui: cva/lucide/radix-slot/clsx/tailwind-merge）。

## 证据链（判据 #5/#6）

| 判据 | 证据 | 结果 |
|------|------|------|
| #5 renderer 依赖图 external 化 | renderer 0.3.7 `index.js` = **187.5 kB**（旧 436.7 kB）· **17 处 `from "@magicvr/schema-ui-*"`**（protocol/lib/ui 子路径·消费端解析）· 产物 js+d.ts **0 处 `@/` 残留**；peer 矩阵声明 | ✅ |
| #5 冻结面 v1.4.0 | `attachments/freeze-face-v1.4.0.md`（六包导出面 + peer/deps 矩阵 + 版本终值） | ✅ |
| #6 ui 纯原子拆分 | `components/ui` = 12 原子组件（无 renderer/protocol/i18n 反向依赖断言）；业务组件（data-table + renderer 面）留在 renderer 包；ui 包独立消费（probe-six/ui 面） | ✅ |
| 消费实证 | golden-field 五探针**全绿**：probe-r5（external 化断言 PASS · 17 imports）/ probe（protocol 2.9）/ probe-render（1573B）/ probe-six（六包可消费）/ token-check（brand=2 ⊆ index=5）；安装 = npmjs 公开无凭据（空 userconfig） | ✅ |

## 历程与知识项（有界口径记录）

- 命名链：初产 `@schema-ui/*`（deprecated 全称）→ 发布实态定为 `@magicvr/schema-ui-*`（用户裁决）→ 构建/重写/finalize 三脚本统一为终名。
- Node ESM 三项铁律（脚本已固化为知识）：exports target 必须 `./` 前缀；相对/包子路径 import 必须带扩展名；JSON import 必须 `with { type: "json" }`（Node 20.10+）。
- 版本推进沉痛史（0.2.x/0.1.x 多轮）＝上述铁律逐一违反的学费；终版以 probe 全绿为准。

## 残余登记

1. **shell 类型面残余**：shell 包 d.ts 7 处 `@/account/*`、`@/host/*`、`@/components/data-table` 引用无对应包子面（host/account 不在六包）；JS 运行时自包含不受影响（五探针绿中 shell 消费正常）→ **消费端 tsc 类型面未验证，登记 R7 复核**（冻结面 §3 注记）。
2. **主题/字体资产面板**（如有后续渲染面）→ R7 复核清单。
3. 早期中间版本（0.3.0~0.3.6 / 0.2.1~0.2.10 / 0.1.1~0.1.8 命名/元数据修正链）已被终版取代（npmjs 保留历史版本；消费指引 = 六包终值）。