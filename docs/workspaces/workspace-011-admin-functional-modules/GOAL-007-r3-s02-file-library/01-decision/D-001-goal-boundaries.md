---
id: D-001
goal: GOAL-007-r3-s02-file-library
title: 立项边界：模块身份、Profile 归属与审计策略
date: 2026-08-14
status: accepted
parent: GOAL-007-r3-s02-file-library
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 立项边界（S-02 文件/附件库）

## 决定

1. **模块身份**：`admin.file-library`（标准 Admin 功能模块 · 共享能力，对齐 admin.data-transfer 先例）；Descriptor 依赖 core.auth-session / core.schema-render / core.navigation-capability / core.operationlog。
2. **Profile 归属（I-003 闭合）**：进入 **admin 默认集**（Profile 内容扩展，沿用 F-01/F-02 声明——不改装配语义）；**mvp / demo 不启用**（mvp 保持精简最小集，demo = mvp 集不变）。
3. **审计策略**：文件库涉上传授权/配额/所有权 → security/data 高影响门禁 → **独立审计**（grok build，P-004 在关门门禁与用户确认 provider）。
4. **列表数据源**：磁盘扫描（复用 uploadStore 单一事实源），不新增业务表；v1 不做引用注册表（I-002 边界见 D-002 §6）。
5. **审计事件**：新增 files.upload / files.download / files.delete 三个操作日志事件 → operationlog CHECK 扩展 **migration 0018**（归属 core.operationlog 台账，重建表模式同 0014/0015）。

## 未选方案

- DB 镜像表（files 表 + 迁移）：引入磁盘/DB 双写一致性问题；文件面单一事实源 = 磁盘 + meta，列表扫描即可满足 admin 工具规模。
- 进入 mvp 默认集：mvp 定位最小集，文件库属常用档管理面（同 data-transfer 仅 admin）。
