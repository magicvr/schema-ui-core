---
title: 决策 · R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.1.0
---

# 决策 · GOAL-007

## D-001 · 用一个端到端目标实施 Root D-010

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 以单一 R4 子目标承载 records 精确契约、SQLite 持久化、CRUD API、Schema 读写交互、权限/错误状态及重启回归；六个成功标准按依赖顺序推进。
  2. 延续 Root D-010：代表实体固定为 `records`；生产默认迁入 SQLite；重启保持为 required 验收；错误响应继续使用 HTTP status + 稳定 `code` + message 的统一 envelope。
  3. 立项只确定路线与门禁，不在缺少证据时枚举新的 error `code`、DDL、并发策略或 Schema action 映射。
  4. `I-007-001`、`I-007-002`、`I-007-003`、`I-007-004` 均为 required；每项必须在表中所列首个受影响实施或验收动作前由证据关闭并记录后续决策。
- **理由**：API、持久化、Schema action 与重启证据共同构成一个可验证的业务生命周期；拆成多个并列目标会形成无法独立验收的中间态。把未知项显式登记为 required，可在保留端到端交付边界的同时防止方案被代码隐式冻结。
- **实施门禁**：当前四项信息均为 `open`。允许开展只读收集和方案设计；不得开始其影响范围内的产品代码变更，也不得据此勾选 S1～S6 或 Root R4。

### 未选方案

- **按 API / DB / Web / 测试拆成四个并列目标**：依赖紧密且成功边界不可独立成立，会增加跨目标门禁与中间态。
- **沿用进程内 records 并只补 Schema 页面**：无法满足 D-010 的 SQLite 与重启保持 required 边界。
- **立项时先猜精确 error code / DDL / action 形状**：会把尚未收集和验证的信息伪装为决定，违反 P-005。
