---
title: I-006-002 · R3 真实菜单 features 投影矩阵
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-006-r3-persistent-rbac-menu
version: 0.1.0
related_info: I-006-002
related_decision: D-003
---

# I-006-002 · R3 真实菜单 features 投影矩阵

> **结论**：本附件与 D-003 关闭 `I-006-002`，作为 S5 的实施输入。当前真实 manifest 和 `/api/accounts/me` 尚未改动，以下矩阵不是运行结果。

## 1. 当前事实

- 真实 manifest 的 `list-edit-lifecycle` 页面已存在，路由为 `/list-edit-lifecycle`，schema 为 `/api/schema/list-edit-lifecycle`；它位于 sidebar 的 `Examples` group。
- 对应 checked-in schema fixture 包含 detail 与 edit form 组合，能代表需要写权限的管理页面；`settings` / `activity` 当前没有对应 checked-in schema fixture。
- `NavigationContext` 支持 `$context.features.<path>`；未知路径求值为 false，`visibleWhen` false 时移除项目，group 全部 children 被移除时会剪枝。
- API `Session.Features` 是 flat `map[string]bool`；Web 表达式中的点号按嵌套对象逐层查找，因此 flat key 不得使用点号表达命名空间。
- `/api/accounts/me` 当前对真实认证返回空 features map；菜单尚未由数据库 grant 投影。

证据：`apps/web/public/.well-known/schema-ui/app-manifest.json`；`apps/api/internal/handler/fixtures/schema/list-edit-lifecycle.json`；`apps/web/src/protocol/app-manifest.ts`（expression / visibility）；`apps/web/src/app/navigation.ts`（project/prune）；`apps/api/internal/account/session.go`；`apps/api/internal/handler/account.go`。

## 2. 固定映射

| 字段 | 值 |
|------|----|
| menu id | `menu-list-edit-lifecycle` |
| `page_ref` | `list-edit-lifecycle` |
| `feature_key` | `menu_list_edit_lifecycle` |
| manifest condition | `$context.features.menu_list_edit_lifecycle == true` |
| static location | `navigation.sidebar` → `Examples` → `list-edit-lifecycle` child |

真实 manifest 只增加：

```json
"visibleWhen": {
  "when": "$context.features.menu_list_edit_lifecycle == true"
}
```

页面声明、route、schemaUrl、label、icon 和 group 结构保持静态，不进入数据库。

## 3. 角色与响应矩阵

| 主体 | `records.read` | `records.write` | menu grant | `/me.features.menu_list_edit_lifecycle` | 导航结果 |
|------|----------------|-----------------|------------|-----------------------------------------|----------|
| admin | true | true | true | true | 显示 `List + edit` |
| viewer | true | false | false | false | 隐藏该 child；其它 Examples 项保留 |
| editor（R2 兼容） | true | false | false | false | 隐藏该 child；不扩大旧写权限 |
| 已认证但无角色 | false | false | false | false | 隐藏；records 读写均 `403` |
| 匿名 | — | — | — | 无 `/me` 响应（`401`） | AuthGate 阻断，不建立导航上下文 |

投影规则：对所有已登记 `menu_items.feature_key` 输出布尔值；仅当菜单 `enabled=1` 且当前用户任一角色存在 grant 时为 true，否则显式 false。多角色按 OR 合并。

## 4. 安全与可达性边界

- visibleWhen 只控制导航展示；直接访问 `/list-edit-lifecycle` 不得视为授权成功。
- 页面内 records GET/PATCH/DELETE 分别由 `records.read` / `records.write` 后端 gate 决定；前端隐藏不能替代 `401` / `403`。
- `overview` 保持 home/top 可达；`catalog` 保持 viewer 的明确只读入口。
- 该 child 隐藏后 `Examples` 仍有其它项目；空组剪枝继续由既有通用投影测试证明，不虚构真实 group 已变空。

## 5. 实施验证矩阵

| ID | 层 | 必须证明 |
|----|----|----------|
| `V-MENU-01` | store | admin/viewer/editor/无角色的 menu feature 投影分别为 true/false/false/false；多角色 OR 合并 |
| `V-MENU-02` | API | 真实 Bearer `/api/accounts/me` 保持 user 形状，并返回完整 bool features；匿名仍 `401` |
| `V-MENU-03` | manifest | checked-in manifest condition 通过解析，pageRef 仍解析到现有 page/route/schema |
| `V-MENU-04` | Web nav | admin 显示、viewer/editor 隐藏；declaration order 与其它 Examples children 不变 |
| `V-MENU-05` | fail closed | features 缺失、key 缺失、false 或错误类型均不显示该 child |
| `V-MENU-06` | projection | 既有空组剪枝测试继续通过；本真实 group 不被误删 |
| `V-MENU-07` | security | 深链不成为授权证据；records API 的 401/403/200 矩阵独立通过 |
| `V-MENU-08` | regression | `loadAccountContext` 映射、manifest/navigation tests 与生产 build 通过 |

建议复用 `navigation.test.ts`、`app-manifest.test.ts`、`context.test.ts`、`account_test.go` 和 `permission_test.go` 的现有断言模式；实际新增测试符号和通过结果在 S5 执行记录中留痕。
