---
id: D-001
goal: GOAL-008-r3-s01-data-dictionary
title: 立项边界：模块身份、Profile 归属与审计策略
date: 2026-08-14
status: accepted
parent: GOAL-008-r3-s01-data-dictionary
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 立项边界（S-01 数据字典）

## 决定

1. **模块身份**：`admin.data-dictionary`（标准 Admin 功能模块）；Descriptor 依赖 core.auth-session / core.schema-render / core.navigation-capability / core.operationlog。
2. **Profile 归属（I-002 闭合）**：进入 **admin 默认集**（Profile 内容扩展，S-02/D-001 §2 与 F-01/F-02 先例）；mvp / demo 不启用。
3. **审计策略**：数据字典 = 常规 CRUD + 枚举数据，低风险 → **self 审计为主**；关门门禁按 P-004 与用户确认是否引入独立审计（预期 self 足够）。
4. **迁移归属**：0019 = 字典表（dict_types/dict_entries，归属 admin.data-dictionary）；0020 = operation_log CHECK 扩展 dictionary.* 事件（归属 core.operationlog，rebuild 模式同 0014/0015/0018）。
5. **删除语义**：类型删除级联删除其条目（v1，文档化；条目为从属数据）。
