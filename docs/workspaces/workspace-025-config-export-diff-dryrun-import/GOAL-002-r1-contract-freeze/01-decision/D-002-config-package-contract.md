---
id: D-002
title: 配置包合同 v0.1.0（冻结 · Admin 功能 · VP-025 R1）
date: 2026-08-30
status: accepted
---

# D-002 · 配置包合同 v0.1.0（2026-08-30 冻结）

> 责任文件（frozen）。实现（R2/R3）与验收以本合同为分母。本波不实施任何超出本合同的改动；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；不改迁移台账；不改 Charter。

## 0. 适用与对象面

- **对象面 = serve 壳配置树**：`apps/api/server/config.default.yaml`（内嵌默认 · `profile: admin`）+ `apps/api/server/config.go` 装载（env 插值 `$VAR` fail-closed / `$VAR:-default`）+ 骨架模板 `config.yaml.tmpl`。
- **CLI 产线**：`schema-ui`（VP-023/024 已交付 create/add/upgrade/migrate-fork）；本波新增 `config` 子命令族。
- **范围外**：管理面 UI（VP-007 Settings 不重开）；配置中心 / 远程分发 / 订阅拉取（Charter 非目标延伸）；运行时热加载；Secret Provider / KMS（RT-K02，仍 `trigger-gated`）；业务域配置；多租户/组织级配置。

## 1. 配置包格式（`config-package` v1）

YAML 文档，结构：

```yaml
package:
  format: schema-ui-config-package   # 固定字面量
  version: 1                         # 包格式版本（本波 = 1）
  app: schema-ui-app                 # 信息性：源实例 app.name
  env: development                   # 信息性：源实例 app.env（导入不得据此改目标 env 语义）
  profile: admin                     # 信息性：包对应 Profile（导入不得改变目标 Profile 默认集）
  exported_at: <RFC3339>             # 信息性：导出时间（diff 忽略该键）
config:                              # 与 config.default.yaml 树形同构的非敏感结构键（键序规范 §2.2）
  app:
    name: schema-ui-app
    env: ${APP_ENV:-development}     # env 引用保留形态，不解析
  http:
    addr: "127.0.0.1:25080"
    read_timeout: 5s
    write_timeout: 10s
    idle_timeout: 60s
    shutdown_timeout: 10s
  db:
    dialect: sqlite
    path: ./data/schema-ui.db
  auth:
    access_ttl: 15m
    refresh_ttl: 720h
    public_base_url: ${AUTH_PUBLIC_BASE_URL:-}
  log:
    level: info
secrets:                             # 敏感键清单（永不出现明文；§3）
  exclude:
    - key: auth.jwt_secret
      env: AUTH_JWT_SECRET
    - key: admin.initial_password
      env: ADMIN_INITIAL_PASSWORD
```

- `config` 段 = 非敏感结构键全集；缺失键 = 使用内嵌默认（导出时按「既有显式值 or 默认值」标定，明细 R2 冻结）。
- `secrets` 段只含排除清单与所需 env 名；**配置包任何位置不得出现敏感明文**。
- 敏感键判定（I-025-001 冻结）：默认清单 = `auth.jwt_secret` / `admin.initial_password`；保守规则 = 键名含 `secret` / `password` / `token` 即视为敏感并列入 `exclude`。

## 2. CLI 命令面（I-025-002 冻结）

### 2.1 子命令与退出码

| 子命令 | 输入 | 输出 | 退出码 |
|--------|------|------|--------|
| `schema-ui config export` | `[-o <path>]`（缺省 stdout） | 配置包 YAML（§1） | `0` 成功 / `1` 失败 |
| `schema-ui config diff` | `<pkg-a> <pkg-b>` 或 `<pkg> --against <运行配置>` | 键级差量（§2.2） | `0` 无差 / `1` 有差 / `2` 错误 |
| `schema-ui config dry-run` | `<pkg>` | 校验报告 + 影响列表（§2.3） | `0` 通过 / `1` 预检失败 / `2` 错误 |
| `schema-ui config import` | `<pkg> [-file <path>]` | 应用结果报告（§2.4） | `0` 成功 / `1` 失败 / `2` 错误 |

