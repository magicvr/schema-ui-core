---
id: D-011-a021-response-metrics-position
doc: decision-entry
goal: GOAL-001-modular-admin-architecture
date: 2026-08-06
status: accepted
---

# D-011 · 响应 A-021：接受动态复审 pass，固定「指标 = 按需」立场

响应 [A-021](../03-audit/A-021-vp003-apps-code-independent-reaudit.md)
（independent · 动态代码复审，verdict `pass`、required 0、recommended 2）。
本决策接受其结论并处理两条 recommended：

## R-021-001（recommended）→ `fixed`

`apps/web/public/.well-known/schema-ui/` 与 `apps/web/dist/.well-known/schema-ui/`
本地空目录残留已删除（2026-08-06；两目录均 0 项、git 未跟踪，删除无版本库影响）。
生产路径不受影响：无静态 manifest 文件、Dockerfile 断言 `test ! -e dist/.../app-manifest.json`、
nginx 精确 `location =` 反代 API。

## R-021-002（recommended · 措辞澄清）→ `fixed`（决策留痕）

**决策**：指标（metrics）在本架构中为**按需能力**；当前实现**无指标贡献契约**，
亦无任何指标基础设施（grep `metric|prometheus|expvar` 于 `apps/api/internal` 零命中）。
已交付的 Observability 范围是：日志（带 `module_id`）、健康诊断（`/healthz` 存活、
`/readyz` 模块图 + system-data 就绪门控）。

因此 [module-architecture.md](../../../../architecture/module-architecture.md) §7 与
VP-003 exit #5 中「日志与指标均携带 `module_id` 语义」的表述应理解为**能力语义**
（指标若引入，必须携带 `module_id`），**不**表示当前已交付指标基础设施。按审计
建议的「Root 决策」通道落盘；未改动 architecture/VP 原文。

**后续触发**：未来若引入任何指标贡献或基础设施，必须新决策并登记信息项，
不得以本决策视为已交付。VP-003 exit #5 措辞的正式修订（若需要）属决策层，归
`/vision`，不阻断本响应。

## 状态影响

- A-021 无 required finding、与 A-018/A-019/A-020 无冲突 → Root `done / 6/6`
  与 VP-003 `active` 均维持现状，**不**回退、**不**放行 VP-003 `closed`。
- 本决策不改变任何信息项状态（I-001～I-007 仍 verified），不新增信息门禁。
