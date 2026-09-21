---
id: GOAL-047-w35-dev-db-init-and-bootstrap
doc: decision
status: active
parent: null
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# 决策记录 · GOAL-047-w35-dev-db-init-and-bootstrap（W35）

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 / 决策 |
|----|------|----------|----------|------|-------------|
| I-047-001 | non-blocking | 初始化快捷方式的 CLI 形态（子命令名、附加库名、多方言参数） | A | **verified（2026-09-21）** | `D-001`：`dev.cmd init-db` → `cmd/dbsetup`，默认 dev+test，支持 `--dev-only`/`--test-only`/附加库名 |
| I-047-002 | non-blocking | e2e 专用库是否并入同一快捷方式 | B | **verified（本目标边界）** | `D-001`：e2e 生命周期独立；`e2e-pgset` 复用自举层但不并入默认 dev/test 集 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-21 | 初始化快捷方式 CLI 形态与自举边界 | accepted | `01-decision/D-001-cli-shape-and-init-boundary.md` |

> 新决策从 `01-decision/D-NNN-<slug>.md` 写入；编号在本目标内单调不复用。

## 已继承的边界事实（不得在本目标重开）

- **应用不自动建库**：`openStore` 只 `Ping` + 应用 catalog；自动 `CREATE DATABASE` 被明确列为非目标（会掩盖拼错的库名/连错实例，且与 `W21` 启动身份判定取向冲突）。
- **VP-040 已关门**：本目标不触碰迁移链、canonical SQL/checksum、codec、wire 合同、备份 Port。
- **`.env` 语义不变**：初始化只读连接参数，不写秘密、不改键名。
- **结构裁决**：用户 2026-09-21 选择「在 workspace-010 新加子目标」（不新开 VP）；slug 经用户确认。
