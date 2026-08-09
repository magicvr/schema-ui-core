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
GOAL-001-localization-and-system-settings [active] (5/6)
├── GOAL-002-s1-locale-core [done] (6/6)
├── GOAL-003-s2-ui-schema-bilingual [done] (5/5)
├── GOAL-004-s3-settings-productization [done] (6/6)
└── GOAL-005-s4-error-localization [done] (5/5)
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| GOAL-001-localization-and-system-settings | 多语种与系统设置产品化 | `null` | **active** | `5/6` | 2026-08-09 |
| GOAL-002-s1-locale-core | S1 · 多语种核心（locale 解析/资源/切换/格式化） | GOAL-001-localization-and-system-settings | **done** | `6/6` | 2026-08-09 |
| GOAL-003-s2-ui-schema-bilingual | S2 · 固定 UI 与 Schema 分母双语化（titleKey/labelKey 真解析） | GOAL-001-localization-and-system-settings | **done** | `5/5` | 2026-08-09 |
| GOAL-004-s3-settings-productization | S3 · 系统设置产品化（四类设置 + 公开启动配置） | GOAL-001-localization-and-system-settings | **done** | `6/6` | 2026-08-09 |
| GOAL-005-s4-error-localization | S4 · 后端反馈本地化（稳定错误码 + 有界服务端协商） | GOAL-001-localization-and-system-settings | **done** | `5/5` | 2026-08-09 |

## 维护说明

- Root `GOAL-001-localization-and-system-settings` 于 2026-08-09 由 `/govern` scaffold：`docs/workspace-007-localization-and-system-settings/`（delivery，`plan_refs`/`primary_plan` = VP-007）。
- 纲领路线图 S0–S5 见 Root `00-meta.md`；`progress: 5/6` 由 6 个阶段检查点等权派生（P-001），仅展示，不放行阶段、不关闭 finding。
- **S0 已完成**（2026-08-09，D-002）：`I-L10N-001`～`005` 全部 `verified`（用户书面裁决）；F-V029 冻结；独立审计 A-001 → A-002 响应闭合。
- **S1 已完成**（2026-08-09）：GOAL-002 done 6/6（locale 运行时 + 切换器）。
- **S2 已完成**（2026-08-09）：GOAL-003 done 5/5（固定 UI + 12 page/schema 并集双语化 + titleKey/labelKey 真解析 + M4）。
- **S3 已完成**（2026-08-09）：GOAL-004 done 6/6（四类设置 + `/api/branding` 扩展 + 权限/审计/刷新闭环）。
- **S4 已完成**（2026-08-09）：GOAL-005 done 5/5（错误码契约 + 有界服务端协商 + 前端保底；I-L10N-004 实施证据齐备）。
- **S5** 下一阶段：双 Profile 验证矩阵、文档、独立关门审计与用户确认。
- 开放 required：Goal 审计 0；Vision Review 0。
