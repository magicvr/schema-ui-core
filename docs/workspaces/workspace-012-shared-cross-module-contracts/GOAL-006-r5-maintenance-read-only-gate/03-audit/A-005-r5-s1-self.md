---
id: A-005-r5-s1-self
goal: GOAL-006-r5-maintenance-read-only-gate
source: self
verdict: pass
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# A-005 · R5 S1 实现自审

## 审计头

| 项 | 值 |
|----|----|
| source | self |
| scope | S1 runtime config、bootstrap/status projection、正常态兼容与非回归 |
| verdict | pass |
| required findings | 0 |

## 核对结论

1. `Config.Load` 明确区分 YAML 缺省 normal 与显式空/未知 `RUNTIME_MODE`，env 覆盖优先级有定向测试；既有零值 Config 测试保持兼容且不削弱 Load fail-closed。
2. bootstrap 的 read-only/degraded 投影不使用 capability narrowing；正常态 bytes/字段保持兼容，非法 runtime mode 在注册时 fail closed。
3. system-monitoring 追加字段而不改既有 status/readiness 字段，生产 provider 显式传入 runtime mode，旧测试构造器保持默认 normal。
4. 组合装配、system-monitoring、config、handler 全量测试和 docscheck 通过；未修改 Profile 默认集、模块闭包、Manifest 聚合算法、协议 pin 或 readiness 探针职责。

## Findings

| ID | 等级 | finding | disposition |
|----|------|---------|-------------|
| F-001 | recommended | S2 需用 core 与 Provider 各一条真实写路径验证门禁和 allowlist，且验证未知路径 404/405 与 request-id。 | implementation gate：S2 |
| F-002 | recommended | S3 需做 Host/客户端对 `SERVICE_*` code 的消费核对，确保不按 503 误判为 Host unavailable 或 ACCOUNT_LOCKED。 | implementation gate：S3 |

## 结论

S1 self 通过，开放 required = 0；S2 获准实施，A-005 的 recommended findings 保留到后续阶段验证。
