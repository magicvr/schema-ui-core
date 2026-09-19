---
id: GOAL-005-r4-result-center-experience
doc: decision
status: active
parent: GOAL-001-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.2.0
---

# 决策记录 · GOAL-005

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账正文在 `00-meta.md` 维护；长决策与独立决策记录放在 `01-decision/D-NNN-<slug>.md`。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-014 | required | 结果中心结构选型（既有 jobs 页收敛 vs 独立页） | C2 | C2 前 | 读页面结构/行操作/导航约定 | **verified** | — | **用户 P-004 裁决 = 方案 A**；`D-001` §0/§3 |
| I-038-015 | required | 管理作用域取消/重试的合同口径 | C1 | C1 前 | 读 repository actor 实现与 runner 取消路径 | **verified** | — | **用户 P-004 裁决 = 逐字镜像**；`D-001` §0/§1 |
| I-038-016 | non-blocking | 结果过期后的呈现与下载失效语义 | C2 | R4 前 | 对照既有错误呈现约定 | **verified** | — | 由派生字段 `downloadable` 关闭；`D-001` §4 |
| I-038-013 | non-blocking | 导出文件名/格式与两入口 UI 文案一致性（承接自 R3） | C3 | R4 内 | 对照既有导出约定 | **verified** | — | 单一实现 `lib/job-result-download.ts` + 4 例守卫；`E-001` §1 |
| I-038-010 | non-blocking | 导航分组与 i18n 键位（承接自 R2） | C3 | R4 内 | 读分组测试与 fragment | **verified** | — | 双目录键集合对齐（1227/1227）+ `schema-keys.structural` / 能力声明 guard 全绿；`E-001` §3 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-19 | R4 结构选型与写面合同冻结（含两项用户 P-004 裁决原文与派生口径） | active | [D-001-r4-structure-and-write-contract-freeze.md](01-decision/D-001-r4-structure-and-write-contract-freeze.md) |

## 约束输入（来自 R1～R3 冻结，非本目标可改）

- **R1 `D-001`**：契约归本地（方案 B）；ADR-0022 同步语义冻结；首波 1 条；`jobs.read` 管理作用域；K-3 结果三段语义（409/410/终态）。
- **R2 `D-001`**：O-1 自定义组件触发；O-2 不声明 `actions.batch.request`；路由前缀 `/api/jobs`；结果地址由 `jobs.ResultURL` 派生；`jobs.write` 归 R3。
- **R3 `D-001`**：`jobs.batch-export` 数据面（users/roles 同分母、上限 500、结果 JSON 含 csv）；`jobs.write` + `data.export` 双重门禁；进度语义；前端不 `reloadList`。

## 待冻结（本目标方案项）——已全部冻结（2026-09-19）

1. **结构选型**（`I-038-014`）：**方案 A** —— 既有 `jobs` 页原地收敛（行操作 + `recordView`），不新增页面/导航。**用户 P-004 裁决**，见 `D-001` §0/§3。
2. **取消/重试合同口径**（`I-038-015`）：**逐字镜像既有 actor 作用域合同**（同状态集合、同冻结错误码、不重置 attempt、`jobs.write`）。**用户 P-004 裁决**，见 `D-001` §0/§1。
3. **过期与下载失效的 UI 呈现**（`I-038-016`）：由服务端派生字段 `downloadable`（仅 `succeeded`）关闭；410 为直接访问兜底。`D-001` §2/§4。
4. **i18n/主题/可访问性收敛范围**（`I-038-013`、`I-038-010`）：下载单一实现 + 双目录对齐 + 整页主题 token 断言；`D-001` §3/§5，证据 `E-001` §3。

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。
