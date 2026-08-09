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
GOAL-001-localization-and-system-settings [active] (2/6)
└── GOAL-002-s1-locale-core [done] (6/6)
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| GOAL-001-localization-and-system-settings | 多语种与系统设置产品化 | `null` | **active** | `2/6` | 2026-08-09 |
| GOAL-002-s1-locale-core | S1 · 多语种核心（locale 解析/资源/切换/格式化） | GOAL-001-localization-and-system-settings | **done** | `6/6` | 2026-08-09 |

## 维护说明

- Root `GOAL-001-localization-and-system-settings` 于 2026-08-09 由 `/govern` scaffold：`docs/workspace-007-localization-and-system-settings/`（delivery，`plan_refs`/`primary_plan` = VP-007）。
- 纲领路线图 S0–S5 见 Root `00-meta.md`；`progress: 2/6` 由 6 个阶段检查点等权派生（P-001），仅展示，不放行阶段、不关闭 finding。
- **S0 已完成**（2026-08-09，D-002）：`I-L10N-001`～`005` 全部 `verified`（用户书面裁决）；F-V029 覆盖表冻结于 Root `attachments/`；独立审计 A-001（conditional）已由 A-002 正式响应，required findings F-001/F-002 全部合法闭合（用户裁决 fixed）——**审计响应后确认 S0 done**。
- **S1 已完成**（2026-08-09）：`GOAL-002-s1-locale-core` done 6/6（locale 运行时 + 切换器 + 45 项新测试，674/674 全绿，build 通过）；自审 A-001 pass。
- **S2** 下一阶段：固定 UI + 双 Profile page/schema 并集双语化、`titleKey`/`labelKey` 真解析。
- 开放 required：Goal 审计 0；Vision Review 0。
