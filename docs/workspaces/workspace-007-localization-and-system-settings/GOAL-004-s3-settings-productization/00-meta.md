---
id: GOAL-004-s3-settings-productization
title: S3 · 系统设置产品化（四类设置 + 公开启动配置）
status: done
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
progress: 6/6
---

# GOAL-004 · S3 · 系统设置产品化

## 概述

承接 Root [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md) 的 **S3 阶段**：把 `admin.settings` 从仅站点标题/Logo 扩展为 **General / Branding / Localization / Appearance** 四类系统设置产品面——字段读写、预览、校验、恢复默认、权限失败与刷新行为可验证；`GET /api/branding` 按 I-L10N-003 兼容扩展为公开启动配置（defaultLocale/supportedLocales/defaultTheme/siteTimezone + 品牌字段）；站点标题/品牌资产/默认语种/时区/默认主题在 Shell、登录页与公开启动配置一致生效；权限沿用 `settings.read`/`settings.write`，变更留下操作日志与配置刷新信号。

**方案依据**：Root D-002（I-L10N-003 兼容扩展 `/api/branding`、I-L10N-005 UTC 存储+显示转换、F-001 用户裁决补全 Branding 字段 `logoUrlLight`/`logoUrlDark`/`faviconUrl`、Settings 四类字段冻结）。本目标只实施与验证，不重新决策。

**范围纪律**：不实现上传/对象存储/媒体管理（Non-goals）；不把 `admin.settings` 加入 `mvp` Profile；双 Profile 共用同一前端 build；不改写协议语义；`mvp` 不暴露 settings 编辑面。

## 成功标准（可验收 · 等权检查点 · 共 6 项）

- [x] **C1**：四类字段读写生效：General（siteTitle）、Branding（logoUrl/logoUrlLight/logoUrlDark/faviconUrl）、Localization（defaultLocale/siteTimezone）、Appearance（defaultTheme）——PATCH 字段级合并、空值清空语义、恢复默认、校验错误（含无效 IANA 时区 → 400 INVALID_TIMEZONE 且不清空原值）可验证（Go 测试驱动真实 handler/仓库）。
- [x] **C2**：公开启动配置：`GET /api/branding` additive 扩展返回全部启动字段（旧字段兼容），缓存语义不变；Shell/登录页消费扩展字段；配置刷新事件（X-Schema-UI-Config-Changed）触发重新拉取并生效。
- [x] **C3**：投影一致：站点标题、favicon/浅深色 Logo、默认语种、默认主题、时区在 Shell、登录页与公开启动配置一致生效（品牌应用、主题解析含系统默认、locale provider 注入系统默认）。
- [x] **C4**：Settings 四类设置面：schema 驱动页面按四类组织（字段、说明、保存、预览、恢复默认按钮）；`admin.settings` 编辑面仅 admin Profile 可达；`mvp` Profile 不可达（N/A 模块边界）。
- [x] **C5**：权限/审计/刷新闭环：`settings.read`/`settings.write` 权限失败返回稳定错误码且前端可呈现；设置变更写操作日志；刷新信号触发配置重新加载。
- [x] **C6**：验证：Go 测试全绿（新增 settings 扩展测试）；vitest 全绿（新增四类设置/启动配置/投影测试）；`npm run build` 通过；证据捕获 `{SCRATCH}`；F-V029 M3 证据回填。

## 派生进度展示

`progress: 6/6` 由上方 6 个等权检查点派生；仅为展示，不放行阶段、不关闭 finding。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 四类字段语义（合并/清空/恢复默认/校验/时区失败）是否齐备 | C1–C6 | 实施前 | 读 Root D-002（I-L10N-003/005、F-001 fixed）+ 现有 settings 迁移/仓库/配置 | **closed** | — | Root D-002 冻结（2026-08-09 用户裁决） |
| I-002 | required | Settings 页 schema 形态与预览/恢复默认的客户端行为承载（协议允许面） | C4 | 实施前 | 盘点 Renderer 对 custom action/状态的现状，冻结方案 | **closed** | — | 见 D-001：预览=表单内即时投影（Branding 图片缩略 + Appearance 主题即时应用，client 白名单 handler）；恢复默认=request action |

## 父目标

- [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md)（Root；本目标为 S3 阶段子目标）

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。
