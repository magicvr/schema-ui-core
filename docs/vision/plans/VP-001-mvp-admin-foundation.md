---
doc_type: vision-plan
id: VP-001-mvp-admin-foundation
title: MVP Admin 基架
status: planned
vision_ref: schema-ui-core-admin-foundation@0.1.0
lead_workspace: null
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
parent: null
---

# VP-001 · MVP Admin 基架

## 意图

初始化可 fork 的 React 前端与 Go 后端 Admin MVP。该波次以固定的 `schema-ui-docs@v2.7.0` 协议为实现边界，完成核心账号与权限能力，并让每一纳入范围的协议功能都有可观察的范例页面和验证路径。

## 协议固定引用

| 项 | 值 |
|----|----|
| source | https://github.com/magicvr/schema-ui-docs |
| release | `v2.7.0` |
| commit | `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` |
| manifest | https://raw.githubusercontent.com/magicvr/schema-ui-docs/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/protocol-manifest.json |

上述引用固定协议来源。本地实施清单与 React/Go/范例/验证映射见 [protocol-inventory-v2.7.0.md](../protocol-inventory-v2.7.0.md)（`F-V001` → `fixed`，2026-07-31）。

**仍未冻结**：哪些协议能力域纳入本 VP 的 MVP 覆盖子集。冻结与实施门禁核验在开区后由 **`/govern`** 完成；在此之前不得主张“支持全部协议功能”。

## 方向级退出判据

在同时满足下列方向时，本 VP 可以提议关门；证据必须在挂接工作区的目标记录中：

1. React 前端与 Go 后端构成可运行、可 fork 的基础 Admin 工程，并以固定协议版本为兼容边界。
2. 由受控的协议清单定义 MVP 覆盖范围；其中每项都有前后端实现、范例页面或场景，以及可执行的验证路径。
3. 核心账号与权限链路具有可验证的前后端集成，不依赖未声明的业务领域模块。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | — | — | `planned`，尚未建立工作区。 |

## 关门记录

仅在 `closed` 或 `abandoned` 时填写。

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-07-31 | 用户确认首个 MVP 意图、技术方向、协议来源和方向级退出判据。 |
| 2026-07-31 | `/vision` 响应 VRev：链接协议清单；明确覆盖子集未冻结。 |
