---
id: GOAL-002-s0-denominator-freeze
doc: execution
status: active
parent: GOAL-001-admin-module-readiness
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
---

# 执行记录 · GOAL-002

## 执行条目索引

| E-ID | 日期 | 动作 | 结果 | 证据 / 文件 |
|------|------|------|------|-------------|
| E-001 | 2026-08-10 | 完成 S0 分母证据盘点（代码/环境/Profile/DB/命令矩阵/模块名册/协议面）；首轮基线：API build/test、Web build/vitest 实测通过 | pass | 候选 commit `852ee7e`（clean） |
| E-002 | 2026-08-10 | 首轮基线补测：`go vet ./...`、smoke（mvp + admin，SM-001~005+SM-007）、compose disposable（SM-006） | pass（见证据） | V-003/V-007/V-008；详见 Root [D-003](../GOAL-001-admin-module-readiness/01-decision/D-003-s0-denominator-freeze.md) 验证命令矩阵 |

## 首轮基线结果（候选 `852ee7e`，2026-08-10）

| 命令 | 结果 | 备注 |
|------|------|------|
| `go build ./...` | ✅ exit 0 | — |
| `go test ./...` | ✅ 全包 pass | 含 cmd/server 重启、store 迁移、kernel、composition |
| `go vet ./...` | ✅ exit 0 | — |
| `npm test`（vitest） | ✅ 40 files / 728 tests | 含协议 conformance、i18n、D-PERM |
| `npm run build` | ✅ vite 6.4.3 · 1834 modules · 3.61s | — |
| `smoke.sh`（mvp，SM-001~005+007） | ✅ SM-007 PASS · exit 8 | 非 disposable 部分绿，符合预期 |
| `smoke.sh`（admin，SM-001~005+007） | ✅ SM-007 PASS · exit 8 | 12 页 manifest 校验通过 |
| `smoke.sh --disposable`（隔离 compose，SM-006） | ✅ exit 0（SM-001~006 完整绿） | V-008；`ci-smoke-s0` 隔离 project，重启后种子断言通过 |
| `npm run test:e2e`（mvp） | ✅ 3 pass + 1 skip（admin-only 用例） | V-006 |
| `npm run test:e2e`（admin） | ✅ 3 pass + 1 skip（mvp-only 用例） | V-006 |

## 记录规则

只写已发生事实；命令、结果、commit/digest 与 Q2 证据路径一并记录。基线证据绑定字段按 Root `I-READINESS-007` 冻结。
