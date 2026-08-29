---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-go-library-consumption
version: 0.1.0
---

# E-003 · S3 泄漏面实测与依赖链盘点（2026-08-29）

## 实验序列

| # | 实验 | 结果 | 结论 |
|---|------|------|------|
| 1 | golden-consumer `import internal/kernel`（E-001，S1） | `use of internal package ... not allowed` | internal 规则阻断**目录树外**导入（当时 kernel 在 internal/ 下） |
| 2 | 外移后（E-002）：`import apps/api/kernel` + `apps/api/modules/dashboard` + `go run` | exit 0 · `kernel=2.0.0 … modules=23` | 升级目录树后 A/B 层可外部消费 |
| 3 | **本条目**：`import apps/api/modules/users`（其依赖链含仍留在 `internal/` 的 `auth`/`handler`） | **exit 0**（go.sum 补齐后） | Go internal 规则以**目录树**为界：`apps/api/...` 前缀内包相互可引用——B 层包引用内部实现**合法** |
| 4 | 本条目：尝试以 `users.New` 签名构造（需 `*auth.Authenticator`） | 推断：下游**无法 import `internal/auth`**，类型名不可命名 → 直接构造签名不可完成 | 泄漏形态 = **类型命名面泄漏**（非包加载失败） |

## 依赖链盘点（模块 → C 层）

- 模块对 `internal/` 的 import：**handler ×18**（路由工厂，模块 Register 内部实现细节——下游不需要命名）、**auth ×17**（`*auth.Authenticator` 进入 `users.New`/`authsession` 构造签名——**下游唯一命名需求**）、jobs ×1。
- `internal/store` / `config` / `mail` / `objectstore` 无模块直接 import：store 经 `authsession.TxRunner`（B 层接口）间接；mail 经 `kernel.MailSender`（A 层接口）间接。
- `auth.New(secret, ttl, ttl, runner authsession.TxRunner, dev)`——**参数全部可构造/可命名**（[]byte、time.Duration、B 层接口、bool）；`runner` 的非 nil 实现来自 `internal/store`（C 层方言）。

## 关键机制发现（决定收敛方案）

**Go 类型推断允许「不命名」消费**：`a := assembly.NewAuthenticator(cfg)` 后可直接 `users.New(a, …)`——接收方无需 import/命名 `*auth.Authenticator`。因此：

- 新增**公开装配工厂包**（如 `apps/api/assembly`，同模块树、可 import `internal/*`），导出以 config/kernel/B 层接口为输入的工厂 → **零契约修订、零外移**即可让下游完成可运行装配。
- 该包为**新增面**（冻结面扩展为 B+ 层，additive，minor 合规，非 breaking）。

## 关联

- F-001（C 层泄漏）→ 收敛方案三候选已在 S3 验证矩阵实测：① 上移契约（kernel 接口化，长期）② **公开装配工厂（推荐，见 D-003 提案）** ③ 有限白名单外移（auth+store 外移，备选）。
- F-002（B 层符号回填）：`attachments/modules-export-inventory-v0.1.md`（22 包导出扫描）已生成。