---
id: A-004-r5-a002-closure-independent
goal: GOAL-006-r5-maintenance-read-only-gate
doc: audit-entry
record_id: A-004
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-002 F-001/F-002 finding-closure；D-003 修订契约；I-002/I-003/I-004；Profile/Manifest/readiness/protocol 不变式
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-006-r5-maintenance-read-only-gate
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
responds_to: A-002
reviews: A-003
---

# A-004 · A-002 F-001/F-002 关闭复核（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：ad-hoc / finding-closure
- **scope**：GOAL-006 S0；仅复核 A-002 required F-001/F-002 是否已按 D-003 可核对闭合；顺带核对 I-002/I-003/I-004 设计充分性，以及 D-003 是否保持 Profile/Manifest/readiness/protocol 不变
- **verdict**：**pass**
- **required findings**：0

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）。
- **covered**：A-002 F-001/F-002 原文与证据；A-003 候选响应；D-002 历史提案 vs D-003 修订契约；I-002/I-003/I-004；A-002 引用的 Host/capability/错误消费实现。
- **excluded**：S1/S2/S3 实施与关门；A-002 F-003～F-006 作为本轮 required 重开；改写 D-003 / `00-meta` / `status` / `progress` / goal-tree / 业务代码；其他工作区上下文；共享资料内容（目录为 `none`）。
- **P-005**：本意见不改 `00-meta` 信息表。I-002～I-004 在登记上仍为 `collecting`；本轮结论是设计证据已足够，编排器可将其改为 `verified`。
- **本轮只读验证**：核对 D-002/D-003 用词差、`host-bootstrap.schema.json` mode enum、`capability-registry.json` 中 `form.controls.readonly` 语义、Host bootstrap/boot/failure/resource/auth-client、以及实现中仍不存在 `RUNTIME_MODE` / `SERVICE_*` / `availabilityMode`（S1 未开工，符合预期）。未运行 `go test` / e2e；不把未写代码当成已交付。

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-006 `parent`、`primary_plan` 一致；`goal-tree.md` 含本目标且 `status: active` |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 既有 Goal 审计 | A-001 self = pass（0 required）；A-002 independent = conditional（F-001 high / F-002 med）；A-003 self = conditional（候选 fixed，未自闭） | `03-audit.md` |
| P-004 冲突 | 无待裁冲突 | A-003 采纳 A-002 必改侧并留下 D-003；本条与 A-002 无 pass/fail 对撞 |

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| D-002 原文仍含被否投影，未被静默改写 | D-002 L31：`degraded/read-only → degraded + disabledCapabilities:["form.controls.readonly"]`；L24：`423 SERVICE_READ_ONLY`；L36：客户端使用 bootstrap recovery |
| D-003 已把 `form.controls.readonly` / 任何协议能力移出 `disabledCapabilities` | D-003 L26：后三者均**省略**该字段；「不把任何协议 capability 当作运行态开关」 |
| `form.controls.readonly` 仍是字段 `readOnly` 协议能力，不是运行态开关 | `docs/schemas/capability-registry.json` L134–140；`form-controls.ts` L336–349（缺能力 → `FORM_CAPABILITY_REQUIRED`）；`wallet.json` L6–16、`dictionary-entries.json` L6–16 已在 page meta 要求该能力 |
| Host `disabledCapabilities` 会从 `effectiveCapabilities` 剔除能力 | `bootstrap.ts` L278–284、L298–300、L317–318 |
| 现网 `READY_DEGRADED` 进应用且不走 Host 终态 | `boot.ts` L201–206、L256；`main.tsx` L150–156（只看 `failure !== null`） |
| D-003 将 degraded/read-only 映射为既有 bootstrap `degraded`，不扩展上游 mode | D-003 L26；`host-bootstrap.schema.json` L48–50 enum 仍为 `normal / maintenance / upgrade-required / degraded`，无 `read-only` |
| D-003 写拒绝改为应用内 503 catalog error，按 code 分流 | D-003 L21：统一 HTTP 503 + `SERVICE_MAINTENANCE` / `SERVICE_DEGRADED` / `SERVICE_READ_ONLY` |
| 现网 Resource 写失败走包络 `ResourceApiError`，不进 Host recovery | `resource.ts` L45–75、L168–182 |
| Host 5xx HTTP 映射仍是 `unavailable` 终态，因此不能拿来消费 degraded/read-only 写拒绝 | `failure.ts` L124–126 |
| maintenance recovery 仅适用于 bootstrap `maintenance` | `bootstrap.ts` L275；`boot.ts` L93–99；`failure.ts` L137–138；`HostFailureScreen.tsx` 标题按 `kind` 查表（L7–21、L92） |
| 现网 login 仍把任意 423 映射为 `ACCOUNT_LOCKED` | `auth.go` L126–128；`auth-client.ts` L328–331；`AuthContext.tsx` L147–152。D-003 将 read-only 从 423 改为 503，避开该碰撞 |
| 词典与实现中无 `host.readOnly` | 全仓检索无匹配；D-003 L27 明确不新增该文案依赖 |
| 空 env 助手仍把空字符串当未设置 | `config.go` L670–674；D-003 L31 禁止 `RUNTIME_MODE` 复用 `envOr` |
| 实现尚未出现 R5 运行态字段/码 | `apps/api` / `apps/web/src` 无 `RUNTIME_MODE` / `SERVICE_*` / `availabilityMode`；`bootstrap.go` L16–18、L45 仍固定 `normal` |
| health/readiness 探针职责未改 | `health.go` L59–77（`/healthz` 不碰库）；L72–107（`/readyz` 只表达存储 + 模块图） |
| D-003 书面保持装配/协议不变式 | D-003 L33：不改变 Profile 默认集、模块依赖闭包、Manifest bytes/聚合算法、协议 pin 或 health/readiness 语义 |
| A-003 未冒充 independent closure | A-003 L22、L37：候选 fixed，待本条复核 |

