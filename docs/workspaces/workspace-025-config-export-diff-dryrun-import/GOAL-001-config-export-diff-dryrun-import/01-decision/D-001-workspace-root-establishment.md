---
status: active
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-001-config-export-diff-dryrun-import
version: 0.1.0
---

# D-001 · 工作区 Root 建立（2026-08-30）

1. **绑定**：lead workspace = `workspace-025-config-export-diff-dryrun-import`（`vision_role: delivery`，单一 lead）；`root_goal` = `GOAL-001-config-export-diff-dryrun-import`（`parent: null`）；`primary_plan` = `VP-025-config-export-diff-dryrun-import`（active v0.2.0）；不改变 Charter `primary_workspace`。
2. **freshness 三字段（V-F089 · VRev-054 recommended → fixed）**：
   - `consumer_vp` = VP-025（消费候选基线：HEAD `055da2fd` · `apps/api/v0.3.0` · 六包）
   - `last_freshness_review_at` = 2026-08-30（VRev-054：`c9122478` → `055da2fd` PASS · 协议 pin / 依赖锁 / 迁移台账 / Profile 装配 / provenance 五域零变更）
   - `next_freshness_review_trigger` = 下一次 VP 激活/消费前主动核对；协议 pin / 依赖锁 / 迁移台账 / Profile 默认集 / provenance 任一变更即触发重验证（VP-008 `go` 消费有效性规则）
3. **审计模式**：阶段关门 default self；实证门禁（R4 证据 / 关门）可按需 independent（grok build 先例）。
4. **信息门禁（P-005 · V-F090 → fixed）**：I-025-001（R1 前置 · 包内容边界/密钥处理）与 I-025-002（R1 前置 · 落地形态）required，**R1 合同冻结前必须经用户裁决**；I-025-004（R3 前置 · 导入失败语义）required；I-025-005 投影 `registered`（Profile 红线冻结不进）——登记于 Root `00-meta`。
5. **现状锚点（V-F091 → fixed）**：R1 合同冻结的「配置包内容边界」以 **serve 壳配置树**为对象面——`apps/api/server/config.default.yaml`（内嵌默认 · `profile: admin`）+ `server/config.go` 装载（env 插值 `$VAR` fail-closed / `$VAR:-default`）+ 骨架模板 `config.yaml.tmpl`；密钥/敏感值按冻结规则排除或脱敏（fail-closed 保持）。
6. **红线（激活即生效）**：不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；热加载不进分母。
7. **Root progress**：0/4（R1～R4 检查点等权；R1 待立项 GOAL-002）。