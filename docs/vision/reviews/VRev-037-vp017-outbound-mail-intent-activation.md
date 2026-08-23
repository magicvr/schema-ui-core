---
doc_type: vision-review
id: VRev-037-vp017-outbound-mail-intent-activation
status: active
source: independent
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
parent: null
---

# VRev-037 · VP-017 出站邮件 · 意图完备 / 激活就绪 / 对齐链（2026-08-22）

| 字段 | 值 |
|------|-----|
| source | independent |
| auditor | DeepSeek Harness（`/vision-audit`） |
| scope | VP-017-outbound-mail · 意图完备性 / 激活就绪门闩 / 对齐链 / 信息需求登记 / 架构边界 |
| audit_type | vision-plan |
| verdict | **pass** |
| 建议 class | editorial |

---

## 范围与结论

本轮独立审视 [VP-017-outbound-mail](../plans/VP-017-outbound-mail.md)（`planned` · `v0.1.0` · 架构 A6 · 0 区）在以下五个维度的合规性：

1. 对齐链（`vision_ref` → Charter → 单愿景）
2. 意图完备性与退出分母可判定性
3. P-005 信息需求登记质量
4. 组合位置与边界清洁度
5. 激活就绪门闩设计

**结论**：五项全部通过。VP-017 是与 VP-013 / VP-014 / VP-015 / VP-016 同一模式的架构 VP，意图清晰、边界显式、信息登记完整。无 required finding；两项 recommended 建议供 `/vision` 激活时参考。

---

## 逐项核查

### 1. 对齐链

| 检查项 | 证据 | 判定 |
|--------|------|------|
| `vision_ref` = `schema-ui-core-admin-foundation@0.2.0` | VP-017 frontmatter 第 6 行 | ✓ 与 Charter `vision_id@version` 精确匹配 |
| Charter `status: active` · 单愿景制 | `charter.md` frontmatter | ✓ 唯一 active Charter |
| VP `status: planned` · 0 区合规 | VP-017 §工作区绑定；`alignment.md` §5 | ✓ `planned` 允许 0 个工作区 |
| 组合编排已登记 | `roadmap.md` 条目 17（架构 A6） | ✓ |
| `workspaces.md` 无虚挂条目 | 扫描 16 行，无 workspace-017 行 | ✓ 符合 0 区预期 |

### 2. 意图完备性与退出分母

| 检查项 | 证据 | 判定 |
|--------|------|------|
| 核心意图清晰（发送端口 + SMTP + capture/log sink） | VP-017 §意图 1–4 | ✓ |
| 退出分母表（5 列能力矩阵）明确「本 VP 交付」vs「不进本 VP」 | VP-017 §首波冻结 | ✓ |
| 方向级退出判据可独立核查（6 条，含「开放 required = 0」） | VP-017 §方向级退出判据 | ✓ |
| 非目标列表完整（SMS、账号 email、邀请、恢复状态机、Notification Transport 等） | VP-017 §非目标 | ✓ |
| 与相邻 VP 边界显式声明（VP-003/012/016/008/009/010 + 后续消费链） | VP-017 §与相邻 VP 的边界 | ✓ |
| 不重开已 closed VP | VP-017 §意图末段；§非目标 | ✓ |

### 3. P-005 信息需求登记

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 审计判定 |
|----|--------------|------|----------|----------|----------|
| I-017-001 | STARTTLS（587）vs 隐式 TLS（465）；只钉一种可核对路径 | required | 方案冻结 / 实施 | R2 接入前 | ✓ 合规；影响门禁与最晚阶段对应 |
| I-017-002 | 配置键名与凭证注入（YAML + env fail-closed；secret 不入库） | required | 方案冻结 | R2 接入前 | ✓ 合规 |
| I-017-003 | 默认 sink：进程内 capture vs 只写结构化日志；测试如何取出报文 | required | 方案冻结 | R1 端口冻结 | ✓ 合规 |
| I-017-004 | 单次 `Send` 的 To 基数（建议单收件人）| required | 方案冻结 | R1 端口冻结 | ✓ 合规；建议意见已内嵌 |
| I-017-005 | HTML/MIME 是否作为可选体（建议不进退出分母） | non-blocking | 关门叙事 | R4 | ✓ 合规；级别与影响匹配 |

全部 5 项已登记，字段完整（问题、级别、门禁、最晚阶段、状态）。初始状态均 `open`，符合「允许带未知立项」规则。

### 4. 组合位置与架构边界清洁度

