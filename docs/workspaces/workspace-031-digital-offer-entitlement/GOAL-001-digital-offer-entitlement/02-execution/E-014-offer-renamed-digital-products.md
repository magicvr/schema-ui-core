---
doc_type: goal-execution
id: E-014-offer-renamed-digital-products
parent: GOAL-001-digital-offer-entitlement
date: 2026-09-06
status: done
version: 1.0.0
---

# E-014 · 关门后维护：「数字 Offer」显示命名改为「数字商品」

## 事实（时间线）

- 2026-09-06 · **用户意见**：「数字 Offer」混排英文，不符合中文直觉，询问能否改为更贴中文语境的说法。
- 2026-09-06 · **用户裁决（P-004）**：选「数字商品」。链路语义：数字商品 → 数字订单 → 数字权益。
- 2026-09-06 · **边界界定**：只改**用户可见显示层**（i18n 值 + schema/manifest 的 fallback label + 导航 label）；**内部标识不动**——pageId `digitaloffer-offers`、路由、`menu_digitaloffer_*`、权限键 `digitaloffer.*`、审计事件 `bizoffer.offer.*`、表/字段 `digital_offers`/`offer_id`、i18n key 名、errorcatalog code 均保持（涉及数据库/合同/既有数据，改即破坏性变更）。
- 2026-09-06 · **实施**：
  - zh-CN：manifest title/nav「数字商品」；搜索「搜索商品」；字段/列「商品 ID / 商品 / 商品名称」；按钮「新建商品 / 创建商品 / 保存商品」；错误提示「商品不存在 / 商品未上架 / 商品币种… / 请求标识已被其他商品使用 / 商品已被并发修改 / 提交的商品信息无效…」。
  - en-US 同步 product 系列（Digital products / Product not found / New product…），中英对照一致。
  - schema/manifest fallback label（Digital products / New product / Create product / Save product）与 `provider.go` 导航 label 同步。
  - 键集双文件一致（1117 键），结构性测试通过。
- 2026-09-06 · **验证**：JSON 校验通过；i18n 键集一致性通过；API kernel/composition/digitaloffer/handler 相关测试全绿；web 全包 vitest 后台跑完登记结果。

## 产物路径

- `apps/web/src/i18n/messages/{zh-CN,en-US}.json`
- `apps/api/modules/digitaloffer/schema/digitaloffer-offers.json`、`manifest/fragment.json`、`provider.go`

## 进度评估

- 关门状态不变（Root done 4/4 · VP-031 closed v0.3.4）。纯显示层重命名，无数据/契约/路由变化；审计模式 none。
- 遗留说明：D-002 合同与代码内部仍用 "offer" 术语；如需整体改内部标识（路由/权限/事件/字段）属破坏性迁移，需另行立项并先评估既有数据与外部调用方。
