---
id: GOAL-001-iam-recovery
doc: decision-entry
record_id: D-003
status: accepted
goal: GOAL-001-iam-recovery
created: 2026-08-26
updated: 2026-08-26
version: 1.0.0
---

# D-003 · 响应 Root A-001：两条 recommended findings 的处置取舍

## 触发

Root `03-audit/A-001-closeout-independent.md`（independent · close-out · verdict `pass` · 0 required）给出两条增量 findings：

- **F-001**（recommended / low）：`password_policy_settings.go` L136 死导入保持行 `var _ = errors.Is`；PATCH 失败路径一律映射 INTERNAL、无 sentinel 细分。
- **F-002**（recommended / info）：`loginRateLimiter` 为进程内内存桶，多实例部署时限流预算按节点各自计算；审计已明确「与既有 login 面同型、不在 workspace-019 边界内」，建议登记为部署拓扑注意项供后续生产化波次评估。

用户本轮明确指示：**修正 2 条 recommended**（`/govern 响应 GOAL-001 A-001 修正2条 recommended`）。

## 决定

1. **F-001 按 `fixed` 路径做真实代码修正**：
   - 删除死导入保持行，让 `errors.Is` 从「占位」变为「在用」；
   - 在领域层新增 sentinel `authsession.ErrPasswordPolicyNotSeeded`：`UpdatePasswordPolicy` 检查 `RowsAffected`，单例行缺失（legacy pre-0057 store）时**fail closed** 返回该 sentinel，不再静默 no-op 后谎报 200；
   - handler PATCH 失败路径按 sentinel 细分：`ErrPasswordPolicyNotSeeded` → **404 + 既有冻结码 `SETTINGS_NOT_FOUND`**（复用 settings 面惯例），其余存储错误维持 INTERNAL——**不新增任何错误码字面量**，不触碰 R1 冻结契约（漂移护栏测试继续通过）。
   - 补双向测试（领域层 unseeded/seeded + handler 层 sentinel 映射/成功路径）。
2. **F-002 按审计自身处方处置为「登记完成」**：
   - 不改限流器实现——分布式化属部署拓扑/生产化范围，越出本工作区边界；
   - 将部署拓扑注意项**持久登记**进本区台账（E-009），并指认后续生产化波次（自然归属 VP-009 production-hardening 程序）为评估责任位；代码侧 process-local 语义文档已存在于 `rate_limit.go` L12–16，无需重复。

## 为什么

- F-001 的「无 sentinel 细分」若以新增错误码方式处理，必须走冻结契约更新，代价与一条 low 级卫生 finding 不成比例；复用既有 `SETTINGS_NOT_FOUND` 既满足细分语义又零契约漂移。顺带修复的「unseeded store PATCH 静默 no-op」是同一路径上的真实潜在缺陷，属 fail-closed 改进而非行为破坏（现行迁移 0057 恒播种单例行，主路径零影响）。
- F-002 的审计原文即把处置定义为登记+后续波次评估而非本区内修码；把它包装成本区代码变更反而制造越界。登记落盘后该 finding 的可核对动作即告完成。
- 两条均非 required，不触发 P-004 强制裁决点；用户指令已覆盖处置方向，无需停等。

## 未选方案

- **F-001 新增专用错误码（如 POLICY_NOT_SEEDED）+ 契约更新**：语义最精确，但需动 Root D-002 appendix A 冻结集、errorcatalog 双语条目与契约测试，改动面远超卫生问题所需。未选。
- **F-002 本区内实现共享存储限流（DB/Redis 桶）**：跨入生产化拓扑决策，涉及新迁移与新依赖，且审计明示不在本区边界。未选；仅登记。
- **对两条均只留聊天响应不修码**：违反用户明确指示且浪费低成本改进窗口。未选。