- **架构 A6**：位于已 closed 的 A5（VP-016）之后，符合建议顺序（roadmap §架构分支建议顺序）。
- **不重开任何 closed VP**：已显式声明不重开 VP-012 / VP-016；不干扰 VP-009 / VP-010 持续程序。
- **三分支归属清晰**：运输实现归架构；消费链（账号邮箱身份 + IAM 恢复状态机）留 Admin 功能分支；roadmap §Admin 功能分支末段已登记消费顺序。
- **Charter 边界不受影响**：VP §非目标显式排除业务域页面与 Notification Transport 产品。

### 5. 激活就绪门闩

| 门闩 | VP 原文位置 | 判定 |
|------|------------|------|
| 激活前须 self Vision Review | VP-017 §状态与门闩 `Vision required` | ✓ |
| 激活须另一次用户确认 + 开区 | VP-017 §状态与门闩 `lead_workspace` | ✓ |
| 本文件落盘 ≠ 方向已稳到可开区 | VP-017 §状态与门闩 | ✓ |
| 未激活前不得以本文件为 `primary_plan` 推进实现 | VP-017 §状态与门闩 `关门门闩` | ✓ |
| 架构类 freshness review 在激活前 | VP-017 §状态与门闩 + §与相邻 VP 边界（VP-008 `go`） | ✓ |

门闩设计与 VP-013 / VP-015 / VP-016 激活模式一致，无跳过证据。

---

## Findings

### V-F070 · recommended · 激活时应显式登记 freshness 检查的候选锚点

**级别**：recommended  
**严重度**：低  
**证据**：VP-017 §与相邻 VP 的边界（VP-008 `go`）描述「激活前仍须架构类 freshness review」，但未明确指出应以哪个 HEAD 对比原 `go` 候选 `ed99e88`，也未指明若比对区间存在 VP-009/VP-010 open finding 时的暂挂规则是否自动适用。  
**背景**：VP-013/014/015/016 均在 Root D-001 中记录「原 `go` 候选 → 现 HEAD → 非业务解锁 → 不暂挂」路径。VP-017 作为下一轮架构 VP，期望沿用此路径，但文本中未明示。  
**关闭要求**：激活时的 self Vision Review 或 Root D-001 显式记录 freshness 结论即视为 fixed；激活前可不修改本 VP 文本。

### V-F071 · recommended · 可选注册「生效方式」为 I-017-006 非阻断信息项

**级别**：recommended  
**严重度**：低  
**证据**：VP-017 §配置面 末段明确「本波默认进程重启后生效。热加载不进退出分母。」这是一个明确的架构决策，与 VP-016 的同类决策结构一致，但未在 P-005 信息需求表中显式登记（VP-016 亦未登记，是同批 pattern）。  
**背景**：该决策无歧义、已冻结，影响仅在「是否允许热加载」这一非阻断问题上。注册为 `I-017-006 non-blocking` 可使信息需求表更具完整追踪性，但不构成必改项。  
**关闭要求**：在 VP 或 Root 信息需求表中补充 `I-017-006 non-blocking`，或在激活 Vision Review 中以书面说明「与 VP-016 pattern 一致、不另注册」为 fixed。

---

## 声明

本意见 `source: independent`，不修改 VP-017 / Charter / Goal status，不自行闭合任何 finding。

- **V-F070 / V-F071**：recommended，非 required。允许在激活 self Vision Review 时同步响应，亦可提前响应。
- 本 VP `planned` 状态与 0 区绑定合规；审视范围内**无 required finding**，**open required = 0**。
- 建议的 `/vision` 响应输入：激活 VP-017、绑定 lead workspace（候选 `workspace-017-outbound-mail`）、触发 self Vision Review 并在 Root D-001 中记录 freshness 结论。

| finding | level | 状态 |
|---------|-------|------|
| V-F070 | recommended | open |
| V-F071 | recommended | open |

### 响应（2026-08-22 · `/vision` 激活与 `/govern` 开区）

不回溯改写原 verdict `pass` 与 finding 正文。

| finding | 闭合 | 证据 |
|---------|------|------|
| V-F070 | **fixed** | [VRev-038](VRev-038-vp017-activation-self.md) 与 Root [D-001](../../workspaces/workspace-017-outbound-mail/GOAL-001-outbound-mail/01-decision/D-001-workspace-root-establishment.md) 架构类 freshness：原 `go` `ed99e88`；HEAD `250cb9c`；非业务解锁；VP-009/010 无新暂挂；F-007 不升格；不暂挂 `go`。 |
| V-F071 | **fixed** | VP-017 v0.2.0 与 Root 信息表登记 `I-017-006` / `I-006` **non-blocking**：本波重启生效，热加载不进退出分母。 |

用户书面「响应独立审计意见，然后激活 vp，然后交 `/govern` 开启工作区」已执行：self VRev-038 `pass`；VP `active`；lead `workspace-017-outbound-mail`（惯例 slug，D-001 留痕）；Root scaffold。本 scope **0 open required、0 open recommended**。
