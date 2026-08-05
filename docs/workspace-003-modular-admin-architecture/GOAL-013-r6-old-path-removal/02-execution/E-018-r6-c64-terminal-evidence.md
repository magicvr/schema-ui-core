---
id: E-018-r6-c64-terminal-evidence
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-018 · R6 C6.4 V01-V07 终态证据包

## 已发生事实

- 终态实现候选固定为 `9409b7176a5a07e60b9b07e3f2e1a2fc07ebf683`；实现提交链为
  `99784bc`、`88a3840`、`9409b71`。
- D-004 C64-V01～V07 已完成逐项动态验证，证据正文见
  [attachments/r6-c64-terminal-evidence.md](../attachments/r6-c64-terminal-evidence.md)。
- clean committed clone 固定候选 revision 后，API full/vet/build、Web `495/495`/build、
  `mvp` 与 `admin` Chromium E2E、clone-context Compose build/start 与 admin disposable
  smoke 全绿；clone 前后均 clean，总耗时 3.56 分钟。
- 候选 clone 构建出的同一 API image
  `sha256:75b9872626809cddf96f56a264843b00d97a549ac08a68011dfedbc9e375a013`
  与同一 Web image
  `sha256:3b89f8da0882104438ccc447c7acd88b708fa364874dfd09b02ddc955c1297bc`
  分别以 `mvp`、`admin` 启动，两个隔离 project 的 disposable smoke 均 SM-001～007
  全绿；另以同一 API image 和同一 volume 完成 `admin → mvp → admin` 回环，用户、
  Settings 与 operation-log 历史均保留。
- 所有 `c64*` disposable Compose project/volume 已删除；两个临时 clone 已送入回收站。

## 状态边界

- C64-V01～V07 已有本地候选 revision 证据；C64-V08 的 self + Grok independent 与
  `/govern` 响应尚未执行。
- 本地绿不等于 Hosted CI、合并 revision、部署或 VP-003 关闭；workflow 只证明已配置，
  不声称 GitHub Actions 已运行。
- R6-I004 保持 `collecting`，C6.4 不勾选，GOAL-013 保持 `active / 3/4`，Root 保持
  `active / 5/6`。

## 下一步（计划）

先执行 GOAL-013 C6.4 self close-out，再由 Grok Build `grok-4.5` / `high` 执行
independent `/audit`；经 `/govern` 响应全部相关意见后才可更新 R6-I004 与状态。
