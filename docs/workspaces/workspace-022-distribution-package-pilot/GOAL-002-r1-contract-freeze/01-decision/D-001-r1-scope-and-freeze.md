---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-contract-freeze
version: 0.1.0
---

# D-001 · R1 范围与冻结面分界决策

## 决策

1. **冻结面三分层**（Go 侧）：
   - **A 层 · 内核冻结面** = `internal/kernel` 全部导出符号（11 非测试文件）。变更受 semver 约束；breaking 必须 major + changelog 迁移说明；`KernelAPIVersion = "2.0.0"`（module.go）为版本锚点，registry 兼容窗校验（`KernelAPIRange`）已存在并继续作为 fail-closed 门禁。
   - **B 层 · 模块契约面** = `internal/modules/*` 向组合根暴露的装配面：`ModuleID` / `Provider` / `New*(...)` 构造 / `Descriptor()` / `CompiledPersistence()` / `Register(ctx, kernel.Registrar)`。下游组合根 import 引用；breaking 与 A 层同纪律（先 deprecate 再删）。
   - **C 层 · 内部实现面** = composition / handler / server / auth / config / jobs / mail / objectstore / obs / manifest / migration / store 方言（kernel 端口之外的全部内部包）。不冻结；**默认禁止下游 import**（除清单 §5 白名单张力项）。
2. **清单草案**：`attachments/freeze-face-v0.1.0.md`；**待用户确认后 v1.0.0 生效**（关键决策：冻结面直接决定契约稳定性税与 R2 打包范围）。
3. **不进本目标**：G2 多模块细版本、CLI、npm 拆包实施（I-002 属 R3）、发布通道落地（I-003 属 R5）。`pkg/version` 为构建注入面不冻结。
4. **审计**：S1–S3 对账 = self（A-001）；S4 用户书面确认关门；R5 发布/兼容门禁与 Root 关门 = independent（grok build）。

## 边界张力（A-001 F-001，随 R2 闭合）

模块构造签名引用 C 层内部类型（例：`users.New(*auth.Authenticator, *authsession.Repository, ...)`）——下游组合根若直接装配将被迫 import `internal/auth`。候选收敛路径（R2 试点时定）：① 将 `auth.Authenticator` 等上移为 A 层契约类型；② 提供公开构造工厂/适配器包（如 `foundation/assembly`）；③ 有限 C 层白名单。**不预先裁定，留给 R2 证据。**