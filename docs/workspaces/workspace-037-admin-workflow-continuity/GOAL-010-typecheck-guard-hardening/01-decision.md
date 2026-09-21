---
id: GOAL-010-typecheck-guard-hardening-decisions
doc: decision
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
---

# 决策台账 · GOAL-010

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-18 | 守卫加固范围与 `-p` 目标判定规则 | accepted | [D-001-guard-hardening-scope-and-rule.md](01-decision/D-001-guard-hardening-scope-and-rule.md) |

## 当前事实

- 2026-09-18，用户指示开设本目标（整改 `GOAL-008 A-002 F-001`），并在自审后以本地 grok build（grok 4.6 · 思考强度 xhigh）执行独立审计，确认无问题后关门。
- 判定规则：`tsc` 调用只有在带 `-b`/`--build`，或 `-p`/`--project` 指向**自身选择源文件**的配置时才计为检查型；否则（含裸调用、`-p` 指向 solution-style 根配置、目标不存在/不可解析）为违规。详见 `D-001`。
