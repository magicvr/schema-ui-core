---
title: E-003 · W15 S5 F-007 落地与全量验证
status: active
created: 2026-08-30
updated: 2026-08-30
parent: null
version: 0.1.0
---

# E-003 · W15 S5 F-007 落地与全量验证

日期：2026-08-30 · checkpoint：`609cd6d6`

## 实施事实

1. **F-007 · LocalStore Unix 权限收紧**（`apps/api/internal/objectstore/local.go`，用户裁决 = fixed，D-002 §1）：
   - `Put`：目录 `MkdirAll 0o700`、对象临时文件与最终文件 `0o600`、sidecar `0o600`；同进程创建/读取、无静态直出路径，零破坏。
   - 既有文件不强制改写（无迁移 churn）；文档化 Docker 非 root 单用户卷拓扑下暴露本就低，此为本波有界加固。
2. **新增测试**（`internal/objectstore/local_test.go`）：`TestLocalPutTightenedUnixPermissions` — 非 Windows 平台断言目录 `0700`、body/sidecar `0600`（Windows 不强制 POSIX 模式，跳过）。

## 全量验证（S5 检查点）

- API：`go vet ./...` 0；`go test ./...`（apps/api 全量）exit 0。
- Web：`tsc -b` 0；vitest **1183/1183**；`vite build` exit 0（仅既存体积警告）。
- 部署检查：serve 默认回环 + 显式 env 语义已由 config 层负例锁定；既有 Docker 非 root 卷拓扑无需改动。

## 状态

S5 完成：F-007 落地；required 分母（F-001～F-006）+ recommended（F-007）全部实现并通过回归；I-005 已随 D-002 关闭（verified）。待 S6 分层审计闭合。