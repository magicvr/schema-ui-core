---
id: R1-C4-EVIDENCE
title: R1 C4 协议继承与模块候选三态矩阵
status: recorded
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-002-r1-contract-migration-baseline
version: 0.1.0
source: self
---

# R1 C4 · 协议继承与模块候选三态矩阵

## 固定输入

本矩阵只读取 Q2 路径 [I-PROTO-001 v0.1.3 覆盖表](../../../../docs/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)。其 `freeze_status` 为 `frozen`，来源固定为 `schema-ui-docs v2.7.0` commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`。本记录不把 workspace-001 的过程、目标或审计状态作为事实。

| disposition | 固定 domain |
|--------------|-------------|
| `include` | `D-NODE`, `D-EXPR`, `D-DATA`, `D-PERM`, `D-APP`, `D-VER`, `D-VAL` |
| `include-partial` | `D-COMP`, `D-ACT`, `D-TABLE`, `D-FORM` |
| `exclude` | `D-UPLOAD` |

## R1 候选模块映射

以下 disposition 是协议域对候选模块的继承范围，不是“当前模块已经实现”或“每个 Profile 必须启用全部候选”的证明。`mvp`/`admin` 的精确集合和 precedence 仍属于 I-004/R2。

| 候选模块 | Profile 候选 | include | include-partial | exclude/边界 | 当前依据 |
|----------|--------------|---------|-----------------|--------------|----------|
| `core.server-registration` | `mvp`/`admin` required candidate | `D-APP`, `D-NODE`, `D-VER` | — | 不承诺动态远程 Manifest | `apps/api/internal/handler/health.go:27-37`; `I-PROTO-001:39,45,49` |
| `core.auth-session` | `mvp`/`admin` required candidate | `D-PERM`, `D-VER` | — | 不承诺完整 IAM | `apps/api/internal/handler/resources.go:185-199`; `I-PROTO-001:44,49` |
| `core.manifest-route` | `mvp`/`admin` required candidate | `D-APP`, `D-NODE`, `D-VER`, `D-VAL` | — | host 精确接受协议 `2.7` | `apps/web/src/protocol/app-manifest.ts:556-617`; `I-PROTO-001:39,45,49-50` |
| `core.navigation-capability` | `mvp`/`admin` required candidate | `D-EXPR`, `D-PERM`, `D-APP`, `D-VER` | — | visible-when 受固定表达式子集约束 | `apps/web/src/app/navigation.ts:137-152`; `apps/web/src/renderer/permissions.ts:94-110`; `I-PROTO-001:40,44-45,49` |
| `core.schema-render` | `mvp`/`admin` required candidate | `D-NODE`, `D-EXPR`, `D-VAL` | `D-COMP`, `D-FORM` | 不纳入完整 registry、upload type | `apps/api/internal/handler/schema.go:24-68`; `I-PROTO-001:39-41,47-50` |
| `admin.users` | `mvp` optional / `admin` required candidate | `D-DATA`, `D-PERM` | `D-COMP`, `D-TABLE`, `D-FORM` | 仅排序/搜索等固定列表子集；无批量 action | `apps/api/internal/handler/health.go:32`; `I-PROTO-001:42-47` |
| `admin.roles` | `mvp` optional / `admin` required candidate | `D-DATA`, `D-PERM` | `D-COMP`, `D-TABLE`, `D-FORM` | 同上；Profile 精确集合未冻结 | `apps/api/internal/handler/health.go:33`; `I-PROTO-001:42-47` |
| `admin.settings` | `mvp` optional / `admin` required candidate | `D-DATA`, `D-PERM`, `D-APP` | `D-COMP`, `D-FORM` | 不把 Shell 私有通知纳入 | `apps/api/internal/handler/health.go:35`; `I-PROTO-001:42-47` |
| `admin.activity` | `mvp` optional / `admin` required candidate | `D-DATA`, `D-PERM` | `D-ACT` | 仅非批量 action；Activity UI 与 operationlog 横切记录分离 | `apps/api/internal/handler/health.go:34`; `docs/architecture/module-architecture.md:103-107`; `I-PROTO-001:43-44` |
| `shell.registry` | target candidate | — | — | 当前 absent；`D-UPLOAD` 继续 exclude | C1 evidence；`I-PROTO-001:48,78-80` |

## Domain 与 fixture 门槛

| domain/suite | R1 inherited disposition | 固定边界 |
|--------------|--------------------------|----------|
| `D-NODE` / structural schemas | include | Node/Page 合法树 |
| `D-EXPR` / `reactions` | include | `$context` visible/disabled 子集；multi-round `$deps` 不纳入当前 MVP |
| `D-DATA` / request-construction、response-mapping、query-serialization、static-data | include | 列表/详情 datasource 与序列化 |
| `D-PERM` / permissions-inheritance | include | 会话、权限继承和后端授权 |
| `D-APP` / app-manifest、app-navigation | include | Manifest、导航壳、路由入口 |
| `D-VER` / version-negotiation、runtime-defaults | include | 生产 host 当前精确接受 `2.7` |
| `D-VAL` / validation schemas | include | 加载/构建结构校验 |
| `D-COMP` / component-format | include-partial | 固定白名单/格式子集，不是完整 registry |
| `D-ACT` / actions、request-lifecycle | include-partial | 仅非批量 action/request |
| `D-TABLE` / table-sort、search-table | include-partial | 排序/搜索表，排除多选批量 |
| `D-FORM` / form/control schemas | include-partial | 基础表单、固定控件与校验回填，排除 upload/未列 type |
| `D-UPLOAD` / uploads | exclude | UI、端点、fixtures 整域排除 |

`scenarios` 只是 `support-only` 信息性场景，不新增三态 domain，也不构成独立 conformance gate。

## 范围变更门槛

任何新增 domain、扩大 `include-partial` 子集、改变 `D-UPLOAD` 排除、引入新上游协议版本，都必须：

1. 新增用户确认的决策记录；
2. 递增 I-PROTO 覆盖表版本并保留旧版追踪；
3. 在受影响的 `/govern` required 信息门禁前完成验证；
4. 更新本工作区的协议矩阵和相关审计证据。

当前没有范围扩张决定；本矩阵沿用 v0.1.3，不修改固定表。

## 检查点结论

C4 的候选模块三态矩阵、fixture 边界和范围升版门槛已形成。Root I-007 仍为 `open`，本记录只证明继承范围可追踪，不证明协议实现或 R1 阶段已通过。
