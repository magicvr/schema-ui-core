---
doc_type: goal-audit
id: A-002-contract-freeze-independent
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: close-out
scope: GOAL-002 R1 合同冻结全量（I-027-001/003/004 P-004 裁决核验 / D-002 §0～§11 ↔ kernel/ratelimit.go 逐节一致性 / 合同级快测覆盖 / 迁移不回归基线 / R1 越界核账 / go.mod 无 redis / P-005 信息门禁）
verdict: pass
open_required: 0
status: active
version: 0.1.0
---

# A-002 · R1 合同冻结独立交叉审计（independent）

> 编排器代贴（本地 grok build · grok-4.6 · 思考强度 high · 单轮 headless 输出），全文证据见 `attachments/audit-A-002-grok-output.md`；`source: independent` 保留，正文未改写。

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：close-out（合同冻结 + 端口落地 + 快测 + 越界）
- **scope**：`workspace-027-rate-limiter-port` / `GOAL-002-r1-contract-freeze`（R1 合同冻结；不含 R2 供应商与使用点迁移）
- **verdict**：**pass**
- **开放 required 计数**：**0**

## 范围与区间

| 项 | 值 |
|----|----|
| 工作区 | `workspace-027-rate-limiter-port`（`root_goal` = `GOAL-001-rate-limiter-port`；`canonical_scope` 已校验；`shared_materials_catalog: none`） |
| 被审目标 | `GOAL-002-r1-contract-freeze`（`parent: GOAL-001-rate-limiter-port`） |
| 冻结分母 | `01-decision/D-002-rate-limiter-port-contract.md` v0.1.0 |
| 裁决证据 | `01-decision/D-001-info-adjudication.md`（`status: accepted`） |
| 被审实现 | `apps/api/kernel/ratelimit.go`、`apps/api/kernel/ratelimit_test.go` |
| 迁移不回归基线 | `apps/api/internal/handler/rate_limit.go`（本波零 diff） |
| 对照 self | `03-audit/A-001-contract-freeze-closeout-self.md`（`verdict: pass`） |

**排除（明确不审、不放行）**：R2 内存供应商行为 / `-race` / 7 处使用点迁移；R3 Redis 接缝；R4 证据矩阵；I-027-002（最晚阶段 R2）。

## 独立复跑（2026-09-01）

| 命令 | 工作目录 | 结果 |
|------|----------|------|
| `go vet ./kernel/...` | `apps/api` | 退出码 **0** |
| `go test ./kernel/... -count=1` | `apps/api` | **ok**（0.756s） |
| `go test ./kernel/ -count=1 -v -run 'RateLimiter|DefaultRateLimiter'` | `apps/api` | 3 个测试 / 15 个子例 **全部 PASS** |
| `go build -o NUL ./kernel/` | `apps/api` | 退出码 **0** |
| `gofmt -l` | `apps/api` | 仅缺文件末尾 newline（无版式改写；已由编排器 A-003 fixed） |
| `git status --short` / `git diff --stat`（仓库根） | 仓库根 | 变更面 ⊆ 允许集 |

## 信息门禁核验（P-005 / P-004）

| ID | 级别 | 最晚阶段 | 独立核验 | 与 D-001 一致性 |
|----|------|----------|----------|-----------------|
| I-027-001 | required | R1 / C1 | D-001 `accepted`：**语义拆分保持**（Allow 不注册 + Record 失败计数 + RetryAfterSeconds + Clear + `now` 注入 + 工厂 `window/max/capacity`，capacity≤0 → `1<<16`） | **一致** |
| I-027-003 | non-blocking | R1 / C1 | D-001 采纳① **滑动窗口保持 + 策略接口独立** | **一致**（合同 §3；`RateLimiterInWindow` = `t.After(now.Add(-window))`） |
| I-027-004 | non-blocking | R1 / C1 | D-001 采纳① **本波不新增复合 key** | **一致**（合同 §2；key 不透明 `string`） |
| I-027-002 | required | **R2** | 待裁决 | **不阻断 R1**（最晚阶段未到） |

## 逐节一致性核验（D-002 §0～§11 ↔ `ratelimit.go`）

§0 范围外清单 ✓ · §1 端口形状逐字一致 ✓ · §2 key 不解析 ✓ · §3 `RateLimiterInWindow` 与既有 `allow` 剪枝代数等价 ✓ · §4 常量 `1<<16` 与 `newLoginRateLimiter` 默认一致 ✓ · §5 Retry-After 与 `retryAfterSeconds` 计算逐字同构；W12 七行常量表与现码行号/窗口/阈值一致 ✓ · §6 并发注释 ✓ · §7 无生命周期 ✓ · §8 红线 ✓ · §9 信息裁决 ✓ · §10 快测四类 ✓ · §11 未选未出现 ✓。

