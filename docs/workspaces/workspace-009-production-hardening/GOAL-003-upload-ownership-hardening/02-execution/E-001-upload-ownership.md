---
id: E-001
goal: GOAL-003-upload-ownership-hardening
title: 上传 owner 绑定 + 下载鉴权 + ReadHeaderTimeout
date: 2026-08-10
status: recorded
---

# E-001 · 上传 owner 绑定 + 下载鉴权 + ReadHeaderTimeout

## 事实

2026-08-10 实施并回归：

### 1. 上传绑定 owner（High IDOR 关闭）

- 文件：`apps/api/internal/handler/upload.go`
- `POST /api/upload`：从 `auth.IdentityFrom` 取 `user.ID`，写入 meta `owner`；无 identity → 401。
- `GET /api/files/{id}`：**仅当 `meta.owner == identity.ID` 才返回字节**；否则 403 `FORBIDDEN`。
- 缺 `owner` 的 legacy/损坏 meta → **fail-closed 403**（不公开读）。
- 既有 XSS 防护保留：主动内容标记、attachment、nosniff、CSP sandbox。

### 2. 回归测试

- `TestUploadOwnerOnlyDownload`：owner 200；另一登录用户 403；legacy 无 owner 403。
- 既有 `TestUploadEndpointContract` / `TestUploadRejectsHtmlAndForcesAttachment` 仍绿。

### 3. ReadHeaderTimeout（Low）

- 文件：`apps/api/internal/server/server.go`
- `http.Server.ReadHeaderTimeout` = `cfg.ReadTimeout`（≤0 时 5s）。
- `apps/api/internal/server/server_test.go` · `TestNewSetsReadHeaderTimeout`。

### 4. 测试命令与结果

```text
cd apps/api
go test ./internal/handler/ ./internal/server/ -count=1
# ok handler  ~13s
# ok server   ~0.7s
```

## 未做（本 E 范围外，见 00-meta residual）

- refresh token 迁出 localStorage
- 新建 `files.write` 权限键
- schema/manifest 匿名策略变更
- 全站 CSP/HSTS
