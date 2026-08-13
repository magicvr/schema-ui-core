---
id: A-002
goal: GOAL-007-r3-s02-file-library
source: self
date: 2026-08-14
scope: S2-S4 实现与验证
verdict: pass
parent: GOAL-007-r3-s02-file-library
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-002 · self 审计（S2–S4）

## 结论

**verdict: pass**（0 required findings）。

## 核对

- 安全面：下载/删除/确认端点全部 fail-closed（401/403）；id 格式白名单防路径穿越；下载头 nosniff + attachment + CSP sandbox + 文件名 sanitize；上传复用中心端点全部加固（大小/类型/活跃内容/配额）。
- 权限面：files.read/files.delete PolicyAdmin（admin-only，与 data-transfer 先例一致）；menu_files visibility PolicyAdmin。
- 审计面：三个 file 事件经 0018 进入冻结 CHECK；上传确认端点幂等（重复 ack 重复记事件为 best-effort 可接受）。
- 数据面：列表磁盘扫描单一事实源；删除 = 对象 + meta 硬删；重复删除 404。
- 迁移面：0018 checksum 由 Go 权威计算并与台账复算一致；17→18 计数断言全量更新；mvp 计数不变。
- 渲染扩展：runCustomAction 行 {id} 解析 fail-closed（无行 id 拒绝）；下载文件名取行 name（客户端 blob，无 header 注入面）；export 行为不变（有单测）。

## Findings

- 无 required。
- 建议（non-blocking）：(1) mtime 平局排序不稳定——D-002 §2 已文档化；(2) 目录扫描 O(files)/请求——配额已限规模，文档化接受。
