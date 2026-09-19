---
id: D-001-r1-contract-and-denominator-freeze
doc: decision-entry
parent: GOAL-002-r1-denominator-and-contract-freeze
status: accepted
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# D-001 · R1 分母与契约冻结（C1/C2/C3；用户 P-004）

用户 2026-09-19 书面裁决三问：maintenance 呈现 = **B**；升级入口 = **A**；Shell 可见性 = **A**。本决策把裁决落实为可实施合同，并区分**用户裁决**与**派生实现口径**。

## 1 · C2 `I-039-002` 四模式投影（用户裁决 B）

**用户裁决**：改 Host **生产者**投影，让 `runtime.mode=maintenance` 时 Admin Shell 仍能加载。

**派生口径（非第二轮产品分叉）**：

1. **不改** Host *消费者* 协议：`evaluateBootstrap` 对 `availability.mode=maintenance` 仍终态。不新增 Host enum，不改 pinned `host-bootstrap` fixtures。
2. **改** 本仓 bootstrap **文档生产者**（`handler/bootstrap.go`）：`runtime.mode` 为 `maintenance` / `degraded` / `read-only` 时，Host 文档 `availability.mode` **一律 `degraded`**（与现行 read-only 折叠同类）。`normal` 仍为 `normal`。
3. 精确模式不走公开 bootstrap。已登录通道：`GET /api/accounts/me` **additive** 字段 `runtimeMode` = 原样 `runtime.mode`（`normal|maintenance|degraded|read-only`）。无需新模块、不改 Profile。
4. 写门禁合同 **不改**：错误码与白名单保持 VP-012。
5. Shell 横幅：已登录且 `runtimeMode ∈ {maintenance, degraded, read-only}` 时显示；文案按精确模式三分，不把三者写成同一个「降级」。

冻结矩阵：`attachments/r1-mode-projection-matrix.md`。

**未选**：A（maintenance 继续 HostFailureScreen）；C（混合匿名终态）。

## 2 · C1 `I-039-001` 版本身份与升级入口（用户裁决 A）

**用户裁决**：授权用户看到版本；升级入口链到既有 QUICKSTART「schema-ui upgrade」说明；不建站内 changelog。

冻结：

| 项 | 合同 |
|----|------|
| 权威 | `pkg/version.Version` / `Commit`（ldflags）；与 `/healthz`、system-monitoring status 同源 |
| `BuiltAt` | 不进首波产品面 |
| npm/cli 六包版本 | **不是**本分母 |
| Shell 版本提示 | 仅 `monitoring.read`（见 §3）展示 `Version`；**不**在 Shell 展示 `Commit` / 模块清单 |
| 监控页 | 保持现有 version/commit（已有） |
| 升级入口 | `monitoring.read` 可见的文档链接；目标 = 仓库 `QUICKSTART.md` 升级节（默认 GitHub blob `main/QUICKSTART.md` 或 R2 冻结的等价公开 URL）。不新建页面、不自动升级、不远程拉发行说明 |

## 3 · C3 `I-039-003` 诊断分母 + Shell 可见性（用户裁决 A）

**用户裁决**：横幅给所有已登录用户；版本仅 `monitoring.read`。

诊断报告 = **既有 `admin.system-monitoring` 页**，字段分母见 `attachments/r1-diagnostic-field-matrix.md`。不新建诊断页。Shell 不复制模块清单/DB 大小。

mvp/demo 无 system-monitoring：仍能看到模式横幅（经 `/me.runtimeMode`），看不到版本 chip 与监控页。

## 4 · 排除与红线

- 不改 `ResolveProfile` / 默认集；不新模块。
- 不改 pinned `docs/schemas/**`、`apps/web/src/protocol/upstream/**`。
- 不热切换 `runtime.mode`；不实施 VP-040；不引入 Redis/MQ/多实例/搜索引擎。
- 不把 Host `availability.mode` 改成 `read-only`（非法 enum）。

## 5 · R2/R3 移交

| ID | 移交 |
|----|------|
| T-1 | bootstrap 生产者：maintenance→Host `degraded` + 回归（含 bootstrap 单测） |
| T-2 | `/api/accounts/me` additive `runtimeMode` + 测试 |
| T-3 | Shell 已登录横幅（三分文案、i18n、浅色深色、不挡登录页） |
| T-4 | `monitoring.read` 版本 chip + QUICKSTART 链接 |
| T-5 | HostFailureScreen maintenance 路径保持协议实现，但本仓生产 bootstrap 不再发出 `maintenance` 模式 |
