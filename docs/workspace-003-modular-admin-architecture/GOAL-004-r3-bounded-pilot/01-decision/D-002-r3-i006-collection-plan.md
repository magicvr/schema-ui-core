---
id: D-002-r3-i006-collection-plan
doc: decision
goal: GOAL-004-r3-bounded-pilot
date: 2026-08-05
status: accepted
---

# D-002 · R3 I-006 信息收集计划

## 决定

在冻结 R3 实施方案前，按三个 required 信息项收集静态/嵌入式 Manifest、
Nginx/Docker、中心注册、Host/Shell 特例、模块禁用和数据保留回滚证据。

## 收集动作

| 信息项 | 收集动作 | 必须形成的证据 |
|--------|----------|----------------|
| R3-I006-01 | 扫描源码入口；核对 `go:embed`、Vite/Nginx 代理、Docker 最终层；列出运行时路由 | 带 `file:line` 的入口清单和保留/移除边界 |
| R3-I006-02 | 明确开发期兼容窗口、告警位置和移除触发 | 可测试的期限/触发记录；生产无静态兜底证明 |
| R3-I006-03 | 明确禁用和清理失败的回滚触发、恢复步骤及数据约束 | 失败演练、数据计数/字段保留和恢复后端点核验 |

## 门禁

任何信息项缺少证据时保持 `collecting`。不得用既有局部测试、计划文本、
截图或候选产物替代同一构建上的实施/运行证据。

收集范围同时覆盖开发和生产路径：嵌入/静态 Manifest fixture、Nginx/Docker、
集中式 handler 注册、Host/Shell 特例、模块禁用以及数据保留回滚。单独的计划
或模块描述不能将任何信息项标为 `verified`。
