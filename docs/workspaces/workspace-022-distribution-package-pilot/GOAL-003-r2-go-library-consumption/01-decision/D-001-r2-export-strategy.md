---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-go-library-consumption
version: 0.1.0
---

# D-001 · internal 外移方案提案（关键决策 · 待用户裁决）

背景事实（E-001）：外部模块 import `apps/api/internal/kernel` 被 Go internal 规则拒绝。A 层（kernel）与 B 层（modules）必须在发布前移出 `internal/`。

## 方案

| 方案 | 内容 | 优点 | 缺点 | 判定 |
|------|------|------|------|------|
| **A · 目录提升（推荐）** | `internal/kernel` → `apps/api/kernel`；`internal/modules` → `apps/api/modules`；C 层（auth/handler/server/composition/store 方言等）保持 `internal/` | 单 go.mod 不变（G1 粗粒度）；改动机械（import 批量替换 + gofmt）；冻结面即公开路径，契约面天然可见 | 全仓数百文件 import 改写（回归由 `go build ./...` + 全量测试兜底）；目录结构一动影响 CI/文档路径引用 | ★ |
| **B · 独立模块（G2 变体）** | kernel/modules 各自 go.mod + 独立 tag | 细粒度版本 | 多模块 tag/发布矩阵成本高；模块间 store 引用成网；与 R1「G1 起步」决策相悖 | 不作为本期选项（R5 go 后评估） |
| **C · 发布拷贝/代码生成** | 发布时将 internal 复制到公开路径 | 不动主仓 | 双份代码漂移、防呆失效 | **否决** |

## 决策点（请用户确认）

1. **外移范围**：推荐方案 A 全文——A 层 + B 层提升为 `apps/api/kernel` / `apps/api/modules`（严格等于冻结面清单 §1/§2）；C 层不动。
2. **F-001 联动**：`users.New(*auth.Authenticator, …)` 的 C 层泄漏在方案 A 下仍存在（auth 不随移）——收敛候选：① auth 构造类型上移 kernel（A 层增契约）② 公开装配工厂 `apps/api/assembly`（新公开包，属 B+ 层）③ 白名单外移 auth。**建议**：R2 S3 用真实组合根实测后再定（三个候选都写进 S3 验证矩阵），本轮只定外移范围。
3. **主仓证据基线**：S2 完成后全量回归（`go test ./...` + 既有 e2e/协议套件）作为外移等价性闸门。

## 未选方案

- 不选 B（多模块）：本期 G1 粗粒度决策保持。
- 不选 C（拷贝发布）：防呆与一致性失败。
- 不因 external 规则修改 Go build 模式（如一直用 replace 本地目录）：released 消费仍要求公开路径。