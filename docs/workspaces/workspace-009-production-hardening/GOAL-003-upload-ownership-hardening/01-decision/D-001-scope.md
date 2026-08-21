---
id: D-001
goal: GOAL-003-upload-ownership-hardening
title: 立项与修复范围（上传所有权 + ReadHeaderTimeout）
date: 2026-08-10
status: recorded
---

# D-001 · 立项与修复范围

## 决策

用户于 2026-08-10 指示：在 **workspace-009** 添加子目标承载治理上下文，并修复本轮安全审视确认的问题。

本目标（GOAL-003）承接：

1. **High — 上传文件认证后 IDOR**：`GET /api/files/{id}` 任意登录用户可读 → 绑定 `owner`，下载校验 identity。
2. **Medium — 上传缺权限门**：保持「已认证即可上传」（产品仍允许通用附件上传），但**必须**绑定 owner，使下载不再全局可读。不在本目标引入新的 `files.write` 权限键（避免扩大 RBAC 面与 schema 联动）。
3. **Low — `ReadHeaderTimeout`**：在 `server.New` 设置，缓解慢速 header 耗连接。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 新建 `files.read`/`files.write` 权限并改 schema | 超出本波最小可验证修复；会牵动 RBAC reconcile 与页面 schema |
| admin 默认可读全部上传 | 扩大隐式提权面；本波默认 owner-only |
| 同期迁移 refresh 到 HttpOnly cookie | 中长期架构债；本目标登记 residual，不实施 |

## 背景

- 上一波 GOAL-002 已关门（16/16），Root/VP-009 曾关门；本波为**新审视发现**的访问控制缺口，不是 C1–D8 重开。
- Root GOAL-001 将重开并增加 S2 检查点以挂接本子目标；VP-009 是否重开由愿景层另议，本决策仅覆盖实现层工作区。
