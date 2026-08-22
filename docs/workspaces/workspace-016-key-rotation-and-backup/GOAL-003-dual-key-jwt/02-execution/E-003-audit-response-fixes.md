---
id: E-003
doc: execution-entry
goal: GOAL-003-dual-key-jwt
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-003 · 审计响应落地（A-002 F-001/F-002/F-003）

## 事实（2026-08-22）

对 A-002（independent pass · 0 required · 3 recommended）的修正记录：

| Finding | 响应 | 证据 |
|---------|------|------|
| F-001 composition 双密钥接线缺自动化钉死 | **fixed**：新增 `TestNewAuthenticatorWiresPreviousSecret`（composition 包）——生产构造路径 `newAuthenticator(cfg, jwtSecret, repository)` 在 `cfg.AuthJWTSecretPrevious` 非空时中间件接受 old-key token；previous 为空时同一路径拒绝 old-key token | `internal/composition/composition_test.go`；`-run TestNewAuthenticatorWiresPreviousSecret -count=1 -v` PASS |
| F-002 执行索引未登记 E-002 | **fixed**：E-002 已补入 `02-execution.md` 索引（连同本条 E-003） | 本文件与索引表 |
| F-003 「整包 exit 0」不可作为已复现事实 | **fixed（措辞收窄）**：E-002 v1.1 改为「vet 0 + JWT 相关包 ok（双方独立复现）」，并注明 store PG 两条集成失败为共享 probe DB 残留（`WasFresh` 对任意用户表敏感），非本切片回归；附 R3 使用专用 DB 的遗留注记 | E-002 v1.1 |

三条均为 recommended（非必改），全部按 fixed 路径闭合；无开放 required finding。
