---
id: GOAL-003-dual-key-jwt
doc: decision
status: done
parent: GOAL-001-key-rotation-and-backup
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
---

# 决策记录 · GOAL-003

## 信息需求镜像（以 Root 00-meta 为权威）

| ID | 级别 | 状态 | 证据 |
|----|------|------|------|
| I-003 | required | **verified** | 本目标 D-001（重叠窗 = 配置存续期；不用 kid；refresh opaque 不受影响） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-22 | R2 双密钥语义冻结（关闭 I-003）与实施方案 | accepted | [D-001-dual-key-semantics.md](01-decision/D-001-dual-key-semantics.md) |
