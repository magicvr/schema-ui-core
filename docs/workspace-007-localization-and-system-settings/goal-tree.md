---
workspace_id: workspace-007-localization-and-system-settings
root_goal: GOAL-001-localization-and-system-settings
canonical_scope: docs/workspace-007-localization-and-system-settings/
status: active
created: 2026-08-09
updated: 2026-08-09
version: 0.5.0
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
GOAL-001-localization-and-system-settings [active] (6/7)
├── GOAL-002-s1-locale-core [done] (6/6)
├── GOAL-003-s2-ui-schema-bilingual [done] (5/5)
├── GOAL-004-s3-settings-productization [done] (6/6)
├── GOAL-005-s4-error-localization [done] (5/5)
├── GOAL-006-s5-evidence-and-closeout [done] (4/4)
└── GOAL-007-s6-settings-form-page [active] (3/4)
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| GOAL-001-localization-and-system-settings | 多语种与系统设置产品化 | `null` | **active** | `6/7` | 2026-08-09 |
| GOAL-002-s1-locale-core | S1 · 多语种核心（locale 解析/资源/切换/格式化） | GOAL-001-localization-and-system-settings | **done** | `6/6` | 2026-08-09 |
| GOAL-003-s2-ui-schema-bilingual | S2 · 固定 UI 与 Schema 分母双语化（titleKey/labelKey 真解析） | GOAL-001-localization-and-system-settings | **done** | `5/5` | 2026-08-09 |
| GOAL-004-s3-settings-productization | S3 · 系统设置产品化（四类设置 + 公开启动配置） | GOAL-001-localization-and-system-settings | **done** | `6/6` | 2026-08-09 |
| GOAL-005-s4-error-localization | S4 · 后端反馈本地化（稳定错误码 + 有界服务端协商） | GOAL-001-localization-and-system-settings | **done** | `5/5` | 2026-08-09 |
| GOAL-006-s5-evidence-and-closeout | S5 · 双 Profile 验证矩阵与关门 | GOAL-001-localization-and-system-settings | **done** | `4/4` | 2026-08-09 |
| GOAL-007-s6-settings-form-page | S6 · 设置页表单/详情页改造（recordSource 预填） | GOAL-001-localization-and-system-settings | **active** | `3/4` | 2026-08-09 |

## 维护说明

- Root `GOAL-001-localization-and-system-settings` 于 2026-08-09 由 `/govern` scaffold；2026-08-09 用户书面确认关门（GOAL-006 D-002），S0–S5 全部完成。
- **2026-08-09 用户指令暂时回退 Root 关门状态**（`done`→`active`，`6/6`→`6/7`）承接 S6 子目标 `GOAL-007-s6-settings-form-page`（设置页表单/详情页改造）；重新开根决策见 GOAL-001 `D-003`。历史关门记录不重写。
- `progress` 由纲领阶段检查点等权派生（P-001），仅展示。
- **S5 关门证据**：F-V029 矩阵 + dual-run branding（admin/mvp）+ web build + playwright admin/mvp + independent A-001 → A-002（required 0）+ 用户确认 D-002。
- 开放 required：Goal 审计 0；Vision Review 0。VP-007 保持 `closed`（S6 为其上增量产品化，不触碰 vision 状态）。
