---
status: active
created: 2026-08-29
updated: 2026-08-30
parent: GOAL-008-r7-topline-and-closeout
version: 0.1.0
---

# D-001 · 残余复核四项定档与置顶宣告（2026-08-29 · 用户书面口径经 Root 关门审计确认）

## 决策

1. **残余 1 · hosted CI 实触发（consumer-regression）**：`核销`（2026-08-30）——用户授权后初始化并推送（origin `main` 首推），hosted 首跑三连闭环：`33286154992`（action-setup 版本源 ❌）→ `33286191334`（探针/断言全绿但清理段 exit 1 ❌）→ `33286302663` **PASS**（1m9s · 四探针 + `shutdown.complete`）；A-002 F-005 → fixed。原「登记」口径的复审触发已满足，证据见 E-004。
2. **残余 2 · shell 类型面**：`登记`——4 文件 7 处 `@/account|@/host` 引用（JS 运行时自包含 · 五探针绿）；消费端 tsc 类型面未验证。**复审触发**：shell 包发布类型面进阶或消费者报告类型错误时。
3. **残余 3 · GH Packages 私有包**：`评述定稿`——**保留不删**（历史消费面 / VP-023 安装凭证链）；新消费一律 npmjs 公开（create/golden-field `.npmrc` 已钉）。
4. **残余 4 · C 类深度定制 fork 包化承载面**：`未来候选`——kernel 契约扩展通道（assembly 扩展 + 六包 external 组合）不在本 VP 边界；`migrate-fork` 对 C 型建议保持 fork（Charter fork 并存）。
5. **置顶宣告**：默认主路径 = 方法 B（cli+包）；fork 保持一等公民路径；QUICKSTART 按「决策块 + 方法 B 首节 + 迁移节 + 方法 A」置顶（已执行 v1）。

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| hosted 实跑作为判据 #3 硬门禁 | 登记 + 等价证据 | 用户环境事实（无 CI 槽位授权不负有限）；不把未做写成 acceptance |
| shell 类型面升级为 required | 登记 | 运行时消费链全绿（五探针 + UI-ONLY）；类型面为进阶面 |
| GH 私有包退役删除 | 保留 | 历史安装/凭证链不可逆破坏风险；npmjs 已是唯一新消费面 |