---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-contract-freeze
version: 0.1.0
---

# A-001 · S1–S3 对账自审（source: self · 2026-08-29）

## scope

GOAL-002 S1 扫描事实、S2 成文（冻结面清单 v0.1.0 / semver 流程 / changelog 模板）、S3 对账：清单 vs 实际代码导出面。

## verdict

**conditional**（证据核对成立；1 条 recommended 随 R2 闭合）

## 核对点

| # | 核对项 | 方法 | 结论 |
|---|--------|------|------|
| 1 | 清单 §1 声称 11 个 kernel 非测试文件 | 目录列举 vs 清单 | ✅ 完全一致（11/11） |
| 2 | 导出锚点抽样：`KernelAPIVersion` / `Store` / `Provider` / `RegisterContributions` / `MigrationChecksum` / `MailSender` / `ObjectStore` 存在 | grep 导出声明 | ✅ 均命中 |
| 3 | B 层模块模式（ModuleID/New/Descriptor/Register 六贡献） | users/provider.go 全文 | ✅ 吻合 |
| 4 | 组合根静态导入面可作为下游组合根模板 | composition.go imports | ✅ 吻合 |
| 5 | semver 流程与 changelog 模板可执行（含 breaking 迁移节） | 文档核对 | ✅ |

## findings

- **F-001（recommended）**：模块 `New*` 构造签名引用 C 层类型（`*auth.Authenticator` 等）+ B 层 store 附件包——包化装配面存在「C 层泄漏」。**关闭路径**：R2 打包试点用真实下游组合根验证，按 D-001 §4 三候选（上移契约 / 公开工厂 / 有限白名单）择一落地后在清单 v1.0.0 定稿。不阻断 R1 关门（Q2 证据不足已登记）。
- **F-002（recommended）**：B 层符号全量清单（逐包导出）尚未枚举，目前只有规则 + 抽样。**关闭路径**：R2 打包时以 AST/导出扫描补全，回填清单 §4。

## 响应

- F-001/F-002：采纳为 R2 前置项（GOAL-003 S1 纳入）；本条目随 R2 关门时以 fixed 回填。

## 响应回填（2026-08-29 · R2 关门）

- **F-001 → `fixed`**：C 层泄漏收敛 = 用户裁决方案 β（GOAL-003 D-003）→ `apps/api/assembly` 公开装配工厂实证（GOAL-003 E-004：users 全链零 internal 命名装配 + 贡献计数 = Descriptor 声明；A-002 核对点 #5）。
- **F-002 → `fixed`**：B 层符号全量盘点 = `GOAL-003 attachments/modules-export-inventory-v0.1.md`（22 包 · 279 行），冻结面 v1.1.0 增列引用（A-002 核对点 #6）。
- 冻结面 v1.0.1 → **v1.1.0**：增列 B+ 层（assembly 导出签名）+ B 层盘点引用；契约内容 additive，无 breaking。