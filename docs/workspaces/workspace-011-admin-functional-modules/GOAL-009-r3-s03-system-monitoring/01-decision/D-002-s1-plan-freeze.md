---
id: D-002
goal: GOAL-009-r3-s03-system-monitoring
title: 方案冻结：系统监控与错误日志（S1）
date: 2026-08-14
status: accepted
parent: GOAL-009-r3-s03-system-monitoring
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · 方案冻结（S-03 系统监控与错误日志）

## 1. 边界（I-001 闭合）

- **监控面**：只读汇总端点，进程内探测（store ping + readiness 门）+ 版本/启动时长/模块清单/数据库大小——不依赖 /healthz /readyz 的 HTTP 往返（进程内等价探测）。
- **错误日志面（v1 边界）**：基架无独立错误日志持久化（slog 走容器 stderr）。v1 **复用 operationlog 事件面**作「最近事件/错误视图」（best-effort，q 可过滤，倒序分页）；**不新增内核错误采集**（文档化残余，B-11 登录日志独立视图另议）。
- **与 activity 切分**：activity（operations.read，PolicyAdminEditor）= 完整审计追溯；监控页（monitoring.read，PolicyAdmin）= 监控语境下的最近事件快照 + 系统状态。两页数据同源（operationlog），语义不同，文档化。

## 2. 端点与权限

| 端点 | 门禁 | 说明 |
|------|------|------|
| GET /api/system-monitoring/status | monitoring.read | {status, ready, version, commit, uptimeSeconds, modules[], dbSizeBytes}（进程内探测） |
| GET /api/system-monitoring/errors | monitoring.read | 只读资源（工厂 ReadOnly）：operationlog 最近事件（q/sort/order/page/pageSize；默认 created desc） |

- 权限键：`monitoring.read`（PolicyAdmin）；导航 `menu_monitoring`（visibility PolicyAdmin，Permission monitoring.read，PageID system-monitoring）。
- 读端点全部 401/403 fail-closed（与既有只读面一致）。

## 3. 页面与 Schema

- 页面 `system-monitoring`：statCard 网格（状态/ready/启动时长/模块数/DB 大小，dataSource = status 端点 + valueField）+ 最近事件表（dataSource = errors 端点）。
- 既有节点：text/grid/statCard/table（dashboard 先例）；无新 renderer 扩展。
- i18n zh/en：manifest.title.systemMonitoring / manifest.nav.systemMonitoring + schema.systemMonitoring.* 键。

## 4. 协议对照（呈现自由）

- 监控页无新协议语义契约；statCard/table/grid 均为既有节点类型；呈现自由留痕。

## 5. 测试与验证

- handler 测试：status 内容（200/401/403、字段齐全、ready 门翻转）、errors 列表契约（分页/排序/q）。
- 组合根：admin 权限 17→18、导航 9→10；mvp 不变；无迁移变化。
- web：fixture（admin + sha 重钉）、schema-keys 分母、s5-denominator、e2e admin 导航断言。
- 冒烟：SM-007 admin 页面集 + system-monitoring。

## 6. 未选方案

- 内核错误日志采集（slog hook → 新表）：侵入内核 + 数据生命周期问题，v1 否决（B-11 语境另议）。
- 调用 /healthz /readyz HTTP 端点拼装：进程内探测等价且少一跳，选进程内。
