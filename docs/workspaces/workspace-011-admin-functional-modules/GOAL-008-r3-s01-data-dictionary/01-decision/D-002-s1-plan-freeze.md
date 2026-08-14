---
id: D-002
goal: GOAL-008-r3-s01-data-dictionary
title: 方案冻结：数据字典设计（S1）
date: 2026-08-14
status: accepted
parent: GOAL-008-r3-s01-data-dictionary
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · 方案冻结（S-01 数据字典）

## 1. 协议对照（I-001 闭合）

- 字典管理无新协议语义契约（枚举/字典为呈现自由）；页面用既有 schema 驱动机制（table/form/request/modal/navigate action）。
- 类型页行操作「条目」= 协议 `NavigateAction`（action.schema.json，type "navigate" + url，ADR-0023 既有导航路径）；无新 renderer 扩展。
- 条目 dict_key 输入 v1 用文本字段（服务端校验存在性；不做动态 options 数据源——无协议 select 动态数据契约，呈现自由留痕）。

## 2. 数据模型（migration 0019，admin.data-dictionary）

`dict_types`：id TEXT PK ｜ key TEXT UNIQUE NOT NULL ｜ name TEXT NOT NULL ｜ enabled INTEGER 1 ｜ description TEXT ｜ sort INTEGER 0 ｜ created_at / updated_at INTEGER NOT NULL

`dict_entries`：id TEXT PK ｜ dict_key TEXT NOT NULL REFERENCES dict_types(key) ON DELETE CASCADE ｜ entry_key TEXT NOT NULL ｜ label TEXT NOT NULL ｜ enabled INTEGER 1 ｜ sort INTEGER 0 ｜ remark TEXT ｜ created_at / updated_at INTEGER NOT NULL ｜ UNIQUE(dict_key, entry_key)

## 3. 资源与端点（资源工厂 × 2）

| 资源 | Path | 排序字段 | 说明 |
|------|------|----------|------|
| dict-types | /api/data-dictionary/types | key/name/sort/updatedAt | 类型 CRUD；删除级联条目 |
| dict-entries | /api/data-dictionary/entries | dictKey/entryKey/label/sort/updatedAt | 条目 CRUD；create/patch 校验 dict_key 存在（DomainError DICT_KEY_NOT_FOUND → 400） |

- 读门禁 dictionary.read；写门禁 dictionary.write（均 PolicyAdmin，admin-only 管理面）。
- q 搜索（类型 key/name；条目 dictKey/entryKey/label）；分页/排序走工厂既有契约。
- 条目页 v1 为全局列表 + q 过滤（不做 dictKey 专用筛选参数——工厂查询契约不扩展，文档化）。

## 4. 页面与 i18n

- 页面 `data-dictionary`（类型列表：key/name/enabled/sort/updatedAt；工具栏新建；行操作 编辑/删除/条目→navigate /dictionary-entries）。
- 页面 `dictionary-entries`（条目列表：dictKey/entryKey/label/enabled/sort/updatedAt；工具栏新建；行操作 编辑/删除）。
- manifest fragment：两页 + `menu_dictionary`（visibility PolicyAdmin，Permission dictionary.read，指向 data-dictionary）。
- i18n zh/en：manifest.title.dataDictionary / manifest.title.dictionaryEntries / manifest.nav.dataDictionary + schema.dataDictionary.* / schema.dictionaryEntries.* + error.dictKeyNotFound。

## 5. 审计（migration 0020，core.operationlog）

- 事件 `dictionary.create / dictionary.update / dictionary.delete`（类型与条目共用；record_id = 行 id）；CHECK rebuild 同 0018。

## 6. 测试与验证

- 模块 provider 测试（注册面）；handler 测试（类型/条目 CRUD、门禁 401/403、级联删除、DICT_KEY_NOT_FOUND、审计事件、q/排序/分页）。
- 组合根：admin 权限 15→17、导航 8→9；mvp 不变；迁移 18→20（fresh/reopen/ownership/ending）。
- web：fixture（admin + sha 重钉）、schema-keys 分母、s5-denominator、e2e admin shell 导航断言 + manifest schema 护栏自动覆盖。
- 冒烟：SM-007 admin 页面集 + data-dictionary。

## 7. 未选方案

- 类型/条目单页 master-detail（recordView 内嵌表）：renderer 无该组合契约，v1 否决。
- 条目页 dictKey 专用筛选参数：工厂查询契约不扩展（q 可替代），文档化。
- 动态 select options 数据源：无协议契约，v1 否决（呈现自由留痕）。
- 类型删除 RESTRICT + DICT_TYPE_IN_USE：v1 级联更可预测（文档化）；如业务需要再引入。
