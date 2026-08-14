---
id: E-004
goal: GOAL-008-w7-yaml-config
date: 2026-08-14
status: recorded
parent: GOAL-008-w7-yaml-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-004 · S3 验证完成

## 事实

- 2026-08-14：S3 验证完成。

### 单元 / 集成（config 包）

- `go test ./internal/config/...` 全绿（含 W7 新增 10 项：CONFIG_FILE 值 + env 覆盖、显式 CONFIG_FILE 缺失 fail-closed、裸 ${VAR} fail-closed、${VAR:-default}、env 插值优先、未知 YAML 键 fail-closed、CONFIG_ENV_FILE 提供 secret / 不覆盖进程 env / 显式缺失 fail-closed、ValidateProd 暴露 LoadError）。
- `go test ./internal/handler/... ./internal/composition/...` 全绿（upload policy 注入路径）。
- `go test ./...`（apps/api 全量）全绿。

### 双路径实测（真实进程）

| 路径 | 配置 | 结果 |
|------|------|------|
| 纯 YAML | configs/config.yaml 存在 + 少量 env | 启动成功（Fx RUNNING） |
| YAML+env 组合 | configs/config.yaml + env 覆盖 + HTTP 冒烟 | readyz 200 + 登录 200 + accessToken（admin profile，seed password 经 env） |
| embed 回退 | 无 configs/config.yaml（临时工作目录） | 启动成功（内置默认生效） |
| fail-closed | APP_ENV=production 且无 AUTH_JWT_SECRET | 启动被拒（32 字符校验错误） |

### 回归判定

- apps/web：无代码变更（配置载体仅 api 侧）；env 覆盖优先保证 smoke.sh / e2e 的 env 注入方式不受影响（静态确认：两脚本均通过进程 env 启动 api）。

## 遗留

- S4 go 影响判定 + self 审计；S5 grok 关门审计（data/部署门禁，用户指定 grok 4.6 high think）。
