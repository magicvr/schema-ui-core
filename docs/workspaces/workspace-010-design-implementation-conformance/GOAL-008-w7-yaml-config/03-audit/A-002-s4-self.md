---
id: A-002
goal: GOAL-008-w7-yaml-config
source: self
date: 2026-08-14
scope: S2/S3 实施与验证
verdict: pass
parent: GOAL-008-w7-yaml-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-002 · self 审计（S2/S3）

## 结论

**verdict: pass**（E-003/E-004）。

## 核对

- S2 交付物齐全：configs/config.yaml、config.default.yaml（embed）、configs/.env.example、Load 分层、上传三字段迁入、RegisterUpload 变参、compose 注释同步。
- 优先级链（进程 env > CONFIG_FILE > embed > 字段默认）与 D-002 冻结一致；CONFIG_ENV_FILE 不覆盖进程 env 有单测。
- 插值行级作用域：注释内 ${...} 不算引用（有测试覆盖文档式示例）。
- fail-closed 三例均有测试：显式 CONFIG_FILE 缺失、裸 ${VAR}、未知 YAML 键（KnownFields）。
- 向后兼容：旧 env-only 部署经 env 覆盖路径零迁移（embed 回退实测 + 既有测试全绿）。
- 回归：apps/api go test ./... 全绿；web 无代码变更。

## Findings

- 无 required。
- 备注（非必改）：configs/config.yaml 与 embed 默认存在双份内容（operator 模板 vs 内置）。已用头部注释说明复制关系；未来若默认值演进，两处需同步（S5 关门审计可复核）。
