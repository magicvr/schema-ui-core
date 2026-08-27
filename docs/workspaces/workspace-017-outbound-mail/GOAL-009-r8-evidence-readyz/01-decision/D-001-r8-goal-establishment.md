---
id: D-001
doc: decision-entry
goal: GOAL-009-r8-evidence-readyz
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-001 · 开设 R8 子目标

## 背景

R7 完成后仅剩 R8（收官证据）。P-001：按阶段开设。

## 决定

1. 创建 `GOAL-009-r8-evidence-readyz`，parent = `GOAL-001-outbound-mail`，status `active`。
2. 探针语义（实现级，沿既有先例）：Resend 探针 = 对所配 endpoint 的轻量可用性检查（GET /domains，Bearer 鉴权，短超时），仅在 boot 显式配置 resend 时挂入 readyz——镜像「仅显式配置后 readyz 扩依赖」的 R4 冻结语义；mock/未配置 probe 恒为 nil。热切换后的探针不随运行时渠道变化（boot 冻结），与 R7 文档声明一致。
3. live 投递测试沿 `smtp_live_test.go` 先例做 env-gated skip 缝；是否补跑 live 由用户在关门问询时裁决（判据 3 已授权等价 harness）。
4. Root 再关门审计模式（self / independent）与愿景层收口由编排器在本目标关门后推进；涉及 P-004 问询处另行问询。
5. 审计模式：scaffold none → 关门 self。

## 未选方案

- 把 Root 再关门审计并进本目标五件套：审计层级混淆（Root 级关门向审计归 Root 台账）。
- 强制 live 投递作为硬门禁：违反 VP 判据 3 的「或等价 harness」已冻结文本。
