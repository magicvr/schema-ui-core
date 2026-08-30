---
id: D-001
title: 信息裁决：I-025-001 / I-025-002 / I-025-003（用户 2026-08-30 采纳建议）
date: 2026-08-30
status: accepted
---

# D-001 · 信息裁决（2026-08-30 · P-004 / P-005）

> 2026-08-30 用户界面裁决：两条 required（I-025-001 / I-025-002）全部**采纳建议**，I-025-003（non-blocking）随 lead 提案确认。对应 I-025-001/002/003 → `verified`。合同正文见 D-002。

## I-025-001 · 配置包内容边界与密钥处理（required）

**现状事实**（只读扫描，未改任何代码）：

1. `apps/api/server/config.default.yaml`（serve 壳内嵌默认）：`app`（name/env）· `http`（addr/read/write/idle/shutdown_timeout）· `db`（dialect/path）· `profile: admin` · `auth`（access_ttl/refresh_ttl/jwt_secret/public_base_url）· `admin`（initial_password）· `log`（level）。
2. env 插值规则（文件注释 + `server/config.go` 装载）：`$VAR`（无默认，未设置 fail-closed）/ `$VAR:-default`（如 `${APP_ENV:-development}`）。
3. 敏感键现状：`auth.jwt_secret`（`${AUTH_JWT_SECRET:-}`）与 `admin.initial_password`（`${ADMIN_INITIAL_PASSWORD:-}`）为凭据类键；W15 加固后默认监听回环地址。

**建议（采纳即冻结）**：

- **包 = 非敏感结构键全集**（`config` 段与 `config.default.yaml` 树形同构，减去敏感键）；`profile` / `package` 元数据为信息性字段（导入不改变目标 Profile 默认集）。
- **env 引用保留 `$VAR` / `$VAR:-default` 形态，不解析为字面值**（可移植 + 不泄密 + 往返一致）。
- **敏感键处理**：默认敏感清单 = `auth.jwt_secret` / `admin.initial_password`；保守规则 = 键名含 `secret`/`password`/`token` 即视为敏感并列入 `secrets.exclude`。导出时值位置 = 占位（`$VAR` 形态），`secrets.exclude` 记录键路径 + 所需 env 名。
- **导入 fail-closed**：按 `secrets.exclude` 检查所需 env 可解析性（缺失 → dry-run 预检失败、拒绝导入；不回落默认、不写入空值）。

## I-025-002 · 落地形态（required）

**建议（采纳即冻结）**：

- **CLI 主路径**：`schema-ui config export|diff|dry-run|import`（接 VP-023/024 CLI 产线，`schema-ui` 已有 create/add/upgrade/migrate-fork）。
- **管理面本波不做**：配置包是运维/自动化动作，非日常设置；不重开 VP-007 Settings 面。
- **输出**：yaml/json 双格式（`-f` 选择，默认 yaml）；diff 输出机器可读。

## I-025-003 · diff 语义与输出（non-blocking）

**lead 提案（随裁决确认）**：键级规范化差量（固定键序、逐键比较，忽略信息性元数据 `exported_at`），变更条目 = `add` / `modify` / `remove` + 键路径 + old/new（敏感键只显示占位）；退出码 `0` 无差 / `1` 有差 / `2` 错误。

## 裁决结果（2026-08-30 · 用户界面裁决，全部采纳建议）

- **I-025-001 → `verified`**：非敏感结构键全集 + env 保留形态 + 敏感键占位/清单 + 导入 fail-closed（上方建议原样冻结）。
- **I-025-002 → `verified`**：CLI 主路径四子命令；管理面不做；yaml/json 双格式。
- **I-025-003 → `verified`（lead 提案 · 随裁决确认）**：规范化键级差量 + 退出码 0/1/2。

合同正文 = `01-decision/D-002-config-package-contract.md`（C2）。