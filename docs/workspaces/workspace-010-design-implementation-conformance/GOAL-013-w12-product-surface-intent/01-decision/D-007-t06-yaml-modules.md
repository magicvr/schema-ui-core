---
id: D-007-t06-yaml-modules
doc: decision-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-007 · T-06 模块启用只认 config.yaml（S2 分项冻结）

## 背景

用户书面（P-004 Other）：**改成用 config.yaml 指定，不再使用命令参数、环境变量和 .env。** 新增一节指定加载哪些模块：可指向另一份 YAML（预设如 mvp / admin / demo / 其它自定义），也可直接写在 config.yaml。

W7 已有 `app.profile` + 逗号串 `app.modules_enabled`，且 **env 覆盖 YAML**。本条改的是**启用集的权威来源与形态**，不是从零做配置加载器。

## 决定

1. **权威**：进程启动后的模块启用集只来自 YAML。废除用 `APP_PROFILE` / `APP_MODULES_ENABLED` / CLI 开关选择启用集。
2. **新节**（实施时定名，候选 `app.modules`）：
   - **引用预设文件**：指向仓库内预设 YAML（至少 `mvp` / `admin` / `demo`；允许自定义路径）；
   - **或内联列表**：直接列出模块 ID。
   - 二者如何叠加（纯引用 / 纯内联 / 引用后再追加）S3 按「引用与内联互斥，避免静默合并」实现，除非 S2 补裁。
3. **预设内容**：`mvp` / `admin` / `demo` 的模块集合与现 `kernel.ResolveProfile` **保持一致**（本波不改默认集本身）。自定义预设 = 新文件，不改内置三档。
4. **闭合 I-004**：这是配置面 + 启用集**载体**变更，不是改 mvp/admin 默认成员。
5. **本条不废除** W7 对密钥等敏感项的 `${VAR}` 进程环境插值（JWT / 初始密码等）。用户「不再用环境变量和 .env」若指**全局**废除 env，见开放项 I-006，未裁前不拆 compose 密钥。

## 理由

- 用户要的是「开哪些功能写在文件里」，不是再记一套环境变量。
- 预设外置 YAML 便于 fork 复制一份 admin 再加减模块，而不改代码。
- 默认三档内容不动 → 不触发「改 Profile 默认集」的 go 暂挂；部署契约（不再认 APP_PROFILE）在 S4 书面判定。

## 未选方案

- **仅卫生、仍保留 env 覆盖启用集**：与「不再用环境变量指定」冲突。
- **改 mvp/admin 默认成员**：用户未要求。
- **只改文档**：达不到「用 config.yaml 指定」。

## 影响

- `apps/api/internal/config`、`configs/` 预设文件、README / QUICKSTART / compose 教学面。
- compose 里的 `APP_PROFILE` 改为挂载/选择 config 文件（如 `CONFIG_FILE` 或默认路径），**不再**用环境变量选模块。
- 审计：`self`；S4 做 go「部署契约变化、默认集不变」判定。

## 后续

- I-006（全局是否废除 env/.env）未裁则不得删除密钥插值。
- T-06 进 S3 P2。
