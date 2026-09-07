---
id: GOAL-001-digital-offer-entitlement
title: 数字 Offer 与权益
status: active
parent: null
created: 2026-09-05
updated: 2026-09-06
version: 0.1.0
---

# GOAL-001-digital-offer-entitlement · 02-execution 索引

| id | date | scope | summary | status |
|----|------|-------|---------|--------|
| [E-001-workspace-establishment](02-execution/E-001-workspace-establishment.md) | 2026-09-05 | 开区建立 | VP-031 激活投影 + workspace scaffold + Root 五件套 + goal-tree；不含业务实现 | active |
| [E-002-r1-subgoal-and-adjudication](02-execution/E-002-r1-subgoal-and-adjudication.md) | 2026-09-05 | R1 启动 | 创建 GOAL-002-r1-contract-freeze（C1）；用户裁决 I-031-001～003；D-001/D-002 落盘；Root/VP-031 台账回写 | done |
| [E-003-r1-closure](02-execution/E-003-r1-closure.md) | 2026-09-05 | R1 关门 | D-002 v1.0.0 accepted（R2 期间两次附录至 v1.2.0）；A-001～A-007 审计闭合（A-006 independent pass）；GOAL-002 done 3/3；Root 1/4 | done |
| [E-004-r2-closure](02-execution/E-004-r2-closure.md) | 2026-09-05 | R2 关门 | GOAL-003 done 4/4：实施 + 双库验收 + A-001～A-008 审计循环（A-008 independent pass 0 required）；Root 2/4 | done |
| [E-005-root-closure](02-execution/E-005-root-closure.md) | 2026-09-05 | R4 关门与 Root 关门 | GOAL-005 done 2/2（证据矩阵 + 关门审计 A-001～A-005，A-004 independent pass 0 required）；Root done 4/4；VP-031 closed | done |
| 第 2 次关门（GOAL-005 [E-007-reclose](../GOAL-005-r4-evidence-closeout/02-execution/E-007-reclose.md)） | 2026-09-05 | 第 2 次关门 | A-006 运行时集成 fail 整改（BuiltinModules 注册 / D-003 裁决 / 组合根验收）经 A-010 independent pass 0 required 确认；Root done 4/4；VP-031 closed v0.3.2 | done |
| 第 3 轮撤回（GOAL-005 [E-008-a012-response](../GOAL-005-r4-evidence-closeout/02-execution/E-008-a012-response.md)） | 2026-09-05 | A-012 关门撤回与证据补齐 | A-012 independent `conditional` 2 required（F-007/F-008 证据缺口）→ D-004 fixed ×2 + A-013 closed ×2（迁移失败/reopen 测试 SQLite+PG；Telegram-enabled 组合根 + 结构化 Manifest）；Root → active 3/4，待 focused independent closure 复审后重新关门 | done |
| 第 3 次关门（GOAL-005 [E-009-a014-response-and-reclose](../GOAL-005-r4-evidence-closeout/02-execution/E-009-a014-response-and-reclose.md)） | 2026-09-05 | 第 3 次关门 | A-014 independent focused finding-closure `pass` 0 required（F-007/F-008 关闭证据成立，真实 PG 未 skip）→ 用户裁决先处理 F-001（recommended）→ A-015 fixed（结构化 Manifest + dispatcher 指针同一断言）→ GOAL-005 done 2/2、Root done 4/4、VP-031 closed v0.3.4 | done |
| 关门后维护（[E-010-post-closure-maintenance](02-execution/E-010-post-closure-maintenance.md)） | 2026-09-06 | 数字权益页报错 + 侧栏图标 + 菜单位置 | entitlements 页 `type: confirm` 非冻结 action 类型导致 D-VAL fail-closed → 改为行操作 confirm + requestMapping（fixed）；shell iconRegistry 注册 `check-badge`（fixed）；DefaultNavigationOrder 将 digitaloffer 两入口移至预付凭证后、操作日志前（fixed）；`go test ./...` + web vitest 1218 全绿 | done |
| 关门后维护（[E-011-create-offer-500-and-i18n](02-execution/E-011-create-offer-500-and-i18n.md)） | 2026-09-06 | 新建 Offer 500 + 中英文对照 | 根因：`operation_log.event` CHECK 白名单未含 `bizoffer.*` 审计事件，生产 `CreateOffer` 事务审计写入失败回滚 → 500（测试因 stub recorder 未触达）。migration 0071 扩 CHECK 枚举（fixed）；zh-CN 字段描述贴合日常直觉（fixed）；全量测试绿 + 实测 201 + 审计行落盘 | done |
| 关门后维护（[E-012-offer-price-yuan-input](02-execution/E-012-offer-price-yuan-input.md)） | 2026-09-06 | Offer 标价以元输入 | inputNumber 新增 Host-local `unit:"yuan"`：显示/输入为元（step 0.01，提交取整分），wire 仍为分——创建/编辑/列表/存储/契约全不变；创建与编辑弹窗 priceAmount 启用；zh「标价（元，最多两位小数）」；web 1219 测试全绿 | done |
| 关门后维护（[E-013-digital-orders-page](02-execution/E-013-digital-orders-page.md)） | 2026-09-06 | 补齐「数字订单（购买记录）」页面 | 用户裁决补三段链路：新增只读 digitaloffer-purchases 页（复用既有 API/权限，纯 schema）+ 权益页加 purchaseId 列；descriptor/manifest/nav/icon/i18n 同步；无新表/路由/权限、不改契约 | done |
| 关门后维护（[E-014-offer-renamed-digital-products](02-execution/E-014-offer-renamed-digital-products.md)） | 2026-09-06 | 「数字 Offer」显示命名 →「数字商品」 | 用户裁决改显示层：zh「数字商品」+ en「Digital products」系列（按钮/字段/错误提示同步）；内部标识（pageId/路由/权限/事件/表字段/i18n key）不变；键集一致 1117 | done |
| 关门后维护（[E-015-systemdata-version-bump](02-execution/E-015-systemdata-version-bump.md)） | 2026-09-06 | dev.cmd start 启动失败修复 | 根因：E-014 改导航 Label 未递增 SystemDataVersion → 持久库 `menu_digitaloffer_offers` ledger checksum 漂移 → reconcile fail-closed 拒启。修复：SystemDataVersion 1→2（设计机制）；dev.cmd start 全流程验证通过后 stop 交还 | done |

## 执行记录（ledger）

`02-execution/` 平铺；编号递增；时间线只记事实。
