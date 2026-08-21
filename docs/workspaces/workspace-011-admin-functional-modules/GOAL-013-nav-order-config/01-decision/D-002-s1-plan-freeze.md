---
id: D-002
goal: GOAL-013-nav-order-config
status: accepted
date: 2026-08-14
scope: S1 方案冻结
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · S1 方案冻结：默认清单 + YAML 覆盖

## 1. 用户裁决（2026-08-14）

- **默认清单顺位确认**：Dashboard → Users → Roles → Settings → Activity → Account → Notifications → File library → Data dictionary → System monitoring → Scheduled tasks → Recycle bin（12 项，产品级冻结）。
- **覆盖载体确认**：W7 YAML 配置段 `navigation.order`（configs/config.yaml 新增 navigation 小节），env 同名覆盖；无效配置回退默认 + 告警。

## 2. 默认清单（NodeID 映射）

| 顺位 | NodeID | 模块 | 现 Order |
|------|--------|------|----------|
| 1 | menu_dashboard | admin.dashboard | 0 |
| 2 | menu_users | admin.users | 1 |
| 3 | menu_roles | admin.roles | 2 |
| 4 | menu_settings | admin.settings | 1 |
| 5 | menu_activity | admin.activity | 2 |
| 6 | menu_account | admin.account | 1 |
| 7 | menu_notifications | admin.notifications | 2 |
| 8 | menu_files | admin.file-library | 3 |
| 9 | menu_dictionary | admin.data-dictionary | 4 |
| 10 | menu_monitoring | admin.system-monitoring | 5 |
| 11 | menu_scheduled_tasks | admin.scheduled-tasks | 6 |
| 12 | menu_recycle_bin | admin.recycle-bin | 8 |

实现为 kernel 包内常量 `DefaultNavigationOrder []string`（go 侧冻结），快照测试锁定。

## 3. 排序语义

`sortNavigation` 改为三层：

1. **清单顺位**：NodeID 在默认清单（或用户覆盖清单）中 → 按清单下标排；
2. **清单外兜底**：未列 NodeID 追加末尾（保持现 Parent → Order → NodeID 字母兜底，新模块不消失）；
3. Parent 分组仍优先（子节点跟随父节点分组，父节点按清单排）。

## 4. 覆盖载体（W7 YAML）

- configs/config.yaml 新增小节：

```yaml
navigation:
  # 全量顺序清单（NodeID 列表）；空 = 用内置默认清单
  order: []
```

- env 覆盖：`NAVIGATION_ORDER`（逗号分隔 NodeID，与 W7 env 覆盖规则一致：已设置优先于 YAML）。
- Config 新增 `NavigationOrder []string` 字段；kernel 排序接收该覆盖（nil/空 = 默认清单）。
- **非法配置**（未知 NodeID、格式错误、非字符串项）：整体回退默认清单 + 启动告警日志（不 fail-closed——排序是 UI 结构非安全门禁，且用户明确倾向回退+告警）。

## 5. 维护规则（I-003）

- 新模块落地时（波次流程）：维护者把新 NodeID 插入默认清单 + 更新快照测试（编译期无法强制，测试锁定兜底）。
- 未更新清单的新模块仍可见（追加末尾），不产生不可用状态。

## 6. 未选方案

| 方案 | 未选原因 |
|------|----------|
| 独立 NAVIGATION_ORDER_FILE | 用户确认走 YAML 配置段（与 W7 单配置源一致） |
| 运行时拖拽重排 | 用户已裁方案 A（重启生效） |
| 差异补丁式覆盖 | 全量清单语义更简单可预测（缺项追加末尾） |
| 非法配置 fail-closed | 排序非安全门禁；回退+告警更易运维 |

## 7. 影响面

- apps/api/internal/kernel（provider.go sortNavigation + 常量 + 签名扩展）
- apps/api/internal/config（Config.NavigationOrder + YAML 小节 + env）
- apps/api/configs/config.yaml + config.default.yaml（navigation 小节）
- 快照测试（kernel 或 composition 排序快照）
- **go（VP-008）判定**：manifest 导航内容扩展，非装配语义/非门禁语义 → 不 held（S4 确认）。
