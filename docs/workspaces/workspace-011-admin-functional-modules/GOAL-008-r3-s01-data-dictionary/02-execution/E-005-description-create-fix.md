---
id: E-005
goal: GOAL-008-r3-s01-data-dictionary
date: 2026-08-14
status: recorded
parent: GOAL-008-r3-s01-data-dictionary
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-005 · 关门后修复：新建类型时 description 始终为空

## 事实

- 2026-08-14（S5 关门后）：用户反馈数据字典新建类型时「描述」输入项不被识别（无论输入什么内容都为空）。

## 根因

- dict types 的 Resource 注册 `CreateFields: ["key", "name"]`、`JSONFields: ["enabled", "sort"]`；工厂 `decodeResourceCreate` 只透传 CreateFields/JSONFields/RawStringFields。
- `description` 不在任何通道 → create 请求体中的 description 被静默丢弃；实体 `dictTypeEntity.Create` 虽读取 `body["description"]` 但永远拿不到值。

## 修复

- `description` 移入 `JSONFields`（可选语义完全匹配：create 缺失 = 实体默认 ""；patch 缺失 = 不动、提供 = 更新、空串 = 清空）；`PatchFields` 收敛为 `["name"]`（name 是唯一必填非空 patch 字段）。
- 新增测试：create 带 description 断言回读；patch 更新 description；patch 清空 description 允许（JSONFields 可选语义）。
- 真实 HTTP 冒烟：create description="描述内容 hello" 回读一致；patch "updated" 生效。

## 验证

- handler + datadictionary 全量测试绿。

## 关联

- 用户同时提出表单验证与表单弹窗样式问题（INVALID_PATCH_FIELD 缺字段级提示；两列布局业界对比）——待评估是否新增波次子目标（未在本条处理）。
