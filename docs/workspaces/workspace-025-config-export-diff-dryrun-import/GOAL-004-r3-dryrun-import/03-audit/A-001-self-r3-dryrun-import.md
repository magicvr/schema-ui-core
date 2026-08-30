---
id: A-001
title: self 审计 · R3 dry-run + 导入（合同一致性 / 方案 A 落实 / 红线 / 证据）
date: 2026-08-30
source: self
scope: GOAL-004-r3-dryrun-import（C1 dry-run → C2 import → 合同↔实现一致性 / 裁决落实 / 红线核对 / 测试证据）
verdict: pass
---

# A-001 · R3 dry-run + import self 审计（2026-08-30）

## 范围

① 实现 ↔ 配置包合同 v0.1.0（GOAL-002 D-002）§2.3/§2.4 + import 语义决策（GOAL-004 D-002 · 方案 A）逐条一致；② 红线未越界；③ 测试/冒烟证据充分；④ 信息项状态。

## 核对结果

1. **dry-run ↔ §2.3**：只读预检三件套（结构校验 = 严格包解码 KnownFields + 多文档/格式/版本拒绝；`secrets.exclude` env fail-closed 不回落默认；包→目标方向影响报告 add/modify/remove）——**符合**；零写副作用（TestDryRunPass 目标文件快照对比）。
2. **import ↔ §2.4 + D-002（方案 A）**：
   - 预检前置：任一 check fail → 拒绝导入、目标零触碰 —— **符合**（TestImportRejectsAndKeepsUntouched 前两段）；
   - 备份：应用前 `<file>.pre-import.bak`（保留；备份失败 fail-closed）—— **符合**（TestImportBackup 内容 = 旧值）；
   - 原子应用：`.tmp` 写 → `server.LoadConfig` 装载校验 → `os.Rename` 替换 —— **符合**；
   - 失败路径：清 tmp、原文件未被触碰 —— **符合**（TestImportRejectsAndKeepsUntouched 第三段：目标原样 + 无 tmp 泄漏）；
   - 敏感键不入文件（env 注入路径）—— **符合**（生成头 + config 树无敏感段）。
3. **红线核对**：延续 `E-004-redline-precheck`（开区~R3 代码面仅 4 文件 · store/profile/upstream/迁移零触碰）；本波新增代码仍在 `cmd/schema-ui`（configpkg.go）与 `server/config.go`（只读导出）—— Profile 默认集/模块矩阵/Manifest 装配零触碰；迁移台账零变更；密钥 fail-closed 装载语义未改。
4. **证据**：cmd/schema-ui 18 用例全绿；`go test ./...` 49 包全绿；CLI 端到端冒烟（往返 diff `[]` exit 0 · 备份生成 · fail-closed 拒绝 exit 1）。
5. **信息项**：I-025-001/002/003/004 全部 `verified`（004 = 用户 2026-08-30 GUI 裁决方案 A · D-002）；I-025-005 `registered` 维持。
6. **判据映射**：判据 #3（dry-run 无副作用）与 #4（导入不破坏）交付面满足；判据 #1 跨实例往返闭环于冒烟实证。

## Findings

无 required；无 recommended。

## 结论

verdict **pass**（0 required）。R3 可关门；Root `progress 2/4 → 3/4`。

## 声明

本意见不直接修改 `status` / `progress`。关门与状态变更走 `/govern` 流程。