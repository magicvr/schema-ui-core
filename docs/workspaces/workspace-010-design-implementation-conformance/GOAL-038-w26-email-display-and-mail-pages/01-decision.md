---
id: GOAL-038-w26-email-display-and-mail-pages
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# 决策记录 · GOAL-038

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 全渠道出站记录的存储演进（`mail_outbox` 加列 vs 新表）、存量 mock 记录兼容、双方言 DDL、投递状态取值集 | S2 方案冻结 | S1 | 迁移设计写入 D-001；核对 0053/0054 ALTER 先例 | **closed（verified 2026-08-26）** | — | D-001 §2.1：0060 portable additive ALTER ×2，存量行默认值即真实语义；取值集冻结 |
| I-002 | required | 邮件控制台/出站记录两独立页面的导航挂载与权限复用方式（沿用 `settings.read`，不新设权限） | S2 方案冻结 | S1 | 勘察导航贡献先例后写入 D-001 | **closed（verified 2026-08-26）** | — | D-001 §2.2：admin.settings 贡献两页 + menu_mail/menu_mail_outbox（Permission=settings.read 复用）；users-invites/data-dictionary 导航先例核实 |
| I-003 | required | 用户列表批量返回邮箱状态的读取路径（避免 N+1） | S2 方案冻结 | S1 | 存储层勘察后写入 D-001 | **closed（verified 2026-08-26）** | — | D-001 §1：0054 列已在 users 表，ListUsers 投影加两列同查询完成，无 N+1 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-26 | W26 方案冻结：三问题修复设计（邮箱读面 / 全渠道出站记录 + 页面化 / 撤销绑定） | accepted | `01-decision/D-001-w26-design-freeze.md` |
