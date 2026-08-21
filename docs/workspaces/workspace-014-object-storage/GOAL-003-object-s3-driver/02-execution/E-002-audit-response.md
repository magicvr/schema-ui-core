---
id: E-002-audit-response
title: A-002 独立审计响应——R-001/R-002/R-003 同批闭合
status: recorded
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-003-object-s3-driver
version: 0.1.0
---

# E-002 · 审计响应实施

独立审计 A-002 verdict **pass**（开放 required 0）。recommended 项同批处理：

1. **R-001（fixed）**：main.go 启动警告改为"readyz 已覆盖后端探针；三类落盘仍走本地至 R3 接线"，消除与现状矛盾。
2. **R-002（fixed）**：`go mod tidy`——直接依赖归位、移除未使用默认链模块（config/imds/sso/sts/signin）。
3. **R-003（fixed）**：
   - fake 增加 put/get/head/delete/list 传输错误注入；新增 `TestS3TransportErrorsPropagate`（Put/Delete/List 传播、Exists 非 404 不误报 false）。
   - composition 探针构造提取为 `newObjectProbe(cfg)` 并新增 `TestNewObjectProbe`（local/零值→nil 探针；driver=s3→非 nil 且对不可达端点快速失败），锁住 driver=s3 接线。
4. **N-001（fixed）**：D-001 §1 文本同步为实际的手写 aws.Config 构造。
5. **N-003（fixed）**：本条即补记——全量 `go test ./...` exit 0（commit 1545134 后复跑）；live MinIO 证据按计划归 R5。
6. **N-002 / N-004**：留痕——stub token 同形问题无影响；R3 调用方不得把 SDK 原文错误直接打日志（R3 立项时引用）。

## 验证

go build exit 0；go test ./internal/{objectstore,composition}/ 全绿；全量套件见 E-001 补记与本条。
## 验证

go build exit 0；objectstore/composition 测试全绿；全量 go test ./... exit 0。

> **修正记录（诚实留痕）**：本条初版在 stray-brace 修复前写了"全量绿"——实际全量首跑 FAIL（s3_test.go 结构体多一个右括号，commit a4c68ef 带入）。修复后复跑全量 exit 0（本修正所在提交）。教训：全量验证必须在提交前完成。
