---
id: A-006-response-to-a005-and-root-closeout-reask
doc: audit-entry
status: active
parent: GOAL-008-r3-exit-matrix-and-root-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# A-006 · 编排器响应（A-005 `pass`）+ 重新提请 Root 关门确认

- **source**: `self`（编排器响应）
- **日期**: 2026-09-21
- **scope**: 响应 `A-005`（independent · grok build · grok-4.6 · high，定向复审 `F-I-101`）的全部内容；重新提请 `I-041-010`
- **verdict**: `pass`

## 1. 对 A-005 判定的接受

A-005 判定 **`F-I-101`（required）已 `fixed`**，并**独立复现**了：旧函数在真 docker 上的相对挂载失败（exit 125 / 原文报错）、新函数绝对化后挂载成功、以及**一次性库 `vp040_a005` 上的端到端**（healthz 200 / readyz 200 / class-B 产物 101,764 B），随后清理并确认工作树干净。它同时确认修复**未触碰**任何冻结物（v1–v87 canonical SQL/checksum、codec、wire 合同、C3 Port 形状、`PgProvider` 命令形状），并把「同类相对路径隐患」逐处扫过（13 处），结论是本缺陷是 C3 PG provider 与相对 `db.path` 的交点，**不是**家族性漏网。

**开放 required = 0**；A-005 **同意**再次把 Root 交给用户确认关门，并**不同意**自行把 Root 标 `done`。编排器接受。

## 2. 逐条响应

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| `F-I-101`（相对 `db.path` → `docker run -v` 启动失败） | required | **维持 `fixed`**（A-005 已独立复现并同意） | 代码 + 两条回归 + 端到端；`E-003`、`A-004`、`A-005` §1 |
| `F-I-102`（`PgProvider` 不强制绝对 `WorkDir`；`Abs` 失败回退相对） | recommended | **fixed（纵深防御）**：① `PgProvider.Create`/`Restore` 在非绝对 `WorkDir` 上 **fail closed**（`KindInvalidRequest`，不再触达 docker）；② `absoluteArtifactDir` 在 `filepath.Abs` 失败时返回 **空**（锚点禁用、由 store 记录），不再把相对路径交出去；③ 函数注释改为「绝对或空」，与实现一致。新增 `TestPgProviderRejectsRelativeWorkDir`（3 例：`./data/recovery`、`data/recovery`、`.`；断言分类为 `InvalidRequest` 且**未调用** docker） | `provider_pg.go`、`provider_pg_test.go`、`composition.go` |
| `F-I-103`（`00-meta` 信息表仍写 `I-041-011` = open） | recommended | **fixed**：`GOAL-008/00-meta.md` 信息表该行改为 **`verified`**（用户 2026-09-21 裁决：只留退出矩阵/R3-C 附件），与 `01-decision.md` 和 `A-004` 一致 | `GOAL-008/00-meta.md` |

## 3. 检查点状态

- 检查点 A/B 已完成（`progress: 2/3`）；`A-001`～`A-006` 全部落盘；跨目标开放 required = 0。
- 唯一未满足的门禁是**判据 6 / `I-041-010`：用户确认关门**（设计如此）。

## 4. 重新提请的用户裁决（2026-09-21）

| # | 事项 | 级别 | 编排器建议 | 备注 |
|--:|------|------|------------|------|
| 1 | `I-041-010`：是否接受退出判据矩阵 + 独立关门审计（`A-002`、`A-005`）与残留清账后**关闭 Root** | **required** | **建议确认**：判据 1–5 有可复跑证据；开放 required = 0；用户报告的 `dev.cmd start` 缺陷已独立复现并闭合（`F-I-101` `fixed`） | 确认后编排器执行：`GOAL-008` → `done · 3/3`、Root → `done · 3/3`、勾选六条成功标准、同步 `goal-tree`/`docs/vision`/`workspace.md` 投影并提交 |
| 2 | `.env` 使 dev 启动器指向**共享常驻实例的 `postgres` 库**（`DB_NAME=postgres` 覆盖 `config.yaml` 的 `schema_ui`） | non-blocking（环境） | **建议为 dev 指定专用库**（如 `schema_ui_dev`）：当前配置下每次 `dev.cmd start` 都会在共享维护库上跑完整迁移链 | 属用户环境决策；不改也不阻断关门，但改造前请注意共享库会被迁移 |

## 5. 边界

- 本条目**不**把 Root 标 `done`：`I-041-010` 仍 open，等待用户书面确认。
- 本轮修复只涉及 `internal/backup/provider_pg.go`、`internal/composition/composition.go` 与其测试；未触碰任何冻结物。
