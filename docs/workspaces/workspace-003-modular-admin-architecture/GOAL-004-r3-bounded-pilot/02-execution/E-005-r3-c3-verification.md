---
id: E-005-r3-c3-verification
doc: execution-entry
goal: GOAL-004-r3-bounded-pilot
date: 2026-08-05
status: recorded
---

# E-005 · R3 C3 验证事实

## 自动化回归

在本地 dirty snapshot 执行：

- `go test ./...`（`apps/api`）：通过。
- `npm test -- --run`（`apps/web`）：24 个测试文件、495 个测试通过。
- `npm run build`（`apps/web`）：Vite production build 通过。
- `gofmt -w internal/composition/composition_test.go` 后，组合根恢复测试通过。

## 快照恢复演练

`TestMVPRecoveryRestoresOptionalModuleDataAndCoreReadiness` 创建 Admin 写入，
关闭数据库并复制已知良好 SQLite 文件；随后修改 live copy 模拟失败发布，
把快照复制到新恢复路径，以 MVP Profile 重启并核对：Settings 关键字段、
已有 operationlog 行、恢复后 operationlog 写入/读取、`readyz`、Manifest、
Settings/Activity Schema 和 HTTP 路由。测试通过，禁用模块的数据没有被删除。

## 同一 Web 镜像双 Profile 运行矩阵

命令序列：

```text
docker build --file apps/web/Dockerfile --tag schema-ui-core-r3-check .
docker build --file apps/api/Dockerfile --tag schema-ui-core-api-r3-check apps/api
```

最终镜像：

- Web：`sha256:fff440c5afb29e20d220dfb008e5c1fdf5dc6af1559fd7de3ed09e3e9606916c`
- API：`sha256:747879ea1df144b7561628cf1277242fc6bc49c31a7bfce87c30c4ce35f62ebc`

以同一个 Web 镜像 `sha256:fff440...` 分别连接 MVP/Admin API 容器，结果为：

| Profile | Manifest pages | settings/activity Schema | settings/operations HTTP | `/readyz` | Web landing |
|---------|----------------|--------------------------|--------------------------|-----------|-------------|
| MVP | 无 `settings`、`activity` | 404 / 404 | 404 / 404 | 200 | 200 |
| Admin | 有 `settings`、`activity` | 200 / 200 | 401 / 401（需要认证） | 200 | 200 |

两次 Manifest 都返回 `X-Schema-UI-Manifest-Source: api`。Web 镜像内静态
`/.well-known/schema-ui/app-manifest.json` 已不存在，且用临时 network alias
`api` 执行 `nginx -t` 成功。临时容器和 network 已清理；镜像保留为本地证据
对象，不代表发布资产。

## 证据边界

上述结果证明当前本地 snapshot 的代码、构建和运行行为；不声称 hosted CI、
clean revision、部署或发布完成。
