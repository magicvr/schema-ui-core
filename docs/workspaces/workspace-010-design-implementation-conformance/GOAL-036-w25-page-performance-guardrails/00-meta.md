---
id: GOAL-036-w25-page-performance-guardrails
title: W25 · 页面性能问题全盘修复与防复发（钱包页 + 全局机制 + 防复发栅栏）
status: active
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-001-design-implementation-conformance
version: 0.2.0
progress: 5/6
---

# GOAL-036 · W25 · 页面性能问题全盘修复与防复发

## 概述

用户报告「我的钱包」页面数据显示很慢（2026-08-23）。初版钱包页修复完成（请求 10→3、SQLite 连接面优化）后，用户升级范围（2026-08-23 书面裁决，D-002）：**全盘修复此类问题，并确保以后不会出现此类问题**。

本波覆盖四类页面性能反模式的**全局消除 + 防复发栅栏**：

1. 同 URL 展示节点重复请求（statCard/chart）→ 渲染层 in-flight 合并（钱包页 3×、系统监控页 **6×**、data-display 示例 3× 全部合 1）；
2. 自定义组件「挂载即写 + 整页重拉」→ 探活后写契约（wallet-ensure）；
3. 页面 schema 每次导航重取 + 全量校验 → App 级文档缓存；
4. 后端 SQLite 单连接全局串行 + 逐提交 fsync → 文件库池 4 + WAL/busy_timeout/synchronous=NORMAL。

防复发机制（S5）：store 连接面白盒回归测试、渲染层合并与定向刷新回归测试、schema 组件注册校验测试、`module-contribution-playbook.md` §6 页面数据面性能规范。

**范围裁决（D-002）**：纳入 monitoring 定向刷新（`refreshList`，只刷 `/status`）+ schema 注册校验；**大表 COUNT(*) 优化出局**（容量课题，索引健康，与本类问题正交，理由见 D-002）。

## 成功标准

1. **C1 诊断**：放大因素定位有文件级证据（D-001/E-001）；
2. **C2 钱包页实施**：后端池/WAL + 前端合并/探活/schema 缓存落地；
3. **C3 全盘扫描台账**：26 页 schema 扫描（D-001/E-001），全部同类问题由全局机制覆盖；
4. **C4 防复发机制**：store 回归测试 + 渲染层回归测试 + 注册校验测试 + Playbook 规范章节（E-002）；
5. **C5 回归全绿**：go test / vitest 全量 / tsc/build（E-001 + E-002）；
6. **C6 验证与关门**：Playwright e2e（**I-001 已完成关闭**）+ 活栈计时复核（I-002，open）+ 自审 + 台账同步（**尚未完成，本波不闭门**）。

## 路线图（P-001 · 分母 = 6）

```text
S1 诊断     → 四因素定位（E-001）✓
S2 方案     → D-001 取舍 ✓
S3 钱包页实施 → 后端 A + 前端 B/C/D（E-001）✓
S4 全盘扫描  → 26 页台账 + 全局机制覆盖确认（E-001）✓
S5 防复发   → 测试栅栏 + 规范章节 + monitoring 定向刷新（E-002）✓
S6 验证关门  → C6：e2e（I-001 ✓ closed）/ 活栈（I-002）/ 自审 / 台账同步（部分完成 —— 不闭门）
```

## 信息需求登记（P-005）

| 编号 | 问题 | 级别 | 影响门禁 | 状态 | 证据/结论 |
|------|------|------|----------|------|-----------|
| I-001 | Playwright e2e（admin/mvp profile × sqlite）在本次改动后是否仍全绿？ | required | C6（关门） | **closed** | 2026-08-23 双 profile 全绿（admin 9/9、mvp 9/9，另各 1 profile 专属跳过，exit 0）；e2e 暴露后端缺陷「删用户遗留 user_roles 孤儿 → 角色永久不可删」（`DeleteUser`/`DeleteUsersBatch` 不清理关联），已修复 + 2 项单元回归；证据 E-004 + attachments/I-001-evidence.md |
| I-002 | 活栈（compose 或本地双进程）体感与请求计时对比（DevTools/慢网络），10→3 请求数是否带来可感提升？ | non-blocking | C6（关门） | **open** | 未做 |

## 边界与审计声明

- 改动仅限页面数据请求层 / schema 文档缓存 / SQLite 连接面；postgres 方言（生产权威）零改动。
- 大表 COUNT(*) 容量优化出局（用户书面裁决，D-002）；monitoring 事件表随手动刷新更新（定向刷新只覆盖 /status，语义取舍留痕 D-002）。
- 审计模式 `self`（低风险可逆、无门禁语义变化、全自动回归兜底）；**本波尚未关门、尚无 A 条目**——S6 完成后追加审计与关门。