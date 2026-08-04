---
id: E-003
title: 收集并核对 C1/C2 现状证据
status: recorded
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-002-r1-contract-migration-baseline
version: 0.1.0
---

# E-003 · C1/C2 现状证据

## 已发生事实

- 通过只读检索核对 API `handler.Register`、通用资源注册、权限 gate、Schema fixture、Web Manifest/route/navigation/permission 消费链，并记录在 `attachments/r1-c1-module-profile-inventory.md`。
- 明确记录了当前未发现 `apps/shell`、Profile registry、Profile 依赖矩阵和独立 Web route registry；这些是现状缺口，不是“无须实现”的结论。
- 核对 `Store.Open` 的迁移与 seed 顺序、0001～0008 编译链、schema ledger/checksum 校验、逐迁移事务、升级前 snapshot、完整性检查和直接测试覆盖，并记录在 `attachments/r1-c2-migration-seed-boundary.md`。
- 明确记录了当前未发现显式 tombstone、独立 system-data reconcile、应用层 rollback runner，以及部分失败路径测试覆盖缺口。
- 主编排器沿子代理返回的 file:line 抽查了上述入口和测试范围；子代理只提供定位线索，未写入文件。

## 状态边界

C1、C2 已完成本子目标的证据收集检查点，当前进度为 `2/4`。Root I-001、I-002 仍为 `open`；未运行测试命令，未宣称 R1 冻结或阶段通过。下一步继续形成 C3 生命周期/Fx 契约和 C4 协议矩阵。
