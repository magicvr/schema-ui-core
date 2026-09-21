---
id: E-004-roadmap-registry-and-vf124-closure
doc: execution-entry
status: recorded
goal_id: GOAL-043-w31-cross-workspace-residual-closeout
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-design-implementation-conformance
version: 1.0.0
---

# E-004 · 路线图统一登记与 `V-F124` 闭合（C3）

## 1. 路线图统一登记（`docs/vision/roadmap.md`）

新增章节 **「未决项统一登记（残余 / 悬置 / trigger-gated）」**，位置在「当前组合焦点」之后，含：

- **用途声明**：本区是已交付范围之外未决事项的**统一登记处**——不是待办、不是承诺、不代表已验证；目的是让任何人一眼看到「它是什么 / 为什么不现在做 / 什么条件下做 / 谁负责 / 证据在哪」。
- **维护约定**：新增或闭合任何残余/悬置/触发项时**必须同步本节**（与 goal-tree、`03-audit` 同级要求）；登记只描述现状与触发条件，禁止把 deferred/recommended 写成已验证或已承诺。
- **三类登记表**：
  1. **有界残余（B 类）**：`GOAL-009 A-001 F-001`（fixed）、`F-002`（fixed）、`GOAL-005 A-002 F-002`/`R5-I-005`（fixed）、`GOAL-008 A-002 F-002`（bounded residual，触发=历史记录被再次引用为类型检查证据）、`V-F124`（fixed）。
  2. **悬置的范围决策（C 类）**：`I-037-005` 跨用户共享/最近/收藏/协作权限（触发=真实协作需求，owner `/vision`）。
  3. **未推进 / trigger-gated 能力（A 类）**：实体全文检索（`RT-X01`/`RT-X02`）、批量结果中心、组织·部门·岗位与 `org` 数据权限、新业务域、Redis/MQ/多实例（+第二持久化栈）、文件扫描/隔离策略、版本与维护提示、以及全部「扩展接缝」（typed domain event、Notification Transport、OIDC/SSO/SCIM、Approval Gate、Entitlement、多组织 context、SSE/WebSocket、外部连接器/Secret 产品面、自定义 metadata/tags、文件预览）——每条含现状、触发条件、责任人与出处。

同时更新 `roadmap.md` 的 VP-037 行与「Admin 功能最近一拍」段：残余不再逐个列在行内，改为指向本登记节；文件版本 `0.90.0 → 0.91.0`。

## 2. `V-F124` 闭合（愿景层）

- 新增 Vision Review **`VRev-097`**（self · `pass` · open required = 0）：逐项对照 `V-F124` 索要的五类内容与 R1 实际交付物（`r1-denominator-matrix.json` 24/58、`r1-form-matrix.json`、`r1-state-feedback-matrix.md` + `D-003`/`D-004`/`D-005`，`I-037-001`～`004` verified），确认**矩阵即 R1 交付物本身**，且其防范意图（未滑向共享视图/实体搜索/第二套基础设施）在关门审视中成立 → `V-F124` **`fixed`**。
- `docs/vision/reviews.md` 索引新增 `VRev-097` 行（含摘要）；`docs/vision/revisions.md` 新增 **`VR-083`**（editorial：VP-037 关门投影 + 未决项统一登记 + `V-F124` 闭合），版本 `0.4.56 → 0.4.57`。

## 3. 各投影同步

`docs/vision/workspaces.md`（workspace-037 行与叙述：残余改为指向登记节）、VP-037 计划（残余段与规划短史）、workspace-037 Root `00-meta` 备注、`GOAL-006 D-002` 追加说明——均改为「已收口/已登记 + 指针」，避免读者仍以为这些项原地开放。

C3 完成（`I-043-003`、`I-043-004` verified）。C4 见 `A-001`、`E-005`。
