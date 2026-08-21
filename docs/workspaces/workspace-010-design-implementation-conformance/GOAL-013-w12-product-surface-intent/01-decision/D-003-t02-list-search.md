---
id: D-003-t02-list-search
doc: decision-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-003 · T-02 列表搜索字段矩阵（S2 分项冻结）

## 背景

W11 U-04 给多页补了通用 `q` + `feedback.search`，语义不明。渲染器控件白名单已齐（`form-controls.ts`）；缺的是各业务 schema 与后端筛选接线。用户本轮采纳「关键词改名 + 补业务筛选」。

## 决定

1. **原则**：禁止无语义的「搜索」占位。每页最多一个关键词输入（`id` 仍可用 `q` 以兼容现有 handler），**label 必须写清搜什么**；离散状态用 `select`，走 `resourceFilter.Extra` / 已有 query。
2. **本波页集**：现有 11 个 `q` 搜索页 + **通知列表**。**不做**：数据权限、系统监控。不加日期范围、不扩协议。
3. **深度**：schema 改标签 + 接线已有 Extra；并对「列已在库、handler 尚未暴露」的筛选补 ExtraQuery / store 条件。
4. **闭合 I-001**。

### 冻结矩阵

| 页 | 关键词 label（意图） | 筛选 | 后端 |
|----|----------------------|------|------|
| users | 用户名 / 显示名 / ID | `enabled`、`locked`（全部/是/否） | 现仅 `Q`；本波补 Extra + `UserFilter` |
| roles | 角色名 / Key | `system`（全部/系统/自定义） | 现仅 `Q`；本波补 Extra |
| activity | 事件 / 操作者 / 资源 | —（关键词覆盖 event/actor/resource；不加日期） | 现 `Q`，保持 |
| wallet | 账户 / 所有者 | `ownerType` | API **已有** `ownerType`，schema 接线 |
| wallet-entries | 备注 / 关联单 | `entryType`（若 list 已按类型可滤；否则仅关键词） | 核对 list 后再接线，不发明类型 |
| file-library | 文件名 | —（现 `q` 已扫 name/type/owner） | 不新增 Extra |
| data-dictionary | 字典 Key / 名称 | — | 现 `Q` 已扫 key/name |
| dictionary-entries | 条目 Key / 名称 | `dictKey` 保持内页上下文 | Extra **已有** `dictKey` |
| recycle-bin | 资源 ID / 操作者 | `resource`（资源类型） | API **已有** `resource`，schema 接线 |
| scheduled-tasks | 任务 Key / 名称 | `enabled` | 现 `Q` 扫 key/name；本波补 Extra |
| task-runs | 任务 Key | `status`（若 run 行有 status 列） | 现 `Q`；有列则补 Extra |
| notifications | 标题 / 正文 | `read`（全部/未读/已读） | 本波补 query；对齐 W11 U-12 未做部分 |

实施时若某 Extra 在 store 无列或无稳定枚举，该筛选降级为「仅关键词」并在 E 条目留痕，不另开目标。

## 理由

- 用户要的是业务语义，不是再做一个万能框。
- ExtraQuery 管道已存在（GOAL-015 `dictKey`），补筛选不必改协议。
- 监控是实时面板、数据权限是策略矩阵，本波硬加搜索收益低。

## 未选方案

- **只改标签、不补后端筛选**：回收站/钱包已有 query 却接不上，用户仍觉得「只有一个框」。
- **再加数据权限页**：范围膨胀；该页以矩阵操作为主。
- **日期范围**：渲染器有 `dateRangePicker`，但多数 list 无时间谓词；留给后续波次。

## 影响

- schema + i18n + 若干 handler/store 筛选；审计模式保持 `self`。
- 不改 Profile / 模块矩阵。

## 后续

- T-02 进 S3 P1（与 T-03）。下一项裁决：T-03 个人中心 Tabs。
