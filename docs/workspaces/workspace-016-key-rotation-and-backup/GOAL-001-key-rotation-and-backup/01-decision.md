---
id: GOAL-001-key-rotation-and-backup
doc: decision
status: active
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

> 状态以 `00-meta.md` 信息表为准（本表为镜像，须保持同号同状态）。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | current/previous 键名、fail-closed、熵 | R1 方案 | R1 冻结 | R1 决策 | **verified**（D-002） | — | D-002 §1；证据 `config.go:202/496/953-964`、`main.go:74-85` |
| I-002 | required | 密钥集合是否仅 JWT | R1 方案 | R1 冻结 | R1 决策 | **verified**（D-002） | — | D-002 §2；服务凭证 SHA-256 opaque 不共用（`auth.go:411-420`） |
| I-003 | required | 重叠窗 / kid / refresh | R2 方案 | R2 接入前 | R2 决策 | **verified**（GOAL-003 D-001） | — | 重叠窗 = previous 配置存续期；不用 kid；refresh opaque 不受影响 |
| I-004 | required | 轮换后恢复剧本 | R3 方案 | R3 接入前 | R3 决策 | collecting | — | VP I-016-004 |
| I-005 | non-blocking | 旧 access 立即失效是否残余 | 退出 1 | R2 | 用户书面残余 | collecting | — | VP I-016-005 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-22 | 开区 scaffold 与 A5 纲领路线图 | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
| D-002 | 2026-08-22 | R1 轮换合同与配置面冻结（关闭 I-001/I-002） | accepted | [D-002-rotation-contract-freeze.md](01-decision/D-002-rotation-contract-freeze.md) |
