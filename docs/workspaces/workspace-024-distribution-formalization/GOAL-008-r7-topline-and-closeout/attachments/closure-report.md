# 收口报告 · VP-024 分发形态正式化（2026-08-29 · GOAL-008 S2）

> 面向现有与潜在下游的官方收口声明：**cli+包 消费 = schema-ui-core 默认主路径**（方法 B）；fork 为深度定制/贡献者路径（方法 A），并给 A/B 型 fork 提供迁移工具。Charter 0.3.0 措辞不变（fork 与包消费并存）。

## 1. 判据核销表（VP-024 八条退出判据 → 证据）

| # | 判据 | 状态 | 证据（本波次内可复核路径） |
|---|------|------|---------------------------|
| 1 | serve 壳闭环 | ✅ R1 | GOAL-002 done 5/5 · `apps/api/server` 公开面（config/serve/测试 13+2）· RT-D02 全序停机 |
| 2 | 公开发布通道 | ✅ R2 | GOAL-003 done 4/4 · npmjs.com `@magicvr/schema-ui-*` 六包实发 · golden-field 免凭据消费 · scope 定稿（用户裁决） |
| 3 | compose/CI 实跑 | ✅ R3（有界） | GOAL-004 done 4/4 · compose 全服务 + consumer-regression 免凭据重构 + linux 容器 harness A/B；hosted 实触发见 §3-1 |
| 4 | fork 对照计时 | ✅ R4 | GOAL-005 done 4/4 · v0.3.0→v0.4.0 同一演进集：fork 1 冲突 + 2 改写点 + ≈13.2s vs 包 0 冲突 + ≈4.8s（fork-comparison-report） |
| 5 | 六包形态细化（renderer external 化 · 冻结面 v1.4.0） | ✅ R5 | GOAL-006 done 5/5 · renderer 187.5kB · 17 处包子路径 import · peer 矩阵实发 · 冻结面 v1.4.0 · 五探针 + UI-ONLY 全绿 |
| 6 | 纯原子拆分 | ✅ R5 | ui 包 = 12 原子组件设计系统面（data-table 归 ui · 用户 P-004 裁决）· 独立消费实证 |
| 7 | 迁移工具化 | ✅ R6 | GOAL-007 done 4/4 · `schema-ui migrate-fork`（A/B/C 判定 · 非破坏 · 9510023 实测 v0.3.0→v0.4.0） |
| 8 | 方法 B 置顶 + 收口报告 | ✅ R7（本文件） | QUICKSTART 首段 = cli+包（决策块 + 方法 B + 迁移节；fork 节顺延）· 本核销表 |

## 2. 公开消费往返实证（end-to-end）

- **生产面**：v0.4.0 tag（`00d97b5b`）+ npmjs 六包终值（protocol 0.2.11 · lib 0.1.10 · renderer 0.3.8 · ui 0.1.8 · shell 0.1.4 · theme 0.1.4）。
- **消费面**（golden-field · 无凭据）：`pnpm install`（npmjs 公开 registry）→ 五探针全绿（external 断言 17 imports / protocol 2.9 / render 1573B / six / token）→ `UI-ONLY`（仅装 ui+peer）PASS。
- **升级往返**：`schema-ui upgrade` = go get @latest + pnpm add @latest + 探针回归（R2/R6 实操）；迁移 = `schema-ui migrate-fork`（R6 实操）。
- **计时口径**：create → 双端绿分钟级；升级秒级；冲突 0（R4 实证）。

## 3. 残余复核清单（四项定稿）

| 项 | 结论 | 处置 |
|----|------|------|
| 1 · hosted CI 实触发（consumer-regression workflow） | 本地等价 + linux 容器 harness A/B 已证；hosted 触发需 CI 槽位授权（用户环境事实） | **登记**：触发脚本/流程已就绪（workflow 已提交）；随用户下一次 CI 槽位授权执行（触发 = `workflow_dispatch`） |
| 2 · shell 类型面（4 文件 7 处 `@/account|@/host`） | JS 运行时自包含（五探针绿中 shell 消费正常）；消费端 tsc 类型面未验证 | **登记**：R7 后候选（shell 包声明类型面进阶或登记为已知面）；不影响运行时消费 |
| 3 · GH Packages 私有包 | 历史消费面（VP-023 时代安装/凭证），保留不删；**新消费一律 npmjs 公开**（F-004 已钉） | **评述定稿**：双 registry 并存策略（GH=私有历史面 · npmjs=公开现面） |
| 4 · C 类深度定制 fork 的包化承载面 | kernel 契约扩展通道（assembly 扩展 + 六包 external 组合）不在本 VP；`migrate-fork` C 型建议保持 fork | **登记**：未来候选（新 VP/波次）；Charter fork 并存维持 |

## 4. 置顶宣告

- **默认主路径 = 方法 B（cli+包）**：`go install …/schema-ui@latest` → `schema-ui create` → `schema-ui upgrade`（零冲突升级）。
- **fork 保持为一等公民路径**（深度定制/贡献者），A/B 型用户由 `schema-ui migrate-fork` 平滑迁移。
- QUICKSTART 已按此结构置顶（决策块 30 秒定位）；Charter 措辞不变。

## 5. 历史脉络

VP-022（哔包试点 go/no-go）→ VP-023（产线化 · 六包冻结面 v1.3.0）→ VP-024（对外正式化：serve 壳 · 公开发布 · CI 实跑 · 对照实证 · 形态细化 · 迁移工具化 · 置顶收口）；全部阶段级经 grok build（grok-4.6 · high）独立审闭环，required findings 三路径闭合。Root 关门审计 = GOAL-008 A-002（independent）。