---
id: E-002-r1-inventory-and-contract-decision
doc: execution-entry
parent_goal: GOAL-001-nav-group-collapsible
status: recorded
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# E-002 · R1 导航盘点与契约决策

## 已发生事实

1. 完成当前 API Provider、Manifest fragment、Profile 默认集合、Web navigation projection 与 Shell 渲染路径的代码盘点。
2. 将当前 sidebar 分组分母、profile/slot 覆盖和有意非 sidebar 面记录到 [R1 导航 / Profile / slot 盘点](../attachments/r1-navigation-profile-slot-matrix.md)。
3. 用户确认 VP-034 五组顺序、组内顺序、Dashboard 顶层单例、Examples 与 top/user slot 边界。
4. 用户确认 R2 采用契约优先方案；决定记录于 [D-002](../01-decision/D-002-r1-baseline-and-group-contract.md)。
5. API 基线验证 `go test ./... -count=1` 通过；Web 直接调用本地 Vitest `97` 个文件 / `1332` 个测试通过，直接调用本地 `tsc -b` 通过。pnpm 包装入口因环境的 `esbuild` ignored-build-script 保护提前退出，未把该环境阻断误报为代码失败；尚未宣称 R2/R3/R4 实现完成。
6. A-001 self 对 D-002 与 R1/R2 方案判定 `pass`；由于 scope 跨 API kernel、Manifest 聚合与 Web Shell，等待项目指定的本地 grok build independent design-plan audit。

## 当前门禁事实

- `I-034-001`：静态代码盘点已记录；运行时 Profile/Manifest 矩阵仍需 R4 harness 核验。
- `I-034-002`：用户产品决策已 verified（decision）；Manifest/Shell 实际呈现留到实现回归验证。
- `I-034-003`：non-blocking，继续 open，R3 决定会话内或浏览器持久化。
- `I-034-004`、`I-034-005`：尚未到最晚阶段，保持 open。

## 未完成边界

- `NavigationContribution` group 契约尚未写入代码。
- Manifest 跨模块分组归一化尚未实施。
- Shell 折叠/展开、键盘可访问、直接 URL 自动展开、状态保持尚未实施。
- 当前注册 sidebar 全量迁移与 default/optional/custom/demo 回归尚未完成。
