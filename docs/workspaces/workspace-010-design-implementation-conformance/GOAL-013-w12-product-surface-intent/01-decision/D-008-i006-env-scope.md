---
id: D-008-i006-env-scope
doc: decision-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-008 · I-006 环境变量范围 + S3 分批

## 决定

1. **只取消模块启用相关 env**：不再用 `APP_PROFILE` / `APP_MODULES_ENABLED`（及等价 CLI）选择启用集。
2. **保留** W7 敏感项与运维项：JWT、管理员初始密码、监听地址等仍可用进程环境或 `${VAR}` 插值。`.env` 不再作为选模块的方式。
3. **S3 分批**（用户未另裁，按建议冻结）：
   - P0：T-05 删除时间 + T-01 顶栏下拉
   - P1：T-03 个人中心 Tabs + T-02 列表搜索
   - P2：T-06 YAML 模块节
4. **闭合 I-006**。S1 / S2 检查点可勾选。

## 理由

全局废除 env 会把密钥推进 YAML 或另造文件协议，超出「用文件指定启用哪些功能」。

## 未选方案

- 配置值全部只写 YAML（含密钥明文）。
- 密钥改走独立文件、完全不用环境变量。
