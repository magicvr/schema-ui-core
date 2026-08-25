---
id: A-001
doc: audit-entry
goal: GOAL-002-iam-contract-freeze
source: self
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# A-001 · self · R1 合同冻结关门审计

## 条目头

| 项 | 值 |
|----|-----|
| source | **self** |
| 日期 | 2026-08-25 |
| scope | GOAL-002 整体：R1 IAM 合同冻结的完备性、可实施性与台账一致性（关门向） |
| verdict | **pass** |
| 开放 required finding | **0** |

## 核对成果

1. **信息门禁**：Root 信息表 I-001～I-009 全部关闭——I-001/I-002/I-009 `verified`（Root D-002，结构化裁决留痕）；I-003/I-004/I-005/I-007/I-008 `verified`（GOAL-002 D-001 输入裁决）；I-006 `registered`（2026-08-22 产品事实）。无 collecting / deferred required 残留；两处镜像表（Root `00-meta.md` ↔ `01-decision.md`）逐行核对同号同状态。
2. **合同完备性**：D-001 §1～§5 覆盖恢复状态机全生命周期（发起资格 / 投递目标恒为已校验邮箱 / 码形态 / TTL·冷却·尝试上限 / MFA 第二因子门 / 完成动作与会话撤销）、密码策略（配置面投影、默认参数、四口强制点、渐进生效）、邀请全链（双形态、邀请即建号、7 天一次性可撤销、重发撤旧发新）、会话语义投影。R2/R3 无需再猜产品语义。
3. **边界**：无 SMS / 模板中心 / 多邮箱 / 组织权限 / OIDC / 业务域越界；未写 DDL、未改应用代码；安全滥用面显式移交 VP-009/VP-010；管理员重置特权路径未动。
4. **检查点证据**：C1 = Root D-002 + 两轮裁决会话记录；C2 = D-001 五节条款 + E-002；本条目即 C3。
5. **对齐递归**：GOAL-002 → Root GOAL-001（R1）→ VP-019（active · VRev-043 pass）→ Charter @0.2.0，无冲突。

## Findings（notes，非必改）

| # | 级别 | 说明 | 处置 |
|---|------|------|------|
| N-1 | note | 恢复端点级请求节流（如 `/recovery/start` 投递轰炸面）未在合同单列数值——由挑战级限制（错 5 次作废 + 重发冷却 60 s）覆盖主体，端点 limiter 口径沿用仓库既有模式在 R2 设计落位 | 移交 R2 方案设计；不阻断关门 |
| N-2 | note | 审计事件（operationlog envelope 复用）与邀请权限键命名等实现侧命名未冻结数值 | 按 VP-012 envelope 与既有 permissions 模式在 R2/R3 设计定稿；不阻断关门 |
| N-3 | note | 「恢复完成是否直接签发会话」已在 D-001 §1 显式标记为 R2 设计定稿项（默认建议回登录页），非静默缺口 | R2 定稿时向用户确认一次即可 |

## 结论

R1 合同冻结达成成功标准：条款可核对、九项信息输入全部闭合、无越界、无开放 required finding。**verdict: pass** —— GOAL-002 可关门（C3 满），Root R1 记完成。
