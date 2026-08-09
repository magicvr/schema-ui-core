---
id: GOAL-001-localization-and-system-settings
title: 多语种与系统设置产品化
status: done
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.9.0
progress: 7/7
plan_refs:
  - VP-007-localization-and-system-settings
primary_plan: VP-007-localization-and-system-settings
serves_summary: 在 I-PROTO-FULL-001 完整协议面与 VP-005 设计系统之上建立 zh-CN/en-US 多语种运行时与 General/Branding/Localization/Appearance 四类系统设置产品面；不扩张 Profile 可见性边界、不重定义上游协议语义。
---

# GOAL-001 · 多语种与系统设置产品化

## 概述

本 Root 承接 [VP-007 · 多语种与系统设置产品化](../../vision/plans/VP-007-localization-and-system-settings.md)（2026-08-09 用户确认激活；lead delivery 工作区 `workspace-007-localization-and-system-settings`）。意图权威 = VP-007 文件；方向级退出判据 1～6 与信息门禁 `I-L10N-001`～`005` 以此为准。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `plan_refs` / `primary_plan` | `VP-007-localization-and-system-settings` |
| Charter | `schema-ui-core-admin-foundation@0.2.0` |
| 工作区 | `workspace-007-localization-and-system-settings` · delivery |

## 成功边界（镜像 VP-007 方向级退出判据）

1. **语种解析与用户控制闭合**：`zh-CN`/`en-US`/`auto` 解析、匹配、回退与切换可验证；优先级一致应用于登录前后；HTML `lang`、日期和数字格式跟随有效 locale。
2. **前端可见文本形成可维护翻译面**：S0 冻结的最小可枚举证据面（固定 UI + 双 Profile Runtime Manifest page/schema 并集 + M1～M4 主流程）不存在依赖硬编码英文的路径；`titleKey`/`labelKey` 真解析 + 文本 fallback + 缺失 key 可观察。
3. **系统设置产品面可用**：General / Branding / Localization / Appearance 四类字段读写生效；预览、校验、恢复默认、权限失败和刷新行为可验证。
4. **Profile 与公开启动边界一致**：`mvp` 与 `admin` 同一前端 build、同一多语种运行时；`admin.settings` 编辑面仅 admin Profile 暴露。
5. **错误与提示本地化保持兼容**：稳定错误码 + 前端按码/key/参数本地化不可降级；`I-L10N-004` 以 (a) 服务端 locale 实施证据或 (b) 用户书面 `accepted-residual` 二选一关闭。
6. **质量与关门证据完整**：两语种 × 两 Profile × 匿名/认证证据矩阵复用同一分母；lead 工作区 Root 完成约定范围；开放 required = 0；用户确认关门。

## 纲领路线图（P-001）

| 阶段 | 目的 | 状态 |
|------|------|------|
| S0 | 差距与契约冻结：关闭 `I-L10N-001`～`005`；冻结 locale、翻译资源、公开 bootstrap、错误兼容和时区语义 | **done**（2026-08-09，D-002 + F-V029 冻结） |
| S1 | 多语种核心：locale resolver/provider、资源装载、缺失 key/fallback、用户切换、HTML lang 与格式化 | **done**（2026-08-09，GOAL-002-s1-locale-core 6/6，A-001 pass） |
| S2 | 前端与 Schema 覆盖：固定 UI 与双 Profile Runtime Manifest page/schema 分母双语化；`titleKey`/`labelKey` 真解析 | **done**（2026-08-09，GOAL-003-s2-ui-schema-bilingual 5/5，A-001 pass） |
| S3 | 系统设置产品化：四类设置、公开启动配置、品牌/locale/theme 生效、权限/审计/刷新闭环 | **done**（2026-08-09，GOAL-004-s3-settings-productization 6/6，A-001 pass） |
| S4 | 后端用户可见反馈：稳定错误码 + 前端本地化保底；按 `I-L10N-004` 结论实施有界服务端 locale 协商 | **done**（2026-08-09，GOAL-005-s4-error-localization 5/5，A-001 pass；I-L10N-004 实施证据齐备） |
| S5 | 双 Profile 验证与关门：两语种 × 两 Profile × 登录前后/权限/失败矩阵、文档、审计与 close-out | **done**（2026-08-09，GOAL-006 done 4/4；A-001→A-002；用户确认 D-002） |
| S6 | 设置页表单/详情页改造：实现 `form.props.recordSource` 预填（ADR-0021）+ settings 页四类内联表单重构 | **done**（2026-08-09，GOAL-007 done 4/4；A-001→A-002；用户确认 D-002；Root 恢复 done） |

## 信息需求（P-005）

完整信息台账（编号、问题、影响门禁、最晚阶段、验证动作、状态、延期/复核）见 [01-decision.md](01-decision.md)；`I-L10N-001`～`005` 全部 `required` / `open`，与 VP-007 同 id 对齐：

| ID | 问题（摘要） | 状态 | 最晚需要阶段 |
|----|--------------|------|--------------|
| `I-L10N-001` | Schema 驱动页面文本如何本地化且不创建宽于/冲突于 `schema-ui-docs@v2.7.0` 的私有语义 | **verified**（前端 key 解析，D-002） | 多语种方案冻结前（S0）✓ |
| `I-L10N-002` | 用户显式语种选择持久化边界与匿名→登录合并规则 | **verified**（localStorage 单通道，D-002） | 用户控制实施前（S1）✓ |
| `I-L10N-003` | 公开品牌/locale 启动配置：兼容扩展 `/api/branding` 还是新 bootstrap 契约 | **verified**（兼容扩展，D-002） | Settings API 方案冻结前（S3）✓ |
| `I-L10N-004` | 错误 envelope 扩展到 locale 协商的真实成本与兼容边界（exit 5 二选一关闭） | **verified**（路径 a 有界服务端协商，D-002） | 后端提示本地化实施前（S4）✓ |
| `I-L10N-005` | 默认时区存储/展示/服务器时间语义（UTC 存储、显示转换、失败语义） | **verified**（UTC + 显示转换，D-002） | Localization 设置实施前（S3）✓ |

## 派生进度

`progress: 7/7`：由上方 7 个纲领阶段检查点等权派生（P-001），仅展示；不放行阶段、不关闭 finding。S0–S5 已由 GOAL-006 D-002 用户书面确认关门（历史记录保留不重写）；2026-08-09 用户指令暂时回退承接 S6（GOAL-001 `D-003`，`done`→`active`、`6/6`→`6/7`）；**S6 C4 用户书面确认（GOAL-007 D-002）后 Root 恢复 `done`（`7/7`），解除临时回退（本文件 `D-004`）**。

## 台账布局

`01-decision/` · `02-execution/` · `03-audit/`（自首条记录起使用平铺 ledger）
