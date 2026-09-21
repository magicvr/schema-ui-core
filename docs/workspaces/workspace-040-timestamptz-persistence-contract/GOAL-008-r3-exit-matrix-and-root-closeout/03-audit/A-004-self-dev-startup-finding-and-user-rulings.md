---
id: A-004-self-dev-startup-finding-and-user-rulings
doc: audit-entry
status: active
parent: GOAL-008-r3-exit-matrix-and-root-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# A-004 · 用户报告缺陷的处置（`F-I-101` fixed）+ 用户裁决留痕

- **source**: `self`（编排器；含用户 2026-09-21 的书面裁决留痕）
- **日期**: 2026-09-21
- **scope**: 用户报告的 `.\dev.cmd start` 启动失败（复现、根因、修复、端到端验证）；用户对 `I-041-011` 与 GOAL-004 `user-overruled` 载体的裁决；Root 关门确认的重新提请
- **verdict**: `pass`（新 finding `F-I-101` 已 `fixed`；开放 required = 0；Root 关门仍待用户确认）

## 1. 用户裁决留痕（2026-09-21，P-004）

| 事项 | 用户裁决 | 处置 |
|------|----------|------|
| `I-041-010`（Root 关门确认） | **未选择任一选项**，改为报告新问题：`.\dev.cmd start` 启动失败，**需先修复** | 记为**关门前置条件**：Root 保持 `active`；缺陷修复并验证后再提请确认（本条目 §4） |
| `I-041-011`（产品级支持组合表述） | **只保留在退出矩阵与 R3-C 附件**（不写入发布说明） | `I-041-011` → **verified（用户裁决）**：退出矩阵保留；不产生产品发布说明表述 |
| GOAL-004 `F-I-002`（`*time.Time` + JSON `null`）| **授权补一条 `D-` 条目** | 保留结论不变，把「用户书面依据」从审计响应自述提升为 `GOAL-004/01-decision/D-002-*` 正式决策条目；`A-003` 的修订段保留 |

## 2. 新 finding `F-I-101`（self 记录，源自用户报告）

| 项 | 内容 |
|----|------|
| **级别** | **required**（阻塞 Root 关门：用户明确要求先修复；且它使 PG 方言下的**真实入口**无法启动） |
| **问题** | `composition.recoveryArtifactsDir` 对相对 `db.path` 返回**相对**目录（`data\recovery`），被交给 `backup.PgProvider.WorkDir` 并用于 `docker run -v`；docker 拒绝相对源路径（exit 125），C3 class-A rollback 产物创建失败 → `openStore` 失败 → API 启动失败 |
| **证据** | `E-003`（复现日志、根因链、覆盖缺口分析、端到端验证） |
| **闭合** | **`fixed`**：`recoveryArtifactsDir` 全分支经 `absoluteArtifactDir`（`filepath.Abs`）解析为绝对路径；新增 `TestRecoveryArtifactsDirIsAbsolute`（5 例）与 `TestRecoveryWiringHandsPgProviderAnAbsoluteWorkDir`；**端到端实测**：一次性库上 `healthz 200`（~5s）、`readyz 200`（~10s）、`data/recovery/*.artifact` 落盘 101,787 B；清理后工作树仅含修复文件 |
| **为什么此前未发现** | 既有 composition PG 启动测试用 `t.TempDir()`（绝对路径），与真实配置的**相对 `db.path`** 不一致；`A-002` 也只在测试夹具证据上复核。**已把该覆盖缺口补进退出矩阵的限定**（见 §3） |

## 3. 对退出矩阵的修订（不重开已落盘结论）

- **判据 3/4 的限定加强**：矩阵 §3 的 `L-1`（真实路径绑定本环境）追加一条：**测试夹具使用绝对临时路径，不能代表「配置里为相对 `db.path`」的真实入口**；`F-I-101` 证明该差异会掩盖启动级缺陷。修订后判据 3/4 的结论不变（修复后真实入口实测通过），但限定更完整。
- **判据 5 不受影响**（无新依赖、无驱动类型泄漏）。
- 修复本身不触碰任何冻结物（迁移/checksum/codec/wire 合同）。

## 4. Root 关门确认（重新提请）

`F-I-101` 修复并端到端验证后，原先的确认包仍然成立，另加一项**环境问题**需用户裁决（属用户环境，不代裁）：

| # | 事项 | 编排器建议 |
|--:|------|------------|
| 1 | `I-041-010` Root 关门确认 | 建议确认（判据 1–5 有可复跑证据；开放 required = 0；`F-I-101` 已 `fixed` 并实测） |
| 2 | `.env` 使 dev 启动器指向**共享常驻实例的 `postgres` 库**（`DB_NAME=postgres` 覆盖 `config.yaml` 的 `schema_ui`） | 建议为 dev 指定**专用库**（例如 `schema_ui_dev`），避免 `dev.cmd start` 在共享维护库上跑完整迁移链；是否改动属用户环境决策 |

## 5. 边界

- 本条目**不**把 Root 标 `done`：`I-041-010` 仍为 open required，等待用户对 §4 的书面确认。
- `I-041-011` 已按用户裁决闭合；GOAL-004 的 `D-` 条目按用户授权补写（结论不变）。
