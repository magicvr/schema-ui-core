---
id: A-001
title: self 审计 · R2 导出+diff（合同一致性 / 红线 / 测试证据）
date: 2026-08-30
source: self
scope: GOAL-003-r2-export-diff（C1 export → C2 diff → 合同↔实现一致性 / 红线核对 / 测试证据）
verdict: pass
---

# A-001 · R2 导出+diff self 审计（2026-08-30）

## 范围

① 实现 ↔ 配置包合同 v0.1.0（GOAL-002 D-002）逐节一致；② 红线未越界；③ 测试/冒烟证据充分；④ R1 信息项状态核对。

## 核对结果

1. **合同一致性**：
   - §1 包格式：package 元数据（format/version/app/env/profile/exported_at）+ config 非敏感键 + secrets.exclude —— **符合**（D-001 §1/§3/§4；cfgTree 段序 = 源文件段序）。
   - §1 env 引用保留：`${APP_ENV:-development}` 原样入包 —— **符合**（TestExportDefaultShape）。
   - §1 敏感键：默认清单 + 宽规则（secret/password/token）+ 永不出现明文 —— **符合**（登记表 + 不变量；TestExportDefaultShape 的 config 段结构断言 + 凭据形态扫描）。
   - §2.1 CLI 面：export 0/1；diff 0/1/2；`-f yaml|json` —— **符合**（cliError；TestExportJSON / TestDiffErrors / 冒烟退出码）。
   - §2.2 diff 语义：add/modify/remove + 路径 + old/new；忽略信息性元数据；机器可读 —— **符合**（TestDiffIdenticalAndIgnoredMeta / TestDiffModify / TestDiffAddRemove；`--against` 同管线）。
   - §2.2 `--against`：包 vs 配置源 —— **符合**（TestDiffAgainst）。
   - §3 fail-closed 装载纪律：export 前 LoadConfig 校验（坏配置 exit 1）—— **符合**（TestExportBadConfigFails）。
2. **红线核对**：Profile 默认集 / 模块矩阵 / Manifest 装配零触碰（仅新增 CLI 命令、只读 `DefaultConfigYAML()`、main 错误处理扩展）；迁移台账零变更；密钥 fail-closed 语义未改；热加载未引入；管理面未做。**git diff 面可控**（见 E-002）。
3. **证据**：10 个单元测试全绿；`go test ./...` 49 包全绿；CLI 端到端冒烟（导出产物形态 + 退出码 0/1）实证留痕。
4. **信息项**：I-025-001/002/003 `verified` 维持；I-025-004 仍属 R3 前置（本合同 §2.4/§7 边界预告未动）。

## Findings

- （recommended）注册占位 `dry-run`/`import` 返回 `cliError{2}` 并标注「R3」——随 R3 实现自然闭合，无需先行处理。

## 结论

verdict **pass**（0 required）。R2 可关门；Root `progress 1/4 → 2/4`。

## 声明

本意见不直接修改 `status` / `progress`。关门与状态变更走 `/govern` 流程。