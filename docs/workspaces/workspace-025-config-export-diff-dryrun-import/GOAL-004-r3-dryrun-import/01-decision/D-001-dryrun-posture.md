---
id: D-001
title: C1 dry-run 实现口径（lead · 合同 §2.3 派生）
date: 2026-08-30
status: accepted
---

# D-001 · dry-run 实现口径（2026-08-30）

配置包合同 v0.1.0（GOAL-002 D-002）§2.3 为分母；本条目落定实现层细节：

1. **结构校验 = 严格包解码**：包文件按 `configPackage` 结构 `KnownFields` 解码（未知键/多文档拒绝）；`package.format` 非空且 ≠ `schema-ui-config-package` → 拒绝；`package.version > 1` → 拒绝（宽容 0 = 未声明）。
2. **env fail-closed**：`secrets.exclude` 逐项——`env` 非空时 `os.LookupEnv` 未设置 → 预检失败（`不回落默认、不写空值`，合同 §3）；`env` 为空（源值无 `${}` 引用）→ 跳过（无所需环境变量可查，由 R3 import 写入路径保证占位语义）。
3. **影响报告 = 包 → 目标方向差量**：`changes = diffLeafMaps(targetLeaf, pkgLeaf)`（`old` = 目标当前值 / `new` = 包将应用值；remove = 目标有而包无）；目标配置 = `-config` 指定或内嵌默认。
4. **报告形态**：`checks[]`（path/status/message）+ `changes[]`（复用 diffEntry）；yaml/json 双格式；退出码 `0` 通过（可含 changes）/ `1` 预检失败 / `2` 工具错误（参数/IO）。
5. **零写副作用**：dry-run 全路径只读（不 create/rename/backup 任何文件）——快测以目标文件快照对比断言。

**未选方案**：宽松解析（跳过未知键）——Host 契约要求严格（未知键 = 拼写失误/格式漂移，属预检失败）。