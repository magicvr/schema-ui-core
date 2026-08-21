---
id: E-008-r4-s4-validation
goal: GOAL-005-r4-async-job-contract
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# E-008 · R4 S4 全量验证与错误码契约修复

第一次 `go test -timeout 15m ./...` 仅在 `internal/handler` 的 error-contract drift guard 失败：R4 D-002 §7 已冻结并实现的五个 HTTP Job 码与两个持久化终态码尚未进入跨工作区权威 Appendix A / pinned set。

checkpoint `425215a` 按原 Appendix A“可新增、不可改语义/复用”规则追加七个稳定码；测试区分 handler 实际发出的 HTTP 码与只存在于 Job representation 的 stored terminal 码。定向 error-contract 测试与 `go test ./internal/docscheck` 通过。

修复后第二次执行 `go test -timeout 15m ./...`：全部 package PASS；其中 `internal/handler` 完整回归 217.133s，`internal/store`、`internal/jobs`、`internal/composition`、wallet 与 operationlog 包均通过。

既有 S3 高风险证据仍有效：Job/wallet race PASS、count=10 PASS、migration 42→43 correlation preservation PASS、`git diff --check` 与未跟踪文件 trailing-whitespace 检查 PASS。
