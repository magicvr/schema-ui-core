---
id: E-004
goal: GOAL-013-nav-order-config
date: 2026-08-14
status: recorded
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-004 · S3 验证完成

## 事实

- 2026-08-14：S3 验证完成。

### 单元测试

- kernel：TestDefaultNavigationOrderSnapshot（12 项冻结）+ TestSortNavigationDefaultOrder / CustomOrder / InvalidOrderFallsBackToDefault / PartialCustomOrderAppendsRest 全过。
- config：TestLoadNavigationOrder 5 项（YAML 列表 / 空=默认 / env 覆盖 / 非列表回退 / 非字符串回退）全过。
- manifest：TestSortNavigationOrdersSlots（features.menu_* 提取 + 未列项追加末尾）过。
- composition / 全量 go test ./... 全绿；web vitest 903 全绿（导航顺序不影响静态 fixture 测试）。

### 覆盖路径实测（真实进程，admin profile）

| 场景 | 配置 | sidebar 结果 |
|------|------|-------------|
| 默认清单 | 无覆盖 | Dashboard | Users | Roles | Activity | File library | Data dictionary | System monitoring | Scheduled tasks | Recycle bin（user 槽：Settings | Account） |
| env 覆盖 | NAVIGATION_ORDER=menu_recycle_bin,menu_dashboard | Recycle bin | Dashboard | 其余字母序尾部 |
| 非法覆盖 | NAVIGATION_ORDER=menu_bogus | 回退默认清单 + WARN（normalize 在 composition 层统一处理 kernel+manifest 两处） |

### S2 补丁（E-003 后）

- 实施中发现：非法覆盖在 kernel 层回退，但 plan.NavigationOrder 原样传给 manifest → manifest 层未回退（排序退化为字母序）。修复：kernel 导出 NormalizeNavigationOrder(order, known)；composition 在 manifest 构建前用注册 NodeID 集规范化一次，两处共用同一结果。

## 遗留

- S4 go 判定 + 自审；S5 grok 关门审计。
