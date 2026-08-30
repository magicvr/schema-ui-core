---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-web-package-consumption
version: 0.1.0
---

# D-001 · Web 拆包路径与 pin 漂移处置（用户裁决）

## 裁决 1 · 拆包实施路径

用户 2026-08-29 书面裁决：**B · Vite lib 产物打包**（不改源码，新增 lib 构建入口产出可发布 ESM+d.ts；先发 protocol 首包；A 源码 monorepo 化留 go 后正式化评估）。

执行：vite.lib.config.ts + tsconfig.protocol.json + `@schema-ui/protocol` v0.1.0 产物 + golden-web 消费验证（E-002）。renderer/shell/ui/theme 包化沿用同一链路（S3 增量）。

## 裁决 2 · 协议 pin 漂移（I-007）

用户 2026-08-29 书面裁决：**登记，留 `/vision` 后续裁决**——本目标内不修改 Charter/roadmap pin（保持 `v2.8.0` 文档事实），不自行推翻代码 2.9 事实；R4 演练基线选择与 R5 发布前提将由 `/vision` 统一审视（pin bump 需 provenance/容器门禁链同步）。

## 未选方案

- 不选 A（源码 monorepo 化）：改动面大，试点期不必要；留 go 后。
- 不选 C（仅 protocol 最小切片）：渲染闭环不达判据 #2；本路径以 B 为骨架并在 S3 增量覆盖。
- 不立即启动 `/vision` pin 升级：避免试点轮期间愿景层大动作；登记触发点 = R4 前。