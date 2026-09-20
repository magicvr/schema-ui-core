---
id: GOAL-007-r3-pg-cross-version-restore-matrix
doc: audit
status: active
parent: null
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# 审计台账 · GOAL-007-r3-pg-cross-version-restore-matrix（R3-C）

> 本文件是唯一正式审计台账索引：`self` 与 `independent` **共用** `A-NNN` 序列。
> 每条意见正文在 `03-audit/A-NNN-<slug>.md`；本文件登记条目头（`source`/日期/scope/`verdict`）。
> 独立审计默认只写意见，不修改 `status`/`progress`/方案正文；响应归编排器。

## 意见索引

| A-ID | source | 日期 | scope | verdict | 摘要 | 文件 |
|------|--------|------|-------|---------|------|------|
| — | — | — | — | — | 尚无意见（本目标刚立项） | — |

## 待复审事项（编排器登记，供独立审计取证）

| # | 事项 | 证据位置 | 说明 |
|--:|------|----------|------|
| 1 | 组合定义是否真的覆盖 server×client 的关键组合，而非只跑「同版本」的舒适组合 | 组合定义 + 驱动脚本 | 反例优先：能否指出一个被跳过的组合及其理由 |
| 2 | 「supported / unsupported」判定口径是否事先冻结、是否可复现（命令、版本串、退出码） | 检查点 A 决策 + 记录 | 禁止事后按结果反推口径 |
| 3 | unsupported 组合是否**逐条**记录原因，而不是被「全部支持」概括 | 矩阵记录 | 检查每一行的证据列是否非空 |
| 4 | 升级后恢复有界核对是否真的跨版本（旧版本产物 → 新版本恢复），且有形状/校验证据 | 记录 + 证据附件 | 需能指回 C3 校验路径或 canonical 形状断言 |
| 5 | 破坏性动作是否严格限于一次性/专用测试 database（`D-017` 约束） | 驱动脚本 + 环境声明 | 宿主既有实例不得被修改 |
| 6 | 是否把 Docker 容器结果误当生产就绪证据（`D-017` §3 约束②） | 记录措辞 | 定位必须是 CI/reproducibility |
| 7 | `I-041-004` 的关闭/residual 是否有对应证据与范围说明 | `01-decision.md` 信息表 | residual 须写明范围与复审触发，且经用户书面接受 |
