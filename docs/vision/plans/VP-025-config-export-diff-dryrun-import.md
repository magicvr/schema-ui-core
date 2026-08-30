---
doc_type: vision-plan
id: VP-025-config-export-diff-dryrun-import
title: 配置包导出 / diff / dry-run / 导入
status: active
vision_ref: schema-ui-core-admin-foundation@0.3.0
lead_workspace: workspace-025-config-export-diff-dryrun-import
created: 2026-08-30
updated: 2026-08-30
version: 0.2.0
parent: null
---

# VP-025 · 配置包导出 / diff / dry-run / 导入

## 状态与激活门禁（2026-08-30 · **active**）

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-08-30 激活 · 用户书面确认；VRev-054 self `pass`） |
| **lead_workspace** | `workspace-025-config-export-diff-dryrun-import`（Root `GOAL-001-config-export-diff-dryrun-import` **active** · 0/4：R1 合同冻结 → R2 导出+diff → R3 dry-run+导入 → R4 证据与关门） |
| **Vision required** | VRev-054 self `pass`（0 required；V-F089/090/091 recommended → 开区事务内闭合义务登记，见 [VRev-054](reviews/VRev-054-vp025-activation.md) 响应节） |
| **freshness 消费候选** | HEAD `055da2fd` · `apps/api/v0.3.0` · 六包（Admin 类 freshness PASS `c9122478` → `055da2fd` · 五域零变更 · 不暂挂 `go`） |
| **组合位置** | Admin 功能分支 · 基架能力剩余 #3（roadmap 明文点名非门控项） |
| **红线（激活即生效）** | 不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go`）；密钥 fail-closed；热加载不进分母 |

## 意图

在 RT-K01 既有配置系统（YAML + env 插值 · 密钥 fail-closed）与 VP-023/024 已交付的 CLI/包产线之上，把「配置包」运维能力收成可核对的 Admin 合同：

1. **导出**：当前生效配置树导出为**可移植配置包**；往返一致可核对（导出 → 干净实例导入 → 与源配置 diff 一致）。
2. **diff**：配置包之间、或包 vs 运行配置的**键级差量对比**；输出可核对、可机器读（快测断言）。
3. **dry-run**：导入前**预检**（结构/校验 + 影响报告），只读无副作用。
4. **导入**：预检通过后应用；导入前后实例可启动、可核对（回归快测）；失败路径不破坏既有配置（快照/回滚语义按冻结面）。

密钥与敏感值按冻结规则**排除或脱敏**，导出不泄密（fail-closed 保持）；**热加载不进退出分母**。

本 VP 属 **Admin 功能分支**（roadmap 明文「其后非门控未立项 = 基架能力剩余 #3」）。不重开 VP-007 / VP-023 / VP-024；不引入业务域配置语义。**不改 Charter**；**不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义**（触及即触发 VP-008 `go` 消费有效性暂挂，属本波红线）。

## 首波冻结（退出分母 = 配置包操作化）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| 导出 | 当前生效配置树 → 可移植配置包（文件产物）；往返可核对（快测/harness） | 配置中心 / 远程分发 / 订阅拉取；密钥明文进出包（保持 fail-closed）；Secret Provider / KMS（RT-K02 仍 gated）；热加载 |
| diff | 两包或包 vs 运行配置的键级差量，输出可核对（快测断言） | UI 配置市场；多租户 / 组织级配置 diff |
| dry-run | 导入前预检：校验 + 影响报告，只读无副作用（快测覆盖成功/失败路径） | 预检之外的写副作用模拟 |
| 导入 | 预检通过后应用；导入前后实例可启动、可核对；失败路径不破坏既有配置 | 改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；超出冻结面的回滚/备份语义（R1～R3 冻结，I-025-004） |
| 落地形态 | CLI（`schema-ui config *` 子命令，消费 VP-023/024 产线）与/或管理面——**R1 冻结**（I-025-002） | 新模块入 Profile 默认集；重开 VP-007 Settings 骨架 |

## 非目标

- **配置中心 / 远程配置分发 / 订阅拉取 / 运行时热加载**（热加载明确不进分母）
- **密钥明文进出包**；Secret Provider / KMS / HSM（架构 RT-K02 仍 `trigger-gated`）
- **改 Profile 默认集 / 模块矩阵 / Manifest 装配语义**（VP-008 `go` 红线）
- **业务域配置语义**（Catalog/支付等域配置归业务域分支）；多租户 / 组织级配置
- 重开 VP-007 / VP-012 / VP-023 / VP-024 已 closed 记录；替代 VP-009 / VP-010；改变 Charter 边界

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003** | 遵守薄内核。配置包是配置工具层能力，不建平行认证或中央注册路径 |
| **VP-007** | 管理面若进 Settings 走既有设置骨架；不重开 Settings / locale 交付 |
| **VP-008 `go`** | Admin 类能力；激活前做 Admin 类 freshness。若实现改变 Profile 默认集 / Manifest 装配，按消费有效性暂挂（本波红线：不改） |
| **VP-009 / VP-010** | 配置相关安全（密钥泄露、注入类 gap）与符合性 gap 归持续程序 |
| **VP-023 / VP-024** | 消费 CLI 产线与包形态（`schema-ui config *`）；迁移/分发工具化语义；不重开其 closed 记录 |
| **架构 RT-K01 / RT-K02** | 密钥 fail-closed 保持；KMS / Secret Provider 仍 gated，不进本波 |
| **业务域** | 域配置语义不进本波；业务域成立后可消费本合同 |

## 方向级退出判据

1. **导出闭环**：当前生效配置可导出为可移植配置包；往返（导出 → 干净实例导入 → 再核对）一致，密钥/敏感值按冻结规则排除或脱敏（快测 + 至少一条 harness/CLI 实证）。
2. **diff 可核对**：两包 / 包 vs 运行配置的差量输出可机器读并可断言（快测覆盖一致、仅差、冲突场景）。
3. **dry-run 无副作用**：预检覆盖校验与影响报告，成功/失败路径均有快测，不产生写副作用。
4. **导入不破坏**：预检通过后应用；导入前后实例可启动、回归快测通过；失败路径不破坏既有配置（快照/回滚语义按 I-025-004 冻结）。
5. **边界保持**：未改 Charter；未改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；热加载不进分母；密钥 fail-closed 保持。
6. **审计闭合**：开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root（P-001）书写：R1 合同冻结（包内容边界 / 密钥处理 / 落地形态）→ R2 导出 + diff → R3 dry-run + 导入 → R4 证据与关门。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

允许带未知立项。下列不影响「本 VP 意图已冻结」，但必须在对应阶段前关闭或经用户接受残余。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-025-001 | 配置包内容边界：包 = 当前生效配置树的哪些键；env 引用（不解析 vs 解析后值）；密钥/敏感值处理（排除 / 脱敏 / 占位 + fail-closed）。 | required | 方案冻结 + 退出判据 1 | R1 合同冻结 | 待裁决 |
| I-025-002 | 落地形态：CLI（`schema-ui config *`）vs 管理面 vs 两者；与 VP-007 Settings 面的关系。 | required | 方案冻结 | R1 合同冻结 | 待裁决 |
| I-025-003 | diff 语义与输出：键级规范化/排序/类型；输出格式（text / yaml / json 合一或分面）。 | non-blocking | 退出判据 2 | R2 | 待确认 |
| I-025-004 | 导入失败语义：预检失败即止 vs 应用期失败快照回滚；与既有升级前快照（VP-013 方言级 / VACUUM INTO）的关系。 | required | 退出判据 4 | R3 | 待裁决 |
| I-025-005 | 是否触及 Profile 默认集 / 模块矩阵 / Manifest 装配？**本 VP 冻结为不进**（VP-008 `go` 红线）。本行只作台账投影。 | required | 退出分母 | R1 | **registered**（VP 已冻结不进；激活时投影至 Root D-001） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-025-config-export-diff-dryrun-import | GOAL-001-config-export-diff-dryrun-import | lead | 2026-08-30（激活开区） | 唯一 delivery；激活审视 VRev-054 self `pass` + Admin 类 freshness PASS（`c9122478`→`055da2fd`）；Root active 0/4（R1～R4；D-001 已落痕） |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-30 | 初创 `planned`：用户确认立项（Admin 功能分支基架能力剩余 #3 · 配置包导出/diff/dry-run/导入；roadmap 明文「其后非门控未立项」点名）；退出分母 = 配置包操作化；不改 Profile 默认集/Manifest（VP-008 `go` 红线）、密钥 fail-closed 保持、热加载不进分母。roadmap 索引原子同步（并修正 VP-024 状态滞后） |
| 2026-08-30 | v0.2.0 · **激活**（用户书面确认）：[VRev-054](reviews/VRev-054-vp025-activation.md) self `pass`（0 required · Admin 类 freshness PASS `c9122478`→`055da2fd`，五域零变更，不暂挂 `go`；V-F089/090/091 recommended → 开区事务内 fixed）；`planned → active`；lead `workspace-025-config-export-diff-dryrun-import` 开区（Root active 0/4 · D-001 落痕） |