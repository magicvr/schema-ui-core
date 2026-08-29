# 产线化报告（R5 · 2026-08-29）

> VP-023 收官报告：判据 #1–5 满足声明 + 实证数据 + go 后清单核销 + 默认主路径建议（**不改 Charter 措辞**——fork 并存维持，用户既定）。

## 1. 判据达成全景（VP-023 六条）

| # | 判据 | 证据 | 状态 |
|---|------|------|------|
| 1 | 真实发布通道 | Go：`apps/api/v0.1.0/v0.2.0` tag + 公共 proxy `go get` 实证；npm：六包 GitHub Packages 发布（protocol 0.2.0 / lib 0.1.0 / theme 0.1.0 / ui 0.1.0 / renderer 0.2.0 / shell 0.1.0）+ golden-field `pnpm add @ver` 安装实证 | ✅ |
| 2 | CLI 闭环 | `schema-ui create/add/upgrade`（Go 单二进制 · 零依赖 · go:embed 模板）；create 双端全绿（双轨同构）；一次 registry 升级零冲突（v0.1.0→v0.2.0）· F-001 核销 | ✅ |
| 3 | 六包细化 + d.ts | 六包独立发布 + TS5056 根治（render/form-controls 改名 · 五包 tsc declaration 全 0）· F-006 核销；冻结面 v1.3.0 | ✅ |
| 4 | 覆盖运维 | PG external 实测（postgres:16 · 64 迁移 apply · 幂等）· F-005 核销；ops-playbook + compose/Dockerfile；consumer-regression workflow | ✅ |
| 5 | 上手与迁移 | QUICKSTART 方法 B（cli+包 起步）；fork→包迁移指南；**从零走查 8.4s**（create 0.5 + 装配 6.9 + web/探针 1.1，依赖缓存预置口径） | ✅ |
| 6 | 产线化报告 | 本报告（下节） | ✅ |

## 2. 实证数据

| 指标 | 值 | 口径 |
|------|-----|------|
| create → 双端绿 | **8.4s** | 依赖缓存预置（QUICKSTART 口径）；golden-field/demo-admin 双仓复现 |
| registry 升级 | **零冲突 · 秒级** | `schema-ui upgrade`（go+pnpm @latest + 探针回归）· 冲突 0 · 无 merge |
| 发布往返（实测） | tag 推送 → sumdb 可消费 | 分钟级（首次发布含 sumdb 索引时延知识） |
| Go tag 命名 | `apps/api/vX.Y.Z`（子目录约束，知识 §R1） | 发布脚本化注意点 |
| d.ts 自动化 | 五包 declaration 全 0 | TS5056 根治（改名方案） |

**breaking 演练**：流程侧已备（semver-breaking-policy §3 + changelog 迁移说明模板 + R2 升级演练对照）；**实演 = go 后首个 major 发布时执行**（D-001 定案：不向 registry 实发 breaking 影响真实消费面）。

## 3. go 后清单核销表

| 项 | 状态 |
|----|------|
| origin tag + Go proxy 发布（F-001-R1 复审项） | ✅ 核销（R1） |
| 配置键/依赖样本补测（F-002-R4） | ⏳ 保留：breaking/配置键样本并入 go 后首个 major 发布回归 |
| CI 接入 + registry 上传（F-003-R4） | ✅ 核销（R1 上传 + consumer-regression workflow） |
| 六包细化 + d.ts（F-006） | ✅ 核销（R3） |
| PG external（F-005） | ✅ 核销（R4） |
| fork 对照计时实验 | ⏳ 保留：fork 同步对照 = go 后（当前定性成立：R4 演练主仓 4 处手工同步 = fork 冲突映射） |
| **go 后新增** | ① `schema-ui serve`（HTTP 壳 + config 装载 + assembly 服务器面）② renderer 依赖图 external 化（ui 包消费）③ 纯原子拆分（业务组件出 ui 包）④ compose 实跑核验（CI）⑤ breaking 实演 ⑥ fork→包迁移工具化 ⑦ npm registry 转为 npmjs/GH Packages 公开可见性决策 ⑧ 明 release 补 sumdb 等待 |

## 4. 默认主路径建议（不改 Charter）

**建议**：cli+包 达「默认主路径」条件——QUICKSTART 方法 B 与 A 并列置顶、`schema-ui serve` 落地后即可正式宣布。**执行层建议**：go 后清单 ①（serve 壳）完成后，将方法 B 设为 QUICKSTART 首段（fork 为第二路径），Charter 0.3.0 措辞维持不变（fork 逃生舱语义）。

## 5. 残余风险

- 公开可发现性：GH Packages 包默认私有（依赖 token 消费；公开化 = 组织决策）。
- breaking 流程首次实战未发生（政策与模板已备）。
- HTTP/serve 壳未交付（运维面部分契约引用）。