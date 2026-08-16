---
id: GOAL-013-w12-product-surface-intent
doc: execution
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.3.0
---

# 执行记录 · GOAL-013

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| E-001 | 2026-08-16 | 目标建立与 as-built 对照（只读，无代码改动） | recorded | `02-execution/E-001-goal-created-and-as-built.md` |
| E-002 | 2026-08-16 | S3 P0：T-05 回收站时间 ISO + T-01 顶栏用户下拉 | recorded | `02-execution/E-002-s3-p0-t05-t01.md` |
| E-003 | 2026-08-16 | S3 P1：T-03 个人中心 Tabs + T-02 列表搜索矩阵 | recorded | `02-execution/E-003-s3-p1-t03-t02.md` |
| E-004 | 2026-08-16 | S3 P2：T-06 模块启用只认 config.yaml | recorded | `02-execution/E-004-s3-p2-t06.md` |
| E-005 | 2026-08-16 | S4 验证与关门（回归证据 + T-06 go 判定） | recorded | `02-execution/E-005-s4-verification-and-go.md` |
| E-006 | 2026-08-16 | A-002 响应：F-001 fixed（台账纠正）/ F-002 复跑复现 / F-003 fixed / F-004 fixed / F-005 accepted | recorded | `02-execution/E-006-a002-response.md` |

## 事实边界

> 只写已经发生且有证据的事实。计划、未知与建议留在决策。

- **2026-08-16**：记录决策 D-002（T-01 顶栏用户下拉冻结；I-005 verified）。无代码改动。
- **2026-08-16**：记录决策 D-003（T-02 搜索矩阵冻结；I-001 verified）。无代码改动。
- **2026-08-16**：记录决策 D-004（T-03 个人中心 Tabs；I-002 verified）。无代码改动。
- **2026-08-16**：记录决策 D-005（T-04 移交 [workspace-011] GOAL-022；I-003 verified）。无本区代码改动。
- **2026-08-16**：记录决策 D-006（T-05 API 改出 ISO-8601）。无代码改动。
- **2026-08-16**：记录决策 D-007（T-06 模块启用只认 YAML；I-004 verified；I-006 仍开放）。无代码改动。
- **2026-08-16**：记录决策 D-008（I-006 verified；S1/S2 勾选，progress 2/4）。无代码改动。
- **2026-08-16**：S3 P0（E-002）——T-05 回收站 deletedAt/restoredAt 改 UTC ISO-8601（handler + 测试断言）；T-01 顶栏用户下拉（UserMenu：头像+姓名触发器、projection.user 序菜单 + 退出末项、抽屉去用户链、shell.userMenu i18n）+ user-menu 测试 4 例。
- **2026-08-16**：S3 P1（E-003）——T-03 account.json 三档 Tabs（资料/安全/会话）+ TabsView labelKey；T-02 搜索矩阵：searchFormSubmit 泛化（非 q 字段入 filters）+ 12 页 schema 改造（含 notifications 新增搜索）+ 后端 ExtraQuery 接线（users/roles/tasks/runs/notifications）+ i18n ~50 键/语言 + search-form-filters 测试。
- **2026-08-16**：S3 P2（E-004）——T-06 模块启用只认 config.yaml：config.go 移除 APP_PROFILE/APP_MODULES_ENABLED 与 modules_enabled，新增 app.modules（preset 内置名/预设文件 + 内联 list，互斥）；kernel source/precedence 更新；compose 挂载 configs + 移除 env；dev.cmd 写 overlay CONFIG_FILE；README/QUICKSTART/playbook 文档同步；config_test 重写（7 子用例）。
- **2026-08-16**：S4（E-005）——回归全绿：Go 全量 0 FAIL、Web 1027/1027、tsc 0、D-VAL 绿；T-06 go 判定「部署契约变化、默认集不变 → 不暂挂」；A-001 self pass；A-002 grok independent **conditional**（required F-001 台账问题，响应见 E-006）。
- **2026-08-16**：A-002 响应（E-006）——F-001 **fixed**（撤回预写 done / 索引按真实 verdict 纠正）；F-002 **fixed**（全量复跑 1027/1027 复现绿）；F-003 **fixed**（usersWhere 增加 id 搜索 + COUNT 别名）；F-004 **fixed**（显示名全断点显示）；F-005 **accepted-residual**（D-003 降级条款）。S4 闭合，W12 关门 4/4。
