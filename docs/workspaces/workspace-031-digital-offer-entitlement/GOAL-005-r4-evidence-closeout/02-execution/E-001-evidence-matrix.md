---
doc_type: goal-execution
id: E-001-evidence-matrix
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: done
version: 1.0.0
---

# E-001 · C1 证据矩阵与边界核账

## 事实（时间线）

- 2026-09-05 · 证据矩阵构建：VP-031 方向级判据 1～8 逐条映射到代码/测试/审计证据（见下表）。
- 2026-09-05 · 边界核账（命令可复核）：
  - Charter 最后一次变更为 2026-09-01（`1694dea7`，本工作区开设前）——R1～R3 未改 Charter（判据 7）。
  - `apps/api/kernel/profile.go` 不含 `biz.digital-offer`——不进默认 Profile（判据 7）。
  - `apps/api/modules/digitaloffer/store/store.go` 无 `admin.users` 关联（仅 subjects 表引用）——判据 4/subject-only。
  - 迁移 0070 仅含 `digital_offers` / `digital_purchases` / `digital_entitlements` 三表——无类目/SKU/税/库存/物流订单/购物车域对象。
  - 构建与测试：在 Go module `apps/api/`（cwd）内执行 `go build ./...` 与 `go test ./...` 均通过（退出码 0；2026-09-05，本机；PostgreSQL 集成测试因 PG_TEST_* 已配置而实际运行）。注意仓库根目录无 go.mod——按字面在仓库根执行这两条命令不可复现，本矩阵统一以 apps/api module 为命令边界。

## VP-031 判据证据矩阵

| 判据 | 证据（代码 / 测试 / 审计） |
|------|----------------------------|
| 1 · Offer CRUD + Admin 协议页面 + 权限键 + 审计 + C 端上架列表 | 代码：`apps/api/modules/digitaloffer/provider.go`（路由/权限键 digitaloffer.read、offer.manage、entitlement.void/页面/manifest/导航）、`apps/api/internal/handler/digitaloffer.go`、`modules/digitaloffer/schema/*.json`；测试：`internal/handler/digitaloffer_test.go`（401 门控、创建/状态变更乐观锁、形态不可变 409、公开目录仅 on_sale 且无内部字段、审计失败回滚、公开目录限流 429）；审计：GOAL-003 A-002/A-008 |
| 2 · 余额不足拒绝、freeze→deduct_frozen、凭证与权益同事务或等价 fail-closed、失败 unfreeze 语义、并发测试 | 代码：`modules/digitaloffer/service/service.go`（§4.2/§4.4 单事务购买、互异幂等键、ref 反链、有界重试）；测试：`modules/digitaloffer/service/purchase_test.go`（余额不足零残留、并发双发恰一凭证/一次扣款/一份权益、余额恰够一次恰一成功无冻结残留、重试耗尽零残留、幂等重放、SQLite + 真 PostgreSQL 双库）；审计：GOAL-002 A-002/A-006、GOAL-003 A-002/A-004/A-008。失败 unfreeze 说明：单事务回滚使冻结从未对外可见（D-002 §4.2 更强 fail-closed 形态）；unfreeze 保留为 wallet Admin 纠错入口 |
| 3 · 权益有效/过期/耗尽可测，服务提供前统一核验 | 代码：`service.go` Check/Consume（§5.1 聚合、§5.2 算法）；测试：`check_consume_test.go`（聚合四态、duration 过期、耗尽、多行最旧优先、并发 void 线性化、SQLite/PG） |
| 4 · 购买与权益只挂 VP-029 subject_id，不创建 admin.users | 代码：store 仅引用 subjects；购买 tx 内 `SubjectExistsInTx` 门控；账户经 `GetOrCreateSubjectAccountInTx`；核账：store 无 users 关联 |
| 5 · Telegram 启用时注册可演示命令；未启用时测试不依赖 Bot API | 代码：`RegisterTelegram`（price/buy/entitlements）+ composition Disabled 注入；测试：Disabled no-op、真实 Dispatcher 回复（price 列表/buy 幂等/身份门控/entitlements 有效行） |
| 6 · 激活门禁留痕 | Root `00-meta.md` freshness 三字段（VRev-080 · 2026-09-05 · H-002 同进程 · RT-Q03/Q05 本波不需要 Redis） |
| 7 · 边界保持 | 边界核账（本文件时间线）；红线清单见 D-002 §10 |
| 8 · 开放 required = 0 | GOAL-002 A-007（A-006 independent 复核 pass）；GOAL-003 A-008（independent pass 0 required）；GOAL-004 A-003（A-002 pass 0 required）——关门审计（GOAL-005 A-001/A-002）最终确认 |

## 产物路径

- 本文件（证据矩阵 + 边界核账记录）

## 进度评估

- C1 关门；C2（关门审计与 Root 关门）待执行。