- `-f yaml|json` 双格式（默认 yaml）；错误走 stderr + 结构化消息；stdout 只写结构化产物（机器可读）。

### 2.2 diff 语义（I-025-003 冻结）

- **对象**：两包之间，或包 vs 运行配置（`--against`）。
- **规范化**：固定键序（按 `config.default.yaml` 段序）；逐键比较；**忽略信息性元数据**（`package.exported_at`、`package.env` 的差异不视为变更）。
- **变更条目**：`add` / `modify` / `remove` + 键路径 + old/new（敏感键只显示占位 `$VAR`）。
- **机器可读**：yaml/json 双格式；空差 = 空列表。

### 2.3 dry-run 语义（R1 基线；明细 R3 冻结）

- **只读预检，无任何写副作用**：① 结构校验（对照 §1 树形 + 类型/区间）；② 敏感键 env 可解析性（`secrets.exclude` 逐项检查；缺失 → 预检失败，**fail-closed**）；③ 影响报告 = 相对目标实例的将变更键列表（add/modify/remove，含是否触碰敏感键占位）。

### 2.4 import 语义（边界预告；R3 冻结明细）

- **前置**：dry-run 预检通过（含 env fail-closed）。
- **应用**：写入目标配置文件（`-file` 指定；缺省 = 骨架 `config.yaml` 位置）。
- **失败路径不破坏既有配置**：快照/回滚语义由 I-025-004 于 R3 冻结（候选 = 应用前快照 + 失败回滚，与 VP-013 方言级/`VACUUM INTO` 升级前快照关系在 R3 核对）。
- **不改 Profile 默认集 / 模块矩阵 / Manifest 装配**：导入只更新配置树，不触及装配语义（VP-008 `go` 红线）。

## 3. 敏感键与 fail-closed（I-025-001 冻结）

- 导出：敏感键值位置 = 占位 `$VAR`；`secrets.exclude` 记键路径 + 所需 env；任何含敏感明文的产物 = 合同违约（快测须断言产物无明文）。
- 导入 / dry-run：按 `secrets.exclude` 检查 env 可解析性；缺失 → 预检失败（fail-closed：不回落默认、不写空值）。
- 与 RT-K01 语义同构：配置装载侧 `$VAR` fail-closed 保留不变。

## 4. 红线与边界（全程）

- 不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义（VP-008 `go` 消费有效性；触碰即暂挂）。
- 密钥 fail-closed；热加载不进分母；管理面本波不做；配置中心/远程分发不进。
- 不改 Charter；不重开 VP-007 / VP-023 / VP-024。

## 5. 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| env 引用解析为字面值 | 未选 | 跨环境不可移植 + 扩大泄密面；往返不一致 |
| 敏感键明文进包 | 未选 | 违反 fail-closed；泄密 |
| 管理面主路径 | 未选 | 配置包是运维/自动化动作；VP-007 Settings 面不重开；CLI 可脚本化 |
| diff text-only 输出 | 未选 | 机器可读（yaml/json）优先，快测断言需要 |
| 导入自动回滚/快照 | R3 按 I-025-004 裁决 | 本合同只作边界预告（§2.4 候选） |

## 6. 验收方式（R2/R3 预告）

- **R2（导出+diff）**：① export 产物通过 §1 结构校验（快测 schema 断言）；② 产物无敏感明文（快测扫描 `jwt_secret`/`initial_password` 值位）；③ 两包（或包 vs 运行配置）diff 可断言（一致 → 空差/退出码 0；单键改 → `modify` 条目/退出码 1）；④ yaml/json 双格式等价。
- **R3（dry-run+导入）**：① dry-run 前后目标配置零变更（快照对比）；② env 缺失场景预检失败（fail-closed 快测）；③ import 往返（导出 → 干净实例导入 → 再导出 → diff 无差）；④ 失败路径按 I-025-004 冻结语义不破坏既有配置。
- **回归锁**：`go test ./...` 与 web 测试全绿；迁移台账 checksum 不变。

---

**引用链**：证据 → `GOAL-002/01-decision/D-001`（信息裁决）；实施责任 → R2（GOAL-003）+ R3（GOAL-004）；验收 → R4（GOAL-005）。