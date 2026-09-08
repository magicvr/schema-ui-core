---
id: VRev-082
doc_type: vision-review
title: VP-034 导航分组折叠体验 · 计划阶段意图审视
source: self
scope: VP-034-nav-group-collapsible · planned
verdict: pass
date: 2026-09-07
auditor: /vision (claude-sonnet-4-6)
created: 2026-09-07
updated: 2026-09-07
parent: null
version: 0.1.0
---

# VRev-082 · VP-034 导航分组折叠体验 · 计划阶段意图审视

## 审视范围

新立 VP-034-nav-group-collapsible（Admin 功能分支 · 导航分组折叠体验）的计划阶段意图审视：

- Charter 对齐检查
- 意图与退出判据合法性
- 结构选型合理性（新 VP vs VP-010 子目标）
- P-005 信息就绪
- 非目标完整性

## 审视结论

**verdict: `pass`**（0 required，1 recommended）

VP-034 的意图、退出判据、结构选型和 Charter 对齐均通过计划阶段审视。可以按 Admin 类 freshness review 流程推进激活。

## Charter 对齐检查

| 项 | 检查结果 |
|----|---------|
| `vision_ref` = `schema-ui-core-admin-foundation@0.4.0` | ✅ 精确匹配现行 Charter |
| 意图落在 Charter 边界内 | ✅ Charter 成功边界 #3「前端产品化·体验参考 Linear/Vercel」和 #5「增减模块不要求修改 Shell 中央注册路径」直接覆盖导航分组体验改进 |
| 不改变 Charter 目的/边界/非目标 | ✅ 导航分组是已有范围内的产品完善，不影响协议、架构、分发、业务域边界 |
| 与 Charter 非目标不冲突 | ✅ 非目标未禁止 Admin 体验增强；非目标中不涉及导航组件范畴 |

## 结构选型合理性

P-006 判定树路径已在用户确认前完整走过（`/vision` 会话，2026-09-07）：

| 问题 | 判断 | 依据 |
|------|------|------|
| 改 Charter 目的/边界？ | 否 | 导航分组是已有能力的体验增强，不改变愿景边界 |
| 同愿景新纲领波次？ | 是 | Admin 功能分支「体验增强」象限下有界交付 |
| 独立 goal-tree / 隔离 / 长期并行？ | 否 | 单次有界整改 |
| 塞入 VP-010 符合性程序？ | 不适合 | VP-010 语义 = as-designed vs as-built 偏差修复；导航分组是「向好演进」非「偏差整改」；混入会模糊 VP-010 长期程序语义 |
| 结论 | **新 VP + 新工作区**（交 /govern 建区） | 用户确认，2026-09-07 |

## 意图与退出判据审视

| 退出判据 | 可判定性 |
|---------|---------|
| 1. Shell 分组渲染（折叠/展开 + 键盘可访问） | ✅ 可量化验收 |
| 2. 模块注册 API 兼容（`group` 字段可选，已有模块零破坏） | ✅ 向后兼容可回归验证 |
| 3. 激活态感知（直接 URL 进入时对应分组自动展开） | ✅ 用户书面补充确认，可端到端验证 |
| 4. 至少一个内置模块示例（端到端验证） | ✅ 可量化 |
| 5. playbook 更新（注册规范含 `group` 说明） | ✅ 文档产物可核查 |

退出分母范围合理，有界，不承诺：强制所有模块迁移分组、服务端持久化折叠状态、多级嵌套、拖拽排序。

## P-005 信息就绪

| 项 | 状态 |
|----|------|
| 折叠状态持久化方案（会话内 vs localStorage vs 服务端） | `non-blocking`（执行阶段决策；退出判据 1 已涵盖「可在本次会话内持久化」作为最低下限）|
| 分组 key 命名约定（模块自声明 vs 中央注册） | `non-blocking`（执行阶段冻结；判据 2/5 确保 playbook 覆盖）|
| 键盘可访问性具体规范（ARIA 模式） | `non-blocking`（判据 1 要求「至少支持回车/空格」，具体 ARIA 方案执行阶段决策）|

无 `required` 级信息门禁阻断激活。

## Findings

### V-F121（recommended · 非阻断）

**分组 key 命名空间建议在 playbook 中明确**：模块自声明 `group` 字段时，建议 key 使用注册式字符串（如 `"iam"`、`"settings"`）而非自由文本，以避免多模块用相近 key 导致意外合并。建议在退出判据 5 的 playbook 更新中明确 key 命名规范（如 snake_case 小写英文、建议列举已有约定 key）。

状态：`open · recommended`（不阻断激活或关门，但执行阶段建议落入 playbook 范畴）

## 激活门禁（后续用 /vision）

激活前须完成：

1. **Admin 类 freshness review**（验证当前代码主线候选身份与 VP-008 `go` 消费有效性）
2. **self Vision Review**（激活就绪确认，可与 freshness review 合并为一次）

本 VRev-082 仅为计划阶段意图审视，**不是激活许可**。

## 声明

- 本 Review source = `self`，由 `/vision` 编排器完成
- 不改变 Charter / VP / Goal 任何 status
- open required = 0；V-F121 为 recommended，不阻断
- 建议下一步：交 `/govern` 建立工作区骨架并推进实现层；激活时经 `/vision` 完成 freshness + 激活就绪 Review
