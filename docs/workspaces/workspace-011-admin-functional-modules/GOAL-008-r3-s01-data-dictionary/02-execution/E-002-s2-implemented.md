---
id: E-002
goal: GOAL-008-r3-s01-data-dictionary
date: 2026-08-14
status: recorded
parent: GOAL-008-r3-s01-data-dictionary
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-002 · S2 实现完成

## 事实

- 2026-08-14：S2 实现完成，覆盖 D-002 §2–§5：
  - **migration 0019**（admin.data-dictionary）：dict_types + dict_entries（UNIQUE(dict_key, entry_key)、类型删除级联条目）；checksum 8f2c2a18…（Go 权威计算）。
  - **migration 0020**（core.operationlog）：CHECK + dictionary.create/update/delete；checksum ac6d2c28…。
  - **store** apps/api/internal/modules/datadictionary/store/：子包仓库（TxRunner 边界、List/Get/Create/Update/Delete ×2、哨兵错误 ErrNotFound/ErrTypeKeyTaken/ErrEntryKeyTaken/ErrDictKeyNotFound）。
  - **handler** apps/api/internal/handler/dictionary.go：双资源（types/entries）经资源工厂；PATCH 缺失字段保持原值；JSONFields 承载 enabled/sort/remark（工厂字符串字段模型）；审计事件 + slog 失败留痕。
  - **模块** provider.go（五面贡献）+ schema ×2（类型页 + 条目页，navigate action 到条目页）+ fragment（两页 + menu_dictionary）。
  - **装配**：profile（admin 默认集 + BuiltinModules）、composition、testsupport、compiled/persistence.go、handler 测试环境。
  - **web**：i18n zh/en（manifest/schema/dictionaryEntries + error.* 5 键）；schema-keys 结构测试跳过 wire-mapping 子树（bodyMapping.label 非用户可见文本）；fixture（admin + sha 重钉 68cdd3ed…）；smoke admin 页面集 + data-dictionary。
  - **测试**：handler dictionary_test（生命周期/门禁/级联/DICT_KEY_NOT_FOUND/审计）、provider_test（注册面 + 端到端）；错误码目录 + 钉住集（5 个字典码入 domain 集）。