## 对照成功标准（本轮 scope）

| 检查项 | 状态 | 证据 |
|--------|------|------|
| F-001：`form.controls.readonly` 已移出 `disabledCapabilities` | **fixed** | D-003 L26 省略该字段；对照 D-002 L31 原文已替换 |
| F-001：degraded/read-only 使用既有 `degraded`，原始 mode 走 status | **fixed** | D-003 L25–26：`availabilityMode` 为权威区别 |
| F-002：degraded/read-only 写拒绝是应用内 API error | **fixed** | D-003 L21；现网 `ResourceApiError` 路径可承接 |
| F-002：maintenance recovery 边界仅 bootstrap `maintenance` | **fixed** | D-003 L21、L27；Host `MAINTENANCE` 仅 `mode === "maintenance"` |
| I-002 四模式写阻断 / HTTP / code / 重试语义 | 设计已闭合 | D-003 L19–21；登记状态仍 `collecting`，待 `/govern` 改 `verified` |
| I-003 认证优先级 / allowlist / 探针豁免 | 设计已闭合 | D-003 L19 保持 D-002 门禁位置与 allowlist；423 碰撞已用 503 消除 |
| I-004 Host availability 投影且不引入未知 capability | 设计已闭合 | 只用已有 `degraded` mode；省略 `disabledCapabilities`；不扩展 schema enum |
| Profile / Manifest / readiness / protocol 不变 | 保持 | D-003 L33；本轮未要求改 registry / bootstrap schema / health 探针 |

## Findings

本轮 **无新 required / recommended finding**。A-002 F-003～F-006 保持 recommended 实施门，已由 D-003 吸收为 S1/S2 验证要求，不阻断本条 closure。

### A-002 disposition

