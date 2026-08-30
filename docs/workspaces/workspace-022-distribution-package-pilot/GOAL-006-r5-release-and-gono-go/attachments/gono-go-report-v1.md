# go/no-go 报告 · 分发形态包化试点（R5 · 2026-08-29）

> 依据 VP-022 v0.2.0「激活前置要求 §2 · go/no-go 触发框架」判定。累计证据 span = workspace-022（R1–R5）。

## 1. 判据达成全景（VP-022 六条退出判据）

| # | 判据 | 证据 | 状态 |
|---|------|------|------|
| 1 | Go 库消费闭环 | GOAL-003 done 4/4：internal 外移 + `assembly` 公开装配工厂 + golden-consumer users 全链（迁移 fresh · 贡献 = Descriptor 声明） | ✅ |
| 2 | Web 包消费闭环 | GOAL-004 done 4/4：`@schema-ui/protocol` 0.2.0 + `@schema-ui/renderer` 0.1.0 + SSR 渲染 + Token 覆盖纪律 | ✅ |
| 3 | 零冲突升级演练 | GOAL-005 done 4/4：V2 演进（A 层/协议/迁移 0063）→ gc 1 行 bump / gw 0 行 diff · 冲突 0 · 无 merge · 全量回归绿 | ✅ |
| 4 | 契约冻结面落盘 | freeze-face v1.1.0（A/B/B+/C 分层 + semver 流程 + changelog 模板）；Web 六包边界 v1.2 定稿（S4 随本报告） | ✅ |
| 5 | 发布可复现 | `scripts/pack-npm-packages.mjs` 一键 tgz（protocol/renderer）；Go tag `v0.0.2`；golden-web **tarball 安装**（registry 语义）回归全绿 | ✅ |
| 6 | go/no-go 报告 | 本报告 | ✅ |

## 2. 触发框架判向（VP-022 §2）

| 指标 | 判定 | 证据 |
|------|------|------|
| 判据 #1–#5 全部达成 | ✅ | 上表 |
| 零冲突演练冲突计数 = 0 | ✅ | GOAL-005 A-001 |
| 包路径升级总耗时 ≤ fork 同步对照基线 | ✅ | 实测：bump + 安装 + 回归 ≈ **分钟级**（R4 全程 < 20 分钟含主仓回归）；fork 对照 = 同一 V2 演变在合并模型下需**手动同步测试夹具等四处**（R4 演练实证：store 测试 4 文件 + 台账头常量手工同步——正是 fork 冲突的直接映射），小时级工程 |
| golden consumer 回归稳定 | ✅ | 三探针 × R4/R5 两轮全绿；tarball 语义安装通过 |

**判定：倾向推进 Charter strategic 修订**（把「构建期包消费」写入成功边界，fork 与包消费并列为单主线两种交付形态）。

## 3. Charter 修订建议（strategic · 候选 0.3.0）

1. **成功边界 #1 追加**：「或经**构建期包消费路径**（Go 模块 `go get` / npm 包组 `@schema-ui/*` + 公开装配工厂 `apps/api/assembly` + changelog 迁移说明）获得同一基架与协议兼容实现；fork 与包消费为同一单主线的两种交付形态，不维护平行代码线。」
2. **非目标追加澄清**：「禁止运行时模块下载 / 热插拔（**构建期包组装允许**——包在编译时进入组合根）；不承诺以 fork 为唯一消费路径（物理裁剪由 fork 或按需装包承担）。」
3. **配套**：VRev（self → 建议 independent 复核）+ 受影响 VP/区 re-align（VP-003/004/008 消费语义第 v0.3.0 注记；历史证据不重写）+ QUICKSTART/playbook 增补包消费路径章节。

## 4. 附带裁决建议（I-007 · 协议 pin 漂移）

- 代码事实：web 支持窗 `2.7–2.9`、`APP_MANIFEST_PROTOCOL_VERSION = 2.9`（`81aa1d8`），R4 演练证明 additive 超集无破坏。
- 建议 `/vision`：**协议来源 `v2.8.0` → `v2.9.0`**（Charter/roadmap/provenance 同步；`I-PROTO-FULL-001` 历史分母保留注记）；对既有闭合 VP 证据不追溯改写。
- 待用户裁决（两项：① 推进 Charter 修订 ② pin 升级）。

## 5. 残余与风险（go 后清单）

| 项 | 处置 |
|----|------|
| F-005 PG external 消费 | 无本地 PG 未实测；rendered 进 go 后发布回归（PG 生产线） |
| F-006 d.ts 自动化（TS5056） | go 后 monorepo 化时修（当前手工声明可用） |
| F-003 PG drain 全量并发时序 | 本日复核通过（无失败）；仍建议 `-short` 分流（VP-009 候选） |
| 六包细化（ui/theme 独立包） | go 后正式化（冻结面 v1.2 起列为目标形态） |
| npm 公开 registry / GitHub Packages 上传 | 外部发布动作（凭据/网络），留 go 后流水线 |
| Grok independent 审计 | 随 R5 关门审计执行（见 GOAL-006 A-001） |