先例对照（`kernel/cache.go`）：同为 R1 端口 + 可执行谓词 + 编译期 stub；限流端口有意更薄（无 context / sentinel / 形状校验），不是漏搬 Cache 形态。

**命名漂移（不改判定）**：VP-027 退出判据 #1、Root 成功标准、workspace.md 仍写 `Allow/Record/Reset/RetryAfter`；冻结合同与代码为 **`Clear`**（见 F-001）。

**§3 剪枝路径注记（不改判定）**：合同把 `RetryAfterSeconds` 列为剪枝路径，与既有 `retryAfterSeconds` 不剪枝不完全同构（见 F-006）。

## 快测覆盖评估

对照合同 §10 预告快测面：编译期端口面断言 ✓ · `DefaultRateLimiterCapacity` 常量断言 ✓ · `RateLimiterInWindow` 表驱动 8 例（含 cutoff 恰等反例 / 零窗口）✓ · `RateLimiterRetryAfterSeconds` 表驱动 7 例（含 remain≤0→1 / 亚秒 Round 双向）✓。未测项（工厂 capacity 回落、allow 不注册、驱逐、`-race`）均属 R2 供应商行为。

## 迁移不回归基线

`git diff --stat -- apps/api/internal/handler`：**空**。7 处构造点零改动，行号/窗口/阈值/容量与 D-002 §5 表逐行一致（登录 15min/20、验证码 1min/10、密码 15min/5、恢复 15min/20、MFA verify 15min/10、MFA step-up 15min/5、邀请 15min/10；容量均 `1<<16`）。

## 越界核账（R1 波）

变更面 = `apps/api/kernel/ratelimit.go`、`apps/api/kernel/ratelimit_test.go`、`docs/workspaces/workspace-027-rate-limiter-port/**`。禁区（handler / go.mod / go.sum / kernel/profile.go / provider.go / internal/manifest / charter.md）全部未触碰；`go.mod` / `go.sum` 对 redis 无匹配；`ratelimit.go` 仅 `import "time"`。

## 对照成功标准（GOAL-002）

1～6 全部**满足**（端口契约冻结 + 快测断言；API 形态按 D-001；滑动窗口 + 无后台协程 + 策略独立；key 不透明；容量默认冻结；未越界）。

## Findings

| ID | 级别 | 严重度 | 主张 | 建议 |
|----|------|--------|------|------|
| F-001 | recommended | med | 上游仍写 `Reset`，冻结合同与代码为 `Clear`，R4 判据对名会误判 | C3 将 VP / Root / workspace 文案统一为 `Allow/Record/Clear/RetryAfterSeconds` |
| F-002 | recommended | med | self 意见未进入 `03-audit.md` 索引，台帐不完整 | `/govern` 登记 A-001 与本 A-002 |
| F-003 | recommended | low | 检查点算术与树注记不一致（C1 关闭但 progress 0/3；goal-tree C2 ✅ 与 meta 不同步） | 关门时按已关门检查点数重算 progress 并同步 goal-tree |
| F-004 | recommended | low | 两份新 Go 文件缺 EOF newline | `gofmt` 补换行 |
| F-005 | informational | low | VP-027 信息表与 workspace.md R1 行未随 D-001 回写 | 关门时回写 VP / workspace |
| F-006 | informational | low | §3 将 RetryAfterSeconds 列为剪枝路径，与既有 `retryAfterSeconds` 不剪枝不完全同构 | R2 实现时写明：剪枝仅 Allow，或 RetryAfter 剪枝后空列表的返回值 |
| F-007 | informational | low | `attachments/audit-A-002-grok-output.md` 为空占位 | 落盘 A-002 正文后替换该附件 |

**required / 必改：无。开放 required finding = 0。**

## 结论 + 建议给编排器

**verdict = pass**（0 required）。R1 合同冻结在 scope 内成立：P-004 三裁决与 D-001/D-002/代码一致；§0–§11 与 `kernel/ratelimit.go` 一致；快测覆盖合同可执行谓词与端口面；7 处构造点与 handler 基线零改动；变更面未越出允许路径；`go.mod` 无 redis。

建议：① 落盘 A-002 并在 `03-audit.md` 同时登记 A-001/A-002；② 合并响应后按检查点关门（progress 重算），回写 Root / VP-027 / workspace.md；③ `gofmt` 补 EOF newline；④ **不把本 pass 当 R2 放行**——I-027-002 仍阻断使用点迁移策略。

## 声明

- `source: independent`。本意见不修改目标 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree。
- 保证等级为框架默认 **L0**（入口分离），不是法定第三方鉴证。