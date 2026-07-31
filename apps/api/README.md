# apps/api · schema-ui-core API 骨架

R1 Go 服务骨架（GOAL-003）。**无业务域路由、无鉴权中间件。**

## 要求

- Go **1.26+**（本仓 R1 在 Windows 实测 `go1.26.0`）
- Module：`github.com/magicvr/schema-ui-core/apps/api`

## 布局

```text
cmd/server/          # 进程入口
internal/config/     # 环境配置
internal/server/     # http.Server 包装
internal/handler/    # 路由（仅 /healthz）
pkg/version/         # 构建版本变量
```

## 运行

```bash
# 可选
cp .env.example .env

make run
# 或
go run ./cmd/server
```

默认监听 `:8080`（`HTTP_ADDR`）。

### 探活

```bash
curl http://localhost:8080/healthz
# {"status":"ok","timestamp":"...","version":"0.1.0","commit":"unknown"}
```

## 其它命令

```bash
make build   # 输出 bin/schema-ui-core-api
make test
make tidy
```

## 非目标（R1）

- 订单 / 钱包 / 通知等业务 API
- JWT / RBAC / SQLite 演示数据
- 协议兼容完成主张（见仓库愿景与 Root 信息门禁）
