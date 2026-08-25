---
id: A-002
doc: audit-entry
goal: GOAL-001-iam-recovery
source: self
auditor: /govern 编排器（A-001 响应复核）
audit_type: finding-response
verdict: pass
created: 2026-08-26
updated: 2026-08-26
version: 1.1.0
---

# A-002 · A-001 findings 响应复核（F-001 fixed · F-002 登记闭合）

> 本条目是编排器对 [A-001](A-001-closeout-independent.md) 两条 recommended findings 的**响应与核验**记录（`source: self`），不改变、不替代独立意见原文。处置取舍见 D-003，实施事实见 E-009。

## 响应汇总

| Finding | 级别 | 处置 | 路径 | 证据 |
|---------|------|------|------|------|
| F-001 死导入保持行 + PATCH 无 sentinel 细分 | recommended / low | **fixed** | 代码修正 | E-009 §F-001：死行删除；sentinel `ErrPasswordPolicyNotSeeded` + `RowsAffected` fail-closed；handler 细分映射 404 `SETTINGS_NOT_FOUND`（复用冻结码，零契约漂移）；新增 4 测试 |
| F-002 进程内限流器多实例预算分摊 | recommended / info | **fixed（按审计处方登记）** | 台账登记 | E-009 §F-002：部署拓扑注意项持久登记（含语义要点与 VP-009 归属位指认）；代码侧 process-local 文档已在 `rate_limit.go` L12–16 |

## 核验（本会话实测，2026-08-26）

1. `gofmt -l` 四个触碰文件：干净。
2. 定向测试：authsession `TestUpdatePasswordPolicy*` ok（1.8s）、handler `TestPatchPasswordPolicy*` 两测 PASS（2.2s）。
3. **全量回归**：`go test ./... -count=1`（apps/api）**exit 0 全绿**——含错误契约漂移护栏 `error_contract_test.go`（未新增字面量；`SETTINGS_NOT_FOUND` 冻结集内既有发射点不变）、store 黄金断言（38.0s）、composition、r4_evidence 三链。
4. 行为影响面核对：迁移 0057 恒播种单例行（migration.go checksum 注释「singleton seed」），现行部署/测试路径零行为变化；仅 legacy pre-0057 未播种 store 从「静默 no-op 假成功」变为 fail-closed 404。

## 开放项

- 无开放 required。
- F-002 指向的**多实例共享限流评估**为遗留跟踪项，归属后续生产化波次（VP-009 程序域），不构成本工作区任何门禁；是否立项由该程序波次规划决定。
- **2026-08-26 更新**：该项已经用户指示正式立项 → [workspace-009] `GOAL-012-w12-multi-instance-rate-limiting`（评估先行 · active 1/4；Q2 路径见 E-009 §F-002 立项回执）。本响应记录的开放项就此全部移交/闭合。
- 建议：无需为此重开 `/audit` 复审（两条均非 required 且核验证据可重复）；若后续生产化波次启动，可在其 scope 内自然覆盖。

## 结论

**verdict: pass。** A-001 两条 recommended findings 已全部按 P-003 可核对路径处理完毕并留痕；Root 关门结论（done 4/4 · 开放 required = 0）不受影响，亦无需变更任何 status / progress。