| 原 finding | 原级别 | 本轮状态 | 闭合路径 | 证据 |
|------------|--------|----------|----------|------|
| A-002 F-001 · `form.controls.readonly` 放入 `disabledCapabilities` 语义反转 | required / high | **fixed** | P-003 `fixed` | D-003 L26–27；A-003 L31；现网 capability 语义未变 |
| A-002 F-002 · degraded/read-only 重试语义误指 bootstrap recovery | required / med | **fixed** | P-003 `fixed` | D-003 L21、L27；`resource.ts` / `failure.ts` / `bootstrap.ts` 对照成立 |
| A-002 F-003 · `423` 与 `ACCOUNT_LOCKED` 碰撞 | recommended | 转为 S2 实施门 | 非本轮 required | D-003 已改 503 + 按 code 分流 |
| A-002 F-004 · `host.readOnly` 无法区分 mode | recommended | 转为 S3 UI 可选消费 | 非本轮 required | D-003 删除该 key，改 `availabilityMode` |
| A-002 F-005 · 空 `RUNTIME_MODE` 不得走 `envOr` | recommended | 转为 S1 实施门 | 非本轮 required | D-003 L31 |
| A-002 F-006 · 已注册方法与 allowlist 精确匹配 | recommended | 转为 S2 实施门 | 非本轮 required | D-003 L19、L32 |

## 信息项核对（P-005）

| ID | 级别 | 最晚阶段 | 登记状态 | 本审计结论 |
|----|------|----------|----------|------------|
| I-001 | required | S0 结束前 | verified | 维持；本轮不重开关 |
| I-002 | required | S0 结束前 | collecting | **设计可转 `verified`**：写矩阵、503/code、无 Retry-After、maintenance 与 degraded/read-only 消费边界已写入 D-003 |
| I-003 | required | S0 结束前 | collecting | **设计可转 `verified`**：门禁先于 handler 认证；allowlist 仍覆盖公开写与强制改密；GET/HEAD/探针放行；423 碰撞已消除 |
| I-004 | required | S0 结束前 | collecting | **设计可转 `verified`**：投影对象改为既有 `degraded` mode 且省略能力裁剪，不再引入未知 capability |
| I-005 | required | S1 实施前 | verified | 维持；本条即 F-001/F-002 的 independent closure |

无 `deferred`。无用户书面 `accepted-residual`。`00-meta` 证据列仍写「D-002 §2/§3，待 A-002」，属编排器更新债务，不是开放 required finding。

## 必改项汇总

**无。** required = 0。

A-002 F-001、F-002 已按 P-003 `fixed` 合法闭合。

## 与既有意见的异同

| 点 | A-002 | A-003 | 本意见 |
|----|-------|-------|--------|
| F-001 | required / open | 候选 fixed | **fixed**（D-003 可重复核对） |
| F-002 | required / open | 候选 fixed | **fixed** |
| I-002～I-004 | 不能 verified | 仍 collecting，待本条 | 设计充分；状态改写归 `/govern` |
| D-003 不变式 | 要求不改 Profile/Manifest/pin/readiness | 写入 D-003 L33 | 同意；现网 schema/探针未被本修订要求改动 |
| verdict | conditional | conditional（不自闭） | **pass** |

不是 P-004.2 冲突：A-003 未宣称已闭合，本条确认其候选成立。

## 结论 + 建议给编排器/用户的下一步

A-002 的两条 required 已由 D-003 完成可核对修正：运行态不再裁剪 `form.controls.readonly`，degraded/read-only 写拒绝改为应用内 catalog error，maintenance Host recovery 边界清楚，装配/协议/readiness 不变式保持。S0 设计门禁可以放行。

建议 `/govern` 下一句：响应 A-004，将 A-002 F-001/F-002 标为 `fixed`，把 I-002～I-004 改为 `verified`（证据改为 D-003 + A-004），将 D-003 标为 `accepted`，并放行 S1。不要改写 D-003 正文。

## 声明

本意见不修改 `status` / `progress` / D-003 / `00-meta` / goal-tree / 业务代码。响应与信息项闭合由 `/govern` 处理。
