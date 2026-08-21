---
id: E-003
goal: GOAL-008-w7-yaml-config
date: 2026-08-14
status: recorded
parent: GOAL-008-w7-yaml-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-003 · S2 实现完成（YAML 主配置）

## 事实

- 2026-08-14：S2 实现完成，提交 df5d440。
- 新增 `apps/api/configs/config.yaml`（operator 模板）与 `apps/api/internal/config/config.default.yaml`（go:embed 内置默认）；`apps/api/configs/.env.example`（敏感值参考，`.env` 保持 gitignore）。
- `config.Load()` 重写为分层加载：CONFIG_ENV_FILE（默认 configs/.env，不覆盖进程 env）→ CONFIG_FILE（默认 configs/config.yaml，显式指定缺失则 fail-closed）→ embed 默认 → 字段默认；YAML 值支持 ${VAR}（fail-closed）/ ${VAR:-default} 插值（行级作用域，注释内 ${...} 不算引用）；进程 env 已设置时覆盖 YAML。
- yaml.v3 解析开启 KnownFields（未知键 fail-closed，防拼写错误）；LoadError 字段经 ValidateProd 使启动失败。
- upload 三字段（UPLOAD_ALLOWED_TYPES / UPLOAD_MAX_FILES_PER_USER / UPLOAD_MAX_BYTES_PER_USER）从 handler 包级 env 读取迁入 Config；RegisterUpload 增加变参 UploadOption（WithAllowedTypes / WithUserLimits），旧调用点零改动（默认值保持 1000 / 256MiB）。
- compose.yaml 注释同步（W7 说明 + upload 透传示例）。

## 验证（S3 前置）

- go build ./... 通过；config/handler/composition 及全量 go test ./... 通过。
- 实测三路径：纯 YAML（configs/config.yaml 存在）启动成功；生产 env 缺 AUTH_JWT_SECRET 启动被拒（fail-closed，32 字符校验）；无 config.yaml（临时目录）回退 embed 默认启动成功。
- 新增单测：CONFIG_FILE 值 + env 覆盖、显式 CONFIG_FILE 缺失 fail-closed、裸 ${VAR} fail-closed、${VAR:-default}、env 插值优先、未知 YAML 键 fail-closed、CONFIG_ENV_FILE 提供 secret / 不覆盖进程 env / 显式缺失 fail-closed、ValidateProd 暴露 LoadError。

## 遗留

- S3 正式记录（含双路径实测条目）与 S4 go 判定、S5 grok 关门审计待办。
