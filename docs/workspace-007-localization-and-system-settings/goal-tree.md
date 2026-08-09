---
workspace_id: workspace-007-localization-and-system-settings
root_goal: GOAL-001-localization-and-system-settings
canonical_scope: docs/workspace-007-localization-and-system-settings/
status: active
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
parent: null
---

# 目标树 · 多语种与系统设置产品化工作区

| 字段 | 值 |
|------|----|
| 工作区 | `workspace-007-localization-and-system-settings` |
| Root Goal | `GOAL-001-localization-and-system-settings` |
| primary plan | `VP-007-localization-and-system-settings` |

## ASCII 树

```text
GOAL-001-localization-and-system-settings [active] (3/6)
├── GOAL-002-s1-locale-core [done] (6/6)
└── GOAL-003-s2-ui-schema-bilingual [done] (5/5)
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| GOAL-001-localization-and-system-settings | 多语种与系统设置产品化 | `null` | **active** | `3/6` | 2026-08-09 |
| GOAL-002-s1-locale-core | S1 · 多语种核心（locale 解析/资源/切换/格式化） | GOAL-001-localization-and-system-settings | **done** | `6/6` | 2026-08-09 |
| GOAL-003-s2-ui-schema-bilingual | S2 · 固定 UI 与 Schema 分母双语化（titleKey/labelKey 真解析） | GOAL-001-localization-and-system-settings | **done** | `5/5` | 2026-08-09 |

## 维护说明

- Root `GOAL-001-localization-and-system-settings` 于 2026-08-09 由 `/govern` scaffold：`docs/workspace-007-localization-and-system-settings/`（delivery，`plan_refs`/`primary_plan` = VP-007）。
- 纲领路线图 S0–S5 见 Root `00-meta.md`；`progress: 3/6` 由 6 个阶段检查点等权派生（P-001），仅展示，不放行阶段、不关闭 finding。
- **S0 已完成**（2026-08-09，D-002）：`I-L10N-001`～`005` 全部 `verified`（用户书面裁决）；F-V029 覆盖表冻结于 Root `attachments/`；独立审计 A-001（conditional）已由 A-002 正式响应，required findings 全部合法闭合（用户裁决 fixed）——审计响应后确认 S0 done。
- **S1 已完成**（2026-08-09）：`GOAL-002-s1-locale-core` done 6/6（locale 运行时 + 切换器；A-001 pass）。
- **S2 已完成**（2026-08-09）：`GOAL-003-s2-ui-schema-bilingual` done 5/5（固定 UI + 12 page/schema 并集双语化 + titleKey/labelKey 真解析 + M4；690/690 测试全绿；A-001 pass）；协议 pin 边界核对：component-registry 未改写，缺口字段以本地文档约定登记。
- **S3** 下一阶段：General/Branding/Localization/Appearance 四类设置 + `/api/branding` 扩展 + 权限/审计/刷新闭环。
- 开放 required：Goal 审计 0；Vision Review 0。
